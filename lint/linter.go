package lint

import (
	"bytes"
	"fmt"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	goversion "github.com/hashicorp/go-version"
	"golang.org/x/mod/modfile"
	"golang.org/x/sync/errgroup"

	"github.com/mgechev/revive/internal/timing"
)

// ReadFile defines an abstraction for reading files.
type ReadFile func(path string) (result []byte, err error)

// Linter is used for linting set of files.
type Linter struct {
	reader         ReadFile
	fileReadTokens chan struct{}
	logger         *slog.Logger
}

// New creates a new Linter.
func New(reader ReadFile, maxOpenFiles int) Linter {
	var fileReadTokens chan struct{}
	if maxOpenFiles > 0 {
		fileReadTokens = make(chan struct{}, maxOpenFiles)
	}
	return Linter{
		reader:         reader,
		fileReadTokens: fileReadTokens,
		logger:         slog.New(slog.DiscardHandler),
	}
}

// SetLogger sets the logger for the linter if the provided logger is not nil.
func (l *Linter) SetLogger(logger *slog.Logger) {
	if logger != nil {
		l.logger = logger
	}
}

func (l *Linter) readFile(path string) ([]byte, error) {
	if l.fileReadTokens != nil {
		// "take" a token by writing to the channel.
		// It will block if no more space in the channel's buffer
		l.fileReadTokens <- struct{}{}
		defer func() {
			// "free" a token by reading from the channel
			<-l.fileReadTokens
		}()
	}

	return l.reader(path)
}

var (
	generatedPrefix  = []byte("// Code generated ")
	generatedSuffix  = []byte(" DO NOT EDIT.")
	packageKeyword   = []byte("package")
	defaultGoVersion = goversion.Must(goversion.NewVersion("1.0"))
)

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(packages [][]string, ruleSet []Rule, config Config) (<-chan Failure, error) {
	failures := make(chan Failure)

	doneVersions := timing.Track(l.logger, "go-version-detection")
	perModVersions := map[string]*goversion.Version{}
	modules := 0
	perPkgVersions := make([]*goversion.Version, len(packages))
	for n, files := range packages {
		if len(files) == 0 {
			continue
		}
		if config.GoVersion != nil {
			perPkgVersions[n] = config.GoVersion
			continue
		}

		dir, err := filepath.Abs(filepath.Dir(files[0]))
		if err != nil {
			doneVersions("error", err)
			return nil, err
		}

		alreadyKnownMod := false
		for d, v := range perModVersions {
			if strings.HasPrefix(dir, d) {
				perPkgVersions[n] = v
				alreadyKnownMod = true
				break
			}
		}
		if alreadyKnownMod {
			continue
		}

		d, v, err := detectGoMod(dir)
		if err != nil {
			// No luck finding the go.mod file thus set the default Go version
			v = defaultGoVersion
			d = dir
			l.logger.Debug("go module not found, using the default Go version", "dir", dir, "version", v, "error", err)
		} else {
			modules++
			l.logger.Debug("go module detected", "dir", d, "version", v)
		}
		perModVersions[d] = v
		perPkgVersions[n] = v
	}

	doneVersions("packages", len(packages), "modules", modules)

	var ruleTimings *timing.Rules
	if timing.Enabled(l.logger) {
		ruleTimings = timing.NewRules()
		ruleSet = timedRules(ruleSet, ruleTimings)
	}

	doneLint := timing.Track(l.logger, "lint")
	var wg errgroup.Group
	for n := range packages {
		wg.Go(func() error {
			pkg := packages[n]
			gover := perPkgVersions[n]
			if err := l.lintPackage(pkg, gover, ruleSet, config, failures); err != nil {
				return fmt.Errorf("error during linting: %w", err)
			}
			return nil
		})
	}

	go func() {
		err := wg.Wait()
		if err != nil {
			failures <- NewInternalFailure(err.Error())
		}
		doneLint("packages", len(packages), "rules", len(ruleSet))
		ruleTimings.Log(l.logger)
		close(failures)
	}()

	return failures, nil
}

func (l *Linter) lintPackage(filenames []string, gover *goversion.Version, ruleSet []Rule, config Config, failures chan Failure) error {
	if len(filenames) == 0 {
		return nil
	}

	// A single directory can contain files from multiple Go packages:
	// source files and internal tests share one package name (e.g. "foo"),
	// while external test files use a different one ("foo_test").
	// Type-checking must be performed per Go package — mixing package "foo" and
	// "foo_test" files into one types.Package makes results non-deterministic
	// because the name used for the types.Package is picked from a Go map
	// (see Package.TypeCheck). Group files by their parsed package name first.
	fset := token.NewFileSet()
	filesByPkg := map[string]map[string]*File{}

	// Timing attributes are boxed into interfaces, which allocates: build them only when they are logged.
	timed := timing.Enabled(l.logger)
	donePackage := timing.Track(l.logger, "package-lint")
	var totalFiles int
	defer func() {
		if timed {
			donePackage("dir", filepath.Dir(filenames[0]), "files", totalFiles)
		}
	}()

	doneReadParse := timing.Track(l.logger, "package-read-parse")
	var generated, invalid int
	for _, filename := range filenames {
		content, err := l.readFile(filename)
		if err != nil {
			doneReadParse("files", len(filenames), "error", err)
			return err
		}
		if !config.IgnoreGeneratedHeader && isGenerated(content) {
			generated++
			continue
		}

		// Parse into a throwaway Package so NewFile can succeed; the real
		// assignment happens below once we know which group the file belongs to.
		tmp := &Package{fset: fset, files: map[string]*File{}, goVersion: gover, logger: l.logger}
		file, err := NewFile(filename, content, tmp)
		if err != nil {
			invalid++
			addInvalidFileFailure(filename, err.Error(), failures)
			continue
		}
		file.logger = l.logger

		pkgName := file.AST.Name.Name
		group, ok := filesByPkg[pkgName]
		if !ok {
			group = map[string]*File{}
			filesByPkg[pkgName] = group
		}
		group[filename] = file
		totalFiles++
	}
	if timed {
		doneReadParse("files", len(filenames), "generated", generated, "invalid", invalid)
	}

	if len(filesByPkg) == 0 {
		return nil
	}

	var eg errgroup.Group
	for _, groupFiles := range filesByPkg {
		pkg := &Package{
			fset:      fset,
			files:     groupFiles,
			goVersion: gover,
			logger:    l.logger,
		}
		for _, f := range groupFiles {
			f.Pkg = pkg
		}
		eg.Go(func() error {
			return pkg.lint(ruleSet, config, failures)
		})
	}

	return eg.Wait()
}

