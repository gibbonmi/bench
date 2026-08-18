# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `a2914fd5`, working tree clean, 2 commits unpushed
Spec: none staged; the light-path ticket at
`specs/light-path-live-root-conformance/tickets/grade-the-live-root-inside-the-dev-gate.md`
is implemented and landed

## State

The 2026-08 capability audit is reconciled and its worktrees are retired.
Results live at `docs/audits/2026-08-bench-capability/results-fable-high/`;
`final-reconciliation.md` (§A, §N, §O) and `next-ticket.md` are the durable
reads. The three audit worktrees (`bench-audit-fable`, `-opus`, `-sol`) and
their `audit/*` branches are gone, discarded through `bench worktree clean`
with recovery refs retained under `refs/bench/recovery/`. One unlanded delta
died with them: the fable tree also turned `design`, `artifact-diagramming`,
and `code-review` off in `.claude/settings.json`'s `skillOverrides`. Not
landed — reviewer's call whether those three should be hidden in this repo.

Audit decisions that stay closed: incremental strategy; no work-state store,
compiler, or claim graph (A12); `/bench` is an adapter over
`bench status --route`; FT100 waits for the measurement harness (A11).

**Audit action A1 is landed** (`a2914fd5`). `bench gate`'s ordinary test phase
now carries the graded root and the dev tier to `TestRootConformance`, so the
conformance registry's 29 checks grade the live tree on every gate rather than
only under `prep-release`. An environment-class skip observed by the oracle is
red, not a footer count, and every skip line names the emitting test. Ten
diagnostics were red at HEAD behind a green gate; all ten are disposed, plus an
eleventh the fix exposed (the `BENCH-reference.md` anchor needle itself carried
the removed `spec build` token its own sweep forbids).

Three calls in that commit are post-hoc veto surface:

- The `.agents/commands/bench-implement-spec.md` prose budget went 60 → 75 in
  `projects/benchkit.md`. The file was at exactly 60 and the restored
  `## Entry orientation` / `## Exit handoff` headings do not fit; 75 matches its
  sibling `bench-write-spec.md` at 73. The alternative was cutting prose the
  `fa4e1f02` slimming had already chosen to keep.
- `bounds.CanaryInnerWidth` is retired rather than rehomed. Its only consumer,
  `internal/canary/canary.go`, was deleted by `3701c4a0` with the parallel
  fixture sweep, and nothing replaced the arithmetic.
- `decisions/spec-build-review-gate-cadence.md`'s two dangling `Sources` rows
  now point at `CHANGELOG.md` and `internal/landing/landing.go`. The larger
  question is untouched and open: that map is `Status: shaping` for a
  provisional spec-build lifecycle the tree deleted wholesale, so whether it is
  still live work belongs in a `/bench-what-next` drain.

Next by the audit's own ranking is **A2** (the verdict reader — changes what the
dashboard reports), then **A3** (`/bench` router), both in
`results-fable-high/action-items.yaml`. A1 was sequenced first so A3's new phase
file, rename, and adapters land graded rather than unenforced.

Standing repo condition, not this session's doing: `bench status` reports the
`pre-push` hook missing (`bench link` installs it) and 62 `bench structure`
issues.

## Next command

`/bench-review-implementation`

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
