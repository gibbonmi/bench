package adopt

import (
	"fmt"
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

// LifecycleInput carries the current and desired managed inventories.
type LifecycleInput struct {
	Current, Desired []Asset
	Preserve         []string
}

// LifecyclePlan is the branch-native link, setup, upgrade, and unlink decision.
type LifecyclePlan struct {
	Operations []Operation
	Refusal    string
}

// PlanLifecycle compares immutable inventories without reading or changing a repository.
func PlanLifecycle(input LifecycleInput) LifecyclePlan {
	current, refusal := assetMap(input.Current)
	if refusal != "" {
		return LifecyclePlan{Refusal: refusal}
	}
	desired, refusal := assetMap(input.Desired)
	if refusal != "" {
		return LifecyclePlan{Refusal: refusal}
	}
	preserve := map[string]bool{}
	for _, name := range input.Preserve {
		if !validAssetPath(name) {
			return LifecyclePlan{Refusal: fmt.Sprintf("managed path %q is invalid", name)}
		}
		preserve[name] = true
	}
	paths := make([]string, 0, len(current)+len(desired))
	seen := map[string]bool{}
	for name := range current {
		seen[name] = true
		paths = append(paths, name)
	}
	for name := range desired {
		if !seen[name] {
			paths = append(paths, name)
		}
	}
	sort.Strings(paths)
	plan := LifecyclePlan{}
	for _, name := range paths {
		before, hadBefore := current[name]
		after, hasAfter := desired[name]
		switch {
		case !hadBefore:
			plan.Operations = append(plan.Operations, Operation{Kind: "add", Path: name})
		case !hasAfter && !preserve[name]:
			plan.Operations = append(plan.Operations, Operation{Kind: "remove", Path: name})
		case hasAfter && before != after:
			plan.Operations = append(plan.Operations, Operation{Kind: "change", Path: name})
		}
	}
	return plan
}

func assetMap(assets []Asset) (map[string]string, string) {
	out := make(map[string]string, len(assets))
	for _, asset := range assets {
		if !validAssetPath(asset.Path) || asset.Fingerprint == "" {
			return nil, fmt.Sprintf("managed asset %q is invalid", asset.Path)
		}
		if _, exists := out[asset.Path]; exists {
			return nil, fmt.Sprintf("managed asset %q is duplicated", asset.Path)
		}
		out[asset.Path] = asset.Fingerprint
	}
	return out, ""
}

func validAssetPath(value string) bool {
	return value != "" && value != "." && !strings.HasPrefix(value, "/") && path.Clean(value) == value &&
		!strings.HasPrefix(value, "../") && strings.IndexFunc(value, unicode.IsControl) < 0
}
