# Resolve the git common dir through one fail-closed owner

Blocked by: none
Writes: internal/git, internal/gittest, internal/intent/intent.go, internal/worktree/subshell.go, internal/worktree/worktree.go, internal/worktree/classifier.go, internal/worktree/lifecycle.go, internal/worktree/resume.go, internal/gate/engine.go, internal/dashboard/render.go, internal/status/status.go

## What to build

`git.Worktrees` fails closed on common-dir resolution instead of trusting it
blindly, per the spec's "One common-dir owner" and "Refusal shape" decisions:
export `CommonDir(root string) (string, error)` over one unexported rev-parse
argv helper (the contract the bound ticket drives through its variant), land
the exported typed error class whose `Error()` text carries path, shape word,
and action clause, migrate the eight production rev-parse sites to
`CommonDir` (deleting `commonGitDir` in `internal/gate/engine.go`), and make
the three class-specific labels (resume's wrap, the dashboard template,
`PruneLandedBranches`' wrap) the neutral `worktree discovery failed`. A
thinner cut strands nothing but delivers nothing either: WE19 and WE20 are
the complete tracer — both red against today's tree, where `Worktrees` runs
no rev-parse at all. The `internal/gittest` stub-git builder lands here
(pure file-writing, argv-aware rev-parse handling, bad-rev-parse and
fail-rev-parse modes) with the package charter comment widened from
repositories-only to the shared git test scaffolds. Tests reuse
`git_test.go`'s `runGit`/`newRepo`: the census caps constructor sites in
`internal/git/*_test.go` at one, and the stub builder lives in
`internal/gittest`, outside that count (it also launches no process — that
avoids racing a real git against millisecond overrides, not the census). "Within the bounded wait" in this spec means a goroutine plus
a `bounds.TestDeadline(0)`-floor deadline — no production bound exists yet.

## Acceptance

- [ ] With the stub in bad-rev-parse mode, `git.Worktrees` refuses naming the garbage path, and the argv log holds no `worktree` invocation (covers WE19)
- [ ] With the stub in fail-rev-parse mode, `git.Worktrees` returns the typed resolution failure naming the rev-parse invocation, and the argv log holds no `worktree` invocation (covers WE20)
- [ ] The three relabeled surfaces render `worktree discovery failed` in place of their old invocation-naming labels, and each migrated caller resolves the same common dir it did before
