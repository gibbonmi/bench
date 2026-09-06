# The bench probe verb

Status: staged

Roadmap: FT303

Decision source: the `software-factory` decision map tickets #10 to #12, answered by the reviewer on 2026-09-06, with `decisions/assets/ft303-cli-assessment.md` as the source evidence

Verification log: 2 iteration(s) to accept — the round on `gpt-6-astra` at high effort returned revise at iteration 1 with six blocking findings. The author folded all six. Iteration 2 returned accept with six prose and accounting corrections, folded into the acceptance. The acceptance grades the spec and the tickets; the build records the runtime reds.

## Problem

A coordinator or a delegate proves a done-claim with a mutation probe. The
probe copies a file aside, changes one line, runs one focused test, compares
the bytes, and restores the file. The FT303 assessment counted 102 such
sequences in three days, at a median of three raw shell calls each.

The census counts a probe that names the pool path as raw shell, so the signal
cannot tell a probe from a leak. A mutation that fails to compile, or a run
pattern that matches no test, looks like a deliberate red. No verb owns the
sequence. `bench worktree path` answers a path that 198 raw calls then used in
a shell step.

## Solution

`bench probe` is a root verb. It takes one subject file, one exact mutation,
and one focused run in the `bench test` selection grammar. It preserves the
subject under the Bench home, applies the mutation once, and runs the focused
test. Then it restores the subject and proves the restore byte-exact against
the bytes it read at the start.

It prints one verdict row and then the focused run's own tables. The verdicts
are `bit` at exit 0, `silent` at exit 1, `invalid` at exit 1, and
`restore-failed` at exit 2 with the preserved copy named. Every validation
refusal runs before any write. `bench worktree path` gains a one-line note on stderr that
the path serves the file tools.

## User stories

### Group A — the verb

Line: opus / medium. The verb composes the focused-test owner and carries a
restore invariant. The scorecard routes gate and adapter logic to the mid tier
at medium effort.

1. As a coordinator, I want one call to apply a swap, run the focused test, and restore the file, so that one call replaces five.
2. As a coordinator, I want `--omit <old>` to delete the one match, so that an omission needs no empty replacement.
3. As a coordinator, I want the verdict `bit` at exit 0 when a focused test fails under the mutation, so that a bite is success.
4. As a coordinator, I want the verdict `silent` at exit 1 when every focused test passes, so that a vacuous probe is red.
5. As a coordinator, I want the verdict `invalid` at exit 1 when the mutated package fails to build, so that a compile failure never bites.
6. As a coordinator, I want the verdict `invalid` at exit 1 when the run runs no test, so that a mistyped pattern never bites.
7. As a coordinator, I want one verdict row that names the subject, mutation, cause, failed count, and restore state, so that I cite it.
8. As a coordinator, I want the focused run's own tables after the verdict row, so that the evidence is what `bench test` prints.
9. As a coordinator, I want `--check <name>` to run the registered check as the focused run, so that a conformance probe needs no hand mutation.
10. As a coordinator, I want `--full` to reach the focused run, so that a long failure diagnostic is not previewed.
11. As a delegate, I want the verb to run under `bench worktree exec <target> --`, so that no pool path appears in the call.
12. As an operator, I want the verb to run in the primary checkout and from a subdirectory, so that a probe needs no worktree.

### Group B — preservation and restore

Line: opus / medium. The restore is the invariant a wrong build breaks in
silence. The scorecard routes gate and adapter logic to the mid tier at medium
effort.

13. As a coordinator, I want the subject bytes preserved under the Bench home before the mutation, so that a crash leaves a copy.
14. As a coordinator, I want the subject restored to the preserved bytes after the focused run, so that the tree is clean.
15. As a coordinator, I want the restore proven by a byte comparison against the start bytes, so that `restored=yes` is evidence.
16. As a coordinator, I want the verdict `restore-failed` at exit 2 to name the preserved copy path, so that I can restore by hand.
17. As a coordinator, I want the preserved copy removed after a proven restore, so that the home holds no litter.
18. As a coordinator, I want the restore to run when the focused run is interrupted, so that an interrupt leaves no mutation.
19. As a coordinator, I want the subject's file mode kept across the mutation, so that an executable stays executable.
20. As a coordinator, I want the verb to write only the preserved copy and the subject, so that no second record owner appears.

### Group C — refusals before any write

Line: opus / medium. A refusal that runs after a write is the failure class
the reviewer excluded. The scorecard routes gate and adapter logic to the mid
tier at medium effort.

21. As a coordinator, I want a refusal unless the old string matches exactly once, so that a probe never mutates the wrong site.
22. As a coordinator, I want a refusal when `--with` equals `--swap`, so that a no-op mutation never runs.
23. As a coordinator, I want a refusal while a gate run holds the tree, so that the gate's subject never changes under it.
24. As a coordinator, I want an absent, directory, symlink, special-file, or outside-root subject refused, so that only a regular tree file mutates.
25. As a coordinator, I want a usage exit 2 for a missing or doubled selection or mutation, so that the grammar has one of each.
26. As a coordinator, I want an unknown `--check` name refused by the `bench test` rule before any write, so that the two verbs agree.
27. As a coordinator, I want `--check prose` and `--check system` refused as probe targets, so that a probe always runs a Go test or a registered check.
28. As a coordinator, I want each refusal as one structured line on stdout, and an unknown check as `bench test`'s usage, so that answers parse.
29. As a coordinator, I want a refusal to touch no subject, write no copy, and start no run child, so that nothing changes.
30. As a session outside a repository, I want the not-in-repo error, so that the verb matches every AXI command.

### Group D — the focused-run owner

Line: opus / medium. The owner change touches the one consumer of the gate's
test-argv producer. The scorecard routes gate and adapter logic to the mid tier
at medium effort.

31. As the probe, I want `bench test`'s parser and runner exposed as a typed outcome, so that the verdict derives from the report.
32. As a `bench test` caller, I want its output and exit codes unchanged, so that the merge and the operator keep their contract.

### Group E — the registry and the routing

Line: opus / low. Each edit is mechanical at a known seam under a covering
golden and three conformance checks. The scorecard routes such a ticket to the
mid tier at low effort.

