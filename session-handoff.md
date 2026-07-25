# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `7485601`, clean tree, 13 unpushed commits
Spec: `specs/session-handoff-emission.md` (Status: staged)
Gate: green at `58a6303` — stale, work tree `6df45ef`

## State

- **FT122 is built and landed green.** All 26 stories, 30 coverage rows, every
  phase green. The header above is the single source for branch, HEAD, tree
  state, spec status, and the gate verdict — do not restate those here.
- **Semantic review has not run — this is the outstanding phase.** The reviewer
  still owes `/bench-review-implementation` against
  `specs/session-handoff-emission.md`, which is why the spec reads `staged` and
  not `implemented`. Nothing downstream should treat FT122 as closed until it
  runs. `/bench-final-check` owns the landing commit and the `Status: implemented`
  flip. Do not flip that line by hand.
- **The Claude 5 skill-realignment batch is merged** (`751d45c`, fast-forward,
  all six changes reviewer-approved; its worktree is released). Dogfooding the
  trimmed descriptions caught one regression, fixed in `7485601`: an unquoted
  ` #` in `craft-line`'s frontmatter description started a YAML comment, so the
  harness loaded only the first of three trigger clauses and never saw the
  delegation or tier-escalation triggers. `craft-skills`, `craft-tdd`, and
  `craft-design-system` fire correctly on their triggers. A gate check for the
  same YAML trap is parked in `IDEAS.md`.
- **Story 8 was reopened and re-decided mid-build.** The board's `Action` is a
  prose hint, not a command — most rows read `fix before commit` or
  `split (craft-seams)` — so the original derivation rendered a hint where the
  field promises an invocation. The reviewer chose *skip to the first invocable
  signal*: walk the board in its own severity order, take the first action that
  is a command, never re-rank. The accepted tradeoff is that a red gate can be
  skipped over in this field; the header's `Gate:` line carries it regardless.
- **Two defects came out of driving the command rather than reading it.** A
  runtime test sliced a subject-reported hash unguarded, which panicked under
  the canary's stub subject, killed the whole test binary, and left an unrelated
  fixture reporting `did not bite` — the canary names the fixture that went
  quiet, not the test that silenced it. And the compound git action
  `/bench-final-check / push` passed a bare prefix test while being two
  commands. Both are fixed and both have recorded red signals.
- **`.bench/learnings.md` carries three open entries** from this build, all
  rule-shaped, all awaiting a `/bench-what-next` verdict: the canary's panic
  reporting, the shared-checkout write-delegation, and `bench idea` voiding an
  in-flight gate verdict.
- **`bench commit` advises "set them aside" for files outside the named set, and
  no such route exists.** Parked in `IDEAS.md` (`1da9ad1`). Until it exists, the
  only exits are naming the file in the commit or reaching for guard-blocked
  git; name it.
- **Never mutate the repository while a gate runs** — it cost a full re-run here
  even though every phase was green. `projects/benchkit.md`'s cold-session notes
  carry this and the `internal/canary` nested-run trap; read them before touching
  either, and note that `dist/bench` must be built with `scripts/go-build.sh`.
- **FT91 was raised to HIGH by the reviewer** (`9cbb138`) — the gate's length is
  the dominant cost of small changes here. A fourth arm is recorded on the row:
  `RunConformance` runs its fifteen checks serially in one test, and that phase
  was ~94% of wall clock across two measured runs. First step is timing each
  check before committing to the arm, because the win depends on how the ~500s
  distributes. The gate reads stale in the header only because that doc commit
  landed after the green run.
- `bench structure` reports 17 issues, all pre-existing.

## Next command

`/bench-review-implementation` against `specs/session-handoff-emission.md` — the
outstanding phase for FT122, owed before the spec can go `implemented`. The
board's leading invocable signal is `/bench-what-next` (`drain`, now 3 ideas and
3 open learnings); take it after the review.

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
