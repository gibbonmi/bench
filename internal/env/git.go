package env

// GitTestConfig returns the environment entries that turn off git auto-maintenance for
// a run of the kit's tests. From git 2.47, a commit starts that maintenance in the
// background. Its worktree prune can then remove a worktree admin entry a test fixture
// planted, while the test still reads it.
//
// The entries set GIT_CONFIG_COUNT, so they replace any configuration the caller's
// environment passes to git the same way. The gate's kit phases, `bench test`, and the
// release preflight phases carry them. A plain `go test` run does not.
func GitTestConfig() []string {
	return []string{
		"GIT_CONFIG_COUNT=1",
		"GIT_CONFIG_KEY_0=maintenance.auto",
		"GIT_CONFIG_VALUE_0=false",
	}
}
