package intent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLedgerCommonDirectoryAndSchemaUpsert(t *testing.T) {
	root := newRepo(t)
	linked := filepath.Join(t.TempDir(), "linked checkout")
	runGit(t, root, "worktree", "add", "-q", "--detach", linked, "HEAD")
	created := time.Date(2026, 7, 11, 12, 0, 0, 0, time.UTC)
	entries := []Entry{
		{Key: "shift-1", Kind: KindShift, CreatedAt: created, Worktree: linked, Branch: "bench/shift-1"},
		{Key: "worktree-1", Kind: KindWorktree, CreatedAt: created},
		{Key: "agent-1", Kind: KindClaudeAgent, CreatedAt: created},
	}
	for i, e := range entries {
		writerRoot := root
		if i == 1 {
			writerRoot = linked
		}
		if err := Upsert(writerRoot, e); err != nil {
			t.Fatalf("Upsert(%s): %v", e.Kind, err)
		}
	}
	path, err := Address(linked)
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(root, ".git", Filename)
	if path != wantPath {
		t.Fatalf("Address = %q, want common-dir %q", path, wantPath)
	}
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	ledger, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if ledger.Schema != Schema || len(ledger.Entries) != 3 {
		t.Fatalf("ledger = %#v", ledger)
	}
	if err := Upsert(root, entries[0]); err != nil {
		t.Fatal(err)
	}
	second, _ := os.ReadFile(path)
	if !bytes.Equal(first, second) {
		t.Fatalf("identical upsert changed bytes\nfirst=%s\nsecond=%s", first, second)
	}
}

func TestReadEvidenceStates(t *testing.T) {
	root := newRepo(t)
	path, _ := Address(root)
	if got, err := Read(root); err != nil || len(got.Entries) != 0 {
		t.Fatalf("absent Read = %#v, %v", got, err)
	}
	cases := []struct{ name, body string }{
		{"empty", ""},
		{"malformed", "{\n"},
		{"missing final newline", `{"schema":1,"entries":[]}`},
		{"duplicate field", "{\"schema\":1,\"schema\":1,\"entries\":[]}\n"},
		{"nested duplicate field", "{\"schema\":1,\"entries\":[{\"key\":\"legacy\",\"key\":\"legacy\",\"kind\":\"shift\",\"created_at\":\"2026-07-11T00:00:00Z\"}]}\n"},
		{"unknown field", "{\"schema\":1,\"entries\":[],\"unknown\":true}\n"},
		{"nested unknown field", "{\"schema\":1,\"entries\":[{\"key\":\"legacy\",\"kind\":\"shift\",\"created_at\":\"2026-07-11T00:00:00Z\",\"unknown\":true}]}\n"},
		{"trailing value", "{\"schema\":1,\"entries\":[]} {}\n"},
		{"trailing bytes", "{\"schema\":1,\"entries\":[]} nope\n"},
		{"wrong field type", "{\"schema\":\"1\",\"entries\":[]}\n"},
		{"unknown schema", "{\"schema\":99,\"entries\":[]}\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(path, []byte(tc.body), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := Read(root); err == nil {
				t.Fatalf("Read accepted %s ledger", tc.name)
			}
		})
	}
	legacy := " \n { \"schema\" : 1, \"entries\" : [ { \"key\" : \"legacy-1\", \"kind\" : \"shift\", \"created_at\" : \"2026-07-11T00:00:00Z\" } ] } \t\n"
	if err := os.WriteFile(path, []byte(legacy), 0o600); err != nil {
		t.Fatal(err)
	}
	ledger, err := Read(root)
	if err != nil || ledger.Schema != LegacySchema || len(ledger.Entries) != 1 || ledger.Entries[0].Key != "legacy-1" {
		t.Fatalf("canonical legacy ledger = %#v, %v", ledger, err)
	}
	if err := os.Chmod(path, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(path, 0o600)
	if _, err := Read(root); err == nil {
		t.Fatal("Read accepted unreadable ledger")
	}
}

