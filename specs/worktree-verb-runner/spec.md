# Run every worktree verb in tests through one verb runner

Status: staged

Decision source: `specs/worktree-verb-runner/decisions/worktree-seams.md` (ready compiled map).

Verification log: 3 iteration(s) to accept — the reviewer raised the cap from 2 to 3. Iteration 1 found 11 findings and iteration 2 found 5. Iteration 3 found 4 findings and used the last iteration under the cap. The coordinator confirmed the iteration-3 folds by a read of the diff and of the cited code, with no fourth review, on the reviewer's approval.

## Problem

The worktree tests restate how to call a verb, how to build its fixture, and how
to read its answer. Each test file declares its own output buffers, calls a verb
entry or its joins form directly, and parses the output with its own helper.
Several run wrappers and fingerprint extractors repeat the same job with small
differences. Fixture builders return positional tuples, so a call site reads
`root, _, creation, _, home` and a reader must count positions.

`axitest.DecodeDocument` exists, but few worktree tests use it. Thus a test can
pass on a substring match while the document around it is malformed.

## Solution

One test helper in the worktree package, the verb runner, calls each verb and
returns one verb result. The verb result carries the exit code, both streams, the
rows read through `axitest`, and the fingerprint. Every direct verb call in the
package's tests moves onto the verb runner. The named run wrappers and
fingerprint extractors go. Each positional fixture tuple becomes a named fixture
value, and one helper file declares those values.

The static census then refuses a direct verb call outside the two runner files.
It also refuses a call to a joins form or a runner-private reader there. The top-level
test count holds, the serial set stays at or below its ceiling, and the joins
value does not change. The seam reduction spec follows this spec and replaces
each stub once, on the verb runner.

## User stories

Line: opus / medium.
Implementation-line reason: VR-C1 is the hardest material chunk, because every later ticket consumes the runner contract. VR-C5 adds a census with a known precedent. The source is exact, and the seam shape is known. A weakened assertion in a mechanical migration stays green, so the partial gate coverage moves the cheap row up one tier.
Harder chunks: VR-C1, VR-C5.

### The verb runner

1. As a test author, I want one verb runner for every worktree verb, so that I do not restate how to call a verb.
2. As a test author, I want the exit code and both streams in the verb result, so that I need no buffers of my own.
3. As a test author, I want table rows read through `axitest.DecodeDocument`, so that a row assertion decodes the whole document.
4. As a test author, I want an error from a rows read on non-TOON stdout, so that a decode failure never passes as no rows.
5. As a test author, I want the fingerprint read from a table's `fingerprint` cell, so that I write no regular expression for each verb.
6. As a test author, I want the fingerprint read from a record's `fingerprint=` cell, so that the reset plan has the same reader.
7. As a test author, I want an error for no fingerprint or two different fingerprints, so that a wrong plan cannot pass as a match.
8. As a test author, I want a numeric-looking fingerprint returned as the producer's text, so that a quoted cell needs no trim helper.
9. As a test author, I want must-form readers that fail my test on a reader error, so that no test drops a reader error.
10. As a test author, I want a must-form check that stdout has no fingerprint, so that an error plan stays testable.
11. As a test author, I want to give the runner a joins value, so that a stubbed test runs through the same runner.
12. As a test author, I want to give the verb runner stdin, so that the exec verb runs through the runner.
13. As a test author, I want the exec assignment in the verb result, so that the dispatcher's assignment report stays testable.
14. As the seam reduction author, I want the runner to accept the kit root, the home, and the clock, so that its signature stays.
15. As a test author, I want a kit value to fail until the seam reduction wires it, so that no test trusts an ignored kit.
16. As a test author, I want a clock value to fail until the seam reduction wires it, so that no test trusts an ignored clock.

### The migration

17. As a maintainer, I want every direct verb call in a worktree test on the runner, so that one way to call a verb exists.
18. As a maintainer, I want the named run wrappers deleted, so that no second way to call a verb survives.
19. As a maintainer, I want the fingerprint extractors and inline matches deleted, so that one reader owns the fingerprint.
20. As a maintainer, I want each positional fixture tuple to become a named fixture value, so that a call site names its fixture part.
21. As a reviewer, I want the top-level test count to hold through the migration, so that no test is removed or merged.
22. As a reviewer, I want the serial set to stay at or below its ceiling, so that the verb runner does not cost parallel eligibility.
23. As a reviewer, I want each migrated test to keep its pass or skip outcome, so that no test turns into a silent skip.
24. As a reviewer, I want each migrated test to keep every assertion it made, so that the migration does not weaken the suite.
25. As a reviewer, I want no line-count target for the worktree tests, so that the acceptance stays structural.

### The census refusal

26. As a maintainer, I want the census to report a verb entry call outside the runner files, so that a new direct call turns red.
27. As a maintainer, I want the census to report a joins form call outside the runner files, so that no stubbed test bypasses the runner.
28. As a maintainer, I want the census to report a verb entry used as a value, so that a function table cannot bypass the runner.
29. As a maintainer, I want the census to report a verb call inside a closure or a subtest, so that nesting is no shelter.
30. As a maintainer, I want the census to report a package-level verb reference, so that a table outside a function is no shelter.
31. As a maintainer, I want the census to allow every reference inside the two runner files, so that the runner itself is legal.
32. As a maintainer, I want the census to derive the verb entries from their signatures, so that a new entry needs no list edit.
33. As a maintainer, I want the census to derive each joins form from its verb entry's body, so that a renamed joins form stays covered.
34. As a maintainer, I want the census to report a core reader call outside the runner files, so that no test discards a reader error.
35. As a maintainer, I want the live-tree census to report no direct verb call, so that the migration's end state is pinned.

### Reviewed exclusions

36. As the seam reduction author, I want the joins value and all production files unchanged, so that the seam reduction decides each field once.
37. As a maintainer of another package, I want my package's test harness left unchanged, so that this spec stays in the worktree package.

## Implementation decisions

**The runner files.** Two test files hold the runner. `verb_runner_test.go` holds the verb runner, the typed verb keys, the verb call value, the verb result, the core readers, and the must-form readers. `verb_runner_check_test.go` holds the runner's own tests. Each file stays at or below 400 lines.

**The verb runner.** The runner selects a verb by a typed verb key. Each key names one verb entry, and a verb with a joins form also names that form. A key exists for each verb that a test runs. The `shell` verb has no key, because no test runs its verb entry.

The verb call value carries the root, the Bench home, the kit root, the clock, stdin, joins, and the arguments. The joins value is optional. Without it, the runner calls the verb entry. With it, the runner calls the verb's joins form with that value unchanged. A verb without a joins form refuses a joins value.

The function `checkVerbCall` returns an error for a call with a kit value or a clock value. The runner fails the test with that error. The kit message is exactly `verb runner: a kit value waits for the seam reduction spec`. The clock message is exactly `verb runner: a clock value waits for the seam reduction spec`. The seam reduction spec wires both values and keeps the call value's fields.

**The verb result.** The verb result carries the exit code, stdout, stderr, and the assignment that the exec verb resolved. A verb that returns its output as a string, for example `list`, fills stdout and leaves stderr empty. The core rows reader, `readVerbRows`, decodes stdout with `axitest.DecodeDocument` and returns the rows of one table block, or an error.

