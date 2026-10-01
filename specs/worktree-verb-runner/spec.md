# Run every worktree verb in tests through one verb runner

Status: staged

Decision source: `specs/worktree-verb-runner/decisions/worktree-seams.md` (ready compiled map).

Review status: the coordinator's independent review round is pending.

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

The static census then refuses a test-file call to a verb entry or to its joins
form outside the verb runner file. The top-level test count holds, the serial set
stays at or below its ceiling, and the joins value does not change. The seam
reduction spec follows this spec and rewrites each stubbed test once, on the verb
runner.

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
9. As a test author, I want to give the runner a joins value, so that a stubbed test runs through the same runner.
10. As a test author, I want to give the verb runner stdin, so that the exec verb runs through the runner.
11. As a test author, I want the exec assignment in the verb result, so that the dispatcher's assignment report stays testable.
12. As the seam reduction author, I want the runner to accept the kit root, the home, and the clock, so that its signature stays.
13. As a test author, I want a kit value to fail until the seam reduction wires it, so that no test trusts an ignored kit.
14. As a test author, I want a clock value to fail until the seam reduction wires it, so that no test trusts an ignored clock.

### The migration

15. As a maintainer, I want every direct verb call in a worktree test on the runner, so that one way to call a verb exists.
16. As a maintainer, I want the named run wrappers deleted, so that no second way to call a verb survives.
17. As a maintainer, I want the named fingerprint extractors deleted, so that one reader owns the fingerprint.
18. As a maintainer, I want each positional fixture tuple to become a named fixture value, so that a call site names its fixture part.
19. As a reviewer, I want the top-level test count to hold through the migration, so that no test is removed or merged.
20. As a reviewer, I want the serial set to stay at or below its ceiling, so that the verb runner does not cost parallel eligibility.
21. As a reviewer, I want each migrated test to keep its pass or skip outcome, so that no test turns into a silent skip.
22. As a reviewer, I want each migrated test to keep every assertion it made, so that the migration does not weaken the suite.
23. As a reviewer, I want no line-count target for the worktree tests, so that the acceptance stays structural.

### The census refusal

24. As a maintainer, I want the census to report a verb entry call outside the runner file, so that a new direct call turns red.
25. As a maintainer, I want the census to report a joins form call outside the runner file, so that no stubbed test bypasses the runner.
26. As a maintainer, I want the census to report a verb entry used as a value, so that a function table cannot bypass the runner.
27. As a maintainer, I want the census to report a verb call inside a closure or a subtest, so that nesting is no shelter.
28. As a maintainer, I want the census to report a package-level verb reference, so that a table outside a function is no shelter.
29. As a maintainer, I want the census to allow every verb reference inside the verb runner file, so that the runner itself is legal.
30. As a maintainer, I want the census to derive the verb entries from their signatures, so that a new entry needs no list edit.
31. As a maintainer, I want the census to derive each joins form from its verb entry's body, so that a renamed joins form stays covered.
32. As a maintainer, I want the live-tree census to report no direct verb call, so that the migration's end state is pinned.

### Reviewed exclusions

33. As the seam reduction author, I want the joins value and all production files unchanged, so that the seam reduction decides each field once.
34. As a maintainer of another package, I want my package's test harness left unchanged, so that this spec stays in the worktree package.

## Implementation decisions

**The verb runner.** One test file, `verb_runner_test.go`, holds the verb runner, the verb call value, the verb result, and the runner's own tests. The runner selects a verb by a typed verb key that the runner file declares. Each key names one verb entry, and a verb with a joins form also names that form. A key exists for each verb that a test runs. The `shell` verb has no key, because no test runs its verb entry.

The verb call value carries the root, the Bench home, the kit root, the clock, stdin, an optional joins value, and the arguments. Without a joins value, the runner calls the verb entry. With a joins value, the runner calls the verb's joins form with that value unchanged. A verb without a joins form refuses a joins value. A call that carries a kit value or a clock value fails, and the message names the seam reduction spec. That spec wires both values and keeps the call value's fields.

