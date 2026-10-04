# Bind worktrees and shifts to one outcome

Blocked by: 03-protect-planning-and-commits.md
Writes: internal/commitment (new), internal/intent, internal/worktree, internal/shift, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/systemtest/otel_verbs_test.go, cmd/bench/main.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: DC23, DC24, DC25, DC50, DC69

## What to build

Integrate assignment creation, sibling inheritance, and shift startup with the shared admission owner. Preserve one outcome across the integration source and its ticket assignments.
Ticket 02 supplies durable claims. Ticket 03 supplies planning scope and preflight facts.
Default new assignments to planning until an explicit start binds delivery. An inherited sibling keeps the exact source outcome and deliverable.

Wire shift admission before shift intent acquisition or adapter execution. Route its explicit outcome operand through the existing grammar.
Update both the root help row in `cmd/bench/main.go` and the shift-owned help in `internal/shift/shift.go`.
Keep their existing help inventory expectation consistent in this checkpoint.
Replace the root help row without line growth; any needed extraction must create headroom in this ticket before the file grows.

Validate invocation and assignment identity first. An old assignment cannot borrow another assignment's approved continuation.
The continuation data shape is the exact assignment, request, and allowed scope; publication consumes it in ticket 05.

Update the existing shared worktree and shift fixture helpers for explicit admission. Enumerate each helper's callers before changing its posture.
Only the selected ordinary CLI is trusted at startup. Do not claim the wrapper manifest authenticates this route.

Read worktree CreateCommand and attribution, lifecycle and ownership helpers, shift grammar and run loop, and intent binding exports. Extract within the fenced package when an existing file lacks headroom. Review the binding seam before the publication consumer starts.

## Acceptance

- [ ] Two ticket assignments for A consume one outcome slot (DC23).
- [ ] A sibling based on A refuses an attempt to inherit B (DC24).
- [ ] Refused shift leaves both adapter marker and shift intent absent (DC25).
- [ ] An unlisted old assignment cannot use a listed assignment identity (DC50).
- [ ] Usage or assignment identity errors retain their precedence and leave no admission effects (DC69).

- [ ] Root help and shift-specific help advertise the same outcome operand, and the help inventory check passes before any successor ticket.

## Checkpoint verification

The existing system journey uses BENCH_KIT through `bench test --check system`. Adapt its fixture at this checkpoint when admission changes its route.

Run `bench test --package ./internal/worktree`, `bench test --package ./internal/shift`, and `bench test --package ./internal/intent`. Retain existing head, tip, and checked-out-ref cases where fixture assignment identity changes.
