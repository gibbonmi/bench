package canary

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
)

func TestRunReportsAcceptedInventoryBindings(t *testing.T) {
	working, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := Run([]string{filepath.Join(working, "../..")}, &stdout, &stderr); code != 0 {
		t.Fatalf("Run() = %d, stderr=%q", code, stderr.String())
	}
	if got := stdout.String(); !strings.HasPrefix(got, "canary inventory ok (") || !strings.HasSuffix(got, " fixture bindings)\n") {
		t.Fatalf("Run() stdout = %q, want accepted inventory schema", got)
	}
}

func TestFixturesRefusesEmptyInventoryWithInventoryOnlyDiagnostic(t *testing.T) {
	_, err := Fixtures(t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "canary fixture inventory is empty") {
		t.Fatalf("Fixtures(empty) error = %v, want inventory-only empty diagnostic", err)
	}
}

func TestFixturesClassifiesExpectMetadataBeforeReading(t *testing.T) {
	for _, test := range []struct {
		name string
		make func(t *testing.T, dir string)
	}{
		{name: "dangling", make: func(t *testing.T, dir string) {
			if err := os.Symlink(filepath.Join(dir, "gone"), filepath.Join(dir, "EXPECT")); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "fifo", make: func(t *testing.T, dir string) {
			if err := syscall.Mkfifo(filepath.Join(dir, "EXPECT"), 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			fixture := filepath.Join(root, "family")
			if err := os.Mkdir(fixture, 0o755); err != nil {
				t.Fatal(err)
			}
			test.make(t, fixture)
			if _, err := Fixtures(root); err == nil || !strings.Contains(err.Error(), "marker") {
				t.Fatalf("Fixtures(%s) error = %v, want metadata refusal", test.name, err)
			}
		})
	}
}

func TestFixturesTreatsMissingExpectAsFamily(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "family", "fixture")
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := Fixtures(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got["fixture"]; !ok {
		t.Fatalf("Fixtures(missing EXPECT) = %#v, want family child inventory", got)
	}
}

func TestUnboundConformanceFamiliesReportsInvalidExpectMetadata(t *testing.T) {
	root := t.TempDir()
	family := filepath.Join(root, "tests", "canary", "family")
	if err := os.MkdirAll(family, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(family, "EXPECT"), 0o600); err != nil {
		t.Fatal(err)
	}
	diagnostics := UnboundConformanceFamilies(root)
	if len(diagnostics) != 1 || !strings.Contains(diagnostics[0], "marker") {
		t.Fatalf("UnboundConformanceFamilies(FIFO) = %v, want metadata diagnostic", diagnostics)
	}
}

// TestFixturePinsEnumeratesLiveInventory covers TG15. A fixture planted in a
// temporary root appears in the pin enumeration through every producer the
// materializer reads: the BASE list, an @ include, the files/ overlay, and
// MUTATE.json. The enumeration derives from the live inventory, so a copied list
// cannot answer for it.
func TestFixturePinsEnumeratesLiveInventory(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "example-family", "planted")
	if err := os.MkdirAll(filepath.Join(fixture, "files", "internal", "example"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFixtureFile(t, filepath.Join(fixture, "EXPECT"), "planted diagnostic\n")
	writeFixtureFile(t, filepath.Join(fixture, "BASE"), "internal/example/base.go\n@base-include.list\n")
	writeFixtureFile(t, filepath.Join(root, "base-include.list"), "internal/example/included.go\n")
	writeFixtureFile(t, filepath.Join(fixture, "files", "internal", "example", "overlay.go"), "package example\n")
	writeFixtureFile(t, filepath.Join(fixture, "MUTATE.json"), `[{"path":"internal/example/mutated.go","old":"a","new":"b"}]`)

	pins, err := FixturePins(root)
	if err != nil {
		t.Fatalf("FixturePins: %v", err)
	}
	const want = "tests/canary/example-family/planted"
	for _, path := range []string{
		"internal/example/base.go",
		"internal/example/included.go",
		"internal/example/overlay.go",
		"internal/example/mutated.go",
	} {
		got := pins[path]
		if len(got) != 1 || got[0] != want {
			t.Errorf("FixturePins()[%q] = %v, want [%q]", path, got, want)
		}
	}
	if got := pins["internal/example/unpinned.go"]; len(got) != 0 {
		t.Errorf("FixturePins()[unpinned] = %v, want no pin", got)
	}
}

// TestFixturePinsAnswersEmptyForAnInventorylessRoot keeps an absent tests/canary
// an answer rather than a fault. A repository with no fixture inventory pins
// nothing, and every caller reads that as an empty closure.
func TestFixturePinsAnswersEmptyForAnInventorylessRoot(t *testing.T) {
	pins, err := FixturePins(t.TempDir())
	if err != nil || len(pins) != 0 {
		t.Fatalf("FixturePins(inventoryless) = %v, %v; want an empty map and no error", pins, err)
	}
}

// TestPinsNameEachProducerLocation grades the typed pins under FixturePins: each producer
// names the file and line that pin its path, and an included list names itself rather
// than the BASE file that includes it.
func TestPinsNameEachProducerLocation(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "example-family", "planted")
	writeFixtureFile(t, filepath.Join(fixture, "EXPECT"), "planted diagnostic\n")
	writeFixtureFile(t, filepath.Join(fixture, "BASE"), "internal/example/base.go\n@base-include.list\n")
	writeFixtureFile(t, filepath.Join(root, "base-include.list"), "internal/example/first.go\ninternal/example/included.go\n")
	writeFixtureFile(t, filepath.Join(fixture, "files", "internal", "example", "overlay.go"), "package example\n")
	writeFixtureFile(t, filepath.Join(fixture, "MUTATE.json"), "[\n  {\"path\":\"internal/example/base.go\",\"old\":\"a\",\"new\":\"b\"},\n  {\"path\":\"internal/example/mutated.go\",\"old\":\"a\",\"new\":\"b\"}\n]\n")
	pins, err := Pins(root)
	if err != nil {
		t.Fatalf("Pins: %v", err)
	}
	const dir = "tests/canary/example-family/planted"
	want := []Pin{
		{"internal/example/base.go", dir, dir + "/BASE", 1, PinBase},
		{"internal/example/first.go", dir, "base-include.list", 1, PinInclude},
		{"internal/example/included.go", dir, "base-include.list", 2, PinInclude},
		{"internal/example/overlay.go", dir, dir + "/files/internal/example/overlay.go", 0, PinOverlay},
		{"internal/example/base.go", dir, dir + "/MUTATE.json", 2, PinMutation},
		{"internal/example/mutated.go", dir, dir + "/MUTATE.json", 3, PinMutation},
	}
	if !slices.Equal(pins, want) {
		t.Fatalf("Pins = %v, want %v", pins, want)
	}
}

// TestPinsRefuseASpecialInventory keeps a present inventory that is no readable directory
// a fault. Only absence answers empty, so a link, a file, or an unreadable directory at
// tests/canary cannot read as a repository with nothing pinned.
func TestPinsRefuseASpecialInventory(t *testing.T) {
	for name, place := range map[string]func(t *testing.T, root string){
		"file": func(t *testing.T, root string) { writeFixtureFile(t, Dir(root), "not a directory\n") },
		"unreadable": func(t *testing.T, root string) {
			if os.Geteuid() == 0 {
				t.Skip("root reads a mode 000 directory")
			}
			writeFixtureFile(t, filepath.Join(Dir(root), "family", "fixture", "BASE"), "internal/example/base.go\n")
			if err := os.Chmod(Dir(root), 0); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = os.Chmod(Dir(root), 0o755) })
		},
		"link": func(t *testing.T, root string) {
			writeFixtureFile(t, filepath.Join(root, "elsewhere", "family", "fixture", "BASE"), "internal/example/base.go\n")
			if err := os.MkdirAll(filepath.Dir(Dir(root)), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(filepath.Join(root, "elsewhere"), Dir(root)); err != nil {
				t.Fatal(err)
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			place(t, root)
			if pins, err := FixturePins(root); err == nil || !strings.Contains(err.Error(), "canary fixture inventory") {
				t.Fatalf("FixturePins(%s) = %v, %v; want the inventory refusal", name, pins, err)
			}
		})
	}
}

// TestFixturePinsReturnsAnInventoryError keeps a malformed inventory a fault: only
// an absent inventory answers empty.
func TestFixturePinsReturnsAnInventoryError(t *testing.T) {
	root := t.TempDir()
	for _, family := range []string{"first-family", "second-family"} {
		writeFixtureFile(t, filepath.Join(Dir(root), family, "twin", "EXPECT"), "twin diagnostic\n")
	}
	if _, err := FixturePins(root); err == nil || !strings.Contains(err.Error(), "appears in multiple families") {
		t.Fatalf("FixturePins(duplicate name) error = %v, want the inventory diagnostic", err)
	}
}

// TestPinnedPathsReadsOneFixtureInBaseThenMutationOrder holds the per-fixture
// accessor to the order the materializer applies: the BASE list first, then the
// MUTATE.json anchors. A caller that filters diagnostics by pin reads this
// order, so a reversal changes which line the reader sees first.
func TestPinnedPathsReadsOneFixtureInBaseThenMutationOrder(t *testing.T) {
	root := t.TempDir()
	fixture := filepath.Join(root, "tests", "canary", "example-family", "planted")
	writeFixtureFile(t, filepath.Join(fixture, "EXPECT"), "planted diagnostic\n")
	writeFixtureFile(t, filepath.Join(fixture, "BASE"), "internal/example/base.go\n")
	writeFixtureFile(t, filepath.Join(fixture, "MUTATE.json"), `[{"path":"internal/example/mutated.go","old":"a","new":"b"}]`)

	paths, err := PinnedPaths(root, fixture)
	if err != nil {
		t.Fatalf("PinnedPaths: %v", err)
	}
	want := []string{"internal/example/base.go", "internal/example/mutated.go"}
	if !slices.Equal(paths, want) {
		t.Fatalf("PinnedPaths = %v, want %v", paths, want)
	}
}

func writeFixtureFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
