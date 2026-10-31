# Read the default-branch tip through git.DefaultTip

Blocked by: none
Writes: internal/preflight/default_branch.go, internal/preflight/default_branch_test.go, internal/commitment/repository/sources.go, internal/commitment/repository/export_test.go, internal/commitment/repository/default_revision_test.go
Covers: none

## What to build

`git.DefaultTip` is the one lookup of the default-branch tip. It answers the commit identity of the tip. Two production sites do not use it. Each site resolves the branch with `git.ResolvedDefault`, and then reads the tip with a plain `rev-parse` of the branch name.

The base-current facts in the build preflight read the tip through `git.DefaultTip`. The merge-base query uses that tip and the same root. The default revision in the commitment repository also reads the tip through `git.DefaultTip`. It keeps `errDefaultUnresolved` as its answer when no default branch resolves.

A plain `rev-parse` of a name resolves a tag before a branch. Thus, an annotated tag named like the default branch gives the tag object and not the commit. With `git.DefaultTip`, each site reads the commit.

## Acceptance

- [ ] When an annotated tag has the name of the default branch, and HEAD contains the default-branch tip, the base-current facts give the branch, resolved, and current.
- [ ] When the default branch does not resolve, the base-current facts give no branch, not resolved, and not current.
- [ ] When an annotated tag has the name of the default branch, the commitment default revision gives the commit of the branch tip and not the tag object.
- [ ] When the default branch does not resolve, a commitment gap write still refuses with the unresolved default-branch error.
