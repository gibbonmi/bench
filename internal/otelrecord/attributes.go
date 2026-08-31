package otelrecord

// This file declares the whole attribute set a Bench span may carry, and the outcome
// vocabulary those spans state. Story 19 limits a span to the seam name, the subject id,
// the outcome, and the measures. No mechanical check enforces the limit, so a review
// grades each new attribute against this declaration. A key that is not declared here
// does not ship.
//
// No attribute carries payload. A subject is identified by its digest — a tree id or a
// commit id — and never by its subject text.

const (
	// AttrSeam names the seam that started the span, such as "gate" or "worktree.land".
	AttrSeam = "bench.seam"

	// AttrSubjectID carries the subject digest: a composed tree id or a commit id.
	AttrSubjectID = "bench.subject.id"

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
)

// RecordStart is the AttrRecord value on the line written at span start. The start line
// carries no end time, so a consumer filters unfinished spans by this marker.
const RecordStart = "start"

// DeclaredAttributes is the complete declared set. A reviewer grades a new span attribute
// against this list, and a later ticket that adds a key adds it here in the same diff.
var DeclaredAttributes = []string{
	AttrSeam,
	AttrSubjectID,
	AttrOutcome,
	AttrOutcomeCheck,
	AttrOutcomeDiagnostic,
	AttrOutcomeBlocker,
	AttrMeasurePathCount,
	AttrMeasureCensusRawCalls,
	AttrRecord,
}

// The outcome vocabulary. A consumer groups runs by these three words, so a seam states
// one of them and never its own spelling of the same idea.
const (
	OutcomeGreen   = "green"
	OutcomeRed     = "red"
	OutcomeSkipped = "skipped"
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
