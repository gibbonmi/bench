# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — clean tree; ~60 unpushed commits
Spec: `specs/remove-spec-build-lifecycle/spec.md` (Status: implemented) —
Spec A of the Pocock-alignment program, done. Program map:
`specs/remove-spec-build-lifecycle/decisions/pocock-alignment.md` (all 13
tickets reviewer-resolved 2026-08-11; its answers are closed decisions).

## State

Spec A landed over five serial commit-on-green tickets plus one composed-review
repair round: the `bench spec build` lifecycle, `worktree recovery`, and all
preservation/provisional ref machinery are gone (~16k lines net deleted);
resume's reconcile sweeps `refs/bench/{specbuild,recovery}` idempotently while
the gate's `refs/bench/green` provably survives; `bench commit --spec` is the
sole `Status: implemented` author; a removed-verb sweep and removed-grammar
tests stand guard. Reviewer authorized direct-to-main commits without the gate
for this program; the terminal landing ran the real gate via `bench commit
--spec`.

Open reviewer decisions, none blocking:
- `bench worktree clean --apply` still preserves dirty work into
  `refs/bench/recovery/`, which the next resume sweeps — recommend a follow-up
  making explicit clean refuse a dirty removal instead.
- Accepted deviations flagged for veto: CHANGELOG.md exempt from the
  removed-verb sweep (append-only history); RM11's predicate corrected to
  family-help + per-verb grammar lines; RM2/RM9 predicates tightened to match
  the tree (`0a6fb1b4`).
- Residual: `projects/benchkit.md` ~311–317 still describes lifecycle cadence
  around a re-pinned needle; `capture/audits/injected-interface-composition.md`
  cites deleted symbols (audit record, left as history). Both fold into Spec C.

## Next command

`/bench-write-spec` for Spec B (`bench preflight`) from
`specs/remove-spec-build-lifecycle/decisions/pocock-alignment.md` ticket #7,
in a fresh session. Then Spec C (doctrine adoption, map #4/#5/#6/#8/#9/#10),
then one `/bench-what-next` drain over the finished tree (roadmap re-verdicts
incl. FT128/FT173/FT184/FT185, mooted learnings, parked ideas, retirement of
the `spec-build-review-gate-cadence` and `parallel-session-landings` maps, and
re-rank of the three parked specs).

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
