# Extract landing policy decisions

Blocked by: 02-add-explicit-effect-inputs.md
Writes: internal/worktree/, internal/worktree/landingpolicy/ (new)

## What to build

Extract landing decisions into a pure child package with typed source,
destination, resume, publication, release, and residue facts. The parent adapter
continues to own Git, filesystem, process, and rendering effects.

Move combinatorial landing partitions to pure tables. Retain representative
public journeys and focused fact-adapter coverage.

## Acceptance

- [ ] LP1: Typed landing tables cover publish, refusal, interruption, resume, and residue decisions.
- [ ] FA1: Focused real-Git adapter tests prove landing fact translation.
- [ ] RJ1: Every named Git-supplied landing journey remains serial and green.
