# Migrate release and landing consumers

Blocked by: 06-migrate-apply-under-lock.md
Writes: internal/worktree/ownership.go, internal/worktree/worktree.go, internal/worktree/land.go, internal/worktree/ownership_test.go, internal/worktree/lifecycle_test.go, internal/worktree/land_test.go, specs/worktree-cleanup-eligibility/tickets/07-migrate-release-and-landing.md

## What to build

Route release, first-run landing release, and resumed landing release through the
automatic verdict and the apply-under-lock contract. Keep their existing receipt
identity, release diagnostics, and first-run/resume output exactly intact; this
ticket changes which decision they consume, not their public command behavior.

## Acceptance

- [ ] CO4: `ReleaseCommand` and both first-run and resumed `LandCommand` release only through the automatic verdict while retaining their exact receipts and diagnostics.
