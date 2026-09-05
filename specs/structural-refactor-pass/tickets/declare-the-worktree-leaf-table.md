# Declare the worktree leaf table

Blocked by: none
Writes: cmd/bench/worktree_leaves.go (new), cmd/bench/main.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/usage/worktree.go, tests/canary/package-core-guard/unrouted-subcommand, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: SR1, SR2, SR3, SR4, SR5, SR6, SR7, SR8, SR9, SR10, SR11, SR12, SR13, SR14

## What to build

The line is opus / low. Give `bench worktree` one leaf table. Each row names
the leaf, its grammar constants, its root need, and its handler. The root
need is required, boundary, or none. The dispatcher reads the root need once
per call and resolves the repository root at that one site.

A required root that fails prints `toon.NotInRepo()` on stderr at exit 1. A
boundary root passes the empty string to the handler. A leaf with no root
need receives no root. The nine root-required leaves are `path`, `exec`,
`show`, `build`, `create`, `release`, `reauthorize`, `merge`, and `land`.
`land` keeps its invoked-executable field and never consults it.

The bare, `help`, `--help`, `-h`, unknown-leaf, and unknown-flag answers stay
generic to the family and keep their current bytes and exit codes. The `help`
match still needs exactly one argument. The worktree registry entry keeps its
literal AXI child call, its `Name` literal, and its three hand-typed help
rows. Two conformance parsers read those literals. The usage text
stays in `internal/usage/worktree.go`, and the leaf table references the
grammar constants there. The table lands in the new file
`cmd/bench/worktree_leaves.go`, so the dispatch file does not grow.

The exit proof for this ticket is the pre-existing suite, green with its test
logic unchanged. A mechanical rename is permitted. A needed assertion change
stops the ticket and reports.

## Acceptance

- [ ] `bench worktree <leaf> --help` prints `usage: ` plus that leaf's grammar at exit 0 for every grammar-bearing leaf.
- [ ] Each of the nine root-required leaves prints `toon.NotInRepo()` on stderr at exit 1 in a non-repository directory.
- [ ] `bench consumers git.Root` lists one call inside the worktree dispatch.
- [ ] Bare `bench worktree` prints `usage.WorktreeUsage()` on stdout at exit 2 and acquires no assignment.
- [ ] `bench worktree --help` exits 0, names every kept grammar and the clean grammar, and ends with the exec-gate trailer.
- [ ] `bench worktree help extra` prints `toon.Usage("bench worktree", "help")` on stderr at exit 2.
- [ ] `bench worktree unknown` and `bench worktree --unknown` print the family usage line on stderr at exit 2 and leave the ledger unchanged.
- [ ] `bench worktree create --unknown` refuses through the create grammar, and `bench worktree clean --help` prints exactly the clean grammar.
- [ ] `bench worktree clean --help` and `bench worktree list` answer identical bytes from the root, from a subdirectory, and outside a repository.
- [ ] `bench worktree shell` with `SHELL=true` creates and releases its assignment at exit 0.
- [ ] `bench worktree land` refuses with repository proofs and never names the invoked executable.
- [ ] `bench help` prints the byte-identical public inventory, and the wrapper and the binary agree.
- [ ] `bench test --check axi-query-registry` and `bench test --check subcommand-routing` stay green.
- [ ] Every argv row in the spec's differential family answers identical stdout, stderr, and exit code at the base and at the tip.
- [ ] Self-probe: omit the root need from the `path` row, and report the observed red.