The core fingerprint reader, `readVerbFingerprint`, returns one value or an error. It reads the `fingerprint` cell of every row in a decoded table, and all cells must agree. When stdout does not decode, it reads the `fingerprint=` cell of the one record line that carries it.

A plan with no applicable fingerprint does not always omit the cell. A failed explicit set writes the placeholder `unapplicableFingerprint` in each row. A faulted `--discard-branch --unclaimed` set writes an empty value in each row. A stale set refusal echoes the fingerprint that the apply requested.

When the agreed value is empty or equals `unapplicableFingerprint`, the reader returns the no-fingerprint error. This rule applies to the table cell and to the record cell. A no-op reset plan writes a literal `none` in its record, and the reader catches it through the same constant. The reader compares with the constant and does not restate its text.

The two core reader names are pinned. No test file in the package declares a local with either name, so the census reports only a real reader call.

**The must forms.** The verb result has three must-form methods: `mustRows`, `mustFingerprint`, and `mustNoFingerprint`. Each takes a `testing.TB` and fails it on a core reader result that it does not accept. `mustNoFingerprint` accepts only the core reader's no-fingerprint error. Migrated tests use only the must forms. Each core reader and `checkVerbCall` is a top-level function whose last result is an `error`, so the census treats it as runner-private.

**Fixture values.** One test file, `verb_fixture_test.go`, declares every named fixture value type. A fixture builder stays in its present file and returns one named value. A fixture value may carry a method that builds a verb call value from its root, home, and joins value. That method builds a call and runs no verb, so it is not a run wrapper.

A positional fixture tuple is a test-file function that returns at least two values and no `error`. One of the values is a `Creation` or a `[]Creation`, or the function returns three or more values. A function that returns exactly a joins value and a probe is not a fixture tuple, and the seam reduction spec owns it.

**The census refusal.** A new census file, `verb_call_census_test.go`, reuses the parse helpers of `parallel_census_test.go`. A verb entry is an exported function of a non-test file with a parameter named `args` of type `[]string`. A verb's joins form is a function that a verb entry's body calls with `defaultJoins()` as its first argument. A runner-private function is a top-level function of a runner file whose last result is an `error`.

The census reports each identifier that names a verb entry, a joins form, or a runner-private function. It reads every test file except the two runner files. It walks function bodies, closures, and package-level declarations. Each report names the file, the line, the enclosing declaration, and the reported name. The census decides by bare identifier, the same way as the serial helper edge.

**Order.** The runner and its readers land first, alone, as the seam chunk. The reset family and the shared fixture values follow, then the other verb families, then the landing family. The census refusal lands last, after every migration ticket, as the contract step.

**Constraints that every ticket keeps.** No production file changes. The top-level test count holds, and a ticket that adds a test raises `worktreeTestCount` in the same change. The serial set stays at or below its ceiling of 46. The lane's growth ratchet refuses an over-budget test file that grows, so each ticket keeps such a file at or below its base line count. The over-budget table in Further notes names each such file and its tickets.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| VR-C1 / `1-add-verb-runner.md` | The verb runner, the verb result, and its readers exist with their own tests | VR1, VR2, VR3, VR4, VR5, VR6, VR7, VR8, VR9, VR10, VR11, VR12, VR13, VR14, VR15, VR16, VR17, VR18, VR19, VR20, VR21, VR22, VR59, VR60, VR61 | The runner tests in `verb_runner_check_test.go` | yes |
| VR-C2 / `2-move-reset-family.md`, `3-name-assignment-fixtures.md`, `4-name-pool-and-residue-fixtures.md` | The reset family runs on the runner, and the shared fixtures return named values | VR23, VR24, VR25, VR26, VR27 | The reset family tests, the package suite, and the count pin | no |
| VR-C3 / `5-migrate-cleanup-verbs.md`, `6-migrate-query-and-create-verbs.md`, `7-migrate-merge-and-reauthorize.md` | Every other non-landing verb runs on the verb runner | VR28, VR29, VR30, VR31, VR32, VR33, VR34, VR35 | The family test files and the count pin | no |
| VR-C4 / `8-name-landing-fixtures.md`, `9-migrate-landing-composition.md`, `10-migrate-landing-effects.md` | Every landing verb runs on the verb runner with named landing fixtures | VR36, VR37, VR38, VR39 | The landing test files and the count pin | no |
| VR-C5 / `11-refuse-direct-verb-calls.md` | The census refuses a direct verb call outside the runner files, and the package-wide end state holds | VR40, VR41, VR42, VR43, VR44, VR45, VR46, VR47, VR48, VR49, VR50, VR51, VR52, VR53, VR54, VR55, VR56, VR57, VR58 | The census tests in `verb_call_census_test.go`, the count pin, and the serial ceiling | yes |

After each chunk, freeze its predecessor and current tips for Standards, Spec, and Coverage review. Each chunk review records three things for that chunk:

1. The differential name sets of VR44 and VR45.
2. The mechanical assertion pre-check of VR46. For each migrated test, the review counts the `t.Fatal`, `t.Fatalf`, `t.Error`, and `t.Errorf` calls at the chunk base and the chunk tip. It names each test whose count fell and accounts for each drop as a must-form replacement.
3. The side-by-side assertion comparison of VR46. The comparison logs each accepted drop by name. The VR-C4 review logs the any-64-hex check of the folded-sibling landing in `land_effects_cleanup_test.go` as one accepted drop.

Ticket 11 owns those rows because it is the last ticket that touches the package.

## Testing decisions

- A good test drives a real verb through the verb runner over a real git fixture and examines the verb result. The runner's own tests compare the result with a direct call inside a runner file.
- The reader tests feed the core readers a real verb output where one exists. A synthetic table derives its text through the same `toon.Table` call that the producer makes. The must forms call only `Helper` and `Fatalf` on their `testing.TB`, and each returns right after its `Fatalf` call. The must-form tests pass a test recorder that embeds a nil `testing.TB`. The recorder overrides `Helper`, `Fatal`, `Fatalf`, and `FailNow`. Its `Fatalf` records the failure and returns.
- The census tests plant synthetic file sets with `plantTestFiles`, which is the precedent of `parallel_census_test.go`. One live-tree test runs the census over the package.
- The gate observes the feature through the worktree package's Go tests in the test phase. No new gate check is added.

