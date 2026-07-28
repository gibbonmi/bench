# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `c56be81`, clean tree, 2 unpushed commits (`9977c27` drain,
`c56be81` handoff)
Spec: `specs/implement-spec-full-run.md` — `Status: implemented`, merged, awaiting
retirement.
Gate: green at `c56be81`.

## State

- **FT152 is closed out.** `/bench-final-check` ran the gate on `main` at
  `c56be81`: green across build, gofmt, vet, test, race, conformance,
  conformance-suite, contract, shellcheck, canary. No commit — the work was
  already landed and the status already flipped, so the phase was an honest
  no-op. Ship-tier verification has not run; `bench prep-release` covers that
  once per release.

- **The FT152 build commits are already pushed.** Only the drain and this
  handoff remain unpushed. No leftover worktrees, no orphaned review pickup.

- **Retirement is the one open FT152 action, and it is the reviewer's.**
  `bench spec retire implement-spec-full-run` deletes the spec and its decision
  map, and the spec header (lines 17-56) carries a 15-item veto list of calls
  decided in the spec rather than carried by the map. Two were made on blanket
  authorization and are the ones most worth a look: the `.bench/structure.budgets`
  grant for `docs_workflow_helpers_test.go 660`, and correcting the spec's
  `.claude/commands/` mirror decisions, which assumed a copied tree where this
  repo has a symlink. Read the list, then retire — retiring first destroys the
  surface.

- **Next build is FT154**, marked unblocked by the 2026-07-28 drain and specs-next
  after FT152. `bench status` also flags 35 structure issues, 2 unresolved
  decision maps, and 1 roadmap row for merged work — the roadmap row belongs to
  `/bench-what-next`, not to a build.

## Next command

Reviewer decision first: read the veto list at
`specs/implement-spec-full-run.md:17-56`, then retire with
`bench spec retire implement-spec-full-run`.

Then `/bench-write-spec` for FT154.

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
