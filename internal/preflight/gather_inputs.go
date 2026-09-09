package preflight

import (
	"github.com/gibbonmi/bench/internal/diff"
	"github.com/gibbonmi/bench/internal/git"
	specref "github.com/gibbonmi/bench/internal/spec"
)

// Gather gathers preflight facts for root, mode, slug, and the first optional
// explicit base. It supplies no source-tip pin. It returns one immutable
// Facts snapshot for Decide or one BootstrapFailure; it never classifies a
// check. Exactly one result is non-zero.
func Gather(root, mode, slug string, explicitBase ...string) (Facts, *BootstrapFailure) {
	base := ""
	if len(explicitBase) > 0 {
		base = explicitBase[0]
	}
	return GatherPinned(root, mode, slug, base, "")
}

// GatherPinned gathers preflight facts for root, mode, slug, an explicit
// base, and a source-tip pin. With an explicit base, it reads one
// movement-checked source snapshot and retries one snapshot drift. It returns
// snapshot drift after a second movement, or one BootstrapFailure for a read
// or pin failure. A resolved pin that names the wrong commit is a Decide row.
func GatherPinned(root, mode, slug, explicitBase, sourceTipPin string) (Facts, *BootstrapFailure) {
	return gatherPinned(root, mode, slug, explicitBase, sourceTipPin, false)
}

func gatherCharge(root, mode, slug string, source *diff.SourceRange, sourcePaths []string, sourceTipPin string) (Facts, *BootstrapFailure) {
	return gather(root, mode, slug, source, sourcePaths, sourceTipPin, true)
}

func gatherPinned(root, mode, slug, explicitBase, sourceTipPin string, noFollow bool) (Facts, *BootstrapFailure) {
	if explicitBase != "" {
		var gathered Facts
		var gatherFailure *BootstrapFailure
		result := diff.MovementCheckedRetry(root, func(snapshot diff.MovementSnapshot) (string, string) {
			var err error
			var resolveKind, resolveHint string
			source, resolveKind, resolveHint := snapshot.ResolveSourceRange(explicitBase)
			if resolveKind != "" {
				return resolveKind, resolveHint
			}
			paths, err := snapshot.SourceSnapshotPaths(source)
			if err != nil {
				return "changed files not readable", err.Error()
			}
			if mode == "review" {
				dirty, statusErr := git.Output("-C", root, "status", "--porcelain")
				if statusErr != nil {
					return "source status unreadable", statusErr.Error()
				}
				if dirty != "" {
					return "source not clean", "review source has uncommitted changes"
				}
			}
			gathered, gatherFailure = gather(root, mode, slug, &source, paths, sourceTipPin, noFollow)
			if gatherFailure != nil {
				return gatherFailure.Kind, gatherFailure.Hint
			}
			return "", ""
		})
		if result.Kind != "" {
			if gatherFailure != nil {
				return Facts{}, gatherFailure
			}
			return Facts{}, &BootstrapFailure{result.Kind, result.Hint}
		}
		if result.DriftKind != "" {
			return Facts{}, &BootstrapFailure{"snapshot drift", result.DriftHint}
		}
		return gathered, nil
	}
	return gather(root, mode, slug, nil, nil, sourceTipPin, noFollow)
}

func resolveSpecInput(root, slug string, noFollow bool) ([]byte, string, []string, bool, error) {
	if noFollow {
		return specref.ResolveNoFollow(root, slug)
	}
	return specref.Resolve(root, slug)
}

func gatherTickets(root, dir, mode, tag string) (ticketFacts, *BootstrapFailure) {
	return gatherTicketsWithPolicy(root, dir, mode, tag, false)
}