**The verb result.** The verb result carries the exit code, stdout, stderr, and the assignment that the exec verb resolved. A verb that returns its output as a string, for example `list`, fills stdout and leaves stderr empty. The rows reader decodes stdout with `axitest.DecodeDocument` and returns the rows of one table block, or an error.

The fingerprint reader returns one value or an error. It reads the `fingerprint` cell of every row in a decoded table, and all cells must agree. When stdout does not decode, it reads the `fingerprint=` cell of the one record line that carries it. Both readers return an error instead of a failed test, so the runner's own tests can grade the error path. A caller fails its own test on the error.

**Fixture values.** One test file, `verb_fixture_test.go`, declares every named fixture value type. A fixture builder stays in its present file and returns one named value. A positional fixture tuple is a test-file function that returns at least two values and no `error`. One of the values is a `Creation` or a `[]Creation`, or the function returns three or more values. A stub builder that returns a joins value with its probe is not a fixture tuple, and the seam reduction spec owns it.

**The census refusal.** A new census file, `verb_call_census_test.go`, reuses the parse helpers of `parallel_census_test.go`. A verb entry is an exported function of a non-test file with a parameter named `args` of type `[]string`. A verb's joins form is a function that a verb entry's body calls with `defaultJoins()` as its first argument. The census reports each identifier that names a verb entry or a joins form in a test file other than `verb_runner_test.go`.

The census walks function bodies, closures, and package-level declarations. Each report names the file, the line, the enclosing declaration, and the verb form. The census decides by bare identifier, the same way as the serial helper edge.

**Order.** The runner and its first family land first, as the seam chunk. The shared fixture values follow, then the verb families, then the landing family. The census refusal lands last, after every migration ticket, as the contract step.

**Constraints that every ticket keeps.** No production file changes. The top-level test count holds, and a ticket that adds a test raises `worktreeTestCount` in the same change. The serial set stays at or below its ceiling of 46. The lane's growth ratchet refuses an over-budget test file that grows, so each ticket keeps such a file at or below its base line count. A new file stays at or below 400 lines.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| VR-C1 / `1-add-verb-runner.md` | The verb runner and verb result exist, and the reset family runs on them | VR1, VR2, VR3, VR4, VR5, VR6, VR7, VR8, VR9, VR10, VR11, VR12, VR13, VR14, VR15, VR16, VR17 | The runner tests in `verb_runner_test.go` and the reset family tests | yes |
| VR-C2 / `2-name-shared-fixtures.md` | The shared assignment and pool fixtures return named values | VR18 | The package suite and the count pin | no |
| VR-C3 / `3-migrate-cleanup-verbs.md`, `4-migrate-query-and-create-verbs.md`, `5-migrate-merge-and-reauthorize.md` | Every non-landing verb runs on the verb runner | VR19, VR20, VR21, VR22, VR23 | The family test files and the count pin | no |
| VR-C4 / `6-name-landing-fixtures.md`, `7-migrate-landing-composition.md`, `8-migrate-landing-effects.md` | Every landing verb runs on the verb runner with named landing fixtures | VR24, VR25, VR26 | The landing test files and the count pin | no |
| VR-C5 / `9-refuse-direct-verb-calls.md` | The census refuses a direct verb call outside the runner file, and the package-wide end state holds | VR27, VR28, VR29, VR30, VR31, VR32, VR33, VR34, VR35, VR36, VR37, VR38, VR39, VR40, VR41, VR42, VR43, VR44 | The census tests in `verb_call_census_test.go`, the count pin, and the serial ceiling | yes |

After each chunk, freeze its predecessor and current tips for Standards, Spec, and Coverage review. Each chunk review also records the differential name sets of VR31 and VR32 and the assertion comparison of VR33 for that chunk. Ticket 9 owns those rows because it is the last ticket that touches the package.

## Testing decisions

- A good test drives a real verb through the verb runner over a real git fixture and examines the verb result. The runner's own tests compare the result with a direct call inside the runner file.
- The reader tests feed the rows and fingerprint readers a real verb output where one exists. A synthetic table derives its text through the same `toon.Table` call that the producer makes.
- The census tests plant synthetic file sets with `plantTestFiles`, which is the precedent of `parallel_census_test.go`. One live-tree test runs the census over the package.
- The gate observes the feature through the worktree package's Go tests in the test phase. No new gate check is added.

