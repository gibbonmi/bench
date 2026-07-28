# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0dd9c83`, clean tree, 6 unpushed commits
Spec: `specs/ft91-artifact-build-tiering.md` (Status: staged)
Gate: green at `85e757e` — stale, work tree `fa03401`

## State

- **`specs/ft91-artifact-build-tiering.md` is staged and reviewer-approved.**
  Ten stories, three seams, sixteen coverage rows, `bench coverage --check`
  green. Nothing of it is built yet. Six commits are unpushed; the reviewer owns
  the push.
- **It re-tiers release-build hermeticity out of the dev gate.**
  `scripts/build-artifacts.sh` gains one posture switch whose *absence* is
  hermetic — that polarity is the whole safety argument, and inverting it would
  let a release build inherit dev posture silently. Under the opt-in the script
  honors ambient Go caches and skips the reproducibility second build; the
  two-build proof stays on the ship tier, where `release-evidence-probe`
  (`Tier: Ship`) already validates it.
- **Four decisions in the spec are the author's, not the map's, and are flagged
  in it for veto:** the token name `BENCH_SHARED_BUILD_CACHE=1`, using the
  `internal/preprelease` unit seam instead of the conformance check the Handoff
  named, *deleting* a stale `reproducibility.json` (the spec's only destructive
  file operation), and setting the opt-in in `TestMain` per package.
- **A top-tier falsification pass returned BLOCK on the first draft and every
  finding verified against the tree.** The revision replaced three coverage rows
  whose red signal could not fire, a citation that matched no test and exited 0,
  and a rollback injection that fired before the block it claimed to test. The
  `GOPROXY=off` lever is now the discriminator for four separate failures. Do not
  re-loosen those rows.
- **Measured baselines, for judging whether the build worked:** artifact suite
  133.5 s and surface suite 115.5 s before; a cold-cache build is 4.79 s against
  0.20 s warm. The artifact suite fell to 73.2 s with a shared cache in a
  throwaway probe.
- **The ~73 s residual is already diagnosed and parked in `IDEAS.md`** — it is
  invocation count, not packing cost. Out of scope for this spec by reviewer
  decision, because which tests may share an artifact set is a test-isolation
  call.
- **One blocker for the build session.** `/bench-implement-spec` requires
  assigning write work to a subagent before the first edit, and this repo's
  recent sessions carry a standing instruction forbidding the Agent tool unless
  the reviewer asks. Get a yes or a no on delegates before starting. The
  unresolved rule conflict is `.bench/learnings.md`'s territory and is already
  queued as FT107's eighth clause.
- **`ft91-gate-phase-split` stays unretired on purpose,** so `bench status` keeps
  reporting one spec awaiting retirement and one roadmap row for merged work
  until the reviewer rules on its stories 4, 5, and 9.

## Next command

`/bench-implement-spec specs/ft91-artifact-build-tiering.md` — in a fresh
mid-tier session, per the spec's approval. Not `bench shift`: the coverage map is
not fully gate-observable (story 8 is honestly marked not-TDD-able, story 9's
proof is ship-tier-only), so it fails `craft-line`'s venue-routing test for an
unattended run.

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
