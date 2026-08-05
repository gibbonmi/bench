# Expose the anchor query

Blocked by: migrate-anchor-table.md
Ownership fence: `cmd/bench/main.go`, `cmd/bench/main_test.go`, `internal/contract/axi/axi_anchors_test.go`, `internal/contract/axi/axi_test.go`, `internal/contract/axi/axi_grammar_test.go`, `internal/conformance/subcommand_routing_test.go`, `internal/conformance/docs_workflow_helpers_test.go`, `bin/bench.sh`, `.bench/BENCH.md`, `projects/benchkit.md`
Contracts: compiled registry rows cross `internal/anchors`→`cmd/bench/main.go`→the AXI stdout table, asserted by AQ1 against the real registry; the `anchors` command name and argv grammar cross `cmd/bench/main.go`→`bin/bench.sh` and `internal/conformance/subcommand_routing_test.go`, asserted by AQ3 and AQ4 through every shipped route

## What to build

Add `bench anchors <path>` over the compiled registry with kind, section, and needle columns, repository-relative path resolution, definitive empty output, structured usage errors, shell/help routing, and the canonical query-surface advertisements. This tracer legitimately spans more than two directories because one public verb must carry its grammar, launcher, contract, routing registry, and user-facing inventory together to stay green.

## Acceptance

- [ ] [AQ1] Querying pinned paths lists every matching registry row with kind, section, and needle, including a section-scoped anchor tied to real file content.
- [ ] [AQ2] Querying a path with no anchors prints the definitive empty table and exits 0.
- [ ] [AQ3] A missing argument, extra argument, or unknown flag prints a structured stdout usage error and exits 2.
- [ ] [AQ4] Invoking from a repository subdirectory resolves the same repository-relative anchor path and every shipped routing surface reaches the same implementation.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| AQ1 | omit the section field or one row for a pinned path | the AXI anchors contract | build the worktree binary, query both pinned paths, and expect the exact row-set comparison to fail |
| AQ2 | render no bytes for an unpinned path | the AXI anchors contract | query the unpinned literal and expect the definitive `anchors[0]` header assertion to fail |
| AQ3 | accept an extra positional argument | the AXI grammar contract | invoke with two paths and expect exit 2 plus structured stdout usage |
| AQ4 | resolve the query path relative to process cwd | the deep-cwd AXI contract and routing checks | invoke from a nested directory and through the shell launcher, expect the root-relative row set and no unrouted-command diagnostic |
