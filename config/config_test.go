package config_test

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	goversion "github.com/hashicorp/go-version"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestGetConfig(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		for name, tc := range map[string]struct {
			confPath   string
			wantConfig lint.Config
		}{
			"default config": {
				wantConfig: lint.Config{
					IgnoreGeneratedHeader: false,
					Confidence:            0.8,
					Severity:              lint.SeverityWarning,
					EnableAllRules:        false,
					EnableDefaultRules:    false,
					Rules: lint.RulesConfig{
						"blank-imports": {
							Severity: lint.SeverityWarning,
						},
						"context-as-argument": {
							Severity: lint.SeverityWarning,
						},
						"context-keys-type": {
							Severity: lint.SeverityWarning,
						},
						"dot-imports": {
							Severity: lint.SeverityWarning,
						},
						"empty-block": {
							Severity: lint.SeverityWarning,
						},
						"error-naming": {
							Severity: lint.SeverityWarning,
						},
						"error-return": {
							Severity: lint.SeverityWarning,
						},
						"error-strings": {
							Severity: lint.SeverityWarning,
						},
						"errorf": {
							Severity: lint.SeverityWarning,
						},
						"exported": {
							Severity: lint.SeverityWarning,
						},
						"increment-decrement": {
							Severity: lint.SeverityWarning,
						},
						"indent-error-flow": {
							Severity: lint.SeverityWarning,
						},
						"package-comments": {
							Severity: lint.SeverityWarning,
						},
						"range": {
							Severity: lint.SeverityWarning,
						},
						"receiver-naming": {
							Severity: lint.SeverityWarning,
						},
						"redefines-builtin-id": {
							Severity: lint.SeverityWarning,
						},
						"superfluous-else": {
							Severity: lint.SeverityWarning,
						},
						"time-naming": {
							Severity: lint.SeverityWarning,
						},
						"unexported-return": {
							Severity: lint.SeverityWarning,
						},
						"unreachable-code": {
							Severity: lint.SeverityWarning,
						},
						"unused-parameter": {
							Severity: lint.SeverityWarning,
						},
						"var-declaration": {
							Severity: lint.SeverityWarning,
						},
						"var-naming": {
							Severity: lint.SeverityWarning,
						},
					},
					ErrorCode:   0,
					WarningCode: 0,
					Directives:  lint.DirectivesConfig{},
					Exclude:     []string{},
					GoVersion:   nil,
				},
			},
			"non-reg issue #470": {
				confPath: "issue-470.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					Severity:   lint.SeverityWarning,
					Rules: lint.RulesConfig{
						"add-constant": {
							Severity: lint.SeverityWarning,
							Arguments: lint.Arguments{
								map[string]any{
									"maxLitCount": "3",
									"allowStrs":   `"`,
									"allowFloats": "0.0,1.0,1.,2.0,2.",
									"allowInts":   "0,1,2",
								},
							},
						},
					},
				},
			},
			"config from file issue #585": {
				confPath: "issue-585.toml",
				wantConfig: lint.Config{
					Confidence: 0.0,
					Severity:   lint.SeverityWarning,
				},
			},
			"config from file default confidence issue #585": {
				confPath: "issue-585-default-confidence.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					Severity:   lint.SeverityWarning,
				},
			},
			"config from file go-version": {
				confPath: "go-version.toml",
				wantConfig: lint.Config{
					Confidence: 0.8,
					GoVersion:  goversion.Must(goversion.NewSemver("1.20.0")),
				},
			},
			"config from file ignore-generated-header": {
				confPath: "ignore-generated-header.toml",
				wantConfig: lint.Config{
					Confidence:            0.8,
					IgnoreGeneratedHeader: true,
				},
			},
			"config from file enable-default-rules": {
				confPath: "enable-default.toml",
				wantConfig: lint.Config{
					Confidence:            0.8,
					IgnoreGeneratedHeader: false,
					EnableDefaultRules:    true,
					Rules: lint.RulesConfig{
						"blank-imports":        {},
						"context-as-argument":  {},
						"context-keys-type":    {},
						"dot-imports":          {},
						"empty-block":          {},
						"error-naming":         {},
						"error-return":         {},
						"error-strings":        {},
						"errorf":               {},
						"exported":             {},
						"increment-decrement":  {},
						"indent-error-flow":    {},
						"package-comments":     {},
						"range":                {},
						"receiver-naming":      {},
						"redefines-builtin-id": {},
						"superfluous-else":     {},
						"time-naming":          {},
						"unexported-return":    {},
						"unreachable-code":     {},
						"unused-parameter":     {},
						"var-declaration":      {},
						"var-naming":           {},
					},
				},
			},
			"config with non-defaults": {
				confPath: "non-defaults.toml",
				wantConfig: lint.Config{
					Confidence:            0.5,
					Severity:              lint.SeverityError,
					IgnoreGeneratedHeader: true,
					EnableDefaultRules:    true,
					ErrorCode:             2,
					WarningCode:           1,
					Rules: lint.RulesConfig{
						"argument-limit": {
							Severity: lint.SeverityWarning,
							Exclude:  []string{"excluded/file.go"},
							Arguments: lint.Arguments{
								[]any{4},
							},
						},
						"blank-imports": {
							Disabled: true,
							Severity: lint.SeverityError,
						},
						"context-as-argument": {
							Severity: lint.SeverityError,
						},
						"context-keys-type": {
							Severity: lint.SeverityError,
						},
						"dot-imports": {
							Severity: lint.SeverityError,
						},
						"empty-block": {
							Severity: lint.SeverityError,
						},
						"error-naming": {
							Severity: lint.SeverityError,
						},
						"error-return": {
							Severity: lint.SeverityError,
						},
						"error-strings": {
							Severity: lint.SeverityError,
						},
						"errorf": {
							Severity: lint.SeverityError,
						},
						"exported": {
							Severity: lint.SeverityError,
							Arguments: lint.Arguments{
								"check-private-receivers", "disable-stuttering-check",
							},
							Exclude: []string{"excluded/file-exported.go"},
						},
						"increment-decrement": {
							Severity: lint.SeverityError,
						},
						"indent-error-flow": {
							Severity: lint.SeverityError,
						},
						"package-comments": {
							Severity: lint.SeverityError,
						},
						"range": {
							Severity: lint.SeverityError,
						},
						"receiver-naming": {
							Severity: lint.SeverityError,
						},
						"redefines-builtin-id": {
							Severity: lint.SeverityError,
						},
						"superfluous-else": {
							Severity: lint.SeverityError,
						},
						"time-naming": {
							Severity: lint.SeverityError,
						},
						"unexported-return": {
							Severity: lint.SeverityError,
						},
						"unreachable-code": {
							Severity: lint.SeverityError,
						},
						"unused-parameter": {
							Severity: lint.SeverityError,
						},
						"var-declaration": {
							Severity: lint.SeverityError,
						},
						"var-naming": {
							Severity: lint.SeverityError,
						},
					},
				},
			},
		} {
			t.Run(name, func(t *testing.T) {
				var cfgPath string
				if tc.confPath != "" {
					cfgPath = filepath.Join("testdata", tc.confPath)
				}

				cfg, err := config.GetConfig(cfgPath)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}
				if cfg.IgnoreGeneratedHeader != tc.wantConfig.IgnoreGeneratedHeader {
					t.Errorf("IgnoreGeneratedHeader: expected %v, got %v", tc.wantConfig.IgnoreGeneratedHeader, cfg.IgnoreGeneratedHeader)
				}
				if cfg.Confidence != tc.wantConfig.Confidence {
					t.Errorf("Confidence: expected %v, got %v", tc.wantConfig.Confidence, cfg.Confidence)
				}
				if cfg.Severity != tc.wantConfig.Severity {
					t.Errorf("Severity: expected %v, got %v", tc.wantConfig.Severity, cfg.Severity)
				}
				if cfg.EnableAllRules != tc.wantConfig.EnableAllRules {
					t.Errorf("EnableAllRules: expected %v, got %v", tc.wantConfig.EnableAllRules, cfg.EnableAllRules)
				}
				if cfg.EnableDefaultRules != tc.wantConfig.EnableDefaultRules {
					t.Errorf("EnableDefaultRules: expected %v, got %v", tc.wantConfig.EnableDefaultRules, cfg.EnableDefaultRules)
				}
				if cfg.ErrorCode != tc.wantConfig.ErrorCode {
					t.Errorf("ErrorCode: expected %v, got %v", tc.wantConfig.ErrorCode, cfg.ErrorCode)
				}
				if cfg.WarningCode != tc.wantConfig.WarningCode {
					t.Errorf("WarningCode: expected %v, got %v", tc.wantConfig.WarningCode, cfg.WarningCode)
				}
				if !tc.wantConfig.GoVersion.Equal(cfg.GoVersion) {
					t.Errorf("GoVersion: expected %v, got %v", tc.wantConfig.GoVersion, cfg.GoVersion)
				}

				if len(cfg.Exclude) != len(tc.wantConfig.Exclude) {
					t.Errorf("Exclude length: expected %v, got %v", len(tc.wantConfig.Exclude), len(cfg.Exclude))
				} else {
					for i, exclude := range tc.wantConfig.Exclude {
						if cfg.Exclude[i] != exclude {
							t.Errorf("Exclude[%d]: expected %v, got %v", i, exclude, cfg.Exclude[i])
						}
					}
				}

				if len(cfg.Rules) != len(tc.wantConfig.Rules) {
					t.Errorf("Rules count: expected %v, got %v", len(tc.wantConfig.Rules), len(cfg.Rules))
				}
				for ruleName, wantRule := range tc.wantConfig.Rules {
					gotRule, exists := cfg.Rules[ruleName]
					if !exists {
						t.Errorf("Rule %q: expected to exist, but not found", ruleName)
						continue
					}
					if gotRule.Disabled != wantRule.Disabled {
						t.Errorf("Rule %q Disabled: expected %v, got %v", ruleName, wantRule.Disabled, gotRule.Disabled)
					}
					if gotRule.Severity != wantRule.Severity {
						t.Errorf("Rule %q Severity: expected %v, got %v", ruleName, wantRule.Severity, gotRule.Severity)
					}
					if len(gotRule.Arguments) != len(wantRule.Arguments) {
						t.Errorf("Rule %q Arguments length: expected %v, got %v", ruleName, len(wantRule.Arguments), len(gotRule.Arguments))
					}
					if len(gotRule.Exclude) != len(wantRule.Exclude) {
						t.Errorf("Rule %q Exclude length: expected %v, got %v", ruleName, len(wantRule.Exclude), len(gotRule.Exclude))
					} else {
						for i, wantExclude := range wantRule.Exclude {
							if gotRule.Exclude[i] != wantExclude {
								t.Errorf("Rule %q Exclude[%d]: expected %v, got %v", ruleName, i, wantExclude, gotRule.Exclude[i])
							}
						}
					}
				}
				// Check for unexpected rules in actual config
				for ruleName := range cfg.Rules {
					if _, exists := tc.wantConfig.Rules[ruleName]; !exists {
						t.Errorf("Rule %q: found in actual config but not expected", ruleName)
					}
				}

				if len(cfg.Directives) != len(tc.wantConfig.Directives) {
					t.Errorf("Directives count: expected %v, got %v", len(tc.wantConfig.Directives), len(cfg.Directives))
				}
				for directiveName, wantDirective := range tc.wantConfig.Directives {
					gotDirective, exists := cfg.Directives[directiveName]
					if !exists {
						t.Errorf("Directive %q: expected to exist, but not found", directiveName)
						continue
					}
					if gotDirective.Severity != wantDirective.Severity {
						t.Errorf("Directive %q Severity: expected %v, got %v", directiveName, wantDirective.Severity, gotDirective.Severity)
					}
				}
				// Check for unexpected directives in actual config
				for directiveName := range cfg.Directives {
					if _, exists := tc.wantConfig.Directives[directiveName]; !exists {
						t.Errorf("Directive %q: found in actual config but not expected", directiveName)
					}
				}
			})
		}

		t.Run("rule-level file filter excludes", func(t *testing.T) {
			cfg, err := config.GetConfig("testdata/rule-level-exclude-850.toml")
			if err != nil {
				t.Fatal("should be valid config")
			}
			r1 := cfg.Rules["r1"]
			if len(r1.Exclude) > 0 {
				t.Fatal("r1 should have empty excludes")
			}
			r2 := cfg.Rules["r2"]
			if len(r2.Exclude) != 1 {
				t.Fatal("r2 should have exclude set")
			}
			if !r2.MustExclude("some/file.go") {
				t.Fatal("r2 should be initialized and exclude some/file.go")
			}
			if r2.MustExclude("some/any-other.go") {
				t.Fatal("r2 should not exclude some/any-other.go")
			}
		})
	})

	t.Run("failure", func(t *testing.T) {
		for name, tc := range map[string]struct {
			confPath  string
			wantError string
		}{
			"unknown file": {
				confPath:  "unknown",
				wantError: "cannot read the config file",
			},
			"malformed file": {
				confPath:  "malformed.toml",
				wantError: "cannot parse the config file",
			},
			"invalid exclude pattern": {
				confPath:  "invalid-exclude-pattern.toml",
				wantError: "error in config of rule [var-naming]",
			},
			"enable-all-rules and enable-default-rules both set": {
				confPath:  "enable-all-and-default.toml",
				wantError: "config options enable-all-rules and enable-default-rules cannot be combined",
			},
			"same option with different casing": {
				confPath:  "duplicate-option.toml",
				wantError: "refer to the same option",
			},
		} {
			t.Run(name, func(t *testing.T) {
				_, err := config.GetConfig(filepath.Join("testdata", tc.confPath))

				if err == nil {
					t.Fatalf("Expected error %q, got none", tc.wantError)
				}
				if !strings.Contains(err.Error(), tc.wantError) {
					t.Errorf("Unexpected error: want %q, got: %q", tc.wantError, err)
				}
			})
		}
	})
}

