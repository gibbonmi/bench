# Make the retrospective race test bite

Blocked by: contain-retrospective-write-race.md
Writes: internal/roadmap/retro_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: none

## What to build

Make the deterministic retrospective race test replace an existing destination
component after the primary root opens. Prove the hook completed the
replacement before accepting the command refusal.

## Acceptance

- [ ] Test setup removes or renames the existing `capture` directory before
      installing the outside symlink.
- [ ] Setup failures fail the test instead of satisfying the expected refusal.
- [ ] The command refuses after the successful swap and writes no outside file.
