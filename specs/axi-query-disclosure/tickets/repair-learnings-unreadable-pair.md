# Repair learnings unreadable pair

Blocked by: none
Writes: `internal/learnings/`

## What to build

Cover the open/stat unreadable class (review finding R6): the QD6 "unreadable
file" state is currently exercised only through the oversized read-error route
(`bounds.Classify` line 81), while the stat route (dangling symlink, classified
unreadable at `classify.go:75,125-129`) has no public-surface pair. Add a
public old-to-new paired fixture that points the journal path at a dangling
symlink and asserts the byte-identical early-refusal contract.

## Acceptance

- [ ] [LU1] (covers QD6) a public `learnings` command test creates a dangling
  symlink at `capture/learnings.md`, asserts the structured unreadable refusal
  on stdout with exit 1 as an exact old/new byte pair, and names the state
  distinctly from the oversized fixture (no borrowed oracle, no alias without
  both states named).
- [ ] [LU2] (covers QD6) the pair is skipped cleanly on a filesystem that
  cannot create symlinks, using the existing capability-skip posture rather
  than a silent pass.
