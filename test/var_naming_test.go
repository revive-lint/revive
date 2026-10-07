package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarNaming(t *testing.T) {
	testRule(t, "var_naming", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}},
	})
	for _, key := range []string{"skip-initialism-name-checks", "skipInitialismNameChecks", "skipinitialismnamechecks"} {
		testRule(t, "var_naming_skip_initialism_name_checks_true", &rule.VarNamingRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{
				[]any{},
				[]any{},
				[]any{map[string]any{key: true}},
			},
		})
	}
	testRule(t, "var_naming_skip_initialism_name_checks_false", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{},
			[]any{},
			[]any{map[string]any{"skip-initialism-name-checks": false}},
		},
	})
	testRule(t, "var_naming_allowlist_blocklist_skip_initialism_name_checks", &rule.VarNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{"ID"},
			[]any{"GRPC"},
			[]any{map[string]any{"skip-initialism-name-checks": true}},
		},
	})

	testRule(t, "var_naming_test", &rule.VarNamingRule{}, &lint.RuleConfig{})

	testRule(t, "var_naming_upper_case_const_false", &rule.VarNamingRule{}, &lint.RuleConfig{})
	for _, key := range []string{"upper-case-const", "upperCaseConst", "uppercaseconst"} {
		testRule(t, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{key: true}}},
		})
	}
}

func TestVarNamingConfigureResetsPreviousConfiguration(t *testing.T) {
	r := &rule.VarNamingRule{}
	testRule(t, "var_naming_config_reset", r, &lint.RuleConfig{
		Arguments: lint.Arguments{
			[]any{"ID"},
			[]any{"GRPC"},
			[]any{map[string]any{"skipInitialismNameChecks": true, "upperCaseConst": true}},
		},
	})
	testRule(t, "var_naming_config_reset_default", r)
}

func BenchmarkUpperCaseConstTrue(b *testing.B) {
	for b.Loop() {
		testRule(b, "var_naming_upper_case_const_true", &rule.VarNamingRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{[]any{}, []any{}, []any{map[string]any{"upperCaseConst": true}}},
		})
	}
}
