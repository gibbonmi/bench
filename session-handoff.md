# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `4a756b3`, 3 dirty paths, 0 unpushed commits
Spec: none staged.
Gate: red at `842b3c0` — current

## State

- Uncommitted here: a reviewer-approved `/bench-what-next` drain. Ten ideas and
  both implementation retros are drained, `IDEAS.md` and `.bench/retros/` are
  empty, FT128 is reconciled out as shipped, and FT173/FT174/FT175 are new rows.
  The reconcile was sampled, not exhaustive — roughly a third of the 49 rows were
  checked against the tree this run.
- `TestHandoffShapeSingleSourced` was red on `4a756b3` before that drain existed,
  reproduced at HEAD in a clean worktree. This file's regenerated `## Shape` is the
  fix, and it rides the drain commit because `bench commit` refuses beside
  unrelated dirty paths.
- The same commit deletes `"FT128": 1` from the occurrence-ledger migration map in
  `internal/conformance/docs_workflow_checks_test.go`, with reviewer approval. That
  map pins the counts a one-time migration produced and is designed to shrink as
  rows retire — its own bite test fails if a retired FT stays required. Reconciling
  a shipped row out of `ROADMAP.md` reds the gate until its entry goes too, and no
  phase instruction covers that today.
- FT168 was reviewer-priced LOW → MEDIUM and now owns the wider question of whether
  the oracle may answer for less than the whole tree.
- FT135's shaping finished in a concurrent session and lands in this same commit:
  `decisions/pre-push-guard-visibility.md` is `Status: ready`, with its research
  asset beside it. FT135 is therefore past shaping — the sequence routes it to
  `/bench-write-spec`, not another grill.
- The canonical full gate was last green on `5ae1540`. No prospective or
  post-commit gate ran for `5e0c347`; `refs/bench/green/main` remains at the last
  lifecycle-confirmed base rather than claiming the final tree.
- Both merged specs are now safe to retire — their retros are drained. Six
  `recovered` worktree assignments with missing trees remain in the pool from the
  cadence build.
- Nothing has been pushed; main is well ahead of `origin/main`.

## Next command

`/bench-shape-idea` — FT135, the top line of the refreshed `## Recommended sequence`.

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
