# Sweep input constructors and give the fence rule one owner

Blocked by: none
Writes: .agents/skills/bench-craft-spec/SKILL.md, .agents/skills/bench-craft-spec/references/map-discipline.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go, tests/canary/workflow-guidance-anchors/craft-spec-fence-after-slice/
Covers: none

## What to build

This ticket is the FT300 residual kit edit. The `slicing-closure` build paid three
repair rounds for readers outside the `Writes:` fences. The reader sweep searched
callers of changed functions. It did not search the constructors of the input that a
new check grades.

The `map-discipline` reference gets four rules under "Before the map locks":

1. A new invariant check sweeps every input constructor.
2. Each input constructor satisfies the invariant or joins a `Writes:` line.
3. A change to a preflight row names every reader of the row count or row list.
4. A change to the output of a fixture builder names every test that asserts or rewrites that output.

The `slicing-checks` reference gets two slicing rules. One rule requires each claim to
name its producer and check at its checkpoint. The other rule requires a ref preserve or
move to check each ref shape.

The `craft-spec` fence paragraph restated the fence-equals-union rule that the slicing
checks state. It now points to the slicing checks. Its anchor row, its test expectation,
and its canary move with the sentence. Each new rule gets one anchor row and one
independent test expectation.

## Acceptance

- [ ] `craft-spec` no longer states the union rule, and it links to the slicing checks for the fence rule.
- [ ] `map-discipline.md` states the four input-constructor rules, and each rule has an anchor row.
- [ ] `slicing-checks.md` states the producer rule and the ref-shape rule, and each rule has an anchor row.
- [ ] If you delete a new rule sentence from the live tree, `bench test --check docs-currency-workflow` goes red.
- [ ] The canary `craft-spec-fence-after-slice` bites on the new pointer sentence.
- [ ] The prose lane and the gate pass.
