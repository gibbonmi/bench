# Extend the root guards to per-component scoping

Blocked by: Skip the build phase on its attested seal
Ownership fence: the root-guard cases in `internal/gate/reduced_run_test.go`
Assumptions: `[RB1]` already pins that a root which is not the kit runs
unreduced, and `[R20]` that an allowlist-confined change with no ancestor runs
the full gate; both are whole-changeset guards written before the partition
existed. Re-derive from the tree at pickup.

## What to build

The existing root guards prove the whole-changeset reduction never fires outside
the kit. The partition is a second way to skip work, and it needs the same
proof: a linked repository gating through the kit's binary, and a root with no
Go module, must never scope any component. The declarations are the kit's own —
matched by directory identity, not by spelling, so a symlinked kit path still
counts and any stat failure runs everything.

This ticket extends the guards rather than replacing them: both original cases
keep their coverage, and each gains its per-component half.

## Acceptance

- [ ] PC18a — a root that is not the kit executes every component, with no slot authored and no skip announced.
- [ ] PC18b — a root with no Go module executes every component its table carries and scopes none.
- [ ] PC18c — a symlinked path to the kit root still counts as the kit, and a stat failure on either path runs every component.
- [ ] PS36 — the existing `[RB1]` and `[R20]` cases still pass unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC18a | compare root and kit by string equality instead of file identity, then drop the guard from the per-component site | `TestForeignRootScopesNoComponent` | point `BENCH_KIT` at a kit fixture, gate a different root, assert every phase marker present and the store empty |
| PC18b | resolve component identities before checking for a Go module | `TestNoGoRootScopesNoComponent` | gate a fixture with no `go.mod`, assert every table phase executed |
| PC18c | treat a stat failure as a match | `TestSymlinkedKitCountsAndStatFailureRunsAll` | gate through a symlink to the kit fixture, then through a path removed mid-run, assert scoping in the first and none in the second |
| PS36 | fold the whole-changeset guard into the per-component one | existing `[RB1]`/`[R20]` cases | run them unchanged |
