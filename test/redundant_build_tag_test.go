package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRedundantBuildTagRule(t *testing.T) {
	testRule(t, "redundant_build_tag", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleNoFailure(t *testing.T) {
	testRule(t, "redundant_build_tag_no_failure", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleGo116(t *testing.T) {
	testRule(t, "redundant_build_tag_go116", &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
}

func TestRedundantBuildTagRuleGo121(t *testing.T) {
	for _, fixture := range []string{
		"go1.21/redundant_build_tag",
		"go1.21/redundant_build_tag_go120",
		"go1.21/redundant_build_tag_go122",
		"go1.21/redundant_build_tag_version_linux",
	} {
		t.Run(fixture, func(t *testing.T) {
			testRule(t, fixture, &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
		})
	}
}

func TestRedundantBuildTagRuleFileName(t *testing.T) {
	for _, fixture := range []string{
		"redundant_build_tag_windows",
		"redundant_build_tag_windows_test",
		"redundant_build_tag_tab_windows",
		"redundant_build_tag_linux_amd64",
		"redundant_build_tag_or_linux",
		"redundant_build_tag_partly_darwin",
		"redundant_build_tag_negated_windows",
		"redundant_build_tag_plus_linux",
		"redundant_build_tag_no_failure_linux",
		"redundant_build_tag_no_failure_amd64",
		"redundant_build_tag_no_failure_implied_illumos",
	} {
		t.Run(fixture, func(t *testing.T) {
			testRule(t, fixture, &rule.RedundantBuildTagRule{}, &lint.RuleConfig{})
		})
	}
}
