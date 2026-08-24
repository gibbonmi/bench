# Refuse bench idea on the primary checkout

Blocked by: none
Writes: internal/usage/worktree.go, internal/commit/commit.go, internal/roadmap/roadmap.go, internal/roadmap/roadmap_test.go, .bench/BENCH.md

## What to build

`bench idea` obeys the same checkout boundary as `bench commit`. If the
command runs in the primary checkout, it refuses with the one shared refusal
line and names the worktree verb. It does not touch `capture/IDEAS.md`. Inside
a linked worktree the command appends as before, and the phase landing carries
the line to `main` through the union rule. `bench commit` and `bench idea`
derive the refusal line from one source.

## Acceptance

- [ ] `bench idea "<text>"` in the primary checkout exits 1, prints the refusal, and leaves `capture/IDEAS.md` unchanged.
- [ ] `bench idea "<text>"` in a linked worktree appends the dated line as before.
- [ ] `bench commit` prints the same refusal text from the same source.
