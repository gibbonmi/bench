package shift

import (
	"bytes"
	"io"
	"os"
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

// LE73 and LE74: a live owner's entry and an unparsable key stay unchanged.
func TestRecoverKeepsALiveOrUnknownEntry(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	live := intent.NewEntry(intent.KindShift)
	// The stamp parses and the owner does not, so only the owner parse refuses this key.
	unknown := intent.Entry{Key: "shift-owner-1", Kind: intent.KindShift, CreatedAt: live.CreatedAt}
	root := seedShiftEntry(t, live)
	seedShiftEntry(t, unknown)
	Recover(root, io.Discard)
	for _, key := range []string{live.Key, unknown.Key} {
		if got := ledgerEntry(t, root, key).Outcome; got != "" {
			t.Fatalf("entry %s outcome = %q, want it unchanged", key, got)
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

// LE77: a shift runs the pass before its acquire.
func TestAShiftRecoversBeforeItsAcquire(t *testing.T) {
	faultFixtureCore(t, greenGate, nil)
	withAgent(t, "true\n")
	entry := intent.EntryOwnedBy(intent.KindShift, deadOwner)
	root := seedShiftEntry(t, entry)
	var stdout bytes.Buffer
	if code := Loop("recovering shift", &stdout, io.Discard); code != exitCodes[OutcomeNoOp] {
		t.Fatalf("Loop = %d, want no-op", code)
	}
	if got := ledgerEntry(t, root, entry.Key).Outcome; got != otelrecord.WorkAbandoned {
		t.Fatalf("outcome = %q, want %q", got, otelrecord.WorkAbandoned)
	}
}
