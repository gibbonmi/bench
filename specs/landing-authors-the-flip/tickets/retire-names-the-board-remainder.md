# Make bench spec retire name the board remainder

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/spec_test.go

## What to build

`bench spec retire` reads the spec's `Roadmap:` value through the existing
metadata parser before it deletes the folder. A value that matches `FT`
followed by digits names the `ROADMAP.md` row `FT<n>` in the `next:` line.
When `roadmap/FT<n>.md` exists as a regular file, the line names that path
too. Any other value, or no line, prints the generic line that names the row
and its `roadmap/FT<n>.md` detail file. The deletion order and every refusal
stay as they are. This ticket is disjoint from the flip-author chain and lands
first on its own green gate.

## Acceptance

- [ ] FC1: retire on a spec with `Roadmap: FT7` and an existing `roadmap/FT7.md` prints a `next:` line that names `FT7` and `roadmap/FT7.md`.
- [ ] FC2: retire on a spec with no `Roadmap:` line prints a `next:` line that names the row and `roadmap/FT<n>.md`.
- [ ] FC3: retire on `Roadmap: FT7` with no `roadmap/FT7.md` on disk names `FT7` and no detail path.
- [ ] FC4: retire on `Roadmap: ft7`, `Roadmap: FT 7`, and `Roadmap:` with no value prints the generic line.
- [ ] FC5: retire on a staged spec still refuses with exit 1 and deletes nothing.
- [ ] FC6: a `Roadmap: FT7` line inside a fenced code block is not read, so the generic line prints.
