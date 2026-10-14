package tickets

import (
	"os"
	"path"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
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

// TestCommandPackageReadsTheBindingRows grades the command-package predicate against the
// binding rows: a path under a command row answers true, and a path under the
// `internal/toon` row, which binds other files, answers false.
func TestCommandPackageReadsTheBindingRows(t *testing.T) {
	for path, want := range map[string]bool{
		"internal/harnesses/harnesses.go": true,
		"cmd/bench/main.go":               true,
		"internal/toon/toon.go":           false,
		"internal/harnessesx/x.go":        false,
	} {
		if got := CommandPackage(path); got != want {
			t.Errorf("CommandPackage(%q) = %t, want %t", path, got, want)
		}
	}
}

// TestCommandRowsBindTheHelpProjectionAndEnvelopeCases keeps both command
// registries in every command row. A verb that ships without one of them reaches
// the gate with a registry nobody updated.
func TestCommandRowsBindTheHelpProjectionAndEnvelopeCases(t *testing.T) {
	for _, file := range []string{"cmd/bench/command_registry.go", "cmd/bench/command_registry_test.go"} {
		if !slices.Contains(BoundFiles("internal/harnesses/harnesses.go"), file) {
			t.Errorf("BoundFiles for a command package omits %q", file)
		}
	}
}

// TestBindingLinesLocateEveryRow reads this registry's own source. BindingRegistry must
// name that file, and BindingLines must place each compiled row on the line that spells
// its prefix, and each grammar root on the line that spells its name, so a closure reason
// cites the row it applied.
func TestBindingLinesLocateEveryRow(t *testing.T) {
	src, err := os.ReadFile(path.Base(BindingRegistry))
	if err != nil {
		t.Fatalf("BindingRegistry %s is not this package's file: %v", BindingRegistry, err)
	}
	lines, err := BindingLines(src)
	if err != nil {
		t.Fatalf("BindingLines: %v", err)
	}
	physical := strings.Split(string(src), "\n")
	for _, row := range Bindings() {
		line := lines[row.Prefix]
		if line < 1 || !strings.Contains(physical[line-1], "Prefix: "+strconv.Quote(row.Prefix)) {
			t.Errorf("BindingLines()[%q] = %d, want the line that declares the row", row.Prefix, line)
		}
	}
	for _, root := range GrammarRoots() {
		line := lines[root.String()]
		if line < 1 || !strings.Contains(physical[line-1], "Name: "+strconv.Quote(root.Name)) {
			t.Errorf("BindingLines()[%q] = %d, want the line that declares the root", root, line)
		}
	}
	if _, err := BindingLines([]byte("package tickets\n")); err == nil {
		t.Error("BindingLines over source without a bindings table succeeded, want an error")
	}
}

// TestBindingRulesProjectBoundFiles keeps the typed rules and the legacy projection on one
// table. A path under the internal/toon row takes one rule for each bound file, a bound
// file takes no rule that requires itself, and a sibling package takes none.
func TestBindingRulesProjectBoundFiles(t *testing.T) {
	for path, want := range map[string][]BindingRule{
		"internal/toon/toon.go": {
			{Prefix: "internal/toon", File: "internal/conformance/data_handling_test.go", Kind: BindingCore},
			{Prefix: "internal/toon", File: "internal/toon/toon_test.go", Kind: BindingCore},
		},
		"internal/toon/toon_test.go": {{Prefix: "internal/toon", File: "internal/conformance/data_handling_test.go", Kind: BindingCore}},
		"internal/toonery/x.go":      nil,
	} {
		rules := BindingRules(path)
		files := BoundFiles(path)
		if !reflect.DeepEqual(rules, want) || len(files) != len(want) {
			t.Errorf("BindingRules(%s) = %v and BoundFiles = %v, want rules %v", path, rules, files, want)
			continue
		}
		for i, rule := range want {
			if files[i] != rule.File {
				t.Errorf("BoundFiles(%s) = %v, want the rules' files in order", path, files)
			}
		}
	}
}