### Seam diagram

    trigger: a worktree test
        │
        ▼
    verb call value  ──▶  [ verb runner: verb entry or joins form ]  ──▶  verb result ──▶ must-form readers
                              ◀ tests attach here: compare with a direct call; feed the core readers

    trigger: go test ./internal/worktree
        │
        ▼
    test files + source files  ──▶  [ verb call census ]  ──▶  report lines
                              ◀ tests attach here: synthetic file sets and the live tree

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| VR1 | 1 | Each key whose verb has a usage grammar returns that verb's usage refusal for a usage-refusing argument list: `create`, `release`, `clean`, `reclaim`, `reauthorize`, `merge`, `reset`, `land`, `land-resume`, `list`, `path`, `show`, `build`, and `exec` | planned `TestVerbRunnerKeyReachesItsOwnVerb` in internal/worktree/verb_runner_check_test.go | A key wired to another verb's entry returns another refusal. The expectation reads the usage constant, or the grammar's `Cmd` field for `show` and `build` |
| VR2 | 1 | The `pool` key with one root argument returns the output of `poolAt` for the call's home and that root, followed by a newline | planned `TestVerbRunnerPoolKeyReturnsThePoolPath` in internal/worktree/verb_runner_check_test.go | A pool key wired to another entry returns other text or a usage line |
| VR3 | 1 | The `lease-file` key with no argument returns the stdout and exit code of a direct `LeaseFileCommand` call with no argument | planned `TestVerbRunnerLeaseFileKeyMatchesItsEntry` in internal/worktree/verb_runner_check_test.go | A lease-file key wired to another entry returns other text, while the expectation comes from the verb itself |
| VR4 | 1 | The `resume-clean` key with one unknown argument returns the streams and exit code of a direct `ResumeCleanCommand` call with that argument | planned `TestVerbRunnerResumeCleanKeyMatchesItsEntry` in internal/worktree/verb_runner_check_test.go | A resume-clean key wired to another entry returns another refusal, while the expectation comes from the verb itself |
| VR5 | 2 | The verb result's exit code, stdout, and stderr equal those of a direct call to the same verb entry with the same input | planned `TestVerbRunnerReturnsBothStreamsAndTheExitCode` in internal/worktree/verb_runner_check_test.go | A runner that drops stderr or returns a fixed exit code differs from the direct call |
| VR6 | 3 | The core rows reader returns the decoded rows of the named table block of a real `list` output | planned `TestVerbResultRowsDecodeTheWholeDocument` in internal/worktree/verb_runner_check_test.go | A reader that returns the wrong block or no rows differs from the rows that `axitest.DecodeDocument` returns |
| VR7 | 4 | The core rows reader returns an error for stdout that holds a table followed by a line outside the TOON grammar | planned `TestVerbResultRowsRefuseAPartialDocument` in internal/worktree/verb_runner_check_test.go | A substring reader or a reader that returns an empty set on a decode failure returns no error |
| VR8 | 9 | `mustRows` fails a test recorder when the core rows reader returns an error | planned `TestVerbResultMustRowsFailsOnAReaderError` in internal/worktree/verb_runner_check_test.go | A must form that ignores the error leaves the recorder without a failure |
| VR9 | 5 | The core fingerprint reader returns the shared `fingerprint` cell of a real `clean --landed` plan | planned `TestVerbResultFingerprintReadsTheTableCell` in internal/worktree/verb_runner_check_test.go | A reader that reads another cell or the first 64-hex run of stdout returns a different value |
| VR10 | 6 | The core fingerprint reader returns the `fingerprint=` cell of a real `reset --to` plan record | planned `TestVerbResultFingerprintReadsTheRecordCell` in internal/worktree/verb_runner_check_test.go | A table-only reader returns an error for the record, and a reader that takes the `next=` cell returns more text |
| VR11 | 7 | The core fingerprint reader returns an error for stdout that carries no fingerprint | planned `TestVerbResultFingerprintRefusesAnAbsentValue` in internal/worktree/verb_runner_check_test.go | A reader that returns the empty string with no error passes a plan that has no fingerprint |
| VR12 | 7 | The core fingerprint reader returns an error for a table whose rows carry two different fingerprints | planned `TestVerbResultFingerprintRefusesConflictingCells` in internal/worktree/verb_runner_check_test.go | A reader that returns the first cell passes a set whose rows disagree |
| VR13 | 9 | `mustFingerprint` fails a test recorder when the core fingerprint reader returns an error | planned `TestVerbResultMustFingerprintFailsOnAReaderError` in internal/worktree/verb_runner_check_test.go | A must form that ignores the error leaves the recorder without a failure |
| VR14 | 8 | The core fingerprint reader returns a numeric-looking value exactly as the producer wrote it before `toon.Table` quoted it | planned `TestVerbResultFingerprintKeepsANumericLookingCell` in internal/worktree/verb_runner_check_test.go | A reader that keeps the quotes or decodes a number returns different text |
| VR15 | 10 | `mustNoFingerprint` leaves a test recorder unfailed for a real error plan whose rows carry the `none` placeholder | planned `TestVerbResultMustNoFingerprintAcceptsAnErrorPlan` in internal/worktree/verb_runner_check_test.go | A must form that treats every reader error as a failure fails the recorder |
| VR16 | 10 | `mustNoFingerprint` fails a test recorder for a real plan that carries a fingerprint | planned `TestVerbResultMustNoFingerprintRefusesAPlan` in internal/worktree/verb_runner_check_test.go | A must form that never fails leaves the recorder without a failure |
| VR17 | 11 | With a joins value, the runner calls the verb's joins form with that value, and the value's stub runs | planned `TestVerbRunnerPassesTheJoinsValue` in internal/worktree/verb_runner_check_test.go | A runner that calls the verb entry runs the default instead, and the stub's probe stays at zero |
| VR18 | 12 | Stdin that the call carries reaches the exec child | planned `TestVerbRunnerFeedsStdinToExec` in internal/worktree/verb_runner_check_test.go | A runner that passes a nil reader gives the child an empty input, and its echo differs |
| VR19 | 13 | The exec verb result carries the assignment ID that the exec target resolved to | planned `TestVerbRunnerReturnsTheExecAssignment` in internal/worktree/verb_runner_check_test.go | A runner that calls `ExecCommand` or drops the second result returns an empty assignment |
| VR20 | 14 | The verb call value declares a kit root field and a clock field beside the root and the home | review-owned: the Standards axis reads the call value in the verb runner file | A call value without the two fields forces the seam reduction spec to change the runner's signature |
| VR21 | 15 | `checkVerbCall` returns an error whose text is exactly `verb runner: a kit value waits for the seam reduction spec` for a call with a kit value | planned `TestVerbCallRefusesAKitValue` in internal/worktree/verb_runner_check_test.go | A check that ignores the kit value returns no error, and another text fails the comparison |
| VR22 | 16 | `checkVerbCall` returns an error whose text is exactly `verb runner: a clock value waits for the seam reduction spec` for a call with a clock value | planned `TestVerbCallRefusesAClockValue` in internal/worktree/verb_runner_check_test.go | A check that ignores the clock value returns no error, and another text fails the comparison |
| VR23 | 18, 19 | No reset-family run wrapper or fingerprint extractor exists after ticket 2 | review-owned: the VR23 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR24 | 20 | `restoreFixture` returns one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan in Further notes omits `restoreFixture` | A positional `restoreFixture` appears in the scan output |
| VR25 | 17 | The ticket 2 reset files call no verb form directly | review-owned: the verb form command over the ticket 2 files prints no line | A surviving direct call prints its file and line |
| VR26 | 20 | `newOwnedAssignment`, `newPendingAssignment`, and `newOwnedSubmoduleAssignment` each return one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits the three builders | A positional builder appears in the scan output |
| VR27 | 20 | `unprovableLandedAssignment`, `newResidueGuardFixture`, `newReclaimPool`, and `poolRootFixture` each return one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits the four builders | A positional builder appears in the scan output |
| VR28 | 18, 19 | No cleanup-family run wrapper, fingerprint extractor, test-side row reader, or inline fingerprint match exists in the ticket 5 files | review-owned: the VR28 command in Further notes prints no line | A surviving declaration or inline match prints its file and line |
| VR29 | 20 | `landedSetFixture`, `retainedMemberFixture`, `removableSetFixture`, and `refusedUnlandedRelease` each return one named fixture value | review-owned: the tuple scan omits the four builders | A positional builder appears in the scan output |
| VR30 | 17 | The ticket 5 files call no verb form directly | review-owned: the verb form command over the ticket 5 files prints no line | A surviving direct call prints its file and line |
| VR31 | 18, 19 | No `runCreate`, no `execAtOwnedTarget`, and no inline fingerprint match exists in the ticket 6 files | review-owned: the VR31 command in Further notes prints no line | A surviving declaration or inline match prints its file and line |
| VR32 | 17 | The ticket 6 files call no verb form directly | review-owned: the verb form command over the ticket 6 files prints no line | A surviving direct call prints its file and line |
| VR33 | 18 | No merge run wrapper exists after ticket 7 | review-owned: the VR33 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR34 | 20 | `mergeFixture` and `reauthorizeFixture` each return one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits both builders | A positional builder appears in the scan output |
| VR35 | 17 | The ticket 7 files call no verb form directly, except the one `LandCommand` call in `delegated_integration_test.go` that ticket 9 moves | review-owned: the verb form command over the ticket 7 files prints only that `LandCommand` line | A surviving direct call prints its file and line |
| VR36 | 20 | Each landing fixture builder returns one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits every landing builder in the Enumerations list | A positional builder appears in the scan output |
| VR37 | 18 | No landing run wrapper exists after ticket 9 | review-owned: the VR37 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR38 | 17 | The ticket 9 files call no verb form directly | review-owned: the verb form command over the ticket 9 files prints no line | A surviving direct call prints its file and line |
| VR39 | 17, 19 | The ticket 10 files call no verb form directly and hold no inline fingerprint match | review-owned: the VR39 command in Further notes prints no line | A surviving direct call or inline match prints its file and line |
| VR40 | 20 | No test-file function in the package returns a positional fixture tuple | review-owned: the tuple scan in Further notes prints no line | A positional builder that any ticket missed appears in the scan output |
| VR41 | 18, 19 | No run wrapper, fingerprint extractor, test-side row reader, or inline fingerprint match exists in the package outside the runner files | review-owned: the VR41 command in Further notes prints no line | A surviving wrapper, extractor, or match prints its file and line |
| VR42 | 21 | The package's top-level test count equals `worktreeTestCount`, which rises only by the tests this spec adds | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | A removed or merged test drops the count below the pin |
| VR43 | 22 | The serial set stays at or below the ceiling of 46 | `internal/worktree/parallel_census_test.go` (`TestSerialSetStaysBelowTheCeiling`) | A runner that binds the process environment adds every caller to the serial set |
| VR44 | 23 | The PASS name set of the package's fresh test run equals the chunk base's set plus the added tests | review-owned: each chunk review compares the `go test -count=1 -json ./internal/worktree` name sets | A migration that makes a test skip moves its name out of the PASS set |
| VR45 | 23 | The SKIP name set of the package's fresh test run equals the chunk base's set | review-owned: each chunk review compares the `go test -count=1 -json ./internal/worktree` name sets | A migration that adds a capability skip adds a name to the SKIP set |
| VR46 | 24 | Each migrated test keeps each assertion on the exit code, the streams, and the repository state that it made at the chunk base | review-owned: each chunk review runs the failure-call count pre-check, then the Spec axis compares each migrated test with its base form and logs each accepted drop | A migration that drops an assertion stays green, so only the count drop and the side-by-side read catch it |
| VR47 | 17 | No test outside the runner files declares an output buffer pair for a verb call | review-owned: the Coverage axis reads each chunk diff at each removed verb call | A test that keeps its buffers and calls the runner restates part of the call |
| VR48 | 26 | The census reports a test-file call to a verb entry outside the runner files with the file, the line, the enclosing declaration, and the entry | planned `TestVerbCallCensusReportsAnEntryCall` in internal/worktree/verb_call_census_test.go | A census that skips verb entries returns no report |
| VR49 | 27 | The census reports a test-file call to a joins form outside the runner files | planned `TestVerbCallCensusReportsAJoinsFormCall` in internal/worktree/verb_call_census_test.go | A census that knows only exported entries returns no report for the joins form |
| VR50 | 28 | The census reports a verb entry that a test function stores as a value without a call | planned `TestVerbCallCensusReportsAnEntryUsedAsAValue` in internal/worktree/verb_call_census_test.go | A census that reads only call expressions returns no report for the value |
| VR51 | 29 | The census reports a verb call inside a subtest closure | planned `TestVerbCallCensusReportsACallInsideASubtest` in internal/worktree/verb_call_census_test.go | A census that stops at a function literal returns no report |
| VR52 | 30 | The census reports a verb reference in a test-file package-level variable declaration | planned `TestVerbCallCensusReportsAPackageLevelReference` in internal/worktree/verb_call_census_test.go | A census that walks only function declarations returns no report |
| VR53 | 31 | The census reports nothing for references in `verb_runner_test.go` and `verb_runner_check_test.go` | planned `TestVerbCallCensusAllowsTheRunnerFiles` in internal/worktree/verb_call_census_test.go | A census without the runner exemption reports the runner's own calls |
| VR54 | 32 | The census reports a call to an exported function that has an `args []string` parameter | planned `TestVerbCallCensusDerivesEntriesFromTheSignature` in internal/worktree/verb_call_census_test.go | A census with a fixed name list returns no report for a new entry in the synthetic set |
| VR55 | 32 | The census reports nothing for a call to an exported function without an `args []string` parameter | planned `TestVerbCallCensusIgnoresAnExportedHelper` in internal/worktree/verb_call_census_test.go | A census that treats every exported function as an entry reports the helper |
| VR56 | 33 | The census reports a call to the function that a synthetic verb entry calls with `defaultJoins()` first | planned `TestVerbCallCensusDerivesTheJoinsForm` in internal/worktree/verb_call_census_test.go | A census with a fixed joins form list returns no report for the synthetic joins form |
| VR57 | 34 | The census reports a call outside the runner files to a synthetic runner-file function whose last result is an `error` | planned `TestVerbCallCensusReportsACoreReaderCall` in internal/worktree/verb_call_census_test.go | A census that knows only verb forms returns no report, so a test can discard a reader error |
| VR58 | 35 | The census reports no line on the live package tree | planned `TestVerbCallCensusOnTheLiveTree` in internal/worktree/verb_call_census_test.go | A surviving direct verb call or core reader call on the live tree is reported with its file and line |
| VR59 | 7 | The core fingerprint reader returns the no-fingerprint error for a real explicit-set `clean` error plan whose rows carry the `none` placeholder, and for a real faulted `clean --discard-branch --unclaimed` plan whose rows carry an empty value | planned `TestVerbResultFingerprintTreatsAPlaceholderAsAbsent` in internal/worktree/verb_runner_check_test.go | A reader that returns the agreed cell gives `none` or the empty string with no error, so an error plan passes as a fingerprint |
| VR60 | 7 | The core fingerprint reader returns the no-fingerprint error for a real no-op `reset` plan whose record carries the `none` placeholder | planned `TestVerbResultFingerprintReadsANoOpRecordAsAbsent` in internal/worktree/verb_runner_check_test.go | A reader that returns the first record cell gives `none` with no error, so a no-op reset passes as a fingerprint |
| VR61 | 7 | The core fingerprint reader returns an error that is not the no-fingerprint error for two record lines with different fingerprints | planned `TestVerbResultFingerprintRefusesConflictingRecords` in internal/worktree/verb_runner_check_test.go | A reader that returns the first record cell passes a plan whose records disagree |

