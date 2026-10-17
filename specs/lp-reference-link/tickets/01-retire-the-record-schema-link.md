# Retire the record schema link in the platform reference

Blocked by: none
Writes: .bench/BENCH-reference.md
Covers: none

## What to build

The `Assessment records` section of `.bench/BENCH-reference.md` links the record schema at `../internal/assessment/README.md`.
Every linked repo receives `.bench/`, but only the Bench kit repository holds `internal/`.
Thus the link is dead in every linked repo.

Replace the link with plain text that names the schema path and states that the path is in the Bench kit repository.
The next paragraph names the same path again for the normalized input.
Change that sentence to refer back to the schema, so that one sentence holds the path.

No check grades a link from a platform file toward a kit-only path.
This ticket adds no such check.

## Acceptance

- [ ] `.bench/BENCH-reference.md` holds no Markdown link to a path under `internal/`.
- [ ] The `Assessment records` section names `internal/assessment/README.md` once, as plain text, and states that the path is in the Bench kit repository.
- [ ] Each other sentence of the section is unchanged.
