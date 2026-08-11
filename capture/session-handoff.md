# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `4b96a7b`, clean tree, 2 unpushed commits
Spec: `specs/axi-coherent-diff/spec.md` (Status: staged), `specs/axi-query-disclosure/spec.md` (Status: staged), `specs/bench-preflight/spec.md` (Status: staged), `specs/single-build-serial-gate/spec.md` (Status: staged)
Gate: green at `a56d225` — stale, work tree `3839fa5`

## State

`/bench-implement-spec --full bench-preflight` phase reached: **tickets staged,
breakdown reviewed, implementation dispatch next**. The 7-ticket breakdown
under `specs/bench-preflight/tickets/` passed the `craft-tickets` breakdown
review (Codex `gpt-5.6-sol`/high: 12 findings, all repaired; verify round
confirmed, two residuals closed after it). Reviewer approvals recorded
2026-08-11: spec fences gained `.bench/structure-accept` and
`tests/canary/package-core-guard/`; map rows PF24/PF25 added; the rename sweep
exempts `specs/`, `CHANGELOG.md`, `capture/`, `ROADMAP.md`, `decisions/`.
Line: ticket write-delegates sonnet (reviewer-directed sonnet-over-opus),
advertise ticket fable/high (leverage override), coordinator fable/medium.

Build order: frontier `rename-release-preflight-package.md`,
`export-diff-review-base.md`, `repair-routing-checker.md` (fence-disjoint,
parallel dispatch, serial landings) → `implement-preflight-review.md` →
`harden-preflight-bootstrap-errors.md` → `implement-preflight-build.md` →
`advertise-preflight-kit-prose.md`. Landings via path-scoped `bench commit`;
`--spec bench-preflight` rides only the final landing. Parked pre-reshape
specs (`axi-coherent-diff`, `axi-query-disclosure`, `single-build-serial-gate`)
await re-rank, not active. Post-spec queue unchanged: Spec C (doctrine
adoption), then a `/bench-what-next` drain.

## Next command

`/bench-implement-spec --full bench-preflight`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
