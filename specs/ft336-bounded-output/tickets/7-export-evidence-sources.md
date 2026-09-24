# 7. Export the verified evidence sources to a directory

Blocked by: 6-summarize-evidence-default.md
Writes: internal/preflight/evidencecmd/, internal/chargeevidence/export.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go
Covers: BO45, BO46, BO47, BO48, BO49, BO50

## What to build

Add `bench preflight evidence <id> --to <dir>` to the evidence operation registry, so the grammar, the preflight help, and the root help show it. A relative directory resolves against the working directory. The directory must be absent or empty, and it must not be a symlink.

The export verifies each page digest and each source digest before it writes that source. It writes `source-<ordinal>` for each source and one `index.toon` table with the ordinal, the source id, the bytes, and the file name. On any failure, it removes each file it created, and it removes the directory if it created it. It prints one `exported{sources=<n>,bytes=<n>,dir=<absolute path>}` line.

## Acceptance

- [ ] `--to <dir>` over 3 sources writes `source-1` to `source-3` with the source bytes and one `index.toon`.
- [ ] `--to` at a directory that holds one file exits 1 and writes nothing.
- [ ] `--to` at a symlink exits 1 and writes nothing.
- [ ] `--to` over a pack with one corrupt page exits 1 and leaves no directory that it created.
- [ ] A source id `../escape` exports to `source-1` inside the directory.
- [ ] A write fault on the second source removes `source-1` and the created directory.
