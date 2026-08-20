# Close the light-path ticket on its landing commit

Blocked by: none
Writes: internal/commit, internal/landing

## What to build

`bench commit --spec <slug>` gains one behavior branch. When `specs/<slug>/`
holds no `spec.md`, the green landing commit deletes that folder as part of the
commit it publishes, rather than flipping a status line. A folder that holds a
`spec.md` keeps today's staged-to-implemented flip exactly. A slug naming no
folder returns a structured error and lands nothing. On a red gate nothing is
deleted.

The contract this ticket sets for `count-tickets-only-folders-in-status.md`: a
tickets-only folder is a direct child of `specs/` that contains no `spec.md`.
Both tickets read that predicate from one place.

## Acceptance

- [ ] a green landing on a tickets-only slug publishes a commit in which the folder is gone (H01).
- [ ] a green landing on a slug holding `spec.md` flips its status and keeps the folder (H02).
- [ ] a slug naming no folder returns a structured error and no commit lands (H03).
- [ ] a red gate leaves the tickets-only folder present in the tree (H04).
