# Expose the test-check inventory

Blocked by: none
Writes: internal/testreport/command.go, internal/testreport/check_test.go
Covers: LF25

## What to build

Give bench test help one canonical named-check inventory. An unknown check names
its operand and prints that same inventory instead of generic usage.

## Acceptance

- [ ] Test help lists every supported named check.
- [ ] An unknown check reports its exact operand.
- [ ] Help and refusal use one inventory owner.