func TestGetConfig_OptionCasing(t *testing.T) {
	want := lint.Config{
		IgnoreGeneratedHeader: true,
		Confidence:            0.5,
		Severity:              lint.SeverityError,
		EnableDefaultRules:    true,
		ErrorCode:             2,
		WarningCode:           1,
		GoVersion:             goversion.Must(goversion.NewSemver("1.20")),
	}

	for _, confPath := range []string{
		"options-camelCase.toml",
		"options-kebab-case.toml",
		"options-lowercase.toml",
	} {
		t.Run(confPath, func(t *testing.T) {
			cfg, err := config.GetConfig(filepath.Join("testdata", confPath))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if cfg.IgnoreGeneratedHeader != want.IgnoreGeneratedHeader {
				t.Errorf("IgnoreGeneratedHeader: expected %v, got %v", want.IgnoreGeneratedHeader, cfg.IgnoreGeneratedHeader)
			}
			if cfg.Confidence != want.Confidence {
				t.Errorf("Confidence: expected %v, got %v", want.Confidence, cfg.Confidence)
			}
			if cfg.Severity != want.Severity {
				t.Errorf("Severity: expected %v, got %v", want.Severity, cfg.Severity)
			}
			if cfg.EnableDefaultRules != want.EnableDefaultRules {
				t.Errorf("EnableDefaultRules: expected %v, got %v", want.EnableDefaultRules, cfg.EnableDefaultRules)
			}
			if cfg.ErrorCode != want.ErrorCode {
				t.Errorf("ErrorCode: expected %v, got %v", want.ErrorCode, cfg.ErrorCode)
			}
			if cfg.WarningCode != want.WarningCode {
				t.Errorf("WarningCode: expected %v, got %v", want.WarningCode, cfg.WarningCode)
			}
			if !want.GoVersion.Equal(cfg.GoVersion) {
				t.Errorf("GoVersion: expected %v, got %v", want.GoVersion, cfg.GoVersion)
			}
		})
	}
}

