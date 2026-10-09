# Probe subjects and mutations

Status: staged

Related roadmap: FT168, partial. This spec settles the mutation-input, subject, verdict, and subject-lock faces. The batch, multi-package, canary, and failing-name faces stay on FT168.

Decision source: named reviewed artifact `roadmap/FT168.md` at aabd462d, with the `/bench-write-spec FT168` invocation of 2026-10-09.

Verification log: pending — the independent review round has not run.

## Problem

An author who proves a test with `bench probe` meets six gaps.
The `--with` replacement is a scalar, so a multiline replacement needs a hand edit.
The probe refuses `--check prose` and `--check system`, so a Markdown or system-row probe returns to copy, edit, rebuild, and restore by hand.
A swap that leaves an import unused gives an `invalid` verdict, and the verdict does not name the compile failure.
A mutant that breaks the private run-binary build reads as a refusal, not as a build failure.
Two probes can mutate one file at the same time, and a bite by an unrelated test reads as `bit`.
Sources: `internal/probe/probe.go:121`, `internal/probe/command.go:19`, `internal/testreport/command.go:152`, `docs/known-issues.md:13`, and the FT168 occurrences.

## Solution

The probe takes five new mutation forms.
`--with-file` reads a multiline replacement, and `--insert-before` puts text before one anchor.
`--occurrence` selects one of several matches, and `--delete` removes the subject for the run.
Repeated `--swap` and `--omit` edits make one combined mutant.

The probe accepts the prose check and the system suite as targets through the registered check owner.
The verdict row names the first failing test, and `--expect` requires a named witness.
A mutant that does not compile reads `build-failed` and names its first build line.
A second probe on the same subject refuses while the first runs.

## User stories

Line: opus / high.
Implementation-line reason: the hardest chunk changes the result envelope that the NRO import reads. The edit rules are exact but numerous. The seams and the canned Go runner already exist and give strong tests.
Harder chunks: PS-C.

### Mutation input

1. As an author, I want `--with-file <path>` to read the replacement for `--swap`, so that a swap can put in several lines.
2. As an author, I want a bad `--with-file` input refused before any write, so that it never mutates the tree.
3. As an author, I want `--insert-before <anchor>` to put text before the one anchor match, so that I add a statement without a hand edit.
4. As an author, I want the insertion text from `--with` or `--with-file`, so that one form takes scalar and multiline text.
5. As an author, I want `--occurrence <n>` to select the nth match, so that I can target one of two identical lines.
6. As an author, I want an out-of-range occurrence refused with the match count, so that a wrong index never mutates another site.
7. As an author, I want `--delete` to remove the subject for the mutated run, so that a reader of an absent file can go red.
8. As an author, I want a deleted subject restored with its bytes and its mode, so that the tree stays as it was.
9. As an author, I want an emptying omission refused with the `--delete` route, so that an empty file never stands in for an absent one.
10. As an author, I want repeated `--swap` and `--omit` edits to make one mutant, so that I can remove an import with its last use.
11. As an author, I want each combined edit to match once in the start bytes, so that a combined mutant stays exact.
12. As an author, I want overlapping combined edits refused before any write, so that one edit never rewrites another.
13. As an author, I want the mutation cell to name each combined form in subject order, so that the row states the mutant.
14. As an author, I want a single-edit flag in a combined mutation refused as usage, so that each form keeps one meaning.

### Subjects

15. As an author, I want `--check prose` accepted as a probe target, so that a Markdown mutation can show a prose red.
16. As an author, I want the prose findings of the mutated run to count as failed tests, so that a prose bite reads `bit`.
17. As an author, I want `--check system` accepted as a probe target, so that a system-row probe needs no hand edit or rebuild.
18. As an author, I want the mutated system run to start after the mutant write, so that the system suite sees the mutant.
19. As an author, I want the probe help to name no hand route, so that the help states only the current forms.
20. As a maintainer, I want each named-check probe to run through the `testreport` check owner, so that its evidence stays producible there.

### Verdict

