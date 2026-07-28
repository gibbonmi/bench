# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree, 11 unpushed commits
Spec: `specs/ft91-artifact-build-tiering.md` (Status: implemented)
Gate: green for all code at `64635b8`; the pin reads stale only because the
spec-status and handoff commits landed after it, both doc-only

## State

- **FT91 artifact build tiering is built, gate-green, and unpushed.** All ten
  stories landed across three commits: the build, a review fix pass closing
  findings on all three axes, and the status flip. Nothing was parked. The
  reviewer owns the push.
- **What it does.** `scripts/build-artifacts.sh` gained
  `BENCH_SHARED_BUILD_CACHE`, matched exactly against `1` so every other value —
  absent, empty, `yes`, `" 1"` — resolves hermetic. Under the opt-in it honors
  the ambient Go build and module caches (resolved before the `HOME` override,
  which is the whole correctness of that story), skips the reproducibility
  second build, and removes a stale `reproducibility.json` as part of the same
  atomic promotion, restoring it if the promotion fails. The two dev contract
  packages set the opt-in in `TestMain`. Byte-reproducibility across independent
  builds stays a ship-tier claim.
- **Measured effect, solo runs against the spec's recorded baselines:** the
  artifact suite 133.5 s → 106 s, the surface suite 115.5 s → 12 s. That is
  249 s → 118 s, and it already absorbs ~50 s of new hermetic-posture tests the
  spec's polarity rows require.
- **The gate wall-clock did not drop by that much — 4m51s → 4m27s — and nobody
  has diagnosed why.** Under gate parallelism the artifact suite inflates to
  ~152 s, so the contract phase is no longer clearly the critical path. The
  spec's motive was gate time, so this gap is worth a look before anyone calls
  FT91 done on its own terms. The remaining artifact-suite lever is already
  diagnosed and parked in `IDEAS.md`: the residual is invocation count, not
  per-invocation cost.
- **Four post-approval spec corrections await veto.** The semantic review found
  four statements in the spec that were factually wrong about the code — story
  4's rollback row described a seam that parks above the promotion backups,
  story 9 over-stated a deletion, and two edge-inventory exclusions rested on
  reasons the opt-in invalidated. Each is marked `**Post-approval correction,
  flagged:**` in the spec. The code is right in all four cases; only the prose
  changed. Whether that edit was mine to make is the open question captured in
  `.bench/learnings.md`.
- **One open learning, unverdicted:** whether a review's Spec-axis finding may
  correct an approved spec in place, or must stop for sign-off. It proposes two
  candidate rules and asks the reviewer to pick.
- **`ft91-gate-phase-split` stays unretired on purpose,** so `bench status`
  keeps reporting one spec awaiting retirement until the reviewer rules on its
  stories 4, 5, and 9.

## Next command

Push, which is the reviewer's to perform. Then `/bench-what-next` in a fresh
mid-tier session — the drain has one parked idea and one open learning, and
`/bench-what-next` is the only path either takes into the roadmap.

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
