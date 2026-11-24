package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestIdenticalSwitchBranches(t *testing.T) {
	testRule(t, "identical_switch_branches", &rule.IdenticalSwitchBranchesRule{})
	for _, key := range []string{"allow-identical-default", "allowIdenticalDefault", "allowidenticaldefault"} {
		testRule(t, "identical_switch_branches_allow_identical_default", &rule.IdenticalSwitchBranchesRule{}, &lint.RuleConfig{
			Arguments: lint.Arguments{map[string]any{key: true}},
		})
	}
}

func TestIdenticalSwitchBranchesConfigureResetsPreviousConfiguration(t *testing.T) {
	r := &rule.IdenticalSwitchBranchesRule{}
	testRule(t, "identical_switch_branches_allow_identical_default", r, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"allowIdenticalDefault": true}},
	})
	testRule(t, "identical_switch_branches", r)
}
