package backgroundstate

// GradePaths answers the exact paths that a detached grade writes for the worktree git directory
// gitdir and the git common directory common: the background record, the log, the running log, the
// stop marker, and the claim directory. A foreground gate reads them as the state of another run.
func GradePaths(gitdir, common string) []string {
	return []string{RecordPath(gitdir), LogPath(gitdir), RunningLogPath(gitdir), stopPath(gitdir), claimsPath(common)}
}
