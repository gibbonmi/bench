# Name the expected commit in three refusals

Blocked by: none
Writes: internal/reviewrecord/coverage.go, internal/gate/review_checkpoint_commits_test.go (new), internal/worktree/land_identity.go, internal/worktree/land_refusal.go, internal/worktree/land_rerun.go (new), internal/worktree/land_folded_base_test.go (new)
Covers: none

## What to build

This ticket is the FT316 fix. Three refusals held the commit that the caller needed and
did not print it. The caller then found each commit by hand.

1. The checkpoint chain-gap refusal names the expected base, which is the predecessor chunk tip. It also names the base that it got, and it states the chain rule.
2. The checkpoint uncovered-source refusal names the range that no chunk review covers. When the reviewed tip is the requested tip, it names the uncommitted work past that commit.
3. The landing fence refusal names the folded merge base when the source merged a later default-branch commit after `--base`. The refusal carries the caller's base as `observed` and the merge base as `wanted`. Its route re-points `--base` at the merge base.

The default-branch fact comes from `git.ResolvedDefault`, and the ancestry fact comes from
`authorization.IsAncestor`. A new `landingBaseFlag` is the one rendering of the re-run's
`--base` argument, and the route swaps that exact text. The re-run rendering moves to its
own file, so `land_refusal.go` does not grow past its line budget.

## Acceptance

- [ ] `bench gate --checkpoint <spec> --chunk 2` over a chain gap prints `expected base <predecessor tip>` and the chain rule.
- [ ] `bench gate --checkpoint <spec> --chunk 1` over an unreviewed repair commit prints `<reviewed tip>..<source tip>`.
- [ ] `bench worktree land` over a source that folded a later default-branch commit prints `wanted=<merge base>`, and its `next=` route names `--base <merge base>`.
- [ ] If `foldedDefaultBase` answers empty, the landing test goes red.
- [ ] The existing checkpoint, fence, and landing registry tests pass unchanged.
