# Repair registry AST single source

Blocked by: none
Writes: `internal/conformance/subcommand_routing_test.go`, `internal/conformance/axi_query_registry_test.go`

## What to build

Collapse the duplicated `commandRegistry` AST locator and entry walk into the existing conformance seam. Keep routing and AXI membership as distinct judgments over the same parsed registry facts; do not weaken either check or broaden the approved AXI set.

## Acceptance

- [ ] [RA1] (covers QD3) one parser locates exactly one named `commandRegistry` composite literal, enumerates its entries in producer order, and exposes the keyed fields consumed by both subcommand routing and AXI membership checks.
- [ ] [RA2] (covers QD3) routing retains its dispatch parity, exact-registry identity, duplicate-name, malformed-literal, const-name, and decoy-registry refusal postures with their attributable diagnostics.
- [ ] [RA3] (covers QD3) AXI membership retains exact disposition cardinality, nested-member parsing, exemption reasons, and both-direction member mutations; removing an approved member and admitting an operational member are independently red.
- [ ] [RA4] (covers local) focused `go test ./internal/conformance -run 'Test(SubcommandRouting|AXIQueryRegistry)' -count=1` is green after the consolidation, and an observed mutation of the shared registry locator makes both consumers' targeted fixtures red.

