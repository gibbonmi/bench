package usage

import "strings"

const (
	WorktreeList     = "bench worktree list"
	WorktreeCreate   = "bench worktree create --request <opaque-id> --label <work-item>"
	WorktreeRelease  = "bench worktree release --request <opaque-id> <path>"
	WorktreeClean    = "bench worktree clean [--discard-ignored] [--full] <path> [--apply <fingerprint>]"
	WorktreeRecovery = "bench worktree recovery <ref> [--apply <fingerprint>]"
)

var worktreeCommands = []string{
	WorktreeList,
	WorktreeCreate,
	WorktreeRelease,
	WorktreeClean,
	WorktreeRecovery,
}

func WorktreeUsage() string {
	return "usage: bench worktree\n       " + strings.Join(worktreeCommands, "\n       ") + "\n"
}
