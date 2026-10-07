package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestAddConstantWithDefaultArguments(t *testing.T) {
	testRule(t, "add_constant_default", &rule.AddConstantRule{}, &lint.RuleConfig{})
}

func TestAddConstantWithArguments(t *testing.T) {
	testRule(t, "add_constant", &rule.AddConstantRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"maxLitCount": "2",
			"allowStrs":   "\"\"",
			"allowInts":   "0,1,2",
			"allowFloats": "0.0,1.0",
			"ignoreFuncs": "os\\.(CreateFile|WriteFile|Chmod|FindProcess),\\.Println,ignoredFunc,\\.Info",
		}},
	})
	testRule(t, "add_constant", &rule.AddConstantRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{
			"max-lit-count": "2",
			"allow-strs":    `""`,
			"allow-ints":    "0,1,2",
			"allow-floats":  "0.0,1.0",
			"ignore-funcs":  `os\.(CreateFile|WriteFile|Chmod|FindProcess),\.Println,ignoredFunc,\.Info`,
		}},
	})
}

func TestAddConstantConfigureResetsIgnoredFunctions(t *testing.T) {
	r := &rule.AddConstantRule{}
	testRule(t, "add_constant_config_reset", r, &lint.RuleConfig{
		Arguments: lint.Arguments{map[string]any{"ignoreFuncs": `fmt\.Print`}},
	})
	testRule(t, "add_constant_config_reset_ok", r)
}
