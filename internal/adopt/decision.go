package adopt

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"unicode"
)

// Asset is one immutable managed-path fingerprint.
type Asset struct {
	Path        string
	Fingerprint string
}

// Operation is one ordered lifecycle effect returned by PlanLifecycle.
type Operation struct {
	Kind string
	Path string
}

// LifecycleInput carries the current and desired managed inventories. A non-nil
// Execution selects the execution contract, and PlanLifecycle then ignores the
// inventory fields.
type LifecycleInput struct {
	Current, Desired []Asset
	Preserve         []string
	Execution        *LifecycleInventory
}

// LifecycleInventory is the immutable execution input: the command action, link or
// unlink, and one observation for each path that the command inspects.
type LifecycleInventory struct {
	Action       string
	Observations []LifecycleObservation
}

// LifecycleObservation holds the facts that a caller observed for one path. It holds
// no lifecycle result.
type LifecycleObservation struct {
	Path  string
	Order int // the caller's original order
	// PlanKind is the desired plan entry kind, or claude for the conditional CLAUDE.md
	// owner. An empty kind marks a path that has no desired entry.
	PlanKind string
	Present  bool
	FileKind string // the destination kind: regular, symlink, or other
	Perm     fs.FileMode
	// Recorded is the manifest fingerprint, empty when the path is unowned. Observed is
	// the destination fingerprint, empty when the destination has none.
	Recorded, Observed string
	Staged, StagedKind string
	StagedPerm         fs.FileMode
	// ResolvedEquivalent reports that the destination resolves to the bytes that a staged
	// symlink names. CanonicalAdapter reports that the destination is the canonical
	// adapter target file.
	ResolvedEquivalent bool
	CanonicalAdapter   bool
	SymlinkParent      bool
	Escapes            bool // the path resolves outside the repository
	Strict             bool // compatibility repair also compares regular-file modes
	KnownForm          bool // CLAUDE.md holds a known Bench form
}

// LifecycleDecision is the complete lifecycle choice for one path. ManifestAction is
// omit, retain, or record, and Fingerprint is the value that it retains or records.
type LifecycleDecision struct {
	Path, Kind, Reason, ManifestAction, Fingerprint string
	Order                                           int
}

// LifecyclePlan is the branch-native link, setup, upgrade, and unlink decision.
// Operations is the effect projection of Decisions.
type LifecyclePlan struct {
	Decisions  []LifecycleDecision
	Operations []Operation
	Refusal    string
}

const (
	lifecycleAdd      = "add"
	lifecycleChange   = "change"
	lifecycleRemove   = "remove"
	lifecyclePreserve = "preserve"

	manifestOmit   = "omit"
	manifestRetain = "retain"
	manifestRecord = "record"

	claudePlanKind = "claude"
)

// PlanLifecycle decides every path from immutable input without reading or changing a
// repository.
func PlanLifecycle(input LifecycleInput) LifecyclePlan {
	var decisions []LifecycleDecision
	var refusal string
	if input.Execution != nil {
		decisions, refusal = executionDecisions(*input.Execution)
	} else {
		decisions, refusal = inventoryDecisions(input)
	}
	if refusal != "" {
		return LifecyclePlan{Refusal: refusal}
	}
	sort.Slice(decisions, func(i, j int) bool { return decisions[i].Path < decisions[j].Path })
	plan := LifecyclePlan{Decisions: decisions}
	for _, d := range decisions {
		switch d.Kind {
		case lifecycleAdd, lifecycleChange, lifecycleRemove:
			plan.Operations = append(plan.Operations, Operation{Kind: d.Kind, Path: d.Path})
		}
	}
	return plan
}

// inventoryDecisions compares the current and desired inventories. Preserve suppresses
// a removal but not a desired change, and a Preserve path outside both inventories gets
// no decision.
func inventoryDecisions(input LifecycleInput) ([]LifecycleDecision, string) {
	current, refusal := assetMap(input.Current)
	if refusal != "" {
		return nil, refusal
	}
	desired, refusal := assetMap(input.Desired)
	if refusal != "" {
		return nil, refusal
	}
	preserve := map[string]bool{}
	for _, name := range input.Preserve {
		if !validAssetPath(name) {
			return nil, fmt.Sprintf("managed path %q is invalid", name)
		}
		preserve[name] = true
	}
	decisions := make([]LifecycleDecision, 0, len(current)+len(desired))
	for name, before := range current {
		after, hasAfter := desired[name]
		switch {
		case !hasAfter && !preserve[name]:
			decisions = append(decisions, LifecycleDecision{Path: name, Kind: lifecycleRemove, ManifestAction: manifestOmit})
		case !hasAfter:
			decisions = append(decisions, preserved(name, "requested-preserve", manifestRetain, before))
		case before != after:
			decisions = append(decisions, recorded(name, lifecycleChange, after))
		default:
			decisions = append(decisions, preserved(name, "unchanged", manifestRetain, before))
		}
	}
	for name, after := range desired {
		if _, hadBefore := current[name]; !hadBefore {
			decisions = append(decisions, recorded(name, lifecycleAdd, after))
		}
	}
	return decisions, ""
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
	}
	adoptAdapter := o.Recorded == "" && o.PlanKind == "adapter" && o.SymlinkParent && o.CanonicalAdapter
	if o.Present && (o.Recorded != "" || adoptAdapter) && converged(o) {
		return preserved(o.Path, "converged", manifestRecord, o.Observed), ""
	}
	if o.SymlinkParent {
		return LifecycleDecision{}, symlinkParentRefusal(o.Path)
	}
	if o.Present && !ownedClean(o) {
		if o.Recorded != "" {
			return preserved(o.Path, "modified-managed", manifestRetain, o.Recorded), ""
		}
		return preserved(o.Path, "project-owned", manifestOmit, ""), ""
	}
	return writeDecision(o), ""
}

