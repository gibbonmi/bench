package lifecyclepolicy

import (
	"errors"
	"testing"
	"time"

	"github.com/gibbonmi/bench/internal/intent"
)

// The tables in this file are the LC1 typed lifecycle matrices: each partition
// varies exactly one named lifecycle fact — ownership, lease, eligibility,
// age, ignored output, preservation, action — over invented typed facts and
// pins the verdict the parent's real-Git characterization matrices observe. A
// missing ownership or preservation branch turns its named partition red.

func ownedAssignment(id string) *intent.Assignment {
	return &intent.Assignment{ID: id, Branch: "refs/heads/bench/a", State: intent.StateActive, Recovery: []intent.Recovery{}}
}

// ownedFacts is a clean, owned, correctly locked, trivially landed fixture the
// partitions below perturb one fact at a time.
func ownedFacts() ExplicitFacts {
	return ExplicitFacts{
		RegistrationBranchRef:  "refs/heads/bench/a",
		RegistrationLockReason: "bench lock",
		AssignmentLockReason:   "bench lock",
		MarkerPresent:          true,
		MatchedAssignment:      ownedAssignment("a1"),
		InitialTracked:         "clean",
		DefaultKnown:           true,
		LandedOK:               true,
		HeadRef:                "refs/heads/bench/a",
		Head:                   "oid",
	}
}

func TestExplicitDecisionTable(t *testing.T) {
	cases := []struct {
		name       string
		mutate     func(*ExplicitFacts)
		action     Action
		reasonCode Reason
		reason     string
	}{
		{"ownership/clean-owned-removes", func(f *ExplicitFacts) {}, ActionRemove, "", ""},
		{"ownership/marker-malformed", func(f *ExplicitFacts) { f.MarkerErr = errors.New("owner marker is malformed") }, ActionRetain, ReasonMalformed, "owner marker is malformed"},
		{"ownership/ledger-unreadable", func(f *ExplicitFacts) { f.AssignmentLedgerErr = errors.New("boom") }, ActionRetain, ReasonMalformed, "assignment ledger is unreadable"},
		{"ownership/ambiguous-assignments", func(f *ExplicitFacts) { f.AssignmentAmbiguous = true }, ActionRetain, ReasonMalformed, "registration has ambiguous assignments"},
		{"ownership/no-matching-assignment", func(f *ExplicitFacts) { f.MatchedAssignment = nil }, ActionRetain, ReasonMalformed, "owner marker has no matching assignment"},
		{"ownership/branch-mismatch", func(f *ExplicitFacts) { f.RegistrationBranchRef = "refs/heads/other" }, ActionRetain, ReasonUncertain, "assignment does not match current branch"},
		{"ownership/lock-mismatch", func(f *ExplicitFacts) { f.RegistrationLockReason = "foreign reason" }, ActionRetain, ReasonUnexpectedLock, "assignment does not match current Bench lock"},
		{"ownership/foreign-lock", func(f *ExplicitFacts) {
			f.MarkerPresent, f.MatchedAssignment, f.RegistrationLocked = false, nil, true
		}, ActionRetain, ReasonUnexpectedLock, "foreign or unexpected lock is retained"},
		{"lease/live-owned", func(f *ExplicitFacts) { f.LeasePresent, f.LeaseState = true, LeaseLive }, ActionRetain, ReasonLiveLease, "assignment has a live lease"},
		{"lease/live-unowned", func(f *ExplicitFacts) {
			f.MarkerPresent, f.MatchedAssignment = false, nil
			f.LeasePresent, f.LeaseState = true, LeaseLive
		}, ActionRetain, ReasonLiveLease, "unowned assignment has an ambiguous lease"},
		{"lease/dead-unowned", func(f *ExplicitFacts) {
			f.MarkerPresent, f.MatchedAssignment = false, nil
			f.LeasePresent, f.LeaseState = true, LeaseDead
		}, ActionRetain, ReasonLiveLease, "unowned assignment has an ambiguous lease"},
		{"lease/dead-owned-removes", func(f *ExplicitFacts) { f.LeasePresent, f.LeaseState = true, LeaseDead }, ActionRemove, "", ""},
		{"lease/unknown", func(f *ExplicitFacts) { f.LeasePresent, f.LeaseState = true, LeaseUnknown }, ActionRetain, ReasonUncertain, UnknownLeaseReason},
		{"lease/stat-error-owned", func(f *ExplicitFacts) { f.LeaseStatErr = errors.New("stat") }, ActionRetain, ReasonUncertain, UnknownLeaseReason},
		{"eligibility/nested-error", func(f *ExplicitFacts) { f.NestedErr = errors.New("scan") }, ActionRetain, ReasonUncertain, "nested repository state is unknown"},
		{"eligibility/nested-dirty", func(f *ExplicitFacts) { f.NestedState = NestedDirty }, ActionRetain, ReasonUncertain, "nested repository or submodule is dirty"},
		{"eligibility/embedded", func(f *ExplicitFacts) { f.NestedState = NestedEmbeddedClean }, ActionRetain, ReasonUncertain, "embedded repository is retained"},
		{"ignored/declaration-malformed", func(f *ExplicitFacts) { f.BuildOutputErr = errors.New("bad json") }, ActionRetain, ReasonMalformed, "build-output declaration is malformed"},
		{"ignored/inventory-uncertain", func(f *ExplicitFacts) { f.IgnoredErr = errors.New("walk") }, ActionRetain, ReasonUncertain, "ignored inventory is uncertain"},
		{"ignored/over-limit", func(f *ExplicitFacts) { f.IgnoredOverLimit = true }, ActionRetain, ReasonIgnored, "ignored inventory exceeds the destructive limit"},
		{"ignored/undeclared-refuses", func(f *ExplicitFacts) { f.IgnoredCount = 1 }, ActionRetain, ReasonIgnored, "ignored residuals require --discard-ignored"},
		{"action/declared-discard-remove", func(f *ExplicitFacts) { f.IgnoredCount, f.DeclaredIgnored = 1, true }, ActionDiscardRemove, "", ""},
		{"action/flagged-discard-remove", func(f *ExplicitFacts) { f.IgnoredCount, f.DiscardIgnored = 1, true }, ActionDiscardRemove, "", ""},
		{"preservation/dirty-recover-remove", func(f *ExplicitFacts) { f.InitialTracked = "dirty" }, ActionRecoverRemove, "", ""},
		{"preservation/detached-registration", func(f *ExplicitFacts) { f.RegistrationDetached = true }, ActionRecoverRemove, "", ""},
		{"preservation/recorded-recovery", func(f *ExplicitFacts) {
			f.MatchedAssignment.Recovery = []intent.Recovery{{Ref: "refs/bench/r1"}}
		}, ActionRecoverRemove, "", ""},
		{"eligibility/unsafe-target-override", func(f *ExplicitFacts) { f.UnsafeTarget = true }, ActionRetain, ReasonUncertain, "target contains unsafe control bytes"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			facts := ownedFacts()
			c.mutate(&facts)
			v := DecideExplicit(facts)
			if v.Action != c.action || v.ReasonCode != c.reasonCode || v.Reason != c.reason {
				t.Fatalf("verdict = %s/%s/%q, want %s/%s/%q", v.Action, v.ReasonCode, v.Reason, c.action, c.reasonCode, c.reason)
			}
		})
	}
}