33. As a session, I want `bench help` to list `bench probe` with its grammar, so that the verb is discoverable.
34. As a session, I want `bench probe --help` to print the usage line on stdout with exit 0, so that help is consistent.
35. As the gate, I want the routing, AXI registry, and parity checks green with the new verb, so that the inventory owners agree.
36. As the lane, I want no over-budget file to gain lines, so that the growth ratchet stays green.
37. As a maintainer, I want the anchors command, the help tests, and the routing table in their own files, so that the registry has headroom.

### Group F — guidance

Line: opus / medium. Guidance prose routes to the mid tier at medium effort
under the scorecard's 2026-08-26 rule. The profile's guidance override names
high effort, and the scorecard's later decision names medium. The spec follows
the scorecard and records the conflict for the reviewer.

38. As a delegate, I want a stderr note on `bench worktree path` that the path serves the file tools, so that shell steps use exec.
39. As a session, I want the reference guide to describe the probe verb, verdicts, and preserved copy, so that a cold session finds it.
40. As a cold session, I want `CONTEXT.md` to define **probe verdict** with its Avoid list, so that the vocabulary does not drift.

## Implementation decisions

- The verb is the root verb `probe` in a new package `internal/probe`. The dispatcher registers it beside `test` with the disposition `axiExempt(axiReasonMutation)`, because the verb writes the tree for the run's span. The routing table routes `probe` to `internal/probe`. The wrapper's default case routes the verb to the Go binary, so the wrapper gains no line.
- The grammar is `bench probe <file> (--swap <old> --with <new> | --omit <old>) (--package <expr> [--run <go-regex>] | --check <name>) [--full]`. Exactly one mutation form and exactly one selection form are required. `--run` needs `--package`. Every value flag refuses an empty value. A missing or doubled form prints the usage line on stdout and exits 2.
- The subject resolves the way `bench anchors` resolves its path: an absolute path stays, and a relative path joins onto the working directory. The verb then requires the path to sit inside the repository root and to be a regular file by `Lstat`. Each hostile subject refuses with `error: probe subject unavailable — <path> is <reason>`. The hostile subjects are a path outside the root, an absent path, a directory, a symlink, and a special file. The reasons are `outside the repository`, `absent`, `a directory`, `a symlink`, and `a special file`. The row and the refusal spell the subject repo-relative with forward slashes.
- The mutation is one exact string. `--swap <old> --with <new>` replaces the one match of `old` with `new`. `--omit <old>` replaces the one match with nothing. The old string must match exactly once. Zero or several matches refuse with `error: probe mutation ambiguous — the old string matches <n> times, want exactly 1`. A `--with` equal to `--swap` refuses with `error: probe mutation empty — --with equals --swap`.
- The verb refuses while a gate run holds the tree. It asks the gate package's new reader `ExecutionInProgress(root)`, which resolves the checkout administration directory and reads the execution lock as the verdict inspection does. A held lock refuses with `error: gate execution in progress — wait for the gate run to finish before you mutate the tree`. An unreadable lock refuses with `error: gate execution state unavailable — <reason>`.
- The verb validates the selection through the focused-run owner before any write. `testreport.Prepare(root, args)` parses the selection with `bench test`'s grammar and returns its usage line and code verbatim. The verb passes `--package`, `--run`, `--check`, and `--full` through unchanged. `--check prose` and `--check system` refuse, because neither runs a Go test through the report. The line is `error: probe focused run unsupported — --check prose and --check system are not probe targets`.
- The refusal order is: usage, not in a repository, the subject, the mutation count, the selection, the unsupported check, then the gate lock. Every validation refusal runs before the preserved copy is written and before any focused-run or run-binary build child starts. The root read and the administration-directory read run Git, which is not a run child.
- The preserved copy lives at `<home>/probe/<repo-key>/<unix-nanoseconds>-<pid>/<basename>` with mode 0600 in a 0700 directory. The home is the one `BENCH_HOME` read, and the repository key is the pool key. A failure to write the copy refuses with `error: probe preservation failed — <reason>` and mutates nothing.
- The mutation write and the restore write each replace the subject atomically. The verb writes a sibling temporary file with the subject's mode and renames it over the subject. A failed mutation write leaves the subject unchanged, removes the copy, and refuses with `error: probe mutation failed — <reason>`.
- The focused run is `testreport.Execute(root, request)`, which returns a typed outcome, the rendered report, and the exit code `bench test` prints. The outcome kinds are `passed`, `failed`, `build-failed`, `no-test-run`, `refused`, and `interrupted`, with the failed-test count. `Command` becomes `Prepare` then `Execute`, so `bench test` keeps its output and exit codes.
- The verdict derives from the outcome. `failed` gives `bit` at exit 0. `passed` gives `silent` at exit 1. `build-failed`, `no-test-run`, `refused`, and `interrupted` give `invalid` at exit 1. A package that fails with no failing test is `build-failed`. A run that runs no test is `no-test-run`.
- The restore runs after every focused run, and also after an interrupted run, through one deferred path. The preserved copy file is the restore source, and the bytes read at the start are the oracle. The verb writes the copy's bytes back, reads the subject again, and compares the read-back with the start bytes. Equal bytes remove the copy directory and set `restored=yes`. Unequal bytes or a failed write keep the copy and give `restore-failed` at exit 2 with the reason. The `restore-failed` verdict overrides the run's verdict.
- The output is content-first. The first stdout block is `probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:`. The `mutation` cell is `swap` or `omit`. The `cause` cell is the outcome kind. The focused run's rendered report follows the row unchanged. A `restore-failed` verdict adds `preserved[1]{path,reason}:` after the row, with the absolute copy path.
- The verb appends no `help[]` envelope, because it is not an approved AXI query.
- The restore runs before the render. A subject path that `toon.Table` cannot represent restores the subject, then prints the shared render error line and exits 1.
- The help row sits at inventory order 22 after `bench test`. Its suffix is the grammar's operand list. Its description is `mutate one file once, run one focused test or check, restore the file, and report bit, silent, invalid, or restore-failed`.
- The three registry files that gain a line sit over the structure budget, so a prefactoring ticket moves one cohesive block out of each. The moved inventory is the four anchors symbols, the two help-inventory tests, and the routing table with its why constants. The anchors grammar, command, and its two helpers move from `cmd/bench/main.go` to `cmd/bench/anchors_command.go`. The two help-inventory tests move from `cmd/bench/main_test.go` to `cmd/bench/help_inventory_test.go`. The `subcommandRouting` table and its why constants move from `internal/conformance/subcommand_routing_test.go` to `internal/conformance/subcommand_routing_table_test.go`.
- The ticket binding registry then names the two new files as the help assertion and the routing census. It also gains the row `internal/probe` with the command registry set. That row is discretionary ownership closure, because the ticket-grammar check binds only AXI query packages and seed owners.
- `bench worktree path` keeps its stdout contract, one path line. After the path, it prints one note line on stderr with the target as typed. The line is `note: the path serves the file tools; run a shell step through bench worktree exec <target> -- <command>`. The help row description becomes `print one active owned worktree's absolute path for the file tools`.
- The reference guide's Command Notes gain one paragraph on `bench probe`. It names the sequence, the four verdicts with their exits, the preserved copy location, and the exec form. `CONTEXT.md` gains **probe verdict**: the one word `bench probe` prints for a run, `bit`, `silent`, `invalid`, or `restore-failed`. Its Avoid list holds `test result`, `probe outcome`, and `mutation score`.

