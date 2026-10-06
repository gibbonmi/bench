# Staleness pass

`/bench-implement-spec` charges this pass after the build preflight. The pass
audits an approved spec and its tickets against the current `main` tree.
Cheap-tier delegates audit the slices in parallel. The orchestrator reviews
their returns and amends the spec without a reviewer stop.

## The charge

- Line: the cheap tier, default effort, one iteration for each delegate.
- Access: read-only. A delegate writes no file.
- Partition: one delegate audits the spec body and its decision map. Each other
  delegate audits two or three tickets of one chunk or of adjacent chunks.
- Inputs: the spec path, the reviewed graph commit, the slice, and each
  preflight red row that the pass took.

Send every delegate in one dispatch, so that the slices run in parallel.

## Procedure for each delegate

1. List the drift. Run `git log --oneline <reviewed graph commit>..main` on
   each path that the slice names in a `Writes:` line or a claim. Add each file
   that a taken preflight red row names.
2. If the drift list is empty, return the verdict current and stop.
3. Examine each current-code claim of the slice against `main`: paths, symbols,
   flags, output shapes, test names, cited lines, and current-behavior premises.
   A claim is a current-code claim under the claim rules of `craft-spec`.
4. Find planned behavior that already shipped, in part or in full.
5. Find each contract on `main` that the spec predates and that changes what a
   ticket must build.
6. Make sure that each `Blocked by:`, `Covers:`, and `Writes:` line still holds
   after the code moved.

## Return

Each delegate returns one verdict for each spec part or ticket of its slice:
current, light refresh, rework, already shipped, or obsolete. For each stale
item, it returns the location, the claim, the current evidence as a line
citation or a commit, and the needed change. It also names what it did not
verify.

## Orchestrator review

The orchestrator reads every return before it changes the spec. When two
returns disagree, it examines the tree itself. It decides each finding by the
spec-contradiction predicate of `.bench/BENCH.md` and does not ask the reviewer.

The amendment corrects stale claims, citations, fences, and line numbers, and
it removes landed work. A late occurrence that the spec neither covers nor
excludes gets an `Out of scope` line or an acceptance row. A defect outside the
spec goes to `bench learning`. A finding outside the plan-expansion scope of
`.bench/BENCH.md` stops the phase.
