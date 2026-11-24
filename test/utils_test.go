package test_test

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
)

// configureRule configures the given rule with the given configuration
// if the rule implements the ConfigurableRule interface.
func configureRule(tb testing.TB, rule lint.Rule, arguments lint.Arguments) {
	tb.Helper()

	cr, ok := rule.(lint.ConfigurableRule)
	if !ok {
		return
	}

	err := cr.Configure(arguments)
	if err != nil {
		tb.Fatalf("Cannot configure rule %s: %v", rule.Name(), err)
	}
}

func testRule(tb testing.TB, filename string, rule lint.Rule, config ...*lint.RuleConfig) {
	tb.Helper()
	var ruleConfig lint.RuleConfig
	if len(config) > 0 {
		ruleConfig = *config[0]
	}

	lintConfig := lint.Config{
		Rules: map[string]lint.RuleConfig{
			rule.Name(): ruleConfig,
		},
	}

	testRuleWithLintConfig(tb, filename, rule, lintConfig)
}

func testRuleWithLintConfig(tb testing.TB, filename string, rule lint.Rule, config lint.Config) {
	tb.Helper()

	baseDir := filepath.Join("..", "testdata", filepath.Dir(filename))
	filename = filepath.Base(filename) + ".go"
	fullFilePath := filepath.Join(baseDir, filename)
	src, err := os.ReadFile(fullFilePath)
	if err != nil {
		tb.Fatalf("Bad filename path in test for %s: %v", rule.Name(), err)
	}

	configureRule(tb, rule, config.Rules[rule.Name()].Arguments)

	ins := parseInstructions(tb, fullFilePath, src)
	if ins == nil {
		assertSuccess(tb, fullFilePath, []lint.Rule{rule}, config)
		return
	}
	assertFailures(tb, fullFilePath, []lint.Rule{rule}, config, ins)
}

func assertSuccess(tb testing.TB, filePath string, rules []lint.Rule, config lint.Config) {
	tb.Helper()

	l := lint.New(os.ReadFile, 0)

	ps, err := l.Lint([][]string{{filePath}}, rules, config)
	if err != nil {
		tb.Errorf("Linting %s: %v", filePath, err)
		return
	}

	var failures []string
	for p := range ps {
		failures = append(failures, p.Failure)
	}
	if len(failures) > 0 {
		tb.Errorf("Expected the rule to pass but got the following failures: %s", strings.Join(failures, ", "))
	}
}

func assertFailures(tb testing.TB, filePath string, rules []lint.Rule, config lint.Config, ins []instruction) {
	tb.Helper()

	l := lint.New(os.ReadFile, 0)

	ps, err := l.Lint([][]string{{filePath}}, rules, config)
	if err != nil {
		tb.Errorf("Linting %s: %v", filePath, err)
		return
	}

	failures := []lint.Failure{}
	for f := range ps {
		if f.Category == "" {
			tb.Errorf("failure %q at %s:%d has no category", f.Failure, filePath, f.Position.Start.Line)
		}
		failures = append(failures, f)
	}

	type simplifiedFailure struct {
		File    string
		Line    int
		Failure string
	}
	var reportedFailures []simplifiedFailure

	for _, in := range ins {
		ok := false
		for i, p := range failures {
			if p.Position.Start.Line != in.Line {
				continue
			}

			if in.Match == p.Failure {
				// check replacement if we are expecting one
				if in.Replacement != "" {
					// ignore any inline comments, since that would be recursive
					r := p.ReplacementLine
					if i := strings.Index(r, " //"); i >= 0 {
						r = r[:i]
					}
					if r != in.Replacement {
						reportedFailures = append(reportedFailures, simplifiedFailure{
							File:    filePath,
							Line:    in.Line,
							Failure: fmt.Sprintf("Replacement: got %q, want %q", r, in.Replacement),
						})
					}
				}

				if in.Confidence > 0 {
					if in.Confidence != p.Confidence {
						reportedFailures = append(reportedFailures, simplifiedFailure{
							File:    filePath,
							Line:    in.Line,
							Failure: fmt.Sprintf("Confidence: got %f, want %f", p.Confidence, in.Confidence),
						})
					}
				}

				// remove this problem from ps
				copy(failures[i:], failures[i+1:])
				failures = failures[:len(failures)-1]

				ok = true
				break
			}
		}
		if !ok {
			reportedFailures = append(reportedFailures, simplifiedFailure{
				File:    filePath,
				Line:    in.Line,
				Failure: "want: " + in.Match,
			})
		}
	}

	for _, p := range failures {
		reportedFailures = append(reportedFailures, simplifiedFailure{
			File:    filePath,
			Line:    p.Position.Start.Line,
			Failure: "got:  " + p.Failure,
		})
	}

	slices.SortFunc(reportedFailures, func(a, b simplifiedFailure) int {
		if a.File != b.File {
			return strings.Compare(a.File, b.File)
		}
		if a.Line != b.Line {
			return a.Line - b.Line
		}

		return strings.Compare(a.Failure, b.Failure)
	})

	errorMessage := ""
	lastFileLine := ""
	for _, p := range reportedFailures {
		currentFileLine := fmt.Sprintf("%s:%d", p.File, p.Line)
		if currentFileLine != lastFileLine {
			if errorMessage != "" {
				tb.Error(errorMessage)
			}
			errorMessage = fmt.Sprintf("problem at %s: ", currentFileLine)
			lastFileLine = currentFileLine
		}

		errorMessage += "\n" + p.Failure
	}

	if errorMessage != "" {
		tb.Error(errorMessage)
	}
}

