package admissionpolicy

import (
	"reflect"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent/ledger"
)

const (
	ownerID   = "0123456789abcdef0123456789abcdef"
	firstID   = "11111111111111111111111111111111"
	secondID  = "22222222222222222222222222222222"
	oldDigest = "aa11111111111111111111111111111111111111111111111111111111111111"
	newDigest = "bb22222222222222222222222222222222222222222222222222222222222222"
	otherHex  = "cc33333333333333333333333333333333333333333333333333333333333333"
	startOID  = "1234567890123456789012345678901234567890"
)

var (
	early = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	late  = time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
)

func entry(key string, kind ledger.Kind, created time.Time, worktree, branch string) ledger.Entry {
	return ledger.Entry{Key: key, Kind: kind, CreatedAt: created, Worktree: worktree, Branch: branch}
}

func assignment(id, request string) ledger.Assignment {
	return ledger.Assignment{
		Schema:   ledger.AssignmentRecordSchema,
		ID:       id,
		OwnerID:  ownerID,
		Request:  request,
		Label:    "lane",
		Start:    startOID,
		Branch:   ledger.AssignmentBranchRef(ownerID, id),
		Worktree: "/pool/" + id,
		State:    ledger.StateActive,
	}
}

func receipt(target, state, phase string) ledger.CleanupReceipt {
	return ledger.CleanupReceipt{
		Schema:      ledger.CleanupReceiptSchema,
		Repo:        "/repo",
		Operation:   "release",
		Target:      target,
		Fingerprint: otherHex,
		State:       state,
		Phase:       phase,
	}
}

