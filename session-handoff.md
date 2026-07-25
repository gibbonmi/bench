# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `149e8d2`, clean tree, 20 unpushed commits
Spec: `specs/ft86-fail-closed-control-records.md` (Status: staged)
Gate: green at `7ca7056` — stale, work tree `6bb335c`

## State

- **FT86 is spec'd and staged, not built** (`149e8d2`). The spec was compiled
  from the closed `decisions/ft86-fail-closed-control-records.md` map, whose
  Handoff is complete with no uncertainty flags. The reviewer approved it. The
  header above is the single source for branch, HEAD, tree state, spec status,
  and the gate verdict — do not restate those here.
- **What FT86 does:** `internal/bounds` becomes the one classified reader (six
  typed states, Lstat-first, size-bounded); `learnings`, `maps`, `roadmap`, and
  `coverage` fail closed on any non-absent read failure; `status` and
  `roadmap --context` stay exit 0 but degrade visibly to `unknown`;
  `coverage --check` requires a map or the historical marker and validates exact
  story membership; `git.DefaultBranch` is deleted with `ResolvedDefault` the
  sole owner.
- **Two calls are flagged in-spec for post-hoc veto and stay open until vetoed.**
  Story 16 deviates from the map — `outline` splits its skip classes into failure
  classes (exit 1) and content classes (binary/oversized, exit 0), because
  exiting 1 on every binary file would make the command unusable. And the
  ROADMAP row's `traversal` state is not adopted; it guards the
  caller-supplied-path surface rather than the control-record reads, so it sits
  in Out of scope at 3 edits, 2 gate runs.
- **`DefaultBranch` has eight call sites, not the two the map implied** — `adopt`
  (×2), `worktree` (×2), `status` (×2), and `git`'s own internals. Each posture
  is decided in the spec's Implementation decisions; read that section before
  touching any of them. This is why story 18 routes mid-tier.
- **The build must not run as a bare `bench shift`.** Stories 1–3, 7, 16, and 18
  route `gpt-5.6-terra`, so the spec fails `craft-line`'s all-cheap venue test.
  Run `/bench-implement-spec` interactively and escalate per story.
- **Step 9's falsification pass never ran.** Story 16's deviation fires the
  trigger, but the authoring session carried a standing no-unasked-delegation
  rule. Running it is still available and is the cheapest second opinion on the
  outline split.
- **`.bench/learnings.md` carries three open entries** from the FT122 build, all
  rule-shaped, all awaiting a `/bench-what-next` verdict: the canary's panic
  reporting, the shared-checkout write-delegation, and `bench idea` voiding an
  in-flight gate verdict.
- **Never mutate the repository while a gate runs.** `projects/benchkit.md`'s
  cold-session notes carry this and the `internal/canary` nested-run trap; read
  them before touching either, and note that `dist/bench` must be built with
  `scripts/go-build.sh`.
- **FT91 was raised to HIGH by the reviewer** (`9cbb138`) — the gate's length is
  the dominant cost of small changes here. A fourth arm is recorded on the row:
  `RunConformance` runs its fifteen checks serially in one test, ~94% of wall
  clock across two measured runs. First step is timing each check before
  committing to the arm.
- `bench structure` reports 17 issues, all pre-existing.
## Next command

`/bench-implement-spec — specs/ft86-fail-closed-control-records.md`

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
