// Package revivelib provides revive's linting functionality as a lib.
package revivelib

import (
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/revive-lint/dots"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/logging"
)

// Revive is responsible for running linters and formatters
// and returning a set of results.
type Revive struct {
	config       *lint.Config
	lintingRules []lint.Rule
	logger       *slog.Logger
	maxOpenFiles int
}

// New creates a new instance of [Revive] lint runner.
func New(
	conf *lint.Config,
	setExitStatus bool,
	maxOpenFiles int,
	extraRules ...ExtraRule,
) (*Revive, error) {
	logger, err := logging.GetLogger()
	if err != nil {
		return nil, fmt.Errorf("initializing revive - getting logger: %w", err)
	}

	if setExitStatus {
		conf.ErrorCode = 1
		conf.WarningCode = 1
	}

	extraRuleInstances := make([]lint.Rule, len(extraRules))
	for i, extraRule := range extraRules {
		extraRuleInstances[i] = extraRule.Rule

		ruleName := extraRule.Rule.Name()

		_, isRuleAlreadyConfigured := conf.Rules[ruleName]
		if !isRuleAlreadyConfigured {
			conf.Rules[ruleName] = extraRule.DefaultConfig
		}
	}

	lintingRules, err := config.GetLintingRules(conf, extraRuleInstances)
	if err != nil {
		return nil, fmt.Errorf("initializing revive - getting lint rules: %w", err)
	}

	logger.Info("Config loaded", "rules", slices.Collect(maps.Keys(conf.Rules)))

	return &Revive{
		logger:       logger,
		config:       conf,
		lintingRules: lintingRules,
		maxOpenFiles: maxOpenFiles,
	}, nil
}

// Lint the included patterns, skipping excluded ones.
func (r *Revive) Lint(patterns ...*LintPattern) (<-chan lint.Failure, error) {
	includePatterns := []string{}
	excludePatterns := []string{}

	for _, lintpkg := range patterns {
		if lintpkg.IsExclude() {
			excludePatterns = append(excludePatterns, lintpkg.Pattern())
		} else {
			includePatterns = append(includePatterns, lintpkg.Pattern())
		}
	}

	if len(excludePatterns) == 0 { // if no excludes were set
		excludePatterns = r.config.Exclude // use those from the configuration
	}

	// by default if no excludes exclude vendor
	if len(excludePatterns) == 0 {
		excludePatterns = []string{"vendor/..."}
	}

	packages, err := getPackages(includePatterns, excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("linting - getting packages: %w", err)
	}

	revive := lint.New(func(file string) ([]byte, error) {
		contents, err := os.ReadFile(file) //nolint:gosec // ignore G304: potential file inclusion via variable
		if err != nil {
			return nil, fmt.Errorf("reading file %v: %w", file, err)
		}

		return contents, nil
	}, r.maxOpenFiles)
	revive.SetLogger(r.logger)

	failures, err := revive.Lint(packages, r.lintingRules, *r.config)
	if err != nil {
		return nil, fmt.Errorf("linting - retrieving failures channel: %w", err)
	}

	return failures, nil
}

// Format gets the output for a given failures channel from [Revive.Lint].
func (r *Revive) Format(
	formatterName string,
	failuresChan <-chan lint.Failure,
) (output string, exitCode int, err error) {
	conf := r.config
	formatChan := make(chan lint.Failure)
	exitChan := make(chan bool)

	formatter, err := config.GetFormatter(formatterName)
	if err != nil {
		return "", 0, fmt.Errorf("formatting - getting formatter: %w", err)
	}

	var (
		out       string
		formatErr error
	)

	go func() {
		out, formatErr = formatter.Format(formatChan, *conf)

		exitChan <- true
	}()

	for failure := range failuresChan {
		if failure.Confidence < conf.Confidence {
			continue
		}

		if exitCode == 0 {
			exitCode = conf.WarningCode
		}

		if failure.SeverityFor(conf) == lint.SeverityError {
			exitCode = conf.ErrorCode
		}

		formatChan <- failure
	}

	close(formatChan)
	<-exitChan

	if formatErr != nil {
		return "", exitCode, fmt.Errorf("formatting: %w", formatErr)
	}

	return out, exitCode, nil
}

func getPackages(includePatterns []string, excludePatterns ArrayFlags) ([][]string, error) {
	globs := normalizeSplit(includePatterns)
	if len(globs) == 0 {
		globs = append(globs, ".")
	}

	packages, err := dots.ResolvePackages(globs, normalizeSplit(excludePatterns))
	if err != nil {
		return nil, fmt.Errorf("getting packages - resolving packages in dots: %w", err)
	}

	return packages, nil
}

func normalizeSplit(strs []string) []string {
	res := []string{}

	for _, s := range strs {
		t := strings.Trim(s, " \t")
		if t != "" {
			res = append(res, t)
		}
	}

	return res
}

// ArrayFlags type for string list.
// Implements [flag.Value] interface, to be used in command line arguments.
type ArrayFlags []string

var _ flag.Value = (*ArrayFlags)(nil)

// String returns the space-separated representation of the ArrayFlags.
func (i *ArrayFlags) String() string {
	return strings.Join([]string(*i), " ")
}

// Set value for array flags.
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)

	return nil
}
