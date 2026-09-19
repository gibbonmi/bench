# 8. Name the window of every production wait

Blocked by: 7-switch-verdict-windows.md
Writes: internal/conformance/bounds_policy_test.go, internal/conformance/fixture_bite_test.go, tests/canary/package-core-guard/, internal/bounds/bounds.go, internal/git/git.go, internal/handoffdoc/, internal/capturetx/, internal/gate/runner.go, internal/runbinary/runbinary.go, internal/worktree/subshell.go, internal/worktree/exec.go, internal/testreport/command.go, internal/freshness/freshness_publish.go, internal/releaseevidence/release_evidence.go, internal/shift/loop.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: TD40, TD41

## What to build

Chunk: TD-C3.

The bounds-policy check grows one rule. In production code outside the bounds package, each timed wait passes its window through one of the two bounds accessors. The rule covers each context timeout or deadline, each timer, each after-channel, each after-function, and each deadline computed from the current time.

Move the ref-check timeout, the handoff lock deadline, and the capture lock wait into the bounds registry. Read each one through the first accessor, and add each one to the owner table with the first accessor around it. Pass each cancel grace, each poll interval, and the operator's wall limit through the second accessor. `TestUpdateRefusesALockAnotherWriterHolds` gets its own lock window through a setter. Two new canaries in the `package-core-guard` family prove the rule.

## Acceptance

- [ ] The bounds-policy check reds a production file that passes a raw local duration to a context timeout.
- [ ] The bounds-policy check reds a production file that computes a deadline from the current time with a raw local duration.
- [ ] The owner table names the first accessor around the ref-check, the handoff lock, and the capture lock entries. The second accessor around any of the three reds the check.
- [ ] The live tree passes the bounds-policy check.
- [ ] `TestUpdateRefusesALockAnotherWriterHolds` still sees the refusal with the switch at `1`.
