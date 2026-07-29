# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `517e996`, 17 dirty paths, 2 unpushed commits
Spec: none staged.
Gate: green at `a0be126` — stale, work tree `c189925`

## State

- **FT154 `--full` reached implementation and stopped at the first ticket's
  cap.** The approved folder spec is staged at
  `specs/craft-tickets/spec.md` in green commit `517e996`; the old binary's
  flat-only `Facts` scan is why the generated pin block above says no spec is
  staged.

- **The partial ticket work is preserved at
  `/home/mgibs/.bench/worktrees/bench-2826441890/b6f4c57d40d54a8053bc12141783227a-4c89fbebe1fec7cbc54d1a3f5077bd76`.**
  `internal/spec`, status, coverage, the focused runtime spec/history tests,
  and the coverage-map canary are green. Runtime gate-proof, handoff,
  status/commit and AXI coverage/roadmap fixtures remain flat-form red; the
  retire interrupt-resume cases remain unwritten. No implementation commit or
  whole gate run exists.

- **The first ticket did not fit one fresh context.** It combined a wide
  folder-layout refactor with every consumer migration. Recommended repair:
  approve an expand–migrate–contract split so each replacement ticket can land
  green, while the final state still refuses live flat specs. Do not continue
  the existing atomic ticket without that reviewer ruling.

- **`IDEAS.md` is unrelated reviewer-owned capture.** It remains modified in
  the main checkout and was deliberately excluded from the FT154 staging
  commit.

## Next command

`$bench-implement-spec --full specs/craft-tickets/spec.md`

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
