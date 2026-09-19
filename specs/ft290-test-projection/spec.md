# `bench test` projection faces

Status: staged

Roadmap: FT290

Decision source: `specs/ft290-test-projection/decisions/ft290-test-projection.md` (ready compiled map).

Verification log: 0 iteration(s) to accept — the review round has not run.

## Problem

An agent that runs `bench test` cannot read some facts from the output. A green
named check does not show that the check ran. A green package looks the same as
a package that ran no test. `bench test --check prose` prints nothing on green.
The failures table shows one diagnostic line for each failed test, so the agent
pays a second run with `--full`. One red in the system suite costs a whole
re-run, because `--check` refuses `--run`.

The unknown-check refusal lists a check set with no owner, so a stale
executable gives a stale list with no attribution. `bench test --changed`
widens its selection and gives no cause. The check inventory and the fixtures
of a check have no projection, so an agent reads test source for them.

## Solution

Each named-check result starts with one `check` row that names the check, its
kind, and its counts. A named check that ran nothing exits 1. Each packages row
carries a `tests_run` cell. The prose check prints its `check` row on green, and
`--full` lists the graded subjects.

The failures table carries a `lines` count, and `--full` prints one row for each
diagnostic line. `--check system` accepts `--run`. The unknown-check refusal
names the running executable and the source digest of its seal.

`bench test --checks` lists the named checks. `bench test --check <name>
--fixtures` lists the fixtures that the check owns. Each `--changed` packages
row carries a `selected_by` cell with one cause.

## User stories

Line: opus / medium.
Implementation-line reason: TP-C1a is the hardest chunk, because its header changes red each test that pins the old tables. The spec fixes each output byte, each seam exists, and the package tests cover each row.
Harder chunks: TP-C1a.

Result evidence:

1. As an agent, I want a `tests_run` cell on each packages row, so that a green package differs from a no-test package.
2. As an agent, I want a `check` row at the start of each named-check result, so that the result names the check that ran.
3. As an agent, I want a named check that ran nothing to exit 1, so that no execution never reads as evidence.
4. As an agent, I want a package run with no test to keep its exit code, so that the zero rule stays with named checks.
5. As an agent, I want the prose check to print its `check` row on green, so that a green prose run is not silent.
6. As an agent, I want `--full` to list each graded prose subject, so that I can see which files the check graded.
7. As an agent, I want a prose run with zero subjects to exit 1, so that an empty grade never reads as a pass.
8. As an agent, I want a red prose run to keep each finding after the `check` row, so that I lose no finding.
9. As an agent, I want a prose grader refusal to print its diagnostic after the `check` row, so that I can repair the policy file.
10. As an agent, I want a `lines` count on each failures row, so that I know when `--full` holds more diagnostics.
11. As an agent, I want `--full` to print one failures row for each diagnostic line, so that no line hides in a joined cell.
12. As a `bench probe` caller, I want the failed-test count to stay a count of tests under `--full`, so that the probe verdict row stays correct.

Run filter and refusal identity:

13. As a delegate, I want `bench test --check system --run <regex>`, so that I read one system red without a whole re-run.
14. As a delegate, I want a system run pattern with no match to exit 1, so that a wrong pattern never reads as green.
15. As a maintainer, I want each other named check to keep the `--run` refusal, so that a scoped root test cannot be filtered to nothing.
16. As an agent, I want the unknown-check refusal to name the running executable, so that I know which executable owns the listed check set.
17. As an agent, I want the refusal to name the source digest of the seal, so that I can compare it with the current source.
18. As an agent, I want the refusal to print `unsealed` when the seal is unreadable, so that the refusal never fails on a missing seal.

Inventory faces:

19. As a charge author, I want `bench test --check <name> --fixtures`, so that I read the fixtures of a check from the CLI.
20. As a charge author, I want a `CHECK` marker to decide the fixture owner, so that the rows agree with the canary inventory.
21. As an agent, I want a check with no fixture to print an empty table at exit 0, so that absence is definite.
22. As an agent, I want the inventory faces to start no Go child, so that an inventory read costs no build.
23. As an agent in a linked repository with no `tests/canary` directory, I want an empty answer at exit 0, not a fault.
24. As a maintainer, I want an invalid canary inventory to refuse with its diagnostic, so that a wrong inventory never reads as an empty one.
25. As an agent, I want `bench test --checks` at exit 0, so that I read the inventory without an unknown-check refusal.
26. As an agent, I want each `--checks` row to give the kind and the family count, so that I know which checks own fixtures.
27. As a maintainer, I want the help list and the `--checks` table to read one source, so that the two never disagree.

Cause cell:

28. As a delegate, I want a `selected_by` cell on each `--changed` packages row, so that I know why the run holds each package.
29. As a delegate, I want one cause by a fixed precedence, so that a package with two causes prints one stable value.
30. As a delegate, I want an `imports` cause to name one selected dependency, so that I can follow the chain to the change.
31. As a delegate, I want a Go metadata change to print `go-metadata` on each row, so that I know why the whole set runs.
32. As an agent, I want the other forms to print no `selected_by` cell, so that the cell appears only where a selection exists.

Grammar:

33. As an agent, I want each invalid flag combination to refuse with usage at exit 2, so that no partial form runs.
34. As an agent, I want `bench help` and `bench test --help` to show each new form, so that I find the forms without a trial.

Reviewed exclusions:

35. As the reviewer, I want the compile error on a non-compiling `--package` kept out, so that the spec adds no row for shipped behavior.
36. As the reviewer, I want the run binary provenance kept out, so that its research closes before its projection is decided.

