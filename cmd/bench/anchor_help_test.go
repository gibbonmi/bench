package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/toon"
)

// The six real AGENTS.md registry needles, in registry order. These tests key an
// isolated fixture repository's own AGENTS.md instead of reading the live tree, so a
// prose edit to the real AGENTS.md cannot move these lines out from under the test.
var anchorsFixtureNeedles = []string{
	"never by polling a self-matching pattern",
	"runs plan-before-apply: print the exact target list, sample it, then apply",
	"repository-wide sweep uses `rg --hidden`",
	"Discover Bench verbs non-interactively",
	"Run the system suite by hand through `bench test --check system`.",
	"`bench handoff` rewrites only the calling worktree's assignment section.",
}

func TestAnchorsReportsNeedleLines(t *testing.T) {
	var body strings.Builder
	body.WriteString("# AGENTS fixture\n")
	wantLines := make([]int, len(anchorsFixtureNeedles))
	line := 1
	for i, needle := range anchorsFixtureNeedles {
		body.WriteString("\n")
		line += 2
		body.WriteString(needle + "\n")
		wantLines[i] = line
	}
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), body.String())

	result := runAXICommandAt(t, root, []string{"anchors", "AGENTS.md"})
	if result.code != 0 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md = %#v, want exit 0 and no stderr", result)
	}
	rows := make([][]any, len(anchorsFixtureNeedles))
	for i, needle := range anchorsFixtureNeedles {
		rows[i] = []any{"require", "", needle, wantLines[i]}
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
	// The fixture keeps the first, third, and fifth needle and drops the rest, so the
	// dropped rows must read 0 while their siblings keep their real lines.
	body := "# AGENTS fixture\n\n" +
		anchorsFixtureNeedles[0] + "\n\n" +
		anchorsFixtureNeedles[2] + "\n\n" +
		anchorsFixtureNeedles[4] + "\n"
	root := newAXIEnvelopeRepo(t)
	writeAXIFixture(t, filepath.Join(root, "AGENTS.md"), body)

	result := runAXICommandAt(t, root, []string{"anchors", "AGENTS.md"})
	if result.code != 0 || result.stderr != "" {
		t.Fatalf("anchors AGENTS.md = %#v, want exit 0 and no stderr", result)
	}
	want, err := toon.TableTyped("anchors", []string{"kind", "section", "needle", "line"}, [][]any{
		{"require", "", anchorsFixtureNeedles[0], 3},
		{"require", "", anchorsFixtureNeedles[1], 0},
		{"require", "", anchorsFixtureNeedles[2], 5},
		{"require", "", anchorsFixtureNeedles[3], 0},
		{"require", "", anchorsFixtureNeedles[4], 7},
		{"require", "", anchorsFixtureNeedles[5], 0},
	})
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
	rows := make([][]any, len(anchorsFixtureNeedles))
	for i, needle := range anchorsFixtureNeedles {
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
