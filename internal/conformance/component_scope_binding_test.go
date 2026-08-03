package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/gate"
)

// declaredComponentRows derives the machine-readable side of the per-component binding
// from the gate's own registry — never from a restated list, so a component added to
// gate.componentInputDeclarations reds this binding until the profile advertises it. The
// row order is the registry's own, so a diagnostic set is stable across runs without a
// second ordering kept here.
func declaredComponentRows() []gate.ComponentInputSource {
	return gate.ComponentInputSources()
}

// componentRowTokens is one component's expected row, both cells combined into a single
// token set: the declared source's parts (a multi-part source like
// "module-test-closure+manifest" renders as one token per part) plus the provenance word
// derived from the source itself, never restated as a second independent fact.
func componentRowTokens(source gate.Source) []string {
	tokens := strings.Split(string(source), "+")
	tokens = append(tokens, provenanceWord(source))
	return tokens
}

// provenanceWord names a source as derived or hand-written. gate.SourceHandDeclared is
// the registry's one hand-declared entry; every other source is computed from a listing.
func provenanceWord(source gate.Source) string {
	if source == gate.SourceHandDeclared {
		return "hand-written"
	}
	return "derived"
}

// checkComponentScopeBinding grades the profile's per-component table against
// gate.ComponentInputSources(), following checkScopeBinding's shape: a prose table
// cross-checked against its machine-readable source, red on divergence in either
// direction. Without it a component's declared source can drift from what the profile
// advertises with nothing to notice.
func checkComponentScopeBinding(root string) []string {
	return componentScopeBindingDiags(declaredComponentRows(), readIfExists(filepath.Join(root, "projects", "benchkit.md")))
}

// componentScopeBindingDiags compares each table row's combined token set against the
// declared component's expected tokens by exact equality, never subset: a subset
// comparison survives a prose-only addition, so equality is what makes a mutation on
// either side alone produce a diagnostic. declared arrives as a parameter so the bite
// test can drive the declaration-side mutation without a real registry entry changing.
func componentScopeBindingDiags(declared []gate.ComponentInputSource, profile string) []string {
	rendered, tabled := profileComponentCells(profile)
	if !tabled {
		return []string{"profile per-component table missing: projects/benchkit.md renders no '| component | ... |' table, so the per-component declarations have no prose cross-check"}
	}
	known := make(map[string]struct{}, len(declared))
	var diags []string
	for _, d := range declared {
		known[d.Component] = struct{}{}
		expected := componentRowTokens(d.Source)
		prose, present := rendered[d.Component]
		if !present {
			diags = append(diags, fmt.Sprintf("profile per-component row missing: projects/benchkit.md's per-component table has no %q row rendering gate.ComponentInputSources()'s [%s]", d.Component, strings.Join(expected, ", ")))
			continue
		}
		if !equalTokenSets(prose, expected) {
			diags = append(diags, fmt.Sprintf("profile per-component row stale: projects/benchkit.md renders %q as [%s], but gate.ComponentInputSources() declares [%s]", d.Component, strings.Join(prose, ", "), strings.Join(expected, ", ")))
		}
	}
	var unknown []string
	for component := range rendered {
		if _, ok := known[component]; !ok {
			unknown = append(unknown, component)
		}
	}
	sort.Strings(unknown)
	for _, component := range unknown {
		diags = append(diags, fmt.Sprintf("profile per-component row unknown: %q renders nothing gate.ComponentInputSources() declares", component))
	}
	return diags
}

// profileComponentCells parses the profile's per-component table into cells[component] =
// combined tokens from the "declares" and "provenance" columns, reporting whether such a
// table was found. The table is found by its `component` corner cell and read row by row
// the way the human it exists for reads it; an anchored table rather than a whole-file
// token search is what keeps a token that survives only in an unrelated paragraph from
// satisfying the cross-check.
func profileComponentCells(profile string) (map[string][]string, bool) {
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
		component := strings.TrimSpace(name)
		if !found {
			found = strings.ToLower(component) == "component"
			continue
		}
		if component == "" || strings.Trim(component, "-:") == "" {
			continue
		}
		declares, provenance, _ := strings.Cut(rest, "|")
		cells[component] = append(scopeTokens(declares), scopeTokens(provenance)...)
	}
	return cells, found
}

