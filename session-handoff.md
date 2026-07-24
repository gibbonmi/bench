# Session handoff

Repository: /home/mgibs/workspace/bench — branch `main`. Everything below is
executable from a cold start; no conversation history is needed.

## State

- **FT91's first arm is spec'd and approved; nothing has landed in code.**
  `specs/gate-concurrency-budget.md` is `Status: staged`, reviewer-approved,
  committed at `c1ac60a` along with the ROADMAP pointer. Compiled from the
  closed `decisions/gate-concurrency.md`, whose Handoff carries the seams.
- **What the spec locks.** Every inner gate gets an explicit `GOMAXPROCS=k` in
  its env, stripped-then-set so an inherited value cannot leak past the cap.
  Canary workers become `fixtureWorkers(runtime.GOMAXPROCS(0), len(fixtures))`
  — divide by k, floor at 1, cap at the fixture count. `k = 2` lives in
  `internal/bounds` as `CanaryInnerWidth`, consumed only by
  `internal/canary/canary.go` and enforced by `checkBoundsPolicy`. No
  Bench-specific knob: `GOMAXPROCS=8 bench gate` is the escape hatch.
  `internal/gate` is untouched.
- **The one map deviation, approved at sign-off.** Handoff item 4 says the
  existing overlap and error-order canary tests stay unchanged. They cannot.
  The conformance phase runs the kit's own `go test` over core packages as a
  subprocess inheriting the phase env, so inside a fixture's inner gate the
  canary tests run at `GOMAXPROCS=2`, where the derived bound is 1 — and both
  tests need two workers to make progress. All three concurrency tests
  (`TestSweepRunsFixturesConcurrently`,
  `TestSweepReportsErrorsInSortedFixtureOrder`,
  `TestSweepBoundsFixtureConcurrencyAtNumCPU`) retarget to the derived bound;
  the two needing overlap gate on it through `capability.CPU`. This is story 7,
  and it is the highest-risk part of the build.
- **Two closed veto points from the map** (batch approval): k = 2 (k = 1 tied
  on wall at 330 s but doubles concurrent inner gates; k = 4 cost +30%), and
  the outer conformance/contract phases stay uncapped this arm.
- **Ship evidence is manual.** Story 1's payoff — full gate wall and load
  against the 10–15 min / load ~123 baseline of 2026-07-22 — is not
  gate-assertable. Measure after the build and compare to the map's k table.
  The prototype measured 332 s at k = 2; the conformance phase (319 s) is now
  the long pole.
- **Build gotcha.** A plain `go build -o dist/bench ./cmd/bench` stamps
  `version=dev` and fails two `internal/contract/surface` contracts.
  Hand-running that package needs
  `go build -ldflags "-X main.version=0.2.0" -o dist/bench ./cmd/bench`.
- **`bench status` flags to ignore or resolve in passing:** one roadmap row
  names a retired spec (a `/bench-what-next` reconcile is the backstop), and
  the gate verdict is stale for the current tree (docs-only commits since).
- **Unpushed:** `main` is well ahead of origin. Pushing is the reviewer's call.

## Next command

`/bench-implement-spec specs/gate-concurrency-budget.md` in a fresh mid-tier
session. `bench shift` is wrong here: the spec fails `craft-line`'s
venue-routing test, because story 1's outcome is not gate-observable and
stories 1 and 6 are routed mid, not cheap. The build is small — one constant in
`internal/bounds`, two functions in `internal/canary`, three registrations in
the bounds-policy conformance check, one new canary mutation fixture, and the
retargeted concurrency tests.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
