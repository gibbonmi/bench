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

// boundedActionFamily is one bounded action guidance family. Each phase states the same
// five prerequisites over its own canonical reader, so the rules below are written once and
// bound to a family rather than copied per phase.
type boundedActionFamily struct {
	// name titles the family's subtests.
	name string
	// diagnosticPrefix enumerates the family's registry rows without a second registry.
	diagnosticPrefix string
	// sentenceLead opens every sentence that states one prerequisite. A paragraph that
	// carries the lead is pinned prose.
	sentenceLead string
	// preflightCommand is the phase's preparation command. The retired full charge form is
	// a flag pair, not one spelling: the guidance may name the command, and it may name the
	// run control, but one sentence may not carry both.
	preflightCommand string
}

// retiredFullFlag is the run control both retired charge forms carried.
const retiredFullFlag = "--full"

// boundedActionFamilies is the one family inventory. A phase that migrates to the bounded
// evidence path joins this table, and every rule below then grades it.
func boundedActionFamilies() []boundedActionFamily {
	return []boundedActionFamily{
		{
			name:             "build",
			diagnosticPrefix: "bounded build action: ",
			sentenceLead:     "Build action requires",
			preflightCommand: "bench preflight build",
		},
		{
			name:             "review",
			diagnosticPrefix: "bounded review action: ",
			sentenceLead:     "Review action requires",
			preflightCommand: "bench preflight review",
		},
	}
}

// TestEvidenceBoundedActionFamilyInventory states the expected families apart from the
// table. Every rule below iterates that table, so a family dropped from it leaves each rule
// vacuously green. The names are compared in document order, so a drop, an addition, and a
// rename all red here, and no count stands in for the names.
func TestEvidenceBoundedActionFamilyInventory(t *testing.T) {
	want := []string{"build", "review"}
	var got []string
	for _, family := range boundedActionFamilies() {
		got = append(got, family.name)
	}
	if !slices.Equal(got, want) {
		t.Errorf("bounded action families = %v, want %v", got, want)
	}
}

// boundedActionFamilyNamed returns the one family the name selects. A case binds by name, so
// a reordered table cannot re-point it at another family.
func boundedActionFamilyNamed(t *testing.T, name string) boundedActionFamily {
	t.Helper()
	for _, family := range boundedActionFamilies() {
		if family.name == name {
			return family
		}
	}
	t.Fatalf("no bounded action family is named %q", name)
	return boundedActionFamily{}
}

func (f boundedActionFamily) anchors() []anchors.Anchor {
	return anchorsWithDiagnosticPrefix(f.diagnosticPrefix)
}

// guidance returns the one canonical reader the family pins. The rows' File fields are that
// path's single source, so a row that names a second file fails here.
func (f boundedActionFamily) guidance(t *testing.T, family []anchors.Anchor) string {
	t.Helper()
	files := map[string]bool{}
	for _, anchor := range family {
		files[anchor.File] = true
	}
	if len(files) != 1 {
		t.Fatalf("%s rows name %d guidance files, want one", f.name, len(files))
	}
	for file := range files {
		return file
	}
	return ""
}

// pinnedSentences returns every sentence of every paragraph that carries the family lead.
// It reads prose.Paragraphs, which is the paragraph rule and the sentence rule the prose
// gate itself applies, so the pin and the gate cannot disagree about where a paragraph or a
// sentence ends. Each sentence arrives with its whitespace collapsed.
func (f boundedActionFamily) pinnedSentences(text string) []string {
	var found []string
	for _, paragraph := range prose.Paragraphs(text) {
		if slices.ContainsFunc(paragraph, func(s string) bool { return strings.Contains(s, f.sentenceLead) }) {
			found = append(found, paragraph...)
		}
	}
	return found
}

