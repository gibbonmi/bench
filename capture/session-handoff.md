# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `843a4b7`, 6 dirty paths, 0 unpushed commits
Spec: `specs/check-level-conformance-scoping/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `62e5458` — stale, work tree `7758ce4`

## State

- **The four-map decision session is closed and dual-reviewed** (opus +
  codex gpt-5.6-sol): `decisions/ft156-anchor-registry.md`,
  `decisions/ft144-post-approval-edits.md`,
  `decisions/ft181-spec-build-preconditions.md`, and
  `decisions/ft183-gate-scoping-residuals.md` are all `Status: ready`, with
  review findings folded in as flagged corrections and four reviewer
  amendments (FT181 #2 authorizes the `internal/worktree` abandon-planner
  change; FT181 #4's mechanism text corrected to the `errRecompose`
  refusal; FT144's rule is now intent-based with a run-boundary timing rule;
  FT183 #3 adds refuse-unknown-rows exhaustiveness and binds the canary row).
- **The maps and `decisions/assets/ft183-derivation-binding.md` are
  UNCOMMITTED** — `bench commit` refuses while the other session's untracked
  `specs/check-level-conformance-scoping/` sits outside the named set. After
  that spec lands (or leaves the tree), commit with:
  `bench commit -m "shape-idea: close the four-map decision session; FT183 ready after research and dual review" decisions/ft144-post-approval-edits.md decisions/ft156-anchor-registry.md decisions/ft181-spec-build-preconditions.md decisions/ft183-gate-scoping-residuals.md decisions/assets/ft183-derivation-binding.md`
- Decisions that stay closed: all rulings in the four maps, including the
  2026-08-02 grill answers and the 2026-08-03 amendments above; the
  stronger-than-substring anchor mechanism stays deferred (FT156 Out of
  scope).
- FT156's roadmap-row anchor sizing is stale against the instrumented count
  (299 needles, ~222 require-direction); the map carries the correction.
- Nothing pushed; push is the reviewer's call.

## Next command

`/bench-write-spec` — start with `decisions/ft181-spec-build-preconditions.md`
(all four faces ruled, discretion bounded, no research dependencies); FT183
and FT156 follow, each from its ready map.

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
