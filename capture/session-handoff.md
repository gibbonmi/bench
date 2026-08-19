# Session handoff

Repository: `9a43f8b9d7384e58ad06b1bcea12e3c0-47e4fe85241ecdf736dfbcb48a8117db` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/.bench/worktrees/bench-2826441890/9a43f8b9d7384e58ad06b1bcea12e3c0-47e4fe85241ecdf736dfbcb48a8117db`
Branch: `bench/assign/9a43f8b9d7384e58ad06b1bcea12e3c0/47e4fe85241ecdf736dfbcb48a8117db` — HEAD `0d32a9b`, clean tree, 4 unpushed commits
Spec: `specs/bench-front-door/spec.md` (Status: staged)
Gate: green at `20dcb24` — stale, work tree `f071524`

## State

**`/bench-implement-spec specs/bench-front-door/spec.md --full` resumed on the retained
integration source.** `bench preflight build bench-front-door` was green. Ticket 01
(`e6d5e623`) extracted the route owner into `internal/status`; ticket 02 (`0d32a9bf`)
normalized all current status actions and the dashboard's empty cell. Each commit passed
the whole gate. Ticket 02's first all-actions test was repaired when an independent
mutation of the live intent producer left it green; it now drives the actual producers.

**Next frontier:** `03-setup-and-staged-signals.md`, dependent on `02-normalize-actions.md`.
The retained source is clean at `0d32a9bf`; `main` remains at `408faf50`. The reviewer-visible
decisions in the spec's Further notes remain open and unchanged. Ticket 01's advisory
`Writes:` list omitted its necessary conformance-location repair; the spec ownership fence
authorizes it, and it is reviewer-veto surface rather than a behavior change.

## Next command

`$bench-implement-spec specs/bench-front-door/spec.md --full --reviewer sol high`

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
