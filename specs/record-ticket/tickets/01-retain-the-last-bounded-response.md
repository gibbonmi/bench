# Retain the last bounded response

Blocked by: none
Writes: cmd/bench/census_output.go, cmd/bench/census_output_test.go, cmd/bench/retain.go (new), cmd/bench/retain_test.go (new), cmd/bench/response_bound_test.go, cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand/, internal/responsebound/owner.go, internal/responsebound/owner_test.go, internal/responsebound/store.go, internal/responsebound/retain.go (new), internal/responsebound/retain_test.go (new), internal/responsebound/retire_test.go, internal/responsebound/responseboundtest/checkout.go, internal/git/status.go, internal/git/facts_test.go, internal/handoff/facts.go, internal/status/status.go, tests/canary/docs-currency-token-diet/, internal/roadmap/context_parse.go, internal/shift/shift.go, internal/worktree/land_identity.go, internal/worktree/response_spill_test.go, internal/gate/gate_go.go, internal/systemtest/exec_nested_bound_test.go
Covers: RT1, RT2, RT3, RT4, RT5, RT6, RT7, RT8, RT65

## What to build

After `Finish`, `runBounded` writes the retained response of the assignment scope that `AssignmentScope` answers.
The retained response holds the complete response bytes in arrival order, the exit code, the typed verb word, and its arguments.
It also holds the `HEAD` commit of the worktree and its dirty paths, read after the verb returns.

`responsebound` gains one accessor for the complete bytes of a finished response.
The accessor answers the retained writes of a response within the bound and the spill file bytes of a spilled response.
It answers no bytes for the spill-failed state and for the `passing` state.
`Finish` sets the retained writes to nil today, so the accessor reads them before that reset.
The lead identity row is not part of the retained bytes.

The only caller of `runBounded` passes the arguments after the verb word.
That caller passes the typed verb word too, so an alias reaches the retained response as the operator typed it.

The dirty-path read moves into one `internal/git` function.
That function runs `git status --porcelain=v1 -z --no-renames` and `ParsePorcelainZ`.
`LandedState`, `Facts`, and retention call it.
The commit comes from `rev-parse HEAD`.
Both reads run only when retention writes.

These calls leave the retained response unchanged:

- a call in the primary checkout or outside a repository
- a call of the `record` family
- a call of the `worktree` family
- a call that a retiring leaf runs

A response with no complete bytes clears the retained response.
Retention writes through one durable replacement, and its error changes neither the response nor the exit code.
The retained file lives under the assignment scope, so `Drop` removes it at release.
Its name does not match the spill name pattern, so the primary prune never reads it.

The system test for RT5 runs a real child process, so this ticket sets `BENCH_KIT` for the system suite through `bench test --check system`.

Before the ticket locks the skip declaration, it reads the headroom of `cmd/bench/command_registry.go` (396 of 400 lines) and `cmd/bench/main.go` (over its budget).
A name comparison in a new `cmd/bench` file adds no line to either file.

## Acceptance

- [ ] A `bench test` call within the bound in a pool worktree leaves a retained response with its exact bytes.
- [ ] A spilled `bench test` response leaves a retained response equal to the spill file bytes, with no `spilled{` frame.
- [ ] The retained response holds the exit code and the argument list after `bench`.
- [ ] A `bench record` call after a `bench test` call leaves the `bench test` response retained.
- [ ] A `bench worktree exec` call leaves its child's `bench test` response retained, in the system suite.
- [ ] A bounded call in the primary checkout writes no retained response.
- [ ] A spill-failed response and a `passing` response each clear the earlier retained response.
- [ ] A failed retention write leaves the printed response and the exit code unchanged.
- [ ] `LandedState` and `Facts` call the shared porcelain reader, and their existing tests stay green.
