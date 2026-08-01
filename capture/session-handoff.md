# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `6f3486a`, 9 unpushed commits
Spec: `specs/reduced-gate-phase-set/spec.md` (built and landed, Status still `staged`),
`specs/pre-push-guard-visibility/spec.md` (Status: staged, unstarted)

## State

- **`reduced-gate-phase-set` is built and landed at `6f3486a`** — gate green over the
  composed tree. All 26 acceptance-coverage rows plus five tests beyond the map. Its
  spec still reads `Status: staged`; nothing authored the flip, because the spec-build
  lifecycle was abandoned mid-run (below) and `promote` is what normally writes it.
  Flipping it is a one-line reviewer call, not pickup work.
- The retro is `capture/retros/reduced-gate-phase-set.md`. Read it before the drain — it
  carries the design findings, three process failures, and what was deliberately left
  behind.
- **A spec-build run for this slug is permanently stuck `active`.** Every one of the
  eight operations routes to `promote` for recomposition, and `promote` checks a clean
  review and released assignments *before* the recompose branch — conditions a
  mid-repair run cannot meet. `abandon` is blocked by the same gate. Six worktrees stay
  registered and its provisional refs retained; no sanctioned operation retires them.
  Parked as a PRIORITIZE idea. This is inert — it blocks only that slug's lifecycle, not
  ordinary commits.
- Also inert, also recorded: five recovery refs from the earlier cadence build that
  `bench worktree recovery` refuses because it cannot prove their payloads landed. The
  refusal is correct.
- `capture/IDEAS.md` carries five parked items, three of them marked PRIORITIZE because
  they fire on every build rather than once. `capture/learnings.md` carries four
  entries, two added this session. All await `/bench-what-next`.
- One follow-up is scoped and ready but unbuilt: `checkContractCaptureReads` and
  `checkScopeBinding` are reachable only through `conformance-suite`'s whole-package run,
  by test-name prefix rather than registry registration. Registering them in
  `registry.Checks` + `conformanceChecks` makes them run in the `conformance` phase by
  construction and inherits the existing unbound-row protection. Two files, light path.
  The reviewer call is whether that composes the registry seam or crosses the
  gate-contract seam.
- Nothing has been pushed; `main` is nine commits ahead of `origin/main`.

## Next command

`/bench-what-next`

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