type instruction struct {
	Line        int     // the line number this applies to
	Match       string  // which pattern to match
	Replacement string  // what the suggested replacement line should be
	RuleName    string  // what rule we use
	Category    string  // which category
	Confidence  float64 // confidence level
}

// JSONInstruction structure used when we parse json object instead of classic MATCH string.
type JSONInstruction struct {
	Match      string  `json:"MATCH"`
	Category   string  `json:"Category"`
	Confidence float64 `json:"Confidence"`
}

// parseInstructions parses instructions from the comments in a Go source file.
// It returns nil if none were parsed.
func parseInstructions(tb testing.TB, filename string, src []byte) []instruction {
	tb.Helper()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		tb.Fatalf("Test file %v does not parse: %v", filename, err)
	}
	var ins []instruction
	for _, cg := range f.Comments {
		ln := fset.Position(cg.Pos()).Line
		raw := cg.Text()
		for line := range strings.Lines(raw) {
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "ignore") {
				continue
			}
			if line == "OK" && ins == nil {
				// so our return value will be non-nil
				ins = []instruction{}
				continue
			}
			switch extractDataMode(line) {
			case "json":
				jsonInst, err := extractInstructionFromJSON(strings.TrimPrefix(line, "json:"), ln)
				if err != nil {
					tb.Fatalf("At %v:%d: %v", filename, ln, err)
				}
				ins = append(ins, jsonInst)
			case "classic":
				match, err := extractPattern(line)
				if err != nil {
					tb.Fatalf("At %v:%d: %v", filename, ln, err)
				}
				matchLine := ln
				if i := strings.Index(line, "MATCH:"); i >= 0 {
					// This is a match for a different line.
					lns := strings.TrimPrefix(line[i:], "MATCH:")
					lns = lns[:strings.Index(lns, " ")] //nolint:gocritic // offBy1: false positive
					matchLine, err = strconv.Atoi(lns)
					if err != nil {
						tb.Fatalf("Bad match line number %q at %v:%d: %v", lns, filename, ln, err)
					}
				}
				var repl string
				if r, ok := extractReplacement(line); ok {
					repl = r
				}
				ins = append(ins, instruction{
					Line:        matchLine,
					Match:       match,
					Replacement: repl,
				})
			}
		}
	}
	return ins
}

func extractInstructionFromJSON(line string, lineNumber int) (instruction, error) {
	// Use the json.Unmarshal function to parse the JSON into the struct
	var jsonInst JSONInstruction
	if err := json.Unmarshal([]byte(line), &jsonInst); err != nil {
		return instruction{}, fmt.Errorf("parsing json instruction: %w", err)
	}

	ins := instruction{
		Match:      jsonInst.Match,
		Confidence: jsonInst.Confidence,
		Category:   jsonInst.Category,
		Line:       lineNumber,
	}
	return ins, nil
}

