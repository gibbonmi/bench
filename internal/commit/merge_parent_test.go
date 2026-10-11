package commit

import (
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/refusalroute"
	"github.com/gibbonmi/bench/internal/refusalroute/routetest"
	"github.com/gibbonmi/bench/internal/sanitize"
)

// These constants name the paths and the branch of a pending merge journey. Both sides
// change mergeConflictPath, and the incoming commit on the branch mergeIncoming modifies
// tracked.txt. The incoming commit can also add mergeAddedPath and mergeHostilePath.
// mergeHostilePath holds a space, a glob character, an ESC byte, and a newline.
const (
	mergeConflictPath = "conflict.txt"
	mergeAddedPath    = "added.txt"
	mergeHostilePath  = "esc\x1b name*\n.txt"
	mergeIncoming     = "incoming"
)

// mergeNameSet is the complete name set of a startPendingMerge journey that adds no path.
var mergeNameSet = []string{mergeConflictPath, "tracked.txt"}

// mergeNameSetRefusal opens the error line of a name set that omits a merged path.
const mergeNameSetRefusal = "error: pending merge changed paths that the commit does not name: "

// startPendingMerge starts a conflicted merge of an incoming commit in checkout, and it
// returns HEAD before the merge. It first resets the staged tracked.txt that
// prepareLandingCheckout leaves. The incoming commit on a side branch modifies tracked.txt,
// changes the conflicted path, and adds each path of added. The checkout's own commit
// changes the conflicted path too. A conflicted merge exits 1, so the merge runs through
// its own call, as stageHandMerge in internal/worktree does.
func startPendingMerge(t *testing.T, checkout string, added ...string) string {
	t.Helper()
	runGit(t, checkout, "reset", "-q", "--hard", "HEAD")
	runGit(t, checkout, "checkout", "-q", "-b", mergeIncoming)
	for _, path := range append([]string{"tracked.txt", mergeConflictPath}, added...) {
		mustWrite(t, filepath.Join(checkout, path), "incoming\n", 0o644)
	}
	runGit(t, checkout, "add", "-A")
	runGit(t, checkout, "commit", "-qm", "incoming")
	runGit(t, checkout, "checkout", "-q", "-")
	mustWrite(t, filepath.Join(checkout, mergeConflictPath), "ours\n", 0o644)
	runGit(t, checkout, "add", mergeConflictPath)
	runGit(t, checkout, "commit", "-qm", "ours")
	merge := exec.Command("git", "-C", checkout, "merge", "--no-commit", mergeIncoming)
	if out, err := merge.CombinedOutput(); err == nil || !strings.Contains(string(out), "CONFLICT") {
		t.Fatalf("merge = %v, %s; want the conflict", err, out)
	}
	return head(t, checkout)
}

// refusedLine fails t unless the commit at root refused with exit 1, printed one error line
// and the route of the registered face, and left HEAD at before. It returns the error line
// and the printed route steps.
func refusedLine(t *testing.T, root, before, face string, code int, stderr string) (string, []string) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(stderr, "\n"), "\n")
	next, printed := printedNext(stderr)
	if code != 1 || len(lines) != 2 || !printed {
		t.Fatalf("commit = (%d, %q), want exit 1, one error line, and one next= line", code, stderr)
	}
	var registered refusalroute.Face
	for _, candidate := range refusalroute.Inventory() {
		if candidate.Name == face {
			registered = candidate
		}
	}
	steps := routetest.Steps(t, registered, next, "", nil)
	if after := head(t, root); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
	return lines[0], steps
}

// TestCommitRefusesAnUnreadableMergeState proves that a directory at MERGE_HEAD refuses the
// commit through the reviewer's handback route, although the name set covers every merged
// path. (Coverage row MP25.)
func TestCommitRefusesAnUnreadableMergeState(t *testing.T) {
	root, _ := landingRepo(t, 0, func(*testing.T, string) {})
	before := startPendingMerge(t, root)
	mergeHead := adminPath(t, root, "MERGE_HEAD")
	mustRemove(t, mergeHead)
	mustMkdirAll(t, mergeHead)
	code, _, stderr := runCommand(t, root, append([]string{"-m", "m", "--"}, mergeNameSet...)...)
	if line, _ := refusedLine(t, root, before, faceHandback, code, stderr); !strings.HasPrefix(line, "error: read pending merge: ") {
		t.Fatalf("error line = %q, want the read pending merge refusal", line)
	}
}

