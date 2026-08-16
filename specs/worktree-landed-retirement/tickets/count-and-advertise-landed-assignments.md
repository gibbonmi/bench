# Count and advertise landed assignments

Blocked by: none
Writes: internal/worktree/landed.go, internal/worktree/landed_test.go, internal/worktree/classifier.go, internal/worktree/resume.go, internal/worktree/worktree.go, internal/worktree/list.go, internal/worktree/list_actions_test.go, internal/worktree/orphan_sweep_test.go, internal/worktree/orphan_test.go

## What to build

One classifier in `internal/worktree` answers "is this assignment landed?" — exactly
`state=active ∧ landed proof=true ∧ lease≠live`: only `active` records; only a proof of
`true` (a `false` proof, an unresolvable default branch, or an errored proof is not
landed); leases `none`, `dead`, and `unknown` (unparseable regular file) qualify and only
`live` disqualifies. The automatic planner returns retain reason `landed` for such a tree
ahead of every other retain reason (`live-lease` stays distinct), so `bench resume-clean`
prints `retained landed=N` beside `active=N`, prints exactly one line advertising the bare
`bench worktree clean --landed` invocation when N > 0 (never a discard flag), and skips
landed rows in the age-based orphan candidate list; the automatic path still removes
nothing. `bench worktree list` calls the same classifier per row and renders exactly one
`bench worktree clean --landed` help action when at least one row is landed, none when
zero, keeping the per-row `path`/`exec` actions. Demo: cut two worktrees, land one, run
`bench resume-clean` and `bench worktree list`.

## Acceptance

- [x] `(covers LR1)` A pool with one landed and one non-landed active assignment prints `retained active=1 landed=1`, removes nothing, and both trees remain.
- [x] `(covers LR2)` A landed tree with a live lease counts `live-lease=1` and not `landed`.
- [x] `(covers LR21)` A landed tree with a dead lease counts `landed=1`.
- [x] `(covers LR22)` A landed tree with an unparseable regular lease file counts `landed=1`.
- [x] `(covers LR3)` A landed tree with undeclared ignored residue and one with dirty tracked state both count `landed`, not `ignored`, `active`, or `orphaned`.
- [x] `(covers LR4)` An aged landed tree counts `landed` and is absent from the orphan lines.
- [x] `(covers LR23)` An aged non-landed tree still prints its orphan line.
- [x] `(covers LR24)` A landed count above zero prints exactly one `bench worktree clean --landed` line and never `--discard-ignored`.
- [x] `(covers LR5)` With no resolvable default branch, a landed-looking assignment counts `active` and `list` renders no `clean --landed` action.
- [x] `(covers LR25)` With an erroring landedness proof, the assignment counts `active` and `list` renders no action.
- [x] `(covers LR26)` `cleanup-pending`, `recovered`, and `complete` records over landed branches are never counted `landed` and render no action.
- [x] `(covers LR6)` Exactly one `bench worktree clean --landed` help action when at least one row is landed, none when zero, per-row `path`/`exec` actions unchanged.
