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
	plan, err = planReset(root, plan.assignment.ID, plan.checkpoint, plan.envelope)
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
	movePlan := plan
	if plan.mode == "restore" {
		movePlan.checkpoint = plan.manifest.Tip
	}
	if err := j.resetMove(plan.assignment.Worktree, plan.assignment.Branch, movePlan.checkpoint); err != nil {
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
	if !resetCheckoutMatches(movePlan) {
		return resetResult(stdout, plan, envelope.Ref, 3)
	}
	if plan.mode == "restore" {
		if plan.manifest.Base != plan.manifest.Tip {
			if _, err := git.Output("-C", plan.assignment.Worktree, "switch", "--detach", plan.manifest.Base); err != nil {
				return resetResult(stdout, plan, envelope.Ref, 3)
			}
		}
		if err := j.resetLayers(plan.assignment.Worktree, plan.manifest); err != nil {
			return resetResult(stdout, plan, envelope.Ref, 3)
		}
		if !resetRestoredLayersMatch(root, plan) {
			return resetResult(stdout, plan, envelope.Ref, 3)
		}
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
	restore := "none"
	if preserved != "none" {
		restore = "bench worktree reset --restore " + preserved + " " + plan.assignment.ID
	}
	if code == 3 {
		next = ",next=" + resetPlanCommand(plan)
		if restore != "none" {
			next = ",next=" + restore
		}
	}
	ref := plan.assignment.Branch
	if plan.mode == "restore" && plan.manifest.Base != plan.manifest.Tip {
		ref = "detached"
		if code == 0 {
			next = ",next=bench worktree reset --to " + plan.manifest.Tip + " " + plan.assignment.ID
		}
	}
	fmt.Fprintf(stdout, "reset{worktree=%s,mode=%s,checkpoint=%s,previous=%s,ref=%s,preserved=%s,restore=%s%s}\n", plan.assignment.ID, plan.mode, plan.checkpoint, plan.head, ref, preserved, restore, next)
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