Not covered: story 25 — the source rules out a line-count target, so no row grades a line count.
Not covered: story 36 — no ticket plans a production edit, and the bound registry paths sit in the fence for the preflight binding only.
Not covered: story 37 — no ticket plans an edit outside `internal/worktree` test files.

### Edge inventory

The census walks a build-tagged test file and skips a special file, through the shared `parseGoFiles` walk. It walks a closure, a subtest, a method of a test type, and a package-level declaration. An absent runner file leaves no file exempt for that name, so the census reports every remaining call. These behaviors serve this repository only, because the worktree package is kit source.

The hostile-input checklist applies to the readers. The numeric-looking cell class is VR14. The control-byte, path, and patch-header classes do not apply, because the runner renders no output and reads no path from git output.

- Won't handle: the three `subshellAt` calls in `subshell_test.go` — that form takes a shell path and an environment, and the census still guards `Subshell`.
- Won't handle: a local identifier that shadows a reported name — the census reports it by bare identifier, and the author renames the local.
- Won't handle: a verb entry whose argument parameter is not named `args` — every current verb entry uses `args`, and review catches a rename.
- Won't handle: `RunTreeChild` — its parameter is `argv`, it serves the tree verbs of other packages, and its worktree callers stay on the exec verb entry.
- Won't handle: direct verb calls in other packages' tests, for example `internal/treetarget` — the map puts them out of scope, and the entries stay exported.
- Won't handle: a 64-hex text outside the two fingerprint cells — `mustNoFingerprint` reads only those cells, and the folded-sibling test keeps its `--apply` check. The VR-C4 review logs this narrower check as an accepted drop under VR46.