// converged reports that the destination already holds what the staged entry writes.
// A fingerprint covers content only, so a regular destination also compares its mode.
func converged(o LifecycleObservation) bool {
	if o.Observed == "" {
		return false
	}
	if o.StagedKind == "symlink" {
		return o.Observed == o.Staged || o.ResolvedEquivalent
	}
	if o.FileKind != "symlink" && o.Perm != o.StagedPerm {
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
	return !(o.Strict && o.FileKind == "regular" && o.StagedKind == "regular" && o.Perm != o.StagedPerm)
}

// seedDecision installs an absent seed once and preserves a present one. A seed adds
// no manifest row, but it keeps an ownership row that already exists.
func seedDecision(o LifecycleObservation) (LifecycleDecision, string) {
	if o.SymlinkParent {
		return LifecycleDecision{}, symlinkParentRefusal(o.Path)
	}
	action, fingerprint := manifestOmit, ""
	if o.Recorded != "" {
		action, fingerprint = manifestRetain, o.Recorded
	}
	if o.Present {
		return preserved(o.Path, "seed-owned", action, fingerprint), ""
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
		return preserved(o.Path, "modified-managed", manifestRetain, o.Recorded)
	}
	return preserved(o.Path, "project-owned", manifestOmit, "")
}

// withdrawalDecision removes a clean manifest row that link no longer plans.
func withdrawalDecision(o LifecycleObservation) (LifecycleDecision, string) {
	switch {
	case o.Escapes || o.SymlinkParent:
		return LifecycleDecision{}, "refused manifest path: " + o.Path
	case !o.Present:
		return preserved(o.Path, "absent", manifestOmit, ""), ""
	case o.FileKind != "regular" && o.FileKind != "symlink":
		return LifecycleDecision{}, "refused non-regular dropped asset: " + o.Path
	case o.Observed == "":
		return LifecycleDecision{}, fmt.Sprintf("dropped asset %q has no fingerprint", o.Path)
	case o.Observed == o.Recorded:
		return LifecycleDecision{Kind: lifecycleRemove, ManifestAction: manifestOmit}, ""
	}
	return preserved(o.Path, "kept-modified-removed", manifestRetain, o.Recorded), ""
}

// unlinkDecision removes a clean manifest row. An unsafe or modified row stays as a
// residual whose manifest row the caller keeps.
func unlinkDecision(o LifecycleObservation) (LifecycleDecision, string) {
	switch {
	case o.Escapes || o.SymlinkParent:
		return preserved(o.Path, "refused", manifestRetain, o.Recorded), ""
	case !o.Present:
		return preserved(o.Path, "absent", manifestOmit, ""), ""
	case o.Observed == "":
		return preserved(o.Path, "refused", manifestRetain, o.Recorded), ""
	case o.Observed != o.Recorded:
		return preserved(o.Path, "modified", manifestRetain, o.Recorded), ""
	}
	return LifecycleDecision{Kind: lifecycleRemove, ManifestAction: manifestOmit}, ""
}

func writeDecision(o LifecycleObservation) LifecycleDecision {
	kind := lifecycleAdd
	if o.Present {
		kind = lifecycleChange
	}
	return recorded(o.Path, kind, o.Staged)
}

func recorded(path, kind, fingerprint string) LifecycleDecision {
	return LifecycleDecision{Path: path, Kind: kind, ManifestAction: manifestRecord, Fingerprint: fingerprint}
}

func preserved(path, reason, action, fingerprint string) LifecycleDecision {
	return LifecycleDecision{Path: path, Kind: lifecyclePreserve, Reason: reason, ManifestAction: action, Fingerprint: fingerprint}
}

func symlinkParentRefusal(path string) string {
	return "conflict: " + path + " has a symlink parent directory"
}

func duplicatedAsset(path string) string {
	return fmt.Sprintf("managed asset %q is duplicated", path)
}

func assetMap(assets []Asset) (map[string]string, string) {
	out := make(map[string]string, len(assets))
	for _, asset := range assets {
		if !validAssetPath(asset.Path) || asset.Fingerprint == "" {
			return nil, fmt.Sprintf("managed asset %q is invalid", asset.Path)
		}
		if _, exists := out[asset.Path]; exists {
			return nil, duplicatedAsset(asset.Path)
		}
		out[asset.Path] = asset.Fingerprint
	}
	return out, ""
}

func validAssetPath(value string) bool {
	return value != "" && value != "." && !strings.HasPrefix(value, "/") && path.Clean(value) == value &&
		!strings.HasPrefix(value, "../") && strings.IndexFunc(value, unicode.IsControl) < 0
}
