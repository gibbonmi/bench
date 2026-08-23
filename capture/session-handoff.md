# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` once this spec phase lands; the source is the `ft113-commit-spec-residuals` worktree on base `d3390d20`
Spec: `specs/landing-authors-the-flip/spec.md`, `Status: staged`, five tickets under `specs/landing-authors-the-flip/tickets/`
Gate: green on `d3390d20` (2026-08-23)

## State

The FT113 spec is staged and awaits the reviewer's sign-off on the approval
table. The review round (opus / high, read-only) ran once; the author folded
its eight blocking findings and the spec names them in its verification log.
`bench coverage --check` passes with 24 rows, and the conformance package is
green on the worktree.

Closed decisions, dated 2026-08-23:

- `bench worktree land --spec` is the one author of the flip and the tickets-only close
- `bench commit` drops `--spec`, and `bench spec implemented` retires
- a retirement that does not complete names its remainder, with no retry verb
- Bench does not remove the board row or the detail file

Ticket order: `retire-names-the-board-remainder.md` lands first on its own
gate; then `retire-the-commit-route-flip-and-close.md`,
`retire-bench-spec-implemented.md`, `commit-exit-3-names-the-remainder.md`,
and `guidance-names-one-author.md` in blocker order. The build runs in a
fresh mid-tier session in a bench worktree and lands through `bench worktree
land --spec landing-authors-the-flip`.

## Next command

`/bench-implement-spec` — specs/landing-authors-the-flip/spec.md

## Shape

Rewritten in full at every phase close, pruned rather than accreted. A fresh
session pays for every line it reads cold; drop anything it would not act on.

Operational gotchas are placed by lifetime, not copied here. One that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes. One scoped to a build
belongs instead in that spec's coverage rows.

This file names at most when you'll hit one, never the command — a second copy
drifts from the source.

Keep the three sections above. **State** holds what is true now, including anything
uncommitted. **Next command** holds the exact harness-native invocation, not a
description of it. This section is the third.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
