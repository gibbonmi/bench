# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT91's first arm shipped, was reviewed, drained, and retired. Nothing from it
  is open.** Four commits: `4ea4880` (the build) and `93eb291` (a capture) are
  pushed; `0d5d935` (the drain) and `7c1ea2d` (the retirement) are **unpushed**.
  Pushing is the reviewer's call.
- **The payoff, for the record:** the full gate went from 10–15 minutes at load
  ~123 to **323–336 s at peak load ~28–33** on 16 cores. The conformance phase
  (~325 s) is now the long pole. FT91's roadmap row carries the detail.
- **`specs/` is empty — no staged spec.** The next work starts at
  `/bench-write-spec`, not at an implement phase. `decisions/gate-concurrency.md`
  was deleted with the retirement; its durable content is in
  `projects/benchkit.md` (the canary concurrency budget in the gate section, and
  the nested-run trap in the cold-session notes). Read those before touching
  `internal/canary` — the nested trap deadlocks rather than fails, and this file
  deliberately does not restate it.
- **The drain added two rows.** FT120 collects two gate/canary test-harness
  defects (the R12 fixture's immortal spinner, and the unasserted per-fixture
  `BENCH_CANARY_PHASE` isolation) and carries the
  `internal/canary/canary_concurrency_test.go` split — the FT91 arm pushed it to
  401 lines against a 400 budget, which is the one `bench structure` violation
  this work created. FT121 is the `bench commit --spec` gate-staleness defect.
  Both want a spec, not a direct fix; FT120's reason is in the row.
- **No open learnings, no parked ideas.** Both capture sources are empty, so
  `bench roadmap` prints its recommended sequence.
- **A concurrent session was writing this repo on 2026-07-24** and edited this
  file mid-gate, which the working-tree tripwire correctly caught. Check
  `bench status` and `ps` for other live writers before starting a gate.

## Next command

`/bench-write-spec` — FT86, fail-closed control records and single-sourced
repository facts. The highest bank-track row still open and the top of the
roadmap's recommended sequence. FT71 is next after it, then FT120.

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
