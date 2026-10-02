# 2. Carry an ambient value below each verb entry

Blocked by: 1-read-the-kit-value-once-in-gate.md
Writes: internal/worktree/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: WS9, WS10, WS11, WS12, WS13, WS14, WS15, WS16

## What to build

Chunk: SR-C2.

Add the `ambient` value: the Bench home, the kit value, one instant, and a warnings
writer. Add one `effects.go` constructor that reads the kit through `gate.KitValue` and
the clock through `currentTime`. The constructor takes the home and the stderr writer from
its caller. Each verb entry with an internal form calls the constructor once. Each
internal form takes the joins value first and the ambient value second.

Remove the `now` and `home` fields from the joins value. The clean discard path reads the
ambient instant. The retirement path reads the ambient home. `resumeCleanCommandWith` and
`releaseAssignment` replan at the ambient instant, so one run reads the clock once.
`ListCommand` uses no joins home.

Change the verb runner. Build the ambient value with the same constructor, and let a
call's kit value and clock value replace its reads. A call with no joins value, no kit
value, and no clock value runs the public entry. Every other call runs the internal form,
with `defaultJoins()` when the call holds no joins value. `checkVerbCall` refuses a kit
value or a clock value for a verb key without an internal form, with
`verb runner: the <key> verb takes no kit value` or
`verb runner: the <key> verb takes no clock value`.

Move the discard tests from `discardJoins` to the call's clock value. Drop the
`BENCH_HOME` bind from the two resume-clean tests in `resume_test.go`, and pass the home
through the verb runner. Add the resume-clean home test and the clock-value test in a new
test file, because `resume_test.go` is over its line budget.

## Acceptance

- [ ] Release, clean, and land drop census records under the home that the verb receives.
- [ ] resume-clean under a non-default home drops the retired assignment's census record under that home.
- [ ] resume-clean with a clock value eight days after an active assignment's creation names the clean command for its path.
- [ ] A discard with the clock value at `discardDay` writes the ref that `intent.DiscardedRef(discardDay, ref)` names.
- [ ] `checkVerbCall` refuses a kit value and a clock value for the `path` key with the two messages above.
- [ ] No over-budget file in the package grew.
