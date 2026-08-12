# FT153 branch-native canary contract research

Assessed 2026-08-12 against `main` at `62617657`. This asset records the
current factual premises for the FT153 decision map; it does not choose the
reviewer-owned contract.

## Current owners

- The ordinary gate has one Go test phase and no nested Go or canary driver
  (`projects/benchkit.md:187-216`; `internal/gate/phases.go:83-90`).
- `bench canary` discovers fixtures, selects one binding for each, and invokes
  `Dispatch` with a callback that always returns an empty diagnostic
  (`internal/canary/inventory.go:90-140`). It materializes no fixture and
  invokes no check implementation.
- The profile assigns planted-red and restored-green proof to an ordinary
  mutation test. `runFixtureBite` is the implementation that reads `EXPECT`,
  invokes the registered conformance owner on the materialized mutation,
  restores the subject, and requires the diagnostic to disappear
  (`projects/benchkit.md:211-216`;
  `internal/conformance/fixture_bite_test.go:545-567`).
- The selected-executable system journey proves CLI routing and repeats the
  inventory-derived count; it does not observe a fixture owner
  (`internal/systemtest/owner_test.go:361-374`).

## Executed probes

In an isolated assignment worktree, replacing
`tests/canary/load-validity-metadata/invalid-json/EXPECT` with
`FT153_SENTINEL_OWNER_MUST_BITE` produced:

```text
go run -buildvcs=false ./cmd/bench canary <worktree>
canary ok (184 fixture owners dispatched)
exit 0

go test -count=1 ./internal/conformance \
  -run TestLoadValidityMetadataFixturesBite/invalid-json
invalid-json did not bite through owner load-validity-metadata; want
"FT153_SENTINEL_OWNER_MUST_BITE"
exit 1
```

This confirms the documented inventory/test split. It is not a kit-gate
false-green.

Replacing
`tests/canary/package-core-guard/bounds-duplicate-owner/EXPECT` with
`FT153_UNGRADED_EXPECT_SENTINEL` then left both the source-built public command
and `go test -count=1 ./internal/conformance` green. Both probe mutations were
restored byte-for-byte; the assignment worktree was clean afterward.

## Incomplete proof inventory

- The tree carries 184 `EXPECT` files.
- `package-core-guard` carries 33; its direct bite list names five and says the
  removed canary sweep proves the rest
  (`internal/conformance/fixture_bite_test.go:222-239`). That leaves 28 retained
  fixture expectations outside the direct list.
- The two `compliance-hardening` fixtures have a synthetic bite test that
  hard-codes diagnostic substrings rather than reading their `EXPECT` files
  (`internal/conformance/compliance_checks_test.go:55-71`).
- `injected-ports/unregistered-port` is described as the permanent red, but its
  recorded bite proof builds a separate synthetic tree and does not read the
  fixture (`internal/conformance/injected_ports_test.go:305-349`).

The reader census therefore leaves 31 retained `EXPECT` files without a direct
comparison to the diagnostic produced from that fixture. Only the
`bounds-duplicate-owner` member was mutation-probed in this session; the other
30 are source-derived members of the same class.

## Linked-repo boundary

`bench init` seeds a fixture whose overlay plants `DO-NOT-SHIP` and whose
`EXPECT` names the example-check diagnostic. The generated gate runs the
example check against the ordinary repo root, then calls only the inventory
command and says that command proves the seed check bites
(`internal/adopt/init.go:51-65,103-136`). No current path materializes the seed
overlay or compares its `EXPECT` with the example-check diagnostic.

## Current architectural constraint

The branch-native profile prohibits a gate, wrapper, `go test`, or `go run`
constructor in a fixture owner (`projects/benchkit.md:385-387`). ADRs 0003 and
0009 still describe the removed nested sweep, while ADR 0001 still promises a
planted-reason tripwire without distinguishing the kit's ordinary tests from a
linked repo. Those documents are not evidence that the removed mechanism still
exists; the resulting contract needs a reviewer decision and current-state
documentation.
