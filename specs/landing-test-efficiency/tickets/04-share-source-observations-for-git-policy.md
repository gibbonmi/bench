# 04 Share source observations for Git ownership checks

Blocked by: 03-consolidate-conformance-executions.md
Writes: internal/conformance/sourcefiles/ (new), internal/conformance/check_bindings_test.go, internal/conformance/checks_test.go, internal/conformance/git_plumbing_owner_test.go, internal/conformance/tier_test.go
Covers: LTE25, LTE26, LTE27, LTE32

## What to build

Review chunk: LTE-C4.

Introduce the run-owned source snapshot and use it for the Git plumbing visitor. Keep named executable bindings and add a typed snapshot function form.
The dispatcher owns snapshots by subject. Direct binding calls create fresh snapshots. Root and kit subjects remain independently routed.

The shared interface supplies deterministic directory observations, canonical source bytes, parse results, file positions, and observation errors.
Visitors select their own domains and error postures. Snapshot lifetime ends with the invocation, and callers treat returned observations as immutable.
Provide counted observation seams for later composed verification. No global cache or second production policy table is permitted.

This small foundation chunk closes before the architecture, bounds, and skip migrations start. Its Git visitor must deliver a real reduction independently.
Keep the independent flag expectation and binding identity, tier, and subject oracles. Preserve fresh-snapshot adapters for direct callers where useful.

## Acceptance

- [ ] Repeated file requests within one snapshot cause one source read and one parse.
- [ ] New invocations observe mutation and restoration at the same path.
- [ ] The Git visitor's diagnostics match the prior implementation across the spec's domain and error cases.
- [ ] Root Go files remain checked, while tests and owner packages retain their exemptions.
- [ ] Function, tier, and subject substitutions still fail metadata checking.
- [ ] The named flag-omission mutation still fails the independent expectation.
- [ ] Go import loading succeeds for the new helper package and its conformance consumer.

## Verification

Compare old and new diagnostic slices on the enumerated synthetic family before removing the old walker. Run sourcefiles tests, all Git plumbing tests, and TestConformanceMetaBites. Record the omitted-flag red and exact restoration. Do not retain a second policy implementation.