## Ownership fences

- `capture/restructure-backlog.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/worktree/admin_readers_test.go`
- `internal/worktree/build_test.go`
- `internal/worktree/classifier_shape_test.go`
- `internal/worktree/clean_branch_test.go`
- `internal/worktree/clean_classes_test.go`
- `internal/worktree/clean_discard_test.go`
- `internal/worktree/clean_discard_transaction_test.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/clean_landed_hostile_test.go`
- `internal/worktree/clean_landed_test.go`
- `internal/worktree/clean_operand_test.go`
- `internal/worktree/clean_set_apply_test.go`
- `internal/worktree/clean_set_command_test.go`
- `internal/worktree/clean_set_outcomes_test.go`
- `internal/worktree/clean_set_refusal_test.go`
- `internal/worktree/clean_set_test.go`
- `internal/worktree/clean_set_wiring_test.go`
- `internal/worktree/clean_unclaimed_test.go`
- `internal/worktree/delegated_integration_test.go`
- `internal/worktree/eligibility_test.go`
- `internal/worktree/exec_pwd_test.go`
- `internal/worktree/exec_test.go`
- `internal/worktree/identifier_operand_test.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/land_bench_home_test.go`
- `internal/worktree/land_census_output_test.go`
- `internal/worktree/land_census_test.go`
- `internal/worktree/land_effects_cleanup_test.go`
- `internal/worktree/land_effects_test.go`
- `internal/worktree/land_empty_sibling_test.go`
- `internal/worktree/land_facts_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_flags_test.go`
- `internal/worktree/land_folded_base_test.go`
- `internal/worktree/land_freshness_test.go`
- `internal/worktree/land_identity_test.go`
- `internal/worktree/land_journey_test.go`
- `internal/worktree/land_local_capture_test.go`
- `internal/worktree/land_prunes_landed_siblings_test.go`
- `internal/worktree/land_reauthorization_test.go`
- `internal/worktree/land_reauthorize_operand_test.go`
- `internal/worktree/land_release_refusal_test.go`
- `internal/worktree/land_resume_refusal_test.go`
- `internal/worktree/land_resume_test.go`
- `internal/worktree/land_spec_amendment_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/land_trace_test.go`
- `internal/worktree/landed_test.go`
- `internal/worktree/lifecycle_acquire_test.go`
- `internal/worktree/lifecycle_facts_test.go`
- `internal/worktree/lifecycle_policy_test.go`
- `internal/worktree/lifecycle_test.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/list_selected_test.go`
- `internal/worktree/live_binary_test.go`
- `internal/worktree/merge_caller_root_test.go`
- `internal/worktree/merge_from_sha_test.go`
- `internal/worktree/merge_test.go`
- `internal/worktree/orphan_render_test.go`
- `internal/worktree/orphan_test.go`
- `internal/worktree/ownership_test.go`
- `internal/worktree/parallel_census_test.go`
- `internal/worktree/path_identifier_test.go`
- `internal/worktree/pool_reclaim_facts_test.go`
- `internal/worktree/pool_reclaim_test.go`
- `internal/worktree/pool_root_test.go`
- `internal/worktree/reauthorize_test.go`
- `internal/worktree/recovery_retry_test.go`
- `internal/worktree/release_inside_test.go`
- `internal/worktree/release_registration_test.go`
- `internal/worktree/request_token_test.go`
- `internal/worktree/reset_apply_test.go`
- `internal/worktree/reset_fingerprint_test.go`
- `internal/worktree/reset_plan_test.go`
- `internal/worktree/reset_refusal_test.go`
- `internal/worktree/reset_repair_test.go`
- `internal/worktree/reset_restore_refusal_test.go`
- `internal/worktree/reset_restore_test.go`
- `internal/worktree/response_spill_test.go`
- `internal/worktree/resume_reconcile_test.go`
- `internal/worktree/resume_test.go`
- `internal/worktree/show_test.go`
- `internal/worktree/subshell_test.go`
- `internal/worktree/unlanded_route_test.go`
- `internal/worktree/verb_call_census_test.go`
- `internal/worktree/verb_fixture_test.go`
- `internal/worktree/verb_runner_check_test.go`
- `internal/worktree/verb_runner_test.go`
- `internal/worktree/worktree_test.go`
- `reviews/worktree-verb-runner.md`

Reviewer disposition: pending the coordinator's review round and the reviewer's sign-off.
The fence is the union of the ticket write lines and the review pickup.
Build preflight binds five registry paths to the worktree package: the command registry file, its test, the help inventory test, and two conformance registry tests. They are in the fence for that binding only, and no ticket plans an edit to them.
The operating guide governs execution-plan changes.

## Ticket graph

