# Grade each anchored guidance file below a directory

Blocked by: none
Writes: cmd/bench/anchors_command.go, cmd/bench/anchors_dir_test.go, internal/anchors/registry.go, cmd/bench/main.go, cmd/bench/help_inventory_test.go, CHANGELOG.md
Covers: none

## What to build

Today `bench anchors <path>` accepts one file. This ticket adds the directory form, which is one face of FT125.

When the operand is a directory, the verb grades each file that the anchor registry names below that directory. The registry supplies the file inventory, so the verb does not walk the directory. The verb prints one row for each file, sorted by path.

Each row gives the file, its anchor count, its diagnostic count, and its verdict. The verdict is `pass`, `fail`, or `refused`. The rows go in one TOON table, and the help block follows the table. The help block names the file form when a row does not pass. The registry size sets the maximum number of rows.

A directory that has no anchored file below it gives an empty table and exit 0. The empty table is the explicit answer. The file form keeps its current output exactly. The help line and the help inventory row spell the operand as `<file|dir>`.

## Acceptance

- [ ] `bench anchors .` in a fixture repository prints one row for each registered anchor file, sorted by path, and no other row.
- [ ] `bench anchors <dir>` for a nested directory prints only the registered files below that directory.
- [ ] A planted file whose anchors all hold reads `pass` with zero diagnostics. An absent file with a required anchor reads `fail`, and the verb exits 1.
- [ ] `bench anchors <dir>` for a directory with no anchored file below it prints `files[0]{file,anchors,diagnostics,verdict}:` and an empty help block, and exits 0.
- [ ] `bench anchors AGENTS.md` prints the same output as before this change.
- [ ] `bench help` shows `bench anchors <file|dir>`, and the usage line of the verb shows the same operand.
