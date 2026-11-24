package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceSliceStyle_any(t *testing.T) {
	testRule(t, "enforce_slice_style_any", &rule.EnforceSliceStyleRule{})
}

func TestEnforceSliceStyle_make(t *testing.T) {
	testRule(t, "enforce_slice_style_make", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"make"},
	})
}

func TestEnforceSliceStyle_literal(t *testing.T) {
	testRule(t, "enforce_slice_style_literal", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"literal"},
	})
}

func TestEnforceSliceStyle_nil(t *testing.T) {
	testRule(t, "enforce_slice_style_nil", &rule.EnforceSliceStyleRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"nil"},
	})
}

func TestEnforceSliceStyleConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.EnforceSliceStyleRule{}
	testRule(t, "enforce_slice_style_nil", r, &lint.RuleConfig{
		Arguments: lint.Arguments{"nil"},
	})
	testRule(t, "enforce_slice_style_any", r)
}
