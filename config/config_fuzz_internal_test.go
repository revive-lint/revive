package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// FuzzGetConfig checks that loading an arbitrary configuration file, including configuring every rule
// with arbitrary arguments, never panics. Errors are expected and ignored.
func FuzzGetConfig(f *testing.F) {
	for _, path := range []string{"../defaults.toml", "../revive.toml"} {
		data, err := os.ReadFile(path)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}
	f.Add([]byte("enableAllRules = true\n[rule.argument-limit]\narguments = [4]\n"))
	f.Add([]byte("[rule.line-length-limit]\narguments = [\"80\"]\n"))
	f.Add([]byte("[rule.var-naming]\narguments = [[\"ID\"], [\"VM\"], [{upperCaseConst = true}]]\n"))

	// Reused across iterations so that `go list std` (run for checkCollisionWithAllStd) runs at most once per worker:
	// Configure resets all of its state except the loaded std package names, which do not depend on the input.
	packageNaming := &rule.PackageNamingRule{}

	f.Fuzz(func(_ *testing.T, data []byte) {
		cfg := &lint.Config{}
		if err := parseConfig(data, cfg); err != nil {
			return
		}
		if err := validateConfig(cfg); err != nil {
			return
		}
		Normalize(cfg)

		for name, ruleConfig := range cfg.Rules {
			r := newConfigurableRule(name, packageNaming)
			if r == nil || ruleConfig.Disabled {
				continue
			}
			_ = r.Configure(ruleConfig.Arguments)
		}
	})
}

// newConfigurableRule returns a new instance of the named rule, so that no rule state is shared between fuzz iterations.
// It returns nil if the rule is unknown or not configurable.
func newConfigurableRule(name string, packageNaming *rule.PackageNamingRule) lint.ConfigurableRule {
	if name == packageNaming.Name() {
		return packageNaming
	}
	for _, r := range allRules {
		if r.Name() != name {
			continue
		}
		c, _ := reflect.TypeAssert[lint.ConfigurableRule](reflect.New(reflect.TypeOf(r).Elem()))
		return c
	}
	return nil
}
