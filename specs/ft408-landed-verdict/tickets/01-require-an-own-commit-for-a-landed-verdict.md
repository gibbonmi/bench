# Require an own commit for a landed verdict on dirty output

Blocked by: none
Writes: internal/worktree/landed.go, internal/worktree/list.go, internal/worktree/clean_landed.go, internal/worktree/subshell.go, internal/worktree/landed_test.go, ROADMAP.md, roadmap/FT408.md, CHANGELOG.md
Covers: none

## What to build

Today an active assignment is landed when its branch tip is an ancestor of `main`. This is also true when the assignment made no commit and its tree holds uncommitted output. The list cell, the clean of the landed set, the session-start count, and the orphan sweep then treat that output as landed work.

An assignment has no commit of its own when its branch tip equals its recorded start. One shared proof in `internal/worktree/landed.go` reports such an assignment as not landed when its tree is dirty. Dirty means the tracked and untracked status of the checkout is not empty.

Every site that today calls `git.LandedInDefault` for an assignment calls this proof instead. These sites are the `landed` cell of `bench worktree list` and the selector of `bench worktree clean --landed`. The explicit planner is the third site. Its `LandedOK` fact feeds the resume sweep and the automatic classifier.

An assignment with no commit and a clean tree stays landed, as the code defines it today. An assignment with an own commit that `main` contains stays landed, with or without dirty output. The reviewer decides the salvage of each unlanded worktree.

## Acceptance

- [ ] `bench worktree list` shows `landed=false` for an active assignment with no own commit and an untracked file.
- [ ] `bench worktree clean --landed` does not select that assignment.
- [ ] The session-start summary does not count that assignment as landed.
- [ ] An assignment with no own commit and a clean tree still shows `landed=true` and still counts as landed.
- [ ] An assignment with an own commit that `main` contains and a dirty tree still counts as landed.
- [ ] The FT408 heading in `ROADMAP.md` and the file `roadmap/FT408.md` are removed, and `CHANGELOG.md` records the change.
