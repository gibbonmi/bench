package landing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/gibbonmi/bench/internal/gate"
	"github.com/gibbonmi/bench/internal/gate/authorization"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// StackRequest is the complete input to one landing stack: the destination snapshot and
// one reviewed request for each member, in file order. The stack fields own the
// destination, so the destination fields of each member request are not read.
type StackRequest struct {
	Root, Destination, DestinationBase  string
	DestinationFingerprint, GreenMarker string
	Members                             []ReviewedRequest
	Stdout, Stderr                      io.Writer
}

// StackAdmission is the publication decision of a landing stack. CheckMember decides
// member i, an index into StackRequest.Members, over that member's own landing tree before
// the grade. PublishStack decides the first k members again over their own landing trees
// and runs publish while that decision still holds. trees holds exactly k trees: the own
// landing tree of each published member, in file order.
type StackAdmission interface {
	CheckMember(i int, tree string) error
	PublishStack(k int, trees []string, publish func() error) error
}

// StackResult is the receipt of a landing stack. Merges holds the published merge of each
// published member, in file order. Kind is the authorization kind of the top tree, or the
// kind that stopped the bisect. Attributed is the 1-based member that the bisect names for
// the red, or 0 when no member carries it.
type StackResult struct {
	DestinationBase string
	Merges          []string
	Kind            authorization.Kind
	Attributed      int
}

// Published is the count of published members.
func (r StackResult) Published() int { return len(r.Merges) }

// StackMemberError is a refusal that one member of a stack causes. Member is the 1-based
// index of that member in file order, as an operator reads the stack file.
type StackMemberError struct {
	Member int
	Err    error
}

func (e StackMemberError) Error() string { return fmt.Sprintf("stack member %d: %v", e.Member, e.Err) }
func (e StackMemberError) Unwrap() error { return e.Err }

// LandStack composes each member onto the merge of the member before it, writes one landing
// merge for each member, and grades the top tree. Each member's admission decides its own
// landing tree: the destination with that member alone. A green publishes every member by
// one compare-and-swap of the destination to the top merge. A candidate red bisects the
// prefixes, attributes the red to the first red prefix above a green one, publishes that
// green prefix, and returns the red refusal. Any other kind on any grade publishes nothing
// and returns the authorization refusal.
func (o Owner) LandStack(ctx context.Context, r StackRequest, admission StackAdmission) (StackResult, error) {
	if admission == nil {
		return StackResult{}, errors.New("landing stack requires a publication admission")
	}
	if len(r.Members) == 0 {
		return StackResult{}, errors.New("landing stack has no members")
	}
	members := slices.Clone(r.Members)
	sources := make([]string, len(members))
	var destination string
	for i := range members {
		m := &members[i]
		// The stack never transitions a staged spec; that member lands alone.
		if m.SpecPath != "" {
			return StackResult{}, StackMemberError{i + 1, errors.New("a staged-spec member lands alone")}
		}
		m.Root, m.Destination, m.DestinationBase = r.Root, r.Destination, r.DestinationBase
		m.DestinationFingerprint, m.GreenMarker = r.DestinationFingerprint, r.GreenMarker
		m.Stdout, m.Stderr = r.Stdout, r.Stderr
		var err error
		if destination, sources[i], err = reviewedSource(*m); err != nil {
			return StackResult{}, StackMemberError{i + 1, err}
		}
	}
	// prefixTrees[i] and merges[i] are the tree and the merge of the first i+1 members.
	merges, prefixTrees := make([]string, len(members)), make([]string, len(members))
	previous := destination
	for i, m := range members {
		tree, err := o.landingTree(m, previous, sources[i])
		if err != nil {
			return StackResult{}, StackMemberError{i + 1, err}
		}
		if merges[i], err = landingMerge(r.Root, tree, previous, sources[i], m.Message); err != nil {
			return StackResult{}, err
		}
		previous, prefixTrees[i] = merges[i], tree
	}
	trees := make([]string, len(members))
	stack := make([]gate.StackMember, len(members))
	for i, m := range members {
		tree, err := o.landingTree(m, destination, sources[i])
		if err == nil {
			err = admission.CheckMember(i, tree)
		}
		if err != nil {
			return StackResult{}, StackMemberError{i + 1, err}
		}
		trees[i], stack[i] = tree, gate.StackMember{Source: sources[i], Folder: m.ClosePath}
	}
	ctx = gate.WithLanding(ctx, destination, r.GreenMarker)
	grade := func(k int) authorization.Result {
		return o.authorize(gate.WithStack(ctx, stack[:k]), r.Root, prefixTrees[k-1], r.Stdout, r.Stderr)
	}
	got := grade(len(members))
	result := StackResult{DestinationBase: destination, Kind: got.Kind}
	lo, refusal := len(members), error(nil)
	if !o.reviewedPublishes.permits(got.Kind) {
		if got.Kind != authorization.Candidate {
			return result, AuthorizationRefusal{got}
		}
		// A candidate red means that the destination held a green, so prefix 0 is green.
		hi := lo
		lo = 0
		for hi-lo > 1 {
			mid := (lo + hi) / 2
			switch prefix := grade(mid); {
			case o.reviewedPublishes.permits(prefix.Kind):
				lo = mid
			case prefix.Kind == authorization.Candidate:
				hi, got = mid, prefix
			default:
				result.Kind = prefix.Kind
				return result, AuthorizationRefusal{prefix}
			}
		}
		result.Attributed, refusal = hi, AuthorizationRefusal{got}
		if lo == 0 {
			return result, refusal
		}
	}
	for i, m := range members[:lo] {
		if err := sourceUnchanged(m, sources[i]); err != nil {
			return result, StackMemberError{i + 1, err}
		}
	}
	if err := destinationUnchanged(r.Root, r.Destination, destination, r.DestinationFingerprint); err != nil {
		return result, err
	}
	if err := admission.PublishStack(lo, trees[:lo], func() error {
		return o.publishMerge(r.Root, r.Destination, merges[lo-1], destination)
	}); err != nil {
		return result, err
	}
	result.Merges = merges[:lo]
	return result, refusal
}

// sourceUnchanged rechecks the two moving identities of a reviewed source after the gate and
// before an otherwise unreachable object publishes. Tree equality is insufficient: review
// binds a commit.
func sourceUnchanged(r ReviewedRequest, source string) error {
	if branchTip, err := benchgit.ResolveCommit(r.Root, r.Source); err != nil || branchTip != source {
		return errors.New("reviewed source tip moved")
	}
	if fingerprint, err := CheckoutFingerprint(r.SourceWorktree); err != nil || fingerprint != r.SourceFingerprint {
		return errors.New("reviewed source checkout changed")
	}
	return nil
}

// landingMerge writes the landing commit of one reviewed source: first the commit it lands
// onto, then the source tip. No ref names it.
func landingMerge(root, tree, first, source, message string) (string, error) {
	created, err := output(root, "commit-tree", tree, "-p", first, "-p", source, "-m", message)
	if err != nil {
		return "", fmt.Errorf("create landing commit: %w", err)
	}
	return created, nil
}

// publishMerge moves ref from base to commit by compare-and-swap.
func (o Owner) publishMerge(root, ref, commit, base string) error {
	if err := o.updateRef(root, ref, commit, base); err != nil {
		return destinationUpdateFailure(root, ref, base, err)
	}
	return nil
}
