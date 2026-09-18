package evidencecmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/chargesource"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// The expected review source order, metadata blocks, and prepared fields in this file are
// stated independently of the format registry and the review source policy, so an omission
// from either turns a case red.

// prepareReview seeds the shared review tree and prepares one review evidence artifact.
func prepareReview(t *testing.T) (root, identity string, row map[string]any, out string) {
	t.Helper()
	root, _, args := preflighttest.SeedReviewEvidence(t, false)
	identity, row, out = prepareEvidence(t, args)
	return root, identity, row, out
}

// TestEvidenceReviewAxes is CE69. One prepared artifact serves every axis, so each axis
// resolves the identical generated source identities. The metadata names the three axes
// against that one artifact, and each shared row points at a declared generated source.
func TestEvidenceReviewAxes(t *testing.T) {
	root, identity, _, _ := prepareReview(t)
	metadata := preflighttest.PublishedPack(t, root, identity).Metadata()

	axes := reviewrecord.Axes()
	if len(metadata.Charge) != len(axes) {
		t.Fatalf("charge rows = %d, want one per axis", len(metadata.Charge))
	}
	for i, axis := range axes {
		row := metadata.Charge[i]
		if row.Axis != axis || row.Access != "read-only" || row.Ticket != "" {
			t.Errorf("charge row %d = %#v, want axis %q read-only with no ticket", i, row, axis)
		}
	}

	manifest := preflighttest.PublishedPack(t, root, identity).Manifest()
	generated := map[string]chargeevidence.ManifestSource{}
	for _, source := range manifest.Sources {
		if source.Kind == chargeevidence.KindGenerated {
			generated[source.ID] = source
		}
	}
	wantKinds := []string{"diff", "consumers", "coverage"}
	if len(metadata.Shared) != len(wantKinds) {
		t.Fatalf("shared rows = %d, want %d", len(metadata.Shared), len(wantKinds))
	}
	for i, kind := range wantKinds {
		row := metadata.Shared[i]
		if row.Kind != kind {
			t.Errorf("shared row %d kind = %q, want %q", i, row.Kind, kind)
		}
		// The binding is the guarantee, not the existence of some generated source: each
		// row names the capture of its own kind, so a rotation of the bindings reds here.
		source, declared := generated[row.Source]
		if !declared {
			t.Errorf("shared row %q names %q, which declares no generated source", kind, row.Source)
			continue
		}
		if source.Role != kind {
			t.Errorf("shared row %q names %q, whose generated source is the %q capture", kind, row.Source, source.Role)
		}
	}
}

// TestEvidenceReviewProvenanceRows is CE73's manifest half. Every generated source declares
// its producer and the exact arguments that produced it, and no repository source does.
func TestEvidenceReviewProvenanceRows(t *testing.T) {
	root, identity, _, _ := prepareReview(t)
	manifest := preflighttest.PublishedPack(t, root, identity).Manifest()

	generated := map[string]bool{}
	for _, source := range manifest.Sources {
		if source.Kind == chargeevidence.KindGenerated {
			generated[source.ID] = true
		}
	}
	if len(manifest.Producers) != len(generated) || len(generated) != 3 {
		t.Fatalf("producers = %d for %d generated sources, want 3 each", len(manifest.Producers), len(generated))
	}
	wantProducers := map[string]string{"diff": "bench diff", "consumers": "bench consumers", "coverage": "bench coverage"}
	byID := map[string]chargeevidence.ProducerRow{}
	for _, producer := range manifest.Producers {
		if !generated[producer.Source] {
			t.Errorf("producer names %q, which declares no generated source", producer.Source)
		}
		byID[producer.Source] = producer
	}
	for _, source := range manifest.Sources {
		if !generated[source.ID] {
			continue
		}
		producer := byID[source.ID]
		if want := wantProducers[source.Role]; producer.Name != want {
			t.Errorf("source %s producer = %q, want %q", source.Role, producer.Name, want)
		}
		arguments := 0
		for _, argument := range manifest.Arguments {
			if argument.Source == source.ID {
				arguments++
			}
		}
		if arguments == 0 {
			t.Errorf("generated source %s declares no producer argument", source.Role)
		}
	}
}

// seedLargeReview seeds the shared review tree with one committed artifact larger than a
// bounded response, and returns the root with the preparation arguments pinned to that
// commit. Every case that needs an oversized review capture shares this one fixture.
func seedLargeReview(t *testing.T) (root string, args []string) {
	t.Helper()
	root, _, args = preflighttest.SeedReviewEvidence(t, false)
	preflighttest.MustWriteFile(t, "notes/large.txt", strings.Repeat("large review evidence\n", 3000))
	preflighttest.RunGit(t, "add", "-A")
	preflighttest.RunGit(t, "commit", "-q", "-m", "large review source")
	args[6] = preflighttest.RunGit(t, "rev-parse", "HEAD")
	return root, args
}

