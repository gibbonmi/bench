# Register the conformance trust kernel

Blocked by: seal-outer-conformance-selection.md
Ownership fence: `internal/conformance`, `internal/conformance/registry`, `tests/canary`
Assumptions: semantic policy checks stay ordinary and scopeable; suite-only live-tree assertions are the existing hidden-enforcement defect; claims re-derived from the tree at pickup

## What to build

Every live-tree conformance assertion is registered exactly once, and the registry names
a minimal always-run meta set that proves selection machinery without absorbing semantic
policy checks.

## Acceptance

- [ ] [CM1] Every executable live-tree assertion is registered exactly once as meta or ordinary, including the current suite-only checks.
- [ ] [CM2] Meta covers registry/function/tier agreement, declaration completeness, profile bindings, canary ownership, and partition completeness.
- [ ] [CM3] Semantic docs, routing, package, bounds, decision-map, and example-agreement checks remain ordinary rather than always-on meta.
- [ ] [CM4] Adding either a prefix-named or an ordinary-named live-tree test without registration turns meta red.
- [ ] [CM5] Swapping two registered function bindings without changing their names turns meta red.
- [ ] [CM6] Moving a dev check to ship-only without changing its executable binding turns meta red.
- [ ] [CM7] Every migrated live-tree check declares whether it grades `root` or `kitRoot`; component-scope binding grades `kitRoot`, and checks over optional fixture surfaces guard absence before grading.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CM1 | remove one registered live-tree check while leaving its executable assertion | registry bijection test | inspect the executable family through the registry contract, expect the unregistered-enforcement diagnostic |
| CM2 | delete one declaration, profile binding, canary ownership validation, or partition validation from meta | conformance-meta mutation suite | apply each omission independently, expect its targeted trust-kernel diagnostic |
| CM3 | mark `docs-currency-workflow` as meta | meta membership contract | compare the enumerated semantic family to meta, expect the ordinary-check-classification failure |
| CM4 | add `TestRootConformanceHiddenPolicy` and, independently, `TestHiddenLiveTreePolicy` without rows | hidden-enforcement contract | make each test grade the live tree, compile the executable live-tree inventory against the registry, expect each test name in its diagnostic |
| CM5 | exchange the functions bound to two unchanged registry names | registry/function agreement test | swap the bindings, expect both names in the binding-mismatch diagnostic |
| CM6 | change one dev registry row to ship while leaving its function bound | tier membership test | resolve the dev executable set, expect the missing dev-check diagnostic |
| CM7 | bind component-scope grading to `root`, or remove an optional-surface guard | subject-binding contract | drive a minimal graded fixture against the real kit root, expect no missing-profile noise; then plant the optional surface and require its targeted diagnostic |