## Implementation decisions

### Late answers from the reviewer, 2026-09-19

- The unknown-check refusal names the running executable path and the source digest of its seal. `freshness.SealDigests` supplies the digest. The refusal prints `unsealed` when the seal is unreadable. It names no source commit, because a seal and an executable record none.
- The check row has one schema, `check[1]{name,kind,tests_run,subjects}`. The prose check prints `tests_run` 0 and its real `subjects` count. A Go-backed check prints `subjects` 0. The zero rule reads `subjects` for the prose check and `tests_run` for each other check.
- The `--fixtures` rows and the `--checks` family count use the one owner that `canary.Fixtures` resolves for each fixture.

### The check row and the zero rule

The `check` row is the first block of each named-check result that reached a
verdict. The `kind` cell is `conformance`, `system`, or `prose`. One owner maps
a check name to its kind, and the `--checks` table reads the same owner.

The `tests_run` count is the count of distinct tests and subtests that emitted a
run event. The report already holds this count for `Outcome.Ran`. A packages
row prints the same count for its one package.

When the zero rule fires, the result prints the `check` row and then one error
line with the title `named check ran nothing`. The exit code is 1. A refusal
that comes before a verdict keeps its current bytes and prints no `check` row.
When `--run` is present and no test ran, the current `go test reported no test
runs` refusal wins, and the zero-rule line does not print. When a package of a
named check does not compile, the outcome is a build failure. The compile
diagnostic prints at exit 1, and the zero-rule line does not print.

### The prose result

The prose grader gives the count of graded subjects and their paths beside its
findings. A green result prints only the `check` row. `--full` adds a
`subjects[N]{path}` table in sorted order. A red result prints the `check` row
and then each finding line as it prints today. A grader refusal prints the
`check` row with `subjects` 0, then the refusal diagnostics, and exits 1.

### The failures table

The header is `failures[N]{package,test,line,lines}` in each mode. The `lines`
cell is the count of diagnostic lines that the test emitted. The default mode
prints one row for each failed test with its first line. `--full` prints one
row for each diagnostic line in emitted order, and each row carries the same
`lines` count. A failed test with no diagnostic prints one row with
`no diagnostic emitted` and `lines` 0 in each mode. A package failure with no
test name obeys the same rules over its package log.

`Outcome.FailedTests` counts distinct failed tests. It does not count rows.

### The system run pattern

`--check system --run <regex>` appends `-run <regex>` after the system suite
operands. `Request.Run()` returns the pattern for that request. The grammar
refuses `--run` with each other named check before a Go child starts.

### The refusal identity

The refusal keeps exit 2 and this line order: the `unknown check: <name>` line,
an `executable: <path>` line, a `seal: <value>` line, and then the check list.
The path is the absolute path of the running executable. One package variable
supplies it, so a test in the same process can set it. The `seal` value is the
source digest of the seal, or `unsealed`. Control characters in the path and in
the caller's check name print escaped through `sanitize.Controls`.

### The inventory faces

`--fixtures` reads `tests/canary` under the graded repository root, the same
root that `bench canary` reads. The header is
`fixtures[N]{family,fixture,path}`. The rows hold each fixture whose inventory
owner is the named check, sorted by path. The `path` cell is relative to the
repository root. A fixture that sits directly under `tests/canary` prints an
empty `family` cell.

An absent or empty `tests/canary` directory gives an empty table at exit 0.
`canary.Fixtures` returns one error for both states today, and only an
unexported message names it. Ticket 8 exports one sentinel error from
`internal/canary`, and `Fixtures` and `FixturePins` both use it. The face reads
the sentinel through `errors.Is`. Each other inventory error refuses at exit 1
with the inventory diagnostic. The face
selects no run binary and starts no Go child.

`--checks` prints `checks[N]{name,kind,families}` in the order of the help
list. The `families` cell counts the distinct family names in which the check
owns one fixture or more. A fixture with no family adds nothing to the count.
An absent `tests/canary` directory gives `families` 0 on each row.

### The cause cell

The changed-package selector returns each selected package with one cause. The
precedence is `go-metadata`, then `changed`, then `embed`, then
`imports <package>`. For `imports`, the cell names the first selected direct
dependency in sorted import-path order. The direct dependencies are the
imports, the test imports, and the external test imports. Only a `--changed`
result prints the `selected_by` cell, and an empty selection prints the same
header with zero rows.

### The structure budget

The growth lane reds a file that is over its line budget and that gains a line.
These fenced files are over the budget at `bcc5543f`: `internal/testreport/command.go`,
`check_test.go`, `selection_test.go`, and `testreport_test.go` of
`internal/testreport`, `cmd/bench/main.go`, and `cmd/bench/command_registry_test.go`.
An edit to one of them keeps or lowers its line count. `internal/canary/inventory.go`
is also over the budget, so ticket 8 does not grow it.

Ticket 1 moves the named-check owner out of `command.go` into a new file, with
no behavior change. Each new test goes in a new test file. An existing test
that needs more lines moves to a new file in the same ticket.

The `internal/testreport/` directory is already over its file-count budget, and
six tickets add a file to it. The growth lane does not grade that count, so
the directory debt grows and stays soft.

### The grammar

The usage line and the `bench help` row both become this text:

`bench test [--full] [--package <expr> | <legacy-package> | --changed] [--base <commit> [--source-tip <commit>]] [--run <go-regex>] | bench test [--full] --check <name> | bench test [--full] --check system --run <go-regex> | bench test --check <name> --fixtures | bench test --checks`