// renderComponentProfile builds a profile whose per-component table renders declared,
// taking its rows from the same slice shape the check compares against so the fixture
// and the expectation have one author here rather than a second hand-written list.
func renderComponentProfile(declared []gate.ComponentInputSource) string {
	var b strings.Builder
	b.WriteString("# benchkit\n\n## Gate\n\nThe per-component declarations, rendered:\n\n| component | declares | provenance |\n|---|---|---|\n")
	for _, d := range declared {
		sourceTokens := strings.Split(string(d.Source), "+")
		entries := make([]string, 0, len(sourceTokens))
		for _, entry := range sourceTokens {
			entries = append(entries, "`"+entry+"`")
		}
		b.WriteString("| " + d.Component + " | " + strings.Join(entries, ", ") + " | `" + provenanceWord(d.Source) + "` |\n")
	}
	b.WriteString("\n## Notes for cold sessions\n\nThe declarations live in the gate package.\n")
	return b.String()
}

// cloneComponentRows deep-copies the declared slice so a case can mutate its own copy
// without the real registry's values moving under it.
func cloneComponentRows(declared []gate.ComponentInputSource) []gate.ComponentInputSource {
	out := make([]gate.ComponentInputSource, len(declared))
	copy(out, declared)
	return out
}

// TestComponentScopeBindingBites is the recorded bite proof for checkComponentScopeBinding
// (per craft-gate), driving each direction of the binding separately: the prose mutated
// alone must produce a diagnostic, and the declaration mutated alone must too, because a
// subset or substring comparison survives one of the two. It runs against a synthetic
// root, with the real registry as the fixed side.
func TestComponentScopeBindingBites(t *testing.T) {
	declared := declaredComponentRows()
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

	if diags := checkComponentScopeBinding(write(t, renderComponentProfile(declared))); len(diags) != 0 {
		t.Fatalf("a profile rendering the declaration exactly got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	for i, d := range declared {
		i, d := i, d
		t.Run(d.Component+" mutated in the prose alone", func(t *testing.T) {
			// Flip the provenance word on this component's row alone: it is the one
			// cell every row carries regardless of how many source parts it has.
			lines := strings.Split(renderComponentProfile(declared), "\n")
			for li, line := range lines {
				if !strings.HasPrefix(line, "| "+d.Component+" | ") {
					continue
				}
				switch {
				case strings.Contains(line, "`derived`"):
					lines[li] = strings.Replace(line, "`derived`", "`hand-written`", 1)
				case strings.Contains(line, "`hand-written`"):
					lines[li] = strings.Replace(line, "`hand-written`", "`derived`", 1)
				}
			}
			diags := checkComponentScopeBinding(write(t, strings.Join(lines, "\n")))
			if !containsDiagnostic(diags, "profile per-component row stale") || !containsDiagnostic(diags, d.Component) {
				t.Fatalf("mutating %s's prose alone: want a stale diagnostic naming the component, got:\n%s", d.Component, strings.Join(diags, "\n"))
			}
		})
		t.Run(d.Component+" widened in the declaration alone", func(t *testing.T) {
			mutated := cloneComponentRows(declared)
			mutated[i] = gate.ComponentInputSource{Component: d.Component, Source: gate.Source(string(d.Source) + "+phantom-part")}
			diags := componentScopeBindingDiags(mutated, renderComponentProfile(declared))
			if !containsDiagnostic(diags, "profile per-component row stale") || !containsDiagnostic(diags, d.Component) {
				t.Fatalf("widening %s's declaration alone: want a stale diagnostic naming the component, got:\n%s", d.Component, strings.Join(diags, "\n"))
			}
		})
	}

	t.Run("a component missing from the table and an unknown row are each named", func(t *testing.T) {
		removed := declared[0].Component
		kept := append(cloneComponentRows(declared)[1:], gate.ComponentInputSource{Component: "unknown-component", Source: gate.SourceHandDeclared})
		profile := renderComponentProfile(kept)
		diags := checkComponentScopeBinding(write(t, profile))
		if !containsDiagnostic(diags, "profile per-component row missing") || !containsDiagnostic(diags, removed) {
			t.Fatalf("removing %s's row was accepted:\n%s", removed, strings.Join(diags, "\n"))
		}
		if !containsDiagnostic(diags, "profile per-component row unknown") || !containsDiagnostic(diags, "unknown-component") {
			t.Fatalf("an unknown row was accepted:\n%s", strings.Join(diags, "\n"))
		}
	})

	t.Run("a deleted table is named rather than scattering the tokens", func(t *testing.T) {
		var b strings.Builder
		b.WriteString("# benchkit\n\n## Gate\n\n")
		for _, d := range declared {
			b.WriteString("The " + d.Component + " component declares " + strings.Join(componentRowTokens(d.Source), ", ") + ".\n")
		}
		diags := checkComponentScopeBinding(write(t, b.String()))
		if len(diags) != 1 || !strings.Contains(diags[0], "profile per-component table missing") {
			t.Fatalf("a profile scattering every token through prose was accepted:\n%s", strings.Join(diags, "\n"))
		}
	})
}

// honestyResidualPhrases are the sentences the profile must carry, verbatim enough to
// survive an editorial rewrite but specific enough that dropping the residual paragraph
// (or the canary-narrowing sentence inside it) cannot pass silently. Each phrase is its
// own diagnostic so a partial deletion still names exactly what went missing.
var honestyResidualPhrases = []string{
	"capture-surface blindness only",
	"reads an undeclared non-capture path skips wrongly",
	"2026-08-01 narrowing",
	"never run against their combined tree",
}

// checkComponentHonestyProse grades the profile for the declaration-honesty residual and
// the canary narrowing (ticket: "bind the per-component table to the declaration"), each
// phrase its own diagnostic when absent so a partial deletion is still named. Matching
// runs against whitespace-collapsed text so a phrase that Markdown line-wraps mid-sentence
// still matches — the wrapping is a rendering accident, not a content difference.
func checkComponentHonestyProse(profile string) []string {
	collapsed := strings.Join(strings.Fields(profile), " ")
	var diags []string
	for _, phrase := range honestyResidualPhrases {
		if !strings.Contains(collapsed, phrase) {
			diags = append(diags, fmt.Sprintf("profile honesty-residual prose missing: projects/benchkit.md does not carry %q", phrase))
		}
	}
	return diags
}

func checkComponentHonestyProfile(kitRoot string) []string {
	return checkComponentHonestyProse(readIfExists(filepath.Join(kitRoot, "projects", "benchkit.md")))
}

// TestProfileStatesTheHonestyResidual is the recorded bite proof (per craft-gate): deleting
// the residual paragraph from a profile that otherwise carries it reds the check.
func TestProfileStatesTheHonestyResidual(t *testing.T) {
	profile := "Declaration-honesty width: " + strings.Join(honestyResidualPhrases, ". ")
	if diags := checkComponentHonestyProse(profile); len(diags) != 0 {
		t.Fatalf("complete residual prose got diagnostics:\n%s", strings.Join(diags, "\n"))
	}

	diags := checkComponentHonestyProse("Declaration-honesty width")
	if len(diags) != len(honestyResidualPhrases) {
		t.Fatalf("deleting the residual paragraph: want a diagnostic per phrase, got:\n%s", strings.Join(diags, "\n"))
	}
}
