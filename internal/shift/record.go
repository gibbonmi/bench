package shift

import (
	"context"
	"os"
	"path/filepath"
	"slices"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// shiftRecord is the shift's one span in the repository's record. It starts right after
// the intent entry persists and ends in finish, so every exit path ends it. No value it
// writes is a path or operator text: the adapter and the recovery pointer reach the record
// as base names, and the line tier only as a word the tier list names.
type shiftRecord struct {
	span    trace.Span
	end     func()
	root    string
	session *session // nil until the worktree is acquired
}

// workStates maps each shift outcome to its work state. A shift that stopped on its own
// terms finished its work, even with no change or at its cap.
var workStates = map[Outcome]string{
	OutcomeComplete:    otelrecord.WorkCompleted,
	OutcomeNoOp:        otelrecord.WorkCompleted,
	OutcomeIncomplete:  otelrecord.WorkCompleted,
	OutcomeFailed:      otelrecord.WorkFailed,
	OutcomeUsage:       otelrecord.WorkFailed,
	OutcomeInterrupted: otelrecord.WorkInterrupted,
}

// beginShiftRecord starts the shift span in root's record with the shift's intent key,
// its adapter, its cap, and its declared tier.
func beginShiftRecord(root string, entry intent.Entry, maxIters int) *shiftRecord {
	attrs := []attribute.KeyValue{
		attribute.String(otelrecord.AttrIntentKey, entry.Key),
		attribute.String(otelrecord.AttrAgent, filepath.Base(os.Getenv("BENCH_AGENT"))),
		attribute.Int(otelrecord.AttrShiftCap, maxIters),
	}
	if tier := os.Getenv("BENCH_MODEL"); slices.Contains(lines.Tiers, tier) {
		attrs = append(attrs, attribute.String(otelrecord.AttrLineTier, tier))
	}
	_, span, end := otelrecord.BeginIn(context.Background(), "", root, "shift", "shift", attrs...)
	return &shiftRecord{span: span, end: end, root: root}
}

// finish ends the span with the shift's result: its outcome, work state, exit outcome,
// head commit, recovery pointer, and cleanup. A nil record, the usage exit before the
// intent entry persists, records nothing.
func (r *shiftRecord) finish(res Result) {
	if r == nil {
		return
	}
	kind, path := splitRecovery(res.Recovery)
	attrs := []attribute.KeyValue{
		attribute.String(otelrecord.AttrShiftOutcome, string(res.Outcome)),
		attribute.String(otelrecord.AttrWorkState, workStates[res.Outcome]),
		attribute.String(otelrecord.AttrOutcome, otelrecord.ExitOutcome(res.ExitCode())),
		attribute.String(otelrecord.AttrCleanup, r.cleanup(kind)),
		attribute.String(otelrecord.AttrRecoveryKind, kind),
	}
	if path != "" {
		attrs = append(attrs, attribute.String(otelrecord.AttrRecoveryKey, filepath.Base(path)))
	}
	if res.Committed > 0 {
		if head, err := git.Output("-C", r.root, "rev-parse", "refs/heads/"+res.Branch+"^{commit}"); err == nil {
			attrs = append(attrs, attribute.String(otelrecord.AttrSubjectID, head))
		}
	}
	r.span.SetAttributes(attrs...)
	r.end()
}

// cleanup names what became of the shift's worktree, given the kind of its recovery
// pointer: retained for recovery, released to the pool, or none when the shift never
// held one or its release failed.
func (r *shiftRecord) cleanup(recoveryKind string) string {
	switch {
	case recoveryKind == recoveryWorktreeKind:
		return otelrecord.CleanupRetained
	case r.session != nil && r.session.released:
		return otelrecord.CleanupReleased
	}
	return otelrecord.CleanupNone
}
