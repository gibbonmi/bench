# Decide whether an unstaged requirement carries

Blocked by: 01-share-the-plan-command-parser.md, 02-share-the-pass-predicate.md
Writes: internal/reviewrecord/carry/carry.go (new), internal/reviewrecord/carry/carry_test.go (new), internal/reviewrecord/carry/carry_reach_test.go (new), internal/reviewrecord/recordtest/carry.go (new), internal/gate/delta_checkpoint.go
Covers: CF4, CF5, CF6, CF7, CF8, CF9, CF10, CF11, CF12, CF13, CF14, CF15, CF16, CF17, CF18, CF19, CF20, CF21, CF22, CF23, CF24, CF25, CF26, CF27, CF28, CF29, CF31, CF32

## What to build

Add the deep module `internal/reviewrecord/carry`. Its one entry takes the root, the spec, the chunk entry after the freeze, and one owed requirement. It answers the carry candidate, or one of four reasons: `no current passing entry`, `the delta is not test-only`, `the delta changes a whole-only test file`, or `the delta reaches the command`.

The entry applies the five steps of "The carry decision" in the spec, in order:

1. Read the carry candidate of `02-share-the-pass-predicate.md`.
2. Carry when the candidate digest equals the chunk digest.
3. Grade the delta with the test-only file list of the review record.
4. Read the delta reach from the gate.
5. Read the requirement command through the parser and the cover predicate of `01-share-the-plan-command-parser.md`.

In the gate, export the delta reach beside the delta set in the delta checkpoint file. The reach is the package operand of each changed test file, plus the conformance package operand, through the existing `deltaPackages` rule. An empty file list gives only the conformance operand. A changed test file in the system suite or in a race-table package answers whole-only. The delta plan keeps its current behavior. Add no new file to `internal/gate`.

Put new fixture helpers in a new `recordtest` file, because the fixture file is near its line budget. A helper commits a prior entry, commits a delta, and saves a hand-built legacy entry with no source commit.

## Acceptance

- [ ] A test-only delta in `./internal/b` carries `bench test --check system`, `bench test --check system --run TestA`, and `bench test --package ./internal/a`.
- [ ] A delta with an assignment-only plan commit and a test change in `./internal/b` carries `bench test --package ./internal/a`.
- [ ] A candidate whose digest equals the chunk digest carries `bench test --check root-conformance`.
- [ ] Each reach row CF8 to CF21 answers its named reason for its named command and delta.
- [ ] Each candidate row CF23 to CF28 answers `no current passing entry`.
- [ ] A carried probe requirement answers its candidate with the probe outcome, exit, and restore.
- [ ] `TestDeltaCheckpointRunsWholeForTheSystemSuite` and the other delta checkpoint tests pass with no edit.
- [ ] `rg -n "SystemSuite\(|raceTests" internal/reviewrecord` finds no hit.
