# Bisect a red stack in the landing owner

Blocked by: 04-compose-and-publish-a-stack.md
Writes: internal/landing/stack.go (new), internal/landing/stack_bisect_test.go (new)
Covers: LS29, LS30, LS31, LS32, LS33, LS34, LS35, LS37, LS99

## What to build

Give `Owner.LandStack` the bisect of the spec section "The bisect". Its result keeps the shape of ticket 04: k published merges, the authorization kind, and now the attributed member.

- The bisect continues only on the kinds `green` and `candidate` of `authorization.Kind`. A `candidate` red on the top tree means that the destination held a green, so the destination is prefix 0 and counts as green.
- The search keeps a green bound `lo` at 0 and a red bound `hi` at n. While `hi - lo > 1`, it grades the prefix `mid = (lo + hi) / 2`, rounded down. A green moves `lo` to `mid`, and a red moves `hi` to `mid`.
- Each prefix grade carries `gate.WithLanding` and a `gate.WithStack` context with one `gate.StackMember` for each member of that prefix and no other.
- The owner attributes the red to member `hi` and publishes the merge of prefix `lo` through `PublishStack(lo, …)`, after the same rechecks for the first `lo` members. When `lo` is 0, it publishes nothing.
- Any other kind on any grade stops the stack with no publication and no attributed member. The result carries that kind, so an `infrastructure` result and an `inherited` red stay apart.

The landing package names no worktree face. The verb of ticket 07 already handles k published members and an optional refusal, so this ticket changes no verb file. Ticket 09 maps an attributed member to `stack-red` and an unattributed `inherited` stop to `land-handback`.

## Acceptance

- [ ] Each grade carries a stack context that names exactly the members of its prefix (LS29).
- [ ] With four members and member 3 red, the stack grades prefixes 4, 2, and 3 in that order, which is 1 + ⌈log2 4⌉ = 3 grades (LS30).
- [ ] With four members and member 3 red, the stack attributes the red to member 3 (LS31).
- [ ] With four members and member 3 red, the stack publishes the merge of member 2 and no later merge (LS32).
- [ ] With four members and member 1 red, the stack grades prefixes 4, 2, and 1, attributes member 1, and moves no ref (LS33).
- [ ] When prefixes 1 to 4 grade green, red, green, and red, the stack grades prefixes 4, 2, and 1, attributes member 2, and publishes the merge of member 1 (LS34).
- [ ] With sixteen members and member 16 red, the stack grades prefixes 16, 8, 12, 14, and 15, which is 1 + ⌈log2 16⌉ = 5 grades (LS35).
- [ ] An infrastructure result on a bisect prefix stops the stack with no further grade and no ref update (LS37).
- [ ] An `inherited` red on the top tree stops the stack after one grade, moves no ref, and attributes no member (LS99).

The rows live in the new file `stack_bisect_test.go` and reuse the multi-member fixture of `stack_test.go`. The authorization stub answers each graded tree from a table of prefix verdicts, and each row compares the exact ordered list of graded prefixes.

## Checkpoint verification

- `bench test --package ./internal/landing`
- `bench test --check root-conformance`

## Mutation proof

Grade the prefix `lo + 1` in place of the midpoint, which makes a linear search. The LS30 test must fail, because the search grades prefixes 4, 1, 2, and 3. Spell the probe against the code that the author writes, and require `bench probe` to report `bit` with a byte-exact restore.
