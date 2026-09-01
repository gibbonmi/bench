# Preserve repeated retrospective capture

Blocked by: add-retrospective-writer.md
Writes: internal/roadmap/retro.go, internal/roadmap/retro_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: none

## What to build

Refuse a retrospective write when its slug already names an existing capture.
The refusal must leave the earlier body byte-for-byte unchanged and must not
weaken the existing distinct-slug accumulation behavior.

## Acceptance

- [ ] A second write to the same slug refuses before mutation.
- [ ] The first body remains byte-for-byte unchanged after the refusal.
- [ ] Writes to distinct slugs still preserve both eligible retrospectives.
