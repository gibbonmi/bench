# Worktree Branch Sweep (FT28)

## #1: Which branches may `bench worktree clean` delete?

Blocked by: —
Type: Grill

### Question
8 orphaned `worktree-agent-*` branches sit in the repo; the status row says
"delete scratch branch" but no bench subcommand deletes branches and the
block-dangerous-git hook denies `git branch -D` to the agent — the recommended
remedy is a dead end. What is safe for the CLI to own?

### Answer
`bench worktree clean` deletes a branch when all three hold: it matches the
delegate scratch naming (`worktree-agent-*`), it is not checked out in any
live worktree, and its tip is fully merged into the default branch (no unique
commits). Unmerged scratch branches are listed by name with "unique commits —
inspect or delete by hand" so nothing with unreviewed work disappears. The
deletion is performed by the CLI itself (the sanctioned owner); the agent-
facing git guard stays exactly as-is. The status action string becomes
`bench worktree clean` — invocable. **Flagged for veto:** the conservative
default (merged-only) means genuinely abandoned unmerged scratch branches
still need the reviewer's hand; the alternative (delete all scratch-named
branches) was rejected because naming alone is not proof the work landed.

## Handoff

1. **Module boundaries.** `internal/worktree` owns detection
   (OrphanedDelegateBranches exists) and gains the sweep; `internal/status`
   only updates the action string; `bin/bench.sh` routing unchanged
   (`worktree clean` exists).
2. **Contracts.** `bench worktree clean` output gains: `deleted branch <name>`
   per merged orphan, `kept branch <name> (unique commits — inspect or delete
   by hand)` per unmerged one; exit 0 in both cases (keeping is not an error).
3. **Deep vs thin.** The merged-into-default check is the deep logic
   (ancestry, not string compare); deletion is a thin `git branch -d`-
   equivalent once safety holds.
4. **Black-box assertables.** In a temp repo: a merged scratch branch
   disappears after clean; an unmerged one survives with the kept line; a
   checked-out one survives; a non-scratch-named branch is never touched.
5. **Gate attachment.** Runtime contract tests against the built binary (the
   existing worktree contract family).
6. **Hostile-input owners.** worktree owns branch names with slashes or
   unicode within the scratch prefix, detached-HEAD default branch, and a
   missing default branch (git.DefaultBranch failure → report, delete
   nothing — the false-empty rule from [FT29] applies here from day one).
7. **Uncertainty flags.** n/a — the seam is settled; detection code exists.
8. **Rejected alternatives.** Deleting unmerged scratch branches; loosening
   the git guard so the agent can delete branches; a separate `bench branch`
   subcommand.
9. **Domain watch-outs.** A branch checked out in a pool worktree cannot be
   deleted by git at all — the sweep must classify before attempting, not
   swallow the git error after.

Dependency order: n/a — single spec.