func TestGetConfig_EnableAllRulesCasing(t *testing.T) {
	for _, confPath := range []string{
		"enable-all-camel-case.toml",
		"enable-all-kebab-case.toml",
		"enable-all-lowercase.toml",
	} {
		t.Run(confPath, func(t *testing.T) {
			cfg, err := config.GetConfig(filepath.Join("testdata", confPath))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if !cfg.EnableAllRules {
				t.Error("EnableAllRules: expected true, got false")
			}
		})
	}
}

func TestGetConfig_RuleOptionCasing(t *testing.T) {
	cfg, err := config.GetConfig(filepath.Join("testdata", "rule-option-casing.toml"))
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	for _, ruleName := range []string{"blank-imports", "dot-imports"} {
		if !cfg.Rules[ruleName].Disabled {
			t.Errorf("Rule %q: expected disabled=true, got false", ruleName)
		}
	}
	for _, ruleName := range []string{"argument-limit", "cyclomatic"} {
		if len(cfg.Rules[ruleName].Arguments) != 1 {
			t.Errorf("Rule %q: expected 1 argument, got %v", ruleName, cfg.Rules[ruleName].Arguments)
		}
	}
	for _, ruleName := range []string{"var-naming", "error-strings"} {
		if len(cfg.Rules[ruleName].Exclude) != 1 {
			t.Errorf("Rule %q: expected 1 exclude, got %v", ruleName, cfg.Rules[ruleName].Exclude)
		}
	}
}

