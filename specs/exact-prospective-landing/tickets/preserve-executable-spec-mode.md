# Preserve executable spec mode

Accepted exact-candidate review repair for finding P1. Prospective composition
currently hardcodes a transitioned spec entry to Git mode `100644`, while
reconciliation preserves the invoking checkout's executable bit before staging.
An executable staged spec can therefore publish one mode and leave the real index
at another.

## What to build

Derive the transitioned spec entry's Git mode from the attributed filesystem entry,
using Git's regular-file modes (`100644` or `100755`) rather than carrying arbitrary
permission bits. Keep the transformed bytes and selected mode identical across the
authorized tree, published commit, and reconciled index/worktree.

## Acceptance

- [ ] [MR1] A staged executable spec is composed and published as `100755`; authorization observes implemented bytes, and the invoking index/worktree is clean after reconciliation with the executable bit preserved.
- [ ] [MR2] A staged non-executable spec remains `100644`; its narrower filesystem permissions remain preserved in the worktree, and existing red/CAS-loss behavior remains unchanged.
- [ ] [MR3] Focused landing spec-transition tests and full `internal/landing` complete green.

Ownership fence: internal/landing/landing.go, internal/landing/landing_test.go

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| MR1 | hardcode the transitioned entry to `100644` again | landing spec-transition test | stage an executable spec, land it, then assert commit/index mode identity and clean reconciliation; expect mode drift |

