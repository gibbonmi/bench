# Parallel-session landings

Status: ready

## Destination

Decide how a reviewed Bench build keeps one explicit Git source and destination
while independently green ticket commits, phase-owned handoff commits, and
other work land concurrently. The result must preserve exact ownership review,
semantic review of the complete build, one authoritative whole-project gate on
the prospective landing tree, and fail-closed publication without recreating
the retired spec-build lifecycle.

## #1: Why can the current branch review base not represent an interleaved build?

Blocked by: none
Type: Research

### Question

Trace the producer and consumers of `branch.<name>.benchBase`, reproduce the
FT198 preflight failure over every relevant candidate base, and identify
whether a correct base value exists without hiding feature work or admitting
phase-owned metadata into its ownership fence.

### Answer

`bench shift` records one mutable branch review base; `bench diff` and
`bench preflight` consume the contiguous changed set from that base to the
current checkout. FT198 landed three authorized ticket commits and then a
roadmap drain committed the phase-owned handoff. Every base that retains a
non-empty FT198 review fails ownership on `capture/session-handoff.md`; the only
base after that commit makes the review subject empty. No correct scalar base
value exists for the interleaved history. Evidence and citations:
`specs/parallel-session-landings/decisions/assets/parallel-session-landings-research.md`.

## #2: What model replaces the inferred contiguous branch range?

Blocked by: #1
Type: Grill

### Question

Choose whether reviewed landing is modeled as an explicit Git source composed
onto an explicit destination, or continues to infer authorship from one mutable
branch base. Define the conflict and publication posture.

### Answer

Model reviewed landing as an explicit Git source and destination. Compose the
source onto the current destination like an ordinary controlled merge: textual
or structural conflicts refuse, a clean composition receives the authoritative
whole-project gate, and publication proceeds only while the destination still
equals the expected tip. The ownership proof grades the source change, not
every commit that happened to land after a mutable branch base.

## #3: Which Git identity represents the source?

Blocked by: #2
Type: Grill

### Question

Choose whether the source is a Git-native integration branch plus its frozen
starting commit, an ordered ledger of landed commit IDs, or another immutable
identity. It must preserve the complete review patch when some ticket commits
are already ancestors of the destination, reject source mutation during held
review, and support the currently interleaved FT198 history without rewriting
`main`.

### Answer

The source is one build-owned Git integration branch plus a frozen base commit.
Every independently green ticket lands through `bench commit` on that branch,
advancing its tip without moving the destination branch. Semantic review pins
the exact `(base, source tip)` pair; any later source-tip movement invalidates
that held review. The destination tip is a separate identity used only for
composition and publication. Do not add an ordered commit ledger or a parallel
run-state machine.

For the already interleaved FT198 build, `0924e02e` is the frozen base and
`c46b135a` is the initial source tip. The remaining doctrine and AXI tickets
advance a build-owned branch from that tip while `main` and all landed commits
remain unchanged.

## #4: Which destination changes invalidate composition or semantic review?

Blocked by: #3
Type: Grill

### Question

Define the response when the destination advances after source review: when to
recompose, when an ordinary merge conflict refuses, and whether review may
survive a destination-only change whose effect on the source is proven absent.

### Answer

Semantic review binds only to the frozen `(base, source tip)` pair. Destination
movement does not invalidate that review while the source identity remains
unchanged. Every destination movement does invalidate the prior prospective
tree and its gate verdict: Bench composes the unchanged source onto the new
destination and runs the authoritative whole-project gate again.

A merge conflict refuses before publication. Resolving that conflict changes
the effective source result and requires fresh semantic review at the new
source tip. A clean composition requires no second semantic review merely
because the destination advanced. For FT198, the later roadmap and handoff
commit therefore require recomposition and a fresh gate but do not erase review
of the unchanged FT198 source.

## #5: Which gate evidence may cross prospective landing trees?

Blocked by: #2
Type: Grill

### Question

Choose whether composed trees may share retained gate evidence, distinguishing
exact-tree reuse, gate-derived component inheritance, and the project-green
transition.

### Answer

A whole-project green verdict is reusable only for the identical Git tree and
identical derived oracle identity. Distinct prospective trees never inherit an
opaque whole-tree green result. The gate may assemble component evidence only
when its authoritative closure proves every component input unchanged and it
executes every missing check. Only the final exact prospective tree, after its
authoritative verdict and successful destination publication, advances
project-green.

## #6: Does this scope recreate multi-coordinator spec-build state?

Blocked by: #2
Type: Grill

### Question

Choose whether the fix is the Git-native landing subject needed by ordinary and
reviewed serial ticket commits, or a new run-state, receipt, recomposition, and
status lifecycle.

### Answer

Keep the scope to a Git-native controlled merge with an explicit source and
destination. Do not recreate the removed spec-build lifecycle, run revisions,
prepared-operation journals, receipt helpers, recomposition command, or AXI run
status surface. Conflict recovery uses Git's merge model plus Bench's ownership,
review, exact-gate, and destination-tip checks.

## #7: How does the already interleaved FT198 build resume?

Blocked by: #3, #4
Type: Grill

### Question

Choose the recovery construction for the three FT198 ticket commits already on
`main`, the uncommitted doctrine ticket in its retained assignment worktree,
and the later phase-handoff commit. Require a complete review patch without
rewriting or dropping any landed commit.

### Answer

The retained `land-index-doctrine` assignment branch is the FT198 integration
source. It remains at `c46b135a` with the nine-path doctrine diff preserved in
its owned worktree; do not copy, reconstruct, or discard that work. Keep
`0924e02e` as the frozen review base. Finish and commit the doctrine ticket on
that source, then start the AXI ticket from the resulting source tip.

Semantic review grades the complete `0924e02e..<source-tip>` change. Compose
that reviewed source onto the current `main`; the first three ticket commits
are already ancestors of the destination, so the composition applies only the
remaining source commits. A clean result receives the authoritative full gate
and final `--spec` landing. A conflict is resolved on the source and requires a
fresh review. Do not rewrite `main` or rerun the three landed tickets.

## Not yet specified

## Spec-writer discretion

- Reversible internal names and storage layout for the chosen Git identities.

## Out of scope

- Recreating the removed `bench spec-build` command family or its run-state
  lifecycle.
- Weakening whole-project gate authority or allowing a partial verdict to
  publish a landing.
- Treating phase-owned metadata as implementation-owned merely to make a fence
  green.
- Rewriting or deleting the landed FT198 or roadmap-maintenance commits.
- Making mixed-authorship commits the normal concurrency mechanism.

## Sources

- Path: `specs/parallel-session-landings/decisions/assets/parallel-session-landings-research.md`
  Supports: #1's current-tree producer trace and exact FT198 base matrix.
  Drift: current-tree research from 2026-08-13; rerun the public preflight and re-check the cited owners if diff, preflight, shift, or handoff changes.
