package conformance

import (
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/landing"
)

func TestLocalCaptureFilesAreNotTracked(t *testing.T) {
	tracked := trackedPaths(NewHarness(t).Root)
	for _, path := range []string{
		"capture/IDEAS.md",
		"capture/learnings.md",
		"capture/session-handoff.md",
	} {
		if slices.Contains(tracked, path) {
			t.Errorf("local capture file is tracked: %s", path)
		}
	}
}

// TestKitGitignoreIgnoresTheBuildLog keeps the kit's own build log out of history.
// The landing constant is the one source of the directory name.
func TestKitGitignoreIgnoresTheBuildLog(t *testing.T) {
	lines := strings.Split(NewHarness(t).ReadRootFile(".gitignore"), "\n")
	if !slices.Contains(lines, landing.BuildLogDir) {
		t.Errorf("the kit .gitignore holds no line %q", landing.BuildLogDir)
	}
}