`--checks` accepts no other flag and no operand. `--fixtures` requires
`--check` and accepts no other flag. An unknown check with `--fixtures` gives
the unknown-check refusal.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| TP-C1a / `1-split-named-check-owner.md`, `2-count-tests-run.md`, `3-prove-named-check-ran.md` | A result proves what ran, and the later tickets get the packages header and the `check` row producer. | TP1, TP2, TP3, TP4, TP5, TP6, TP18, TP47, TP51, TP55 | `bench test --package ./internal/testreport` | yes |
| TP-C1b / `4-print-prose-check-result.md`, `5-show-each-failure-diagnostic.md` | The prose check prints a result, and the failures table shows each diagnostic. | TP7, TP8, TP9, TP10, TP11, TP12, TP13, TP14, TP15, TP16, TP17, TP49 | `bench test --package ./internal/testreport`, `bench test --package ./internal/prose`, `bench test --package ./internal/probe` | no |
| TP-C2 / `6-filter-system-suite.md`, `7-name-running-executable.md` | A delegate filters the system suite, and the refusal has an owner. | TP19, TP20, TP21, TP22, TP23, TP24, TP25, TP26, TP54 | `bench test --package ./internal/testreport`, `bench test --package ./cmd/bench` | no |
| TP-C3 / `8-list-check-fixtures.md`, `9-list-check-inventory.md` | The CLI lists the fixtures of a check and the check inventory. | TP27, TP28, TP29, TP30, TP31, TP32, TP33, TP34, TP35, TP36, TP37, TP38, TP39, TP40, TP50, TP52, TP53 | `bench test --package ./internal/testreport`, `bench test --package ./cmd/bench` | no |
| TP-C4 / `10-explain-changed-selection.md` | Each `--changed` row names its cause. | TP41, TP42, TP43, TP44, TP45, TP46, TP48 | `bench test --package ./internal/testreport` | no |

Stable chunk IDs: the first review round split TP-C1 into TP-C1a and TP-C1b. Tickets 2 and 3 create the header and the `check` row that tickets 4 and 5 consume, so their chunk review closes first. TP-C2, TP-C3, and TP-C4 keep their IDs.

## Testing decisions

- A good test calls `testreport.Command` and compares the printed bytes and the exit code. It does not read the report struct.
- The one seam is the existing `Command` seam. A test puts a canned `go` script first on `PATH` and installs a run binary factory. `TestNamedCheckOwnsConformanceEnvironment` in `internal/testreport/check_test.go` is the prior art, and it was read in this session.
- Each row that compares rendered tables uses canned events, because two real runs differ in elapsed time.
- TP18 is the one composition row. It extends `TestNamedCheckRunsOnlyRegisteredDevScope`, which runs the real Go child and the real root conformance test.
- The package tests of `internal/testreport` run in the gate's `test` phase, so that phase observes each row.

### Posture change: tests that the new output reds

The zero rule reds each named-check test whose canned `go` script emits no run
event and expects exit 0. The fixture builders are `writeCheckGo` and the
inline script of `TestNamedCheckRunsFromKitAgainstLinkedConsumer` in
`internal/testreport/check_test.go`. `check_test.go` calls `writeCheckGo` at seven lines.

| line | test | expected exit |
| --- | --- | --- |
| 26 | `TestNamedCheckOwnsConformanceEnvironment` | 0 |
| 156 | `TestFocusedRequestGrammarRefusals` | 2 |
| 182 | `TestNamedCheckRefusalMatrix` | 2 |
| 228 | `TestNamedCheckRefusesCorruptInheritedSelection` | 1 |
| 293 | `TestNamedChecksWriteNoGateOwnedRecords` | 0 |
| 352 | `TestSystemCheckOwnsTheGateEnvironment` | 0 |
| 405 | `TestSystemCheckRefusesAForeignRoot` | 1 |

`internal/testreport/testreport_test.go` line 294 is one more caller, and its
`--check system` run expects exit 0. The three exit 0 rows, that caller, and
the inline script red under the zero rule. Ticket 3 adds a run event to
`writeCheckGo` and to the inline script, inside their current lines.

The header changes red the exact header matches in `outcome_test.go`,
`testreport_test.go`, `selection_test.go`, `cancel_test.go`, and
`check_test.go` of `internal/testreport`. `TestUnknownNamedCheckReportsOperandAndInventory`
compares the whole refusal, so ticket 7 moves the test and rewrites its expectation.
`internal/testreport/selection_facts_test.go` line 22 expects `AllTests` for
the system check, and ticket 6 adds the pattern case beside it.

