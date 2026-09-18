package conformance

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/prose"
)

const boundedBuildActionDiagnosticPrefix = "bounded build action: "

// buildActionSentenceLead opens every guidance sentence that states one bounded build
// action prerequisite. A paragraph that carries the lead is pinned prose.
const buildActionSentenceLead = "Build action requires"

// The retired build full charge form is a flag pair, not one spelling. The guidance may
// name the preflight build command, and it may name the run control, but one sentence may
// not carry both.
const (
	preflightBuildCommand = "bench preflight build"
	retiredBuildFullFlag  = "--full"
)

func boundedBuildActionAnchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(boundedBuildActionDiagnosticPrefix)
}

// boundedBuildActionGuidance returns the one canonical reader the family pins. The rows'
// File fields are that path's single source, so a row that names a second file fails here.
func boundedBuildActionGuidance(t *testing.T, family []anchors.Anchor) string {
	t.Helper()
	files := map[string]bool{}
	for _, anchor := range family {
		files[anchor.File] = true
	}
	if len(files) != 1 {
		t.Fatalf("bounded build action rows name %d guidance files, want one", len(files))
	}
	for file := range files {
		return file
	}
	return ""
}

// TestEvidenceBuildGuidance grades the bounded build action rows for CE101, CE102, CE103,
// CE140, and CE141. The expectation table below is the independent half of this pair: it
// states each row's kind and diagnostic apart from the registry, so a reworded or deleted
// row bites. The count comes from that table's length, and the mismatch reports rather than
// stops, so a deleted row also reaches the table.
func TestEvidenceBuildGuidance(t *testing.T) {
	actionAnchors := boundedBuildActionAnchors()
	want := []struct {
		kind       anchors.Kind
		diagnostic string
	}{
		{anchors.Require, "permits build action without verified delivery"},
		{anchors.Require, "permits build action without available required context"},
		{anchors.Require, "permits build action without a current binding"},
		{anchors.Require, "permits build action without reviewer approval"},
		{anchors.Require, "permits build action without the complete task supplement"},
		{anchors.Forbid, "re-advertises the retired build full charge form"},
	}
	if got := len(actionAnchors); got != len(want) {
		t.Errorf("bounded build action anchor count = %d, want %d", got, len(want))
	}
	for _, expected := range want {
		found := false
		for _, anchor := range actionAnchors {
			if anchor.Kind == expected.kind && containsDiagnostic([]string{anchor.Diagnostic}, expected.diagnostic) {
				found = true
			}
		}
		if !found {
			t.Errorf("no row of the expected kind states %q", expected.diagnostic)
		}
	}
	runAnchorBites(t, actionAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
}

// TestEvidenceBuildActionSentencesArePinned binds the guidance to the registry. A paragraph
// that states one bounded build action prerequisite is pinned prose, so exactly one Require
// needle covers each of its sentences. A needle that keeps only its lead clause, a needle
// that repeats another row's text, and an unpinned sentence added to that paragraph all
// bite here. One parser reads the guidance and the needle, so a reflow that changes no word
// changes nothing here.
func TestEvidenceBuildActionSentencesArePinned(t *testing.T) {
	family := boundedBuildActionAnchors()
	needles := map[string]int{}
	for _, anchor := range family {
		if anchor.Kind == anchors.Require {
			needles[anchor.Needle]++
		}
	}
	for needle, rows := range needles {
		if rows != 1 {
			t.Errorf("%d Require rows share the needle %q, want one row each", rows, needle)
		}
	}
	guidance := boundedBuildActionGuidance(t, family)
	sentences := pinnedParagraphSentences(guidanceText(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(guidance))))
	if len(sentences) == 0 {
		t.Fatalf("%s states no bounded build action prerequisite", guidance)
	}
	for _, sentence := range sentences {
		pinned := 0
		for needle := range needles {
			if strings.Contains(needleText(needle), sentence) {
				pinned++
			}
		}
		if pinned != 1 {
			t.Errorf("guidance sentence %q is pinned by %d Require rows, want one", sentence, pinned)
		}
	}
}

// TestEvidenceBuildGuidanceRejectsThePreflightBuildFullPair grades the guidance half of
// CE173. The retired form is the pair, so one guidance sentence may not name the preflight
// build command and also carry the run-control flag. Any flag order and any wrap between the
// two halves re-advertises the retired form and bites here. The phase's own full-run control
// names no preflight command, so it stays legal.
func TestEvidenceBuildGuidanceRejectsThePreflightBuildFullPair(t *testing.T) {
	guidance := boundedBuildActionGuidance(t, boundedBuildActionAnchors())
	pairs := preflightBuildFullPairs(guidanceText(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(guidance))))
	for _, sentence := range pairs {
		t.Errorf("%s states %q, which names %q and carries %q, and that re-advertises the retired build full charge form",
			guidance, sentence, preflightBuildCommand, retiredBuildFullFlag)
	}
}

