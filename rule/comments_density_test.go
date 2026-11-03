package rule_test

import (
	"errors"
	"go/ast"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestCommentsDensityRule_Configure(t *testing.T) {
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
			name:      "integer argument",
			arguments: lint.Arguments{int64(15)},
		},
		{
			name:      "invalid argument type",
			arguments: lint.Arguments{"invalid"},
			wantErr:   errors.New(`invalid argument for "comments-density" rule: argument should be an int, got string`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.CommentsDensityRule

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

func TestCommentsDensityRule_Apply(t *testing.T) {
	t.Run("empty AST or empty content returns nil", func(t *testing.T) {
		var r rule.CommentsDensityRule
		_ = r.Configure(lint.Arguments{int64(50)})

		f := &lint.File{}
		failures := r.Apply(f, nil)
		if len(failures) != 0 {
			t.Errorf("expected 0 failures, got %d", len(failures))
		}

		f2 := &lint.File{AST: &ast.File{}}
		failures2 := r.Apply(f2, nil)
		if len(failures2) != 0 {
			t.Errorf("expected 0 failures, got %d", len(failures2))
		}
	})

	tests := []struct {
		name         string
		content      string
		minDensity   int64
		wantFailures int
		wantMsg      string
	}{
		{
			name:         "only code",
			content:      "package p\n",
			minDensity:   50,
			wantFailures: 1,
			wantMsg:      "the file has a comment density of  0% (0 comment lines for 1 code lines) but expected a minimum of 50%",
		},
		{
			name: "mixed comments and code meets threshold",
			content: `package p // line 1: code + comment
// line 2: comment
var x int // line 3: code + comment
`,
			minDensity:   60,
			wantFailures: 0, // 3 comment lines, 2 code lines: 3/(2+3) = 60% >= 60%
		},
		{
			name: "mixed comments and code below threshold",
			content: `package p // line 1: code + comment
var x int // line 2: code + comment
var y int
`,
			minDensity:   50,
			wantFailures: 1, // 2 comment lines, 3 code lines: 2/(3+2) = 40% < 50%
			wantMsg:      "the file has a comment density of 40% (2 comment lines for 3 code lines) but expected a minimum of 50%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r rule.CommentsDensityRule
			if err := r.Configure(lint.Arguments{tt.minDensity}); err != nil {
				t.Fatalf("Configure failed: %v", err)
			}

			reader := func(_ string) ([]byte, error) {
				return []byte(tt.content), nil
			}
			l := lint.New(reader, 0)
			config := lint.Config{
				Rules: map[string]lint.RuleConfig{
					r.Name(): {
						Arguments: lint.Arguments{tt.minDensity},
					},
				},
			}

			failures, err := l.Lint([][]string{{"test.go"}}, []lint.Rule{&r}, config)
			if err != nil {
				t.Fatalf("Lint failed: %v", err)
			}

			var gotFailures []string
			for f := range failures {
				gotFailures = append(gotFailures, f.Failure)
			}

			if len(gotFailures) != tt.wantFailures {
				t.Fatalf("got %d failures, want %d: %v", len(gotFailures), tt.wantFailures, gotFailures)
			}

			if tt.wantFailures > 0 && gotFailures[0] != tt.wantMsg {
				t.Errorf("failure message got %q, want %q", gotFailures[0], tt.wantMsg)
			}
		})
	}
}
