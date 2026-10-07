# Audit existing and planned test functions precisely

Blocked by: 02-audit-symbol-identities.md
Writes: internal/preflight/decision.go, internal/preflight/gather.go, internal/preflight/gather_inputs.go, internal/preflight/preflighttest/fixture.go, internal/preflight/staleness (new)
Covers: SS20, SS21, SS22, SS23, SS24, SS25, SS26, SS27

## What to build

The actual audit distinguishes existing evidence from declared future test work.
Consume the source and symbol values.
Add tests.go and TestTestIdentity in identity_test.go.
Reuse the landed coverage file/function and supported named-subtest grammar.

An absent existing file differs from a missing function in a present file.
A planned function is absent at both requested initial identities; its file may hold unrelated tests.
An already-present supposedly planned function is already-shipped drift.
Unsupported syntax and wrong context remain diagnostic.
Declaration existence claims neither execution nor effectiveness.

Apply planned-function absence only to explicit initial audit-spec.
Post-ticket preflight and charge stay usable after a predecessor creates legitimate tests.
A later explicit audit compares its own pair and may report those tests as drift.
No automatic continuation exemption or evidence record value is created.
Incomplete coverage/completeness still keeps manual_required true.

This ticket belongs to SS-C2.
The explicit predecessor edge orders every overlapping write.
No unimplemented successor is required for this ticket's owned acceptance rows.

## Acceptance

- [ ] SS20: An absent existing test file reaches red cited-tests.
- [ ] SS21: A removed function in an existing file reaches red cited-tests.
- [ ] SS22: An absent planned test function remains valid future work.
- [ ] SS23: An existing supposedly planned function reports already-shipped drift.
- [ ] SS24: Other functions in a planned test's existing file remain allowed.
- [ ] SS25: A named supported subtest uses the canonical citation grammar.
- [ ] SS26: An unsupported test citation remains diagnostic.
- [ ] SS27: Wrong-context test evidence remains diagnostic.

## Verification

Run these future checks through their actual production entries.
Use a few sufficient multi-case witnesses; do not create a redundant test per row.
Retain the named omission mutation, its behavioral failure, restoration, and passing command result.
Compilation failure supplies no behavioral evidence.

- `bench test --package ./internal/preflight/...`
- `bench test --package ./internal/coverage`

Focused producer command: `bench test --package ./internal/preflight/... --run TestTestIdentity`.
Named mutation identity: `ss-03-omission`.
This obligation uses the named production seam and observable omission in this ticket, not a compile-failure substitute.

Mutation obligation: Require whole planned-file absence, repeat initial absence on continuation, or substitute file existence for function identity. The real initial and continuation witnesses must fail.

The completion plan names this ticket's focused command and mutation identity.
Execute bench probe against the specified production seam after its code exists, then restore and rerun that exact focused command.
If the planned mutation cannot reach that producer, amend the enabling plan before verification; no substitute earns completion.

Read the canonical structure report before dispatch.
Remeasure Growth from the actual preceding tip and whole-tree crowding.
Keep cohesive extraction and necessary pin relocation in this same ticket.
