package evidencecmd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/chargeevidence"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// seedLargeEvidence builds the expansion fixture: a 40 KB spec with 5000 fence entries, a
// ticket above the response bound, and escaped and multibyte text.
func seedLargeEvidence(t *testing.T) (root string, args []string, fence []string, ticket string) {
	t.Helper()
	root, slug := preflighttest.SeedConformant(t)
	var lines []string
	for i := 0; i < 5000; i++ {
		entry := fmt.Sprintf("internal/example/generated/%04d \"q\" 雪/", i)
		fence = append(fence, entry)
		lines = append(lines, "- `"+entry+"` (generated fence)")
	}
	spec := preflighttest.SpecBody(slug, lines...)
	if len(spec) < 40000 {
		t.Fatalf("spec fixture holds %d bytes, want at least 40 KB", len(spec))
	}
	preflighttest.MustWriteFile(t, "specs/"+slug+"/spec.md", spec)
	ticket = preflighttest.TicketDoc("One", "PF1", "PF2") + strings.Repeat("tab\there \"quote\" back\\slash 雪🚀\n", 2000)
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", ticket)
	return root, preflighttest.LegacyCommitted(t, root, slug, "large evidence"), fence, ticket
}

// TestEvidenceResponseBudget is CE13, CE131, CE132, CE133, CE134, CE135, CE138, and CE139.
func TestEvidenceResponseBudget(t *testing.T) {
	_, args, _, _ := seedLargeEvidence(t)
	identity, _, prepared := prepareEvidence(t, args)
	cases := map[string]string{"CE13 build preparation": prepared}
	manifestPages, sourcePages := 0, 0
	for _, page := range traverseEvidence(t, identity) {
		if page.row["stream"] == "manifest" {
			manifestPages++
			cases[fmt.Sprintf("CE131 manifest read %d", manifestPages)] = page.raw
		} else {
			sourcePages++
			cases[fmt.Sprintf("CE132 source read %d", sourcePages)] = page.raw
		}
	}
	if manifestPages < 1 || sourcePages < 10 {
		t.Fatalf("expansion fixture produced %d manifest and %d source pages", manifestPages, sourcePages)
	}
	verified, code := preflight.Command([]string{"evidence", identity, "--verify"})
	if code != 0 {
		t.Fatalf("verify exit = %d:\n%s", code, verified)
	}
	cases["CE133 verify"] = verified
	current, code := preflight.Command([]string{"evidence", identity, "--check-current"})
	if code != 0 {
		t.Fatalf("check-current exit = %d:\n%s", code, current)
	}
	cases["CE134 check-current"] = current
	usage, code := preflight.Command([]string{"build", strings.Repeat("s", 70000)})
	if code != 2 {
		t.Fatalf("oversized usage exit = %d", code)
	}
	cases["CE138 usage refusal"] = usage
	refusal, code := preflight.Command([]string{"evidence", "sha256:" + strings.Repeat("0", 64), "--cursor", "v1." + strings.Repeat("0", 64) + ".m.0.0"})
	if code != 1 {
		t.Fatalf("operational refusal exit = %d:\n%s", code, refusal)
	}
	cases["CE139 operational refusal"] = refusal
	// The review fixture seeds its own repository, so it runs after every read this
	// build artifact needs.
	cases["CE135 review preparation"] = largeReviewPreparation(t)
	for name, out := range cases {
		if len(out) > preflighttest.ResponseBudget || len(out) == 0 {
			t.Errorf("%s holds %d encoded bytes, want 1 to %d", name, len(out), preflighttest.ResponseBudget)
		}
	}
}

// TestEvidenceResponseBound is the guard half of CE13, CE131, CE132, CE133, CE134, CE138,
// and CE139. Under
// a lowered in-process limit, every bounded path's response becomes the bounded defect
// refusal, so a path that skips the shared guard turns its case red.
func TestEvidenceResponseBound(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	args := preflighttest.ChargeArgs(t, root, slug)
	identity, _, _ := prepareEvidence(t, args)
	source := "v1." + strings.TrimPrefix(identity, "sha256:") + ".s.1.0"
	restore := evidencecmd.SetResponseLimitForTest(16)
	defer restore()
	for _, test := range []struct {
		name string
		args []string
	}{
		{"CE13 build preparation", args},
		{"CE131 manifest read", []string{"evidence", identity}},
		{"CE132 source read", []string{"evidence", identity, "--cursor", source}},
		{"CE132 selected source read", []string{"evidence", identity, "--source", "s2"}},
		{"CE133 verify", []string{"evidence", identity, "--verify"}},
		{"CE134 check-current", []string{"evidence", identity, "--check-current"}},
		{"CE138 oversized operand usage", []string{"build", strings.Repeat("s", 2000)}},
		{"CE138 grammar usage", []string{"build", slug, "--unknown"}},
		{"CE138 operation usage", []string{"build", slug, "--full"}},
		{"CE139 operational refusal", []string{"evidence", "sha256:" + strings.Repeat("0", 64)}},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, code := preflight.Command(test.args)
			if code != 1 || !strings.HasPrefix(out, "error: response bound exceeded — ") || strings.Count(out, "\n") != 1 {
				t.Fatalf("%s under a 16-byte limit = (%d):\n%.300s", test.name, code, out)
			}
		})
	}
}

