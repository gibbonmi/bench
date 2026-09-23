# Refuse only tracked changes and collisions at the landing destination

Blocked by: none
Writes: internal/worktree/landingpolicy, internal/worktree/land.go, internal/worktree/land_identity.go, internal/worktree/land_refusal.go, internal/worktree/land_resume.go, internal/worktree/reset.go, internal/worktree/identity_component_test.go, internal/worktree/land_facts_test.go, internal/worktree/land_local_capture_test.go, internal/worktree/land_release_refusal_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_surface_test.go, internal/worktree/merge_test.go, CHANGELOG.md, projects/benchkit.md, roadmap/FT243.md
Covers: none

## What to build

The landing destination refuses a tracked change, as before. An untracked or ignored file in the landing checkout is the operator's own and never reaches a commit, so it does not refuse the landing. A file such as `.env` needs no build-output declaration.

One exception stays. Git refuses to overwrite an untracked file where the landing writes, and it overwrites an ignored file there without a warning. So the landing refuses, before the gate, an untracked or ignored file at a path that the reviewed source tree holds. The resume applies the same collision rule against the published tree before its reset of the checkout. The reset verb's ignored-collision check reads the same rule, so the rule has one source.

## Acceptance

- [ ] A landing with an untracked file and an undeclared ignored file in the destination publishes and keeps both files.
- [ ] A landing with an untracked or ignored file at a path the source adds refuses before the gate and names only that path.
- [ ] A landing with a tracked change in the destination refuses and names the changed path.
- [ ] A resume completes with an untracked or undeclared ignored file outside the published tree.
- [ ] A resume refuses an untracked or ignored file at a path the published tree holds.
