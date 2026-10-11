# Grant one landing a reviewer override from a terminal

Blocked by: 6-open-and-close-chunk-gaps.md
Writes: .bench/hooks/block-bench-follow-on.sh, DATA_HANDLING.md, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/commitment_command.go, cmd/bench/commitment_test.go, cmd/bench/guards.go, cmd/bench/help_inventory_test.go, cmd/bench/isolated_command_fixtures_test.go, cmd/bench/main.go, cmd/bench/response_bound_exempt_test.go, internal/anchors/registry_commitment.go, internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go, internal/benchguard/reviewer.go (new), internal/commitment/admission.go, internal/commitment/commitcmd/command.go, internal/commitment/commitcmd/override.go (new), internal/commitment/commitcmd/override_route_test.go (new), internal/commitment/commitcmd/override_test.go (new), internal/commitment/commitcmd/refusal_route_test.go, internal/commitment/repository/admission.go, internal/commitment/repository/override.go (new), internal/commitment/repository/publication.go, internal/commitment/repository/publication_lock_test.go, internal/commitment/repository/publication_test.go, internal/commitment/repository/window.go (new), internal/commitment/repository/window_test.go (new), internal/conformance/axi_query_registry_test.go, internal/conformance/refusal_route_guard_test.go, internal/conformance/subcommand_routing_table_test.go, internal/guards/guards_test.go, internal/intent/ledger/commitment.go, internal/refusalroute/faces_commitment.go, internal/systemtest/bench_follow_on_test.go, internal/worktree/land.go, internal/worktree/land_override_test.go (new), internal/worktree/parallel_census_test.go, tests/canary/data-handling-derivation/undocumented-passlist-var, tests/canary/package-core-guard/unrouted-subcommand
Covers: BW71, BW72, BW73, BW74, BW75, BW76, BW77, BW78, BW79, BW80, BW81, BW82, BW83, BW84

## What to build

A reviewer grants one landing a recorded override with a reason, from a terminal in the primary checkout.
The override admits that one landing through every open window, and its landing consumes it at publish.
No agent route grants it.

`bench commitment override --assignment <id> --reason <text>` is a new commitment form.
The reason follows the one rule that `bench commitment block` applies; extract that rule from `internal/commitment/admission.go:151` into one exported rule that both verbs call.
Each refusal renders through the commitment form's refusal renderer, so it prints `error: bench commitment override refused — <reason>`.
The verb refuses in this order, each with exit 1:

1. A root that is not the primary checkout: `commitment override runs from the primary checkout`.
2. A stdin that `terminal.IsTerminal` rejects: `commitment override needs the reviewer at a terminal; run it in your own terminal`.
3. An assignment that is not active: `commitment override assignment <id> is not active`.
4. A confirmation line that differs from `<id>`, or EOF: `commitment override not confirmed`.

Before the confirmation read, the verb prints one `override` row with the assignment, its label, the refusing open builds, and the reason.
It then prints the prompt `type the assignment id to confirm:`.
On success it writes the override record in one ledger decision and prints `overridden{assignment=<id>}`.

The prompt must reach the terminal before the read.
So give the form a stdin-aware entry and a bound exemption instead of the buffered output adapter.
Add a separate stdin-aware entry rather than a new `commitcmd.Command` signature, because that signature has about 29 test callers.
`cmd/bench/main.go` is over its line budget, so a line added there moves the same headroom out in this ticket.

`PublishAdmitted` removes the override record of the publication assignment in the publish decision.
A landing that refuses before publish keeps the record.
Ticket 5's window already reads the override record.

The follow-on guard (`internal/benchguard`) denies a Bash tool call whose Bench segment runs `commitment override`.
It also denies `bench commitment --in primary override` and the child form under `bench worktree exec`.
The `--in primary` form from a worktree passes the primary-checkout test, so the guard is its only layer.
`internal/benchguard/benchguard.go` has about 13 lines of headroom, so put the rule in `internal/benchguard/reviewer.go`.
The guard manifest line in `.bench/hooks/block-bench-follow-on.sh` names the reviewer verb.

`internal/systemtest/bench_follow_on_test.go` is system-tagged; if it gains a real-hook row, run it with `BENCH_KIT` set through `bench test --check system --run <name>`.
Ticket 8 updates the follow-on guard's glossary entry.
Raise `worktreeTestCount` in `internal/worktree/parallel_census_test.go` in place for each new top-level test.

Test the terminal path through an injected predicate and an injected confirmation reader, as `bench setup` injects `isTerminal` (`internal/adopt/setup.go:25`).
A non-terminal row drives a pipe, because `/dev/null` reads as a character device.

## Acceptance

- [ ] An override for an active assignment writes one record with its reason and prints `overridden{assignment=<id>}`.
- [ ] A landing whose assignment holds an override publishes while two builds are outside their gaps.
- [ ] A landing of another assignment still refuses while the override exists.
- [ ] After the overridden landing publishes, the ledger holds no override for its assignment.
- [ ] After the overridden landing refuses red at its gate, the override record stays.
- [ ] With a pipe on stdin, the verb refuses with `needs the reviewer at a terminal` and writes nothing.
- [ ] With `/dev/null` on stdin, the verb refuses with `not confirmed` and writes nothing.
- [ ] With an injected terminal and a confirmation that differs in one byte, the verb refuses with `not confirmed`.
- [ ] From a Bench worktree, the verb refuses with `runs from the primary checkout` before the terminal check.
- [ ] An override for an unknown assignment refuses with `is not active`.
- [ ] A reason with a newline, and an empty reason, each refuse with the block verb's reason refusal.
- [ ] The guard denies `bench commitment override --assignment x --reason y`, `bench commitment --in primary override`, and `bench worktree exec t -- bench commitment override`.
- [ ] The help inventory lists the `bench commitment override` form.