| Ticket | Blocked by | Delivered outcome |
| --- | --- | --- |
| [1. Add the verb runner and its readers](tickets/1-add-verb-runner.md) | none | The verb runner, the verb result, the core and must-form readers, and their own tests |
| [2. Move the reset family onto the verb runner](tickets/2-move-reset-family.md) | 1-add-verb-runner.md | Reset tests on the runner with a named `restoreFixture` |
| [3. Name the assignment fixtures](tickets/3-name-assignment-fixtures.md) | 2-move-reset-family.md | The three assignment builders return named values |
| [4. Name the pool and residue fixtures](tickets/4-name-pool-and-residue-fixtures.md) | 3-name-assignment-fixtures.md | Four pool and residue builders return named values |
| [5. Move the cleanup verbs onto the verb runner](tickets/5-migrate-cleanup-verbs.md) | 4-name-pool-and-residue-fixtures.md | Release, clean, resume-clean, reclaim, list, and path calls in the cleanup files on the runner |
| [6. Move the query and create verbs onto the verb runner](tickets/6-migrate-query-and-create-verbs.md) | 5-migrate-cleanup-verbs.md | Create, list, path, show, exec, build, pool, and lease-file tests on the runner |
| [7. Move the merge and reauthorize verbs onto the verb runner](tickets/7-migrate-merge-and-reauthorize.md) | 6-migrate-query-and-create-verbs.md | Merge and reauthorize tests on the runner with named fixtures |
| [8. Name the landing fixtures](tickets/8-name-landing-fixtures.md) | 7-migrate-merge-and-reauthorize.md | The landing fixture builders return named values |
| [9. Move the landing composition tests onto the verb runner](tickets/9-migrate-landing-composition.md) | 8-name-landing-fixtures.md | Landing flag, refusal, and identity tests on the runner |
| [10. Move the landing effect and resume tests onto the verb runner](tickets/10-migrate-landing-effects.md) | 8-name-landing-fixtures.md | Landing effect, resume, census, and journey tests on the runner |
| [11. Refuse a direct verb call outside the verb runner](tickets/11-refuse-direct-verb-calls.md) | 9-migrate-landing-composition.md, 10-migrate-landing-effects.md | The census refusal and its live-tree pin |

## Out of scope

- The seam reduction spec from the same map is a separate capability with its own future spec: approximately 60 edits, 4 gate runs. It keeps 15 joins fields, reads each value once per exported function, and lowers the serial ceiling.
- A census refusal for positional fixture tuples is a separate capability that the map did not decide: approximately 3 edits, 1 gate run.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 1: two specs, the verb runner spec first | This spec; Out of scope names the seam reduction spec |
| Ticket 5: delete the run wrappers and the fingerprint extractors | VR23, VR28, VR31, VR33, VR37, VR39, VR41 |
| Ticket 5: every direct verb call moves onto the verb runner | VR25, VR30, VR32, VR35, VR38, VR39, VR47, VR58 |
| Ticket 5: each positional fixture tuple becomes a named value in one helper file | VR24, VR26, VR27, VR29, VR34, VR36, VR40 |
| Ticket 5: the count pin holds | VR42 |
| Ticket 6: the worktree package only, on `axitest.DecodeDocument` | VR6, VR7, and story 37 |
| Ticket 7: the census reports a test-file call to a verb entry outside the runner file, with file and line | VR48, VR51, VR52, VR53, VR58 |
| Ticket 8: no size target; the acceptance is structural | Story 25 |
| Map destination: the runner accepts the kit root, the Bench home, and the clock as values | VR20, VR21, VR22 |
| Map term: the verb result carries the exit code, the rows, and the fingerprint | VR5, VR6, VR9, VR10 |

### Enumerations

A run wrapper is a test-file helper that calls a verb form and returns its output. The run wrappers are these:

- reset, ticket 2: `runReset`, `runResetWith`
- cleanup, ticket 5: `runCleanup`, `runCleanupWith`, `runDiscard`, `planAndApply`, `runResume`, `runResumeAt`, `mustResumeClean`, `mustReclaim`
- query and create, ticket 6: `runCreate`, `execAtOwnedTarget`
- merge, ticket 7: `runMerge`
- landing, ticket 9: `landIn`

The fingerprint extractors are `resetFingerprint` and `restoreFingerprint` in ticket 2, and `reclaimFingerprint` and `cleanupRowFingerprint` in ticket 5. The test-side row readers are `cleanupRows`, `cleanupRowFields`, `cleanupRowValue`, `cleanupRowsField`, `rowForTarget`, and `unclaimedVerdicts`, all in ticket 5. Each splits rendered rows on commas, so none can read the `[]any` rows that `Document.Rows` returns. `mustRows` replaces them, and `mustFingerprint` replaces the fingerprint-column read of `cleanupRowsField`.

The inline fingerprint matches are ten sites. Each goes to the ticket that writes its file:

| File and site | Ticket | Replacement |
| --- | --- | --- |
| `clean_set_test.go`: the `setFingerprint` variable and its uses in the clean files | 5 | `mustFingerprint`, or `mustNoFingerprint` for an error plan |
| `clean_landed_test.go`: the match inside `cleanupRowFingerprint` | 5 | `mustFingerprint` |
| `clean_landed_test.go`: the shared-fingerprint match | 5 | `mustFingerprint` |
| `clean_landed_test.go`: the apply-fingerprint match | 5 | `mustFingerprint` |
| `pool_reclaim_test.go`: the apply-invocation match | 5 | `mustFingerprint` |
| `orphan_render_test.go`: the aggregate-row match in `planReclaimableCount` | 5 | `mustRows` for the count |
| `clean_operand_test.go`: the plan match | 5 | `mustFingerprint` |
| `identifier_operand_test.go`: the plan match | 6 | `mustFingerprint` |
| `worktree_test.go`: the plan match | 6 | `mustFingerprint` |
| `land_effects_cleanup_test.go`: the no-cleanup-fingerprint check | 10 | `mustNoFingerprint` |

The commands run from the repository root. `<files>` is the list of `internal/worktree` test files on that ticket's `Writes:` line.

- VR23: `rg -n '^func (runReset|runResetWith|resetFingerprint|restoreFingerprint)\(' internal/worktree`
- VR28: `rg -n '^func (runCleanup|runCleanupWith|runDiscard|planAndApply|runResume|runResumeAt|mustResumeClean|mustReclaim|cleanupRowFingerprint|reclaimFingerprint|cleanupRows|cleanupRowFields|cleanupRowValue|cleanupRowsField|rowForTarget|unclaimedVerdicts)\(|\bsetFingerprint\b|\[0-9a-f\]\{64\}' <files>`
- VR31: `rg -n '^func (runCreate|execAtOwnedTarget)\(|\[0-9a-f\]\{64\}' <files>`
- VR33: `rg -n '^func runMerge\(' internal/worktree`
- VR37: `rg -n '^func landIn\(' internal/worktree`
- VR39: the verb form command over `<files>`, then `rg -n '\[0-9a-f\]\{64\}' <files>`
- VR41: `rg -n '^func (runCleanup|runCleanupWith|runDiscard|planAndApply|runReset|runResetWith|runMerge|runCreate|runResume|runResumeAt|landIn|execAtOwnedTarget|mustResumeClean|mustReclaim|resetFingerprint|restoreFingerprint|reclaimFingerprint|cleanupRowFingerprint|cleanupRows|cleanupRowFields|cleanupRowValue|cleanupRowsField|rowForTarget|unclaimedVerdicts)\(|\bsetFingerprint\b|\[0-9a-f\]\{64\}' internal/worktree --glob '*_test.go' --glob '!verb_runner_test.go' --glob '!verb_runner_check_test.go'`
- The verb form command: `rg -n '\b(BuildCommand|ExecCommand|ExecCommandResolving|LandCommand|ResumeLandCommand|ListCommand|MergeCommand|PathCommand|ReclaimCommand|ReauthorizeCommand|ResetCommand|ResumeCleanCommand|ShowCommand|Subshell|CleanCommand|ReleaseCommand|CreateCommand|PoolCommand|LeaseFileCommand|buildWith|landWith|resumeLandWith|mergeWith|reauthorizeWith|resetWith|resumeCleanCommandWith|cleanCommandWith|releaseCommandWith)\(' <files>`. That name list is the census's derived set on the base tree, and the census replaces the command after ticket 11.

