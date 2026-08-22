package handoff

import (
	"testing"

	"github.com/gibbonmi/bench/internal/prose"
)

// TestTemplatesPassProse renders every Markdown template this package writes into a
// tree and grades each render through internal/prose. This pins ASD-STE100 conformance
// at the source rather than at the live-tree sweep alone, so a regression here reds at
// the template that caused it.
func TestTemplatesPassProse(t *testing.T) {
	templates := map[string]string{
		"ShapeSection":     ShapeSection,
		"scaffoldGuidance": scaffoldGuidance,
	}
	for name, doc := range templates {
		if findings := prose.Findings(doc); len(findings) != 0 {
			t.Errorf("%s: prose.Findings = %v, want none", name, findings)
		}
	}
}

// TestTemplatesPassProseCatchesALongSentence plants an over-bound sentence in a template
// and asserts prose.Findings reds it (PD42): the check is not a pass-through.
func TestTemplatesPassProseCatchesALongSentence(t *testing.T) {
	planted := scaffoldGuidance + "\n" +
		"This single planted sentence carries far more than the twenty five word bound the rule file states because it keeps adding clause after clause after clause without ever stopping to close.\n"
	if findings := prose.Findings(planted); len(findings) == 0 {
		t.Fatal("prose.Findings on a planted long sentence = no findings, want at least one")
	}
}
