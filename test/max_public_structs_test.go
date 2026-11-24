package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestMaxPublicStructs(t *testing.T) {
	testRule(t, "max_public_structs", &rule.MaxPublicStructsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(1)},
	})
}

func TestMaxPublicStructsDefaultConfig(t *testing.T) {
	testRule(t, "max_public_structs_ok", &rule.MaxPublicStructsRule{}, &lint.RuleConfig{})
}

func TestMaxPublicStructsConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.MaxPublicStructsRule{}
	testRule(t, "max_public_structs", r, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(1)},
	})
	testRule(t, "max_public_structs_ok", r)
}
