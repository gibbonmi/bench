# Delete the contract package and its marker-wait check

Blocked by: none
Writes: internal/contract/marker_wait.go, internal/contract/marker_wait_test.go, internal/conformance/marker_wait_deadline_test.go, internal/conformance/wait_deadline_literal_test.go, internal/conformance/checks_test.go, internal/conformance/bounds_policy_test.go, internal/conformance/registry/checks.go, internal/conformance/registry/packages.go, internal/freshness/publication_topology_test.go, internal/commit/commit_test.go, projects/benchkit.md, tests/canary/package-core-guard/bounds-dot-import-wait/BASE, tests/canary/package-core-guard/bounds-dot-import-wait/EXPECT, tests/canary/package-core-guard/bounds-dot-import-wait/MUTATE.json, tests/canary/package-core-guard/bounds-raw-elapsed-wait/BASE, tests/canary/package-core-guard/bounds-raw-elapsed-wait/EXPECT, tests/canary/package-core-guard/bounds-raw-elapsed-wait/MUTATE.json, tests/canary/package-core-guard/bounds-raw-injected-wait/BASE, tests/canary/package-core-guard/bounds-raw-injected-wait/EXPECT, tests/canary/package-core-guard/bounds-raw-injected-wait/MUTATE.json
Covers: none

## What to build

The quality survey of 2026-09-29, card 03, found that the `internal/contract` package has no production caller. The package holds one two-leg marker wait and its unit test. The `marker-wait-deadlines` conformance check grades the call sites of that wait, and no call site exists. So the check grades a shape that is gone.

Delete the package and the check. Move the shared duration-literal scanner into the `wait-deadline-literals` check, which is its last reader. Remove the check from the registry, from the check table in the conformance tests, and from the input table in the kit profile. The `bench test --help` check list derives from the registry, so it loses one line.

Remove the two exemptions that only skip the deleted package. These are the core test exclusion in the registry and the skip in the publication topology test. Correct the commit unit test comment that names a retired contract path.

Three `bounds-policy` canary fixtures use the deleted wait file as their subject. Move them to `internal/chargeevidence/store.go`, which is a live file. Each mutation adds a new unclassified wait of the same kind, so each fixture keeps its red and its restore proof.

Keep the guards that refuse a retired shape. These are the retired gate fragment in the gate entry check and the contract path rules in the build census.

The deleted wait was the only live code that showed two classified wait shapes read green under `bounds-policy`. These shapes are a classified duration argument to an injected duration function and a classified window in a current-time deadline compare. Add a unit test that plants each shape in a temporary root and expects no diagnostic. Pair each shape with its raw form and expect the named diagnostic.

## Acceptance

- [ ] No Go source names `WaitForTwoLegMarkers`, `MarkerWaitMiss`, `checkMarkerWaitDeadlines`, or `marker-wait-deadlines`.
- [ ] The `internal/contract` directory has no file.
- [ ] The registry, the conformance check table, and the kit profile input table have no `marker-wait-deadlines` row.
- [ ] The `wait-deadline-literals` check keeps its duration-literal scanner, and its bite test passes.
- [ ] The three `bounds-policy` fixtures bite on `internal/chargeevidence/store.go`, and each red goes away after restore.
- [ ] The gate entry check keeps its retired contract fragment, and the build census keeps its contract path rules.
- [ ] A `bounds-policy` unit test gives no diagnostic for a classified injected duration argument or a classified current-time deadline compare. Each raw form gives its diagnostic. A probe that rejects each classified form turns the test red.
- [ ] `go vet ./...`, the changed-package tests, and the touched conformance checks pass.