### Seam diagram

    trigger: `bench test <args>`, `bench probe`, the commit lane, `bench worktree merge`
        │
        ▼
    args, tree, canned or real Go events  ──▶  [ testreport.Command ]  ──▶  TOON bytes, exit code
                      ◀ tests attach here: a canned `go` on PATH, an installed run binary
                        factory, a temporary tree, and a byte comparison of the output

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| TP1 | 1 | A package with two run events prints `tests_run` 2 under the header `packages[1]{package,status,elapsed_ms,tests_run}` | planned TestPackagesRowCountsRunEvents in internal/testreport, through `Command` with canned events | The old header has no `tests_run` cell, so the exact header match is red. |
| TP2 | 1 | A package with the `no-tests` status prints `tests_run` 0 | planned TestPackagesRowNoTestsCountsZero in internal/testreport, through `Command` with canned events | A count that reads the package status in place of the run events prints a wrong value. |
| TP3 | 2 | `--check line-routing` with one run event prints the first block `check[1]{name,kind,tests_run,subjects}:` with the row `line-routing,conformance,1,0` | planned TestNamedCheckPrintsCheckRowFirst in internal/testreport, through `Command` with a canned `go` | A result with no `check` row, or with the row after the packages table, fails the prefix match. |
| TP4 | 2 | `--check system` prints the `kind` cell `system` | planned TestSystemCheckRowKind in internal/testreport, through `Command` with a canned `go` | One fixed kind for each Go-backed check prints `conformance`. |
| TP5 | 3 | `--check line-routing` with a package pass and no run event exits 1 and prints `tests_run` 0 and the title `named check ran nothing` | planned TestNamedCheckRanNothingExitsOne in internal/testreport, through `Command` with a canned `go` | The current code exits 0 for this input. |
| TP51 | 3 | `--check line-routing` with a `build-fail` event and no run event exits 1, prints the compile diagnostic in the failures table, and does not print the title `named check ran nothing` | planned TestNamedCheckBuildFailureWinsOverZeroRule in internal/testreport, through `Command` with canned events | A zero rule that reads only the count prints its title over the build failure. |
| TP6 | 4 | `--package chosen` with a package pass and no run event exits 0 | planned TestPackageRunWithNoTestKeepsExitZero in internal/testreport, through `Command` with a canned `go` | A zero rule on each form changes this exit code to 1. |
| TP55 | 4 | `--changed` over one changed Go file with a package pass and no run event exits 0 | planned TestChangedRunWithNoTestKeepsExitZero in internal/testreport, through `Command` with a canned `go list` and `go test` | `bench worktree merge` folds this exit code into its retry, and a zero rule on each form makes it 1. |
| TP7 | 5 | A green prose run over two subjects prints exactly the row `prose,prose,0,2` under the `check` header, with no `packages[` text, at exit 0 | planned TestProseGreenPrintsOnlyCheckRow in internal/testreport, through `Command` over a temporary tree | The current output is empty, and a packages table fails the exact match. |
| TP8 | 6 | A green prose run with `--full` adds `subjects[2]{path}:` with the two paths in sorted order | planned TestProseFullListsSubjects in internal/testreport, through `Command` over a temporary tree | A result that ignores `--full` prints no `subjects` table. |
| TP9 | 7 | A prose run over a tree with zero graded subjects exits 1 and prints `subjects` 0 and the title `named check ran nothing` | planned TestProseZeroSubjectsExitsOne in internal/testreport, through `Command` over a temporary tree | The current code returns a pass for zero subjects. |
| TP10 | 8 | A red prose run over two subjects prints the row `prose,prose,0,2` and then each finding line, at exit 1 | planned TestProseRedKeepsFindingsAfterCheckRow in internal/testreport, through `Command` over a temporary tree | A red path that skips the `check` row, or that drops a finding, fails the match. |
| TP11 | 9 | A prose run over a tree with no `.bench/prose-exclusions` file prints `subjects` 0 and then the grader diagnostic at exit 1, and does not print the title `named check ran nothing` | planned TestProseGraderRefusalPrintsCheckRow in internal/testreport, through `Command` over a temporary tree | A refusal path with no `check` row fails the prefix match, and a zero rule that reads only `subjects` prints its title over the diagnostic. |
| TP12 | 6 | The prose grader returns the graded subject paths beside its findings for a tree with one excluded file and one graded file | planned TestGradeReportsGradedSubjects in internal/prose, through the prose grade entry | A count that includes the excluded file returns two paths. |
| TP13 | 10 | A failed test with three diagnostic lines prints one default row with its first line and `lines` 3 under `failures[1]{package,test,line,lines}` | planned TestFailuresRowCountsLines in internal/testreport, through `Command` with canned events | The old header has no `lines` cell. |
| TP14 | 10 | A failed test with no diagnostic prints `no diagnostic emitted` and `lines` 0 | planned TestFailuresRowWithNoDiagnosticCountsZero in internal/testreport, through `Command` with canned events | A count that reads the printed cell gives 1. |
| TP15 | 11 | `--full` prints three failures rows for that test in emitted order, and each row carries `lines` 3 | planned TestFullFailuresPrintOneRowPerLine in internal/testreport, through `Command` with canned events | The current code prints one row with a joined cell. |
| TP16 | 11 | No `--full` failures cell holds an escaped newline | planned TestFullFailuresHoldNoJoinedCell in internal/testreport, through `Command` with canned events | The old joined cell survives beside the new rows. |
| TP17 | 11 | A package failure with no test name and two package log lines prints two `--full` rows with an empty `test` cell | planned TestFullPackageFailurePrintsOneRowPerLogLine in internal/testreport, through `Command` with canned events | A row split that handles only named tests leaves the package log joined. |
| TP18 | 2 | The real `--check ordinary-build-census` run prints `tests_run` 1 in its `check` row | TestNamedCheckRunsOnlyRegisteredDevScope in internal/testreport, through the real Go child, with one assertion changed in place | A canned event stream can differ from the real Go child, and only the real child proves the count. |
| TP49 | 12 | A `--full` run with one failed test and three diagnostic lines gives `Outcome.FailedTests` 1 | planned TestFullFailedTestsCountsTests in internal/testreport, through `Execute` with canned events | A count of failures rows gives 3. |
| TP19 | 13 | `--check system --run ^TestX$` starts the Go child with `-run ^TestX$` after the system suite operands | planned TestSystemRunPatternReachesGoArgv in internal/testreport, through `Command` with a canned `go` that records its argv | The current grammar refuses the form at exit 2. |
| TP20 | 13 | `Request.Run()` returns `^TestX$` for that request | planned TestSystemRequestRunFact in internal/testreport, through `Prepare` | The current switch returns `AllTests` for each system request. |
| TP21 | 14 | `--check system --run ^TestNone$` with no run event exits 1 with the title `go test reported no test runs` and without the title `named check ran nothing` | planned TestSystemRunPatternNoMatchRefusalWins in internal/testreport, through `Command` with a canned `go` | Two refusals compete, and the wrong winner prints the zero-rule title. |
| TP22 | 15 | `--check line-routing --run ^TestX$` exits 2 with usage and starts no Go child | TestNamedCheckRefusalMatrix in internal/testreport, through the existing matrix, read in this session | A grammar that accepts `--run` for each check starts the child. |
| TP23 | 15 | `--check prose --run ^TestX$` exits 2 with usage | planned TestProseRefusesRunPattern in internal/testreport, through `Command` | The prose path returns before the Go path, so a late refusal never fires for it. |
| TP24 | 16 | `--check not-registered` prints `unknown check: not-registered`, then `executable: <path>` from the package variable, then the `seal:` line, then the check list, at exit 2 | TestUnknownNamedCheckReportsOperandAndInventory in internal/testreport, through the whole-output match, moved to a new test file | The whole-output match fails when a line is absent or out of order. |
| TP25 | 17 | An executable with a regular seal file beside it prints `seal: <source digest>` with the `sources` value of that file | planned TestUnknownCheckNamesSealSources in internal/testreport, through `Command` with a temporary executable and a seal file | A refusal that prints the executable digest, or a fixed word, fails the match. |
| TP26 | 18 | An executable with no seal file prints `seal: unsealed` at exit 2 | planned TestUnknownCheckPrintsUnsealed in internal/testreport, through `Command` with a temporary executable | A seal read error that becomes the refusal changes the text and the exit code. |
| TP54 | 16 | `--check` with a name that holds U+0001 prints the escaped name in the `unknown check:` line and no raw U+0001 byte | planned TestUnknownCheckEscapesName in internal/testreport, through `Command` | The current code prints the caller's bytes unchanged. |
| TP27 | 19 | `--check package-core-guard --fixtures` over a tree with fixtures `a` and `b` in that family prints `fixtures[2]{family,fixture,path}:` with `package-core-guard,a,tests/canary/package-core-guard/a` first | planned TestFixturesFaceListsOwnedFixtures in internal/testreport, through `Command` over a temporary canary tree | The form does not exist today, and an absolute path fails the exact row. |
| TP28 | 20 | A fixture in the `package-core-guard` family with a `CHECK` file that names `default-branch-single-source` prints under `--check default-branch-single-source --fixtures` | planned TestFixturesFaceHonorsCheckMarker in internal/testreport, through `Command` over a temporary canary tree | A filter by `registry.CanaryFamilies` gives this check zero rows. |
| TP29 | 20 | The same fixture does not print under `--check package-core-guard --fixtures` | planned TestFixturesFaceOmitsReassignedFixture in internal/testreport, through `Command` over a temporary canary tree | A filter by family prints the fixture under the family owner. |
| TP30 | 21 | `--check system --fixtures` prints `fixtures[0]{family,fixture,path}:` at exit 0 | planned TestFixturesFaceEmptyForSystem in internal/testreport, through `Command` over a temporary canary tree | A silent result, or exit 1, fails the match. |
| TP31 | 22 | `--fixtures` and `--checks` write no canned `go` marker and call no run binary builder | planned TestInventoryFacesStartNoChild in internal/testreport, through `Command` with a canned `go` and a counting factory | A face that goes through the run path builds a run binary first. |
| TP32 | 23 | A tree with no `tests/canary` directory prints the empty fixtures table at exit 0, and `--checks` prints `families` 0 on each row | planned TestInventoryFacesAbsentCanaryDirectory in internal/testreport, through `Command` over a temporary tree | `canary.Fixtures` returns an error for an absent directory, so a face that ignores the sentinel exits 1. |
| TP53 | 23 | A present and empty `tests/canary` directory prints the empty fixtures table at exit 0 | planned TestFixturesFaceEmptyCanaryDirectory in internal/testreport, through `Command` over a temporary tree | A directory check that tests only for absence sends the empty directory to the exit 1 branch. |
| TP33 | 24 | A fixture whose `CHECK` file names `no-such-check` makes `--fixtures` exit 1 with the text `names unknown check` | planned TestFixturesFaceRefusesInvalidInventory in internal/testreport, through `Command` over a temporary canary tree | An error that reads as absent prints an empty table at exit 0. |
| TP52 | 24 | An owned fixture whose directory name holds U+0001 makes `--fixtures` exit 1 with the `toon.RenderError` text | planned TestFixturesFaceRefusesUnprintableName in internal/testreport, through `Command` over a temporary canary tree | A face that drops the encoder error prints a partial table at exit 0. |
| TP34 | 25, 27 | `--checks` prints `checks[N]{name,kind,families}:` at exit 0, and its `name` cells equal the help check list in order | planned TestChecksFaceEqualsHelpList in internal/testreport, through `Command`, and the one-source half is review-owned | A table that omits `system` or `prose` differs from the help list. Review confirms that the producer calls `namedChecks()`. |
| TP35 | 26 | The `--checks` rows print `conformance` for `line-routing`, `system` for `system`, and `prose` for `prose` | planned TestChecksFaceKinds in internal/testreport, through `Command` | One fixed kind fails two of the three cells. |
| TP36 | 26 | A check that owns fixtures in the families `package-core-guard` and `guard-classifier-table` prints `families` 2 | planned TestChecksFaceCountsFamilies in internal/testreport, through `Command` over a temporary canary tree | A count of fixtures prints a larger value. |
| TP37 | 26 | A check that owns only one `CHECK`-marked fixture prints `families` 1 | planned TestChecksFaceCountsMarkedFixtureFamily in internal/testreport, through `Command` over a temporary canary tree | A count from `registry.CanaryFamilies` prints 0. |
| TP38 | 26 | A fixture directly under `tests/canary` prints an empty `family` cell and adds nothing to `families` | planned TestChecksFaceIgnoresEmptyFamily in internal/testreport, through `Command` over a temporary canary tree | A count of distinct values that includes the empty name prints 1. |
| TP39 | 33 | Each of `--checks --full`, `--checks --check prose`, `--fixtures`, `--check prose --fixtures --full`, and `--check system --fixtures --run ^X$` exits 2 with usage | planned TestInventoryGrammarRefusals in internal/testreport, through `Command` with a canned `go` marker | A parser that ignores the extra flag runs a partial form. |
| TP40 | 33 | `--check not-registered --fixtures` gives the unknown-check refusal at exit 2 | planned TestFixturesFaceUnknownCheck in internal/testreport, through `Command` | An empty table for an unknown name hides a typing error. |
| TP50 | 34 | `bench help` and `bench test --help` each hold the exact grammar text of this spec | TestTestHelpNamesOnlyRunnableFocusedForms in cmd/bench, through the help render and the help inventory golden | An old help row omits the new forms. |
| TP41 | 28 | A `--changed` run over one changed Go file prints `packages[1]{package,status,elapsed_ms,tests_run,selected_by}:` with the cause `changed` | planned TestChangedRowsCarrySelectedBy in internal/testreport, through `Command` with a canned `go list` and `go test` | The current header has no `selected_by` cell. |
| TP42 | 29 | A package that changed and that imports a changed package prints `changed` | planned TestCauseChangedBeatsImports in internal/testreport, through the changed-package selector with a canned loader | A cause that the closure loop writes last prints `imports`. |
| TP43 | 29 | A package with only a changed embed file prints `embed` | planned TestCauseEmbedOnly in internal/testreport, through the changed-package selector with a canned loader | A selector that folds embed into `changed` fails the cell. |
| TP44 | 30 | A package that imports the selected packages `example/a` and `example/b` prints `imports example/a` | planned TestCauseImportsNamesFirstSortedDependency in internal/testreport, through the changed-package selector with a canned loader | Map iteration order gives `example/b` on some runs. |
| TP45 | 30 | A package that imports only `example/c`, where `example/c` imports the changed `example/a`, prints `imports example/c` | planned TestCauseImportsNamesDirectDependency in internal/testreport, through the changed-package selector with a canned loader | A cause that names the root change prints `imports example/a`. |
| TP46 | 31 | A changed `go.mod` and a changed Go file print `go-metadata` on each row, the changed package included | planned TestCauseGoMetadataWins in internal/testreport, through the changed-package selector with a canned loader | A precedence with `changed` first prints `changed` on one row. |
| TP47 | 32 | A `--package` run prints the header `packages[1]{package,status,elapsed_ms,tests_run}` with no `selected_by` text | planned TestPackageFormHasNoSelectedBy in internal/testreport, through `Command` with canned events | One header for each form prints an empty cause cell. |
| TP48 | 28 | A `--changed` run with no selected package prints `packages[0]{package,status,elapsed_ms,tests_run,selected_by}:` at exit 0 | TestChangedNonGoSubjectRendersExplicitEmpty in internal/testreport, through the existing empty-selection test, with its header expectation changed in place | An empty report that keeps the old header differs from the populated one. |

