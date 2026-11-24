package lint_test

import (
	"bytes"
	"go/token"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"github.com/mgechev/revive/lint"
)

type fakeRule struct {
	name     string
	failures []lint.Failure
}

var _ lint.Rule = (*fakeRule)(nil)

func (r *fakeRule) Name() string { return r.name }

func (r *fakeRule) Apply(*lint.File, lint.Arguments) []lint.Failure {
	return slices.Clone(r.failures)
}

// lintFile lints src as a single file named filename and returns the reported failures.
func lintFile(t *testing.T, filename, src string, rules []lint.Rule, config lint.Config, logger *slog.Logger) []lint.Failure {
	t.Helper()

	l := lint.New(func(string) ([]byte, error) { return []byte(src), nil }, 0)
	l.SetLogger(logger)

	failures, err := l.Lint([][]string{{filename}}, rules, config)
	if err != nil {
		t.Fatal("unexpected error from linting:", err)
	}

	var got []lint.Failure
	for failure := range failures {
		got = append(got, failure)
	}
	return got
}

func TestLint_internalFailureDoesNotAbortOtherRules(t *testing.T) {
	rules := []lint.Rule{
		&fakeRule{
			name: "typecheck-internal-failure-rule",
			failures: []lint.Failure{
				lint.NewInternalFailure("simulated type-check failure"),
			},
		},
		&fakeRule{
			name: "normal-rule",
			failures: []lint.Failure{
				{
					Confidence: 1,
					Failure:    "must reach the channel",
				},
			},
		},
	}

	cfg := lint.Config{
		Confidence: 0.8,
		Rules: lint.RulesConfig{
			"typecheck-internal-failure-rule": {},
			"normal-rule":                     {},
		},
	}

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelWarn}))

	got := lintFile(t, "test.go", "package foo\n", rules, cfg, logger)

	if len(got) != 1 {
		t.Fatalf("expected exactly 1 failure to be reported, got %d: %+v", len(got), got)
	}
	if got[0].RuleName != "normal-rule" {
		t.Errorf("expected failure from %q, got %q", "normal-rule", got[0].RuleName)
	}

	logged := logBuf.String()
	for _, want := range []string{
		"level=WARN",
		`msg="rule skipped due to internal failure"`,
		"rule=typecheck-internal-failure-rule",
		"file=test.go",
		`failure="simulated type-check failure"`,
	} {
		if !strings.Contains(logged, want) {
			t.Errorf("expected log output to contain %q, got %q", want, logged)
		}
	}
}

func positionAtLine(line int) token.Position {
	return token.Position{Filename: "test.go", Line: line}
}

func TestLint_disableDirectives(t *testing.T) {
	const lines = 6

	tests := []struct {
		name      string
		src       string
		wantLines []int
	}{
		{
			name:      "no directives",
			src:       "package foo\n// some comment\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:      "disable rule",
			src:       "package foo\n//revive:disable:rule1\n\n// 4\n\n// 6\n",
			wantLines: []int{1},
		},
		{
			name:      "disable rule twice",
			src:       "package foo\n//revive:disable:rule1\n\n//revive:disable:rule1\n\n// 6\n",
			wantLines: []int{1},
		},
		{
			name:      "enable rule",
			src:       "package foo\n//revive:enable:rule1\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:      "disable and enable rule",
			src:       "package foo\n//revive:disable:rule1\n\n//revive:enable:rule1\n\n// 6\n",
			wantLines: []int{1, 5, 6},
		},
		{
			name:      "disable-line rule",
			src:       "package foo\n//revive:disable-line:rule1\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 3, 4, 5, 6},
		},
		{
			name:      "enable-line rule",
			src:       "package foo\n//revive:enable-line:rule1\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:      "disable-next-line rule",
			src:       "package foo\n//revive:disable-next-line:rule1\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 2, 4, 5, 6},
		},
		{
			name:      "enable-next-line rule",
			src:       "package foo\n//revive:enable-next-line:rule1\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:      "directives for another rule are ignored",
			src:       "package foo\n//revive:disable:rule2\n\n// 4\n\n// 6\n",
			wantLines: []int{1, 2, 3, 4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The fake rule reports one failure on every line, so the surviving lines show which intervals are disabled.
			// Directives are separated by blank lines because a directive inside a comment group applies from the group's last line.
			r := &fakeRule{name: "rule1"}
			for line := 1; line <= lines; line++ {
				r.failures = append(r.failures, lint.Failure{
					Confidence: 1,
					Failure:    "failure",
					Position:   lint.FailurePosition{Start: positionAtLine(line), End: positionAtLine(line)},
				})
			}

			got := lintFile(t, "test.go", tt.src, []lint.Rule{r}, lint.Config{}, nil)

			gotLines := make([]int, 0, len(got))
			for _, failure := range got {
				gotLines = append(gotLines, failure.Position.Start.Line)
			}
			slices.Sort(gotLines)
			if !slices.Equal(gotLines, tt.wantLines) {
				t.Errorf("reported lines = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

// typeCheckRule triggers the type check of the package, as typed rules do.
type typeCheckRule struct{}

func (typeCheckRule) Name() string { return "type-check-rule" }

func (typeCheckRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	_ = file.Pkg.TypeCheck()
	return nil
}

func TestLint_timingLogs(t *testing.T) {
	src := "package foo\n\nfunc Foo() {}\n"
	rules := []lint.Rule{
		&fakeRule{name: "rule1", failures: []lint.Failure{{Confidence: 1, Failure: "failure"}}},
		typeCheckRule{},
	}

	t.Run("debug level reports every phase", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

		lintFile(t, "test.go", src, rules, lint.Config{}, logger)

		got := buf.String()
		want := []string{
			"phase=go-version-detection",
			"phase=package-read-parse",
			"phase=package-lint",
			"phase=scan-sortable",
			"phase=type-check",
			"phase=disabled-intervals",
			"phase=file-lint",
			"phase=lint",
			`msg="rule timing" rule=rule1`,
			`msg="rule timing" rule=type-check-rule`,
		}
		for _, w := range want {
			if !strings.Contains(got, w) {
				t.Errorf("logs do not contain %q:\n%s", w, got)
			}
		}
	})

	t.Run("no timing logs above debug level", func(t *testing.T) {
		var buf bytes.Buffer
		logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

		lintFile(t, "test.go", src, rules, lint.Config{}, logger)

		if got := buf.String(); got != "" {
			t.Errorf("expected no output, got %q", got)
		}
	})
}
