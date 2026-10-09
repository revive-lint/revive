// Package analyzer provides revive as a golang.org/x/tools/go/analysis analyzer.
package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"golang.org/x/tools/go/analysis"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

// New returns an analyzer applying the rules enabled by conf and the given extra rules.
func New(conf *lint.Config, extraRules ...revivelib.ExtraRule) (*analysis.Analyzer, error) {
	config.Normalize(conf)

	extraRuleInstances := make([]lint.Rule, len(extraRules))
	for i, extraRule := range extraRules {
		extraRuleInstances[i] = extraRule.Rule

		ruleName := extraRule.Rule.Name()
		if _, ok := conf.Rules[ruleName]; !ok {
			conf.Rules[ruleName] = extraRule.DefaultConfig
		}
	}

	rules, err := config.GetLintingRules(conf, extraRuleInstances)
	if err != nil {
		return nil, fmt.Errorf("creating analyzer - getting lint rules: %w", err)
	}

	return &analysis.Analyzer{
		Name: "revive",
		Doc:  "fast, configurable, extensible, flexible, and beautiful linter for Go",
		Run: func(pass *analysis.Pass) (any, error) {
			return nil, run(pass, *conf, rules)
		},
	}, nil
}

func run(pass *analysis.Pass, conf lint.Config, rules []lint.Rule) error {
	pkg := lint.NewPackage(pass.Fset, goVersion(pass, conf), pass.Pkg, pass.TypesInfo)

	for _, astFile := range pass.Files {
		if !conf.IgnoreGeneratedHeader && ast.IsGenerated(astFile) {
			continue
		}

		filename := pass.Fset.File(astFile.Pos()).Name()
		content, err := pass.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", filename, err)
		}

		pkg.AddFile(lint.NewFileFromAST(filename, content, astFile, pkg))
	}

	failures := make(chan lint.Failure)
	lintErr := make(chan error, 1)
	go func() {
		lintErr <- pkg.Lint(rules, conf, failures)
		close(failures)
	}()

	for failure := range failures {
		pass.Report(toDiagnostic(pass.Fset, failure))
	}

	return <-lintErr
}

func goVersion(pass *analysis.Pass, conf lint.Config) *goversion.Version {
	if conf.GoVersion != nil {
		return conf.GoVersion
	}

	if pass.Module == nil {
		return nil
	}
	version := strings.TrimPrefix(pass.Module.GoVersion, "go")
	if version == "" {
		return nil
	}

	return goversion.Must(goversion.NewVersion(version))
}

func toDiagnostic(fset *token.FileSet, failure lint.Failure) analysis.Diagnostic {
	diagnostic := analysis.Diagnostic{
		Pos:      failure.Pos,
		End:      failure.End,
		Category: failure.RuleName,
		Message:  failure.Failure,
	}

	if fix := toSuggestedFix(fset, failure); fix != nil {
		diagnostic.SuggestedFixes = []analysis.SuggestedFix{*fix}
	}

	return diagnostic
}

func toSuggestedFix(fset *token.FileSet, failure lint.Failure) *analysis.SuggestedFix {
	if failure.ReplacementLine == "" {
		return nil
	}

	tokenFile := fset.File(failure.Pos)
	line := tokenFile.Line(failure.Pos)
	lineStart := tokenFile.LineStart(line)
	lineEnd := token.Pos(tokenFile.Base() + tokenFile.Size())
	if line < tokenFile.LineCount() {
		lineEnd = tokenFile.LineStart(line+1) - 1
	}

	return &analysis.SuggestedFix{
		Message: "Replace the line by: " + failure.ReplacementLine,
		TextEdits: []analysis.TextEdit{{
			Pos:     lineStart,
			End:     lineEnd,
			NewText: []byte(failure.ReplacementLine),
		}},
	}
}
