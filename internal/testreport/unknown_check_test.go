package testreport

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/freshness"
	"github.com/gibbonmi/bench/internal/treetarget/kittest"
)

// The refusal lines here are authored apart from unknownCheck. These rows grade the line
// order, the `unknown` and `unsealed` words, the seal's source digest, and the escape of a
// control character. A missing or reordered line, the executable digest in the seal line,
// or a raw control byte reds these rows.

// useRunningExecutable makes path the running executable for the rest of the test.
func useRunningExecutable(t *testing.T, path string) {
	t.Helper()
	answerRunningExecutable(t, func() (string, error) { return path, nil })
}

// answerRunningExecutable makes answer the running-executable lookup for the rest of the test.
func answerRunningExecutable(t *testing.T, answer func() (string, error)) {
	t.Helper()
	previous := runningExecutable
	runningExecutable = answer
	t.Cleanup(func() { runningExecutable = previous })
}

// checkInventory is the check list that ends each refusal.
func checkInventory() string {
	return "checks:\n  " + strings.Join(namedChecks(), "\n  ") + "\n"
}

// unsealedExecutable writes a temporary running executable with no seal file beside it.
func unsealedExecutable(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "bench")
	if err := os.WriteFile(path, []byte("Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestUnknownNamedCheckReportsOperandAndInventory(t *testing.T) {
	executable := unsealedExecutable(t)
	useRunningExecutable(t, executable)
	inventory := checkInventory()
	for _, unknown := range []string{"not-registered", "release-evidence-probe"} {
		output, code := Command(t.TempDir(), []string{"--check", unknown})
		if code != 2 {
			t.Errorf("unknown named check %q = %d, want 2\n%s", unknown, code, output)
			continue
		}
		want := "unknown check: " + unknown + "\nexecutable: " + executable + "\nseal: unsealed\n" + inventory
		if output != want {
			t.Errorf("unknown named check %q = %q, want %q", unknown, output, want)
		}
	}
}

// TestUnknownCheckNamesSealSources publishes a marker executable over a kit tree, so the
// seal beside it is a real seal and its sources value is the digest of that tree.
func TestUnknownCheckNamesSealSources(t *testing.T) {
	root := t.TempDir()
	kittest.WriteTree(t, root)
	staged := filepath.Join(t.TempDir(), "staged")
	if err := os.WriteFile(staged, []byte("Bench executable"), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := freshness.PublishedExecutable(root)
	if err := freshness.Publish(root, staged, executable, filepath.Dir(executable), "unknown-check"); err != nil {
		t.Fatal(err)
	}
	sources, err := freshness.Digest(root)
	if err != nil {
		t.Fatal(err)
	}
	useRunningExecutable(t, executable)

	output, code := Command(t.TempDir(), []string{"--check", "not-registered"})
	if want := "\nseal: " + sources + "\n"; code != 2 || !strings.Contains(output, want) {
		t.Fatalf("sealed running executable = %d, %q; want exit 2 and %q", code, output, want)
	}
}

func TestUnknownCheckPrintsUnsealed(t *testing.T) {
	for name, malformed := range map[string]bool{"no seal file": false, "malformed seal file": true} {
		t.Run(name, func(t *testing.T) {
			executable := unsealedExecutable(t)
			if malformed {
				if err := os.WriteFile(executable+".seal", []byte("{}"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			useRunningExecutable(t, executable)
			output, code := Command(t.TempDir(), []string{"--check", "not-registered"})
			if code != 2 || !strings.Contains(output, "\nseal: unsealed\n") {
				t.Fatalf("unreadable seal = %d, %q; want exit 2 and seal: unsealed", code, output)
			}
		})
	}
}

func TestUnknownCheckNamesUnknownExecutable(t *testing.T) {
	answerRunningExecutable(t, func() (string, error) { return "", os.ErrNotExist })
	output, code := Command(t.TempDir(), []string{"--check", "not-registered"})
	want := "unknown check: not-registered\nexecutable: unknown\nseal: unsealed\n" + checkInventory()
	if code != 2 || output != want {
		t.Fatalf("unnamed running executable = %d, %q; want exit 2 and %q", code, output, want)
	}
}

func TestUnknownCheckEscapesName(t *testing.T) {
	useRunningExecutable(t, "/tmp/bench\x01exe")
	output, code := Command(t.TempDir(), []string{"--check", "not\x01registered"})
	if code != 2 || strings.ContainsRune(output, '\x01') {
		t.Fatalf("control-bearing refusal = %d, %q; want exit 2 and no raw U+0001", code, output)
	}
	for _, want := range []string{"unknown check: not\\u0001registered\n", "\nexecutable: /tmp/bench\\u0001exe\n"} {
		if !strings.Contains(output, want) {
			t.Errorf("control-bearing refusal = %q, want %q", output, want)
		}
	}
}
