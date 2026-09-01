# Remove anchor-test provenance labels

Blocked by: verify-done-claim-owners.md
Writes: internal/anchors/registry_data_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: none

## What to build

Remove LF identifiers from anchor-test comments while retaining any timeless
constraint explanation that helps a future reader understand the assertion.

## Acceptance

- [ ] New anchor-test comments contain no LF identifiers or implementation
      history.
- [ ] Every retained comment explains only a current-code constraint.
- [ ] Test behavior is unchanged.