21. As an author, I want the verdict row to name the first failing test, so that I can cite the red test at once.
22. As an author, I want the first-failure cell null when no test failed, so that no value looks like a false test name.
23. As an author, I want `--expect <test>` to give `bit` only when that test failed, so that another test is not my witness.
24. As an author, I want a run without the expected red to read `invalid` with `expect-missed`, so that the wrong red proves nothing.
25. As an author, I want `--expect` refused with the prose check, so that a test name that cannot occur never misses in silence.
26. As an author, I want a broken private run-binary build to read `build-failed`, so that a compile failure never reads as a refusal.
27. As an author, I want a `mutant` table to name the first build line, so that I see an unused import at once.
28. As a consumer, I want the schema token to change with the verdict-row fields, so that an older reader refuses by name.
29. As a consumer, I want `probe.ReadResult` to read the new verdict-row fields, so that native verification import continues to work.

### Subject lock

30. As an author, I want a second probe on one subject refused while the first runs, so that two probes never mutate one file.
31. As an author, I want the subject-lock refusal before any write or run child, so that a refused probe changes nothing.
32. As an author, I want probes on two different subjects to run at the same time, so that the lock serializes only one file.
33. As an author, I want the lock released when the probe ends, an interrupt included, so that a later probe can run.
34. As a maintainer, I want a completed probe to leave no lock entry under the Bench home, so that the probe writes no record.

### Grammar and guidance

35. As a reader, I want the help inventory's probe line to come from the probe grammar, so that the two lines cannot drift.
36. As a reader, I want the reference guide, the glossary, and the delegate skill to name the new forms, so that guidance matches the verb.
37. As a maintainer, I want the unused-import known-issue row removed, so that the issues table lists only open issues.

### Reviewed exclusions

38. As a reviewer, I want multi-package probes and a shared-baseline batch left to a successor, so that a probe keeps one verdict.
39. As a reviewer, I want `bench canary <fixture>` and the `--names` comparison left to successors, so that this spec changes only `bench probe`.
40. As a reviewer, I want the mutated-run exit to stay in NRO's `probe_execution` table, so that one owner reports it.

## Implementation decisions

### Prerequisite

The build starts after the NRO chunk NR-C3 lands on `main`.
NR-C3 adds the `probe_execution` table, `probe.ReadResult`, the `bench-probe-result-v1` schema token, and the shared canned runner in `internal/probe/probetest`.
This spec writes against that shape, and it does not restate it.
The commitment order governs delivery: this staged spec waits for `bench commitment start`.

### Mutation input

The mutation input resolves to an ordered list of edits over the subject's start bytes.
Each edit is a byte range and its replacement bytes.
One owner applies the list, so every form shares the match, overlap, and empty rules.

`--with-file <path>` and `--omit-file <path>` read through one bounded regular-file reader.
That reader is the current `--omit-file` reader, extracted without a change to its refusals.
The refusal text stays `error: probe mutation file unavailable — <path> is <reason>`.

`--insert-before <anchor>` takes its text from `--with` or `--with-file`.
The anchor must match exactly once, unless `--occurrence` selects one match.
The insertion adds the text immediately before the first byte of the selected match.
The mutation cell reads `insert`.

`--occurrence <n>` takes a positive decimal integer.
It applies to one `--swap`, `--omit`, `--unwrap`, or `--insert-before` edit.
Without it, the exact-once rule and its current ambiguous refusal stay unchanged.
With it, the match count must be at least `n`.
A count below `n` refuses with `error: probe mutation occurrence out of range — the old string matches <count> times, --occurrence is <n>`.
A zero, negative, or non-decimal value is a usage refusal with exit 2.

`--delete` removes the subject for the mutated run.
The mutation cell reads `delete`.
The restore writes the preserved copy back through the current atomic replace, and it proves the bytes and the mode.

An omission whose result has zero bytes refuses before any write.
The refusal is `error: probe mutation empties the subject — use --delete to remove <subject> for the run`.
The rule applies to `--omit`, `--omit-file`, and a combined mutant.

A combined mutant takes two or more edits from repeated `--swap` and `--omit` flags.
The i-th `--swap` pairs with the i-th `--with`, and the two counts must agree.
Each old string must match exactly once in the start bytes.
Two edits whose byte ranges overlap refuse with `error: probe mutations overlap — two edits touch the same bytes`.
The mutation cell joins each edit's form with `+` in the order of its start offset, for example `omit+swap`.
A combined mutant refuses `--occurrence`, `--with-file`, `--insert-before`, `--unwrap`, `--omit-file`, and `--delete` with the usage line.

