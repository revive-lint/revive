package test_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestDefer(t *testing.T) {
	testRule(t, "defer", &rule.DeferRule{})
}

func TestDeferLoopDisabled(t *testing.T) {
	testRule(t, "defer_loop_disabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"return", "recover", "callChain", "methodCall", "exitAfterDefer"}},
	})
	testRule(t, "defer_loop_disabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"return", "recover", "call-chain", "method-call", "exit-after-defer"}},
	})
}

func TestDeferOthersDisabled(t *testing.T) {
	testRule(t, "defer_only_loop_enabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"loop"}},
	})
}

func TestDeferExitAfterDefer(t *testing.T) {
	testRule(t, "defer_exit_after_defer", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{[]any{"exit-after-defer"}},
	})
}
