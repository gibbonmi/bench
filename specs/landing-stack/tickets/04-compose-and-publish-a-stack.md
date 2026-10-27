# Compose, grade once, and publish a green stack

Blocked by: 02-prove-a-stack-in-the-gate.md
Writes: internal/landing/stack.go (new), internal/landing/stack_test.go (new), internal/landing/landing.go, internal/gate/stack_completion.go
Covers: LS24, LS25, LS26, LS27, LS28, LS36, LS38, LS39

## What to build

Add `Owner.LandStack(ctx, request, admission)` to `internal/landing` in the new file `stack.go`. The request names the destination ref, the destination base, the green marker, the destination fingerprint, and one `ReviewedRequest` for each member in file order. The admission is a new `StackAdmission` interface with two methods. `CheckMember(i, tree)` decides member `i` over its own tree. `PublishStack(k, trees, publish)` decides the first `k` members again and runs `publish`.

`LandStack` does these steps:

1. It composes each member onto the previous member merge with `Owner.Compose`, and the first member composes onto the destination base. A conflict returns a `ConflictError` that names the member index, before any grade.
2. It applies the tickets-only folder close of each member with `published.Tree`. It then writes one merge commit with `commit-tree`: the previous merge or the destination as first parent, the member tip as second parent, and the member's message. The commits stay unreachable until publication.
3. It computes the own landing tree of each member: the composition of the destination base and that member alone, after its transform. `CheckMember` decides each member over that tree, never over the top tree.
4. It grades the top tree once with `gate.WithLanding` and `gate.WithStack`, and the context holds one `gate.StackMember` for each member.
5. On a green, it rechecks the branch tip and the checkout fingerprint of each member and the destination, as `landReviewed` does. Then `PublishStack` runs one `updateRef` of the destination from the destination base to the top merge.
6. On any other kind, it publishes nothing and returns the authorization result. A red top tree stops here until ticket 08 adds the bisect.

A member request with a staged-spec path returns an error before composition. The stack never transitions a staged spec, and ticket 10 gives that member the `stack-member-alone` face in the preflight.

The result names each published merge, the published count, and the authorization kind. Ticket 07 consumes this result and the `StackAdmission` interface. Ticket 08 adds the bisect behind the same result.

`landing.go` has 32 lines of headroom, so the stack takes its own file. `landing_helpers_test.go` has 12 lines of headroom, so the multi-member fixtures live in `stack_test.go`, and ticket 08 adds its bisect rows in a second test file.

## Acceptance

- [ ] A green three-member stack calls the authorization once, with the top tree (LS24).
- [ ] A green stack publishes a chain whose first parents run from the destination base through each member merge in file order, with each member tip as second parent and each member's message (LS25).
- [ ] A green stack moves the destination ref once, from the destination base to the top merge (LS26).
- [ ] The tree of the first merge lacks the tickets folder of the first member, and the tree of the second merge lacks both folders (LS27).
- [ ] The admission receives the own landing tree of each member and never the top tree (LS28).
- [ ] An infrastructure result on the top tree stops the stack after one grade with no ref update (LS36).
- [ ] A conflict of member 2 with member 1 refuses with a conflict error that names member 2, and the stack makes no grade (LS38).
- [ ] A member whose branch tip moves during the grade refuses the publication with no ref update (LS39).

The tests stub the `authorize` field of the owner, as `greenOwner` in `internal/landing/landing_helpers_test.go` does, and record each graded tree. A multi-member fixture in `stack_test.go` builds two or more source worktrees from one review base, and ticket 08 reuses it.

## Checkpoint verification

- `bench test --package ./internal/landing`
- `bench test --check root-conformance`

## Mutation proof

Pass the top tree to `CheckMember` in place of the member's own tree. The LS28 test must fail. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
