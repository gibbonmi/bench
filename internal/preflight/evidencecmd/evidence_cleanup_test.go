package evidencecmd_test

import (
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/preflight"
	"github.com/gibbonmi/bench/internal/preflight/evidencecmd"
	"github.com/gibbonmi/bench/internal/preflight/preflighttest"
)

// cleanupPlan runs one cleanup plan page and returns its orientation row and target rows.
func cleanupPlan(t *testing.T, args ...string) (map[string]any, []map[string]any) {
	t.Helper()
	out, code := preflight.Command(append([]string{"evidence-clean"}, args...))
	if code != 0 {
		t.Fatalf("cleanup plan = (%d):\n%s", code, out)
	}
	decoded := preflighttest.DecodeMap(t, out)
	rows := preflighttest.TableRows(t, decoded, "cleanup")
	if len(rows) != 1 {
		t.Fatalf("cleanup plan carries %d orientation rows:\n%s", len(rows), out)
	}
	return rows[0], preflighttest.TableRows(t, decoded, "targets")
}

// TestEvidenceCleanupSchema is CE169. The plan orientation row carries every registered
// cleanup field, and it names the successor that authorizes the apply.
func TestEvidenceCleanupSchema(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))

	row, _ := cleanupPlan(t)
	for _, field := range []string{"fingerprint", "targets", "bytes", "response_complete", "stream_end", "next"} {
		if _, present := row[field]; !present {
			t.Errorf("cleanup row carries no %s field: %#v", field, row)
		}
	}
	fingerprint, _ := row["fingerprint"].(string)
	if !strings.HasPrefix(fingerprint, "sha256:") || len(fingerprint) != 71 {
		t.Errorf("cleanup fingerprint = %q, want a prefixed digest", fingerprint)
	}
	if row["response_complete"] != true || row["stream_end"] != true {
		t.Errorf("a one-page plan = %#v, want a complete terminal page", row)
	}
	if want := evidencecmd.CleanCommand + " --apply " + fingerprint; row["next"] != want {
		t.Errorf("cleanup next = %#v, want %q", row["next"], want)
	}
	// CE80: the plan enumerates, so the artifact it named is still published.
	if len(preflighttest.PublishedPacks(t, root)) != 1 {
		t.Errorf("planning changed the store to %v", preflighttest.PublishedPacks(t, root))
	}
}

// TestEvidenceCleanupTargetsSchema is CE170 and the plan half of CE78. Each target row
// carries every registered field and a retrievable identity, and the row count matches the
// orientation count.
func TestEvidenceCleanupTargetsSchema(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))

	row, targets := cleanupPlan(t)
	if len(targets) != 1 {
		t.Fatalf("cleanup planned %d targets, want the one published artifact", len(targets))
	}
	for _, field := range []string{"id", "kind", "bytes"} {
		if _, present := targets[0][field]; !present {
			t.Errorf("target row carries no %s field: %#v", field, targets[0])
		}
	}
	if targets[0]["id"] != identity || targets[0]["kind"] != "published" {
		t.Errorf("target = %#v, want the published artifact %s", targets[0], identity)
	}
	if got := preparedCount(t, row, "targets"); got != len(targets) {
		t.Errorf("orientation counts %d targets, want %d", got, len(targets))
	}
}

// TestEvidenceCleanupAppliedSchema is CE171 and the command half of CE81. Apply carries
// every registered applied field, and the removed artifact then refuses retrieval.
func TestEvidenceCleanupAppliedSchema(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	plan, _ := cleanupPlan(t)
	fingerprint, _ := plan["fingerprint"].(string)

	out, code := preflight.Command([]string{"evidence-clean", "--apply", fingerprint})
	if code != 0 {
		t.Fatalf("cleanup apply = (%d):\n%s", code, out)
	}
	rows := preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "applied")
	if len(rows) != 1 {
		t.Fatalf("apply carries %d rows:\n%s", len(rows), out)
	}
	for _, field := range []string{"fingerprint", "removed", "remaining", "complete", "next"} {
		if _, present := rows[0][field]; !present {
			t.Errorf("applied row carries no %s field: %#v", field, rows[0])
		}
	}
	if rows[0]["fingerprint"] != fingerprint || rows[0]["complete"] != true || rows[0]["next"] != "" {
		t.Errorf("applied = %#v, want a complete apply of %s with no successor", rows[0], fingerprint)
	}
	if len(preflighttest.PublishedPacks(t, root)) != 0 {
		t.Fatalf("apply left %v", preflighttest.PublishedPacks(t, root))
	}
	// CE81: the removed handle refuses retrieval rather than substituting checkout bytes.
	read, code := preflight.Command([]string{"evidence", identity})
	if code != 1 || !strings.Contains(read, "absent-artifact") {
		t.Fatalf("read of the removed artifact = (%d):\n%s", code, read)
	}
}

