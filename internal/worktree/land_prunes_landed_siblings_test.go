package worktree

import (
	"bytes"
	"errors"
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
	root, creation, base, _, _, home := specLessLandingFixture(t, request)

	// The sibling assignment branch: one commit off the landing base, under the
	// assignment namespace, with no record and no checkout, so it reads as unclaimed.
	sibling := strings.TrimPrefix(intent.AssignmentBranchPrefix(), "refs/heads/") + "deadbeefdeadbeefdeadbeefdeadbeef/cafebabecafebabecafebabecafebabe"
	gitRun(t, root, "branch", sibling, base)
	siblingPath := t.TempDir()
	gitRun(t, root, "worktree", "add", "-q", "--detach", siblingPath, sibling)
	commitInWorktree(t, siblingPath, "sibling.txt", "sibling work\n", "sibling: the folded work")
	gitRun(t, root, "branch", "-f", sibling, gitOutput(t, siblingPath, "rev-parse", "HEAD"))
	gitRun(t, root, "worktree", "remove", "--force", siblingPath)

	// The fold is a squash, the shape a raw-git composition leaves: the sibling's content
	// reaches the landing source as a single-parent commit, so ancestry never proves it.
	gitRun(t, creation.Path, "cherry-pick", "--no-commit", "main.."+sibling)
	gitRun(t, creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "merge: sibling into the integration source")
	tip := gitOutput(t, creation.Path, "rev-parse", "HEAD")

	var stdout, stderr bytes.Buffer
	code := LandCommand(root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("land = (%d, %q, %q), want 0", code, stdout.String(), stderr.String())
	}
	unclaimed, err := UnclaimedAssignmentBranchRefs(root)
	mustNoError(t, err)
	if len(unclaimed) != 0 {
		t.Fatalf("landed sibling branch survived the landing as unclaimed: %q", unclaimed)
	}
}

// A prune that cannot read the repository leaves the landing durable but incomplete,
// and the record names the step so the resume repairs it.
func TestLandCommandReportsIncompletePrune(t *testing.T) {
	t.Parallel()
	request := "land-prune-incomplete"
	root, creation, base, tip, _, home := specLessLandingFixture(t, request)
	broken := defaultJoins()
	broken.pruneLandedBranches = func(string) (int, error) { return 0, errors.New("prune refused") }
	var stdout, stderr bytes.Buffer
	code := landWith(broken, root, home, "", specLessLandArgs(request, base, tip, creation.Path), &stdout, &stderr)
	if code != 3 || !strings.Contains(stdout.String(), "worktree=incomplete:prune") {
		t.Fatalf("interrupted prune = (%d, %q, %q), want exit 3 and incomplete:prune", code, stdout.String(), stderr.String())
	}
}
