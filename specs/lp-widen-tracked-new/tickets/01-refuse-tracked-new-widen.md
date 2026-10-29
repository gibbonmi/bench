# Widen refuses a new entry on a tracked path

Blocked by: none
Writes: internal/spec/widen.go, internal/spec/widen_test.go, internal/tickets/writes.go, CHANGELOG.md
Covers: none

## What to build

`bench spec widen` adds a `Writes:` entry with the `(new)` marker even when git tracks its path.
An author then believed that a tracked test file was new and overwrote it.
The `background-grade` retro asks for the refusal.

For each requested `(new)` entry that the ticket does not already cover, read the checkout HEAD through `git.ListTree`.
When HEAD holds the path as a file or a directory, refuse the whole request at exit 1.
The refusal names each such path and tells the caller to drop the `(new)` marker, and the verb writes nothing.
An entry that the ticket already covers keeps the `current` state.
An untracked working file is not a tracked path, so its `(new)` entry still widens.

Export the marker spelling from `internal/tickets` so that the refusal and the test read one source.
Add a `CHANGELOG.md` entry under `[Unreleased]`.

## Acceptance

- [ ] A `(new)` entry on a file that HEAD tracks refuses, names the path, and leaves the ticket unchanged.
- [ ] A `(new)` entry on a directory that HEAD tracks refuses the same way.
- [ ] A `(new)` entry on a file that only the worktree branch commits refuses the same way.
- [ ] A `(new)` entry on an untracked working file widens.
- [ ] A covered `(new)` request entry still reports `current`.
- [ ] The `bench commit` lane passes.
