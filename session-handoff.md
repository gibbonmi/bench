# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree, three commits ahead of origin
Spec: none active — `specs/ft91-gate-phase-split.md` (Status: implemented) is
held unretired on purpose, see below
Gate: green for all code at `64635b8`; the pin reads stale only because
doc-only commits landed after it

## State

- **The drain closed 2026-07-28, approved as one batch diff.** Both inboxes are
  empty, the roadmap parses clean, and `bench roadmap` extracts the sequence.
- **FT91 is retargeted on the seventh arm.** `ft91-artifact-build-tiering`
  shipped and retired: a `BENCH_SHARED_BUILD_CACHE` opt-in, dev-tier only, lets
  the contract suites honor the ambient Go caches and skip the reproducibility
  second build. Measured solo — artifact suite 133.5 s → 106 s, surface suite
  115.5 s → 12 s. **The whole gate absorbed only 24 s of that 131 s** (4m51s →
  4m27s), because under gate parallelism the artifact suite inflates to ~152 s.
  Nobody has diagnosed the gap, and it is the first question the next arm
  answers.
- **The remaining artifact lever is diagnosed but is a decision, not a build.**
  The cost is invocation count — ~20 host-only generator runs at ~3.7 s each,
  not packing. The `BENCH_TEST_PREPARED_ARTIFACTS` seam exists but is scoped
  per-test. Hoisting it to package scope needs a ruling on which tests may share
  one artifact set without losing independence.
- **The spec-edit learning became FT144's third instance.** The question
  generalized: when a phase's finding lands on an *approved* spec rather than on
  the code, what may it do under batch approval? FT144 already carried the
  build-side face; the review-side face now rides it, with two candidate rules
  for the reviewer to pick between. One decision closes both phases.
- **The four post-approval spec corrections were not individually vetoed.** They
  were surfaced at the drain close, the reviewer approved the batch including the
  retirement, and the spec file went with it. Recoverable via
  `bench spec history ft91-artifact-build-tiering`. The code was right in all
  four cases; only prose moved.
- **`ft91-gate-phase-split` stays unretired on purpose,** so `bench status` keeps
  reporting one spec awaiting retirement until the reviewer rules on its stories
  4, 5, and 9 — two shipped as probed phases rather than the manifest the spec
  named, and story 9 was dropped as unsatisfiable.
- **`decisions/cost-follows-project-size.md` is kept, not deleted.** Ticket #6 is
  open and `## Not yet specified` still parks `-count=1` freshness semantics as a
  reviewer-led oracle decision, so `bench status` reports one unresolved map by
  design.

## Next command

`/bench-shape-idea` for FT91's eighth arm.

Three open decisions want one map rather than three patches: why the gate
absorbed only 24 s of a 131 s suite win, which tests may share one prepared
artifact set, and whether `decisions/gate-pipeline.md` — closed on a premise the
slice-C measurement falsified — reopens. Inputs are the FT91 row,
`decisions/gate-concurrency.md`'s watch-outs, and
`decisions/cost-follows-project-size.md`.

Behind it: `/bench-write-spec` for FT98 (the one preserve-then-discard
primitive), then FT71 (versioned local shift evidence, the remaining HIGH
bank-track row).

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
