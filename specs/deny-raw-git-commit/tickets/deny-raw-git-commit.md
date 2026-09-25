# Deny a raw git commit and a raw git reset

Blocked by: none
Writes: internal/gitguard/gitguard.go, internal/gitguard/verdict.go, internal/gitguard/verdict_test.go, internal/gitguard/gitguard_test.go, internal/gitguard/scan_test.go, .bench/hooks/block-dangerous-git.sh, .bench/BENCH-reference.md
Covers: none

## What to build

A ticket author ran a raw `git commit` and a `git reset --soft` in its pool worktree. These
calls went around the `bench commit` lane. The destructive-git guard allowed each call, because
it denied only `git commit --amend` and `git reset --hard`.

The guard denies `git commit` and `git reset` in every option form. One `commit` class
replaces the `amend` class. The `reset` class takes the label `git reset`. Each refusal names
the Bench route. The commit refusal names `bench commit -m <msg> -- <path>...`. The reset
refusal reads the grammar line of `bench worktree reset`, so the route cannot drift.

The guard continues to allow an index write such as `git add` or `git restore --staged`. An
index write discards no work and moves no ref. A Bench commit composes through its own index.
The shim header and the hook-layer note in the reference state the new deny surface.

## Acceptance

- [ ] The verdict rows for `git commit -m x` and `git reset --soft HEAD~1` fail on the base guard, and pass after the guard edit.
- [ ] A probe that omits the commit deny turns the `commit -m` row red, and the restore returns it to green.
- [ ] The commit refusal names `bench commit -m <msg> -- <path>...`, and the reset refusal names the grammar line of `bench worktree reset`.
- [ ] The verdict rows for `git add -A` and `git restore --staged .` stay allowed.
- [ ] The verdicts for every other git verb stay unchanged.