func TestGetLintingRules(t *testing.T) {
	const (
		// len of defaultRules
		defaultRulesCount = 23
		// len of allRules: update this when adding new rules
		allRulesCount = 105
	)

	tt := map[string]struct {
		confPath          string
		wantRulesCount    int
		wantEnabledRules  []string
		wantDisabledRules []string
		wantErr           string
	}{
		"no rules": {
			confPath:       "no-rules.toml",
			wantRulesCount: 0,
			wantDisabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
		},
		"enable-all-rules without disabled rules": {
			confPath:       "enable-all.toml",
			wantRulesCount: allRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
		},
		"enable-all-rules with 2 disabled rules": {
			confPath:       "enable-all-but2.toml",
			wantRulesCount: allRulesCount - 2,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
			wantDisabledRules: []string{
				"exported",   // default rule
				"cyclomatic", // non-default rule
			},
		},
		"enable-default-rules without disabled rules": {
			confPath:       "enable-default.toml",
			wantRulesCount: defaultRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
			},
			wantDisabledRules: []string{
				"deep-exit", // non-default rule
			},
		},
		"enable-default-rules with 2 disabled rules": {
			confPath:       "enable-default-but2.toml",
			wantRulesCount: defaultRulesCount - 2,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
			},
			wantDisabledRules: []string{
				"exported",          // default rule
				"indent-error-flow", // default rule
			},
		},
		"enable-default-rules plus 1 non-default rule": {
			confPath:       "enable-default-plus1.toml",
			wantRulesCount: defaultRulesCount + 1,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"cyclomatic",       // non-default rule
			},
			wantDisabledRules: []string{
				"deep-exit", // non-default rule
			},
		},
		"enable-default-rules plus rule already in defaults": {
			confPath:       "enable-default-plus-default-rule.toml",
			wantRulesCount: defaultRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"exported",         // default rule
			},
			wantDisabledRules: []string{
				"deep-exit", // non-default rule
			},
		},
		"enable-all-rules plus rule already in all": {
			confPath:       "enable-all-with-rule.toml",
			wantRulesCount: allRulesCount,
			wantEnabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
				"cyclomatic",       // non-default rule
			},
		},
		"enable 2 rules": {
			confPath:       "enable2.toml",
			wantRulesCount: 2,
			wantEnabledRules: []string{
				"exported",   // default rule
				"cyclomatic", // non-default rule
			},
			wantDisabledRules: []string{
				"var-declaration",  // default rule
				"package-comments", // default rule
				"deep-exit",        // non-default rule
			},
		},
		"former name of a renamed rule": {
			confPath: "renamed-rule.toml",
			wantErr:  `rule "imports-blacklist" was renamed to "imports-blocklist", update the configuration`,
		},
		"var-naming configure error": {
			confPath: "var-naming-configure-error.toml",
			wantErr:  `cannot configure rule: "var-naming": invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got string`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := config.GetConfig(filepath.Join("testdata", tc.confPath))
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := config.GetLintingRules(cfg, []lint.Rule{})
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Fatalf("Expected error %q, got %q", tc.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error\n\t%v", err)
			}

			ruleNames := make([]string, len(rules))
			for i, rule := range rules {
				ruleNames[i] = rule.Name()
			}
			slices.Sort(ruleNames)

			if len(rules) != tc.wantRulesCount {
				t.Errorf("Expected %v enabled linting rules got: %v. Got rules: %v", tc.wantRulesCount, len(rules), ruleNames)
			}
			for _, wantEnabledRule := range tc.wantEnabledRules {
				if !slices.Contains(ruleNames, wantEnabledRule) {
					t.Errorf("Expected enabled rule %q not found. Got enabled rules: %v", wantEnabledRule, ruleNames)
				}
			}
			for _, wantDisabledRule := range tc.wantDisabledRules {
				if slices.Contains(ruleNames, wantDisabledRule) {
					t.Errorf("Expected disabled rule %q not found. Got enabled rules: %v", wantDisabledRule, ruleNames)
				}
			}
		})
	}
}

