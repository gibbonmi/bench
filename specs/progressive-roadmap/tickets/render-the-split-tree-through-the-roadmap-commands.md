# Render the split tree through bench roadmap and its context snapshot

Blocked by: split-the-board-parser-and-migration-in-one-green.md
Writes: internal/roadmap

## What to build

`bench roadmap`, `bench roadmap --context [--full]`, and `--context --row <ID,...>`
run on the loader: the sources block gains a `roadmap/` row with state and byte
total, each integrity diagnostic renders as a `parse_failures` row whose source is
the offending path with `ROADMAP.md`'s state `malformed` and `sequence_trusted`
false, `--row` returns the row-file body and byte count, and the bare command keeps
its exit-1 postures (empty file, failed read, unsupported schema) while rendering
top rows from the split tree. A row file the classifier cannot read renders a `parse_failures` row naming it with `sequence_trusted` false. Schema stays 4. Coverage rows PR11, PR12, PR13, PR14, PR27.

## Acceptance

- [ ] `--context` on a split tree renders `roadmap/,parsed,<bytes>` under sources and `4` in the context row.
- [ ] `--context` over a tree missing `roadmap/FT7.md` renders a `parse_failures` row sourced `roadmap/FT7.md`, `ROADMAP.md,malformed`, and `sequence_trusted` false.
- [ ] `--context --row FT7` renders the row-file body as `body` and its length as `body_bytes`.
- [ ] `--context` over a tree whose `roadmap/FT7.md` is a directory renders a `parse_failures` row sourced `roadmap/FT7.md` and `sequence_trusted` false.
- [ ] `bench roadmap` renders `title` from the index and `spec`, `occurrence_count`, `occurrence_keys` from row-file bodies; it renders top rows from index headings and detail bodies; an empty `ROADMAP.md` still exits 1 with the record error; a root with neither file renders the absent posture.
