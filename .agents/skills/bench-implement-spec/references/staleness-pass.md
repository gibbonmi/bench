# Staleness pass

`/bench-implement-spec` charges this pass after its initial `--audit-spec`
preflight, unless the omission rule below holds. The pass audits an approved
spec and its tickets against the current `main` tree.
Cheap-tier delegates audit the slices in parallel. The orchestrator reviews
their returns and amends the spec without a reviewer stop.

## Omission

Omit the pass only when the audit response states `complete` true and
`manual_required` false. An exit zero, an empty table, or the audit result of
another consumer does not permit that omission. An audit result never supplies
reviewer approval or the commitment.

## The charge

- Line: the cheap tier, default effort, one iteration for each delegate.
- Access: read-only. A delegate writes no file.
- Partition: one delegate audits the spec body and its decision map. Each other
  delegate audits two or three tickets of one chunk or of adjacent chunks.
- Inputs: the spec path, the reviewed graph commit, the slice, and each
  preflight red row that the pass took. The inputs also hold the audit items of
  the slice: every drift path, red check, and unsupported claim location of the
  audit response.
  Before the charge, the orchestrator resolves each `audit_diagnostics` subject
  to the location of that item in the planning inventory.

Send every delegate in one dispatch, so that the slices run in parallel.

## Procedure for each delegate

1. List the drift: each audit item of the charge. Add each file that a taken
   preflight red row names.
2. If the drift list is empty, return the verdict current and stop.
3. Examine each current-code claim of the slice against `main`: paths, symbols,
   flags, output shapes, test names, cited lines, and current-behavior premises.
   A claim is a current-code claim under the claim rules of `craft-spec`.
4. Find planned behavior that already shipped, in part or in full.
5. Find each contract on `main` that the spec predates and that changes what a
   ticket must build.
6. Make sure that each `Blocked by:`, `Covers:`, and `Writes:` line still holds
   after the code moved.

A delegate reads any deeper source that it needs to resolve one listed item.

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
`.bench/BENCH.md`, including a material acceptance change, stops the phase.
A taken `map-sources` red repoints each named record to a source that the
published tree holds. When that repoint edits a path outside the signed-off
`Writes:` lines, the path takes the widen rule of `.bench/BENCH.md`.

A gap in the planning metadata is not an approval. An amendment of that
metadata moves the reviewed subject, so it takes the enabling plan commit. The
orchestrator then reruns the audit. The reviewed graph commit stays the
`--base`, and the plan commit is the `--source-tip`, so the audit shows the
amendment as drift. The amendment never marks an unsupported claim green.
