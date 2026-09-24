package shift

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/otelrecord"
	"github.com/gibbonmi/bench/internal/worktree"
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

// worktreeState is what a pass must leave unchanged: the lease bytes, the dirty file,
// and the lock state of the entry's worktree.
func worktreeState(t *testing.T, root string, entry intent.Entry) string {
	t.Helper()
	lease, _ := worktree.LeaseFile(entry.Worktree)
	// Only a regular lease is read, so a FIFO lease never blocks the comparison.
	leaseBytes := []byte("special")
	if info, err := os.Lstat(lease); err == nil && info.Mode().IsRegular() {
		leaseBytes, _ = os.ReadFile(lease)
	}
	work, _ := os.ReadFile(filepath.Join(entry.Worktree, "work.txt"))
	return fmt.Sprintf("lease %q, work %q, locked %v", leaseBytes, work, locked(t, root, entry.Worktree))
}

// locked reports whether git lists wt as locked.
func locked(t *testing.T, root, wt string) bool {
	t.Helper()
	listing := runGitOutput(t, root, "worktree", "list", "--porcelain")
	block, _, _ := strings.Cut(listing[strings.Index(listing, "worktree "+wt):], "\n\n")
	return strings.Contains(block, "\nlocked")
}

// rewriteLease replaces the lease owner of entry's worktree with pid.
func rewriteLease(t *testing.T, entry intent.Entry, pid int) {
	t.Helper()
	lease, _ := worktree.LeaseFile(entry.Worktree)
	_, stamp, _ := strings.Cut(entry.Lease, " ")
	if err := os.WriteFile(lease, []byte(fmt.Sprintf("%d %s\n", pid, stamp)), 0o600); err != nil {
		t.Fatal(err)
	}
}

// requireUntouched fails when the pass changed the entry or its worktree, or recorded.
func requireUntouched(t *testing.T, root string, entry intent.Entry, before string) {
	t.Helper()
	if got := ledgerEntry(t, root, entry.Key); got != entry {
		t.Fatalf("entry = %+v, want %+v", got, entry)
	}
	if after := worktreeState(t, root, entry); after != before {
		t.Fatalf("worktree %s, want %s", after, before)
	}
	if spans := recoverySpans(t, root); len(spans) != 0 {
		t.Fatalf("the record holds %d recovery spans, want none", len(spans))
	}
}

// LE64, LE65, LE91, LE66, LE67, and LE100: a killed dirty shift is recovered, and a
// second pass changes nothing.
func TestRecoverRecoversAKilledDirtyShift(t *testing.T) {
	root, entry := crashedShift(t, true, true)
	Recover(root, io.Discard)
	spans := recoverySpans(t, root)
	if len(spans) != 1 {
		t.Fatalf("the record holds %d recovery spans, want 1", len(spans))
	}
	requireAttr(t, spans[0], otelrecord.AttrWorkState, otelrecord.WorkRecovered)
	requireAttr(t, spans[0], otelrecord.AttrIntentKey, entry.Key)
	requireAttr(t, spans[0], otelrecord.AttrMemoryState, otelrecord.MemoryRetained)
	got := ledgerEntry(t, root, entry.Key)
	if got.Outcome != otelrecord.WorkRecovered || !strings.HasPrefix(got.Recovery, recoveryWorktreeKind+":") {
		t.Fatalf("entry = %+v, want recovered with a worktree pointer", got)
	}
	if _, live := shiftIntent(t, entry.Key); !live {
		t.Fatal("the recovered dirty entry is not live")
	}
	if files := memoryFiles(t); len(files) != 1 || string(files[0]) != "MEMMARK\n" {
		t.Fatalf("memory files = %q, want the crashed notes", files)
	}
	if !locked(t, root, entry.Worktree) || !strings.Contains(worktreeState(t, root, entry), `work "work\n"`) {
		t.Fatalf("worktree %s, want it locked with its dirty file", worktreeState(t, root, entry))
	}
	before := worktreeState(t, root, got)
	Recover(root, io.Discard)
	if after := worktreeState(t, root, got); after != before || ledgerEntry(t, root, entry.Key) != got || len(recoverySpans(t, root)) != 1 || len(memoryFiles(t)) != 1 {
		t.Fatalf("a second pass changed the finished recovery: %s", after)
	}
}

