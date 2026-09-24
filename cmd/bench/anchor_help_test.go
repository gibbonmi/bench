package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/toon"
)

// anchorsFixtureRows are the AGENTS.md registry rows in registry order. The fixture file at
// testdata/anchors/fixture-repo/AGENTS.md plants each Require needle of that list and no
// other needle, so the fixture stays the only hand copy of a needle. That copy is the
// fixture-bite shape this package uses: a Require needle that a maintainer edits without
// updating the fixture stops the test here. It is cheaper to keep honest than a live-tree
// read that drifts silently.
func anchorsFixtureRows(t *testing.T) []anchors.Anchor {
	t.Helper()
	planted := map[string]bool{}
	for _, line := range strings.Split(readAnchorsFixture(t, "AGENTS.md"), "\n") {
		planted[line] = true
	}
	var rows []anchors.Anchor
	for _, anchor := range anchors.Entries() {
		if anchor.File != "AGENTS.md" {
			continue
		}
		if planted[anchor.Needle] != (anchor.Kind == anchors.Require) {
			t.Fatalf("fixture AGENTS.md planted=%t for the kind %d needle %q; the fixture plants each Require needle and no other", planted[anchor.Needle], anchor.Kind, anchor.Needle)
		}
		rows = append(rows, anchor)
	}
	return rows
}

// anchorsFixtureKindNames is the expected kind cell for each kind that the AGENTS.md rows
// use. The names are literals, not anchorKindName, so a swapped kind name turns the
// table tests red.
var anchorsFixtureKindNames = map[anchors.Kind]string{
	anchors.Require: "require",
	anchors.Forbid:  "forbid",
}

// anchorsFixtureTable is the expected anchors table for the AGENTS.md rows. Each row takes
// its kind cell from anchorsFixtureKindNames and its line from lines, which holds 0 for a
// row whose needle the file does not carry.
func anchorsFixtureTable(t *testing.T, rows []anchors.Anchor, lines []int) string {
	t.Helper()
	table := make([][]any, len(rows))
	for i, anchor := range rows {
		kind, ok := anchorsFixtureKindNames[anchor.Kind]
		if !ok {
			t.Fatalf("AGENTS.md row %q has kind %d, which anchorsFixtureKindNames does not name", anchor.Needle, anchor.Kind)
		}
		table[i] = []any{kind, anchor.Section, anchor.Step, anchor.Needle, lines[i]}
	}
	want, err := toon.TableTyped("anchors", []string{"kind", "section", "step", "needle", "line"}, table)
	if err != nil {
		t.Fatal(err)
	}
	return want
}

func TestAnchorsReportsNeedleLines(t *testing.T) {
	rows := anchorsFixtureRows(t)
	body := readAnchorsFixture(t, "AGENTS.md")
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), body)

	result := runAXICommandAt(t, root, []string{"anchors", "AGENTS.md"})
	if result.code != 0 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md = %#v, want exit 0 and no stderr", result)
	}
	// The fixture file places each Require needle on its own line, two lines apart,
	// starting at line 3 (line 1 is the title, line 2 is blank). A Forbid row reads 0.
	lines := make([]int, len(rows))
	next := 3
	for i, anchor := range rows {
		if anchor.Kind == anchors.Require {
			lines[i] = next
			next += 2
		}
	}
	want := anchorsFixtureTable(t, rows, lines)
	want += "help[0]{cmd,why}:\n"
	if result.stdout != want {
		t.Fatalf("anchors AGENTS.md stdout = %q, want %q", result.stdout, want)
	}
}

