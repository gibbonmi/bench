package registry

// The family→check table's own invariants: every binding resolves to a registered
// check, the table is populated at all, and the derived family list is deterministic.

import (
	"slices"
	"sort"
	"testing"
)

// TestFamilyCheckTableBindsRegisteredChecks closes the phantom-check gap: a value
// naming no registry row leaves that family's fixtures without a resolved binding.
func TestFamilyCheckTableBindsRegisteredChecks(t *testing.T) {
	if len(familyChecks) == 0 {
		t.Fatal("the family→check table is empty, so no conformance fixture resolves a binding")
	}
	for family, name := range familyChecks {
		check, found := Find(name)
		if !found {
			t.Errorf("family %q binds check %q, which no registry row carries", family, name)
			continue
		}
		// A fixture with no CHECK file resolves a dev-tier binding from this table. A
		// ship-tier family binding therefore leaves that fixture outside its declared
		// tier, which the conformance driver rejects.
		if check.Tier != Dev {
			t.Errorf("family %q binds check %q at the %s tier; its CHECK-less fixtures require a dev-tier binding", family, name, check.Tier)
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

// TestTierForDefaultsToDev pins the asymmetry that keeps the default un-overridable by
// accident: only the exact ship name widens a run, so an unset, misspelled, or
// differently-cased value grades the dev tier instead of quietly reaching ship work.
func TestTierForDefaultsToDev(t *testing.T) {
	for _, test := range []struct {
		value string
		want  Tier
	}{
		{"", Dev},
		{"Ship", Dev},
		{"dev", Dev},
		{"anything", Dev},
		{string(Ship), Ship},
	} {
		if got := TierFor(test.value); got != test.want {
			t.Errorf("TierFor(%q) = %q, want %q", test.value, got, test.want)
		}
	}
}

func TestFamilyCheckReportsUnboundFamily(t *testing.T) {
	if _, bound := FamilyCheck("no-such-family"); bound {
		t.Error("FamilyCheck reported a binding for a family the table does not carry")
	}
}
