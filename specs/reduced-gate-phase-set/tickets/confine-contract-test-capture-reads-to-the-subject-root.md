# Confine contract-test capture reads to the subject root

Blocked by: Declare the reduced gate scope

Ownership fence: `internal/conformance/contract_capture_reads_test.go`
Assumptions: the contract phase runs the suite from the kit checkout against a separate subject root, and that root is what the stripped materialization replaces

## What to build

A static conformance check over the contract package: a contract test may read an
allowlisted path only relative to the subject root, never relative to the kit
checkout. This is the one place the design asserts rather than constructs, and it is
scoped deliberately — stripping reaches the subject root only, so a test resolving
`ROADMAP.md` from the kit checkout would keep reading the real tree and stay
observably coupled to a capture path with nothing looking.

Keep the check scoped to that one package and its resolution helper. Generalising it
to the whole tree turns a checkable assertion into a broad prohibition the rest of
the kit never agreed to, and the narrowness is what makes the assertion honest about
what it does and does not close.

## Acceptance

- [ ] [R12] A planted kit-relative read of an allowlisted path inside the contract package produces a diagnostic, while the same read resolved from the subject root does not.
