package anchors

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

type anchorErrors []string

func TestEvaluatePathReportsRegistryDiagnostics(t *testing.T) {
	for _, anchor := range Entries() {
		if anchor.File != "AGENTS.md" || anchor.Kind != Require {
			continue
		}
		h := anchorHarness{rules: []anchorRule{{file: anchor.File, needle: anchor.Needle}}}
		result := EvaluatePath(h.write(t, 0), anchor.File)
		if !slices.Contains(result.Diagnostics, anchor.Diagnostic) {
			t.Fatalf("path diagnostics = %v, want missing-anchor diagnostic %q", result.Diagnostics, anchor.Diagnostic)
		}
		return
	}
	t.Fatal("registry has no required AGENTS.md anchor")
}

func TestEvaluatePathAnchorKinds(t *testing.T) {
	for _, tc := range []struct {
		kind      Kind
		forbidden bool
	}{
		{Require, false}, {Forbid, true},
		{RequireInSection, false}, {ForbidInSection, true},
		{ForbidCaseFoldedEmphasis, true},
		{RequireInStep, false},
	} {
		t.Run(fmt.Sprint(tc.kind), func(t *testing.T) {
			anchor := pathTestAnchor(t, tc.kind)
			h := anchorHarness{rules: []anchorRule{{file: anchor.File, section: anchor.Section, step: anchor.Step, needle: anchor.Needle, forbidden: tc.forbidden}}}
			for _, broken := range []int{-1, 0} {
				result := EvaluatePath(h.write(t, broken), anchor.File)
				if got := slices.Contains(result.Diagnostics, anchor.Diagnostic); got != (broken == 0) {
					t.Fatalf("broken=%d diagnostics = %v, want violation=%t for %q", broken, result.Diagnostics, broken == 0, anchor.Diagnostic)
				}
				index := slices.IndexFunc(result.Locations, func(location Location) bool { return location.Anchor == anchor })
				if index < 0 || (result.Locations[index].Line > 0) != ((broken == 0) == tc.forbidden) {
					t.Fatalf("broken=%d locations = %v, want the registered anchor's presence", broken, result.Locations)
				}
			}
		})
	}
}

