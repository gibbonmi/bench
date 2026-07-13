package runtime

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

func projectionProof(id, state, surface string) actionProof {
	return proof(id, func(t *testing.T) { proveTypedProjection(t, state, surface) })
}

func purityProof(id, state, surface string) actionProof {
	return proof(id, func(t *testing.T) { proveProjectionPurity(t, state, surface) })
}

func presentationProof(id, variant string) actionProof {
	return proof(id, func(t *testing.T) { provePresentation(t, variant) })
}

var ft78Story6Proofs = []actionProof{
	projectionProof("R18/absent-status", "absent", "status"), projectionProof("R18/absent-dashboard", "absent", "dashboard"), projectionProof("R18/absent-roadmap-context", "absent", "roadmap"),
	projectionProof("R18/reusable-green-status", "reusable-green", "status"), projectionProof("R18/reusable-green-dashboard", "reusable-green", "dashboard"), projectionProof("R18/reusable-green-roadmap-context", "reusable-green", "roadmap"),
	projectionProof("R18/red-status", "red", "status"), projectionProof("R18/red-dashboard", "red", "dashboard"), projectionProof("R18/red-roadmap-context", "red", "roadmap"),
	projectionProof("R18/stale-status", "stale", "status"), projectionProof("R18/stale-dashboard", "stale", "dashboard"), projectionProof("R18/stale-roadmap-context", "stale", "roadmap"),
	projectionProof("R18/locked-pending-status", "locked-pending", "status"), projectionProof("R18/locked-pending-dashboard", "locked-pending", "dashboard"), projectionProof("R18/locked-pending-roadmap-context", "locked-pending", "roadmap"),
	projectionProof("R18/interrupted-pending-status", "interrupted-pending", "status"), projectionProof("R18/interrupted-pending-dashboard", "interrupted-pending", "dashboard"), projectionProof("R18/interrupted-pending-roadmap-context", "interrupted-pending", "roadmap"),
	projectionProof("R18/invalid-status", "invalid", "status"), projectionProof("R18/invalid-dashboard", "invalid", "dashboard"), projectionProof("R18/invalid-roadmap-context", "invalid", "roadmap"),
	projectionProof("R18/unavailable-status", "unavailable", "status"), projectionProof("R18/unavailable-dashboard", "unavailable", "dashboard"), projectionProof("R18/unavailable-roadmap-context", "unavailable", "roadmap"),
	purityProof("R19/pending-status", "locked-pending", "status"), purityProof("R19/pending-dashboard", "locked-pending", "dashboard"), purityProof("R19/pending-roadmap-context", "locked-pending", "roadmap"),
	purityProof("R19/invalid-status", "invalid", "status"), purityProof("R19/invalid-dashboard", "invalid", "dashboard"), purityProof("R19/invalid-roadmap-context", "invalid", "roadmap"),
	purityProof("R19/legacy-status", "legacy", "status"), purityProof("R19/legacy-dashboard", "legacy", "dashboard"), purityProof("R19/legacy-roadmap-context", "legacy", "roadmap"),
	presentationProof("R20/absent-signal-severity", "absent"), presentationProof("R20/reusable-green-signal-severity", "reusable-green"),
	presentationProof("R20/red-signal-severity", "red"), presentationProof("R20/stale-signal-severity", "stale"),
	presentationProof("R20/locked-pending-signal-severity", "locked-pending"), presentationProof("R20/interrupted-pending-signal-severity", "interrupted-pending"),
	presentationProof("R20/invalid-signal-severity", "invalid"), presentationProof("R20/unavailable-signal-severity", "unavailable"),
	presentationProof("R20/self-contained-html", "html"), presentationProof("R20/axi-gate-cache-fields", "axi"),
	presentationProof("R20/hostile-bytes-status", "hostile-status"), presentationProof("R20/hostile-bytes-dashboard", "hostile-dashboard"), presentationProof("R20/hostile-bytes-roadmap-context", "hostile-roadmap"),
}

var ft78Story6ExpectedIDs = []string{
	"R18/absent-status", "R18/absent-dashboard", "R18/absent-roadmap-context", "R18/reusable-green-status", "R18/reusable-green-dashboard", "R18/reusable-green-roadmap-context", "R18/red-status", "R18/red-dashboard", "R18/red-roadmap-context", "R18/stale-status", "R18/stale-dashboard", "R18/stale-roadmap-context", "R18/locked-pending-status", "R18/locked-pending-dashboard", "R18/locked-pending-roadmap-context", "R18/interrupted-pending-status", "R18/interrupted-pending-dashboard", "R18/interrupted-pending-roadmap-context", "R18/invalid-status", "R18/invalid-dashboard", "R18/invalid-roadmap-context", "R18/unavailable-status", "R18/unavailable-dashboard", "R18/unavailable-roadmap-context",
	"R19/pending-status", "R19/pending-dashboard", "R19/pending-roadmap-context", "R19/invalid-status", "R19/invalid-dashboard", "R19/invalid-roadmap-context", "R19/legacy-status", "R19/legacy-dashboard", "R19/legacy-roadmap-context",
	"R20/absent-signal-severity", "R20/reusable-green-signal-severity", "R20/red-signal-severity", "R20/stale-signal-severity", "R20/locked-pending-signal-severity", "R20/interrupted-pending-signal-severity", "R20/invalid-signal-severity", "R20/unavailable-signal-severity", "R20/self-contained-html", "R20/axi-gate-cache-fields", "R20/hostile-bytes-status", "R20/hostile-bytes-dashboard", "R20/hostile-bytes-roadmap-context",
}

func TestFT78Story6ProofLedgerCompleteness(t *testing.T) {
	contract.NoteContractFailure(t, ft78LedgerFailure)
	seen := map[string]int{}
	for _, c := range ft78Story6Proofs {
		seen[c.id]++
		if c.driver == nil {
			t.Fatalf("%s: nil real driver", c.id)
		}
	}
	if len(seen) != len(ft78Story6ExpectedIDs) {
		t.Fatalf("%s: registered IDs = %d, want %d", ft78LedgerFailure, len(seen), len(ft78Story6ExpectedIDs))
	}
	for _, id := range ft78Story6ExpectedIDs {
		if seen[id] != 1 {
			t.Fatalf("%s: %s registrations = %d, want 1", ft78LedgerFailure, id, seen[id])
		}
	}
}

func TestFT78Story6ProofLedger(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	for _, c := range ft78Story6Proofs {
		c := c
		t.Run(c.id, c.driver)
	}
}
