package test_test

import (
	"os"
	"testing"

	goversion "github.com/hashicorp/go-version"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
)

// FuzzLint checks that running all rules over arbitrary Go source never panics.
// Failures reported by the rules are expected and ignored.
func FuzzLint(f *testing.F) {
	for _, name := range []string{"argument_limit", "var_naming", "unused_param", "early_return", "struct_tag"} {
		data, err := os.ReadFile("../testdata/" + name + ".go")
		if err != nil {
			f.Fatal(err)
		}
		f.Add(data)
	}

	cfg := &lint.Config{EnableAllRules: true, GoVersion: goversion.Must(goversion.NewVersion("1.26"))}
	config.Normalize(cfg)
	rules, err := config.GetLintingRules(cfg, nil)
	if err != nil {
		f.Fatal(err)
	}

	f.Fuzz(func(t *testing.T, src []byte) {
		l := lint.New(func(string) ([]byte, error) { return src, nil }, 0)
		failures, err := l.Lint([][]string{{"fuzz.go"}}, rules, *cfg)
		if err != nil {
			t.Fatal(err)
		}
		count := 0
		for range failures {
			count++ // drain the channel so that linting completes
		}
		t.Logf("%d failures", count)
	})
}
