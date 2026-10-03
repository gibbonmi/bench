package gate

import (
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/canonicalpath"
)

// KitValue answers the raw BENCH_KIT value, which is empty when the variable is unset. It
// is the gate's one read of the variable. Each form that takes a kit argument takes this
// value from its caller, so the form itself never reads the process environment.
func KitValue() string {
	return os.Getenv("BENCH_KIT")
}

// KitDir resolves the kit directory an adoption-side reader installs from. Its fallback is
// the executable's parent, then the current directory, where kitRoot in phases.go falls
// back to the graded root. The two are not one derivation: this one answers where the kit
// lives for a caller that may hold no root, and kitRoot answers which kit grades the root
// in hand.
func KitDir() string {
	return kitDirAt(KitValue())
}

// kitDirAt answers kit, or KitDir's fallback when kit is empty.
func kitDirAt(kit string) string {
	if kit != "" {
		return kit
	}
	if exe, err := os.Executable(); err == nil {
		// Development layout: <kit>/dist/bench. Platform package layout:
		// <pkg>/bin/bench, where callers should pass BENCH_KIT from the wrapper.
		return filepath.Clean(filepath.Join(filepath.Dir(exe), ".."))
	}
	return "."
}

// KitSourceCheckout reports whether root is the source tree of the process kit.
func KitSourceCheckout(root string) bool {
	return KitSourceCheckoutAtKit(root, KitValue())
}

// KitSourceCheckoutAtKit reports whether root is the source tree of kit, where an empty
// kit takes KitDir's fallback. The kit repo is where the managed AGENTS.md block and the
// bin/bench.sh launcher are authored, so it never carries the consumer-side copy of
// either, and a row that sent its reader to bench link would name a remedy that breaks
// the shim route and the land route. A consumer repo never satisfies the predicate,
// because its kit resolves to a package or cache directory outside the repository. Both
// paths take their canonical absolute spelling first, so a repository reached by one
// spelling and a kit set to another, or to a relative path, still match. A path with no
// canonical spelling matches nothing.
func KitSourceCheckoutAtKit(root, kit string) bool {
	if root == "" {
		return false
	}
	resolvedRoot, ok := resolvedPath(root)
	if !ok {
		return false
	}
	resolvedKit, ok := resolvedPath(kitDirAt(kit))
	return ok && resolvedRoot == resolvedKit
}

func resolvedPath(path string) (string, bool) {
	resolved, err := canonicalpath.Resolve(path)
	return resolved, err == nil
}

// LaneForCommit resolves the lane a worktree commit at root runs under the process kit.
func LaneForCommit(root string) (*Lane, error) {
	return LaneForCommitAtKit(root, KitValue())
}

// LaneForCommitAtKit resolves the lane a worktree commit at root runs under kit, where an
// empty kit takes the graded root as kitRoot does. It answers a nil lane for a root that
// declares none, so one read tells a caller both whether a lane exists and what it is.
func LaneForCommitAtKit(root, kit string) (*Lane, error) {
	kit = kitRootAt(root, kit)
	checks, err := LaneFor(root, kit)
	if err != nil || checks == nil {
		return nil, err
	}
	// The one built-in-versus-manifest decision. LaneFor answers the built-in lane under
	// this same predicate, so the selection switch and the lane come from one call.
	if sameDirectory(root, kit) {
		return &Lane{Checks: checks, Selective: true}, nil
	}
	return &Lane{Checks: checks, Kit: kit}, nil
}