// extraRule is an externally provided rule, as revivelib users can register.
type extraRule struct {
	name string
}

func (r *extraRule) Name() string { return r.name }

func (*extraRule) Apply(_ *lint.File, _ lint.Arguments) []lint.Failure { return nil }

func TestGetLintingRules_FormerRuleName(t *testing.T) {
	cfg, err := config.GetConfig(filepath.Join("testdata", "renamed-rule.toml"))
	if err != nil {
		t.Fatalf("Unexpected error while loading conf: %v", err)
	}

	// An extra rule registered under the former name takes over it:
	// embedders can provide their own rule to keep configurations using the former name working.
	rules, err := config.GetLintingRules(cfg, []lint.Rule{&extraRule{name: "imports-blacklist"}})
	if err != nil {
		t.Fatalf("Unexpected error\n\t%v", err)
	}
	if want := []string{"imports-blacklist"}; !slices.Equal(ruleNames(rules), want) {
		t.Errorf("Expected rules %v, got %v", want, ruleNames(rules))
	}
}

func TestGetGlobalSeverity(t *testing.T) {
	tt := map[string]struct {
		confPath               string
		wantGlobalSeverity     string
		particularRule         lint.Rule
		wantParticularSeverity string
	}{
		"enable 2 rules with one specific severity": {
			confPath:               "testdata/enable2-one-specific-severity.toml",
			wantGlobalSeverity:     "warning",
			particularRule:         &rule.CyclomaticRule{},
			wantParticularSeverity: "error",
		},
		"enable-all-rules with one specific severity": {
			confPath:               "testdata/enable-all-one-specific-severity.toml",
			wantGlobalSeverity:     "error",
			particularRule:         &rule.DeepExitRule{},
			wantParticularSeverity: "warning",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			cfg, err := config.GetConfig(tc.confPath)
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			rules, err := config.GetLintingRules(cfg, []lint.Rule{})
			if err != nil {
				t.Fatalf("Unexpected error while loading conf: %v", err)
			}
			for _, r := range rules {
				ruleName := r.Name()
				ruleCfg := cfg.Rules[ruleName]
				ruleSeverity := string(ruleCfg.Severity)
				switch ruleName {
				case tc.particularRule.Name():
					if tc.wantParticularSeverity != ruleSeverity {
						t.Fatalf("Expected Severity %v for rule %v, got %v", tc.wantParticularSeverity, ruleName, ruleSeverity)
					}
				default:
					if tc.wantGlobalSeverity != ruleSeverity {
						t.Fatalf("Expected Severity %v for rule %v, got %v", tc.wantGlobalSeverity, ruleName, ruleSeverity)
					}
				}
			}
		})
	}
}

