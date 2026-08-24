# Gate and publish under the stable owner

Blocked by: 01-pin-the-stable-promotion-owner.md
Writes: internal/landing, internal/gate, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/land_freshness_test.go, internal/systemtest, CHANGELOG.md, projects/benchkit.md

## What to build

Materialize the prospective tree under the stable owner.
Run the baseline gate policy with a gate binary from that exact tree.
Let only the stable owner validate evidence and publish the destination.

## Acceptance

- [ ] SOL09 records a gate executable built from the prospective tree.
- [ ] SOL10 keeps the baseline phase set after a candidate policy omission.
- [ ] SOL11 leaves the destination and marker unchanged after a red gate.
- [ ] SOL12 refuses a destination race at the final compare-and-swap.
- [ ] SOL13 binds evidence reuse to the tree and baseline runner identity.
- [ ] SOL14 resumes each post-publication failure without another publication.
- [ ] SOL15 removes every temporary prospective binary after completion or refusal.