Not covered: story 35 — commit `4ff47076` shipped the behavior, and `TestRunPatternReportsCompilerDiagnostic` already grades it.
Not covered: story 36 — `decisions/run-binary-provenance.md` owns the outcome, and its research is open.

### Edge inventory

Audience: each behavior serves every repository that links the kit, except the
inventory faces. Those faces serve this repository, because a linked repository
holds no conformance fixtures. TP32 fixes the linked-repository answer.

Handled edges, each with its row:

- The absent `tests/canary` directory: TP32.
- The present and empty `tests/canary` directory: TP53.
- The named check whose package does not compile: TP51.
- The `--changed` run with no run event: TP55.
- The invalid canary inventory: TP33.
- The unreadable seal: TP26.
- The two competing refusals for a system run pattern: TP21.
- The grader refusal in the prose check: TP11.
- The failed test with no diagnostic: TP14.
- The fixture with no family: TP38.

A linked repository with no Markdown subject gets exit 1 from
`bench test --check prose`, by TP9. `prose.Grade` keeps its pass for zero
subjects, so its conformance caller does not change. The commit lane grades
prose through `bench gate-prose`, not through this verb.

Tests that swap a package variable: the run binary selector `selectRunBinary`,
and the new running-executable variable. Each test calls `Command` in its own
process, so each swap reaches the code under test.

