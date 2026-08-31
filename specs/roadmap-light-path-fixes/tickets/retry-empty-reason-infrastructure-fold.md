# Retry one verified infrastructure fold

Blocked by: none
Writes: internal/worktree/merge.go, internal/worktree/merge_test.go
Covers: LF5

## What to build

Give an infrastructure-attributed merge or fold refusal with no reason one
focused verification. Retry the unchanged operation once only when that check
is green.

## Acceptance

- [ ] An empty-reason infrastructure refusal runs focused verification.
- [ ] A green verification permits exactly one unchanged retry.
- [ ] Reasoned, diff-owned, or twice-refused failures do not retry.

