package conformance

import (
	"fmt"
	"slices"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

func TestImplementationContinuation(t *testing.T) {
	wantPredicates := map[string]struct {
		file    string
		section string
		needle  string
	}{
		"implementation continuation: craft-delegate restored its copy of the diagnostic helper boundary that the operating guide owns": {
			file: ".agents/skills/bench-craft-delegate/SKILL.md", section: "",
			needle: "A diagnostic helper can inspect evidence, but it receives no implementation or repair assignment.",
		},
		"implementation continuation: bounded repair restored the retired mode list": {
			file: ".agents/skills/bench-craft-line/references/bounded-repair-policy.md", section: "",
			needle: "This allowance applies to retained, full, delegated, unattended, and light-path implementation runs.",
		},
	}
	family := anchorsWithDiagnosticPrefix("implementation continuation: ")
	for diagnostic, predicate := range wantPredicates {
		index := slices.IndexFunc(family, func(anchor anchors.Anchor) bool { return anchor.Diagnostic == diagnostic })
		if index < 0 {
			t.Errorf("implementation-continuation anchor is absent: %s", diagnostic)
			continue
		}
		anchor := family[index]
		if anchor.File != predicate.file || anchor.Section != predicate.section || anchor.Needle != predicate.needle {
			t.Errorf("implementation-continuation predicate drifted for %q: %s", diagnostic, fmt.Sprintf("%#v", anchor))
		}
	}
	runAnchorBites(t, family, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}