Hostile input, shell CLI surface: a control character in the executable path
prints escaped. `canary.Fixtures` does not call `canary.Select`, so it does not refuse a control
character in a fixture name. The `toon.Table` encoder refuses a cell below
U+0020, and the face prints `toon.RenderError` at exit 1: TP52. U+007F passes
`toon.Representable` and prints as it is. The caller's check name in the
unknown-check refusal prints escaped: TP54.
An operand after `--checks` refuses through TP39's branch.

**Won't handle:**

- A `--run` filter for a named check other than `system` — each such check runs one root test with no subtest. TP22 and TP23 keep the refusal for the surviving callers.
- The pinned repository paths of a fixture in the `--fixtures` rows — `canary.FixturePins` keeps that answer for the build preflight, which is the surviving caller.
- A source commit in the refusal — no seal and no executable records one, and TP25 gives the digest that the freshness check compares.
- An executed-check list from the conformance timing file in the `check` row — a wrong scope already reds the root test. TP18 keeps the real run as the surviving proof.
- The project-owned fixtures of a linked repository in `--fixtures` — their owner is `project:<family>`, which is not a named check, and `bench canary` stays their caller.

## Ownership fences

- `internal/testreport/`
- `internal/prose/walk.go`
- `internal/prose/walk_test.go`
- `internal/canary/inventory.go`
- `internal/canary/inventory_test.go`
- `cmd/bench/main.go`
- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `tests/canary/package-core-guard/unrouted-subcommand`
- `reviews/ft290-test-projection.md`

