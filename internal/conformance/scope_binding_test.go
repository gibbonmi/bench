package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// scopeBindingRows is the fixed row set of the profile's reduced-scope table, in the
// order diagnostics report them. Each row names one accessor of the declaration, so
// the table renders the whole of gate.ReducedScope() and nothing else.
var scopeBindingRows = []string{"directories", "files", "excludable phases", "included phases"}

// declaredScopeSets derives the machine-readable side of the binding from the gate's
// declaration — never from a restated literal, so this check can only ever disagree
// with the prose, not with the oracle. The included row comes from the phase table the
// kit's reduced run actually executes, so a non-excludable phase added to that table
// reds this binding until the profile advertises it.
func declaredScopeSets() map[string][]string {
	scope := gate.ReducedScope()
	return map[string][]string{
		"directories":       scope.Directories(),
		"files":             scope.Files(),
		"excludable phases": scope.ExcludablePhases(),
		"included phases":   includedPhaseRow(),
	}
}

// includedPhaseRow resolves the kit's own phase table and derives the included set from
// it. A resolution failure is rendered as its own token rather than an empty row, so the
// binding fails closed with a diagnostic naming the cause instead of comparing the
// profile against silence.
func includedPhaseRow() []string {
	kit, err := findKitRoot()
	if err != nil {
		return []string{"(kit root unresolvable: " + err.Error() + ")"}
	}
	names, err := gate.IncludedPhaseNames(kit, kit)
	if err != nil {
		return []string{"(kit phase table unresolvable: " + err.Error() + ")"}
	}
	return names
}

// checkScopeBinding grades the profile's rendering of the reduced-scope declaration
// against gate.ReducedScope(), following checkLineBinding's shape: a prose table
// cross-checked against its machine-readable source, red on divergence. Without it a
// future edit widens the fast path in one place only, and the document and the oracle
// disagree with nothing to notice.
func checkScopeBinding(root string) []string {
	return scopeBindingDiags(declaredScopeSets(), readIfExists(filepath.Join(root, "projects", "benchkit.md")))
}

// scopeBindingDiags compares each table row's token set against the declared set by
// exact equality, never subset or substring: a subset comparison survives a
// prose-only addition, so equality is what makes a mutation on either side alone
// produce a diagnostic. The declaration arrives as a parameter so the bite test can
// drive the declaration-side mutation without editing the real one.
func scopeBindingDiags(declared map[string][]string, profile string) []string {
	rendered, tabled := profileScopeCells(profile)
	if !tabled {
		return []string{"profile reduced-scope table missing: projects/benchkit.md renders no '| reduced scope | ... |' table, so the declaration has no prose cross-check"}
	}
	var diags []string
	for _, row := range scopeBindingRows {
		prose, present := rendered[row]
		if !present {
			diags = append(diags, fmt.Sprintf("profile reduced-scope row missing: projects/benchkit.md's reduced-scope table has no '%s' row rendering gate.ReducedScope()'s [%s]", row, strings.Join(declared[row], ", ")))
			continue
		}
		if !equalTokenSets(prose, declared[row]) {
			diags = append(diags, fmt.Sprintf("profile reduced-scope row stale: projects/benchkit.md renders %s as [%s], but gate.ReducedScope() declares [%s]", row, strings.Join(prose, ", "), strings.Join(declared[row], ", ")))
		}
	}
	var unknown []string
	for row := range rendered {
		if _, known := declared[row]; !known {
			unknown = append(unknown, row)
		}
	}
	sort.Strings(unknown)
	for _, row := range unknown {
		diags = append(diags, fmt.Sprintf("profile reduced-scope row unknown: '%s' renders nothing gate.ReducedScope() declares (rows: %s)", row, strings.Join(scopeBindingRows, ", ")))
	}
	return diags
}

// profileScopeCells parses the profile's reduced-scope table into cells[row] = tokens,
// reporting whether such a table was found. The table is found by its `reduced scope`
// corner cell and read row by row the way the human it exists for reads it; an
// anchored table rather than a whole-file token search is what keeps a token that
// survives only in an unrelated paragraph from satisfying the cross-check.
func profileScopeCells(profile string) (map[string][]string, bool) {
	cells := map[string][]string{}
	found := false
	for _, line := range strings.Split(profile, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") {
			if found {
				break
			}
			continue
		}
		body := strings.TrimSuffix(strings.TrimPrefix(trimmed, "|"), "|")
		name, rest, _ := strings.Cut(body, "|")
		row := strings.ToLower(strings.TrimSpace(name))
		if !found {
			found = row == "reduced scope"
			continue
		}
		if row == "" || strings.Trim(row, "-:") == "" {
			continue
		}
		cells[row] = scopeTokens(rest)
	}
	return cells, found
}

