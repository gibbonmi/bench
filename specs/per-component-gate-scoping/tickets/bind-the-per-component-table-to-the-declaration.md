# Bind the per-component table to the declaration

Blocked by: Derive the toolchain and contract input sets; Declare shellcheck's
and canary's input sets
Ownership fence: `internal/conformance/component_scope_binding_test.go`, `internal/gate/component_inputs.go`, `projects/benchkit.md`
Assumptions: `checkScopeBinding` in `scope_binding_test.go` binds the existing
reduced-scope table to `gate.ReducedScope()` and is read — not edited — by this
ticket; that table and `ReducedScope()` both stay, as the capture-surface floor
the per-component declarations layer on top of. Re-derive from the tree at
pickup.

## What to build

The profile renders the per-component declarations, and a conformance check
binds the rendering to the source so drift reds the gate rather than waiting for
review. One row per scoped component: what it declares, and whether that
declaration is derived or hand-written.

The prose also states the declaration-honesty width, with the same candor the
current construction prose carries: the stripped-worktree construction proves
capture-surface blindness only. For the per-component declarations, honesty
rests on mandatory derivation plus this binding, and a component that reads an
undeclared non-capture path skips wrongly. That residual is recorded, not
hidden. The `canary` row states the reviewer's 2026-08-01 narrowing — the binary
is excluded, two changes graded separately may land together with the canary
never run against the combined tree, and `bench gate --fresh` and the ship tier
are what re-prove the tripwire.

## Acceptance

- [ ] [PC6a] the profile renders one row per scoped component and the check compares each row against the declaration by exact set equality.
- [ ] [PC6b] a row mutated in the prose alone reds, and a declaration mutated alone reds, each naming the component.
- [ ] [PC6c] a component missing from the table and a table row naming no declared component are each named as their own diagnostic.
- [ ] [PS31] the profile states the declaration-honesty residual and the canary narrowing in prose.
- [ ] [PS32] the existing reduced-scope table and its binding still render and pass unchanged.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC6a | compare rows by subset instead of equality | `TestComponentScopeBindingBites` | render a profile whose row carries an extra token, run the check, expect a diagnostic |
| PC6b | read the expected sets from a literal restated in the check | `TestComponentScopeBindingBites` | mutate the prose alone and the declaration alone in two subtests, expect one diagnostic each naming the component |
| PC6c | treat an absent row as an empty set | `TestComponentScopeBindingBites` | delete one component's row and rename another, expect a missing diagnostic and an unknown diagnostic |
| PS31 | delete the residual paragraph | `TestProfileStatesTheHonestyResidual` | read the profile, assert the residual and the canary-narrowing sentences are present |
| PS32 | fold the reduced-scope table into the per-component table | existing `TestDeclaredAllowlistBindingBites` | run the existing binding cases unchanged |
