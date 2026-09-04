package handoff

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/status"
)

// TestRenderPathAbbreviates pins the `~` form, including exactly at $HOME where the
// remainder is empty and the abbreviation is the whole path.
func TestRenderPathAbbreviates(t *testing.T) {
	cases := []struct{ root, home, want string }{
		{"/home/a", "/home/a", "~"},
		{"/home/a/", "/home/a", "~"},
		{"/home/a", "/home/a/", "~"},
		{"/home/a/workspace/bench", "/home/a", "~/workspace/bench"},
		{"/home/a/x", "/home/a", "~/x"},
		{"/", "/", "~"},
		{"/srv/x", "/", "~/srv/x"},
	}
	for _, tc := range cases {
		if got := renderPath(tc.root, tc.home); got != tc.want {
			t.Fatalf("renderPath(%q, %q) = %q, want %q", tc.root, tc.home, got, tc.want)
		}
	}
}

// TestRenderPathOutsideHome pins the other side of the boundary. A prefix match on the
// raw string would turn /home/abc into ~bc: a path that resolves nowhere.
func TestRenderPathOutsideHome(t *testing.T) {
	cases := []struct{ root, home, want string }{
		{"/home/abc", "/home/a", "/home/abc"},
		{"/home/abc/deep", "/home/a", "/home/abc/deep"},
		{"/srv/checkouts/bench", "/home/a", "/srv/checkouts/bench"},
		{"/home/a", "", "/home/a"},
		{"/home/a", "relative/home", "/home/a"},
		{"/home", "/home/a", "/home"},
	}
	for _, tc := range cases {
		if got := renderPath(tc.root, tc.home); got != tc.want {
			t.Fatalf("renderPath(%q, %q) = %q, want %q", tc.root, tc.home, got, tc.want)
		}
	}
}

// HS2 and HS10, stories 4 and 13. The six pins carry the values a resuming session acts
// on. The request token is the plain caller token byte for byte, because a landing passes
// that value back and the digest is not accepted in its place. The tip is the assignment's
// own worktree HEAD, read after a commit that only that tree carries.
func TestSectionRendersTheAssignmentPins(t *testing.T) {
	root := benchRepo(t)
	const token = "hs-plain-token-20260904"
	a := activeAssignment(t, root, token, "hs-resolve")
	write(t, filepath.Join(a.Worktree, "specs", "alpha", "spec.md"), "# Alpha\n\nStatus: staged\n")
	commitIn(t, a.Worktree, "worktree commit")
	tip := gitOut(t, a.Worktree, "rev-parse", "HEAD")

	runIn(t, a.Worktree, nil)
	got := sectionBytes(t, filepath.Join(root, status.HandoffFile), a.Request)
	for _, want := range []string{
		handoffdoc.LabelRequestToken + ": " + token,
		handoffdoc.LabelLabel + ": hs-resolve",
		handoffdoc.LabelWorktreeTip + ": " + tip,
		handoffdoc.LabelRecordedBase + ": " + a.Start,
		handoffdoc.LabelSpec + ": specs/alpha/spec.md",
		handoffdoc.LabelSpecStatus + ": staged",
	} {
		if !strings.Contains(got, want+"\n") {
			t.Errorf("section carries no %q line\n%s", want, got)
		}
	}
	if tip == a.Start {
		t.Fatalf("the fixture left the tip equal to the base %q, so the tip pin proves nothing", tip)
	}
}

// HS3, story 5. A worktree that holds two staged specs renders both pairs. One pair would
// lose a spec the owning phase is still building.
func TestSectionRendersOnePairPerLiveSpec(t *testing.T) {
	root := benchRepo(t)
	a := activeAssignment(t, root, "hs-two-specs", "hs-two-specs")
	write(t, filepath.Join(a.Worktree, "specs", "alpha", "spec.md"), "# Alpha\n\nStatus: staged\n")
	write(t, filepath.Join(a.Worktree, "specs", "beta", "spec.md"), "# Beta\n\nStatus: drafting\n")

	runIn(t, a.Worktree, nil)
	got := sectionBytes(t, filepath.Join(root, status.HandoffFile), a.Request)
	if n := strings.Count(got, handoffdoc.LabelSpec+": "); n != 2 {
		t.Fatalf("section carries %d %q lines, want 2\n%s", n, handoffdoc.LabelSpec, got)
	}
	for _, want := range []string{
		handoffdoc.LabelSpec + ": specs/alpha/spec.md\n" + handoffdoc.LabelSpecStatus + ": staged\n",
		handoffdoc.LabelSpec + ": specs/beta/spec.md\n" + handoffdoc.LabelSpecStatus + ": drafting\n",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("section carries no pair %q\n%s", want, got)
		}
	}
}

