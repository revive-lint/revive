package test_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStringFormatConcurrentFiles(t *testing.T) {
	const fileCount = 32
	dir := t.TempDir()
	files := make([]string, 0, fileCount)
	for i := range fileCount {
		filename := filepath.Join(dir, fmt.Sprintf("file%d.go", i))
		source := fmt.Sprintf("package fixtures\nimport \"fmt\"\nfunc f%d() { fmt.Errorf(\"bad\") }\n", i)
		if err := os.WriteFile(filename, []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
		files = append(files, filename)
	}

	r := &rule.StringFormatRule{}
	arguments := lint.Arguments{[]any{"fmt.Errorf[0]", "/^good$/"}}
	if err := r.Configure(arguments); err != nil {
		t.Fatal(err)
	}
	config := lint.Config{
		Rules: map[string]lint.RuleConfig{"string-format": {Arguments: arguments}},
	}
	linter := lint.New(os.ReadFile, 0)
	failures, err := linter.Lint([][]string{files}, []lint.Rule{r}, config)
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for range failures {
		count++
	}
	if count != fileCount {
		t.Fatalf("got %d findings, want %d", count, fileCount)
	}
}

func TestStringFormat(t *testing.T) {
	testRule(t, "string_format", &rule.StringFormatRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{
				"stringFormatMethod1", // The first argument is checked by default
				"/^[A-Z]/",
				"must start with a capital letter",
			},

			[]any{
				"stringFormatMethod2[2].d",
				"/[^\\.]$/",
			}, // Must not end with a period
			[]any{
				"s.Method3[2]",
				"!/^[Tt][Hh]/",
				"must not start with 'th'",
			},
			[]any{
				"s.Method4", // same as before, but called from a struct
				"!/^[Ot][Tt]/",
				"must not start with 'ot'",
			},
		},
	})
}

func TestStringFormatDuplicatedStrings(t *testing.T) {
	testRule(t, "string_format_issue_1063", &rule.StringFormatRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{
			"fmt.Errorf[0],errors.New[0]",
			"/^([^A-Z]|$)/",
			"must not start with a capital letter",
		}},
	})
}