## Testing decisions

- A good verb test drives `probe.Command` over a temporary Go module with the real `go` on `PATH`. It reads the exact stdout, the exit code, the subject bytes, and the preserved directory. The prior art is `TestRunPatternRefusesZeroMatches` in `internal/testreport`, which drives a real focused run over `focusedTestModule`. The probe package builds its own module fixture `probefixture` with `clamp.go` and `clamp_test.go`, one module constructor for the package.
- The fixture is `Clamp(n int) int` with a comment `// Clamp keeps n at or above zero.`, a guard `if n < 0 {` that returns `0`, and a final `return n`. `TestClampNegative` asserts `Clamp(-1) == 0`, and `TestClampPositive` asserts `Clamp(3) == 3`.
- The outcome kinds run through a stub `go` on `PATH` that emits canned `-json` events. The prior art is `writeCheckGo` in `internal/testreport/check_test.go`. The stub records a marker file when it starts, so a refusal row proves that no run child started.
- The restore-failure row makes the subject's directory read-only from inside the stub `go`, so the restore's temporary file cannot be created.
- The interrupt row runs the verb in a child process and signals it, in the shape of `TestFocusedGoTestDrainsGroupOnCancelSignal` in `internal/testreport/cancel_test.go`.
- The gate-lock rows hold the execution lock from a child process, in the shape of `runGateLockHolder` in `internal/gate/run_failure_outcomes_test.go`.
- The `--check` form is proven at the parse seam and by one recorded run over this repository. A registered check builds a Bench executable and grades the kit, so a fixture module cannot host it. The `internal/probe` tests reach no test-only helper of `internal/testreport`; the `refused` outcome case runs inside `internal/testreport` with that package's selection factory.
- The focused-run tables and the `--full` rows use canned events, because a real run's `elapsed_ms` differs between two runs. The `--full` row uses a diagnostic longer than `bounds.PreviewRuneLimit`, because `--full` lifts the preview limit on the first diagnostic line.
- The `bench test` compatibility row records the base commit's exact output for each canned set before the refactor, and the candidate matches those goldens. The goldens cover the package form and the run form. The merge caller's `--changed` form keeps its existing changed-selection tests, so the comparison is not a differential over every caller form.
- The gate's test phase runs the new package tests and the root conformance entry, which observes the routing, the AXI registry, and the parity checks. The fast lane's growth check observes story 36 on each worktree commit.

