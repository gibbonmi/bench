package worktree

import (
	"bytes"
	"testing"
)

// TestTargetRefusalEscapesOnlyAnUnsafeLine: the shared refusal printer prints a line with no
// control rune as is, so a pasted repair names the label the operator typed. It escapes a
// line with a control rune, and only that line. Each row prints through
// PrintTreeBuildRefusal, so each row also pins the label as one shell word in the next line.
// Each wanted stderr is exact, so it holds no raw control byte.
func TestTargetRefusalEscapesOnlyAnUnsafeLine(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, detail, label, want string
	}{
		{name: "a plain label", detail: "worktree build is missing", label: "alpha",
			want: "bench status --in: worktree build is missing\nnext=bench worktree build alpha\n"},
		{name: "a label with a backslash", detail: "worktree build is missing", label: `a\b`,
			want: "bench status --in: worktree build is missing\nnext=bench worktree build 'a\\b'\n"},
		{name: "a control byte in the detail", detail: "worktree build\x1b[2J is missing", label: `a\b`,
			want: "bench status --in: worktree build\\u001b[2J is missing\nnext=bench worktree build 'a\\b'\n"},
		{name: "a control byte in the label", detail: "worktree build is missing", label: "a\x07b",
			want: "bench status --in: worktree build is missing\nnext=bench worktree build 'a\\u0007b'\n"},
		// The lookup refuses a label with a control byte before the build check, so no call of
		// Run reaches this refusal with one. This row is the only pin of that case.
		{name: "a control byte on either line", detail: "worktree build\x1b[2J is missing", label: "a\x07b",
			want: "bench status --in: worktree build\\u001b[2J is missing\nnext=bench worktree build 'a\\u0007b'\n"},
	} {
		t.Run(row.name, func(t *testing.T) {
			var stderr bytes.Buffer
			code := PrintTreeBuildRefusal(&stderr, "bench status --in", row.label, row.detail)
			requireTest(t, code == 1 && stderr.String() == row.want, "refusal = (%d, %q), want (1, %q)", code, stderr.String(), row.want)
		})
	}
}
