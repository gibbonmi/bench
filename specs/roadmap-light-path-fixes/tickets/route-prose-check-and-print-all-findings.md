# Route the prose check and print all findings

Blocked by: expose-test-check-inventory.md
Writes: internal/testreport/command.go, internal/testreport/check_test.go
Covers: LF26

## What to build

Add prose to the named-check inventory and route it to the canonical grader.
Preserve the existing full finding loop without a second prose parser.

## Acceptance

- [ ] bench test --check prose reaches the canonical prose grader.
- [ ] A multi-file failure prints every finding.
- [ ] The named route adds no second prose policy owner.

