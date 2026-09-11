package worktree

import (
	"github.com/gibbonmi/bench/internal/landing"
	"testing"

	"github.com/gibbonmi/bench/internal/reviewrecord/recordtest"
)

// These fixtures model a completed review after their named source amendment.
func refreshLandingEvidence(t *testing.T, source, base string) {
	t.Helper()
	recordtest.RetainSingleChunk(t, source, "specs/x/spec.md", base)
}

// Retain the actual composed bytes as reviewed source before trying publication.
func foldCompletionComposition(t *testing.T, root, source, base, tip, destination string) string {
	t.Helper()
	result, err := landing.New().Compose(landing.CompositionRequest{Root: root, Destination: destination, Source: tip, ReviewBase: base})
	if err != nil || result.Tree == "" {
		t.Fatalf("compose fixture source: %+v, %v", result, err)
	}
	folded := gitOutput(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit-tree", result.Tree, "-p", tip, "-p", destination, "-m", "review composed capture")
	gitRun(t, source, "reset", "--hard", folded)
	refreshLandingEvidence(t, source, base)
	return gitOutput(t, source, "rev-parse", "HEAD")
}
