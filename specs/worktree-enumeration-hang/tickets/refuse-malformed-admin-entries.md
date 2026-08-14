# Refuse malformed worktree admin entries before enumerating

Blocked by: resolve-git-common-dir.md
Writes: internal/git, internal/gittest, internal/worktree/list.go, internal/worktree/list_actions_test.go, internal/worktree/resume_test.go, internal/status/status.go, internal/status/status_test.go

## What to build

A discovering command against a repository whose shared admin dir holds a
non-regular entry returns promptly with the attributable typed refusal
instead of hanging — from the primary checkout and from a linked worktree —
per the spec's "One shape owner" and "Refusal shape" decisions. Export
`ScanWorktreeAdmin(commonDir string) error` (the contract the doctor ticket
calls), call it from `git.Worktrees` between the predecessor's fail-closed
resolution and the porcelain launch, and route the typed class at the two
action-cell surfaces: `appendWorktree` renders typed detail and the type's
action field, and `bench worktree list` replaces its fixed
message-and-retry pair with the typed detail and action. The predecessor's
obligations this ticket completes: the FT29 prior-art test's detail
assertion moves to the typed framing (its `chmod 000` fixture now fails at
the rev-parse) and its "rather than a PATH-shimmed git" comment is
rewritten; the scanner carries the retirement comment naming the git 2.43.0
blocking-open behavior. The stub gains block-worktree and fail-worktree
modes. Tests reuse `git_test.go`'s `runGit`/`newRepo` — the census caps
`internal/git` test constructors at one. "Within the bounded wait" means a
goroutine plus a `bounds.TestDeadline(0)`-floor deadline; the production
bound is the bound ticket's. WE4 lands last: its red needs the scanner
already refusing.

## Acceptance

- [ ] A FIFO `gitdir` yields a `git.Worktrees` refusal containing `worktrees/<id>/gitdir`, `fifo`, and `inspect and remove it` within the bounded wait (covers WE1)
- [ ] First-level symlink, first-level FIFO, and symlinked `gitdir` with a regular target are each refused naming the offending path and shape word (covers WE2)
- [ ] Absent or FIFO `worktrees/` enumerates normally (covers WE5)
- [ ] A FIFO named `stray` under a space-and-glob `<id>` is refused naming the id verbatim (covers WE6)
- [ ] Prunable, gitdir-less, empty-`<id>`, and depth-3 `logs/HEAD` FIFO states all still enumerate (covers WE7)
- [ ] From a linked worktree's root the same attributable refusal arrives within the bounded wait (covers WE12)
- [ ] With the stub logging argv, the name-shape-action refusal lands and the log holds no `worktree` invocation (covers WE13)
- [ ] A symlinked `worktrees/` is refused naming `worktrees` and `symlink` (covers WE18)
- [ ] The session-start resume output contains the path, shape, and action, and not `git worktree list failed` (covers WE24)
- [ ] `appendWorktree`: FIFO fixture row carries the typed detail and `inspect and remove it`; fail-rev-parse fixture's detail does **not** contain `git worktree list failed` and its action is `investigate the git failure`, not `inspect and remove it`; fail-worktree fixture keeps the generic detail and re-run action (covers WE3)
- [ ] `bench worktree list`: FIFO fixture error contains `worktrees/<id>/gitdir` and `inspect and remove it`; bad-rev-parse fixture carries the resolution detail with `investigate the git failure`, not `inspect and remove it`; fail-worktree fixture keeps the fixed "cannot read registered worktrees" / retry message (covers WE4)
