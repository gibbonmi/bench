package gate

import (
	"os"
	"path/filepath"
)

// KitDir resolves the kit directory an adoption-side reader installs from. Its fallback is
// the executable's parent, then the current directory, where kitRoot in phases.go falls
// back to the graded root. The two are not one derivation: this one answers where the kit
// lives for a caller that may hold no root, and kitRoot answers which kit grades the root
// in hand.
func KitDir() string {
	if kit := os.Getenv("BENCH_KIT"); kit != "" {
		return kit
	}
	if exe, err := os.Executable(); err == nil {
		// Development layout: <kit>/dist/bench. Platform package layout:
		// <pkg>/bin/bench, where callers should pass BENCH_KIT from the wrapper.
		return filepath.Clean(filepath.Join(filepath.Dir(exe), ".."))
	}
	return "."
}

// KitSourceCheckout reports whether root is the kit's own source tree. The kit repo is
// where the managed AGENTS.md block and the bin/bench.sh launcher are authored, so it
// never carries the consumer-side copy of either, and a row that sent its reader to bench
// link would name a remedy that breaks the shim route and the land route. A consumer repo
// never satisfies the predicate, because its kit resolves to a package or cache directory
// outside the repository. Both paths resolve through symlinks first, so a repository
// reached by one spelling and a BENCH_KIT set to another still match.
func KitSourceCheckout(root string) bool {
	kit := KitDir()
	if root == "" || kit == "" {
		return false
	}
	return resolvedPath(root) == resolvedPath(kit)
}

func resolvedPath(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}
	return filepath.Clean(path)
}
