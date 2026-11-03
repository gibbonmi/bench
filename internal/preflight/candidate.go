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
// exit 1.
func BuildVerdictAt(root, slug, base, candidate string) (string, int) {
	tip, err := git.ResolveCommit(root, candidate)
	if err != nil {
		return toon.Errorf("cannot resolve candidate", "'"+candidate+"' does not name a commit reachable in this repository") + "\n", 1
	}
	args := []string{modeBuild, slug}
	if base != "" {
		args = append(args, evidencecmd.FlagBase, base)
	}
	args = append(args, evidencecmd.FlagTip, tip)
	return verdictCommand(root, evidencecmd.KindVerdict, modeBuild, args, func() (Facts, *BootstrapFailure) {
		return gatherSource(root, modeBuild, slug, base, tip, tip, false, false)
	})
}