// detectGoMod looks for the go.mod file in dir and its parent directories.
// It returns the directory containing the go.mod file and the Go version declared in it.
// It returns an error if no go.mod file is found, or if it can't be read, parsed or has no Go version.
func detectGoMod(dir string) (string, *goversion.Version, error) {
	modFileName, err := retrieveModFile(dir)
	if err != nil {
		return "", nil, fmt.Errorf("%q doesn't seem to be part of a Go module", dir)
	}

	mod, err := os.ReadFile(modFileName) //nolint:gosec // ignore G304: potential file inclusion via variable
	if err != nil {
		return "", nil, fmt.Errorf("failed to read %q, got %w", modFileName, err)
	}

	modAst, err := modfile.ParseLax(modFileName, mod, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse %q, got %w", modFileName, err)
	}

	if modAst.Go == nil {
		return "", nil, fmt.Errorf("%q does not specify a Go version", modFileName)
	}

	ver, err := goversion.NewVersion(modAst.Go.Version)
	return filepath.Dir(modFileName), ver, err
}

func retrieveModFile(dir string) (string, error) {
	const lookingForFile = "go.mod"
	for {
		// filepath.Dir returns 'C:\' on Windows, and '/' on Unix
		isRootDir := (dir == filepath.VolumeName(dir)+string(filepath.Separator))
		if dir == "." || isRootDir {
			return "", fmt.Errorf("did not found %q file", lookingForFile)
		}

		lookingForFilePath := filepath.Join(dir, lookingForFile)
		info, err := os.Stat(lookingForFilePath)
		if err != nil || info.IsDir() {
			// lets check the parent dir
			dir = filepath.Dir(dir)
			continue
		}

		return lookingForFilePath, nil
	}
}

// isGenerated reports whether the source file is generated code
// according to the rules from https://go.dev/s/generatedcode.
// The header must appear before the package clause, so scanning stops there.
func isGenerated(src []byte) bool {
	for line := range bytes.Lines(src) {
		b := bytes.TrimSuffix(bytes.TrimSuffix(line, []byte{'\n'}), []byte{'\r'})
		if bytes.HasPrefix(b, packageKeyword) {
			return false
		}
		if bytes.HasPrefix(b, generatedPrefix) && bytes.HasSuffix(b, generatedSuffix) && len(b) >= len(generatedPrefix)+len(generatedSuffix) {
			return true
		}
	}
	return false
}

// addInvalidFileFailure adds a failure for an invalid formatted file.
func addInvalidFileFailure(filename, errStr string, failures chan Failure) {
	position := getPositionInvalidFile(filename, errStr)
	failures <- Failure{
		Confidence: 1,
		Failure:    fmt.Sprintf("invalid file %s: %v", filename, errStr),
		Category:   failureCategoryValidity,
		Position:   position,
	}
}

// errPosRegexp matches with a NewFile error message:
//
//	corrupted.go:10:4: expected '}', found 'EOF
//
// The first group matches the line, and the second group matches the column.
var errPosRegexp = regexp.MustCompile(`.*:(\d*):(\d*):.*$`)

// getPositionInvalidFile gets the position of the error in an invalid file.
func getPositionInvalidFile(filename, s string) FailurePosition {
	pos := errPosRegexp.FindStringSubmatch(s)
	if len(pos) < 3 {
		return FailurePosition{}
	}
	line, err := strconv.Atoi(pos[1])
	if err != nil {
		return FailurePosition{}
	}
	column, err := strconv.Atoi(pos[2])
	if err != nil {
		return FailurePosition{}
	}

	return FailurePosition{
		Start: token.Position{
			Filename: filename,
			Line:     line,
			Column:   column,
		},
	}
}

// timedRules wraps each rule to accumulate the time spent applying it into timings.
func timedRules(rules []Rule, timings *timing.Rules) []Rule {
	result := make([]Rule, len(rules))
	for i, r := range rules {
		result[i] = timedRule{Rule: r, stat: timings.Stat(r.Name())}
	}
	return result
}

type timedRule struct {
	Rule

	stat *timing.RuleStat
}

func (r timedRule) Apply(file *File, arguments Arguments) []Failure {
	start := time.Now()
	defer func() { r.stat.Add(time.Since(start)) }()

	return r.Rule.Apply(file, arguments)
}
