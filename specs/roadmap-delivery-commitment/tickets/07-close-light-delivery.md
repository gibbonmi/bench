# Close verified light-path and rowless delivery

Blocked by: 06-close-verified-spec-delivery.md
Writes: internal/commitment (new), internal/intent, internal/worktree, internal/landing, internal/roadmap, internal/spec, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DC67, DC68

## What to build

Apply the reviewed exact closure contract to tickets-only deliverables and outcomes with no roadmap row.
Ticket 06 supplies the complete transform and recovery path. The light-path acceptance and gate evidence replace the full-spec completion record only where the existing light path permits it.
A `Roadmap:` label alone is never evidence of full delivery.

An approved complete row closes in the light-path publication. A rowless outcome records only its bound obligation and produces no invented board entry.
Keep partial completion and unrelated obligations under the same owner rule as spec delivery. Do not copy a second closure algorithm.

Read the spec metadata and tickets-only retirement readers, then consume the reviewed ticket 06 closure seam. The checkpoint is a second real delivery route, not a test-only follow-up.

## Acceptance

- [ ] A tickets-only delivery closes its approved complete roadmap obligation in the published commit (DC67).
- [ ] A rowless delivery completes only its bound obligation without creating or deleting roadmap rows (DC68).
- [ ] Both routes retain ticket 06 refusal and resume behavior through the shared closure owner.

## Checkpoint verification

Run `bench test --package ./internal/worktree`, `bench test --package ./internal/landing`, and `bench test --package ./internal/spec`. Use actual tickets-only and rowless candidate fixtures, with successor milestone verification still absent.
