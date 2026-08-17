# Contract superseded eligibility policy

Blocked by: 05-migrate-landed-set-planning.md, 07-migrate-release-and-landing.md
Writes: internal/worktree/eligibility.go, internal/worktree/subshell.go, internal/worktree/classifier.go, internal/worktree/clean_landed.go, internal/worktree/resume.go, internal/worktree/ownership.go, internal/worktree/worktree.go, internal/worktree/land.go, internal/worktree/landed.go, internal/worktree/list.go, internal/worktree/lifecycle.go, docs/adr/0005-worktree-cleanup-requires-verifiable-ownership.md, specs/worktree-cleanup-eligibility/tickets/08-contract-eligibility-policy.md

## What to build

Remove every superseded producer-side eligibility decision after the migrated
consumers are green, leaving the verdict module as the indispensable owner of
ordered eligibility policy. Rewrite ADR 0005 as resulting-state documentation:
one eligibility verdict, still fail-closed on the conjunctive marker,
assignment, lock, landedness, recovery, and preservation requirements. Do not
alter pre-existing test logic except a mechanical rename if one is unavoidable.

## Acceptance

- [ ] OS1: the pre-existing suite passes with unchanged test logic, apart from any mechanical rename.
- [ ] OS2: ADR 0005 names the eligibility verdict and retains verifiable conjunctive ownership and preservation requirements.
- [ ] OS3: deleting the eligibility module would require ordered policy to return to explicit, automatic, and landed-set consumers.
