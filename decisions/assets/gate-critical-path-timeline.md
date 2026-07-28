# Gate critical-path timeline (2026-07-28)

## Post-stage-2 measurement (2026-07-28, ft91-canary-compiled-bites built)

Same box (16 cores) and method, warm caches, idle box, on the worktree
carrying the full stage-2 build: behavior-owned fixtures graded by compiled
contract test binaries, no nested gates. A vacuity baseline carries no phase
pin, matching what shipped post-stage-1. An earlier post-stage-2 pass
measured a build that silently narrowed every non-contract vacuity baseline
to a single gate phase; those baselines were cheaper than they should be, and
restoring correct behavior gave back roughly 6 s of gate and about 1 s of
canary against what that pass reported.

- Solo canary **152.6 s** (post-stage-1: 151.5 s). **No meaningful change.**
- Full dev gate **167.3 s** mean over three consecutive green runs
  (169.3 / 166.0 / 166.6 s), against 172 s post-stage-1 — roughly 5 s better,
  and the wall is now bound by the contract phase alone rather than jointly by
  contract and canary.
- Contract per-package (`ok` lines), carried over from the pre-fix runs and
  not re-measured; the vacuity-baseline defect does not touch the contract
  phase: artifact **140.9 s**, runtime 42.8 s, surface 36.1 s,
  publication 16.5 s, axi 4.1 s, preprelease 4.0 s, contract 0.4 s.
- All 33 behavior-owned fixtures still bite non-vacuously under the compiled
  shape with **zero fixture edits** — the migration-casualty class the
  decision map named (an EXPECT observable only in gate-level framing) did
  not materialize.

**Why the canary did not move.** Stage 1 had already scoped each
behavior-owned fixture's inner gate to the single contract package owning its
EXPECT. What stage 2 removes on top of that is per-fixture gate startup and
the per-fixture `go test` compile — on the order of a second each, and
absorbed by the 8-worker pool. The canary's wall is set by the longest single
fixture, and that fixture still runs the `surface/artifact` contract suite in
full, because the compiled binary runs every `Test*` in the package. The
artifact multiplier is untouched by design: the spec's own "Out of scope"
section names the prepared-artifact hoist (map #2/#3) as the next slice and
the thing that attacks this tail. Until it lands, this arm cannot move the
canary wall.

**Mechanism confirmed, not assumed.** Process sampling during a live sweep
showed roughly ten concurrent processes executing out of the sweep-owned
`bench-canary-bin-*` temporary directory (including the encoded
`surface_sartifact.test` binary), against a single `gate.sh` — the one
non-behavior-owned family that still spawns an inner gate correctly. The
sweep is genuinely taking the compiled-bite path; the flat canary number is
not a silent fallback to nested gates.

**Open observation — one unreproduced red.** The very first full gate run
against this build came back red on the `test` phase; the canary and contract
phases were green in that same run. Its output was not captured, so the
failing package is unknown. Seven subsequent full gate runs were green —
three before the vacuity-baseline fix, four after — and
`go test ./internal/...` run separately was clean. The red is therefore
unreproduced and unattributed: it is not known to be caused by this diff, and
it is equally not known to be unrelated.

## Post-stage-1 measurement (2026-07-28, ft91-canary-contract-scoping landed)

Same box and method, main at the scoping commits (`bbfe233` + `f10faae`):
instrumented `bench gate` then solo `bench canary .`, warm caches.

- Gate wall **172 s** (was 267 s), all phases green. Contract phase ended
  t+143 (artifact 133.8 s of it); canary summary at t+172 — the gate is now
  jointly contract/canary-bound instead of canary-bound by 105 s.
- Solo canary **151.5 s** (was 250 s). Against the spec's ≤100 s acceptance
  this is a **miss**: real reduction, not the degenerate-migration signature,
  but above threshold.
- Why above the 60–80 s floor estimate: the estimate keyed on the 45 s
  surface straggler, but the five `surface/artifact`-bound fixtures each pay
  the artifact suite (~134 s at full width, more at inner width 2) — the
  artifact suite is now the canary tail exactly as it is the contract-phase
  tail. Stage 2 (in-process bite checks) and the prepared-artifact hoist
  (map #3) both attack that same tail.

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
