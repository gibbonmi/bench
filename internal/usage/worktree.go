package usage

import "strings"

const (
	WorktreeList        = "bench worktree list"
	WorktreePath        = "bench worktree path <target>"
	WorktreeExec        = "bench worktree exec <target> -- <command> [args...]"
	WorktreeCreate      = "bench worktree create [--refresh] --request <opaque-id> --label <work-item>"
	WorktreeRelease     = "bench worktree release --request <opaque-id> <path>"
	WorktreeClean       = "bench worktree clean [--discard-ignored] [--discard-branch] [--full] (<path> | --landed) [--apply <fingerprint>]"
	WorktreeReauthorize = "bench worktree reauthorize --assignment <assignment-id> --request <opaque-id> --base <commit> --source-tip <commit> <path>"
	WorktreeLand        = "bench worktree land --request <opaque-id> --base <commit> --source-tip <commit> --spec <slug> -m <message> <path>"
	WorktreeLandResume  = "bench worktree land --resume <published-commit> --request <opaque-id> --base <commit> --source-tip <commit> --spec <slug> <path>"
)

var worktreeCommands = []string{
	WorktreeList,
	WorktreePath,
	WorktreeExec,
	WorktreeCreate,
	WorktreeRelease,
	WorktreeClean,
	WorktreeReauthorize,
	WorktreeLand,
	WorktreeLandResume,
}

func WorktreeUsage() string {
	return "usage: bench worktree [--refresh] [objective...]\n       " + strings.Join(worktreeCommands, "\n       ") + "\n       bash bin/bench.sh gate --fresh\n"
}
