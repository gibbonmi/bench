# Widen the invalid-mutant known issue to an unused variable

Blocked by: none
Writes: docs/known-issues.md
Covers: none

## What to build

A coordinator reads `docs/known-issues.md` before it writes a delegate charge.
The row for an `invalid` probe mutant names only an unused import.
The commit-merge-parent retro of 2026-10-11 recorded three `invalid` probes from an unused variable.
Two orchestrator swaps and one author omission left a variable unused.
The row names an unused import or an unused variable as the cause of the `invalid` verdict.
The workaround tells the author to keep each import and each variable in use.

## Acceptance

- [ ] The `bench probe` row of `docs/known-issues.md` names a mutant that leaves a variable unused.
- [ ] The workaround of that row tells the author to keep each import and each variable in use.
