package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestFunctionLength_Configure(t *testing.T) {
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
			name:      "valid arguments",
			arguments: lint.Arguments{int64(50), int64(100)},
		},
		{
			name:      "too few arguments",
			arguments: lint.Arguments{int64(50)},
			wantErr:   errors.New(`invalid configuration for "function-length" rule, expected 2 arguments but got 1`),
		},
		{
			name:      "too many arguments",
			arguments: lint.Arguments{int64(50), int64(100), int64(150)},
			wantErr:   errors.New(`invalid configuration for "function-length" rule, expected 2 arguments but got 3`),
		},
		{
			name:      "invalid max statements type",
			arguments: lint.Arguments{"invalid", int64(100)},
			wantErr:   errors.New(`invalid configuration value for max statements in "function-length" rule; need int64 but got string`),
		},
		{
			name:      "negative max statements",
			arguments: lint.Arguments{int64(-1), int64(100)},
			wantErr:   errors.New(`the configuration value for max statements in "function-length" rule cannot be negative, got -1`),
		},
		{
			name:      "invalid max lines type",
			arguments: lint.Arguments{int64(50), "invalid"},
			wantErr:   errors.New(`invalid configuration value for max lines in "function-length" rule; need int64 but got string`),
		},
		{
			name:      "negative max lines",
			arguments: lint.Arguments{int64(50), int64(-1)},
			wantErr:   errors.New(`the configuration value for max lines in "function-length" rule cannot be negative, got -1`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.FunctionLength

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
