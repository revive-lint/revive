package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestLintFilenameFormat(t *testing.T) {
	testRule(t, "filename_ok_default", &rule.FilenameFormatRule{}, &lint.RuleConfig{})

	testRule(t, "filenamе_with_non_ascii_char", &rule.FilenameFormatRule{}, &lint.RuleConfig{})

	testRule(t, "filename_with_underscores", &rule.FilenameFormatRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"^[A-Za-z][A-Za-z0-9]*.go$"}})
}

func TestFilenameFormatConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.FilenameFormatRule{}
	testRule(t, "filename_with_underscores", r, &lint.RuleConfig{
		Arguments: lint.Arguments{"^[A-Za-z][A-Za-z0-9]*.go$"},
	})
	testRule(t, "filename_ok_default", r)
}
