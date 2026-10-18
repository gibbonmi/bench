# Add bench spec widen

Blocked by: 02-delete-ownership-closure.md
Writes: internal/spec/spec.go, internal/spec/widen.go (new), internal/spec/widen_test.go (new), internal/spec/retire.go (new), internal/spec/spec_test.go, internal/spec/retire_test.go (new), cmd/bench/main.go, cmd/bench/help_inventory_test.go, cmd/bench/command_registry_test.go, CHANGELOG.md, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go
Covers: CR23, CR24, CR25, CR26, CR27, CR28, CR29, CR30, CR31, CR32, CR33, CR34, CR35, CR36, CR37, CR38, CR39, CR40, CR41, CR42, CR43, CR44, CR45, CR46, CR47, CR48, CR49, CR50, CR51, CR52, CR53, CR54, CR92

## What to build

`bench spec widen <slug> --ticket <basename> [--decision <reference>] <entry>...` is a new leaf of the existing `bench spec` command. It adds each in-scope entry to the one `Writes:` line of `specs/<slug>/tickets/<basename>` in the current worktree, and it makes no commit.

The approved scope comes from the default-branch tip, which `git.ResolvedDefault` names. The verb reads the slug's ticket files at that tip through the git adapter. The scope is the set of directory prefixes of every signed-off `Writes:` entry. The prefix of entry E is E when the tip holds E as a directory, otherwise E's parent directory, and a root-level file is its own prefix. `tickets.Covers` decides membership.

The verb grades each entry through `tickets.ParseWrite`, which refuses a `..` segment after ticket 02. It also refuses an entry that holds `, ` and an entry that TOON cannot render. An entry that the ticket's own `Writes:` already covers reports `current` and writes nothing. Each other in-scope entry reports `added`. An out-of-scope entry exits 1 with a refusal that names the entry and the `--decision <reference>` rerun. With a non-empty `--decision`, an out-of-scope entry reports `decided`. Any refusal writes nothing for the whole request. The verb prints one `widen[N]{ticket,entry,state}` table through `toon.Table`.

The verb refuses the primary checkout through the shared refusal, a default tip with no ticket of the slug, a ticket file that is a symbolic link or a special file, and a ticket with no `Writes:` line. It keeps every byte outside the `Writes:` value, including a CRLF ending and a missing final newline. It resolves the ticket from the repository root, so a deeper working directory edits the same file. It writes through `durablefile`.

`bench help` and `bench spec --help` list the leaf, and the changelog's `[Unreleased]` section names it. `internal/spec/spec.go`, `internal/spec/spec_test.go`, and `cmd/bench/main.go` are over their structure budget, so each edit there is net zero. The retire leaf moves to `retire.go`, and its tests move to `retire_test.go`, to give the switch and the pins their headroom. The help test that names only retire and history takes a name that includes widen, with no added line.

Contract for ticket 04: the `source-not-fenced` route names this verb by its exact spelling `bench spec widen`.

## Acceptance

- [ ] An in-scope entry is appended to the `Writes:` value, prints state `added` at exit 0, and leaves `HEAD` unchanged.
- [ ] An out-of-scope entry exits 1 with unchanged bytes, and the refusal names the entry and the `--decision <reference>` rerun.
- [ ] `--decision r1` admits an out-of-scope entry with state `decided`, and an empty `--decision` exits 2.
- [ ] The scope reads the default-branch tip, so a directory that only a worktree ticket names stays out of scope.
- [ ] The directory-prefix rule admits a sibling subtree and refuses a sibling directory, another directory entry's path, and another root file.
- [ ] A covered entry reports `current` with unchanged bytes, and a mixed request writes nothing.
- [ ] The verb keeps CRLF and missing-final-newline bytes and refuses each hostile entry and ticket class in the coverage rows.
- [ ] `bench help`, `bench spec --help`, and the changelog name the leaf.
- [ ] `bench structure --growth` reports no growth for `spec.go`, `spec_test.go`, and `main.go`.