// retiredFormPairs returns every sentence that names the family's preparation command and
// also carries the retired run control. The unit is one sentence of the same parser, so a
// wrap between the two halves hides nothing and a lawful run control in another sentence
// stays lawful. A heading and a frontmatter field hold no sentence; the literal Forbid row
// grades those, because an anchor reads the whole file.
func (f boundedActionFamily) retiredFormPairs(text string) []string {
	var found []string
	for _, paragraph := range prose.Paragraphs(text) {
		for _, sentence := range paragraph {
			if strings.Contains(sentence, f.preflightCommand) && strings.Contains(sentence, retiredFullFlag) {
				found = append(found, sentence)
			}
		}
	}
	return found
}

// wantedRows is the independent half of the registry pair: each family states the same five
// prerequisite diagnostics and the one retired-form refusal, apart from the registry, so a
// reworded or deleted row bites.
func (f boundedActionFamily) wantedRows() []struct {
	kind       anchors.Kind
	diagnostic string
} {
	action := strings.ToLower(f.name)
	return []struct {
		kind       anchors.Kind
		diagnostic string
	}{
		{anchors.Require, "permits " + action + " action without verified delivery"},
		{anchors.Require, "permits " + action + " action without available required context"},
		{anchors.Require, "permits " + action + " action without a current binding"},
		{anchors.Require, "permits " + action + " action without reviewer approval"},
		{anchors.Require, "permits " + action + " action without the complete task supplement"},
		{anchors.Forbid, "re-advertises the retired " + action + " full charge form"},
	}
}

// TestEvidenceBoundedActionGuidance grades the bounded action rows for CE101, CE102, CE103,
// CE140, CE141, CE142, CE143, CE144, CE145, and CE146. The count comes from the expectation
// table's length, and the mismatch reports rather than stops, so a deleted row also reaches
// the table.
func TestEvidenceBoundedActionGuidance(t *testing.T) {
	for _, family := range boundedActionFamilies() {
		t.Run(family.name, func(t *testing.T) {
			actionAnchors := family.anchors()
			want := family.wantedRows()
			if got := len(actionAnchors); got != len(want) {
				t.Errorf("%s anchor count = %d, want %d", family.name, got, len(want))
			}
			for _, expected := range want {
				found := false
				for _, anchor := range actionAnchors {
					if anchor.Kind == expected.kind && containsDiagnostic([]string{anchor.Diagnostic}, expected.diagnostic) {
						found = true
					}
				}
				if !found {
					t.Errorf("no %s row of the expected kind states %q", family.name, expected.diagnostic)
				}
			}
			runAnchorBites(t, actionAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
		})
	}
}

// TestEvidenceBoundedActionSentencesArePinned binds each guidance to its registry. A
// paragraph that states one prerequisite is pinned prose, so exactly one Require needle
// covers each of its sentences. A needle that keeps only its lead clause, a needle that
// repeats another row's text, and an unpinned sentence added to that paragraph all bite
// here. One parser reads the guidance and the needle, so a reflow that changes no word
// changes nothing here.
func TestEvidenceBoundedActionSentencesArePinned(t *testing.T) {
	for _, family := range boundedActionFamilies() {
		t.Run(family.name, func(t *testing.T) {
			rows := family.anchors()
			needles := map[string]int{}
			for _, anchor := range rows {
				if anchor.Kind == anchors.Require {
					needles[anchor.Needle]++
				}
			}
			for needle, count := range needles {
				if count != 1 {
					t.Errorf("%d Require rows share the needle %q, want one row each", count, needle)
				}
			}
			guidance := family.guidance(t, rows)
			sentences := family.pinnedSentences(guidanceText(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(guidance))))
			if len(sentences) == 0 {
				t.Fatalf("%s states no bounded %s action prerequisite", guidance, family.name)
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
		})
	}
}

// TestEvidenceBoundedActionRejectsTheRetiredPair grades the guidance half of CE110 and
// CE173. The retired form is the pair, so one guidance sentence may not name the phase's
// preparation command and also carry the run-control flag. Any flag order and any wrap
// between the two halves re-advertises the retired form and bites here. A phase's own
// full-run control names no preflight command, so it stays legal.
func TestEvidenceBoundedActionRejectsTheRetiredPair(t *testing.T) {
	for _, family := range boundedActionFamilies() {
		t.Run(family.name, func(t *testing.T) {
			guidance := family.guidance(t, family.anchors())
			pairs := family.retiredFormPairs(guidanceText(t, filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(guidance))))
			for _, sentence := range pairs {
				t.Errorf("%s states %q, which names %q and carries %q, and that re-advertises the retired full charge form",
					guidance, sentence, family.preflightCommand, retiredFullFlag)
			}
		})
	}
}

