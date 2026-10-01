# Move the landing effect and resume tests onto the verb runner

Blocked by: 6-name-landing-fixtures.md
Writes: internal/worktree/land_effects_cleanup_test.go, internal/worktree/land_effects_test.go, internal/worktree/land_empty_sibling_test.go, internal/worktree/land_resume_refusal_test.go, internal/worktree/land_resume_test.go, internal/worktree/land_reauthorization_test.go, internal/worktree/land_census_test.go, internal/worktree/land_census_output_test.go, internal/worktree/land_trace_test.go, internal/worktree/land_journey_test.go
Covers: VR26

## What to build

Move every `land` and `land-resume` verb call in the listed files onto the verb runner, together with each joins form call. Each stubbed landing passes its joins value through the runner, so each effect stub still runs. `resumeLandArgs` builds an argument list and stays.

Keep every test name and every assertion.

## Acceptance

- [ ] The verb form `rg` in the spec's Further notes prints no line for the listed files.
- [ ] `TestPackageTestCountPin` passes with no change to `worktreeTestCount`, and the serial ceiling holds.
- [ ] No over-budget test file grows past its base line count.
