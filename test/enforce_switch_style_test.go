package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSwitchStyle(t *testing.T) {
	testRule(t, "enforce_switch_style", &rule.EnforceSwitchStyleRule{})
	testRule(t, "enforce_switch_style_allow_no_default", &rule.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-no-default"},
	})
	testRule(t, "enforce_switch_style_allow_not_last", &rule.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-default-not-last"},
	})
	testRule(t, "enforce_switch_style_allow_no_default_allow_not_last", &rule.EnforceSwitchStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-no-default", "allow-default-not-last"},
	})
}

func TestEnforceSwitchStyleConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.EnforceSwitchStyleRule{}
	testRule(t, "enforce_switch_style_allow_no_default_allow_not_last", r, &lint.RuleConfig{
		Arguments: lint.Arguments{"allow-no-default", "allow-default-not-last"},
	})
	testRule(t, "enforce_switch_style", r)
}