// TestExplicitVerdictEvidence pins the non-action evidence the verdict carries:
// landed branch-deletion authority and the recovery lookup kinds.
func TestExplicitVerdictEvidence(t *testing.T) {
	t.Run("action/landed-branch-deletion", func(t *testing.T) {
		v := DecideExplicit(ownedFacts())
		if !v.DeleteBranch || v.BranchRef != "refs/heads/bench/a" || v.BranchOID != "oid" {
			t.Fatalf("landed verdict = %+v, want branch-deletion authority for the proven head", v)
		}
		if v.Landed.Kind != LandednessProven || !v.Landed.ProvenLanded() {
			t.Fatalf("landed typed = %+v, want proven landed", v.Landed)
		}
	})
	t.Run("action/unlanded-keeps-branch", func(t *testing.T) {
		facts := ownedFacts()
		facts.LandedOK = false
		if v := DecideExplicit(facts); v.DeleteBranch || v.Landed.ProvenLanded() {
			t.Fatalf("unlanded verdict = %+v, want no branch-deletion authority", v)
		}
	})
	t.Run("preservation/owned-lookup", func(t *testing.T) {
		facts := ownedFacts()
		facts.InitialTracked = "dirty"
		if v := DecideExplicit(facts); v.RecoveryLookup != RecoveryLookupOwned || v.Recovery != "none" {
			t.Fatalf("owned dirty verdict = %+v, want an owned recovery lookup", v)
		}
	})
	t.Run("preservation/recorded-ref-needs-no-lookup", func(t *testing.T) {
		facts := ownedFacts()
		facts.InitialTracked = "dirty"
		facts.MatchedAssignment.Recovery = []intent.Recovery{{Ref: "refs/bench/r1"}}
		if v := DecideExplicit(facts); v.RecoveryLookup != RecoveryNoLookup || v.Recovery != "refs/bench/r1" {
			t.Fatalf("recorded-recovery verdict = %+v, want the recorded ref without a lookup", v)
		}
	})
	t.Run("preservation/foreign-lookup", func(t *testing.T) {
		facts := ownedFacts()
		facts.MarkerPresent, facts.MatchedAssignment = false, nil
		facts.InitialTracked = "dirty"
		if v := DecideExplicit(facts); v.RecoveryLookup != RecoveryLookupForeign {
			t.Fatalf("foreign dirty verdict = %+v, want a foreign recovery lookup", v)
		}
	})
}

