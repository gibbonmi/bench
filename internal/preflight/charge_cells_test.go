package preflight

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/consumers"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// TestDecideChargeRows covers CR06. A build charge and its current-evidence check decide
// the same rows as a build verdict over the same source pair, so neither charge operation
// adds a premise row. Each case holds the facts that its real gather holds: both read an
// explicit source range, and only the charge pins its source tip, which adds tip-current.
// A prepared charge prints no check row, so the test reads Decide directly.
func TestDecideChargeRows(t *testing.T) {
	verdictRows := []string{
		"paths-authorized",
		"tickets-parse", "completion-plan", "blockers-resolve", "writes-resolve", "writes-landed",
		"kit-pin", "binary-seal", "map-sources",
		"rows-owned", "rows-membership", "chunk-rows", "diff-nonempty",
	}
	for _, test := range []struct {
		kind   evidencecmd.Kind
		pinned string
		want   []string
	}{
		{evidencecmd.KindPrepareEvidence, "beef", append([]string{"base-current", "tip-current"}, verdictRows...)},
		{evidencecmd.KindCurrentEvidence, "", append([]string{"base-current"}, verdictRows...)},
	} {
		f := baseFacts()
		f.Mode = modeBuild
		f.TicketsDirExists = true
		f.ExplicitSourceRange = true
		f.SourceBase, f.SourceTip, f.PinnedSourceTip = "cafe", "beef", test.pinned
		f.Operation = test.kind
		verdict := Decide(f)
		var got []string
		for _, c := range verdict.Checks {
			got = append(got, c.Check)
		}
		if verdict.Red || !slices.Equal(got, test.want) {
			t.Errorf("Decide under operation %d = red %t, checks %v, want green %v", test.kind, verdict.Red, got, test.want)
		}
	}
}

// TestDecideHasNoCheckFloorRow covers BC32. A build verdict over tickets prints no
// check-floor row, because the chunk checkpoint is the only grade of a plain check.
func TestDecideHasNoCheckFloorRow(t *testing.T) {
	f := baseFacts()
	f.Mode = modeBuild
	f.TicketsDirExists = true
	if row, ok := checkRow(Decide(f), "check-floor"); ok {
		t.Fatalf("build verdict holds the retired row %+v", row)
	}
}

// TestCompleteConsumerEvidenceMatchesLiveConsumerMeta grades the one review refusal
// branch that guards a live cross-package contract: completeConsumerEvidence reads the
// meta field names internal/consumers declares in metaFields. A rename there refuses
// every review charge, so this test reds on the rename instead of leaving it to a
// confusing charge failure.
func TestCompleteConsumerEvidenceMatchesLiveConsumerMeta(t *testing.T) {
	_, _, args := seedReviewEvidence(t, false)
	out, code := consumers.CommandWithVersion("fixture-version")([]string{
		"--changed", "--base", args[4], "--source-tip", args[6], "--full",
	})
	if code != 0 {
		t.Fatalf("live consumer evidence = (%d):\n%s", code, out)
	}
	complete, err := completeConsumerEvidence(out)
	if err != nil || !complete {
		t.Fatalf("live consumer meta = (%t, %v), want complete with no error:\n%s", complete, err, out)
	}
	renamed := strings.Replace(out, ",truncated}", ",omitted}", 1)
	if renamed == out {
		t.Fatalf("live consumer meta header carries no truncated field:\n%s", out)
	}
	if _, err := completeConsumerEvidence(renamed); err == nil ||
		!strings.Contains(err.Error(), "meta.truncated") {
		t.Fatalf("renamed meta field = %v, want a meta.truncated refusal", err)
	}
}

// TestPreparedFormsIgnoreNestedWorkingDirectory closes the hostile-input class of a cwd
// deeper than the repo root. Every prepared form resolves the root itself, and the review
// preparation reaches three collectors that each resolve their own root, so a nested run
// must render the identical bytes.
func TestPreparedFormsIgnoreNestedWorkingDirectory(t *testing.T) {
	t.Run("review charge", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		assertNestedRunMatches(t, root, filepath.Join(root, "target"), args)
	})

	t.Run("build preparation", func(t *testing.T) {
		root, slug := preflighttest.SeedConformant(t)
		args := preflighttest.ChargeArgs(t, root, slug)
		assertNestedRunMatches(t, root, filepath.Join(root, "internal", slug), args)
	})
}

func assertNestedRunMatches(t *testing.T, root, nested string, args []string) {
	t.Helper()
	fromRoot, code := Command(args)
	if code != 0 {
		t.Fatalf("run at %s = (%d):\n%s", root, code, fromRoot)
	}
	t.Chdir(nested)
	fromNested, nestedCode := Command(args)
	if nestedCode != code || fromNested != fromRoot {
		t.Fatalf("run from %s = (%d), want the root run (%d):\nnested:\n%s\nroot:\n%s",
			nested, nestedCode, code, fromNested, fromRoot)
	}
}