### Seam diagram

    trigger: a worktree test
        │
        ▼
    verb call value  ──▶  [ verb runner: verb entry or joins form ]  ──▶  verb result
                              ◀ tests attach here: compare with a direct call; read rows and fingerprint

    trigger: go test ./internal/worktree
        │
        ▼
    test files + source files  ──▶  [ verb call census ]  ──▶  report lines
                              ◀ tests attach here: synthetic file sets and the live tree

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| VR1 | 1 | Each verb key returns the usage refusal of its own verb for a usage-refusing argument list | planned `TestVerbRunnerKeyReachesItsOwnVerb` in internal/worktree/verb_runner_test.go | A key wired to another verb's entry returns the other verb's usage line, which the expectation reads from that verb's usage constant |
| VR2 | 2 | The verb result's exit code, stdout, and stderr equal those of a direct call to the same verb entry with the same input | planned `TestVerbRunnerReturnsBothStreamsAndTheExitCode` in internal/worktree/verb_runner_test.go | A runner that drops stderr or returns a fixed exit code differs from the direct call |
| VR3 | 3 | The rows reader returns the decoded rows of the named table block of a real `list` output | planned `TestVerbResultRowsDecodeTheWholeDocument` in internal/worktree/verb_runner_test.go | A reader that returns the wrong block or no rows differs from the rows that `axitest.DecodeDocument` returns |
| VR4 | 4 | The rows reader returns an error for stdout that holds a table followed by a line outside the TOON grammar | planned `TestVerbResultRowsRefuseAPartialDocument` in internal/worktree/verb_runner_test.go | A substring reader or a reader that returns an empty set on a decode failure returns no error |
| VR5 | 5 | The fingerprint reader returns the shared `fingerprint` cell of a real `clean --landed` plan | planned `TestVerbResultFingerprintReadsTheTableCell` in internal/worktree/verb_runner_test.go | A reader that reads another cell or the first 64-hex run of stdout returns a different value |
| VR6 | 6 | The fingerprint reader returns the `fingerprint=` cell of a real `reset --to` plan record | planned `TestVerbResultFingerprintReadsTheRecordCell` in internal/worktree/verb_runner_test.go | A table-only reader returns an error for the record, and a reader that takes the `next=` cell returns more text |
| VR7 | 7 | The fingerprint reader returns an error for stdout that carries no fingerprint | planned `TestVerbResultFingerprintRefusesAnAbsentValue` in internal/worktree/verb_runner_test.go | A reader that returns the empty string with no error passes a plan that has no fingerprint |
| VR8 | 7 | The fingerprint reader returns an error for a table whose rows carry two different fingerprints | planned `TestVerbResultFingerprintRefusesConflictingCells` in internal/worktree/verb_runner_test.go | A reader that returns the first cell passes a set whose rows disagree |
| VR9 | 8 | The fingerprint reader returns a numeric-looking value exactly as the producer wrote it before `toon.Table` quoted it | planned `TestVerbResultFingerprintKeepsANumericLookingCell` in internal/worktree/verb_runner_test.go | A reader that keeps the quotes or decodes a number returns different text |
| VR10 | 9 | With a joins value, the runner calls the verb's joins form with that value, and the value's stub runs | planned `TestVerbRunnerPassesTheJoinsValue` in internal/worktree/verb_runner_test.go | A runner that calls the verb entry runs the default instead, and the stub's probe stays at zero |
| VR11 | 10 | Stdin that the call carries reaches the exec child | planned `TestVerbRunnerFeedsStdinToExec` in internal/worktree/verb_runner_test.go | A runner that passes a nil reader gives the child an empty input, and its echo differs |
| VR12 | 11 | The exec verb result carries the assignment ID that the exec target resolved to | planned `TestVerbRunnerReturnsTheExecAssignment` in internal/worktree/verb_runner_test.go | A runner that calls `ExecCommand` or drops the second result returns an empty assignment |
| VR13 | 12 | The verb call value declares a kit root field and a clock field beside the root and the home | review-owned: the Standards axis reads the call value in the verb runner file | A call value without the two fields forces the seam reduction spec to change the runner's signature |
| VR14 | 13 | A call that carries a kit value fails with a message that names the seam reduction spec | planned `TestVerbCallRefusesAKitValue` in internal/worktree/verb_runner_test.go | A runner that ignores the kit value runs the verb and returns no refusal |
| VR15 | 14 | A call that carries a clock value fails with a message that names the seam reduction spec | planned `TestVerbCallRefusesAClockValue` in internal/worktree/verb_runner_test.go | A runner that ignores the clock value runs the verb and returns no refusal |
| VR16 | 16, 17 | No reset-family run wrapper or fingerprint extractor exists after ticket 1 | review-owned: the VR16 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR17 | 18 | `restoreFixture` returns one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan in Further notes omits `restoreFixture` | A positional `restoreFixture` appears in the scan output |
| VR18 | 18 | Each shared assignment and pool fixture builder returns one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits `newOwnedAssignment`, `newPendingAssignment`, `newOwnedSubmoduleAssignment`, `newResidueGuardFixture`, `unprovableLandedAssignment`, `newReclaimPool`, and `poolRootFixture` | A positional builder appears in the scan output |
| VR19 | 16, 17 | No cleanup-family run wrapper or fingerprint extractor exists after ticket 3 | review-owned: the VR19 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR20 | 18 | Each cleanup-family fixture builder returns one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits `landedSetFixture`, `retainedMemberFixture`, `removableSetFixture`, and `refusedUnlandedRelease` | A positional builder appears in the scan output |
| VR21 | 16, 17 | No create or exec run wrapper and no inline 64-hex fingerprint match exists after ticket 4 | review-owned: the VR21 command in Further notes prints no line | A surviving declaration or inline match prints its file and line |
| VR22 | 16 | No merge run wrapper exists after ticket 5 | review-owned: the VR22 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR23 | 18 | `mergeFixture` and `reauthorizeFixture` each return one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits both builders | A positional builder appears in the scan output |
| VR24 | 18 | Each landing fixture builder returns one named fixture value that `verb_fixture_test.go` declares | review-owned: the tuple scan omits every landing builder in the Enumerations list | A positional builder appears in the scan output |
| VR25 | 16 | No landing run wrapper exists after ticket 7 | review-owned: the VR25 command in Further notes prints no line | A surviving declaration prints its file and line |
| VR26 | 15 | The landing effect, resume, census, and journey test files call no verb form directly after ticket 8 | review-owned: the verb form `rg` in Further notes prints no line for the ticket 8 files | A surviving direct call prints its file and line |
| VR27 | 18 | No test-file function in the package returns a positional fixture tuple | review-owned: the tuple scan in Further notes prints no line | A positional builder that any ticket missed appears in the scan output |
| VR28 | 16, 17 | No function in the run wrapper list or the fingerprint extractor list exists in the package | review-owned: the VR28 `rg` in Further notes prints no line | A surviving wrapper or extractor prints its file and line |
| VR29 | 19 | The package's top-level test count equals `worktreeTestCount`, which rises only by the tests this spec adds | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | A removed or merged test drops the count below the pin |
| VR30 | 20 | The serial set stays at or below the ceiling of 46 | `internal/worktree/parallel_census_test.go` (`TestSerialSetStaysBelowTheCeiling`) | A runner that binds the process environment adds every caller to the serial set |
| VR31 | 21 | The PASS name set of the package's fresh test run equals the chunk base's set plus the added tests | review-owned: each chunk review compares the `go test -count=1 -json ./internal/worktree` name sets | A migration that makes a test skip moves its name out of the PASS set |
| VR32 | 21 | The SKIP name set of the package's fresh test run equals the chunk base's set | review-owned: each chunk review compares the `go test -count=1 -json ./internal/worktree` name sets | A migration that adds a capability skip adds a name to the SKIP set |
| VR33 | 22 | Each migrated test keeps each assertion on the exit code, the streams, and the repository state that it made at the chunk base | review-owned: the Spec axis compares each migrated test with its base form | A migration that drops an assertion stays green, so only the side-by-side read catches it |
| VR34 | 15 | No test outside `verb_runner_test.go` declares an output buffer pair for a verb call | review-owned: the Coverage axis reads each chunk diff at each removed verb call | A test that keeps its buffers and calls the runner restates part of the call |
| VR35 | 24 | The census reports a test-file call to a verb entry outside the runner file with the file, the line, the enclosing declaration, and the entry | planned `TestVerbCallCensusReportsAnEntryCall` in internal/worktree/verb_call_census_test.go | A census that skips verb entries returns no report |
| VR36 | 25 | The census reports a test-file call to a joins form outside the runner file | planned `TestVerbCallCensusReportsAJoinsFormCall` in internal/worktree/verb_call_census_test.go | A census that knows only exported entries returns no report for the joins form |
| VR37 | 26 | The census reports a verb entry that a test function stores as a value without a call | planned `TestVerbCallCensusReportsAnEntryUsedAsAValue` in internal/worktree/verb_call_census_test.go | A census that reads only call expressions returns no report for the value |
| VR38 | 27 | The census reports a verb call inside a subtest closure | planned `TestVerbCallCensusReportsACallInsideASubtest` in internal/worktree/verb_call_census_test.go | A census that stops at a function literal returns no report |
| VR39 | 28 | The census reports a verb reference in a test-file package-level variable declaration | planned `TestVerbCallCensusReportsAPackageLevelReference` in internal/worktree/verb_call_census_test.go | A census that walks only function declarations returns no report |
| VR40 | 29 | The census reports nothing for verb references in `verb_runner_test.go` | planned `TestVerbCallCensusAllowsTheRunnerFile` in internal/worktree/verb_call_census_test.go | A census without the runner exemption reports the runner's own calls |
| VR41 | 30 | The census reports a call to an exported function that has an `args []string` parameter | planned `TestVerbCallCensusDerivesEntriesFromTheSignature` in internal/worktree/verb_call_census_test.go | A census with a fixed name list returns no report for a new entry in the synthetic set |
| VR42 | 30 | The census reports nothing for a call to an exported function without an `args []string` parameter | planned `TestVerbCallCensusIgnoresAnExportedHelper` in internal/worktree/verb_call_census_test.go | A census that treats every exported function as an entry reports the helper |
| VR43 | 31 | The census reports a call to the function that a synthetic verb entry calls with `defaultJoins()` first | planned `TestVerbCallCensusDerivesTheJoinsForm` in internal/worktree/verb_call_census_test.go | A census with a fixed joins form list returns no report for the synthetic joins form |
| VR44 | 32 | The census reports no line on the live package tree | planned `TestVerbCallCensusOnTheLiveTree` in internal/worktree/verb_call_census_test.go | A surviving direct verb call on the live tree is reported with its file and line |

