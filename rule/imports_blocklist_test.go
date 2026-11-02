package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportsBlocklistRule_Configure(t *testing.T) {
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
			name:      "string arguments",
			arguments: lint.Arguments{"crypto/md5", "**/internal/**"},
		},
		{
			name: "path-scoped map",
			arguments: lint.Arguments{map[string]any{
				"paths":    []any{"**/cmd/http/**"},
				"packages": []any{"crypto/md5", "github.com/example/db/**"},
			}},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument to the imports-blocklist rule. Expecting a string or map, got int"),
		},
		{
			name:      "argument that does not compile",
			arguments: lint.Arguments{"("},
			wantErr:   errors.New("invalid package pattern \"(\": error parsing regexp: missing closing ): `(?m)\"(\"$`"),
		},
		{
			name: "map missing packages",
			arguments: lint.Arguments{map[string]any{
				"paths": []any{"**/http/**"},
			}},
			wantErr: errors.New("invalid argument #0 to the imports-blocklist rule: missing required option \"packages\""),
		},
		{
			name: "map with unknown option",
			arguments: lint.Arguments{map[string]any{
				"packages": []any{"crypto/md5"},
				"foo":      "bar",
			}},
			wantErr: errors.New("invalid argument #0 to the imports-blocklist rule: unknown option \"foo\""),
		},
		{
			name: "map with invalid path filter",
			arguments: lint.Arguments{map[string]any{
				"packages": []any{"crypto/md5"},
				"paths":    []any{"a**b.go"},
			}},
			wantErr: errors.New("invalid argument #0 to the imports-blocklist rule: invalid path filter \"a**b.go\": invalid file filter [a**b.go], invalid glob pattern"),
		},
		{
			name: "map with non-string package",
			arguments: lint.Arguments{map[string]any{
				"packages": []any{1},
			}},
			wantErr: errors.New("invalid argument #0 to the imports-blocklist rule: invalid value for option \"packages\": expected string, got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.ImportsBlocklistRule

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
