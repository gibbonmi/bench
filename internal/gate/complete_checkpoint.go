package gate

import (
	"context"
	"io"

	"github.com/gibbonmi/bench/internal/landing/published"
)

// executeCompleteCheckpoint grades the tree that the landing publishes from HEAD of root.
// The landing publishes only committed bytes, so committedSource refuses an uncommitted
// change before the oracle runs. The landing's closure transform applies to the HEAD tree,
// and the oracle grades the result under the completion obligation for spec at HEAD. Each
// refusal here names the route of its cause.
func executeCompleteCheckpoint(ctx context.Context, root, spec string, stdout, stderr io.Writer, arm postAcquireContextArm, mode runMode) Result {
	tip, sourceTree, refusal, refused := committedSource(ctx, root, stderr, mode)
	if refused {
		return refusal
	}
	tree, err := published.Tree(root, sourceTree, spec, tip)
	if err != nil {
		return refuse(ctx, root, stderr, mode, faceComposition, "complete checkpoint unavailable: "+err.Error())
	}
	return executeTreeWithOwner(WithCompletion(ctx, spec, tip), root, tree, stdout, stderr, arm, mode, nil)
}
