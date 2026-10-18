package staleness_test

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight/ownership"
	"github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/preflight/staleness"
)

// specPath is the repository path of the fixture spec.
func specPath(slug string) string {
	return "specs/" + slug + "/spec.md"
}

// itemLocation is the inventory location of an item that the fixture spec's title line holds.
func itemLocation(slug string) string {
	return `"location":{"path":"` + specPath(slug) + `","line":1,"needle":"# ` + slug + `","source":"spec"}`
}

// seamClaim is one inventory test seam whose path is the claim source path.
func seamClaim(slug, path string) string {
	return `{"id":"claim",` + itemLocation(slug) + `,"path":"` + path + `","package":"example.com/outside","function":"TestClaim","state":"existing","rows":["PF1"]}`
}

// inventorySpec is the fixture spec whose inventory holds each given premise and seam.
func inventorySpec(slug string, premises, seams []string) string {
	return preflighttest.SpecBodyNoInventory(slug) + inventorytest.Section(inventorytest.Manifest("declared-complete", premises, seams, nil, nil))
}

// TestAuditCommittedInventory covers the committed tip inventory. It calls Collect directly,
// because the gather takes the source tip from the head of a clean checkout, so no command
// input makes a checkout inventory differ from the committed tip spec.
func TestAuditCommittedInventory(t *testing.T) {
	root, slug := preflighttest.SeedAdmitted(t)
	write(t, map[string]string{"outside/claim_test.go": "package outside\n"})
	base := commitAll(t, "unclaimed source")
	write(t, map[string]string{"outside/claim_test.go": "package outside\n\n// changed\n"})
	tip := commitAll(t, "unclaimed source changes")
	// The clean checkout moves past the tip to a spec whose inventory alone names the source.
	write(t, map[string]string{specPath(slug): inventorySpec(slug, nil, []string{seamClaim(slug, "outside/claim_test.go")})})
	commitAll(t, "checkout claims the source")
	facts, err := staleness.Collect(root, ownership.Source{Base: base, Tip: tip, SpecPath: specPath(slug)}, ownership.Snapshot{})
	if claimed := slices.ContainsFunc(facts.Paths, func(state staleness.PathState) bool { return state.Path == "outside/claim_test.go" }); err != nil || claimed {
		t.Fatalf("checkout inventory = (%v), scope holds its claim %t; want the committed tip scope alone", err, claimed)
	}
}
