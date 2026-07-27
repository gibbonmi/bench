package registry

// The family→check table's own invariants: every binding resolves to a registered
// check, the table is populated at all, and the derived family list is deterministic.

import (
	"slices"
	"sort"
	"testing"
)

// TestFamilyCheckTableBindsRegisteredChecks closes the phantom-check gap: a value
// naming no registry row would scope that family's fixtures to a check that never
// runs, and every one of them would report did-not-bite forever.
func TestFamilyCheckTableBindsRegisteredChecks(t *testing.T) {
	if len(familyChecks) == 0 {
		t.Fatal("the family→check table is empty, so no conformance fixture resolves a scope")
	}
	for family, name := range familyChecks {
		check, found := Find(name)
		if !found {
			t.Errorf("family %q binds check %q, which no registry row carries", family, name)
			continue
		}
		// A fixture with no CHECK file of its own is swept at the dev tier and takes its
		// scope from this table. Bind a family to a ship-tier check and every such
		// fixture in it is scoped to a check the dev sweep never executes, which the
		// conformance driver reds outright as a scope the tier does not run.
		if check.Tier != Dev {
			t.Errorf("family %q binds check %q, which runs at the %s tier; its CHECK-less fixtures are swept at dev and would red for naming a check that tier does not run", family, name, check.Tier)
		}
	}
}

func TestFamiliesIsSortedAndComplete(t *testing.T) {
	families := Families()
	if len(families) != len(familyChecks) {
		t.Fatalf("Families() lists %d of the table's %d entries", len(families), len(familyChecks))
	}
	if !sort.StringsAreSorted(families) {
		t.Fatalf("Families() is unsorted, so the derived list varies run to run:\n%v", families)
	}
	for family := range familyChecks {
		if !slices.Contains(families, family) {
			t.Errorf("Families() omits table entry %q", family)
		}
	}
}

func TestFamilyCheckReportsUnboundFamily(t *testing.T) {
	if _, bound := FamilyCheck("no-such-family"); bound {
		t.Error("FamilyCheck reported a binding for a family the table does not carry")
	}
}
