package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/axi"
	"github.com/gibbonmi/bench/internal/axi/axitest"
	"github.com/gibbonmi/bench/internal/gittest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/responsebound"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

// TestWorktreeLeafHelpAnswersItsOwnGrammar grades SR1. Every row that names a help
// grammar must answer that grammar, so a row paired with another leaf's handler shows
// up as a foreign grammar rather than as a silent reroute. The two sinks are joined
// because a leaf answers its grammar on stdout or on stderr and the row does not grade
// which.
func TestWorktreeLeafHelpAnswersItsOwnGrammar(t *testing.T) {
	for _, leaf := range worktreeLeaves {
		if leaf.Grammar == "" {
			continue
		}
		t.Run(leaf.Name, func(t *testing.T) {
			root := newAXIEnvelopeRepo(t)
			t.Setenv("BENCH_HOME", t.TempDir())
			result := runAXICommandAt(t, root, []string{"worktree", leaf.Name, "--help"})
			out := result.stdout + result.stderr
			want := "usage: " + leaf.Grammar
			if result.code != 0 || !strings.Contains(out, want) {
				t.Fatalf("worktree %s --help = output=%q exit=%d, want exit 0 and %q", leaf.Name, out, result.code, want)
			}
		})
	}
}

// TestWorktreeRootRequiredLeavesRefuseOutsideARepository grades SR2. A leaf whose root
// need was lost would run its handler against an empty root and answer that handler's
// own refusal, so the row demands the one not-in-repo line on stderr at exit 1.
func TestWorktreeRootRequiredLeavesRefuseOutsideARepository(t *testing.T) {
	names := make([]string, 0, len(worktreeLeaves))
	for _, leaf := range worktreeLeaves {
		if leaf.Root == rootRequired {
			names = append(names, leaf.Name)
		}
	}
	want := []string{"exec", "path", "show", "build", "create", "release", "reauthorize", "merge", "reset", "land"}
	if !equalStringSets(names, want) {
		t.Fatalf("root-required leaves = %q, want %q", names, want)
	}
	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			t.Setenv("BENCH_HOME", t.TempDir())
			result := runAXICommandAt(t, t.TempDir(), []string{"worktree", name})
			wantErr := toon.NotInRepo() + "\n"
			if result.stdout != "" || result.stderr != wantErr || result.code != 1 {
				t.Fatalf("worktree %s outside a repository = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 1", name, result.stdout, result.stderr, result.code, wantErr)
			}
		})
	}
}

// ER29: the worktree reader is a tree-scoped query leaf. A target that names no active
// assignment refuses through the tree resolver before any child starts. In the caller's
// tree the identity row leads, the framed page fits the response bound, and the
// continuation names the tree and the root-relative path. No call here follows the
// continuation, because its tree target starts a child.
func TestWorktreeReadIsATreeScopedQuery(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "page.txt"), strings.Repeat("page text\n", 600))
	writeAXIFixture(t, filepath.Join(root, "sub", "other.txt"), "")
	t.Chdir(filepath.Join(root, "sub"))
	var stdout, stderr bytes.Buffer
	if code := (Command{Stdout: &stdout, Stderr: &stderr}).Run([]string{"worktree", "read", "--in", "beta", "page.txt"}); code != 1 || stdout.String() != "" || stderr.String() != "bench worktree read --in: target is unassigned\nnext=bench worktree list\n" {
		t.Fatalf("worktree read --in beta = %q %q exit %d, want the unassigned target refusal", stdout.String(), stderr.String(), code)
	}
	stdout.Reset()
	code := (Command{Stdout: &stdout}).Run([]string{"worktree", "read", "../page.txt"})
	if out := stdout.String(); code != 0 || !strings.HasPrefix(out, "tree[1]{target,head,dirty}:\n  primary,") || !responsebound.Fits(out) {
		t.Fatalf("worktree read = %q exit %d, want the identity row first and the framed page within the bound", out, code)
	}
	next, err := axitest.RecoverHelpCommandArgv(withoutTreeRow(stdout.String()))
	if want := []string{"bench", "worktree", "read", "--in", "primary", "page.txt", "--offset"}; err != nil || len(next) < len(want) || !slices.Equal(next[:len(want)], want) {
		t.Fatalf("worktree read continuation = %q (%v), want it to begin %q", next, err, want)
	}
	if !slices.Contains(axiRegistryMemberNames(t), "worktree read") {
		t.Fatal("worktree read is no AXI query member")
	}
}

