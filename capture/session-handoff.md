# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` once the landing runs; the source is the `ft113-build` worktree (request id in the scratchpad of the building session), base `0e17d428`
Spec: `specs/landing-authors-the-flip/spec.md`, `Status: staged`, five tickets built
Gate: green on every ticket commit (2026-08-23)

## State

The FT113 build is implemented in the `ft113-build` integration worktree. Ticket
commits, in order: `624ffb5b` (Group C), `9d84030a` (Group D, merged as
`41be490c`), `66bbc334` (Group A commit route), `c2e7951c` (Group A spec verb),
`0b29bd9c` (Group B, merged as `42eb3bec`). The `/bench-implement-spec --full`
run stands at the review boundary.

Open flags for the reviewer:

- the spec fence list gained `internal/spec/history.go`,
  `internal/sanitize/sanitize.go`, and `internal/worktree/worktree.go`: the
  build exported `sanitize.LineSafe` as the one line-safe predicate and
  fixed a stale comment
- `capture/IDEAS.md` holds one parked idea about the landing verb's `next=`
  sanitizer

Closed decisions, dated 2026-08-23:

- `bench worktree land --spec` is the one author of the flip and the tickets-only close
- `bench commit` has no `--spec`, and `bench spec implemented` is gone
- a retirement that does not complete names its remainder, with no retry verb
- Bench does not remove the board row or the detail file

## Next command

`/bench-review-implementation` — specs/landing-authors-the-flip/spec.md over base `0e17d428` and the `ft113-build` source tip; then `bench worktree land --spec landing-authors-the-flip`; then `/bench-final-check`

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
