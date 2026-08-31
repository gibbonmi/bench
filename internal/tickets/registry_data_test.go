package tickets

import (
	"reflect"
	"sort"
	"testing"
)

// TestRegistryBindingRowsAreSorted keeps the declared binding registry in one
// readable order with one row per package. An unsorted or duplicated table hides
// the row a reader looks for and lets two rows disagree about one prefix.
func TestRegistryBindingRowsAreSorted(t *testing.T) {
	rows := Bindings()
	if len(rows) == 0 {
		t.Fatal("Bindings() is empty, want the declared binding registry")
	}
	prefixes := make([]string, 0, len(rows))
	seen := map[string]bool{}
	for _, row := range rows {
		if seen[row.Prefix] {
			t.Errorf("binding registry declares prefix %q twice", row.Prefix)
		}
		seen[row.Prefix] = true
		prefixes = append(prefixes, row.Prefix)
		if len(row.Files) == 0 {
			t.Errorf("binding row %q binds no file", row.Prefix)
		}
		if !sort.StringsAreSorted(row.Files) {
			t.Errorf("binding row %q files = %v, want sorted", row.Prefix, row.Files)
		}
	}
	if !sort.StringsAreSorted(prefixes) {
		t.Errorf("binding prefixes = %v, want sorted", prefixes)
	}
}

// TestBindingRegistryCarriesTheSeedOwners keeps a binding row for every declared
// seed owner. An owner with no row is an owner a ticket author reconstructs by
// hand.
func TestBindingRegistryCarriesTheSeedOwners(t *testing.T) {
	for _, prefix := range SeedOwners() {
		found := false
		for _, row := range Bindings() {
			if row.Prefix == prefix {
				found = true
			}
		}
		if !found {
			t.Errorf("binding registry has no row for the seed owner %q", prefix)
		}
	}
}

// TestBoundFilesResolvesAtASegmentBoundary is the lookup contract: a bound
// package covers itself and anything under it, never a sibling that shares its
// leading bytes, and never the written file itself.
func TestBoundFilesResolvesAtASegmentBoundary(t *testing.T) {
	got := BoundFiles("internal/toon/toon_test.go")
	want := []string{"internal/conformance/data_handling_test.go"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BoundFiles(internal/toon/toon_test.go) = %v, want %v", got, want)
	}
	if got := BoundFiles("internal/toonery/x.go"); len(got) != 0 {
		t.Errorf("BoundFiles(internal/toonery/x.go) = %v, want none", got)
	}
}

// TestCommandRowsBindTheHelpProjectionAndEnvelopeCases keeps both command
// registries in every command row. A verb that ships without one of them reaches
// the gate with a registry nobody updated.
func TestCommandRowsBindTheHelpProjectionAndEnvelopeCases(t *testing.T) {
	for _, file := range []string{"cmd/bench/command_registry.go", "cmd/bench/command_registry_test.go"} {
		if !holdsString(BoundFiles("internal/harnesses/harnesses.go"), file) {
			t.Errorf("BoundFiles for a command package omits %q", file)
		}
	}
}
