# A kit landing that changes no Go file runs the reader set

Status: staged

Roadmap: FT449

Decision source: `specs/no-go-schedule/decisions/no-go-schedule.md` (ready compiled map).

Verification log: 4 iteration(s) to accept — Opus/medium accepted the spec in 2 iterations, after repairs of the chunk rows, the missing-log refusal, and the `-exec` quoting. Opus/medium accepted the tickets in 2 iterations, after repairs of the old-argv pins and the reader-record query.

## Problem

From 2026-10-18 to 2026-10-25, 48 of 71 landings on `main` changed no Go file. Each of them still ran the whole gate.
On 2026-10-25 the whole gate ran its phases one after another in about 268 seconds, and its `test` phase took 195 seconds.
The research asset of the decision map holds these measurements as Q1 and Q2.

Most of that cost buys nothing for such a landing. A Go test can change its outcome only through a non-Go path that it reads.
The gate has no record of which test package reads which path, so it cannot run less than the whole gate.

## Solution

Each whole run of the kit's built-in table records one Go test log for each test package. A green whole run derives the reader table from those logs.
The reader table maps each covered path family to the test packages that read it.

A kit landing whose every changed path sits in a covered family then runs the reader set instead of the whole gate.
The reader set runs the readers of the changed families, the root-conformance entry, and the `system` phase. It adds the `race` phase only when a race-table package is a reader.

A green reader set writes a reader record and authorizes its own landing. It never moves the green marker. Five such landings, a red whole run, or a changed oracle send the next landing to the whole gate.
The ambient dashboard shows a quiet `reader green (<n> of 5)` row for such a tip.

## User stories

Line: opus / medium.
Implementation-line reason: the spec fixes every predicate and the gate tests cover each one, but the trace and the landing schedule are new seams in the oracle. The harder chunks move one effort step up. NG-C6 edits guidance prose, so the leverage override also routes it to high effort.
Harder chunks: NG-C2, NG-C3, NG-C6.

### Trace and reader table

1. As a kit maintainer, I want each whole run to record one Go test log for each test package, so that the reader table comes from observed reads.
2. As a kit maintainer, I want the `test` phase to stay one `go test` call, so that Go keeps its package scheduling.
3. As a kit maintainer, I want a green whole run to map each package to the covered families it read, so that a new file in a known family is covered.
4. As a kit maintainer, I want a read through a tracked symbolic link to count as a read of the target family, so that `.claude/` reads cover `.agents/`.
5. As a kit maintainer, I want a read inside the package's own directory to add no family, so that package-local data stays package-local.
6. As a kit maintainer, I want a malformed or empty test log to derive no table, so that a broken trace cannot authorize a reduced landing.
7. As a kit maintainer, I want a table that does not name `internal/conformance` as a reader of every covered family to be rejected, so that a short trace yields no table.
8. As a kit maintainer, I want a red whole run to derive no table, so that only a green tree supplies readers.
9. As a kit maintainer, I want the declared child-reader rows to join the table, so that a reader inside a child process runs too.
10. As a kit maintainer, I want the table to record the oracle identity of its run, so that a later landing detects a changed oracle.
11. As a kit maintainer, I want a whole run to remove its test logs after the derivation, so that the run leaves no residue.
12. As a kit maintainer, I want the table store to stay bounded, so that many whole runs do not fill the Git directory.
13. As a linked-repository maintainer, I want my `test` phase argv to stay unchanged, so that the kit trace never reaches my gate.

### Oracle identity

14. As a kit maintainer, I want one function to compute the tree-free oracle identity beside the oracle digest, so that the two cannot drift.
15. As a kit maintainer, I want the identity to include the capability posture, so that a change of `BENCH_REQUIRE_CAPABILITIES` forces a whole run.
16. As a kit maintainer, I want the full-verdict oracle digest to keep its value, so that every current green stays reusable.

### The reader set at a landing

17. As a landing operator, I want a kit landing whose every changed path sits in a covered family to run the reader set, so that a no-Go landing does not pay the whole gate.
18. As a landing operator, I want a changed path outside the covered families to run the whole gate, so that an unmapped input never escapes grading.
19. As a landing operator, I want a changed embed target to run the whole gate, so that a compile-time read never escapes grading.
20. As a landing operator, I want a changed shellcheck-inventory path to run the whole gate, so that the shell lint never escapes grading.
21. As a landing operator, I want a changed symbolic link or gitlink to run the whole gate, so that a path whose bytes name another place never escapes grading.
22. As a landing operator, I want a rename out of a covered family to run the whole gate, so that the new path is graded.
23. As a landing operator, I want the reader set to run the readers of each changed family and `internal/conformance`, so that every dev-tier check runs.
24. As a landing operator, I want the reader set to run the declared child readers of each changed family, so that a child-process read is graded.
25. As a landing operator, I want the reader set to run the `system` phase always, so that the Bench children in the system journeys grade the change.
26. As a landing operator, I want the reader set to run the whole `race` table when a race-table package is a reader, so that the race sentinels stay whole.
27. As a landing operator, I want the reader set to run no `gofmt`, `vet`, or `shellcheck` phase, so that the reader set costs only what the change can affect.
28. As a landing operator, I want a missing reader table for the marker's tree to run the whole gate, so that the first landing after the install builds the table.
29. As a landing operator, I want a changed oracle identity to run the whole gate, so that a toolchain change never trusts an old table.
30. As a landing operator, I want the sixth landing after five reader-record landings to run the whole gate, so that the chain of trust stays short.
31. As a landing operator, I want a chain commit without a reader record to run the whole gate, so that an unrecorded link never extends the chain.
32. As a landing operator, I want an absent green marker to run the whole gate, so that no chain starts without a whole green.
33. As a landing operator, I want the gate to name the first case that refused the reader set, so that I know why a landing paid the whole gate.
34. As a landing operator, I want only `bench worktree land` to run the reader set, so that a merge, a plain gate, and a checkpoint keep their schedules.
35. As a linked-repository maintainer, I want my landing to run my whole table, so that the kit's reader table never grades my tree.
36. As a landing operator, I want a green reader set to write a reader record and no full verdict, so that no other reader reuses it.
37. As a landing operator, I want a green reader set to authorize its landing, so that the landing publishes.
38. As a landing operator, I want the reader landing to leave the green marker at the last whole green, so that the marker keeps its meaning.
39. As a landing operator, I want a resumed reader landing to leave the green marker unchanged, so that a resume does not promote a reader green.
40. As a landing operator, I want a red reader set to refuse the landing as a red, so that the repair route stays the same.
41. As a landing operator, I want the green verdict line to name the reader set, so that nobody reads it as a whole green.
42. As a release operator, I want a reader record never to count as reusable green, so that the release preflight and the delta base keep their whole-green requirement.
43. As the reviewer, I want one measured reader landing beside one measured whole landing, so that the premise is checked before the build completes.

### Red after reader landings

44. As a landing operator, I want a red whole landing run to write a reader hold that names the green marker, so that every later landing runs the whole gate until the marker moves.
45. As a landing operator, I want a whole green landing to retire the hold, so that the reader set returns after a repair.
46. As a landing operator, I want the red refusal to name each reader-record landing since the green marker, so that I know the suspects.
47. As a landing operator, I want a red reader set to write no hold, so that a candidate red does not close the reader set for other landings.

### The ambient dashboard

48. As a session, I want the gate signal to show `reader green (<n> of 5)` with no command for a reader-record tip, so that I do not run a whole gate by reflex.
49. As a session, I want a tip with no reader record to keep the stale row, so that a real staleness still routes to `bench gate`.

### Docs

