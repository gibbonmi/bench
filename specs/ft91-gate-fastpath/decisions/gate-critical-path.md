# Gate assessment — where the wall lives, and the levers left

Assessed 2026-07-29 at working tree of `187bc36`, re-measured on the landed
tree of `d91b709` once moves 1 and 2 shipped. Question: how fast can the dev
gate get without sacrificing quality, including the oracle-semantics levers
FT91 parked as reviewer decisions — verdict reuse for the exact tree, canary
narrowing, and freshness relaxation. This record incorporates the retired
FT91 timeline's historical measurements, the instrumented run that located
the wall, the post-build re-measurement, and the resulting decision analysis.

## Result

The changed-tree gate is **one package wide**: `internal/contract/surface/artifact`
(~109 s solo) floors the contract phase and with it the whole run. Everything
else in the gate finishes by ~t+62, and the canary sweep — formerly the tail,
because it re-ran that same package once per `surface/artifact` fixture — is
now 25 s.

Three moves take the gate from its ~135–172 s baseline; two are measured,
the third is an estimate:

| after | changed-tree gate | unchanged-tree gate |
|---|---|---|
| baseline | 135–172 s | 135–172 s |
| 1. per-test canary bites *(landed)* | **128 s** | same |
| 2. gate-level verdict reuse + freshness ruling *(landed)* | same | **0.6 s** |
| 3. artifact-suite restructure (parallelism) *(open)* | **~60–75 s** | 0.6 s |

Moves 1 and 2 landed in FT91's gate-fastpath build: move 1 on semantics the
reviewer had already ruled (stage 2's own wording), move 2 on an authority
ADR 0002 already granted `bench commit` and now extends to gate execution.
Move 3 is ordinary spec work with no semantic content, and is the one slice
left that moves the changed-tree wall. A full gate re-write is **not**
warranted — the scheduler is already maximally concurrent and the bound is
elsewhere (see "Re-design verdict").

## Where the wall was (2026-07-29, this box, warm caches, session load present)

The decomposition that located the wall, measured at `187bc36`. Its canary
figures are pre-lever-1 and are superseded by the post-build measurement
below; the contract-phase figure still stands.

Full `bash .bench/gate.sh`, per-line timestamps. Wall **172 s**, green.
The idle-box equivalent measured for the same subject family was 135.3 s
(timeline asset); the delta is background load, not a regression.

| phase | span | note |
|---|---|---|
| test | t+1 → t+62 | almost entirely `(cached)` — Go test cache already serves this phase |
| race | t+8 → t+9 | |
| conformance | t+15 | 12.6 s package; `package-core-guard` is 7.6 s of it |
| conformance-suite | t+35 | |
| contract | t+2 → t+140 | artifact **132.7 s**, runtime 52.1 s, surface 38.6 s, publication 20.1 s, axi 4.5 s |
| canary | summary at t+172 | the gate's tail; decomposition below |

Canary accounting (174 fixtures; 33 behavior-owned: 8 axi, 6 runtime,
10 surface, 5 surface/artifact, 4 surface/publication): each behavior-owned
fixture invokes its owning package's compiled test binary over the fixture
tree (`internal/canary/canary.go`, `RunBite`), so the sweep's worker-seconds
are ≈ 5×133 + 6×52 + 10×39 + 4×20 + 8×4.5 ≈ **1,480**, plus the ~140
conformance-scoped fixtures and baselines. At 8 workers that is ~185 s — which
is the observed wall. The floor is the longest single fixture: an artifact
fixture at `GOMAXPROCS=2`, ~110–133 s. The canary and the contract phase were
floored by the same suite; lever 1 broke the canary half of that.

## Post-build measurement (2026-07-29, `d91b709`, WSL2 host, warm caches)

Measured on the landed tree, in a `main` checkout whose own `dist/bench` was
rebuilt via `scripts/go-build.sh` first.

| what | command | wall | before |
|---|---|---|---|
| solo canary sweep | `time ./dist/bench canary .` | **25.2 s** (user 1m58.8s, sys 51.8s) | ~172 s |
| full gate, real run forced | `time bench gate --fresh` | **128.0 s** | ~135–172 s |
| unchanged-tree reuse | `time bench gate` immediately after a green | **0.57 s** | the full wall |

**Stop rule.** FT91's ≤60 s rule is **met on the canary critical path** the
rule was written against — the sweep is 25 s. It is **not met on the
changed-tree full gate**, which stands at 128 s, floored by the contract
phase's `internal/contract/surface/artifact` package (~109 s solo). The
unchanged-tree path beats its expected ~2–5 s by an order of magnitude: it
prints `gate: green (fresh verdict reused for this tree)` and does no phase
work.

**The toolchain closure declares `go` and `node`, not `npm`.** Every subject
hash walks each declared binary and its shebang chain, and the 0.57 s reuse
timing covers that whole hash-and-look-up path, so on warm caches the hashing
is an observed non-issue at this measurement. A cold cache or a slower
filesystem is untested. `npm` is deliberately outside the closure on the
reviewer's ruling, for two reasons: a host with node but no npm (the Debian
split packaging) would hold a permanently open subject, which disables reuse
and locks `bench prep-release`'s dev-green entry check with no operator
escape; and npm's shebang chain re-hashes `node`, doubling the closure's cost
for a signal already carried. Nothing is lost, because an npm upgrade ships
with its node — a toolchain upgrade still changes the subject through `node`.

