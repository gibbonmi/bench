package worktree

import (
	"fmt"
	"io"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
)

func applyReset(j joins, root, home string, plan resetPlan, fingerprint string, stdout io.Writer) (exit int) {
	finish := beginVerbSpan(home, root, otelResetSeam)
	defer func() { finish(exit, plan.assignment.ID) }()
	unlock, err := lockCleanupRegistration(j, root, plan.assignment.Worktree)
	if err != nil {
		return landRefusalError(stdout, err)
	}
	defer unlock()
	plan, err = planReset(root, plan.assignment.ID, plan.checkpoint)
	if err != nil {
		return landRefusalError(stdout, err)
	}
	if plan.action == "none" || fingerprint != plan.fingerprint {
		return landRefusalError(stdout, refusalError{refusal{detail: "reset plan is stale", wanted: plan.fingerprint, next: resetPlanCommand(plan)}})
	}
	envelope := intent.Recovery{Ref: "none"}
	if plan.preserve == "envelope" {
		envelope, err = j.resetEnvelope(root, plan)
		if err != nil {
			return landRefusalError(stdout, err)
		}
		if !resetEnvelopeValid(root, envelope) {
			return landRefusal(stdout, "reset envelope failed verification")
		}
	}
	if err := j.resetMove(plan.assignment.Worktree, plan.assignment.Branch, plan.checkpoint); err != nil {
		return resetResult(stdout, plan, envelope.Ref, 3)
	}
	if plan.lock == "repair" {
		if plan.registrationLocked {
			if err := unlockWorktree(root, plan.assignment.Worktree); err != nil {
				return resetResult(stdout, plan, envelope.Ref, 3)
			}
		}
		if err := lockWorktree(root, plan.assignment.Worktree, lockReason(plan.assignment)); err != nil {
			return resetResult(stdout, plan, envelope.Ref, 3)
		}
	}
	if !resetCheckoutMatches(plan) {
		return resetResult(stdout, plan, envelope.Ref, 3)
	}
	return resetResult(stdout, plan, envelope.Ref, 0)
}

func resetCheckoutMatches(plan resetPlan) bool {
	head, err := git.Output("-C", plan.assignment.Worktree, "rev-parse", "HEAD^{commit}")
	if err != nil || head != plan.checkpoint {
		return false
	}
	ref, err := git.Output("-C", plan.assignment.Worktree, "symbolic-ref", "--quiet", "HEAD")
	if err != nil || ref != plan.assignment.Branch {
		return false
	}
	status, err := resetStatus(plan.assignment.Worktree)
	return err == nil && len(status) == 0
}

func resetResult(stdout io.Writer, plan resetPlan, preserved string, code int) int {
	next := ""
	if code == 3 {
		next = ",next=" + resetPlanCommand(plan)
	}
	fmt.Fprintf(stdout, "reset{worktree=%s,mode=reset,checkpoint=%s,previous=%s,ref=%s,preserved=%s,restore=none%s}\n", plan.assignment.ID, plan.checkpoint, plan.head, plan.assignment.Branch, preserved, next)
	return code
}

func moveResetCheckout(path, branch, checkpoint string) error {
	if _, err := git.Output("-C", path, "symbolic-ref", "HEAD", branch); err != nil {
		return err
	}
	if _, err := git.Output("-C", path, "reset", "--hard", checkpoint); err != nil {
		return err
	}
	_, err := git.Output("-C", path, "clean", "-fd")
	return err
}