// HS9, story 12. The header states the repository the document lives in, so it pins main's
// HEAD and never the caller worktree's. A header written from the caller's tree would put
// one phase's branch above every phase's section.
func TestHeaderPinsTheMainHeadAndNotTheWorktreeHead(t *testing.T) {
	root := benchRepo(t)
	a := activeAssignment(t, root, "hs-header", "hs-header")
	commitIn(t, a.Worktree, "worktree only")
	mainHead := status.Short(gitOut(t, root, "rev-parse", "HEAD"))
	worktreeHead := status.Short(gitOut(t, a.Worktree, "rev-parse", "HEAD"))

	runIn(t, a.Worktree, nil)
	document := read(t, filepath.Join(root, status.HandoffFile))
	head, _, _ := strings.Cut(document, "\n## ")
	for _, want := range []string{
		"Repository: `" + filepath.Base(root) + "`",
		"Path: `",
		"HEAD `" + mainHead + "`",
		"Gate: ",
	} {
		if !strings.Contains(head, want) {
			t.Errorf("header carries no %q\n%s", want, head)
		}
	}
	if strings.Contains(head, worktreeHead) {
		t.Fatalf("header pins the worktree HEAD %s rather than main's %s\n%s", worktreeHead, mainHead, head)
	}
}

func TestCommandUsesCleanBoardRouteFallback(t *testing.T) {
	root := benchRepo(t)
	runIn(t, root, nil)
	assertNextCommand(t, root, "`bench roadmap`")
}

func TestCommandUsesCodexRouteForPhase(t *testing.T) {
	root := benchRepo(t)
	write(t, filepath.Join(root, "capture", "IDEAS.md"), "- 2026-08-18  pending\n")
	runIn(t, root, []string{"--harness", "codex"})
	assertNextCommand(t, root, "`$bench-drain`")
}

// HC22. `--harness none` is accepted, and the routed shell command is the one a cold
// session without a harness reads. The board is led by `git push`: a clean tree one commit
// ahead of its upstream.
func TestCommandUsesShellRouteForNoHarness(t *testing.T) {
	root := benchRepo(t)
	remote := filepath.Join(t.TempDir(), "remote.git")
	for _, args := range [][]string{
		{"init", "-q", "--bare", remote},
		{"-C", root, "remote", "add", "origin", remote},
		{"-C", root, "push", "-q", "-u", "origin", "HEAD"},
		{"-C", root, "commit", "-q", "--allow-empty", "-m", "ahead"},
	} {
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", args, err, out)
		}
	}
	runIn(t, root, nil)
	assertNextCommand(t, root, "`git push`")
}

// assertNextCommand pins the routed value on the label line the leaf grammar renders.
func assertNextCommand(t *testing.T, root, want string) {
	t.Helper()
	got := sectionBytes(t, filepath.Join(root, status.HandoffFile), handoffdoc.MainKey)
	if !strings.Contains(got, handoffdoc.LabelNextCommand+": "+want) {
		t.Fatalf("next command line does not carry %s\n%s", want, got)
	}
}

// commitIn lands one commit in the tree at dir, so that tree's HEAD moves past the base
// every other checkout still sits on.
func commitIn(t *testing.T, dir, message string) {
	t.Helper()
	write(t, filepath.Join(dir, "note.md"), message+"\n")
	gitRun(t, dir, "add", "-A")
	gitRun(t, dir, "-c", "user.email=t@example.com", "-c", "user.name=t", "commit", "-qm", message)
}
