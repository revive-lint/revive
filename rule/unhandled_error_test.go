package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestUnhandledErrorRule_Configure(t *testing.T) {
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
			arguments: lint.Arguments{"fmt.Print.*", `os\.(Un)?Setenv`},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument to the unhandled-error rule. Expecting a string, got int"),
		},
		{
			name:      "empty argument",
			arguments: lint.Arguments{"  "},
			wantErr:   errors.New("invalid argument to the unhandled-error rule, expected regular expression must not be empty"),
		},
		{
			name:      "argument that does not compile",
			arguments: lint.Arguments{"("},
			wantErr:   errors.New("invalid argument to the unhandled-error rule: regexp \"(\" does not compile: error parsing regexp: missing closing ): `(`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.UnhandledErrorRule

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