// scopeTokens splits a table cell of comma-separated code-span tokens into bare
// strings, stripping the backticks the profile renders each entry inside.
func scopeTokens(cell string) []string {
	var tokens []string
	for _, raw := range strings.Split(cell, ",") {
		if token := strings.Trim(strings.TrimSpace(raw), "`"); token != "" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

// equalTokenSets reports order-insensitive equality with duplicates significant: the
// table may order a row for the reader, but every declared entry must appear exactly
// once, so a dropped entry and a phantom entry both break it.
func equalTokenSets(a, b []string) bool {
	a, b = slices.Clone(a), slices.Clone(b)
	sort.Strings(a)
	sort.Strings(b)
	return slices.Equal(a, b)
}

// TestRootConformanceScopeBinding grades the real kit checkout: the profile's
// reduced-scope table must match gate.ReducedScope(). It carries the entry point's
// name as a prefix so a `-run TestRootConformance` invocation reaches it alongside
// the registered suite, while the inner skip pattern — anchored to exactly
// TestRootConformance — still runs it in the conformance-suite phase.
func TestRootConformanceScopeBinding(t *testing.T) {
	kitRoot, err := findKitRoot()
	if err != nil {
		t.Fatalf("resolve kit root: %v", err)
	}
	for _, diag := range checkScopeBinding(kitRoot) {
		t.Errorf("gate: %s", diag)
	}
}

// renderScopeProfile builds a profile whose reduced-scope table renders sets, taking
// its cells from the same map shape the check compares against so the fixture and the
// expectation have one author here rather than a second hand-written list.
func renderScopeProfile(sets map[string][]string) string {
	var b strings.Builder
	b.WriteString("# benchkit\n\n## Gate\n\nThe reduced run's declaration, rendered:\n\n| reduced scope | declared |\n|---|---|\n")
	for _, row := range scopeBindingRows {
		entries := make([]string, 0, len(sets[row]))
		for _, entry := range sets[row] {
			entries = append(entries, "`"+entry+"`")
		}
		b.WriteString("| " + row + " | " + strings.Join(entries, ", ") + " |\n")
	}
	b.WriteString("\n## Notes for cold sessions\n\nThe declaration lives in the gate package.\n")
	return b.String()
}

// cloneScopeSets deep-copies a declaration map so a case can mutate its own copy.
func cloneScopeSets(sets map[string][]string) map[string][]string {
	out := make(map[string][]string, len(sets))
	for row, entries := range sets {
		out[row] = slices.Clone(entries)
	}
	return out
}

// TestDeclaredAllowlistBindingBites is the recorded bite proof for checkScopeBinding
// (per craft-gate), driving each direction of the binding separately: the prose
// mutated alone must produce a diagnostic, and the declaration mutated alone must
// too, because a subset or substring comparison survives one of the two. It runs
// against a synthetic root, with the real declaration as the fixed side.
func TestDeclaredAllowlistBindingBites(t *testing.T) {
	declared := declaredScopeSets()
	write := func(t *testing.T, profile string) string {
		t.Helper()
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "projects"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "projects", "benchkit.md"), []byte(profile), 0o644); err != nil {
			t.Fatal(err)
		}
		return root
	}

	if diags := checkScopeBinding(write(t, renderScopeProfile(declared))); len(diags) != 0 {
		t.Fatalf("a profile rendering the declaration exactly got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	for _, row := range scopeBindingRows {
		t.Run(row+" dropped from the prose alone", func(t *testing.T) {
			mutated := cloneScopeSets(declared)
			mutated[row] = mutated[row][1:]
			diags := checkScopeBinding(write(t, renderScopeProfile(mutated)))
			if len(diags) != 1 || !strings.Contains(diags[0], "profile reduced-scope row stale") || !strings.Contains(diags[0], row) {
				t.Fatalf("dropping a %s entry from the prose alone: want one stale diagnostic naming the row, got:\n%s", row, strings.Join(diags, "\n"))
			}
		})
		t.Run(row+" widened in the declaration alone", func(t *testing.T) {
			mutated := cloneScopeSets(declared)
			mutated[row] = append(mutated[row], "phantom-entry")
			diags := scopeBindingDiags(mutated, renderScopeProfile(declared))
			if len(diags) != 1 || !strings.Contains(diags[0], "profile reduced-scope row stale") || !strings.Contains(diags[0], row) {
				t.Fatalf("widening the %s declaration alone: want one stale diagnostic naming the row, got:\n%s", row, strings.Join(diags, "\n"))
			}
		})
	}

	t.Run("an entry rendered under the wrong row fires both rows", func(t *testing.T) {
		mutated := cloneScopeSets(declared)
		moved := mutated["files"][0]
		mutated["files"] = mutated["files"][1:]
		mutated["directories"] = append(mutated["directories"], moved)
		diags := checkScopeBinding(write(t, renderScopeProfile(mutated)))
		for _, row := range []string{"directories", "files"} {
			if !containsDiagnostic(diags, "renders "+row) {
				t.Fatalf("moving %q between rows did not fire the %s row:\n%s", moved, row, strings.Join(diags, "\n"))
			}
		}
	})

	t.Run("a missing row is named, not silently equal to empty", func(t *testing.T) {
		profile := strings.Replace(renderScopeProfile(declared), "| included phases |", "| retired phases |", 1)
		diags := checkScopeBinding(write(t, profile))
		if !containsDiagnostic(diags, "profile reduced-scope row missing") || !containsDiagnostic(diags, "'included phases'") {
			t.Fatalf("a renamed row was accepted:\n%s", strings.Join(diags, "\n"))
		}
		if !containsDiagnostic(diags, "profile reduced-scope row unknown") || !containsDiagnostic(diags, "'retired phases'") {
			t.Fatalf("an unknown row was accepted:\n%s", strings.Join(diags, "\n"))
		}
	})

	t.Run("a deleted table is named rather than scattering the tokens", func(t *testing.T) {
		var b strings.Builder
		b.WriteString("# benchkit\n\n## Gate\n\n")
		for _, row := range scopeBindingRows {
			b.WriteString("The " + row + " are " + strings.Join(declared[row], ", ") + ".\n")
		}
		diags := checkScopeBinding(write(t, b.String()))
		if len(diags) != 1 || !strings.Contains(diags[0], "profile reduced-scope table missing") {
			t.Fatalf("a profile scattering every token through prose was accepted:\n%s", strings.Join(diags, "\n"))
		}
	})
}