### Subjects

The probe stops its refusal of `--check prose` and `--check system`.
`testreport.HandProbeRoute` retires, and `testreport.ProbeNotes` drops the two hand-route notes.
Each named check runs through `testreport.Execute`, which owns the prose grade, the system suite, and the private run-binary selection.

A prose probe grades the whole tree twice, once for the baseline and once for the mutant.
The prose runner already counts each finding as a failed test, so a finding over the mutant reads `bit`.

A system probe runs the gate's system suite with the caller's `--run` pattern.
Each `testreport.Execute` call selects its own run binary.
Without an inherited `BENCH_RUN_BINARY`, the mutated run builds that binary from the mutated tree.
The normal shell and `bench worktree exec` set no inherited value.

### Verdict

The verdict row gains the cell `first_failure` after `failed_tests`.
The cell holds the test name of the first failed test in the report's package-then-test order.
It reads null when no named test failed, which includes a build failure, a prose run, and every refusal.
The schema token changes to `bench-probe-result-v2`, and `probe.ReadResult` reads the new field set.

`testreport.Outcome` gains the failed test names in the order that the failures table prints them.
`failed_tests`, `first_failure`, and `--expect` all read that one list.

`--expect <test>` names one Go test.
A failed test matches when its name equals the value or starts with the value and `/`.
A failed mutated run without a match reads verdict `invalid`, cause `expect-missed`, and exit 1.
A passed mutated run still reads `silent`.
`--expect` with `--check prose` is a usage refusal with exit 2.

`runbinary` reports a build-script failure as a typed build error that carries the builder output.
`testreport` maps that error to the outcome kind `build-failed` and keeps its current refusal text.
The `bench test` output and exit stay unchanged.

When the baseline passed and the mutated run reads `build-failed`, a `mutant[1]{first_build_line}` table follows `probe_execution`.
For a Go test build failure, the cell is the first diagnostic line that the failures table prints for that package.
For a run-binary build failure, the cell is the first builder output line that contains `.go:`, or else the first non-empty line.
The table sits inside the counted suffix, so `trailing_bytes` includes it.

### Subject lock

The probe takes an exclusive, non-blocking record lock before it preserves the subject.
The lock file is `$BENCH_HOME/probe/locks/<sha256 of the subject's absolute path>`.
After the lock succeeds, the probe confirms that the path still names the locked file, and it retries once when it does not.
A held lock refuses with `error: probe subject in use — another probe holds <subject>`.
The probe removes the lock file and then releases the lock when it ends.
The refusal order is: usage, repository, subject, mutation input, mutation match, selection, gate lock, subject lock, then baseline.

### Grammar and guidance

`internal/probe` exports its usage line as `probe.Usage`.
The help inventory derives the probe suffix from it, as `cmd/bench/test_command.go` derives the `bench test` suffix.
The reference guide, `CONTEXT.md`, the delegate skill, and `CHANGELOG.md` name the new forms.
`docs/known-issues.md` loses the unused-import row.

### Trust chain

The trust chain is unchanged.
The probe starts no new executable kind.
The focused-run owner selects the run binary and launches `go` as `bench test` does today.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| PS-A / sliced after review | The five new mutation forms and the derived help line | PS1-PS19, PS44 | planned probe edit tests, existing omit-file and refusal tests, help inventory test | no |
| PS-B / sliced after review | Prose and system subjects through the check owner | PS20-PS26 | planned subject-check tests, existing selection-facts tests | no |
| PS-C / sliced after review | First failure, expected witness, self-naming build failure, schema v2 | PS27-PS37 | planned verdict and mutant-build tests, NRO result tests, testreport and runbinary tests | yes |
| PS-D / sliced after review | Subject lock | PS38-PS43 | planned lock tests, existing no-record test | no |
| PS-E / sliced after review | Guidance and known issues | PS45-PS47 | review-owned reads, prose gate | no |

PS-A comes first, because every later chunk reads its grammar.
PS-B, PS-C, and PS-D each write `internal/probe/probe.go`, so they run in series.
PS-E comes last, because its guidance names every form.

