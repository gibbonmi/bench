package registry

import "strings"

// ConformancePackage is graded by the filtered second invocation instead of the
// unfiltered list. It is the package whose entry point is this very run.
const ConformancePackage = "internal/conformance"

// ReleaseOnlyPackages are reached only from cmd/bench's dispatch switch on the release
// path. Their suites rebuild the binary and shell out to node. On a cold Go test cache,
// they exceed the 600 s package timeout and present as a hung gate. The ship tier is
// where they earn their keep.
var ReleaseOnlyPackages = []string{"internal/releasepreflight", "internal/releaseevidence", "internal/publication"}

// IsExcludedTestPackage reports whether the unfiltered core `go test` leaves a package to
// some other surface. The filtered invocation runs the conformance package. The ship
// tier runs the release-only packages.
func IsExcludedTestPackage(pkg string, tier Tier) bool {
	if isPackage(pkg, ConformancePackage) {
		return true
	}
	if tier == Ship {
		return false
	}
	for _, releaseOnly := range ReleaseOnlyPackages {
		if isPackage(pkg, releaseOnly) {
			return true
		}
	}
	return false
}

// isPackage matches a module-relative package path against a `go list` import path. The
// import path carries the module prefix in every repo except a bare fixture module.
func isPackage(pkg, rel string) bool {
	return pkg == rel || strings.HasSuffix(pkg, "/"+rel)
}
