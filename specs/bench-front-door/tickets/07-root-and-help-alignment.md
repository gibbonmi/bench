# Align the wrapper's bare root and the binary's no-arg/help

Blocked by: 06-front-door-adapters.md
Writes: bin/bench.sh, cmd/bench, internal/conformance/subcommand_routing_test.go

## What to build

Bare `bin/bench.sh` runs `status --route` (route_porcelain) and returns its exit code;
the binary's no-arg does the same. Move the inventory heredoc into a Go `help` verb
(registered in `commandRegistry`; `subcommandRouting` gains an exempt row: takes no
arguments); the wrapper's `help|--help|-h` arm routes to it. Reviewer decision recorded
in the spec: if the alternative is chosen, keep the heredoc and make the binary's `help`
print a one-line pointer instead — the acceptance below then reads "pointer" for R38's
help half.

Covers: R37, R38, R39, R40, R48

## Acceptance

- [ ] Bare wrapper prints the `next[1]` table; `bin/bench.sh help` prints the inventory starting `bench — Pocock pipeline`.
- [ ] `dist/bench` (rebuilt) with no args prints the same table; `dist/bench help` is byte-identical to `bin/bench.sh help`.
- [ ] `rg "Pocock pipeline"` finds one production source; routing conformance green with the `help` row.
- [ ] Gate green.
