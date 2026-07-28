# Gate critical-path timeline (2026-07-28)

Instrumented full dev gate plus a solo canary run, kit repo, 16 cores, idle
box, warm caches. Method: `bench gate` piped through a per-line epoch
timestamper; a 5 s sampler recording `/proc/loadavg` and the top processes by
CPU; then `bench canary .` solo under the same sampler. Raw logs were
session-scratch; every number they support is inlined here.

## Headline

The dev gate is canary-bound end to end: gate wall 267 s, canary solo 250 s.
The contract phase finished at t+162 s — 105 s before the gate — so the
seventh arm's 131 s contract-suite win could only move the wall by the
~24 s it did. Canary was underneath the whole time.

## Phase timeline (gate run, wall 267 s)

Phase output spans (first→last prefixed line; all `phase X: green` summaries
print together at t+267):

| phase             | first output | last output |
|-------------------|--------------|-------------|
| contract          | t+0          | t+162       |
| test              | t+0          | t+64        |
| race              | t+0          | t+7         |
| conformance       | t+13         | t+13        |
| conformance-suite | t+30         | t+30        |
| canary            | t+267 (buffered single summary line) | t+267 |

Contract per-package (`ok` lines): artifact **158.2 s**, runtime 54.5 s,
surface 42.2 s, publication 21.8 s, axi 5.6 s, preprelease 5.3 s,
contract 0.5 s. Artifact is 158 s of the phase's 162 s.

Load: peak 23.4 at t+49, ~15–22 through t+120, decaying 10→3 over
t+160→267 — the tail is canary draining, nothing else left running.

## Canary decomposition (solo run, 250 s)

Sampler process mix during solo canary: 8 workers saturated by
`runtime.test` and `surface.test` — the kit's own contract test binaries —
for essentially the full 250 s, ending in a single `surface.test` straggler
from t+215 to t+246 (~45 s at inner width 2, `bounds.CanaryInnerWidth`).
Conformance-scoped fixture runs (`conformance.test`) barely register
(2 samples).

Accounting: 8 workers × 250 s ≈ 2,000 worker-seconds. The 34
`behavior-owned` fixtures at ~50 s of full contract suite each ≈ 1,700 of
them. The other ~123 fixtures (conformance-scoped since the fifth arm, plus
one-fixture toolchain families) are the remainder.

## Mechanism

- `FixturePhase` maps family `behavior-owned` → phase `contract`
  (`internal/canary/canary.go:36`), and the contract phase argv is
  `go test -count=1 ./internal/contract/...` (`internal/gate/phases.go:96`,
  `goTestArgv` at `phases.go:187`) — so each of the 34 fixtures' inner gates
  re-runs all seven contract packages, `-count=1` defeating the test cache,
  at `GOMAXPROCS=2` per the concurrency arm's pin.
- The suites are subprocess-heavy integration tests (they exec `dist/bench`,
  `git`, `npm pack`), so the cost is process trees, not Go compute.
- Artifact-suite inflation under the gate (106 s solo → 158 s in-gate) is CPU
  oversubscription: load 15–23 on 16 cores while contract, test, and canary's
  nested runs overlap through t+120. No lock or disk signal observed.

## Floor arithmetic

- Today: wall = canary ≈ 250–267 s regardless of contract.
- After stage 1 (package-scoped behavior-owned nesting): canary ≈ 60–80 s
  (floor: slowest bound package at width 2, cf. the 45 s surface straggler);
  wall becomes the contract phase ≈ 162 s.
- After the prepared-artifact hoist (map #3): contract ≈ 30–50 s; after
  stage 2 (in-process bite checks): canary ≈ 20–40 s. Wall ≈ max(test 64 s,
  contract, canary) ≈ **60–90 s** with oracle semantics untouched.
- Below ~60 s requires oracle-semantics changes (verdict caching, `-count=1`)
  — reviewer-owned, pulled into FT91 by the ≤60 s stop ruling.
