package evidencecmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/capability"
	"github.com/gibbonmi/bench/internal/chargeevidence"
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

// TestEvidenceCleanupCursorStream is CE78 at the cursor seam. One cursor grammar carries two
// streams, so the cleanup page refuses an artifact cursor, the artifact read refuses a
// cleanup cursor, and each path still accepts its own. The rendered stream marker is the one
// token that separates them, so it is pinned here.
func TestEvidenceCleanupCursorStream(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	identity, _, _ := prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	plan, _ := cleanupPlan(t)
	fingerprint, _ := plan["fingerprint"].(string)

	// Each cursor names the stream it does not belong to, under the identity its target path
	// accepts, so the refusal below is the stream check and never an identity mismatch.
	foreignToClean := chargeevidence.Cursor{Identity: fingerprint}.String()
	foreignToRead := chargeevidence.Cursor{Identity: identity, Clean: true}.String()
	if got := strings.Split(foreignToRead, ".")[2]; got != "c" {
		t.Errorf("cleanup cursor stream marker = %q, want c", got)
	}
	if got := strings.Split(foreignToClean, ".")[2]; got != "m" {
		t.Errorf("artifact cursor stream marker = %q, want m", got)
	}
	for _, test := range []struct {
		name string
		args []string
	}{
		{"cleanup page from an artifact cursor", []string{"evidence-clean", "--cursor", foreignToClean}},
		{"artifact read from a cleanup cursor", []string{"evidence", identity, "--cursor", foreignToRead}},
	} {
		out, code := preflight.Command(test.args)
		if code != 1 || !strings.Contains(out, chargeevidence.RefuseCursor) {
			t.Errorf("%s = (%d):\n%s", test.name, code, out)
		}
	}
	// Each path accepts its own stream, so the refusals above are not a blanket refusal.
	own := chargeevidence.Cursor{Identity: fingerprint, Clean: true}.String()
	if out, code := preflight.Command([]string{"evidence-clean", "--cursor", own}); code != 0 {
		t.Errorf("cleanup page from its own cursor = (%d):\n%s", code, out)
	}
	if out, code := preflight.Command([]string{"evidence", identity, "--cursor", chargeevidence.Cursor{Identity: identity}.String()}); code != 0 {
		t.Errorf("artifact read from its own cursor = (%d):\n%s", code, out)
	}
}

// TestEvidenceCleanupStopped is CE82 and CE122 at the command surface. A store that admits no
// deletion stops the apply: it exits one, deletes nothing, and names a fresh plan. That named
// plan is the only recovery, and it applies completely once the store admits deletions again.
func TestEvidenceCleanupStopped(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	plan, _ := cleanupPlan(t)
	fingerprint, _ := plan["fingerprint"].(string)

	// The apply holds its locks through already opened files, so a store directory that
	// admits no unlink fails exactly the deletion the plan authorized.
	store := preflighttest.StoreDir(t, root)
	t.Cleanup(func() { _ = os.Chmod(store, 0o700) })
	if err := os.Chmod(store, 0o500); err != nil {
		capability.Capability(t, capability.Privilege, fmt.Sprintf("cannot strip directory permissions: %v", err))
	}
	// A privileged run writes a mode 0o500 directory anyway, so this case proves the
	// refusal it needs before it grades the apply.
	probe := filepath.Join(store, "admits-writes")
	if err := os.WriteFile(probe, nil, 0o600); err == nil {
		_ = os.Remove(probe)
		capability.Capability(t, capability.Privilege, "mode 0o500 store directory is still writable by this user")
	}
	out, code := preflight.Command([]string{"evidence-clean", "--apply", fingerprint})
	if err := os.Chmod(store, 0o700); err != nil {
		t.Fatal(err)
	}
	rows := preflighttest.TableRows(t, preflighttest.DecodeMap(t, out), "applied")
	if code != 1 || len(rows) != 1 {
		t.Fatalf("stopped apply = (%d):\n%s", code, out)
	}
	if rows[0]["complete"] != false || preparedCount(t, rows[0], "removed") != 0 || preparedCount(t, rows[0], "remaining") != 1 {
		t.Errorf("stopped apply = %#v, want nothing removed and one remaining", rows[0])
	}
	if len(preflighttest.PublishedPacks(t, root)) != 1 {
		t.Fatalf("the stopped apply changed the store to %v", preflighttest.PublishedPacks(t, root))
	}
	next, _ := rows[0]["next"].(string)
	if !strings.Contains(next, evidencecmd.CleanCommand) {
		t.Fatalf("the stopped apply names no recovery: %#v", rows[0])
	}
	// The recovery it names produces a plan that applies completely.
	fresh, _ := cleanupPlan(t)
	freshFingerprint, _ := fresh["fingerprint"].(string)
	if applied, code := preflight.Command([]string{"evidence-clean", "--apply", freshFingerprint}); code != 0 {
		t.Fatalf("apply of the fresh plan = (%d):\n%s", code, applied)
	}
	if got := preflighttest.PublishedPacks(t, root); len(got) != 0 {
		t.Fatalf("the fresh apply left %v", got)
	}
}

