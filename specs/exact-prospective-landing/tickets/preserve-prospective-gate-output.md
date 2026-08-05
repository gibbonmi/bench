# Preserve prospective gate output through landing

Blocked by: build-exact-landing-owner.md
Ownership fence: `internal/gate/authorization`, `internal/landing`
Integration surfaces: exact-tree gate writers→`internal/gate/authorization` asserted by PO1; landing request output contract→`internal/landing` asserted by PO2; commit adapter consumer→adopt-exact-landing-in-commit.md; spec-build authorization consumer→existing `internal/specbuild/authorization_gate.go` plus PO3
Contracts: stdout and stderr writers cross `internal/commit`→`internal/landing`→`internal/gate/authorization`, asserted by PO1-PO2 against the real prospective gate; authorization result crosses `internal/gate/authorization`→`internal/landing`, asserted by PO1-PO3 against the existing classification and evidence contract

## What to build

Let an exact landing stream the prospective gate's public output to its caller without duplicating authorization policy, while retaining the discard-output authorization entry used by lifecycle consumers that do not expose gate output.

## Acceptance

- [ ] [PO1] The writer-aware authorization entry sends the supplied stdout and stderr to the real exact-tree gate and returns the same green, candidate, inherited, or infrastructure classification and evidence contract.
- [ ] [PO2] A landing request carries stdout and stderr through its owner to writer-aware authorization, so reused-green and fresh gate diagnostics reach the caller without a second gate invocation.
- [ ] [PO3] The existing discard-output authorization entry remains the single compatibility route for spec-build authorization and preserves its current behavior.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PO1 | replace the supplied gate writers with `io.Discard` | authorization writer-contract test | authorize a real prospective tree with capturing writers and expect the gate verdict text to disappear under mutation |
| PO2 | omit either landing request writer when calling authorization | landing output-plumbing test | land through a capturing request and expect the exact gate diagnostic on its original stream without an additional authorization call |
| PO3 | make the compatibility entry maintain a second classification path | authorization compatibility test | compare discard-entry and writer-aware results for the same enumerated outcomes and expect identical kind and evidence |
