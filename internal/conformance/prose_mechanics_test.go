package conformance

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/bounds"
	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/prose"
)

// checkProseMechanics grades prose mechanics and learning-journal structure.
func checkProseMechanics(root string) []string {
	diags := prose.Grade(root)
	return append(diags, learningJournalDiagnostics(root)...)
}

func learningJournalDiagnostics(root string) []string {
	path := filepath.Join(root, filepath.FromSlash(learnings.JournalPath))
	c := bounds.ClassifyNoFollow(path)
	switch c.State {
	case bounds.StateAbsent, bounds.StateEmpty:
		return nil
	case bounds.StateParsed:
		if reason := learnings.UnsupportedSchemaReason(c.Data); reason != "" {
			return []string{fmt.Sprintf("learning journal: %q: refused %s: %s", learnings.JournalPath, bounds.StateUnsupportedSchema, reason)}
		}
		_, malformed := learnings.Parse(c.Data)
		diags := make([]string, 0, len(malformed))
		for _, entry := range malformed {
			diags = append(diags, fmt.Sprintf("learning journal: %q line %d: %s", learnings.JournalPath, entry.Line, entry.Reason))
		}
		return diags
	default:
		return []string{fmt.Sprintf("learning journal: %q: refused %s: %s", learnings.JournalPath, c.State, c.Reason)}
	}
}

// approvedProseExclusionRows is the reviewed set of paths .bench/prose-exclusions may
// name. It holds exactly the four permanent rows. The set stays independent of the
// file it grades: a set derived from the file would accept every row somebody adds.
// A later added row needs a visible edit here, or this test reds.
var approvedProseExclusionRows = map[string]bool{
	"tests/canary/":    true,
	"docs/audits/":     true,
	"CHANGELOG.md":     true,
	"capture/IDEAS.md": true,
}

// TestProseMechanicsHoldsOnTheLiveTree is the check's live-tree assertion and the
// delegates' focused seam: a migration batch runs this one test to learn whether the
// files it rewrote clear the bounds.
func TestProseMechanicsHoldsOnTheLiveTree(t *testing.T) {
	h := NewHarness(t)
	if diags := checkProseMechanics(h.KitRoot); len(diags) != 0 {
		t.Fatalf("the kit's authored Markdown is over the prose mechanics bounds:\n%s", strings.Join(diags, "\n"))
	}
}

