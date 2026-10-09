package test_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// TestUnexportedReturnMixedPackages verifies that the unexported-return rule produces
// stable and correct results when a directory contains source files,
// internal test files (package pub) and external test files (package pub_test).
//
// This mirrors what dots.ResolvePackages returns at runtime: GoFiles + TestGoFiles + XTestGoFiles
// collapsed into a single []string passed to lint.Linter.Lint.
func TestUnexportedReturnMixedPackages(t *testing.T) {
	baseDir := filepath.Join("..", "testdata", "unexported_return_mixed_pkg")

	files := []string{
		filepath.Join(baseDir, "pub.go"),
		filepath.Join(baseDir, "pub_internal_test.go"),
		filepath.Join(baseDir, "pub_external_test.go"),
	}

	expectedFailures := []string{
		"exported func New returns unexported type *pub.impl[T], which can be annoying to use",
		"exported func NewExportedWithUnexportedParam returns unexported type *pub.Exported[pub.hidden], which can be annoying to use",
	}

	const iterations = 50

	for i := 0; i < iterations; i++ {
		l := lint.New(os.ReadFile, 0)

		ps, err := l.Lint([][]string{files}, []lint.Rule{&rule.UnexportedReturnRule{}}, lint.Config{})
		if err != nil {
			t.Fatalf("iteration %d: Lint failed: %v", i, err)
		}

		var failures []lint.Failure
		for p := range ps {
			failures = append(failures, p)
		}

		var prodFailures []lint.Failure
		for _, f := range failures {
			if !strings.HasSuffix(f.Filename(), "_test.go") {
				prodFailures = append(prodFailures, f)
			}
		}

		if len(prodFailures) != len(expectedFailures) {
			t.Fatalf("iteration %d: expected %d failures on pub.go, got %d: %+v",
				i, len(expectedFailures), len(prodFailures), prodFailures)
		}

		for idx, exp := range expectedFailures {
			got := prodFailures[idx].Failure
			if got != exp {
				t.Fatalf("iteration %d, failure %d: wrong failure message\n got: %s\nwant: %s",
					i, idx, got, exp)
			}
		}
	}
}

// TestUnexportedReturnConcurrentExecution verifies that concurrent linting across
// multiple packages and goroutines does not exhibit race conditions or shared mutable state issues.
func TestUnexportedReturnConcurrentExecution(t *testing.T) {
	baseDir := filepath.Join("..", "testdata", "unexported_return_mixed_pkg")

	files := []string{
		filepath.Join(baseDir, "pub.go"),
		filepath.Join(baseDir, "pub_internal_test.go"),
		filepath.Join(baseDir, "pub_external_test.go"),
	}

	const goroutines = 10
	var wg sync.WaitGroup
	for g := range goroutines {
		wg.Go(func() {
			l := lint.New(os.ReadFile, 0)
			ps, err := l.Lint([][]string{files}, []lint.Rule{&rule.UnexportedReturnRule{}}, lint.Config{})
			if err != nil {
				t.Errorf("goroutine %d: Lint failed: %v", g, err)
				return
			}

			count := 0
			for f := range ps {
				if !strings.HasSuffix(f.Filename(), "_test.go") {
					count++
				}
			}

			if count != 2 {
				t.Errorf("goroutine %d: expected 2 failures, got %d", g, count)
			}
		})
	}

	wg.Wait()
}
