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
	root, creation, home := newOwnedAssignment(t, "release-inside")
	var stdout, stderr strings.Builder
	code := ReleaseCommand(creation.Path, home, []string{"--request", "landed-release-inside", creation.Path}, &stdout, &stderr)
	requireTest(t, code == 0 && strings.Contains(stdout.String(), ",complete,removed\n"), "release from inside = (%d, %q, %q), want exit 0 and a removed row", code, stdout.String(), stderr.String())
	_, statErr := os.Lstat(creation.Path)
	requireTest(t, os.IsNotExist(statErr), "release left %s: %v", creation.Path, statErr)
	_, err := assignmentByID(root, creation.Assignment.ID)
	requireTest(t, err != nil, "assignment %s survived a completed release", creation.Assignment.ID)
}
