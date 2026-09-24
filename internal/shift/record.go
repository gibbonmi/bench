package shift

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/lines"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// The shift's seams: the shift itself, each main iteration, and each refactor pass.
const (
	shiftSeam     = "shift"
	iterationSeam = "shift.iteration"
	refactorSeam  = "shift.refactor"
)

// shiftRecord is the shift's span in the repository's record, with the pass span open
// below it. The shift span starts right after the intent entry persists and ends in
// finish, so every exit path ends it. No value it writes is a path or operator text: the
// adapter and the recovery pointer reach the record as base names, and the line tier only
// as a word the tier list names.
type shiftRecord struct {
	ctx     context.Context // carries the tracer and the shift span
	span    trace.Span
	pass    trace.Span // the open pass span, or nil between passes
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
	ctx, span, end := otelrecord.BeginIn(context.Background(), "", root, shiftSeam, shiftSeam, attrs...)
	return &shiftRecord{ctx: ctx, span: span, end: end, root: root}
}

// beginPass ends the open pass and starts a pass span of seam under the shift span. The
// returned context parents the pass's gate run, so the gate span is the pass's child. A
// nil record returns a context with no span.
func (r *shiftRecord) beginPass(seam string) context.Context {
	if r == nil {
		return context.Background()
	}
	r.endPass()
	ctx, pass := otelrecord.TracerFrom(r.ctx).Start(r.ctx, seam, trace.WithAttributes(attribute.String(otelrecord.AttrSeam, seam)))
	r.pass = pass
	return ctx
}

// adapterRan records how the pass's adapter run ended: an exit with its code, or a
// failure to spawn.
func (r *shiftRecord) adapterRan(err error) {
	if r == nil || r.pass == nil {
		return
	}
	var exitErr *exec.ExitError
	switch {
	case err == nil:
		r.pass.SetAttributes(attribute.String(otelrecord.AttrAdapterResult, otelrecord.AdapterExited), attribute.Int(otelrecord.AttrAdapterExit, 0))
	case errors.As(err, &exitErr):
		r.pass.SetAttributes(attribute.String(otelrecord.AttrAdapterResult, otelrecord.AdapterExited), attribute.Int(otelrecord.AttrAdapterExit, exitErr.ExitCode()))
	default:
		r.pass.SetAttributes(attribute.String(otelrecord.AttrAdapterResult, otelrecord.AdapterSpawnFailed))
	}
}

// passCommitted records the commit the pass made in the worktree at wt.
func (r *shiftRecord) passCommitted(wt string) {
	if r == nil || r.pass == nil {
		return
	}
	if head, err := git.Output("-C", wt, "rev-parse", "HEAD"); err == nil {
		r.pass.SetAttributes(attribute.String(otelrecord.AttrSubjectID, head))
	}
}

// endPass ends the open pass span, if one is open.
func (r *shiftRecord) endPass() {
	if r == nil || r.pass == nil {
		return
	}
	r.pass.End()
	r.pass = nil
}

// finish ends the open pass and then the shift span, with the shift's result: its
// outcome, work state, exit outcome, head commit, recovery pointer, and cleanup. A nil
// record, the usage exit before the intent entry persists, records nothing.
func (r *shiftRecord) finish(res Result) {
	if r == nil {
		return
	}
	r.endPass()
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