// TestCommitRefusesAnIncompleteMergeNameSet proves that a commit during a pending merge
// refuses a name set that omits a modified or an added merged path. The refusal quotes each
// omitted path on its one error line, and its route ends with the caller's own commit.
// (Coverage rows MP32-MP34, MP41, and MP62.)
func TestCommitRefusesAnIncompleteMergeNameSet(t *testing.T) {
	quoted := sanitize.ShellQuote
	rerun := func(paths ...string) string {
		return callerCommit(quoted(landingLabel), append([]string{"-m", quoted("m"), "--"}, paths...)...)
	}
	for _, row := range []struct {
		name  string
		paths []string
		// omitted is each omitted path as the refusal quotes it.
		omitted, rerun string
	}{
		{"modified path omitted", []string{mergeConflictPath, mergeAddedPath, mergeHostilePath}, `"tracked.txt"`,
			rerun(quoted(mergeConflictPath), quoted(mergeAddedPath), "<path>")},
		{"added path omitted", []string{mergeConflictPath, "tracked.txt", mergeHostilePath}, `"added.txt"`,
			rerun(quoted(mergeConflictPath), quoted("tracked.txt"), "<path>")},
		{"two paths omitted", []string{mergeConflictPath, mergeAddedPath}, `"esc\x1b name*\n.txt" "tracked.txt"`,
			rerun(quoted(mergeConflictPath), quoted(mergeAddedPath))},
	} {
		t.Run(row.name, func(t *testing.T) {
			root, _ := landingRepo(t, 0, func(*testing.T, string) {})
			before := startPendingMerge(t, root, mergeAddedPath, mergeHostilePath)
			code, _, stderr := runCommand(t, root, append([]string{"-m", "m", "--"}, row.paths...)...)
			line, steps := refusedLine(t, root, before, faceNamedPath, code, stderr)
			if line != mergeNameSetRefusal+row.omitted || steps[len(steps)-1] != row.rerun {
				t.Fatalf("refusal = (%q, %q), want (%q, a route that ends with %q)", line, steps, mergeNameSetRefusal+row.omitted, row.rerun)
			}
		})
	}
}

// TestDryRunRefusesAnIncompleteMergeNameSet proves that a rehearsal refuses the name set
// that the commit refuses. (Coverage row MP39.)
func TestDryRunRefusesAnIncompleteMergeNameSet(t *testing.T) {
	root, _ := landingRepo(t, 0, func(*testing.T, string) {})
	before := startPendingMerge(t, root)
	code, _, stderr := runCommand(t, root, "--dry-run", "-m", "m", "--", mergeConflictPath)
	if line, _ := refusedLine(t, root, before, faceNamedPath, code, stderr); line != mergeNameSetRefusal+`"tracked.txt"` {
		t.Fatalf("error line = %q, want the refusal of the omitted tracked.txt", line)
	}
}

// TestMergeNameSetRefusalPrecedesCommitment proves that in an unbound planning worktree,
// whose production change the commitment refuses, the merge name-set refusal prints in
// place of the commitment refusal. (Coverage row MP40.)
func TestMergeNameSetRefusalPrecedesCommitment(t *testing.T) {
	root := planningCommitRepo(t)
	before := startPendingMerge(t, root)
	code, stdout, stderr := runCommand(t, root, "-m", "m", "--", mergeConflictPath)
	if code != 1 || !strings.Contains(stderr, mergeNameSetRefusal+`"tracked.txt"`) || strings.Contains(stderr, "commitment") {
		t.Fatalf("commit = (%d, %q, %q), want exit 1 with the merge name-set refusal and no commitment refusal", code, stdout, stderr)
	}
	if after := head(t, root); after != before {
		t.Fatalf("HEAD moved from %s to %s", before, after)
	}
}

// TestCommitRecordsTheMergeParent proves that after a conflicted merge, a commit of the
// complete name set publishes two parents. The parents are the previous tip and the incoming
// commit. Its record keeps the committed fields with the count of the name set. After the
// commit, the worktree status is clean, and a later commit of a new change publishes one
// parent. (Coverage rows MP19, MP44, MP46, and MP54.)
func TestCommitRecordsTheMergeParent(t *testing.T) {
	root, _ := landingRepo(t, 0, func(*testing.T, string) {})
	before := startPendingMerge(t, root)
	headParents := func() string {
		return strings.TrimSpace(string(runGit(t, root, "rev-list", "--parents", "-n", "1", "HEAD")))
	}
	incoming := strings.TrimSpace(string(runGit(t, root, "rev-parse", mergeIncoming)))
	mustWrite(t, filepath.Join(root, mergeConflictPath), "resolution\n", 0o644)
	code, stdout, stderr := runCommand(t, root, append([]string{"-m", "merge", "--"}, mergeNameSet...)...)
	if code != 0 {
		t.Fatalf("commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	merge := head(t, root)
	if parents, want := headParents(), merge+" "+before+" "+incoming; parents != want {
		t.Fatalf("merge commit and parents = %q, want %q", parents, want)
	}
	lines := strings.Split(strings.TrimSuffix(stdout, "\n"), "\n")
	name, fields, order := recordFields(t, lines[len(lines)-1])
	if name != "committed" || !reflect.DeepEqual(order, []string{"published_commit", "path_count"}) ||
		fields["published_commit"] != merge || fields["path_count"] != strconv.Itoa(len(mergeNameSet)) {
		t.Fatalf("stdout = %q, want committed{published_commit=%s,path_count=%d}", stdout, merge, len(mergeNameSet))
	}
	if status := string(runGit(t, root, "status", "--porcelain=v1")); status != "" {
		t.Fatalf("status after the merge commit = %q, want a clean worktree", status)
	}
	mustWrite(t, filepath.Join(root, "a.txt"), "after the merge\n", 0o644)
	if code, stdout, stderr := runCommand(t, root, "-m", "after", "--", "a.txt"); code != 0 {
		t.Fatalf("second commit = (%d, %q, %q), want exit 0", code, stdout, stderr)
	}
	if parents, want := headParents(), head(t, root)+" "+merge; parents != want {
		t.Fatalf("second commit and parents = %q, want %q", parents, want)
	}
}
