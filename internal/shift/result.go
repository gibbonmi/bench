package shift

import (
	"fmt"
	"io"
	"strconv"

	"github.com/gibbonmi/bench/internal/toon"
)

// Outcome is one of the six truthful shift result states. It is the only thing a Loop
// exit path chooses; the exit code and the shift_result block are both derived from it,
// never picked separately, so the taxonomy stays one source.
type Outcome string

const (
	OutcomeComplete    Outcome = "complete"
	OutcomeFailed      Outcome = "failed"
	OutcomeUsage       Outcome = "usage"
	OutcomeIncomplete  Outcome = "incomplete"
	OutcomeNoOp        Outcome = "no-op"
	OutcomeInterrupted Outcome = "interrupted"
)

// exitCodes is the one map from outcome to process exit code — the taxonomy pinned by
// the spec. Every exit path resolves through Result.ExitCode rather than returning a
// bare int, so this map is the single place the mapping can drift.
var exitCodes = map[Outcome]int{
	OutcomeComplete:    0,
	OutcomeFailed:      1,
	OutcomeUsage:       2,
	OutcomeIncomplete:  3,
	OutcomeNoOp:        4,
	OutcomeInterrupted: 130,
}

// resultFields is the shift_result block's field order, pinned by the spec.
var resultFields = []string{"outcome", "exit", "branch", "committed", "iterations_used", "recovery", "detail"}

// RecoveryNone is the one recovery-pointer sentinel meaning "nothing to preserve" — used
// by the shift_result block, the intent-ledger record, and bench status's rendering, so
// the three surfaces never drift on what "no recovery" looks like.
const RecoveryNone = "none"

// recoveryRefNamespace is the one ref namespace a shift's dirty-tree snapshot lives
// under, mirroring internal/intent/assignment.go's named-constant idiom for ref
// namespaces — every site that builds or reads this prefix uses this constant rather
// than repeating the literal.
const recoveryRefNamespace = "refs/bench/recovery/"

// recoveryRef and recoveryWorktree are the recovery pointer's only two non-"none"
// constructors. Every preserving failure builds its pointer through exactly one of
// these, so the ref-vs-worktree encoding lives in one place.
func recoveryRef(branch string) string    { return "ref:" + recoveryRefNamespace + branch }
func recoveryWorktree(path string) string { return "worktree:" + path }

// Result is the one value computed at every shift exit path: outcome, branch, committed
// count, iterations used, recovery pointer, and a short human-readable detail. Recovery
// is "none" whenever nothing beyond scratch was dirty, and otherwise one of
// recoveryRef/recoveryWorktree's pointer strings.
type Result struct {
	Outcome        Outcome
	Branch         string
	Committed      int
	IterationsUsed int
	Recovery       string
	Detail         string
}

// ExitCode resolves the outcome to its pinned process exit code.
func (r Result) ExitCode() int {
	return exitCodes[r.Outcome]
}

// Emit writes the shift_result TOON block to stdout exactly once. An empty branch or
// recovery reads as "none"; numbers render as decimal strings. toon.Table makes the
// block control-byte-safe by construction, so a hostile detail cannot corrupt the
// surface — it degrades to the AXI error line instead of a crash or a lossy block.
func (r Result) Emit(stdout io.Writer) {
	branch := r.Branch
	if branch == "" {
		branch = "none"
	}
	recovery := r.Recovery
	if recovery == "" {
		recovery = RecoveryNone
	}
	row := []string{
		string(r.Outcome),
		strconv.Itoa(r.ExitCode()),
		branch,
		strconv.Itoa(r.Committed),
		strconv.Itoa(r.IterationsUsed),
		recovery,
		r.Detail,
	}
	block, err := toon.Table("shift_result", resultFields, [][]string{row})
	if err != nil {
		fmt.Fprintln(stdout, toon.RenderError(err))
		return
	}
	fmt.Fprint(stdout, block)
}

// evidenceOutcome splits a post-mutation failure by committed evidence: failed when
// nothing landed this shift, incomplete when at least one green iteration did. Every
// commit, staging, gate, or adapter failure after the first adapter run resolves
// through this one rule — the gate stays the oracle for work, committed count is the
// evidence for progress.
func evidenceOutcome(committed int) Outcome {
	if committed == 0 {
		return OutcomeFailed
	}
	return OutcomeIncomplete
}
