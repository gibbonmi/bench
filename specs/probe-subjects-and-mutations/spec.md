# Probe subjects and mutations

Status: staged

Related roadmap: FT168, partial. This spec settles the mutation-input, subject, verdict, and subject-lock faces. The batch, multi-package, canary, and failing-name faces stay on FT168.

Decision source: named reviewed artifact `roadmap/FT168.md` at aabd462d, with the `/bench-write-spec FT168` invocation and the reviewer decisions D1 and D2 of 2026-10-09.

Verification log: 2 iteration(s) to accept — an opus xhigh spec review found five blocking defects in iteration 1, and the author folded them. Iteration 2 accepted after folds. A sonnet high ticket review accepted after folds in one iteration, and the slicer added the missing Writes paths.

## Problem

An author who proves a test with `bench probe` meets six gaps.
The `--with` replacement is a scalar, so a multiline replacement needs a hand edit.
The probe refuses `--check prose` and `--check system`, so a Markdown or system-row probe returns to copy, edit, rebuild, and restore by hand.
A swap that leaves an import unused gives an `invalid` verdict, and the verdict does not name the compile failure.
A mutant that breaks the private run-binary build reads as a refusal, not as a build failure.
Two probes can mutate one file at the same time, and a bite by an unrelated test reads as `bit`.
Sources: `mutationForm` and `gradeCheckTarget` in `internal/probe/probe.go`, `internal/probe/command.go:19`, `internal/testreport/command.go:157`, the invalid-mutant row of `docs/known-issues.md`, and the FT168 occurrences.

## Solution

The probe takes five new mutation forms.
`--with-file` reads a multiline replacement, and `--insert-before` puts text before one anchor.
`--occurrence` selects one of several matches, and `--delete` removes the subject for the run.
Repeated `--swap` and `--omit` edits make one combined mutant.

The probe accepts the prose check and the system suite as targets through the registered check owner.
The verdict row names the first failing test, and `--expect` requires a named witness.
A mutant that does not compile reads `build-failed` and names its first build line.
A second probe on the same subject refuses while the first runs.

## Reviewer decisions

- D1 (2026-10-09): `first_failure` goes into the verdict row, and the schema token becomes `bench-probe-result-v2`. NRO's "byte-identical" sentence binds only NRO's own build, so it is a flagged change, not a contradiction.
- D2 (2026-10-09): the batch with multi-package, `bench canary <fixture>`, and `--names` split into successor specs, and FT168 stays partial. The batch is FT168's next slice.

## User stories

Line: opus / high.
Implementation-line reason: two chunks are hard. PS-C changes the result envelope that the NRO import reads. PS-D adds a cross-process lock with an identity recheck. The canned Go runner and the gate-holder precedent give strong tests.
Harder chunks: PS-C, PS-D.

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
13. As an author, I want the mutation cell to name each combined form, so that the verdict row states the mutant.
14. As an author, I want a single-edit flag in a combined mutation refused as usage, so that each form keeps one meaning.
15. As an author, I want each form's partner flags enforced as usage, so that no flag is ignored in silence.
16. As an author, I want a zero-match edit to read the current substring miss, so that absent text reads the same under every form.

### Subjects

17. As an author, I want `--check prose` accepted as a probe target, so that a Markdown mutation can show a prose red.
18. As an author, I want the prose findings of the mutated run to count as failed tests, so that a prose bite reads `bit`.
19. As an author, I want `--check system` accepted as a probe target, so that a system-row probe needs no hand edit or rebuild.
20. As an author, I want `--run` accepted with `--check system`, so that a system probe can select one test.
21. As an author, I want `--run` with any other check refused as usage, so that a pattern is never dropped in silence.
22. As an author, I want each focused run to build its own run binary, so that the mutated system run tests the mutant.
23. As an author, I want the probe help to name no hand route and no rebuild step, so that it states current forms.
24. As a maintainer, I want each named-check probe to run through the `testreport` check owner, so that its evidence stays producible there.
25. As an author, I want a substring-miss answer to name the selection that the request parsed, so that the row never guesses.

### Verdict

26. As an author, I want the verdict row to name the first failing test, so that I can cite the red test at once.
27. As an author, I want the first-failure cell null when no test failed, so that no value looks like a false test name.
28. As an author, I want `--expect <test>` to give `bit` only when that test failed, so that another test is not my witness.
29. As an author, I want a run without the expected red to read `invalid` with `expect-missed`, so that the wrong red proves nothing.
30. As an author, I want `--expect` refused with the prose check, so that a test name that cannot occur never misses in silence.
31. As an author, I want a failed run-binary build script to read `build-failed`, so that a compile failure never reads as a refusal.
32. As an author, I want a cancelled build or an absent Go to keep its current kind, so that only a compile failure reads `build-failed`.
33. As a maintainer, I want the build error text unchanged, so that every reader of `Factory.Own` prints the same bytes.
34. As an author, I want a `mutant` table to name the first build line, so that I see an unused import at once.
35. As a consumer, I want the schema token to change with the verdict-row fields, so that an older reader refuses by name.
36. As a consumer, I want `probe.ReadResult` to read the new verdict-row fields, so that native verification import continues to work.

### Subject lock

37. As an author, I want a second probe on one subject refused while the first runs, so that two probes never mutate one file.
38. As an author, I want a subject that changed before the lock refused, so that a probe never restores another probe's mutant.
39. As an author, I want the subject-lock refusal before any write or run child, so that a refused probe changes nothing.
40. As an author, I want probes on two different subjects to run at the same time, so that the lock serializes only one file.
41. As an author, I want the lock released when the probe ends, an interrupt included, so that a later probe can run.
42. As an author, I want a lock file without a live holder to grant the lock, so that a killed probe never blocks its subject.
43. As an author, I want a lock that cannot be set up refused by name before the baseline, so that I see the cause.
44. As a maintainer, I want a completed probe to leave no lock entry under the Bench home, so that the probe writes no record.

### Grammar and guidance

45. As a reader, I want the help inventory's probe line to come from the probe grammar, so that the two lines cannot drift.
46. As a reader, I want the reference guide, the glossary, the delegate skill, and the changelog to name the new forms.
47. As a maintainer, I want the unused-import and unused-variable known-issue row removed, so that the issues table lists only open issues.
48. As a maintainer, I want the probe-cost row to name the batch obligation, so that the combined mutant does not retire it.

### Reviewed exclusions

49. As a reviewer, I want multi-package probes and a shared-baseline batch left to a successor, so that a probe keeps one verdict.
50. As a reviewer, I want `bench canary <fixture>` and the `--names` comparison left to successors, so that this spec changes only `bench probe`.
51. As a reviewer, I want the mutated-run exit to stay in NRO's `probe_execution` table, so that one owner reports it.

## Implementation decisions

### Prerequisite

The build starts after the NRO chunk NR-C3 lands on `main`.
NR-C3 adds the `probe_execution` table, `probe.ReadResult`, the `bench-probe-result-v1` schema token, and the shared canned runner in `internal/probe/probetest`.
This spec writes against that shape, and it does not restate it.
The commitment order governs delivery: this staged spec waits for `bench commitment start`.

### Owner files

`internal/probe/probe.go` has 332 lines on `main`, and three chunks change it.
Each chunk therefore puts its new logic in its own file.
`internal/probe/edit.go` owns the edit list and the form rules.
`internal/probe/verdict.go` owns `first_failure`, `--expect`, and the `mutant` table.
`internal/probe/lock.go` owns the subject lock and the re-read.
`probe.go` keeps the sequence and calls these owners.

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
A match count of zero gives the current substring-miss answer.
A match count from 1 to `n - 1` refuses with `error: probe mutation occurrence out of range — the old string matches <count> times, --occurrence is <n>`.
A zero, negative, or non-decimal value is a usage refusal with exit 2.

