# Session handoff

Repository: `bench` (origin `github.com/gibbonmi/bench`) — branch `main`,
checked out at `~/workspace/bench` on the machine that wrote this. Everything
below is executable from a cold start; no conversation history is needed.

## State

- **FT122 is specced, approved, and ready to build.**
  `specs/session-handoff-emission.md` (Status: staged) adds `bench handoff`:
  it derives the cold-session pin block from the tree, prints it, and rewrites
  this file — preserving `## State` byte-for-byte while regenerating the header,
  `## Next command`, and `## Shape`. 24 stories, 24 coverage rows, two seams
  (the command boundary against a fixture repo, and unit tests for the section
  splitter). Nine stories route to mid; the rest is cheap-tier plumbing.
- **Its map was closed in the same session as the spec**, under the
  reviewer-closed path in `/bench-write-spec`'s entry contract.
  `decisions/session-handoff-emission.md` carries three **[veto]** marks —
  recorded alternatives, not open questions. **No falsification pass ran**: the
  same-session map is a mandatory trigger for one, but that session was
  configured not to spawn subagents unasked. Running it before the build is
  cheap insurance and remains available.
- **Three unpushed commits, and the gate has not run on this tree.**
  `99139d0` (capture), `faaca55` (the spec and map), `849ff72` (ADR 0009). All
  three are doc-only and landed with plain `git commit`, so the last green gate
  is still the one at `4ea4880`. Pushing is the reviewer's call.
- **Five CLI candidates are parked in `IDEAS.md`, awaiting a drain.** FT122 came
  out of that survey; the other four — `bench worktree path` / `exec` (which
  must render `~`-relative paths for cross-machine portability), `bench test`,
  `bench spec show`, `bench outline --symbol` — have no maps yet. They graduate
  through `/bench-what-next`, not directly into specs.
- **FT120 and FT121 stay open** from the 2026-07-24 drain: gate/canary
  test-harness defects plus a `canary_concurrency_test.go` split, and the
  `bench commit --spec` gate-staleness defect. Both want specs.
- `bench structure` reports 17 issues, all pre-existing. Read
  `projects/benchkit.md`'s cold-session notes before touching
  `internal/canary` — the nested-run trap deadlocks rather than fails, and this
  file deliberately does not restate the commands.

## Next command

`/bench-implement-spec specs/session-handoff-emission.md` — interactive rather
than `bench shift`, because nine stories route to mid tier and the coverage map
is not uniformly cheap, which fails `craft-line`'s venue-routing test for an
unattended loop. In Codex, that is `$bench-implement-spec`.

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
