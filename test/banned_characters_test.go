package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestBannedCharactersDefault(t *testing.T) {
	testRule(t, "banned_characters_default", &rule.BannedCharsRule{}, &lint.RuleConfig{})
}

// Test banned characters in a const, var and func names.
// One banned character is in the comment and should not be checked.
// One banned character from the list is not present in the fixture file.
func TestBannedCharacters(t *testing.T) {
	testRule(t, "banned_characters", &rule.BannedCharsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"Ω", "Σ", "σ", "1"},
	})
}

func TestBannedCharactersConfigureResetsPreviousList(t *testing.T) {
	r := &rule.BannedCharsRule{}
	testRule(t, "banned_characters", r, &lint.RuleConfig{
		Arguments: lint.Arguments{"Ω", "Σ", "σ", "1"},
	})
	testRule(t, "banned_characters_default", r)
}