// TestEvidenceBoundedActionRulesBiteOnSyntheticText grades the two rules above apart from
// the shipped guidance. Each case feeds text through the rule's own predicate, so a rule
// that is weakened in place reds here while the shipped guidance still satisfies it. This
// guard is a case table rather than a canary fixture, because a canary mutates the guidance
// file and a weakened rule accepts every mutation of it. The name carries the family prefix,
// so the chunk's declared TestEvidence run reaches it.
func TestEvidenceBoundedActionRulesBiteOnSyntheticText(t *testing.T) {
	build := boundedActionFamilyNamed(t, "build")
	review := boundedActionFamilyNamed(t, "review")
	for _, tt := range []struct {
		name string
		rule func(string) []string
		text string
		want []string
	}{
		{
			name: "the lead opens the second sentence of a pinned paragraph",
			rule: build.pinnedSentences,
			text: "The ticket lands first. Build action requires reviewer approval.\n",
			want: []string{"The ticket lands first.", "Build action requires reviewer approval."},
		},
		{
			name: "a wrap inside a pinned sentence collapses",
			rule: build.pinnedSentences,
			text: "Build action requires verified delivery: act only after\nthis session reads it.\n",
			want: []string{"Build action requires verified delivery: act only after this session reads it."},
		},
		{
			name: "a neighbouring paragraph stays outside the pin",
			rule: build.pinnedSentences,
			text: "The reviewer may waive it.\n\nBuild action requires reviewer approval.\n",
			want: []string{"Build action requires reviewer approval."},
		},
		{
			name: "a paragraph with no lead states no prerequisite",
			rule: build.pinnedSentences,
			text: "The ticket lands first. It commits green.\n",
		},
		{
			name: "the review lead pins its own paragraph",
			rule: review.pinnedSentences,
			text: "The axis reads first. Review action requires reviewer approval.\n",
			want: []string{"The axis reads first.", "Review action requires reviewer approval."},
		},
		{
			name: "the build lead leaves the review paragraph alone",
			rule: review.pinnedSentences,
			text: "Build action requires reviewer approval.\n",
		},
		{
			name: "one sentence names the command and the flag",
			rule: build.retiredFormPairs,
			text: "Run `bench preflight build <slug> --full` now.\n",
			want: []string{"Run `bench preflight build <slug> --full` now."},
		},
		{
			name: "a wrap splits the command from the flag",
			rule: build.retiredFormPairs,
			text: "Run `bench preflight build <slug>` with the retired\n`--full` control now.\n",
			want: []string{"Run `bench preflight build <slug>` with the retired `--full` control now."},
		},
		{
			name: "the command alone advertises nothing",
			rule: build.retiredFormPairs,
			text: "Run `bench preflight build <slug>` now.\n",
		},
		{
			name: "the run control alone advertises nothing",
			rule: build.retiredFormPairs,
			text: "The `--full` run reaches the landing.\n",
		},
		{
			name: "the review command pairs with the retired control",
			rule: review.retiredFormPairs,
			text: "Run `bench preflight review <slug> --charge --full` now.\n",
			want: []string{"Run `bench preflight review <slug> --charge --full` now."},
		},
		{
			name: "a lawful diff run control stays lawful",
			rule: review.retiredFormPairs,
			text: "A historical review keeps `bench diff --full --commit <sha>` for the landed commit.\n",
		},
		{
			name: "the build command does not trip the review rule",
			rule: review.retiredFormPairs,
			text: "Run `bench preflight build <slug> --full` now.\n",
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
