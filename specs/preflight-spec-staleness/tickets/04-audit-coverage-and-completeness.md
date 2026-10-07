# Audit canonical coverage and complete supported inventory

Blocked by: 03-audit-test-citations.md
Writes: internal/preflight/completion_plan_test.go, internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/preflighttest/fixture.go, internal/preflight/staleness (new)
Covers: SS28, SS29, SS30, SS31, SS32, SS33, SS34, SS49, SS50, SS57, SS58

## What to build

The audit combines source, symbol, test, and canonical coverage collection.
Consume the reviewed ownership metadata parser and all predecessor collectors.
Add completeness.go, TestCoverageIdentity, and TestAuditCompleteness in identity_test.go.
Keep canonical rows-owned, rows-membership, duplicate-token, and completion-plan obligations.
Predicate bytes and changed ticket Covers join exact drift.

Missing, partial, duplicate, unsupported, or citation-incomplete metadata keeps Complete false.
A present empty inventory differs from missing metadata but cannot self-certify scope completeness.
Cross-check every supported citation.
Manual claims and unsupported collectors remain diagnostic; use no prose scanning fallback.
Authored authority or completeness grants no approval or semantic guarantee.

The command now supplies complete decision inputs.
The documented workflow remains manual until SS-C3.
Every collector supplies explicit completion or failure, never silent empty green.
Preserve ordinary charge and after-ticket applicability without changing ownership policy.

Use the landed NativeFuture descriptors without reclassifying them as supported Go seams.
The five exact pending native items keep Audit.Complete false and ManualRequired true, even when current premises permit charge.
TestAuditCompleteness drives this state through the real command.
Dropping these items or using diagnostic-only exit zero to omit manual review must fail.

This ticket belongs to SS-C2.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] SS28: An unowned coverage row retains red rows-owned.
- [ ] SS29: A phantom Covers token retains red rows-membership.
- [ ] SS30: Changed coverage predicate bytes join drift.
- [ ] SS31: Changed ticket Covers membership joins drift.
- [ ] SS32: Missing planning metadata reports incomplete audit evidence.
- [ ] SS33: Partial declared metadata keeps complete false.
- [ ] SS34: An omitted supported citation keeps complete false.
- [ ] SS49: Duplicate metadata refuses complete audit collection.
- [ ] SS50: A present empty declared claim set differs from missing metadata.
- [ ] SS57: A pending native-future manual entry keeps audit Complete false.
- [ ] SS58: A pending native-future manual entry keeps manual_required true.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/coverage`
- `bench test --check ticket-grammar`

Focused producer command: `bench test --package ./internal/preflight/... --run 'TestCoverageIdentity|TestAuditCompleteness'`.
Named mutation identity: `ss-04-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Omit the second citation, mask rows-owned with green symbols, or infer completeness from an empty container. The command inventory witness must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
