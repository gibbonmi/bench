package landing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gibbonmi/bench/internal/gate/authorization"
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
	if got := o.authorize(ctx, r.Root, tree, r.Stdout, r.Stderr); !o.publishes.permits(got.Kind) {
		return MergeResult{}, fmt.Errorf("prospective authorization refused: %s", got.Kind)
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

func mergeTipUnmoved(root, branch, previous string) error {
	tip, err := output(root, "rev-parse", "--verify", branch+"^{commit}")
	if err != nil || tip != previous {
		return errors.New("merge target tip moved")
	}
	return nil
}