50. As a teammate, I want ADR 0017 to state the reader-set exception, so that the decision record matches the gate.
51. As a teammate, I want the profile and the reference to state the reader set, the reader record, the reader hold, and the traced `test` argv, so that the docs match the gate.
52. As a teammate, I want `CONTEXT.md` to define the covered family and the reader hold, so that every session uses one name for each.

## Implementation decisions

### Terms and boundaries

`CONTEXT.md` owns the reader table, the reader set, and the reader record. This build adds two terms there. A **covered family** is one of the eleven path families in the reader-family declaration. A **reader hold** is the record that a red whole landing run writes.

The reader-family declaration is one Go table in `internal/gate`. It holds eleven families: `specs/`, `roadmap/`, `ROADMAP.md`, `capture/`, `reviews/`, `decisions/`, `CHANGELOG.md`, `.bench/commitment.json`, `.agents/`, `.bench/BENCH.md`, and `docs/`. A family that ends with `/` claims every path below it. Each other family claims its exact path. The commit lane's path classes stay a separate fact, because they select lane checks.

The reader set applies only when the graded root is the kit and the gate runs the kit's built-in table. The test is the one that the delta set uses today: `narrowsKitTable` and `RootConformanceRuns`.

### The oracle identity

`buildSubjectForTree` computes one more digest in the same pass. The identity digest receives every frame of the oracle digest except the tree, and then the `BENCH_REQUIRE_CAPABILITIES` posture. The oracle digest keeps its current frames and order, so its value does not change.

The posture enters the identity only. A green recorded under one posture therefore stays reusable under the other posture, as today. The flagged additions list this choice.

### The trace

The trace applies when `RootConformanceRuns(root, kit)` holds, the test that `rootConformanceEnv` already uses. The kit's whole `test` phase argv then becomes these seven elements: `go`, `test`, `-trimpath`, `-count=1`, `-timeout=10m0s`, `-exec=<bench-run-binary> gate-go exec`, and `./...`. The sixth element is one element. A linked repository's `test` phase, the delta set, and the reader set keep their current argv.

`cmd/go` splits the `-exec` value with its quoted-string rules. `withRunBinary` therefore replaces the run-binary token inside the `-exec=` element with the selection path in quoted form, so a path with a space stays one word. Its second caller, the lane at `internal/gate/lane.go:306`, carries no `-exec=` element, so the lane is unchanged.

`bench gate-go exec <binary> <args...>` is the shim. When `BENCH_GATE_TESTLOG` names a directory, the shim creates the file `<dir>/<package key>.log`. It then replaces its image with `<binary> -test.testlogfile=<file> <args...>`. The package key is the package directory relative to the graded root, in an encoding that maps two different directories to two different names. Without `BENCH_GATE_TESTLOG`, the shim replaces its image with `<binary> <args...>`.

`internal/subprocess` owns the image replacement, as it owns the start gate. The shim creates the log file before the replacement, so a test binary that never starts leaves an empty file, not a missing one.

The gate owner creates the log directory for a whole run of the kit's built-in table and passes it in the phase environment. After the run, the owner derives the table, persists it on green, and removes the directory.

`gateEnv` strips `BENCH_GATE_TESTLOG` and the reader schedule from every inherited environment, as it strips the delta schedule (`internal/gate/gate.go:160-167`). The owner sets them again on the one phase child that runs the table. So a nested gate run inside a test, such as a system journey, never writes into the outer run's logs and never runs the outer schedule.

### The reader table

The derivation is one pure function. Its output is a table or a refusal reason. Its inputs are these:

- the package logs
- the list of test packages: each package of the graded root that holds a test file, as `go list` reports it after the run with the build flags and tags of the traced `test` call
- the graded root
- the tracked symbolic links of the graded tree
- the reader-family declaration
- the declared child-reader rows

For each log, the derivation applies these rules in order:

- The first line must be `# test log`. Each other line must be `<op> <path>`, where `<op>` is `open`, `stat`, `chdir`, or `getenv`. Any other line refuses the table.
- A relative path resolves against the package directory, and each `chdir` line changes that directory.
- A path outside the graded root adds nothing. A path inside the package's own directory adds nothing.
- A path whose prefix is a tracked symbolic link resolves to the link's target. The target's family counts.
- Each remaining `open` or `stat` path adds its covered family, if it has one, to the package's reader set. A `chdir` line and a `getenv` line add no family.

The derivation refuses the table in each of these cases:

- A listed test package has no log.
- A log is empty.
- `internal/conformance` is not a reader of every covered family.

Go's test log drops a name that holds a newline, so such a read leaves no line. Its family still counts through the package's other reads, or the witness refuses the table.

The declared child-reader list sits beside the derivation. Each row is one pair of a test package and a covered family. Today it holds one row: `internal/releaseevidence` with `CHANGELOG.md`.

A green whole run persists the table in the common Git directory, keyed by the graded tree. The record holds the tree, the identity digest of the run, and the readers of each family. The store keeps at most 32 tables. It removes the oldest first, and it never removes the table of the green marker's tree.

A failed derivation persists no table and leaves the whole green unchanged. A red whole run derives nothing.

### The reader set schedule

The gate owner passes the reader set to the phase child in one environment entry, as the delta set does. The phase child then runs these phases in this order:

- `test` over the reader packages, with the root-conformance environment
- `race` with its whole table, when the schedule flags it
- `system`

The reader packages are the union of the table's readers of each changed family, the child readers of each changed family, and `./internal/conformance`. Each package appears once. The race flag holds when a package of the race-test table is among the reader packages.

A green reader set prints its phase table and then `gate: green (reader set)`. A malformed reader schedule refuses before any phase.

### The reader record

The reader record is a new record class with its own reuse refusal. It carries the graded tree, the destination tree, the green marker commit, the reader schedule, the run binary, and the record time. The gate writes one reader record for each graded tree, in a store keyed by the tree. A green reader set writes no full verdict, no pending record, and no evidence.

One exported query, `gate.ReaderRecorded(root, tree)`, answers whether the store holds a valid reader record for a tree. The resume marker policy and the gate inspection read the store only through it.

### The landing plan

`landReviewed` puts the destination commit and the green marker in the gate context, as `gate.WithCompletion` does for a completion. Only that call sets the value. A merge, a plain `bench gate`, a `--fresh` run, and a checkpoint carry no landing value, so they keep their schedules.

The reader plan reads the landing value and runs these cases in order. The first case that holds sends the run to the whole gate:

1. The run is not a kit run of the built-in table.
2. The run carries no landing value.
3. The green marker is absent, or it is not an ancestor of the destination.
4. A changed path between the destination tree and the graded tree sits outside every covered family.
5. A changed path matches an embed directive of any Go file of the composed tree, test files included, or it is a path of the shellcheck inventory. A path matches a directive when it equals the directive's path, sits below it, or matches its glob. Today `EmbedTargets` reads only non-test files below the root, `cmd/`, and `internal/` (`internal/packagesurface/assets.go:122-160`), so the plan needs a scan that also reads test files and every other directory.
6. A changed path has a symbolic-link or gitlink mode on its old side or on its new side.
7. The store holds no valid reader table for the green marker's tree.
8. The table's identity differs from the identity of this run.
9. A reader hold names the current green marker.
10. A first-parent commit after the green marker, up to the destination, has no reader record of its own tree.
11. Five first-parent commits follow the green marker.

When a case holds, the gate prints `gate: whole gate (reader set refused: <case>)` before the phases. Otherwise the run is the reader set.

### Authorization and publication

`gate.Result` states the schedule that ran and the suspects of a red whole landing run. `authorization` maps a green reader set to the new kind `reader green`, and a red reader set to `candidate`. The kind `reader green` carries no evidence token.

