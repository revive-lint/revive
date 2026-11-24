package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCommentsDensity(t *testing.T) {
	testRule(t, "comments_density_1", &rule.CommentsDensityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(60)},
	})

	testRule(t, "comments_density_2", &rule.CommentsDensityRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(90)},
	})

	testRule(t, "comments_density_3", &rule.CommentsDensityRule{}, &lint.RuleConfig{})
}

func TestCommentsDensityConfigureResetsPreviousOptions(t *testing.T) {
	r := &rule.CommentsDensityRule{}
	testRule(t, "comments_density_2", r, &lint.RuleConfig{
		Arguments: lint.Arguments{int64(90)},
	})
	testRule(t, "comments_density_3", r)
}
