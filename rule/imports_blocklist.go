package rule

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/mgechev/revive/lint"
)

// importsBlocklistEntry is a set of package patterns, optionally limited to matching files.
type importsBlocklistEntry struct {
	packages []*regexp.Regexp
	paths    []*lint.FileFilter
}

// ImportsBlocklistRule disallows importing the specified packages.
type ImportsBlocklistRule struct {
	entries []importsBlocklistEntry
}

var replaceImportRegexp = regexp.MustCompile(`/?\*\*/?`)

var _ lint.ConfigurableRule = (*ImportsBlocklistRule)(nil)

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
//
// Arguments are strings (global package bans) or maps with "packages"
// (required) and optional "paths" (file filters; same syntax as rule Exclude).
func (r *ImportsBlocklistRule) Configure(arguments lint.Arguments) error {
	r.entries = nil

	for i, arg := range arguments {
		switch v := arg.(type) {
		case string:
			entry, err := newImportsBlocklistEntry([]string{v}, nil)
			if err != nil {
				return err
			}
			r.entries = append(r.entries, entry)
		case map[string]any:
			entry, err := parseImportsBlocklistMap(v)
			if err != nil {
				return fmt.Errorf("invalid argument #%d to the imports-blocklist rule: %w", i, err)
			}
			r.entries = append(r.entries, entry)
		default:
			return fmt.Errorf("invalid argument to the imports-blocklist rule. Expecting a string or map, got %T", arg)
		}
	}
	return nil
}

func parseImportsBlocklistMap(args map[string]any) (importsBlocklistEntry, error) {
	var (
		packages []string
		paths    []string
		havePkgs bool
	)

	for k, v := range args {
		switch {
		case isRuleOption(k, "packages"):
			havePkgs = true
			list, err := stringListOption(v, "packages")
			if err != nil {
				return importsBlocklistEntry{}, err
			}
			packages = list
		case isRuleOption(k, "paths"):
			list, err := stringListOption(v, "paths")
			if err != nil {
				return importsBlocklistEntry{}, err
			}
			paths = list
		default:
			return importsBlocklistEntry{}, fmt.Errorf("unknown option %q", k)
		}
	}

	if !havePkgs {
		return importsBlocklistEntry{}, errors.New("missing required option \"packages\"")
	}
	if len(packages) == 0 {
		return importsBlocklistEntry{}, errors.New("option \"packages\" must not be empty")
	}

	return newImportsBlocklistEntry(packages, paths)
}

func stringListOption(v any, name string) ([]string, error) {
	raw, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid value for option %q: expected []string, got %T", name, v)
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value for option %q: expected string, got %T", name, item)
		}
		out = append(out, s)
	}
	return out, nil
}

func newImportsBlocklistEntry(packages, paths []string) (importsBlocklistEntry, error) {
	entry := importsBlocklistEntry{
		packages: make([]*regexp.Regexp, 0, len(packages)),
		paths:    make([]*lint.FileFilter, 0, len(paths)),
	}

	for _, pkg := range packages {
		re, err := regexp.Compile(fmt.Sprintf(`(?m)"%s"$`, replaceImportRegexp.ReplaceAllString(pkg, `(\W|\w)*`))) //nolint:gocritic // regexpSimplify: false positive
		if err != nil {
			return importsBlocklistEntry{}, fmt.Errorf("invalid package pattern %q: %w", pkg, err)
		}
		entry.packages = append(entry.packages, re)
	}

	for _, path := range paths {
		ff, err := lint.ParseFileFilter(path)
		if err != nil {
			return importsBlocklistEntry{}, fmt.Errorf("invalid path filter %q: %w", path, err)
		}
		entry.paths = append(entry.paths, ff)
	}

	return entry, nil
}

func (e importsBlocklistEntry) appliesToFile(name string) bool {
	if len(e.paths) == 0 {
		return true
	}
	for _, ff := range e.paths {
		if ff.MatchFileName(name) {
			return true
		}
	}
	return false
}

func (e importsBlocklistEntry) isBlocklisted(importPath string) bool {
	for _, regex := range e.packages {
		if regex.MatchString(importPath) {
			return true
		}
	}
	return false
}

// Apply applies the rule to given file.
func (r *ImportsBlocklistRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, is := range file.AST.Imports {
		path := is.Path
		if path == nil {
			continue
		}
		for _, entry := range r.entries {
			if !entry.appliesToFile(file.Name) {
				continue
			}
			if entry.isBlocklisted(path.Value) {
				failures = append(failures, lint.Failure{
					Confidence: 1,
					Failure:    "should not use the following blocklisted import: " + path.Value,
					Node:       is,
					Category:   lint.FailureCategoryImports,
				})
				break
			}
		}
	}

	return failures
}

// Name returns the rule name.
func (*ImportsBlocklistRule) Name() string {
	return "imports-blocklist"
}