A red reader set is `candidate` because the chain proof stands in for the green baseline. Each commit from the green marker to the destination passed its own readers, so a red of those readers on the composed tree belongs to the diff. The flagged additions list this choice.

The reviewed landing publishes on `green` and on `reader green`. The worktree-commit landing keeps `green` and `lane pass`. After a `reader green` publication, the landing skips the green-marker advance. The resume marker policy gains one fact: when the published tree holds a reader record, the resume accepts the marker and advances nothing.

### The reader hold

A red whole run with a landing value writes the reader hold. The hold names the green marker commit of that run. A later reader plan reads the hold, and the hold applies only while the marker is that commit. A whole green landing moves the marker, so the hold stops applying with no write.

The same red run names its suspects: each first-parent commit after the green marker, up to the destination, that holds a reader record. `refusalMessage` appends `; reader landings since the green marker: <commits>` when the list is not empty. A red reader set writes no hold.

### The ambient dashboard

`gate.Inspect` reports two more facts: whether the current tree holds a reader record, and the count of reader-record commits after the green marker up to HEAD. When the gate row would read stale and the current tree holds a reader record, the row reads `reader green (<n> of 5)` with no command. It takes the stale row's rank, and no stale row prints.

### Structure headroom

`internal/gate/run_transaction.go` and `internal/worktree/land.go` hold 397 lines, `internal/landing/landing.go` holds 392, and `internal/gate/phases.go` holds 388. The structure growth check reds a file over 400 lines that grows. Each ticket therefore puts new logic in a new file and keeps each of these files at 400 lines or fewer.

## Implementation chunks

