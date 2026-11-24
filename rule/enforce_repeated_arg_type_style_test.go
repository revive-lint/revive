package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceRepeatedArgTypeStyleRule_Configure(t *testing.T) {
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
			name: "valid arguments: short",
			arguments: lint.Arguments{
				"short",
			},
		},
		{
			name: "valid arguments",
			arguments: lint.Arguments{
				map[string]any{
					"funcArgStyle":    "full",
					"funcRetValStyle": "short",
				},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"funcargstyle":    "full",
					"funcretvalstyle": "short",
				},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				map[string]any{
					"func-arg-style":     "full",
					"func-ret-val-style": "short",
				},
			},
		},
		{
			name: "unrecognized key",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "someValue",
				},
			},
			wantErr: errors.New("invalid map key for 'enforce-repeated-arg-type-style' rule. Expecting 'funcArgStyle' or 'funcRetValStyle', got unknownKey"),
		},
		{
			name: "invalid argument type",
			arguments: lint.Arguments{
				123,
			},
			wantErr: errors.New("invalid argument '123' for 'enforce-repeated-arg-type-style' rule. Expecting string or map[string]string, got int"),
		},
		{
			name: "invalid argument when string",
			arguments: lint.Arguments{
				"invalid_argument",
			},
			wantErr: errors.New("invalid argument to the enforce-repeated-arg-type-style rule: invalid repeated arg type style: invalid_argument (expecting one of [any short full])"),
		},
		{
			name: "invalid funcArgStyle value",
			arguments: lint.Arguments{
				map[string]any{
					"funcArgStyle": 123,
				},
			},
			wantErr: errors.New("invalid map value type for 'enforce-repeated-arg-type-style' rule. Expecting string, got int"),
		},
		{
			name: "invalid funcRetValStyle value",
			arguments: lint.Arguments{
				map[string]any{
					"funcRetValStyle": 123,
				},
			},
			wantErr: errors.New("invalid map value '123' for 'enforce-repeated-arg-type-style' rule. Expecting string, got int"),
		},
		{
			name: "invalid funcArgStyle value: wrong string",
			arguments: lint.Arguments{
				map[string]any{
					"funcArgStyle": "invalid",
				},
			},
			wantErr: errors.New("invalid argument to the enforce-repeated-arg-type-style rule: invalid repeated arg type style: invalid (expecting one of [any short full])"),
		},
		{
			name: "invalid funcRetValStyle value: wrong string",
			arguments: lint.Arguments{
				map[string]any{
					"funcRetValStyle": "invalid",
				},
			},
			wantErr: errors.New("invalid argument to the enforce-repeated-arg-type-style rule: invalid repeated arg type style: invalid (expecting one of [any short full])"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.EnforceRepeatedArgTypeStyleRule

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
