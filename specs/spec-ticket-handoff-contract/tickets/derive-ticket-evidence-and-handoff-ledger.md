# Derive ticket evidence and handoff ledger

Blocked by: author-identified-spec-handoffs.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`, `internal/conformance/docs_workflow_helpers_test.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: authored map and fence contract→author-identified-spec-handoffs.md; accountable handoff ledger and fence-drift result→require-approved-handoff-before-lifecycle.md; workflow predicates→`internal/anchors/registry_data.go`; fixture owner→`internal/conformance/fixture_bite_test.go`; focused helper→`internal/conformance/docs_workflow_helpers_test.go`; omission fixtures→`tests/canary/workflow-guidance-anchors`
Contracts: the handoff ledger crosses `.agents/skills/bench-craft-tickets/SKILL.md`→`.agents/commands/bench-implement-spec.md`; type is an ordered build-plan disposition, domain is every story, coverage row, named seam, edge disposition, and approved fence, order is the approved blocker-before-consumer plan whose first claimant is accountable, absence is an incomplete breakdown that cannot start, asserted by TL4 and TL5 against the real spec and ticket files
Closure: TL1/observed-red-route, TL2/already-covered-route, TL3/not-tdd-able-route, TL4/ledger-totality, TL5/fence-drift-stop

## What to build

Ticket derivation transforms every approved row classification into acceptance plus independent mutation evidence and emits one deterministic handoff ledger. Current-tree contract refresh may narrow a ticket, but a required path outside the approved fences returns to spec approval.

## Acceptance

- [ ] [TL1] (covers SH5) an observed-red row carries its failing public operation into the accountable ticket and adds a distinct post-implementation subject mutation.
- [ ] [TL2] (covers SH6) an already-covered row keeps its named control and adds a subject mutation proving the changed route reaches that control.
- [ ] [TL3] (covers SH7) a not-TDD-able row names its blocker, maps to the first ticket where its seam exists, and receives its subject mutation there.
- [ ] [TL4] (covers SH8) the handoff ledger accounts for every story, row, named seam, edge row or `Won't handle`, and approved fence before assignment.
- [ ] [TL5] (covers SH9) a required current-tree path outside every approved fence stops ticketing and returns the exact path and reason to `$bench-write-spec`.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TL1/observed-red-route | replace the distinct subject-mutation requirement with permission to copy the obsolete absence probe | the registered section-sensitive workflow anchor | apply the subject swap, run the focused workflow-anchor conformance test, require the observed-red-route diagnostic, restore the subject |
| TL2/already-covered-route | delete the changed-route subject mutation while retaining the named positive control | the registered section-sensitive workflow anchor | apply the subject omission, run the focused workflow-anchor conformance test, require the already-covered-route diagnostic, restore the subject |
| TL3/not-tdd-able-route | retain permanent mutation exemption after the seam exists | the registered section-sensitive workflow anchor | apply the subject swap, run the focused workflow-anchor conformance test, require the not-tdd-able-route diagnostic, restore the subject |
| TL4/ledger-totality | remove the approved-fence disposition from the ledger's enumerated domain | the registered section-sensitive workflow anchor | apply the subject omission, run the focused workflow-anchor conformance test, require the ledger-totality diagnostic, restore the subject |
| TL5/fence-drift-stop | replace the return-to-spec route with ticket-local fence widening | the registered section-sensitive workflow anchor | apply the subject swap, run the focused workflow-anchor conformance test, require the fence-drift-stop diagnostic, restore the subject |