Not covered: story 23 — the source rules out a line-count target, so no row grades a line count.
Not covered: story 33 — no ticket plans a production edit, and the bound registry paths sit in the fence for the preflight binding only.
Not covered: story 34 — no ticket plans an edit outside `internal/worktree` test files.

### Edge inventory

The census walks a build-tagged test file and skips a special file, through the shared `parseGoFiles` walk. It walks a closure, a subtest, a method of a test type, and a package-level declaration. An absent `verb_runner_test.go` leaves no file exempt, so the census reports every remaining verb call. These behaviors serve this repository only, because the worktree package is kit source.

The hostile-input checklist applies to the readers. The numeric-looking cell class is VR9. The control-byte, path, and patch-header classes do not apply, because the runner renders no output and reads no path from git output.

- Won't handle: the three `subshellAt` calls in `subshell_test.go` — that form takes a shell path and an environment, and the census still guards `Subshell`.
- Won't handle: a local identifier that shadows a verb form name — the census reports it by bare identifier, and the author renames the local.
- Won't handle: a verb entry whose argument parameter is not named `args` — every current verb entry uses `args`, and review catches a rename.
- Won't handle: `RunTreeChild` — its parameter is `argv`, it serves the tree verbs of other packages, and its worktree callers stay on the exec verb entry.
- Won't handle: direct verb calls in other packages' tests, for example `internal/treetarget` — the map puts them out of scope, and the entries stay exported.

