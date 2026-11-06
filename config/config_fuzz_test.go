package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mgechev/revive/config"
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

	f.Fuzz(func(t *testing.T, data []byte) {
		path := filepath.Join(t.TempDir(), "revive.toml")
		if err := os.WriteFile(path, data, 0o600); err != nil {
			t.Fatal(err)
		}
		cfg, err := config.GetConfig(path)
		if err != nil {
			return
		}
		_, _ = config.GetLintingRules(cfg, nil)
	})
}
