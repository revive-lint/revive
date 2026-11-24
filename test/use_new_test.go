package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUseNew(t *testing.T) {
	testRule(t, "use_new", &rule.UseNewRule{})
	testRule(t, "go1.26/use_new", &rule.UseNewRule{})
}