// TestEvidenceCleanupEmptyPlan is CE100 at the command surface. An empty store plans no
// target, so it advertises no apply: an apply of an empty plan would authorize nothing.
func TestEvidenceCleanupEmptyPlan(t *testing.T) {
	preflighttest.SeedConformant(t)
	row, targets := cleanupPlan(t)
	if len(targets) != 0 || preparedCount(t, row, "targets") != 0 || preparedCount(t, row, "bytes") != 0 {
		t.Fatalf("the unprepared store planned %d targets: %#v", len(targets), row)
	}
	if row["next"] != "" {
		t.Errorf("the empty plan advertises %#v, want no successor", row["next"])
	}
}

// TestEvidenceCleanupZeroByteTarget is CE85 and CE100 at the command surface. A zero-byte
// orphan temporary is still a target the apply deletes, so the plan that names it advertises
// its apply. The target count decides that successor; the planned byte total does not.
func TestEvidenceCleanupZeroByteTarget(t *testing.T) {
	root, slug := preflighttest.SeedConformant(t)
	prepareEvidence(t, preflighttest.ChargeArgs(t, root, slug))
	plan, _ := cleanupPlan(t)
	fingerprint, _ := plan["fingerprint"].(string)
	if out, code := preflight.Command([]string{"evidence-clean", "--apply", fingerprint}); code != 0 {
		t.Fatalf("apply = (%d):\n%s", code, out)
	}
	// The emptied store keeps the lock file every plan needs, so a killed writer's empty
	// temporary is the one target it holds.
	name := chargeevidence.TempPrefix + "orphan" + chargeevidence.TempSuffix
	if err := os.WriteFile(filepath.Join(preflighttest.StoreDir(t, root), name), nil, 0o600); err != nil {
		t.Fatal(err)
	}

	row, targets := cleanupPlan(t)
	if len(targets) != 1 || preparedCount(t, row, "targets") != 1 || preparedCount(t, row, "bytes") != 0 {
		t.Fatalf("the orphaned store planned %d target rows: %#v", len(targets), row)
	}
	orphaned, _ := row["fingerprint"].(string)
	if want := evidencecmd.CleanCommand + " --apply " + orphaned; row["next"] != want {
		t.Errorf("the zero-byte plan advertises %#v, want %q", row["next"], want)
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
		// A non-empty malformed fingerprint passes the flag parser, so only the command's own
		// shape check refuses it. It is a usage fault, not a plan the store could ever hold.
		{"apply with a malformed fingerprint", []string{"evidence-clean", "--apply", "sha256:" + strings.Repeat("g", 64)}},
		{"duplicate cursor", []string{"evidence-clean", "--cursor", "v1." + strings.Repeat("0", 64) + ".c.0.0", "--cursor", "v1." + strings.Repeat("0", 64) + ".c.0.1"}},
		{"extra operand", []string{"evidence-clean", "example"}},
		{"evidence source selector", []string{"evidence-clean", "--source", "s2"}},
		{"evidence verify selector", []string{"evidence-clean", "--verify"}},
		{"build ticket flag", []string{"evidence-clean", "--ticket", "one.md"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			out, code := preflight.Command(test.args)
			if code != 2 || !strings.HasPrefix(out, "usage: ") {
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
