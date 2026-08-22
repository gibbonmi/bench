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
// name. It is an independent literal on purpose: a set derived from the file it grades
// would accept every row somebody adds. A migration ticket removes rows and stays inside
// the set; a build that parks a hard file behind a new row reds here.
var approvedProseExclusionRows = map[string]bool{
	"tests/canary/":    true,
	"docs/audits/":     true,
	"CHANGELOG.md":     true,
	"capture/IDEAS.md": true,

	".bench/BENCH.md":           true,
	".bench/BENCH-reference.md": true,

	".agents/skills/bench-craft-spec/references/ste-prose.md":           true,
	".agents/skills/bench-craft-skills/":                                true,
	".agents/skills/bench-craft-comments/":                              true,
	".agents/skills/bench-craft-spec/SKILL.md":                          true,
	".agents/skills/bench-craft-spec/references/bootstrap-authority.md": true,
	".agents/skills/bench-craft-seams/":                                 true,
	".agents/skills/bench-craft-line/":                                  true,
	".agents/skills/bench-craft-domain/":                                true,
	".agents/skills/bench-craft-gate/":                                  true,
	".agents/skills/bench-craft-tickets/":                               true,
	".agents/skills/bench-craft-tdd/":                                   true,
	".agents/skills/bench-craft-delegate/":                              true,
	".agents/skills/bench-craft-grill/":                                 true,
	".agents/skills/bench-craft-adr/":                                   true,
	".agents/skills/bench-craft-review/":                                true,
	".agents/skills/bench-craft-synthesis/":                             true,
	".agents/skills/bench-craft-cli/":                                   true,
	".agents/skills/bench-craft-design-system/":                         true,
	".agents/skills/prototype/":                                         true,

	".agents/commands/bench.md":                       true,
	".agents/commands/bench-write-spec.md":            true,
	".agents/commands/bench-implement-spec.md":        true,
	".agents/commands/bench-review-implementation.md": true,
	".agents/commands/bench-final-check.md":           true,
	".agents/commands/bench-debug.md":                 true,
	".agents/commands/bench-what-next.md":             true,
	".agents/commands/bench-drain.md":                 true,
	".agents/commands/bench-setup-repo.md":            true,
	".agents/commands/bench-shape-idea.md":            true,
	".agents/commands/bench-deepen.md":                true,
	".agents/commands/bench-assess.md":                true,
	".agents/commands/bench-update-kit.md":            true,

	".agents/skills/bench/":                       true,
	".agents/skills/bench-assess/":                true,
	".agents/skills/bench-debug/":                 true,
	".agents/skills/bench-deepen/":                true,
	".agents/skills/bench-drain/":                 true,
	".agents/skills/bench-final-check/":           true,
	".agents/skills/bench-implement-spec/":        true,
	".agents/skills/bench-review-implementation/": true,
	".agents/skills/bench-setup-repo/":            true,
	".agents/skills/bench-shape-idea/":            true,
	".agents/skills/bench-update-kit/":            true,
	".agents/skills/bench-what-next/":             true,
	".agents/skills/bench-write-spec/":            true,

	"AGENTS.md":                         true,
	"CLAUDE.md":                         true,
	"CONTEXT.md":                        true,
	".claude/":                          true,
	"DATA_HANDLING.md":                  true,
	"SECURITY.md":                       true,
	"projects/gl-axi.md":                true,
	"docs/greenfield-build-sequence.md": true,
	"docs/release-runbook.md":           true,
	"docs/reporesident-distillation.md": true,
	"README.md":                         true,
	"projects/benchkit.md":              true,
	"ASSESSMENT.md":                     true,
	"skills-assessment.md":              true,
	"docs/adr/":                         true,
	"ROADMAP.md":                        true,
	"roadmap/":                          true,

	"decisions/gate-budget.md":                               true,
	"decisions/assets/gate-budget-cpu-wall-census.md":        true,
	"decisions/assets/gate-budget-memory-profile.md":         true,
	"decisions/gate-pipeline.md":                             true,
	"decisions/gate-critical-path.md":                        true,
	"decisions/gate-concurrency.md":                          true,
	"decisions/assets/gate-pipeline-fixture-inventory.md":    true,
	"decisions/ft183-gate-scoping-residuals.md":              true,
	"decisions/assets/ft183-derivation-binding.md":           true,
	"decisions/byte-preserving-axi-foundation/":              true,
	"decisions/bounded-network-resource-cli.md":              true,
	"decisions/craft-research.md":                            true,
	"decisions/assets/craft-research-research.md":            true,
	"decisions/cost-follows-project-size.md":                 true,
	"decisions/diff-visual.md":                               true,
	"decisions/spec-build-review-gate-cadence.md":            true,
	"decisions/worktree-orphan-retirement.md":                true,
	"decisions/ft144-post-approval-edits.md":                 true,
	"decisions/assets/ft171-shared-fixture-staged-binary.md": true,

	"tickets/":                             true,
	"capture/agent-performance/":           true,
	"capture/audits/":                      true,
	"capture/FIXES.md":                     true,
	"capture/parallel-session-friction.md": true,
	"capture/learnings.md":                 true,
	"capture/session-handoff.md":           true,
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
	rows := readProseExclusionRows(t, filepath.Join(h.KitRoot, filepath.FromSlash(prose.ExclusionFile)))
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

// readProseExclusionRows returns the subject of every row, under the same grammar the
// engine parses: a `#` starts a comment, a blank line is skipped, and the subject ends
// at the first space.
func readProseExclusionRows(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", prose.ExclusionFile, err)
	}
	var rows []string
	for _, line := range strings.Split(string(data), "\n") {
		if index := strings.IndexByte(line, '#'); index >= 0 {
			line = line[:index]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		subject := trimmed
		if index := strings.IndexAny(trimmed, " \t"); index >= 0 {
			subject = trimmed[:index]
		}
		rows = append(rows, subject)
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
