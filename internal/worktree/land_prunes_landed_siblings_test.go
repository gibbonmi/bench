package worktree

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/intent"
)

// The landing is the one moment the repository can prove a squash-folded sibling landed:
// its content is on main and no later commit has moved the files it touched. A sibling
// assignment branch whose content the landing carries must not survive the landing, or
// the next retirement, which deletes the spec the sibling also edited, leaves a branch
// nothing can prove landed and status routes it to a destructive discard.
func TestLandCommandPrunesSquashFoldedSiblingBranch(t *testing.T) {
	t.Parallel()
	request := "land-prunes-sibling"
	f := specLessLandingFixture(t, request)

	// The sibling assignment branch: one commit off the landing base, under the
	// assignment namespace, with no record and no checkout, so it reads as unclaimed.
	sibling := strings.TrimPrefix(intent.AssignmentBranchPrefix(), "refs/heads/") + "deadbeefdeadbeefdeadbeefdeadbeef/cafebabecafebabecafebabecafebabe"
	gitRun(t, f.root, "branch", sibling, f.base)
	siblingPath := t.TempDir()
	gitRun(t, f.root, "worktree", "add", "-q", "--detach", siblingPath, sibling)
	commitInWorktree(t, siblingPath, "sibling.txt", "sibling work\n", "sibling: the folded work")
	gitRun(t, f.root, "branch", "-f", sibling, gitOutput(t, siblingPath, "rev-parse", "HEAD"))
	gitRun(t, f.root, "worktree", "remove", "--force", siblingPath)

	// The fold is a squash, the shape a raw-git composition leaves: the sibling's content
	// reaches the landing source as a single-parent commit, so ancestry never proves it.
	gitRun(t, f.creation.Path, "cherry-pick", "--no-commit", "main.."+sibling)
	gitRun(t, f.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "merge: sibling into the integration source")
	tip := gitOutput(t, f.creation.Path, "rev-parse", "HEAD")

	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, tip, f.creation.Path)...))
	if r.exit != 0 {
		t.Fatalf("land = (%d, %q, %q), want 0", r.exit, r.stdout, r.stderr)
	}
	unclaimed, err := CountUnclaimedRefs(f.root)
	mustNoError(t, err)
	if unclaimed.Rows() != 0 {
		t.Fatalf("landed sibling branch survived the landing as unclaimed: %+v", unclaimed)
	}
}

// A prune that cannot delete a landed branch leaves the landing durable but incomplete,
// and the record names the step so the resume repairs it.
func TestLandCommandReportsIncompletePrune(t *testing.T) {
	t.Parallel()
	request := "land-prune-incomplete"
	f := specLessLandingFixture(t, request)
	lockLandedSibling(t, f.root, f.base)
	r := runVerb(t, verbLand, f.call(specLessLandArgs(request, f.base, f.tip, f.creation.Path)...))
	if r.exit != 3 || !strings.Contains(r.stdout, "worktree=incomplete:prune") {
		t.Fatalf("interrupted prune = (%d, %q, %q), want exit 3 and incomplete:prune", r.exit, r.stdout, r.stderr)
	}
}

// lockLandedSibling creates a sibling branch at landed, a commit the destination already
// holds, and plants a stale lock for that branch. The landing's prune then reads the
// branch as landed, and the real branch delete fails on the lock.
func lockLandedSibling(t *testing.T, root, landed string) {
	t.Helper()
	const sibling = "landed-sibling"
	gitRun(t, root, "branch", sibling, landed)
	common := gitOutput(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir")
	mustWrite(t, filepath.Join(common, "refs", "heads", sibling+".lock"), nil, 0o644)
}