// ER29: a worktree read whose identity-row target does not resolve back to the read tree
// refuses before any content and names root read, so no page offers a continuation that
// reads another tree or refuses. The primary keyword in the primary checkout still pages
// and continues, which TestWorktreeReadIsATreeScopedQuery grades.
func TestWorktreeReadRefusesATargetThatNamesAnotherTree(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree func(t *testing.T, root string, active intent.Assignment) string
	}{
		{"an unassigned linked tree", func(t *testing.T, root string, _ intent.Assignment) string {
			linked := filepath.Join(t.TempDir(), "unassigned")
			gittest.MustRun(t, root, "worktree", "add", "-q", "-b", "unassigned", linked)
			return linked
		}},
		{"two active assignments with one label", func(t *testing.T, root string, active intent.Assignment) string {
			twin := putRelatedAssignment(t, root, relatedID(101), intent.StateActive, filepath.Join(root, "absent", "twin"))
			twin.Label = active.Label
			if err := intent.PutAssignment(root, twin); err != nil {
				t.Fatal(err)
			}
			return active.Worktree
		}},
		{"an assignment labelled primary", func(t *testing.T, root string, active intent.Assignment) string {
			active.Label = "primary"
			if err := intent.PutAssignment(root, active); err != nil {
				t.Fatal(err)
			}
			return active.Worktree
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, active := relatedFixture(t, 0)
			tree := tc.tree(t, root, active)
			writeAXIFixture(t, filepath.Join(tree, "page.txt"), strings.Repeat("page text\n", 600))
			t.Chdir(tree)
			var stdout bytes.Buffer
			code := (Command{Stdout: &stdout}).Run([]string{"worktree", "read", "page.txt"})
			if body := withoutTreeRow(stdout.String()); code != 1 || !strings.HasPrefix(stdout.String(), "tree[1]{target,head,dirty}:\n") || !strings.HasPrefix(body, "error: ") || !strings.Contains(body, "bench read <path>") || strings.Contains(body, "page text") {
				t.Fatalf("worktree read in %s = %q exit %d, want the identity row and a refusal that names bench read <path> with no content", tc.name, stdout.String(), code)
			}
		})
	}
}

// TestWorktreeHelpWithAnExtraArgumentRefuses grades SR6. The help form matches on
// exactly one argument, so a table that matched `help` by name alone would answer the
// family usage at exit 0 for a two-word call.
func TestWorktreeHelpWithAnExtraArgumentRefuses(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := Command{Stdout: &stdout, Stderr: &stderr, Executable: "bench"}.Run([]string{"worktree", "help", "extra"})
	want := toon.Usage("bench worktree", "help") + "\n"
	if stdout.String() != "" || stderr.String() != want || code != 2 {
		t.Fatalf("worktree help extra = stdout=%q stderr=%q exit=%d, want stderr=%q and exit 2", stdout.String(), stderr.String(), code, want)
	}
}

func equalStringSets(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	seen := make(map[string]int, len(got))
	for _, name := range got {
		seen[name]++
	}
	for _, name := range want {
		seen[name]--
	}
	for _, count := range seen {
		if count != 0 {
			return false
		}
	}
	return true
}

func TestResetHelpIsAvailableAtEveryLevel(t *testing.T) {
	for _, argv := range [][]string{{"worktree", "reset", "--help"}, {"worktree", "--help"}, {"help"}} {
		result := runAXICommandAt(t, newAXIEnvelopeRepo(t), argv)
		if result.code != 0 || !strings.Contains(result.stdout+result.stderr, usage.WorktreeReset) {
			t.Fatalf("%v = %d %s %s, want %s", argv, result.code, result.stdout, result.stderr, usage.WorktreeReset)
		}
	}
}

// keptWorktreeGrammars are the pool operations the worktree family help must keep naming.
// The surviving family route says nothing about the operations under it. Each operation is
// reachable only through that dispatcher, so the grammar line is where its survival shows.
var keptWorktreeGrammars = []string{
	// WF42: the full create grammar, so the family help pins `--from` rather than the
	// bare verb.
	usage.WorktreeCreate,
	"bench worktree path",
	"bench worktree exec",
	usage.WorktreeShow,
	usage.WorktreeBuild,
	"bench worktree release",
	"bench worktree clean",
	"bench worktree reauthorize",
	"bench worktree merge",
	usage.WorktreeReset,
	"bench worktree land",
}

