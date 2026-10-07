package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestSuperfluousElse(t *testing.T) {
	testRule(t, "superfluous_else", &rule.SuperfluousElseRule{})
	testRule(t, "superfluous_else_scope", &rule.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "superfluous_else_scope", &rule.SuperfluousElseRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}

func TestSuperfluousElseConfigureResetsPreviousConfiguration(t *testing.T) {
	r := &rule.SuperfluousElseRule{}
	testRule(t, "superfluous_else_config_reset", r, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "superfluous_else_config_reset_default", r)
}
