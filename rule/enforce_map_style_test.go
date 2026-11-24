package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceMapStyleRule_Configure(t *testing.T) {
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
			name:      "empty string argument",
			arguments: lint.Arguments{""},
		},
		{
			name:      "any argument",
			arguments: lint.Arguments{"any"},
		},
		{
			name:      "make argument",
			arguments: lint.Arguments{"make"},
		},
		{
			name:      "literal argument",
			arguments: lint.Arguments{"literal"},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{123},
			wantErr:   errors.New("invalid argument '123' for 'enforce-map-style' rule. Expecting string, got int"),
		},
		{
			name:      "unknown style",
			arguments: lint.Arguments{"invalid"},
			wantErr:   errors.New("invalid argument to the enforce-map-style rule: invalid map style: invalid (expecting one of [any make literal])"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.EnforceMapStyleRule

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
