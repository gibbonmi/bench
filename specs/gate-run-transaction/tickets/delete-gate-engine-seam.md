# Delete the gateEngine seam and its dead wrappers

Blocked by: extract-gate-run-transaction.md
Writes: internal/gate (production files; test edits limited to mechanical renames)

## What to build

The contract half of story 2. Remove `gateEngine`, `productionGateEngine`,
`engineEvaluation`, `newEngineEvaluationAtKit`, `executeWithEngine`,
`executeWithEngineAfterAcquire`, `executeSubjectWithEngine`, and
`newWorkingTreeEvaluation`; `durableReplaceWithEngine`, `inspectWithEngine`,
`operationalWithEngine`, and `persistInterruptedIfGreen` lose their engine
parameter and call the clock, filesystem, lock, and subject functions directly
(collapsing into their engine-less twins where one exists). The run module
the extract ticket landed stays the sole transaction owner — this ticket removes
the seam around it and moves no orchestration back into the entry points.
Behavior is unchanged; the exit test grades it.

## Acceptance

- [ ] None of the eight identifiers occurs anywhere in `internal/gate` (covers GT2)
- [ ] The run module remains the only production file referencing the six transaction markers named in the extract ticket
- [ ] Every characterized outcome and consumer-visible gate behavior is unchanged with test assertions unmodified (renames only)
