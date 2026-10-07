package preflight

import (
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/tickets"
)

func pathsAuthorizedCheck(f Facts) CheckResult {
	if !f.ReviewBaseResolved {
		return red("paths-authorized", "review base does not resolve: "+f.ReviewBaseHint)
	}
	if unauthorized := unauthorizedPaths(f); len(unauthorized) > 0 {
		return red("paths-authorized", "not authorized by any ownership fence: "+strings.Join(unauthorized, ", "))
	}
	return green("paths-authorized")
}

// unauthorizedPaths is every changed path that no authorizing entry covers and that the
// default-branch tip does not already hold, in changed-path order. The row's detail sentence
// and the landing's typed path list both read this one answer, so a report and a refusal can
// never name a different set.
func unauthorizedPaths(f Facts) []string {
	var unauthorized []string
	entries := authorizingEntries(f)
	for _, p := range f.ChangedPaths {
		if !slices.Contains(f.DefaultTipPaths, p) && !fenceAuthorizes(p, entries) {
			unauthorized = append(unauthorized, p)
		}
	}
	return unauthorized
}

// captureEntry is the phase-owned capture folder. Every phase close writes the
// handoff, the learnings, or the ideas there, so a reviewed range always carries it,
// and no spec fences it. It is authorized for every range.
const captureEntry = "capture"

// authorizingEntries is every entry paths-authorized consults: the spec's
// declared fence entries plus the implicit entries. They are appended to a
// copy so the gathered FenceEntries slice is never mutated. Mode is
// deliberately not consulted: build preflight, review preflight, and the
// landing's final source authorization all get the same answer.
func authorizingEntries(f Facts) []string {
	return append(append([]string{}, f.FenceEntries...), implicitEntries(f)...)
}

// specFolder is the directory containing the resolved spec, empty when the
// spec path carries no directory at all. The result is an ordinary fence
// entry, so tickets.Covers grades it, never a second prefix rule.
func specFolder(specPath string) string {
	i := strings.LastIndex(specPath, "/")
	if i < 0 {
		return ""
	}
	return specPath[:i]
}

// fenceAuthorizes reports whether one of the spec's declared fence entries
// covers path. Each fence entry takes the one `Writes:` split, so a trailing
// directory slash is never itself an extra path segment.
func fenceAuthorizes(path string, fences []string) bool {
	for _, fence := range fences {
		if entry, _ := tickets.WritesPath(fence); tickets.Covers(entry, path) {
			return true
		}
	}
	return false
}
