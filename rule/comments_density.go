package rule

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/mgechev/revive/lint"
)

// CommentsDensityRule enforces a minimum comment / code relation.
type CommentsDensityRule struct {
	minimumCommentsDensity int64
}

const defaultMinimumCommentsPercentage = 0

var _ lint.ConfigurableRule = (*CommentsDensityRule)(nil)

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *CommentsDensityRule) Configure(arguments lint.Arguments) error {
	r.minimumCommentsDensity = defaultMinimumCommentsPercentage
	if len(arguments) < 1 {
		return nil
	}

	var ok bool
	r.minimumCommentsDensity, ok = arguments[0].(int64)
	if !ok {
		return fmt.Errorf("invalid argument for %q rule: argument should be an int, got %T", r.Name(), arguments[0])
	}
	return nil
}

// Apply applies the rule to given file.
func (r *CommentsDensityRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	commentLines, codeLines := countLines(file)
	total := codeLines + commentLines
	if total == 0 {
		return nil
	}

	density := (float32(commentLines) / float32(total)) * 100

	if density < float32(r.minimumCommentsDensity) {
		return []lint.Failure{
			{
				Category:   lint.FailureCategoryMaintenance,
				Node:       file.AST,
				Confidence: 1,
				Failure: fmt.Sprintf("the file has a comment density of %2.f%% (%d comment lines for %d code lines) but expected a minimum of %d%%",
					density, commentLines, codeLines, r.minimumCommentsDensity),
			},
		}
	}

	return nil
}

// Name returns the rule name.
func (*CommentsDensityRule) Name() string {
	return "comments-density"
}

func countLines(file *lint.File) (commentLines, codeLines int) {
	content := file.Content()
	if len(content) == 0 {
		return 0, 0
	}

	lineStarts := []int{0}
	for i, b := range content {
		if b == '\n' {
			lineStarts = append(lineStarts, i+1)
		}
	}

	masked := make([]byte, len(content))
	copy(masked, content)

	commentLinesSet := map[int]struct{}{}

	if file.AST != nil {
		for _, cg := range file.AST.Comments {
			for _, c := range cg.List {
				offset := file.ToPosition(c.Slash).Offset
				if offset < 0 || offset >= len(content) {
					if file.AST.FileStart > 0 {
						offset = int(c.Slash - file.AST.FileStart)
					}
				}
				if offset < 0 || offset >= len(content) {
					continue
				}

				end := offset
				if bytes.HasPrefix(content[offset:], []byte("/*")) {
					if idx := bytes.Index(content[offset+2:], []byte("*/")); idx >= 0 {
						end = offset + 2 + idx + 2
					} else {
						end = len(content)
					}
				} else if bytes.HasPrefix(content[offset:], []byte("//")) {
					if idx := bytes.IndexByte(content[offset+2:], '\n'); idx >= 0 {
						end = offset + 2 + idx
					} else {
						end = len(content)
					}
				} else {
					end = offset + len(c.Text)
					if end > len(content) {
						end = len(content)
					}
				}

				for i := offset; i < end && i < len(masked); i++ {
					if masked[i] != '\n' {
						masked[i] = ' '
					}
				}

				startLine := offsetToLine(lineStarts, offset)
				endLine := startLine
				if end > offset {
					endLine = offsetToLine(lineStarts, end-1)
				}
				for l := startLine; l <= endLine; l++ {
					commentLinesSet[l] = struct{}{}
				}
			}
		}
	}

	for line := range bytes.Lines(masked) {
		if len(bytes.TrimSpace(line)) > 0 {
			codeLines++
		}
	}

	return len(commentLinesSet), codeLines
}

func offsetToLine(lineStarts []int, offset int) int {
	i, exact := slices.BinarySearch(lineStarts, offset)
	if exact {
		return i + 1
	}
	return i
}
