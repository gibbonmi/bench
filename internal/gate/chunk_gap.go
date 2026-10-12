package gate

import (
	"fmt"
	"io"

	commitrepo "github.com/gibbonmi/bench/internal/commitment/repository"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/reviewrecord"
)

// openChunkGap opens the build's gap after the green checkpoint of a chunk. A failed write
// exits 1, and a rerun reuses the green verdict and writes the record.
func openChunkGap(root string, checkpoint Checkpoint, stderr io.Writer) int {
	if err := writeChunkGap(root, checkpoint); err != nil {
		fmt.Fprintln(stderr, "chunk gap not recorded: "+err.Error())
		return 1
	}
	return 0
}

// writeChunkGap reads the plan at HEAD, so the gap follows the plan that the checkpoint graded.
func writeChunkGap(root string, checkpoint Checkpoint) error {
	head, err := benchgit.ResolveCommit(root, "HEAD")
	if err != nil {
		return err
	}
	_, plan, err := reviewrecord.ReadSpecAtCommit(root, checkpoint.Spec, head)
	if err != nil {
		return err
	}
	return commitrepo.Store{Root: root}.OpenGap(plan, checkpoint.Chunk)
}
