package adopt

import "fmt"

// decideObservation submits one execution observation for action to the lifecycle owner. It
// returns the complete decision of that observation, or the owner's hard refusal.
func decideObservation(action string, o LifecycleObservation) (LifecycleDecision, string) {
	plan := PlanLifecycle(LifecycleInput{Execution: &LifecycleInventory{Action: action, Observations: []LifecycleObservation{o}}})
	if plan.Refusal != "" {
		return LifecycleDecision{}, plan.Refusal
	}
	return plan.Decisions[0], ""
}

// executionDecisions decides each observation for its action. A hard refusal refuses
// the whole plan, and the refusal of the earliest caller order wins, as it does in the
// caller's own checks.
func executionDecisions(inventory LifecycleInventory) ([]LifecycleDecision, string) {
	if len(inventory.Observations) == 0 {
		return nil, ""
	}
	var decide func(LifecycleObservation) (LifecycleDecision, string)
	switch inventory.Action {
	case "link":
		decide = linkDecision
	case "unlink":
		decide = unlinkDecision
	default:
		return nil, fmt.Sprintf("lifecycle action %q is invalid", inventory.Action)
	}
	decisions := make([]LifecycleDecision, 0, len(inventory.Observations))
	seen := map[string]bool{}
	refusal, refusalOrder := "", 0
	for _, o := range inventory.Observations {
		if seen[o.Path] {
			return nil, duplicatedAsset(o.Path)
		}
		seen[o.Path] = true
		d, refused := decide(o)
		if refused != "" && (refusal == "" || o.Order < refusalOrder) {
			refusal, refusalOrder = refused, o.Order
		}
		d.Path, d.Order = o.Path, o.Order
		decisions = append(decisions, d)
	}
	if refusal != "" {
		return nil, refusal
	}
	return decisions, ""
}

// linkDecision decides one link path. A converged destination outranks the
// symlink-parent refusal, and that refusal outranks the preservation of a collision.
func linkDecision(o LifecycleObservation) (LifecycleDecision, string) {
	switch o.PlanKind {
	case "":
		return withdrawalDecision(o)
	case "seed":
		return seedDecision(o)
	case claudePlanKind:
		return claudeDecision(o), ""
	case bespokePlanKind:
		return bespokeDecision(o), ""
	}
	adoptAdapter := o.Recorded == "" && o.PlanKind == "adapter" && o.SymlinkParent && o.CanonicalAdapter
	if o.Present && (o.Recorded != "" || adoptAdapter) && converged(o) {
		return preserved(o.Path, reasonConverged, manifestRecord, o.Observed), ""
	}
	if o.SymlinkParent {
		return LifecycleDecision{}, symlinkParentRefusal(o.Path)
	}
	if o.Present && !ownedClean(o) {
		if o.Recorded != "" {
			return preserved(o.Path, reasonModifiedManaged, manifestRetain, o.Recorded), ""
		}
		return preserved(o.Path, reasonProjectOwned, manifestOmit, ""), ""
	}
	return writeDecision(o), ""
}

// converged reports that the destination already holds what the staged entry writes.
// It compares the staged bytes rather than the recorded fingerprint, so a stale owned
// asset still takes the current release. A fingerprint covers content only, so a regular
// destination also compares its mode.
func converged(o LifecycleObservation) bool {
	if o.Observed == "" {
		return false
	}
	if o.StagedKind == fileSymlink {
		return o.Observed == o.Staged || o.ResolvedEquivalent
	}
	if o.FileKind != fileSymlink && o.Perm != o.StagedPerm {
		return false
	}
	return o.Observed == o.Staged
}

// ownedClean reports that an owned destination still holds its recorded bytes. Strict
// repair also refuses a regular file whose mode differs from the staged mode.
func ownedClean(o LifecycleObservation) bool {
	if o.Recorded == "" || o.Observed != o.Recorded {
		return false
	}
	return !(o.Strict && o.FileKind == fileRegular && o.StagedKind == fileRegular && o.Perm != o.StagedPerm)
}

// seedDecision installs an absent seed once and preserves a present one. A seed adds
// no manifest row, because a recorded row would read a later project edit back as a
// modified-managed conflict. It keeps an ownership row that already exists.
func seedDecision(o LifecycleObservation) (LifecycleDecision, string) {
	if o.SymlinkParent {
		return LifecycleDecision{}, symlinkParentRefusal(o.Path)
	}
	action, fingerprint := manifestOmit, ""
	if o.Recorded != "" {
		action, fingerprint = manifestRetain, o.Recorded
	}
	if o.Present {
		return preserved(o.Path, reasonSeedOwned, action, fingerprint), ""
	}
	return LifecycleDecision{Kind: lifecycleAdd, ManifestAction: action, Fingerprint: fingerprint}, ""
}