A hit inside an expected-output literal that the test compares whole is not a reader. The Coverage axis names each such hit.

The tuple scan is a Go AST scan of the package's test files under the predicate in Implementation decisions. Ticket 2 records the scan program in its review record, and every later review reruns it. The run wrappers also match the predicate and go under VR41. At the decision-source commit, the scan lists these fixture builders:

- reset, ticket 2: `restoreFixture`
- assignments, ticket 3: `newOwnedAssignment`, `newPendingAssignment`, `newOwnedSubmoduleAssignment`
- pool and residue, ticket 4: `unprovableLandedAssignment`, `newResidueGuardFixture`, `newReclaimPool`, `poolRootFixture`
- cleanup, ticket 5: `landedSetFixture`, `retainedMemberFixture`, `removableSetFixture`, `refusedUnlandedRelease`
- merge and reauthorize, ticket 7: `mergeFixture`, `reauthorizeFixture`
- landing, ticket 8: `publicLandingFixture`, `publicLandingFixtureAtHome`, `specLessLandingFixture`, `foldedLandingFixture`, `ticketsOnlyLandingFixture`, `landingFixtureAtHome`
- landing, ticket 8: `redProspectiveGateLanding`, `brokerChangingLanding`, `brokerDestinationFixture`, `landSurface`, `foldLandingSibling`

### Over-budget test files

The default line budget is 400, and `.bench/structure.budgets` grants `worktree_test.go` 533 lines. The table names each test file over its budget at the base and each ticket that writes it. Each listed ticket keeps the file at or below its base line count.

| File | Base lines | Budget | Tickets |
| --- | --- | --- | --- |
| `parallel_census_test.go` | 1125 | 400 | 1, 11 |
| `worktree_test.go` | 1031 | 533 | 3, 6, 7 |
| `merge_test.go` | 932 | 400 | 7 |
| `resume_test.go` | 598 | 400 | 3, 5 |
| `pool_reclaim_test.go` | 546 | 400 | 4, 5 |
| `identity_component_test.go` | 509 | 400 | 3, 8, 9 |
| `lifecycle_test.go` | 436 | 400 | 3, 5 |
| `exec_test.go` | 432 | 400 | 3, 6 |
| `land_journey_test.go` | 431 | 400 | 8, 10 |
| `ownership_test.go` | 421 | 400 | 3 |

Tickets 1 and 11 change only the `worktreeTestCount` value in `parallel_census_test.go`, so its line count stays.

### Reader sweep

- `worktreeTestCount` has one reader, `TestPackageTestCountPin`. Tickets 1 and 11 raise it.
- `worktreeSerialCeiling` has one reader, `TestSerialSetStaysBelowTheCeiling`. No ticket changes it.
- The run wrapper names have readers only in `internal/worktree` test files and in `capture/restructure-backlog.md`, which names `runCreate` as a split anchor. Ticket 6 updates that row.
- The joins forms have no reader outside `internal/worktree`. The verb entries have production readers in `cmd/bench/worktree_leaves.go`, `cmd/bench/main.go`, `cmd/bench/commit_chain.go`, `internal/harness/worktree.go`, and `internal/sessioninspect/sessioninspect.go`, and one test reader in `internal/treetarget/run_test.go`. No reader changes, because no production signature changes.
- The usage lines that VR1 reads come from `internal/usage/worktree.go`. The `show` and `build` grammars refuse with their command name only, so VR1 reads that name from each grammar's `Cmd` field. The `pool`, `lease-file`, and `resume-clean` verbs have no usage constant, so VR2 to VR4 read their expectation from the verb itself. No literal moves into a test.
- The map's `## Sources` name `parallelCensus` as a drift trigger. This spec adds a census file beside it and does not change `parallelCensus`.

### Pre-review proof checklist

- Cited symbols, by file:
  - `internal/axi/axitest/document.go`: `axitest.DecodeDocument`, `Document.Rows`
  - `internal/worktree/parallel_census_test.go`: `parseGoFiles`, `parseTestFiles`, `parseSourceFiles`, `plantTestFiles`, `TestPackageTestCountPin`, `TestSerialSetStaysBelowTheCeiling`
  - `internal/worktree/joins.go`: `defaultJoins`
  - `internal/worktree/exec.go`: `ExecCommandResolving`
  - `internal/worktree/worktree.go`: `poolAt`, `LeaseFileCommand`
  - `internal/worktree/clean_set.go`: `unapplicableFingerprint`
  - `internal/worktree/resume.go`: `ResumeCleanCommand`
  - `internal/toon/toon.go`: `toon.Table`
- Import edges: `internal/worktree` tests already import `internal/axi/axitest` in `list_selected_test.go`, `list_actions_test.go`, `landed_test.go`, `unlanded_route_test.go`, and `path_identifier_test.go`. The package already imports `internal/toon` and `internal/usage`. No new edge.
- Source-row clauses and occurrences: the Source trace table quotes each clause; the map is the only occurrence.
- Promised field labels: `fingerprint` as a table field; `fingerprint=` as a record cell, emitted by `reset.go`. The two pinned refusal messages in Implementation decisions.
- Changed-function callers: no production function changes. The deleted test helpers' callers are the test files in the fence.
- Copy survival: VR41 fails when one run wrapper, one fingerprint extractor, one test-side row reader, or one inline match survives. VR58 fails when one direct verb call or core reader call survives.
- Rendered-shape readers: none; no rendered output shape changes.

### Probes run during authoring

