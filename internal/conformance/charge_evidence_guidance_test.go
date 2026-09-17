package conformance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

const boundedBuildActionDiagnosticPrefix = "bounded build action: "

// buildActionSentenceLead opens every guidance sentence that states one bounded build
// action prerequisite.
const buildActionSentenceLead = "Build action requires"

// buildPhaseGuidance is the canonical build reader the bounded build action rows pin.
const buildPhaseGuidance = ".agents/commands/bench-implement-spec.md"

func boundedBuildActionAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(boundedBuildActionDiagnosticPrefix)
}

// TestEvidenceBuildGuidance grades the bounded build action rows for CE101, CE102, CE103,
// CE140, and CE141. It checks the row inventory, each row's kind, and that every row bites
// alone.
func TestEvidenceBuildGuidance(t *testing.T) {
	actionAnchors := boundedBuildActionAnchors()
	if got, want := len(actionAnchors), 6; got != want {
		t.Fatalf("bounded build action anchor count = %d, want %d", got, want)
	}
	for _, want := range []struct {
		kind       anchors.Kind
		diagnostic string
	}{
		{anchors.Require, "permits build action without verified delivery"},
		{anchors.Require, "permits build action without available required context"},
		{anchors.Require, "permits build action without a current binding"},
		{anchors.Require, "permits build action without reviewer approval"},
		{anchors.Require, "permits build action without the complete task supplement"},
		{anchors.Forbid, "re-advertises the retired build full charge form"},
	} {
		found := false
		for _, anchor := range actionAnchors {
			if anchor.Kind == want.kind && containsDiagnostic([]string{anchor.Diagnostic}, want.diagnostic) {
				found = true
			}
		}
		if !found {
			t.Errorf("no row of the expected kind states %q", want.diagnostic)
		}
	}
	runAnchorBites(t, actionAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

// TestEvidenceBuildActionSentencesArePinned binds the guidance to the registry. The
// guidance is the one source of each prerequisite sentence, and one Require row pins each
// one. A needle that keeps only its lead clause, or a needle that repeats another row's
// text, leaves a guidance sentence with no row and bites here.
func TestEvidenceBuildActionSentencesArePinned(t *testing.T) {
	needles := map[string]int{}
	for _, anchor := range boundedBuildActionAnchors() {
		if anchor.Kind == anchors.Require {
			needles[anchor.Needle]++
		}
	}
	for needle, rows := range needles {
		if rows != 1 {
			t.Errorf("%d Require rows share the needle %q, want one row each", rows, needle)
		}
	}
	sentences := buildActionSentences(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(buildPhaseGuidance)))
	if len(sentences) == 0 {
		t.Fatalf("%s states no bounded build action prerequisite", buildPhaseGuidance)
	}
	for _, sentence := range sentences {
		pinned := 0
		for needle := range needles {
			if needle == sentence || strings.HasPrefix(needle, sentence+" ") {
				pinned++
			}
		}
		if pinned != 1 {
			t.Errorf("guidance sentence %q is pinned by %d Require rows, want one", sentence, pinned)
		}
	}
}

// buildActionSentences returns every guidance sentence that opens with the bounded build
// action lead. A sentence runs to the first period a space follows, or to the line end.
func buildActionSentences(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var found []string
	for _, line := range strings.Split(string(data), "\n") {
		for rest := line; ; {
			lead := strings.Index(rest, buildActionSentenceLead)
			if lead < 0 {
				break
			}
			sentence := rest[lead:]
			if end := strings.Index(sentence, ". "); end >= 0 {
				sentence = sentence[:end+1]
			}
			found = append(found, sentence)
			rest = rest[lead+len(buildActionSentenceLead):]
		}
	}
	return found
}
