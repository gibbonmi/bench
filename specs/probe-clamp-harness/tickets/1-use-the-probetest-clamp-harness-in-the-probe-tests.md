# Use the probetest clamp harness in the probe tests

Blocked by: none
Writes: internal/probe/baseline_test.go, internal/probe/combined_test.go, internal/probe/delete_test.go, internal/probe/edit_test.go, internal/probe/invalid_cause_test.go, internal/probe/lock_test.go, internal/probe/omit_file_test.go, internal/probe/outcome_test.go, internal/probe/package_test.go, internal/probe/probe_test.go, internal/probe/probetest/fixture.go, internal/probe/refusal_test.go, internal/probe/replacement_test.go, internal/probe/result_test.go, internal/probe/set_input_test.go, internal/probe/set_test.go, internal/probe/subject_check_test.go, internal/probe/unwrap_test.go, internal/probe/verdict_test.go
Covers: none

## What to build

Consolidate the clamp repository fixture in `internal/probe/probetest/fixture.go`. Route the `internal/probe` package tests through the shared `probetest` clamp repository, source, selection, and canned events. Remove the duplicate fixture builders and canned constants from `internal/probe/probe_test.go`.

## Acceptance

- [ ] `rg -n "func newFixture|func clampSource" internal/probe/*_test.go` returns no match.
- [ ] `bench test --package ./internal/probe/...` passes with the same test count as before the change.
