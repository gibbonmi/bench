# Bench handoff

This is a pickup note for the next kit-development session. Durable product facts
live in `README.md`, `AGENTS.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`,
`CONTEXT.md`, and `projects/benchkit.md`; do not treat this file as a second
inventory of commands, skills, hooks, or CLI surfaces.

## Current state

Active phase: `/bench-implement-spec` on the what-next spec. The review/spec-guidance
learnings slice and the merged-spec retirement slice are complete and committed locally.
The working tree is clean, and `bench gate` was green before each commit.

The branch has local commits not yet pushed. Recent local work includes:

- `spec-retire: implemented specs`
- `learnings: promote review spec guidance`
- `Promote line governance learnings`
- `Clean up assessment findings`
- handoff refresh commits

`git push` was attempted after the spec-retirement commit and blocked by the
repo's PreToolUse hook: pushing is reviewer-owned authority. Do not bypass that
hook; the reviewer or an approved push path must publish the local commits.

## Committed work

Retired in `c76459a`:

- The 16 implemented specs in `specs/` were deleted under the
  `spec-retire: implemented specs` commit.
- `ROADMAP.md` no longer points at the retired spec files as live paths.

Promoted in `cdd8826`:

- `craft-synthesis` now requires a fresh-session dogfood run when a candidate
  changes skill or command triggers.
- `craft-delegate` now says read-only review findings are verified and fixed by
  the invoking session in the checkout that owns the diff.
- `/bench-review-implementation` now falls back to inline axes when a harness
  forbids unsolicited sub-agents.
- `/bench-write-spec` now checks external format/library divergence and runnable
  byte/wire compatibility probes.
- `/bench-shape-idea` now requires Research assets that claim byte or wire
  compatibility to include a runnable probe.
- `projects/benchkit.md` now treats real CLI, linked by-path CLI, hooks, and
  adapters as explicit hostile-input invocation surfaces.

Gate/conformance fix in `cdd8826`:

- `.bench/BENCH.md` is now the canonical CLI inventory source.
- `.bench/BENCH-reference.md` points back to that inventory instead of carrying a
  second command list.
- `HANDOFF.md` is no longer checked for inventory completeness, but command names
  it mentions are still checked for unknown/stale CLI references.
- The stale CLI canary now plants its stale reference in `.bench/BENCH.md`, and a
  new missing-inventory canary proves the BENCH.md inventory check still bites.

Recorded/pruned in `cdd8826`:

- `CHANGELOG.md` has the 2026-07-04 review/spec guidance learnings entry.
- `.bench/learnings.md` now has 2 open entries.

## Verification already run

- `git diff --check`
- `go test -count=1 ./internal/conformance -run '^TestDocsCurrencyTokenDietAndWorkflowFixturesBite$'`
- `go test -count=1 ./internal/conformance -run 'Fixture|Registry'`
- `BENCH_CONFORMANCE_ROOT=/home/devuser/workspace/bench go test -count=1 ./internal/conformance -run '^TestRootConformance$'`
- `go test -count=1 ./...`
- `bench gate`
- Retired-spec reference sweep over current docs after deleting `specs/*.md`
- `bench status` after `c76459a`: specs row cleared

## Remaining work

`bench learnings` reports 2 open entries, both needing product decisions before
implementation:

- Session-start stale gate: decide how `bench status` classifies benign drift vs
  real untrusted code drift.
- Review findings persistence: decide the artifact location and lifecycle before
  changing `/bench-review-implementation`.

Other `bench status` items:

- `structure` reports 7 oversized Go test files. This needs a `craft-seams` split
  plan, not a blind line-count cleanup.
- 14 roadmap ideas are parked; they are capture-only until shaped.

Recommended next action: publish the four local commits through an approved push
path. After that, continue with one decision at a time: first choose the review
findings artifact location/lifecycle or the stale-gate benign/real split, then
shape the structure split plan.
