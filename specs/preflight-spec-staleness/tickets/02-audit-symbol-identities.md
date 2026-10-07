# Audit cited symbol identities through the typed loader

Blocked by: 01-audit-source-drift.md
Writes: internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/preflighttest/fixture.go, internal/preflight/staleness (new)
Covers: SS13, SS14, SS15, SS16, SS17, SS18, SS19

## What to build

The real audit response grades symbol discrepancies through the landed consumers loader.
Consume SS-C1's reviewed AuditFacts and explicit source bytes.
Add identity.go and TestSymbolIdentity in identity_test.go.
Compare declared identity against baseline AST syntax and current typed origin.
Keep owner path, package, declaration kind, receiver, and signature exact.

Missing, moved, renamed, wrong-kind, and same-named replacements report their actual discrepancy.
Aliases follow canonical origin.
Ambiguous bare names, reflection, plugins, non-Go claims, and non-default contexts stay diagnostic.
Syntax or constant equality proves no semantic equivalence or caller effectiveness.
No arbitrary historical-context loader or lexical fallback is authorized.

Return cited-symbols through the established command projection.
Test and coverage collectors remain incomplete, so manual_required remains true.
No successor test is required to prove this ticket.

This ticket belongs to SS-C2.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] SS13: An absent cited symbol reaches a red cited-symbols result.
- [ ] SS14: A cited signature change reaches a red cited-symbols result.
- [ ] SS15: A cited symbol moved to another owner reports that discrepancy.
- [ ] SS16: A wrong declaration kind reports its exact mismatch.
- [ ] SS17: An ambiguous bare symbol remains diagnostic.
- [ ] SS18: An alias resolves through the canonical origin rule.
- [ ] SS19: A non-default context claim remains diagnostic.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/consumers`

Focused producer command: `bench test --package ./internal/preflight/... --run TestSymbolIdentity`.
Named mutation identity: `ss-02-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Resolve a retained name to the wrong owner or drop receiver/signature comparison. The typed-source identity witness must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
