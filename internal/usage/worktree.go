package usage

import "strings"

const (
	WorktreeList     = "bench worktree list"
	WorktreePath     = "bench worktree path <target>"
	WorktreeExec     = "bench worktree exec <target> -- <command> [args...]"
	WorktreeCreate   = "bench worktree create [--refresh] --request <opaque-id> --label <work-item>"
	WorktreeRelease  = "bench worktree release --request <opaque-id> <path>"
	WorktreeClean    = "bench worktree clean [--discard-ignored] [--discard-branch] [--full] <path> [--apply <fingerprint>]"
	WorktreeRecovery = "bench worktree recovery <ref> [--apply <fingerprint>] [--discard <fingerprint>]"
)

var worktreeCommands = []string{
	WorktreeList,
	WorktreePath,
	WorktreeExec,
	WorktreeCreate,
	WorktreeRelease,
	WorktreeClean,
	WorktreeRecovery,
}

func WorktreeUsage() string {
	return "usage: bench worktree [--refresh] [objective...]\n       " + strings.Join(worktreeCommands, "\n       ") + "\n       bash bin/bench.sh gate --fresh\n"
}
