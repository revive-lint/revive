package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportsBlocklistOriginal(t *testing.T) {
	testRule(t, "imports_blocklist_original", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"crypto/md5", "crypto/sha1"},
	})
}

func TestImportsBlocklist(t *testing.T) {
	testRule(t, "imports_blocklist", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"},
	})
}

func TestImportsBlocklistPaths(t *testing.T) {
	cfg := &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"paths":    []any{"**/http/**"},
			"packages": []any{"crypto/md5"},
		}},
	}
	testRule(t, "imports_blocklist_paths/http/banned", &rule.ImportsBlocklistRule{}, cfg)
	testRule(t, "imports_blocklist_paths/grpc/allowed", &rule.ImportsBlocklistRule{}, cfg)
}

func BenchmarkImportsBlocklist(b *testing.B) {
	for b.Loop() {
		testRule(b, "imports_blocklist", &rule.ImportsBlocklistRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{"github.com/full/match", "wildcard/**/between", "wildcard/backward/**", "**/wildcard/forward", "full"},
		})
	}
}

func TestImportsBlocklistConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.ImportsBlocklistRule{}
	testRule(t, "imports_blocklist_original", r, &lint.RuleConfig{
		Arguments: lint.Arguments{
			"crypto/md5",
			map[string]any{
				"packages": []any{"crypto/sha1"},
			},
		},
	})
	testRule(t, "imports_blocklist_default", r)
}
