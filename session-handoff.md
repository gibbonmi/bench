# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT91's first arm is mapped, measured, and spec-ready.**
  `decisions/gate-concurrency.md` is closed with a full Handoff: product
  concurrency budget (inner gates pinned to `GOMAXPROCS=k`, canary workers =
  `max(1, GOMAXPROCS(0)/k)`), k = 2 picked from live measurement, no
  bench-specific knob. Single spec, seams and assertables enumerated in the
  map's Handoff section.
- **Two veto points, chosen under the reviewer's standing batch approval.**
  k = 2 (k = 1 tied on wall-clock at 330 s but doubles the concurrent inner
  gates; k = 4 cost +30%); and outer conformance/contract phases stay
  uncapped this arm. Both are recorded in the map with the measurement table.
- **Measured payoff:** full gate 332 s at k = 2 versus the 10–15 min / load
  ~123 baseline — on the prototype worktree. Nothing has landed in code;
  the throwaway worktree is released. The conformance phase (319 s nested kit
  test suite) is now the gate's long pole.
- **FT91's roadmap row stays** until the arm ships; the remaining arms
  (test-result caching, hermetic build cache) are out of this map's scope.
- **`bench status` flags to ignore or resolve in passing:** one roadmap row
  names a retired spec (a `/bench-what-next` reconcile is the backstop), and
  the gate verdict is stale for the current tree (docs-only commits since).
- **Build gotcha.** A plain `go build -o dist/bench ./cmd/bench` stamps
  `version=dev` and fails two `internal/contract/surface` contracts.
  Hand-running that package needs
  `go build -ldflags "-X main.version=0.2.0" -o dist/bench ./cmd/bench`.
- **Unpushed:** `main` is well ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-write-spec` in a fresh mid-tier session, seeded with
`decisions/gate-concurrency.md` — the map's Handoff carries the seams, so the
spec needs no interview. The build after it is small: one package
(`internal/canary`), one constant in `internal/bounds`, and a retargeted
concurrency-bound test.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
