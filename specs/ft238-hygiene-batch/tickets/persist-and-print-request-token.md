# Persist and print the request token

Blocked by: none
Writes: internal/intent/assignment.go, internal/worktree/ownership.go, internal/worktree/list.go, tests beside each

## What to build

`bench worktree create` persists the original `--request` token in the
assignment record, beside the digest that stays the authorization identity.
`bench worktree list` prints the token in a new column. A resumed landing can
then reuse the token without a refusal round-trip. Records written before the
field existed carry none and list an empty cell.

## Acceptance

- [ ] After create, the assignment record carries the plain token and the digest still validates.
- [ ] `bench worktree list` prints the token beside the assignment.
- [ ] `bench worktree land` succeeds with the token that list printed.
- [ ] A pre-field record still loads, and its list row shows an empty token cell.
