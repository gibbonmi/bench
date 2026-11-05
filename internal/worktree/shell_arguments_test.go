package worktree

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/intent"
)

// TestTreeBuildArgument: the tree-build refusal names the repair with the label as one
// always-quoted shell argument. Each wanted operand is authored from the label, so a former
// spelling or a bypassed owner reds the exact line. A POSIX shell then reads the next line
// back, so a label that loses a quote or splits reds the recovered argv, and a sentinel file
// reds a label that the shell executes.
func TestTreeBuildArgument(t *testing.T) {
	t.Parallel()
	sentinel := axitest.NewShellSentinel(t)
	for _, row := range []struct {
		name, label, operand string
	}{
		{name: "an embedded quote", label: "it's", operand: `'it'\''s'`},
		{name: "spaces", label: "my  alpha ", operand: "'my  alpha '"},
		{name: "glob and dollar", label: "*.go $HOME", operand: "'*.go $HOME'"},
		{name: "a command substitution", label: "$(" + sentinel.Touch() + ")", operand: "'$(" + sentinel.Touch() + ")'"},
		{name: "a command list", label: "a; " + sentinel.Touch(), operand: "'a; " + sentinel.Touch() + "'"},
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
	sentinel.RequireUnexecuted(t)
}

// TestRecoveryRouteArgument: each real recovery line with a hostile path equals both the
// command authored from the path and the command of its help row, and a POSIX shell reads
// the exact path back from the line. A path that no help row discloses keeps the line's
// quoted path, and the help block omits that row. A sentinel file reds a path that the
// shell executes.
func TestRecoveryRouteArgument(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	sentinel := axitest.NewShellSentinel(t)
	for _, row := range []struct {
		name, path, operand string
		disclosed           bool
	}{
		{name: "an embedded quote", path: "/gone/it's", operand: `'/gone/it'"'"'s'`, disclosed: true},
		{name: "spaces", path: "/gone/my  tree ", operand: "'/gone/my  tree '", disclosed: true},
		{name: "a command substitution", path: "/gone/$(" + sentinel.Touch() + ")", operand: "'/gone/$(" + sentinel.Touch() + ")'", disclosed: true},
		{name: "a command list", path: "/gone/a; " + sentinel.Touch(), operand: "'/gone/a; " + sentinel.Touch() + "'", disclosed: true},
		{name: "an angle bracket", path: "/gone/a<b>", operand: "'/gone/a<b>'"},
	} {
		gone := intent.Assignment{ID: "gone", RequestToken: "req-gone", Branch: "refs/heads/bench/gone", Worktree: row.path}
		for _, route := range []struct {
			name, line string
			words      []string
			help       axi.Action
		}{
			{name: "discard", line: releaseNext(ReasonUnmerged, row.path, gone.ID), words: []string{"bench", "worktree", "clean", "--discard-branch"}, help: discardUnlanded(row.path, gone.ID).action()},
			{name: "missing tree", line: missingTreeRoute(root, gone), words: []string{"bench", "worktree", "release", "--request", "req-gone"}, help: recoverMissingTree(false, gone.RequestToken, row.path).action()},
		} {
			t.Run(row.name+"/"+route.name, func(t *testing.T) {
				rendered, err := axi.RenderHelp([]axi.Action{route.help})
				mustNoError(t, err)
				document, err := axitest.DecodeDocument(rendered)
				mustNoError(t, err)
				rows, err := document.HelpActions()
				mustNoError(t, err)
				var help []string
				for _, action := range rows {
					help = append(help, action.Cmd)
				}
				want := strings.Join(route.words, " ") + " " + row.operand
				wantHelp := []string(nil)
				if row.disclosed {
					wantHelp = []string{want}
				}
				requireTest(t, route.line == want && slices.Equal(help, wantHelp), "recovery line = %q and help commands = %q, want %q and %q", route.line, help, want, wantHelp)
				argv, err := axitest.HelpAction{Cmd: route.line}.Argv()
				wantArgv := append(slices.Clone(route.words), row.path)
				requireTest(t, err == nil && slices.Equal(argv, wantArgv), "argv of %q = (%q, %v), want %q", route.line, argv, err, wantArgv)
			})
		}
	}
	sentinel.RequireUnexecuted(t)
}
