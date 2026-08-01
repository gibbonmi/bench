# Spec-integration gate cadence — implementation retro

Spec: `specs/spec-integration-gate-cadence/spec.md` (implemented)

## A ticket must be an executable contract, not an umbrella claim

The largest source of rework was not missing prose in the spec; it was tickets that
said a broad behavior was covered without naming the mutation and oracle that proved
each acceptance row. Future tickets should bind every row to one concrete red
mutation, its independent owner, and the exact end-to-end lifecycle sequence. A
claim such as "checks every fact" is insufficient when delete-only probes miss
swaps, stale retained identities, or a second-call failure.

Ticket fences must include every independent registry, release, and conformance
owner needed to make the ticket green. They must also be re-derived from the tree
after earlier tickets land: defect inventories, expected base hashes, and ownership
assumptions age immediately. The ticket should state clean-checkout prerequisites,
ignored build artifacts, receipt vocabulary, source-binding expectations, structure
headroom, and the exact public operation sequence it must prove.

## Lifecycle tests must cross process and phase boundaries

Unit-level success hid failures that appeared only after state was serialized and a
fresh CLI process loaded it. Recomposition tests need to cover the whole sequence:
working-branch advance, compatible same-path replay, project-green advancement,
retained assignment provenance, fresh reload, fresh review, and the second Promote.
Stopping after the first Promote missed both the mutable-base run-identity defect and
the Review/Promote precondition deadlock.

Compatibility should be decided by applying the exact patch, not by rejecting every
shared pathname. Genuine content conflicts still fail closed. Stable lifecycle
identity is the canonical SHA-256 run ID bound to its exact candidate ref; mutable
base tips must not redefine it.

## The final gate has one canonical owner

`gate-run --fresh` prints a valid phase result but does not publish the project-green
evidence consumed by promotion; `bench gate` is the canonical evidence-producing
entry. Treating worktree, main-root, and prospective checks as interchangeable caused
several redundant full runs. A ticket that changes gate cadence must name which
command authors evidence and which phase consumes it.

Prospective execution must reproduce a clean linked checkout, including construction
of ignored bootstrap artifacts from the exact unpublished source. Its source-binding
test must commit version A and gate unpublished version B. Direct ordinary gate entry
must remain sealed, policy identity must change when execution semantics change, and
bounded diagnostics must survive a red prospective run.

## Release outcome

The canonical gate was green on `5ae1540`. The exact reviewed Story 9 patch was then
committed directly as `5e0c347` and the spec status set to `implemented`, following
the explicit instruction that this be the last gate. The final commit therefore did
not receive a separate prospective gate. Retained spec-build bookkeeping was closed
as terminal after the Review/Promote precondition deadlock; project-green remains an
honest marker for the last lifecycle-confirmed base rather than being advanced by
hand.

Broad Git-command blocking is not the remedy. Keep destructive hooks narrow and put
workflow authority in durable lifecycle preconditions, with process-boundary tests
that prove those preconditions compose.
