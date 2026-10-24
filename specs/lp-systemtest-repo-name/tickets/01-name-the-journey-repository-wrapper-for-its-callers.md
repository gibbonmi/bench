# Name the journey-repository wrapper for its callers

Blocked by: none
Writes: internal/systemtest/adoption_test.go, internal/systemtest/bench_follow_on_test.go, internal/systemtest/compatibility_test.go, internal/systemtest/guard_projection_differential_test.go, internal/systemtest/guard_projection_test.go, internal/systemtest/guard_refusal_test.go, internal/systemtest/land_route_fixture_test.go, internal/systemtest/land_route_test.go, internal/systemtest/owner_landing_fixture_test.go, internal/systemtest/owner_selection_test.go, internal/systemtest/owner_teardown_test.go, internal/systemtest/owner_test.go, internal/systemtest/process_lifetime_test.go, internal/systemtest/session_start_route_test.go, internal/systemtest/session_start_test.go, internal/systemtest/wrapper_pool_home_test.go, capture/restructure-backlog.md, capture/retros/landing-gate-long-poles.md
Covers: none

## What to build

The system suite has a private wrapper that makes a journey repository. The wrapper calls the system owner's `journeyRepository` method and fails the test on an error. Its old name tied it to the land route. Tests in many other files call it.

Rename the wrapper to `journeyRepo`. The name states what the wrapper returns. It also matches the method that it wraps and the helper of the same role in the worktree package.

Update every call site. Add a doc comment to the definition. Update the capture documents that name the old identifier. Do not change behavior.

Do not change `specs/shared-test-fixtures/inventory.json`. That file records candidate sites against a pinned base, and its spec refreshes it before implementation.

## Acceptance

- [ ] The system suite defines `journeyRepo` and calls it from every former call site.
- [ ] A hidden-file search for the old identifier finds no hit outside `specs/shared-test-fixtures/inventory.json`.
- [ ] `bench test --check system` passes.