// TestEvidenceLargeReview is CE70. A review capture larger than one bounded response
// reconstructs byte for byte through the ordinary read stream.
func TestEvidenceLargeReview(t *testing.T) {
	root, args := seedLargeReview(t)

	identity, _, _ := prepareEvidence(t, args)
	_, sources := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	manifest := preflighttest.PublishedPack(t, root, identity).Manifest()

	oversized := 0
	for _, source := range manifest.Sources {
		if source.Bytes <= chargeevidence.ResponseLimit {
			continue
		}
		oversized++
		if got := sources[source.ID]; len(got) != source.Bytes || sha(got) != source.SHA256 {
			t.Fatalf("source %s reconstructed %d of %d bytes", source.Role, len(got), source.Bytes)
		}
	}
	if oversized == 0 {
		t.Fatal("no review source exceeded one bounded response")
	}
}

// largeReviewPreparation prepares one review artifact whose committed captures exceed the
// shared response bound, and returns the preparation response the budget case grades.
func largeReviewPreparation(t *testing.T) string {
	t.Helper()
	root, args := seedLargeReview(t)

	identity, _, out := prepareEvidence(t, args)
	committed := 0
	for _, source := range preflighttest.PublishedPack(t, root, identity).Manifest().Sources {
		committed += source.Bytes
	}
	if committed <= chargeevidence.ResponseLimit {
		t.Fatalf("review fixture committed %d bytes, want more than one bounded response", committed)
	}
	return out
}

// TestEvidenceReviewPreparedSchema is CE165. The response carries every registered prepared
// field, with the review mode and the frozen pair it prepared.
func TestEvidenceReviewPreparedSchema(t *testing.T) {
	root, _, args := preflighttest.SeedReviewEvidence(t, false)
	identity, row, _ := prepareEvidence(t, args)

	for _, field := range []string{
		"evidence", "mode", "base", "source_tip", "assignment", "selection", "metadata",
		"sources", "pages", "manifest_bytes", "response_complete", "delivery", "next",
	} {
		if _, present := row[field]; !present {
			t.Errorf("prepared row carries no %s field: %#v", field, row)
		}
	}
	if row["mode"] != "review" || row["base"] != args[4] || row["source_tip"] != args[6] {
		t.Errorf("prepared selection = %#v, want review over the frozen pair", row)
	}
	if row["delivery"] != "unverified" || row["response_complete"] != true {
		t.Errorf("prepared state = %#v, want an unverified complete response", row)
	}
	manifest := preflighttest.PublishedPack(t, root, identity).Manifest()
	if got := preparedCount(t, row, "sources"); got != len(manifest.Sources) {
		t.Errorf("prepared sources = %d, want %d", got, len(manifest.Sources))
	}
	if got := preparedCount(t, row, "pages"); got != len(manifest.Pages) {
		t.Errorf("prepared pages = %d, want %d", got, len(manifest.Pages))
	}
	if preparedCount(t, row, "manifest_bytes") == 0 {
		t.Error("prepared row declares no manifest bytes")
	}
}

// preparedCount reads one declared integer cell of the prepared row. The shared TOON
// adapter decodes every number as a float, so the cell converts once here.
func preparedCount(t *testing.T, row map[string]any, field string) int {
	t.Helper()
	value, ok := row[field].(float64)
	if !ok {
		t.Fatalf("prepared %s = %#v, want an integer", field, row[field])
	}
	return int(value)
}

// TestEvidenceReviewPreparationNext is CE166. The preparation returns the manifest-first
// retrieval command for the artifact it just published, never a source or verify form.
func TestEvidenceReviewPreparationNext(t *testing.T) {
	_, identity, row, _ := prepareReview(t)
	want := "bench preflight evidence " + identity
	if row["next"] != want {
		t.Fatalf("prepared next = %q, want %q", row["next"], want)
	}
	if !strings.HasPrefix(identity, chargeevidence.IdentityPrefix) {
		t.Fatalf("prepared identity %q carries no %q prefix", identity, chargeevidence.IdentityPrefix)
	}
}

