package repository_test

import (
	"reflect"
	"testing"

	"github.com/gibbonmi/bench/internal/commitment"
	"github.com/gibbonmi/bench/internal/commitment/commitmenttest"
	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/intent"
)

// stackReason is the reason of each override that a stack test grants.
const stackReason = "back out a landed change"

// stackRepo is the publicationRepo with a second bound assignment of the same delivery, so a
// stack of two bound members passes the open chunk gap.
func stackRepo(t *testing.T) (root string, bound, second, unbound commitrepo.Publication, tree string) {
	t.Helper()
	root, bound, unbound, tree = publicationRepo(t)
	second = publication(t, root, commitmenttest.BindBuild(t, root, "second", publicationDeliverable))
	return root, bound, second, unbound, tree
}

// grantOverrides grants an override to the assignment of each member.
func grantOverrides(t *testing.T, root string, members ...commitrepo.Publication) {
	t.Helper()
	for _, member := range members {
		if err := (commitrepo.Store{Root: root}).Override(member.Assignment, stackReason); err != nil {
			t.Fatal(err)
		}
	}
}

// publishStack publishes members, each over tree, and returns how many times publish ran.
func publishStack(t *testing.T, root, tree string, members ...commitrepo.Publication) (int, error) {
	t.Helper()
	trees := make([]string, len(members))
	for i := range trees {
		trees[i] = tree
	}
	published := 0
	err := commitrepo.Store{Root: root}.PublishStackAdmitted(members, trees, func() error {
		published++
		return nil
	})
	return published, err
}

// requirePublished fails unless the stack publication ran publish once and succeeded.
func requirePublished(t *testing.T, published int, err error) {
	t.Helper()
	if err != nil || published != 1 {
		t.Fatalf("stack publication = %v, publish ran %d times; want one publish", err, published)
	}
}

// requireStackOverrides fails unless the ledger at root holds exactly want.
func requireStackOverrides(t *testing.T, root string, want []intent.LandingOverride) {
	t.Helper()
	got, err := commitmenttest.Overrides(root)
	if err != nil || len(got)+len(want) > 0 && !reflect.DeepEqual(got, want) {
		t.Fatalf("overrides after the stack publication = %+v (%v), want %+v", got, err, want)
	}
}

// A closed gap refuses the second member, whose assignment holds no override. The first
// member passes through its override, so only a decision that reads each member's window
// raises the window refusal.
func TestStackPublicationRefusesARefusedMember(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	putGap(t, root, intent.ChunkGap{Outcome: commitmenttest.DeliveryOutcome, Chunk: "C1"})
	grantOverrides(t, root, bound)
	bound.Destination, unbound.Destination = leasedDestination, leasedDestination
	published, err := publishStack(t, root, tree, bound, unbound)
	if published != 0 {
		t.Fatalf("publish ran %d times for a refused stack, want none", published)
	}
	requireWindowRefusal(t, err, commitmenttest.DeliveryOutcome+" ("+publicationDeliverable+")")
}

// A stack with no member, or with no tree for one member, is refused before publish.
func TestStackPublicationRefusesAStackWithoutATreeForEachMember(t *testing.T) {
	root, bound, _, tree := publicationRepo(t)
	for _, row := range []struct {
		name    string
		members []commitrepo.Publication
		trees   []string
	}{
		{name: "empty"},
		{name: "missing-tree", members: []commitrepo.Publication{bound, bound}, trees: []string{tree}},
	} {
		t.Run(row.name, func(t *testing.T) {
			published := false
			err := commitrepo.Store{Root: root}.PublishStackAdmitted(row.members, row.trees, func() error {
				published = true
				return nil
			})
			if err == nil || published {
				t.Fatalf("stack publication = %v (published=%t), want a refusal before publish", err, published)
			}
		})
	}
}

func TestStackPublicationConsumesEachOverride(t *testing.T) {
	root, bound, second, _, tree := stackRepo(t)
	grantOverrides(t, root, bound, second)
	published, err := publishStack(t, root, tree, bound, second)
	requirePublished(t, published, err)
	requireStackOverrides(t, root, nil)
}

func TestStackPublicationKeepsAnUnpublishedOverride(t *testing.T) {
	root, bound, unbound, tree := publicationRepo(t)
	grantOverrides(t, root, bound, unbound)
	published, err := publishStack(t, root, tree, bound)
	requirePublished(t, published, err)
	requireStackOverrides(t, root, []intent.LandingOverride{{Assignment: unbound.Assignment, Reason: stackReason}})
}

// The first member holds the one lease of the stack. The decision drops it, so the next
// landing takes the lease at once.
func TestStackPublicationDropsTheLease(t *testing.T) {
	root, bound, second, unbound, tree := stackRepo(t)
	published, err := publishStack(t, root, tree, leased(t, root, bound), second)
	requirePublished(t, published, err)
	if leases := ledgerLeases(t, root); len(leases) != 0 {
		t.Fatalf("leases after the stack publication = %+v, want none", leases)
	}
	requireAcquire(t, root, unbound.Assignment, "after the stack publication")
}

// requireBinding fails unless the binding query for assignment answers want.
func requireBinding(t *testing.T, root, assignment string, want bool) {
	t.Helper()
	got, err := commitrepo.Store{Root: root}.HoldsDeliveryBinding(assignment)
	if err != nil || got != want {
		t.Fatalf("HoldsDeliveryBinding(%s) = %t, %v; want %t", assignment, got, err, want)
	}
}

// A started tickets-only deliverable binds its assignment but opens no build.
func TestBindingQueryNamesATicketsOnlyAssignment(t *testing.T) {
	root := commitmenttest.Staged(t, "A")
	identity := commitmenttest.WriteTickets(t, root)
	commitmenttest.EditPolicy(t, root, func(policy *commitment.Policy) {
		delivery := &policy.Milestones[0].Outcomes[0]
		delivery.Deliverables = append(delivery.Deliverables, commitmenttest.TicketsBinding(identity))
	})
	commitmenttest.Commit(t, root, "approve the tickets-only folder")
	light := startBuild(t, root, "light", "A", commitmenttest.TicketsFolder)
	requireBinding(t, root, light.Assignment, true)
}

func TestBindingQueryLeavesALightPathAssignment(t *testing.T) {
	root := commitmenttest.Staged(t, "A")
	startBuild(t, root, "build", "A", approvedSpec(t, root, "A"))
	unbound := windowSource(t, root, commitmenttest.Assignment(t, root, "unbound"))
	requireBinding(t, root, unbound.Assignment, false)
}
