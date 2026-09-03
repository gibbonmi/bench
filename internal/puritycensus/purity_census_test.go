package puritycensus

import "testing"

// TestPurePackageSourceCensus runs the census over this package's own directory under
// the leaf policy, so the package that grades every other pure owner is not itself a
// blind spot. (Coverage row LQ16.)
func TestPurePackageSourceCensus(t *testing.T) {
	Scan(t, ".", LeafPackage()).MustHold(t, "census.go")
}
