package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/toon"
)

// anchorsFixtureNeedles are the AGENTS.md registry needles that the fixture file at
// testdata/anchors/fixture-repo/AGENTS.md carries, in registry order. The list is the
// intersection of the registry and the fixture, so the fixture stays the only hand copy of
// a needle. That copy is the fixture-bite shape this package uses: a registry needle a
// maintainer edits without updating the fixture leaves the registry row and the fixture
// line unmatched, so the derived list is shorter than the rows the command prints and
// TestAnchorsReportsNeedleLines reds on purpose. It is cheaper to keep honest than a
// live-tree read that drifts silently.
func anchorsFixtureNeedles(t *testing.T) []string {
	t.Helper()
	planted := map[string]bool{}
	for _, line := range strings.Split(readAnchorsFixture(t, "AGENTS.md"), "\n") {
		planted[line] = true
	}
	var needles []string
	for _, anchor := range anchors.Entries() {
		if anchor.File == "AGENTS.md" && planted[anchor.Needle] {
			needles = append(needles, anchor.Needle)
		}
	}
	return needles
}

func TestAnchorsReportsNeedleLines(t *testing.T) {
	needles := anchorsFixtureNeedles(t)
	body := readAnchorsFixture(t, "AGENTS.md")
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), body)

	result := runAXICommandAt(t, root, []string{"anchors", "AGENTS.md"})
	if result.code != 0 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md = %#v, want exit 0 and no stderr", result)
	}
	// The fixture file places each needle on its own line, two lines apart, starting at
	// line 3 (line 1 is the title, line 2 is blank).
	rows := make([][]any, len(needles))
	for i, needle := range needles {
		rows[i] = []any{"require", "", needle, 3 + 2*i}
	}
	want, err := toon.TableTyped("anchors", []string{"kind", "section", "needle", "line"}, rows)
	if err != nil {
		t.Fatal(err)
	}
	want += "help[0]{cmd,why}:\n"
	if result.stdout != want {
		t.Fatalf("anchors AGENTS.md stdout = %q, want %q", result.stdout, want)
	}
}

func TestAnchorsReportsAbsentNeedles(t *testing.T) {
	needles := anchorsFixtureNeedles(t)
	// Drop every second needle's line (and its blank separator) from the on-disk fixture,
	// so those rows must read 0 while their kept siblings — unmoved by the removal — close
	// up to lines 3, 5, and 7.
	var drop []string
	for i, needle := range needles {
		if i%2 == 1 {
			drop = append(drop, needle)
		}
	}
	body := dropFixtureLines(readAnchorsFixture(t, "AGENTS.md"), drop...)
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), body)

	result := runAXICommandAt(t, root, []string{"anchors", "AGENTS.md"})
	if result.code != 0 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md = %#v, want exit 0 and no stderr", result)
	}
	kept := make([][]any, len(needles))
	line := 3
	for i, needle := range needles {
		if i%2 == 1 {
			kept[i] = []any{"require", "", needle, 0}
			continue
		}
		kept[i] = []any{"require", "", needle, line}
		line += 2
	}
	want, err := toon.TableTyped("anchors", []string{"kind", "section", "needle", "line"}, kept)
	if err != nil {
		t.Fatal(err)
	}
	want += "help[0]{cmd,why}:\n"
	if result.stdout != want {
		t.Fatalf("anchors AGENTS.md stdout = %q, want %q", result.stdout, want)
	}

	// An absent file gives 0 on every row, siblings included.
	empty := newAXIEnvelopeRepo(t)
	result = runAXICommandAt(t, empty, []string{"anchors", "AGENTS.md"})
	if result.code != 0 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md (absent file) = %#v, want exit 0 and no stderr", result)
	}
	rows := make([][]any, len(needles))
	for i, needle := range needles {
		rows[i] = []any{"require", "", needle, 0}
	}
	want, err = toon.TableTyped("anchors", []string{"kind", "section", "needle", "line"}, rows)
	if err != nil {
		t.Fatal(err)
	}
	want += "help[0]{cmd,why}:\n"
	if result.stdout != want {
		t.Fatalf("anchors AGENTS.md (absent file) stdout = %q, want %q", result.stdout, want)
	}
}

