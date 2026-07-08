# Review pickup — defect-batch-ft43-49 (FT44 slice, commit 9e626cb)

Scope: stories 3–4 only (sweep landed-proof + status action split).

## Standards

3 findings — all judgment calls, 0 hard violations. Worst: #1 (drift risk on a
future edit).

1. **Landed-orphan row derived twice** — `internal/status/status.go`,
   `appendWorktree`: the `plural(..., "orphaned worktree branch", ...)` +
   `"bench worktree clean"` row is built identically in the
   unresolvable-default branch and the landed-subset branch. The two branches
   are semantically distinct and may diverge, so collapsing is a reviewer
   call, not a rule violation.
2. **Advertisement/enforcement coupling** — status's kept-action string
   ("… bench worktree clean keeps them") advertises behavior enforced in
   `internal/worktree/clean.go`. Display copy and sweep note differ and both
   are contract-pinned; collapsing would likely be worse. Noted for awareness.
3. **Test fixture divergence** — `testRuntimeStatusSalvageBranchActions`
   rebuilds the on-main fixture by hand instead of sharing `onMainFixture`
   (different package file). Minor; full reuse isn't clean because the status
   test needs extra branches.

## Spec

0 findings. All four FT44 coverage rows verified red-capable; flagged
defaults honored exactly (cherry all-`-` proof, kept-action text char-for-char,
unresolvable-default guard preserved).

## Coverage

2 findings. Worst: #1 — destructive false positive, **verified by live repro**.

1. **`git cherry` is blind to merge commits, so merge-only content is deleted
   as "landed".** `LandedInDefault` (`internal/worktree/branches.go`) reads
   only `git cherry` lines, and `git cherry` enumerates non-merge commits
   only. A `worktree-*` orphan whose non-merge commits all squash-landed in
   default (all `-`) but whose merge commit carries unique content (conflict
   resolution / evil merge) classifies landed; the sweep force-deletes it and
   the content survives nowhere but the reflog. Verified repro: branch B off
   A; default gains C; branch merges default adding `evil.txt` in the merge
   commit; default then cherry-picks B. `bench worktree clean` deletes the
   branch (`landed by content`) while `evil.txt` is absent from default.
   Plausible agent shape: worktree ran `git merge main` with a resolution,
   feature commits later squash-merged. Fix sketch: treat any merge commit in
   `def..branch` as not-landed (e.g. `git rev-list --merges` probe) before
   trusting the cherry proof. Test attaches beside
   `testRuntimeWorktreeSweepKeepsUniquePatch` in
   `internal/contract/runtime/runtime_worktree_test.go`: keeps-evil-merge
   fixture asserting survival + kept line.
2. **Row-budget interaction untested (minor, low confidence)** — the orphan
   row can now be two sev-2 rows; under the five-row default budget the extra
   row can push a sibling signal behind `+N more`. By-design truncation with
   `--all` escape hatch; no test exercises the split on a full board.
