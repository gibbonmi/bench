# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `0f4a861`, 8 dirty paths, 7 unpushed commits
Spec: none staged.
Gate: green at `523e0ff` — stale, work tree `c611032`

## State

- **The `/bench-what-next` drain is committed at `0f4a861`; both capture
  sources are empty.** ft91-canary-check-scoping shipped, merged, and is now
  retired, so FT91 carries it as the fifth landed arm and the check-scoping
  prerequisite paragraph is gone. Two new rows: **FT143** (the family→check
  binding reds only after each unbound fixture pays a full unscoped inner gate,
  and only on the `bench gate` path — `bench canary` never reaches the
  conformance registry) and **FT144** (a `craft-spec` edge-inventory prompt for
  kit code with two audiences, plus your decision on an "intent stands, seam
  moves" lighter case in the workflow). Both are open to post-hoc veto.
- **`decisions/gate-pipeline.md` stays closed** — its Handoff carries the
  seams for slices B (manifest + DAG runner) and C (`checkGoCore` split +
  fixture migration), which spec now that A has shipped.
- **Decisions that stay closed:** baseline grouping key is the resolved check
  name alone (unscoped fixtures share today's single full baseline); the
  live sweep's did-not-bite verdict is the binding's enforcement; no
  fixture merging; the family→check table stays in
  `internal/conformance/registry` as the imported-by layer.
- Codex CLI note: `codex exec` must run with stdin closed (`</dev/null`) or
  it blocks reading the pipe forever — cost two dead attempts this session.
- `bench prep-release` stays shelved — blocked by FT116's race and FT142's
  ship-track findings; both are board rows, not handoff state.
- The branch/worktree sweep (23 non-`main` branches, 19 worktrees) remains
  proposed, not executed — reviewer's call.

## Next command

`/bench-shape-idea` — FT91's pipeline arm, the top of the refreshed
recommended sequence. Inputs: the FT91 row and `decisions/gate-concurrency.md`'s
watch-outs.

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
