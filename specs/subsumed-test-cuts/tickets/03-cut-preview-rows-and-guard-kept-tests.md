# Cut the preview rows and guard the kept tests

Blocked by: none
Writes: internal/adopt/lifecycle_preview_test.go
Covers: TC20, TC21, TC22, TC23, TC24, TC41, TC42, TC44

## What to build

The Upgrade preview test loses seven e2e rows that its pure `plan kinds` table subsumes. The two review rows whose probes stayed silent keep their tests, and their guard probes prove it. The spec rows T1-12, T1-8, and T1-13 own this ticket.

- In `TestLifecycleUpgradePreviewCompatibility`, remove the rows `owned inline drift`, `owned adapter drift`, `unowned file`, `unowned inline`, `unowned adapter`, `dropped row`, and `preserved CLAUDE.md`.
- Keep the rows `converged`, `owned file drift`, `hook refresh`, and `hook refresh at equal version`, and keep the `plan kinds` subtest.
- Remove each helper that only a removed row used.
- Change no line of `internal/systemtest/adoption_test.go` or of the Link partition tables. Run the guard probes `tc41` and `tc42` on those tests unchanged.

The build changes no production file. Run each probe through `bench probe` with `--expect` on the kept test, and read the `restored` cell of each answer.

## Acceptance

- [ ] The seven rows no longer exist, and the four kept rows and `plan kinds` remain (TC44, review-owned).
- [ ] The completion-plan probes `tc20`, `tc21`, `tc22`, and `tc23` read `bit` with the selection `^TestLifecycleUpgradePreviewCompatibility$/^plan_kinds$`.
- [ ] The completion-plan probe `tc24` reads `bit` on the kept `hook refresh` row.
- [ ] The completion-plan probe `tc41` reads `bit` on `TestDetectedProjectGateRejectsIgnoredDeclaredInput` through `bench test --check system`.
- [ ] The completion-plan probe `tc42` reads `bit` on the `dangling leaf collision` row of `TestLifecycleLinkPartitions`.
- [ ] `bench test --package ./internal/adopt` passes.
