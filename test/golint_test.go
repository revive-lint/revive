package test_test

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

var lintMatch = flag.String("lint.match", "", "restrict fixtures matches to this pattern")

var rules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.VarNamingRule{},
	&rule.IndentErrorFlowRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorNamingRule{},
	&rule.ErrorStringsRule{},
	&rule.ReceiverNamingRule{},
	&rule.IncrementDecrementRule{},
	&rule.ErrorReturnRule{},
	&rule.UnexportedReturnRule{},
	&rule.TimeNamingRule{},
	&rule.ContextKeysType{},
}

func TestAll(t *testing.T) {
	baseDir := "../testdata/golint/"

	for _, r := range rules {
		configureRule(t, r, nil)
	}

	rx, err := regexp.Compile(*lintMatch)
	if err != nil {
		t.Fatalf("Bad -lint.match value %q: %v", *lintMatch, err)
	}

	fis, err := os.ReadDir(baseDir)
	if err != nil {
		t.Fatalf("os.ReadDir: %v", err)
	}
	if len(fis) == 0 {
		t.Fatalf("no files in %v", baseDir)
	}
	for _, fi := range fis {
		if !rx.MatchString(fi.Name()) {
			continue
		}
		t.Run(fi.Name(), func(t *testing.T) {
			filePath := filepath.Join(baseDir, fi.Name())
			src, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed reading %s: %v", fi.Name(), err)
			}

			ins := parseInstructions(t, filePath, src)

			assertFailures(t, filePath, rules, lint.Config{}, ins)
		})
	}
}
