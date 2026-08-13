# Compose reviewed source commits with the destination

Blocked by: none
Writes: internal/landing

## What to build

Extend the existing deep landing owner with a committed-source composition
request. It accepts immutable destination and source-tip commits, asks Git for
the real merge base, and returns either the ordinary merged tree or one bounded
conflict classification without touching an index, worktree, ref, or merge-state
file. The frozen review base remains authorship metadata; it is never substituted
for Git's merge base.

Exercise true divergence, the FT198 partial-ancestry graph, and Git's conflict
families at this seam. This ticket is independently green as a reusable landing
primitive; it does not expose mutating porcelain.

## Acceptance

- [ ] Diverged committed histories preserve destination-only content and apply
      source-only content once (covers PL6).
- [ ] An FT198-shaped history applies only the non-ancestral source commits once
      in the returned tree (covers the composition half of PL7).
- [ ] Textual, modify/delete, rename/rename, file/directory, mode, symlink, and
      gitlink conflicts return bounded classifications without mutation (covers
      the composition half of PL8).