// helpRows answers the help rows of one decoded document.
func helpRows(t *testing.T, document axitest.Document) []any {
	t.Helper()
	rows, err := document.Rows("help")
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// ER32, ER34: the related view offers the actions of the default view, which one action
// owner derives from each row state, and adds one read action for each active row with a
// present tree. The fixture's one active row has a present tree and a landed branch, so
// the default view offers path, exec, and the batch clean. The default view keeps its
// actions without the read action.
func TestRelatedWorktreeActions(t *testing.T) {
	root, active := relatedFixture(t, 0)
	result := runAXICommandAt(t, root, []string{"worktree", "list"})
	defaults, err := axitest.DecodeDocument(result.stdout)
	if result.code != 0 || err != nil {
		t.Fatalf("worktree list = %#v, %v", result, err)
	}
	want := helpRows(t, defaults)
	read := map[string]any{"cmd": "bench worktree read --in " + axi.ShellQuote(active.Label) + " <path>", "why": "read a file in an active worktree by its label"}
	if len(want) != 3 || !strings.HasPrefix(want[0].(map[string]any)["cmd"].(string), "bench worktree path ") || !strings.HasPrefix(want[1].(map[string]any)["cmd"].(string), "bench worktree exec ") || want[2].(map[string]any)["cmd"] != "bench worktree clean --landed" {
		t.Fatalf("default actions = %#v, want path, exec, and the batch clean", want)
	}
	path, exec, clean := want[0], want[1], want[2]
	_, related := relatedDocument(t, root, []string{"worktree", "list", "--spec", "x"})
	if got, want := helpRows(t, related), []any{path, exec, read, clean}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("related actions = %#v, want %#v", got, want)
	}

	// ER34: the read action appears only when its label names the row's tree alone. A row
	// whose label is the primary keyword, another active assignment's label, or a value that
	// the action owner rejects keeps its path and exec actions, and the other row keeps its
	// read action.
	for _, tc := range []struct {
		name  string
		label func(t *testing.T, root string, active intent.Assignment) string
	}{
		{"the primary keyword", func(*testing.T, string, intent.Assignment) string { return usage.PrimaryTreeTarget }},
		{"an unrelated active twin's label", func(t *testing.T, root string, active intent.Assignment) string {
			twin := putRelatedAssignment(t, root, relatedID(101), intent.StateActive, filepath.Join(root, "absent", "twin"))
			twin.Label = active.Label
			if err := intent.PutAssignment(root, twin); err != nil {
				t.Fatal(err)
			}
			return active.Label
		}},
		{"a label the action owner rejects", func(*testing.T, string, intent.Assignment) string { return "a <label>" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, active := relatedFixture(t, 0)
			active.Label = tc.label(t, root, active)
			if err := intent.PutAssignment(root, active); err != nil {
				t.Fatal(err)
			}
			other := putLinkedRelated(t, root, 102)
			read := map[string]any{"cmd": "bench worktree read --in " + axi.ShellQuote(other.Label) + " <path>", "why": read["why"]}
			_, related := relatedDocument(t, root, []string{"worktree", "list", "--spec", "x"})
			if got, want := helpRows(t, related), []any{path, exec, read, clean}; fmt.Sprint(got) != fmt.Sprint(want) {
				t.Fatalf("related actions with %s = %#v, want %#v", tc.name, got, want)
			}
		})
	}
}

// ER32: a --spec path names its spec from the caller's directory, an absolute path inside
// the repository names the same spec, and a slug names its live folder spec from anywhere.
// The projection spans several pages, and the walk follows every continuation from the
// caller's directory, so a continuation must repeat the caller's own operand.
func TestRelatedSpecOperandResolvesFromTheCaller(t *testing.T) {
	root, _ := relatedFixture(t, 40)
	sub := filepath.Join(root, "sub")
	writeAXIFixture(t, filepath.Join(sub, "keep.txt"), "")
	for _, tc := range []struct{ dir, operand string }{
		{sub, "../specs/x/spec.md"},
		{sub, filepath.Join(root, "specs", "x", "spec.md")},
		{sub, "x"},
		{root, "specs/x/spec.md"},
	} {
		content, pages := readPages(t, tc.dir, []string{"worktree", "list", "--spec", tc.operand})
		if pages < 2 {
			t.Fatalf("--spec %q from %s fits %d page, want a continuation", tc.operand, tc.dir, pages)
		}
		document, err := axitest.DecodeDocument(string(content))
		if err != nil {
			t.Fatalf("decode %q: %v", content, err)
		}
		summary, err := document.Rows("related")
		if err != nil || len(summary) != 1 {
			t.Fatalf("--spec %q from %s summary = %#v, %v", tc.operand, tc.dir, summary, err)
		}
		if cells := summary[0].(map[string]any); fmt.Sprint(cells["spec"], ",", cells["matched"]) != "specs/x/spec.md,41" {
			t.Fatalf("--spec %q from %s summary = %#v, want specs/x/spec.md with 41 matches", tc.operand, tc.dir, cells)
		}
	}
}
