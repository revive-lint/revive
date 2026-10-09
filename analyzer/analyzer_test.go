package analyzer_test

import (
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/mgechev/revive/analyzer"
	"github.com/mgechev/revive/lint"
)

func TestAnalyzer(t *testing.T) {
	conf := &lint.Config{
		Confidence: 0.8,
		Rules: lint.RulesConfig{
			"errorf":            {},
			"line-length-limit": {Arguments: lint.Arguments{int64(150)}},
			"package-comments":  {},
			"unhandled-error":   {},
			"var-naming":        {},
		},
	}

	revive, err := analyzer.New(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results := analysistest.Run(t, analysistest.TestData(), revive, "a")

	const wantFixPrefix = "\treturn fmt.Errorf(\"failed: %d\", 42)"
	foundErrorfFix := false
	for _, result := range results {
		for _, diagnostic := range result.Diagnostics {
			if diagnostic.Category != "errorf" {
				continue
			}

			if len(diagnostic.SuggestedFixes) != 1 {
				t.Fatalf("expected 1 suggested fix for errorf diagnostic, got %d", len(diagnostic.SuggestedFixes))
			}
			edits := diagnostic.SuggestedFixes[0].TextEdits
			if len(edits) != 1 {
				t.Fatalf("expected 1 text edit for errorf suggested fix, got %d", len(edits))
			}
			if got := string(edits[0].NewText); !strings.HasPrefix(got, wantFixPrefix) {
				t.Errorf("expected errorf suggested fix to start with %q, got %q", wantFixPrefix, got)
			}
			foundErrorfFix = true
		}
	}
	if !foundErrorfFix {
		t.Error("expected an errorf diagnostic with a suggested fix, found none")
	}
}
