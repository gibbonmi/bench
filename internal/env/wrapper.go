package env

import "strings"

// WrapperRouting names the variables bin/bench.sh exports before it execs the Bench
// binary, so the binary can find the kit its wrapper belongs to. They describe the
// wrapper that started this process, not the tree a child runs in. A child launched
// against a different tree has to resolve its own kit, and an inherited value points
// it silently back at the caller's.
var WrapperRouting = []string{"BENCH_KIT", "BENCH_WRAPPER"}

// WithoutWrapperRouting returns base with every WrapperRouting assignment removed,
// along with any extra names the caller's own child contract adds. Every other
// assignment survives in order.
func WithoutWrapperRouting(base []string, extra ...string) []string {
	filtered := make([]string, 0, len(base))
	for _, entry := range base {
		if assignsAny(entry, WrapperRouting) || assignsAny(entry, extra) {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func assignsAny(entry string, names []string) bool {
	for _, name := range names {
		if strings.HasPrefix(entry, name+"=") {
			return true
		}
	}
	return false
}
