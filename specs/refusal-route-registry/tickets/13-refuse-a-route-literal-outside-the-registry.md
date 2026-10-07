# Refuse a route literal outside the registry

Blocked by: 12-print-the-commitment-verb-routes-from-the-registry.md
Writes: internal/conformance/refusal_route_bypass_test.go (new), CHANGELOG.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_retained_workflow.go, tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns, tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: RR49, RR50, RR65

## What to build

Add the static bypass check in `internal/conformance`.
It scans the string literals in the production Go files of the write-verb packages.
The scanned packages include `internal/commitment` and `internal/commitment/repository`.
It refuses a literal outside `internal/refusalroute` that holds `next=`, `next[`, `run bench `, or `; then `, or that equals `"next"`.

The check holds one reviewed allowlist of files whose routes serve a non-write verb.
The allowlist holds `internal/worktree/path.go`, `internal/worktree/build.go`, and `internal/worktree/tree_target.go`.
Each other entry needs a non-write caller.

The check bites.
Its bite test plants a `next=` literal and a `next[1]:` literal in a write-verb production file, and a `; run bench` tail in an `internal/commitment/repository` production file.
Each planted literal turns the check red.

This ticket comes last, because the check turns red on any route that an earlier ticket has not moved.
If the check finds a literal in a file outside this ticket's `Writes:` line, stop and report the file.

Write the build's `CHANGELOG.md` entry.
It states the shared registry, the agent and reviewer authority, and `bench recovery`.
It also states the red-source fold exit, the commit exit 3 reset route, and the checkpoint route of each cause.

## Acceptance

- [ ] The bypass check passes on the tree after every earlier ticket lands.
- [ ] A planted `next=` literal and a planted `next[1]:` literal in a write-verb production file each make the bypass check fail.
- [ ] A planted `; run bench` tail in an `internal/commitment/repository` production file makes the bypass check fail.
