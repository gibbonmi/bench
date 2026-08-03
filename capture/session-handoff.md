# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `6a3ea99` before the pending spec-and-capture commit; FT183 implemented and committed
Spec: `specs/check-level-conformance-scoping/spec.md` (Status: staged, current commit batch), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `6a3ea99` (both FT183 commits landed through `bench commit`'s full gate)

## State

- **FT183 is done.** Ticket 1 (`b41b4d2`) retired the whole-changeset reduced gate path: gate/verdict/status/prep-release code and tests, the runtime and prep-release contract tests rewritten to the full-run and invalid-cache-record expectations, and `projects/benchkit.md`'s reduced-run prose replaced (ReducedScope, stripped-worktree enforcement, and status softening survive; all conformance anchors kept). Ticket 2 (`6a3ea99`) added `internal/gate/component_inputs_identity_test.go`: the Source → function identity check over all registry rows including the hand-declared canary row, per-row exhaustiveness refusal, and the method-expression/pointer-identity guard; swap A, swap B, and the exhaustiveness red were each demonstrated and reverted.
- **The check-level-conformance-scoping spec is staged and included in the current commit batch.** FT183 was its sequencing blocker and both implementation tickets have landed.
- `capture/IDEAS.md` carries three parked ideas; `capture/learnings.md` has one open entry (set-aside dance) for the next `/bench-what-next` drain.
- Decisions that stay closed: everything in the four maps landed at `449eb2a`, including the 2026-08-03 amendments.
- Pushed through `6a3ea99` (reviewer, 2026-08-03).

## Next command

`$bench-implement-spec --full specs/check-level-conformance-scoping/spec.md`

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
