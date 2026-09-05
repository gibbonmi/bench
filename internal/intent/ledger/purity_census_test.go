package ledger

import (
	"testing"

	"github.com/gibbonmi/bench/internal/puritycensus"
)

// TestPurePackageSourceCensus scans this package's own directory under the leaf
// package policy, which internal/puritycensus owns. The scanned set must hold this
// package's own source, so a census pointed at another directory reds.
// (Coverage row LS13.)
func TestPurePackageSourceCensus(t *testing.T) {
	puritycensus.Scan(t, ".", puritycensus.LeafPackage()).MustHold(t, "ledger.go")
}