## Testing decisions

A good test drives `probe.Command` with a real fixture repository and the canned Go runner.
It observes the stdout rows, the exit, the subject bytes, and the Bench home listing.
The prior art is `internal/probe/probe_test.go` (`newFixture`), `internal/probe/omit_file_test.go`, and the NR-C3 canned runner.
Rows that derive an expected row build it through `toon.TableTyped`, as `probeRow` does.

A `testreport` test drives `Execute` through its current fixtures for the failed-name list and the build-failed mapping.
A `runbinary` test drives `Factory.Own` with a failing `Build` seam for the typed build error.
The gate observes every row through `bench test` of `./internal/probe/...`, `./internal/testreport/...`, `./internal/runbinary/...`, and `./cmd/bench/...`.

### Seam diagram

    trigger: an author or a coordinator runs bench probe <file> <mutation> <selection>
        │
        ▼
    argv, start bytes  ──▶  [ internal/probe: edits, subject lock, preserve, testreport.Execute twice, restore ]  ──▶  probe, selection, probe_execution, mutant, report
                      ◀ tests attach here: probe.Command over a fixture repository with the canned Go runner

    trigger: a focused run selects its run binary
        │
        ▼
    source root  ──▶  [ runbinary.Factory.Own → testreport outcome kind ]  ──▶  build-failed with builder output
                      ◀ tests attach here: Factory with a failing Build seam, and testreport.Execute

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PS1 | 1 | A swap with `--with-file` replaces the match with the file's exact multiline bytes and reads `bit` | planned internal/probe/edit_test.go (TestWithFileSwapsMultilineBytes) | A reader that trims or drops the second line changes the mutant bytes that the stub records. |
| PS2 | 2 | `--with-file` over an absent path refuses before any write | planned internal/probe/edit_test.go (TestWithFileRefusesBadInputs) | A reader that skips the check mutates the subject with empty text. |
| PS3 | 2 | `--with-file` over a symlink, a directory, a FIFO, an oversized file, and an empty file each refuse with its reason | planned internal/probe/edit_test.go (TestWithFileRefusesBadInputs) | A second reader with different refusals drifts from the `--omit-file` reader. |
| PS4 | 3 | `--insert-before` puts the text immediately before the one anchor match | planned internal/probe/edit_test.go (TestInsertBeforePutsTextBeforeTheAnchor) | An insertion after the anchor or a replacement of the anchor gives other mutant bytes. |
| PS5 | 3 | `--insert-before` over an anchor that matches twice refuses as ambiguous | planned internal/probe/edit_test.go (TestInsertBeforePutsTextBeforeTheAnchor) | An insertion at the first match mutates a site that the caller did not name. |
| PS6 | 4 | `--insert-before` takes its text from `--with-file` | planned internal/probe/edit_test.go (TestInsertBeforePutsTextBeforeTheAnchor) | A form that reads only `--with` cannot insert several lines. |
| PS7 | 5 | `--occurrence 2` swaps only the second of two identical matches | planned internal/probe/edit_test.go (TestOccurrenceSelectsOneMatch) | An index that counts from zero or picks the first match changes the other line. |
| PS8 | 6 | `--occurrence 3` over two matches refuses with the out-of-range line that names the count 2 | planned internal/probe/edit_test.go (TestOccurrenceSelectsOneMatch) | A clamp to the last match mutates a site that the caller did not name. |
| PS9 | 6 | `--occurrence 0` and `--occurrence x` each print the usage line with exit 2 | planned internal/probe/edit_test.go (TestOccurrenceSelectsOneMatch) | A parser that accepts zero selects no match and reports a miss instead. |
| PS10 | 7 | `--delete` leaves the subject absent while the mutated run runs | planned internal/probe/delete_test.go (TestDeleteRemovesTheSubjectForTheRun) | The stub hook finds the file when the form only empties it. |
| PS11 | 8 | After `--delete`, the subject holds its start bytes and its start mode | planned internal/probe/delete_test.go (TestDeleteRemovesTheSubjectForTheRun) | A restore that recreates the file at the default mode loses an executable bit. |
| PS12 | 8 | A failed restore after `--delete` reads `restore-failed` and names the preserved copy | planned internal/probe/delete_test.go (TestDeleteReportsARestoreFailure) | A deletion path outside the preserve owner leaves no copy to restore by hand. |
| PS13 | 9 | `--omit-file` with the whole subject refuses with the empties line that names `--delete` | planned internal/probe/delete_test.go (TestOmissionThatEmptiesTheSubjectRefuses) | The FT290 omission empties the file and the red cannot occur. |
| PS14 | 10 | Two `--swap` pairs and one `--omit` make one mutant whose bytes hold all three edits | planned internal/probe/combined_test.go (TestCombinedEditsMakeOneMutant) | A form that applies only the first edit leaves the import in use. |
| PS15 | 11 | A combined edit whose old string matches twice refuses as ambiguous | planned internal/probe/combined_test.go (TestCombinedEditsMatchOnce) | A combined path that skips the count mutates an unnamed site. |
| PS16 | 11 | Each combined old string matches against the start bytes, not against the result of an earlier edit | planned internal/probe/combined_test.go (TestCombinedEditsMatchOnce) | A sequential apply lets a later old string match text that an earlier edit added. |
| PS17 | 12 | Two edits with overlapping byte ranges refuse with the overlap line before any write | planned internal/probe/combined_test.go (TestCombinedEditsRefuseOverlap) | An apply in argv order rewrites the bytes of the first edit. |
| PS18 | 13 | The mutation cell of an omit before a swap in the subject reads `omit+swap` | planned internal/probe/combined_test.go (TestCombinedEditsMakeOneMutant) | A cell in argv order or a fixed word hides what the mutant was. |
| PS19 | 14 | A combined mutant with `--occurrence`, `--with-file`, `--insert-before`, `--unwrap`, `--omit-file`, or `--delete` prints the usage line with exit 2 | planned internal/probe/combined_test.go (TestCombinedEditsRefuseSingleForms) | A silent pair by position gives one form two meanings. |
| PS20 | 15 | A Markdown mutation under `--check prose` that adds a long sentence reads `bit` | planned internal/probe/subject_check_test.go (TestProbeTakesTheProseCheck) | The current refusal answers before the baseline. |
| PS21 | 16 | The prose bite's `failed_tests` cell equals the count of prose findings over the mutant | planned internal/probe/subject_check_test.go (TestProbeTakesTheProseCheck) | A probe that counts zero for a non-Go run reads `invalid`. |
| PS22 | 17 | `--check system --run <pattern>` passes the baseline and runs the mutant through the system-suite argv | planned internal/probe/subject_check_test.go (TestProbeTakesTheSystemCheck) | The current refusal answers before the baseline. |
| PS23 | 18 | The system run's stub hook reads the mutant bytes from the subject | planned internal/probe/subject_check_test.go (TestProbeTakesTheSystemCheck) | A run that starts before the write tests the unmutated tree. |
| PS24 | 19 | The probe help names no hand route and no not-a-target sentence | `internal/probe/refusal_test.go` (`TestProbeHelpSpellings`) | A surviving note still sends authors to the copy route. |
| PS25 | 19 | `testreport.ProbeNotes` holds no copy-aside route and no not-a-target sentence | `internal/testreport/selection_facts_test.go` (`TestProbeNotesDeriveTheirTokens`) | A kept hand route lets the old copy route return to the probe help. |
| PS26 | 20 | The selection row of a system probe reads form `check` and target `system` | planned internal/probe/subject_check_test.go (TestProbeTakesTheSystemCheck) | A private system path in the probe reports another form. |
| PS27 | 21 | A bite with two failed tests names the first in package-then-test order in `first_failure` | planned internal/probe/verdict_test.go (TestFirstFailureNamesTheFirstFailedTest) | A cell from the last event or from map order names the wrong test. |
| PS28 | 22 | A silent run, a build failure, and a baseline refusal each read `first_failure` null | planned internal/probe/verdict_test.go (TestFirstFailureNamesTheFirstFailedTest) | An empty string or a package name reads as a test name. |
| PS29 | 23 | `--expect TestClampNegative` over a failure of that test reads `bit` | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | A form that ignores the value cannot show that it compares. |
| PS30 | 23 | `--expect TestClampNegative` over a failure of only `TestClampNegative/zero` reads `bit` | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | An exact-only compare misses a parent whose subtest failed. |
| PS31 | 24 | `--expect TestClampPositive` over a failure of only `TestClampNegative` reads `invalid`, `expect-missed`, and exit 1 | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | A bite by another test counts as the witness. |
| PS32 | 25 | `--expect` with `--check prose` prints the usage line with exit 2 | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | A prose run never names a test, so the expectation misses in silence. |
| PS33 | 26 | A failing `Build` seam gives `Factory.Own` an error that `errors.As` reads as the typed build error with the builder output | planned internal/runbinary/runbinary_test.go (TestOwnReportsATypedBuildError) | An untyped wrap cannot be told from a temp-directory failure. |
| PS34 | 26 | `testreport.Execute` over a typed build error returns kind `build-failed` with its current refusal bytes | planned internal/testreport/outcome_test.go (TestExecutableBuildFailureIsBuildFailed) | The current refused kind reads as a refusal at the probe. |
| PS35 | 27 | A mutant that leaves an import unused prints the `mutant` table with the compiler line that names the import | planned internal/probe/verdict_test.go (TestMutantBuildFailureNamesItself) | A verdict without the table hides the cause in the report. |
| PS36 | 28 | The probe result schema cell reads `bench-probe-result-v2` | planned internal/probe/result_test.go (TestProbeResultSchemaVersion) | A v1 token over a changed field set lets an older reader fail on a field error. |
| PS37 | 29 | `probe.ReadResult` reads a v2 envelope and refuses a v1 envelope by its schema | planned internal/probe/result_test.go (TestProbeResultSchemaVersion) | A reader with the old field list refuses every new answer. |
| PS38 | 30 | A probe started while another probe holds the subject lock refuses with the in-use line and exit 1 | planned internal/probe/lock_test.go (TestSubjectLockRefusesASecondProbe) | Two probes mutate one file and one restore certifies the other's damage. |
| PS39 | 31 | The in-use refusal leaves the subject bytes, the Bench home, and the go marker unchanged | planned internal/probe/lock_test.go (TestSubjectLockRefusesASecondProbe) | A lock taken after the baseline starts a run child first. |
| PS40 | 31 | With the gate lock and the subject lock both held, the gate-lock line wins | planned internal/probe/lock_test.go (TestSubjectLockFollowsTheGateLock) | A reversed order shows the subject line under a live gate. |
| PS41 | 32 | A probe on a second subject runs while the first subject's lock is held | planned internal/probe/lock_test.go (TestSubjectLockAllowsAnotherSubject) | A repository-wide lock serializes unrelated probes. |
| PS42 | 33 | After an interrupted probe, a new probe on the same subject reaches its verdict | planned internal/probe/lock_test.go (TestSubjectLockReleasesOnInterrupt) | A lock that survives the interrupt blocks the subject until a manual clean. |
| PS43 | 34 | A completed probe leaves the record surfaces unchanged | `internal/probe/outcome_test.go` (`TestProbeWritesNoRecord`) | A kept lock file adds an entry under the Bench home. |
| PS44 | 35 | The help inventory's probe line names `--unwrap`, `--omit-file`, `--with-file`, `--insert-before`, `--occurrence`, `--delete`, and `--expect` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The hand-copied suffix already omits `--unwrap` and `--omit-file`. |
| PS45 | 36 | The reference guide names each new form, the subject lock, and the `mutant` table | review-owned: the reviewer reads the probe paragraphs | The prose check grades sentences, not content. |
| PS46 | 36 | `CONTEXT.md` defines **combined mutant** and **expected witness** with Avoid lists | review-owned: the reviewer reads the entries | The prose check grades sentences, not terms. |
| PS47 | 37 | `docs/known-issues.md` holds no row that names an unused probe import | review-owned: the reviewer reads the table | No check reads the issue table. |

