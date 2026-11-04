package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRedundantVarDecl(t *testing.T) {
	testRule(t, "redundant_var_decl", &rule.RedundantVarDeclRule{})
}

func TestRedundantVarDeclRealWorld(t *testing.T) {
	testRule(t, "redundant_var_decl_real_world", &rule.RedundantVarDeclRule{})
}

func TestRedundantVarDeclOnlyTypes(t *testing.T) {
	onlyTypes := []any{"error", "[]byte", "io.Reader", "fixtures.flag"}
	testRule(t, "redundant_var_decl_only_types", &rule.RedundantVarDeclRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"onlyTypes": onlyTypes}},
	})
	testRule(t, "redundant_var_decl_only_types", &rule.RedundantVarDeclRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"only-types": onlyTypes}},
	})
}

func TestRedundantVarDeclConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.RedundantVarDeclRule{}
	testRule(t, "redundant_var_decl_only_types", r, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"onlyTypes": []any{"error", "[]byte", "io.Reader", "fixtures.flag"}}},
	})
	testRule(t, "redundant_var_decl", r)
}
