package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestPackageDirectoryMismatchRule_Configure(t *testing.T) {
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
			arguments: lint.Arguments{map[string]any{
				"ignoreDirectories": []any{"testdata", "vendor"},
			}},
		},
		{
			name: "valid kebab-cased arguments",
			arguments: lint.Arguments{map[string]any{
				"ignore-directories": []any{"testdata"},
			}},
		},
		{
			name: "empty ignored directories",
			arguments: lint.Arguments{map[string]any{
				"ignoreDirectories": []any{},
			}},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument type: expected map[string]any, got int"),
		},
		{
			name: "unknown argument",
			arguments: lint.Arguments{map[string]any{
				"unknown": []any{"testdata"},
			}},
			wantErr: errors.New("unknown argument unknown for package-directory-mismatch rule"),
		},
		{
			name: "invalid ignored directories type",
			arguments: lint.Arguments{map[string]any{
				"ignoreDirectories": "testdata",
			}},
			wantErr: errors.New("invalid value testdata for argument ignoreDirectories of rule package-directory-mismatch, expected []string got string"),
		},
		{
			name: "invalid ignored directory element type",
			arguments: lint.Arguments{map[string]any{
				"ignoreDirectories": []any{123},
			}},
			wantErr: errors.New("invalid value in ignoreDirectories argument of rule package-directory-mismatch: expected string, got int"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.PackageDirectoryMismatchRule

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