func TestEvaluatePathIgnoresCaseFoldedEmphasisInsideHTMLComment(t *testing.T) {
	anchor := pathTestAnchor(t, ForbidCaseFoldedEmphasis)
	root := t.TempDir()
	path := filepath.Join(root, filepath.FromSlash(anchor.File))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("<!-- An executable **red** is mandatory. -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result := EvaluatePath(root, anchor.File)
	if slices.Contains(result.Diagnostics, anchor.Diagnostic) {
		t.Fatalf("comment-only phrase raised %q", anchor.Diagnostic)
	}
	for _, location := range result.Locations {
		if location.Anchor == anchor && location.Line != 0 {
			t.Fatalf("comment-only phrase located at line %d, want 0", location.Line)
		}
	}
}

func TestEvaluatePathRefusesInvalidSubjects(t *testing.T) {
	anchor := pathTestAnchor(t, RequireInSection)
	for _, tc := range []struct {
		name string
		body string
		want string
	}{
		{"missing", "", "section-scoped anchor file missing: " + anchor.File},
		{"empty", "", fmt.Sprintf("%s is missing the %q section that owns a scoped anchor", anchor.File, anchor.Section)},
		{"duplicated", strings.Repeat("## "+anchor.Section+"\n\n"+anchor.Needle+"\n", 2), fmt.Sprintf("%s carries 2 %q sections; a scoped anchor needs exactly one owning section", anchor.File, anchor.Section)},
		{"directory", "", RefusalPrefix + anchor.File},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			path := filepath.Join(root, filepath.FromSlash(anchor.File))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if tc.name == "directory" {
				if err := os.Mkdir(path, 0o755); err != nil {
					t.Fatal(err)
				}
			} else if tc.name != "missing" {
				if err := os.WriteFile(path, []byte(tc.body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			result := EvaluatePath(root, anchor.File)
			if !strings.Contains(strings.Join(result.Diagnostics, "\n"), tc.want) {
				t.Fatalf("diagnostics = %v, want %q", result.Diagnostics, tc.want)
			}
		})
	}
}

// TestResolveStepRefusesAStepScopedAnchorWithNoStep pins the kind as the one owner of step
// narrowing. The evaluator narrows on the kind alone, so a step-scoped anchor that names no
// step must say so rather than read its whole section as if it were unscoped.
func TestResolveStepRefusesAStepScopedAnchorWithNoStep(t *testing.T) {
	const file, title = "guide.md", "Process"
	want := fmt.Sprintf("%s carries a step-scoped anchor with no step in the %q section; a step-scoped anchor names one step", file, title)
	if got := resolveStep(file, title, 0, "1. step one\n   needle\n"); got.diagnostic != want || got.body != "" {
		t.Fatalf("resolveStep with no step = %+v, want the unnamed-step diagnostic %q and no body", got, want)
	}
}

// TestRegistryBindsStepToItsKind refuses the authoring mistake the evaluator no longer reads
// around: a Step on a kind that narrows no step, or a step-scoped kind with no Step. The kind
// and the field must agree, because the kind alone decides the narrowing.
func TestRegistryBindsStepToItsKind(t *testing.T) {
	for _, anchor := range Entries() {
		if anchor.Kind.stepScoped() != (anchor.Step != 0) {
			t.Errorf("anchor %q on %s has kind %d with step %d; a step belongs to a step-scoped kind and to no other", anchor.Needle, anchor.File, anchor.Kind, anchor.Step)
		}
	}
}

func pathTestAnchor(t *testing.T, kind Kind) Anchor {
	t.Helper()
	entries := Entries()
	counts := map[string]int{}
	for _, anchor := range entries {
		counts[anchor.Diagnostic]++
	}
	for _, anchor := range entries {
		if anchor.Kind == kind && counts[anchor.Diagnostic] == 1 {
			return anchor
		}
	}
	t.Fatalf("registry has no anchor of kind %d", kind)
	return Anchor{}
}

func (*anchorErrors) Helper() {}
func (e *anchorErrors) Errorf(format string, args ...any) {
	*e = append(*e, fmt.Sprintf(format, args...))
}

func TestAnchorHarnessAttributesFailureDiagnostics(t *testing.T) {
	const subject = ".agents/skills/bench-craft-spec/SKILL.md"
	const other = ".bench/BENCH.md"
	var own, foreign Anchor
	for _, entry := range Entries() {
		if entry.Group != AfterImplementSpec || !strings.HasPrefix(entry.Diagnostic, "retained workflow:") {
			continue
		}
		if entry.File == subject && own.File == "" {
			own = entry
		}
		if entry.File == other && foreign.File == "" {
			foreign = entry
		}
	}
	if own.File == "" || foreign.File == "" {
		t.Fatal("custom-prefix fixture anchors are absent")
	}
	h := anchorHarness{group: AfterImplementSpec, rules: []anchorRule{
		{file: own.File, section: own.Section, needle: own.Needle, want: own.Diagnostic},
		{file: foreign.File, section: foreign.Section, needle: foreign.Needle, want: foreign.Diagnostic},
	}}
	root := h.write(t, 0)
	// Removing the other subject gives the full group a diagnostic from another file.
	// The failed assertion must still display only its own subject.
	if err := os.Remove(filepath.Join(root, filepath.FromSlash(other))); err != nil {
		t.Fatal(err)
	}
	diags := EvaluateGroup(root, h.group)
	if !slices.Contains(diags, own.Diagnostic) || len(diags) < 2 {
		t.Fatalf("fixture did not produce own and unrelated diagnostics: %v", diags)
	}
	h.rules[0].want = "deliberately mismatched expectation"
	var failures anchorErrors
	h.checkBroken(&failures, root, 0, diags)
	if len(failures) != 1 || !strings.Contains(failures[0], own.Diagnostic) || strings.Contains(failures[0], other) {
		t.Fatalf("failed assertion mixed diagnostic ownership: %v", failures)
	}
	// Cross-talk still reads the full slice even though failure display is scoped.
	h.rules[0].want = own.Diagnostic
	failures = nil
	h.checkBroken(&failures, root, 0, append(diags, foreign.Diagnostic))
	if len(failures) != 1 || !strings.Contains(failures[0], foreign.Diagnostic) {
		t.Fatalf("cross-talk escaped the assertion: %v", failures)
	}
}
