package gate

import (
	"context"
	"io"

	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/landing/published"
)

// cleanCheckoutRefusal names the recovery for a complete checkpoint on a dirty checkout.
const cleanCheckoutRefusal = "complete checkpoint requires a clean checkout: commit or remove each uncommitted change"

// executeCompleteCheckpoint grades the tree that the landing publishes from HEAD of root.
// The landing publishes only committed bytes, so an uncommitted change refuses before the
// oracle runs, and the review record has no exemption. The landing's closure transform
// applies to the HEAD tree, and the oracle grades the result under the completion
// obligation for spec at HEAD.
func executeCompleteCheckpoint(ctx context.Context, root, spec string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	tip, err := benchgit.ResolveCommit(root, "HEAD")
	if err != nil {
		return operational(root, 0, stderr, "complete checkpoint source unavailable: "+err.Error())
	}
	sourceTree, err := benchgit.Output("-C", root, "rev-parse", "--verify", tip+"^{tree}")
	if err != nil {
		return operational(root, 0, stderr, "complete checkpoint source unavailable: "+err.Error())
	}
	working, err := workingTreeSource{root: root}.tree()
	if err != nil {
		return operational(root, 0, stderr, "complete checkpoint source unavailable: "+err.Error())
	}
	if working != sourceTree {
		return operational(root, 0, stderr, cleanCheckoutRefusal)
	}
	tree, err := published.Tree(root, sourceTree, spec, tip)
	if err != nil {
		return operational(root, 0, stderr, "complete checkpoint unavailable: "+err.Error())
	}
	return executeTreeWithOwner(WithCompletion(ctx, spec, tip), root, tree, stdout, stderr, arm, mode, nil)
}
