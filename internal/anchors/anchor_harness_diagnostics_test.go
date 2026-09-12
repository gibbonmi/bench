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