Each chunk is one coherent outcome with one review checkpoint. A chunk's predecessor supplies the seam that the chunk consumes. The tickets of a chunk commit in their listed order.

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| NG-C1 / `01-compute-the-oracle-identity.md` | One pass computes the tree-free oracle identity beside the oracle digest | NG1, NG2, NG3, NG4 | Completion-plan checks | no |
| NG-C2 / `02-record-a-test-log-per-package.md`, `03-derive-the-reader-table.md` | Each whole kit run records test logs and persists a reader table on green | NG5, NG6, NG7, NG8, NG9, NG10, NG11, NG12, NG13, NG14, NG15, NG16, NG17, NG18, NG19, NG20, NG21, NG22, NG23, NG24, NG25, NG26, NG27, NG28, NG71, NG72, NG81 | Completion-plan checks | yes |
| NG-C3 / `04-run-the-reader-set-schedule.md`, `05-plan-the-reader-set-at-a-landing.md`, `06-publish-a-reader-landing.md` | A covered kit landing runs the reader set, publishes, and keeps the marker | NG29, NG30, NG31, NG32, NG33, NG34, NG35, NG36, NG37, NG38, NG39, NG40, NG41, NG42, NG43, NG44, NG45, NG46, NG47, NG48, NG49, NG50, NG51, NG52, NG53, NG54, NG55, NG56, NG57, NG58, NG59, NG60, NG73, NG74, NG75, NG76, NG77, NG79, NG80 | Completion-plan checks | yes |
| NG-C4 / `07-hold-the-whole-gate-after-a-red.md` | A red whole landing run holds the whole gate and names its suspects | NG61, NG62, NG63, NG64, NG78 | Completion-plan checks | no |
| NG-C5 / `08-show-reader-green-on-the-dashboard.md` | The gate signal shows a reader-record tip as reader green | NG65, NG66, NG67 | Completion-plan checks | no |
| NG-C6 / `09-amend-the-landing-rule-docs.md` | The ADR, the profile, the reference, and the glossary state the reader set | NG68, NG69, NG70 | Completion-plan checks | yes |

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"ft449-no-go-schedule-full-20261025","orchestrator_session":"claude:ngs-orchestrator-20261025","author_limit":1,"assignments":{"01-compute-the-oracle-identity.md":[{"session":"claude:ngs-t01-author-20261025","assignment":"760876ee2e49120ba1b88c516f065700","model":"opus","effort":"medium","source":"f80bbf194a15f02108b0be096e6140b7c5942553","native_ref":"claude-agent:bench-writer-medium:ngs-t01-author-20261025"}],"02-record-a-test-log-per-package.md":[],"03-derive-the-reader-table.md":[],"04-run-the-reader-set-schedule.md":[],"05-plan-the-reader-set-at-a-landing.md":[],"06-publish-a-reader-landing.md":[],"07-hold-the-whole-gate-after-a-red.md":[],"08-show-reader-green-on-the-dashboard.md":[],"09-amend-the-landing-rule-docs.md":[]}},"chunks":[{"id":"NG-C1","tickets":["01-compute-the-oracle-identity.md"],"verification":[{"id":"ticket-1-check-1","command":"bench test --package ./internal/gate","ticket":"01-compute-the-oracle-identity.md"},{"id":"ticket-1-check-2","command":"bench test --check root-conformance","ticket":"01-compute-the-oracle-identity.md"},{"id":"ticket-1-probe","command":"bench test --package ./internal/gate","probe":"In buildSubjectForTree, frame the tree into the identity digest too. TestOracleIdentityIgnoresTheTree must fail, then pass after the restore.","ticket":"01-compute-the-oracle-identity.md"}]},{"id":"NG-C2","tickets":["02-record-a-test-log-per-package.md","03-derive-the-reader-table.md"],"verification":[{"id":"ticket-2-check-1","command":"bench test --package ./internal/gate","ticket":"02-record-a-test-log-per-package.md"},{"id":"ticket-2-check-2","command":"bench test --package ./internal/subprocess","ticket":"02-record-a-test-log-per-package.md"},{"id":"ticket-2-check-3","command":"bench test --check ordinary-build-census","ticket":"02-record-a-test-log-per-package.md"},{"id":"ticket-2-check-4","command":"bench test --package ./internal/conformance","ticket":"02-record-a-test-log-per-package.md"},{"id":"ticket-2-check-5","command":"bench test --check root-conformance","ticket":"02-record-a-test-log-per-package.md"},{"id":"ticket-2-probe","command":"bench test --package ./internal/gate","probe":"In the gate-go exec shim, drop the -test.testlogfile operand. TestGateGoExecAddsTheTestLog must fail, then pass after the restore.","ticket":"02-record-a-test-log-per-package.md"},{"id":"ticket-3-check-1","command":"bench test --package ./internal/gate","ticket":"03-derive-the-reader-table.md"},{"id":"ticket-3-check-2","command":"bench test --package ./internal/git","ticket":"03-derive-the-reader-table.md"},{"id":"ticket-3-check-3","command":"bench test --check root-conformance","ticket":"03-derive-the-reader-table.md"},{"id":"ticket-3-probe","command":"bench test --package ./internal/gate","probe":"In the reader-table derivation, remove the internal/conformance witness refusal. TestReaderTableRequiresTheConformanceWitness must fail, then pass after the restore.","ticket":"03-derive-the-reader-table.md"}]},{"id":"NG-C3","tickets":["04-run-the-reader-set-schedule.md","05-plan-the-reader-set-at-a-landing.md","06-publish-a-reader-landing.md"],"verification":[{"id":"ticket-4-check-1","command":"bench test --package ./internal/gate","ticket":"04-run-the-reader-set-schedule.md"},{"id":"ticket-4-check-2","command":"bench test --package ./internal/git","ticket":"04-run-the-reader-set-schedule.md"},{"id":"ticket-4-check-3","command":"bench test --check root-conformance","ticket":"04-run-the-reader-set-schedule.md"},{"id":"ticket-4-probe","command":"bench test --package ./internal/gate","probe":"In the reader run, write the reader record before the status check. TestReaderRedWritesNoRecord must fail, then pass after the restore.","ticket":"04-run-the-reader-set-schedule.md"},{"id":"ticket-5-check-1","command":"bench test --package ./internal/gate","ticket":"05-plan-the-reader-set-at-a-landing.md"},{"id":"ticket-5-check-2","command":"bench test --package ./internal/git","ticket":"05-plan-the-reader-set-at-a-landing.md"},{"id":"ticket-5-check-3","command":"bench test --package ./internal/packagesurface","ticket":"05-plan-the-reader-set-at-a-landing.md"},{"id":"ticket-5-check-4","command":"bench test --check root-conformance","ticket":"05-plan-the-reader-set-at-a-landing.md"},{"id":"ticket-5-probe","command":"bench test --package ./internal/gate","probe":"In the reader plan, change the bound of five first-parent commits to six. TestReaderPlanRunsWholeAfterFiveReaderLandings must fail, then pass after the restore.","ticket":"05-plan-the-reader-set-at-a-landing.md"},{"id":"ticket-6-check-1","command":"bench test --package ./internal/gate/authorization","ticket":"06-publish-a-reader-landing.md"},{"id":"ticket-6-check-2","command":"bench test --package ./internal/landing","ticket":"06-publish-a-reader-landing.md"},{"id":"ticket-6-check-3","command":"bench test --package ./internal/worktree","ticket":"06-publish-a-reader-landing.md"},{"id":"ticket-6-check-4","command":"bench test --package ./internal/worktree/landingpolicy","ticket":"06-publish-a-reader-landing.md"},{"id":"ticket-6-check-5","command":"bench test --package ./cmd/bench","ticket":"06-publish-a-reader-landing.md"},{"id":"ticket-6-check-6","command":"bench test --check root-conformance","ticket":"06-publish-a-reader-landing.md"},{"id":"ticket-6-probe","command":"bench test --package ./internal/worktree","probe":"In the reviewed landing, advance the green marker after a reader green publication. TestReaderLandingLeavesTheGreenMarker must fail, then pass after the restore.","ticket":"06-publish-a-reader-landing.md"}]},{"id":"NG-C4","tickets":["07-hold-the-whole-gate-after-a-red.md"],"verification":[{"id":"ticket-7-check-1","command":"bench test --package ./internal/gate","ticket":"07-hold-the-whole-gate-after-a-red.md"},{"id":"ticket-7-check-2","command":"bench test --package ./internal/gate/authorization","ticket":"07-hold-the-whole-gate-after-a-red.md"},{"id":"ticket-7-check-3","command":"bench test --package ./internal/landing","ticket":"07-hold-the-whole-gate-after-a-red.md"},{"id":"ticket-7-check-4","command":"bench test --package ./internal/git","ticket":"07-hold-the-whole-gate-after-a-red.md"},{"id":"ticket-7-check-5","command":"bench test --check root-conformance","ticket":"07-hold-the-whole-gate-after-a-red.md"},{"id":"ticket-7-probe","command":"bench test --package ./internal/gate","probe":"In the reader plan, skip the reader-hold case. TestReaderPlanRunsWholeUnderAHold must fail, then pass after the restore.","ticket":"07-hold-the-whole-gate-after-a-red.md"}]},{"id":"NG-C5","tickets":["08-show-reader-green-on-the-dashboard.md"],"verification":[{"id":"ticket-8-check-1","command":"bench test --package ./internal/status","ticket":"08-show-reader-green-on-the-dashboard.md"},{"id":"ticket-8-check-2","command":"bench test --package ./internal/gate","ticket":"08-show-reader-green-on-the-dashboard.md"},{"id":"ticket-8-check-3","command":"bench test --check root-conformance","ticket":"08-show-reader-green-on-the-dashboard.md"},{"id":"ticket-8-probe","command":"bench test --package ./internal/status","probe":"In appendGateInfo, drop the reader-green branch. TestGateRowShowsReaderGreen must fail, then pass after the restore.","ticket":"08-show-reader-green-on-the-dashboard.md"}]},{"id":"NG-C6","tickets":["09-amend-the-landing-rule-docs.md"],"verification":[{"id":"ticket-9-check-1","command":"bench test --check docs-currency-workflow","ticket":"09-amend-the-landing-rule-docs.md"},{"id":"ticket-9-check-2","command":"bench test --check prose-mechanics","ticket":"09-amend-the-landing-rule-docs.md"},{"id":"ticket-9-check-3","command":"bench test --check root-conformance","ticket":"09-amend-the-landing-rule-docs.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/no-go-schedule/spec.md"},{"id":"census","command":"bench test --check ordinary-build-census"},{"id":"system","command":"bench test --check system"}]}
```

## Testing decisions

A good test drives the owner seam with exact inputs and observes its output: a table, a schedule, a record, a kind, or a rendered row. The derivation takes canned logs, so its tests need no Go test run. The phase-child tests follow `TestDeltaCheckpointRunsTheDeltaSet`, which drives the phase child with an environment schedule over a fixture module.

The plan tests build a real Git fixture with a `main` branch, a green-marker ref, chain commits, and stored records. Each refusal case gets its own fixture change, so each row reaches exactly one case.

The landing tests replace the authorizer, as `internal/landing/landing_test.go` does today. The end-to-end proof of the real trace and the real table is the measured landing of row NG60. It runs in a disposable clone of the integration source, because the open build holds the landing window of `main`.

The gate's own whole run grades this build at every chunk checkpoint and at the landing. After the build lands, the first covered landing on `main` runs the whole gate, because no table exists yet. That run derives the first table.

### Seam diagram

    whole kit run ──▶ test phase ──▶ gate-go exec shim ──▶ one test log per package
        │                                                        │
        ▼                                                        ▼
    green verdict ──▶ [ reader-table derivation ] ──▶ reader table keyed by tree
                          ◀ tests attach here: canned logs in, table or refusal out

    bench worktree land ──▶ landReviewed ──▶ gate context: destination, marker
        │
        ▼
    [ reader plan ] ──▶ whole gate (first refusing case named) | reader set schedule
        ◀ tests attach here: Git fixture with marker, chain, records, table
        │
        ▼
    reader set green ──▶ reader record ──▶ authorization `reader green` ──▶ publish, marker stays

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
| --- | --- | --- | --- | --- |
| NG1 | 14 | Two subjects that differ only in their tree have one identity digest and two oracle digests | planned TestOracleIdentityIgnoresTheTree in internal/gate/oracle_identity_test.go | An identity that frames the tree differs across the two subjects. |
| NG2 | 15 | Setting `BENCH_REQUIRE_CAPABILITIES=1` changes the identity digest of one subject | planned TestOracleIdentityCarriesTheCapabilityPosture in internal/gate/oracle_identity_test.go | An identity without the posture stays equal across the two postures. |
| NG3 | 14 | Changing the bytes of a declared tool changes both the identity digest and the oracle digest | planned TestOracleIdentityAndOracleShareTheClosure in internal/gate/oracle_identity_test.go | A second derivation that skips the closure keeps one digest equal. |
| NG4 | 16 | The oracle digest keeps every frame and its order from before this build | review-owned: the diff of `buildSubjectForTree` adds frames only to the identity digest | A reordered or added oracle frame changes the digest of every current green. |
| NG5 | 1 | With `BENCH_GATE_TESTLOG` set, the shim runs the binary with `-test.testlogfile=<dir>/<package key>.log` before its arguments | planned TestGateGoExecAddsTheTestLog in internal/gate/testlog_exec_test.go | A shim that drops the flag leaves the probe binary's recorded argv without it. |
| NG6 | 1 | Without `BENCH_GATE_TESTLOG`, the shim runs the binary with its arguments unchanged | planned TestGateGoExecWithoutALogDirectory in internal/gate/testlog_exec_test.go | A shim that always adds the flag changes the recorded argv. |
| NG7 | 1 | The shim creates the log file before it replaces its image, so a binary that exits at once leaves an empty file | planned TestGateGoExecCreatesTheLogFirst in internal/gate/testlog_exec_test.go | A shim that leaves creation to the binary leaves no file for a binary that exits at once. |
| NG8 | 1 | Two package directories that differ only by a separator map to two different log names | planned TestTestLogKeysAreDistinct in internal/gate/testlog_exec_test.go | A key that maps `/` to `__` gives `a/b` and `a__b` one name. |
| NG9 | 2 | The kit's whole `test` phase argv is the seven elements `go`, `test`, `-trimpath`, `-count=1`, `-timeout=10m0s`, `-exec=<bench-run-binary> gate-go exec`, and `./...` | planned TestOrdinaryBuildCensusRequiresTheTracedTestArgv in internal/conformance/ordinary_build_census_test.go | The census check compares the elements exactly, so the old argv and a split `-exec` element both red. |
| NG10 | 2 | `withRunBinary` replaces the run-binary token inside the `-exec=` operand | planned TestWithRunBinaryReplacesTheExecToken in internal/gate/testlog_exec_test.go | A substitution of the first element alone leaves the token in the operand. |
| NG11 | 13 | A linked repository's `test` phase argv carries no `-exec=` operand | planned TestLinkedTestPhaseCarriesNoTrace in internal/gate/testlog_exec_test.go | A trace added to every Go module reaches the linked root. |
| NG12 | 1 | A phase-child run of the whole kit table with `BENCH_GATE_TESTLOG` set leaves one log that starts with `# test log` for each test package | planned TestWholeTestPhaseWritesOneLogPerPackage in internal/gate/testlog_exec_test.go | A missing shim or a shared log name leaves fewer logs than packages. |
| NG13 | 3 | A log that opens `roadmap/FT1.md` makes its package a reader of `roadmap/` | planned TestReaderTableMapsAReadToItsFamily in internal/gate/reader_table_test.go | A table keyed by exact paths has no `roadmap/` row. |
| NG14 | 3 | A table derived from a read of `roadmap/FT1.md` names that package for a change of `roadmap/FT2.md` | planned TestReaderTableCoversANewPathInAKnownFamily in internal/gate/reader_table_test.go | An exact-path table names no reader for the new path. |
| NG15 | 4 | A read of `.claude/commands/x.md`, under a tracked link to `../.agents/commands`, makes the package a reader of `.agents/` | planned TestReaderTableResolvesATrackedSymlink in internal/gate/reader_table_test.go | Mapping the logged path as written finds no covered family. |
| NG16 | 5 | A read inside the package's own directory adds no family | planned TestReaderTableSkipsTheOwnDirectory in internal/gate/reader_table_test.go | A derivation without the rule adds the family of the package's own path. |
| NG17 | 3 | A relative path after a `chdir` line resolves against the new directory | planned TestReaderTableFollowsChdir in internal/gate/reader_table_test.go | A derivation that ignores `chdir` maps the path to no family. |
| NG18 | 3 | A path outside the graded root adds no family | planned TestReaderTableIgnoresPathsOutsideTheRoot in internal/gate/reader_table_test.go | A derivation that matches by suffix maps an outside `docs/` path. |
| NG19 | 6 | A log whose first line is not `# test log` refuses the table and names its package | planned TestReaderTableRefusesAMalformedLog in internal/gate/reader_table_test.go | A lenient parser returns a table from the malformed log. |
| NG20 | 6 | A log line with an operation other than `open`, `stat`, `chdir`, or `getenv` refuses the table | planned TestReaderTableRefusesAnUnknownOperation in internal/gate/reader_table_test.go | A parser that skips unknown lines returns a table. |
| NG21 | 6 | An empty log refuses the table | planned TestReaderTableRefusesAnEmptyLog in internal/gate/reader_table_test.go | A parser that reads an empty file as no reads returns a table. |
| NG22 | 7 | A table in which `internal/conformance` reads no `docs/` path is refused | planned TestReaderTableRequiresTheConformanceWitness in internal/gate/reader_table_test.go | A derivation without the witness returns the short table. |
| NG23 | 9 | The child row of `internal/releaseevidence` makes it a reader of `CHANGELOG.md` with no logged read | planned TestReaderTableJoinsTheChildReaders in internal/gate/reader_table_test.go | A table of logged reads alone names no reader of `CHANGELOG.md` there. |
| NG24 | 10 | A green whole kit run persists one table keyed by its tree, with the identity digest of the run | planned TestWholeGreenPersistsTheReaderTable in internal/gate/reader_store_test.go | A run that skips the persistence leaves the store without the tree's table. |
| NG25 | 8 | A red whole kit run persists no table for its tree | planned TestRedWholeRunPersistsNoReaderTable in internal/gate/reader_store_test.go | A persistence that ignores the status stores the red tree's table. |
| NG26 | 6 | A whole run whose derivation fails stays green and persists no table | planned TestFailedDerivationKeepsTheGreen in internal/gate/reader_store_test.go | A failed derivation that reds the run changes the verdict. |
| NG27 | 12 | The store holds at most 32 tables and keeps the table of the green marker's tree | planned TestReaderTableStoreBoundsItsTables in internal/gate/reader_store_test.go | An unbounded store keeps a 33rd table, and a pure age prune removes the marker's table. |
| NG28 | 11 | A whole run removes its log directory after the derivation | planned TestWholeRunRemovesItsTestLogs in internal/gate/reader_store_test.go | A run that keeps the directory leaves it after the verdict. |
| NG29 | 23 | A phase child with a reader schedule runs `test` over the named packages and `./internal/conformance`, then `system` | planned TestReaderScheduleRunsTheReaderSet in internal/gate/reader_set_test.go | A schedule that drops `./internal/conformance` or `system` prints another phase list. |
| NG30 | 41 | A green reader set prints `gate: green (reader set)` as its last line | planned TestReaderScheduleRunsTheReaderSet in internal/gate/reader_set_test.go | A renderer without the reader branch prints `gate: green`. |
| NG31 | 26 | A reader schedule with the race flag also runs `race` with every race-table sentinel | planned TestReaderScheduleRunsTheRaceTable in internal/gate/reader_set_test.go | A schedule that ignores the flag prints no `race` row. |
| NG32 | 27 | A reader schedule runs no `gofmt`, `vet`, or `shellcheck` phase | planned TestReaderScheduleRunsNoSourcePhase in internal/gate/reader_set_test.go | A schedule built from the whole table prints those rows. |
| NG33 | 23 | A malformed reader schedule refuses before any phase | planned TestReaderScheduleRefusesAMalformedSchedule in internal/gate/reader_set_test.go | A lenient schedule reader runs phases from a malformed entry. |
| NG34 | 36 | A green reader set writes one reader record for its tree and leaves the verdict cache and the evidence store unchanged | planned TestReaderGreenWritesOnlyTheReaderRecord in internal/gate/reader_set_test.go | A reader run that writes a full verdict changes the cache bytes. |
| NG35 | 42 | A reader record never makes its tree reusable green | planned TestReaderRecordIsNeverReusableGreen in internal/gate/reader_set_test.go | A class without a reuse refusal answers reusable green. |
| NG36 | 40 | A red reader set writes no reader record | planned TestReaderRedWritesNoRecord in internal/gate/reader_set_test.go | A record written before the status check exists after a red. |
| NG37 | 17 | A landing value whose changes all sit in covered families, with a valid table, an equal identity, and a clean chain, runs the reader set | planned TestLandingRunsTheReaderSet in internal/gate/reader_plan_test.go | A plan that never selects the reader set runs the whole table. |
| NG38 | 23 | The reader packages are the changed families' readers, their child readers, and `./internal/conformance`, each once | planned TestReaderPlanSelectsTheChangedFamilyReaders in internal/gate/reader_plan_test.go | A union that skips one source or repeats a package changes the package list. |
| NG39 | 26 | The race flag holds when a race-table package reads a changed family | planned TestReaderPlanFlagsTheRaceTable in internal/gate/reader_plan_test.go | A plan without the flag rule leaves the flag off. |
| NG40 | 18 | A change to `AGENTS.md` runs the whole gate | planned TestReaderPlanRunsWholeForAnUncoveredPath in internal/gate/reader_plan_test.go | A plan that ignores uncovered paths selects the reader set. |
| NG41 | 22 | A rename of `docs/x.md` to `x.md` runs the whole gate | planned TestReaderPlanRunsWholeForARenameOutOfAFamily in internal/gate/reader_plan_test.go | A plan that reads only the old path selects the reader set. |
| NG42 | 19 | A change to `docs/x.txt`, which a fixture Go file embeds by its exact path, runs the whole gate | planned TestReaderPlanRunsWholeForAnEmbedTarget in internal/gate/reader_plan_test.go | A plan without the embed case selects the reader set. |
| NG43 | 20 | A change to a path of an injected shellcheck inventory runs the whole gate | planned TestReaderPlanRunsWholeForAShellcheckPath in internal/gate/reader_plan_test.go | A plan without the inventory case selects the reader set. |
| NG44 | 21 | A change that adds `docs/link` as a symbolic link runs the whole gate | planned TestReaderPlanRunsWholeForAnAddedSymlink in internal/gate/reader_plan_test.go | A plan without the link case selects the reader set. |
| NG45 | 21 | A change that adds `docs/module` as a gitlink runs the whole gate | planned TestReaderPlanRunsWholeForAnAddedGitlink in internal/gate/reader_plan_test.go | A plan without the gitlink case selects the reader set. |
| NG46 | 28 | A store without a table for the green marker's tree runs the whole gate | planned TestReaderPlanRunsWholeWithoutATable in internal/gate/reader_plan_test.go | A plan that falls back to an empty table selects the reader set. |
| NG47 | 29 | A table whose identity digest differs from the run's runs the whole gate | planned TestReaderPlanRunsWholeForANewOracle in internal/gate/reader_plan_test.go | A plan that skips the identity comparison selects the reader set. |
| NG48 | 30 | Five reader-record commits after the green marker run the whole gate | planned TestReaderPlanRunsWholeAfterFiveReaderLandings in internal/gate/reader_plan_test.go | An off-by-one bound of six selects the reader set. |
| NG49 | 30 | Four reader-record commits after the green marker run the reader set | planned TestReaderPlanRunsTheReaderSetAfterFourReaderLandings in internal/gate/reader_plan_test.go | An off-by-one bound of four runs the whole gate. |
| NG50 | 31 | A first-parent commit after the green marker without a reader record runs the whole gate | planned TestReaderPlanRunsWholeForAChainGap in internal/gate/reader_plan_test.go | A plan that counts commits alone selects the reader set. |
| NG51 | 32 | An absent green marker runs the whole gate | planned TestReaderPlanRunsWholeWithoutAMarker in internal/gate/reader_plan_test.go | A plan that treats an absent marker as the destination selects the reader set. |
| NG52 | 34 | The same covered change without a landing value runs the whole gate | planned TestReaderPlanNeedsALandingContext in internal/gate/reader_plan_test.go | A plan that ignores the landing value selects the reader set for a merge. |
| NG53 | 35 | A covered change in a linked root runs the whole table | planned TestReaderPlanRunsWholeForALinkedRoot in internal/gate/reader_plan_test.go | A plan without the kit test selects the reader set. |
| NG54 | 33 | A refused plan prints `gate: whole gate (reader set refused: <case>)` with the first case that held | planned TestReaderPlanNamesTheFirstRefusingCase in internal/gate/reader_plan_test.go | A plan that names the last case, or none, prints another line when two cases hold. |
| NG55 | 37 | A green reader set attributes the kind `reader green` with no evidence token | planned TestAuthorizeAttributesAReaderGreen in internal/gate/authorization/reader_kind_test.go | An authorizer that reads only the verdict cache attributes `infrastructure`. |
| NG56 | 40 | A red reader set attributes the kind `candidate` | planned TestAuthorizeAttributesAReaderRedToTheCandidate in internal/gate/authorization/reader_kind_test.go | An authorizer that reads only the verdict cache attributes `infrastructure`. |
| NG57 | 37 | The reviewed landing publishes on `reader green`, and the worktree-commit landing refuses it | planned TestOnlyTheReviewedLandingPublishesOnAReaderGreen in internal/landing/reader_publish_test.go | A policy that adds the kind to both lists publishes the worktree commit. |
| NG58 | 38 | After a `reader green` publication, `refs/bench/green/main` names the prior marker | planned TestReaderLandingLeavesTheGreenMarker in internal/worktree/land_reader_test.go | A landing that advances the marker after every publication moves the ref. |
| NG59 | 39 | The resume marker policy accepts a published tree with a reader record and advances nothing | planned TestResumeMarkerAcceptsAReaderLanding in internal/worktree/landingpolicy/reader_marker_test.go | A policy that reads only the destination position answers advance. |
| NG60 | 43 | A disposable clone of the integration source lands one `ROADMAP.md` change through the reader set after one whole-green landing, and the review pickup records both times | review-owned: the measurement procedure of ticket 06 | A reader set no faster than the whole gate refutes the premise, and the build returns to the spec. |
| NG61 | 44 | A red whole run with a landing value writes a reader hold that names the run's green marker | planned TestRedWholeLandingWritesTheReaderHold in internal/gate/reader_hold_test.go | A red path without the write leaves no hold. |
| NG62 | 44 | A reader hold that names the current green marker runs the whole gate | planned TestReaderPlanRunsWholeUnderAHold in internal/gate/reader_hold_test.go | A plan that ignores the hold selects the reader set. |
| NG63 | 45 | A reader hold that names an older green marker leaves the reader set selectable | planned TestReaderPlanIgnoresAHoldOfAnOlderMarker in internal/gate/reader_hold_test.go | A plan that reads any hold runs the whole gate. |
| NG64 | 46 | A refusal whose result carries two suspect commits names both in the refusal line, and an empty list keeps the current line | planned TestRefusalMessageNamesTheReaderLandings in internal/landing/reader_refusal_test.go | A renderer without the suffix prints neither commit, and an unconditional suffix changes the current line. |
| NG65 | 48 | A tip whose tree holds a reader record, two reader-record commits after the marker, shows `reader green (2 of 5)` with no command | planned TestGateRowShowsReaderGreen in internal/status/status_reader_test.go | A row without the reader branch prints the stale row and its `bench gate` route. |
| NG66 | 48 | The reader-green tip prints no `stale (` row | planned TestGateRowShowsReaderGreen in internal/status/status_reader_test.go | A row that adds the reader line beside the stale row prints both. |
| NG67 | 49 | A stale tip with no reader record keeps the stale row and its `bench gate` route | planned TestGateRowKeepsTheStaleRowWithoutAReaderRecord in internal/status/status_reader_test.go | A reader branch that reads any record replaces a real stale row. |
| NG68 | 50 | ADR 0017 states that a kit landing inside the reader table runs the reader set | review-owned: the ADR diff and `bench gate-prose` | An unamended ADR keeps the whole-gate-only rule. |
| NG69 | 51 | The profile and the reference state the reader set, the reader record, the reader hold, and the traced `test` argv | review-owned: the docs diff and `bench test --check docs-currency-workflow` | An unamended profile keeps the no-partition statement and the old argv. |
| NG70 | 52 | `CONTEXT.md` defines the covered family and the reader hold, each with an Avoid list | review-owned: the glossary diff | A missing entry leaves sessions to invent a name. |
| NG71 | 6 | Three listed test packages with two logs refuse the table and name the package without a log | planned TestReaderTableRefusesAMissingLog in internal/gate/reader_table_test.go | A derivation that reads only the logs it finds returns a table. |
| NG72 | 2 | With a run-binary path that holds a space, `go test` runs the shim at the whole path | planned TestTraceExecSurvivesASpaceInThePath in internal/gate/testlog_exec_test.go | An unquoted substitution makes `cmd/go` split the path and fail to start the shim. |
| NG73 | 19 | A change that deletes `docs/assets/a.txt`, below a directory that a fixture Go file embeds, runs the whole gate | planned TestReaderPlanRunsWholeBelowAnEmbeddedDirectory in internal/gate/reader_plan_test.go | An exact-path match misses a file below the embedded directory. |
| NG74 | 19 | A change to `docs/a.txt`, which a fixture embed glob `docs/*.txt` matches, runs the whole gate | planned TestReaderPlanRunsWholeForAnEmbedGlob in internal/gate/reader_plan_test.go | An exact-path match misses a glob match. |
| NG75 | 21 | A change that replaces the symbolic link `docs/link` with a regular file runs the whole gate | planned TestReaderPlanRunsWholeForARemovedSymlink in internal/gate/reader_plan_test.go | A plan that reads only the new mode misses the old link. |
| NG76 | 21 | A change that removes the gitlink `docs/module` runs the whole gate | planned TestReaderPlanRunsWholeForARemovedGitlink in internal/gate/reader_plan_test.go | A plan that reads only the new mode misses the old gitlink. |
| NG77 | 19 | A change to `docs/t.txt`, which a fixture `_test.go` file embeds, runs the whole gate | planned TestReaderPlanRunsWholeForATestFileEmbed in internal/gate/reader_plan_test.go | A scan of non-test sources alone misses the test-file directive. |
| NG78 | 46 | A red whole landing run after two reader-record commits lists both commits as suspects in its result | planned TestRedWholeLandingListsTheSuspects in internal/gate/reader_hold_test.go | A red path without the chain walk returns no suspect. |
| NG79 | 39 | `gate.ReaderRecorded` answers true for a tree with a valid reader record and false for a tree without one | planned TestReaderRecordedAnswersForItsTree in internal/gate/reader_set_test.go | A query that reads any record, or none, answers the same for both trees. |
| NG80 | 34 | A phase child's environment holds no reader schedule that it inherited | planned TestPhaseChildrenDropTheReaderSchedule in internal/gate/reader_set_test.go | Without the strip, a nested gate run inside a reader set runs the reduced schedule. |
| NG81 | 1 | A phase child's environment holds no `BENCH_GATE_TESTLOG` that it inherited | planned TestPhaseChildrenDropTheTestLogDirectory in internal/gate/testlog_exec_test.go | Without the strip, a nested gate run writes into the outer run's logs. |

Not covered: story 24 — row NG38 covers the child readers inside the reader-package union, so no separate row exists.
Not covered: story 25 — row NG29 covers the `system` phase of every reader schedule, so no separate row exists.
Not covered: story 47 — row NG36 shows that a red reader set writes no record, and the hold writer runs only on the whole red path of row NG61.

### Edge inventory

| edge | disposition and evidence |
| --- | --- |
| The first covered landing after the install | NG46: no table exists, so the landing runs the whole gate. That green derives the first table. |
| A spec-flip completion landing whose every path is covered | The plan treats it like any landing. The completion proof check keeps running before the plan. |
| A landing with no changed path | The landing refuses `nothing to commit` before the gate, as today. |
| A deleted file in a covered family | The deletion is a changed path in its family, so its family readers run. NG13 and NG37 cover the family rule. |
| A path with a space in a covered family | Git's `-z` change list and the test log both keep the space. The plan tests use `docs/a file.md` as one changed path. |
| A run-binary path with a space | NG72: the quoted substitution keeps the path as one word. |
| A read of a name with a newline | Go's test log drops the name, so no line appears. The family still counts through the package's other reads, or the witness of NG22 refuses the table. |
| A test that runs its own binary again with its own arguments | The child inherits `-test.testlogfile` and writes the parent's log file. The derivation then refuses the table (NG19, NG21), or the child's reads drop out and the bound of five catches the miss. |
| An absent table store, or an empty one | NG46: both read as no table. |
| An absent hold, or a hold of an older marker | NG63: the plan reads no hold that applies. |
| Two landings at once | The landing lease and the gate lock serialize them, as today. Each record write happens under the gate lock. |
| A run interrupted after the derivation and before the verdict | The table persists only after the green verdict persists, so an interrupt leaves no table. |
| A linked repository | NG11 and NG53: the trace and the reader set stay kit-only. |
| A shellcheck-inventory path inside a covered family | No such path exists today. NG43 proves the case through an injected inventory. |
| A malicious candidate that forges its own test logs | Won't handle — the candidate's own tests already grade the candidate, and review owns a hostile change to the shim, as it owns a hostile test. |
| A reader that tolerates a failed read, or a read that depends on another family's content | Won't handle — the bound of five landings and the whole run after it detect such a miss, and `/bench-debug` repairs it (map decision #11). |

## Out of scope

- A reader set for a linked repository. It needs a table that each repository derives: about 12 edits and 3 gate runs.
- A reader set for `bench worktree merge`, a plain `bench gate`, or a checkpoint: about 6 edits and 2 gate runs.
- A per-check selection inside the conformance entry: about 15 edits and 3 gate runs.
- A tracer for child processes, such as `strace`: about 10 edits and 3 gate runs, and a host dependency.
- The reuse of a merge verdict by its landing, which FT314 owns.
- A landing queue, which FT450 owns.
- The reader-green wording in the handoff pin block (`internal/handoff/facts.go:247`) and in the dashboard page (`internal/dashboard/render.go:86`). Map decision #12 names the ambient dashboard only: about 4 edits and 1 gate run.

## Further notes

### Source verification

The author re-read every `## Sources` entry of the compiled map in this session, at `main` commit `8250a689`. No entry changed after the map landed. The research asset's measurements stay as the map recorded them, and their invalidation triggers stay with them.

The author read these current-code facts in this session:

- The delta seam runs before the pending record at `internal/gate/run_transaction.go:191-193`.
- The reviewed landing authorizes at `internal/landing/landing.go:272` and publishes only `green` at `:119`.
- `gate.WithCompletion` puts landing data in the gate context at `internal/landing/landing.go:270`.
- Only `internal/worktree/land.go:220` and `internal/worktree/land_resume.go:82` advance the green marker.
- The resume marker policy answers advance whenever the destination sits at the published commit (`internal/worktree/landingpolicy/landingpolicy.go:162-170`).
- The authorizer grants `green` only on a reusable full verdict (`internal/gate/authorization/authorization.go:149-169`).
- The kit `test` phase is one `go test ./...` (`internal/gate/gate_go.go:121-123`, `internal/gate/phases.go:101-111`).
- `withRunBinary` replaces the run-binary token in the first argv element only (`internal/gate/phases.go:166-181`).
- The oracle digest frames the tree (`internal/gate/subject.go:73`), and the identity frames end at `:123`.
- `cmd/go` passes `-test.testlogfile` only with the cache on and no `-exec` (`$GOROOT/src/cmd/go/internal/test/test.go:1546-1549`, Go 1.25.14).
- The green verdict line has one renderer (`internal/gate/report.go:94-101`).
- The gate row reads stale and routes to `bench gate` (`internal/status/status.go:424-427`, `:463-468`).
- `internal/subprocess` owns image replacement (`internal/subprocess/startgate.go:46`).
- `cmd/go` splits the `-exec` value with its quoted-string rules (`$GOROOT/src/cmd/go/internal/test/testflag.go:39`, `$GOROOT/src/cmd/go/internal/base/flag.go:20-27`).
- The prospective gate builds its run binary below the artifact root, a temporary directory (`internal/gate/prospective.go:28-40`).
- Go's test log drops a name that holds a newline (`$GOROOT/src/testing/internal/testdeps/deps.go:93-96`).
- `EmbedTargets` returns each directive's resolved path or pattern (`internal/packagesurface/assets.go:87-92`).
- The traced pass of 2026-10-25 shows `internal/conformance` reading all eleven covered families in process, as the research asset's T1 records. So the witness holds on the tree of that date.

### Reader and writer sweep

| fact the build changes | readers and writers | disposition |
| --- | --- | --- |
| The kit `test` argv | `internal/conformance/ordinary_build_census_test.go:39`, `internal/gate/branch_native_phases_test.go:104`, `internal/gate/lane_test.go:72-87`, `internal/gate/delta_checkpoint_test.go:34-50` with its stub run binary, `projects/benchkit.md:346` | Ticket 02 writes all five. |
| The record-class list | `internal/gate/verdict_registry_guard_test.go:30-42` pins the class count | Ticket 04 writes it. |
| The reviewed publish kinds | `internal/landing/lane_test.go:56-59` | Ticket 06 keeps both assertions and adds NG57 in a new file. |
| The authorization kinds | `internal/landing/merge.go:127`, `internal/worktree/merge_refusal.go:43-96`, `internal/commit/commit.go:223`, `internal/worktree/land.go:252` | A merge and a commit never receive `reader green`, so their switches stay unchanged. |
| The green-marker advance | `internal/worktree/land.go:220`, `internal/worktree/land_resume.go:82` | Ticket 06 writes both sites. |
| The gate verdict lines | `internal/gate/report.go:94-101`, `.bench/BENCH-reference.md:624-655` | Ticket 04 writes the renderer, and ticket 09 writes the reference. |
| The gate inspection | `internal/status/status.go:437-458` | Ticket 08 writes it. |
| The landing rule prose | `docs/adr/0017-the-worktree-commit-runs-the-fast-lane.md:17-20`, `:30`; `projects/benchkit.md:351-353`, `:393-397`, `:476-479`; `.bench/BENCH-reference.md:548` | Ticket 09 writes all of them. |
| The gate-go usage | `internal/gate/gate_go.go:44` | Ticket 02 writes it. |

The writers of the new stores are the gate owner alone, under the gate lock. The landing, the status reader, and the plan read them.

### Source clauses and coverage

| map decision | rows |
| --- | --- |
| #1 the terms | NG70 and the existing `CONTEXT.md` entries |
| #2 the covered families | NG13, NG37, NG40 |
| #3 the kit only | NG11, NG53 |
| #4 the landing only | NG52, NG80 |
| #5 the complete conformance entry | NG29, NG38 |
| #6 the reader record and the marker | NG34, NG35, NG58, NG59, NG79 |
| #7 the derived table | NG5 to NG28, NG71, NG72, NG81 |
| #8 the bounds | NG46 to NG51, NG62 |
| #9 the red after reader landings | NG61 to NG64, NG78 |
| #10 the child readers | NG23 |
| #11 the child-reader list and the `system` phase | NG23, NG29, NG38 |
| #12 the dashboard row | NG65, NG66, NG67 |
| #13 the forced whole gate and the phases | NG29 to NG32, NG41 to NG45, NG73 to NG77 |
| #14 the ADR and the profile | NG68, NG69 |

### Pre-review proof checklist

- `Cited symbols`: every cited symbol resolves at `8250a689`: `buildSubjectForTree`, `withRunBinary`, `testPhase`, `toolchainPhases`, `scheduledTable`, `deltaPlan`, `greenVerdict`, `narrowsKitTable`, `RootConformanceRuns`, `shellcheckFiles`, `raceTests`, `AuthorizeWithWriters`, `landReviewed`, `ResumeMarker`, `appendGateInfo`, and `GateVerdict`.
- `Import edges`: no new cross-package edge beyond `internal/gate` to `internal/subprocess`, which exists today.
- `Source-row clauses and occurrences`: the table above maps each map decision to its rows.
- `Promised field labels`: `gate: green (reader set)`, `gate: whole gate (reader set refused: <case>)`, `reader green (<n> of 5)`, the kind `reader green`, and the refusal suffix `; reader landings since the green marker: <commits>`.
- `Changed-function callers`: `withRunBinary` has two callers, `phasesCommandAtKitWithSelection` (`internal/gate/phases.go:295`) and the lane (`internal/gate/lane.go:306`). The lane's checks carry no `-exec=` element, so the lane keeps its argv. `refusalMessage` serves every authorization caller. `ResumeMarker` has one caller (`internal/worktree/land_resume.go:80`).
- `Copy survival`: none.
- `Rendered-shape readers`: the old `test` argv has the three readers of the sweep table. The `gate: green` line keeps its text, and the reader line is new.
- `Pin operators`: the census check compares the `test` argv with `slices.Equal` (`internal/conformance/ordinary_build_census_test.go:39`).
- `Entry reads`: the reader plan reads the landing value from the context, and its tests inject the store root and the shellcheck inventory.
- `Derived expectations`: none. Each row's expected value comes from its fixture, not from the code under test.
- `Consolidated rules`: none.
- `Quantified obligations`: eleven covered families, one child-reader row, a bound of five, and at most 32 tables. Each count names its declaration.
- `Workflow-step writes`: none. The build adds no review-record or ledger field.

### Flagged additions and delegated decisions

- The refusal line `gate: whole gate (reader set refused: <case>)` (NG54). The map did not ask for it. It makes each whole-gate case observable to the operator and to the tests.
- The capability posture enters the identity digest only (NG2, NG4). The full-verdict oracle digest keeps its value, so current greens stay reusable.
- The table store bound of 32 tables and the marker exception (NG27).
- The shim lives in the existing internal verb `bench gate-go` as `exec`, not in a new verb.
- A red reader set attributes `candidate` (NG56). No map ticket names the kind. The chain proof is the reason, as the authorization decision states.
- The embed case matches a directive's exact path, a path below it, and a glob match (NG42, NG73, NG74). The map names "the embed targets of the composed tree", and a directory or glob directive names more than one file.
- The oracle identity keeps the `PATH` frame of the subject identity (`internal/gate/subject.go:73`), as map decision #8 reads. A landing from a session with another `PATH` therefore runs the whole gate. The NG60 measurement runs from the session shape of a real landing. If `PATH` drift defeats the reader set, the build returns to the spec.
- The ticket order puts the hold in its own chunk after the landing chunk. Until NG-C4 lands, a red whole run writes no hold, but the build window holds every landing of `main` until the build completes.

### Slice accounting

Nine tickets in six chunks. Each ticket's `Writes:` line names every file of its rows' seams. The author context of each ticket stays under 150k tokens by the byte total of its `Writes:` files.
