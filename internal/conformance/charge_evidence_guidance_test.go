package conformance

import (
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
	"github.com/gibbonmi/bench/internal/prose"
)

// boundedActionFamily is one bounded action guidance family. Each phase pins its
// artifact-integrity check and its current-action binding over its own canonical reader, and
// it forbids the retired full charge form, so the rules below are written once and bound to a
// family rather than copied per phase.
type boundedActionFamily struct {
	// name titles the family's subtests.
	name string
	// diagnosticPrefix enumerates the family's registry rows without a second registry.
	diagnosticPrefix string
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
			preflightCommand: "bench preflight build",
		},
		{
			name:             "review",
			diagnosticPrefix: "bounded review action: ",
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

// wantedRows is the independent half of the registry pair: each family states its integrity
// check, its current binding, and the one retired-form refusal, apart from the registry, so a
// reworded or deleted row bites. A diagnostic names the reader by its base name, as every
// registered row does, and guidance supplies that reader.
func (f boundedActionFamily) wantedRows(guidance string) []wantedRow {
	action := strings.ToLower(f.name)
	reader := f.diagnosticPrefix + path.Base(guidance)
	return []wantedRow{
		{anchors.Require, reader + " dropped the separate bench preflight evidence --verify integrity check"},
		{anchors.Require, reader + " permits " + action + " action without a current binding"},
		{anchors.Forbid, reader + " re-advertises the retired " + action + " full charge form"},
	}
}

// wantedRow is one expected registry row: its kind and its complete diagnostic.
type wantedRow struct {
	kind       anchors.Kind
	diagnostic string
}

// TestEvidenceBoundedActionGuidance grades the bounded action rows for CE101, CE102, CE103,
// CE140, CE141, CE142, CE143, CE144, CE145, and CE146. The count comes from the expectation
// table's length, and the mismatch reports rather than stops, so a deleted row also reaches
// the table.
func TestEvidenceBoundedActionGuidance(t *testing.T) {
	for _, family := range boundedActionFamilies() {
		t.Run(family.name, func(t *testing.T) {
			actionAnchors := family.anchors()
			guidance := family.guidance(t, actionAnchors)
			want := family.wantedRows(guidance)
			if got := len(actionAnchors); got != len(want) {
				t.Errorf("%s anchor count = %d, want %d", family.name, got, len(want))
			}
			for _, expected := range want {
				requireRegisteredRow(t, actionAnchors, expected.kind, guidance, expected.diagnostic)
			}
			runAnchorBites(t, actionAnchors, func(anchor anchors.Anchor) string { return anchor.Diagnostic })
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
			h := NewHarness(t)
			pairs := family.retiredFormPairs(guidanceText(t, filepath.Join(h.KitRoot, filepath.FromSlash(guidance))))
			for _, sentence := range pairs {
				t.Errorf("%s states %q, which names %q and carries %q, and that re-advertises the retired full charge form",
					guidance, sentence, family.preflightCommand, retiredFullFlag)
			}
		})
	}
}

// TestEvidenceBoundedActionRulesBiteOnSyntheticText grades the retired-pair rule above apart
// from the shipped guidance. Each case feeds text through the rule's own predicate, so a rule
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

// guidanceText reads one canonical reader with its line endings normalized.
func guidanceText(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Clean(name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return strings.ReplaceAll(string(data), "\r\n", "\n")
}
