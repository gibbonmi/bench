# Split every map with the migration program

Blocked by: parse-ticket-files-beside-the-inline-shape.md
Writes: scripts/split-decision-maps/main.go (new), decisions/, specs/decision-map-split/decisions/, docs/research/ (new), .gitignore
Covers: DS34, DS35, DS36, DS37, DS38, DS39

## What to build

Build a one-shot Go program at `scripts/split-decision-maps/main.go`. Run it as `go run ./scripts/split-decision-maps [--apply]`. Without `--apply` it prints every target and writes nothing. The targets are each map with its ticket files, each asset with its destination, each tracked file it will rewrite, and the ignore line. With `--apply` it does the work in one pass over `decisions/` and every `specs/*/decisions/`.

For each map it parses the inline shape with the expand-era parser and writes the index and one file per ticket. The index keeps the title, Status, Destination, and the four terminal sections. It adds an empty `## Notes`. It adds a `## Decisions so far` with one gist per resolved ticket, seeded from the answer's first sentence.

The program moves each asset that the map's Sources names into `decisions/<topic>/assets/`. It rewrites every reference to the old path in tracked files under `decisions/`, `specs/*/decisions/`, and `docs/`. It never edits `specs/*/spec.md` or `specs/*/tickets/`. It moves an asset that no map names to `docs/research/<slug>.md`, then removes `decisions/assets/`. It replaces the `research/` ignore line with `/research/` under a comment that names the shift-scratch folder.

Before the commit, run `bench maps` and save its rows. After `--apply`, run it again, drop the `path` column if the projection ticket has not landed, and diff the two. Record the dry-run plan, the two row sets, the diff, the `rg decisions/assets/` result, the `git ls-files docs/research` result, and both `git check-ignore` results in the ticket's return. Report any remaining hit for the old paths outside the rewrite scope.

## Acceptance

- [ ] [M1] The dry run prints the plan and leaves `git status` clean.
- [ ] [M2] After `--apply`, no index under `decisions/` or `specs/*/decisions/` holds an `## #` line, and every map validates.
- [ ] [M3] Every index holds `## Notes` and a `## Decisions so far` with one gist per resolved ticket.
- [ ] [M4] Every Sources Path names a file under that map's assets folder, and no reference under `decisions/`, `specs/*/decisions/`, or `docs/` names `decisions/assets/`.
- [ ] [M5] `docs/research/ft191-resolved-reader-research.md` is tracked, and `decisions/assets/` does not exist.
- [ ] [M6] `git check-ignore docs/research/x.md` exits 1, and `git check-ignore research/x` exits 0.
- [ ] [M7] The `bench maps` rows before and after the migration are equal once the `path` column is dropped.
- [ ] [M8] The return cites the program's source file and pastes the dry-run output.