func pendingOutcome() ExplicitOutcome {
	return ExplicitOutcome{
		Action:          ActionRemove,
		HasAssignment:   true,
		Owned:           true,
		AssignmentID:    "a1",
		AssignmentState: intent.StateCleanupPending,
		Tracked:         "clean",
		Landed:          Landedness{Kind: LandednessProven, Landed: true},
	}
}

func TestAutomaticDecisionTable(t *testing.T) {
	cases := []struct {
		name       string
		facts      func() AutomaticFacts
		action     Action
		reasonCode Reason
		reason     string
		assignment string
	}{
		{"eligibility/explicit-error-uncertain", func() AutomaticFacts {
			return AutomaticFacts{ExplicitErr: errors.New("git broke")}
		}, ActionRetain, ReasonUncertain, "git broke", ""},
		{"eligibility/explicit-error-malformed", func() AutomaticFacts {
			return AutomaticFacts{ExplicitErr: errors.New("intent ledger is unreadable")}
		}, ActionRetain, ReasonMalformed, "intent ledger is unreadable", ""},
		{"eligibility/missing-branch-salvage", func() AutomaticFacts {
			return AutomaticFacts{ExplicitErr: errors.New("git broke"), MissingBranchAssignmentID: "m1"}
		}, ActionRetain, ReasonActive, "assignment landedness is unknown", "m1"},
		{"lease/missing-branch-live-lease", func() AutomaticFacts {
			return AutomaticFacts{ExplicitErr: errors.New("git broke"), MissingBranchAssignmentID: "m1", MissingBranchLiveLease: true}
		}, ActionRetain, ReasonLiveLease, "assignment has a live lease", "m1"},
		{"lease/live-overrides-everything", func() AutomaticFacts {
			f := AutomaticFacts{Explicit: pendingOutcome(), RecoveryMatches: true}
			f.LiveLease = true
			return f
		}, ActionRetain, ReasonLiveLease, "assignment has a live lease", ""},
		{"eligibility/retain-passthrough", func() AutomaticFacts {
			o := pendingOutcome()
			o.Action, o.ReasonCode, o.Reason = ActionRetain, ReasonIgnored, "ignored residuals require --discard-ignored"
			return AutomaticFacts{Explicit: o}
		}, ActionRetain, ReasonIgnored, "ignored residuals require --discard-ignored", ""},
		{"eligibility/retain-landed-swap", func() AutomaticFacts {
			o := pendingOutcome()
			o.Action, o.ReasonCode, o.Reason = ActionRetain, ReasonIgnored, "ignored residuals require --discard-ignored"
			return AutomaticFacts{Explicit: o, Landed: true}
		}, ActionRetain, ReasonLanded, "assignment branch has landed", ""},
		{"ownership/unowned-refused", func() AutomaticFacts {
			o := pendingOutcome()
			o.HasAssignment, o.Owned = false, false
			return AutomaticFacts{Explicit: o}
		}, ActionRetain, ReasonForeign, "registration is not a verified owned assignment", ""},
		{"age/young-active", func() AutomaticFacts {
			o := pendingOutcome()
			o.AssignmentState = intent.StateActive
			return AutomaticFacts{Explicit: o}
		}, ActionRetain, ReasonActive, "assignment is not cleanup-pending", "a1"},
		{"age/orphaned-active", func() AutomaticFacts {
			o := pendingOutcome()
			o.AssignmentState = intent.StateActive
			return AutomaticFacts{Explicit: o, OrphanedActive: true}
		}, ActionRetain, ReasonOrphaned, "assignment is not cleanup-pending", "a1"},
		{"age/landed-active", func() AutomaticFacts {
			o := pendingOutcome()
			o.AssignmentState = intent.StateActive
			return AutomaticFacts{Explicit: o, Landed: true}
		}, ActionRetain, ReasonLanded, "assignment is not cleanup-pending", "a1"},
		{"eligibility/recovered-state-uncertain", func() AutomaticFacts {
			o := pendingOutcome()
			o.AssignmentState = intent.StateRecovered
			return AutomaticFacts{Explicit: o}
		}, ActionRetain, ReasonUncertain, "assignment is not cleanup-pending", "a1"},
		{"eligibility/recovery-mismatch", func() AutomaticFacts {
			return AutomaticFacts{Explicit: pendingOutcome()}
		}, ActionRetain, ReasonMalformed, "assignment recovery metadata does not match refs", "a1"},
		{"eligibility/landedness-unknown", func() AutomaticFacts {
			o := pendingOutcome()
			o.Landed = Landedness{Kind: LandednessUnknownNoDefault}
			return AutomaticFacts{Explicit: o, RecoveryMatches: true}
		}, ActionRetain, ReasonUncertain, "assignment landedness is unknown", "a1"},
		{"eligibility/unmerged", func() AutomaticFacts {
			o := pendingOutcome()
			o.Landed = Landedness{Kind: LandednessProven, Landed: false}
			return AutomaticFacts{Explicit: o, RecoveryMatches: true}
		}, ActionRetain, ReasonUnmerged, "assignment branch has not landed", "a1"},
		{"preservation/dirty-refused", func() AutomaticFacts {
			o := pendingOutcome()
			o.Action, o.Tracked = ActionRecoverRemove, "dirty"
			return AutomaticFacts{Explicit: o, RecoveryMatches: true}
		}, ActionRetain, ReasonDirty, "automatic cleanup does not preserve uncommitted work", "a1"},
		{"action/clean-pending-removes", func() AutomaticFacts {
			return AutomaticFacts{Explicit: pendingOutcome(), RecoveryMatches: true}
		}, ActionRemove, "", "", "a1"},
		{"action/declared-discard-remove-passes", func() AutomaticFacts {
			o := pendingOutcome()
			o.Action = ActionDiscardRemove
			return AutomaticFacts{Explicit: o, RecoveryMatches: true}
		}, ActionDiscardRemove, "", "", "a1"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := DecideAutomatic(c.facts())
			if v.Action != c.action || v.ReasonCode != c.reasonCode || v.Reason != c.reason || v.AssignmentID != c.assignment {
				t.Fatalf("verdict = %s/%s/%q/%q, want %s/%s/%q/%q", v.Action, v.ReasonCode, v.Reason, v.AssignmentID, c.action, c.reasonCode, c.reason, c.assignment)
			}
		})
	}
}

