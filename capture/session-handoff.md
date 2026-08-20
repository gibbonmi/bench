# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main`

## State

FT225 is specified and signed off. `specs/land-spec-amendments/spec.md`
(`Status: staged`, 13 coverage rows, verification log: 2 iterations to accept)
and its three tickets are committed at `58bac422`. The reviewer approved the
spec, the ticket breakdown, and the `fable / high` line for the guidance
ticket. FT225's roadmap row stays until shipped.

Closed decisions that stay closed (also recorded in the spec): the landing
publishes `Implemented(source-tip spec bytes)` and the composition neutralizes
the spec path so a destination-side amendment — overlapping included — never
conflicts; `paths-authorized` implicitly authorizes `specs/<slug>/` via the
existing spec-path fact; the reviewer's typed `--base` and `--source-tip` are
the acceptance; the in-range fence-widening exposure is a named Won't handle
with the typed range as the control.

Ticket order on one retained integration source: 1
`authorize-the-spec-folder-implicitly`, 2 `publish-the-source-spec-bytes`
(blocked by 1), 3 `name-the-amendment-cadence-in-review-guidance`. The spec's
explicit self-fence entry is deliberate bootstrap — this build is the last one
that needs it.

## Next command

`/bench-implement-spec` — land-spec-amendments, in a fresh mid-tier session on
one retained integration source

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
