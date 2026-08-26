# Add the bench harnesses verb

Blocked by: 01-add-the-harness-record-package.md
Writes: internal/harnesses/command.go (new), internal/harnesses/command_test.go (new), cmd/bench/main.go, cmd/bench/main_test.go, bin/bench.sh, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go, projects/benchkit.md, .agents/skills/bench-craft-cli/SKILL.md

## What to build

`bench harnesses` projects the record as TOON, so the matrix is one call
away.

Bare, the verb prints `schema: 1` first, then
`harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}`
with one row per record row. The `none` row appears like the rest, because
the degraded path must stay visible.

With one harness argument, the verb prints
`cells[12]{field,value,source,checked}` for that row. Each cell carries its
source and its ISO date, and an `unknown` cell leaves both empty. So
`bench harnesses codex` names the Codex hooks docs as the `delegation_guard`
source.

An unknown harness argument prints the usage line and exits 2. Two positional
arguments are also a usage error, because the verb takes at most one.

The verb joins five registries, and the gate reds if any one misses it:

- `approvedAXIQueries` in `internal/conformance/axi_query_registry_test.go`, with a nil child list.
- The AXI seam bullet in `projects/benchkit.md`.
- The disclosure table in `.agents/skills/bench-craft-cli/SKILL.md`, with a terminal `help[0]` disposition.
- The `harnesses)` case label in `bin/bench.sh`, routed as a porcelain command.
- The `subcommandRouting` registry in `internal/conformance/subcommand_routing_test.go`.

Add the command to the `commandRegistry` in `cmd/bench/main.go` with a public
inventory help row. This ticket writes the profile's AXI seam bullet only;
ticket 07 writes the profile's conformance table, so the two edits do not
collide.

## Acceptance

- [ ] `bench harnesses` prints `schema: 1` and then `harnesses[4]{harness,provider,phase_form,hooks,delegation_guard,headless,checked}` with four rows. (covers HC24)
- [ ] `bench harnesses codex` prints `cells[12]{field,value,source,checked}` with the `delegation_guard` source naming the Codex hooks docs. (covers HC25)
- [ ] `bench harnesses cursor` prints the usage line and exits 2. (covers HC26)
- [ ] `checkAXIQueryRegistry` over the live root reports no diagnostic. (covers HC27)
- [ ] `bin/bench.sh harnesses` prints the same output as the direct verb. (covers HC28)
- [ ] `bench harnesses codex claude` prints the usage line and exits 2.
- [ ] `checkSubcommandRouting` and `checkBenchShRoutes` over the live root report no diagnostic.
