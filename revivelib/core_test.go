package revivelib_test

import (
	"strings"
	"testing"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
	"github.com/mgechev/revive/rule"
)

func TestReviveLint(t *testing.T) {
	// ARRANGE
	revive := getMockRevive(t)

	// ACT
	failures, err := revive.Lint(revivelib.Include("../testdata/if_return.go"))
	if err != nil {
		t.Fatal(err)
	}

	// ASSERT
	failureList := []lint.Failure{}

	for failure := range failures {
		failureList = append(failureList, failure)
	}

	const expected = 5

	got := len(failureList)
	if got != expected {
		t.Fatalf("Expected failures to have %d failures, but it has %d.", expected, got)
	}
}

func TestReviveInstancesHaveIndependentRuleConfiguration(t *testing.T) {
	_, err := revivelib.New(&lint.Config{
		Confidence: 0.8,
		Rules: map[string]lint.RuleConfig{
			"package-naming": {Arguments: lint.Arguments{map[string]any{"skipDefaultBadNameCheck": true}}},
		},
	}, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	second, err := revivelib.New(&lint.Config{
		Confidence: 0.8,
		Rules:      map[string]lint.RuleConfig{"package-naming": {}},
	}, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	failures, err := second.Lint(revivelib.Include("../testdata/package_naming_bad_default.go"))
	if err != nil {
		t.Fatal(err)
	}
	var found bool
	for failure := range failures {
		if failure.RuleName == "package-naming" && strings.Contains(failure.Failure, `don't use "util"`) {
			found = true
		}
	}
	if !found {
		t.Fatal("the second Revive instance inherited the first instance's package-naming configuration")
	}
}

func TestReviveFormat(t *testing.T) {
	t.Setenv("NO_COLOR", "true")
	// ARRANGE
	revive := getMockRevive(t)

	failuresChan, err := revive.Lint(revivelib.Include("../testdata/if_return.go"))
	if err != nil {
		t.Fatal(err)
	}

	// ACT
	failures, exitCode, err := revive.Format("stylish", failuresChan)
	// ASSERT
	if err != nil {
		t.Fatal(err)
	}

	errorMsgs := []string{
		"(91, 3)  https://revive.run/r#unreachable-code  unreachable code after this statement",
		"(98, 3)  https://revive.run/r#unreachable-code  unreachable code after this statement",
		"(15, 2)  https://revive.run/r#if-return         redundant if ...; err != nil check, just return error instead.",
		"(88, 3)  https://revive.run/r#if-return         redundant if ...; err != nil check, just return error instead.",
		"(95, 3)  https://revive.run/r#if-return         redundant if ...; err != nil check, just return error instead.",
	}
	for _, errorMsg := range errorMsgs {
		if !strings.Contains(failures, errorMsg) {
			t.Fatalf("Expected formatted failures\n'%s'\nto contain\n'%s', but it didn't.", failures, errorMsg)
		}
	}

	const expected = 1
	if exitCode != expected {
		t.Fatalf("Expected exit code to be %d, but it was %d.", expected, exitCode)
	}
}

type mockRule struct{}

func (*mockRule) Name() string {
	return "mock-rule"
}

func (*mockRule) Apply(_ *lint.File, _ lint.Arguments) []lint.Failure {
	return nil
}

func getMockRevive(t *testing.T) *revivelib.Revive {
	t.Helper()

	conf, err := config.GetConfig("../defaults.toml")
	if err != nil {
		t.Fatal(err)
	}

	revive, err := revivelib.New(
		conf,
		true,
		2048,
		revivelib.NewExtraRule(&rule.IfReturnRule{}, lint.RuleConfig{}),
		revivelib.NewExtraRule(&mockRule{}, lint.RuleConfig{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	return revive
}
