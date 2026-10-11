# State what Bench does in the help title

Blocked by: none
Writes: cmd/bench/command_registry.go, cmd/bench/help_inventory_test.go, cmd/bench/main_test.go, internal/conformance/help_inventory_single_source_test.go, CHANGELOG.md
Covers: none

## What to build

The first line of `bench help` read "bench — Pocock pipeline meets Kun Chen substrate, gated by your invariants." That line names the projects that Bench borrows from. It does not tell a reader what Bench does or where to start. On 2026-10-11 the reviewer directed that the line become meaningful and helpful.

The new title states what Bench does: agent work in isolated worktrees, graded by one gate, and landed only on green. The title also names `bench status` as the command that gives the next action. The title does not name `bench <verb> --help`, because a verb can take `--help` as an operand.

The production constant stays the one source of the title. The binary routing test derives its prefix from that constant. The conformance check keeps its own copy of the title, because it cannot import the `main` package. Its bite test reads that copy. The completeness test keeps its independent expectation, because that expectation is the omission oracle.

## Acceptance

- [ ] `bench help`, `bench --help`, and `bench -h` print the new title as the first line, and the inventory rows follow it unchanged.
- [ ] The help inventory single-source check finds the new title exactly once on the production surfaces.
- [ ] `CHANGELOG.md` has one entry for the changed title.
