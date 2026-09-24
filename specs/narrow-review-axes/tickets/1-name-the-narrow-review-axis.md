# 1. Name the narrow review axis shape

Blocked by: none
Writes: .agents/commands/bench-review-implementation.md, .agents/skills/bench-craft-review/SKILL.md, internal/anchors/registry_ft311_review_dispatch.go, internal/conformance/charge_evidence_guidance_test.go
Covers: none

## What to build

The review phase names the narrow axis shape that FT71 measured. In FT71 each full-retrieval axis used about 337k to 733k tokens. Each fresh narrow axis used about 20k to 45k tokens, and its findings held.

Change the review guidance so that each review round dispatches fresh axis sessions. Each axis does these steps:

- It runs `bench preflight evidence <id> --check-current` once, to bind the artifact to the assignment and the source pair.
- It reads the chunk delta, or the repair delta in a confirming round, through one `git diff` of the frozen pair.
- It reads the spec rows, the tickets, the standards, and the surrounding code with targeted reads.
- It returns a bounded report with one line for each finding.

A confirming round reads only the repair delta, and its charge names the folds to confirm. A resumed axis session keeps every earlier stream in its context, so the phase does not resume an axis for a later round.

The full-retrieval rule ("act only after this session retrieves every required source") stops applying to a review axis. The coordinator keeps the manifest, the metadata page, and the current binding, and it hands the evidence identity to each axis. The build phase keeps its own delivery rule unchanged.

Record in the guidance that the shape is provisional. Each narrow round records what it read and what it found. One full-retrieval control review of the same diff decides whether the shape becomes the permanent rule.

Keep the ASD-STE100 prose rules. Change each anchor or conformance pin that names the old sentence in the same commit, and run the owning check.

## Acceptance

- [ ] `.agents/commands/bench-review-implementation.md` names the narrow axis shape, the fresh session for each round, the repair-delta confirming round, and the bounded return.
- [ ] No review guidance still requires an axis to retrieve every evidence page.
- [ ] The build phase delivery rule in `.agents/commands/bench-implement-spec.md` is unchanged.
- [ ] `bench gate` is green on the worktree.