func extractDataMode(line string) string {
	if strings.HasPrefix(line, "json") {
		return "json"
	}
	if strings.Contains(line, "MATCH") {
		return "classic"
	}

	return ""
}

func extractPattern(line string) (string, error) {
	a, b := strings.Index(line, "/"), strings.LastIndex(line, "/")
	if a == -1 || a == b {
		return "", fmt.Errorf("malformed match instruction %q", line)
	}
	return line[a+1 : b], nil
}

func extractReplacement(line string) (string, bool) {
	// Look for this:  / -> `
	// (the end of a match and start of a backtick string),
	// and then the closing backtick.
	const start = "/ -> `"
	a, b := strings.Index(line, start), strings.LastIndex(line, "`")
	if a < 0 || a > b {
		return "", false
	}
	return line[a+len(start) : b], true
}

func srcLine(src []byte, p token.Position) string {
	// Run to end of line in both directions if not at line start/end.
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}

// TestLine tests srcLine function.
func TestLine(t *testing.T) {
	tests := []struct {
		src    string
		offset int
		want   string
	}{
		{"single line file", 5, "single line file"},
		{"single line file with newline\n", 5, "single line file with newline\n"},
		{"first\nsecond\nthird\n", 2, "first\n"},
		{"first\nsecond\nthird\n", 9, "second\n"},
		{"first\nsecond\nthird\n", 14, "third\n"},
		{"first\nsecond\nthird with no newline", 16, "third with no newline"},
		{"first byte\n", 0, "first byte\n"},
	}
	for _, test := range tests {
		got := srcLine([]byte(test.src), token.Position{Offset: test.offset})
		if got != test.want {
			t.Errorf("srcLine(%q, offset=%d) = %q, want %q", test.src, test.offset, got, test.want)
		}
	}
}

// exportedType reports whether typ is an exported type.
// It is imprecise, and will err on the side of returning true,
// such as for composite types.
func exportedType(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Named:
		// Builtin types have no package.
		return t.Obj().Pkg() == nil || t.Obj().Exported()
	case *types.Map:
		return exportedType(t.Key()) && exportedType(t.Elem())
	case interface {
		Elem() types.Type
	}: // array, slice, pointer, chan
		return exportedType(t.Elem())
	}
	// Be conservative about other types, such as struct, interface, etc.
	return true
}

// TestExportedType tests exportedType function.
func TestExportedType(t *testing.T) {
	tests := []struct {
		typString string
		exp       bool
	}{
		{"int", true},
		{"string", false}, // references the shadowed builtin "string"
		{"T", true},
		{"t", false},
		{"*T", true},
		{"*t", false},
		{"map[int]complex128", true},
	}
	for _, test := range tests {
		src := `package foo; type T int; type t int; type string struct{}`
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, "foo.go", src, 0)
		if err != nil {
			t.Fatalf("Parsing %q: %v", src, err)
		}
		// use the package name as package path
		config := &types.Config{}
		pkg, err := config.Check(file.Name.Name, fset, []*ast.File{file}, nil)
		if err != nil {
			t.Fatalf("Type checking %q: %v", src, err)
		}
		tv, err := types.Eval(fset, pkg, token.NoPos, test.typString)
		if err != nil {
			t.Errorf("types.Eval(%q): %v", test.typString, err)
			continue
		}
		if got := exportedType(tv.Type); got != test.exp {
			t.Errorf("exportedType(%v) = %t, want %t", tv.Type, got, test.exp)
		}
	}
}

// mkdirTempDotGit adds a temporary .git directory to the given root directory in testdata.
// We can't commit .git directly because of the Git restrictions.
func mkdirTempDotGit(t *testing.T, root string) {
	t.Helper()

	baseDir := filepath.Join("..", "testdata", root)
	dir, err := filepath.Abs(baseDir)
	if err != nil {
		t.Fatalf("Failed to resolve abs path: %v", err)
	}

	gitDir := filepath.Join(dir, ".git")
	if err := os.MkdirAll(gitDir, 0o750); err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}

	t.Cleanup(func() {
		err = os.RemoveAll(gitDir)
		if err != nil {
			t.Logf("Failed to remove %v: %v", gitDir, err)
		}
	})
}
