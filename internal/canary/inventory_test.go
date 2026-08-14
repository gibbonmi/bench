package canary

import (
	"bytes"
	"os"
	"path/filepath"
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
	if got, want := stdout.String(), "canary inventory ok (191 fixture bindings)\n"; got != want {
		t.Fatalf("Run() stdout = %q, want %q", got, want)
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