func TestActionRemovesTable(t *testing.T) {
	removes := map[Action]bool{
		ActionRemove: true, ActionRecoverRemove: true, ActionDiscardRemove: true, ActionReleaseRemove: true,
		ActionRetain: false, ActionRemoved: false, ActionError: false, ActionReleaseLeftover: false,
	}
	for action, want := range removes {
		if action.Removes() != want {
			t.Fatalf("%s.Removes() = %v, want %v", action, !want, want)
		}
	}
}

func TestPreservesTable(t *testing.T) {
	cases := []struct {
		name     string
		action   Action
		tracked  string
		detached bool
		want     bool
	}{
		{"preservation/recover-remove", ActionRecoverRemove, "clean", false, true},
		{"preservation/discard-remove-dirty", ActionDiscardRemove, "dirty", false, true},
		{"preservation/discard-remove-clean", ActionDiscardRemove, "clean", false, false},
		{"preservation/detached-registration", ActionRemove, "clean", true, true},
		{"preservation/plain-remove", ActionRemove, "clean", false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Preserves(c.action, c.tracked, c.detached); got != c.want {
				t.Fatalf("Preserves(%s,%s,%v) = %v, want %v", c.action, c.tracked, c.detached, got, c.want)
			}
		})
	}
}

func TestLeaseOwnerPIDTable(t *testing.T) {
	cases := []struct {
		name    string
		content string
		pid     int
		ok      bool
	}{
		{"lease/well-formed", "42 2026-07-15T00:00:00Z\n", 42, true},
		{"lease/empty", "", 0, false},
		{"lease/no-newline", "42 2026-07-15T00:00:00Z", 0, false},
		{"lease/two-lines", "42 2026-07-15T00:00:00Z\nx\n", 0, false},
		{"lease/extra-field", "42 2026-07-15T00:00:00Z extra\n", 0, false},
		{"lease/zero-pid", "0 2026-07-15T00:00:00Z\n", 0, false},
		{"lease/padded-pid", "042 2026-07-15T00:00:00Z\n", 0, false},
		{"lease/bad-stamp", "42 yesterday\n", 0, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			pid, ok := LeaseOwnerPID([]byte(c.content))
			if ok != c.ok || (ok && pid != c.pid) {
				t.Fatalf("LeaseOwnerPID(%q) = (%d,%v), want (%d,%v)", c.content, pid, ok, c.pid, c.ok)
			}
		})
	}
}