func TestProseMechanicsGradesLearningJournal(t *testing.T) {
	owner, bound := conformanceChecks["prose-mechanics"]
	if !bound {
		t.Fatal("prose-mechanics conformance owner is not bound")
	}
	for _, tc := range []struct {
		name, journal, want  string
		directory, oversized bool
	}{
		{name: "absent"},
		{name: "empty", journal: ""},
		{name: "valid", journal: "## 2026-08-24 — valid entry [open]\n"},
		{name: "unsupported schema", journal: "not a learnings journal\n", want: "learning journal: \"capture/learnings.md\": refused unsupported-schema: no dated heading found"},
		{name: "malformed", journal: "## malformed heading\n", want: "capture/learnings.md\" line 1: malformed learning heading"},
		{name: "unaccounted content", journal: learnings.JournalSchemaHeading + "\n\n" + learnings.JournalEntriesMarker + "\n\norphaned content\n", want: "capture/learnings.md\" line 5: learning content below the entries marker is not an entry"},
		{name: "invalid UTF-8", journal: learnings.JournalSchemaHeading + "\xff\n", want: "learning journal: \"capture/learnings.md\": refused malformed: invalid UTF-8"},
		{name: "oversized", oversized: true, want: "learning journal: \"capture/learnings.md\": refused unreadable: read limit exceeded"},
		{name: "wrong type", directory: true, want: "learning journal: \"capture/learnings.md\": refused wrong-type:"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			journal := filepath.Join(root, "capture", "learnings.md")
			if err := os.MkdirAll(filepath.Dir(journal), 0o755); err != nil {
				t.Fatal(err)
			}
			if tc.directory {
				if err := os.Mkdir(journal, 0o755); err != nil {
					t.Fatal(err)
				}
			} else if tc.name != "absent" {
				content := []byte(tc.journal)
				if tc.oversized {
					content = bytes.Repeat([]byte("a"), int(bounds.ControlRecordLimit)+1)
				}
				if err := os.WriteFile(journal, content, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if err := os.MkdirAll(filepath.Join(root, ".bench"), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(root, ".bench", "prose-exclusions"), nil, 0o644); err != nil {
				t.Fatal(err)
			}
			diags := owner.run(root, root, registry.Dev)
			if tc.want == "" {
				if len(diags) != 0 {
					t.Fatalf("clean learning journal produced %q", diags)
				}
				return
			}
			if !containsDiagnostic(diags, tc.want) {
				t.Fatalf("learning journal produced %q, want %q", diags, tc.want)
			}
		})
	}
}

// TestProseExclusionRowsStayInsideTheApprovedSet holds the no-new-row rule. The
// approved set above is independent of the file, so an added row reds here even when
// the row itself is well formed and its path exists.
func TestProseExclusionRowsStayInsideTheApprovedSet(t *testing.T) {
	h := NewHarness(t)
	rows := readProseExclusionRows(t, h.KitRoot)
	if len(rows) == 0 {
		t.Fatalf("%s names no row", prose.ExclusionFile)
	}
	var unapproved []string
	for _, row := range rows {
		if !approvedProseExclusionRows[row] {
			unapproved = append(unapproved, row)
		}
	}
	sort.Strings(unapproved)
	if len(unapproved) != 0 {
		t.Fatalf("%s names rows outside the approved set: %s", prose.ExclusionFile, strings.Join(unapproved, ", "))
	}
}

// readProseExclusionRows returns the subject of every row under root. The grammar lives
// in internal/prose, so this test reads the rows through that package rather than
// parsing the file a second time.
func readProseExclusionRows(t *testing.T, root string) []string {
	t.Helper()
	rows, err := prose.ExclusionRows(root)
	if err != nil {
		t.Fatalf("read %s: %v", prose.ExclusionFile, err)
	}
	return rows
}

// TestProseMechanicsCanaryFixturesBite runs every fixture in the family through the
// registered owner. Each fixture plants one `*.md` subject, so the restore that removes
// it takes the whole grade away and the red belongs to the mutation rather than to
// ambient state in the materialized tree.
func TestProseMechanicsCanaryFixturesBite(t *testing.T) {
	h := NewHarness(t)
	owner, bound := conformanceChecks["prose-mechanics"]
	if !bound {
		t.Fatal("prose-mechanics conformance owner is not bound")
	}
	familyDir := h.KitPath("tests", "canary", "prose-mechanics")
	entries, err := os.ReadDir(familyDir)
	if err != nil {
		t.Fatalf("read the prose-mechanics fixture family: %v", err)
	}
	seen := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		seen++
		t.Run(entry.Name(), func(t *testing.T) {
			fixture := filepath.Join(familyDir, entry.Name())
			want := strings.TrimSpace(readIfExists(filepath.Join(fixture, "EXPECT")))
			if want == "" {
				t.Fatalf("fixture %s has no EXPECT diagnostic", entry.Name())
			}
			root := t.TempDir()
			if err := canary.MaterializeMutationFixture(h.KitRoot, fixture, root); err != nil {
				t.Fatalf("materialize %s: %v", entry.Name(), err)
			}
			if diags := owner.run(root, h.KitRoot, registry.Dev); !containsDiagnostic(diags, want) {
				t.Fatalf("the mutated fixture did not raise %q:\n%s", want, strings.Join(diags, "\n"))
			}
			if err := canary.RestoreMutationFixture(h.KitRoot, fixture, root); err != nil {
				t.Fatalf("restore %s: %v", entry.Name(), err)
			}
			if diags := owner.run(root, h.KitRoot, registry.Dev); containsDiagnostic(diags, want) {
				t.Fatalf("the restored fixture still raises %q:\n%s", want, strings.Join(diags, "\n"))
			}
		})
	}
	if seen == 0 {
		t.Fatal("the prose-mechanics fixture family is empty")
	}
}
