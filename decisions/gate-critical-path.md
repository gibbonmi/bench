# Gate assessment — where the wall lives, and the levers left

Assessed 2026-07-29 at working tree of `187bc36`. Question: how fast can the
dev gate get without sacrificing quality, including the oracle-semantics
levers FT91 parked as reviewer decisions — verdict reuse for the exact tree,
canary narrowing, and freshness relaxation. This record incorporates the
retired FT91 timeline's historical measurements, one fresh instrumented run,
and the resulting decision analysis.

## Result

The gate is **one package wide**: `internal/contract/surface/artifact`
(~109 s solo, ~133 s under gate load) sets both remaining walls — the
contract phase runs it once, and the canary sweep runs it five more times
(once per `surface/artifact` fixture, at inner width 2). Everything else in
the gate finishes by ~t+62.

Three moves, in order, take the gate from ~135–172 s to roughly:

| after | changed-tree gate | unchanged-tree gate |
|---|---|---|
| today | 135–172 s | 135–172 s |
| 1. per-test canary bites | ~110–135 s | same |
| 2. gate-level verdict reuse + freshness ruling | same | **~2–5 s** |
| 3. artifact-suite restructure (parallelism) | **~60–75 s** | ~2–5 s |

Move 1 follows semantics the reviewer already ruled (stage 2's own wording);
move 2 extends an authority ADR 0002 already grants `bench commit`; move 3 is
ordinary spec work with no semantic content. A full gate re-write is **not**
warranted — the scheduler is already maximally concurrent and the bound is
elsewhere (see "Re-design verdict").

## Fresh measurement (2026-07-29, this box, warm caches, session load present)

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
fixture at `GOMAXPROCS=2`, ~110–133 s. The canary and the contract phase are
floored by the same suite.

## What already shipped (do not re-litigate)

- Concurrency budget (ADR 0009): 10–15 min → ~5.5 min.
- Dev/ship tier split, host-only dev generator, shared dev caches
  (`decisions/cost-follows-project-size.md` #3, #8–#10): removed the
  release-evidence and hermeticity cost from the per-commit path.
- Canary stage 1 (family/package scoping) + stage 2 (compiled bites, no
  nested gates): canary 250 s → ~152 s.
- Prepared-artifact hoist (FT91 #2/#3): artifact suite ~141 s → ~109 s.
- Commit-time verdict reuse for a fresh green on the identical closed subject
  (ADR 0002 posture 5; `internal/commit/commit.go:118`).

Closed rulings that stay closed: diff-scoped gating is unsound; no check is
weakened or dropped for wall-clock; canary skip-on-"inputs unchanged" key
rejected as fail-dangerous (ADR 0003); fixture batching rejected
(cross-contamination, ADR 0003).

## Lever 1 — per-test canary bites (recommended; not oracle-weakening)

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

**What it buys:** the five artifact fixtures stop paying ~133 s each; every
behavior-owned fixture drops to binary startup plus one test (seconds; the
worst single tests do one real build, ~5–15 s). Sweep worker-seconds fall
from ~1,480 to low hundreds. Canary wall lands at ~30–60 s (bounded by the
conformance-scoped remainder, the baselines, and the slowest single test),
and the gate becomes contract-phase-bound at ~110–135 s.

**Fail postures (all loud, matching stage 1/2):**
- Binding: an explicit per-fixture `TEST` file naming the owning test,
  mirroring the existing `CHECK` convention — a name the package does not
  declare is a red at sweep start (same shape as `CHECK`'s unknown-check
  error), never a silent fallback to the full package.
- A renamed or deleted owning test → did-not-bite red.
- Vacuity baseline runs the same `-test.run` shape over the empty tree
  (scoped-vs-scoped, #6 precedent). Spec-time detail: a single scoped test
  can be near-silent on the empty tree, and `scopeBaselines` rejects an
  empty baseline — the scoped baseline likely needs `-test.v` so its output
  is never empty.

**Reviewer decision needed:** sign-off that per-test scope rides the existing
scoping precedent. Recommended: yes — it aligns execution with semantics the
map already declared.

## Lever 2 — gate-level verdict reuse + freshness (recommended, with closure fix)

**What exists:** every gate run already builds a *closed oracle subject* —
working-tree hash (untracked files included: `internal/git/tree.go:16` stages
`add -A` into a throwaway index), resolved gate, execution-policy version,
`PATH`, the gate script's full launcher closure, and every declared
environment/path/tool input from `.bench/gate-inputs.json`
(`internal/gate/subject.go:53`). The verdict is recorded durably under lock,
pending-before-run, subject-rechecked after. `Inspect` computes
`ReusableGreen` = ready + green + same tree + same oracle + < 10 min old
(`internal/gate/verdict.go:243`). Today only `bench commit` short-circuits on
it; `bench gate` (`gate.Execute`) always re-runs.

**The change:** let gate execution itself return the reused green for a
`ReusableGreen` subject, printing an explicit line (the `commit.go:118`
precedent: `gate: green (fresh verdict reused for this tree)`), with an
opt-out flag for a deliberately fresh run. Because the shift loop, the stop
hook, and `/bench-final-check` all route through this package, the reuse
covers every repeat-subject re-judging in one place: final-check → commit →
push chains, post-review re-gates on an untouched tree, and an operator
re-typing `bench gate`.

**What "green" comes to mean:** "this exact closed subject was judged green
within the freshness window," rather than "the oracle just ran." ADR 0002
posture 5 already accepted exactly this claim for the commit path, with the
same reopen trigger (a reused verdict authorizing what a run would have
refused). Extending it to `bench gate` widens *where* the claim is made, not
*what* it claims.

**Honest residuals:**
- A nondeterministically-passing test's lucky green lives for the freshness
  window. Already true on the commit path; a longer window lengthens it.
  Reds are never reused — a red always re-judges.
- The subject closure does not cover undeclared ambient inputs: the Go
  toolchain, `node`, `npm`, and `shellcheck` are not in
  `.bench/gate-inputs.json` `tools`, so a toolchain upgrade does not change
  the subject. At 10 minutes this is negligible; before relaxing freshness,
  **declare them** (one-line manifest edit; the collector hashes each binary
  and its shebang chain — tens of ms). Clock- and network-dependence stay
  outside any closure; they are why expiry should not go to infinity.

**Reviewer decisions needed:** (a) grant gate-level reuse; (b) pick the
window. Recommended: grant it; extend freshness to ~60 min *after* the tools
closure is completed. That covers the real review-then-commit gap the 10 min
window currently defeats, while keeping a daily bound on ambient drift.

**What this does not buy:** the first gate after any edit. The FT91 ≤60 s
stop rule is about a changed tree — levers 1 and 3 are what move that;
lever 2 makes the *second* judgment of the same tree free.

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

## Lever 4 — after the levers: the artifact suite itself (ordinary spec work)

Once lever 1 lands, the wall is the contract phase ≈ the artifact package
(~109 s solo). The retired FT91 inventory established that its remaining cost
is legitimate: the collapsible generations already share one build; what
remains is posture subjects (4 hermetic double-builds, 3 GOPROXY-off, 2
stale-record), Distributable's host+non-host build, and npm/node/git-heavy
offline-smoke work — none cuttable without changing what green proves.

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

Re-check afterwards, per the map's dormant item: the outer
conformance/contract width cap, if oversubscription symptoms persist once the
canary no longer overlaps the contract phase.

## Re-design verdict

A re-write of gate sequencing is not the fast path. The current design is
already: one thin entry (`.bench/gate.sh` → `gate-phases`), a DAG scheduler
that runs every phase concurrently with only real edges (build → readers of
`dist/bench`), per-phase attribution, durable subject-bound verdicts, and a
self-defending canary. FT91 #4 already falsified "pipeline structure is the
wall" once — the measured wall is (1) one serial test package and (2)
re-judging unchanged subjects, and no re-ordering of phase calls touches
either.

The one architectural idea a re-design would surface — memoizing verdicts at
finer grain (per phase or per package) — collapses on inspection: the only
*closed* input key the system has is the whole-subject key (tree + oracle
closure). Any finer key is a file→test map, which is diff-scoped gating,
ruled unsound. So sound memoization is exactly lever 2, at whole-gate
granularity, and the machinery for it already exists.

## Recommended sequence

1. **Per-test canary bites** (lever 1) — spec against stage 2's semantics;
   canary tail collapses; gate becomes contract-bound.
2. **Tools-closure completion + gate-level reuse + freshness ruling**
   (lever 2) — small diff in `internal/gate` + one manifest line; two
   reviewer sign-offs recorded (reuse grant, window).
3. **Artifact package split** (lever 4) — ordinary spec; takes the
   changed-tree gate to ~60–75 s.
4. **Re-measure** against the FT91 ≤60 s stop rule and update this record;
   re-check the dormant width cap.

Reject: `-count=1` removal (lever 3), canary input-key skip, fixture
batching, diff-scoped gating — all unsound or already ruled.
