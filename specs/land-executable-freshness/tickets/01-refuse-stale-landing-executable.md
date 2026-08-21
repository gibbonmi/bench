# Refuse a stale landing executable

Blocked by: none
Writes: internal/worktree/land.go, internal/worktree/land_test.go, cmd/bench/main.go, cmd/bench/command_registry_test.go, ROADMAP.md, roadmap/FT242.md

## What to build

`bench worktree land` (first run only, never `--resume`) proves its own
executable before any landing proof. When `scripts/go-build.inputs` exists at
the landing root in any form, the command runs `freshness.Verify(root,
executable)` through a package-var seam and refuses on any error with the
existing `refused{detail=...}` line, exit 1, passing the owner's message (it
carries the rebuild remedy) through the existing sanitized refusal path. Only a
not-exist `Lstat` skips the check. `LandCommand` gains the invoked-executable
parameter, supplied by the registry closure from `Command.Executable` — the
same source `freshness-check` uses. The ticket also rewrites `roadmap/FT242.md`
and its `ROADMAP.md` index line to the re-scoped decision recorded in the spec.

## Acceptance

- [ ] LF1: manifest present plus an erroring owner → exit 1 and a `refused{detail=...}` line carrying the owner's message.
- [ ] LF2: after that refusal, the destination ref, the green marker, and the assignment state are unchanged.
- [ ] LF3: no manifest at the root → the check is never consulted and the landing reaches its next proof.
- [ ] LF4: `--resume` completes marker, reconcile, and release while the check seam errors.
- [ ] LF5: manifest present plus a sealless fixture executable → the real `freshness.Verify` refuses through the land surface.
- [ ] LF8: a present-but-empty manifest still routes to the check.
- [ ] LF9: with an erroring check and a dirty destination, the emitted refusal is the owner's message, not the destination proof's.
- [ ] LF10: the registry closure forwards `Command.Executable` into the check (command registry test).
- [ ] The board records the FT242 re-scope (roadmap/FT242.md body and ROADMAP.md index line).
