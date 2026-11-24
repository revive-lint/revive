package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageNaming_conventionName(t *testing.T) {
	testRule(t, "package_naming_mixed_caps", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_mixed_caps", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": false},
		},
	})
	testRule(t, "package_naming_mixed_caps_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": true},
		},
	})
	testRule(t, "package_naming_underscore", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_underscore", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": false},
		},
	})
	testRule(t, "package_naming_underscore_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-convention-name-check": true},
		},
	})
	testRule(t, "package_naming_convention_name_regex", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"skip-convention-name-check":  false,
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	testRule(t, "package_naming_convention_name_regex", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	testRule(t, "package_naming_convention_name_regex_test", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	testRule(t, "package_naming_convention_name_regex_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
	testRule(t, "package_naming_convention_name_regex_test_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"convention-name-check-regex": "^[a-z][a-zA-Z0-9]*$", // allow camel case
			},
		},
	})
}

func TestPackageNaming_topLevel(t *testing.T) {
	testRule(t, "package_naming_top_level_pkg", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_top_level_pkg", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-top-level-check": false},
		},
	})
	testRule(t, "package_naming_top_level_pkg_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-top-level-check": true},
		},
	})
}

func TestPackageNaming_badNames(t *testing.T) {
	testRule(t, "package_naming_bad_default", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_bad_default", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-default-bad-name-check": false},
		},
	})
	testRule(t, "package_naming_bad_default_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-default-bad-name-check": true},
		},
	})

	testRule(t, "package_naming_bad_extra", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-extra-bad-name": true},
		},
	})
	testRule(t, "package_naming_bad_extra_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_bad_extra_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-extra-bad-name": false},
		},
	})

	testRule(t, "package_naming_bad_user_defined", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"user-defined-bad-names": []any{"data"}},
		},
	})
	testRule(t, "package_naming_bad_user_defined_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_bad_user_defined_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"user-defined-bad-names": []any{}},
		},
	})
}

func TestPackageNaming_stdLibConflict(t *testing.T) {
	testRule(t, "package_naming_std_common_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_std_common_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-collision-with-common-std": false},
		},
	})
	testRule(t, "package_naming_std_common_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"skip-collision-with-common-std": true},
		},
	})
	testRule(t, "package_naming_std_all_conflict", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": true},
		},
	})
	testRule(t, "package_naming_std_all_common", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": true},
		},
	})
	testRule(t, "package_naming_std_all_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{})
	testRule(t, "package_naming_std_all_conflict_skip", &rule.PackageNamingRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{"check-collision-with-all-std": false},
		},
	})
}

func TestPackageNamingConfigureResetsPreviousConfiguration(t *testing.T) {
	tests := []struct {
		name              string
		configuredFixture string
		defaultFixture    string
		arguments         lint.Arguments
	}{
		{
			name:              "skip convention check",
			configuredFixture: "package_naming_mixed_caps_skip",
			defaultFixture:    "package_naming_mixed_caps",
			arguments:         lint.Arguments{map[string]any{"skipConventionNameCheck": true}},
		},
		{
			name:              "convention regex",
			configuredFixture: "package_naming_mixed_caps_skip",
			defaultFixture:    "package_naming_mixed_caps",
			arguments:         lint.Arguments{map[string]any{"conventionNameCheckRegex": "^[A-Z]"}},
		},
		{
			name:              "skip top-level check",
			configuredFixture: "package_naming_top_level_pkg_skip",
			defaultFixture:    "package_naming_top_level_pkg",
			arguments:         lint.Arguments{map[string]any{"skipTopLevelCheck": true}},
		},
		{
			name:              "skip default bad-name check",
			configuredFixture: "package_naming_bad_default_skip",
			defaultFixture:    "package_naming_bad_default",
			arguments:         lint.Arguments{map[string]any{"skipDefaultBadNameCheck": true}},
		},
		{
			name:              "check extra bad names",
			configuredFixture: "package_naming_bad_extra",
			defaultFixture:    "package_naming_bad_extra_skip",
			arguments:         lint.Arguments{map[string]any{"checkExtraBadName": true}},
		},
		{
			name:              "user-defined bad names",
			configuredFixture: "package_naming_bad_user_defined",
			defaultFixture:    "package_naming_bad_user_defined_skip",
			arguments:         lint.Arguments{map[string]any{"userDefinedBadNames": []any{"data"}}},
		},
		{
			name:              "skip common standard-library collision",
			configuredFixture: "package_naming_std_common_conflict_skip",
			defaultFixture:    "package_naming_std_common_conflict",
			arguments:         lint.Arguments{map[string]any{"skipCollisionWithCommonStd": true}},
		},
		{
			name:              "check all standard-library collisions",
			configuredFixture: "package_naming_std_all_conflict",
			defaultFixture:    "package_naming_std_all_conflict_skip",
			arguments:         lint.Arguments{map[string]any{"checkCollisionWithAllStd": true}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &rule.PackageNamingRule{}
			testRule(t, tt.configuredFixture, r, &lint.RuleConfig{Arguments: tt.arguments})
			testRule(t, tt.defaultFixture, r)
		})
	}
}
