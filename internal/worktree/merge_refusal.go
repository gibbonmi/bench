// Merge refusal routing: the merge's faces, the caller's re-run that each face's route
// ends with, and the typed reads that pick a face or a retry.
package worktree

import (
	"context"
	"errors"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/landing"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// The merge refusal face names. The shared registry declares each face, and a raising site
// names the face it raises.
const (
	faceMergeTargetRed             = "merge-target-red"
	faceMergeFoldRed               = "merge-fold-red"
	faceMergeTargetNotClean        = "merge-target-not-clean"
	faceMergeSiblingNotClean       = "merge-sibling-not-clean"
	faceMergeConflict              = "merge-conflict"
	faceMergeInfrastructure        = "merge-infrastructure"
	faceMergePublishedUnreconciled = "merge-published-unreconciled"
	faceMergeHandback              = "merge-handback"
)

// mergeRerun is the caller's own merge command, with the values it passed. A value that is
// not line-safe prints its placeholder.
func mergeRerun(spelling, operand string) string {
	return "bench worktree merge --from " + refusalroute.Arg("from", spelling) + " " + refusalroute.Arg("target", operand)
}

// mergeFaceRoute attaches the caller's own re-run to a merge refusal. A refusal that names
// a registered face takes that face, and an authorization refusal of an infrastructure
// kind takes the infrastructure face. A refusal that carries a route of its own keeps it,
// and every other cause hands back to the reviewer under its own sentence.
func mergeFaceRoute(err error, rerun string) error {
	raised := raisedRefusal(err)
	var refused landing.AuthorizationRefusal
	switch {
	case raised.face != "":
	case errors.As(err, &refused) && refused.Result.Kind == authorization.Infrastructure:
		raised.face = faceMergeInfrastructure
	case raised.next != "":
		return err
	default:
		raised.face = faceMergeHandback
	}
	return landingFaceRefusal(raised.face, raised, rerun, "")
}

// mergeRedRefusal picks the face of a fold that the gate ran red, by the red's cause, and
// passes every other error through. Only the candidate kind attributes the red to the fold.
// Neither an inherited kind nor a lane fail tells a target red from a fold red, so the face
// then follows the grade of the target tip alone. A red target is the agent's own repair,
// a grade that infrastructure stopped takes the infrastructure face, and a green target
// means that the fold adds the red, which FT342 decides.
func mergeRedRefusal(err error, grade func() authorization.Result, spelling, label string) error {
	var refused landing.AuthorizationRefusal
	if !errors.As(err, &refused) || !landing.RedKind(refused.Result.Kind) {
		return err
	}
	face := faceMergeFoldRed
	if refused.Result.Kind != authorization.Candidate {
		switch target := grade().Kind; {
		case landing.RedKind(target):
			face = faceMergeTargetRed
		case target == authorization.Infrastructure:
			face = faceMergeInfrastructure
		}
	}
	return refusalError{refusal{detail: err.Error(), face: face, values: map[string]string{refusalroute.FactLabel: label, refusalroute.FactFrom: spelling}}}
}

// mergeTargetGrade grades the target tip alone with the lane that graded the fold, or with
// the whole gate when the target declares no lane. The lane measures the target against
// the incoming commit, so a selective lane selects each path where the two differ, and
// that covers every path the fold changed.
func mergeTargetGrade(kit string, request landing.MergeRequest) authorization.Result {
	owner, err := mergeOwner(kit, request.Worktree, request.Incoming)
	if err != nil {
		return authorization.Result{Kind: authorization.Infrastructure}
	}
	return owner.GradeTarget(context.Background(), request)
}

// retryEmptyReasonInfrastructureFold retries a fold once when the gate refused it with an
// infrastructure outcome and no reason, and only after the focused verification proves
// the unchanged range green. It reads the typed refusal, not its sentence.
func retryEmptyReasonInfrastructureFold(first error, verify func() int, retry func() error) error {
	var refused landing.AuthorizationRefusal
	if !errors.As(first, &refused) || refused.Result.Kind != authorization.Infrastructure || refused.Result.Reason != "" {
		return first
	}
	if verify() != 0 {
		return first
	}
	return retry()
}

// mergeReconcileNext names the repair the exit-3 boundary leaves the operator: the reset
// verb's plan at the published commit, which reconciles the checkout under a preserved
// envelope. The guard denies an agent's raw git reset, so the repair is the Bench route.
func mergeReconcileNext(target intent.Assignment, tip string) string {
	return resetCommand("--to", tip, target.ID)
}