// LE68: a killed shift with only scratch files is released.
func TestRecoverReleasesAKilledCleanShift(t *testing.T) {
	root, entry := crashedShift(t, false, true)
	Recover(root, io.Discard)
	lease, _ := worktree.LeaseFile(entry.Worktree)
	for _, path := range []string{lease, filepath.Join(entry.Worktree, notesFile)} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("%s survived the release: %v", path, err)
		}
	}
	if spans := recoverySpans(t, root); len(spans) != 1 {
		t.Fatalf("the record holds %d recovery spans, want 1", len(spans))
	} else {
		requireAttr(t, spans[0], otelrecord.AttrCleanup, otelrecord.CleanupReleased)
	}
	if files := memoryFiles(t); len(files) != 1 || string(files[0]) != "MEMMARK\n" {
		t.Fatalf("memory files = %q, want the crashed notes", files)
	}
}

// LE69 and LE70: another identity with a dead owner abandons the entry and touches nothing.
func TestRecoverAbandonsAnotherDeadIdentity(t *testing.T) {
	root, entry := crashedShift(t, true, true)
	rewriteLease(t, entry, reapedPID(t))
	before := worktreeState(t, root, entry)
	Recover(root, io.Discard)
	if got := ledgerEntry(t, root, entry.Key).Outcome; got != otelrecord.WorkAbandoned {
		t.Fatalf("outcome = %q, want abandoned", got)
	}
	if spans := recoverySpans(t, root); len(spans) != 1 {
		t.Fatalf("the record holds %d recovery spans, want 1", len(spans))
	} else {
		requireAttr(t, spans[0], otelrecord.AttrWorkState, otelrecord.WorkAbandoned)
	}
	if after := worktreeState(t, root, entry); after != before {
		t.Fatalf("worktree %s, want %s", after, before)
	}
}

// LE89: an absent lease abandons the entry and leaves every worktree file.
func TestRecoverAbandonsAnAbsentLease(t *testing.T) {
	root, entry := crashedShift(t, true, true)
	lease, _ := worktree.LeaseFile(entry.Worktree)
	if err := os.Remove(lease); err != nil {
		t.Fatal(err)
	}
	before := worktreeState(t, root, entry)
	Recover(root, io.Discard)
	if got := ledgerEntry(t, root, entry.Key).Outcome; got != otelrecord.WorkAbandoned {
		t.Fatalf("outcome = %q, want abandoned", got)
	}
	if after := worktreeState(t, root, entry); after != before {
		t.Fatalf("worktree %s, want %s", after, before)
	}
}

