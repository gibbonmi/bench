package landingpolicy

import (
	"testing"

	"github.com/gibbonmi/bench/internal/puritycensus"
)

// TestPurePackageSourceCensus enforces the pure-package boundary over this package s
// own directory: no import of the parent package, the Git owner, the bounds
// constants, process execution, or the system-call surface, and no ambient
// environment, directory, clock, or descendant effect. internal/puritycensus owns the
// policy. The scanned set must hold this package's own source, so a census pointed at
// another directory reds. (Coverage rows LQ12, LQ14.)
func TestPurePackageSourceCensus(t *testing.T) {
	puritycensus.Scan(t, ".", puritycensus.PolicyPackage()).MustHold(t, "landingpolicy.go")
}
