package otelrecord

import "sync/atomic"

// This file declares the whole attribute set a Bench span may carry, the resource block
// every line carries, and the outcome vocabulary those spans state. The declared set is
// the redaction contract: the encoder drops each span attribute whose key is not declared
// here, so a key ships only when the diff that sets it also declares it. A review still
// grades the value source of each declared key.
//
// No attribute carries payload. A subject is identified by its digest — a tree id or a
// commit id — and never by its subject text.

// The resource block. The encoder writes these keys itself for every line and never the
// SDK resource, so no detector and no OTEL_* environment value reaches the record.
const (
	// ResourceServiceName names the service that wrote the line; its value is always bench.
	ResourceServiceName = "service.name"

	// ResourceServiceVersion carries the Bench version that wrote the line. A release
	// producer selects its own events by this value.
	ResourceServiceVersion = "service.version"

	// ResourceRecordSchema carries the record schema version of the line.
	ResourceRecordSchema = "bench.record.schema"

	// ServiceName is the one ResourceServiceName value.
	ServiceName = "bench"

	// RecordSchema is the schema version this package writes and the one it reads. A
	// line with another schema value is malformed; a line with none is a legacy line.
	RecordSchema = "1"
)

// recordVersion is the Bench version the command layer set at process start. A process
// that sets none writes no ResourceServiceVersion key.
var recordVersion atomic.Pointer[string]

// SetVersion hands the process's stamped Bench version to the record. The command layer
// calls it once at process start, and every later line names that version.
func SetVersion(version string) {
	recordVersion.Store(&version)
}

const (
	// AttrSeam names the seam that started the span, such as "gate" or "worktree.land".
	AttrSeam = "bench.seam"

	// AttrSubjectID carries the subject digest: a composed tree id or a commit id.
	AttrSubjectID = "bench.subject.id"

	// AttrAssignmentID identifies the worktree assignment independently of a published subject.
	AttrAssignmentID = "bench.assignment.id"

	// AttrOutcome carries the seam's exit, such as "green", "red", or "skipped".
	AttrOutcome = "bench.outcome"

	// AttrOutcomeCheck names the first failing check when the outcome is red.
	AttrOutcomeCheck = "bench.outcome.check"

	// AttrOutcomeDiagnostic carries the first diagnostic line of the failing check.
	AttrOutcomeDiagnostic = "bench.outcome.diagnostic"

	// AttrOutcomeBlocker names the need that blocked a skipped phase.
	AttrOutcomeBlocker = "bench.outcome.blocker"

	// AttrMeasurePathCount carries the composed path count at a publication seam.
	AttrMeasurePathCount = "bench.measure.path_count"

	// AttrMeasureCensusRawCalls carries the census raw-call count at the landing seam.
	AttrMeasureCensusRawCalls = "bench.measure.census_raw_calls"

	// AttrRecord marks the record line itself rather than the seam.
	AttrRecord = "bench.record"

	// AttrIntentKey carries the intent entry key of the run that opened the span.
	AttrIntentKey = "bench.intent.key"

	// AttrAgent carries the base name of the adapter a shift runs, never its path.
	AttrAgent = "bench.agent"

	// AttrShiftCap carries the iteration cap of a shift.
	AttrShiftCap = "bench.shift.cap"

	// AttrLineTier carries the declared line tier, and only a tier the tier list names.
	AttrLineTier = "bench.line.tier"

	// AttrShiftOutcome carries the shift's own outcome word, such as "complete" or "no-op".
	AttrShiftOutcome = "bench.shift.outcome"

	// AttrWorkState carries the work-state word of the run.
	AttrWorkState = "bench.work.state"

	// AttrCleanup carries the cleanup word of the run's worktree.
	AttrCleanup = "bench.cleanup"

	// AttrRecoveryKind names the kind of recovery pointer a run left.
	AttrRecoveryKind = "bench.recovery.kind"

	// AttrRecoveryKey carries the base name of the recovery pointer, never its path.
	AttrRecoveryKey = "bench.recovery.key"

	// AttrAdapterResult carries how a pass's adapter run ended.
	AttrAdapterResult = "bench.adapter.result"

	// AttrAdapterExit carries the adapter's exit code when the adapter exited.
	AttrAdapterExit = "bench.adapter.exit"

	// AttrLineHarness carries the harness a line resolution named, and only a known one.
	AttrLineHarness = "bench.line.harness"

	// AttrLineModel carries the resolved model, and only a safe model token.
	AttrLineModel = "bench.line.model"

	// AttrMemoryState carries what became of a shift's notes.
	AttrMemoryState = "bench.memory.state"

	// AttrMemoryBytes carries the byte count of retained notes.
	AttrMemoryBytes = "bench.memory.bytes"

	// AttrMemoryDigest carries the SHA-256 digest of the retained memory file, which
	// references the notes without carrying their text.
	AttrMemoryDigest = "bench.memory.digest"
)

