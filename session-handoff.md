# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `4ffd607`, tree clean, 33 unpushed commits
Spec: `specs/ft126-recurrence-tallying/spec.md` (Status: staged), `specs/ft128-agent-line-binding/spec.md` (Status: staged)
Gate: green at `4ffd607` (full run during the map commit; dist/bench rebuilt and resealed)

## State

- **The diff-visual decision map is closed and committed at `4ffd607`.** All 12
  tickets resolved, `## Handoff` written and placeholder-free, `bench maps`
  shows no row. Assets committed beside it: `decisions/diff-visual-prototype.html`
  (approved working visual, iterated live against the real FT131 change) and
  `decisions/diff-visual-ft131-report.html` (shareable, fully offline sample).
  `/bench-write-spec` compiles the map and moves both assets into
  `specs/<slug>/decisions/`. Recommended slicing in the map's Handoff:
  schema core (emitter + validator) → renderer → final-check integration.
- **Decisions that stay closed:** report form and sections (#10), edge kinds
  three-drawn-two-relational (#11), no layout dependency (#12, supersedes #9),
  context nodes admitted as their own class (#2 amendment), three-valued
  coverage badge and justified heat (#6 amendment).
- **A parallel session migrated the decision map corpus at `2b6ee58`** and
  committed the previously dirty reviewer files; the tree is clean.
- **FT131 remains implemented and retired on `main`.** `.bench/retros/` holds
  no ft131 file — a prior handoff claimed one pending; the tree wins.
- **Drain still pending:** parked ideas and the staged FT126/FT128 specs are
  `/bench-what-next` territory, untouched by this phase.

## Next command

`/bench-write-spec Compile the closed decisions/diff-visual.md into the build spec; the Handoff section carries seams, contracts, and uncertainty flags.`

Run it on a fresh mid-tier session.

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
