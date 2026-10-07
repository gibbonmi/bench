# Print the commitment verb routes from the registry

Blocked by: 11-route-the-commitment-policy-refusals-through-faces.md
Writes: internal/refusalroute/registry.go (new), internal/commitment/commitcmd/command.go, internal/commitment/commitcmd/admission.go, internal/commitment/commitcmd/refusal_route_test.go (new), internal/commitment/verification_test.go, cmd/bench/commitment_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RR45, RR46, RR47, RR48

## What to build

The commitment verb keeps its `next[1]{command}` table.
The cell holds the route of the face that refused.
A refusal that only the reviewer clears prints a route that starts with `reviewer: `.
A refusal that the agent clears prints an agent route.

The `commitment-decision` face covers each start, block, plan, and approve refusal whose clear changes the active commitment.
Declare the `commitment-plan-input` face and the `commitment-handback` face in the registry.
A refusal with no typed face takes `commitment-handback`.
The registry owns the commitment `next` table, so the verb composes the table label from the registry export.
The verify success table also composes its label from the registry export.

Add `TestCommitmentFacesFollowTheirRoutes`.
It walks the commitment faces of the registry and drives one producing fixture for each face.
Each fixture follows its printed route and reruns the verb.

## Acceptance

- [ ] A `bench commitment start` refusal for an outcome outside the active milestone prints a `next` cell that starts with `reviewer: `.
- [ ] A `bench commitment start` refusal without an owned assignment prints a `next` cell that contains `bench worktree create --request`.
- [ ] A `bench commitment verify` refusal keeps its `bench commitment verify --milestone` cell.
- [ ] Each commitment face has exactly one producing fixture that follows its route out of the face.