`--delete` removes the subject for the mutated run through the `Remove` operation of the injected `durablefile.Replacer`.
The mutation cell reads `delete`.
A failed removal can leave the subject absent, so it takes the recovery of a published replace failure.
The copy stays until a restore proves the start bytes, and a failed restore keeps the copy at exit 2.
The restore writes the preserved copy back through `Replacer.Replace`, as every restore does, and it proves the bytes and the mode.

An omission whose result has zero bytes refuses before any write.
The refusal is `error: probe mutation empties the subject — use --delete to remove <subject> for the run`.
The rule applies to `--omit`, `--omit-file`, and a combined mutant.

A combined mutant takes two or more edits from repeated `--swap` and `--omit` flags.
The i-th `--swap` pairs with the i-th `--with`, and the two counts must agree.
Each pair whose `--with` equals its `--swap` refuses with the current `error: probe mutation empty — --with equals --swap`.
Each old string must match exactly once in the start bytes.
Two edits whose byte ranges overlap refuse with `error: probe mutations overlap — two edits touch the same bytes`.

The combined mutation cell lists each `--swap` edit and then each `--omit` edit, joined by `+`, for example `swap+swap+omit`.
The cell is the same on every answer path, a substring miss included, because it derives from the two flag groups alone.
The cell does not follow argv order: an argv that leads with `--omit` still lists the swaps first.
A zero-match edit in a combined mutant gives the current substring-miss answer with that cell.

The form matrix is closed.
An argv names exactly one mutation group: the swap and omit group, `--unwrap`, `--omit-file`, `--insert-before`, or `--delete`.
Any other count is a usage refusal with exit 2.
So `--insert-before a --with b --swap c` refuses, because it names two groups.
Each of these argv shapes is also a usage refusal with exit 2:

- `--with` together with `--with-file`.
- `--with-file` without `--swap` or `--insert-before`.
- `--insert-before` without `--with` or `--with-file`.
- `--occurrence` with `--omit-file`, `--delete`, or a combined mutant.
- `--delete` with `--with` or `--with-file`.
- A combined mutant with `--with-file`, `--insert-before`, `--unwrap`, `--omit-file`, or `--delete`.

### Selection

`--run` is valid with `--package` or with `--check system`.
`--run` with any other check stays a usage refusal, as the current grammar answers.
The usage line names both `--run` positions.
`selectionArgs` passes `--run` to `testreport.Prepare` for both forms.

The probe parses the selection before it answers a substring miss.
The substring-miss answer then reads its form, target, and run cells from the prepared request.
That request is the one source of those cells, so a check-form miss names the named-check run pattern, not `all`.
An unknown check therefore wins over a substring miss.

The refusal order is: usage, repository, subject, mutation input, mutation match, selection, substring miss, gate lock, subject lock, subject re-read, then baseline.

### Subjects

The probe stops its refusal of `--check prose` and `--check system`.
`testreport.HandProbeRoute` retires.
`testreport.ProbeNotes` keeps two notes: the `--package` expression rule and the `--check` inventory note without the not-a-target clause.
The rebuild note retires, because each focused run without an inherited selection builds its own run binary from the tree.
Each named check runs through `testreport.Execute`, which owns the prose grade, the system suite, and the private run-binary selection.

A prose probe grades the whole tree twice, once for the baseline and once for the mutant.
The prose runner already counts each finding in `FailedTests`, so a finding over the mutant reads `bit`.

A system probe runs the gate's system suite with the caller's `--run` pattern.
Each `testreport.Execute` call selects its own run binary.
Without an inherited `BENCH_RUN_BINARY`, the mutated run builds that binary from the mutated tree.
The normal shell sets no inherited value, and `bench worktree exec` removes it.

### Verdict

The verdict row gains the cell `first_failure` after `failed_tests`.
The cell holds the first entry of the failed-name list.
It reads null when that list is empty, which includes a build failure, a prose run, and every refusal.
The schema token changes to `bench-probe-result-v2`, and `probe.ReadResult` reads the new field set.

`testreport.Outcome` gains `FailedNames`, the failed Go test names in the order that the failures table prints them.
For a Go run, `FailedTests` and `FailedNames` derive from the one `report.failed()` walk.
A prose run keeps `FailedTests` as its finding count and leaves `FailedNames` empty.
`first_failure` and `--expect` read `FailedNames`; `failed_tests` reads `FailedTests`.
`Outcome` stops being comparable with `!=`, so the tests that compare it move to a field-wise helper.

`--expect <test>` names one Go test.
A failed test matches when its name equals the value or starts with the value and `/`.
A failed mutated run without a match reads verdict `invalid`, cause `expect-missed`, and exit 1.
A passed mutated run still reads `silent`.
`--expect` with the prose check is a usage refusal with exit 2.
The probe reads the prose check name from one `testreport` export, and its own `proseCheck` constant retires.

### Build failure

`runbinary` builds a typed `BuildError` only on the build-script exit-failure branch.
That type carries the builder output, and its `Error()` text equals the current `<exit error>: <output>` bytes.
A cancelled build and an absent Go keep their current untyped errors.
`testreport` maps a `BuildError` to the outcome kind `build-failed` and keeps its current refusal text.
The `OutcomeBuildFailed` doc comment in `internal/testreport/outcome.go` then names both sources: a package failure with no failing test, and a failed run-binary build script.
Every other selection failure stays `refused`, and the `bench test` output and exit stay unchanged.

`testreport.Outcome` gains `BuildLine`.
For a Go test build failure, it is the first diagnostic line that the failures table prints for that package.
For a `BuildError`, it is the first builder output line that contains `.go:`, or else the first non-empty line.

When the baseline passed and the mutated run reads `build-failed`, a `mutant[1]{first_build_line}` table follows `probe_execution`.
The table sits inside the counted suffix, so `trailing_bytes` includes it.
A baseline whose run-binary build fails now reads the cause `baseline-build-failed`, not `baseline-refused`.

### Subject lock

The lock directory is `$BENCH_HOME/probe-locks/`, a sibling of the preserve directory `$BENCH_HOME/probe/`.
The lock file is `<sha256 of the subject's absolute path>` in that directory.
The probe takes the lock after the gate-lock read and before the baseline.

The acquire protocol makes at most three attempts:

1. Create the directory at mode 0700 if it is absent.
2. Open the lock file read-write with create, no-follow, and close-on-exec at mode 0600. A write lock needs a writable descriptor. If the open answers ENOENT, start the next attempt.
3. Take an exclusive, non-blocking write lock with `gocache.RecordLock`. If another process holds it, refuse with the in-use line.
4. Compare the held descriptor's file identity with the path's file identity through `fstat` and `lstat`. Open no second descriptor. If they differ, or if `lstat` answers ENOENT, close the descriptor and start the next attempt.

A held lock refuses with `error: probe subject in use — another probe holds <subject>`.
Four faults refuse with `error: probe subject lock unavailable — <reason>`.
They are a directory that cannot be created, a third ENOENT at the open, any other open error, and a third identity mismatch.
ELOOP and EACCES are examples of another open error.
A lock file without a live holder grants the lock, because a record lock dies with its process.

After the lock succeeds, the probe reads the subject again.
The re-read classifies the path through `resolveSubject`, so an absent, symlinked, special, or unreadable subject gives its current `error: probe subject unavailable — <subject> is <reason>` line.
If the bytes or the mode differ from the start read, the probe refuses with `error: probe subject changed — <subject> changed before the subject lock; run the probe again`.

The release protocol removes the lock file while the probe holds the lock.
It then closes the descriptor, which releases the lock.
Last, it removes the lock directory and ignores ENOENT and ENOTEMPTY.
A waiter that opened the old file fails the identity check and retries.
The deferred release also runs on an interrupt and on every refusal after the lock.

