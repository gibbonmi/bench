# Prove the installed candidate lands its own new fields end to end

Blocked by: 2-install-the-candidate.md, 3-hold-the-refresh.md
Writes: internal/systemtest/candidate_install_test.go (new), internal/systemtest/owner_test.go
Covers: CB15, CB21, CB23, CB24, CB25, CB50

## What to build

Deliver the end-to-end proof of CB-C2 with two real executables built from two trees.

The fixture copies the kit source and adds one new ledger field and one new review-record field in a candidate tree. Its spec declares both fields on a `Persisted fields:` line on the fixture's default branch. The test builds the `main` executable and the candidate, binds the build, and runs the install through the real wrapper. The file carries the `system` build tag. Shared fixture helpers may go in `owner_test.go`, which has 75 lines of headroom under its budget row.

`TestCandidateInstallLandsANewLedgerField` installs, writes the new ledger key through the candidate, lands the build through the unchanged land route, and reads the manifest and the primary executable after the landing. `TestCandidateInstallSurvivesAGapLanding` opens a chunk gap, lands an unrelated change through the broker, and reads the primary digest and `bench status` after that landing.

This ticket adds no production code. A failure here that needs a production change goes back to ticket 2 or ticket 3 as a repair.

## Acceptance

- [ ] After an install whose `doctor --fix` child starts beside a built worktree `dist/bench`, the manifest `path` is the primary executable and its digest is that executable's digest (CB15).
- [ ] An install whose running executable is a worktree `dist/bench` refuses before its build and writes no file (CB21).
- [ ] After the install, `bench status` from the primary checkout exits 0 with a ledger that holds a `gaps` key (CB23).
- [ ] The installed candidate lands a build whose published tree holds a review record with a `folds` key (CB24).
- [ ] After that landing, the primary manifest digest equals the primary `dist/bench` digest (CB25).
- [ ] After an unrelated landing during an open gap, the primary digest is unchanged and `bench status` exits 0 (CB50).

## Checkpoint verification

Run the `system` command of CB-C2's verification plan. The orchestrator then freezes the CB-C2 delta for independent review.
