# Run the gap tests as root and pin the tip fault

Blocked by: none
Writes: internal/gate/chunk_gap_test.go, internal/commitment/repository/gap_test.go
Covers: none

## What to build

When the `0o200` ledger stays readable, `TestUnrecordedGapExitsOne` in `internal/gate/chunk_gap_test.go` skips through `capability.Capability` with `capability.Privilege` and a reason. When the `0o200` ledger is unreadable, the test still verifies the exit 1 refusal and the rerun. In `internal/commitment/repository/gap_test.go`, a test pins that `OpenGap` and `CloseGap` return an error and write no gap when the default-branch tip does not resolve.

## Acceptance

- [ ] `TestUnrecordedGapExitsOne` routes a readable `0o200` ledger through `capability.Capability` with `capability.Privilege` and a reason instead of failing.
- [ ] `OpenGap` and `CloseGap` return an error and write no gap when the default-branch tip fails to resolve.
