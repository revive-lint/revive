package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFilenameFormatRule_Configure(t *testing.T) {
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
			name:      "valid regexp argument",
			arguments: lint.Arguments{`^[_a-z][_a-z0-9]*\.go$`},
		},
		{
			name:      "too many arguments",
			arguments: lint.Arguments{`^a`, `^b`},
			wantErr:   errors.New(`rule "filename-format" expects only one argument, got 2 [^a ^b]`),
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New(`rule "filename-format" expects a string argument, got 123 of type int`),
		},
		{
			name:      "argument that does not compile",
			arguments: lint.Arguments{"("},
			wantErr:   errors.New("rule \"filename-format\" expects a valid regexp argument, got error for (: error parsing regexp: missing closing ): `(`"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.FilenameFormatRule

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
