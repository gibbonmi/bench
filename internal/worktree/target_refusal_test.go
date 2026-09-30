package worktree

import (
	"bytes"
	"testing"
)

// TestTargetRefusalEscapesOnlyAnUnsafeLine: the shared refusal printer prints a line with no
// control rune as is, so a pasted repair names the label the operator typed. It escapes a
// line with a control rune, and only that line. Each wanted stderr is exact, so it holds no
// raw control byte.
func TestTargetRefusalEscapesOnlyAnUnsafeLine(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name, detail, next, want string
	}{
		{name: "a plain label", detail: "worktree build is missing", next: "bench worktree build 'alpha'",
			want: "bench status --in: worktree build is missing\nnext=bench worktree build 'alpha'\n"},
		{name: "a label with a backslash", detail: "worktree build is missing", next: `bench worktree build 'a\b'`,
			want: "bench status --in: worktree build is missing\nnext=bench worktree build 'a\\b'\n"},
		{name: "a control byte in the detail", detail: "worktree build\x1b[2J is missing", next: `bench worktree build 'a\b'`,
			want: "bench status --in: worktree build\\u001b[2J is missing\nnext=bench worktree build 'a\\b'\n"},
		{name: "a control byte in the label", detail: "worktree build is missing", next: "bench worktree build 'a\x07b'",
			want: "bench status --in: worktree build is missing\nnext=bench worktree build 'a\\u0007b'\n"},
	} {
		t.Run(row.name, func(t *testing.T) {
			var stderr bytes.Buffer
			code := printTargetRefusal(&stderr, "bench status --in", refusalError{refusal{detail: row.detail, next: row.next}})
			requireTest(t, code == 1 && stderr.String() == row.want, "refusal = (%d, %q), want (1, %q)", code, stderr.String(), row.want)
		})
	}
}
