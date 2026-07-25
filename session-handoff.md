# Session handoff

Repository: `bench` (origin `github.com/gibbonmi/bench`) — branch `main`,
checked out at `~/workspace/bench` on the machine that wrote this. Everything
below is executable from a cold start; no conversation history is needed.

## State

- **FT122 is specced, falsified, amended, and ready to build.**
  `specs/session-handoff-emission.md` (Status: staged) adds `bench handoff`:
  it derives the cold-session pin block from the tree, prints it, and rewrites
  this file — preserving `## State` byte-for-byte while regenerating the header,
  `## Next command`, and `## Shape`. 26 stories, 29 coverage rows, two seams
  (the command boundary against a fixture repo, and unit tests for the section
  splitter, each unit seam paired with a runtime row that proves the command
  actually calls it).
- **One story carries a reviewer-approved top-tier bump.** Story 15 — the
  scaffolded skeleton's guidance prose — routes `gpt-5.6-sol` / high under the
  leverage override, bounded to that prose. Story 14, moving today's `## Shape`
  text into the binary, stays cheap because it is transcription. Nine stories
  route to mid, the rest cheap; every deviation from the profile's CLI-plumbing
  routing is named in its story.
- **A falsification pass ran and its eleven findings are folded in** (`2f51fe9`).
  The one worth carrying forward: the spec had introduced a gate-verdict field
  with no decision behind it, and a bare `gate: green` would have emitted a
  cached verdict whose tree had already moved. The map now decides it at #7 and
  the field renders verdict, cached tree, and staleness together. Four other
  blockers were degenerate-implementation holes — one-fixture rows a hardcoded
  constant would pass, and a `## Shape` row that rewarded preserve-and-append.
  `decisions/session-handoff-emission.md` carries the **[veto]** marks.
- **Five unpushed commits, and the gate has not run on this tree.**
  `99139d0` (capture), `faaca55` (the spec and map), `849ff72` (ADR 0009),
  `8cbb34a` and this one (handoff rewrites), plus `2f51fe9` (the amendment). All
  are doc-only and landed with plain `git commit`, so the last green gate is
  still the one at `4ea4880`. Pushing is the reviewer's call.
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
than `bench shift`, because nine stories route to mid and one to top, which
fails `craft-line`'s venue-routing test for an unattended loop. In Codex, that
is `$bench-implement-spec`.

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
