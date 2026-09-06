# Note the file-tool route on the path verb

Blocked by: 04-add-the-bench-probe-verb.md
Writes: internal/worktree/path.go, internal/worktree/path_identifier_test.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/help_inventory_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, tests/canary/package-core-guard/unrouted-subcommand
Covers: PB34, PB35

## What to build

Verify the premise first. Read `PathCommand` in internal/worktree/path.go. Read
`TestListPathActionRunsAsAdvertised` and `TestPathResolvesTheLabelAndTheIdAlike` in
internal/worktree/path_identifier_test.go, which read stderr only on failure. Read
the `bench worktree path` help row in `commandRegistry` and its golden line in
cmd/bench/help_inventory_test.go.

After the path line on stdout, print on stderr the one line
`note: the path serves the file tools; run a shell step through bench worktree exec <target> -- <command>`,
with `<target>` replaced by the operand as typed and `<command>` kept literal. Change
the help row description to `print one active owned worktree's absolute path for the
file tools` and the golden line with it. Keep stdout as one path line, so a `$(...)`
capture stays intact.

Write `TestPathNotesTheFileToolRouteOnStderr`: stdout equals the path line alone, and
stderr equals the note with the typed target.

Self-probe: print the note on stdout and show the stdout equality red.

## Acceptance

- [ ] `bench worktree path <target>` prints the path alone on stdout and the note on stderr with the typed target.
- [ ] `bench help` prints the path row with the new description.
- [ ] `go test ./internal/worktree -run 'Path' ./cmd/bench -run 'Help'` passes.
