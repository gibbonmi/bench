# Cost follows project size

## Destination

One complaint — gate, context, and delegation cost grow with the tree while the
value of a small change does not — decided across its three angles: FT91 (gate
wall-clock), FT101 (ambient-surface scope), FT136 (delegate slicing). Each angle
leaves this map spec-ready or explicitly deferred with a revive trigger; the
roadmap rows stay separate because the owners differ.

## #1: Which FT91 arms are in scope for the next build cycle?

Type: Grill

### Question
Four arms remained: parallelizing the fifteen conformance checks, removing the
hardcoded `-count=1`, cache infrastructure (hermetic build cache / pinned-subject
verdicts), and reviving the outer phase-concurrency cap.

### Answer
Conformance arm only (reviewer, 2026-07-26). Conformance is ~94% of measured
wall clock and its checks run strictly serially, so the pure-scheduling win is
taken before any oracle-semantics or cache-infrastructure decision is spent.
`-count=1` and caching are deferred until re-measurement after this arm lands
shows conformance is no longer the long pole; the capping arm stays dormant per
the roadmap. Superseded by #3 (2026-07-26): #2's timing showed one composite
check owns 99.8% of the phase, so the conformance arm became the two-tier gate
split instead of check fan-out.

## #2: How does conformance wall clock distribute across the fifteen checks?

Type: Research

### Question
Instrument the conformance driver (`internal/conformance`) with per-check
timing and run the gate once. If one check dominates the ~400–520 s phase,
parallelization buys little and #1's chosen arm changes shape; if the time
spreads, an indexed-slice fan-out recovers most of it. Each check is a pure
function over a read-only tree, so ordering is recoverable by collecting into
an indexed slice.

### Answer
Measured 2026-07-26 (throwaway probe on this tree). The phase is one composite
check: `checkPackageCoreAndGuards` is ~99.8% of ~826 s; the other fourteen
checks total ~1.3 s, so fifteen-check fan-out buys nothing. Inside the
composite: `checkReleasePreflight` ≈372 s (release artifact matrix build plus
preflight verify) and `checkGoCore` owns the rest (inner `go test` over all
non-contract packages, worktree race test, cross-compile matrix, build+vet).
The inner `go test` runs without `-count=1` and leans on Go's test cache; on a
cache miss `internal/preflight` alone exceeds the 600 s go-test default
package timeout (its subtests rebuild the preflight binary and exercise real
archives — slow, not hung; >900 s observed uncached), pushing the phase past
1000 s. The levers are staging (which checks run when) and cache behavior, not
scheduling.

## #3: Given the timings, does the parallelization spec proceed, and what ships with it?

Blocked by: #2
Type: Grill

### Question
Go/no-go on the fan-out given the measured distribution; worker-width policy and
its interaction with the canary concurrency budget (`bounds.CanaryInnerWidth`);
whether per-check timing stays as permanent gate output or was a throwaway probe.

### Answer
No-go on the fan-out; the gate splits into two tiers instead (reviewer,
2026-07-26). **Dev tier** (`bench gate`, the shift loop, final-check): drops
`checkReleasePreflight` and the cross-compile matrix, and the inner `go test`
excludes release-only packages (`internal/preflight` and kin) the way it
already excludes `contract` — dev green means the kit works from the tree, and
is immune to the test-cache-miss blowup. **Ship tier** (`bench prep-release`,
a new command): artifact matrix build, cross-compile matrix, release preflight
verify, and the release-only package tests; the release path refuses without
its evidence. Final-check on green prints a one-line ship-tier reminder, never
a prompt; the pre-push hook stays fast-tier. Restaging is not check-weakening:
every check keeps full authority at a boundary no release can bypass, and dev
green is explicitly the narrower claim. Per-check timing becomes permanent
gate observability, owned by the `RunConformance` driver in a stable format
(the probe file is deleted). Worker-width policy and `CanaryInnerWidth`
interaction: n/a — no fan-out.

## #4: Does the FT136 slicing rule wait for the cheap-tier retest?

Type: Grill

### Question
The FT86 evidence confounds tier, charge weight, and slicing, and the row named
a cheap-tier retest as its acceptance trigger.

### Answer
Rule lands now; retest runs separately (reviewer, 2026-07-26). The fence rule
and shared-primitives-first are tier-independent — the FT86 evidence (zero
conflicts on fence-aligned slices vs ~25 min / 184k tokens on the theme-cut
slice) supports them regardless of tier. The retest gates only whether
mid-tier-by-default for build delegates is settled; it does not block the kit
edit.

## #5: Where does the FT136 rule land in the kit?

Type: Grill

### Question
Slicing is authored at spec time, fence alignment is checked at charge time,
and fence-boundary duplication is hunted at review time — one skill or three?

### Answer
Three surfaces, one source (reviewer, 2026-07-26). `craft-spec` owns the rule:
the slice boundary and the ownership fence must be the same line, and shared
primitives are named up front and land as a deep-unit slice before the
consuming seams. `craft-review`'s Standards axis adds fence-boundary
duplication as an explicit hunt. `craft-delegate` gets a one-line charge-time
cross-reference pointing at `craft-spec`. Kit edit under `craft-synthesis`.

## #6: Does the cheap tier hold on a genuinely seam-shaped slice?

Type: Task