// TestAdmissionPolicyDecidesEveryMutator drives each of the seven mutators' rules
// with a literal ledger value and observes the returned ledger and the changed
// report. No case builds a repository and no case takes a lock, so a rule left in
// the effect adapter cannot satisfy this test. (Coverage row LS15.)
func TestAdmissionPolicyDecidesEveryMutator(t *testing.T) {
	present := LivenessFacts{WorktreeExists: map[string]bool{"/pool/live": true}}
	cases := []struct {
		name        string
		start       ledger.Ledger
		run         func(*testing.T, ledger.Ledger) (ledger.Ledger, bool, error)
		want        ledger.Ledger
		wantChanged bool
		wantErr     string
	}{
		{
			name:  "upsert appends an unseen key",
			start: ledger.Ledger{Entries: []ledger.Entry{entry("a", ledger.KindShift, early, "", "")}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return Upsert(l, entry("b", ledger.KindWorktree, late, "", ""))
			},
			want: ledger.Ledger{Entries: []ledger.Entry{
				entry("a", ledger.KindShift, early, "", ""),
				entry("b", ledger.KindWorktree, late, "", ""),
			}},
			wantChanged: true,
		},
		{
			name:  "upsert keeps the stored stamp and reports no change",
			start: ledger.Ledger{Entries: []ledger.Entry{entry("a", ledger.KindShift, early, "", "")}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return Upsert(l, entry("a", ledger.KindShift, late, "", ""))
			},
			want:        ledger.Ledger{Entries: []ledger.Entry{entry("a", ledger.KindShift, early, "", "")}},
			wantChanged: false,
		},
		{
			name:  "upsert enriches a stored key",
			start: ledger.Ledger{Entries: []ledger.Entry{entry("a", ledger.KindShift, early, "", "")}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return Upsert(l, entry("a", ledger.KindShift, late, "/pool/a", ""))
			},
			want:        ledger.Ledger{Entries: []ledger.Entry{entry("a", ledger.KindShift, early, "/pool/a", "")}},
			wantChanged: true,
		},
		{
			name: "compaction drops the entry no fact proves live",
			start: ledger.Ledger{Entries: []ledger.Entry{
				entry("live", ledger.KindWorktree, early, "/pool/live", ""),
				entry("gone", ledger.KindWorktree, late, "/pool/gone", ""),
			}},
			run:         func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) { return Compact(l, present) },
			want:        ledger.Ledger{Entries: []ledger.Entry{entry("live", ledger.KindWorktree, early, "/pool/live", "")}},
			wantChanged: true,
		},
		{
			name:  "cleanup receipt appends and reports a change",
			start: ledger.Ledger{},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return PutCleanupReceipt(l, receipt("/target", ledger.ReceiptInFlight, ledger.ReceiptPhasePlanned))
			},
			want:        ledger.Ledger{CleanupReceipts: []ledger.CleanupReceipt{receipt("/target", ledger.ReceiptInFlight, ledger.ReceiptPhasePlanned)}},
			wantChanged: true,
		},
		{
			name:  "cleanup receipt replaces its own identity",
			start: ledger.Ledger{CleanupReceipts: []ledger.CleanupReceipt{receipt("/target", ledger.ReceiptInFlight, ledger.ReceiptPhasePlanned)}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return PutCleanupReceipt(l, receipt("/target", ledger.ReceiptComplete, ledger.ReceiptPhaseTerminal))
			},
			want:        ledger.Ledger{CleanupReceipts: []ledger.CleanupReceipt{receipt("/target", ledger.ReceiptComplete, ledger.ReceiptPhaseTerminal)}},
			wantChanged: true,
		},
		{
			name:  "cleanup receipt refuses an invalid record",
			start: ledger.Ledger{},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return PutCleanupReceipt(l, receipt("relative", ledger.ReceiptInFlight, ledger.ReceiptPhasePlanned))
			},
			want:    ledger.Ledger{},
			wantErr: "cleanup receipt has invalid identity",
		},
		{
			name:  "assignment appends an unseen ID",
			start: ledger.Ledger{},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return PutAssignment(l, assignment(firstID, oldDigest))
			},
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantChanged: true,
		},
		{
			name:  "assignment identical to the stored one reports no change",
			start: ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return PutAssignment(l, assignment(firstID, oldDigest))
			},
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantChanged: false,
		},
		{
			name:  "assignment refuses a request another assignment owns",
			start: ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return PutAssignment(l, assignment(secondID, oldDigest))
			},
			want:    ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantErr: "assignment request already belongs to another assignment",
		},
		{
			name:  "reauthorization swaps the expected digest",
			start: ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return ReauthorizeAssignment(l, firstID, oldDigest, newDigest)
			},
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, newDigest)}},
			wantChanged: true,
		},
		{
			name:  "reauthorization refuses a digest that moved under it",
			start: ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, newDigest)}},
			run: func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				return ReauthorizeAssignment(l, firstID, oldDigest, newDigest)
			},
			want:    ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, newDigest)}},
			wantErr: "assignment request changed during reauthorization",
		},
		{
			name:  "purge drops the record keep rejects",
			start: ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest), assignment(secondID, newDigest)}},
			run: func(t *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				next, dropped, changed, err := PurgeAssignments(l, len(l.Assignments), func(a ledger.Assignment) bool { return a.ID == firstID })
				if dropped != 1 {
					t.Fatalf("purge reported %d dropped records, want 1", dropped)
				}
				return next, changed, err
			},
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantChanged: true,
		},
		{
			name:  "purge that keeps every record reports no change",
			start: ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			run: func(t *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) {
				next, dropped, changed, err := PurgeAssignments(l, len(l.Assignments), func(ledger.Assignment) bool { return true })
				if dropped != 0 {
					t.Fatalf("purge reported %d dropped records, want 0", dropped)
				}
				return next, changed, err
			},
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantChanged: false,
		},
		{
			name:        "delete removes the named assignment",
			start:       ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest), assignment(secondID, newDigest)}},
			run:         func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) { return DeleteAssignment(l, secondID) },
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantChanged: true,
		},
		{
			name:        "delete of an absent assignment reports no change",
			start:       ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			run:         func(_ *testing.T, l ledger.Ledger) (ledger.Ledger, bool, error) { return DeleteAssignment(l, secondID) },
			want:        ledger.Ledger{Assignments: []ledger.Assignment{assignment(firstID, oldDigest)}},
			wantChanged: false,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			next, changed, err := testCase.run(t, testCase.start)
			if testCase.wantErr != "" {
				if err == nil || err.Error() != testCase.wantErr {
					t.Fatalf("the rule answered error %v, want %q", err, testCase.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("the rule answered an unexpected error: %v", err)
			}
			if changed != testCase.wantChanged {
				t.Fatalf("the rule reported changed=%t, want %t", changed, testCase.wantChanged)
			}
			if !sameRecords(next, testCase.want) {
				t.Fatalf("the rule answered %+v, want %+v", next, testCase.want)
			}
		})
	}
}

// sameRecords compares the three record sets a rule can change, treating an empty
// set and a nil set as the same answer, because a rule never promises which one it
// returns.
func sameRecords(got, want ledger.Ledger) bool {
	return len(got.Entries) == len(want.Entries) &&
		len(got.Assignments) == len(want.Assignments) &&
		len(got.CleanupReceipts) == len(want.CleanupReceipts) &&
		(len(want.Entries) == 0 || reflect.DeepEqual(got.Entries, want.Entries)) &&
		(len(want.Assignments) == 0 || reflect.DeepEqual(got.Assignments, want.Assignments)) &&
		(len(want.CleanupReceipts) == 0 || reflect.DeepEqual(got.CleanupReceipts, want.CleanupReceipts))
}

