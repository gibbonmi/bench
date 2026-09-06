# Retire shift residue branches through the unclaimed cleanup

Blocked by: none
Writes: internal/intent/ledger/validate.go, internal/intent/ledger_aliases.go, internal/shift/loop.go, internal/worktree/clean_unclaimed.go, internal/worktree/clean_unclaimed_test.go
Covers: none

## What to build

`bench shift` creates a `bench/shift-<timestamp>` branch and never deletes it. When the
shift's worktree is released, the branch stays as an ordinary unique ref. No Bench verb
retires it unless the content proof succeeds, so the reviewer deletes it by hand.

The unclaimed cleanup, `bench worktree clean --discard-branch --unclaimed`, extends its
selection to the shift namespace. Bench created these refs, so Bench retires them under
the same plan-before-apply fingerprint. The planner keeps every protection it has: a
recorded assignment branch, a checked-out branch, and the default branch stay out of
the set. The plan row for a shift branch carries the reason `shift residue branch`. The
plan row for an assignment branch keeps its reason.

The shift namespace has one source. A `ShiftBranchPrefix` function in the ledger leaf
returns `refs/heads/bench/shift-`. The shift loop derives its branch name from it, and
the planner reads it. The bare literal leaves `internal/shift/loop.go`.

## Acceptance

- [ ] A released shift's `bench/shift-<timestamp>` branch appears in the `--unclaimed` plan with the reason `shift residue branch`, and `--apply <fingerprint>` deletes it.
- [ ] A shift branch that a worktree has checked out stays out of the plan.
- [ ] An assignment branch in the plan keeps the reason `unclaimed assignment branch`.
- [ ] A branch outside both namespaces, for example `archive/x`, stays out of the plan.
- [ ] The shift loop still creates its branch under `bench/shift-`, derived from the one prefix source.
