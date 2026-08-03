# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7dbb686`, clean tree, 21 unpushed commits
Spec: `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `e11c4c7` — stale, work tree `e11c4c7`

## State

- **The `/bench-what-next` drain is committed** (`626cf92`), and
  `ft164-ticket-contracts` is retired (`7dbb686`). All three capture sources
  are empty: 3 ideas, 3 learnings, 1 retro drained.
- **FT164's row is now its residual only** — the repair-and-experiment lane
  the spec cut (model-comparison, inventory-currency, shared-cache charges),
  plus two rules from its own build: a done-claim resolves every owner named
  in a ticket's Red-mutations table to a real artifact, and a repair round
  trues up any enumeration it extends. It names no spec path.
- Rulings recorded during the drain: FT173's principle-9 byte-relaxation
  applies to the `help[]` spec alone; FT175's spec-start gate is
  owners-consumed, and its three ledger decisions are deferred behind the
  foundation and `help[]`. FT130 gained the tip-move-under-active-run face;
  FT156 gained the anchor-file budget signal (864 vs 660).
- Two restructure candidates were named but not applied (a default run):
  FT130 now carries three faces of one decision, and FT164's residual could
  fold into the FT179 + FT111 visit. Both are open for a `--restructure` pass.
- Nothing pushed; push is the reviewer's call.

## Next command

`/bench-shape-idea` — the board's leading invocable signal (`decisions`).

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
