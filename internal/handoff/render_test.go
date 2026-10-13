package handoff

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/handoffdoc"
	"github.com/gibbonmi/bench/internal/intent"
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

// The pins carry the values a resuming session acts on. The request token is the plain
// caller token byte for byte, because a landing passes that value back and the digest is
// not accepted in its place. The tip is the assignment's own worktree HEAD, read after a
// commit that only that tree carries.
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

// The line after a live spec's status names the outcomes that the policy on main binds to
// that path, in policy order. An outcome of a milestone that is not active still binds. A
// binding that only the worktree commits reads unbound.
func TestSectionRendersTheSpecBinding(t *testing.T) {
	const alpha, beta = "specs/alpha/spec.md", "specs/beta/spec.md"
	for _, tc := range []struct {
		name string
		seed func(t *testing.T, root string) (a intent.Assignment, path, want string)
	}{
		{"SB23 one binding outcome", func(t *testing.T, root string) (intent.Assignment, string, string) {
			id := commitBinding(t, root, alpha)
			return activeAssignment(t, root, "sb23", "sb23"), alpha, id
		}},
		{"SB23 binding outside the active milestone", func(t *testing.T, root string) (intent.Assignment, string, string) {
			commitmenttest.Write(t, root, alpha, commitmenttest.StagedBody)
			commitmenttest.SeedRowless(t, root, alpha)
			var id string
			commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
				bound := &policy.Milestones[0]
				active := commitment.Milestone{ID: "later", Outcomes: bound.Outcomes[1:]}
				id, bound.Outcomes = bound.Outcomes[0].ID, bound.Outcomes[:1]
				policy.Milestones, policy.ActiveMilestone = append(policy.Milestones, active), active.ID
			})
			commitmenttest.Commit(t, root, "activate a later milestone")
			return activeAssignment(t, root, "sb23-later", "sb23-later"), alpha, id
		}},
		{"SB24 another path bound", func(t *testing.T, root string) (intent.Assignment, string, string) {
			commitBinding(t, root, alpha)
			a := activeAssignment(t, root, "sb24", "sb24")
			commitmenttest.Write(t, a.Worktree, beta, commitmenttest.StagedBody)
			return a, beta, unbound
		}},
		{"SB25 bound only in the worktree", func(t *testing.T, root string) (intent.Assignment, string, string) {
			commitBinding(t, root, alpha)
			a := activeAssignment(t, root, "sb25", "sb25")
			commitBinding(t, a.Worktree, beta)
			return a, beta, unbound
		}},
		{"SB26 two binding outcomes", func(t *testing.T, root string) (intent.Assignment, string, string) {
			commitmenttest.Write(t, root, alpha, commitmenttest.StagedBody)
			commitmenttest.SeedRowless(t, root, alpha)
			var first, second string
			commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
				outcomes := policy.Milestones[0].Outcomes
				outcomes[1].Deliverables = outcomes[0].Deliverables
				first, second = outcomes[0].ID, outcomes[1].ID
			})
			commitmenttest.Commit(t, root, "bind two outcomes")
			return activeAssignment(t, root, "sb26", "sb26"), alpha, first + " " + second
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := benchRepo(t)
			a, path, want := tc.seed(t, root)
			runIn(t, a.Worktree, nil)
			requireSpecBinding(t, root, a.Request, path, want)
		})
	}
}

