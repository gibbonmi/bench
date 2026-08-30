# Name the gate form in both inventories

Blocked by: none
Writes: internal/usage/worktree.go, cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go

## What to build

A cold agent reads either inventory and finds one gate form for a worktree. The
worktree usage trailer becomes `bench worktree exec <target> -- bench gate`, so
`bench worktree --help` names no raw gate path. The `bench help` row
`bash bin/bench.sh gate --fresh` moves out of the gate family. It becomes a
worktree row that names the same exec form. The gate family keeps its two other
rows, and `.bench/gate.sh` keeps its own wrapper-less refusal.

This ticket edits the exact `bench help` fixture in `cmd/bench/main_test.go`: it
adds the gate-form row and removes the `bash bin/bench.sh gate --fresh` row. The
build ticket and the create ticket edit the same fixture, `internal/usage/worktree.go`,
and `cmd/bench/command_registry_test.go`. That is a write conflict, not a
capability dependency, so the coordinator serializes the three worktree tickets
in the order build, create, gate form. In that order this ticket is the last one
to write `internal/usage` and `cmd/bench`, so it carries both packages' gate
invariants.

## Acceptance

- [ ] WF44: `bench help` prints the `exec <target> -- bench gate` row and no `bash bin/bench.sh` row, byte-equal to the fixture.
- [ ] WF14: `bench worktree --help` ends with `bench worktree exec <target> -- bench gate` and holds no `bin/bench.sh`.
- [ ] WF15: `.bench/gate.sh` still refuses a wrapper-less entry, and it still names `bash bin/bench.sh gate`.
- [ ] The gate `test` phase stays green for the whole `internal/usage` and `cmd/bench` packages.
