package lint_test

import (
	"testing"

	"github.com/mgechev/revive/lint"
)

func TestFileFilter(t *testing.T) {
	t.Run("whole file name", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("a/b/c.go")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("a/b/c.go") {
			t.Fatal("should match a/b/c.go")
		}
		if ff.MatchFileName("a/b/d.go") {
			t.Fatal("should not match")
		}
	})

	t.Run("whole file name with regexp metacharacters", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("pkg/a+(b)[c].go")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("pkg/a+(b)[c].go") {
			t.Fatal("should match pkg/a+(b)[c].go")
		}
		if ff.MatchFileName("pkg/aa(b)c.go") {
			t.Fatal("should not match pkg/aa(b)c.go")
		}
	})

	t.Run("glob with regexp metacharacters", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("pkg/a|b*.go")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("pkg/a|bc.go") {
			t.Fatal("should match pkg/a|bc.go")
		}
		if ff.MatchFileName("pkg/alpha/unrelated.go") {
			t.Fatal("should not match pkg/alpha/unrelated.go")
		}
		if ff.MatchFileName("zzz/xb.go") {
			t.Fatal("should not match zzz/xb.go")
		}
	})

	t.Run("glob with windows separators", func(t *testing.T) {
		ff, err := lint.ParseFileFilter(`pkg\x*.go`)
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("pkg/xy.go") {
			t.Fatal("should match pkg/xy.go")
		}
		if ff.MatchFileName("pkg/y.go") {
			t.Fatal("should not match pkg/y.go")
		}
	})

	t.Run("whole file name with invalid UTF-8", func(t *testing.T) {
		_, err := lint.ParseFileFilter("pkg/\xff.go")
		if err == nil {
			t.Fatal("should return an error, not panic")
		}
	})

	t.Run("regex", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("~b/[cd].go$")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("a/b/c.go") {
			t.Fatal("should match a/b/c.go")
		}
		if !ff.MatchFileName("b/d.go") {
			t.Fatal("should match b/d.go")
		}
		if ff.MatchFileName("b/x.go") {
			t.Fatal("should not match b/x.go")
		}
	})

	t.Run("TEST well-known", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("TEST")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("a/b/c_test.go") {
			t.Fatal("should match a/b/c_test.go")
		}
		if ff.MatchFileName("a/b/c_test_no.go") {
			t.Fatal("should not match a/b/c_test_no.go")
		}
	})

	t.Run("glob *", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("a/b/*.pb.go")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("a/b/xxx.pb.go") {
			t.Fatal("should match a/b/xxx.pb.go")
		}
		if !ff.MatchFileName("a/b/yyy.pb.go") {
			t.Fatal("should match a/b/yyy.pb.go")
		}
		if ff.MatchFileName("a/b/xxx.nopb.go") {
			t.Fatal("should not match a/b/xxx.nopb.go")
		}
	})

	t.Run("glob **", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("a/**/*.pb.go")
		if err != nil {
			t.Fatal(err)
		}
		if !ff.MatchFileName("a/x/xxx.pb.go") {
			t.Fatal("should match a/x/xxx.pb.go")
		}
		if !ff.MatchFileName("a/xxx.pb.go") {
			t.Fatal("should match a/xxx.pb.go")
		}
		if !ff.MatchFileName("a/x/y/z/yyy.pb.go") {
			t.Fatal("should match a/x/y/z/yyy.pb.go")
		}
		if ff.MatchFileName("a/b/xxx.nopb.go") {
			t.Fatal("should not match a/b/xxx.nopb.go")
		}
	})

	t.Run("empty", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("")
		if err != nil {
			t.Fatal(err)
		}
		fileNames := []string{"pb.go", "a/pb.go", "a/x/xxx.pb.go", "a/x/xxx.pb_test.go"}
		for _, fn := range fileNames {
			if ff.MatchFileName(fn) {
				t.Fatalf("should not match %s", fn)
			}
		}
	})

	t.Run("just *", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("*")
		if err != nil {
			t.Fatal(err)
		}
		fileNames := []string{"pb.go", "a/pb.go", "a/x/xxx.pb.go", "a/x/xxx.pb_test.go"}
		for _, fn := range fileNames {
			if !ff.MatchFileName(fn) {
				t.Fatalf("should match %s", fn)
			}
		}
	})

	t.Run("just ~", func(t *testing.T) {
		ff, err := lint.ParseFileFilter("~")
		if err != nil {
			t.Fatal(err)
		}
		fileNames := []string{"pb.go", "a/pb.go", "a/x/xxx.pb.go", "a/x/xxx.pb_test.go"}
		for _, fn := range fileNames {
			if !ff.MatchFileName(fn) {
				t.Fatalf("should match %s", fn)
			}
		}
	})
}
