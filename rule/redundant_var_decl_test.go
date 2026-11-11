package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestRedundantVarDeclRule_Configure(t *testing.T) {
	tests := []struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}{
		{
			name: "no arguments",
		},
		{
			name:      "valid arguments",
			arguments: lint.Arguments{map[string]any{"onlyTypes": []any{"error", "io.Reader"}}},
		},
		{
			name:      "valid kebab-cased arguments",
			arguments: lint.Arguments{map[string]any{"only-types": []any{"error"}}},
		},
		{
			name:      "empty list",
			arguments: lint.Arguments{map[string]any{"onlyTypes": []any{}}},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{"error"},
			wantErr:   errors.New("invalid argument to the redundant-var-decl rule. Expecting a k,v map, got string"),
		},
		{
			name:      "invalid onlyTypes value",
			arguments: lint.Arguments{map[string]any{"onlyTypes": "error"}},
			wantErr:   errors.New(`invalid argument to the redundant-var-decl rule, []string expected for "onlyTypes". Got 'error' (string)`),
		},
		{
			name:      "invalid onlyTypes item",
			arguments: lint.Arguments{map[string]any{"onlyTypes": []any{"error", int64(1)}}},
			wantErr:   errors.New(`invalid argument to the redundant-var-decl rule, string expected in "onlyTypes". Got '1' (int64)`),
		},
		{
			name:      "unrecognized key",
			arguments: lint.Arguments{map[string]any{"types": []any{"error"}}},
			wantErr:   errors.New(`invalid argument to the redundant-var-decl rule: unknown option "types"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.RedundantVarDeclRule

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
