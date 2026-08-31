# Bind citations to phase package scope

Blocked by: none
Writes: CHANGELOG.md, internal/gate/tag_census.go, internal/gate/tag_census_test.go, internal/coverage/citations.go, internal/coverage/citations_test.go, internal/coverage/citation_execution.go (new), internal/coverage/citation_execution_test.go (new), internal/bounds/bounds.go, internal/conformance/bounds_policy_test.go, tests/canary/coverage-map-validation/unexecuted-tag-citation/
Covers: PS1, PS2, PS3, PS4, PS5, PS6, PS7, PS8, PS9, PS10, PS11, PS12

## What to build

Make citation execution use one complete Go test phase entry. Preserve the
phase's package operands beside its tags, and use the real Go package loader to
resolve the selected packages.

The coverage validator accepts a cited file only when one entry selects its
package and accepts its build constraints. Extend the existing canary so the
new ignored-package diagnostic reaches the gate.

## Acceptance

- [ ] Different phases cannot donate package scope and build tags to one citation. (covers PS1)
- [ ] Equal tag sets retain each phase's package operands. (covers PS2)
- [ ] A manifest test phase contributes its package operands. (covers PS3)
- [ ] A recursive package pattern excludes a cited `testdata` package. (covers PS4)
- [ ] A recursive package pattern excludes a cited underscore-prefixed package. (covers PS5)
- [ ] An explicit package operand can select cited evidence. (covers PS6)
- [ ] An absent package list selects the effective phase directory. (covers PS7)
- [ ] A scope violation names its coverage row and cited file. (covers PS8)
- [ ] A package expansion error rejects the citation. (covers PS9)
- [ ] Package selection does not weaken build-constraint matching. (covers PS10)
- [ ] A root with no Go test phase keeps the execution check inapplicable. (covers PS11)
- [ ] The extended canary reports the ignored-package citation. (covers PS12)