// LE103, LE71, LE75, and LE90: a live owner, a live helper, and a malformed or special
// lease leave the entry and its worktree unchanged.
func TestRecoverKeepsAnUnprovenLease(t *testing.T) {
	for _, tc := range []struct {
		name string
		kill bool
		set  func(t *testing.T, entry intent.Entry, lease string)
	}{
		{"live other owner", true, func(t *testing.T, entry intent.Entry, _ string) { rewriteLease(t, entry, os.Getpid()) }},
		{"live helper", false, func(*testing.T, intent.Entry, string) {}},
		{"malformed", true, func(t *testing.T, _ intent.Entry, lease string) { _ = os.WriteFile(lease, []byte("garbage\n"), 0o600) }},
		{"no final newline", true, func(t *testing.T, entry intent.Entry, lease string) {
			_ = os.WriteFile(lease, []byte(entry.Lease), 0o600)
		}},
		{"fifo", true, func(t *testing.T, _ intent.Entry, lease string) {
			_ = os.Remove(lease)
			if err := syscall.Mkfifo(lease, 0o600); err != nil {
				t.Fatal(err)
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, entry := crashedShift(t, true, tc.kill)
			lease, _ := worktree.LeaseFile(entry.Worktree)
			tc.set(t, entry, lease)
			before := worktreeState(t, root, entry)
			done := make(chan struct{})
			go func() { Recover(root, io.Discard); close(done) }()
			select {
			case <-done:
			case <-time.After(helperWindow):
				t.Fatal("the pass never returned")
			}
			requireUntouched(t, root, entry, before)
			if files := memoryFiles(t); len(files) != 0 {
				t.Fatalf("an unproven lease kept %d memory files, want none", len(files))
			}
		})
	}
}

// LE98 and LE105: a conceded claim changes nothing, and the notes are already kept.
func TestRecoverConcedesAFaultedClaim(t *testing.T) {
	root, entry := crashedShift(t, true, true)
	armFault(t, func(step shiftStep) error {
		if step == stepRecoveryClaim {
			return fmt.Errorf("injected claim loss")
		}
		return nil
	})
	before := worktreeState(t, root, entry)
	Recover(root, io.Discard)
	requireUntouched(t, root, entry, before)
	if files := memoryFiles(t); len(files) != 1 || string(files[0]) != "MEMMARK\n" {
		t.Fatalf("memory files = %q, want the crashed notes", files)
	}
}

// recoveredBeforeTheAct runs a first pass in a helper that exits between its entry write
// and its act, and returns the entry that pass wrote.
func recoveredBeforeTheAct(t *testing.T, dirty bool) (string, intent.Entry) {
	t.Helper()
	root, entry := crashedShift(t, dirty, true)
	startHelper(t, root, "recover-act").wait(t)
	got := ledgerEntry(t, root, entry.Key)
	if got.Outcome != otelrecord.WorkRecovered || locked(t, root, got.Worktree) {
		t.Fatalf("entry = %+v, want recovered before the act", got)
	}
	return root, got
}

// LE99 and LE101: a second pass finishes an interrupted recovery with one memory file, and
// its span carries no memory keys. A tree already locked keeps its lock, and a clean tree
// is released.
func TestRecoverFinishesAnInterruptedRecovery(t *testing.T) {
	for _, tc := range []struct {
		name          string
		dirty, locked bool
	}{{"dirty", true, false}, {"already locked", true, true}, {"clean", false, false}} {
		t.Run(tc.name, func(t *testing.T) {
			root, entry := recoveredBeforeTheAct(t, tc.dirty)
			if tc.locked {
				runGitOutput(t, root, "worktree", "lock", entry.Worktree)
			}
			Recover(root, io.Discard)
			lease, _ := worktree.LeaseFile(entry.Worktree)
			if _, err := os.Stat(lease); !os.IsNotExist(err) || locked(t, root, entry.Worktree) != tc.dirty {
				t.Fatalf("worktree %s, want no lease and lock %v", worktreeState(t, root, entry), tc.dirty)
			}
			if files := memoryFiles(t); len(files) != 1 {
				t.Fatalf("the record holds %d memory files, want 1", len(files))
			}
			if spans := recoverySpans(t, root); len(spans) != 1 {
				t.Fatalf("the record holds %d recovery spans, want 1", len(spans))
			} else {
				requireNoAttr(t, spans[0], otelrecord.AttrMemoryState)
			}
		})
	}
}

// LE104: a recovered entry whose lease names another dead owner stays unchanged.
func TestRecoverKeepsARecoveryUnderAnotherLease(t *testing.T) {
	root, entry := recoveredBeforeTheAct(t, true)
	rewriteLease(t, entry, reapedPID(t))
	before := worktreeState(t, root, entry)
	Recover(root, io.Discard)
	requireUntouched(t, root, entry, before)
}
