# Deny a raw git commit and a raw git reset

Blocked by: none
Writes: internal/gitguard/gitguard.go, internal/gitguard/verdict.go, internal/gitguard/verdict_test.go, internal/gitguard/gitguard_test.go, internal/gitguard/scan.go, internal/gitguard/scan_test.go, internal/worktree/merge.go, internal/worktree/merge_test.go, .bench/hooks/block-dangerous-git.sh, .bench/BENCH-reference.md
Covers: none

## What to build

A ticket author ran a raw `git commit` and a `git reset --soft` in its pool worktree. These
calls went around the `bench commit` lane. The destructive-git guard allowed each call, because
it denied only `git commit --amend` and `git reset --hard`.

The guard denies `git commit` and `git reset` in every option form. One `commit` class
replaces the `amend` class. The `reset` class takes the label `git reset`. Each refusal names
the Bench route. The commit refusal names `bench commit -m <msg> -- <path>...`. The reset
refusal reads the grammar line of `bench worktree reset`, so the route cannot drift.

A pool worktree runs git through `bench worktree exec <target> -- <argv>`. The scanner grades
that child argv one level deep, the same depth as a `bash -c` string. So every deny class
applies to the exec child. The child runs in another checkout, so a push through exec fails
closed as an unresolved destination.

The guard continues to allow an index write such as `git add` or `git restore --staged`. An
index write discards no work and moves no ref. A Bench commit composes through its own index.
The shim header and the hook-layer note in the reference state the new deny surface.

The exit-3 record of `bench worktree merge` named a raw `git reset --merge` as its repair.
The guard now denies that command for an agent. The record names the reset plan of
`bench worktree reset --to <tip> <assignment>` instead, and that verb reconciles an unfinished
merge checkout.

## Acceptance

- [ ] The verdict rows for `git commit -m x` and `git reset --soft HEAD~1` fail on the base guard, and pass after the guard edit.
- [ ] The verdict rows for `bench worktree exec X -- git commit -m x`, `-- git reset --soft HEAD~1`, and `-- git reset --hard` fail on the base scanner, and pass after the scanner edit.
- [ ] The exec rows for `git status`, `git log`, `git diff`, `git rev-parse HEAD`, `git apply --index`, and `git format-patch` stay allowed.
- [ ] A probe that omits the commit deny turns the `commit -m` row red, and the restore returns it to green.
- [ ] A probe that omits the exec recursion turns the exec deny rows red, and the restore returns them to green.
- [ ] The commit refusal names `bench commit -m <msg> -- <path>...`, and the reset refusal names the grammar line of `bench worktree reset`.
- [ ] The verdict rows for `git add -A` and `git restore --staged .` stay allowed.
- [ ] The exit-3 merge record names the reset plan at the published commit.
- [ ] The verdicts for every other direct git verb stay unchanged.