## Ticket graph

| ticket | blocked by | chunk |
| --- | --- | --- |
| `1-split-named-check-owner.md` | none | TP-C1a |
| `2-count-tests-run.md` | none | TP-C1a |
| `3-prove-named-check-ran.md` | `1-split-named-check-owner.md`, `2-count-tests-run.md` | TP-C1a |
| `4-print-prose-check-result.md` | `3-prove-named-check-ran.md` | TP-C1b |
| `5-show-each-failure-diagnostic.md` | `2-count-tests-run.md` | TP-C1b |
| `6-filter-system-suite.md` | `3-prove-named-check-ran.md` | TP-C2 |
| `7-name-running-executable.md` | `1-split-named-check-owner.md` | TP-C2 |
| `8-list-check-fixtures.md` | `6-filter-system-suite.md` | TP-C3 |
| `9-list-check-inventory.md` | `8-list-check-fixtures.md` | TP-C3 |
| `10-explain-changed-selection.md` | `2-count-tests-run.md` | TP-C4 |

Tickets 6, 8, and 9 each add one form to the pinned grammar text, so they land
in that order. Ticket 9 completes the text that TP50 compares.

## Out of scope

- The run binary provenance in a result: `decisions/run-binary-provenance.md` owns it. Estimate: 6 edits, 2 gate runs.
- A `--run` filter inside a conformance check, which needs subtests in the root conformance test. Estimate: 12 edits, 3 gate runs.
- The pinned paths of a fixture as a projection. Estimate: 4 edits, 1 gate run.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| #1: the destination holds `--run` with `--check system` and the failures row unit | TP13 to TP17, TP19 to TP21 |
| #1: the compile error is a reviewed exclusion | story 35, Not covered |
| #2, amended 2026-09-19: the refusal names the running executable path and the source digest of its seal, or `unsealed` | TP24, TP25, TP26, TP54 |
| #3, fixed 2026-09-19: the row `check[1]{name,kind,tests_run,subjects}`, the zero rule, and the packages `tests_run` cell | TP1 to TP6, TP18, TP49, TP51, TP55 |
| #4, approved 2026-09-19: one row for each fixture by the inventory owner, no test run, and an empty table at exit 0 | TP27 to TP33, TP40, TP52, TP53 |
| #5: the prose `check` row, no packages table, the subjects under `--full`, and zero subjects exits 1 | TP7 to TP12 |
| #6: one `selected_by` cause by precedence, and `imports` names the first selected dependency in sorted order | TP41 to TP48 |
| #7: one row for each named check with kind and family count at exit 0, from one source with the help text | TP34 to TP38 |
| #8: the default `lines` count, and one row for each diagnostic line under `--full` | TP13 to TP17 |
| #9: `system` only, and a pattern with no match exits 1 | TP19 to TP23 |
| #10: the provenance moves to its own map | story 36, Not covered |
| Derived from #4, #7, and #9: each new form has a grammar, so the help text shows it and a wrong combination refuses | TP39, TP50 |

### Reader sweep and proof checklist

Readers of the rendered report and of the `testreport` interface:

- `cmd/bench/main.go` line 185 calls `testreport.Command`, and line 109 holds the help row.
- `internal/probe/probe.go` appends the report text and reads `Outcome.FailedTests` and `Outcome.Ran`. `internal/probe/baseline.go` reads `Outcome.Kind`. `internal/probe/command.go` reads `ProbeNotes`. The probe refuses the `system` and `prose` checks, so only the conformance `check` row reaches it.
- `internal/worktree/merge.go` line 101 prints the `--changed` output and reads only the exit code.
- `internal/gate/lane_select.go` line 292 runs `test --check <name>` for each lane check and reads the exit code. A real conformance run has `tests_run` 1, so the zero rule does not red the lane.
- `CHANGELOG.md` lines 191 and 192 name the old `packages` header. That entry is a historical record, and no row changes it.
- `internal/anchors/registry_data.go` lines 419 and 420 pin the `phases` and the `failures[N]{phase,line}` tables of the gate. They do not pin a `bench test` table.
- The sweep with `rg --hidden` found no other guidance file, script, or workflow file that names the two headers.

Pinned literals that the grammar change moves:

- `cmd/bench/main.go` line 109, the `Suffix` of the `test` help row.
- `cmd/bench/command_registry_test.go` line 764, the whole grammar text, and line 771, the text `bench test [--full] --check <name>`, which the new grammar keeps.
- `cmd/bench/help_inventory_test.go` line 84, the whole help row.
- `internal/testreport/command.go`, the `Cmd` and `Help` fields of the grammar.

