# Document the harness record

Blocked by: 06-add-the-bench-harnesses-verb.md, 07-add-the-harness-record-conformance-check.md, 08-add-the-entry-point-parity-conformance-check.md
Writes: .bench/BENCH-reference.md, projects/benchkit.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, CHANGELOG.md

## What to build

A cold session reads one account of the record and finds the data with one
command.

The Hook Layers bullet on the agent-line guard keeps its one-clause why and
names `bench harnesses codex` as the verdict's source. Delete the dated
verdict sentence, because the record's cell now holds that fact with its
source and its date. The Files section's adapter bullet names the record as
the adapter list's source.

One anchor rule holds the new pointer. Add a `Require` entry to
`internal/anchors/registry_data.go` in the shape that file already uses. The
entry names `.bench/BENCH-reference.md`, the `bench harnesses codex` needle,
and a diagnosis. Use the `AfterSpecAuthorization` group, which already holds
the reference's other entries. Add a second entry for the `CHANGELOG.md`
mention of `bench harnesses`.

The profile's conformance table lists both `harness-record` and
`entry-point-parity`. Tickets 07 and 08 each add their own row, so confirm
both rows and their input classes rather than adding a duplicate. Add one
`CHANGELOG.md` entry under Unreleased for the verb and for the `none` route.
Leave `ROADMAP.md` and `roadmap/FT239.md` alone: the landing flips the spec's
status, and a later `bench spec retire` removes the row.

Write `.agents/` paths only. `.claude/commands` is a symlink to
`.agents/commands`, so no `.claude/` path belongs in a rule or in a diff.

## Acceptance

- [ ] The Hook Layers agent-line bullet contains `bench harnesses codex`. (covers HC45)
- [ ] The profile's conformance table lists `harness-record` and `entry-point-parity`. (covers HC46)
- [ ] `CHANGELOG.md` names `bench harnesses` under Unreleased. (covers HC47)
- [ ] A removed pointer sentence in `.bench/BENCH-reference.md` reds the anchors registry test.
- [ ] The Files section's adapter bullet names `internal/harnesses`.
- [ ] The diff writes no `.claude/` path.
