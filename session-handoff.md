# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `35c9a73`, 9 dirty paths, 3 unpushed commits
Spec: `specs/ft96-delegation-discipline.md` (Status: staged)
Gate: green at `f28e44d` — stale, work tree `b4b63d6`

## State

- **The FT96 spec is staged and reviewer-approved.** Ten stories, three seams,
  31 coverage rows. Everything the build needs is in the spec; do not re-derive
  it from the roadmap row.
- **No decision map backs it.** It was compiled from `ROADMAP.md`'s FT96 row
  under `/bench-write-spec`'s reviewer-directed batch-drain override. The spec's
  `## Flagged for reviewer veto` section lists all six defaulted decisions; the
  reviewer approved with those flags in place, so they are **closed** — do not
  reopen them mid-build.
- **A `fable` falsification pass returned block; all three findings are folded
  in.** They grew the spec from nine stories to ten: the veto-flag section was
  missing, story 1's cheapest wrong build passed every row (fixed by story 10's
  anchor rows), and clause 3's concrete substitute plus clause 2's rejected arm
  had been dropped without record.
- **Stories 8, 9, and 10 are the only gate-observable arms.** Stories 1–7 are
  guidance prose in one owner file
  (`.agents/skills/bench-craft-delegate/SKILL.md`); the reviewer is their oracle,
  which is why they route top.
- **Do not run `bench shift` on this spec.** It fails `craft-line`'s
  venue-routing test — stories 1–7 route top and are not gate-observable. The
  build is interactive.
- Verified at spec time: `bench coverage --check` green, and
  `BENCH_CONFORMANCE_ROOT=$PWD go test ./internal/conformance -run
  '^TestRootConformance$'` green (228s) with the spec staged, so the docs layer
  is clean going in.
- **Both capture sources are at zero.** `IDEAS.md` empty, `.bench/learnings.md`
  holds no open entries.
- Known ambient facts: the gate is stale from this commit onward
  (`session-handoff.md` is not on the capture-only allowlist — FT113's gate
  face); the FT91 conformance-phase long pole and the structure budget
  violations are unchanged.

## Next command

`/bench-implement-spec specs/ft96-delegation-discipline.md` — in a **fresh
session on `opus`** (the mid binding), which escalates per story rather than
running the whole build on the top model's big-context iteration.

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
