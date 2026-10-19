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
// the same rows as the build verdict, so neither charge operation adds a premise row. A
// prepared charge prints no check row, so the test reads Decide directly.
func TestDecideChargeRows(t *testing.T) {
	wantChecks := []string{
		"base-current", "paths-authorized",
		"tickets-parse", "completion-plan", "check-floor", "blockers-resolve", "writes-resolve",
		"fence-writes", "kit-pin", "binary-seal",
		"rows-owned", "rows-membership", "diff-nonempty",
	}
	for _, kind := range []evidencecmd.Kind{evidencecmd.KindPrepareEvidence, evidencecmd.KindCurrentEvidence} {
		f := baseFacts()
		f.Mode = modeBuild
		f.TicketsDirExists = true
		f.Operation = kind
		var got []string
		for _, c := range Decide(f).Checks {
			got = append(got, c.Check)
		}
		if !slices.Equal(got, wantChecks) {
			t.Errorf("Decide under operation %d checks = %v, want %v", kind, got, wantChecks)
		}
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