// TestEvidenceReviewMetadataSchema is CE164. The metadata source carries the registered
// review blocks: one charge row per axis, the declared fence, the check and return sources,
// the shared capture bindings, and the completion facts. Review selects no ticket, so its
// writes and coverage blocks stay empty.
func TestEvidenceReviewMetadataSchema(t *testing.T) {
	root, identity, _, _ := prepareReview(t)
	pack := preflighttest.PublishedPack(t, root, identity)
	metadata := pack.Metadata()
	manifest := pack.Manifest()

	// The fence is graded against the entries the seeded spec declares, so a truncated or
	// reordered fence reds rather than passing on a non-empty list.
	if got, want := strings.Join(metadata.Fence, "\n"), strings.Join(preflighttest.ReviewFence(), "\n"); got != want {
		t.Errorf("metadata fence =\n%s\nwant\n%s", got, want)
	}
	if len(metadata.Writes) != 0 || len(metadata.Coverage) != 0 {
		t.Errorf("review metadata declares %d writes and %d coverage rows, want none",
			len(metadata.Writes), len(metadata.Coverage))
	}
	if len(metadata.Completion) != 1 {
		t.Fatalf("completion rows = %d, want one", len(metadata.Completion))
	}
	if record := metadata.Completion[0].Record; record != "reviews/example.md" {
		t.Errorf("completion record = %q, want reviews/example.md", record)
	}

	byID := map[string]chargeevidence.ManifestSource{}
	for _, source := range manifest.Sources {
		byID[source.ID] = source
	}
	if len(metadata.Checks) != 1 || byID[metadata.Checks[0]].Path != chargesource.ReviewSkill {
		t.Errorf("checks = %v, want the review skill source", metadata.Checks)
	}
	wantReturns := []string{chargesource.ReviewPhase, chargesource.DelegateSkill, chargesource.DelegateProcedure}
	if len(metadata.Returns) != len(wantReturns) {
		t.Fatalf("returns = %d, want %d", len(metadata.Returns), len(wantReturns))
	}
	for i, want := range wantReturns {
		if got := byID[metadata.Returns[i]].Path; got != want {
			t.Errorf("return %d = %q, want %q", i, got, want)
		}
	}
}

// TestEvidenceReviewGrammar is CE163. Review preparation accepts its declared form and
// refuses every other spelling before any collector starts.
func TestEvidenceReviewGrammar(t *testing.T) {
	_, slug, valid := preflighttest.SeedReviewEvidence(t, false)
	for _, test := range []struct {
		name string
		args []string
		want string
	}{
		{"ticket selector", append(append([]string{}, valid...), "--ticket", "one.md"), "unknown argument: --ticket"},
		{"read-only flag", append(append([]string{}, valid...), "--verify"), "cannot be combined"},
		{"source flag", append(append([]string{}, valid...), "--source", "s2"), "cannot be combined"},
		{"cursor flag", append(append([]string{}, valid...), "--cursor", "v1.x"), "unknown argument: --cursor"},
		{"retired full", append(append([]string{}, valid...), "--full"), "unknown argument: --full"},
		{"duplicate base", append(append([]string{}, valid...), "--base", valid[4]), "unknown argument: --base"},
		{"extra operand", append(append([]string{}, valid...), "extra"), "usage"},
		{"absent base", []string{"review", slug, "--charge", "--source-tip", valid[6]}, "--charge requires review"},
		{"absent tip", []string{"review", slug, "--charge", "--base", valid[4]}, "--charge requires review"},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, code := preflight.Command(test.args)
			if code != 2 || !strings.Contains(out, test.want) {
				t.Fatalf("%s = (%d), want exit 2 containing %q:\n%s", test.name, code, test.want, out)
			}
		})
	}
}

// TestEvidenceRemovedFull is CE108 and CE109. The run control is gone from the flag
// registry, so every charge form that names it refuses through the bounded usage path, and
// no executable next action any preparation prints can name it.
func TestEvidenceRemovedFull(t *testing.T) {
	_, slug, reviewArgs := preflighttest.SeedReviewEvidence(t, false)
	_, row, out := prepareEvidence(t, reviewArgs)
	if strings.Contains(out, "--full") {
		t.Errorf("the review preparation response advertises --full:\n%s", out)
	}
	if next := row["next"].(string); strings.Contains(next, "--full") {
		t.Errorf("the review next action advertises --full: %q", next)
	}

	buildRoot, buildSlug := preflighttest.SeedConformant(t)
	_, buildRow, buildOut := prepareEvidence(t, preflighttest.ChargeArgs(t, buildRoot, buildSlug))
	if strings.Contains(buildOut, "--full") {
		t.Errorf("the build preparation response advertises --full:\n%s", buildOut)
	}
	if next := buildRow["next"].(string); strings.Contains(next, "--full") {
		t.Errorf("the build next action advertises --full: %q", next)
	}

	for _, args := range [][]string{
		append(append([]string{}, reviewArgs...), "--full"),
		{"review", slug, "--full"},
	} {
		if out, code := preflight.Command(args); code != 2 || !strings.Contains(out, "unknown argument: --full") {
			t.Fatalf("%v = (%d):\n%s", args, code, out)
		}
	}
}

// TestEvidenceHelpInventory is CE167. The public help projects the review preparation and
// names no retired charge route, and its rows come from the operation registry.
func TestEvidenceHelpInventory(t *testing.T) {
	rows := evidencecmd.HelpRows()
	var review string
	for _, row := range rows {
		if strings.Contains(row.Suffix, "--full") {
			t.Errorf("help row advertises the retired run control: %q", row.Suffix)
		}
		if strings.HasPrefix(row.Suffix, " review ") && strings.Contains(row.Suffix, "--charge") {
			review = row.Suffix
		}
	}
	if review == "" {
		t.Fatalf("help rows project no review preparation: %#v", rows)
	}
	for _, want := range []string{"--base", "--source-tip"} {
		if !strings.Contains(review, want) {
			t.Errorf("review help row %q omits %q", review, want)
		}
	}
}
