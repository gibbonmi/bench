// Tests for reference-file producer states, control-rune refusal, marker spans, and rename residue.
package skillsindex

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/capability"
)

// referenceSnapshot records what the reference is without ever reading a path the
// classifier would refuse. A FIFO opened here would block the test itself, which is
// the failure the production reader exists to avoid.
func referenceSnapshot(t *testing.T, path string) string {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		return "absent"
	}
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		target, err := os.Readlink(path)
		if err != nil {
			t.Fatal(err)
		}
		return "symlink " + target
	case info.Mode().IsRegular():
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return "regular " + string(data)
	}
	return "other " + info.Mode().Type().String()
}

// await runs one reference-reading call under a deadline. A reader that reopens the
// path directly blocks in open(2) on the FIFO row forever. A hung package test reports
// as a suite-wide timeout rather than as this row, so the bound is here.
func await[T any](t *testing.T, what string, call func() T) T {
	t.Helper()
	done := make(chan T, 1)
	go func() { done <- call() }()
	select {
	case got := <-done:
		return got
	case <-time.After(bounds.TestDeadline(0)):
		t.Fatalf("%s blocked on the reference, so it opened the path before classifying it", what)
	}
	panic("unreachable")
}

// TestReferenceProducerStatesStayDistinctAndBlockWrite is HI4: every hostile or
// degenerate reference disposition keeps its own diagnostic. Empty is not missing,
// and an untrustworthy object is neither. None of them lets Write touch the bytes it
// could not read.
func TestReferenceProducerStatesStayDistinctAndBlockWrite(t *testing.T) {
	for _, tc := range []struct {
		name    string
		plant   func(*testing.T, string)
		want    string
		refusal bool
	}{
		{name: "missing", want: referenceRel + " missing (skills index unverifiable)"},
		{
			name:  "empty",
			plant: func(t *testing.T, path string) { writeAt(t, path, "") },
			want:  referenceRel + " empty (skills index unverifiable)",
		},
		{
			name: "oversized",
			plant: func(t *testing.T, path string) {
				writeAt(t, path, strings.Repeat("a", int(bounds.ControlRecordLimit)+1))
			},
			refusal: true,
		},
		{
			name:    "invalid UTF-8",
			plant:   func(t *testing.T, path string) { writeAt(t, path, reference("")+"\xff\xfe") },
			refusal: true,
		},
		{
			name: "directory",
			plant: func(t *testing.T, path string) {
				if err := os.MkdirAll(path, 0o755); err != nil {
					t.Fatal(err)
				}
			},
			refusal: true,
		},
		{
			name: "fifo",
			plant: func(t *testing.T, path string) {
				if err := syscall.Mkfifo(path, 0o644); err != nil {
					capability.Capability(t, capability.Fifo, fmt.Sprintf("FIFOs unavailable on this filesystem: %v", err))
				}
			},
			refusal: true,
		},
		{
			name: "live symlink",
			plant: func(t *testing.T, path string) {
				target := filepath.Join(filepath.Dir(path), "elsewhere.md")
				writeAt(t, target, reference(""))
				if err := os.Symlink(target, path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
			refusal: true,
		},
		{
			name: "dangling symlink",
			plant: func(t *testing.T, path string) {
				if err := os.Symlink(filepath.Join(filepath.Dir(path), "gone.md"), path); err != nil {
					capability.Capability(t, capability.Symlink, fmt.Sprintf("symlinks unavailable on this filesystem: %v", err))
				}
			},
			refusal: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
			path := filepath.Join(root, filepath.FromSlash(referenceRel))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if tc.plant != nil {
				tc.plant(t, path)
			}
			before := referenceSnapshot(t, path)

			diags := await(t, "Check", func() []string { return Check(root) })
			err := await(t, "Write", func() error { return Write(context.Background(), root) })
			if err == nil {
				t.Fatal("Write accepted a reference it could not read")
			}
			if got := referenceSnapshot(t, path); got != before {
				t.Fatalf("Write changed the reference it refused: %.60q, want %.60q", got, before)
			}
			if tc.refusal {
				prefix := ReferenceRefusalPrefix()
				if !hasPrefixed(diags, prefix) {
					t.Fatalf("Check reported %q, want a diagnostic prefixed %q", diags, prefix)
				}
				if !strings.HasPrefix(err.Error(), prefix) || strings.TrimSpace(strings.TrimPrefix(err.Error(), prefix)) == "" {
					t.Fatalf("Write reported %q, want an attributed refusal prefixed %q", err, prefix)
				}
				return
			}
			if !hasExact(diags, tc.want) {
				t.Fatalf("Check reported %q, want %q", diags, tc.want)
			}
			if err.Error() != tc.want {
				t.Fatalf("Write reported %q, want %q", err, tc.want)
			}
		})
	}
}

// TestControlRunesNeverReachTheRenderedLine is the sink contract: the index line is
// line-structured markdown. Any control rune in a rendered field could split or forge
// an entry, and is refused before rendering. Ordinary graphic Unicode, accented Latin,
// CJK, emoji, arrows, still renders as exactly one line. The permitted half is
// asserted alongside the refused one because an ASCII-only fix passes the refusals and
// quietly loses the rest of the control partition.
func TestControlRunesNeverReachTheRenderedLine(t *testing.T) {
	for _, row := range []struct {
		name    string
		key     string
		value   string
		refused bool
	}{
		{"tab in the trigger", "index", "doing\tthings", true},
		{"carriage return forging a second line", "index", "safe → `x`\r- forged → `y`", true},
		{"escape in the trigger", "index", "doing\x1b[31m things", true},
		{"bell in the trigger", "index", "doing\a things", true},
		{"nul in the trigger", "index", "doing\x00 things", true},
		{"delete in the trigger", "index", "doing\x7f things", true},
		{"c1 next-line in the trigger", "index", "doing\u0085things", true},
		{"tab in the note", "index-note", "a\tnote", true},
		{"carriage return in the note", "index-note", "a note\r- forged", true},
		{"escape in the note", "index-note", "a\x1bnote", true},
		{"graphic unicode in the trigger", "index", "café 文書 🚀 →", false},
		{"graphic unicode in the note", "index-note", "註記 ✨ →", false},
	} {
		t.Run(row.name, func(t *testing.T) {
			root := t.TempDir()
			trigger, note := "doing probe things", "a probe note"
			if row.key == "index" {
				trigger = row.value
			} else {
				note = row.value
			}
			writeFile(t, root, ".agents/skills/probe/SKILL.md", "---\nindex: "+trigger+"\nindex-note: "+note+"\n---\n")
			original := reference("")
			writeFile(t, root, ".bench/BENCH-reference.md", original)

			block := renderedBlock(t, root)
			if !row.refused {
				want := "- " + trigger + " → `.agents/skills/probe/SKILL.md` + " + note
				if block != want {
					t.Fatalf("rendered block =\n%q\nwant\n%q", block, want)
				}
				return
			}
			if block != "" {
				t.Fatalf("refused field rendered a line: %q", block)
			}
			wantPrefix := ".agents/skills/probe/SKILL.md refused: "
			diags := Check(root)
			if len(diags) != 1 || !strings.HasPrefix(diags[0], wantPrefix) {
				t.Fatalf("Check = %v, want one diagnostic starting %q", diags, wantPrefix)
			}
			if err := Write(context.Background(), root); err == nil || !strings.HasPrefix(err.Error(), wantPrefix) {
				t.Fatalf("Write error = %v, want one starting %q", err, wantPrefix)
			}
			after, err := os.ReadFile(filepath.Join(root, ".bench", "BENCH-reference.md"))
			if err != nil {
				t.Fatal(err)
			}
			if string(after) != original {
				t.Fatalf("Write changed the reference:\n%q", string(after))
			}
		})
	}
}

func writeAt(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func hasExact(diags []string, want string) bool {
	for _, diag := range diags {
		if diag == want {
			return true
		}
	}
	return false
}

func hasPrefixed(diags []string, prefix string) bool {
	for _, diag := range diags {
		if strings.HasPrefix(diag, prefix) && strings.TrimSpace(strings.TrimPrefix(diag, prefix)) != "" {
			return true
		}
	}
	return false
}

// A newline cannot survive the fence scan's line split, so the sink table above cannot
// construct one from a file. The predicate is asserted directly, so a later parser
// that does yield multi-line values still finds the sink closed.
func TestControlRefusalCoversTheNewlineTheParserCannotYield(t *testing.T) {
	if controlRefusal("index", "safe\n- forged") == "" {
		t.Fatal("newline accepted in an index field")
	}
	if got := controlRefusal("index", "café 文書 🚀 →"); got != "" {
		t.Fatalf("graphic unicode refused: %s", got)
	}
}

// TestExactlyOneMarkerSpanIsRequired walks the marker cardinality table through both
// public consumers. A parser that stops at the first end marker accepts most of these
// rows, and each acceptance would let Write rewrite bytes it cannot place.
func TestExactlyOneMarkerSpanIsRequired(t *testing.T) {
	const (
		start = "<!-- bench:skills-index:start -->"
		end   = "<!-- bench:skills-index:end -->"
	)
	malformed := []struct {
		name string
		ref  string
	}{
		{name: "zero markers", ref: "# Reference\n\nno block here\n"},
		{name: "reversed order", ref: "# Reference\n\n" + end + "\n\n" + start + "\n"},
		{name: "unclosed start", ref: "# Reference\n\n" + start + "\n- a line\n"},
		{name: "duplicate start", ref: "# Reference\n\n" + start + "\n" + start + "\n" + end + "\n"},
		{name: "duplicate end", ref: "# Reference\n\n" + start + "\n" + end + "\n" + end + "\n"},
		{name: "two complete spans", ref: "# Reference\n\n" + start + "\n" + end + "\n" + start + "\n" + end + "\n"},
	}
	want := referenceRel + " skills-index markers missing (bench:skills-index)"
	for _, tc := range malformed {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, ".bench/BENCH-reference.md", tc.ref)
			if got := Check(root); len(got) != 1 || got[0] != want {
				// Errorf, not Fatalf: the byte-preservation clause below is a separate
				// acceptance and must still be exercised when Check accepts the span.
				t.Errorf("check = %v, want [%s]", got, want)
			}
			err := Write(context.Background(), root)
			if err == nil || err.Error() != want {
				t.Errorf("write = %v, want %s", err, want)
			}
			after, readErr := os.ReadFile(filepath.Join(root, filepath.FromSlash(referenceRel)))
			if readErr != nil {
				t.Fatal(readErr)
			}
			if string(after) != tc.ref {
				t.Fatalf("write changed bytes:\n%q\nwant\n%q", after, tc.ref)
			}
		})
	}

	t.Run("exactly one ordered span", func(t *testing.T) {
		root := t.TempDir()
		writeFile(t, root, ".bench/BENCH-reference.md", reference(""))
		if got := Check(root); len(got) != 0 {
			t.Fatalf("check = %v, want none", got)
		}
		if err := Write(context.Background(), root); err != nil {
			t.Fatal(err)
		}
	})
}

