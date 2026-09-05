# Own the operand path match in intent

Blocked by: none
Writes: internal/intent/worktree_owner.go (new), internal/intent/worktree_owner_test.go (new), internal/intent/assignment.go, internal/worktree/path.go, internal/worktree/landed.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SR24, SR25, SR26, SR27, SR28, SR29, SR60

## What to build

The line is opus / medium. Export a canonical-path match from
`internal/intent` over a given slice of assignments, per decision (e). The
match answers every row whose canonical worktree equals the canonical path,
in ledger order, in any state. The exported match lands in the new file
`internal/intent/worktree_owner.go` with its test, per decision (o), so
`internal/intent/assignment.go` does not grow.

The documented lookup composes the exported match and keeps its active-only
answer. The lookup keeps its one-argument shape, so the three existing intent
tests stay as written.

Move only the path matches the spec's census names. The selector's path arm
in `internal/worktree/path.go` calls the exported match over the slice it
already holds, then applies its own ambiguity, state, and missing-tree
refusals. The missing-branch scan in `internal/worktree/landed.go` calls the
same match and keeps its active-state filter, per decision (p). A new unit
test beside the landed classifier tests pins that filter. The four
exact-string identity scans stay unmoved, per decision (f), because a
symlinked pool home changes their answer.

The exit proof for this ticket is the pre-existing suite, green with its test
logic unchanged. A mechanical rename is permitted. A needed assertion change
stops the ticket and reports.

## Acceptance

- [ ] The documented lookup matches through a symlink, ignores a retired assignment, and answers no owner outside the pool.
- [ ] The exported match answers a retired row and a symlinked spelling at one path, in ledger order.
- [ ] Every path-taking verb resolves a label, an id, and an 8-12 character prefix, and an ambiguous prefix names every colliding id.
- [ ] A `--from` that names a sibling in a non-active state refuses with the state component named.
- [ ] `list` names one clean landed row for an assignment whose tree is missing, and a target verb names the missing-tree reason.
- [ ] `bench consumers intent.Assignments` lists the fifteen production sites, and the two moved sites call the exported match.
- [ ] The missing-branch scan answers no assignment for a lone retired row, and answers the active row when a retired row shares its path.
- [ ] Self-probe: omit the state filter from the documented lookup, and report the observed red.