- A probe test fed `axitest.DecodeDocument` four shapes. A `reset_plan{...}` record and a `merged{...}` record decode as a string, not an object, so the decode returns an error. A lane line before a record fails with a missing colon. A `worktree_cleanup[1]{...}:` table decodes with its `fingerprint` cell. Thus the fingerprint reader needs the record branch, and the rows reader returns an error on record output. The probe file was deleted.
- The joins value still holds 29 fields, and `parallel_census_test.go` still pins 664 tests and a serial ceiling of 46. No source file of the map changed between the asset commit `8ab52861` and the base `6d0c1e79`.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"worktree-verb-runner-20261001","orchestrator_session":"claude:session-8e50287c-d1de-4951-ba3d-d6632e597e5d","author_limit":1,"assignments":{"1-add-verb-runner.md":[{"session":"claude:bench-writer/vr-t1-author","assignment":"vr-t1-author","model":"opus","effort":"medium","source":"0c95c9447c20189f3f2155719ef965bffc339856","native_ref":"claude:agent/vr-t1-author-20261001@0c95c9447c20189f3f2155719ef965bffc339856"},{"session":"claude:bench-writer/vr-t1-repair-1","assignment":"vr-t1-repair-1","model":"opus","effort":"medium","source":"179e2eb8f945c2c878471e7d4fe5ac51fc9e01c2","native_ref":"claude:agent/vr-t1-repair-1-20261001@179e2eb8f945c2c878471e7d4fe5ac51fc9e01c2","predecessor":"claude:bench-writer/vr-t1-author","trigger":"user-directed","stopped":"The author returned its final report after record commit 4479b2fa and holds no write.","preserved":"179e2eb8f945c2c878471e7d4fe5ac51fc9e01c2"},{"session":"claude:bench-writer/vr-t1-verify-2","assignment":"vr-t1-verify-2","model":"opus","effort":"medium","source":"b02ef7a3aabe58cb0d68ed8a32d7123f7f9910a3","native_ref":"claude:agent/vr-t1-verify-2-20261001@b02ef7a3aabe58cb0d68ed8a32d7123f7f9910a3","predecessor":"claude:bench-writer/vr-t1-repair-1","trigger":"user-directed","stopped":"The repair session returned its final report after record commit d44d28ce and holds no write.","preserved":"b02ef7a3aabe58cb0d68ed8a32d7123f7f9910a3"}],"2-move-reset-family.md":[{"session":"claude:bench-writer/vr-t2-author","assignment":"vr-t2-author","model":"opus","effort":"medium","source":"2fc7c9c079f47794850a79a3f00b81812584d168","native_ref":"claude:agent/vr-t2-author-20261001@2fc7c9c079f47794850a79a3f00b81812584d168"}],"3-name-assignment-fixtures.md":[{"session":"claude:bench-writer/vr-t3-author","assignment":"vr-t3-author","model":"opus","effort":"medium","source":"933de8625f5ef6bbe2500a6f70384396686862f6","native_ref":"claude:agent/vr-t3-author-20261001@933de8625f5ef6bbe2500a6f70384396686862f6"}],"4-name-pool-and-residue-fixtures.md":[],"5-migrate-cleanup-verbs.md":[],"6-migrate-query-and-create-verbs.md":[],"7-migrate-merge-and-reauthorize.md":[],"8-name-landing-fixtures.md":[],"9-migrate-landing-composition.md":[],"10-migrate-landing-effects.md":[],"11-refuse-direct-verb-calls.md":[]}},"chunks":[{"id":"VR-C1","tickets":["1-add-verb-runner.md"],"verification":[{"id":"1-worktree","command":"bench test --package ./internal/worktree","ticket":"1-add-verb-runner.md"},{"id":"1-probe","command":"bench probe internal/worktree/verb_runner_test.go --swap 'unapplicableFingerprint' --with '\"probe-placeholder\"' --package ./internal/worktree --run TestVerbResultFingerprintTreatsAPlaceholderAsAbsent","probe":"swap","ticket":"1-add-verb-runner.md"}]},{"id":"VR-C2","tickets":["2-move-reset-family.md","3-name-assignment-fixtures.md","4-name-pool-and-residue-fixtures.md"],"verification":[{"id":"2-worktree","command":"bench test --package ./internal/worktree","ticket":"2-move-reset-family.md"},{"id":"3-worktree","command":"bench test --package ./internal/worktree","ticket":"3-name-assignment-fixtures.md"},{"id":"4-worktree","command":"bench test --package ./internal/worktree","ticket":"4-name-pool-and-residue-fixtures.md"}]},{"id":"VR-C3","tickets":["5-migrate-cleanup-verbs.md","6-migrate-query-and-create-verbs.md","7-migrate-merge-and-reauthorize.md"],"verification":[{"id":"5-worktree","command":"bench test --package ./internal/worktree","ticket":"5-migrate-cleanup-verbs.md"},{"id":"6-worktree","command":"bench test --package ./internal/worktree","ticket":"6-migrate-query-and-create-verbs.md"},{"id":"7-worktree","command":"bench test --package ./internal/worktree","ticket":"7-migrate-merge-and-reauthorize.md"}]},{"id":"VR-C4","tickets":["8-name-landing-fixtures.md","9-migrate-landing-composition.md","10-migrate-landing-effects.md"],"verification":[{"id":"8-worktree","command":"bench test --package ./internal/worktree","ticket":"8-name-landing-fixtures.md"},{"id":"9-worktree","command":"bench test --package ./internal/worktree","ticket":"9-migrate-landing-composition.md"},{"id":"10-worktree","command":"bench test --package ./internal/worktree","ticket":"10-migrate-landing-effects.md"}]},{"id":"VR-C5","tickets":["11-refuse-direct-verb-calls.md"],"verification":[{"id":"11-worktree","command":"bench test --package ./internal/worktree","ticket":"11-refuse-direct-verb-calls.md"},{"id":"11-probe","command":"bench probe internal/worktree/verb_call_census_test.go --swap '\"args\"' --with '\"argv\"' --package ./internal/worktree --run TestVerbCallCensusDerivesEntriesFromTheSignature","probe":"swap","ticket":"11-refuse-direct-verb-calls.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/worktree-verb-runner/spec.md"},{"id":"worktree","command":"bench test --package ./internal/worktree"}]}
```

### Contradiction flagged for reviewer veto

Ticket 1 of the map says that the seam reduction spec rewrites each stubbed test once, onto the verb runner. Ticket 5 says that every direct verb call moves onto the verb runner in this spec, and a stubbed test calls a joins form directly. This spec follows ticket 5: a stubbed test moves onto the runner with its joins value now. The seam reduction spec then replaces the stub with a real fixture once, without a second runner rewrite. This reading is non-behavioral.

### Flagged additions

- The census also refuses a joins form, not only a verb entry. Without it, a test can call a joins form with its own buffers, which is the case ticket 7 closes.
- The census also refuses a verb entry used as a value and a package-level reference, so a table cannot bypass the runner.
- The census also refuses a runner-private function outside the runner files, so a migrated test cannot discard a core reader error.
- The runner fails a call that carries a kit value or a clock value until the seam reduction spec wires them.
- The core readers return an error, and the must forms fail the test. The split gives the error paths a test seam and keeps every migrated read fail-closed.
- The runner uses two files, because the runner and its own tests exceed one 400-line file.
- VR44 and VR45 add a differential run per chunk, because the count pin cannot see a test that turns into a skip.

### Map handling

This spec compiled the map: the map, its tickets, and its asset moved from `decisions/` into `specs/worktree-verb-runner/decisions/`. The seam reduction spec reads the same compiled map in place. This spec carries no `Roadmap:` line, because `FT356` stays open until the seam reduction spec lands. One map source names `kitRoot` under `internal/gate/kit_source.go`; the function lives in `internal/gate/phases.go`. That drift touches only the seam reduction spec.

The promote-then-delete retirement of a spec removes its compiled topic folders. The retirement of this spec must therefore keep `specs/worktree-verb-runner/decisions/` until the seam reduction spec stages and moves or reads the map. The reviewer decides that order at this spec's retirement.
