# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5614726`, 8 dirty paths, 1 unpushed commit
Spec: none staged.
Gate: green at `78c917e` — stale, work tree `a5f10c1`

## State

- **A `/bench-what-next` drain closed 2026-07-27 (`5614726`) and both inboxes
  are empty.** 49 rows in, 49 out — nothing was removed, because FT91 has arms
  left. `IDEAS.md` is empty and `.bench/learnings.md` carries no open entry.
- **FT91 slice B landed and its spec is retired.** The `.bench/phases.json`
  loader and the DAG scheduler are in; nothing consumes them yet, because the
  kit ships no manifest of its own. Slice C is what first does.
- **Slice B's four open reviewer items are now board rows, not handoff state.**
  The two coverage-map corrections went to FT140 alongside FT86's identical
  provenance question — one decision closes all three instances — and the
  falsified `os/exec` premise went to FT99 as its third instance. The residual
  `runPhases` fail-open was closed by decision (unreachable: `Execute` overrides
  to 124 and a bare `bench gate-phases` carries no deadline). `internal/gate/`
  at 21 files against a granted 16 is a live `bench structure` advisory, so it
  surfaces on its own.
- **Slice C is the next build** — `checkGoCore` split, fixture migration,
  parity test, probed fallback phases (decision map #3/#6/#7) — and its spec
  also carries FT143's kit-root family→check binding assertion. Consumer-facing
  manifest docs land with it. FT143 and FT144 stay open until then; FT144's
  workflow question is the reviewer's, unmade.
- **Two folds in the drain are flagged for veto.** FT104 was widened past its
  original title to carry both faces of a load-induced commit refusal, and the
  review residuals were folded onto existing rows rather than opening FT145 —
  both are FT137's not-yet-shipped restructuring moves, applied on the
  2026-07-26 precedent.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows.
- The branch/worktree sweep remains proposed, not executed — reviewer's call.

## Next command

`/bench-write-spec` — FT91 slice C: split `checkGoCore` into first-class gate
phases. Inputs: the FT91 row and `decisions/gate-pipeline.md`.

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
