# One green-marker reader

Blocked by: none

Ownership fence: `internal/gate` (a new marker package beneath it), `internal/gate/authorization`, `internal/worktree/land.go`
Assumptions: decision source `decisions/deepening-2026-08.md` #4 (exit-test rule) and #11 (candidate 6); the public seams `gate.ValidateProjectGreen`, `authorization.CheckMarker`, `authorization.AdvanceMarker`, and `worktree.landingMarker`'s callers do not move

## What to build

`refs/bench/green/<branch>` is read three ways today. `gate.ValidateProjectGreen`
resolves the ref without peeling, so a marker that points at a non-commit object
mismatches instead of resolving; `worktree.landingMarker` peels `^{commit}` and turns
every read failure into an empty marker; `authorization.markerCommit` peels and
additionally classifies a dangling symbolic ref as an error rather than absence.
Compare-and-advance lives only in `authorization.checkMarker` alongside its
already-at-destination tolerance.

Collapse the three into one deep module that owns the ref name, the read (peel plus
dangling-symbolic-ref classification), and the compare-and-advance. The three sites
call it; none spells `refs/bench/green/` or the peel again. Callers keep learning one
question — does the marker authorize this tip — through the unchanged
`authorization.CheckMarker` / `AdvanceMarker` seam.

Refactor exit test (map #4): the pre-existing suite passes with test logic unmodified;
mechanical renames are the only permitted test edit, and a changed assertion reverts
the move.

## Acceptance

- [ ] [GM1] Exactly one production source spells the `refs/bench/green/` prefix and the marker peel; `engine.go`, `land.go`, and `authorization.go` all read through it.
- [ ] [GM2] Behavior delta named: `ValidateProjectGreen` now peels, so a marker pointing at a non-commit object that peels to the tip validates instead of reporting "project-green marker changed"; a dangling symbolic marker still reports that reason.
- [ ] [GM3] Every pre-existing test in `internal/gate`, `internal/worktree`, `internal/landing`, and `internal/systemtest` passes with test logic unmodified.
