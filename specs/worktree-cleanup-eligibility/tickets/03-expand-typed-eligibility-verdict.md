# Expand a typed eligibility verdict

Blocked by: 02-characterize-automatic-eligibility-outcomes.md
Writes: internal/worktree/eligibility.go, internal/worktree/subshell.go, internal/worktree/classifier.go, internal/worktree/eligibility_test.go, specs/worktree-cleanup-eligibility/tickets/03-expand-typed-eligibility-verdict.md

## What to build

Introduce the in-process eligibility owner beside the unchanged plan contract.
It accepts typed ownership and safety facts, decides the current explicit outcome
once, and returns a verdict carrying the decision plus all evidence that projection
and later consumers require. `PlanExplicitWithOptions` gathers facts and only
projects that verdict, retaining the current fingerprint inputs and the
derived-after `DiscardBranch` boundary. The two characterization matrices stay
independent of the new policy implementation.

## Acceptance

- [ ] EV1: one eligibility call returns the decision and typed evidence, and explicit projection does not make a second decision.
- [ ] EV2: eligibility action and refusal ordering lives only in the eligibility module before execution.
- [ ] EV3: fingerprints bind every removal-relevant evidence byte and verdict projection, including typed landedness.
