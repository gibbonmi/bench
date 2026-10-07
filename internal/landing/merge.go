package landing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/gibbonmi/bench/internal/gate/authorization"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// The three merge kinds one target-and-incoming pair can take. Ancestry is reflexive,
// so an equal pair satisfies both tests; MergeKindCurrent is decided first.
const (
	MergeKindCurrent     = "current"
	MergeKindFastForward = "fast-forward"
	MergeKindMerge       = "merge"
)

// MergeRequest names the immutable pair one merge composes and the target checkout the
// publication is bound to. Both commits are exact; the caller resolves every spelling.
// Root is the checkout the Git operations run in, and Worktree is the root the composed
// tree is graded from.
type MergeRequest struct {
	Root, Branch, PreviousTip, Incoming string
	Worktree, Fingerprint               string
	Subject                             string
	Stdout, Stderr                      io.Writer
}

// MergeResult is the publication receipt one merge returns. Resolved lists every capture
// path the composition policy settled, as "<path>:<side>", so the caller discloses what
// the merge did not decide.
type MergeResult struct {
	Kind                   string
	PreviousTip, Tip, Tree string
	Resolved               []string
}

// Merge publishes the incoming commit into the target branch under the owner's
// worktree-commit policy. It decides the kind by ancestry, composes a diverged pair
// through the capture rule table, authorizes the exact graded tree, and moves the branch
// by compare-and-swap on the previous tip. It does not reconcile the target checkout:
// the caller owns that, because only the caller knows the checkout is still its own.
func (o Owner) Merge(ctx context.Context, r MergeRequest) (MergeResult, error) {
	if r.Root == "" || r.Branch == "" || r.PreviousTip == "" || r.Incoming == "" || r.Worktree == "" || r.Fingerprint == "" || strings.TrimSpace(r.Subject) == "" {
		return MergeResult{}, errors.New("merge request is incomplete")
	}
	previous, err := compositionCommit(r.Root, r.PreviousTip, "previous tip")
	if err != nil || previous != r.PreviousTip {
		return MergeResult{}, errors.New("merge previous tip is not an exact commit")
	}
	incoming, err := compositionCommit(r.Root, r.Incoming, "incoming")
	if err != nil || incoming != r.Incoming {
		return MergeResult{}, errors.New("merge incoming commit is not exact")
	}
	if err := mergeTipUnmoved(r.Root, r.Branch, previous); err != nil {
		return MergeResult{}, err
	}
	// A failed ancestry query is an error, never a classification: a merge composed from
	// an unanswered question would publish an object the pair never justified.
	contains, err := authorization.IsAncestor(r.Root, incoming, previous)
	if err != nil {
		return MergeResult{}, fmt.Errorf("check whether the target contains the incoming commit: %w", err)
	}
	if contains {
		tree, err := output(r.Root, "rev-parse", previous+"^{tree}")
		if err != nil {
			return MergeResult{}, fmt.Errorf("read current target tree: %w", err)
		}
		return MergeResult{Kind: MergeKindCurrent, PreviousTip: previous, Tip: previous, Tree: tree}, nil
	}
	ahead, err := authorization.IsAncestor(r.Root, previous, incoming)
	if err != nil {
		return MergeResult{}, fmt.Errorf("check whether the incoming commit descends from the target: %w", err)
	}
	kind, tree, resolved := MergeKindFastForward, "", []string(nil)
	if ahead {
		if tree, err = output(r.Root, "rev-parse", incoming+"^{tree}"); err != nil {
			return MergeResult{}, fmt.Errorf("read incoming tree: %w", err)
		}
	} else {
		kind = MergeKindMerge
		composition, err := o.Compose(CompositionRequest{Root: r.Root, Destination: previous, Source: incoming})
		if err != nil {
			return MergeResult{}, err
		}
		if composition.Conflict.Kind != "" {
			return MergeResult{}, ConflictError{composition.Conflict}
		}
		tree, resolved = composition.Tree, composition.Resolved
	}
	// The authority grades from the target checkout, as `bench commit` in that checkout
	// does: the caller resolved the lane there, so the lane's file anchors name the target.
	// The caller's root can be another checkout of the same repository, where no anchor of
	// that lane points at the composed tree.
	if got := o.authorize(ctx, r.Worktree, tree, r.Stdout, r.Stderr); !o.publishes.permits(got.Kind) {
		return MergeResult{}, AuthorizationRefusal{got}
	}
	// Recheck both moving identities after the lane and before creating an otherwise
	// unreachable object.
	if err := mergeTipUnmoved(r.Root, r.Branch, previous); err != nil {
		return MergeResult{}, err
	}
	if fingerprint, fingerprintErr := CheckoutFingerprint(r.Worktree); fingerprintErr != nil || fingerprint != r.Fingerprint {
		return MergeResult{}, errors.New("merge target checkout changed")
	}
	tip := incoming
	if kind == MergeKindMerge {
		// The previous tip is the first parent, so the target's first-parent history
		// stays the history the review reads.
		if tip, err = output(r.Root, "commit-tree", tree, "-p", previous, "-p", incoming, "-m", r.Subject); err != nil {
			return MergeResult{}, fmt.Errorf("create merge commit: %w", err)
		}
	}
	if err := o.updateRef(r.Root, r.Branch, tip, previous); err != nil {
		return MergeResult{}, destinationUpdateFailure(r.Root, r.Branch, previous, err)
	}
	return MergeResult{Kind: kind, PreviousTip: previous, Tip: tip, Tree: tree, Resolved: resolved}, nil
}

// redKinds are the authorization kinds of a gate that ran red on a composed tree. The
// landing and the merge each route a red by this one set.
var redKinds = []authorization.Kind{authorization.Inherited, authorization.Candidate, authorization.LaneFail}

// RedKind reports whether the gate ran red on the composed tree, as opposed to a pass or an
// infrastructure outcome.
func RedKind(kind authorization.Kind) bool { return slices.Contains(redKinds, kind) }

// GradeTarget grades the target tip alone under the owner's authority. A merge reads it on
// its refusal path, to tell a red that the target already holds from a red that the fold
// adds. It publishes nothing, and it writes no output: the fold's own grade already printed
// the failing check.
func (o Owner) GradeTarget(ctx context.Context, r MergeRequest) authorization.Result {
	tree, err := output(r.Root, "rev-parse", r.PreviousTip+"^{tree}")
	if err != nil {
		return authorization.Result{Kind: authorization.Infrastructure, Reason: "merge target tree is unreadable"}
	}
	return o.authorize(ctx, r.Worktree, tree, io.Discard, io.Discard)
}

func mergeTipUnmoved(root, branch, previous string) error {
	tip, err := benchgit.ResolveCommit(root, branch)
	if err != nil || tip != previous {
		return errors.New("merge target tip moved")
	}
	return nil
}