### Seam diagram

    bench probe <file> --swap|--omit ... --package|--check ...
        │
        ▼
    args, cwd  ──▶  [ internal/probe: resolve subject, count the match, Prepare the selection, read the gate lock ]  ──▶  refusal line | usage line
                          ◀ tests attach here: stdout, exit, subject bytes, an empty preserved directory, an absent stub marker

    accepted request
        │
        ▼
    subject bytes  ──▶  [ preserve under <home>/probe, replace atomically, testreport.Execute, restore, compare ]  ──▶  probe[1] row + report tables | preserved[1] row
                          ◀ tests attach here: a temporary module with the real go, a stub go for the outcome kinds

    bench test [...]
        │
        ▼
    args  ──▶  [ testreport.Prepare, testreport.Execute ]  ──▶  Outcome{kind, failed_tests}, rendered report, exit
                          ◀ tests attach here: canned -json event sets through the stub go, and Command equality

    bench worktree path <target>
        │
        ▼
    target  ──▶  [ PathCommand ]  ──▶  stdout: the path; stderr: the file-tools note
                          ◀ tests attach here: stdout equality and the stderr line

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| PB1 | 1, 3, 7 | over the fixture module, `probe clamp.go --swap "n < 0" --with "n > 0" --package ./ --run ^TestClampNegative$` prints the first line `probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:` and the row `bit,clamp.go,swap,failed,1,yes`, and exits 0 | a new test `TestProbeBitesWhenTheFocusedTestFails` in `internal/probe` with the real `go` | a verb that skips the mutation reports `silent`, and a verb that skips the run has no failed count |
| PB2 | 2, 3 | over the fixture module, `--omit "return 0"` with the same selection prints the row `bit,clamp.go,omit,failed,1,yes` and exits 0 | a new test `TestProbeOmitsTheMatchOnce` in `internal/probe` | a verb that ignores `--omit` runs an unchanged file and reports `silent` |
| PB3 | 4 | over the fixture module, `--swap "keeps n at" --with "holds n at"` with the same selection prints the row `silent,clamp.go,swap,passed,0,yes` and exits 1 | a new test `TestProbeIsSilentWhenTheFocusedTestPasses` in `internal/probe` | a verb that exits 0 on every completed run reds the exit assertion |
| PB4 | 5 | over the fixture module, `--swap "return n" --with "return"` with `--package ./` and no `--run` prints the row `invalid,clamp.go,swap,build-failed,0,yes` and exits 1 | a new test `TestProbeIsInvalidWhenTheMutationDoesNotCompile` in `internal/probe` | a verb that reads a nonzero `go` exit as a bite reports `bit` |
| PB5 | 6 | over the fixture module, the PB1 mutation with `--run ^TestNoSuch$` prints the row `invalid,clamp.go,swap,no-test-run,0,yes` and exits 1 | a new test `TestProbeIsInvalidWhenNoTestRuns` in `internal/probe` | a verb that reads an empty failure table as `silent` reds the cause cell |
| PB6 | 8 | over the stub `go` with a canned failing event set, stdout after the probe row equals the output `testreport.Command` prints for the same selection over the same canned set | a new test `TestProbeCarriesTheFocusedRunTables` in `internal/probe` over the stub `go` | a verb that prints the tables before the row, or drops them, reds the equality |
| PB7 | 9 | for `--check line-routing`, the parse step's selection request equals the request `testreport.Prepare` returns for `--check line-routing`, and the existing named-check test pins that request's argv | a new in-package test `TestProbeSelectsTheCheckForm` in `internal/probe` over the parse step, and `internal/testreport/check_test.go` (`TestNamedCheckOwnsConformanceEnvironment`) | a verb that maps `--check` onto `--package` builds a different request |
| PB8 | 9 | `bench probe internal/conformance/subcommand_routing_table_test.go --omit '"probe": routed("internal/probe"),' --check subcommand-routing` over this repository prints a `bit` row | review-owned: a recorded run in the build's integration worktree, cited under Further notes | a check form that never reaches the runner has no run to record |
| PB9 | 10 | with a canned failing diagnostic longer than `bounds.PreviewRuneLimit`, the `failures[` row's `line` cell carries the whole diagnostic under `--full` and its preview without `--full` | a new test `TestProbeForwardsFullToTheFocusedRun` in `internal/probe` over the stub `go` | a verb that drops `--full` prints the preview twice |
| PB10 | 12 | with the working directory set to a subdirectory of the fixture module, `probe ../clamp.go` resolves the subject and prints the row with the subject `clamp.go` | a new test `TestProbeResolvesTheSubjectFromTheWorkingDirectory` in `internal/probe` | a verb that joins the operand onto the root refuses the path as absent |
| PB11 | 11 | `bench worktree exec <label> -- bench probe internal/probe/probe.go --swap <old> --with <new> --package ./internal/probe --run <regex>` over the integration worktree prints a `bit` row, and the command names no pool path | review-owned: a recorded run in the build's integration worktree, cited under Further notes | the exec route is unchanged, so the recorded run is the proof that the verb needs no path operand |
| PB12 | 13 | while the stub `go` runs, a file with the subject's original bytes exists at `<home>/probe/<repo-key>/<stamp>/clamp.go` with mode 0600 | a new test `TestProbePreservesTheSubjectBeforeTheRun` in `internal/probe`, whose stub `go` lists the preserved directory into its marker | a verb that preserves in memory only leaves no file for the stub to list |
| PB13 | 14, 15 | after the PB1 run, the subject bytes equal the bytes the test wrote, and the row's `restored` cell is `yes` | the PB1 test | a verb that writes the mutated bytes back reds the byte comparison |
| PB14 | 16 | with the stub `go` making the subject's directory read-only, the verb prints the row `restore-failed,clamp.go,swap,failed,1,no`, then `preserved[1]{path,reason}:` with the copy's absolute path, exits 2, and the copy keeps the original bytes | a new test `TestProbeReportsARestoreFailure` in `internal/probe` | a verb that removes the copy on every exit, or exits 1, reds the assertions |
| PB15 | 17 | after the PB1 run, `<home>/probe/<repo-key>/` holds no entry | the PB1 test | a verb that keeps every copy leaves the stamp directory |
| PB16 | 18 | with a stub `go` that sleeps until a signal, an interrupt to the verb's process restores the subject before the verb exits, and the verb prints the row `invalid,clamp.go,swap,interrupted,0,yes` | a new test `TestProbeRestoresOnInterrupt` in `internal/probe`, run as a child-process helper | a verb whose restore is not deferred leaves the mutation on disk |
| PB17 | 19 | a subject with mode 0755 keeps mode 0755 after the PB1 run | a new test `TestProbeKeepsTheSubjectMode` in `internal/probe` | a restore through a 0644 write reds the mode assertion |
| PB18 | 21 | `--swap "return" --with "x"` prints `error: probe mutation ambiguous — the old string matches 2 times, want exactly 1` and exits 1, and `--swap "absent text" --with "x"` prints the same line with `0 times` | a new test `TestProbeRefusesAnAmbiguousMutation` in `internal/probe` | a verb that replaces the first match mutates the wrong site and runs |
| PB19 | 22 | `--swap "n < 0" --with "n < 0"` prints `error: probe mutation empty — --with equals --swap` and exits 1 | a new test `TestProbeRefusesAnEmptyMutation` in `internal/probe` | a verb that runs the test over an unchanged file reports `silent` |
| PB20 | 23 | with a child process holding the gate execution lock on the fixture repository, the PB1 call prints `error: gate execution in progress — wait for the gate run to finish before you mutate the tree` and exits 1 | a new test `TestProbeRefusesUnderALiveGateRun` in `internal/probe` with a lock-holder child | a verb that never reads the lock mutates under the gate |
| PB21 | 24 | a subject that is absent, a directory, a symlink to a regular file, a FIFO, or a path outside the root each prints `error: probe subject unavailable — <path> is <reason>` with its reason and exits 1 | a new test `TestProbeRefusesAHostileSubject` in `internal/probe`, five cases | a verb that follows the symlink mutates the target, and a verb that skips the root check mutates outside the tree |
| PB22 | 25 | no selection, `--package` with `--check`, `--run` without `--package`, `--swap` with `--omit`, and no mutation each print the usage line on stdout and exit 2 | a new test `TestProbeUsageNamesOneSelectionAndOneMutation` in `internal/probe`, five cases | a grammar that defaults the package runs the whole tree |
| PB23 | 26 | `--check no-such-check` prints the `bench test` unknown-check refusal with its check inventory on stdout and exits 2 | a new test `TestProbeRefusesAnUnknownCheckBeforeAnyWrite` in `internal/probe` | a verb that parses the selection after the mutation leaves a copy behind |
| PB24 | 27 | `--check prose` and `--check system` each print `error: probe focused run unsupported — --check prose and --check system are not probe targets` and exit 1 | a new test `TestProbeRefusesProseAndSystemChecks` in `internal/probe` | a verb that runs the prose check over the mutated tree reports `invalid` |
| PB25 | 28, 29 | for each refusal in PB18 to PB24, the subject bytes are unchanged, `<home>/probe` holds no entry, and no focused-run or run-binary build child started, which the absent stub `go` marker proves | the same tests assert the three facts | a verb that preserves before it validates leaves a copy, and a verb that starts the run first writes the marker |
| PB26 | 30 | outside a repository, the verb prints `error: not in a git repository — run inside a Bench-linked repo` and exits 1 | a new test `TestProbeOutsideARepository` in `internal/probe` | a verb that resolves no root panics on the subject join |
| PB27 | 31 | over a stub `go` that emits a failing-test event set, a passing set, a build-fail set, and a no-run set, `testreport.Execute` answers the outcome kinds `failed` with count 1, `passed`, `build-failed`, and `no-test-run`, an interrupted run answers `interrupted`, and a run whose executable selection fails answers `refused` | a new test `TestExecuteClassifiesTheOutcome` in `internal/testreport`, with the package's selection factory for the `refused` case | an outcome derived from the exit code alone cannot tell `build-failed` from `failed`, and a classifier that skips the pre-run refusals answers `passed` for a refused start |
| PB28 | 32 | for the four canned event sets in the package form and the run form, `testreport.Command` prints the exact output and exit that the base commit's `Command` printed, recorded as goldens before the refactor, and equals `Prepare` then `Execute` | a new test `TestCommandKeepsItsBaseOutput` in `internal/testreport` with goldens recorded at the base commit, `TestCommandIsThePrepareExecuteProjection`, and the package's existing tests | a shared rendering regression reds the base goldens, and a second render path reds the projection equality |
| PB29 | 33 | `bench help` prints, on the line after the `bench test` row, a row that starts with `  bench probe <file> (--swap <old> --with <new>` and ends with `or restore-failed` | `TestHelpInventoryIsComplete` in `cmd/bench`, in its own file after ticket 01 | a registration with `internalInventory` reds the golden |
| PB30 | 34 | `bench probe --help`, `bench probe -h`, and `bench probe help` each print the grammar's usage line, which starts with `usage: bench probe <file> (--swap`, on stdout and exit 0 | a new test `TestProbeHelpSpellings` in `internal/probe` | a hand-rolled parser that treats `help` as the subject refuses it as absent |
| PB31 | 35 | the `subcommand-routing`, `axi-query-registry`, and `entry-point-parity` checks pass with `probe` routed to `internal/probe` and exempt as a mutation | the gate's test phase, through the root conformance entry | an entry point that hand-rolls its parse reds the routing check, and an approved disposition reds the AXI registry check |
| PB32 | 36 | `bench structure --growth <base>` over the landing source reports no over-budget file that gained lines | the fast lane's structure growth check on each worktree commit | a registry line added without the relocation reds the lane |
| PB33 | 37 | after the relocation, `cmd/bench/main.go`, `cmd/bench/main_test.go`, and `internal/conformance/subcommand_routing_test.go` each count fewer lines than at the base, and the four anchors symbols, the two help tests, and the routing table with its constants run from their new files | review-owned: the reviewer reads the three moves as pure moves, and the gate's test phase runs the moved tests | a move that drops a test leaves the golden unproven |
| PB34 | 38 | `bench worktree path <target>` prints the path alone on stdout and the line `note: the path serves the file tools; run a shell step through bench worktree exec <target> -- <command>` on stderr, with the target as typed, and exits 0 | a new test `TestPathNotesTheFileToolRouteOnStderr` in `internal/worktree` | a note on stdout breaks a `$(...)` capture, and the stdout equality reds it |
| PB35 | 38 | the `bench help` row for `bench worktree path` reads `print one active owned worktree's absolute path for the file tools` | `TestHelpInventoryIsComplete` in `cmd/bench`, in its own file after ticket 01 | an unchanged description reds the golden |
| PB36 | 39 | the reference guide's Command Notes hold a paragraph that names `bench probe`, the four verdicts with their exits, the preserved copy under `$BENCH_HOME/probe/<repo-key>/`, and the `bench worktree exec` form | review-owned: the reviewer reads the paragraph | the prose check grades sentences, not content |
| PB37 | 40 | `CONTEXT.md` defines **probe verdict** with the Avoid list `test result`, `probe outcome`, and `mutation score` | review-owned: the reviewer reads the entry | the prose check grades sentences, not terms |
| PB38 | 3, 4, 5, 16 | the verdict mapping answers `bit` at 0 for `failed`, `silent` at 1 for `passed`, `invalid` at 1 for `build-failed`, `no-test-run`, `refused`, and `interrupted`, and `restore-failed` at 2 over every kind | a new table-driven test `TestVerdictExitCodes` in `internal/probe` | a mapping that exits 0 on `silent` reds one case |
| PB39 | 7 | a subject at `my pkg/clamp.go` renders its cell through `toon.Table`, and the expectation derives through the same call | a new test `TestProbeRendersASubjectWithASpace` in `internal/probe` | a hand-joined row disagrees with the encoder's quoting |
| PB40 | 14 | a subject whose last line has no newline restores byte-exact, and the row's `restored` cell is `yes` | a new test `TestProbeRestoresAFileWithoutTrailingNewline` in `internal/probe` | a restore that appends a newline reds the byte comparison |
| PB41 | 7, 14 | a subject whose name carries a BEL byte completes the run, restores the subject, prints `error: unrepresentable TOON cell — <reason>`, and exits 1 | a new test `TestProbeRestoresBeforeARenderRefusal` in `internal/probe` | a verb that renders before it restores leaves the mutation on disk |
| PB43 | 15, 16 | with the stub `go` truncating the preserved copy during the run, the restore writes the truncated bytes, the read-back differs from the start bytes, and the verb prints the row `restore-failed,clamp.go,swap,failed,1,no`, the `preserved[1]` row with the reason, and exits 2 | a new test `TestProbeReportsAReadBackMismatch` in `internal/probe` | a verb that omits the byte comparison prints `restored=yes` after a wrong write |
| PB44 | 29 | with `<home>/probe` present as a regular file, the PB1 call prints `error: probe preservation failed — <reason>` and exits 1, the subject is unchanged, and no run child started | a new test `TestProbeRefusesWhenPreservationFails` in `internal/probe` | a verb that mutates before it preserves leaves the mutation |
| PB45 | 29 | with the subject's directory read-only before the call, the PB1 call prints `error: probe mutation failed — <reason>` and exits 1, the subject is unchanged, `<home>/probe` holds no entry, and no run child started | a new test `TestProbeRefusesWhenTheMutationWriteFails` in `internal/probe` | a verb that starts the run over an unmutated file reports `silent`, and a verb that keeps the copy leaves an entry |
| PB46 | 20 | after the PB1 run, the repository's administration directory, the Bench home, and the working tree hold no new entry | a new test `TestProbeWritesNoRecord` in `internal/probe`, which lists the three before and after | a verb that writes a ledger or a census entry adds an entry |
| PB42 | 23 | `gate.ExecutionInProgress(root)` answers true while a child process holds the execution lock, false after the child releases it, and an error over a directory that is not a repository | a new test `TestExecutionInProgressReadsTheLock` in `internal/gate` with `startGateLockHolder` | a reader that stats the lock file answers true after a finished run |

