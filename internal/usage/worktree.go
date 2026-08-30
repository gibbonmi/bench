package usage

import (
	"strings"

	"github.com/gibbonmi/bench/internal/toon"
)

const (
	WorktreeList        = "bench worktree list"
	WorktreePath        = "bench worktree path <target>"
	WorktreeExec        = "bench worktree exec <target> [--env KEY=VALUE]... -- <command> [args...]"
	WorktreeShow        = "bench worktree show <target> <rev>:<path>"
	WorktreeBuild       = "bench worktree build <target>"
	WorktreeCreate      = "bench worktree create [--refresh] --request <opaque-id> --label <work-item> [--from <target>]"
	WorktreeRelease     = "bench worktree release --request <opaque-id> <path>"
	WorktreeClean       = "bench worktree clean [--discard-ignored] [--discard-branch] [--full] (<path> | --landed) [--apply <fingerprint>] | bench worktree clean --discard-branch --unclaimed [--apply <fingerprint> | --apply-current]"
	WorktreeReclaim     = "bench worktree reclaim [--apply <fingerprint>]"
	WorktreeReauthorize = "bench worktree reauthorize --assignment <assignment-id> --request <opaque-id> --base <commit> --source-tip <commit> <path>"
	WorktreeMerge       = "bench worktree merge --from <commit|target> <target>"
	WorktreeLand        = "bench worktree land --request <opaque-id> --base <commit> --source-tip <commit> [--spec <slug>] -m <message> <path>"
	WorktreeLandResume  = "bench worktree land --resume <published-commit> --request <opaque-id> --base <commit> --source-tip <commit> [--spec <slug>] <path>"
)

// WorktreeExecGate is the one gate form a worktree reader gets. It is an exec composition
// rather than a worktree route, so it trails the grammar list instead of joining it, and
// both this usage trailer and the `bench help` row read it from here.
const WorktreeExecGate = "bench worktree exec <target> -- bench gate"

var worktreeCommands = []string{
	WorktreeList,
	WorktreePath,
	WorktreeExec,
	WorktreeShow,
	WorktreeBuild,
	WorktreeCreate,
	WorktreeRelease,
	WorktreeClean,
	WorktreeReclaim,
	WorktreeReauthorize,
	WorktreeMerge,
	WorktreeLand,
	WorktreeLandResume,
}

func WorktreeUsage() string {
	return "usage: bench worktree [--refresh] [objective...]\n       " + strings.Join(worktreeCommands, "\n       ") + "\n       " + WorktreeExecGate + "\n"
}

// PrimaryCheckoutRefusal is the one refusal a Bench write verb prints from the primary
// checkout. Within Bench, main receives writes only through landings, so every verb
// that writes the tree redirects to a worktree with this line.
func PrimaryCheckoutRefusal() string {
	return toon.Errorf("primary checkout is read-only for Bench phases", "run "+WorktreeCreate)
}