// TestEvidenceCleanupStalePlan is the command half of CE79 and CE122. Apply revalidates the
// fingerprint against the current store before any deletion, so a plan taken before the
// store changed refuses through the public command and deletes nothing. The plan's named
// mutation, skipping that revalidation, reds here at the command the plan names.
func TestEvidenceCleanupStalePlan(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	stale, _ := cleanupPlan(t)
	fingerprint, _ := stale["fingerprint"].(string)

	// A second artifact changes the target inventory the fingerprint committed to.
	preflighttest.MustWriteFile(t, "specs/"+slug+"/tickets/one.md", preflighttest.TicketDoc("One", "PF1", "PF2")+"second artifact\n")
	prepareEvidence(t, preflighttest.LegacyCommitted(t, root, slug, "second artifact"))
	before := preflighttest.PublishedPacks(t, root)
	if len(before) != 2 {
		t.Fatalf("the store holds %v, want two artifacts", before)
	}

	out, code := preflight.Command([]string{"evidence-clean", "--apply", fingerprint})
	if code != 1 || !strings.Contains(out, "stale-plan") {
		t.Fatalf("apply of a stale plan = (%d):\n%s", code, out)
	}
	if got := preflighttest.PublishedPacks(t, root); strings.Join(got, ",") != strings.Join(before, ",") {
		t.Fatalf("the refused apply changed the store to %v, want %v", got, before)
	}
	// The refusal names a fresh plan, and that fresh plan covers both artifacts.
	if !strings.Contains(out, evidencecmd.CleanCommand) {
		t.Errorf("the stale-plan refusal names no fresh plan: %s", out)
	}
	fresh, targets := cleanupPlan(t)
	if fresh["fingerprint"] == fingerprint || len(targets) != 2 {
		t.Errorf("fresh plan = %#v with %d targets, want a new fingerprint over both", fresh["fingerprint"], len(targets))
	}
}

// TestEvidenceCleanupGrammar is CE168. Cleanup accepts exactly its declared forms, and every
// refused form exits two before the store is touched.
func TestEvidenceCleanupGrammar(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	plan, _ := cleanupPlan(t)
	fingerprint, _ := plan["fingerprint"].(string)
	before := preflighttest.PublishedPacks(t, root)

	for _, test := range []struct {
		name string
		args []string
	}{
		{"apply with cursor", []string{"evidence-clean", "--apply", fingerprint, "--cursor", "v1." + strings.Repeat("0", 64) + ".c.0.0"}},
		{"apply without a fingerprint", []string{"evidence-clean", "--apply"}},
		{"apply with an empty fingerprint", []string{"evidence-clean", "--apply", ""}},
		{"duplicate cursor", []string{"evidence-clean", "--cursor", "v1." + strings.Repeat("0", 64) + ".c.0.0", "--cursor", "v1." + strings.Repeat("0", 64) + ".c.0.1"}},
		{"extra operand", []string{"evidence-clean", "example"}},
		{"evidence source selector", []string{"evidence-clean", "--source", "s2"}},
		{"evidence verify selector", []string{"evidence-clean", "--verify"}},
		{"build ticket flag", []string{"evidence-clean", "--ticket", "one.md"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, code := preflight.Command(test.args)
			if code != 2 || !strings.HasPrefix(out, "usage: ") && !strings.Contains(out, "usage") {
				t.Fatalf("%s = (%d):\n%s", test.name, code, out)
			}
			if got := preflighttest.PublishedPacks(t, root); strings.Join(got, ",") != strings.Join(before, ",") {
				t.Fatalf("%s changed the store to %v, want %v", test.name, got, before)
			}
		})
	}
	// The two declared forms stay accepted, so the refusals above are not a blanket refusal.
	if _, code := preflight.Command([]string{"evidence-clean"}); code != 0 {
		t.Errorf("the declared plan form exits %d", code)
	}
}

// TestEvidenceCleanupHelpInventory is CE172. The public help projects exactly the
// implemented cleanup grammar, and an omission oracle names each form independently.
func TestEvidenceCleanupHelpInventory(t *testing.T) {
	var clean []string
	for _, row := range evidencecmd.HelpRows() {
		if strings.HasPrefix(row.Suffix, " evidence-clean") {
			clean = append(clean, strings.TrimSpace(row.Suffix))
		}
	}
	want := []string{
		"evidence-clean [--cursor <cursor>]",
		"evidence-clean --apply <fingerprint>",
	}
	if strings.Join(clean, "\n") != strings.Join(want, "\n") {
		t.Errorf("cleanup help rows =\n%s\nwant\n%s", strings.Join(clean, "\n"), strings.Join(want, "\n"))
	}
	for _, row := range evidencecmd.HelpRows() {
		if strings.Contains(row.Suffix, "evidence-clean") && row.Description == "" {
			t.Errorf("cleanup help row %q carries no description", row.Suffix)
		}
	}
}
