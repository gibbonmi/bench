# Link a landed fix to the open learning it resolves

Blocked by: none
Writes: internal/learnings/entry.go, internal/learnings/learnings.go, internal/learnings/entry_test.go, internal/roadmap/learning.go, internal/roadmap/learning_test.go, internal/roadmap/context_parse.go, internal/roadmap/context_render.go, internal/roadmap/context_command_test.go, cmd/bench/main.go, cmd/bench/help_inventory_test.go, .agents/commands/bench-drain.md, CHANGELOG.md, ROADMAP.md, roadmap/FT394.md
Covers: none

## What to build

`bench learning --resolves "<title>" --commit <sha>` marks one existing open entry as resolved by a commit.
The verb adds one `Resolved by` bullet with the full commit ID to that entry in `capture/learnings.md`.
The heading keeps its `[open]` marker.

The journal parser reads the bullet. The `learnings` block of `bench roadmap --context` shows it in a `resolved_by` column.
The drain reads that column and closes the entry by that commit. It no longer pairs a second `Resolved:` entry by hand.
The same change retires the roadmap row FT394.

## Acceptance

- [ ] `bench learning --resolves "<title>" --commit <sha>` writes the full commit ID into the one open entry with that exact title. The entry stays open and parses with no malformed record. Other entries are unchanged.
- [ ] The verb exits 1 and leaves the journal unchanged when no open entry has the title. It does the same when two open entries have it, when the entry already has a link, or when the commit does not resolve.
- [ ] The verb exits 2 with the usage line when `--resolves` has no `--commit`, or `--commit` has no `--resolves`. It does the same when a title, `--what`, `--right`, or `--rule` joins `--resolves`.
- [ ] `bench roadmap --context` reports the commit in the `resolved_by` column of the entry. An unlinked entry has an empty cell. The index and the full snapshot both show the column.
- [ ] The root help inventory lists the new form, and the drain guidance tells the drain to read `resolved_by`.
