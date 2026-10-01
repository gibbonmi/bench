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
	var stdout, stderr strings.Builder
	code := ReleaseCommand(f.creation.Path, f.home, []string{"--request", "landed-release-inside", f.creation.Path}, &stdout, &stderr)
	requireTest(t, code == 0 && strings.Contains(stdout.String(), ",complete,removed\n"), "release from inside = (%d, %q, %q), want exit 0 and a removed row", code, stdout.String(), stderr.String())
	_, statErr := os.Lstat(f.creation.Path)
	requireTest(t, os.IsNotExist(statErr), "release left %s: %v", f.creation.Path, statErr)
	_, err := assignmentByID(f.root, f.creation.Assignment.ID)
	requireTest(t, err != nil, "assignment %s survived a completed release", f.creation.Assignment.ID)
}
