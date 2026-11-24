package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestDeferRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name:      "no arguments",
			arguments: lint.Arguments{},
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				[]any{
					"loop",
					"callChain",
					"methodCall",
					"return",
					"recover",
					"immediateRecover",
					"exitAfterDefer",
				},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				[]any{
					"loop",
					"callchain",
					"methodcall",
					"return",
					"recover",
					"immediaterecover",
					"exitAfterDefer",
				},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				[]any{
					"loop",
					"call-chain",
					"method-call",
					"return",
					"recover",
					"immediate-recover",
					"exit-after-defer",
				},
			},
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument 'invalid_argument' for 'defer' rule. Expecting []string, got string"),
		},
		{
			name: "invalid subcase type",
			arguments: lint.Arguments{
				[]any{123},
			},
			wantErr: errors.New("invalid argument '123' for 'defer' rule. Expecting string, got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.DeferRule

			err := r.Configure(tt.arguments)

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Configure() unexpected non-nil error %q", err)
				}
				return
			}
			if err == nil || err.Error() != tt.wantErr.Error() {
				t.Errorf("Configure() unexpected error: got %q, want %q", err, tt.wantErr)
			}
		})
	}
}
