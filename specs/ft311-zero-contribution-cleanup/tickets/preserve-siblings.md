# Preserve siblings with no contribution

Blocked by: none
Writes: internal/worktree, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, CHANGELOG.md
Covers: none

## What to build

Require contribution evidence before automatic landing cleanup selects an assignment.
Use the assignment's recorded Start and its resolved branch tip.
Retain an assignment when those commits are equal.
Retain an assignment when the evidence is absent, unresolved, or inconsistent.

Apply this condition through the existing scoped selector and its requalification paths.
Keep explicit cleanup unchanged.
Keep the existing ancestry, lease, preservation, and identity conditions.
Add a Fixed entry under a dedicated Empty sibling cleanup changelog heading.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] Automatic cleanup retains a sibling whose branch has no commit since Start.
- [ ] A sibling created at the published tip survives the landing tail.
- [ ] The same sibling survives a resumed landing.
- [ ] Unknown contribution evidence cannot authorize removal.
- [ ] A folded sibling with contribution remains eligible for cleanup.
- [ ] Explicit cleanup retains its current selection behavior.
