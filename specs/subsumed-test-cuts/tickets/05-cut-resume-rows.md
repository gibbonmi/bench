# Cut the subsumed resume rows

Blocked by: none
Writes: internal/worktree/land_resume_test.go, internal/worktree/land_resume_refusal_test.go
Covers: TC27, TC28, TC34

## What to build

Two resume tests each lose one table row that a kept test subsumes. No top-level test goes, so the count pin does not change. The spec rows WT-4 and WT-10 own this ticket.

- In `TestResumeLandCommandAcceptsSpecSlugAndPath`, remove the `x` row (WT-4). The `./specs/x/spec.md` row stays, and it is the only row that exercises the path form. `TestResumeLandCommandCompletesAnInterruptedMarker` is the kept test for the slug form.
- In `TestResumeLandCommandPublicRefusesDestructiveDestinationState`, stop only the `PL30/clean` cell (WT-10). The table crosses two journeys with the cases, so skip that one pair in the loop and keep both tables whole. `PL25/clean` stays as the control, and every other `PL30` row stays. `TestResumeLandCommandPublicBindsPublishedLandingIdentity` is the kept test.

The build changes no production file. Run each probe through `bench probe` with `--expect` on the kept test, and read the `restored` cell of each answer.

## Acceptance

- [ ] The `x` row and the `PL30/clean` row no longer run, and every other row of both tests still runs.
- [ ] The completion-plan probe `tc27` reads `bit` on `TestResumeLandCommandCompletesAnInterruptedMarker`.
- [ ] The completion-plan probe `tc28` reads `bit` on the kept path row of `TestResumeLandCommandAcceptsSpecSlugAndPath`.
- [ ] The completion-plan probe `tc34` reads `bit` on `TestResumeLandCommandPublicBindsPublishedLandingIdentity`.
- [ ] `bench test --package ./internal/worktree --run 'TestResumeLandCommand'` passes.
