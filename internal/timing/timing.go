// Package timing measures the duration of revive's execution phases.
//
// Measurements are reported as debug log records, so they are only collected when debug logging is enabled
// (see the REVIVE_LOG_LEVEL environment variable). Otherwise every function here does nothing.
package timing

import (
	"cmp"
	"context"
	"log/slog"
	"slices"
	"sync/atomic"
	"time"
)

// Done ends a measurement started by [Track].
// The given attributes, if any, are appended to the emitted log record.
type Done func(attrs ...any)

var noop Done = func(...any) {}

// Value defers computing an attribute value until the record is logged.
type Value func() any

// LogValue implements [slog.LogValuer].
func (v Value) LogValue() slog.Value { return slog.AnyValue(v()) }

// Enabled reports whether logger records debug messages.
func Enabled(logger *slog.Logger) bool {
	return logger != nil && logger.Enabled(context.Background(), slog.LevelDebug)
}

// Track starts measuring the duration of the given phase.
func Track(logger *slog.Logger, phase string, attrs ...any) Done {
	if !Enabled(logger) {
		return noop
	}

	start := time.Now()
	return func(extra ...any) {
		args := slices.Concat(
			[]any{"phase", phase, "duration", time.Since(start)},
			attrs,
			extra,
		)
		logger.Debug("phase completed", args...)
	}
}

// Rules accumulates the time spent in each rule across all linted files.
type Rules struct {
	stats map[string]*RuleStat
}

// RuleStat accumulates the time spent in one rule. It is safe for concurrent use.
type RuleStat struct {
	totalNs atomic.Int64
	calls   atomic.Int64
}

// NewRules creates an empty accumulator.
func NewRules() *Rules {
	return &Rules{stats: map[string]*RuleStat{}}
}

// Stat returns the accumulator of the named rule, creating it if needed.
// It is not safe for concurrent use: get every accumulator before applying the rules.
func (r *Rules) Stat(name string) *RuleStat {
	stat, ok := r.stats[name]
	if !ok {
		stat = &RuleStat{}
		r.stats[name] = stat
	}
	return stat
}

// Add records one application of the rule that took d.
func (s *RuleStat) Add(d time.Duration) {
	s.totalNs.Add(int64(d))
	s.calls.Add(1)
}

// Log reports the accumulated per-rule durations, slowest rule first.
func (r *Rules) Log(logger *slog.Logger) {
	if r == nil {
		return
	}

	type entry struct {
		name  string
		total time.Duration
		calls int64
	}

	entries := make([]entry, 0, len(r.stats))
	for name, stat := range r.stats {
		calls := stat.calls.Load()
		if calls == 0 {
			continue
		}
		entries = append(entries, entry{
			name:  name,
			total: time.Duration(stat.totalNs.Load()),
			calls: calls,
		})
	}

	slices.SortFunc(entries, func(a, b entry) int {
		return cmp.Or(cmp.Compare(b.total, a.total), cmp.Compare(a.name, b.name))
	})

	for _, e := range entries {
		logger.Debug("rule timing", "rule", e.name, "total", e.total, "calls", e.calls, "avg", e.total/time.Duration(e.calls))
	}
}
