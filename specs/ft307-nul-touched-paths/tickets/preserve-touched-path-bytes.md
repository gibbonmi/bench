# Preserve touched path bytes

Blocked by: none
Writes: internal/structure/structure.go, internal/structure/touched_paths_test.go (new), CHANGELOG.md
Covers: none

## What to build

Read the touched-file query through git.Raw with NUL framing.
Pass each complete filename to the existing structure checker.
Keep the base..HEAD range and the ACMR filter unchanged.
Keep working-tree-only changes outside this query.
Keep the existing source-extension filter.

Add a Fixed entry under a dedicated Touched source paths changelog heading.
The coordinator reconciles the roadmap after the fix lands.

## Acceptance

- [ ] Committed non-ASCII source paths reach the checker unchanged.
- [ ] Committed newline-bearing source paths reach the checker unchanged.
- [ ] Source names with spaces, quotes, tabs, or whitespace before their extension survive the query.
- [ ] Unchanged, deleted, and working-tree-only paths retain their existing scope.
- [ ] A failed Git query retains the existing error behavior.
