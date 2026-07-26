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
the roadmap.

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
— (open)

## #3: Given the timings, does the parallelization spec proceed, and what ships with it?

Blocked by: #2
Type: Grill

### Question
Go/no-go on the fan-out given the measured distribution; worker-width policy and
its interaction with the canary concurrency budget (`bounds.CanaryInnerWidth`);
whether per-check timing stays as permanent gate output or was a throwaway probe.

### Answer
— (open)

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
  subject) — only if re-measurement after the conformance arm still hurts.
- `-count=1` freshness semantics — same trigger; oracle decision, reviewer-led.
- FT101 docs-half and profile-half design — dim until the revive trigger fires.

## Out of scope

- Diff-scoped gating in any form — ruled unsound (no file→test map for contract
  and canary); the ruling stands.
- Reviving the outer conformance/contract concurrency cap — dormant unless
  contention flakes persist.
- Gate scope as a speed lever — FT101 guardrail; never reopened for wall-clock.
- Weakening any check to buy wall clock — green keeps meaning what it means.

## Handoff

1. **Module boundaries.** FT136 edit: `.agents/skills/bench-craft-spec` owns the
   slicing rule; `bench-craft-review` and `bench-craft-delegate` carry pointers
   (`.claude/skills/*` are symlinks — edit `.agents/skills/` only). FT91 edit:
   `internal/conformance` (check registry and driver) is the whole surface;
   `internal/gate` phase wiring stays untouched.
2. **Contracts.** FT91: same fifteen checks over the same read-only tree, same
   pass/fail semantics, findings collected into an indexed slice so output
   ordering is byte-stable; timing lands per #3's call. FT136: prose guidance,
   no runtime contract.
3. **Deep vs thin.** The conformance driver is the deep unit (scheduling hides
   behind it; checks stay pure functions). The two pointer skills in FT136 are
   deliberately thin — the rule lives once in `craft-spec`.
4. **Black-box assertables.** FT91: gate exit code unchanged on the existing
   fixture set; deterministic finding order across runs; a seeded failing
   fixture still reds the phase. FT136: n/a — prose; graded by review, not
   asserted.
5. **Gate attachment.** FT91 attaches directly — the conformance phase is the
   thing under change, and the existing suite plus canary observe it. FT136 the
   gate sees only structurally (skill-shape checks); semantic quality rides the
   `craft-synthesis` loops and review.
6. **Hostile-input owners.** FT91: a check that panics or deadlocks under
   fan-out must fail the phase loudly, owned by the driver; contention with
   canary's budget is bounded by the width policy #3 sets. FT136: n/a — prose.
7. **Uncertainty flags.** FT91 seams are provisional until #2's timing data and
   #3's go/no-go — the spec-writer does not start the FT91 spec before #3
   closes. Cheap-tier viability for build delegates stays open under #6.
8. **Rejected alternatives.** Deciding `-count=1` or caching now (#1); shipping
   the slicing rule only after the retest, or dropping the retest (#4); a
   single-skill home for the slicing rule (#5); shaping FT101 fully or deferring
   even its guardrails (#7); diff-scoped gating (standing FT91 ruling).
9. **Domain watch-outs.** Seam-slicing optimizes local coherence and can
   manufacture cross-fence knowledge duplication — when no delegate owns a
   shared primitive, each fence writes its own copy; the deep-unit-first rule
   exists to prevent exactly this. Conformance findings today may rely on
   incidental serial ordering; the indexed slice is what keeps green
   byte-comparable.

Dependency order: FT136 kit edit first (fully unblocked, smallest); FT91 timing
research (#2) can run in parallel, then #3, then the FT91 spec; FT101 builds
nothing this cycle.
