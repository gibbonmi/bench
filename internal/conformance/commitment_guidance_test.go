package conformance

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/anchors"
)

// commitmentGuidanceMutation edits one live guidance file: it replaces old, which occurs
// once, with replacement, and names the one commitment diagnostic that the edit must raise.
type commitmentGuidanceMutation struct {
	name, file, old, replacement, want string
}

// TestCommitmentGuidance grades the live guidance against the commitment family. The two
// spec-stage command lines must stay: deleting one raises its own diagnostic. Each retired
// grant must stay out: restoring its old sentence beside the live text raises its own
// diagnostic and no other. The sentences and diagnostics are written here independently of
// the registry, so a dropped or weakened row fails its case.
func TestCommitmentGuidance(t *testing.T) {
	const (
		guide      = ".bench/BENCH.md"
		drain      = ".agents/commands/bench-drain.md"
		finalCheck = ".agents/commands/bench-final-check.md"
		router     = ".agents/commands/bench.md"
		writeSpec  = ".agents/commands/bench-write-spec.md"
		prefix     = "commitment guidance: "
	)
	ticket := "The coordinator verifies the done-claim against the ticket's acceptance rows and the gate."
	sequence := "The drain does not write the `## Recommended sequence` section."
	learningFix := "**Delegate each light-path fix.**"
	remove := func(name, file, sentence, want string) commitmentGuidanceMutation {
		return commitmentGuidanceMutation{name: name, file: file, old: sentence, want: prefix + want}
	}
	restore := func(name, file, after, retired, want string) commitmentGuidanceMutation {
		return commitmentGuidanceMutation{name: name, file: file, old: after, replacement: after + "\n" + retired, want: prefix + want}
	}
	for _, mutation := range []commitmentGuidanceMutation{
		remove("staged spec", writeSpec, " A staged spec is planning work: `bench preflight build <slug> --plan-only` validates it, and its delivery waits for `bench commitment start`.", "spec authoring dropped the planning-only staged spec"),
		remove("start command as the handoff next command", writeSpec, "From the primary checkout, run `bench handoff --next` with the start command that `bench status` prints.", "spec authoring dropped the start command as the handoff next command"),

		restore("default implementation", drain, ticket, "For a drained item that meets the light-path observables, build the item in this session (\"implement now\") by default.", "drain restored the default implementation of a light-path item"),
		restore("declined row", drain, ticket, "Open a `ROADMAP.md` row only when the reviewer declines.", "drain restored the roadmap row only for a declined item"),
		restore("commitment admission", drain, ticket, "A drained light-path item is implement-now work only after `bench commitment approve` admits it and `bench commitment start` binds its worktree.", "drain restored the commitment admission of a light-path item"),
		restore("sequence rewrite", drain, sequence, "Rewrite the `## Recommended sequence` section: two or three numbered lines, each naming the item and the phase command to run.", "drain restored the sequence rewrite"),
		restore("severity rank", drain, sequence, "Rank rows by severity.", "drain restored the severity rank"),
		restore("actionable preference", drain, sequence, "Within an equal-severity class, choose actionable work over blocked work.", "drain restored the actionable-work preference"),
		restore("actionability tiebreaker", drain, sequence, "Only when rows are equally actionable, apply literal dependencies, then explicit reviewer pricing.", "drain restored dependencies and pricing as an actionability tiebreaker"),
		restore("occurrence tiebreaker", drain, sequence, "Only when all four stronger inputs tie, rank by descending occurrence count.", "drain restored the occurrence-count tiebreaker"),
		restore("defect and cost tiebreaker", drain, sequence, "When occurrence count also ties, apply the existing reproduced defect-over-feature rule, then cheapest-first cost rule.", "drain restored the defect and cost rules as an occurrence tiebreaker"),
		restore("refreshed sequence command", drain, sequence, "The recommended next command is the top line of the refreshed `## Recommended sequence`.", "drain restored the refreshed sequence as its next command"),
		restore("later drain closure", finalCheck, "Verified closure is part of delivery.", "Leave the roadmap and capture rows to `/bench-drain`; that phase owns the reconcile and the drain, and this duty never restates it.", "final check restored roadmap closure as a later drain"),
		restore("unadmitted learning fix", guide, learningFix, "A `bench learning` entry can have a light-path fix that needs no reviewer decision.", "operating guide restored the learning fix without admission"),
		restore("light-path commitment start", guide, "**Deliver only the committed outcome.**", "Delivery starts only through `bench commitment start` for the eligible outcome that `bench status` names, light path and fixes included.", "operating guide restored the commitment start for light-path work"),
		restore("active-outcome learning fix", guide, learningFix, "At any point in the workflow, `/bench-implement-spec` included, a `bench learning` entry can have a light-path fix that the active committed outcome needs.", "operating guide restored the learning fix for the active outcome only"),
		restore("in-build learning fix", guide, learningFix, "At any point in the workflow, `/bench-implement-spec` included, a `bench learning` entry's light-path fix ships at once only when the active outcome needs it.", "operating guide restored the in-build learning fix"),
		restore("mid-work defect fix", guide, "**Fix, don't park.**", "A small defect you find mid-work is not roadmap work: the fix lands in the active workflow as its own commit.", "operating guide restored the fix of every mid-work defect without admission"),
		restore("first sequence row", router, "State its outcome, then", "For roadmap work, run `bench roadmap` and take the first `sequence` row.", "router restored the first sequence row as its work"),
	} {
		t.Run(mutation.name, func(t *testing.T) {
			live := readIfExists(filepath.Join(NewHarness(t).KitRoot, filepath.FromSlash(mutation.file)))
			if count := strings.Count(live, mutation.old); count != 1 {
				t.Fatalf("%s holds %q %d times, want once", mutation.file, mutation.old, count)
			}
			root := t.TempDir()
			if got := commitmentGuidanceDiagnostics(t, root, mutation.file, live); len(got) != 0 {
				t.Fatalf("live %s raised %q", mutation.file, got)
			}
			got := commitmentGuidanceDiagnostics(t, root, mutation.file, strings.Replace(live, mutation.old, mutation.replacement, 1))
			if !reflect.DeepEqual(got, []string{mutation.want}) {
				t.Fatalf("mutated %s raised %q, want only %q", mutation.file, got, mutation.want)
			}
			t.Logf("observed red: %s", mutation.want)
		})
	}

	t.Run("one canonical owner", func(t *testing.T) {
		kit := NewHarness(t).KitRoot
		root := t.TempDir()
		for _, file := range []string{".bench/BENCH.md", "AGENTS.md", "README.md"} {
			body := readIfExists(filepath.Join(kit, filepath.FromSlash(file)))
			if file == "AGENTS.md" {
				body += "\n" + anchors.CommitmentRuleMarker + " Start only the committed outcome.\n"
			}
			writeGuidance(t, root, file, body)
		}
		want := fmt.Sprintf("shared rule duplicated in AGENTS.md (it must live only in .bench/BENCH.md): %q", anchors.CommitmentRuleMarker)
		if got := checkSharedRuleSingleSource(root); !reflect.DeepEqual(got, []string{want}) {
			t.Fatalf("restated rule raised %q, want only %q", got, want)
		}
	})
}

// commitmentGuidanceDiagnostics writes body as the one file below root and returns the
// commitment family's diagnostics for it.
func commitmentGuidanceDiagnostics(t *testing.T, root, file, body string) []string {
	t.Helper()
	writeGuidance(t, root, file, body)
	var diags []string
	for _, diag := range anchors.EvaluatePath(root, file).Diagnostics {
		if strings.HasPrefix(diag, anchors.CommitmentDiagnosticPrefix) {
			diags = append(diags, diag)
		}
	}
	return diags
}

func writeGuidance(t *testing.T, root, file, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(file))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}
