# Repair the sha route and its tests after review

Blocked by: 01-name-the-refusal-routes.md
Writes: internal/worktree/merge.go, internal/worktree/merge_from_sha_test.go, internal/worktree/merge_test.go, internal/worktree/land_identity_test.go, internal/worktree/land_fixtures_test.go
Covers: none

## What to build

The review found that a fresh sibling sits at the default-branch tip. So the
sha route made `--from <default-branch tip sha>` refuse as ambiguous, and it
broke the ordinary fold. The repair orders the lookups. The label lookup and
the default-branch commit lookup keep their order and their ambiguity
refusal. The sha route runs only for a commit the default branch does not
own. It lives in the merge verb alone, so the create verb's `--from` keeps
its label-only contract.

On the sha route, a commit equal to the target's own tip refuses with the
same detail the target's own label gives. A commit equal to two or more
active siblings' tips refuses and names the sibling ids. Exactly one match
selects that sibling, and the sibling checks run once. The sha match returns
the tip it read, so the selected sibling's tip is read one time.

The fold tests for the label spelling and the sha spelling share one helper
that takes the `--from` spelling. The sha fold asserts the kind, the subject,
and the publication, the same as the label fold. The off-branch refusal test
asserts the full detail sentence with both routes. The sha route also covers
a sibling whose checkout is behind its branch tip.

The land refusal test passes a fold commit as `--base` with a later source
tip. The destination advances by one commit, the source folds it, and one
more commit follows the fold. If the fixture cannot be built inside the cap,
the current test stays and the return says so.

## Acceptance

- [ ] `bench worktree merge --from <default-branch tip sha> <target>` composes the default-branch fold when a fresh sibling sits at that tip.
- [ ] `bench worktree merge --from <sha> <target>` with a sha equal to the target's own tip refuses with `--from resolves to the target itself`.
- [ ] `bench worktree merge --from <sha> <target>` with a sha at two active siblings' tips refuses and names both sibling ids.
- [ ] `bench worktree create --from <sha>` still refuses with `--from names no active assignment`.
- [ ] The off-branch refusal test reds when the route clause leaves the detail.
- [ ] The sha fold test reds when the fold's kind or subject differs from the label fold's.
- [ ] `bench worktree land --base <fold commit> --source-tip <later tip>` refuses with the two-bases detail.
