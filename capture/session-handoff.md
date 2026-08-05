# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c3d8e47`, 7 dirty paths, 42 unpushed commits
Spec: `specs/authoring-hardening/spec.md` (Status: staged), `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `1d9aef2` — stale, work tree `96b2ce4`

## State

**Phase reached: `/bench-final-check` closed on `gate-evaluation-snapshot` — promoted
terminal, retired.** Candidate `60b44de` published as squash `c3d8e47` on `main`; the
run's two accepted review findings (S1 optimistic-reuse bypass, C1 malformed listing
metadata) were repaired as lifecycle tickets, re-reviewed clean, and promoted. Exact
green retained at `refs/bench/green/main`.

- Uncommitted: the `specs/gate-evaluation-snapshot/` deletion staged by `bench spec
  retire` and the new `capture/retros/gate-evaluation-snapshot.md`, awaiting the
  `spec-retire: gate-evaluation-snapshot` gated commit — the next session's first act if
  this one stopped before landing it.
- No roadmap row names the spec (it came from the `decisions/gate-budget.md` ordering);
  the retirement commit carries only the deletions and the retro.
- Pending capture: one retro (`capture/retros/gate-evaluation-snapshot.md`) and 4 open
  learnings — the reviewed `/bench-what-next` drain owns both.
- 42 unpushed commits on `main`; pushing is the reviewer's call.
- Closed decisions stay closed: gate-budget/gate-concurrency work preserved; promotion
  remains the sole whole-project gate and landing author.

## Next command

`/bench-what-next`

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