func TestReclaimableTable(t *testing.T) {
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	staleAfter := time.Hour
	dead := func(int) bool { return false }
	live := func(int) bool { return true }
	wellFormed := []byte("42 2026-07-15T00:00:00Z\n")
	cases := []struct {
		name    string
		content []byte
		mtime   time.Time
		alive   func(int) bool
		want    bool
	}{
		{"lease/dead-owner", wellFormed, now, dead, true},
		{"lease/live-owner", wellFormed, now.Add(-100 * time.Hour), live, false},
		{"lease/malformed-fresh", []byte("junk\n"), now.Add(-time.Minute), dead, false},
		{"lease/malformed-stale", []byte("junk\n"), now.Add(-100 * time.Hour), live, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Reclaimable(c.content, c.mtime, now, c.alive, staleAfter); got != c.want {
				t.Fatalf("Reclaimable = %v, want %v", got, c.want)
			}
		})
	}
}

func TestOrphanedTable(t *testing.T) {
	now := time.Date(2026, 7, 15, 12, 0, 0, 0, time.UTC)
	staleAfter := 7 * 24 * time.Hour
	stamp := func(t time.Time) *string { s := t.Format(time.RFC3339); return &s }
	unparseable := "yesterday"
	cases := []struct {
		name string
		a    intent.Assignment
		want bool
	}{
		{"age/non-active-never-orphaned", intent.Assignment{State: intent.StateCleanupPending}, false},
		{"age/absent-stamp-aged", intent.Assignment{State: intent.StateActive}, true},
		{"age/young", intent.Assignment{State: intent.StateActive, CreatedAt: stamp(now.Add(-time.Hour))}, false},
		{"age/aged", intent.Assignment{State: intent.StateActive, CreatedAt: stamp(now.Add(-8 * 24 * time.Hour))}, true},
		{"age/future-stamp-not-aged", intent.Assignment{State: intent.StateActive, CreatedAt: stamp(now.Add(time.Hour))}, false},
		{"age/unparseable-stamp-unknown", intent.Assignment{State: intent.StateActive, CreatedAt: &unparseable}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Orphaned(c.a, now, staleAfter); got != c.want {
				t.Fatalf("Orphaned = %v, want %v", got, c.want)
			}
		})
	}
}

func TestResidualTable(t *testing.T) {
	if !Residual(intent.Assignment{}) {
		t.Fatal("empty recovery set is residue")
	}
	if Residual(intent.Assignment{Recovery: []intent.Recovery{{Ref: "refs/bench/r1"}}}) {
		t.Fatal("recorded recovery is preserved work, never residue")
	}
}

func TestLandednessStringTable(t *testing.T) {
	cases := []struct {
		name string
		l    Landedness
		want string
	}{
		{"action/detached", Landedness{Kind: LandednessDetached}, "detached"},
		{"action/no-default", Landedness{Kind: LandednessUnknownNoDefault}, "unknown"},
		{"action/query-error", Landedness{Kind: LandednessUnknownError, Err: "boom"}, "unknown:boom"},
		{"action/proven-ancestry", Landedness{Kind: LandednessProven, Landed: true}, "true:ancestry"},
		{"action/proven-patch", Landedness{Kind: LandednessProven, Landed: true, ByContent: true}, "true:patch"},
		{"action/proven-unmerged", Landedness{Kind: LandednessProven}, "false:ancestry"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.l.String(); got != c.want {
				t.Fatalf("String() = %q, want %q", got, c.want)
			}
		})
	}
}
