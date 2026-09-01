# Contain retrospective writes across path races

Blocked by: harden-retrospective-writer.md
Writes: internal/roadmap/retro.go, internal/roadmap/retro_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: none

## What to build

Keep retrospective creation rooted inside the primary checkout when a path
component changes after validation. Preserve the static live-link and
dangling-link refusals and the create-only destination contract.

## Acceptance

- [ ] Retrospective creation cannot escape the opened primary root when a
      destination component changes concurrently.
- [ ] Existing live and dangling outside symlinks remain refused.
- [ ] Ordinary create-only capture and same-slug preservation remain green.
