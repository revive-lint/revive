// Package cli implements the revive command line application.
package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/afero"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/revivelib"
)

const (
	defaultVersion = "dev"
	defaultCommit  = "none"
	defaultDate    = "unknown"
	defaultBuilder = "unknown"
)

var (
	version = defaultVersion
	commit  = defaultCommit
	date    = defaultDate
	builtBy = defaultBuilder
	// AppFs is used for operations related to user config files.
	AppFs = afero.NewOsFs()
)

func fail(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1) //revive:disable-line:deep-exit
}

// RunRevive runs the CLI for revive.
func RunRevive(extraRules ...revivelib.ExtraRule) {
	// Move parsing flags outside of init(); otherwise, tests don't work properly.
	// More info: https://github.com/golang/go/issues/46869#issuecomment-865695953
	initConfig()

	if versionFlag {
		fmt.Print(getVersion(builtBy, date, commit, version))
		return
	}

	conf, err := config.GetConfig(configPath)
	if err != nil {
		fail(err.Error())
	}

	revive, err := revivelib.New(
		conf,
		setExitStatus,
		maxOpenFiles,
		extraRules...,
	)
	if err != nil {
		fail(err.Error())
	}

	files := flag.Args()
	packages := []*revivelib.LintPattern{}

	for _, file := range files {
		packages = append(packages, revivelib.Include(file))
	}

	for _, file := range excludePatterns {
		packages = append(packages, revivelib.Exclude(file))
	}

	failures, err := revive.Lint(packages...)
	if err != nil {
		fail(err.Error())
	}

	output, exitCode, err := revive.Format(formatterName, failures)
	if err != nil {
		fail(err.Error())
	}

	if output != "" {
		fmt.Println(output)
	}

	os.Exit(exitCode) //revive:disable-line:deep-exit
}

var (
	configPath      string
	excludePatterns revivelib.ArrayFlags
	formatterName   string
	versionFlag     bool
	setExitStatus   bool
	maxOpenFiles    int
)

var originalUsage = flag.Usage

func logo() string {
	return color.YellowString(` _ __ _____   _(_)__  _____
| '__/ _ \ \ / / \ \ / / _ \
| | |  __/\ V /| |\ V /  __/
|_|  \___| \_/ |_| \_/ \___|`)
}

func call() string {
	return color.MagentaString("revive -config c.toml -formatter friendly -exclude a.go -exclude b.go ./...")
}

func banner() string {
	return fmt.Sprintf(`
%s

Example:
  %s
`, logo(), call())
}

func buildDefaultConfigPath() string {
	var result string
	var homeDirFile string
	configFileName := "revive.toml"
	var configDirFile string
	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		configDirFile = filepath.Join(configDir, configFileName)
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		homeDirFile = filepath.Join(homeDir, configFileName)
	}

	switch {
	case configDirFile != "" && fileExist(configDirFile):
		result = configDirFile
	case fileExist(homeDirFile):
		result = homeDirFile
	default:
		result = ""
	}

	return result
}

func initConfig() {
	if os.Getenv("REVIVE_FORCE_COLOR") == "1" {
		color.NoColor = false //nolint:reassign // We want to reassign the default value of NoColor to force colorizing for non-TTY environments.
	}

	flag.Usage = func() { //nolint:reassign // We want to reassign the default usage function to print our banner.
		fmt.Println(banner())
		originalUsage()
	}

	// command line help strings
	const (
		configUsage       = "path to the configuration TOML file, defaults to $XDG_CONFIG_HOME/revive.toml or $HOME/revive.toml, if present (i.e. -config myconf.toml)"
		excludeUsage      = "list of globs which specify files to be excluded (i.e. -exclude foo/...)"
		formatterUsage    = "formatter to be used for the output (i.e. -formatter stylish)"
		versionUsage      = "get revive version"
		exitStatusUsage   = "set exit status to 1 if any issues are found, overwrites error-code and warning-code in config"
		maxOpenFilesUsage = "maximum number of open files at the same time"
	)

	defaultConfigPath := buildDefaultConfigPath()

	flag.StringVar(&configPath, "config", defaultConfigPath, configUsage)
	flag.Var(&excludePatterns, "exclude", excludeUsage)
	flag.StringVar(&formatterName, "formatter", "", formatterUsage)
	flag.BoolVar(&versionFlag, "version", false, versionUsage)
	flag.BoolVar(&setExitStatus, "set_exit_status", false, exitStatusUsage)
	flag.IntVar(&maxOpenFiles, "max_open_files", 0, maxOpenFilesUsage)
	flag.Parse() //revive:disable-line:deep-exit
}

// getVersion returns build info (version, commit, date, and builtBy).
func getVersion(builtBy, date, commit, version string) string {
	var buildInfo string
	if date != defaultDate && builtBy != defaultBuilder {
		buildInfo = fmt.Sprintf("Built\t\t%s by %s\n", date, builtBy)
	}

	if commit != defaultCommit {
		buildInfo = fmt.Sprintf("Commit:\t\t%s\n%s", commit, buildInfo)
	}

	if version == defaultVersion {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			version = strings.TrimPrefix(bi.Main.Version, "v")
			if buildInfo == "" {
				return fmt.Sprintf("version %s\n", version)
			}
		}
	}

	return fmt.Sprintf("Version:\t%s\n%s", version, buildInfo)
}

func fileExist(path string) bool {
	_, err := AppFs.Stat(path)
	return err == nil
}
