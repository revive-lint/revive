package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportAliasNaming(t *testing.T) {
	testRule(t, "import_alias_naming", &rule.ImportAliasNamingRule{})
	testRule(t, "import_alias_naming", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{},
		},
	})
}

func TestImportAliasNaming_CustomConfig(t *testing.T) {
	testRule(t, "import_alias_naming_custom_config", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{`^[a-z]+$`},
	})
}

func TestImportAliasNaming_CustomConfigWithMultipleRules(t *testing.T) {
	testRule(t, "import_alias_naming_custom_config_with_multiple_values", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"allowRegex": `^[a-z][a-z0-9]*$`,
				"denyRegex":  `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
	testRule(t, "import_alias_naming_custom_config_with_multiple_values", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"allow-regex": `^[a-z][a-z0-9]*$`,
				"deny-regex":  `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
}

func TestImportAliasNaming_CustomConfigWithOnlyDeny(t *testing.T) {
	testRule(t, "import_alias_naming_custom_config_with_only_deny", &rule.ImportAliasNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"denyRegex": `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
}

func TestImportAliasNaming_Reconfigure(t *testing.T) {
	r := &rule.ImportAliasNamingRule{}
	testRule(t, "import_alias_naming_custom_config_with_only_deny", r, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"denyRegex": `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
	testRule(t, "import_alias_naming", r)
	testRule(t, "import_alias_naming_custom_config_with_only_deny", r, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"denyRegex": `^((v\d+)|(v\d+alpha\d+))$`,
			},
		},
	})
	testRule(t, "import_alias_naming_custom_config", r, &lint.RuleConfig{
		Arguments: lint.Arguments{`^[a-z]+$`},
	})
}
