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
// and runs publish while that decision still holds.
type StackAdmission interface {
	CheckMember(i int, tree string) error
	PublishStack(k int, trees []string, publish func() error) error
}

// StackResult is the receipt of a landing stack. Merges holds the published merge of each
// published member, in file order, and Kind is the authorization kind of the top tree.
type StackResult struct {
	DestinationBase string
	Merges          []string
	Kind            authorization.Kind
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
// merge for each member, and grades the top tree once. Each member's admission decides its
// own landing tree: the destination with that member alone. A green publishes every member
// by one compare-and-swap of the destination to the top merge. Any other kind publishes
// nothing and returns the authorization refusal.
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
	merges := make([]string, len(members))
	previous, top := destination, ""
	for i, m := range members {
		tree, err := o.landingTree(m, previous, sources[i])
		if err != nil {
			return StackResult{}, StackMemberError{i + 1, err}
		}
		if merges[i], err = landingMerge(r.Root, tree, previous, sources[i], m.Message); err != nil {
			return StackResult{}, err
		}
		previous, top = merges[i], tree
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
	ctx = gate.WithStack(gate.WithLanding(ctx, destination, r.GreenMarker), stack)
	got := o.authorize(ctx, r.Root, top, r.Stdout, r.Stderr)
	result := StackResult{DestinationBase: destination, Kind: got.Kind}
	if !o.reviewedPublishes.permits(got.Kind) {
		return result, AuthorizationRefusal{got}
	}
	for i, m := range members {
		if err := sourceUnchanged(m, sources[i]); err != nil {
			return result, StackMemberError{i + 1, err}
		}
	}
	if err := destinationUnchanged(r.Root, r.Destination, destination, r.DestinationFingerprint); err != nil {
		return result, err
	}
	if err := admission.PublishStack(len(members), trees, func() error {
		return o.publishMerge(r.Root, r.Destination, previous, destination)
	}); err != nil {
		return result, err
	}
	result.Merges = merges
	return result, nil
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
