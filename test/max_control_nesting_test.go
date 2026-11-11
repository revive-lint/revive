package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestMaxControlNestingDefault(t *testing.T) {
	testRule(t, "max_control_nesting_default", &rule.MaxControlNestingRule{}, &lint.RuleConfig{})
}

func TestMaxControlNesting(t *testing.T) {
	testRule(t, "max_control_nesting", &rule.MaxControlNestingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(2)},
	})
}

func TestMaxControlNestingRange(t *testing.T) {
	testRule(t, "max_control_nesting_range", &rule.MaxControlNestingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(3)},
	})
}

func TestMaxControlNestingConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.MaxControlNestingRule{}
	testRule(t, "max_control_nesting", r, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(2)},
	})
	testRule(t, "max_control_nesting_default", r)
}
