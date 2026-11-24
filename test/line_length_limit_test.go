package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLineLengthLimitDefault(t *testing.T) {
	testRule(t, "line_length_limit_default", &rule.LineLengthLimitRule{}, &lint.RuleConfig{})
}

func TestLineLengthLimit(t *testing.T) {
	testRule(t, "line_length_limit", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(100)},
	})
}

func TestLineLengthLimitWithMaxOption(t *testing.T) {
	testRule(t, "line_length_limit", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"max": int64(100)}},
	})
}

func TestLineLengthLimitExcludes(t *testing.T) {
	testRule(t, "line_length_limit_excludes", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"max":      int64(80),
			"excludes": []any{`^\s*//go:generate `, `https?://`},
		}},
	})
	testRule(t, "line_length_limit_excludes", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"Max":      int64(80),
			"excludes": []any{`^\s*//go:generate `, `https?://`},
		}},
	})
}

func TestLineLengthLimitLongLine(t *testing.T) {
	testRule(t, "line_length_limit_long_line", &rule.LineLengthLimitRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(100)},
	})
}
