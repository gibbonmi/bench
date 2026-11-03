package preflight

import (
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/toon"
)

// headRef is the head of every published grade. The candidate grade passes its candidate
// commit in its place.
const headRef = "HEAD"

// BuildVerdictAt grades candidate, a commit that no ref has to name yet. It prints the
// verdict text and exit that `bench preflight build <slug> [--base <base>] --source-tip
// <candidate>` prints after HEAD moves to that candidate, because the gather reads the
// candidate as its head in place of HEAD. The candidate is also the source-tip pin, so the
// verdict holds the tip-current row. A candidate that does not name a commit refuses with
// exit 1. A snapshot drift refuses with no retry action, because the published command
// reds tip-current until HEAD names the candidate.
func BuildVerdictAt(root, slug, base, candidate string) (string, int) {
	tip, err := git.ResolveCommit(root, candidate)
	if err != nil {
		return toon.Errorf("cannot resolve candidate", "'"+candidate+"' does not name a commit reachable in this repository") + "\n", 1
	}
	return verdictCommand(root, evidencecmd.KindVerdict, modeBuild, nil, func() (Facts, *BootstrapFailure) {
		return gatherSource(root, modeBuild, slug, base, tip, tip, false, false)
	})
}
