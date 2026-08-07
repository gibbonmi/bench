# Scope direct conformance fixture bites

Blocked by: expose-resolved-canary-fixture-checks.md, repair-fixture-bite-dependency-metadata.md
Ownership fence: `internal/conformance/fixture_bite_test.go`, `internal/conformance/data_handling_test.go`, `internal/conformance/decision_map_integrity_test.go`, `internal/conformance/example_agreement_test.go`
Integration surfaces: canary fixture inventory and resolved `Fixture.Check`→`internal/conformance/fixture_bite_test.go` + FB1/FB2; resolved-check producer→expose-resolved-canary-fixture-checks.md + FB1/FB2; executable registry rebind→existing `internal/conformance/registry/registry.go` + FB1/FB4; singular conformance scope and timing record→existing `internal/conformance/checks_test.go` + FB2/FB5; direct fixture-bite callers→the four fenced files + FB3; full-table controls→existing `internal/conformance/tier_test.go` and `internal/conformance/gate_entry_test.go` + FB6
Contracts: discovered fixture name, resolved non-empty `Fixture.Check`, and non-empty EXPECT cross canary fixture inventory→`internal/conformance/fixture_bite_test.go`, asserted by FB1/FB2/FB3/FB4 against the real inventory; resolved check name crosses expose-resolved-canary-fixture-checks.md→`internal/conformance/fixture_bite_test.go`, asserted by FB1/FB2 against the real canary producer and scoped runner; registry rebind crosses the executable registry→`internal/conformance/fixture_bite_test.go`, asserted by FB1 as consumer proof that the resolved value reaches the helper; timing check identity crosses the conformance runner→`internal/conformance/fixture_bite_test.go`, asserted by FB2/FB5 against the real timing record
Closure: FB1/ten-direct-families, FB1/registry-rebind, FB2/one-check, FB2/one-run, FB3/expect-diagnostic, FB4/missing-fixture, FB4/unbound-family, FB4/unknown-check, FB4/meta-check, FB4/wrong-tier, FB4/empty-expect, FB5/separate-subjects, FB5/timing-clear, FB6/full-entry, FB6/full-meta, FB6/ordered-selection

## What to build

Every direct Go conformance fixture-bite journey discovers its fixture and
consumes the resolved `Fixture.Check` from the canary inventory, runs exactly
that ordinary dev check, retains the fixture's existing EXPECT diagnostic, and
records only that check. The registry rebind remains a consumer proof that the
resolved value reaches the helper; canary alone owns CHECK-over-family
precedence. The common helper migrates the existing literal decision-map scope
too. Full-table completeness and public-entry controls remain unscoped. The
helper and all direct callers land together because either half alone leaves the
over-wide journey or an uncompilable call site.

## Acceptance

- [ ] [FB1] (covers CH1) the real direct journeys for the ten enumerated fixture families discover their fixture and consume its resolved `Fixture.Check`, including following a temporary registry-source rebind mutation.
- [ ] [FB2] (covers CH2) each mutated fixture tree invokes `RunConformance` once with the fixture's resolved ordinary dev check and records exactly that one check.
- [ ] [FB3] (covers CH3) every migrated fixture retains its independently authored non-empty EXPECT diagnostic under the registered check.
- [ ] [FB4] (covers CH4) missing fixture, unbound family, unknown check, meta check, wrong-tier check, and empty EXPECT each refuse before any conformance check executes.
- [ ] [FB5] (covers CH5) multiple fixtures sharing one check remain separate mutated-tree runs and each run clears then records its own timing identity.
- [ ] [FB6] (covers CH6) public full-tier entry, registry/meta integrity, ordered selection, timing ordering, and absent/unbound family controls remain unscoped and green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FB1/ten-direct-families | leave one current direct fixture caller on an empty or literal scope | the direct-fixture focused test set | apply the omission, run the nine named fixture-bite top-level tests, require the omitted caller's timing identity to differ from its registered check |
| FB1/registry-rebind | rebind `docs-currency-token-diet` from `docs-currency-workflow` to another valid ordinary dev check in an overlay | the coordinator-owned registry overlay probe | edit only the overlay registry row, run that family's real direct fixture journey, require the recorded check to follow the rebound resolved `Fixture.Check`, then discard the overlay |
| FB2/one-check | replace the helper's resolved `Fixture.Check` scope with an empty scope | the timing-identity assertion | apply the mutation, run one fixture bite, require multiple timing names instead of exactly the fixture's resolved check |
| FB2/one-run | invoke the scoped runner a second time for one fixture | the timing execution-count assertion | apply the mutation, run the helper control, require the call/record count to exceed one even though the writer clears between runs |
| FB3/expect-diagnostic | route one real fixture to a different valid ordinary check | the fixture's existing EXPECT assertion | apply the scope swap, run that fixture's direct journey, require its promised diagnostic to disappear |
| FB4/missing-fixture | remove the fixture from the discovered inventory presented to the resolver | the resolver refusal table | apply the omission, run the missing-fixture case, require refusal before the runner recorder fires |
| FB4/unbound-family | make the injected family lookup return unbound | the resolver refusal table | inject the result, run the unbound-family case, require refusal and zero runner calls |
| FB4/unknown-check | return a check name absent from `registry.Checks` | the resolver refusal table | inject the result, run the unknown-check case, require refusal and zero runner calls |
| FB4/meta-check | return `conformance-meta` for a direct fixture | the resolver refusal table | inject the result, run the meta-check case, require refusal and zero runner calls |
| FB4/wrong-tier | return the ship-only `release-evidence-probe` check | the resolver refusal table | inject the result, run the wrong-tier case at dev, require refusal and zero runner calls |
| FB4/empty-expect | present an existing fixture with an empty EXPECT value | the resolver refusal table | inject the empty value, run the case, require refusal before substring matching or runner execution |
| FB5/separate-subjects | reuse the first materialized root for a second same-check fixture | the same-check fixture control | apply the reuse, run the pair, require the second fixture's independent EXPECT/subject assertion to fail |
| FB5/timing-clear | retain a stale timing line before the scoped run instead of clearing it | the timing-clear control | seed a different check line, run one scoped fixture, require only the current resolved check afterward |
| FB6/full-entry | scope `TestRootConformance` to one ordinary check in an overlay | the public-entry integration control | apply the scope, run the full-entry control, require the expected full registry/timing inventory to shrink and fail |
| FB6/full-meta | scope a registry/family integrity control through the fixture helper | the meta-integrity control | apply the scope, remove one family binding in an overlay, require the missing integrity diagnostic to expose the mutation |
| FB6/ordered-selection | replace the ordered-set control with the singular helper | the ordered-selection control | apply the substitution, run the control, require the meta-plus-selected registry order assertion to fail |
