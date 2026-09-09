package preflight

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/consumers"
)

// The charge row is the packet's contract, and a Contains needle over the whole packet
// grades none of it: the co-located sources, omitted, and evidence tables satisfy every
// handle needle on their own. The tests below decode the packet and compare each named
// cell against a value the test computes, so an emptied or swapped cell reds.

const chargeFixtureAssignment = "00000000000000000000000000000001"

// conformantFence is the ownership fence seedConformant's spec declares, in document
// order. reviewFixtureFence adds the entries seedReviewEvidence appends.
var conformantFence = []string{
	"internal/example/",
	"reviews/example.md",
	".agents/skills/bench-craft-delegate/",
	".agents/commands/bench-implement-spec.md",
}

var reviewFixtureFence = append(append([]string{}, conformantFence...),
	"target/", "edited/", "outside/", "notes/",
	".agents/skills/bench-craft-review/",
	".agents/commands/bench-review-implementation.md",
)

// fixtureHandle recomputes one source handle from the file on disk, independently of the
// renderer. A cell that names the wrong source therefore reds.
func fixtureHandle(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %q: %v", path, err)
	}
	return chargeSource{path: path, data: data}.handle()
}

func fixtureHandles(t *testing.T, paths ...string) string {
	t.Helper()
	handles := make([]string, len(paths))
	for i, path := range paths {
		handles[i] = fixtureHandle(t, path)
	}
	return strings.Join(handles, "; ")
}

func chargeCell(t *testing.T, row map[string]any, name string) string {
	t.Helper()
	value, present := row[name]
	if !present {
		t.Fatalf("charge row carries no %s cell: %#v", name, row)
	}
	if value == nil {
		return ""
	}
	text, ok := value.(string)
	if !ok {
		t.Fatalf("charge cell %s = %T, want a string", name, value)
	}
	return text
}

func assertChargeCells(t *testing.T, row map[string]any, want map[string]string) {
	t.Helper()
	for name, expected := range want {
		if got := chargeCell(t, row, name); got != expected {
			t.Errorf("charge cell %s = %q, want %q", name, got, expected)
		}
	}
}

func TestBuildChargeRowCells(t *testing.T) {
	root, slug := seedConformant(t)
	args := chargeArgs(t, root, slug, true)
	out, code := Command(args)
	if code != 0 {
		t.Fatalf("build charge exit = %d:\n%s", code, out)
	}
	rows := tableRows(t, decodeMap(t, out), "charge")
	if len(rows) != 1 {
		t.Fatalf("charge rows = %d, want 1", len(rows))
	}
	ticket := "specs/" + slug + "/tickets/one.md"
	spec := "specs/" + slug + "/spec.md"
	assertChargeCells(t, rows[0], map[string]string{
		"assignment": chargeFixtureAssignment,
		"checkout":   root,
		"base":       args[6],
		"source_tip": args[8],
		"fence":      strings.Join(conformantFence, ", "),
		"ticket":     fixtureHandle(t, ticket),
		"writes":     "specs",
		"evidence":   fixtureHandles(t, ticket, spec, delegateSkill, buildPhase, delegateProcedure),
		"checks":     fixtureHandles(t, ticket, buildPhase),
		"return":     fixtureHandles(t, delegateSkill, delegateProcedure),
		"complete":   "true",
		"next":       "",
	})
}

func TestBuildChargeCellsCarryDistinctFacts(t *testing.T) {
	root, slug := seedConformant(t)
	out, code := Command(chargeArgs(t, root, slug, true))
	if code != 0 {
		t.Fatalf("build charge exit = %d:\n%s", code, out)
	}
	rows := tableRows(t, decodeMap(t, out), "charge")
	assertDistinctCells(t, rows[0], "fence", "ticket", "evidence", "checks", "return")
}

func TestReviewChargeRowCells(t *testing.T) {
	_, slug, args := seedReviewEvidence(t, false)
	out, code := Command(args)
	if code != 0 {
		t.Fatalf("review charge exit = %d:\n%s", code, out)
	}
	document := decodeMap(t, out)
	rows := tableRows(t, document, "charge")
	if len(rows) != 3 {
		t.Fatalf("charge rows = %d, want one per axis", len(rows))
	}
	root, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// The shared evidence handle is the concatenation of the three collector payloads.
	// The test recomputes it from the packet's own retrieved evidence, so it never reads
	// the renderer's value back.
	shared := chargeSource{path: "shared-review-evidence", data: []byte(
		evidenceContent(t, out, "diff") +
			evidenceContent(t, out, "consumers") +
			evidenceContent(t, out, "coverage"))}
	spec := "specs/" + slug + "/spec.md"
	for i, axis := range []string{"Standards", "Spec", "Coverage"} {
		assertChargeCells(t, rows[i], map[string]string{
			"axis":       axis,
			"assignment": chargeFixtureAssignment,
			"checkout":   root,
			"base":       args[4],
			"source_tip": args[6],
			"fence":      strings.Join(reviewFixtureFence, ", "),
			"ticket":     fixtureHandle(t, spec),
			"writes":     "read-only",
			"evidence":   shared.handle(),
			"checks":     fixtureHandle(t, reviewSkill),
			"return":     fixtureHandles(t, reviewPhase, delegateSkill, delegateProcedure),
			"complete":   "true",
			"next":       "",
		})
		assertDistinctCells(t, rows[i], "fence", "ticket", "evidence", "checks", "return")
	}
}

// assertDistinctCells proves no two named cells hold one value. Two cells that agree
// carry one fact between them, and the second grades nothing.
func assertDistinctCells(t *testing.T, row map[string]any, names ...string) {
	t.Helper()
	seen := map[string]string{}
	for _, name := range names {
		value := chargeCell(t, row, name)
		if value == "" {
			t.Errorf("charge cell %s is empty", name)
			continue
		}
		if other, duplicate := seen[value]; duplicate {
			t.Errorf("charge cells %s and %s hold the identical value %q", other, name, value)
			continue
		}
		seen[value] = name
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
// packet reaches three collectors that each resolve their own root, so a nested run must
// render the identical bytes.
func TestPreparedFormsIgnoreNestedWorkingDirectory(t *testing.T) {
	t.Run("review charge", func(t *testing.T) {
		root, _, args := seedReviewEvidence(t, false)
		assertNestedRunMatches(t, root, filepath.Join(root, "target"), args)
	})

	t.Run("build charge", func(t *testing.T) {
		root, slug := seedConformant(t)
		args := chargeArgs(t, root, slug, true)
		assertNestedRunMatches(t, root, filepath.Join(root, "internal", slug), args)
	})

	t.Run("compact build charge", func(t *testing.T) {
		root, slug := seedConformant(t)
		args := chargeArgs(t, root, slug, false)
		assertNestedRunMatches(t, root, filepath.Join(root, "internal", slug), args)
	})

	t.Run("propose writes", func(t *testing.T) {
		root, slug := seedConformant(t)
		activeAssignment(t, root, root)
		args := []string{
			"build", slug, "--propose-writes", "--ticket", "one.md",
			"--base", runGit(t, "rev-parse", "main"),
			"--source-tip", runGit(t, "rev-parse", "HEAD"),
		}
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
