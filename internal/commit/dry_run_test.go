package commit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// A dry run answers "would this composed set land green" without paying a junk
// commit: the same compose-and-authorize half a landing runs, then a full stop.
func TestDryRunReportsGreenAndMovesNothing(t *testing.T) {
	root, before := landingRepo(t, 0, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "prospective\n", 0o644)
	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "m", "a.txt")
	if code != 0 {
		t.Fatalf("exit = %d, want 0; stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout, "nothing committed") {
		t.Fatalf("stdout = %q, want the dry-run verdict line", stdout)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatal("dry run moved HEAD")
	}
	if headHasPrefix(t, root, "a.txt") {
		t.Fatalf("dry run published the path: %v", headPaths(t, root))
	}
}

// A red composed set reports the refusal as the diagnosis and still moves nothing.
func TestDryRunRedReportsRefusalAndMovesNothing(t *testing.T) {
	root, before := landingRepo(t, 1, func(t *testing.T, root string) {})
	runGit(t, root, "reset", "-q", "--hard", "HEAD")
	mustWrite(t, filepath.Join(root, "a.txt"), "prospective\n", 0o644)
	code, stdout, stderr := runCommand(t, root, "--dry-run", "-m", "m", "a.txt")
	if code != 1 {
		t.Fatalf("exit = %d, want 1; stdout=%q stderr=%q", code, stdout, stderr)
	}
	want := "prospective authorization refused: inherited (the gate ran red on the composed tree and no green baseline attributes the red to this diff)"
	if !strings.Contains(stderr, want) {
		t.Fatalf("stderr = %q, want the operator-facing inherited refusal %q", stderr, want)
	}
	if after := strings.TrimSpace(string(runGit(t, root, "rev-parse", "HEAD"))); after != before {
		t.Fatal("red dry run moved HEAD")
	}
}

// A grammar error prints the one-line usage alone; the example stays a help-only cost.
func TestGrammarErrorPrintsNoExample(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, _, stderr := runCommand(t, root, "a.txt")
	if code != 2 || strings.Contains(stderr, "example:") {
		t.Fatalf("grammar error = (%d, %q), want exit 2 with no example line", code, stderr)
	}
}

// The help text advertises the flag the grammar accepts. Its --dry-run line names the
// declared lane that a worktree commit runs and the gate as the fallback when no lane
// is declared. The expectation is an independent literal, so a negated or truncated
// clause reds it.
func TestHelpAdvertisesDryRun(t *testing.T) {
	root, _ := landingRepo(t, 0, func(t *testing.T, root string) {})
	code, stdout, _ := runCommand(t, root, "--help")
	if code != 0 || !strings.Contains(stdout, "--dry-run") {
		t.Fatalf("help = (%d, %q), want --dry-run advertised", code, stdout)
	}
	if strings.Contains(stdout, "gate the exact composed snapshot") {
		t.Errorf("help = %q, want no claim that --dry-run gates the exact composed snapshot", stdout)
	}
	var dryRunLine string
	for _, line := range strings.Split(stdout, "\n") {
		if strings.HasPrefix(line, "--dry-run:") {
			dryRunLine = line
		}
	}
	if !strings.HasPrefix(dryRunLine, "--dry-run: run the declared lane (or the gate when no lane is declared) on ") {
		t.Errorf("help --dry-run line = %q, want it to name the declared lane and the gate only when no lane is declared", dryRunLine)
	}
}

// unownedCommitSet is the commit fixture whose checkout is a linked worktree that no
// assignment owns.
func unownedCommitSet(t *testing.T) commitSet {
	t.Helper()
	f := primaryCommitSet(t)
	f.checkout = filepath.Join(t.TempDir(), "linked")
	runGit(t, f.primary, "worktree", "add", "-q", "-b", "topic", f.checkout)
	return f
}

// TestNamedPathRefusalReRunsTheCallersCommit holds the agent-clearable causes and the
// hostile arguments of the commit re-run. A commit whose named path or file the caller
// corrects prints an agent route, and the route ends with the caller's own commit: each
// line-safe value shell-quoted, each value that is not line-safe as its slot placeholder,
// and the label placeholder when no assignment owns the checkout.
func TestNamedPathRefusalReRunsTheCallersCommit(t *testing.T) {
	quoted := sanitize.ShellQuote
	const label = "rerun"
	assigned := func(t *testing.T) commitSet { return assignedCommitSet(t, label, gateScript("exit 0")) }
	// written is an assigned checkout that also holds path with content and mode.
	written := func(path, content string, mode os.FileMode) func(t *testing.T) commitSet {
		return func(t *testing.T) commitSet {
			f := assigned(t)
			mustWrite(t, filepath.Join(f.checkout, path), content, mode)
			return f
		}
	}
	// rerun is the caller's commit of path with message m, the re-run of most rows.
	rerun := func(path string) string { return callerCommit(quoted(label), "-m", quoted("m"), "--", path) }
	for _, row := range []struct {
		name  string
		build func(t *testing.T) commitSet
		args  []string
		rerun string
	}{
		{"missing path", assigned, []string{"-m", "m", "--", "missing.txt"}, rerun(quoted("missing.txt"))},
		{"path outside the repository", assigned, []string{"-m", "m", "--", "../outside"}, rerun(quoted("../outside"))},
		{"repository root", assigned, []string{"-m", "m", "--", "."}, rerun(quoted("."))},
		{"nothing to commit", assigned, []string{"-m", "m", "--", "tracked.txt"}, rerun(quoted("tracked.txt"))},
		{"unreadable path", written("secret.txt", "x\n", 0), []string{"-m", "m", "--", "secret.txt"}, rerun(quoted("secret.txt"))},
		{"Go file that does not parse", written("bad.go", "package\n", 0o644), []string{"-m", "m", "--", "bad.go"}, rerun(quoted("bad.go"))},
		{"dry run", assigned, []string{"--dry-run", "-m", "m", "--", "missing.txt"}, callerCommit(quoted(label), "--dry-run", "-m", quoted("m"), "--", quoted("missing.txt"))},
		{"preflight build", assigned, []string{PreflightBuildFlag, "commit-fixture", "-m", "m", "--", "missing.txt"}, callerCommit(quoted(label), PreflightBuildFlag, quoted("commit-fixture"), "-m", quoted("m"), "--", quoted("missing.txt"))},
		{"multi-line message", assigned, []string{"-m", "one\ntwo", "--", "missing.txt"}, callerCommit(quoted(label), "-m", "<msg>", "--", quoted("missing.txt"))},
		{"path with a space", assigned, []string{"-m", "m", "--", "missing file.txt"}, rerun(quoted("missing file.txt"))},
		{"path with a control byte", assigned, []string{"-m", "m", "--", "esc\x1b.txt"}, rerun("<path>")},
		{"no owning assignment", unownedCommitSet, []string{"-m", "m", "--", "missing.txt"}, callerCommit("<label>", "-m", quoted("m"), "--", quoted("missing.txt"))},
	} {
		t.Run(row.name, func(t *testing.T) {
			f := row.build(t)
			code, stdout, stderr := runCommand(t, f.checkout, row.args...)
			next, printed := printedNext(stderr)
			steps := refusalroute.Steps(next)
			if code != 1 || !printed || strings.HasPrefix(next, routetest.ReviewerMarker) || steps[len(steps)-1] != row.rerun {
				t.Fatalf("commit = (%d, %q, %q), want exit 1 and an agent next= route that ends with %q", code, stdout, stderr, row.rerun)
			}
		})
	}
}