func TestGetFormatter(t *testing.T) {
	t.Run("default formatter", func(t *testing.T) {
		formatter, err := config.GetFormatter("")
		if err != nil {
			t.Fatalf("Unexpected error %q", err)
		}
		if formatter == nil || formatter.Name() != "default" {
			t.Errorf("Expected formatter %q, got %v", "default", formatter)
		}
	})
	t.Run("unknown formatter", func(t *testing.T) {
		_, err := config.GetFormatter("unknown")
		if err == nil || err.Error() != "unknown formatter unknown" {
			t.Errorf("Expected error %q, got: %q", "unknown formatter unknown", err)
		}
	})
	t.Run("checkstyle formatter", func(t *testing.T) {
		formatter, err := config.GetFormatter("checkstyle")
		if err != nil {
			t.Fatalf("Unexpected error: %q", err)
		}
		if formatter == nil || formatter.Name() != "checkstyle" {
			t.Errorf("Expected formatter %q, got %v", "checkstyle", formatter)
		}
	})
}

func ruleNames(rules []lint.Rule) []string {
	names := make([]string, len(rules))
	for i, r := range rules {
		names[i] = r.Name()
	}
	slices.Sort(names)
	return names
}

func TestAllRules(t *testing.T) {
	names := ruleNames(config.AllRules())

	if len(names) == 0 {
		t.Fatal("AllRules returned no rules")
	}
	if compacted := slices.Compact(slices.Clone(names)); len(compacted) != len(names) {
		t.Errorf("AllRules should not contain duplicates, got %v", names)
	}
	for _, want := range []string{"argument-limit", "cyclomatic", "exported", "var-naming"} {
		if !slices.Contains(names, want) {
			t.Errorf("AllRules should contain %q", want)
		}
	}
	// Every default rule is also part of all rules.
	for _, name := range ruleNames(config.DefaultRules()) {
		if !slices.Contains(names, name) {
			t.Errorf("AllRules is missing default rule %q", name)
		}
	}
}

