package canary

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// FixtureOwner binds one mutation fixture to the production check that grades it.
type FixtureOwner struct {
	Fixture string
	Owner   string
}

// Selection is the branch-native canary ownership and aggregation decision.
type Selection struct {
	Accepted    bool
	Owners      []FixtureOwner
	Diagnostics []string
}

// DispatchResult is the aggregated result of invoking every selected owner once.
type DispatchResult struct {
	Accepted    bool
	Dispatched  []FixtureOwner
	Diagnostics []string
}

// Select assigns every immutable fixture value to exactly one owning check.
func Select(fixtures []Fixture) Selection {
	decision := Selection{}
	seen := map[string]bool{}
	for _, fixture := range fixtures {
		name := filepath.Base(fixture.Dir)
		owner := resolveFixtureOwner(fixture)
		switch {
		case name == "." || name == "" || unsafeDecisionText(name):
			decision.Diagnostics = append(decision.Diagnostics, "canary fixture has an invalid identity")
		case seen[name]:
			decision.Diagnostics = append(decision.Diagnostics, fmt.Sprintf("canary fixture %q has more than one owner", name))
		case owner == "" || unsafeDecisionText(owner):
			decision.Diagnostics = append(decision.Diagnostics, fmt.Sprintf("canary fixture %q has no valid production owner", name))
		default:
			seen[name] = true
			decision.Owners = append(decision.Owners, FixtureOwner{Fixture: name, Owner: owner})
		}
	}
	if len(fixtures) == 0 {
		decision.Diagnostics = append(decision.Diagnostics, "canary fixture inventory is empty")
	}
	sort.Slice(decision.Owners, func(i, j int) bool { return decision.Owners[i].Fixture < decision.Owners[j].Fixture })
	sort.Strings(decision.Diagnostics)
	decision.Accepted = len(decision.Diagnostics) == 0 && len(decision.Owners) == len(fixtures)
	return decision
}

// Dispatch invokes each selected owner in stable order and retains every diagnostic.
func Dispatch(selection Selection, invoke func(FixtureOwner) string) DispatchResult {
	result := DispatchResult{}
	if !selection.Accepted || invoke == nil {
		result.Diagnostics = append(result.Diagnostics, selection.Diagnostics...)
		if len(result.Diagnostics) == 0 {
			result.Diagnostics = append(result.Diagnostics, "canary dispatch is unavailable")
		}
		return result
	}
	for _, owner := range selection.Owners {
		result.Dispatched = append(result.Dispatched, owner)
		if diagnostic := invoke(owner); diagnostic != "" {
			result.Diagnostics = append(result.Diagnostics, fmt.Sprintf("canary %q (%s): %s", owner.Fixture, owner.Owner, diagnostic))
		}
	}
	result.Accepted = len(result.Diagnostics) == 0 && len(result.Dispatched) == len(selection.Owners)
	return result
}

func resolveFixtureOwner(fixture Fixture) string {
	if fixture.Check != "" {
		return "check:" + fixture.Check
	}
	return ""
}

func unsafeDecisionText(value string) bool {
	return strings.IndexFunc(value, unicode.IsControl) >= 0
}