// RecordStart is the AttrRecord value on the line written at span start. The start line
// carries no end time, so a consumer filters unfinished spans by this marker.
const RecordStart = "start"

// DeclaredAttributes is the complete declared set. The encoder writes only these span
// attributes, and a later ticket that adds a key adds it here in the same diff.
var DeclaredAttributes = []string{
	AttrSeam,
	AttrSubjectID,
	AttrAssignmentID,
	AttrOutcome,
	AttrOutcomeCheck,
	AttrOutcomeDiagnostic,
	AttrOutcomeBlocker,
	AttrMeasurePathCount,
	AttrMeasureCensusRawCalls,
	AttrRecord,
	AttrIntentKey,
	AttrAgent,
	AttrShiftCap,
	AttrLineTier,
	AttrShiftOutcome,
	AttrWorkState,
	AttrCleanup,
	AttrRecoveryKind,
	AttrRecoveryKey,
	AttrAdapterResult,
	AttrAdapterExit,
	AttrLineHarness,
	AttrLineModel,
	AttrMemoryState,
	AttrMemoryBytes,
	AttrMemoryDigest,
}

// declared answers whether the encoder may write a span attribute key.
var declared = func() map[string]bool {
	keys := make(map[string]bool, len(DeclaredAttributes))
	for _, key := range DeclaredAttributes {
		keys[key] = true
	}
	return keys
}()

// The outcome vocabulary. A consumer groups runs by these three words, so a seam states
// one of them and never its own spelling of the same idea.
const (
	OutcomeGreen   = "green"
	OutcomeRed     = "red"
	OutcomeSkipped = "skipped"
)

// The work-state vocabulary. A consumer separates finished work from failed and cut-off
// work by these words, whatever each seam calls its own outcomes.
const (
	WorkCompleted   = "completed"
	WorkFailed      = "failed"
	WorkInterrupted = "interrupted"
)

// The cleanup vocabulary: the run released its worktree, retained it for recovery, or
// held none.
const (
	CleanupReleased = "released"
	CleanupRetained = "retained"
	CleanupNone     = "none"
)

// The memory-state vocabulary: the notes were kept, there were none, their file was a
// link, a special file, or too large to keep, or the store write failed.
const (
	MemoryRetained = "retained"
	MemoryAbsent   = "absent"
	MemoryRefused  = "refused"
	MemoryFailed   = "failed"
)

// The adapter-result vocabulary: the adapter process exited, with a code, or it never
// started.
const (
	AdapterExited      = "exited"
	AdapterSpawnFailed = "spawn-failed"
)

// ExitOutcome is the outcome for a seam whose zero alone is green. A guard, a gate, and
// a worktree child all speak this policy: any nonzero exit the operator saw reads red.
func ExitOutcome(exit int) string {
	if exit == 0 {
		return OutcomeGreen
	}
	return OutcomeRed
}

// PublishedExitOutcome is the outcome for a publication seam, where exit 3 published its
// commit and named its own remaining step on the verb's output. The record reads that
// publication as green, and the remainder stays on the verb's output rather than here.
func PublishedExitOutcome(exit int) string {
	if exit == 3 {
		return OutcomeGreen
	}
	return ExitOutcome(exit)
}