func TestAnchorsRefusesAnUnreadableFile(t *testing.T) {
	t.Run("symlink", func(t *testing.T) {
		root := newAXIEnvelopeRepo(t)
		target := filepath.Join(root, "target.md")
		writeAXIFixture(t, target, "# target\n")
		link := filepath.Join(root, "link.md")
		requireAnchorSymlink(t, target, link)
		result := runAXICommandAt(t, root, []string{"anchors", "link.md"})
		if result.code != 1 || result.stderr != "" {
			t.Fatalf("anchors link.md = %#v, want exit 1 and no stderr", result)
		}
		if !strings.HasPrefix(result.stdout, "error:") || !strings.Contains(result.stdout, "link.md is wrong-type") {
			t.Fatalf("anchors link.md stdout = %q, want a structured wrong-type refusal", result.stdout)
		}
	})

	t.Run("fifo", func(t *testing.T) {
		root := newAXIEnvelopeRepo(t)
		requireAnchorFifo(t, filepath.Join(root, "pipe.md"))
		result := runAXICommandAt(t, root, []string{"anchors", "pipe.md"})
		if result.code != 1 || result.stderr != "" {
			t.Fatalf("anchors pipe.md = %#v, want exit 1 and no stderr", result)
		}
		if !strings.HasPrefix(result.stdout, "error:") || !strings.Contains(result.stdout, "pipe.md is wrong-type") {
			t.Fatalf("anchors pipe.md stdout = %q, want a structured wrong-type refusal", result.stdout)
		}
	})

	t.Run("unreadable", func(t *testing.T) {
		root := newAXIEnvelopeRepo(t)
		path := filepath.Join(root, "denied.md")
		writeAXIFixture(t, path, "secret\n")
		requireAnchorUnreadable(t, path)
		result := runAXICommandAt(t, root, []string{"anchors", "denied.md"})
		if result.code != 1 || result.stderr != "" {
			t.Fatalf("anchors denied.md = %#v, want exit 1 and no stderr", result)
		}
		if !strings.HasPrefix(result.stdout, "error:") || !strings.Contains(result.stdout, "denied.md is unreadable") {
			t.Fatalf("anchors denied.md stdout = %q, want a structured unreadable refusal", result.stdout)
		}
	})
}

// TestAnchorsUsageRefusesAMissingArgument keeps the argument-count refusal pinned
// directly against anchorsCommand, independent of the isolated-fixture tests above.
func TestAnchorsUsageRefusesAMissingArgument(t *testing.T) {
	got, code := anchorsCommand(nil)
	if code != 2 || got != "usage: bench anchors (missing argument: argument)\n" {
		t.Fatalf("anchorsCommand(nil) = (%q, %d), want unchanged usage refusal", got, code)
	}
}

// readAnchorsFixture reads name from the isolated fixture repository under
// testdata/anchors/fixture-repo, so the anchors query is graded against fixture content
// instead of the live tree.
func readAnchorsFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "anchors", "fixture-repo", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// dropFixtureLines removes each line in body that exactly matches one of drop, along
// with the blank line immediately before it, so every line that survives keeps its
// original line number.
func dropFixtureLines(body string, drop ...string) string {
	remove := make(map[string]bool, len(drop))
	for _, needle := range drop {
		remove[needle] = true
	}
	lines := strings.Split(body, "\n")
	kept := make([]string, 0, len(lines))
	for _, line := range lines {
		if remove[line] {
			if len(kept) > 0 && kept[len(kept)-1] == "" {
				kept = kept[:len(kept)-1]
			}
			continue
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n")
}

func requireAnchorSymlink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err != nil {
		capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
	}
}

func requireAnchorFifo(t *testing.T, path string) {
	t.Helper()
	if err := syscall.Mkfifo(path, 0o644); err != nil {
		capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
	}
}

// requireAnchorUnreadable strips a fixture's permissions and proves the strip bit, since
// root ignores the mode entirely and would otherwise read the file straight through the
// assertion.
func requireAnchorUnreadable(t *testing.T, path string) {
	t.Helper()
	t.Cleanup(func() { _ = os.Chmod(path, 0o600) })
	if err := os.Chmod(path, 0o000); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip permissions: %v", err))
	}
	f, err := os.Open(path)
	if err == nil {
		f.Close()
		capability.Capability(t, capability.Privilege, "mode 0o000 is still readable by this user")
	}
}
