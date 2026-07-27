# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `fc99aee`, 8 dirty paths, 5 unpushed commits
Spec: `specs/ft91-canary-check-scoping.md` — implemented, gate-green, unpushed.
Gate: green at `fc99aee`

## State

- **ft91-canary-check-scoping is built and landed on `main`, awaiting your
  merge/push** — five commits `4f046f0`, `f9bb33d`, `95981b9`, `31ab979`,
  `fc99aee`. All seven stories built; coverage map valid at 17 rows; the gate
  ran green on every commit. Three review axes produced 12 findings: Standards
  and Coverage findings are all fixed, one Standards finding was rejected as
  closed by the spec, and the Spec finding is the amendment below.
- **Story 4's seam was amended during the build — open to your veto.** The
  spec placed the unbound-family red in the sweep, during fixture selection.
  Built that way it reddened every adopting repo, because `bench init`
  scaffolds a seed canary family a kit-owned table can never bind. The sweep
  now resolves an unbound family to no scope, and the conformance layer's
  kit-scoped family check raises the red. The spec records the amendment;
  `.bench/learnings.md` journals the deviation.
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

`/bench-what-next` — the board's leading invocable signal (`drain`): two parked
ideas and one open learning, all from this build. Push `main` first if you
accept story 4's amended seam.

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