// claudeDecision installs the current CLAUDE.md form over an absent or known form. Other
// content, including a special file, is project-owned, except owned content that differs
// from its recorded bytes.
func claudeDecision(o LifecycleObservation) LifecycleDecision {
	switch {
	case !o.Present || o.KnownForm:
		return writeDecision(o)
	case o.Recorded != "" && o.Observed != "" && o.Observed != o.Recorded:
		return preserved(o.Path, reasonModifiedManaged, manifestRetain, o.Recorded)
	}
	return preserved(o.Path, reasonProjectOwned, manifestOmit, "")
}

// bespokeDecision publishes the replacement that a domain helper prepared, such as the
// AGENTS.md managed block or the managed pre-push hook. The helper owns the transform, so
// the owner reads only whether a replacement exists. A bespoke path has no manifest row.
func bespokeDecision(o LifecycleObservation) LifecycleDecision {
	if o.Staged == "" {
		return preserved(o.Path, reasonUnchanged, manifestOmit, "")
	}
	d := writeDecision(o)
	d.ManifestAction, d.Fingerprint = manifestOmit, ""
	return d
}

// withdrawalDecision removes a clean manifest row that link no longer plans.
func withdrawalDecision(o LifecycleObservation) (LifecycleDecision, string) {
	switch {
	case o.Escapes || o.SymlinkParent:
		return LifecycleDecision{}, "refused manifest path: " + o.Path
	case !o.Present:
		return preserved(o.Path, reasonAbsent, manifestOmit, ""), ""
	case o.FileKind != fileRegular && o.FileKind != fileSymlink:
		return LifecycleDecision{}, "refused non-regular dropped asset: " + o.Path
	case o.Observed == "":
		return LifecycleDecision{}, fmt.Sprintf("dropped asset %q has no fingerprint", o.Path)
	case o.Observed == o.Recorded:
		return LifecycleDecision{Kind: lifecycleRemove, ManifestAction: manifestOmit}, ""
	}
	return preserved(o.Path, reasonKeptModifiedRemoved, manifestRetain, o.Recorded), ""
}

// unlinkDecision removes a clean manifest row. An unsafe or modified row stays as a
// residual whose manifest row the caller keeps. Unlink never refuses the whole plan.
func unlinkDecision(o LifecycleObservation) (LifecycleDecision, string) {
	switch {
	case o.PlanKind == bespokePlanKind:
		return unlinkBespokeDecision(o), ""
	case o.Escapes || o.SymlinkParent:
		return preserved(o.Path, reasonRefused, manifestRetain, o.Recorded), ""
	case !o.Present:
		return preserved(o.Path, reasonAbsent, manifestOmit, ""), ""
	case o.Observed == "":
		return preserved(o.Path, reasonRefused, manifestRetain, o.Recorded), ""
	case o.Observed != o.Recorded:
		return preserved(o.Path, reasonModified, manifestRetain, o.Recorded), ""
	}
	return LifecycleDecision{Kind: lifecycleRemove, ManifestAction: manifestOmit}, ""
}

// unlinkBespokeDecision publishes the transform that a domain helper prepared for unlink:
// the stripped AGENTS.md, or the removal of the managed hook or of an AGENTS.md that held
// only the managed block. A path whose transform failed stays as a refused residual, and
// it keeps the manifest as a refused row does. Any other path stays in place.
func unlinkBespokeDecision(o LifecycleObservation) LifecycleDecision {
	switch {
	case o.Unprepared:
		return preserved(o.Path, reasonRefused, manifestRetain, "")
	case !o.Present:
		return preserved(o.Path, reasonAbsent, manifestOmit, "")
	case o.StagedKind == fileAbsent:
		return LifecycleDecision{Kind: lifecycleRemove, ManifestAction: manifestOmit}
	}
	return bespokeDecision(o)
}

func writeDecision(o LifecycleObservation) LifecycleDecision {
	kind := lifecycleAdd
	if o.Present {
		kind = lifecycleChange
	}
	return recorded(o.Path, kind, o.Staged)
}

func symlinkParentRefusal(path string) string {
	return "conflict: " + path + " has a symlink parent directory"
}
