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

func resolveFixtureOwner(fixture Fixture) string {
	if fixture.Check != "" {
		return "check:" + fixture.Check
	}
	return ""
}

func unsafeDecisionText(value string) bool {
	return strings.IndexFunc(value, unicode.IsControl) >= 0
}
