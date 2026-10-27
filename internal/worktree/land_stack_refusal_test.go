// Stack refusal tests: the stack proofs refuse a member before the lease, so the gate never
// runs, and the build window decides each member on its own.
package worktree

import (
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	"github.com/gibbonmi/bench/internal/intent"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// requirePreflightRefusal requires that the stack r refused exactly once, for member 1-based
// member m, with route want, and that it took no lease and ran no grade.
func requirePreflightRefusal(t *testing.T, f stackFixture, r verbResult, member string, m stackFixtureMember, want string) {
	t.Helper()
	if next := stackRefusalNext(t, r.stdout, member, m); r.exit != 1 || next != want {
		t.Fatalf("stack = (%d, %q, %q), want exit 1 and next %q", r.exit, r.stdout, r.stderr, want)
	}
	if got := strings.Count(r.stdout, "refused{"); got != 1 {
		t.Fatalf("stack stdout = %q, want one refusal, got %d", r.stdout, got)
	}
	requireNoLease(t, f.root)
	if _, err := os.Stat(f.tally); !os.IsNotExist(err) {
		t.Fatalf("stack ran the gate before its refusal: %v", err)
	}
}

// singleStackRerun is the single landing of member m at its own tip, which a route that
// leaves m alone in the stack re-runs.
func singleStackRerun(f stackFixture, m stackFixtureMember) string {
	return landingRerunAt(m.request, f.base, landingSourceTipFlag(m.tip), m.slug, m.creation.Path, m.creation.Assignment.ID)
}

// stackRepairRerun is the stack re-run that sets the repaired slot of member m.
func stackRepairRerun(file string, m stackFixtureMember) string {
	return "bench worktree land --stack '" + file + "' --source-tip '" + m.request + "'=<repaired-tip>"
}

// foldMember merges the destination into member i of f and moves its tip there.
func foldMember(t *testing.T, f *stackFixture, i int) {
	t.Helper()
	m := &f.members[i]
	gitRun(t, m.creation.Path, "-c", "user.name=bench", "-c", "user.email=bench@local", "merge", "-q", "--no-ff", "-m", "fold the destination", "main")
	m.tip = gitOutput(t, m.creation.Path, "rev-parse", "HEAD")
}

// TestLandStackRefusesABehindMember is LS70. A member whose tip lacks the destination tip
// refuses with its fold route before the lease. A destination that moves after the grade
// refuses with the same face for member 1.
func TestLandStackRefusesABehindMember(t *testing.T) {
	t.Parallel()
	moved := func(m stackFixtureMember) map[string]string {
		return map[string]string{refusalroute.FactLabel: m.creation.Assignment.Label, refusalroute.FactFrom: "main"}
	}
	t.Run("preflight", func(t *testing.T) {
		t.Parallel()
		f := stackLandingFixture(t, "stack-behind", 2, nil)
		advanceDestination(t, f.root)
		foldMember(t, &f, 0)
		file, args := f.stack(t)
		r := runVerb(t, verbLand, f.call(args...))
		m := f.members[1]
		requirePreflightRefusal(t, f, r, "2", m, landingRoute(faceLandDestinationMoved, stackRepairRerun(file, m), moved(m)))
	})
	t.Run("after-the-grade", func(t *testing.T) {
		t.Parallel()
		f := stackLandingFixture(t, "stack-moved", 2, nil)
		j := defaultJoins()
		j.publicationGap = func(root string) {
			commitInWorktree(t, root, "moved.txt", "destination bytes\n", "destination moves in the publication gap")
		}
		file, args := f.stack(t)
		r := runVerb(t, verbLand, f.callWith(j, args...))
		r.mustViaJoins(t)
		m := f.members[0]
		want := landingRoute(faceLandDestinationMoved, stackRepairRerun(file, m), moved(m))
		if next := stackRefusalNext(t, r.stdout, "1", m); r.exit != 1 || next != want {
			t.Fatalf("stack = (%d, %q, %q), want exit 1 and next %q", r.exit, r.stdout, r.stderr, want)
		}
		requireOneRun(t, f.tally, "after the moved stack")
	})
}

// TestLandStackRefusesAnOverlap is LS71. Two members that change one path outside capture/
// refuse the later member, and the refusal names both labels and the path. Its route leaves
// member 1 alone, so the re-run is member 1's single landing.
func TestLandStackRefusesAnOverlap(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-overlap", 2, nil)
	first, second := f.members[0], &f.members[1]
	commitInWorktree(t, second.creation.Path, first.file, "member 2 bytes\n", "member 2 changes the file of member 1")
	second.tip = gitOutput(t, second.creation.Path, "rev-parse", "HEAD")
	_, args := f.stack(t)
	r := runVerb(t, verbLand, f.call(args...))
	requirePreflightRefusal(t, f, r, "2", *second, landingRoute(faceStackMembersOverlap, singleStackRerun(f, first), labelOf(second.creation)))
	for _, named := range []string{first.creation.Assignment.Label, second.creation.Assignment.Label, first.file} {
		if !strings.Contains(r.stdout, named) {
			t.Fatalf("stack stdout = %q, want it to name %q", r.stdout, named)
		}
	}
}

// TestLandStackRefusesAStagedSpecMember is LS72. A member whose spec names a staged spec
// lands alone, so the re-run is the other member's single landing.
func TestLandStackRefusesAStagedSpecMember(t *testing.T) {
	t.Parallel()
	f := stackLandingFixture(t, "stack-staged", 2, nil)
	members := slices.Clone(f.members)
	members[0].slug = "x"
	file := f.writeStack(t, members, f.tips()...)
	r := runVerb(t, verbLand, f.call("--stack", file))
	requirePreflightRefusal(t, f, r, "1", f.members[0], landingRoute(faceStackMemberAlone, singleStackRerun(f, f.members[1]), labelOf(f.members[0].creation)))
}

// TestLandStackRefusesABoundMember is LS73 and LS88. A member whose assignment holds a
// delivery binding lands alone: a started tickets-only deliverable, which opens no build,
// and a staged spec that an open build binds.
func TestLandStackRefusesABoundMember(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name string
		bind func(t *testing.T, f *stackFixture, prefix string)
	}{
		{"tickets-only", func(t *testing.T, f *stackFixture, prefix string) {
			commitmenttest.ApproveTickets(t, f.root)
			commitmenttest.Commit(t, f.root, "approve the tickets-only folder")
			for i := range f.members {
				foldMember(t, f, i)
			}
			m := f.members[1]
			commitmenttest.Admit(t, m.creation.Path, m.request, commitmenttest.TicketsFolder)
			// The fixture's own build leaves, so no open build binds any member.
			build := intent.RequestDigest(prefix + "-build")
			commitmenttest.EditCommitment(t, f.root, func(state *intent.CommitmentState) {
				state.Bindings = slices.DeleteFunc(state.Bindings, func(b intent.DeliveryBinding) bool { return b.Request == build })
			})
		}},
		{"open-build", func(t *testing.T, f *stackFixture, _ string) {
			m := f.members[1]
			commitmenttest.Admit(t, m.creation.Path, m.request, closureSpec)
		}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			prefix := "stack-bound-" + row.name
			f := stackLandingFixture(t, prefix, 2, nil)
			row.bind(t, &f, prefix)
			_, args := f.stack(t)
			r := runVerb(t, verbLand, f.call(args...))
			requirePreflightRefusal(t, f, r, "2", f.members[1], landingRoute(faceStackMemberAlone, singleStackRerun(f, f.members[0]), labelOf(f.members[1].creation)))
		})
	}
}