### Question
Run one build delegate on the cheap tier against a slice whose boundary is a
true ownership fence, charged per current `craft-delegate` (exemplar files
included), and compare outcome quality against the mid-tier norm. The reviewer
picks the slice; the worker runs and reports. Opportunistic — waits for the
next genuinely seam-shaped slice in normal work rather than manufacturing one.
Resolves only mid-tier-by-default (the tier-binding memory and `craft-line`
guidance stay as-is until then).

### Answer
— (open)

## #7: How much of FT101 does this map decide now?

Type: Grill

### Question
No linked repo is a monorepo yet — regroup-app, the first external validation
target, is a single context — but the gate-scoping half was flagged as the
contested part needing a reviewer decision.

### Answer
Guardrails now, build deferred (reviewer, 2026-07-26). Decided and closed: a
scoped gate is legitimate only on a reviewer-declared package boundary, never
derived from a diff (FT91's diff-scoped ruling stands); a change touching two
profiles takes the whole-tree gate; wall-clock is never the justification for
scope — a scoped verdict is explicit evidence, never a silent skip. The docs
half (`CONTEXT-MAP.md` layout, setup question, consumer teaching) and profile
half (path ownership, ambient-surface scoping) are deferred undesigned. Revive
trigger: a linked repo with more than one bounded context.

## Not yet specified

- FT91 cache arms (hermetic build cache; verdicts keyed on the pinned gate
  subject) — only if re-measurement after the tier split lands still hurts.
- `-count=1` freshness semantics — same trigger; oracle decision, reviewer-led.
  Now carries a measured price: uncached `internal/preflight` alone is 10+ min,
  so blanket `-count=1` on the inner suite is off the table without the split.
- FT101 docs-half and profile-half design — dim until the revive trigger fires.

## Out of scope

- Diff-scoped gating in any form — ruled unsound (no file→test map for contract
  and canary); the ruling stands.
- Reviving the outer conformance/contract concurrency cap — dormant unless
  contention flakes persist.
- Gate scope as a speed lever — FT101 guardrail; never reopened for wall-clock.
- Weakening any check to buy wall clock — green keeps meaning what it means.

## Handoff

FT136 shipped (spec retired 2026-07-26); FT101 builds nothing this cycle. The
handoff below is the FT91 tier-split spec.

1. **Module boundaries.** `internal/conformance` owns tier membership (which
   sub-checks run in the dev gate vs `prep-release`), the release-only package
   exclusion in `goCoreTestPackages`, and driver-emitted per-check timing.
   `internal/gate` phase table stays the dev tier. `bin/bench.sh` plus the Go
   core gain the `prep-release` route. `scripts/release-preflight.sh` and the
   release path own the ship-evidence refusal. Final-check's green report
   carries the one-line ship-tier reminder.
2. **Contracts.** Dev gate: same pass/fail semantics minus ship-tier checks;
   green claims "the kit works from the tree", nothing about release
   packaging. `bench prep-release`: exit 0 means ship tier green with evidence
   written (`dist/preflight/release-index.json` and artifacts); nonzero with
   diagnostics otherwise. Release path: refuses without current ship evidence.
   Timing: one stable-format line per check on every gate run; ordering
   byte-stable (indexed), values free to vary.
3. **Deep vs thin.** The conformance driver stays the deep unit — tier
   membership and timing hide behind it; checks stay pure functions.
   `prep-release` is a thin route over existing scripts plus the ship-tier
   checks; it invents no new machinery.
4. **Black-box assertables.** Dev gate green writes nothing under
   `dist/artifacts`; `prep-release` produces the artifact set and
   `release-index.json`; a seeded ship-tier failure reds `prep-release` while
   the dev gate stays green; a seeded dev-tier failure reds both; the release
   path exits nonzero without ship evidence; timing lines present, one per
   check, stable order.
5. **Gate attachment.** Dev tier attaches exactly as today (suite plus canary).
   Ship tier attaches at `prep-release` and the release path's refusal — the
   gate does not see it per-commit; that narrowing is the decided tradeoff.
6. **Hostile-input owners.** Unchanged owners, with one move: the preflight
   archive-hostility suites (hostile package archives, aggregate budgets)
   ride to the ship tier with their package — dev green explicitly does not
   cover them.
7. **Uncertainty flags.** Whether `prep-release` requires/reruns a current
   dev-green verdict or runs ship checks alone (recommend: require dev green
   via the existing `bench gate pin` machinery); which tier the canary's inner
   gates run (recommend: dev); the exact release-only package list —
   `internal/preflight` is decided, `releaseevidence`/`publication` need an
   import read. Cheap-tier viability for build delegates stays open under #6.
8. **Rejected alternatives.** Fifteen-check fan-out (killed by #2's data: one
   composite check is 99.8% of the phase); a build-approval prompt after green
   final-check (nags at commit cadence); ship tier on the pre-push hook;
   blanket `-count=1` now (measured 10+ min price); timing via a probe test
   file (driver owns it instead); deciding caching now (#1).
9. **Domain watch-outs.** The inner `go test` leans on Go's test cache — any
   change that defeats it (env perturbation, `-count=1`, cold cache) makes
   `internal/preflight` blow the 600 s default package timeout, which presents
   as a gate hang; the release-only exclusion is what removes this failure
   mode from dev runs. Timing output must keep finding order byte-stable while
   values vary.

Dependency order: n/a — single spec (the FT91 tier split).