## What already shipped (do not re-litigate)

- Concurrency budget (ADR 0009): 10–15 min → ~5.5 min.
- Dev/ship tier split, host-only dev generator, shared dev caches
  (`decisions/cost-follows-project-size.md` #3, #8–#10): removed the
  release-evidence and hermeticity cost from the per-commit path.
- Canary stage 1 (family/package scoping) + stage 2 (compiled bites, no
  nested gates): canary 250 s → ~152 s.
- Prepared-artifact hoist (FT91 #2/#3): artifact suite ~141 s → ~109 s.
- Commit-time verdict reuse for a fresh green on the identical closed subject
  (ADR 0002 posture 5).
- Per-test canary bites (lever 1): canary ~152 s → 25 s.
- Gate-level verdict reuse at a 60 min freshness window, over a subject
  closure that now declares `go` and `node` (lever 2): a re-judged unchanged
  tree → 0.6 s.

Closed rulings that stay closed: diff-scoped gating is unsound; no check is
weakened or dropped for wall-clock; canary skip-on-"inputs unchanged" key
rejected as fail-dangerous (ADR 0003); fixture batching rejected
(cross-contamination, ADR 0003).

## Lever 1 — per-test canary bites (landed; not oracle-weakening)

**What:** a behavior-owned fixture's bite run invokes its compiled package
binary with `-test.run` pinned to the one contract test whose failure message
its EXPECT is, instead of running every `Test*` in the package.

**Why it is not a semantics change:** stage 2's ruling (FT91 map #7) already
narrowed the claim — *"a behavior-owned fixture's bite is proven at its
owning contract test, not by a nested gate."* The EXPECT is, by the same
ruling, "the owning contract test's own failure message." Running the other
~25 tests in the package per fixture is incidental breadth left over from the
compile-once implementation, not part of the declared claim. Scoping execution
to the declared owner is the fifth-arm/#6/#7 precedent ("scoping a fixture to
the check its EXPECT names is not weakening") applied one level deeper.

**What it bought:** the five artifact fixtures no longer pay ~133 s each;
every behavior-owned fixture is binary startup plus one test (seconds; the
worst single tests do one real build, ~5–15 s). Measured at `d91b709`: the
solo sweep is 25.2 s wall against ~171 s of CPU, bounded by the
conformance-scoped remainder, the baselines, and the slowest single test —
and the gate is contract-phase-bound at 128 s.

**Fail postures (all loud, matching stage 1/2):**
- Binding: an explicit per-fixture `TEST` file naming the owning test,
  mirroring the existing `CHECK` convention — a name the package does not
  declare is a red at sweep start (same shape as `CHECK`'s unknown-check
  error), never a silent fallback to the full package.
- A renamed or deleted owning test → did-not-bite red.
- A contract group's vacuity baseline runs its package binary whole over the
  empty tree — no `-test.run`, no `-test.v` — so it stays deliberately wider
  than the scoped runs it grades, and no baseline carries a per-test twin.
  The two directions fail differently: a baseline narrower than its runs lets
  a vacuous EXPECT clear the screen in silence, while a wider one at worst
  calls a sound EXPECT vacuous, which is a red someone reads. This is the
  wide-baseline principle `scopeBaselines` already holds for phase pins, and
  it is the one place scoping stops.

**Reviewer ruling:** granted — per-test scope rides the existing scoping
precedent, aligning execution with semantics the map had already declared.

## Lever 2 — gate-level verdict reuse + freshness (landed)

**What exists:** every gate run already builds a *closed oracle subject* —
working-tree hash (untracked files included: `internal/git/tree.go:16` stages
`add -A` into a throwaway index), resolved gate, execution-policy version,
`PATH`, the gate script's full launcher closure, and every declared
environment/path/tool input from `.bench/gate-inputs.json`
(`internal/gate/subject.go:53`). The verdict is recorded durably under lock,
pending-before-run, subject-rechecked after. `Inspect` computes
`ReusableGreen` = ready + green + same tree + same oracle + inside the
freshness window (`internal/gate/verdict.go:243`).

**The shipped shape:** gate execution itself returns the reused green for a
`ReusableGreen` subject, printing `gate: green (fresh verdict reused for this
tree)`; `bench gate --fresh` forces a real run. The reuse check lives in the
gate home, and `bench commit` consults it there rather than carrying its own.
Because the shift loop, the stop hook, and `/bench-final-check` all route
through this package, the reuse covers every repeat-subject re-judging in one
place: final-check → commit → push chains, post-review re-gates on an
untouched tree, and an operator re-typing `bench gate`. Measured cost of a
reused verdict: 0.57 s.

**What "green" comes to mean:** "this exact closed subject was judged green
within the freshness window," rather than "the oracle just ran." ADR 0002
posture 5 already accepted exactly this claim for the commit path, with the
same reopen trigger (a reused verdict authorizing what a run would have
refused). Extending it to `bench gate` widens *where* the claim is made, not
*what* it claims.

**Honest residuals:**
- A nondeterministically-passing test's lucky green lives for the freshness
  window, and the 60 min window lengthens that exposure. Reds are never
  reused — a red always re-judges.
- `shellcheck` stays outside the declared closure, so an in-place upgrade of
  it within the window hides behind a reused green. That is the accepted
  residual ADR 0002 posture 5 records: declaring it would disable reuse on
  every host that legitimately lacks it. Clock- and network-dependence stay
  outside any closure too; they are why expiry does not go to infinity.
- The `go` and `node` binaries *are* declared, so a toolchain upgrade changes
  the subject. The collector hashes each binary and its shebang chain on every
  subject build; the post-build measurement shows that cost is not material on
  warm caches. `npm` is not declared — see the toolchain-closure ruling above —
  and needs none of its own, since it upgrades with the node it ships beside.

**Reviewer rulings:** gate-level reuse granted, and the freshness window set
to 60 min once the tools closure landed. That covers the real
review-then-commit gap the old 10 min window defeated, while keeping a bound
on ambient drift.

**What this does not buy:** the first gate after any edit. The FT91 ≤60 s
stop rule is about a changed tree — lever 1 took that to 128 s and lever 4 is
what moves it further; lever 2 makes the *second* judgment of the same tree
free, at 0.57 s.

## Lever 3 — `-count=1` on the contract phase (recommend: reject)

Removing `-count=1` from `go test ./internal/contract/...`
(`internal/gate/phases.go:100`) would let Go's per-package test cache skip
unchanged packages — but Go's cache key sees the tested package's own inputs,
not subprocess effects. These suites exec `dist/bench`, `git`, `npm`; a
change in, say, `internal/adopt` changes `dist/bench` behavior without
invalidating `surface`'s cached pass. That is a *stale green with an
incomplete key* — precisely the unsoundness ADR 0003 rejected for the canary
skip-key. Lever 2 delivers the same class of saving with a **complete** key
(the whole tree plus the oracle closure), so per-package caching buys nothing
sound on top of it. The `test` phase's existing use of the Go cache stays as
is: its packages are unit-level and the posture is already accepted.

## Lever 4 — the artifact suite itself (the live next slice; ordinary spec work)

With lever 1 landed, the wall is the contract phase ≈ the artifact package
(~109 s solo) inside a 128 s changed-tree gate — which is why this is the
next slice, and why FT91 ruled it out of its own scope as a separate
capability (test-package architecture, not oracle semantics). The retired
FT91 inventory established that the package's remaining cost is legitimate:
the collapsible generations already share one build; what remains is posture
subjects (4 hermetic double-builds, 3 GOPROXY-off, 2 stale-record),
Distributable's host+non-host build, and npm/node/git-heavy offline-smoke
work — none cuttable without changing what green proves.

The non-semantic lever left is **scheduling within the phase**: the package
runs serially (no `t.Parallel`; `go test` parallelizes across packages, not
within one). Splitting it by subject into 3–4 sub-packages (posture / offline
journey / promotion+staging / distributable) lets the existing `go test`
scheduler overlap them. Watch-outs from the inventory: `TestMain`'s
shared-cache posture is package-wide and must travel with each split;
posture tests must keep stripping it; the shared prepared-artifact singleton
is per-package and would need one owner package or per-split sets (the
sharers all live on the inspection/promotion side, so keeping them together
preserves the hoist). Alternative — introducing `t.Parallel` — changes the
hazard analysis from ordering to races (inventory's explicit warning) and is
the worse trade. Expected: artifact wall ~109 s → ~40–60 s, gate wall →
~60–75 s (then jointly bound with the `test` phase at ~62 s).

The dormant outer conformance/contract width cap
(`decisions/gate-concurrency.md`) rides the same 128 s figure: the canary no
longer overlaps the contract phase, and the wall did not drop below the
contract phase's own length, so there is no oversubscription symptom to act
on yet. Re-check it on the measurement that follows this split.

## Re-design verdict

A re-write of gate sequencing is not the fast path. The current design is
already: one thin entry (`.bench/gate.sh` → `gate-phases`), a DAG scheduler
that runs every phase concurrently with only real edges (build → readers of
`dist/bench`), per-phase attribution, durable subject-bound verdicts, and a
self-defending canary. FT91 #4 already falsified "pipeline structure is the
wall" once — the measured wall was (1) one serial test package and (2)
re-judging unchanged subjects, and no re-ordering of phase calls touches
either. Lever 2 closed the second; the first is all that is left at 128 s.

The one architectural idea a re-design would surface — memoizing verdicts at
finer grain (per phase or per package) — collapses on inspection: the only
*closed* input key the system has is the whole-subject key (tree + oracle
closure). Any finer key is a file→test map, which is diff-scoped gating,
ruled unsound. So sound memoization is exactly lever 2, at whole-gate
granularity, and the machinery for it already exists.

## Recommended sequence

Done: per-test canary bites (lever 1); tools-closure completion, gate-level
reuse, and the 60 min freshness window (lever 2). What is left:

1. **Artifact package split** (lever 4) — ordinary spec; the one slice that
   takes the changed-tree gate from 128 s to ~60–75 s.
2. **Re-measure** against the FT91 ≤60 s stop rule and update this record;
   re-check the dormant width cap on the same run.

Reject: `-count=1` removal (lever 3), canary input-key skip, fixture
batching, diff-scoped gating — all unsound or already ruled.
