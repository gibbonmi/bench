# Exit a published commit chain with its own code

Blocked by: none
Writes: cmd/bench/commit_chain.go, cmd/bench/commit_chain_test.go, CHANGELOG.md
Covers: none

## What to build

`bench commit --help` states that exit 1 means the commit was refused before
publication and nothing was committed. Before this ticket, the
`--preflight-build` chain also exited 1 when the commit published and the build
or the build preflight then failed. That exit told the caller that nothing was
committed, which is false.

After a commit that published, a failed build step or a failed preflight step
exits with one named code that no other `bench commit` exit uses. The chain
help line states that code, and the exit and the help line read the same named
constant. The commit-chain line still names the published commit. A commit that
refuses keeps exit 1 and runs no later step.

## Acceptance

- [ ] A commit that publishes, followed by a preflight that exits 1, makes the chain exit with the new code, not 1.
- [ ] In that case the commit-chain line names the published commit.
- [ ] A commit that publishes, followed by a build that exits non-zero, makes the chain exit with the new code, and no preflight runs.
- [ ] A commit that refuses exits 1, and no build and no preflight run.
- [ ] A commit that publishes and exits 3 keeps exit 3 and runs no later step.
- [ ] The chain help line states the new code and the meaning "published; a later step failed".
- [ ] A green build and a green preflight exit 0.
- [ ] `CHANGELOG.md` records the changed exit.