### Grammar and guidance

`internal/probe` exports its usage line as `probe.Usage`.
The help inventory derives the probe suffix from it, as `cmd/bench/test_command.go` derives the `bench test` suffix.
The reference guide, `CONTEXT.md`, the delegate skill, and `CHANGELOG.md` name the new forms.
The delegate skill's anchored needle at `.agents/skills/bench-craft-delegate/SKILL.md:81` stays byte-identical, because a central-property self-probe still reports one omission or one swap.
`docs/known-issues.md` loses the unused-import row.
The probe-cost row's retire condition changes to the batch obligation, and FT168's body keeps that obligation.

### Trust chain

The trust chain is unchanged.
The probe starts no new executable kind.
The focused-run owner selects the run binary and launches `go` as `bench test` does today.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| PS-A / 1-add-probe-mutation-forms.md | The five new mutation forms, the closed form matrix, and the derived help line | PS1-PS25, PS76 | planned edit, delete, and combined tests, extended omit-file and refusal tests, help inventory test | no |
| PS-B / 2-probe-prose-and-system-checks.md | Prose and system subjects through the check owner, and the request-derived selection cells | PS26-PS37 | planned subject-check tests, extended check-form and selection-facts tests, a testreport build-count test | no |
| PS-C / 3-name-the-first-failure-and-expected-witness.md, 4-name-a-mutant-build-failure.md | First failure, expected witness, self-naming build failure, and schema v2 | PS38-PS57 | planned verdict, mutant-build, result, testreport, and runbinary tests | yes |
| PS-D / 5-lock-the-probe-subject.md | Subject lock with re-read and identity recheck | PS58-PS68, PS75 | planned lock tests, existing no-record and preservation tests | yes |
| PS-E / 6-document-the-probe-forms.md | Guidance and known issues | PS69-PS74 | review-owned reads, prose gate, anchors | no |

PS-A comes first, because every later chunk reads its grammar.
PS-B, PS-C, and PS-D each call new owners from `internal/probe/probe.go`, so they run in series.
PS-E comes last, because its guidance names every form.

### Completion plan

