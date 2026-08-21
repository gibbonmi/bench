# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `b2deece`, 2 dirty paths, 0 unpushed commits
Spec: `specs/learnings-dated-line-visibility/spec.md` (Status: staged)
Gate: green at `4cd77bc` — stale, work tree `7352eeb`

## State

`specs/learnings-dated-line-visibility/spec.md` is revised, re-sliced into two
tickets, and awaiting reviewer sign-off. Nothing is implemented.

The spec now carries both halves of FT243. Ticket 01 reports a dated line that
misses heading shape; ticket 02 exports `<!-- entries below -->` from
`internal/learnings` and reports content below it that belongs to no entry.
Ticket 02 is blocked by ticket 01. `bench coverage --check` reads 34 rows valid
and all four `bench preflight build learnings-dated-line-visibility` checks are
green.

The reviewer's three entry-round dispositions are closed and stay closed: the
Line is mid (`opus` / medium) against the cached cheap-tier routing, the scope
is widened past dated lines to undated content below the marker, and the ticket
row-ID citation was repaired here rather than routed back.

Two dispositions are open and named in the spec. The review round observed that
ticket 01 is a narrower capability that could ship on its own gate, with ticket
02 deferred. The kit's own journal has been pruned past its marker, so ticket
02's rule is inert in this repo until the marker is restored — restoring it also
narrows `internal/conformance`'s stale-slash-reference scan, so the spec puts
that pair in Out of scope rather than in the build.

`specs/` is untracked, which is the veto surface. `capture/learnings.md` holds
two open entries, both about this phase.

## Next command

`/bench-implement-spec specs/learnings-dated-line-visibility/spec.md --full`

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
