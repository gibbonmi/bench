# Union-merge the phase-owned journals

Blocked by: make-spec-optional-on-the-landing.md
Writes: internal/landing/landing.go, internal/landing/composition_test.go, internal/worktree/land_surface_test.go
Line: opus / medium — the composition owner's conflict path is the known seam, with one plumbing uncertainty in the three-blob union.

## What to build

The composition owner's `CaptureSide` becomes a rule table with three verbs:
`source`, `destination`, and `union`. `capture/session-handoff.md` takes
`source`. `capture/learnings.md` and `capture/IDEAS.md` take `union`. Every
other `capture/` path takes `destination`. A path outside the table has no
rule, so the whole conflict refuses and names every conflicted path.

The union
is Git's three-way union merge over the three stage blobs; with one stage
absent, the present side's blob is the result. Only regular-file stages
qualify; any other object kind under the table refuses with its kind. The
landing's stderr disclosure names each settled path with its verb.

## Acceptance

- [ ] A landing whose source and destination both appended distinct lines to `capture/learnings.md` publishes a file that holds both sides' lines (covers WL11).
- [ ] The same shape on `capture/IDEAS.md` publishes the union (covers WL12).
- [ ] A conflicted `capture/session-handoff.md` publishes the source bytes (covers WL13).
- [ ] A conflicted capture file outside the named rules publishes the destination bytes (covers WL14).
- [ ] A union path deleted on one side and appended on the other publishes the appended side's bytes (covers WL15).
- [ ] A conflict on a code path and a capture path together refuses and names both paths (covers WL17).
- [ ] A landing that settled a journal by union prints one `landing composition{resolved=...}` line that names the path and `union` (covers WL19).
- [ ] A symlink or gitlink conflict under `capture/` refuses with the conflict kind and the path (covers WL20).
- [ ] Both sides adding a journal absent from the merge base publishes the union of the two present stages (edge under WL11).
- [ ] A path with spaces in a conflict record is split on the record's tab, not on spaces (edge under WL11).
