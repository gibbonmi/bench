# Red the gate on a copied list instead of a derivation

Blocked by: Derive the toolchain and contract input sets; Declare shellcheck's
and canary's input sets
Ownership fence: new `internal/conformance/derivation_source_test.go`
Assumptions: the registry in `internal/gate/component_inputs.go` records how each
entry was sourced; `canary` is the only hand-declared entry. Re-derive both from
the tree at pickup rather than from this line.

## What to build

A table check proves the profile and the declaration agree. It cannot prove the
declaration is honest: a hand-copied path list satisfies the table and drifts
from what the component actually grades the moment a package is added. This
check closes that gap — every component with a derivable source must call its
named derivation, and a literal path list standing in for one reds the gate.

`build` must take its set from `freshness.BuildInputs`; `gofmt`, `vet`, `test`,
`race`, and `contract` from the module-wide `-test` closure derivation;
`shellcheck` from its own argv enumeration. `canary` is the sole permitted hand
declaration, and the check names it as such rather than exempting it silently —
a second hand-declared entry appearing later reds until someone decides it is
allowed.

Per `craft-gate`, this check ships with its recorded bite proof: the mutation
that should red it, run against a synthetic declaration, with the real one as
the fixed side.

## Acceptance

- [ ] PC7a — each of `build`, `gofmt`, `vet`, `test`, `race`, `contract`, and `shellcheck` is verified to resolve through its named derivation.
- [ ] PC7b — replacing any of those with a literal path list produces a diagnostic naming the component and its expected derivation.
- [ ] PC7c — a second hand-declared entry beyond `canary` produces a diagnostic.
- [ ] PS33 — a derivation that cannot resolve is reported as its own diagnostic rather than read as an empty, satisfied set.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC7a | check only that each entry is non-empty | `TestDerivationSourceCheckBites` | run the check against the real registry, expect no diagnostics |
| PC7b | compare against the recorded source label without exercising the derivation | `TestDerivationSourceCheckBites` | substitute a literal-list entry for `vet` in a synthetic registry, run the check, expect one diagnostic naming `vet` and its derivation |
| PC7c | allow any entry to declare itself hand-written | `TestOnlyCanaryMayBeHandDeclared` | mark `shellcheck` hand-declared in a synthetic registry, run the check, expect a diagnostic |
| PS33 | treat a derivation error as no diagnostic | `TestUnresolvableDerivationIsNamed` | point the check at a root with no `go.mod`, expect a diagnostic naming the failure |
