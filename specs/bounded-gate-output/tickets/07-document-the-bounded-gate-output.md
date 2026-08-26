# Document the bounded gate output

Blocked by: 04-print-the-green-phase-table.md
Writes: projects/benchkit.md, .bench/BENCH-reference.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, CHANGELOG.md

## What to build

A cold session reads one account of the bounded shape. The profile's
`## Gate (.bench/gate.sh)` section and `.bench/BENCH-reference.md` each
state three facts: what a green run prints, what a red run prints, and where
the complete stream goes. `.bench/BENCH-reference.md` also says today that a
phase name "appears in summary lines and output prefixes". Repair that sentence, because the summaries go away and the prefix moves into the
`.out` file.

One anchor rule per file holds the new sentence. Add two `Require` entries to
`internal/anchors/registry_data.go`, in the shape that file already uses.
Each entry names the file, an empty section, the sentence needle, and a
diagnosis. Use the `AfterSpecAuthorization` group. The reference already has
gate-shape entries there; the profile entry is new to that group.

Write `.agents/` paths only. `.claude/commands` is a symlink to
`.agents/commands`, so no `.claude/` path belongs in a rule or in a diff.
Add one `CHANGELOG.md` entry for the new output shape. This ticket writes no
file that ticket 06 writes, so the two run in parallel.

## Acceptance

- [ ] The profile's gate section and the reference each carry the bounded-output sentence needle. (BG27)
- [ ] A removed sentence in either file reds the anchors registry test.
- [ ] The reference no longer says a phase name appears in summary lines.
- [ ] The diff writes no `.claude/` path.