// TestCompactionRuleReadsTypedLivenessFacts observes that the compaction rule reads
// its three typed facts and nothing else: a present worktree keeps an entry, an
// absent worktree drops one, a landed branch drops one, and the candidate probe
// alone decides an uncorrelated Claude entry. The test builds no repository and
// takes no lock. (Coverage row LS16.)
func TestCompactionRuleReadsTypedLivenessFacts(t *testing.T) {
	worktreeLive := entry("worktree-live", ledger.KindWorktree, early, "/pool/live", "")
	worktreeGone := entry("worktree-gone", ledger.KindWorktree, early, "/pool/gone", "")
	branchLanded := entry("branch-landed", ledger.KindShift, early, "", "topic-landed")
	branchOpen := entry("branch-open", ledger.KindShift, early, "", "topic-open")
	claudeLoose := entry("claude-loose", ledger.KindClaudeAgent, early, "", "")

	cases := []struct {
		name  string
		start []ledger.Entry
		facts LivenessFacts
		want  []string
	}{
		{
			name:  "a present worktree keeps its entry",
			start: []ledger.Entry{worktreeLive},
			facts: LivenessFacts{WorktreeExists: map[string]bool{"/pool/live": true}},
			want:  []string{"worktree-live"},
		},
		{
			name:  "an absent worktree drops its entry",
			start: []ledger.Entry{worktreeLive, worktreeGone},
			facts: LivenessFacts{WorktreeExists: map[string]bool{"/pool/live": true}},
			want:  []string{"worktree-live"},
		},
		{
			name:  "a landed branch drops its entry and an open branch keeps its own",
			start: []ledger.Entry{branchLanded, branchOpen},
			facts: LivenessFacts{Landed: map[string]bool{"topic-landed": true}},
			want:  []string{"branch-open"},
		},
		{
			name:  "a true candidate probe keeps an uncorrelated Claude entry",
			start: []ledger.Entry{claudeLoose},
			facts: LivenessFacts{Candidates: true},
			want:  []string{"claude-loose"},
		},
		{
			name:  "a false candidate probe drops an uncorrelated Claude entry",
			start: []ledger.Entry{claudeLoose},
			facts: LivenessFacts{},
			want:  nil,
		},
		{
			name:  "the candidate probe decides only the uncorrelated Claude entry",
			start: []ledger.Entry{claudeLoose, worktreeGone, branchLanded, branchOpen},
			facts: LivenessFacts{Candidates: true, Landed: map[string]bool{"topic-landed": true}},
			want:  []string{"branch-open", "claude-loose"},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			start := ledger.Ledger{Entries: append([]ledger.Entry(nil), testCase.start...)}
			if got := keysOf(Live(start, testCase.facts)); !equalKeys(got, testCase.want) {
				t.Fatalf("the liveness rule kept %v, want %v", got, testCase.want)
			}
			next, changed, err := Compact(start, testCase.facts)
			if err != nil {
				t.Fatalf("the compaction rule answered an unexpected error: %v", err)
			}
			if want := len(testCase.start) != len(testCase.want); changed != want {
				t.Fatalf("the compaction rule reported changed=%t, want %t", changed, want)
			}
			// The compaction keeps the stored order of the survivors, while the liveness
			// rule answers oldest first, so the compaction's expectation is the stored
			// order filtered by the same survivor set.
			wantCompact := storedOrder(testCase.start, testCase.want)
			if got := keysOf(next.Entries); !equalKeys(got, wantCompact) {
				t.Fatalf("the compaction rule kept %v, want %v", got, wantCompact)
			}
		})
	}
}

func storedOrder(stored []ledger.Entry, survivors []string) []string {
	keep := map[string]bool{}
	for _, key := range survivors {
		keep[key] = true
	}
	order := make([]string, 0, len(survivors))
	for _, e := range stored {
		if keep[e.Key] {
			order = append(order, e.Key)
		}
	}
	return order
}

func keysOf(entries []ledger.Entry) []string {
	keys := make([]string, 0, len(entries))
	for _, e := range entries {
		keys = append(keys, e.Key)
	}
	return keys
}

func equalKeys(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}
