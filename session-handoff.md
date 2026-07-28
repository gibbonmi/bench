# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree after this commit, ahead of origin
Spec: `specs/ft91-canary-contract-scoping.md` — Status: staged, reviewer-approved
2026-07-28 (approval was conditional on clean review findings; both passes ran
and their surviving findings are folded in)
Gate: green, measured 2026-07-28 at 267 s wall

## State

- **Stage 1 of FT91's eighth arm is staged and approved.** The spec scopes
  each `behavior-owned` canary fixture's nested gate to the one contract
  package owning its EXPECT via subfamily directories
  (`tests/canary/behavior-owned/<package path>/<fixture>/`), with
  `BENCH_CANARY_CONTRACT_PACKAGE` narrowing the inner contract argv, loud
  reds for every failure mode, per-group scoped vacuity baselines, a
  kit-only guard pinning the flat set to the two named relocations, and a
  ≤100 s solo-canary acceptance threshold on the ship evidence. Build
  inputs: the spec, `decisions/gate-critical-path.md` (Handoff),
  `decisions/assets/gate-critical-path-timeline.md`, and
  `decisions/assets/behavior-owned-package-bindings.md` (traced fixture →
  package inventory; two fixtures relocate to legacy flat).
- **Two falsification passes ran:** mid-tier (13 findings, all verified and
  folded — manifest bypass, baseline key desync, parent-basename call sites,
  strays) and a reviewer-directed Codex `gpt-5.6-sol` xhigh pass (2 real
  findings — degenerate all-flat migration, single-baseline degenerate —
  both closed with new red-capable rows; verdicts were advisory BLOCKs on
  pre-fix drafts, and every surviving finding is addressed in the committed
  spec).
- **Map state:** `gate-critical-path` tickets #2 (artifact-suite
  build-vs-inspect research) and #3 (prepared-artifact sharing ruling,
  blocked by #2) are the only open tickets; they gate the artifact-hoist
  slice, not this build. FT91's stop condition is a measured dev gate ≤60 s.
- **Roadmap's FT91 row is stale** on the phase-split spec (retired this
  session) and this new staging; the next `/bench-what-next` reconcile owns
  the rewrite. The tree wins meanwhile.

## Next command

`/bench-implement-spec specs/ft91-canary-contract-scoping.md` on a fresh
mid-tier session — the spec is approved, seams are pre-agreed, and stories
route `gpt-5.6-terra` (alias `opus` for Claude Code delegation) per the
cached Lines.

After it lands green: stage 2 (`/bench-write-spec` from map #7), and
`/bench-shape-idea` resume on `gate-critical-path` #2 for the artifact-hoist
slice.

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
