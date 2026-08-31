# Reject invalid prospective temporary roots

Blocked by: single-source-prospective-checkout-name.md
Writes: internal/gate/prospectiveartifact/prospectiveartifact_test.go
Covers: LF22

## What to build

Drive Open against an absent temporary root and a dangling temporary-root
link. Pin fail-closed behavior and absence of partial publication.

## Acceptance

- [ ] Open rejects an absent temporary root.
- [ ] Open rejects a dangling temporary-root link.
- [ ] Neither refusal leaves a published record.

