//go:build stress

package runtime

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

// TestManifestEntryLimitFullScale exercises the real 100_000-entry boundary through the
// built gate. It is heavy (~13s) and gold-plating on every commit: the every-gate ledger
// (TestFT78Story2ProofLedger) proves the boundary logic at tiny scale, and
// TestManifestEntryLimitConstant pins the shipped ceiling. So the full-scale walk runs
// only under `go test -tags stress` (CI), never the gate.
func TestManifestEntryLimitFullScale(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	t.Run("at-limit", func(t *testing.T) { manifestCollectorLimitProof(t, 100_000, 100_000, "", true) })
	t.Run("over-limit", func(t *testing.T) {
		manifestCollectorLimitProof(t, 100_000, 100_001, "declared path unavailable", false)
	})
}