func TestDefaultRules(t *testing.T) {
	names := ruleNames(config.DefaultRules())

	if len(names) == 0 {
		t.Fatal("DefaultRules returned no rules")
	}
	if compacted := slices.Compact(slices.Clone(names)); len(compacted) != len(names) {
		t.Errorf("DefaultRules should not contain duplicates, got %v", names)
	}
	for _, want := range []string{"blank-imports", "exported", "var-naming"} {
		if !slices.Contains(names, want) {
			t.Errorf("DefaultRules should contain %q", want)
		}
	}
	// Default rules are a strict subset of all rules.
	if len(names) >= len(config.AllRules()) {
		t.Errorf("DefaultRules (%d) should be fewer than AllRules (%d)", len(names), len(config.AllRules()))
	}
}

func TestEnabledRules(t *testing.T) {
	t.Run("returns only enabled rules", func(t *testing.T) {
		cfg := &lint.Config{
			Rules: lint.RulesConfig{
				"var-naming":     {},
				"exported":       {Disabled: true},
				"argument-limit": {},
				"cyclomatic":     {Disabled: true},
			},
		}

		got := ruleNames(config.EnabledRules(cfg))
		want := []string{"argument-limit", "var-naming"}
		if !slices.Equal(got, want) {
			t.Errorf("EnabledRules: expected %v, got %v", want, got)
		}
	})

	t.Run("ignores unknown rule names", func(t *testing.T) {
		cfg := &lint.Config{
			Rules: lint.RulesConfig{
				"imports-blacklist": {}, // former name of imports-blocklist
			},
		}

		if got := config.EnabledRules(cfg); len(got) != 0 {
			t.Errorf("EnabledRules: expected none, got %v", ruleNames(got))
		}
	})

	t.Run("empty config has no enabled rules", func(t *testing.T) {
		if got := config.EnabledRules(&lint.Config{}); len(got) != 0 {
			t.Errorf("EnabledRules: expected none, got %v", got)
		}
	})

	t.Run("nil config has no enabled rules", func(t *testing.T) {
		if got := config.EnabledRules(nil); len(got) != 0 {
			t.Errorf("EnabledRules: expected none, got %v", got)
		}
	})
}

