# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `85c55fc`, 8 dirty paths, 3 unpushed commits
Spec: `specs/ft91-gate-phase-split.md` (Status: staged)
Gate: green at `78c917e` — stale, work tree `a68a7e1`

## State

- **The FT91 slice C spec is staged, reviewer-approved, and committed
  (`85c55fc`).** `specs/ft91-gate-phase-split.md` — the `checkGoCore` split into
  first-class gate phases. It is the next build, and it is the first consumer of
  slice B's manifest loader and DAG scheduler, which shipped with nothing using
  them.
- **Two closed decisions the build must not reopen.** Canary fixture EXPECTs
  carry the tool's own diagnostic only, *not* the `phase <name>: red` half the
  map prescribed — inner mode emits no summary lines, and that byte shape is
  pinned. And the seven stray `CHECK` files the map says to add already shipped
  with slice A; the map's "six other checks" is stale (it is five). Both
  corrections are written into the spec.
- **The spec's veto list is (a)–(m) in its Implementation decisions.** Two are
  worth a second look before merge: the `prep-release` ship-parity step, which
  overlaps an item FT142 already carries for veto, and story 20's `fable` line
  for the consumer-facing manifest docs.
- **FT143 is bundled into this spec** (kit-root family→check binding assertion)
  and closes with it. FT144 stays open — its workflow question is the reviewer's,
  unmade; the spec applies FT144's edge-inventory half without waiting on it.
- **A `/bench-what-next` drain closed 2026-07-27 (`5614726`).** `.bench/learnings.md`
  carries no open entry. `IDEAS.md` holds one parked item: `bench handoff`'s pin
  block reports a stale dirty-path count (printed 8 against an actual 1, twice),
  so read that line as advisory, not as tree truth.
- `internal/gate/` at 21 files against a granted 16 is a live `bench structure`
  advisory, and this slice adds to it. It surfaces on its own.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows.
- The branch/worktree sweep remains proposed, not executed — reviewer's call.

## Next command

`/bench-implement-spec` against `specs/ft91-gate-phase-split.md`, in a fresh
mid-tier session. Interactive, not `bench shift`: stories 14 and 15 author canary
fixtures whose EXPECT strings settle at build time, which fails `craft-line`'s
venue-routing test for an unattended loop. Read `decisions/gate-pipeline.md`
#3/#6/#7 and `decisions/assets/gate-pipeline-fixture-inventory.md` first. Adding
`.bench/phases.json` changes the `.bench` tree, so the pre-push hook wants a fresh
`bench gate pin` before any push.

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