func TestAnchorsReportsAbsentNeedles(t *testing.T) {
	rows := anchorsFixtureRows(t)
	// Drop every second Require needle's line (and its blank separator) from the on-disk
	// fixture, so those rows must read 0 while their kept siblings — unmoved by the removal
	// — close up to lines 3, 5, and 7. A Forbid row reads 0 and raises no diagnostic.
	var drop []string
	lines := make([]int, len(rows))
	planted, next := 0, 3
	for i, anchor := range rows {
		if anchor.Kind != anchors.Require {
			continue
		}
		if planted%2 == 1 {
			drop = append(drop, anchor.Needle)
		} else {
			lines[i] = next
			next += 2
		}
		planted++
	}
	body := dropFixtureLines(readAnchorsFixture(t, "AGENTS.md"), drop...)
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), body)

	result := runAXICommandAt(t, root, []string{"anchors", "AGENTS.md"})
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md = %#v, want exit 1 and no stderr", result)
	}
	want := anchorsFixtureTable(t, rows, lines)
	for _, anchor := range rows {
		if slices.Contains(drop, anchor.Needle) {
			want += toon.Errorf("anchor", anchor.Diagnostic) + "\n"
		}
	}
	want += "help[0]{cmd,why}:\n"
	if result.stdout != want {
		t.Fatalf("anchors AGENTS.md stdout = %q, want %q", result.stdout, want)
	}

	// An absent file gives 0 on every row, siblings included, and one missing-file
	// diagnostic for each Require row.
	empty := newAXIEnvelopeRepo(t)
	result = runAXICommandAt(t, empty, []string{"anchors", "AGENTS.md"})
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md (absent file) = %#v, want exit 1 and no stderr", result)
	}
	want = anchorsFixtureTable(t, rows, make([]int, len(rows)))
	want += strings.Repeat(toon.Errorf("anchor", "acceptance coverage anchor file missing: AGENTS.md")+"\n", planted)
	want += "help[0]{cmd,why}:\n"
	if result.stdout != want {
		t.Fatalf("anchors AGENTS.md (absent file) stdout = %q, want %q", result.stdout, want)
	}
}

func TestAnchorsReportsForbiddenNeedles(t *testing.T) {
	for _, anchor := range anchors.Entries() {
		if anchor.Kind != anchors.Forbid {
			continue
		}
		root := newAXIEnvelopeRepo(t)
		writeAXIFixture(t, filepath.Join(root, filepath.FromSlash(anchor.File)), anchor.Needle+"\n")
		result := runAXICommandAt(t, root, []string{"anchors", anchor.File})
		if result.code != 1 || result.stderr != "" || !strings.Contains(result.stdout, toon.Errorf("anchor", anchor.Diagnostic)+"\n") {
			t.Fatalf("forbidden anchor = %#v, want exit 1 and %q", result, anchor.Diagnostic)
		}
		return
	}
	t.Fatal("registry has no forbidden anchor")
}

func TestAnchorsReportsCaseFoldedEmphasisViolation(t *testing.T) {
	var target anchors.Anchor
	for _, anchor := range anchors.Entries() {
		if anchor.Kind == anchors.ForbidCaseFoldedEmphasis {
			target = anchor
			break
		}
	}
	if target.File == "" {
		t.Fatal("registry has no case-folded emphasis anchor")
	}
	root := newAXIEnvelopeRepo(t)
	body := "<!-- hidden violation -->\nKELVIN before the match\nAn EXECUTABLE **RED** IS MANDATORY before specification.\n"
	writeAXIFixture(t, filepath.Join(root, filepath.FromSlash(target.File)), body)
	result := runAXICommandAt(t, root, []string{"anchors", target.File})
	if result.code != 1 || result.stderr != "" {
		t.Fatalf("anchors emphasized violation = %#v, want exit 1 and no stderr", result)
	}
	if !strings.Contains(result.stdout, toon.Errorf("anchor", target.Diagnostic)+"\n") {
		t.Fatalf("anchors emphasized violation stdout = %q, want diagnostic %q", result.stdout, target.Diagnostic)
	}
	want, err := toon.TableTyped("anchors", []string{"kind", "section", "step", "needle", "line"}, [][]any{{"forbid-case-folded-emphasis", "", 0, target.Needle, 3}})
	if err != nil {
		t.Fatal(err)
	}
	wantRow := strings.Split(want, "\n")[1]
	if !strings.Contains(result.stdout, wantRow+"\n") {
		t.Fatalf("anchors emphasized violation stdout = %q, want row %q", result.stdout, wantRow)
	}
}

func TestAnchorsLeavesUnregisteredPathEmpty(t *testing.T) {
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "unregistered.md"), "# An ordinary file\n")
	result := runAXICommandAt(t, root, []string{"anchors", "unregistered.md"})
	if result.code != 0 || result.stderr != "" || result.stdout != "anchors[0]{kind,section,step,needle,line}:\nhelp[0]{cmd,why}:\n" {
		t.Fatalf("unregistered path = %#v, want an empty successful query", result)
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
