# Spec-build review and gate cadence

Status: shaping

## Destination

A successor spec-build workflow in which each ownership-fenced ticket earns
focused executable evidence and a durable orchestrator review before provisional
integration. The exact composed candidate is pinned for one fresh semantic
review, and promotion alone pays the whole-project gate and authors the
landing. The successor consumes FT173's full AXI lifecycle contract and
FT130's capture accounting contract after both have shipped. It does not
reopen the terminal single-build serial-gate work.

## Notes

## Decisions so far

- [What evidence admits a ticket into the composed candidate?](spec-build-review-gate-cadence/tickets/1.md): Each assignment worktree and its assignment branch are one provisional construct.
- [What may ticket-local checks compile and execute?](spec-build-review-gate-cadence/tickets/2.md): Ticket scope is defined by its acceptance rows and declared integration surfaces, not by a file-only compilation boundary.
- [How are ticket-local conformance and canary checks selected?](spec-build-review-gate-cadence/tickets/3.md): Ticket-local conformance runs only the exact checks that own the ticket's changed inputs or declared integration surfaces.
- [What transition pins the composed candidate for final review?](spec-build-review-gate-cadence/tickets/4.md): Final review is a two-step lifecycle transition.
- [Which harness phase owns final review and promotion?](spec-build-review-gate-cadence/tickets/5.md): `$bench-implement-spec --full` remains the end-to-end wrapper and composes the standalone `$bench-review-implementation` phase immediately before promotion.
- [What does review, final-check, commit, and promote each own?](spec-build-review-gate-cadence/tickets/6.md): `$bench-review-implementation` asks whether the exact diff is good on Standards, Spec, and Coverage.
- [Which owner supplies AXI for the lifecycle family?](spec-build-review-gate-cadence/tickets/7.md): FT173 is a separate required predecessor.
- [Which owner accounts for capture after landing?](spec-build-review-gate-cadence/tickets/8.md): FT130 is a separate required predecessor and the one owner of capture queuing, pending-capture visibility.
- [How is the revised cadence accounted for?](spec-build-review-gate-cadence/tickets/9.md): The revised cadence receives a successor spec after FT173 and FT130 ship.

## Not yet specified

## Spec-writer discretion

- Internal type placement and naming for the pinned review generation, provided
  the public `review --begin` and `review --evidence` behavior and invalidation
  predicates remain exact.
- Exact focused Go test filters and package groupings within each ticket's
  acceptance rows and declared integration surfaces. They must not widen into
  a whole-project gate or omit an affected registered owner.
- Receipt field and TOON table layout only where FT173 leaves reversible freedom;
  the cadence spec must consume, not fork, FT173's shared AXI contract.

## Out of scope

- Reopening, repairing, or amending the terminal single-build serial-gate run or
  its spec.
- A whole-project gate or ordinary commit in an assignment worktree.
- Treating semantic review as deterministic gate authority or project-green
  evidence.
- Replacing FT173's AXI foundation, FT185's gate-result payload, or FT130's
  capture-accounting contract inside the successor.
- Durable cross-run Bench executable storage, later-process executable reuse,
  release-tier proofs, or a broader gate-scheduler rewrite.

## Sources

- Path: `ROADMAP.md`
  Supports: #7 and #8 existing FT173 and FT130 ownership.
  Drift: re-read after either roadmap item is reshaped, staged, or retired.
- Path: `CHANGELOG.md`
  Supports: #1's checkpoint receipt and #4/#5's candidate-bound review receipt, whose provisional lifecycle was removed wholesale and now survives only as this removal record.
  Drift: re-read before specifying the successor; its seams re-derive from the surviving landing path, not the deleted lifecycle.
- Path: `internal/landing/landing.go`
  Supports: #4 and #5's surviving composition-and-publication operation, now owned by `bench worktree land`.
  Drift: re-read if composition, gating, or publication ordering changes.
- Path: `projects/benchkit.md`
  Supports: #2 and #3 current process, conformance, and canary ownership.
  Drift: re-read if the gate architecture, process-seam inventory, or conformance registry contract changes.
- Path: `.agents/commands/bench-final-check.md`
  Supports: #6 current workflow-dependent final-check behavior.
  Drift: re-read if final-check or either landing route changes before spec authoring.
