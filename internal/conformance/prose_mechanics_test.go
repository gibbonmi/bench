package conformance

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/canary"
	"github.com/gibbonmi/bench/internal/conformance/registry"
	"github.com/gibbonmi/bench/internal/prose"
)

// checkProseMechanics is the registered wrapper over internal/prose. The package owns
// the parser, the walk, the exclusion grammar, and the classification, so the check is
// the binding and nothing else.
func checkProseMechanics(root string) []string {
	return prose.Grade(root)
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