// siblingTemps names the replacement temps left beside the reference. It reads the
// directory rather than opening any child: a residue count must not itself become a
// reader of a path the classifier would refuse.
func siblingTemps(t *testing.T, root string) []string {
	t.Helper()
	dir, err := os.ReadDir(filepath.Join(root, ".bench"))
	if err != nil {
		t.Fatal(err)
	}
	var residue []string
	for _, child := range dir {
		if strings.HasPrefix(child.Name(), ".skills-index-") {
			residue = append(residue, child.Name())
		}
	}
	return residue
}

func TestRenameFailureLeavesNoResidueAndKeepsReferenceBytes(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".agents/skills/alpha/SKILL.md", "---\nname: alpha\nindex: doing alpha things\n---\n")
	refPath := filepath.Join(root, ".bench", "BENCH-reference.md")
	original := reference("")
	writeFile(t, root, ".bench/BENCH-reference.md", original)

	injected := fmt.Errorf("injected rename failure")
	restore := renameFile
	renameFile = func(string, string) error { return injected }
	defer func() { renameFile = restore }()

	if err := Write(context.Background(), root); err == nil || !strings.Contains(err.Error(), injected.Error()) {
		t.Fatalf("write with a failing rename = %v, want %v", err, injected)
	}
	kept, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(kept) != original {
		t.Fatalf("reference after a failed rename =\n%q\nwant\n%q", kept, original)
	}
	if residue := siblingTemps(t, root); len(residue) != 0 {
		t.Fatalf("failed rename left %v, want no .bench/.skills-index-* residue", residue)
	}
}
