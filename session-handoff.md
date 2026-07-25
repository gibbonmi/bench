# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `bc64686`, 4 dirty paths, 16 unpushed commits
Spec: none staged.
Gate: green at `7ca7056` — stale, work tree `f077ff4`

## State

- **FT122 is built, reviewed, fixed, and landed green** (`bc64686`). The spec is
  retired and its decision map deleted; the roadmap row is gone. The header above
  is the single source for branch, HEAD, tree state, spec status, and the gate
  verdict — do not restate those here.
- **The three-axis review found 18 issues and all are resolved in `bc64686`.**
  The three that mattered: a runtime expectation still held the pre-fix invocable
  rule and would have failed green work; `bench handoff` silently discarded
  everything below any `## ` heading a reviewer wrote inside their own State
  section, on exit 0; and `validate` borrowed `toon.Representable`, which permits
  the tab, newline, and return a line-structured document cannot survive. The
  review pickup was deleted in the same commit that closed its findings.
- **Story 18's idempotence is fixed for the pin block and parked for the board.**
  The dirty clause now excludes the handoff itself, so the command no longer
  prints a fact its own write falsifies. Strict byte-identity on a *tracked* tree
  also needs the board to ignore that write, which needs a `bench status` path
  exclusion that does not exist — recorded as out of scope on the retired spec
  and parked in `IDEAS.md`. This is the one call from the fix pass most worth a
  second look.
- **The status board's handoff row now reads `bench handoff`, not prose.** It is
  therefore selectable as a next command, which it never was before.
- **`.bench/learnings.md` carries three open entries** from the FT122 build, all
  rule-shaped, all awaiting a `/bench-what-next` verdict: the canary's panic
  reporting, the shared-checkout write-delegation, and `bench idea` voiding an
  in-flight gate verdict.
- **Uncommitted right now:** the profile's two new hostile-input classes, the
  roadmap row removal, and the spec/decision-map deletions. They are the
  spec-retire commit and nothing else.
- **`bench commit` advises "set them aside" for files outside the named set, and
  no such route exists.** Parked in `IDEAS.md`. Until it exists, the only exits
  are naming the file in the commit or reaching for guard-blocked git; name it.
- **Never mutate the repository while a gate runs.** `projects/benchkit.md`'s
  cold-session notes carry this and the `internal/canary` nested-run trap; read
  them before touching either, and note that `dist/bench` must be built with
  `scripts/go-build.sh`.
- **FT91 was raised to HIGH by the reviewer** (`9cbb138`) — the gate's length is
  the dominant cost of small changes here, and this session paid it four times. A
  fourth arm is recorded on the row: `RunConformance` runs its fifteen checks
  serially in one test, ~94% of wall clock across two measured runs. First step
  is timing each check before committing to the arm.
- `bench structure` reports 17 issues, all pre-existing.

## Next command

`/bench-what-next` — the board's leading invocable signal (`drain`).

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
