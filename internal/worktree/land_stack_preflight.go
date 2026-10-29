// Stack preflight proofs: the proofs that a stack runs for each member after the source
// proofs of one landing, the owner rule over the member trees, and the build window.
package worktree

import (
	"cmp"
	"fmt"
	"slices"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/stackpolicy"
	"github.com/gibbonmi/bench/internal/refusalroute"
)

// stackMemberProofs runs the stack proofs of one member after its source proofs. A member
// that names a staged spec or whose assignment holds a delivery binding lands alone, because
// its landing transitions that spec or closes that deliverable. A member whose tip lacks the
// destination tip folds the destination first. The first fault stops the rest. An empty
// destination is one that its own proofs refused, so no ancestry proof runs over it.
func stackMemberProofs(run *git.Run, root, destination, branch string, l stackMemberLanding) error {
	if l.source.specPath != "" {
		return stackAloneRefusal("its spec is a staged spec")
	}
	bound, err := commitrepo.Store{Root: root, Run: run}.HoldsDeliveryBinding(l.assignment.ID)
	if err != nil {
		return err
	}
	if bound {
		return stackAloneRefusal("its assignment holds a delivery binding")
	}
	if destination == "" {
		return nil
	}
	return stackBehind(root, destination, branch, l)
}

// stackBehind is the behind proof of member l: a member whose tip lacks the destination tip
// refuses with the fold route of one landing, which merges from branch.
func stackBehind(root, destination, branch string, l stackMemberLanding) error {
	holds, err := authorization.IsAncestor(root, destination, l.source.tip)
	if err != nil {
		return err
	}
	if !holds {
		return destinationMovedRefusal(fmt.Errorf("source tip %s lacks destination tip %s", l.source.tip, destination), branch)
	}
	return nil
}

// stackRaisedRefusal routes raised, a refusal that the stack raised for member i, with the
// re-run of fact and the member label as its route facts.
func stackRaisedRefusal(route stackRoute, i int, label string, raised error, fact stackRerun) error {
	values := map[string]string{refusalroute.FactRerun: route.rerun(fact), refusalroute.FactLabel: label}
	return stackMemberRefusal(i, label, landingFaceRouteWith(raised, values, false))
}

// stackAloneRefusal raises the face of a member that lands alone, with cause after the
// face's sentence.
func stackAloneRefusal(cause string) refusalError {
	return refusalError{refusal{detail: refusalroute.Sentence(faceStackMemberAlone) + causeJoiner + cause, face: faceStackMemberAlone}}
}

// stackOverlaps runs the owner rule over the destination tree and the tree of each member
// that cleared its proofs. Each pair of members that changes one path outside capture/
// refuses the later member of the pair, and the refusal names both labels and the paths.
func stackOverlaps(root, destination string, route stackRoute, landings []stackMemberLanding) []error {
	if destination == "" || len(landings) < 2 {
		return nil
	}
	base, err := gate.TreeEntries(root, destination)
	if err != nil {
		return []error{landingFaceRoute(err, route.rerun(stackUnattributed), false)}
	}
	trees := make([]stackpolicy.Entries, len(landings))
	for i, l := range landings {
		if trees[i], err = gate.TreeEntries(root, l.source.tip); err != nil {
			return []error{stackMemberRefusal(l.index, l.assignment.Label, landingFaceRoute(err, route.rerun(stackRefusedBy(l.index, err)), false))}
		}
	}
	// The owner rule names each member by its 1-based position in trees.
	pairs := map[[2]int][]string{}
	for path, owners := range stackpolicy.Overlaps(stackpolicy.Owners(base, trees)) {
		for a := range owners {
			for _, b := range owners[a+1:] {
				pair := [2]int{owners[a], b}
				pairs[pair] = append(pairs[pair], path)
			}
		}
	}
	keys := make([][2]int, 0, len(pairs))
	for pair := range pairs {
		keys = append(keys, pair)
	}
	slices.SortFunc(keys, func(x, y [2]int) int { return cmp.Or(cmp.Compare(x[0], y[0]), cmp.Compare(x[1], y[1])) })
	var refusals []error
	for _, pair := range keys {
		paths := pairs[pair]
		slices.Sort(paths)
		earlier, later := landings[pair[0]-1], landings[pair[1]-1]
		cause := fmt.Sprintf("members %d (%s) and %d (%s) change the refusal_paths entries", earlier.index+1, earlier.assignment.Label, later.index+1, later.assignment.Label)
		raised := refusalError{refusal{detail: refusalroute.Sentence(faceStackMembersOverlap) + causeJoiner + cause, paths: paths, face: faceStackMembersOverlap}}
		refusals = append(refusals, stackRaisedRefusal(route, later.index, later.assignment.Label, raised, stackRefusedBy(later.index, raised)))
	}
	return refusals
}

// stackWindow decides the build window for each member on its own before the lease, as for
// one landing. Each member that the window refuses refuses with its index.
func stackWindow(run *git.Run, root string, route stackRoute, landings []stackMemberLanding, publications []commitrepo.Publication) []error {
	var refusals []error
	for i, l := range landings {
		if err := (commitrepo.Store{Root: root, Run: run}).AdmitWindow(publications[i]); err != nil {
			refusals = append(refusals, stackMemberRefusal(l.index, l.assignment.Label, landingFaceRoute(err, route.rerun(stackRefusedBy(l.index, err)), false)))
		}
	}
	return refusals
}
