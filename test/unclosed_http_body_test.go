package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnclosedHTTPBody(t *testing.T) {
	testRule(t, "unclosed_http_body", &rule.UnclosedHTTPBodyRule{})
}
