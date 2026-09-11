package conformance

import (
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/lines"
)

func TestRetainedWorkflow(t *testing.T) {
	h := NewHarness(t)
	binding := lines.ParseBinding([]byte(h.ReadRootFile(".bench", "lines.env")))
	if got, want := binding.Cell("codex", "mid"), "gpt-5.6-sol"; got != want {
		t.Fatalf("Codex mid binding = %q, want %q", got, want)
	}
	if got, want := binding.Cell("codex", "top"), "gpt-6-astra"; got != want {
		t.Fatalf("Codex top binding = %q, want %q", got, want)
	}
	if got, want := binding.Cell("codex", "cheap"), "gpt-5.6-terra"; got != want {
		t.Fatalf("Codex cheap binding = %q, want %q", got, want)
	}

	wantDiagnostics := []string{
		"retained workflow: craft-spec dropped the one implementation-line template",
		"retained workflow: craft-spec dropped the implementation-line reason factors",
		"retained workflow: craft-spec dropped the harder-chunk annotation",
		"retained workflow: craft-line dropped the conditional review stage default",
		"retained workflow: craft-line dropped the Codex mid-to-top review branch",
		"retained workflow: craft-line dropped the default mid review branch",
		"retained workflow: review phase dropped the conditional review-line reference",
		"retained workflow: project profile dropped the conditional review-line reference",
		"retained workflow: craft-line dropped the user-directed model-switch boundary",
	}
	wantPredicates := map[string]struct {
		file    string
		section string
		needle  string
	}{
		"retained workflow: craft-spec dropped the implementation-line reason factors": {
			file:    ".agents/skills/bench-craft-spec/SKILL.md",
			section: "Template",
			needle:  "Implementation-line reason: <hardest material chunk, spec precision, seam uncertainty, and test strength>.",
		},
		"retained workflow: craft-line dropped the Codex mid-to-top review branch": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "A Codex mid implementation sends each axis to the Codex top binding.",
		},
		"retained workflow: craft-line dropped the default mid review branch": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "Every other implementation sends each axis to the invoking harness's mid binding.",
		},
		"retained workflow: craft-line dropped the user-directed model-switch boundary": {
			file:   ".agents/skills/bench-craft-line/SKILL.md",
			needle: "A different implementation model or session requires user direction. The author can adjust effort in the retained session and reports the change.",
		},
	}
	family := anchorsWithDiagnosticPrefix("retained workflow: ")
	for _, want := range wantDiagnostics {
		index := slices.IndexFunc(family, func(anchor anchors.Anchor) bool { return anchor.Diagnostic == want })
		if index < 0 {
			t.Errorf("retained-workflow anchor is absent: %s", want)
			continue
		}
		if predicate, ok := wantPredicates[want]; ok {
			anchor := family[index]
			if anchor.File != predicate.file || anchor.Section != predicate.section || anchor.Needle != predicate.needle {
				t.Errorf("retained-workflow predicate drifted for %q: %#v", want, anchor)
			}
		}
	}
	if t.Failed() {
		return
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}
