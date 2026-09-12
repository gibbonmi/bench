package preflight

import (
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// withCompletionPlan reads the spec's completion plan at the gathered source tip
// and returns facts carrying that one answer. The checkpoint the landing runs
// parses the same fence through the same reader, so a plan this row calls valid
// is the plan the landing later accepts.
//
// The read is the gatherer's whole contribution: the digest for the review
// charge's completion table, and the reader's message for the verdict row. An
// unresolvable tip answers as an unreadable plan, because neither state can
// produce the evidence the checkpoint requires.
func withCompletionPlan(root string, facts Facts) Facts {
	tree, err := benchgit.Output("-C", root, "rev-parse", "--verify", facts.SourceTip+"^{tree}")
	if err != nil {
		facts.CompletionPlanError = err.Error()
		return facts
	}
	plan, err := reviewrecord.ReadPlan(root, tree, facts.SpecPath)
	if err != nil {
		facts.CompletionPlanError = err.Error()
		return facts
	}
	facts.CompletionPlanDigest = plan.Digest
	return facts
}

// completionPlanCheck grades that read. Every spec landing is bound to the
// complete checkpoint, and the checkpoint cannot advance a chunk the plan does
// not name. So a spec with no valid plan fence refuses at the start of the
// phase, rather than after the whole build is paid for.
func completionPlanCheck(f Facts) CheckResult {
	if f.CompletionPlanError != "" {
		return red("completion-plan", "spec carries no valid bench-completion-plan fence at "+f.SourceTip+
			": "+f.CompletionPlanError+"; see .bench/BENCH-reference.md, bench gate --checkpoint")
	}
	return green("completion-plan")
}
