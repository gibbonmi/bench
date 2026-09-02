# Gate critical path — where the wall lives, and the levers left

Status: shaping

Assessed 2026-07-29 at the working tree of `187bc36`. Re-measured on the
landed tree of `d91b709` once moves 1 and 2 shipped, and again on the
artifact-split landed tree of `6016be6`. Question: how fast can the dev gate
get without sacrificing quality? This includes the oracle-semantics levers
FT91 parked as reviewer decisions — verdict reuse for the exact tree, canary
narrowing, and freshness relaxation. This durable record incorporates the
retired FT91 timeline's historical measurements, the instrumented run that located
the wall, the post-build re-measurement, and the resulting decision analysis.

## Result

The artifact split is landed. Its four package processes overlap: a focused
recursive run spans about 51 s while their package spans total about 127 s.
The changed-tree gate is now 89.91 s, down from 128 s. Its critical path is
still the contract phase: `posture` takes 85.415 s under the full gate. The
focused posture span is 50.917 s, so concurrent outer phases materially extend
the process that otherwise bounds the split suite.

Three moves take the gate from its ~135–172 s baseline; all are measured:

| after | changed-tree gate | unchanged-tree gate |
|---|---|---|
| baseline | 135–172 s | 135–172 s |
| 1. per-test canary bites *(landed)* | **128 s** | same |
| 2. gate-level verdict reuse + freshness ruling *(landed)* | same | **0.6 s** |
| 3. artifact-suite restructure (parallelism) *(landed)* | **89.91 s** | 0.6 s |

