# Retain every ledger-accountability member

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/anchors/registry_data.go`, `internal/conformance/fixture_bite_test.go`, `tests/canary/workflow-guidance-anchors`
Integration surfaces: ledger member clauses→`internal/anchors/registry_data.go`; hostile mutations→`tests/canary/workflow-guidance-anchors`; retained membership→`internal/conformance/fixture_bite_test.go`
Contracts: the handoff ledger's independent members cross `.agents/skills/bench-craft-tickets/SKILL.md`→`internal/anchors/registry_data.go`→`tests/canary/workflow-guidance-anchors`; type is section-scoped Markdown predicates, domain is current spec/ticket artifacts plus deterministic first claimants and every approved fence, order is build-plan order, absence fails the member fixture and aggregate
Closure: LM1/current-artifacts, LM2/accountable-first-claimant, LM3/fence-owner-or-unused, LM4/retained-memberships

## What to build

Close accepted review finding `COV-002-ledger-member-mutations` by giving SH8's three independently removable ledger-accountability members their own registered section-sensitive mutations and retaining all three fixtures in the existing aggregate. Keep `craft-tickets` as sole ledger owner; add no parser or lifecycle authority.

## Acceptance

- [ ] [LM1] (covers SH8) the ledger must derive from current spec and ticket artifacts rather than copied totals.
- [ ] [LM2] (covers SH8) deterministic build-plan order must identify the accountable first claimant while later tickets may reinforce it.
- [ ] [LM3] (covers SH8) every approved fence must identify one owning ticket or an explicit unused disposition.
- [ ] [LM4] (covers SH11) the aggregate retains all three new fixture names and exact diagnostics with an updated cardinality.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| LM1/current-artifacts | replace the current-artifact/no-copied-totals predicate with permission to copy counts | the registered section-sensitive workflow anchor | apply the swap, run its focused fixture, require the current-artifact diagnostic, restore the clause |
| LM2/accountable-first-claimant | remove the accountable first-claimant predicate while leaving multiple-claimant permission | the registered section-sensitive workflow anchor | apply the omission, run its focused fixture, require the accountable-claimant diagnostic, restore the clause |
| LM3/fence-owner-or-unused | relocate the owner-or-explicit-unused predicate outside ledger review | the registered section-sensitive workflow anchor | apply the relocation, run its focused fixture, require the fence-disposition diagnostic, restore the clause |
| LM4/retained-memberships | remove each new member separately from the aggregate list | `TestSpecTicketHandoffWorkflowFixturesAreComplete` | apply one omission at a time, run the named aggregate, require its cardinality red, restore before the next omission |
