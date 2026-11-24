package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestModifiesValRec(t *testing.T) {
	testRule(t, "modifies_value_receiver", &rule.ModifiesValRecRule{})
}

func TestModifiesValRec_typeCheckDoesNotPanic(t *testing.T) {
	testRule(t, "modifies_value_receiver_issue_59", &rule.ModifiesValRecRule{})
}