Not covered: story 38 — a reviewed exclusion; the Out of scope section prices it.
Not covered: story 39 — a reviewed exclusion; the Out of scope section prices it.
Not covered: story 40 — NRO owns `focused_exit_code`, and this spec adds no row for it.

### Edge inventory

The audience is this repository and every linked repository.
A linked repository has no system suite, so its system probe reads `invalid` with the cause `baseline-refused`.
An absent Bench home directory differs from an empty one: the lock owner creates its directory at mode 0700.

The hostile-input walk covers these classes.

A `--with-file` or subject path with spaces or glob characters stays literal.
A control byte in the subject name keeps the current render refusal.
A control byte in the first build line passes through `sanitize.Preview` before the encoder.
A test name that reads as a number keeps the encoder's quotes, because the expected row derives through `toon.TableTyped`.
A FIFO operand for `--with-file` refuses without a block, through the shared non-blocking open.
A lock path derives from a digest, so no operand byte reaches the lock file name.

**Won't handle** — an inherited `BENCH_RUN_BINARY` under a build-input mutant — the freshness refusal answers it, and the normal shell sets no inherited value.
**Won't handle** — a line-number anchor — `--occurrence` covers the case of two identical lines, and a line number moves with every edit.
**Won't handle** — a combined mutant across two files — the subject lock and the preserve owner hold one file, and two serial probes remain the route.
**Won't handle** — a wall-time cell in the verdict row — the packages table already prints `elapsed_ms` for each package.
**Won't handle** — a first-failure name for a prose run — a finding names a sentence, not a test, and the report lists each finding.
**Won't handle** — a lock that blocks and waits — a waiting probe hides a concurrent author, and the refusal names the holder's subject.

