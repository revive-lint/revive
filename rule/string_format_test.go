package rule_test

import (
	"errors"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStringFormatRule_Configure(t *testing.T) {
	type argumentsTest struct {
		name      string
		arguments lint.Arguments
		wantErr   error
	}
	tests := []argumentsTest{
		{
			name: "Not a Slice",
			arguments: lint.Arguments{
				"this is not a slice",
			},
			wantErr: errors.New("invalid configuration for string-format: argument is not a slice [argument 0, option 0]"),
		},
		{
			name: "Missing Regex",
			arguments: lint.Arguments{
				[]any{
					"method[0]",
				},
			},
			wantErr: errors.New("invalid configuration for string-format: less than two slices found in argument, scope and regex are required [argument 0, option 0]"),
		},
		{
			name: "Bad Argument Type",
			arguments: lint.Arguments{
				[]any{
					1,
				},
			},
			wantErr: errors.New("invalid configuration for string-format: less than two slices found in argument, scope and regex are required [argument 0, option 0]"),
		},
		{
			name: "Empty Scope",
			arguments: lint.Arguments{
				[]any{
					"",
					"//",
				},
			},
			wantErr: errors.New("invalid configuration for string-format: empty scope provided [argument 0, option 0]"),
		},
		{
			name: "Small or Empty Regex",
			arguments: lint.Arguments{
				[]any{
					"method[1].a",
					"-",
				},
			},
			wantErr: errors.New("invalid configuration for string-format: regex is too small (regexes should begin and end with '/') [argument 0, option 1]"),
		},
		{
			name: "Bad Scope",
			arguments: lint.Arguments{
				[]any{
					"1.a",
					"//",
				},
			},
			wantErr: errors.New("failed to parse configuration for string-format: unable to parse rule scope [argument 0, option 0, scope index 0]"),
		},
		{
			name: "Bad Regex",
			arguments: lint.Arguments{
				[]any{
					"method[1].a",
					"/(/",
				},
			},
			wantErr: errors.New("failed to parse configuration for string-format: unable to compile /(/ as regexp [argument 0, option 1]"),
		},
		{
			name: "Sample Config",
			arguments: lint.Arguments{
				[]any{
					"core.WriteError[1].Message", "/^([^A-Z]$)/", "must not start with a capital letter",
				},
				[]any{
					"fmt.Errorf[0]", "/^|[^\\.!?]$/", "must not end in punctuation",
				},
				[]any{
					"panic", "/^[^\\n]*$/", "must not contain line breaks",
				},
			},
		},
		{
			name: "Underscores in Scope",
			arguments: lint.Arguments{
				[]any{
					"some_pkg._some_function_name[5].some_member",
					"//",
				},
			},
		},
		{
			name: "Underscores in Multiple Scopes",
			arguments: lint.Arguments{
				[]any{
					"fmt.Errorf[0],core.WriteError[1].Message",
					"//",
				},
			},
		},
		{
			name: "', ' Delimiter",
			arguments: lint.Arguments{
				[]any{
					"abc, mt.Errorf",
					"//",
				},
			},
		},
		{
			name: "' ,' Delimiter",
			arguments: lint.Arguments{
				[]any{
					"abc ,mt.Errorf",
					"//",
				},
			},
		},
		{
			name: "',   ' Delimiter",
			arguments: lint.Arguments{
				[]any{
					"abc,   mt.Errorf",
					"//",
				},
			},
		},
		{
			name: "',   ' Delimiter",
			arguments: lint.Arguments{
				[]any{
					"abc,   mt.Errorf",
					"//",
				},
			},
		},
		{
			name: "Empty Middle Scope",
			arguments: lint.Arguments{
				[]any{
					"abc, ,mt.Errorf",
					"//",
				},
			},
			wantErr: errors.New("failed to parse configuration for string-format: empty scope in rule scopes: [argument 0, option 0, scope index 1]"),
		},
		{
			name: "Empty First Scope",
			arguments: lint.Arguments{
				[]any{
					",mt.Errorf",
					"//",
				},
			},
			wantErr: errors.New("failed to parse configuration for string-format: empty scope in rule scopes: [argument 0, option 0, scope index 0]"),
		},
		{
			name: "Bad First Scope",
			arguments: lint.Arguments{
				[]any{
					"1.a,fmt.Errorf[0]",
					"//",
				},
			},
			wantErr: errors.New("failed to parse configuration for string-format: unable to parse rule scope [argument 0, option 0, scope index 0]"),
		},
		{
			name: "Bad Second Scope",
			arguments: lint.Arguments{
				[]any{
					"fmt.Errorf[0],1.a",
					"//",
				},
			},
			wantErr: errors.New("failed to parse configuration for string-format: unable to parse rule scope [argument 0, option 0, scope index 1]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.StringFormatRule

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
