# Project ready, stale, and path rows with the map argument

Blocked by: refuse-the-inline-shape-and-render-both-templates.md
Writes: internal/maps/maps.go, internal/maps/freshness.go (new), internal/maps/freshness_test.go (new), internal/maps/maps_command_test.go, internal/maps/testdata/, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: DS22, DS23, DS24, DS25, DS26, DS27, DS28, DS29, DS30, DS31, DS32

## What to build

`bench maps [<map>] [--count|--template|--ticket-template]`. The row shape becomes `maps[N]{map,title,type,state,blockers,path}`. A ticket row's path is its ticket file. A fog, invalid, or ready row's path is the map index file. Each ready active map adds `[name, title, map, ready, "", path]` and the help line `/bench-write-spec <path>`. A compiled map adds no row.

For each Sources Path of an active map, read `git log -1 --format=%ct -- <asset>` and the same for the map index file. Use the git output helper with the repository root as the working directory. When the asset's commit is newer, add `[name, locator, source, stale, "", asset path]`. When either has no commit, add nothing. A URL source adds nothing. Stale rows change no exit code and no count.

The map operand filters the rows by map name after the scan. A well-formed name with no active map exits 1 with `maps: no active map named "<name>"`. An operand with `/`, a `.md` suffix, or a byte below 0x20 is a grammar error. Regenerate the golden files under `testdata/` and any help snapshot the command-registry test pins.

## Acceptance

- [ ] [P1] The regenerated frontier-plus-invalid golden output differs from the old one only by the `path` column.
- [ ] [P2] A ready active map projects its row and its write-spec help line, and a ready compiled map projects nothing.
- [ ] [P3] `bench maps alpha` prints only alpha's rows and help lines, and `bench maps gamma` exits 1 with the refusal.
- [ ] [P4] `bench maps decisions/alpha.md`, `bench maps alpha.md`, and an operand with `\x1b` each exit 2.
- [ ] [P5] An asset committed after its map index prints the stale row, and the reverse order prints none.
- [ ] [P6] An uncommitted asset or index prints no stale row, the exit code stays 0, and `--count` is unchanged.
