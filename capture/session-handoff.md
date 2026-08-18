# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `f1135fd6`, working tree carries the reviewed `/bench-what-next` drain batch, about to land as one commit
Spec: none staged; `specs/` is empty
Gate: green at `f1135fd6` (fresh, via the landing `bench commit` itself)

## State

`/bench-debug` found and fixed a gate defect the drain's own landing attempt
surfaced: `bench commit`'s prospective authorization runs `TestRootConformance`
against an ephemeral `git worktree add` + `read-tree` checkout that never
carries `dist/` (gitignored), and `package-core-guard`'s `checkPackageFiles`
demanded it exist unconditionally — a check meant for a built release payload,
newly reachable on every ordinary commit since A1 wired the conformance
registry into the dev-gate test phase. Every `bench commit` was red for this
reason, independent of diff content. Fixed in `f1135fd6`
(`internal/conformance/package_core_checks_test.go`,
`internal/conformance/fixture_bite_test.go`): a missing, gitignored `files[]`
entry is exempt below ship tier (using the `tier` parameter
`checkPackageCoreAndGuards` already threaded but never used), ship tier keeps
the strict check, and a genuinely missing tracked entry still flags
everywhere. New regression test:
`TestCheckPackageFilesExemptsGitignoredEntryBelowShipTier`. Built in an
isolated `bench worktree create` assignment, gated green there, then
fast-forwarded onto `main` (the assignment predates a spec and
`bench worktree land` requires one, so the merge was a manual
`git merge --ff-only`, not a Bench-owned landing verb — worth a roadmap row if
this recurs). Assignment released and removed.

`/bench-what-next` separately reconciled the roadmap against the tree
(nothing shipped since the last drain beyond the already-annotated FT120/A1
work; `specs/` and `spec_history` are both empty, so no `bench spec retire`
was owed) and drained the one open capture source:

- `capture/IDEAS.md` was already empty; `capture/retros/` was already empty.
- `capture/learnings.md`'s one open entry (`bench commit` refused twice with
  `prospective authorization refused: inherited`, printing a misleading
  `gate: red`) reproduces FT6's own 2026-08-12 graduation trigger — a second
  refusal through `bench commit` itself — so it graduated out of FT6's parked
  tier into a new decision-required row, `roadmap/FT223.md`: reviewer chooses
  between rewording the refusal in operator terms (name `bench gate --fresh`,
  stop printing `gate: red` for a partial verdict) or having `bench commit`
  escalate to a fresh prospective gate itself instead of refusing.
- Caught in the same pass: `## Reds the diff doesn't own` had drifted to
  "Four rows" against five listed rows before FT223 landed; corrected, and
  the section is genuinely five rows again with FT223 added.

`## Recommended sequence` is unchanged — FT223 is LOW severity and does not
outrank the existing top three (FT100, FT207, FT213).

Standing repo condition, not this pass's doing: `bench status` reports the
`pre-push` hook missing (`bench link` installs it) and 62 `bench structure`
issues.

## Next command

`/bench-final-check` — the board's leading invocable signal (`git`), once this
batch and the prior unpushed commits are ready to land together.

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
