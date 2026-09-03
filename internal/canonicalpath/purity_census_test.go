package canonicalpath

import (
	"testing"

	"github.com/gibbonmi/bench/internal/puritycensus"
)

// TestPurePackageSourceCensus scans this package's own directory under the leaf policy,
// which internal/puritycensus owns. The scanned set must hold this package's own source,
// so a census pointed at another directory reds. (Coverage rows LQ3, LQ12, LQ14.)
func TestPurePackageSourceCensus(t *testing.T) {
	puritycensus.Scan(t, ".", puritycensus.LeafPackage()).MustHold(t, "canonicalpath.go")
}
