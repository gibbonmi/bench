package reviewrecord

import "errors"

// The record forms. Each form writes the entry that a refusal naming it lacks, and the
// record command declares each form under its name here.
const (
	FormAssignment   = "assignment"
	FormAmendment    = "amendment"
	FormChunk        = "chunk"
	FormVerification = "verification"
	FormReview       = "review"
	FormConfirmation = "confirmation"
	FormFold         = "fold"
	FormCompletion   = "completion"
)

// Repair is a refusal that one record form repairs. Form names that form, and Err is the
// refusal itself: Repair wraps it, so errors.Is and errors.As still reach its identity.
// The other fields hold the values that the refusing owner knows, and an empty field is
// unknown. Chunk is a planned chunk ID, and Fold is a recorded fold ID. Requirement is a
// planned requirement ID, Ticket a planned ticket basename, and Axis a review axis. Source is
// the full commit ID that the form reads. A verification repair with no Chunk names the
// final list.
type Repair struct {
	Form                                           string
	Chunk, Fold, Requirement, Ticket, Axis, Source string
	Err                                            error
}

func (r Repair) Error() string { return r.Err.Error() }
func (r Repair) Unwrap() error { return r.Err }

// Step is one record write in completion order. Final selects the final verification
// list.
type Step struct {
	Form  string
	Final bool
}

// CompletionOrder is the order of the record writes that a completion needs. A changed
// plan takes its amendment before the next chunk entry, and each chunk then retains its
// verification and its reviews. A test-only repair of a reviewed chunk then retains its
// confirmation. A fold of the default branch follows a chunk and precedes the next chunk or
// the final verification. The final verification precedes the completion.
func CompletionOrder() []Step {
	return []Step{{Form: FormAmendment}, {Form: FormChunk}, {Form: FormVerification}, {Form: FormReview}, {Form: FormConfirmation}, {Form: FormFold}, {Form: FormVerification, Final: true}, {Form: FormCompletion}}
}

// currentPlan refuses a record whose plan digest is not the digest of plan, the plan at
// source. Only an amendment moves the record plan digest, so a changed plan takes its
// amendment before another entry reads the plan.
func currentPlan(record Record, plan Plan, source string) error {
	if record.PlanDigest == plan.Digest {
		return nil
	}
	return Repair{Form: FormAmendment, Source: source, Err: errors.New("stale plan digest; review and record the plan delta")}
}
