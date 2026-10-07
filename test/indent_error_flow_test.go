package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestIndentErrorFlow(t *testing.T) {
	testRule(t, "indent_error_flow", &rule.IndentErrorFlowRule{})
	testRule(t, "indent_error_flow_scope", &rule.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "indent_error_flow_scope", &rule.IndentErrorFlowRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"preserve-scope"}})
}

func TestIndentErrorFlowConfigureResetsPreviousConfiguration(t *testing.T) {
	r := &rule.IndentErrorFlowRule{}
	testRule(t, "indent_error_flow_config_reset", r, &lint.RuleConfig{Arguments: lint.Arguments{"preserveScope"}})
	testRule(t, "indent_error_flow_config_reset_default", r)
}
