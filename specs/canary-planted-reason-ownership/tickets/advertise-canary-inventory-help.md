# Advertise canary inventory help

Blocked by: report-truthful-canary-inventory.md
Writes: `cmd/bench/main_test.go`, `bin/bench.sh`

## What to build

Update the real wrapper command list and top-level route expectation to advertise fixture inventory validation while retaining the existing canary positional grammar, help forms, arity refusals, and exit classes.

## Acceptance

- [ ] (covers CI6) The root positional grammar, help forms, arity, and exit 0/1/2 classes remain unchanged, while real top-level help advertises fixture inventory validation.
