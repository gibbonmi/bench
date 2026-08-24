package conformance

import (
	"slices"
	"testing"
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
