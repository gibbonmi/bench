package tickets

import (
	"sort"
	"testing"
)

// TestCommandPackagesAreSorted keeps the command-package list in one readable order with
// one entry per package. An unsorted or duplicated list hides the entry a reader looks for.
func TestCommandPackagesAreSorted(t *testing.T) {
	if len(commandPackages) == 0 || !sort.StringsAreSorted(commandPackages) {
		t.Fatalf("commandPackages = %v, want a nonempty sorted list", commandPackages)
	}
	for i := 1; i < len(commandPackages); i++ {
		if commandPackages[i] == commandPackages[i-1] {
			t.Errorf("commandPackages declares %q twice", commandPackages[i])
		}
	}
}

// TestCommandPackageCoversAtASegmentBoundary grades the command-package predicate: a path
// under a command package answers true, and a path under another package or a sibling that
// shares a command package's leading bytes answers false.
func TestCommandPackageCoversAtASegmentBoundary(t *testing.T) {
	for path, want := range map[string]bool{
		"internal/harnesses/harnesses.go": true,
		"cmd/bench/main.go":               true,
		"internal/spec/read.go":           true,
		"internal/outline/declaration.go": true,
		"internal/toon/toon.go":           false,
		"internal/harnessesx/x.go":        false,
	} {
		if got := CommandPackage(path); got != want {
			t.Errorf("CommandPackage(%q) = %t, want %t", path, got, want)
		}
	}
}