Moves 1 and 2 landed in FT91's gate-fastpath build. Move 1 rests on semantics
the reviewer had already ruled (stage 2's own wording). Move 2 rests on an
authority ADR 0002 already granted `bench commit`, now extended to gate
execution.

Move 3 exposed package-level parallelism without changing oracle semantics
and removed 38.09 s from the measured changed-tree wall. The split's focused
trace shows that repeated package processes are not the dominant residual
cost. A reviewer decision is needed before reviving the dormant outer-width
cap. The full gate's 85.415 s posture span is 34.498 s longer than the
50.917 s focused span. The 69.506 s test phase and other outer phases
overlap it. A full gate re-write is **not** warranted (see "Re-design
verdict").

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

Canary accounting: 174 fixtures, 33 of them behavior-owned (8 axi, 6 runtime,
10 surface, 5 surface/artifact, 4 surface/publication). Each behavior-owned
fixture invokes its owning package's compiled test binary over the fixture
tree (`internal/canary/canary.go`, `RunBite`). The sweep's worker-seconds are
therefore ≈ 5×133 + 6×52 + 10×39 + 4×20 + 8×4.5 ≈ **1,480**, plus the ~140
conformance-scoped fixtures and baselines. At 8 workers that is ~185 s, which
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
changed-tree full gate**. That gate stands at 89.91 s, floored by the
contract phase's `posture` package (85.415 s under the gate). The
unchanged-tree path beats its expected ~2–5 s by an order of magnitude: it
prints `gate: green (fresh verdict reused for this tree)` and does no phase
work.

## Artifact-split measurement (2026-07-29, `6016be6`, WSL2 host, warm caches)

The worktree binary was rebuilt with `bash scripts/go-build.sh "$PWD"
dist/bench`. The focused command was `/usr/bin/time -p go test -json -count=1
./internal/contract/surface/artifact/...`; its JSON timestamps provide the
package intervals below. The public worktree-local gate command was
`/usr/bin/time -p bash bin/bench.sh gate --fresh`, which routes through that
rebuilt binary. Wall clock is evidence, not a timeout assertion.

| package | JSON start → pass | span |
|---|---|---:|
| `distributable` | 20:09:16.526 → 20:09:30.830 | 14.304 s |
| `offline` | 20:09:16.577 → 20:09:49.471 | 32.893 s |
| `posture` | 20:09:16.577 → 20:10:07.494 | 50.917 s |
| `prepared` | 20:09:16.578 → 20:09:44.826 | 28.249 s |

All four starts occur within 52 ms and every interval overlaps the others; a
serialized trace would not satisfy this record. The focused envelope is 50.97
s. Its 126.363 s of package spans is 17.363 s (16%, 15.93% rounded) above the
former 109 s single-process package. 75.393 s of that work overlaps. The
repeated process overhead is therefore visible but not the remaining wall.

The fresh gate is green in **89.91 s**, a 38.09 s (30%) reduction from the
128 s baseline. Its longest phase is contract: posture 85.415 s, offline
70.327 s, prepared 68.353 s, and distributable 45.497 s. The concurrently
running test phase is 69.506 s, and conformance-suite is 36.696 s. The
remaining critical path is contract/posture plus gate setup.

The focused 50.917 s posture span becomes 85.415 s under those outer spans,
an oversubscription signal. The dormant outer-width cap should therefore
return as a **reviewer decision**. It is not implemented here: this single
host-sensitive observation establishes the contention to price, not a
scheduler policy or width.

**The toolchain closure declares `go` and `node`, not `npm`.** Every subject
hash walks each declared binary and its shebang chain. The 0.57 s reuse
timing covers that whole hash-and-look-up path. On warm caches the hashing is
an observed non-issue at this measurement. A cold cache or a slower
filesystem is untested.

`npm` is deliberately outside the closure, on the reviewer's ruling, for two
reasons. First, a host with node but no npm (the Debian split packaging)
would hold a permanently open subject. This disables reuse and locks
`bench prep-release`'s dev-green entry check with no operator escape.
Second, npm's
shebang chain re-hashes `node`, doubling the closure's cost for a signal
already carried. Nothing is lost, because an npm upgrade ships with its
node — a toolchain upgrade still changes the subject through `node`.

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
  closure declaring `go` and `node` (lever 2): a re-judged unchanged tree →
  0.6 s.

Closed rulings that stay closed: diff-scoped gating is unsound. No check is
weakened or dropped for wall-clock. The canary skip-on-"inputs unchanged" key
is rejected as fail-dangerous (ADR 0003). Fixture batching is rejected too
(cross-contamination, ADR 0003).

## Lever 1 — per-test canary bites (landed; not oracle-weakening)

**What:** a behavior-owned fixture's bite run invokes its compiled package
binary with `-test.run` pinned to the one contract test whose failure message
its EXPECT is. This replaces running every `Test*` in the package.

**This is not a semantics change.** Stage 2's ruling (FT91 map #7) narrowed
the claim — *"a behavior-owned fixture's bite is proven at its owning
contract test, not by a nested gate."* The EXPECT is, by the same ruling,
"the owning contract test's own failure message." Running the other ~25
tests in the package per fixture is incidental breadth left over from the
compile-once implementation, not part of the declared claim. Scoping
execution to the declared owner is the fifth-arm/#6/#7 precedent ("scoping a
fixture to the check its EXPECT names is not weakening"). Here it applies
one level deeper.

**What it bought:** the five artifact fixtures no longer pay ~133 s each.
Every behavior-owned fixture is binary startup plus one test (seconds; the
worst single tests do one real build, ~5–15 s). Measured at `d91b709`: the
solo sweep is 25.2 s wall against ~171 s of CPU. It is bounded by the
conformance-scoped remainder, the baselines, and the slowest single test. The
gate is contract-phase-bound at 128 s.

**Fail postures (all loud, matching stage 1/2):**
- Binding: an explicit per-fixture `TEST` file names the owning test,
  mirroring the existing `CHECK` convention. A name the package does not
  declare is a red at sweep start (same shape as `CHECK`'s unknown-check
  error). This never falls back silently to the full package.
- A renamed or deleted owning test → did-not-bite red.
- A contract group's vacuity baseline runs its package binary whole over the
  empty tree — no `-test.run`, no `-test.v`. It stays deliberately wider than
  the scoped runs it grades, and no baseline carries a per-test twin. The two
  directions fail differently. A baseline narrower than its runs lets a
  vacuous EXPECT clear the screen in silence. A wider one at worst calls a
  sound EXPECT vacuous, which is a red someone reads. This is the
  wide-baseline principle `scopeBaselines` already holds for phase pins, and
  it is the one place scoping stops.

**Reviewer ruling:** granted — per-test scope rides the existing scoping
precedent, aligning execution with semantics the map had already declared.

## Lever 2 — gate-level verdict reuse + freshness (landed)

**What exists:** every gate run already builds a *closed oracle subject*.
This subject includes the working-tree hash (untracked files included:
`internal/git/tree.go:16` stages `add -A` into a throwaway index), the
resolved gate, the execution-policy version, and `PATH`. It also includes
the gate script's full launcher closure and every declared
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
same reopen trigger. A reused verdict authorizes what a run would have
refused. Extending it to `bench gate` widens *where* the claim is made, not
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
stop rule remains about a changed tree. The split took the observed first
run to 89.91 s. Lever 2 makes the *second* judgment of the same tree free,
at 0.57 s.

## Lever 3 — `-count=1` on the contract phase (recommend: reject)

Removing `-count=1` from `go test ./internal/contract/...`
(`internal/gate/phases.go:100`) would let Go's per-package test cache skip
unchanged packages. But Go's cache key sees the tested package's own inputs,
not subprocess effects. These suites exec `dist/bench`, `git`, `npm`; a
change in, say, `internal/adopt` changes `dist/bench` behavior without
invalidating `surface`'s cached pass. That is a *stale green with an
incomplete key* — precisely the unsoundness ADR 0003 rejected for the canary
skip-key.

Lever 2 delivers the same class of saving with a **complete** key
(the whole tree plus the oracle closure). So per-package caching buys
nothing sound on top of it. The `test` phase's existing use of the Go cache
stays as is: its packages are unit-level and the posture is already
accepted.


## Destination

Record the measured FT91 critical path and leave the only remaining policy choice — whether to bound outer-phase concurrency — explicitly with the reviewer.

## #1: How should the artifact suite expose its remaining independent work?

Blocked by: none
Type: Task

### Question

How should the artifact contract suite remove its one-package serial wall
without weakening a check or changing oracle semantics?

### Answer

With lever 1 landed, the wall was the contract phase ≈ the artifact package
(~109 s solo) inside a 128 s changed-tree gate. This is why the split became
a separate capability (test-package architecture, not oracle semantics). The
retired FT91 inventory established that the package's remaining cost is
legitimate. The collapsible generations already share one build. What
remains is posture subjects (4 hermetic double-builds, 3 GOPROXY-off, 2
stale-record), Distributable's host+non-host build, and npm/node/git-heavy
offline-smoke work. None of it is cuttable without changing what green
proves.

The non-semantic lever left is **scheduling within the phase**: the package
runs serially (no `t.Parallel`; `go test` parallelizes across packages, not
within one). Splitting it by subject into 3–4 sub-packages (posture / offline
journey / promotion+staging / distributable) lets the existing `go test`
scheduler overlap them. Watch-outs from the inventory: `TestMain`'s
shared-cache posture is package-wide and must travel with each split.
Posture tests must keep stripping it. The shared prepared-artifact singleton
is per-package and would need one owner package or per-split sets. The
sharers all live on the inspection/promotion side, so keeping them together
preserves the hoist.

Alternative — introducing `t.Parallel` — changes the hazard analysis from
ordering to races (inventory's explicit warning), and is the worse trade.
Expected: artifact wall ~109 s → ~40–60 s, gate wall → ~60–75 s (then
jointly bound with the `test` phase at ~62 s).

The split's measurement changes the dormant outer conformance/contract
width-cap posture. The focused package trace completes posture in 50.917 s.
The full gate leaves posture at 85.415 s while the 69.506 s test phase
overlaps it. That is enough contention evidence to ask the reviewer whether
a bounded outer scheduler is worth its policy and implementation cost. It is
not enough authority to choose or install a cap.

## Re-design verdict

A re-write of gate sequencing is not the fast path. The current design is
already: one thin entry (`.bench/gate.sh` → `gate-phases`), a DAG scheduler
that runs every phase concurrently with only real edges (build → readers of
`dist/bench`), per-phase attribution, durable subject-bound verdicts, and a
self-defending canary. FT91 #4 already falsified "pipeline structure is the
wall" once. The measured wall was (1) one serial test package and (2)
re-judging unchanged subjects, and no re-ordering of phase calls touches
either. Lever 2 closed the second; the first is all that is left at 128 s.

The one architectural idea a re-design would surface — memoizing verdicts at
finer grain (per phase or per package) — collapses on inspection. The only
*closed* input key the system has is the whole-subject key (tree + oracle
closure). Any finer key is a file→test map, which is diff-scoped gating,
ruled unsound. So sound memoization is exactly lever 2, at whole-gate
granularity, and the machinery for it already exists.

## Not yet specified

- Whether the reviewer wants a bounded outer concurrency cap after pricing its policy and implementation cost against the observed focused versus in-gate posture spans.

## Spec-writer discretion

## Out of scope

- Rewriting gate sequencing, per-package verdict reuse, `-count=1` removal, canary input-key skipping, fixture batching, and diff-scoped gating. The record rejects or closes each as unsound or already decided.

## Sources
