package landingpolicy

import "testing"

func supply(t *testing.T, value bool, called *bool) func() bool {
	t.Helper()
	return func() bool {
		if called != nil {
			*called = true
		}
		return value
	}
}

// refuse names a fact supplier the decision must never consult on this path.
func refuse(t *testing.T, name string) func() bool {
	t.Helper()
	return func() bool {
		t.Fatalf("decision consulted %s on a path that must not need it", name)
		return false
	}
}

// TestResidueDecisionTable is the residue half of LP1: every destination
// residue class maps to its exact refusal detail, and the safe partitions
// return the empty accept.
func TestResidueDecisionTable(t *testing.T) {
	clean := func() ResidueFacts {
		return ResidueFacts{NestedClean: true, StatusReadable: true, StatusWellFormed: true}
	}
	for _, tc := range []struct {
		name  string
		facts func(t *testing.T) ResidueFacts
		want  string
	}{
		{"clean-destination", func(t *testing.T) ResidueFacts {
			f := clean()
			f.IgnoredDeclared = refuse(t, "IgnoredDeclared")
			f.StagedMatchesPublished = refuse(t, "StagedMatchesPublished")
			return f
		}, ""},
		{"nested-repositories", func(t *testing.T) ResidueFacts {
			f := clean()
			f.NestedClean = false
			return f
		}, "landing destination has nested repositories"},
		{"status-unreadable", func(t *testing.T) ResidueFacts {
			f := clean()
			f.StatusReadable = false
			return f
		}, "landing destination status is unreadable"},
		{"status-malformed", func(t *testing.T) ResidueFacts {
			f := clean()
			f.StatusWellFormed = false
			return f
		}, "landing destination status is malformed"},
		{"undeclared-ignored-residue", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "!!", Path: "private/output"}}
			f.IgnoredDeclared = supply(t, false, nil)
			return f
		}, "landing destination has ignored residue"},
		{"declared-ignored-residue", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "!!", Path: "dist/bench"}}
			f.IgnoredDeclared = supply(t, true, nil)
			return f
		}, ""},
		{"untracked-collision", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "??", Path: "stray.txt"}}
			return f
		}, "landing destination has untracked collisions"},
		{"tracked-worktree-change", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: " M", Path: "tracked.txt"}}
			return f
		}, "landing destination has tracked-worktree changes"},
		{"staged-off-published-destination", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "M ", Path: "tracked.txt"}}
			f.DestinationAtPublished = false
			f.StagedMatchesPublished = refuse(t, "StagedMatchesPublished")
			return f
		}, "landing destination has staged changes"},
		{"staged-beyond-published-content", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "A ", Path: "extra.txt"}}
			f.DestinationAtPublished = true
			f.StagedMatchesPublished = supply(t, false, nil)
			return f
		}, "landing destination has staged changes"},
		{"staged-exactly-published-content", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "M ", Path: "tracked.txt"}}
			f.DestinationAtPublished = true
			f.StagedMatchesPublished = supply(t, true, nil)
			return f
		}, ""},
		{"empty-status-entry-skipped", func(t *testing.T) ResidueFacts {
			f := clean()
			f.Entries = []StatusEntry{{Status: "", Path: "noise"}}
			return f
		}, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Residue(tc.facts(t)); got != tc.want {
				t.Fatalf("Residue = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestResidueConsultsIgnoredAllowanceOnce proves the allowance fact is read
// at most once across many ignored entries.
func TestResidueConsultsIgnoredAllowanceOnce(t *testing.T) {
	calls := 0
	f := ResidueFacts{NestedClean: true, StatusReadable: true, StatusWellFormed: true,
		Entries: []StatusEntry{{Status: "!!", Path: "dist/a"}, {Status: "!!", Path: "dist/b"}, {Status: "!!", Path: "dist/c"}},
		IgnoredDeclared: func() bool {
			calls++
			return true
		},
	}
	if got := Residue(f); got != "" || calls != 1 {
		t.Fatalf("Residue = %q with %d allowance reads, want accept with 1", got, calls)
	}
}

// TestResumeMarkerDecisionTable is the resume half of LP1 for the marker.
func TestResumeMarkerDecisionTable(t *testing.T) {
	for _, tc := range []struct {
		name  string
		facts func(t *testing.T) MarkerFacts
		want  MarkerAction
	}{
		{"destination-at-published-advances", func(t *testing.T) MarkerFacts {
			return MarkerFacts{DestinationAtPublished: true, MarkerReachesPublished: refuse(t, "MarkerReachesPublished")}
		}, MarkerAdvance},
		{"moved-destination-with-covering-marker-accepts", func(t *testing.T) MarkerFacts {
			return MarkerFacts{MarkerPresent: true, MarkerReachesPublished: supply(t, true, nil)}
		}, MarkerAccept},
		{"moved-destination-without-marker-refuses", func(t *testing.T) MarkerFacts {
			return MarkerFacts{MarkerPresent: false, MarkerReachesPublished: refuse(t, "MarkerReachesPublished")}
		}, MarkerRefuse},
		{"moved-destination-with-divergent-marker-refuses", func(t *testing.T) MarkerFacts {
			return MarkerFacts{MarkerPresent: true, MarkerReachesPublished: supply(t, false, nil)}
		}, MarkerRefuse},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := ResumeMarker(tc.facts(t)); got != tc.want {
				t.Fatalf("ResumeMarker = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestPublicationDecisionTable is the publication half of LP1: each
// authentication fact flip produces its exact refusal, and the fully
// authenticated partitions accept.
func TestPublicationDecisionTable(t *testing.T) {
	authentic := func() PublicationFacts {
		return PublicationFacts{
			Requested: "1111", Resolved: true, Published: "1111",
			ReachableFromDestination: true, ParentsReadable: true,
			Parents: []string{"1111", "2222", "3333"}, RequestedSource: "3333",
			RangeAuthenticates: true,
		}
	}
	for _, tc := range []struct {
		name  string
		facts func() PublicationFacts
		want  Refusal
	}{
		{"authenticated-spec-less", authentic, Refusal{}},
		{"unresolved-value", func() PublicationFacts {
			f := authentic()
			f.Resolved = false
			return f
		}, Refusal{Detail: "published commit is not an exact commit identity"}},
		{"abbreviation-drift", func() PublicationFacts {
			f := authentic()
			f.Requested = "11"
			return f
		}, Refusal{Detail: "published commit is not an exact commit identity", Observed: "11", Wanted: "1111"}},
		{"unreachable-from-destination", func() PublicationFacts {
			f := authentic()
			f.ReachableFromDestination = false
			return f
		}, Refusal{Detail: "published commit is not reachable from the destination"}},
		{"unreadable-parents", func() PublicationFacts {
			f := authentic()
			f.ParentsReadable = false
			return f
		}, Refusal{Detail: "published commit parents are unreadable"}},
		{"non-merge-published-commit", func() PublicationFacts {
			f := authentic()
			f.Parents = []string{"1111", "2222"}
			return f
		}, Refusal{Detail: "published commit does not authenticate the reviewed source parent"}},
		{"foreign-source-parent", func() PublicationFacts {
			f := authentic()
			f.Parents[2] = "4444"
			return f
		}, Refusal{Detail: "published commit does not authenticate the reviewed source parent", Observed: "3333", Wanted: "4444"}},
		{"unauthenticated-review-base", func() PublicationFacts {
			f := authentic()
			f.RangeAuthenticates = false
			return f
		}, Refusal{Detail: "review base does not authenticate the published source"}},
		{"matching-spec-transition", func() PublicationFacts {
			f := authentic()
			f.Spec = SpecTransition{Named: true, TransitionMatches: true}
			return f
		}, Refusal{}},
		{"missing-spec-transition", func() PublicationFacts {
			f := authentic()
			f.Spec = SpecTransition{Named: true}
			return f
		}, Refusal{Detail: "published commit does not carry the source staged spec transition"}},
		{"closed-tickets-only-folder", func() PublicationFacts {
			f := authentic()
			f.Spec = SpecTransition{Named: true, TicketsOnlyClose: true}
			return f
		}, Refusal{}},
		{"unclosed-tickets-only-folder", func() PublicationFacts {
			f := authentic()
			f.Spec = SpecTransition{Named: true, TicketsOnlyClose: true, PublishedHasFolder: true}
			return f
		}, Refusal{Detail: "published commit does not close the source tickets-only folder"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Publication(tc.facts()); got != tc.want {
				t.Fatalf("Publication = %+v, want %+v", got, tc.want)
			}
		})
	}
}

// TestTerminalDecisionTable is the publish, release, and interruption half of
// LP1: each follow-up outcome maps to its exact worktree state and exit code.
func TestTerminalDecisionTable(t *testing.T) {
	for _, tc := range []struct {
		name  string
		facts TerminalFacts
		want  TerminalOutcome
	}{
		{"released", TerminalFacts{Active: true}, TerminalOutcome{ExitCode: 0, WorktreeState: "released"}},
		{"already-complete", TerminalFacts{Active: false}, TerminalOutcome{ExitCode: 0, WorktreeState: "already-complete"}},
		{"marker-interruption", TerminalFacts{FailedStep: "marker", Active: true}, TerminalOutcome{ExitCode: 3, WorktreeState: "incomplete:marker"}},
		{"reconcile-interruption", TerminalFacts{FailedStep: "reconcile", Active: true}, TerminalOutcome{ExitCode: 3, WorktreeState: "incomplete:reconcile"}},
		{"release-interruption", TerminalFacts{FailedStep: "release", Active: true}, TerminalOutcome{ExitCode: 3, WorktreeState: "incomplete:release"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := Terminal(tc.facts); got != tc.want {
				t.Fatalf("Terminal = %+v, want %+v", got, tc.want)
			}
		})
	}
}
