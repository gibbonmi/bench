# Bench handoff

This is a pickup note for the next kit-development session. Durable product facts
live in `README.md`, `AGENTS.md`, `.bench/BENCH.md`, `.bench/BENCH-reference.md`,
`CONTEXT.md`, and `projects/benchkit.md`; do not treat this file as a second
inventory of commands, skills, hooks, or CLI surfaces.

## Current state

Active phase: `$bench-integrate-learnings`. The current green slice is ready to
commit: six review/spec-guidance learnings were promoted and pruned, and the gate
policy mismatch around CLI inventory was fixed.

`bench gate` is green on the dirty tree. The branch is still ahead of
`origin/main` by the two prior local commits until this slice is committed and
pushed.

## Current diff

Promoted in the dirty tree:

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

Gate/conformance fix in the dirty tree:

- `.bench/BENCH.md` is now the canonical CLI inventory source.
- `.bench/BENCH-reference.md` points back to that inventory instead of carrying a
  second command list.
- `HANDOFF.md` is no longer checked for inventory completeness, but command names
  it mentions are still checked for unknown/stale CLI references.
- The stale CLI canary now plants its stale reference in `.bench/BENCH.md`, and a
  new missing-inventory canary proves the BENCH.md inventory check still bites.

Recorded/pruned in the dirty tree:

- `CHANGELOG.md` has the 2026-07-04 review/spec guidance learnings entry.
- `.bench/learnings.md` now has 2 open entries.

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
- `specs` reports 16 merged specs awaiting promote-then-delete retirement.
- 14 roadmap ideas are parked; they are capture-only until shaped.

Recommended next action after committing this green slice: retire the merged specs
mechanically, then rerun `bench gate` before pushing all local commits.
