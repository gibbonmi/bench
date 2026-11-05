package worktree

import (
	"bytes"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/worktree/reclaimpolicy"

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
// quoted words and path, and the help block omits that row. A discard path that no line
// can carry leaves the route to the assignment id. A sentinel file reds a path that the
// shell executes.
func TestRecoveryRouteArgument(t *testing.T) {
	t.Parallel()
	root := newWorktreeRepo(t)
	sentinel := axitest.NewShellSentinel(t)
	for _, row := range []struct {
		name, path, operand string
		disclosed           bool
	}{
		{name: "an embedded quote", path: "/gone/it's", operand: `'/gone/it'\''s'`, disclosed: true},
		{name: "spaces", path: "/gone/my  tree ", operand: "'/gone/my  tree '", disclosed: true},
		{name: "a command substitution", path: "/gone/$(" + sentinel.Touch() + ")", operand: "'/gone/$(" + sentinel.Touch() + ")'", disclosed: true},
		{name: "a command list", path: "/gone/a; " + sentinel.Touch(), operand: "'/gone/a; " + sentinel.Touch() + "'", disclosed: true},
		{name: "an angle bracket", path: "/gone/a<b>", operand: "'/gone/a<b>'"},
	} {
		gone := intent.Assignment{ID: "gone", RequestToken: "req-gone", Branch: "refs/heads/bench/gone", Worktree: row.path}
		for _, route := range []struct {
			name, line, quoted string
			words              []string
			help               axi.Action
		}{
			{name: "discard", line: releaseNext(ReasonUnmerged, row.path, gone.ID), quoted: "bench 'worktree' 'clean' '--discard-branch'", words: []string{"bench", "worktree", "clean", "--discard-branch"}, help: discardUnlanded(row.path, gone.ID).action()},
			{name: "missing tree", line: missingTreeRoute(root, gone), quoted: "bench 'worktree' 'release' '--request' 'req-gone'", words: []string{"bench", "worktree", "release", "--request", "req-gone"}, help: recoverMissingTree(false, gone.RequestToken, row.path).action()},
		} {
			t.Run(row.name+"/"+route.name, func(t *testing.T) {
				rendered, err := axi.RenderHelp([]axi.Action{route.help})
				mustNoError(t, err)
				// A disclosed route prints its help-row command; a route that no help row
				// discloses prints the same quoted words and path, and no help row.
				want := route.quoted + " " + row.operand
				wantArgv := append(slices.Clone(route.words), row.path)
				wantHelp, wantHelpArgv := []string(nil), [][]string(nil)
				if row.disclosed {
					wantHelp, wantHelpArgv = []string{want}, [][]string{wantArgv}
				}
				axitest.RequireHelpCommands(t, rendered, wantHelp, wantHelpArgv)
				requireTest(t, route.line == want, "recovery line = %q, want %q", route.line, want)
				argv, err := axitest.HelpAction{Cmd: route.line}.Argv()
				requireTest(t, err == nil && slices.Equal(argv, wantArgv), "argv of %q = (%q, %v), want %q", route.line, argv, err, wantArgv)
			})
		}
	}
	t.Run("a newline/discard", func(t *testing.T) {
		path := "/gone/a\nb"
		line := releaseNext(ReasonUnmerged, path, "gone")
		const want = "bench 'worktree' 'clean' '--discard-branch' '--target' 'gone'"
		wantArgv := []string{"bench", "worktree", "clean", "--discard-branch", "--target", "gone"}
		rendered, err := axi.RenderHelp([]axi.Action{discardUnlanded(path, "gone").action()})
		mustNoError(t, err)
		axitest.RequireHelpCommands(t, rendered, []string{want}, [][]string{wantArgv})
		requireTest(t, line == want, "recovery line = %q, want %q", line, want)
	})
	sentinel.RequireUnexecuted(t)
}

// TestListHelpArguments is S27: the list actions of a cleanup-pending row, an active row, and
// an orphan quote every known argument, a hostile path included, and keep each future slot.
// A sentinel file reds a path that the shell runs.
func TestListHelpArguments(t *testing.T) {
	t.Parallel()
	sentinel := axitest.NewShellSentinel(t)
	pending, orphan := "/gone/it's; "+sentinel.Touch(), "/orphan/$("+sentinel.Touch()+")"
	rows := []listRow{
		{values: []any{"p", "pending", "req-p", string(intent.StateCleanupPending), "assignment", "present", "none", "unknown"}, assignmentPath: pending},
		activeListRow("a", "req-a", "present", true, "/tmp/a"),
		{values: []any{"foreign", "orphan", "", "foreign", "foreign", "missing"}, orphanPath: orphan},
	}
	help, err := axi.RenderHelp(actionsForRows(rows, nil))
	mustNoError(t, err)
	axitest.RequireHelpCommands(t, help, []string{
		`bench 'worktree' 'release' '--request' 'req-p' '/gone/it'\''s; ` + sentinel.Touch() + `'`,
		"bench 'worktree' 'path' <target>",
		"bench 'worktree' 'exec' <target> '--' <command>",
		"bench 'worktree' 'clean' '/orphan/$(" + sentinel.Touch() + ")'",
	}, [][]string{
		{"bench", "worktree", "release", "--request", "req-p", pending},
		{"bench", "worktree", "path"},
		{"bench", "worktree", "exec"},
		{"bench", "worktree", "clean", orphan},
	})
	sentinel.RequireUnexecuted(t)
}

// TestSelectedListHelpArguments is S28: a selected list names the complete inventory with its
// two words quoted.
func TestSelectedListHelpArguments(t *testing.T) {
	t.Parallel()
	root, assignments := selectedFixture(t, intent.StateActive)
	r := runVerb(t, verbList, verbCall{root: root, args: []string{"--view", "paths", "--target", assignments[0].ID}})
	axitest.RequireHelpCommands(t, r.stdout, []string{"bench 'worktree' 'list'"}, [][]string{{"bench", "worktree", "list"}})
}

// TestExplicitCleanHelpArguments is S29: the apply action of an explicit set quotes the verb,
// each modifier, each selector, the apply flag, and the fingerprint.
func TestExplicitCleanHelpArguments(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	id := f.first.Assignment.ID
	plan := runVerb(t, verbClean, f.call("--full", "--target", id))
	requireTest(t, plan.exit == 0 && plan.stderr == "", "plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	axitest.RequireHelpCommands(t, plan.stdout,
		[]string{"bench 'worktree' 'clean' '--full' '--target' '" + id + "' '--apply' '" + fingerprint + "'"},
		[][]string{{"bench", "worktree", "clean", "--full", "--target", id, "--apply", fingerprint}})
}

// TestCleanRefusalHelpArguments is S30: a stale set refusal re-plans the same selection with
// every modifier and selector quoted and no digest. A sentinel file reds a selector that the
// shell runs.
func TestCleanRefusalHelpArguments(t *testing.T) {
	t.Parallel()
	sentinel := axitest.NewShellSentinel(t)
	selector := "it's; " + sentinel.Touch()
	var out bytes.Buffer
	mustNoError(t, renderStale(&out, nil, cleanArguments(CleanupOptions{DiscardIgnored: true}, "--target", selector)))
	_, help, _ := strings.Cut(out.String(), "help[")
	axitest.RequireHelpCommands(t, "help["+help,
		[]string{`bench 'worktree' 'clean' '--discard-ignored' '--target' 'it'\''s; ` + sentinel.Touch() + `'`},
		[][]string{{"bench", "worktree", "clean", "--discard-ignored", "--target", selector}})
	sentinel.RequireUnexecuted(t)
}

// TestLandedCleanHelpArguments is S31: the landed plan applies with its selector, apply flag,
// and fingerprint quoted, and the retained dirty member names its path remedy quoted.
func TestLandedCleanHelpArguments(t *testing.T) {
	t.Parallel()
	f := landedSetFixture(t)
	plan := runVerb(t, verbClean, f.call("--landed"))
	requireTest(t, plan.exit == 0, "plan = (%d, %q, %q), want exit 0", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	axitest.RequireHelpCommands(t, plan.stdout, []string{
		"bench 'worktree' 'clean' '--landed' '--apply' '" + fingerprint + "'",
		"bench 'worktree' 'clean' '" + f.dirty.Path + "'",
	}, [][]string{
		{"bench", "worktree", "clean", "--landed", "--apply", fingerprint},
		{"bench", "worktree", "clean", f.dirty.Path},
	})
}

// TestUnclaimedCleanHelpArguments is S32: the unclaimed plan applies with both selectors, the
// apply flag, and the fingerprint quoted.
func TestUnclaimedCleanHelpArguments(t *testing.T) {
	t.Parallel()
	root, home := unclaimedBranchFixture(t, "a", "b")
	plan := runVerb(t, verbClean, repoHome{root, home}.call("--discard-branch", "--unclaimed"))
	requireTest(t, plan.exit == 0 && plan.stderr == "", "plan = (%d, %q, %q), want one applicable plan", plan.exit, plan.stdout, plan.stderr)
	fingerprint := plan.mustFingerprint(t)
	axitest.RequireHelpCommands(t, plan.stdout,
		[]string{"bench 'worktree' 'clean' '--discard-branch' '--unclaimed' '--apply' '" + fingerprint + "'"},
		[][]string{{"bench", "worktree", "clean", "--discard-branch", "--unclaimed", "--apply", fingerprint}})
}

// TestReclaimHelpArguments is S33: the pool plan applies with its verb, flag, and fingerprint
// quoted, and the stale refusal re-plans with its two words quoted.
func TestReclaimHelpArguments(t *testing.T) {
	t.Parallel()
	fingerprint := strings.Repeat("ab", 32)
	out, err := renderPoolReclaim(poolReclaimPlan{verdicts: []poolKeyVerdict{{Key: "dead", Verdict: reclaimpolicy.VerdictReclaim, Reason: "gone"}}, fingerprint: fingerprint})
	mustNoError(t, err)
	axitest.RequireHelpCommands(t, out,
		[]string{"bench 'worktree' 'reclaim' '--apply' '" + fingerprint + "'"},
		[][]string{{"bench", "worktree", "reclaim", "--apply", fingerprint}})
	stale, err := renderPoolReclaimStale()
	mustNoError(t, err)
	_, help, _ := strings.Cut(stale, "help[")
	axitest.RequireHelpCommands(t, "help["+help, []string{"bench 'worktree' 'reclaim'"}, [][]string{{"bench", "worktree", "reclaim"}})
}

// TestReadActionArguments is S43: the read action quotes its words and the label, and keeps
// its path slot. A sentinel file reds a label that the shell runs.
func TestReadActionArguments(t *testing.T) {
	t.Parallel()
	sentinel := axitest.NewShellSentinel(t)
	label := "it's; " + sentinel.Touch()
	help, err := axi.RenderHelp([]axi.Action{readAction(label)})
	mustNoError(t, err)
	axitest.RequireHelpCommands(t, help,
		[]string{`bench 'worktree' 'read' '--in' 'it'\''s; ` + sentinel.Touch() + `' <path>`},
		[][]string{{"bench", "worktree", "read", "--in", label}})
	sentinel.RequireUnexecuted(t)
}
