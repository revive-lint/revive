package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnhandledError(t *testing.T) {
	testRule(t, "unhandled_error", &rule.UnhandledErrorRule{})
}

func TestUnhandledErrorWithIgnoreList(t *testing.T) {
	testRule(t, "unhandled_error_w_ignorelist", &rule.UnhandledErrorRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			`unhandledError1`,
			`fmt\.Print`,
			`os\.(Create|WriteFile|Chmod)`,
			`net\..*`,
			`bytes\.Buffer\.Write`,
			`fixtures\.unhandledErrorStruct2\.reterr`,
		},
	})
}

func TestUnhandledErrorConfigureResetsPreviousIgnoreList(t *testing.T) {
	r := &rule.UnhandledErrorRule{}
	testRule(t, "unhandled_error_config_reset", r, &lint.RuleConfig{
		Arguments: lint.Arguments{`fmt\.Print`},
	})
	testRule(t, "unhandled_error_config_reset_ok", r)
}