### Edge inventory

- Error paths: each refusal is one structured line on stdout with exit 1 (PB18 to PB21, PB24, PB26, PB44, PB45). The unknown-check usage exits 2 with the inventory (PB23). A preservation failure refuses before the mutation (PB44). A failed mutation write leaves the subject unchanged and refuses (PB45).
- Empty input: an empty subject file matches the old string zero times and refuses (PB18). An empty flag value is a usage error (PB22).
- Boundary values: the old string matches exactly once (PB1) or refuses (PB18). One failing test is `bit` (PB1), and zero failing tests with a run is `silent` (PB3).
- Hostile paths: a symlink, a FIFO, a directory, and a path outside the root refuse (PB21). A path with a space renders through the encoder (PB39). A path with a control byte restores before the render error (PB41). The subject spelling in the row is repo-relative (PB1, PB10).
- Hostile values: a mutation string is one argv token, so the shell owns its quoting, and the verb compares bytes. A subject without a trailing newline restores byte-exact (PB40).
- Currency: the verb reads the subject once, and the preserved copy is that read. A subject that changes during the run is restored to the read bytes, which is the restore promise (PB13).
- Re-run idempotency: a completed run leaves the tree and the home as they were (PB13, PB15). A second run over the same subject starts from the restored bytes.
- Partial implementation: a build that registers the verb without the routing row reds PB31. A build that restores without the comparison passes PB13 and PB14 and reds PB43. A build that preserves after the lock check but before the selection parse reds PB25. A build that writes a record reds PB46.
- Interrupt: an interrupt during the run restores (PB16). An interrupt during the restore write completes the write, because the context cancels only the child.
- Concurrency: a gate run refuses the probe (PB20, PB42). Two probes on one subject serialize by the coordinator, and the second sees a mutated file (Won't handle).
- Audience: the verb serves this repository and every repository that links the kit. A linked repository with no conformance checks uses the package form. The preserved copy under an unset `BENCH_HOME` lands under the user home, which is the one home reader's rule.
- Package-variable swaps: the `internal/probe` tests swap nothing, and the stub `go` sits on `PATH`. The `refused` case of PB27 swaps the selection factory inside `internal/testreport`, the venue that already owns that swap.
- Absent versus empty: an absent subject refuses (PB21), and an empty subject refuses on the match count (PB18). An absent home directory is created on the first preserve.
- In-flight tree states: the routing row, the registry line, and the golden row land in one ticket, so no commit carries two of the three.

**Won't handle** — two probes on one subject at the same time — the coordinator serializes probes per subject, and the match count refuses a mutated file.

**Won't handle** — a gate run started after the lock check — the gate refuses a subject that changed under it, and `bench gate` stays the caller.

**Won't handle** — a default package derived from the subject's directory — `--package` stays the caller's spelling, and every observed probe named its package.

**Won't handle** — a multi-site or regex mutation — the reviewer rejected `--sed` and `--patch` in ticket #11, and a coordinator runs one probe per site.

**Won't handle** — a verdict record in the census or a ledger — ticket #12 decided that the terminal row is the evidence the caller cites.

**Won't handle** — the exec child's `PWD` — FT254 owns exec comfort, and the verb resolves the root from the working directory the exec route sets.

**Won't handle** — a `bench worktree probe <target>` face — the reviewer chose the root verb in ticket #11, and `bench worktree exec` stays the worktree form.

**Won't handle** — a prose or system probe — `bench test --check prose` and `bench test --check system` stay the callers for those runs after a hand mutation.

## Ownership fences

- `specs/bench-probe/`
- `reviews/bench-probe.md`
- `internal/probe/` (new)
- `internal/testreport/command.go`
- `internal/testreport/testreport.go`
- `internal/testreport/outcome.go` (new)
- `internal/testreport/outcome_test.go` (new)
- `internal/gate/run_transaction.go`
- `internal/gate/execution_probe_test.go` (new)
- `cmd/bench/main.go`
- `cmd/bench/anchors_command.go` (new)
- `cmd/bench/main_test.go`
- `cmd/bench/help_inventory_test.go` (new)
- `internal/conformance/subcommand_routing_test.go`
- `internal/conformance/subcommand_routing_table_test.go` (new)
- `internal/tickets/registry_data.go`
- `internal/worktree/path.go`
- `internal/worktree/path_identifier_test.go`
- `.bench/BENCH-reference.md`
- `CONTEXT.md`
- `cmd/bench/command_registry.go` — closure headroom only
- `cmd/bench/command_registry_test.go` — closure headroom only
- `internal/conformance/axi_query_registry_test.go` — closure headroom only
- `tests/canary/docs-currency-token-diet/signal-vocabulary-drift` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-acceptance-row-vocabulary` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-map-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-row-parts` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-coverage-row-vocabulary` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-decision-map-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-reader-sweep-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/context-ticket-vocabulary` — closure headroom only
- `tests/canary/package-core-guard/unrouted-subcommand` — closure headroom only
- `tests/canary/docs-currency-token-diet/benchref-imported` — closure headroom only
- `tests/canary/docs-currency-token-diet/benchref-pointer-dropped` — closure headroom only
- `tests/canary/docs-currency-token-diet/benchref-section-duplicated` — closure headroom only
- `tests/canary/skills-index-command-adapters/adapter-inert-invocation-key` — closure headroom only
- `tests/canary/skills-index-command-adapters/command-invocation-disabled-against-policy` — closure headroom only
- `tests/canary/skills-index-command-adapters/dangling-index` — closure headroom only
- `tests/canary/skills-index-command-adapters/debug-implicit-invocation-reverted` — closure headroom only
- `tests/canary/skills-index-command-adapters/missing-index-field` — closure headroom only
- `tests/canary/skills-index-command-adapters/stale-index-wording` — closure headroom only
- `tests/canary/skills-index-command-adapters/unindexed-skill` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-agent-push-rule` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-bench-operational-layer` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-category-context` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-category-oracle` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-category-setup` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-category-work` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-gate-authority` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-kit-only-ship` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-no-path-fallback` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-progressive-loading-term` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-refusal-route-shape` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-retro-capture-owner` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-retro-drain-owner` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-skills-guidance` — closure headroom only
- `tests/canary/workflow-guidance-anchors/reference-upgrade-route` — closure headroom only

A closure headroom entry creates no blocker edge and takes no edit. The three
command registries close the `cmd/bench` and `internal/worktree` bindings. The
fixture directories close the pins on `CONTEXT.md`, on `.bench/BENCH-reference.md`,
and on the dispatcher.

The fence is the union of the tickets' `Writes:` lines, plus the spec folder
and the review pickup. `bench preflight build` closes it over the fixture and
registry pins. The wrapper `bin/bench.sh` is outside the fence, because its
default case already routes the verb.

## Out of scope

- A default package derived from the subject's directory. 3 edits, 1 gate run.
- A probe ledger or a census entry for verdicts. Ticket #12 decided against it. 8 edits, 2 gate runs.
- A `bench worktree probe <target>` face. Ticket #11 chose the root verb. 6 edits, 2 gate runs.
- `--sed` and `--patch` mutation forms. Ticket #11 chose the exact string. 5 edits, 1 gate run.
- A production-or-test projection on `bench consumers` and `bench outline`. Candidate 2 of the FT303 assessment. 6 edits, 2 gate runs.
- A `bench structure --path <prefix>` filter. Candidate 3 of the FT303 assessment. 3 edits, 1 gate run.
- The exec child's `PWD`. FT254 owns it. 2 edits, 1 gate run.
- A `bench test` face that names `no-test-run` as a compile attribution. FT290 owns it. 3 edits, 1 gate run.

## Further notes

Flagged additions beyond the decision source:

- The `invalid` verdict also covers a run that runs no test, a refused start, and an interrupted run. The source named the build failure. A probe that ran zero tests proves nothing, which is the `craft-delegate` rule.
- The `--with` equals `--swap` refusal. The source named the exact-match rule only.
- The `--check prose` and `--check system` refusal. The source named the two `bench test` selection forms.
- The preserved copy location, the atomic replacement, and the `cause` column. The source granted the copy location as discretion.
- The `--full` pass-through.
- The subject restrictions, the refusal order, the copy permissions, the mode retention, and the post-run render refusal. The profile's hostile-input checklist demands the restrictions and the render refusal, and the copy location grant covers the permissions. The refusal order and the mode retention are bounded discretion.
- The no-record row PB46, which gives ticket #12's decision a red-capable owner.
- The three relocations for headroom, which the growth ratchet forces, and the optional binding registry row as ownership closure.
- The stderr channel of the path note and the help description clause. The source said one guidance ticket adds the note.
- The **probe verdict** glossary term and the reference guide paragraph.

Build decisions recorded for reviewer veto:

- The registry line, the routing row, and the golden row land in one ticket after the relocation ticket.
- The path-note ticket runs after the verb ticket, because both write the dispatcher and the golden.
- Row PB4 drops `--run` from its selection. Under a run pattern, `bench test` reports `no-test-run` before it reports the build failure, so the compile row under `--run` answers `no-test-run`. The verdict and the exit stay `invalid` at 1. The Out of scope list already parks the compile attribution with FT290.
- The verdict row renders through `toon.TableTyped`, so the `failed_tests` cell is an integer. `toon.Table` quotes a numeric-looking string, and row PB1 spells a bare `1`. Every test expectation derives through the same call.
- The subject refusal has one reason beyond the spec's list: a regular file that cannot be read answers `is unreadable`.
- Ticket 02 took one fence amendment: one line in `internal/testreport/check_test.go`, because that test calls the changed `runGoTest` directly.
- `cmd/bench/command_registry.go` took one pure move: `outputCommand` and `adoptCommand` sit beside the `commandHandler` type they construct. The lane grades growth against the current tip, so the two registry lines in `cmd/bench/main.go` needed headroom at the same commit.
- `internal/probe` carries its own `prose` check name, because `testreport` does not export its constant. The one-source rule prefers an export; the ticket fence did not cover `internal/testreport/command.go`.

Source-sentence-to-row table:

| source sentence | rows |
|---|---|
| the root verb `bench probe`, and the exec route is the only worktree form | PB11, PB12, PB29, PB30, PB31 |
| an exact string pair or an omission on one named file, and the old string matches exactly once | PB1, PB2, PB18, PB19 |
| `bit` at exit 0, `silent` and `invalid` at exit 1, `restore-failed` at exit 2 with the preserved copy named | PB1, PB3, PB4, PB5, PB14, PB38 |
| the verb runs in any checkout and refuses while a gate run holds the tree | PB10, PB20, PB26, PB42 |
| the restore is proven byte-exact before exit | PB12, PB13, PB14, PB15, PB16, PB17, PB40, PB41, PB43 |
| the focused run composes the `bench test` selection grammar with both forms | PB6, PB7, PB8, PB9, PB22, PB23, PB27, PB28 |
| no record beyond the terminal verdict row | PB46 and the fifth Won't handle |
| one guidance ticket adds the one-line note to `bench worktree path` | PB34, PB35 |
| the `--time` face and the heredoc route are out of scope | the Out of scope list |

Pre-review proof checklist:

- Cited symbols in `internal/testreport`: `Command`, `parseFocusedRequest`, `runFocusedRequest`, `runGoTest`, `report`, `focusedTestModule`, `writeCheckGo`, `installTestSelectionFactory`, `TestRunPatternRefusesZeroMatches`, and `TestFocusedGoTestDrainsGroupOnCancelSignal`.
- Cited symbols in `internal/gate`: `lockHeld`, `acquireExecutionLock`, `startGateLockHolder`, and `runGateLockHolder`.
- Cited adapter symbols: `benchgit.AdminDir`, `benchhome.Dir`, `poolkey.Key`, `git.Root`, `toon.Table`, `toon.Errorf`, `toon.NotInRepo`, `toon.RenderError`, `usage.Parse`, and `bounds.PreviewRuneLimit`.
- Cited dispatcher and registry symbols: `anchorQueryPath`, `anchorsCommand`, `anchorKindName`, `anchorsGrammar`, `TestHelpInventoryIsComplete`, `TestHelpRendersPublicCommandRegistryRows`, `subcommandRouting`, and `commandRegistries`.
- Cited worktree symbols: `PathCommand` and `TestListPathActionRunsAsAdvertised`.
- New names: `testreport.Prepare`, `testreport.Execute`, `testreport.Outcome`, `gate.ExecutionInProgress`, and the `internal/probe` package.
- Import edges: `internal/probe` imports `internal/testreport`, `internal/gate`, `internal/git`, `internal/benchhome`, `internal/poolkey`, `internal/toon`, and `internal/usage`. `internal/testreport` already imports `internal/gate`, so no cycle forms. `cmd/bench` gains the import of `internal/probe`.
- Source-row clauses and occurrences: the source is the `## #10`, `## #11`, and `## #12` answers of `decisions/software-factory.md`. The table above lists each clause once.
- Promised field labels: the row header `probe[1]{verdict,subject,mutation,cause,failed_tests,restored}` and the `preserved[1]{path,reason}` row.
- Promised verdict words: `bit`, `silent`, `invalid`, and `restore-failed`. Promised cause words: `passed`, `failed`, `build-failed`, `no-test-run`, `refused`, and `interrupted`.
- The subject and mutation error kinds are `probe subject unavailable`, `probe mutation ambiguous`, `probe mutation empty`, `probe preservation failed`, and `probe mutation failed`. The run error kinds are `gate execution in progress`, `gate execution state unavailable`, and `probe focused run unsupported`.
- Changed-function callers: `testreport.Command` keeps its two callers, the dispatcher and the worktree merge. `PathCommand` keeps its dispatcher caller and its two tests. `lockHeld` keeps its inspection caller.
- Copy survival: none. The verb replaces no copy.

Reader sweep of the changed facts:

- `testreport.Command`'s signature and output: the dispatcher and `internal/worktree/merge.go` keep their calls (PB28).
- The help golden: it moves to `cmd/bench/help_inventory_test.go`, and the ticket binding registry names the new file.
- The routing table: it moves to `internal/conformance/subcommand_routing_table_test.go`, and the ticket binding registry names the new file. The routing check keeps reading `cmd/bench/main.go` as the dispatch file.
- `bench worktree path` stdout: unchanged. Its two tests read stderr only on failure (PB34 reads it on success).
- `CONTEXT.md` and `.bench/BENCH-reference.md`: the added entry and paragraph keep every anchored needle. The forbid needles `progressive loading` and `may inherit exact per-check evidence` do not appear in the new text.
- The census hook: the exec-child form names no pool path, so the census records nothing for a probe.

The shipped-surface claim words: the reference guide paragraph names no repo-only path beside a claim word.

The trust chain is unchanged. The wrapper resolves the Bench executable as it does for every verb. The verb starts no run child before its refusals pass. The focused run's owner selects the run binary and launches `go` from `PATH` through the gate's one test-argv producer, as `bench test` does today. The gate-lock refusal, the subject refusal, and the selection parse run before the preserved copy is written and before any run child starts (PB25).

Recorded runs:

The build ran the two review-owned rows in the integration worktree. Each run used
the worktree's own `dist/bench`, because the installed wrapper resolves the landed
executable, which does not hold the verb before the landing. Neither command names a
pool path.

- PB8, the check form. The command was `bench probe
  internal/conformance/subcommand_routing_table_test.go --omit '"probe":
  routed("internal/probe"),' --check subcommand-routing`. The omission removes the
  routing row, and the check reports the unrouted verb. The stdout was:

      probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
        bit,internal/conformance/subcommand_routing_table_test.go,omit,failed,1,yes
      packages[1]{package,status,elapsed_ms}:
        github.com/gibbonmi/bench/internal/conformance,fail,15
      failures[1]{package,test,line}:
        github.com/gibbonmi/bench/internal/conformance,TestRootConformance,"gate_entry_test.go:33: gate: cmd/bench/main.go dispatches \"probe\" with no entry in the subcommand argument-routing registry; record it as routed through usage.Parse or as an exemption with its reason"
      skips[0]{package,test,reason}:

- PB11, the exec form. The command was `bench worktree exec bench-probe --
  ./dist/bench probe internal/probe/probe.go --swap 'return "bit", 0' --with 'return
  "bit", 1' --package ./internal/probe --run '^TestVerdictExitCodes$'`. The swap
  breaks the verdict mapping, and the focused test reports it. The stdout was:

      probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
        bit,internal/probe/probe.go,swap,failed,1,yes
      packages[1]{package,status,elapsed_ms}:
        github.com/gibbonmi/bench/internal/probe,fail,3
      failures[1]{package,test,line}:
        github.com/gibbonmi/bench/internal/probe,TestVerdictExitCodes/failed,"outcome_test.go:325: verdictFor(\"failed\") = (\"bit\", 1), want (\"bit\", 0)"
      skips[0]{package,test,reason}:

The review round runs `codex exec` with the reviewer-named model `gpt-6-astra` at high effort with a cap of two iterations. The reviewer named it for this run on 2026-09-06. Every subagent runs `opus` at low or medium effort.
