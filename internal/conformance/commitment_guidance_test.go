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

// TestCommitmentGuidance grades the live guidance against the commitment family. The
// canonical rule and each phase route must stay: deleting one sentence raises its own
// diagnostic. Each retired grant must stay out: restoring its old sentence beside the live
// text raises its own diagnostic and no other. The sentences and diagnostics are written
// here independently of the registry, so a dropped or weakened row fails its case.
func TestCommitmentGuidance(t *testing.T) {
	const (
		guide      = ".bench/BENCH.md"
		drain      = ".agents/commands/bench-drain.md"
		finalCheck = ".agents/commands/bench-final-check.md"
		router     = ".agents/commands/bench.md"
		prefix     = "commitment guidance: "
	)
	ticket := "Verify the diff against the ticket's acceptance rows and the gate."
	sequence := "The drain does not write the `## Recommended sequence` section."
	learning := "A learning entry with a light-path fix follows the learning-fix rule in `.bench/BENCH.md`."
	remove := func(name, file, sentence, want string) commitmentGuidanceMutation {
		return commitmentGuidanceMutation{name: name, file: file, old: sentence, want: prefix + want}
	}
	restore := func(name, file, after, retired, want string) commitmentGuidanceMutation {
		return commitmentGuidanceMutation{name: name, file: file, old: after, replacement: after + "\n" + retired, want: prefix + want}
	}
	for _, mutation := range []commitmentGuidanceMutation{
		remove("canonical rule", guide, "**Deliver only the committed outcome.** ", "operating guide dropped the canonical delivery commitment rule"),
		remove("start route", guide, "Delivery starts only through `bench commitment start` for the eligible outcome that `bench status` names, light path and fixes included. ", "operating guide dropped the committed-outcome start route"),
		remove("uncommitted intake", guide, "Any other finding, idea, learning, or drained item stays uncommitted intake, and minimal support that the active outcome needs stays in that outcome. ", "operating guide dropped uncommitted intake"),
		remove("explicit displacement", guide, "Only my explicit direction changes the commitment, through `bench commitment plan` and then `bench commitment approve`; no drain, label, score, or count displaces it. ", "operating guide dropped explicit displacement"),
		remove("purpose priority", guide, "When you propose work, put confirmed defects first, then refactors, then features, by purpose rather than label; dependencies and the approved order govern execution. ", "operating guide dropped purpose-based defect and refactor priority"),
		remove("blocker report", guide, " When the active outcome cannot continue, run `bench commitment block` with the reason and tell me; the obligation stays.", "operating guide dropped the blocker report"),
		remove("drain admission route", drain, "The commitment rule in `.bench/BENCH.md` decides whether intake starts; a drain approval never admits it.", "drain dropped its route to the commitment rule"),
		remove("drain implement-now admission", drain, "A drained light-path item is implement-now work only after `bench commitment approve` admits it and `bench commitment start` binds its worktree.", "drain dropped the admission route for implement-now work"),
		remove("drain sequence owner", drain, sequence, "drain dropped the commitment ownership of the recommended sequence"),
		remove("drain classification limit", drain, "It informs a commitment proposal and never reorders committed work.", "drain dropped the classification limit"),
		remove("final-check closure", finalCheck, "Verified closure is part of delivery. ", "final check dropped verified closure from delivery"),
		remove("implementation start", ".agents/commands/bench-implement-spec.md", "It declares the line, starts its committed outcome through `bench commitment start`, and works vertical slices at the pre-agreed seams.", "implementation dropped its commitment start"),
		remove("staged spec", ".agents/commands/bench-write-spec.md", " A staged spec is planning work: `bench preflight build <slug> --plan-only` validates it, and its delivery waits for `bench commitment start`.", "spec authoring dropped the planning-only staged spec"),
		remove("debug route", ".agents/commands/bench-debug.md", " Diagnosis needs no commitment start, but each commit of the repro or the fix follows the commitment rule in `.bench/BENCH.md`.", "debug dropped its route to the commitment rule"),
		remove("router outlook", router, "take its `commitment_outlook` row", "router dropped the commitment outlook"),
		remove("setup commitment", ".agents/commands/bench-setup-repo.md", "Delivery work waits for the initial commitment. ", "setup dropped the initial commitment"),

		restore("default implementation", drain, ticket, "For a drained item that meets the light-path observables, build the item in this session (\"implement now\") by default.", "drain restored the default implementation of a light-path item"),
		restore("declined row", drain, ticket, "Open a `ROADMAP.md` row only when the reviewer declines.", "drain restored the roadmap row only for a declined item"),
		restore("direct learning fix", drain, learning, "A learning entry with a light-path fix goes to the write delegate that `.bench/BENCH.md` names, and its verdict closes the entry by implementation.", "drain restored the direct implementation of each light-path learning fix"),
		restore("sequence rewrite", drain, sequence, "Rewrite the `## Recommended sequence` section: two or three numbered lines, each naming the item and the phase command to run.", "drain restored the sequence rewrite"),
		restore("severity rank", drain, sequence, "Rank rows by severity.", "drain restored the severity rank"),
		restore("actionable preference", drain, sequence, "Within an equal-severity class, choose actionable work over blocked work.", "drain restored the actionable-work preference"),
		restore("actionability tiebreaker", drain, sequence, "Only when rows are equally actionable, apply literal dependencies, then explicit reviewer pricing.", "drain restored dependencies and pricing as an actionability tiebreaker"),
		restore("occurrence tiebreaker", drain, sequence, "Only when all four stronger inputs tie, rank by descending occurrence count.", "drain restored the occurrence-count tiebreaker"),
		restore("defect and cost tiebreaker", drain, sequence, "When occurrence count also ties, apply the existing reproduced defect-over-feature rule, then cheapest-first cost rule.", "drain restored the defect and cost rules as an occurrence tiebreaker"),
		restore("refreshed sequence command", drain, sequence, "The recommended next command is the top line of the refreshed `## Recommended sequence`.", "drain restored the refreshed sequence as its next command"),
		restore("later drain closure", finalCheck, "Verified closure is part of delivery.", "Leave the roadmap and capture rows to `/bench-drain`; that phase owns the reconcile and the drain, and this duty never restates it.", "final check restored roadmap closure as a later drain"),
		restore("unadmitted learning fix", guide, "**Delegate a light-path fix for a learning.**", "A `bench learning` entry can have a light-path fix that needs no reviewer decision.", "operating guide restored the learning fix without admission"),
		restore("drain implementation", guide, "Parked ideas land in `capture/IDEAS.md`.", "They graduate to the board only through a reviewed `/bench-drain` drain, or close by implementation during that same drain.", "operating guide restored implementation inside a drain"),
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
