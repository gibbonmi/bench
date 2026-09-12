# Advertise landing in root help

Blocked by: none
Writes: cmd/bench, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, tests/canary/package-core-guard/unrouted-subcommand, CHANGELOG.md
Covers: none

## What to build

Add landing and landing-resume rows to the root help inventory.
Derive their grammar from the existing usage constants.
Include land in the worktree help summary.
Keep the command registry as the root inventory owner.
Move existing worktree-specific helpers if main.go needs headroom.
Add a Fixed entry under a dedicated Worktree landing help changelog heading.

The coordinator reconciles the roadmap after this fix lands.

## Acceptance

- [ ] Root help displays the existing landing grammar.
- [ ] Root help displays the existing landing-resume grammar.
- [ ] The worktree summary names land.
- [ ] All root help spellings display the same inventory.
- [ ] Removing a landing row fails the independent inventory expectation.
