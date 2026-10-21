# Read the explicit plan's locations once

Blocked by: 04-carry-the-run-value-and-read-owner-markers-once.md
Writes: internal/worktree/subshell.go, internal/worktree/classifier.go, internal/worktree/clean_landed.go, internal/worktree/clean_set.go, internal/worktree/worktree.go, internal/worktree/resume.go, internal/worktree/journey_test.go, internal/worktree/parallel_census_test.go
Covers: LGP16, LGP21, LGP22, LGP34

## What to build

An explicit cleanup plan reads its location answers through the run value, so its six location starts become one.

`planExplicitWith` (`internal/worktree/subshell.go:160`) takes the run value as a parameter. Its seven callers pass the value from their `ambient`, or a fresh run value at a public entry:

- `classifier.go:267`;
- `subshell.go:141` and `:155` (`PlanExplicit`, `PlanExplicitWithOptions`);
- `clean_landed.go:145`;
- `clean_set.go:191`;
- `worktree.go:344`;
- `resume.go:86`.

Inside the plan, `canonicalRoot`, `Worktrees`, `CommonDir`, and the nested owner-marker validation read through the run value. The two ledger reads through `intent.Address` stay live, and so do the default-branch reads.

`classifier.go` has 463 lines and `clean_landed.go` has 400, and neither has a grant. Change only their call lines, with no net line gain. `subshell.go` has 432 lines against its grant of 497.

Write the LGP16 test first in `journey_test.go`, and record its red count before the change. Call `planExplicitWith` directly, for an owned target under a canonical primary root. Count only the `--git-common-dir` lines. The plan reads 8 today. If the red shows another count, stop and report a wrong spec.

The new test is serial, because `journeyStubGit` binds PATH. Set `worktreeTestCount` to 807 and `worktreeSerialCeiling` to 47, and record that arithmetic.

Record `bench test --full --package ./internal/worktree` before and after the change, for LGP34.

## Acceptance

- [ ] A direct `planExplicitWith` call for an owned target under a canonical primary root starts the common-directory query exactly three times (LGP16).
- [ ] The serial set equals 47, and the test count equals 807 (LGP21, LGP22).
- [ ] The probe "read the common directory again in planExplicitWith" bites.
- [ ] The verification record holds the package time before and after the change (LGP34).
