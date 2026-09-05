# Deny raw git composition in the destructive-git guard

Blocked by: none
Writes: internal/gitguard/gitguard.go, internal/gitguard/verdict.go, internal/gitguard/verdict_test.go, .bench/BENCH-reference.md
Covers: none

## What to build

The destructive-git guard denies `git merge` in every form and `git cherry-pick` in
every form. Bench owns composition: `bench worktree merge` folds a sibling, and
`bench worktree land` lands a source. A raw fold leaves a single-parent commit that
ancestry never proves landed. The sibling branch then survives as an unclaimed ref, and
`bench status` routes it to a destructive discard.

The deny table gains two rows, `merge` and `cherry-pick`. The `merge` row carries the
advice sentence that names the two Bench verbs, because the agent can type that fix.
The `cherry-pick` row carries the same advice. The classifier returns the row's label
for the bare verb and for every option form, including `--squash`, `--no-commit`,
`--abort`, and `--continue`. The guard manifest advertises the two new classes through
the existing table, with no second list.

The reference's hook-layer paragraph states the new posture in one sentence beside the
push sentence.

## Acceptance

- [ ] `bench guard-git` denies `git merge <ref>`, `git merge --squash <ref>`, `git merge --no-commit <ref>`, and `git merge --abort` with the label `git merge`.
- [ ] `bench guard-git` denies `git cherry-pick <ref>`, `git cherry-pick --no-commit <range>`, and `git cherry-pick --continue` with the label `git cherry-pick`.
- [ ] The block message for both labels ends with one advice sentence that names `bench worktree merge` and `bench worktree land`.
- [ ] `git revert <ref>`, `git commit -m x`, and `git log --merges` stay allowed.
- [ ] The guard manifest lists the two new classes, derived from the deny table.
- [ ] `.bench/BENCH-reference.md` states that the guard denies a raw merge and a raw cherry-pick, and names the Bench verbs that compose instead.
