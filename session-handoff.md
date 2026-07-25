# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `149e8d2`, clean tree, 20 unpushed commits
Spec: `specs/ft86-fail-closed-control-records.md` (Status: staged)
Gate: green at `7ca7056` — stale, work tree `6bb335c`

## State

- **FT86 is spec'd, amended, and staged — not built** (`5bee3ce`). The spec was
  compiled from the closed `decisions/ft86-fail-closed-control-records.md` map,
  then corrected against a falsification pass that returned BLOCK. The reviewer
  approved both. The header above is the single source for branch, HEAD, tree
  state, spec status, and the gate verdict — do not restate those here.
- **What FT86 does:** `internal/bounds` becomes the one classified reader (six
  typed states, Lstat-first, size-bounded); `learnings`, `maps`, `roadmap`, and
  `coverage` fail closed on any non-absent read failure; `status` and
  `roadmap --context` stay exit 0 but degrade visibly to `unknown`;
  `coverage --check` requires a map or the historical marker and validates exact
  story membership; `git.DefaultBranch` is deleted with `ResolvedDefault` the
  sole owner across **eleven** enumerated call sites.
- **`bench outline` is deliberately OUT of the classifier migration.** The map
  grouped it with the fail-closed commands; the spec reverses that. Migrating it
  would break `TestAXIOutlineSymlinkSkipped`, which encodes a recorded decision:
  following a tracked symlink indexes the target's symbols under the symlink's
  path, emitting anchors that don't hold. Story 5's follow-the-link rule binds
  control records only. Do not "fix" outline to match the map.
- **Three calls are flagged in-spec for post-hoc veto** and stay open until
  vetoed: the outline reversal above, declining the `traversal` state the
  ROADMAP row names, and `git.Facts`'s unresolved-default shape (a new
  `DefaultResolved` field with `Ahead`/`Behind` zeroed).
- **The build must not run as a bare `bench shift`.** Stories 1-3, 7, and 18
  route `gpt-5.6-terra`, so the spec fails `craft-line`'s all-cheap venue test.
  Run `/bench-implement-spec` interactively and escalate per story.
- **Several red signals are capability-gated** (chmod, FIFO, symlink). A skip and
  a deleted assertion both look green, so the build is not done until the
  affected suites pass with `BENCH_REQUIRE_CAPABILITIES=1` and zero skips.
- **`/bench-write-spec` step 9 changed this session** (`54e4cae`). The
  falsification pass is no longer conditional: every draft gets one at the mid
  binding, spawned without asking. The five former triggers now nominate a draft
  for a *top*-tier pass, which pauses and asks. A new canary fixture
  (`write-spec-review-made-conditional`) guards the always-on rule via the
  `Every draft gets the pass` anchor.
- **`.bench/learnings.md` carries three open entries** from the FT122 build, all
  rule-shaped, all awaiting a `/bench-what-next` verdict: the canary's panic
  reporting, the shared-checkout write-delegation, and `bench idea` voiding an
  in-flight gate verdict.
- **Never mutate the repository while a gate runs.** `projects/benchkit.md`'s
  cold-session notes carry this and the `internal/canary` nested-run trap; read
  them before touching either, and note that `dist/bench` must be built with
  `scripts/go-build.sh`.
- **FT91 was raised to HIGH by the reviewer** (`9cbb138`) — the gate's length is
  the dominant cost of small changes here; this session paid ~13 minutes for one
  kit edit. A fourth arm is recorded on the row: `RunConformance` runs its
  fifteen checks serially, ~94% of wall clock. First step is timing each check.
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
