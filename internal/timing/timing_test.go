package timing_test

import (
	"bytes"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"testing"
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
		var buf bytes.Buffer
		logger := debugLogger(&buf)
		rules := timing.NewRules()
		slow := rules.Stat("slow-rule")
		fast := rules.Stat("fast-rule")
		rules.Stat("unused-rule")

		const calls = 10
		var wg sync.WaitGroup
		for range calls {
			wg.Go(func() {
				slow.Add(time.Second)
				fast.Add(time.Millisecond)
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

	t.Run("returns the same accumulator for the same rule", func(t *testing.T) {
		var buf bytes.Buffer
		rules := timing.NewRules()
		rules.Stat("some-rule").Add(time.Second)
		rules.Stat("some-rule").Add(time.Second)
		rules.Log(debugLogger(&buf))

		if got, want := buf.String(), "calls=2"; !strings.Contains(got, want) {
			t.Errorf("record %q does not contain %q", got, want)
		}
	})

	t.Run("nil logs nothing", func(t *testing.T) {
		var buf bytes.Buffer
		var rules *timing.Rules
		rules.Log(debugLogger(&buf))
		if got := buf.String(); got != "" {
			t.Errorf("expected no output, got %q", got)
		}
	})
}

func debugLogger(out *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
