package worktree

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi/axitest"
)

// TestTreeBuildArgument: the tree-build refusal names the repair with the label as one
// always-quoted shell argument. Each wanted operand is authored from the label, so a former
// spelling or a bypassed owner reds the exact line. A POSIX shell then reads the next line
// back, so a label that loses a quote or splits reds the recovered argv, and a sentinel file
// reds a label that the shell executes.
func TestTreeBuildArgument(t *testing.T) {
	t.Parallel()
	sentinel := filepath.Join(t.TempDir(), "executed")
	for _, row := range []struct {
		name, label, operand string
	}{
		{name: "an embedded quote", label: "it's", operand: `'it'\''s'`},
		{name: "spaces", label: "my  alpha ", operand: "'my  alpha '"},
		{name: "glob and dollar", label: "*.go $HOME", operand: "'*.go $HOME'"},
		{name: "a command substitution", label: "$(touch " + sentinel + ")", operand: "'$(touch " + sentinel + ")'"},
		{name: "a command list", label: "a; touch " + sentinel, operand: "'a; touch " + sentinel + "'"},
	} {
		t.Run(row.name, func(t *testing.T) {
			var stderr bytes.Buffer
			code := PrintTreeBuildRefusal(&stderr, "bench status --in", row.label, "worktree build is missing")
			want := "bench status --in: worktree build is missing\nnext=bench worktree build " + row.operand + "\n"
			requireTest(t, code == 1 && stderr.String() == want, "refusal = (%d, %q), want (1, %q)", code, stderr.String(), want)
			_, next, _ := strings.Cut(strings.TrimSuffix(stderr.String(), "\n"), "\nnext=")
			argv, err := axitest.HelpAction{Cmd: next}.Argv()
			wantArgv := []string{"bench", "worktree", "build", row.label}
			requireTest(t, err == nil && slices.Equal(argv, wantArgv), "argv of %q = (%q, %v), want %q", next, argv, err, wantArgv)
		})
	}
	if _, err := os.Stat(sentinel); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a refusal label ran as a command: %v", err)
	}
}
