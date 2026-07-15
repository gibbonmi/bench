package runtime

import (
	"testing"

	"github.com/gibbonmi/bench/internal/contract"
)

const ft78LedgerFailure = "FT78 proof ledger completeness contract failed"

type actionProof struct {
	id     string
	driver func(*testing.T)
}

func proof(id string, driver func(*testing.T)) actionProof { return actionProof{id, driver} }
func commitProof(id, variant string) actionProof {
	return proof(id, func(t *testing.T) { proveCommitResult(t, variant) })
}
func shiftProof(id, variant string) actionProof {
	return proof(id, func(t *testing.T) { proveShiftResult(t, variant) })
}
func stopProof(id, variant string) actionProof {
	return proof(id, func(t *testing.T) { proveStopResult(t, variant) })
}
func r17Proof(id, variant string) actionProof {
	return proof(id, func(t *testing.T) { proveR17Result(t, variant) })
}

var ft78Story5Proofs = []actionProof{
	commitProof("R14/exact-green-reuse", "reuse"), commitProof("R14/stale-tree-rerun", "stale"),
	commitProof("R14/ordinary-red", "red"), commitProof("R14/oracle-mismatch", "oracle-mismatch"),
	commitProof("R14/absent-inspection", "absent"), commitProof("R14/ready-red-inspection", "ready-red"),
	commitProof("R14/stale-green-inspection", "stale-green"), commitProof("R14/open-subject-green-inspection", "open-green"),
	commitProof("R14/locked-pending-inspection", "locked-pending"), commitProof("R14/interrupted-pending-inspection", "interrupted"),
	commitProof("R14/invalid-inspection", "invalid"), commitProof("R14/unavailable-inspection", "unavailable"),
	commitProof("R14/lock-open-result", "lock-open"), commitProof("R14/lock-acquire-result", "lock-acquire"),
	commitProof("R14/pending-persistence-result", "pending-persistence"), commitProof("R14/final-persistence-result", "final-persistence"),
	commitProof("R14/subject-build-result", "subject-build"), commitProof("R14/subject-recheck-drift-result", "drift"),
	commitProof("R14/cancellation-result", "cancellation"), commitProof("R14/start-failure-result", "start-failure"),
	commitProof("R14/no-gate-result", "no-gate"),
	shiftProof("R15/ordinary-red", "red"), shiftProof("R15/lock", "lock"),
	shiftProof("R15/persistence", "persistence"), shiftProof("R15/drift", "drift"),
	shiftProof("R15/cancellation", "cancellation"),
	stopProof("R16/green", "green"), stopProof("R16/red", "red"),
	stopProof("R16/unarmed", "unarmed"), stopProof("R16/active", "active"),
	stopProof("R16/lock-open", "lock-open"), stopProof("R16/lock-acquire", "lock-acquire"),
	stopProof("R16/pending-persistence", "pending-persistence"), stopProof("R16/final-persistence", "final-persistence"),
	stopProof("R16/subject-build", "subject-build"), stopProof("R16/subject-recheck-drift", "drift"),
	stopProof("R16/cancellation", "cancellation"), stopProof("R16/start-failure", "start-failure"),
	stopProof("R16/no-gate", "no-gate"),
	r17Proof("R17/drift", "drift"), r17Proof("R17/cancellation", "cancellation"),
	r17Proof("R17/invalid-zero-byte", "invalid-zero-byte"), r17Proof("R17/invalid-malformed", "invalid-malformed"),
	r17Proof("R17/invalid-legacy", "invalid-legacy"), r17Proof("R17/invalid-oversized", "invalid-oversized"),
	r17Proof("R17/invalid-wrong-mode", "invalid-wrong-mode"), r17Proof("R17/invalid-symlink", "invalid-symlink"),
	r17Proof("R17/invalid-directory", "invalid-directory"),
	r17Proof("R17/fault-lock-open", "fault-lock-open"), r17Proof("R17/fault-lock-acquisition", "fault-lock-acquisition"),
	r17Proof("R17/fault-temporary-create", "fault-temporary-create"), r17Proof("R17/fault-mode-establishment", "fault-mode-establishment"),
	r17Proof("R17/fault-write", "fault-write"), r17Proof("R17/fault-file-sync", "fault-file-sync"),
	r17Proof("R17/fault-file-close", "fault-file-close"), r17Proof("R17/fault-atomic-rename", "fault-atomic-rename"),
	r17Proof("R17/fault-directory-open", "fault-directory-open"), r17Proof("R17/fault-directory-sync", "fault-directory-sync"),
	r17Proof("R17/fault-directory-close", "fault-directory-close"), r17Proof("R17/fault-post-run-subject-rebuild", "fault-post-run-subject-rebuild"),
}

var ft78Story5ExpectedIDs = []string{
	"R14/exact-green-reuse", "R14/stale-tree-rerun", "R14/ordinary-red", "R14/oracle-mismatch", "R14/absent-inspection", "R14/ready-red-inspection", "R14/stale-green-inspection", "R14/open-subject-green-inspection", "R14/locked-pending-inspection", "R14/interrupted-pending-inspection", "R14/invalid-inspection", "R14/unavailable-inspection", "R14/lock-open-result", "R14/lock-acquire-result", "R14/pending-persistence-result", "R14/final-persistence-result", "R14/subject-build-result", "R14/subject-recheck-drift-result", "R14/cancellation-result", "R14/start-failure-result", "R14/no-gate-result",
	"R15/ordinary-red", "R15/lock", "R15/persistence", "R15/drift", "R15/cancellation",
	"R16/green", "R16/red", "R16/unarmed", "R16/active", "R16/lock-open", "R16/lock-acquire", "R16/pending-persistence", "R16/final-persistence", "R16/subject-build", "R16/subject-recheck-drift", "R16/cancellation", "R16/start-failure", "R16/no-gate",
	"R17/drift", "R17/cancellation", "R17/invalid-zero-byte", "R17/invalid-malformed", "R17/invalid-legacy", "R17/invalid-oversized", "R17/invalid-wrong-mode", "R17/invalid-symlink", "R17/invalid-directory",
	"R17/fault-lock-open", "R17/fault-lock-acquisition", "R17/fault-temporary-create", "R17/fault-mode-establishment", "R17/fault-write", "R17/fault-file-sync", "R17/fault-file-close", "R17/fault-atomic-rename", "R17/fault-directory-open", "R17/fault-directory-sync", "R17/fault-directory-close", "R17/fault-post-run-subject-rebuild",
}

func TestFT78Story5ProofLedgerCompleteness(t *testing.T) {
	contract.NoteContractFailure(t, ft78LedgerFailure)
	seen := map[string]int{}
	for _, c := range ft78Story5Proofs {
		seen[c.id]++
		if c.driver == nil {
			t.Fatalf("%s: nil real driver", c.id)
		}
	}
	if len(seen) != len(ft78Story5ExpectedIDs) {
		t.Fatalf("%s: registered IDs = %d, want %d", ft78LedgerFailure, len(seen), len(ft78Story5ExpectedIDs))
	}
	for _, id := range ft78Story5ExpectedIDs {
		if seen[id] != 1 {
			t.Fatalf("%s: %s registrations = %d, want 1", ft78LedgerFailure, id, seen[id])
		}
	}
}

func TestFT78Story5ProofLedger(t *testing.T) {
	contract.SkipIfSubjectBenchMissing(t)
	r17BridgeBinary = buildR17Bridge(t)
	runProofLedger(t, ft78Story5Proofs, serialProof)
}
