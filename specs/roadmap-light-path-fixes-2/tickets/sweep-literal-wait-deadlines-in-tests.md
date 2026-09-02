# Sweep literal wait deadlines in tests

Blocked by: derive-test-wait-deadlines-from-bounds.md
Writes: internal/conformance/wait_deadline_literal_test.go (new), internal/conformance/marker_wait_deadline_test.go, internal/conformance/checks_test.go, internal/conformance/registry/registry.go
Covers: LP11, LP12, LP13

## What to build

Verify the premise first: `checkMarkerWaitDeadlines` in
internal/conformance/marker_wait_deadline_test.go walks `cmd/` and `internal/`,
and it reds a numeric literal only in the marker-wait helper's deadline
argument. Add one dev-tier check that keys on a test wait's deadline argument.
That argument is a `time.After` or `time.Now().Add` inside a `_test.go` wait.
The check reds a numeric literal there.

Reuse `containsNumericLiteral` from the marker-wait file and `expressionText`
from `bounds_policy_test.go` in the same package; share them rather than copy
them. Allow a poll-interval literal inside the wait loop. Register the
check beside the marker-wait row at the same tier, subject, and inputs. The
live tree, after the blocker ticket, is the green set.

## Acceptance

- [ ] The bite test reds a synthetic test file whose wait deadline is `5 * time.Second`.
- [ ] The bite test stays green for a synthetic wait with a derived deadline and a `10ms` tick.
- [ ] The check body calls the marker-wait literal scanner and defines no second scanner.
- [ ] The check runs at the dev tier and passes on the live worktree.
- [ ] Self-probe: add a literal deadline to one migrated helper, and report the check red with that file named.
