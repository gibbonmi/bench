# Share the pass predicate of a chunk verification entry

Blocked by: none
Writes: internal/reviewrecord/check.go, internal/reviewrecord/coverage.go, internal/reviewrecord/record_test.go, internal/reviewrecord/export_test.go, internal/reviewrecord/carry_candidate_test.go (new)
Covers: CF30

## What to build

The chunk verification grade, `checkVerification`, holds the only rule for a passing entry. Extract the part of that rule that does not read the plan owner into one exported pass predicate. The predicate holds when the state is `completed`, the outcome is `pass`, and the exit code is present and zero. The command must equal the requirement command. For a probe requirement, the probe mutation must equal the planned probe, and `ProbeBites` must report a bite and a passed restore. The predicate reads no performer and no role.

`checkVerification` then reads this predicate, and it keeps its performer and role test, its refusal messages, and their order. Also export the carry candidate. The candidate is the last entry of one requirement in one chunk entry, when that entry passes the predicate, holds a source commit, and is current under `EntryCurrent`. The carry decision of `03-decide-the-carry.md` consumes the candidate.

## Acceptance

- [ ] For each entry variant under the plan performer, the carry candidate and the chunk verification grade give one verdict. The variants are a pass, a failed outcome, a missing exit code, a changed command, a probe that did not bite, and a probe whose restore failed.
- [ ] `TestProbeChecksAgree` and each checkpoint test that names `verification <id> is incomplete, failed, or stale` pass with no change to the expected text.
- [ ] `carry_candidate_test.go` holds a case where the candidate refuses an entry with no source commit. It also holds a case where an earlier entry passed and the last entry failed, and the candidate refuses.

The candidate reads `EntryCurrent`, which queries Git. An external test can build its repository with `recordtest` and reach `checkVerification` through `export_test.go`. An in-package test can build it with `internal/gittest` instead.