// A policy that main does not publish reads unbound, and a policy that does not parse reads
// unknown. The handoff still writes, so a commitment fault never blocks a phase close.
func TestSectionRendersAnUnavailableBinding(t *testing.T) {
	const alpha = "specs/alpha/spec.md"
	t.Run("SB27 no policy on main", func(t *testing.T) {
		root := benchRepo(t)
		a := activeAssignment(t, root, "sb27", "sb27")
		commitmenttest.Write(t, a.Worktree, alpha, commitmenttest.StagedBody)
		runIn(t, a.Worktree, nil)
		requireSpecBinding(t, root, a.Request, alpha, unbound)
	})
	t.Run("SB28 unparseable policy on main", func(t *testing.T) {
		root := benchRepo(t)
		commitmenttest.Write(t, root, commitment.PolicyPath, "")
		commitmenttest.Commit(t, root, "empty policy")
		a := activeAssignment(t, root, "sb28", "sb28")
		commitmenttest.Write(t, a.Worktree, alpha, commitmenttest.StagedBody)
		runIn(t, a.Worktree, nil)
		requireSpecBinding(t, root, a.Request, alpha, unknownBinding)
	})
	// Two local branches, neither one main, and no origin/HEAD leave no default branch.
	t.Run("SB44 no default branch resolves", func(t *testing.T) {
		root := benchRepo(t)
		commitBinding(t, root, alpha)
		gitRun(t, root, "branch", "-m", "main", "trunk")
		gitRun(t, root, "branch", "other")
		runIn(t, root, []string{"--next", "bench status"})
		requireSpecBinding(t, root, handoffdoc.MainKey, alpha, unbound)
	})
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

// HS11, story 14. A mid-build resume invocation carries flags no board can derive, so a
// run without `--next` keeps it byte for byte. The board here leads with `/bench-drain`,
// which is the value a regenerating run would have written over it.
//
// The whole line is compared. A `Contains` check passes a rewrite that kept the command
// and dropped its flags, which is the exact loss this row exists to catch.
func TestCommandKeepsANonEmptyNextCommand(t *testing.T) {
	root := benchRepo(t)
	document := filepath.Join(root, status.HandoffFile)
	write(t, filepath.Join(root, "capture", "IDEAS.md"), "- 2026-08-18  pending\n")
	const planted = "`bench worktree exec hs-next -- bench gate --check prose --fix`"
	seedSection(t, document, planted, "The build is live.")

	runIn(t, root, nil)
	if got := nextCommandLine(t, root); got != handoffdoc.LabelNextCommand+": "+planted {
		t.Fatalf("next command line = %q, want %q", got, handoffdoc.LabelNextCommand+": "+planted)
	}
}

// HS12, story 15. A section that names no command gets the board's. The rendered form of a
// named command is a backticked span, so an empty pair of backticks says nothing at all;
// a rule that read them as content would leave that section routeless forever.
func TestCommandRoutesOnlyAnEmptyNextCommand(t *testing.T) {
	committed := time.Date(2026, 8, 18, 12, 0, 0, 0, time.UTC)
	written := committed.Add(time.Minute)
	for _, tc := range []struct{ name, planted string }{
		{"empty", ""},
		{"whitespace only", "   "},
		{"empty backticks", "``"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("GIT_AUTHOR_DATE", committed.Format(time.RFC3339))
			t.Setenv("GIT_COMMITTER_DATE", committed.Format(time.RFC3339))
			root := benchRepo(t)
			document := filepath.Join(root, status.HandoffFile)
			seedSection(t, document, tc.planted, "The build is live.")
			if err := os.Chtimes(document, written, written); err != nil {
				t.Fatal(err)
			}
			runIn(t, root, nil)
			assertNextCommand(t, root, "`bench roadmap`")
		})
	}
}

// nextCommandLine returns main's whole `Next command` line, so a test can compare the
// value rather than probe it for a substring.
func nextCommandLine(t *testing.T, root string) string {
	t.Helper()
	got := sectionBytes(t, filepath.Join(root, status.HandoffFile), handoffdoc.MainKey)
	for line := range strings.SplitSeq(got, "\n") {
		if strings.HasPrefix(line, handoffdoc.LabelNextCommand+":") {
			return line
		}
	}
	t.Fatalf("section carries no %q line\n%s", handoffdoc.LabelNextCommand, got)
	return ""
}

// assertNextCommand pins the routed value on the label line the leaf grammar renders.
func assertNextCommand(t *testing.T, root, want string) {
	t.Helper()
	got := sectionBytes(t, filepath.Join(root, status.HandoffFile), handoffdoc.MainKey)
	if !strings.Contains(got, handoffdoc.LabelNextCommand+": "+want) {
		t.Fatalf("next command line does not carry %s\n%s", want, got)
	}
}

// commitBinding writes a staged spec at path in dir and commits a policy whose one outcome
// binds that spec. It returns the id of that outcome.
func commitBinding(t *testing.T, dir, path string) string {
	t.Helper()
	commitmenttest.Write(t, dir, path, commitmenttest.StagedBody)
	commitmenttest.SeedAdmission(t, dir, path)
	commitmenttest.Commit(t, dir, "bind "+path)
	return commitmenttest.DeliveryOutcome
}

// requireSpecBinding requires the section key to carry the live spec at path, its staged
// status line, and then its binding line with the value want.
func requireSpecBinding(t *testing.T, root, key, path, want string) {
	t.Helper()
	got := sectionBytes(t, filepath.Join(root, status.HandoffFile), key)
	lines := handoffdoc.LabelSpec + ": " + path + "\n" + handoffdoc.LabelSpecStatus + ": staged\n" + handoffdoc.LabelSpecBinding + ": " + want + "\n"
	if !strings.Contains(got, lines) {
		t.Fatalf("section carries no %q\n%s", lines, got)
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
