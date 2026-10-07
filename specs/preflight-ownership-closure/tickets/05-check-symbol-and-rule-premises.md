# Check symbol reachability and exact rule premises

Blocked by: 04-check-sizing-premises.md
Writes: internal/preflight/closure.go, internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/ownership (new), internal/preflight/preflighttest/fixture.go, internal/preflight/proposal.go
Covers: OC39, OC40, OC41, OC42, OC43, OC44, OC45, OC46, OC47

## What to build

The real proposal checks exact symbol, import, anchor-kind, sentinel, error, and bound premises.
Consume the reviewed loader, Source, and sizing values.
Add symbols.go and TestSymbolPremises in premises_test.go.
Normalize declarations through the consumers AST and Go formatter owner.
Keep owner path, package, receiver, kind, spelling, and signature exact.

Grade import reachability from the named calling seam, refusing private cross-package access and cycles.
Use canonical anchor kinds and actual sentinel/bound declarations.
A familiar name or same-valued local literal cannot substitute for its owner.
Opaque expressions, dynamic errors, and ambiguous references remain diagnostic.
Exact construction does not prove the failure reaches its intended caller.

Report supported results through the established proposal route.
Do not certify unbuilt record/predicate collectors for charge readiness.
Add no parser, policy registry, or source-authentication path.

This ticket belongs to OC-C3.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] OC39: A named symbol resolves to its declared owner file.
- [ ] OC40: A declared signature mismatch refuses that premise.
- [ ] OC41: A calling seam cannot import an unexported cross-package helper.
- [ ] OC42: A cycle in the proposed import edge refuses that premise.
- [ ] OC43: A supported anchor kind resolves from the canonical inventory.
- [ ] OC44: A sentinel premise resolves its exact owner declaration.
- [ ] OC45: An unresolved failure expression stays diagnostic.
- [ ] OC46: A bound premise resolves its owning static value.
- [ ] OC47: An opaque bound expression remains diagnostic.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/consumers`
- `bench test --package ./internal/anchors`

Focused producer command: `bench test --package ./internal/preflight/... --run TestSymbolPremises`.
Named mutation identity: `oc-05-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Resolve by name alone, allow a private edge or cycle, or certify an opaque bound. The actual calling-seam or unknown-state witness must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
