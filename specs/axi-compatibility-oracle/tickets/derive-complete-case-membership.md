# Derive complete compatibility case membership

Blocked by: authenticate-baseline-manifest.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`
Integration surfaces: root registry→`cmd/bench/main.go` exercised unchanged; wrapper grammar→`bin/bench.sh` exercised unchanged; nested grammars→existing production packages exercised unchanged; comparator→compare-four-observations.md
Contracts: member name and argv-class enum with exact not-applicable reason cross production registries/wrapper→`internal/axi/compatibility`, membership is every owner/class in `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md`, production order is preserved, and absence is an error, asserted by MC1
Closure: MC1/root-members, MC1/wrapper-members, MC1/nested-members, MC1/argv-classes, MC1/not-applicable

## What to build

case membership closes per production member over all roots, wrapper surfaces, nested grammars, and applicable required argv classes or exact not-applicable reasons.

## Acceptance

- [ ] [MC1] (covers CO3) case membership closes per production member over all roots, wrapper surfaces, nested grammars, and applicable required argv classes or exact not-applicable reasons.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MC1/root-members | omit one named root registry member | membership closure test | derive cases and require the missing root name |
| MC1/wrapper-members | omit wrapper-only repair or one version alias | membership closure test | derive wrapper cases and require the missing surface |
| MC1/nested-members | omit one named nested grammar operation | membership closure test | derive nested cases and require the missing operation |
| MC1/argv-classes | drop empty-required-value from one applicable member | class closure test | validate the member/class matrix and require the missing pair |
| MC1/not-applicable | accept an empty not-applicable reason | class closure test | validate the disposition and require a concrete reason |