## Ownership fences

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
| [1. Add the verb runner and move the reset family onto it](tickets/1-add-verb-runner.md) | none | The verb runner, the verb result, and the reset family on the runner |
| [2. Name the shared assignment and pool fixtures](tickets/2-name-shared-fixtures.md) | 1-add-verb-runner.md | The shared fixture builders return named values |
| [3. Move the cleanup verbs onto the verb runner](tickets/3-migrate-cleanup-verbs.md) | 2-name-shared-fixtures.md | Release, clean, resume-clean, and reclaim tests on the runner |
| [4. Move the query and create verbs onto the verb runner](tickets/4-migrate-query-and-create-verbs.md) | 3-migrate-cleanup-verbs.md | Create, list, path, show, exec, build, pool, and lease-file tests on the runner |
| [5. Move the merge and reauthorize verbs onto the verb runner](tickets/5-migrate-merge-and-reauthorize.md) | 4-migrate-query-and-create-verbs.md | Merge and reauthorize tests on the runner with named fixtures |
| [6. Name the landing fixtures](tickets/6-name-landing-fixtures.md) | 5-migrate-merge-and-reauthorize.md | The landing fixture builders return named values |
| [7. Move the landing composition tests onto the verb runner](tickets/7-migrate-landing-composition.md) | 6-name-landing-fixtures.md | Landing flag, refusal, and identity tests on the runner |
| [8. Move the landing effect and resume tests onto the verb runner](tickets/8-migrate-landing-effects.md) | 6-name-landing-fixtures.md | Landing effect, resume, census, and journey tests on the runner |
| [9. Refuse a direct verb call outside the verb runner](tickets/9-refuse-direct-verb-calls.md) | 7-migrate-landing-composition.md, 8-migrate-landing-effects.md | The census refusal and its live-tree pin |

