# Prove a squash-landed branch landed

Blocked by: none

Ownership fence: `internal/git`, `internal/worktree/squash_landed_test.go`
Assumptions: every cleanup path already routes its decision through `LandedInDefault`; none of them needs to change

## What to build

`LandedInDefault` proves a branch landed three ways — ancestry, then merge detection,
then `git cherry` patch containment. A **squash-landing**, where a branch's commits were
composed into one commit on the default branch, passes none of them: it is not an
ancestor, it carries no merge, and no individual patch-id survives the squash.

Every cleanup path asks this one question, so every cleanup path refuses a squash-landed
branch. `bench resume` reports `pruned branches 0`; `bench worktree clean` retains the
branch and mints a recovery ref; the ref itself then cannot be retired. A ten-ticket build
left ten branches, eleven recovery refs and eleven pool entries that no sanctioned command
would remove, and the operator had to hand-write deletion loops. It reads as three
separate leaks and is one blind spot.

Add a fourth proof: **reverse-applicability**. If the branch's cumulative diff against its
merge base can be applied in reverse to the default branch's tree, that content is
demonstrably already present — however it got there. Compute it against a throwaway index
seeded from the default branch, the idiom `TreeHash` already uses, so no working tree is
touched and no real index is disturbed.

## This proof authorizes destruction — build it to say no

A wrong `true` here deletes a branch whose work is not in fact landed, and no other check
stands behind it. Every uncertainty resolves to *not landed*: a diff that fails to
generate, an apply that errors for any reason other than a clean reverse application, a
branch with no merge base, an unreadable index, a binary or rename form the comparison
cannot represent faithfully. The existing three proofs already have this posture — a
merge-carrying branch is deliberately kept rather than guessed at — and the fourth must
match it rather than relax it.

Prefer being wrong in the direction that leaves an orphan over being wrong in the
direction that loses work.

## Acceptance

- [ ] [RL1] A branch whose commits were squashed into one commit on the default branch is proven landed, and the existing three proofs still prove what they proved before.
- [ ] [RL2] A branch carrying any change not present in the default branch is not landed, including one whose changes are a strict superset of what landed.
- [ ] [RL3] Every failure and ambiguity — no merge base, unreadable diff, an apply that does not cleanly reverse — reports not landed rather than an error-shaped true.
- [ ] [RL4] The proof touches no working tree and no real index, and reports the same verdict when run twice against an unchanged repository.
