# Add a diff path-list bullet to the hostile-input checklist

Blocked by: none
Writes: projects/benchkit.md
Covers: none

## What to build

The hostile-input checklist in the project profile gets one bullet for a predicate over a `git diff` path list. The bullet tells the spec author that the predicate must meet added, edited, deleted, and renamed paths. It also tells the author that a porcelain `git diff` detects renames by default, so the author reads the list with `--no-renames` when each path must be its own change. It tells the author to read `--name-only` and `--name-status` output with `-z`, because git quotes a path that has unusual characters. The FT406 diff-test seam missed deletions and renames, and its repair added `-z` to a name-status diff.

## Acceptance

- [ ] The hostile-input checklist has one bullet that names the add, edit, delete, and rename kinds for a `git diff` path-list predicate.
- [ ] The bullet states the default porcelain rename detection and the `--no-renames` read.
- [ ] The bullet states the path quoting of `--name-only` and `--name-status` and the `-z` read.
- [ ] Each git fact in the bullet agrees with `git diff --help`.
