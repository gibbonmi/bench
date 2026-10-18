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
// because the gather refuses a dirty or moved checkout before a checkout inventory can differ
// from the committed tip spec.
func TestAuditCommittedInventory(t *testing.T) {
	root, slug := preflighttest.SeedAdmitted(t)
	write(t, map[string]string{"outside/claim_test.go": "package outside\n"})
	base := commitAll(t, "unclaimed source")
	write(t, map[string]string{"outside/claim_test.go": "package outside\n\n// changed\n"})
	tip := commitAll(t, "unclaimed source changes")
	// Only this inventory names the source, and no committed spec holds it.
	checkout := ownership.ParseManifest([]byte(inventorySpec(slug, nil, []string{seamClaim(slug, "outside/claim_test.go")})))
	facts, err := staleness.Collect(root, ownership.Source{Base: base, Tip: tip, SpecPath: specPath(slug)}, checkout, ownership.Snapshot{})
	claimed := slices.ContainsFunc(facts.Paths, func(state staleness.PathState) bool { return state.Path == "outside/claim_test.go" })
	held := slices.ContainsFunc(facts.Diagnostics, func(d staleness.Diagnostic) bool { return d.Subject == ownership.MetadataLabel })
	if err != nil || claimed || !held {
		t.Fatalf("checkout inventory = (%v), scope holds its claim %t, diagnostic %t; want the committed tip scope and a %s diagnostic", err, claimed, held, ownership.MetadataLabel)
	}
}
