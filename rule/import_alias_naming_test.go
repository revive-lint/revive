package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestImportAliasNamingRule_Configure(t *testing.T) {
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
			name:      "valid string argument",
			arguments: lint.Arguments{"^[a-z][a-z0-9]*$"},
		},
		{
			name: "valid map arguments",
			arguments: lint.Arguments{map[string]any{
				"allowRegex": "^[a-z][a-z0-9]*$",
				"denyRegex":  "^v\\d+$",
			}},
		},
		{
			name: "valid map lowercased arguments",
			arguments: lint.Arguments{map[string]any{
				"allowregex": "^[a-z][a-z0-9]*$",
				"denyregex":  "^v\\d+$",
			}},
		},
		{
			name: "valid map kebab-cased arguments",
			arguments: lint.Arguments{map[string]any{
				"allow-regex": "^[a-z][a-z0-9]*$",
				"deny-regex":  "^v\\d+$",
			}},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New(`invalid argument '123' for 'import-alias-naming' rule. Expecting string or map[string]string, got int`),
		},
		{
			name:      "invalid string argument regex",
			arguments: lint.Arguments{"["},
			wantErr: errors.New(`invalid argument to the "allowRegex" option of the "import-alias-naming" rule. ` +
				"Expecting \"[\" to be a valid regular expression, got: error parsing regexp: missing closing ]: `[`"),
		},
		{
			name: "invalid map key",
			arguments: lint.Arguments{
				map[string]any{
					"unknownKey": "value",
				},
			},
			wantErr: errors.New(`invalid map key for 'import-alias-naming' rule. Expecting 'allowRegex' or 'denyRegex', got unknownKey`),
		},
		{
			name: "invalid allowRegex type",
			arguments: lint.Arguments{map[string]any{
				"allowRegex": 123,
			}},
			wantErr: errors.New(`invalid argument '123' for the "allowRegex" option of the "import-alias-naming" rule. Expecting string, got int`),
		},
		{
			name: "invalid denyRegex type",
			arguments: lint.Arguments{map[string]any{
				"denyRegex": 123,
			}},
			wantErr: errors.New(`invalid argument '123' for the "denyRegex" option of the "import-alias-naming" rule. Expecting string, got int`),
		},
		{
			name: "invalid denyRegex regex",
			arguments: lint.Arguments{map[string]any{
				"denyRegex": "[",
			}},
			wantErr: errors.New(`invalid argument to the "denyRegex" option of the "import-alias-naming" rule. ` +
				"Expecting \"[\" to be a valid regular expression, got: error parsing regexp: missing closing ]: `[`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.ImportAliasNamingRule

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
