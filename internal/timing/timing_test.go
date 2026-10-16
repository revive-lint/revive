package timing_test

import (
	"bytes"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/mgechev/revive/internal/timing"
)

func TestEnabled(t *testing.T) {
	var buf bytes.Buffer
	tests := map[string]struct {
		logger *slog.Logger
		want   bool
	}{
		"nil logger":     {logger: nil, want: false},
		"discard logger": {logger: slog.New(slog.DiscardHandler), want: false},
		"info logger":    {logger: slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})), want: false},
		"debug logger":   {logger: debugLogger(&buf), want: true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := timing.Enabled(tt.logger); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrack(t *testing.T) {
	t.Run("logs the phase with its attributes", func(t *testing.T) {
		var buf bytes.Buffer
		timing.Track(debugLogger(&buf), "some-phase", "files", 3)("error", nil)

		got := buf.String()
		for _, want := range []string{`level=DEBUG`, `msg="phase completed"`, `phase=some-phase`, `duration=`, `files=3`, `error=<nil>`} {
			if !strings.Contains(got, want) {
				t.Errorf("log %q does not contain %q", got, want)
			}
		}
	})

	t.Run("logs nothing when debug logging is disabled", func(t *testing.T) {
		var buf bytes.Buffer
		timing.Track(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})), "some-phase")()
		timing.Track(nil, "some-phase")()

		if got := buf.String(); got != "" {
			t.Errorf("expected no output, got %q", got)
		}
	})
}

func TestValue(t *testing.T) {
	t.Run("resolves the value when the record is logged", func(t *testing.T) {
		var buf bytes.Buffer
		timing.Track(debugLogger(&buf), "some-phase")("files", timing.Value(func() any { return 3 }))

		if got, want := buf.String(), "files=3"; !strings.Contains(got, want) {
			t.Errorf("log %q does not contain %q", got, want)
		}
	})

	t.Run("is not resolved when debug logging is disabled", func(t *testing.T) {
		resolved := false
		timing.Track(slog.New(slog.DiscardHandler), "some-phase")("files", timing.Value(func() any {
			resolved = true
			return 3
		}))

		if resolved {
			t.Error("value resolved although debug logging is disabled")
		}
	})
}

func TestRules(t *testing.T) {
	t.Run("accumulates concurrent rule applications", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			var buf bytes.Buffer
			logger := debugLogger(&buf)
			rules := timing.NewRules(logger, []string{"slow-rule", "fast-rule", "unused-rule"})

			const calls = 10
			var wg sync.WaitGroup
			for range calls {
				wg.Go(func() {
					done := rules.Track("slow-rule")
					time.Sleep(time.Second)
					done()

					rules.Track("fast-rule")()
				})
			}
			wg.Wait()

			rules.Log(logger)

			got := buf.String()
			lines := slices.Collect(strings.Lines(got))
			if len(lines) != 2 {
				t.Fatalf("expected 2 records (rules without calls are not reported), got %d: %q", len(lines), got)
			}
			for _, want := range []string{`rule=slow-rule`, `total=10s`, `calls=10`, `avg=1s`} {
				if !strings.Contains(lines[0], want) {
					t.Errorf("record %q does not contain %q", lines[0], want)
				}
			}
			if want := `rule=fast-rule`; !strings.Contains(lines[1], want) {
				t.Errorf("records are not sorted by total duration: %q", got)
			}
		})
	})

	t.Run("ignores unknown rules", func(t *testing.T) {
		var buf bytes.Buffer
		logger := debugLogger(&buf)
		rules := timing.NewRules(logger, []string{"known-rule"})
		rules.Track("unknown-rule")()
		rules.Log(logger)

		if got := buf.String(); got != "" {
			t.Errorf("expected no output, got %q", got)
		}
	})

	t.Run("is nil when debug logging is disabled", func(t *testing.T) {
		var buf bytes.Buffer
		rules := timing.NewRules(slog.New(slog.DiscardHandler), []string{"some-rule"})
		if rules != nil {
			t.Fatalf("expected nil Rules, got %+v", rules)
		}

		rules.Track("some-rule")()
		rules.Log(debugLogger(&buf))
		if got := buf.String(); got != "" {
			t.Errorf("expected no output, got %q", got)
		}
	})
}

func debugLogger(out *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
