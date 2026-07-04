# Bench handoff

This is a pickup note for the next kit-development session. Durable product facts
live in `README.md`, `AGENTS.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`,
`CONTEXT.md`, and `projects/benchkit.md`; do not treat this file as a second
inventory of commands, skills, hooks, or CLI surfaces.

## Current state

Active phase: `$bench-integrate-learnings`. The current green slice is ready to
commit: six review/spec-guidance learnings were promoted and pruned, and the gate
policy mismatch around CLI inventory was fixed.

The review/spec-guidance slice is committed as `cdd8826`. The merged-spec
retirement slice is now staged in the dirty tree and still needs verification,
commit, and push.

## Committed work

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

## Current diff

- The 16 implemented specs in `specs/` are staged for promote-then-delete
  retirement.
- `ROADMAP.md` no longer points at the retired spec files as live paths.
- This handoff reflects that the review/spec-guidance slice is committed and the
  spec-retirement slice is in progress.

## Verification already run

- `git diff --check`
- `go test -count=1 ./internal/conformance -run '^TestDocsCurrencyTokenDietAndWorkflowFixturesBite$'`
- `go test -count=1 ./internal/conformance -run 'Fixture|Registry'`
- `BENCH_CONFORMANCE_ROOT=/home/devuser/workspace/bench go test -count=1 ./internal/conformance -run '^TestRootConformance$'`
- `go test -count=1 ./...`
- `bench gate`

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
- `specs` should clear after the staged retirement commit.
- 14 roadmap ideas are parked; they are capture-only until shaped.

Recommended next action: verify the staged spec-retirement slice with `bench gate`,
commit it as `spec-retire: implemented specs`, then push all local commits. After
that, the only remaining status items should be the two open learnings that need
product decisions and the structure split plan.