// TestEvidenceLargeMetadata is CE14: 5000 fence entries reconstruct without omission.
func TestEvidenceLargeMetadata(t *testing.T) {
	_, args, fence, ticket := seedLargeEvidence(t)
	identity, _, _ := prepareEvidence(t, args)
	_, sources := reconstructEvidence(t, identity, traverseEvidence(t, identity))
	var got []string
	for _, row := range preflighttest.TableRows(t, preflighttest.DecodeMap(t, sources["s1"]), "fence") {
		got = append(got, row["path"].(string))
	}
	want := append(append([]string{}, preflighttest.ConformantFence...), fence...)
	if strings.Join(got, "\n") != strings.Join(want, "\n") || sources["s2"] != ticket {
		t.Fatalf("metadata fence holds %d of %d entries; ticket equal=%t", len(got), len(want), sources["s2"] == ticket)
	}
}

// TestEvidenceBoundedErrors is CE15: an oversized operand is named by length and digest.
func TestEvidenceBoundedErrors(t *testing.T) {
	hostile := strings.Repeat("$(hostile)", 7000)
	for _, args := range [][]string{
		{"build", "example", "--charge", "--ticket", hostile},
		{"evidence", hostile},
		{"evidence", "sha256:" + strings.Repeat("a", 64), "--cursor", hostile},
	} {
		out, code := preflight.Command(args)
		want := "oversized operand bytes=70000 sha256=" + sha(hostile)
		if code != 2 || len(out) > preflighttest.ResponseBudget || !strings.Contains(out, want) || strings.Contains(out, "$(hostile)") {
			t.Fatalf("oversized operand = (%d, %d bytes):\n%.300s", code, len(out), out)
		}
	}
	out, code := preflight.Command([]string{"evidence", "sha256:" + strings.Repeat("g", 64)})
	if code != 2 || !strings.Contains(out, "invalid evidence identifier bytes=71 sha256=") {
		t.Fatalf("invalid identifier = (%d):\n%s", code, out)
	}
}

// publishCraftedEvidence publishes one valid review-mode artifact whose single generated
// source declares an argument holding scalar, so its manifest spans several fragments.
func publishCraftedEvidence(t *testing.T, root, scalar string) string {
	t.Helper()
	pack, err := chargeevidence.Build(chargeevidence.Candidate{
		Selection: chargeevidence.Selection{Mode: "review", Spec: "specs/example/spec.md", Base: strings.Repeat("a", 40), SourceTip: strings.Repeat("b", 40)},
		Metadata:  chargeevidence.Metadata{Checks: []string{"s2"}},
		Sources: []chargeevidence.SourceInput{{
			Role: "diff", Kind: chargeevidence.KindGenerated, Path: "diff", Required: true, Data: []byte("diff body\n"),
			Producer: &chargeevidence.Producer{Name: "diff", Version: "test", Cwd: ".", Arguments: []string{scalar}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	common, err := git.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	staged, err := chargeevidence.OpenStore(common, chargeevidence.StoreOptions{}).Stage(pack, chargeevidence.DefaultQuota, 1)
	if err != nil {
		t.Fatal(err)
	}
	identity, err := staged.Publish(1)
	if err != nil {
		t.Fatal(err)
	}
	return identity
}

// TestEvidenceLargeManifestScalar is CE6: a manifest cell above the response bound
// reconstructs exactly through bounded manifest fragments.
func TestEvidenceLargeManifestScalar(t *testing.T) {
	root, _ := preflighttest.SeedConformant(t)
	scalar := strings.Repeat("argument 雪 \"q\" ", 5000)
	identity := publishCraftedEvidence(t, root, scalar)
	pages := traverseEvidence(t, identity)
	manifest, sources := reconstructEvidence(t, identity, pages)
	arguments := preflighttest.TableRows(t, preflighttest.DecodeMap(t, manifest), "arguments")
	if len(scalar) <= preflighttest.ResponseBudget || len(arguments) != 1 || arguments[0]["value"] != scalar || sources["s2"] != "diff body\n" {
		t.Fatalf("large scalar did not reconstruct: %d arguments", len(arguments))
	}
	for _, page := range pages {
		if len(page.raw) > preflighttest.ResponseBudget {
			t.Fatalf("fragment holds %d bytes", len(page.raw))
		}
	}
}
