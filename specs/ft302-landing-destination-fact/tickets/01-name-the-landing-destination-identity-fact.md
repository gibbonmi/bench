# Name the landing destination identity fact

Blocked by: none
Writes: internal/worktree/land_identity.go, internal/worktree/land_resume.go, internal/worktree/land_facts_test.go

## What to build

One named fact in `internal/worktree/land_identity.go` carries the landing
destination identity: the default branch, the attachment of the checkout to
that branch, the destination commit, and the green marker. The commit resolves
through the one spelling `refs/heads/<branch>^{commit}`. The landing verb reads
the identity fact first, then composes the cleanliness proof, the residue
proof, and the fingerprint proof on top of it. The resume reads the identity
fact alone and never runs the cleanliness proof.

The refusal sentences, the
refusal faces, and the values each verb returns do not change. Every existing
test in `internal/worktree` keeps its logic. One new fact-adapter test in
`internal/worktree/land_facts_test.go` crosses the resume identity: a clean
default-branch checkout translates into the exact commit, branch, and marker.

## Acceptance

- [ ] `internal/worktree/land_identity.go` holds one function that returns the branch, the commit, and the marker, and both the landing verb and the resume call it.
- [ ] `internal/worktree/land_resume.go` resolves the destination through that function alone, with no status read of the checkout.
- [ ] The landing verb resolves the destination through that function and then runs the cleanliness, residue, and fingerprint proofs.
- [ ] The `refs/heads/<branch>^{commit}` spelling is the one commit read for the destination in both verbs.
- [ ] A new `TestLandingResumeDestinationFactAdapterTranslatesCleanCheckout` test in `internal/worktree/land_facts_test.go` proves the resume identity fact on a real clean checkout.
- [ ] `go test ./internal/worktree/...` passes with no change to an existing test body.
