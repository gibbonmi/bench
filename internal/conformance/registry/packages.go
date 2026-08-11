package registry

import "strings"

// ConformancePackage is graded by the filtered second invocation instead of the
// unfiltered list: it is the package whose entry point is this very run.
const ConformancePackage = "internal/conformance"

// ReleaseOnlyPackages are reached only from cmd/bench's dispatch switch on the
// release path. Their suites rebuild the binary and shell out to node, so on a cold
// Go test cache they blow the 600 s package timeout and present as a hung gate; the
// ship tier is where they earn their keep.
var ReleaseOnlyPackages = []string{"internal/releasepreflight", "internal/releaseevidence", "internal/publication"}

// IsExcludedTestPackage reports whether the unfiltered core `go test` leaves a
// package to some other surface. internal/contract is run by the gate's own contract
// phase with the subject root pinned; the conformance package by the filtered
// invocation; the release-only packages by the ship tier.
func IsExcludedTestPackage(pkg string, tier Tier) bool {
	if isContractPackage(pkg) || isPackage(pkg, ConformancePackage) {
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

func isContractPackage(pkg string) bool {
	return isPackage(pkg, "internal/contract") ||
		strings.HasPrefix(pkg, "internal/contract/") ||
		strings.Contains(pkg, "/internal/contract/")
}

// isPackage matches a module-relative package path against a `go list` import path,
// which carries the module prefix in every repo but a bare fixture module.
func isPackage(pkg, rel string) bool {
	return pkg == rel || strings.HasSuffix(pkg, "/"+rel)
}