func TestConcurrentWritersKeepEveryEntryAndStaleLockReclaims(t *testing.T) {
	root := newRepo(t)
	const n = 12
	var wg sync.WaitGroup
	errCh := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errCh <- Upsert(root, Entry{Key: string(rune('a' + i)), Kind: KindShift, CreatedAt: time.Unix(int64(i+1), 0).UTC()})
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
	got, err := Read(root)
	if err != nil || len(got.Entries) != n {
		t.Fatalf("concurrent Read = %d entries, %v", len(got.Entries), err)
	}
	path, _ := Address(root)
	if err := os.WriteFile(path+".lock", []byte("999999 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Upsert(root, Entry{Key: "reclaimed", Kind: KindWorktree, CreatedAt: time.Now().UTC()}); err != nil {
		t.Fatalf("stale lock not reclaimed: %v", err)
	}
}

func TestSnapshotProofLifecycle(t *testing.T) {
	root := newRepo(t)
	created := time.Unix(1, 0).UTC()
	missing := filepath.Join(root, "gone")
	for _, e := range []Entry{
		{Key: "no-proof", Kind: KindShift, CreatedAt: created},
		{Key: "missing", Kind: KindWorktree, CreatedAt: created, Worktree: missing},
		{Key: "landed", Kind: KindShift, CreatedAt: created, Branch: "landed"},
	} {
		if err := Upsert(root, e); err != nil {
			t.Fatal(err)
		}
	}
	runGit(t, root, "branch", "landed")
	live, err := Snapshot(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(live) != 1 || live[0].Key != "no-proof" {
		t.Fatalf("live = %#v", live)
	}
	if err := Compact(root); err != nil {
		t.Fatal(err)
	}
	ledger, _ := Read(root)
	if len(ledger.Entries) != 1 {
		t.Fatalf("compacted = %#v", ledger.Entries)
	}
}

func TestUncorrelatedEntriesUseCandidateSet(t *testing.T) {
	root := newRepo(t)
	for _, key := range []string{"agent-a", "agent-b"} {
		if err := Upsert(root, Entry{Key: key, Kind: KindClaudeAgent, CreatedAt: time.Unix(1, 0).UTC()}); err != nil {
			t.Fatal(err)
		}
	}
	if got, _ := Snapshot(root); len(got) != 0 {
		t.Fatalf("zero candidates kept %#v", got)
	}
	if err := Upsert(root, Entry{Key: "agent-c", Kind: KindClaudeAgent, CreatedAt: time.Unix(1, 0).UTC()}); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "branch", "worktree-agent-candidate")
	if got, _ := Snapshot(root); len(got) != 3 {
		t.Fatalf("one candidate live = %#v", got)
	}
}

func TestCleanupReceiptWindowKeepsExactlyLast256Completions(t *testing.T) {
	root := newRepo(t)
	repo := filepath.Join(root, ".git")
	before, err := LifecycleEvidence(root)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < MaxCleanupReceipts+2; i++ {
		receipt := CleanupReceipt{
			Schema: CleanupReceiptSchema, Repo: repo, Operation: "worktree-clean",
			Target: filepath.Join(root, fmt.Sprintf("target-%03d", i)), Fingerprint: fmt.Sprintf("%064x", i+1),
			State: ReceiptComplete, Phase: ReceiptPhaseTerminal, Action: "removed",
			Tracked: "clean", Ignored: "count=0 bytes=0 shown=0 truncated=false", Recovery: "none",
		}
		if err := PutCleanupReceipt(root, receipt); err != nil {
			t.Fatal(err)
		}
	}
	ledger, err := Read(root)
	if err != nil || len(ledger.CleanupReceipts) != MaxCleanupReceipts {
		t.Fatalf("completion window = %d, %v", len(ledger.CleanupReceipts), err)
	}
	if got := filepath.Base(ledger.CleanupReceipts[0].Target); got != "target-002" {
		t.Fatalf("first retained completion = %q", got)
	}
	if got := filepath.Base(ledger.CleanupReceipts[MaxCleanupReceipts-1].Target); got != "target-257" {
		t.Fatalf("last retained completion = %q", got)
	}
	after, err := LifecycleEvidence(root)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatalf("receipts changed assignment evidence: %q -> %q, %v", before, after, err)
	}
}

func TestUnstampedAssignmentRoundTripsThroughLedger(t *testing.T) {
	root := newRepo(t)
	want := activeAssignment()
	if err := PutAssignment(root, want); err != nil {
		t.Fatal(err)
	}
	got, err := Assignments(root)
	if err != nil || len(got) != 1 {
		t.Fatalf("Assignments = %#v, %v", got, err)
	}
	if !reflect.DeepEqual(got[0], want) {
		t.Fatalf("round-tripped assignment = %#v, want %#v", got[0], want)
	}
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(body, []byte("created_at")) {
		t.Fatalf("unstamped assignment serialized a created_at key: %s", body)
	}
}