Pinned literals that stay: the anchor needles and canary fixtures that name
`bench test --check system`, `bench test --check <owning-check>`,
`bench test --check skip-ownership`, and `bench test --package ./internal/conformance`.
No row changes those bytes.

Proof checklist:

- Cited symbols: each symbol below resolves in the tree at `bcc5543f`.
  - In `internal/testreport`: `Command`, `Prepare`, `Execute`, `Request.Run`, `Outcome.FailedTests`, `Outcome.Ran`, `selectRunBinary`, `namedChecks`, and `selectCurrentPackages`.
  - In `internal/canary`: `Fixtures`, `Select`, and `FixturePins`.
  - In other packages: `registry.CanaryFamilies`, `registry.Names`, `freshness.SealDigests`, `prose.Grade`, and `gate.SystemSuite`.
- Import edges: `internal/testreport` gains `internal/canary` and `internal/freshness`. `internal/canary` imports neither `internal/testreport` nor `internal/gate`. `internal/runbinary` already imports `internal/freshness`.
- Source-row clauses and occurrences: the source trace table above.
- Promised field labels: `check{name,kind,tests_run,subjects}`, `packages{package,status,elapsed_ms,tests_run}`, `packages{package,status,elapsed_ms,tests_run,selected_by}`, `failures{package,test,line,lines}`, `subjects{path}`, `fixtures{family,fixture,path}`, and `checks{name,kind,families}`.
- Changed-function callers: `prose.Grade` has two callers, `internal/testreport/command.go` and `internal/conformance/prose_mechanics_test.go`. The build keeps `prose.Grade` and adds the subject answer beside it, so the conformance caller does not change. `selectCurrentPackages` has one caller, `resolveChangedPackagesWithLoader`.
- Copy survival: none.

Sources re-read in the authoring session: `roadmap/FT290.md`,
`internal/testreport/command.go`, `testreport.go`, `outcome.go`,
`selection.go`, `selection_facts.go`, `check_test.go`,
`internal/conformance/registry/registry.go`, `internal/canary/inventory.go`,
`internal/freshness/freshness.go`, `internal/prose/walk.go`,
`cmd/bench/command_registry_test.go`, and `internal/conformance/checks_test.go`.
Not re-read: `cmd/bench/help_inventory_test.go` past line 84, and the
`internal/probe` tests. The repair pass also read `internal/toon/toon.go` and
`internal/testreport/outcome.go` lines 80 to 100.

### Fence disposition

Reviewer disposition of the ownership fences: open. The freshness package stays
outside the fence, because the build only calls `freshness.SealDigests`.

The fence holds two canary files as a flagged expansion for reviewer veto.
Ticket 8 exports the no-fixtures sentinel there, so that the empty answer has
one owner.

The build preflight binds four more paths to the help row file. They are the
command registry file, the two conformance registry tests, and the
unrouted-subcommand fixture. They are in the fence for that reason only. The spec authoring commit added
the running executable term to the glossary, so no ticket writes the glossary.
No ticket writes the changelog, because that write pulls the anchor registry files into the fence. The reviewer decides where the changelog entry lands.

### Completion plan

```bench-completion-plan
{"version":1,"chunks":[{"id":"TP-C1a","tickets":["1-split-named-check-owner.md","2-count-tests-run.md","3-prove-named-check-ran.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"}]},{"id":"TP-C1b","tickets":["4-print-prose-check-result.md","5-show-each-failure-diagnostic.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"prose","command":"bench test --package ./internal/prose"},{"id":"probe","command":"bench test --package ./internal/probe"}]},{"id":"TP-C2","tickets":["6-filter-system-suite.md","7-name-running-executable.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"cmd","command":"bench test --package ./cmd/bench"}]},{"id":"TP-C3","tickets":["8-list-check-fixtures.md","9-list-check-inventory.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"canary","command":"bench test --package ./internal/canary"}]},{"id":"TP-C4","tickets":["10-explain-changed-selection.md"],"verification":[{"id":"testreport","command":"bench test --package ./internal/testreport"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/ft290-test-projection/spec.md"},{"id":"testreport","command":"bench test --package ./internal/testreport"},{"id":"cmd","command":"bench test --package ./cmd/bench"},{"id":"probe","command":"bench test --package ./internal/probe"},{"id":"canary","command":"bench test --package ./internal/canary"}]}
```

### Flagged additions

Each addition below is not in a ticket answer. The reviewer can veto each one.

- The error title `named check ran nothing` for the zero rule.
- The `seal:` and `executable:` line labels and their order in the refusal.
- The `lines` cell on each `--full` failures row, so that one header serves each mode.
- The `subjects[N]{path}` table name for the prose `--full` list.
- The `--fixtures` root is the graded repository root, as for `bench canary`.
- The empty `family` cell for a fixture directly under `tests/canary`.
- The usage refusal for each extra flag beside `--checks` or `--fixtures`.
- The exit 1 refusal for an invalid canary inventory.
- The fence expansion to `internal/canary/inventory.go` and its test file, for the exported no-fixtures sentinel.
- The build-failure disposition of TP51: the compile diagnostic wins over the zero-rule title.
- The grader-refusal disposition of TP11: the grader diagnostic wins over the zero-rule title.
- The escape of the caller's check name in the unknown-check refusal, TP54.
- The exit 1 render refusal for an unprintable fixture name, TP52.
- The exit 0 answer for a `--changed` run with no run event, TP55.
- The grammar rows TP39 and TP50, which derive from decisions #4, #7, and #9.