## Out of scope

- The seam reduction spec from the same map is a separate capability with its own future spec: approximately 60 edits, 4 gate runs. It keeps 15 joins fields, reads each value once per exported function, and lowers the serial ceiling.
- A census refusal for positional fixture tuples is a separate capability that the map did not decide: approximately 3 edits, 1 gate run.

## Further notes

### Source trace

| Source clause | Coverage |
| --- | --- |
| Ticket 1: two specs, the verb runner spec first | This spec; Out of scope names the seam reduction spec |
| Ticket 5: delete the run wrappers and the fingerprint extractors | VR16, VR19, VR21, VR22, VR25, VR28 |
| Ticket 5: every direct verb call moves onto the verb runner | VR26, VR34, VR44 |
| Ticket 5: each positional fixture tuple becomes a named value in one helper file | VR17, VR18, VR20, VR23, VR24, VR27 |
| Ticket 5: the count pin holds | VR29 |
| Ticket 6: the worktree package only, on `axitest.DecodeDocument` | VR3, VR4, and story 34 |
| Ticket 7: the census reports a test-file call to a verb entry outside the runner file, with file and line | VR35, VR38, VR39, VR40, VR44 |
| Ticket 8: no size target; the acceptance is structural | Story 23 |
| Map destination: the runner accepts the kit root, the Bench home, and the clock as values | VR13, VR14, VR15 |
| Map term: the verb result carries the exit code, the rows, and the fingerprint | VR2, VR3, VR5, VR6 |

### Enumerations

A run wrapper is a test-file helper that calls a verb form and returns its output. The run wrapper list is `runCleanup`, `runCleanupWith`, `runDiscard`, `planAndApply`, `runReset`, `runResetWith`, `runMerge`, `runCreate`, `runResume`, `runResumeAt`, `landIn`, `execAtOwnedTarget`, `mustResumeClean`, and `mustReclaim`.

The fingerprint extractor list is `resetFingerprint`, `restoreFingerprint`, `reclaimFingerprint`, and `cleanupRowFingerprint`, plus the inline 64-hex match in `worktree_test.go`.

The per-ticket deletion commands run from the repository root:

- VR16: `rg -n '^func (runReset|runResetWith|resetFingerprint|restoreFingerprint)\(' internal/worktree`
- VR19: `rg -n '^func (runCleanup|runCleanupWith|runDiscard|planAndApply|runResume|runResumeAt|mustResumeClean|mustReclaim|cleanupRowFingerprint|reclaimFingerprint)\(' internal/worktree`
- VR21: `rg -n '^func (runCreate|execAtOwnedTarget)\(|\[0-9a-f\]\{64\}' internal/worktree --glob '*_test.go' --glob '!verb_runner_test.go'`
- VR22: `rg -n '^func runMerge\(' internal/worktree`
- VR25: `rg -n '^func landIn\(' internal/worktree`

The VR28 command is `rg -n '^func (runCleanup|runCleanupWith|runDiscard|planAndApply|runReset|runResetWith|runMerge|runCreate|runResume|runResumeAt|landIn|execAtOwnedTarget|mustResumeClean|mustReclaim|resetFingerprint|restoreFingerprint|reclaimFingerprint|cleanupRowFingerprint)\(|\[0-9a-f\]\{64\}' internal/worktree --glob '*_test.go' --glob '!verb_runner_test.go'`. A hit inside an expected-output literal that the test compares whole is not a reader; the Coverage axis names each such hit.

The verb form command for VR26 is `rg -n '\b(BuildCommand|ExecCommand|ExecCommandResolving|LandCommand|ResumeLandCommand|ListCommand|MergeCommand|PathCommand|ReclaimCommand|ReauthorizeCommand|ResetCommand|ResumeCleanCommand|ShowCommand|Subshell|CleanCommand|ReleaseCommand|CreateCommand|PoolCommand|LeaseFileCommand|buildWith|landWith|resumeLandWith|mergeWith|reauthorizeWith|resetWith|resumeCleanCommandWith|cleanCommandWith|releaseCommandWith)\(' <files>`. That name list is the census's derived set on the base tree. After ticket 9, the census replaces the command.

The tuple scan is a Go AST scan of the package's test files under the predicate in Implementation decisions. Ticket 1 records the scan program in its review record, and every later review reruns it. The run wrappers also match the predicate and go under VR28. At the decision-source commit, the scan lists these fixture builders:

- shared, ticket 2: `newOwnedAssignment`, `newPendingAssignment`, `newOwnedSubmoduleAssignment`, `newResidueGuardFixture`, `unprovableLandedAssignment`, `newReclaimPool`, `poolRootFixture`
- reset, ticket 1: `restoreFixture`
- cleanup, ticket 3: `landedSetFixture`, `retainedMemberFixture`, `removableSetFixture`, `refusedUnlandedRelease`
- merge and reauthorize, ticket 5: `mergeFixture`, `reauthorizeFixture`
- landing, ticket 6: `publicLandingFixture`, `publicLandingFixtureAtHome`, `specLessLandingFixture`, `foldedLandingFixture`, `ticketsOnlyLandingFixture`, `landingFixtureAtHome`
- landing, ticket 6: `redProspectiveGateLanding`, `brokerChangingLanding`, `brokerDestinationFixture`, `landSurface`, `foldLandingSibling`

### Reader sweep

- `worktreeTestCount` has one reader, `TestPackageTestCountPin`. Tickets 1 and 9 raise it.
- `worktreeSerialCeiling` has one reader, `TestSerialSetStaysBelowTheCeiling`. No ticket changes it.
- The run wrapper names have readers only in `internal/worktree` test files and in `capture/restructure-backlog.md`, which names `runCreate` as a split anchor. Ticket 4 updates that row.
- The joins forms have no reader outside `internal/worktree`. The verb entries have production readers in `cmd/bench/worktree_leaves.go`, `cmd/bench/main.go`, `internal/harness/worktree.go`, and `internal/sessioninspect/sessioninspect.go`, and one test reader in `internal/treetarget/run_test.go`. No reader changes, because no production signature changes.
- The map's `## Sources` name `parallelCensus` as a drift trigger. This spec adds a census file beside it and does not change `parallelCensus`.

### Pre-review proof checklist