// TestEvidenceBuildGuidanceRulesBiteOnSyntheticText grades the two rules above apart from the
// shipped guidance. Each case feeds text through the rule's own predicate, so a rule that is
// weakened in place reds here while the shipped guidance still satisfies it. This guard is a
// case table rather than a canary fixture, because a canary mutates the guidance file and a
// weakened rule accepts every mutation of it. The name carries the family prefix, so the
// chunk's declared TestEvidence run reaches it.
func TestEvidenceBuildGuidanceRulesBiteOnSyntheticText(t *testing.T) {
	for _, tt := range []struct {
		name string
		rule func(string) []string
		text string
		want []string
	}{
		{
			name: "the lead opens the second sentence of a pinned paragraph",
			rule: pinnedParagraphSentences,
			text: "The ticket lands first. Build action requires reviewer approval.\n",
			want: []string{"The ticket lands first.", "Build action requires reviewer approval."},
		},
		{
			name: "a wrap inside a pinned sentence collapses",
			rule: pinnedParagraphSentences,
			text: "Build action requires verified delivery: act only after\nthis session reads it.\n",
			want: []string{"Build action requires verified delivery: act only after this session reads it."},
		},
		{
			name: "a neighbouring paragraph stays outside the pin",
			rule: pinnedParagraphSentences,
			text: "The reviewer may waive it.\n\nBuild action requires reviewer approval.\n",
			want: []string{"Build action requires reviewer approval."},
		},
		{
			name: "a paragraph with no lead states no prerequisite",
			rule: pinnedParagraphSentences,
			text: "The ticket lands first. It commits green.\n",
		},
		{
			name: "one sentence names the command and the flag",
			rule: preflightBuildFullPairs,
			text: "Run `bench preflight build <slug> --full` now.\n",
			want: []string{"Run `bench preflight build <slug> --full` now."},
		},
		{
			name: "a wrap splits the command from the flag",
			rule: preflightBuildFullPairs,
			text: "Run `bench preflight build <slug>` with the retired\n`--full` control now.\n",
			want: []string{"Run `bench preflight build <slug>` with the retired `--full` control now."},
		},
		{
			name: "the command alone advertises nothing",
			rule: preflightBuildFullPairs,
			text: "Run `bench preflight build <slug>` now.\n",
		},
		{
			name: "the run control alone advertises nothing",
			rule: preflightBuildFullPairs,
			text: "The `--full` run reaches the landing.\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule(tt.text); !slices.Equal(got, tt.want) {
				t.Errorf("rule = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestEvidenceNeedleTextReadsTheGuidanceProjection grades the needle half of the pin. A
// needle reaches the comparison through the same projection the guidance reaches it through,
// so a wrap inside a needle, a run of spaces, and a needle of two sentences all read as the
// guidance reads them. A needle that skips that projection keeps its own bytes and reds here.
func TestEvidenceNeedleTextReadsTheGuidanceProjection(t *testing.T) {
	for _, tt := range []struct {
		name   string
		needle string
		want   string
	}{
		{
			name:   "a wrap inside a needle collapses",
			needle: "Build action requires reviewer approval: act only after\nthe reviewer approves it.\n",
			want:   "Build action requires reviewer approval: act only after the reviewer approves it.",
		},
		{
			name:   "a run of spaces collapses",
			needle: "Build action  requires  reviewer approval.\n",
			want:   "Build action requires reviewer approval.",
		},
		{
			name:   "two sentences join as the paragraph joins them",
			needle: "Build action requires reviewer approval. The ticket lands first.\n",
			want:   "Build action requires reviewer approval. The ticket lands first.",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := needleText(tt.needle); got != tt.want {
				t.Errorf("needleText = %q, want %q", got, tt.want)
			}
		})
	}
}

// guidanceText reads one canonical reader with its line endings normalized.
func guidanceText(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n")
}

// pinnedParagraphSentences returns every sentence of every paragraph that carries the bounded
// build action lead. It reads prose.Paragraphs, which is the paragraph rule and the sentence
// rule the prose gate itself applies, so the pin and the gate cannot disagree about where a
// paragraph or a sentence ends. Each sentence arrives with its whitespace collapsed.
func pinnedParagraphSentences(text string) []string {
	var found []string
	for _, paragraph := range prose.Paragraphs(text) {
		if slices.ContainsFunc(paragraph, func(s string) bool { return strings.Contains(s, buildActionSentenceLead) }) {
			found = append(found, paragraph...)
		}
	}
	return found
}

// preflightBuildFullPairs returns every sentence that names the preflight build command and
// also carries the retired run control. The unit is one sentence of the same parser, so a wrap
// between the two halves hides nothing and a lawful run control in another sentence stays
// lawful. A heading and a frontmatter field hold no sentence; the literal Forbid row grades
// those, because an anchor reads the whole file.
func preflightBuildFullPairs(text string) []string {
	var found []string
	for _, paragraph := range prose.Paragraphs(text) {
		for _, sentence := range paragraph {
			if strings.Contains(sentence, preflightBuildCommand) && strings.Contains(sentence, retiredBuildFullFlag) {
				found = append(found, sentence)
			}
		}
	}
	return found
}

// needleText reads one needle through the projection that reads the guidance, so the pin
// compares two texts that one parser normalized. A needle states whole sentences, and they
// join here as the guidance paragraph joins them.
func needleText(needle string) string {
	var sentences []string
	for _, paragraph := range prose.Paragraphs(needle) {
		sentences = append(sentences, paragraph...)
	}
	return strings.Join(sentences, " ")
}
