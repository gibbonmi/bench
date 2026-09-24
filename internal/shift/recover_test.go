package shift

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
)

// deadOwner is a process id above every kernel's pid_max, so no process owns it.
const deadOwner = 1 << 30

// seedShiftEntry writes one open shift entry into the working repository's ledger.
func seedShiftEntry(t *testing.T, entry intent.Entry) (root string) {
	t.Helper()
	root, err := git.Root()
	if err != nil {
		t.Fatal(err)
	}
	if err := intent.Upsert(root, entry); err != nil {
		t.Fatal(err)
	}
	return root
}

// ledgerEntry returns the stored entry of key, live or not.
func ledgerEntry(t *testing.T, root, key string) intent.Entry {
	t.Helper()
	current, err := intent.Read(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range current.Entries {
		if entry.Key == key {
			return entry
		}
	}
	t.Fatalf("the ledger holds no entry %s", key)
	return intent.Entry{}
}

// recoverySpans returns the shift.recovery spans of root's record.
func recoverySpans(t *testing.T, root string) []otelrecord.Span {
	t.Helper()
	spans, _ := otelrecord.ReadSpans(os.Getenv("BENCH_HOME"), root)
	var found []otelrecord.Span
	for _, span := range spans {
		if span.Seam == recoverySeam {
			found = append(found, span)
		}
	}
	return found
}

func TestKeyOwnerReadsTheKeyNewEntryWrites(t *testing.T) {
	if pid, ok := intent.KeyOwner(intent.KindShift, intent.NewEntry(intent.KindShift).Key); !ok || pid != os.Getpid() {
		t.Fatalf("KeyOwner = %d, %v, want this process", pid, ok)
	}
}

// LE72 and LE78: a lease-less entry of a dead owner is abandoned, recorded, and reported.
func TestRecoverAbandonsADeadOwnersEntry(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	entry := intent.EntryOwnedBy(intent.KindShift, deadOwner)
	root := seedShiftEntry(t, entry)
	var out bytes.Buffer
	Recover(root, &out)
	if got := ledgerEntry(t, root, entry.Key).Outcome; got != otelrecord.WorkAbandoned {
		t.Fatalf("outcome = %q, want %q", got, otelrecord.WorkAbandoned)
	}
	spans := recoverySpans(t, root)
	if len(spans) != 1 {
		t.Fatalf("the record holds %d recovery spans, want 1", len(spans))
	}
	requireAttr(t, spans[0], otelrecord.AttrIntentKey, entry.Key)
	requireAttr(t, spans[0], otelrecord.AttrWorkState, otelrecord.WorkAbandoned)
	requireAttr(t, spans[0], otelrecord.AttrCleanup, otelrecord.CleanupNone)
	if want := "bench shift recovery: recovered 0, abandoned 1\n"; out.String() != want {
		t.Fatalf("recovery output = %q, want %q", out.String(), want)
	}
}

// LE73 and LE74: a live owner's entry, an unparsable key, and a closed entry stay
// unchanged.
func TestRecoverKeepsALiveOrUnknownEntry(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	live := intent.NewEntry(intent.KindShift)
	entries := []intent.Entry{live}
	// Each key fails one part of the parse: the owner, the owner's range, or the stamp.
	for _, key := range []string{"shift-owner-1", "shift-5368709120-1", fmt.Sprintf("shift-%d-stamp", deadOwner)} {
		entries = append(entries, intent.Entry{Key: key, Kind: intent.KindShift, CreatedAt: live.CreatedAt})
	}
	closed := intent.EntryOwnedBy(intent.KindShift, deadOwner)
	closed.Outcome = otelrecord.WorkCompleted
	root := ""
	for _, entry := range append(entries, closed) {
		root = seedShiftEntry(t, entry)
	}
	Recover(root, io.Discard)
	for _, entry := range append(entries, closed) {
		if got := ledgerEntry(t, root, entry.Key).Outcome; got != entry.Outcome {
			t.Fatalf("entry %s outcome = %q, want it unchanged", entry.Key, got)
		}
	}
	if spans := recoverySpans(t, root); len(spans) != 0 {
		t.Fatalf("the record holds %d recovery spans, want none", len(spans))
	}
}

// LE79: a pass that acted on nothing prints nothing.
func TestRecoverWithNothingToDoPrintsNothing(t *testing.T) {
	root := faultFixtureCore(t, greenGate, nil)
	var out bytes.Buffer
	Recover(root, &out)
	if out.Len() != 0 {
		t.Fatalf("recovery output = %q, want none", out.String())
	}
}

// LE77: a shift runs the pass before its acquire, so the pass acts even when the
// acquire fails.
func TestAShiftRecoversBeforeItsAcquire(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	withAgent(t, "true\n")
	entry := intent.EntryOwnedBy(intent.KindShift, deadOwner)
	root := seedShiftEntry(t, entry)
	// A regular file where the worktree pool directory goes makes the acquire fail.
	home := os.Getenv("BENCH_HOME")
	if err := os.MkdirAll(home, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "worktrees"), nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if code := Loop("recovering shift", io.Discard, io.Discard); code != exitCodes[OutcomeUsage] {
		t.Fatalf("Loop = %d, want usage from the failed acquire", code)
	}
	if got := ledgerEntry(t, root, entry.Key).Outcome; got != otelrecord.WorkAbandoned {
		t.Fatalf("outcome = %q, want %q", got, otelrecord.WorkAbandoned)
	}
}