// An identical re-write must preserve the ledger bytes. PutAssignment decides that
// with reflect.DeepEqual, so a stamp held behind a pointer has to compare by its
// value, not by its address.
func TestIdenticalStampedAssignmentWritePreservesBytes(t *testing.T) {
	root := newRepo(t)
	first, second := "2026-07-27T00:00:00Z", "2026-07-27T00:00:00Z"
	assignment := activeAssignment()
	assignment.CreatedAt = &first
	if err := PutAssignment(root, assignment); err != nil {
		t.Fatal(err)
	}
	path, err := Address(root)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	assignment.CreatedAt = &second
	if err := PutAssignment(root, assignment); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatalf("identical stamped upsert changed bytes\nbefore=%s\nafter=%s", before, after)
	}
}

func TestAssignmentCreatedAtRejectsMalformed(t *testing.T) {
	cases := []struct{ name, value string }{
		{"empty string", ""},
		{"non-timestamp text", "yesterday"},
		{"date without time", "2026-07-27"},
		{"control byte", "2026-07-27T00:00:00Z\x1b"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assignment := activeAssignment()
			assignment.CreatedAt = &tc.value
			if err := ValidateAssignment(assignment); err == nil {
				t.Fatalf("ValidateAssignment accepted created_at %q", tc.value)
			}
		})
	}
}

func TestAssignmentCreatedAtAcceptsFuture(t *testing.T) {
	future := "2126-07-27T00:00:00Z"
	assignment := activeAssignment()
	assignment.CreatedAt = &future
	if err := ValidateAssignment(assignment); err != nil {
		t.Fatalf("ValidateAssignment rejected a future created_at: %v", err)
	}
}

func TestCompareAndSwapRequestDigestRefusesConcurrentMovementAndPreservesOtherFields(t *testing.T) {
	assignment := activeAssignment()
	original := assignment
	assignment.Request = strings.Repeat("e", 64)
	if err := compareAndSwapRequestDigest(&assignment, original.Request, strings.Repeat("f", 64)); err == nil {
		t.Fatal("compareAndSwapRequestDigest accepted a moved request digest")
	}
	if assignment.Request != strings.Repeat("e", 64) {
		t.Fatalf("refused CAS changed request to %q", assignment.Request)
	}

	assignment = original
	replacement := strings.Repeat("f", 64)
	if err := compareAndSwapRequestDigest(&assignment, original.Request, replacement); err != nil {
		t.Fatal(err)
	}
	if assignment.Request != replacement {
		t.Fatalf("CAS request = %q, want %q", assignment.Request, replacement)
	}
	assignment.Request = original.Request
	if !reflect.DeepEqual(assignment, original) {
		t.Fatalf("CAS changed fields beyond request: got %#v, want %#v", assignment, original)
	}
}

func TestReauthorizeAssignmentRefusesRequestDigestCollision(t *testing.T) {
	root := newRepo(t)
	first := activeAssignment()
	first.Request = RequestDigest("first-request")
	second := activeAssignment()
	second.ID, second.OwnerID = strings.Repeat("c", 32), strings.Repeat("d", 32)
	second.Request = RequestDigest("replacement-request")
	second.Branch = AssignmentBranchRef(second.OwnerID, second.ID)
	second.Worktree = "/pool/other"
	if err := PutAssignment(root, first); err != nil {
		t.Fatal(err)
	}
	if err := PutAssignment(root, second); err != nil {
		t.Fatal(err)
	}
	before, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ReauthorizeAssignment(root, first.ID, "replacement-request", func(Assignment) error { return nil }, noReauthorizeTransition, nil); err == nil {
		t.Fatal("ReauthorizeAssignment accepted another assignment's request digest")
	}
	after, err := Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(after, before) {
		t.Fatalf("collision changed ledger: before=%#v after=%#v", before, after)
	}
}

func noReauthorizeTransition(Assignment, Assignment) (func(), error) { return func() {}, nil }

func activeAssignment() Assignment {
	owner, id := strings.Repeat("a", 32), strings.Repeat("b", 32)
	return Assignment{
		Schema:   AssignmentRecordSchema,
		ID:       id,
		OwnerID:  owner,
		Request:  strings.Repeat("c", 64),
		Label:    "delegate",
		Start:    strings.Repeat("d", 40),
		Branch:   AssignmentBranchRef(owner, id),
		Worktree: "/pool/delegate",
		State:    StateActive,
	}
}

func newRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init", "-q", "-b", "main")
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit(t, root, "add", ".")
	runGit(t, root, "-c", "user.name=bench", "-c", "user.email=bench@local", "commit", "-qm", "init")
	return root
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
