package worktree

import (
	"os"
	"strings"
	"testing"
)

// TestReleaseFromInsideItsOwnTreeCompletes is the FT345 release face. The CLI roots a
// release that runs inside its own tree at that tree, so the root is gone after the
// removal. The release must still complete the assignment and exit 0.
func TestReleaseFromInsideItsOwnTreeCompletes(t *testing.T) {
	t.Parallel()
	f := newOwnedAssignment(t, "release-inside")
	release := runVerb(t, verbRelease, repoHome{f.creation.Path, f.home}.call("--request", "landed-release-inside", f.creation.Path))
	requireTest(t, release.exit == 0 && strings.Contains(release.stdout, ",complete,removed\n"), "release from inside = (%d, %q, %q), want exit 0 and a removed row", release.exit, release.stdout, release.stderr)
	_, statErr := os.Lstat(f.creation.Path)
	requireTest(t, os.IsNotExist(statErr), "release left %s: %v", f.creation.Path, statErr)
	_, err := assignmentByID(f.root, f.creation.Assignment.ID)
	requireTest(t, err != nil, "assignment %s survived a completed release", f.creation.Assignment.ID)
}
