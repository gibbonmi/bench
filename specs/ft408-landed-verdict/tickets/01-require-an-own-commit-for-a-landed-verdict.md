# Require an own commit for a landed verdict on dirty output

Blocked by: none
Writes: internal/worktree/landed.go, internal/worktree/list.go, internal/worktree/clean_landed.go, internal/worktree/classifier.go, internal/worktree/resume.go, internal/worktree/landed_test.go, internal/worktree/landed_verdict_test.go, internal/worktree/parallel_census_test.go, ROADMAP.md, roadmap/FT408.md, CHANGELOG.md
Covers: none

## What to build

Today an active assignment is landed when its branch tip is an ancestor of `main`. This is also true when the assignment made no commit and its tree holds uncommitted output. The list cell, the clean of the landed set, the session-start count, and the orphan sweep then treat that output as landed work.

Take the first-parent chain from the branch tip back to the recorded start. The assignment has an own commit when that chain holds a commit that is not on the first-parent history of the default branch. A landing publishes a two-parent commit with the assignment tip as the second parent, so a landed own commit is never on that history. A fast-forward merge of `main` moves the tip past the start over commits that are all on that history, so it adds no own commit.

One shared fact in `internal/worktree/landed.go` reports an assignment with no own commit as holding unlanded output when its tree is dirty. Dirty means the tracked and untracked status of the checkout is not empty. The scoped selector of the landed set reads the same own-commit fact as its no-contribution test.

The landed verdict of an active assignment reads this fact, and so does the `landed` cell of `bench worktree list`. The verdict feeds the list banner, the selector of `bench worktree clean --landed`, the resume sweep, and the automatic classifier. The explicit planner keeps its ancestry proof, because its branch removal and its dirty retain reasons depend on it.

An assignment with no commit and a clean tree stays landed, as the code defines it today. An assignment with an own commit that `main` contains stays landed, with or without dirty output. The reviewer decides the salvage of each unlanded worktree.

## Acceptance

- [ ] `bench worktree list` shows `landed=false` for an assignment with no own commit and an untracked file.
- [ ] `bench worktree clean --landed` does not select that assignment.
- [ ] The session-start summary does not count that assignment as landed.
- [ ] An assignment with no own commit and a clean tree still shows `landed=true` and still counts as landed.
- [ ] An assignment with an own commit that `main` contains and a dirty tree still counts as landed.
- [ ] A dirty assignment with no own commit shows `landed=false` and does not count as landed after a fast-forward merge of `main`. A landing merge in that range changes nothing.
- [ ] The orphan sweep reports an aged assignment with no own commit and a dirty tree, and it omits the clean one.
- [ ] The FT408 heading in `ROADMAP.md` and the file `roadmap/FT408.md` are removed, and `CHANGELOG.md` records the change.
