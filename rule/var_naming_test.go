package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestVarNamingRule_Configure(t *testing.T) {
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
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"skipInitialismNameChecks": true,
					"upperCaseConst":           true,
					"initialismsAsWords":       true,
				}},
			},
		},
		{
			name: "valid lowercased arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"skipinitialismnamechecks": true,
					"uppercaseconst":           true,
					"initialismsaswords":       true,
				}},
			},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{
				[]any{"ID"},
				[]any{"VM"},
				[]any{map[string]any{
					"skip-initialism-name-checks": true,
					"upper-case-const":            true,
					"initialisms-as-words":        true,
				}},
			},
		},
		{
			name:      "invalid allowlist type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument to the var-naming rule. Expecting a allowlist of type slice with initialisms, got int"),
		},
		{
			name:      "invalid allowlist value type",
			arguments: lint.Arguments{[]any{123}},
			wantErr:   errors.New("invalid 123 values of the var-naming rule. Expecting slice of strings but got element of type []interface {}"),
		},
		{
			name:      "invalid blocklist type",
			arguments: lint.Arguments{[]any{"ID"}, 123},
			wantErr:   errors.New("invalid argument to the var-naming rule. Expecting a blocklist of type slice with initialisms, got int"),
		},
		{
			name:      "invalid third argument type",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, 123},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, got int"),
		},
		{
			name:      "invalid third argument slice size",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, []any{}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, of len==1, but 0"),
		},
		{
			name:      "invalid third argument first element type",
			arguments: lint.Arguments{[]any{"ID"}, []any{"VM"}, []any{123}},
			wantErr:   errors.New("invalid third argument to the var-naming rule. Expecting a options of type slice, of len==1, with map, but int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.VarNamingRule

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
