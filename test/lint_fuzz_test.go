package test_test

import (
	"os"
	"reflect"
	"testing"

	goversion "github.com/hashicorp/go-version"
	"golang.org/x/mod/modfile"

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
		f.Add(data, false)
		f.Add(data, true)
	}

	cfg := &lint.Config{EnableAllRules: true, GoVersion: moduleGoVersion(f)}
	config.Normalize(cfg)

	f.Fuzz(func(t *testing.T, src []byte, isTest bool) {
		fileName := "fuzz.go"
		if isTest {
			fileName = "fuzz_test.go"
		}
		l := lint.New(func(string) ([]byte, error) { return src, nil }, 0)
		failures, err := l.Lint([][]string{{fileName}}, newRules(t), *cfg)
		if err != nil {
			t.Fatal(err)
		}
		for range failures { // drain the channel so that linting completes
		}
	})
}

// newRules returns new instances of all rules, so that no rule state is shared between fuzz iterations.
func newRules(t *testing.T) []lint.Rule {
	t.Helper()
	allRules := config.AllRules()
	rules := make([]lint.Rule, 0, len(allRules))
	for _, r := range allRules {
		r, _ := reflect.TypeAssert[lint.Rule](reflect.New(reflect.TypeOf(r).Elem()))
		if c, ok := r.(lint.ConfigurableRule); ok {
			if err := c.Configure(nil); err != nil {
				t.Fatal(err)
			}
		}
		rules = append(rules, r)
	}
	return rules
}

// moduleGoVersion returns the Go version declared in revive's go.mod.
func moduleGoVersion(tb testing.TB) *goversion.Version {
	tb.Helper()
	data, err := os.ReadFile("../go.mod")
	if err != nil {
		tb.Fatal(err)
	}
	mod, err := modfile.ParseLax("go.mod", data, nil)
	if err != nil {
		tb.Fatal(err)
	}
	return goversion.Must(goversion.NewVersion(mod.Go.Version))
}
