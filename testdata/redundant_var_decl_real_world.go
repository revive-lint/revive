// Cases taken from real-world fixes referenced in https://github.com/revive-lint/revive/issues/1098:
// go-delve/delve#3855, go-gitea/gitea#32381 and vektra/mockery#829.
// Project-specific packages are replaced by local stand-ins, so the type checker knows every type.

package fixtures

import (
	"context"
	"go/format"
	"net/http"
	"os"
	"strings"
)

// Stand-ins for project types.

type asmInstructions []string

type debugger struct{}

func (*debugger) Disassemble(goroutineID int64, startPC, endPC uint64) (asmInstructions, error) {
	return nil, nil
}

type rpcServer struct {
	debugger *debugger
}

type disassembleRequest struct {
	goroutineID    int64
	startPC, endPC uint64
}

type stackframe struct{}

func threadStacktrace(depth int) ([]stackframe, error) { return nil, nil }

type evalScope struct{}

func (evalScope) LocalVariables(cfg int) ([]string, error) { return nil, nil }

func frameToScope() evalScope { return evalScope{} }

type committer interface{ Close() error }

func txContext(ctx context.Context) (context.Context, committer, error) { return ctx, nil, nil }

type searchOptions struct{}

func searchRepository(opts searchOptions) ([]string, int64, error) { return nil, 0, nil }

type treeEntry struct{}

func (treeEntry) ListEntries() ([]treeEntry, error) { return nil, nil }

func toUTF8(content []byte) (string, error) { return string(content), nil }

// Reported cases.

// delve: service/rpc2/server.go (Disassemble), the example from the issue.
func (s *rpcServer) disassemble(arg disassembleRequest) (asmInstructions, error) {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 63; remove the var declaration/
	insts, err := s.debugger.Disassemble(arg.goroutineID, arg.startPC, arg.endPC)
	if err != nil {
		return nil, err
	}
	return insts, nil
}

// delve: pkg/proc/breakpoints.go, declaration inside an if block.
func nextDeferOk(kind int) bool {
	ok := true
	if kind != 0 {
		var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 75; remove the var declaration/
		frames, err := threadStacktrace(2)
		if err == nil {
			ok = len(frames) > 0
		}
	}
	return ok
}

// delve: service/debugger/debugger.go, an unrelated statement between the declaration and :=.
func convertStacktrace(cfg *int) ([]string, error) {
	if cfg != nil {
		var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 88; remove the var declaration/
		scope := frameToScope()
		locals, err := scope.LocalVariables(*cfg)
		if err != nil {
			return nil, err
		}
		return locals, nil
	}
	return nil, nil
}

// delve: cmd/dlv/dlv_test.go, the declaration is separated by several statements, including a defer.
func dlvBinary() error {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 105; remove the var declaration/

	tmpdir := os.TempDir()
	defer os.Remove(tmpdir)

	listenPath := strings.Join([]string{tmpdir, "delve_test"}, "/")
	f, err := os.Create(listenPath)
	if err != nil {
		return err
	}
	return f.Close()
}

// gitea: build/generate-emoji.go, a call into the standard library.
func generate(gemojiURL string) (*http.Response, error) {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 116; remove the var declaration/

	res, err := http.Get(gemojiURL)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// gitea: models/git/lfs.go, the := also redeclares a parameter.
func newLFSMetaObject(ctx context.Context) error {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 127; remove the var declaration/

	ctx, committer, err := txContext(ctx)
	if err != nil {
		return err
	}
	_ = ctx
	return committer.Close()
}

// gitea: routers/api/v1/repo/repo.go, a call with three results.
func search(opts searchOptions) (int64, error) {
	var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 138; remove the var declaration/
	repos, count, err := searchRepository(opts)
	if err != nil {
		return 0, err
	}
	_ = repos
	return count, nil
}

// gitea: routers/web/repo/view.go, declaration inside a loop body.
func findReadme(entries []treeEntry) ([]treeEntry, error) {
	for _, subTree := range entries {
		var err error // MATCH /redundant declaration of "err": it is redeclared by := on line 150; remove the var declaration/
		childEntries, err := subTree.ListEntries()
		if err != nil {
			return nil, err
		}
		if len(childEntries) > 0 {
			return childEntries, nil
		}
	}
	return nil, nil
}

// mockery: pkg/generator_test.go, a non-error variable.
func checkGeneration(buf []byte) ([]byte, error) {
	var actual []byte // MATCH /redundant declaration of "actual": it is redeclared by := on line 164; remove the var declaration/
	actual, fmtErr := format.Source(buf)
	if fmtErr != nil {
		return nil, fmtErr
	}
	return actual, nil
}

// Not reported cases: the fixes in these PRs rewrote later assignments, the declaration alone is not redundant.

// gitea: models/issues/label_test.go, assignment with =.
func getLabels() error {
	var err error
	_, _, err = searchRepository(searchOptions{})
	return err
}

// gitea: modules/charset/charset_test.go, both variables are assigned with =.
func convert() (string, error) {
	var res string
	var err error

	res, err = toUTF8([]byte{0x41, 0x42, 0x43})
	return res, err
}

// mockery: pkg/generator_test.go, compound assignment.
func expectedStruct() string {
	var expected string
	expected += "*struct {"
	expected += "}"
	return expected
}
