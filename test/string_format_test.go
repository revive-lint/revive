package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

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

func TestStringFormatConfigureResetsPreviousRules(t *testing.T) {
	r := &rule.StringFormatRule{}
	testRule(t, "string_format_config_reset", r, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"fmt.Errorf[0]", "/^good$/", "old configuration was applied"}},
	})
	testRule(t, "string_format_config_reset_ok", r, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"errors.New[0]", "/^good$/"}},
	})
}