This version-2 completion plan carries the run identity, the author limit, and one verification set for each ticket.
Each code ticket names one central-property probe, and its commands are future obligations, not executed proof.

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"probe-subjects-build-20261012","orchestrator_session":"claude:ps-orchestrator-20261012","author_limit":1,"assignments":{"1-add-probe-mutation-forms.md":[{"session":"claude:ps_t01","assignment":"e7c800827d5fd2eed7e68741a4ec1720","model":"opus","effort":"high","source":"8209612990d250568855d50b2735518e31d55223","native_ref":"claude-agent:ps_t01"}],"2-probe-prose-and-system-checks.md":[],"3-name-the-first-failure-and-expected-witness.md":[],"4-name-a-mutant-build-failure.md":[],"5-lock-the-probe-subject.md":[],"6-document-the-probe-forms.md":[]}},"chunks":[{"id":"PS-A","tickets":["1-add-probe-mutation-forms.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"1-add-probe-mutation-forms.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"1-add-probe-mutation-forms.md"},{"id":"package-core-guard","command":"bench test --check package-core-guard","ticket":"1-add-probe-mutation-forms.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"1-add-probe-mutation-forms.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"1-add-probe-mutation-forms.md"},{"id":"central-property","command":"bench test --package ./internal/probe/...","probe":"Apply each combined edit over the result of the earlier edit instead of over the start bytes. TestCombinedEditsMatchOnce must turn red.","ticket":"1-add-probe-mutation-forms.md"}]},{"id":"PS-B","tickets":["2-probe-prose-and-system-checks.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"2-probe-prose-and-system-checks.md"},{"id":"testreport","command":"bench test --package ./internal/testreport","ticket":"2-probe-prose-and-system-checks.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"2-probe-prose-and-system-checks.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"2-probe-prose-and-system-checks.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"2-probe-prose-and-system-checks.md"},{"id":"central-property","command":"bench test --package ./internal/probe/...","probe":"Answer the substring miss before the selection parse with a hard-coded all run cell. TestSelectionPrecedesTheSubstringMiss must turn red.","ticket":"2-probe-prose-and-system-checks.md"}]},{"id":"PS-C","tickets":["3-name-the-first-failure-and-expected-witness.md","4-name-a-mutant-build-failure.md"],"verification":[{"id":"t3-probe-family","command":"bench test --package ./internal/probe/...","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-testreport","command":"bench test --package ./internal/testreport","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-runbinary","command":"bench test --package ./internal/runbinary","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-record-forms","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-gate-readers","command":"bench test --package ./internal/gate","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-command-package","command":"bench test --package ./cmd/bench","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-conformance-package","command":"bench test --package ./internal/conformance","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-root-conformance","command":"bench test --check root-conformance","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t3-central-property","command":"bench test --package ./internal/probe/...","probe":"Match --expect by exact name only and drop the parent-and-slash match. TestExpectRequiresTheNamedWitness must turn red.","ticket":"3-name-the-first-failure-and-expected-witness.md"},{"id":"t4-probe-family","command":"bench test --package ./internal/probe/...","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-testreport","command":"bench test --package ./internal/testreport","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-runbinary","command":"bench test --package ./internal/runbinary","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-record-forms","command":"bench test --package ./internal/reviewrecord/recordcmd","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-gate-readers","command":"bench test --package ./internal/gate","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-command-package","command":"bench test --package ./cmd/bench","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-conformance-package","command":"bench test --package ./internal/conformance","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-root-conformance","command":"bench test --check root-conformance","ticket":"4-name-a-mutant-build-failure.md"},{"id":"t4-central-property","command":"bench test --package ./internal/testreport","probe":"Map every builder error to build-failed instead of only a BuildError. TestOnlyAScriptFailureIsBuildFailed must turn red.","ticket":"4-name-a-mutant-build-failure.md"}]},{"id":"PS-D","tickets":["5-lock-the-probe-subject.md"],"verification":[{"id":"probe-family","command":"bench test --package ./internal/probe/...","ticket":"5-lock-the-probe-subject.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"5-lock-the-probe-subject.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"5-lock-the-probe-subject.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"5-lock-the-probe-subject.md"},{"id":"central-property","command":"bench test --package ./internal/probe/...","probe":"Omit the fstat and lstat identity recheck after the record lock. TestSubjectLockRechecksTheIdentity must turn red.","ticket":"5-lock-the-probe-subject.md"}]},{"id":"PS-E","tickets":["6-document-the-probe-forms.md"],"verification":[{"id":"guidance","command":"bench test --check docs-currency-workflow","ticket":"6-document-the-probe-forms.md"},{"id":"conformance-meta","command":"bench test --check conformance-meta","ticket":"6-document-the-probe-forms.md"},{"id":"agent-definitions","command":"bench test --check claude-agent-definitions","ticket":"6-document-the-probe-forms.md"},{"id":"skills-index","command":"bench test --check skills-index-command-adapters","ticket":"6-document-the-probe-forms.md"},{"id":"prose","command":"bench test --check prose","ticket":"6-document-the-probe-forms.md"},{"id":"anchors-package","command":"bench test --package ./internal/anchors","ticket":"6-document-the-probe-forms.md"},{"id":"command-package","command":"bench test --package ./cmd/bench","ticket":"6-document-the-probe-forms.md"},{"id":"conformance-package","command":"bench test --package ./internal/conformance","ticket":"6-document-the-probe-forms.md"},{"id":"root-conformance","command":"bench test --check root-conformance","ticket":"6-document-the-probe-forms.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/probe-subjects-and-mutations/spec.md"},{"id":"probe-family","command":"bench test --package ./internal/probe/..."},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"runbinary","command":"bench test --package ./internal/runbinary"},{"id":"record-forms","command":"bench test --package ./internal/reviewrecord/recordcmd"},{"id":"gate-readers","command":"bench test --package ./internal/gate"},{"id":"command-package","command":"bench test --package ./cmd/bench"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

A good test drives `probe.Command` with a real fixture repository and the canned Go runner.
It observes the stdout rows, the exit, the subject bytes, and the Bench home listing.
The prior art is `internal/probe/probe_test.go` (`newFixture`), `internal/probe/omit_file_test.go`, and the NR-C3 canned runner.
Rows that derive an expected row build it through `toon.TableTyped`, as `probeRow` does.

A lock row holds the lock from another process, as `TestProbeRefusesUnderALiveGateRun` holds the gate lock.
A record lock belongs to its process, so an in-process holder cannot refuse the probe.
Two package hooks in `internal/probe` give the lock rows their windows. The start-read hook runs between the start read and the lock. The lock-owner hook runs between the open and the identity check. The probe runs in the test process, so both hooks reach it.

A `testreport` row drives `Execute` through `installTestSelectionFactory` with a counting or failing builder.
A `BuildError` row drives the canonical builder over a fixture `scripts/go-build.sh`, as `TestBuildSubjectLeavesTheManifestDirectoryToTheScript` does.
The gate observes every row through `bench test` of `./internal/probe/...`, `./internal/testreport/...`, `./internal/runbinary/...`, and `./cmd/bench/...`.

### Seam diagram

    trigger: an author or a coordinator runs bench probe <file> <mutation> <selection>
        │
        ▼
    argv, start bytes  ──▶  [ internal/probe: edits, selection, subject lock, preserve, testreport.Execute twice, restore ]  ──▶  probe, selection, probe_execution, mutant, report
                      ◀ tests attach here: probe.Command over a fixture repository with the canned Go runner and a second-process lock holder

    trigger: a focused run selects its run binary
        │
        ▼
    source root  ──▶  [ runbinary build script → BuildError → testreport outcome kind ]  ──▶  build-failed with BuildLine
                      ◀ tests attach here: a fixture go-build.sh, and testreport.Execute with a test selection factory

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PS1 | 1 | A swap with `--with-file` replaces the match with the file's exact multiline bytes and reads `bit` | planned internal/probe/edit_test.go (TestWithFileSwapsMultilineBytes) | A reader that trims or drops the second line changes the mutant bytes that the stub records. |
| PS2 | 2 | Each bad file input refuses with its reason before any write, for `--omit-file` and for `--with-file` | `internal/probe/omit_file_test.go` (`TestOmitFileRefusesInvalidInputsBeforeRunning`) | A second reader with other refusals drifts from the shared table. |
| PS3 | 3 | `--insert-before` puts the text immediately before the one anchor match | planned internal/probe/edit_test.go (TestInsertBeforePutsTextBeforeTheAnchor) | An insertion after the anchor or a replacement of the anchor gives other mutant bytes. |
| PS4 | 3 | `--insert-before` over an anchor that matches twice refuses as ambiguous | planned internal/probe/edit_test.go (TestInsertBeforePutsTextBeforeTheAnchor) | An insertion at the first match mutates a site that the caller did not name. |
| PS5 | 4 | `--insert-before` takes its text from `--with-file` | planned internal/probe/edit_test.go (TestInsertBeforePutsTextBeforeTheAnchor) | A form that reads only `--with` cannot insert several lines. |
| PS6 | 5 | `--occurrence 2` swaps only the second of two identical matches | planned internal/probe/edit_test.go (TestOccurrenceSelectsOneMatch) | An index that counts from zero or picks the first match changes the other line. |
| PS7 | 6 | `--occurrence 3` over two matches refuses with the out-of-range line that names the count 2 | planned internal/probe/edit_test.go (TestOccurrenceSelectsOneMatch) | A clamp to the last match mutates a site that the caller did not name. |
| PS8 | 6 | `--occurrence 0` and `--occurrence x` each print the usage line with exit 2 | planned internal/probe/edit_test.go (TestOccurrenceSelectsOneMatch) | A parser that accepts zero selects no match and reports a miss instead. |
| PS9 | 7 | `--delete` leaves the subject absent while the mutated run runs | planned internal/probe/delete_test.go (TestDeleteRemovesTheSubjectForTheRun) | The stub hook finds the file when the form only empties it. |
| PS10 | 8 | After `--delete`, the subject holds its start bytes and its start mode | planned internal/probe/delete_test.go (TestDeleteRemovesTheSubjectForTheRun) | A restore that recreates the file at the default mode loses an executable bit. |
| PS11 | 8 | A failed restore after `--delete` reads `restore-failed` and names the preserved copy, and a failed removal whose restore fails keeps the copy at exit 2 | planned internal/probe/delete_test.go (TestDeleteReportsARestoreFailure) | A deletion path outside the preserve owner leaves no copy to restore by hand. |
| PS12 | 9 | `--omit-file` with the whole subject refuses with the empties line that names `--delete` | planned internal/probe/delete_test.go (TestOmissionThatEmptiesTheSubjectRefuses) | The FT290 omission empties the file and the red cannot occur. |
| PS13 | 10 | Two `--swap` pairs and one `--omit` make one mutant whose bytes hold all three edits | planned internal/probe/combined_test.go (TestCombinedEditsMakeOneMutant) | A form that applies only the first edit leaves the import in use. |
| PS14 | 11 | A combined edit whose old string matches twice refuses as ambiguous | planned internal/probe/combined_test.go (TestCombinedEditsMatchOnce) | A combined path that skips the count mutates an unnamed site. |
| PS15 | 11 | Each combined old string matches the start bytes, not the result of an earlier edit | planned internal/probe/combined_test.go (TestCombinedEditsMatchOnce) | A sequential apply lets a later old string match text that an earlier edit added. |
| PS16 | 12 | Two edits with overlapping byte ranges refuse with the overlap line before any write | planned internal/probe/combined_test.go (TestCombinedEditsRefuseOverlap) | An apply in argv order rewrites the bytes of the first edit. |
| PS17 | 13 | An argv that leads with `--omit` and then gives two `--swap` pairs reads the mutation cell `swap+swap+omit` | planned internal/probe/combined_test.go (TestCombinedEditsMakeOneMutant) | An argv-order cell reads `omit+swap+swap`, and a fixed word hides the mutant. |
| PS18 | 14 | A combined mutant with `--with-file`, `--insert-before`, `--unwrap`, `--omit-file`, or `--delete` prints the usage line with exit 2 | planned internal/probe/combined_test.go (TestCombinedEditsRefuseSingleForms) | A silent pair by position gives one form two meanings. |
| PS19 | 15 | `--with` together with `--with-file` prints the usage line with exit 2 | planned internal/probe/edit_test.go (TestFormMatrixRefusesPartnerFaults) | A parser that prefers one value drops the other in silence. |
| PS20 | 15 | `--with-file` without `--swap` or `--insert-before`, and `--insert-before` without text, each print the usage line | planned internal/probe/edit_test.go (TestFormMatrixRefusesPartnerFaults) | A missing partner gives an empty insertion or an unused file. |
| PS21 | 15 | `--occurrence` with `--omit-file`, `--delete`, or a combined mutant prints the usage line | planned internal/probe/edit_test.go (TestFormMatrixRefusesPartnerFaults) | An ignored index lets an author believe one match was chosen. |
| PS22 | 15 | `--delete` with `--with` prints the usage line | planned internal/probe/edit_test.go (TestFormMatrixRefusesPartnerFaults) | A deletion that ignores its text hides a mistyped insertion. |
| PS23 | 15 | A combined pair whose `--with` equals its `--swap` refuses with the empty-mutation line | planned internal/probe/combined_test.go (TestCombinedEditsMatchOnce) | A per-pair check that covers only the first pair passes a no-op edit. |
| PS24 | 16 | Zero matches under `--occurrence 2`, and a zero-match edit in a combined mutant, each give the substring-miss row | planned internal/probe/edit_test.go (TestZeroMatchesStayASubstringMiss) | An out-of-range answer for zero matches hides that the text is absent. |
| PS25 | 45 | The help inventory's probe line names `--unwrap`, `--omit-file`, `--with-file`, `--insert-before`, `--occurrence`, and `--delete` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The hand-copied suffix already omits `--unwrap` and `--omit-file`. |
| PS26 | 17 | A Markdown mutation under `--check prose` that adds a long sentence reads `bit` | planned internal/probe/subject_check_test.go (TestProbeTakesTheProseCheck) | The current refusal answers before the baseline. |
| PS27 | 18 | The prose bite's `failed_tests` cell equals the count of prose findings over the mutant | planned internal/probe/subject_check_test.go (TestProbeTakesTheProseCheck) | A probe that counts zero for a non-Go run reads `invalid`. |
| PS28 | 19 | `--check system --run <pattern>` passes the baseline and runs the mutant through the system-suite argv | planned internal/probe/subject_check_test.go (TestProbeTakesTheSystemCheck) | The current refusal answers before the baseline. |
| PS29 | 19 | Over a root that is not the kit, a system probe reads `invalid` with the cause `baseline-refused` | planned internal/probe/subject_check_test.go (TestProbeTakesTheSystemCheck) | A probe that skips the owner's kit refusal runs the suite in a linked repository. |
| PS30 | 20, 24 | The probe request for `--check system --run ^X$` equals `testreport.Prepare` of the same selection | `internal/probe/outcome_test.go` (`TestProbeSelectsTheCheckForm`) | A dropped `--run` pattern gives a different request. |
| PS31 | 21 | `--run` with `--check line-routing` stays the usage refusal | `internal/probe/refusal_test.go` (`TestProbeUsageNamesOneSelectionAndOneMutation`) | A widened rule passes a pattern that the named-check owner ignores. |
| PS32 | 22 | Two `Execute` calls without an inherited selection call the builder twice | planned internal/testreport/runbinary_test.go (TestEachExecuteBuildsItsOwnSelection) | A cached selection runs the mutant against the baseline binary. |
| PS33 | 22 | Over a fixture without an inherited selection, the mutated run's build script reads the mutant bytes | planned internal/probe/subject_check_test.go (TestProbeBuildsTheMutantRunBinary) | A build that starts before the write tests the unmutated tree. |
| PS34 | 23 | `testreport.ProbeNotes` holds exactly the package-expression note and the check-inventory note | `internal/testreport/selection_facts_test.go` (`TestProbeNotesDeriveTheirTokens`) | A kept hand route or rebuild note returns the old route to the help. |
| PS35 | 23 | The probe help names no copy-aside route and no not-a-target sentence | `internal/probe/refusal_test.go` (`TestProbeHelpSpellings`) | A surviving note still sends authors to the copy route. |
| PS36 | 25 | A substring miss with an unknown check answers the unknown-check refusal with exit 2 | planned internal/probe/subject_check_test.go (TestSelectionPrecedesTheSubstringMiss) | The current order answers the miss before the selection parse. |
| PS37 | 25 | A check-form substring miss names the named-check run pattern in its run cell | planned internal/probe/subject_check_test.go (TestSelectionPrecedesTheSubstringMiss) | A hard-coded `all` cell is a second source of the selection. |
| PS38 | 26 | A bite with two failed tests names the first in package-then-test order in `first_failure` | planned internal/probe/verdict_test.go (TestFirstFailureNamesTheFirstFailedTest) | A cell from the last event or from map order names the wrong test. |
| PS39 | 27 | A silent run, a build failure, a prose bite, and a baseline refusal each read `first_failure` null | planned internal/probe/verdict_test.go (TestFirstFailureNamesTheFirstFailedTest) | An empty string or a package name reads as a test name. |
| PS40 | 28 | `--expect TestClampNegative` over a failure of that test reads `bit` | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | A form that ignores the value cannot show that it compares. |
| PS41 | 28 | `--expect TestClampNegative` over a failure of only `TestClampNegative/zero` reads `bit` | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | An exact-only compare misses a parent whose subtest failed. |
| PS42 | 29 | `--expect TestClampPositive` over a failure of only `TestClampNegative` reads `invalid`, `expect-missed`, and exit 1 | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | A bite by another test counts as the witness. |
| PS43 | 30 | `--expect` with `--check prose` prints the usage line with exit 2 | planned internal/probe/verdict_test.go (TestExpectRequiresTheNamedWitness) | A prose run never names a test, so the expectation misses in silence. |
| PS44 | 31 | A fixture `scripts/go-build.sh` that exits 1 gives `Factory.Own` an error that `errors.As` reads as `BuildError` with the output | planned internal/runbinary/runbinary_test.go (TestOwnReportsATypedBuildError) | An untyped wrap cannot be told from a temp-directory failure. |
| PS45 | 31 | `testreport.Execute` over that failing script returns kind `build-failed` with its current refusal bytes | planned internal/testreport/outcome_test.go (TestOnlyAScriptFailureIsBuildFailed) | The current refused kind reads as a refusal at the probe. |
| PS46 | 32 | A build cancelled through its context stays `refused` | planned internal/testreport/outcome_test.go (TestOnlyAScriptFailureIsBuildFailed) | A type built at the generic wrap marks an interrupt as a compile failure. |
| PS47 | 32 | A build with no Go on `PATH` stays `refused` | planned internal/testreport/outcome_test.go (TestOnlyAScriptFailureIsBuildFailed) | A type built at the generic wrap marks a missing toolchain as a compile failure. |
| PS48 | 32 | A test factory whose `Build` returns a plain error stays `refused` | `internal/testreport/outcome_test.go` (`TestExecuteClassifiesTheOutcome`) | A mapping on any builder error reds the current refused subtest. |
| PS49 | 33 | `BuildError.Error()` equals the bytes of the current exit-status wrap over the same output | planned internal/runbinary/runbinary_test.go (TestOwnReportsATypedBuildError) | A reworded error changes the refusal that the gate and lane readers print. |
| PS50 | 34 | A canned go test build failure for an unused import prints the `mutant` table with that compiler line | planned internal/probe/verdict_test.go (TestMutantBuildFailureNamesItself) | A verdict without the table hides the cause in the report. |
| PS51 | 34 | A `BuildError` output with a header line before a `.go:` line gives `BuildLine` the `.go:` line | planned internal/testreport/outcome_test.go (TestOnlyAScriptFailureIsBuildFailed) | A first-line rule names the package header, not the cause. |
| PS52 | 34 | A `BuildError` output with no `.go:` line gives `BuildLine` its first non-empty line | planned internal/testreport/outcome_test.go (TestOnlyAScriptFailureIsBuildFailed) | A rule that requires `.go:` leaves the cell empty for a script fault. |
| PS53 | 31 | A baseline whose run-binary build script fails reads the cause `baseline-build-failed` | planned internal/probe/verdict_test.go (TestMutantBuildFailureNamesItself) | A baseline path that keeps the old mapping reads `baseline-refused`. |
| PS54 | 35 | The probe result schema cell reads `bench-probe-result-v2` | planned internal/probe/result_test.go (TestProbeResultSchemaVersion) | A v1 token over a changed field set lets an older reader fail on a field error. |
| PS55 | 36 | `probe.ReadResult` reads a v2 envelope and refuses a v1 envelope by its schema | planned internal/probe/result_test.go (TestProbeResultSchemaVersion) | A reader with the old field list refuses every new answer. |
| PS56 | 26 | `testreport.Outcome.FailedNames` lists the failed tests in the failures-table order, and a prose run leaves it empty | planned internal/testreport/outcome_test.go (TestOutcomeNamesTheFailedTests) | A second derivation of the names drifts from the printed table. |
| PS57 | 45 | The help inventory's probe line names `--expect` | `cmd/bench/help_inventory_test.go` (`TestHelpInventoryIsComplete`) | The expected inventory line omits the new verdict flag. |
| PS58 | 37 | With another process that holds the subject lock, a probe refuses with the in-use line and exit 1 | planned internal/probe/lock_test.go (TestSubjectLockRefusesASecondProbe) | Two probes mutate one file and one restore certifies the other's damage. |
| PS59 | 38 | A holder's mutant is on disk at the start read, and the holder restores it inside the start-read window, so the re-read differs and the probe refuses with the changed line, and the tree ends unchanged | planned internal/probe/lock_test.go (TestSubjectLockRereadsTheSubject) | A probe without the re-read preserves another probe's mutant as its start bytes. |
| PS60 | 39 | The in-use refusal leaves the subject bytes, the Bench home, and the go marker unchanged | planned internal/probe/lock_test.go (TestSubjectLockRefusesASecondProbe) | A lock taken after the baseline starts a run child first. |
| PS61 | 39 | With the gate lock and the subject lock both held by other processes, the gate-lock line wins | planned internal/probe/lock_test.go (TestSubjectLockFollowsTheGateLock) | A reversed order shows the subject line under a live gate. |
| PS62 | 40 | A probe on a second subject runs while another process holds the first subject's lock | planned internal/probe/lock_test.go (TestSubjectLockAllowsAnotherSubject) | A repository-wide lock serializes unrelated probes. |
| PS63 | 41 | After an interrupted probe, a new probe on the same subject reaches its verdict | planned internal/probe/lock_test.go (TestSubjectLockReleasesOnInterrupt) | A lock that survives the interrupt blocks the subject until a manual clean. |
| PS64 | 42 | A seeded lock file with no holder grants the lock, and the probe removes the file at its end | planned internal/probe/lock_test.go (TestSubjectLockIgnoresAFileWithoutAHolder) | A presence lock created with exclusive create refuses forever after a kill. |
| PS65 | 42 | A lock file that the lock-owner hook replaces between the open and the identity check gives a retry that holds the new file | planned internal/probe/lock_test.go (TestSubjectLockRechecksTheIdentity) | A lock on an unlinked file lets a second probe lock the new file too. |
| PS66 | 43 | A lock-owner hook that replaces the lock file on all three attempts gives the lock-unavailable line before the baseline | planned internal/probe/lock_test.go (TestSubjectLockRechecksTheIdentity) | An unbounded retry hangs the probe. |
| PS67 | 43 | A regular file at `$BENCH_HOME/probe-locks` refuses with the lock-unavailable line and starts no run child | planned internal/probe/lock_test.go (TestSubjectLockRefusesAnUnusableDirectory) | A lock setup fault without a refusal starts the baseline unlocked. |
| PS68 | 44 | A completed probe leaves the record surfaces unchanged | `internal/probe/outcome_test.go` (`TestProbeWritesNoRecord`) | A kept lock file or lock directory adds an entry under the Bench home. |
| PS69 | 46 | The reference guide names each new form, the subject lock, and the `mutant` table | review-owned: the reviewer reads the probe paragraphs | The prose check grades sentences, not content. |
| PS70 | 46 | `CONTEXT.md` defines **combined mutant** and **expected witness** with Avoid lists | review-owned: the reviewer reads the entries | The prose check grades sentences, not terms. |
| PS71 | 46 | The delegate skill's probe sentence at line 104 names the new forms | review-owned: the reviewer reads the sentence; the anchored line 81 is gate-owned through `internal/anchors/registry_data.go:214` and `internal/conformance/fixture_bite_test.go:360` | No check reads the content of line 104. |
| PS72 | 46 | `CHANGELOG.md` names each new form and the schema v2 change | review-owned: the reviewer reads the entry | No check reads the changelog content. |
| PS73 | 47 | `docs/known-issues.md` holds no row that names an unused probe import or an unused probe variable | review-owned: the reviewer reads the table | No check reads the issue table. |
| PS74 | 48 | The probe-cost row names the batch obligation as its retire condition, and `roadmap/FT168.md` keeps that obligation | review-owned: the reviewer reads both files | The current condition holds once the combined mutant ships and forces an early removal. |
| PS75 | 38 | A subject that the start-read hook deletes reads the absent-subject line at the re-read, and no run child starts | planned internal/probe/lock_test.go (TestSubjectLockRereadsTheSubject) | A re-read that compares only bytes fails on a read error with no named cause. |
| PS76 | 14 | Each argv that names two mutation groups, `--delete` or `--insert-before` included, prints the usage line with exit 2 | `internal/probe/refusal_test.go` (`TestProbeUsageNamesOneSelectionAndOneMutation`) | A pair rule that names only the four current forms lets a new form combine in silence. |

Not covered: story 49 — a reviewed exclusion under D2; the Out of scope section prices it.
Not covered: story 50 — a reviewed exclusion under D2; the Out of scope section prices it.
Not covered: story 51 — NRO owns `focused_exit_code`, and this spec adds no row for it.

### Edge inventory

The audience is this repository and every linked repository.
A linked repository has no system suite, so its system probe reads `invalid` with the cause `baseline-refused` (PS29).
An absent lock directory differs from a path that is not a directory: the first is created, and the second refuses (PS67).

The hostile-input walk covers these classes.

A `--with-file` or subject path with spaces or glob characters stays literal.
A control byte in the subject name keeps the current render refusal.
A control byte in the first build line passes through `sanitize.Preview` before the encoder.
A test name that reads as a number keeps the encoder's quotes, because the expected row derives through `toon.TableTyped`.
A FIFO operand for `--with-file` refuses without a block, through the shared non-blocking open (PS2).
A lock path derives from a digest, so no operand byte reaches the lock file name.

**Won't handle** — an inherited `BENCH_RUN_BINARY` under a build-input mutant — the freshness refusal answers it, and the normal shell sets no inherited value.
**Won't handle** — a line-number anchor — `--occurrence` covers the case of two identical lines, and a line number moves with every edit.
**Won't handle** — a combined mutant across two files — the subject lock and the preserve owner hold one file, and two serial probes remain the route.
**Won't handle** — a wall-time cell in the verdict row — the packages table already prints `elapsed_ms` for each package.
**Won't handle** — a first-failure name for a prose run — a finding names a sentence, not a test, and the report lists each finding.
**Won't handle** — a lock that blocks and waits — a waiting probe hides a concurrent author, and the refusal names the subject.

**Won't handle** — a read of a v1 probe result after the upgrade — `probe.ReadResult` refuses it by schema, and a fresh probe gives a v2 result.
**Won't handle** — concurrent probes on two files of one package — each run can see the other's mutant, and the caller serializes probes per package.

## Ownership fences

The ticket slice derives this fence, so it equals the union of the ticket Writes lines, plus the spec folder and the review pickup.
It holds the owner and test files, the command registry closure, the anchor registry closure, and each canary unit that pins an edited file.
The NRO import test in the record command package reads the schema token.
The named-check file of the testreport package holds the prose name that becomes the one export.
Every reader of that name is in the fence.

- `.agents/skills/bench-craft-delegate/SKILL.md`
- `.bench/BENCH-reference.md`
- `CHANGELOG.md`
- `CONTEXT.md`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/test_command.go`
- `docs/known-issues.md`
- `internal/anchors/registry_calibration.go`
- `internal/anchors/registry_calibration_test.go`
- `internal/anchors/registry_charge_binding.go`
- `internal/anchors/registry_charge_binding_test.go`
- `internal/anchors/registry_chunk_chain.go`
- `internal/anchors/registry_chunk_chain_test.go`
- `internal/anchors/registry_data.go`
- `internal/anchors/registry_data_test.go`
- `internal/anchors/registry_debug_loop.go`
- `internal/anchors/registry_decision_maps.go`
- `internal/anchors/registry_decision_maps_test.go`
- `internal/anchors/registry_ft311_preparation.go`
- `internal/anchors/registry_ft311_review_dispatch.go`
- `internal/anchors/registry_retained_workflow.go`
- `internal/anchors/registry_ticket_passes.go`
- `internal/anchors/registry_ticket_passes_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/probe/baseline.go`
- `internal/probe/baseline_test.go`
- `internal/probe/combined_test.go`
- `internal/probe/command.go`
- `internal/probe/delete_test.go`
- `internal/probe/edit.go`
- `internal/probe/edit_test.go`
- `internal/probe/invalid_cause_test.go`
- `internal/probe/lock.go`
- `internal/probe/lock_holder_test.go`
- `internal/probe/lock_test.go`
- `internal/probe/mutation_input.go`
- `internal/probe/omit_file_test.go`
- `internal/probe/outcome_test.go`
- `internal/probe/probe.go`
- `internal/probe/probe_test.go`
- `internal/probe/refusal_test.go`
- `internal/probe/replacement_test.go`
- `internal/probe/result.go`
- `internal/probe/result_test.go`
- `internal/probe/subject.go`
- `internal/probe/subject_check_fixture_test.go`
- `internal/probe/subject_check_test.go`
- `internal/probe/unwrap_test.go`
- `internal/probe/verdict.go`
- `internal/probe/verdict_test.go`
- `internal/reviewrecord/recordcmd/verification_import_test.go`
- `internal/runbinary/build_script_test.go`
- `internal/runbinary/runbinary.go`
- `internal/runbinary/runbinary_test.go`
- `internal/testreport/canned_go_test.go`
- `internal/testreport/check_test.go`
- `internal/testreport/command.go`
- `internal/testreport/failure_lines_test.go`
- `internal/testreport/full_failure_test.go`
- `internal/testreport/named_check.go`
- `internal/testreport/outcome.go`
- `internal/testreport/outcome_test.go`
- `internal/testreport/runbinary_test.go`
- `internal/testreport/selection_facts.go`
- `internal/testreport/selection_facts_test.go`
- `internal/testreport/system_run_test.go`
- `internal/testreport/testreport.go`
- `reviews/probe-subjects-and-mutations.md`
- `roadmap/FT168.md`
- `specs/probe-subjects-and-mutations/`
- `tests/canary/claude-agent-definitions/agent-unnamed-in-skill`
- `tests/canary/claude-agent-definitions/model-declared`
- `tests/canary/claude-agent-definitions/name-mismatch`
- `tests/canary/claude-agent-definitions/shell-tool-absent`
- `tests/canary/claude-agent-definitions/skill-names-missing-agent`
- `tests/canary/claude-agent-definitions/spawning-tool`
- `tests/canary/claude-agent-definitions/tools-absent`
- `tests/canary/docs-currency-token-diet/benchref-imported`
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped`
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated`
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key`
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy`
- `tests/canary/skills-index-command-adapters/dangling-index`
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted`
- `tests/canary/skills-index-command-adapters/missing-index-field`
- `tests/canary/skills-index-command-adapters/stale-index-wording`
- `tests/canary/skills-index-command-adapters/unindexed-skill`
- `tests/canary/workflow-guidance-anchors/agents-handoff-section-rule`
- `tests/canary/workflow-guidance-anchors/calibration-delegate-pointer`
- `tests/canary/workflow-guidance-anchors/changelog-reduced-schema-columns`
- `tests/canary/workflow-guidance-anchors/changelog-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts`
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary`
- `tests/canary/workflow-guidance-anchors/context-decision-map-term`
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term`
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary`
- `tests/canary/workflow-guidance-anchors/delegate-cap-change-pinning-package`
- `tests/canary/workflow-guidance-anchors/delegate-charge-effort-cap`
- `tests/canary/workflow-guidance-anchors/delegate-coverage-row-charge`
- `tests/canary/workflow-guidance-anchors/delegate-coverage-row-red-green`
- `tests/canary/workflow-guidance-anchors/delegate-cross-harness-reviewer-pointer`
- `tests/canary/workflow-guidance-anchors/delegate-exec-only-every-caller`
- `tests/canary/workflow-guidance-anchors/delegate-model-id-escalation`
- `tests/canary/workflow-guidance-anchors/delegate-own-family-native-surface`
- `tests/canary/workflow-guidance-anchors/delegate-parallel-route-anchor`
- `tests/canary/workflow-guidance-anchors/delegate-release-at-acceptance`
- `tests/canary/workflow-guidance-anchors/delegate-resume-handoff-contents`
- `tests/canary/workflow-guidance-anchors/delegate-self-probe-missing-row`
- `tests/canary/workflow-guidance-anchors/delegate-stash-refusal-anchor`
- `tests/canary/workflow-guidance-anchors/delegated-per-ticket-author`
- `tests/canary/workflow-guidance-anchors/dg-6-owner-reference`
- `tests/canary/workflow-guidance-anchors/fix-pass-sentinel-anchor`
- `tests/canary/workflow-guidance-anchors/prepared-review-delegate-handoff-route`
- `tests/canary/workflow-guidance-anchors/prepared-triage-bounds`
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule`
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer`
- `tests/canary/workflow-guidance-anchors/reference-category-context`
- `tests/canary/workflow-guidance-anchors/reference-category-oracle`
- `tests/canary/workflow-guidance-anchors/reference-category-setup`
- `tests/canary/workflow-guidance-anchors/reference-category-work`
- `tests/canary/workflow-guidance-anchors/reference-gate-authority`
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship`
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback`
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term`
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape`
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner`
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner`
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance`
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route`
- `tests/canary/workflow-guidance-anchors/shared-worktree-path-pin`

## Out of scope

- The batch is FT168's next slice: several packages in one probe, and independent mutants that share one baseline and one private build. Estimate: 10 edits and 4 gate runs. It serves the probe-cost known-issue row and the 2026-10-09 learning "A per-expectation probe loop multiplies the per-call recompile of bench probe". It needs an FT318 decision on how a record imports several verdicts. Single-mutant probes remain the surviving caller.
- `bench canary <fixture>` runs one canary fixture: estimated 6 edits and 3 gate runs. The form collides with the current `bench canary [root]` grammar (`cmd/bench/help_inventory_test.go:97`), so its successor decides the operand shape. `bench test --check <name>` remains the surviving caller.
- A `--names` form that compares failing test names with a base set: estimated 5 edits and 3 gate runs. Its successor designs the base-set input and reads PS-C's `FailedNames` list. The failures table remains the surviving caller.
- A repeatable `--package` on `bench test` (IDEAS, 2026-10-09): estimated 4 edits and 2 gate runs, with the batch successor. One call per package remains the surviving caller.
- A probe that prints the mutated text and reruns the green after the restore (IDEAS, 2026-10-09): estimated 5 edits and 2 gate runs. It is not in the decision source. The verdict row and a `bench test` call remain the surviving callers.

## Further notes

Flagged additions:

- The emptying-omission refusal reads "replaces the emptying `--omit-file`" as a refusal that names `--delete`.
- The `mutant` table is a new block after `probe_execution`, not a new verdict-row cell.
- The help inventory derives its probe suffix from `probe.Usage`, because the hand copy has drifted.
- The removal of the unused-import and unused-variable known-issue row follows from the self-naming build failure and the combined form. The first build line names either compile error.
- The cause values `expect-missed` and `baseline-build-failed` are new under the existing verdict `invalid`. The reviewer accepted `baseline-build-failed` on 2026-10-09, because no reader outside `internal/probe` parses the baseline cause.
- The selection parse moves before the substring-miss answer, so an unknown check now wins over a miss.
- The subject-changed and lock-unavailable refusals are new refusal lines.
- The rebuild note leaves `ProbeNotes`.

Posture changes that red existing tests:

- `TestProbeUsageNamesOneSelectionAndOneMutation` loses its "two mutations" subtest, because that argv is now a valid combined mutant.
- `TestProbeRefusesProseAndSystemChecks` retires, because both checks become probe targets.
- `TestProbeNotesDeriveTheirTokens` expects two notes, not four.
- `TestProbeHelpSpellings` and `TestHelpInventoryIsComplete` expect the new usage line.
- `TestExecuteClassifiesTheOutcome`, `TestFullFailedTestsCountsTests`, and the full-failure table test compare `Outcome` field-wise.
- The literal verdict-row pins at `internal/probe/invalid_cause_test.go:26`, `internal/probe/refusal_test.go:58`, and `internal/reviewrecord/recordcmd/verification_import_test.go:46` on the NRO tip gain the `first_failure` cell.
- `TestProbeRefusesWhenPreservationFails` and `TestProbePreservesTheSubjectBeforeTheRun` stay green, because the lock directory is a sibling of `$BENCH_HOME/probe`.

Boundaries:

- NRO owns `probe_execution`, `focused_exit_code`, `command_exit_code`, `trailing_bytes`, and `probe.ReadResult`. Under D1, this spec changes the verdict-row field set after NRO lands and moves the schema token to v2.
- FT318 owns the consumption of the probe result by the verification record. This spec changes no record form.
- FT98 owns preservation support. Each new form restores through the current preserve owner, and the deliberate-red verdict stays.

Response bound: a biting answer already spills past the 10-line bound after NR-C3, and an open learning asks the reviewer for that decision. This spec adds no line to a biting or silent answer. The `mutant` table adds two lines only to a build-failed answer.

Known issues: the probe-cost row is on `main` at `docs/known-issues.md:12` since 23c59775. Its retire condition is "FT168: one probe takes several mutations in one private build". That condition holds once the combined mutant ships, so PS74 rewords it to the batch obligation. The NRO build removed the probe-exit row, and `main` has no probe-exit row.

Source-clause mapping:

- "multiline replacement text for `--with` from a file" -> PS1, PS2.
- "Markdown files and registered system checks" -> PS26-PS29.
- "Named fixture and class evidence must remain producible through the registered check owner" -> PS30.
- "builds through the worktree build" -> PS32, PS33.
- "One probe takes several mutations" -> PS13-PS18, PS23, PS24.
- "`--insert-before` for an insertion" -> PS3-PS5.
- "`--expect` for the expected red" -> PS40-PS43.
- "One probe grades several packages in one private build" -> Out of scope under D2.
- "`bench canary <fixture>`" -> Out of scope under D2.
- "A `--names` form" -> Out of scope under D2.
- "Each new face must restore the original subject" -> PS10, PS11, PS59, and the current restore tests.
- "The verdict prints the mutated run's exit" -> NRO.
- "a deletion form replaces the emptying `--omit-file`" -> PS9-PS12.
- "a mutant compile error names itself" -> PS44-PS53.
- "A second probe on the same file refuses" -> PS58-PS68.
- "The verdict row also names the first failing test" -> PS38, PS39, PS56.
- "A swap accepts a line or occurrence anchor" -> PS6-PS8.

The pre-review proof checklist follows.

- Cited symbols: `probe.Command`, `mutationInput`, `mutate`, `resolveSubject`, `preserve`, `restore`, `release`, `gradeForms`, `selectionArgs`, `renderInvalidMutation`, `gradeCheckTarget`, `gradeGateLock`, `testreport.Execute`, `testreport.Outcome`, `testreport.HandProbeRoute`, `testreport.ProbeNotes`, `runbinary.Factory.Own`, `runBuildScript`, `gocache.RecordLock`, and `testHelpSuffix` resolve.
  They resolve at aabd462d. `probe.ReadResult` and `probetest` resolve on the NRO build tip.
- Import edges: `internal/testreport` already imports `internal/runbinary`, and `internal/probe` already imports `internal/testreport`. The lock owner adds the production edge from `internal/probe` to `internal/gocache`. `go list -deps ./internal/gocache` names neither `internal/probe` nor `internal/testreport`, so the edge makes no cycle.
- Source-row clauses and occurrences: the full FT168 body and its 24 occurrences were read, through the retros in Git history and `capture/learnings.md`. The light-path repair pass retro text was not found.
- Promised field labels: `first_failure`, `mutant`, `first_build_line`, `FailedNames`, `BuildLine`, `BuildError`, `bench-probe-result-v2`, `expect-missed`, `baseline-build-failed`, `insert`, and `delete`.
- Changed-function callers: `testreport.Execute` has the callers `gradeBaseline`, `probe`, and `testreport.Command`. `HandProbeRoute` has the readers `gradeCheckTarget` and `ProbeNotes`. `ProbeNotes` has the reader `probeGrammar`. `runbinary.Factory.Own` has the readers `productionRunBinaryOwner`, `internal/gate/prospective.go:37`, the `ReuseOrOwn` callers at `internal/gate/prospective.go:39` and `internal/gate/lane.go:274`, and `internal/testreport/command.go:43`; PS49 keeps their error bytes.
- Copy survival: PS34 reds a surviving hand route or rebuild note. PS25 and PS57 red a hand-copied help suffix, and PS37 reds a hard-coded run cell.
- Rendered-shape readers: the strings "not probe targets" and "probe focused run unsupported" occur in `internal/probe/probe.go`, `internal/probe/refusal_test.go`, `internal/testreport/selection_facts.go`, and `internal/testreport/selection_facts_test.go`. The NRO import tests read the schema token. `bench anchors .agents/skills/bench-craft-delegate/SKILL.md` at 23c59775 lists the needle "report the observed result and the mutation's kind (omission or swap)" at line 81, which `internal/anchors/registry_data.go` and `internal/conformance/fixture_bite_test.go` pin.
- Pin operators: the rows compare exact rendered rows, exact subject bytes, exact modes, and exact error bytes.
- Entry reads: the lock owner reads `BENCH_HOME` through `benchhome.Dir`, as `preserve` does.
- Derived expectations: expected rows derive through `toon.TableTyped`; no expectation reads the code under test.
- Consolidated rules: the `--omit-file` and `--with-file` readers consolidate into one reader with the current refusal set, and PS2 runs one table over both flags.
- Quantified obligations: "each new face must restore" covers `--with-file`, `--insert-before`, `--occurrence`, `--delete`, combined mutants, prose, and system.
- Workflow-step writes: none.

PS33 copies a sealed executable pair from a fixture build script. `freshness.VerifyAt` accepts that copy under four conditions, and the PS33 fixture holds all four:

1. The bytes stay unchanged, and the seal sits at `<output>.seal` (`internal/freshness/freshness_verify.go:51-71`).
2. The subject is outside the digest closure (`internal/freshness/freshness_buildinputs.go:50-104`). `clamp.go` qualifies, because the fixture's `cmd/bench/main.go` imports nothing. A subject inside the closure makes the seal stale, and the selection refuses.
3. No path component is a symlink (`rejectSymlinkComponents` in `internal/freshness/freshness_buildinputs.go`).
4. `BENCH_RUN_BINARY` is unset, not empty. `newFixture` sets it, and `ReuseOrOwn` tests its presence (`internal/runbinary/runbinary.go:71`).

Sources not re-read: the light-path repair pass retro, and the learning texts behind six occurrence lines that exist only in `roadmap/FT168.md`.
