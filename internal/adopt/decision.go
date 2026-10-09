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
	// PlanKind is the desired plan entry kind, claude for the conditional CLAUDE.md owner,
	// or bespoke for a replacement that a domain helper prepared. An empty kind marks a
	// path that has no desired entry.
	PlanKind string
	Present  bool
	FileKind string // the destination kind that fileKind names
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
	// Unprepared reports that a domain helper could not prepare the transform of a bespoke
	// path, such as an AGENTS.md special file or a malformed managed block.
	Unprepared bool
}

const (
	fileRegular = "regular"
	fileSymlink = "symlink"
	fileOther   = "other"
	fileAbsent  = "absent" // the staged kind of a prepared transform that removes the path
)

// fileKind names a file mode in the observation vocabulary.
func fileKind(mode fs.FileMode) string {
	switch {
	case mode.IsRegular():
		return fileRegular
	case mode&fs.ModeSymlink != 0:
		return fileSymlink
	}
	return fileOther
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

	claudePlanKind  = "claude"
	bespokePlanKind = "bespoke"

	// The preserve reasons. A caller that reports a conflict matches and prints these.
	reasonConverged           = "converged"
	reasonModifiedManaged     = "modified-managed"
	reasonProjectOwned        = "project-owned"
	reasonSeedOwned           = "seed-owned"
	reasonKeptModifiedRemoved = "kept-modified-removed"
	reasonAbsent              = "absent"
	reasonRefused             = "refused"
	reasonModified            = "modified"
	reasonUnchanged           = "unchanged"
	reasonRequestedPreserve   = "requested-preserve"
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
			decisions = append(decisions, preserved(name, reasonRequestedPreserve, manifestRetain, before))
		case before != after:
			decisions = append(decisions, recorded(name, lifecycleChange, after))
		default:
			decisions = append(decisions, preserved(name, reasonUnchanged, manifestRetain, before))
		}
	}
	for name, after := range desired {
		if _, hadBefore := current[name]; !hadBefore {
			decisions = append(decisions, recorded(name, lifecycleAdd, after))
		}
	}
	return decisions, ""
}

func recorded(path, kind, fingerprint string) LifecycleDecision {
	return LifecycleDecision{Path: path, Kind: kind, ManifestAction: manifestRecord, Fingerprint: fingerprint}
}

func preserved(path, reason, action, fingerprint string) LifecycleDecision {
	return LifecycleDecision{Path: path, Kind: lifecyclePreserve, Reason: reason, ManifestAction: action, Fingerprint: fingerprint}
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
