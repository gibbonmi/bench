# Name the reauthorize recovery when the request token is lost

Blocked by: enrich-refusals-through-one-emitter.md
Writes: internal/worktree/land.go, internal/worktree/land_test.go, internal/worktree/ownership.go, internal/worktree/ownership_test.go, internal/worktree/worktree.go, internal/worktree/worktree_test.go, internal/worktree/lifecycle_test.go
Line: opus / medium — one refusal branch plus a hint; the authority design does not move.

## What to build

When a land or release refuses because the request digest matches no
assignment, and the target path holds exactly one active assignment, the
refusal names that assignment id and a `bench worktree reauthorize`
continuation template as `next` through the emitter ticket's typed refusal:
assignment id and path filled, placeholders for the new token and any identity
the refusing command does not hold, and placeholder text naming the
full-identity requirement (land refuses abbreviations even though reauthorize
resolves them). The stored digest is never accepted as the caller token — a
ledger record cannot authenticate itself; reauthorize's tree-derived proofs
stay the trust root. No hint appears when the path holds no assignment, or
more than one.

## Acceptance

- [ ] A land with an unknown request against a path holding one active assignment refuses and names the assignment id and the reauthorize continuation template (covers LR20).
- [ ] A release with an unknown request against the same shape names the same recovery in its stderr refusal fields (covers LR21).
- [ ] A land passing the stored digest itself as --request refuses and publishes nothing (covers LR22).
- [ ] A land with an unknown request against a path with no assignment refuses without a reauthorize hint (covers LR23).
- [ ] A path holding more than one active assignment gets no hint (edge from the inventory).
