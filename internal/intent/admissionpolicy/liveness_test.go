package admissionpolicy

import (
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent/ledger"
)

// shiftEntry is one shift entry whose worktree the facts prove present.
func shiftEntry(outcome, recovery string) (ledger.Ledger, LivenessFacts) {
	entry := ledger.Entry{Key: "shift-1", Kind: ledger.KindShift, CreatedAt: time.Unix(1, 0).UTC(),
		Worktree: "/pool/wt", Branch: "bench/shift-1", Outcome: outcome, Recovery: recovery}
	return ledger.Ledger{Entries: []ledger.Entry{entry}}, LivenessFacts{WorktreeExists: map[string]bool{"/pool/wt": true}}
}

// LE60: a shift that finished with nothing to recover is done, even while its pool
// worktree exists.
func TestLiveDropsAFinishedShiftWithNoRecovery(t *testing.T) {
	current, facts := shiftEntry("complete", ledger.RecoveryNone)
	if live := Live(current, facts); len(live) != 0 {
		t.Fatalf("Live = %v, want the finished entry dropped", live)
	}
}

// LE62: a shift that left its worktree for recovery stays live while the worktree exists,
// even when its branch holds no commit past the default branch and so reads as landed.
func TestLiveKeepsAShiftWithAWorktreeRecovery(t *testing.T) {
	current, facts := shiftEntry("failed", "worktree:/pool/wt")
	facts.Landed = map[string]bool{"bench/shift-1": true}
	if live := Live(current, facts); len(live) != 1 {
		t.Fatalf("Live = %v, want the recovery entry kept", live)
	}
}