func TestDefaultConfidence(t *testing.T) {
	if config.DefaultConfidence != 0.8 {
		t.Errorf("DefaultConfidence: expected 0.8, got %v", config.DefaultConfidence)
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()

	if cfg.Confidence != config.DefaultConfidence {
		t.Errorf("Confidence: expected %v, got %v", config.DefaultConfidence, cfg.Confidence)
	}
	if cfg.Severity != lint.SeverityWarning {
		t.Errorf("Severity: expected %v, got %v", lint.SeverityWarning, cfg.Severity)
	}

	var names []string
	for name := range cfg.Rules {
		names = append(names, name)
	}
	slices.Sort(names)
	if want := ruleNames(config.DefaultRules()); !slices.Equal(names, want) {
		t.Errorf("Default config rules: expected %v, got %v", want, names)
	}
}

func TestNormalize(t *testing.T) {
	t.Run("enable-default-rules adds default rule entries", func(t *testing.T) {
		cfg := &lint.Config{EnableDefaultRules: true}

		config.Normalize(cfg)

		if got, want := ruleNames(config.EnabledRules(cfg)), ruleNames(config.DefaultRules()); !slices.Equal(got, want) {
			t.Errorf("expected default rules %v, got %v", want, got)
		}
	})

	t.Run("enable-all-rules adds all rule entries", func(t *testing.T) {
		cfg := &lint.Config{EnableAllRules: true}

		config.Normalize(cfg)

		if got, want := ruleNames(config.EnabledRules(cfg)), ruleNames(config.AllRules()); !slices.Equal(got, want) {
			t.Errorf("expected all rules %v, got %v", want, got)
		}
	})

	t.Run("severity is propagated to rules and directives without their own", func(t *testing.T) {
		cfg := &lint.Config{
			Severity: lint.SeverityError,
			Rules: lint.RulesConfig{
				"exported":   {},
				"var-naming": {Severity: lint.SeverityWarning},
			},
			Directives: lint.DirectivesConfig{
				"specify-disable-reason": {},
			},
		}

		config.Normalize(cfg)

		if got := cfg.Rules["exported"].Severity; got != lint.SeverityError {
			t.Errorf("exported severity: expected %q, got %q", lint.SeverityError, got)
		}
		if got := cfg.Rules["var-naming"].Severity; got != lint.SeverityWarning {
			t.Errorf("var-naming severity should be preserved: expected %q, got %q", lint.SeverityWarning, got)
		}
		if got := cfg.Directives["specify-disable-reason"].Severity; got != lint.SeverityError {
			t.Errorf("directive severity: expected %q, got %q", lint.SeverityError, got)
		}
	})

	t.Run("nil rules map is initialized", func(t *testing.T) {
		cfg := &lint.Config{}

		config.Normalize(cfg)

		if cfg.Rules == nil {
			t.Error("Rules map should be initialized")
		}
	})
}