## Ownership fences

- `internal/probe/`
- `internal/testreport/selection_facts.go`
- `internal/testreport/selection_facts_test.go`
- `internal/testreport/outcome.go`
- `internal/testreport/outcome_test.go`
- `internal/testreport/testreport.go`
- `internal/testreport/command.go`
- `internal/runbinary/runbinary.go`
- `internal/runbinary/runbinary_test.go`
- `internal/reviewrecord/recordcmd/`
- `cmd/bench/main.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/record_operations_test.go`
- `.bench/BENCH-reference.md`
- `CONTEXT.md`
- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.agents/skills/bench-craft-delegate/references/delegation-discipline.md`
- `docs/known-issues.md`
- `CHANGELOG.md`
- `specs/probe-subjects-and-mutations/`
- `reviews/probe-subjects-and-mutations.md`

The ticket slice narrows each ticket's `Writes:` line inside these fences.
`internal/reviewrecord/recordcmd/` and `cmd/bench/record_operations_test.go` hold the NRO import tests that read the schema token.

## Out of scope

- Several packages in one probe, and a batch of independent mutants that share one baseline: estimated 10 edits and 4 gate runs. The batch needs an FT318 decision on how a record imports several verdicts. Single-mutant probes remain the surviving caller.
- `bench canary <fixture>` runs one canary fixture: estimated 6 edits and 3 gate runs. `bench test --check <name>` remains the surviving caller.
- A `--names` form that compares failing test names with a base set: estimated 5 edits and 3 gate runs. The failures table remains the surviving caller.
- A repeatable `--package` on `bench test` (IDEAS, 2026-10-09): it belongs with the multi-package successor. One call per package remains the surviving caller.
- A probe that prints the mutated text and reruns the green after the restore (IDEAS, 2026-10-09): it is not in the decision source. The verdict row and a `bench test` call remain the surviving callers.

## Further notes

Flagged additions:

- The emptying-omission refusal reads "replaces the emptying `--omit-file`" as a refusal that names `--delete`.
- The `mutant` table is a new block after `probe_execution`, not a new verdict-row cell.
- The help inventory derives its probe suffix from `probe.Usage`, because the hand copy has drifted.
- The removal of the unused-import known-issue row follows from the self-naming build failure and the combined form.
- The cause `expect-missed` is a new cause value under the existing verdict `invalid`.

Boundaries:

- NRO owns `probe_execution`, `focused_exit_code`, `command_exit_code`, `trailing_bytes`, and `probe.ReadResult`. NRO's spec keeps the probe and selection field sets byte-identical for its build. This spec changes the verdict-row field set after NRO lands, and it moves the schema token to v2.
- FT318 owns the consumption of the probe result by the verification record. This spec changes no record form.
- FT98 owns preservation support. Each new form restores through the current preserve owner, and the deliberate-red verdict stays.

Response bound: a biting answer already spills past the 10-line bound after NR-C3, and an open learning asks the reviewer for that decision. This spec adds no line to a biting or silent answer. The `mutant` table adds two lines only to a build-failed answer.

Known issues: the unlanded probe-cost row on the `ft168-probe-row` worktree names "FT168: one probe takes several mutations in one private build" as its retire condition. The combined mutant does not retire that row, because the cost face is the batch successor.

Source-clause mapping:

- "multiline replacement text for `--with` from a file" -> PS1-PS3.
- "Markdown files and registered system checks" -> PS20-PS23.
- "Named fixture and class evidence must remain producible through the registered check owner" -> PS26.
- "builds through the worktree build" -> PS22, PS23.
- "One probe takes several mutations" -> PS14-PS19.
- "`--insert-before` for an insertion" -> PS4-PS6.
- "`--expect` for the expected red" -> PS29-PS32.
- "One probe grades several packages in one private build" -> Out of scope.
- "`bench canary <fixture>`" -> Out of scope.
- "A `--names` form" -> Out of scope.
- "Each new face must restore the original subject" -> PS11, PS12, and the current restore tests.
- "The verdict prints the mutated run's exit" -> NRO.
- "a deletion form replaces the emptying `--omit-file`" -> PS10-PS13.
- "a mutant compile error names itself" -> PS33-PS35.
- "A second probe on the same file refuses" -> PS38-PS43.
- "The verdict row also names the first failing test" -> PS27, PS28.
- "A swap accepts a line or occurrence anchor" -> PS7-PS9.

The pre-review proof checklist follows.

- Cited symbols: `probe.Command`, `mutationInput`, `mutate`, `resolveSubject`, `preserve`, `restore`, `gradeCheckTarget`, `gradeGateLock`, `testreport.Execute`, `testreport.Outcome`, `testreport.HandProbeRoute`, `testreport.ProbeNotes`, `runbinary.Factory.Own`, `gocache.RecordLock`, and `testHelpSuffix` resolve at aabd462d. `probe.ReadResult` and `probetest` resolve on the NRO build tip.
- Import edges: `internal/testreport` already imports `internal/runbinary`. `internal/probe` already imports `internal/testreport`. The lock owner imports `internal/gocache` for `RecordLock`, a new edge to validate with `go list` at build entry.
- Source-row clauses and occurrences: the full FT168 body and its 24 occurrences were read, through the retros in Git history and `capture/learnings.md`. The light-path repair pass retro text was not found.
- Promised field labels: `first_failure`, `mutant`, `first_build_line`, `bench-probe-result-v2`, `expect-missed`, `insert`, and `delete`.
- Changed-function callers: `testreport.Execute` has the callers `gradeBaseline`, `probe`, and `testreport.Command`. `HandProbeRoute` has the readers `gradeCheckTarget` and `ProbeNotes`. `ProbeNotes` has the reader `probeGrammar`. `runbinary.Own` has the reader `productionRunBinaryOwner`.
- Copy survival: PS25 reds a surviving hand route, and PS44 reds a surviving hand-copied help suffix.
- Rendered-shape readers: the strings "not probe targets" and "probe focused run unsupported" occur in `internal/probe/probe.go`, `internal/probe/refusal_test.go`, `internal/testreport/selection_facts.go`, and `internal/testreport/selection_facts_test.go`. The NRO import tests read the schema token.
- Pin operators: the rows compare exact rendered rows, exact subject bytes, and exact modes.
- Entry reads: the lock owner reads `BENCH_HOME` through `benchhome.Dir`, as `preserve` does.
- Derived expectations: expected rows derive through `toon.TableTyped`; no expectation reads the code under test.
- Consolidated rules: the `--omit-file` and `--with-file` readers consolidate into one reader with the current refusal set.
- Quantified obligations: "each new face must restore" covers `--with-file`, `--insert-before`, `--occurrence`, `--delete`, combined mutants, prose, and system.
- Workflow-step writes: none.

Sources not re-read: the light-path repair pass retro, and the learning texts behind six occurrence lines that exist only in `roadmap/FT168.md`.
