package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCommentSpacings(t *testing.T) {
	testRule(t, "comment_spacings", &rule.CommentSpacingsRule{}, &lint.RuleConfig{})
	testRule(t, "comment_spacings_custom_directive", &rule.CommentSpacingsRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"myOwnDirective:", "+optional"},
	})
}

func TestCommentSpacingsConfigureResetsPreviousConfiguration(t *testing.T) {
	r := &rule.CommentSpacingsRule{}
	testRule(t, "comment_spacings_custom_directive", r, &lint.RuleConfig{
		Arguments: lint.Arguments{"myOwnDirective:", "+optional"},
	})
	testRule(t, "comment_spacings_config_reset", r)
}