// TestLandStackAdmitsEachMemberThroughTheWindow is LS74 and LS75. With an open build and no
// chunk gap, the window refuses each light-path member by its index before the lease. An
// override for member 1 alone admits member 1 and still refuses member 2.
func TestLandStackAdmitsEachMemberThroughTheWindow(t *testing.T) {
	t.Parallel()
	for _, row := range []struct {
		name     string
		override bool
		refused  []int
	}{
		{"no-gap", false, []int{0, 1}},
		{"override-for-member-1", true, []int{1}},
	} {
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()
			f := stackLandingFixture(t, "stack-window-"+row.name, 2, nil)
			commitmenttest.PutGap(t, f.root, intent.ChunkGap{Outcome: commitmenttest.DeliveryOutcome, Chunk: "C1"})
			if row.override {
				commitmenttest.EditCommitment(t, f.root, func(state *intent.CommitmentState) {
					state.Overrides = append(state.Overrides, intent.LandingOverride{Assignment: f.members[0].creation.Assignment.ID, Reason: "land member 1"})
				})
			}
			file, args := f.stack(t)
			r := runVerb(t, verbLand, f.call(args...))
			want := landingRoute(faceLandBuildWindow, "bench worktree land --stack '"+file+"'", nil)
			for _, i := range row.refused {
				if next := stackRefusalNext(t, r.stdout, strconv.Itoa(i+1), f.members[i]); next != want {
					t.Fatalf("member %d next = %q, want %q", i+1, next, want)
				}
			}
			if got := strings.Count(r.stdout, "refused{"); r.exit != 1 || got != len(row.refused) {
				t.Fatalf("stack = (%d, %q), want exit 1 and %d refusals", r.exit, r.stdout, len(row.refused))
			}
			requireNoLease(t, f.root)
		})
	}
}
