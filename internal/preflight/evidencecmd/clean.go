package evidencecmd

import (
	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/toon"
)

// CleanCommand is the public cleanup command the responses name in their successors.
const CleanCommand = "bench preflight " + ModeClean

// freshPlan is the recovery action every stale or stopped cleanup names. A plan is the one
// authorization cleanup accepts, so recovery is always a new plan and never a retry.
const freshPlan = "run " + CleanCommand + " for a fresh plan"

// CleanPlan prints one bounded page of the exact deletion targets and the fingerprint that
// commits to them. It deletes nothing and creates nothing.
func CleanPlan(root string, flags map[string]string) (string, int) {
	store, refusal := evidenceStore(root)
	if refusal != "" {
		return refusal, 1
	}
	plan, err := store.Plan()
	if err != nil {
		return storeRefusal(err), 1
	}
	index := 0
	if text, ok := flags[flagCursor]; ok {
		cursor, err := chargeevidence.ParseCursor(text, plan.Fingerprint)
		if err != nil {
			return storeRefusal(err), 1
		}
		if !cursor.Clean {
			return toon.Errorf("evidence "+chargeevidence.RefuseCursor+": a cleanup page takes a cleanup cursor", freshPlan) + "\n", 1
		}
		index = cursor.Index
	}
	targets, next, err := plan.Page(index)
	if err != nil {
		return storeRefusal(err), 1
	}
	out, err := plan.Encode(targets, next == nil, planSuccessor(plan, next))
	if err != nil {
		return storeRefusal(err), 1
	}
	return out, 0
}

// planSuccessor is the exact command that continues or authorizes one plan. A plan with more
// targets continues at its next cursor; a completely delivered plan names its apply. An empty
// plan authorizes nothing, so it names no successor.
func planSuccessor(plan chargeevidence.CleanupPlan, next *chargeevidence.Cursor) string {
	switch {
	case next != nil:
		return CleanCommand + " " + flagCursor + " " + next.String()
	case len(plan.Targets) == 0:
		return ""
	}
	return CleanCommand + " " + flagApply + " " + plan.Fingerprint
}

// CleanApply deletes exactly the targets one fingerprinted plan named. A changed store
// refuses before any deletion, and a stopped apply reports its exact disposition.
func CleanApply(root string, flags map[string]string) (string, int) {
	store, refusal := evidenceStore(root)
	if refusal != "" {
		return refusal, 1
	}
	// A plan fingerprint is an artifact identity over the target list, so the manifest owner
	// of that shape validates it before any store access.
	fingerprint := flags[flagApply]
	if !chargeevidence.ValidIdentity(fingerprint) {
		return toon.Usage(Grammar.Cmd, flagApply+" needs a plan fingerprint the cleanup plan printed"), 2
	}
	applied, err := store.Apply(fingerprint)
	if err != nil {
		return storeRefusal(err), 1
	}
	next := ""
	code := 0
	if !applied.Complete {
		next, code = freshPlan, 1
	}
	out, err := applied.Encode(next)
	if err != nil {
		return storeRefusal(err), 1
	}
	return out, code
}
