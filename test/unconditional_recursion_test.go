package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnconditionalRecursion(t *testing.T) {
	testRule(t, "unconditional_recursion", &rule.UnconditionalRecursionRule{})
	testRule(t, "unconditional_recursion_import_alias", &rule.UnconditionalRecursionRule{})
}
