package worktree

import "fmt"

func expectedResumeSummary(removed, sweptRefs int, retained string, pruned, reconciled, failed, open int) string {
	return fmt.Sprintf("bench resume: removed %d, swept refs %d%s; pruned branches %d; reconciled %d; failed %d; open assignments %d\n", removed, sweptRefs, retained, pruned, reconciled, failed, open)
}