- Cited symbols: `axitest.DecodeDocument` and `Document.Rows` in `internal/axi/axitest/document.go`; `parseGoFiles`, `parseTestFiles`, `parseSourceFiles`, `plantTestFiles`, `TestPackageTestCountPin`, and `TestSerialSetStaysBelowTheCeiling` in `internal/worktree/parallel_census_test.go`; `defaultJoins` in `internal/worktree/joins.go`; `ExecCommandResolving` in `internal/worktree/exec.go`; `toon.Table` in `internal/toon/toon.go`.
- Import edges: `internal/worktree` tests already import `internal/axi/axitest` in `list_selected_test.go` and `path_identifier_test.go`. The package already imports `internal/toon`. No new edge.
- Source-row clauses and occurrences: the Source trace table quotes each clause; the map is the only occurrence.
- Promised field labels: `fingerprint` as a table field; `fingerprint=` as a record cell, emitted by `reset.go`.
- Changed-function callers: no production function changes. The deleted test helpers' callers are the test files in the fence.
- Copy survival: VR28 fails when one run wrapper or one fingerprint extractor survives. VR44 fails when one direct verb call survives.
- Rendered-shape readers: none; no rendered output shape changes.

### Probes run during authoring

- A probe test fed `axitest.DecodeDocument` four shapes. A `reset_plan{...}` record and a `merged{...}` record decode as a string, not an object, so the decode returns an error. A lane line before a record fails with a missing colon. A `cleanup[1]{...}:` table decodes with its `fingerprint` cell. Thus the fingerprint reader needs the record branch, and the rows reader returns an error on record output.
- The joins value still holds 29 fields, and `parallel_census_test.go` still pins 664 tests and a serial ceiling of 46. No source file of the map changed between the asset commit `8ab52861` and the base `6d0c1e79`.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"VR-C1","tickets":["1-add-verb-runner.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"}]},{"id":"VR-C2","tickets":["2-name-shared-fixtures.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"}]},{"id":"VR-C3","tickets":["3-migrate-cleanup-verbs.md","4-migrate-query-and-create-verbs.md","5-migrate-merge-and-reauthorize.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"}]},{"id":"VR-C4","tickets":["6-name-landing-fixtures.md","7-migrate-landing-composition.md","8-migrate-landing-effects.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"}]},{"id":"VR-C5","tickets":["9-refuse-direct-verb-calls.md"],"verification":[{"id":"worktree","command":"bench test --package ./internal/worktree"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/worktree-verb-runner/spec.md"},{"id":"worktree","command":"bench test --package ./internal/worktree"}]}
```

### Contradiction flagged for reviewer veto

Ticket 1 of the map says that the seam reduction spec rewrites each stubbed test once, onto the verb runner. Ticket 5 says that every direct verb call moves onto the verb runner in this spec, and a stubbed test calls a joins form directly. This spec follows ticket 5: a stubbed test moves onto the runner with its joins value now. The seam reduction spec then replaces the stub with a real fixture once, without a second runner rewrite. This reading is non-behavioral.

### Flagged additions

- The census also refuses a joins form, not only a verb entry. Without it, a test can call a joins form with its own buffers, which is the case ticket 7 closes.
- The census also refuses a verb entry used as a value and a package-level reference, so a table cannot bypass the runner.
- The runner fails a call that carries a kit value or a clock value until the seam reduction spec wires them.
- The readers return an error instead of a failed test, so their error paths have a test seam.
- VR31 and VR32 add a differential run per chunk, because the count pin cannot see a test that turns into a skip.

### Map handling

This spec compiled the map: the map, its tickets, and its asset moved from `decisions/` into `specs/worktree-verb-runner/decisions/`. The seam reduction spec reads the same compiled map in place. This spec carries no `Roadmap:` line, because `FT356` stays open until the seam reduction spec lands. One map source names `kitRoot` under `internal/gate/kit_source.go`; the function lives in `internal/gate/phases.go`. That drift touches only the seam reduction spec.

The promote-then-delete retirement of a spec removes its compiled topic folders. The retirement of this spec must therefore keep `specs/worktree-verb-runner/decisions/` until the seam reduction spec stages and moves or reads the map. The reviewer decides that order at this spec's retirement.
