# Worktree joins seam reduction

Status: staged

Roadmap: FT356

Decision source: `specs/worktree-seam-reduction/decisions/worktree-seams.md` (ready compiled map, the second of its two specs).

Verification log: 2 iteration(s) to accept — round 1 folded blockers B1 (exact pin), B2 (`createAttributed` read), and B3 (multi-kind reads), with A1 to A13. Round 2 folded N1 (self-derived expectations), N2 (WS84 red), and N3 (failed-probe wording). The reviewer re-sliced 8 and 9 at sign-off.

## Problem

The worktree package injects 29 fields through its joins value. Most fields replace a
function that a real git or temp-directory fixture can fail. A test that stubs such a
field grades the stub, not the verb. The joins value also carries the Bench home and the
clock, and two fields exist only because the gate reads `BENCH_KIT` from the process.

Several functions below a verb entry read the Bench home or the clock themselves. The
resume-clean verb drops census records under the operator's home, not under the home it
receives. One resume-clean run reads the clock twice. The verb runner refuses a kit value
and a clock value, because no verb takes either value yet.

The package keeps 46 tests in its serial set. A test that asserts that a stub was not
called can pass when the test never routes the stub to the verb. Test files spell the
production table names as string literals, so a renamed table leaves the tests on a
second source.

## Solution

The joins value keeps the 15 fields that the decision source names. Each of the other 12
fields leaves the value, and each test that stubbed one converts to a real fixture under
its own name. Four fields depend on a build probe. A failed probe keeps that field and its
tests, and the build records the result for reviewer veto.

Each verb entry reads the kit root and the clock once. It passes one ambient value down,
with the Bench home and the stderr writer that it receives. The gate has one reader of the raw
`BENCH_KIT` value and kit-taking forms of its lane and kit-source functions. The verb
runner passes a kit value and a clock value to a verb.

A static census refuses every read of the Bench home, the clock, or the kit root below a
census entry. The serial ceiling drops to the live serial count. A production constant
names each table, and a test census refuses a table-name literal in a test. A verb result
reports whether the run took the joins route, and each not-called assertion reads it.

## User stories

Line: opus / high.
Implementation-line reason: SR-C2 is the hardest chunk, because the ambient value changes the signature of every internal verb form and of the retirement path. The decision source fixes the field list and the census edge, but four fixtures depend on unprobed git behavior. The package tests and the two censuses cover each row cheaply.
Harder chunks: SR-C2, SR-C4, SR-C6.

Line comparison: the reviewer directed a three-arm comparison on the SR-C5 and SR-C6 tickets. Arm one runs opus / high on the pre-slice tickets 8 and 9 at commit `9c9e9b69`. Arm two runs opus / high on the sliced tickets. Arm three runs sonnet / high with one fresh author per sliced ticket. Fable / high reviews the finish of each delegate.

The gate kit reader:

1. As a kit maintainer, I want one gate function to read the raw `BENCH_KIT` value, so that the kit root has one reader.
2. As the merge verb entry, I want a lane function that takes the kit value, so that the entry reads the kit once.
3. As the land verb entry, I want a kit-source predicate that takes the kit value, so that the install notice uses the entry's read.
4. As a commit or adopt caller, I want `LaneForCommit`, `KitSourceCheckout`, and `KitDir` to keep their signatures, so that my calls do not change.
5. As a kit maintainer, I want an empty kit value to keep each consumer's own fallback, so that each consumer answers as before.

The ambient value:

6. As a worktree verb, I want the entry to read the kit and the clock once, so that nothing below the entry reads the process.
7. As a test author, I want the verb runner to pass a kit value, so that a test selects the kit without a process bind.
8. As a test author, I want the verb runner to pass a clock value, so that a test fixes the plan day without a seam.
9. As a test author, I want the runner to refuse an unused kit or clock value, so that a test never trusts an ignored value.
10. As an operator, I want resume-clean to drop census records under the home it receives, so that a non-default home keeps no stale records.
11. As an operator, I want one resume-clean run to judge each plan and replan at one instant, so that the two plans agree.
12. As an operator, I want release, clean, and land to drop census records under the received home, so that retirement has one home.
13. As a test author, I want the two resume-clean tests to pass their home through the runner, so that they leave the serial set.

The kit seams leave:

14. As a test author, I want the merge tests to declare a manifest lane and pass a kit value, so that the real lane runs.
15. As the merge verb, I want a declared lane check that fails to refuse the merge, so that the converted fixture still grades the refusal.
16. As a worktree test author, I want the land install-notice tests to pass a kit value, so that the real kit-source predicate decides the notice.

The landing seams leave:

17. As a test author, I want the marker tests to delete the green marker from a gate script, so that the real marker swap fails.
18. As a test author, I want the reconcile tests to plant a nested destination repository, so that the real residue guard fails.
19. As a worktree test author, I want the prune test to plant a stale branch lock, so that the real branch delete fails.
20. As a test author, I want the release-refusal tests to use the public landing fixture, so that the real source authority supplies the fences.
21. As a test author, I want the stubbed-landing tests to run real landings, so that no fake commit forces a stub.
22. As a reviewer, I want a spy test to assert the real marker state instead of recorded arguments, so that the test grades the outcome.

The cleanup and reset seams leave:

23. As a test author, I want the ignored-stat test to use a directory without search permission, so that the real `os.Lstat` fails.
24. As a test author, I want retirement warnings on the verb's own stderr, so that no test replaces a warning writer.
25. As a test author, I want the special-path test to grade the real planner's reason, so that no spy replaces the planner.
26. As a worktree test author, I want the reauthorize rollback test to deny writes to the admin directory, so that the real unlock fails.
27. As a worktree test author, I want the merge-reconcile tests to plant a stale index lock, so that the real `git reset --merge` fails.
28. As a worktree test author, I want the reset-move tests to plant a stale `HEAD.lock`, so that the real move fails.

The conversion rules:

29. As a reviewer, I want each converted test to keep its name, so that the test count pin holds.
30. As a reviewer, I want a failed probe to keep its field and tests with one recorded learning, so that the build need not stop.
31. As a reviewer, I want a test that cannot convert to stop the build and name the test, so that I decide each loss.
32. As a reviewer, I want the joins value to keep the 15 named fields and each failed-probe field, so that it matches the decision.

The single-read census:

33. As a census entry, I want one read of each kind in my own body to pass, so that an entry can resolve its values.
34. As a reviewer, I want the census to refuse a read in an unexported function, so that no value is read below an entry.
35. As a reviewer, I want the census to refuse a read in a function literal or loop body, so that one read cannot run twice.
36. As a reviewer, I want the census to refuse a second read of one kind, so that one entry reads each value once.
37. As a reviewer, I want the census to refuse an unexported call to a reading entry, so that no helper hides a read.
38. As a reviewer, I want the census to refuse a read reference that is not a call, so that no function value hides a read.
39. As a reviewer, I want the census to refuse a read in a package-level declaration, so that no read runs at package initialization.
40. As a kit maintainer, I want the census to derive its read set from the source, so that a new read needs no census edit.
41. As a caller in another package, I want `ClaimRecordedLease` to read the clock in its own body, so that its unexported form takes the instant.
42. As an operator, I want `Subshell` to read the clock once and release through the internal form, so that the session reads nothing below.
43. As a reviewer, I want the census to report nothing on the live tree, so that every read sits at a census entry.
44. As a reviewer, I want the serial ceiling to equal the live serial count, so that a test leaving the serial set lowers the pin.

The one-source test items:

45. As a test author, I want a production constant to name each table that the worktree package renders, so that the tests read one source.
46. As a reviewer, I want a census to refuse a table-name literal in a worktree test, so that a renamed table leaves no stale literal.
47. As a test author, I want the verb result to report the joins route, so that a test can prove its stub was reachable.
48. As a reviewer, I want each not-called assertion to require the joins route, so that a dropped route turns the test red.

The reviewed exclusions:

49. As a reviewer, I want a census refusal for a positional fixture tuple priced out of scope, so that this spec stays one capability.
50. As a reviewer, I want the lock-hook fold into `Fault` steps priced out of scope, so that the joins value keeps the decided list.
51. As a reviewer, I want the test harness of every other package to stay out of scope, so that the spec changes one package's tests.

## Implementation decisions

### The gate kit reader

`gate.KitValue` is the one function that reads the raw `BENCH_KIT` value. `kitRoot` and
`KitDir` call it. `gate.LaneForCommitAtKit(root, kit)` and
`gate.KitSourceCheckoutAtKit(root, kit)` take the kit value. An empty kit makes the lane
form fall back to the graded root. An empty kit makes the kit-source form fall back to
the parent of the running executable, as `KitDir` does today.

That fallback does not call
`KitDir`. Neither form reaches
`KitValue`, so the census in ticket 12 can tell a kit form from a kit read.

`LaneForCommit` and `KitSourceCheckout` stay as wrappers that pass `KitValue()`, so
`internal/commit`, `internal/adopt`, and `cmd/bench` do not change. The new forms live in `kit_source.go`,
because `lane_select.go` sits near its line budget.

### The ambient value

The `ambient` value carries the Bench home, the kit value, one instant, and a warnings
writer. A verb entry builds it with one call to an `effects.go` constructor. The
constructor reads the kit through `gate.KitValue` and the clock through `currentTime`, and
it takes the home and the stderr writer from the entry's parameters. Each internal verb
form takes the joins value first and the ambient value second. The verb call census
derives a joins form from the `defaultJoins()` first argument, so that rule still holds.

The `now` and `home` fields leave the joins value. `defaultJoins` then reads no process
value. The retirement path takes the home and the warnings writer from the ambient value.
The cleanup transaction replans at the entry's instant. An older instant only makes an
assignment and a lease look younger, so the planner retains rather than removes.

`Subshell` plans its release with the instant that it reads at its entry, after the shell
session ends. The same conservative direction holds.

The verb runner builds the ambient value with the same constructor. A call's kit value and
clock value replace the constructor's reads. A call with no joins value, no kit value, and
no clock value runs the public entry. Every other call runs the internal form, with
`defaultJoins()` when the call holds no joins value. `checkVerbCall` refuses a kit value
or a clock value for a verb key without an internal form. Its messages are
`verb runner: the <key> verb takes no kit value` and
`verb runner: the <key> verb takes no clock value`.

### The converted fields

| field | real fixture | field sites | probe |
| --- | --- | --- | --- |
| `mergeLane` | a committed phase manifest lane and a kit value apart from the target | 7 | no |
| `kitSourceCheckout` | a kit value that names the destination, or another directory | 1 | no |
| `advanceLandingMarker` | a gate script in a real landing that deletes `refs/bench/green/<branch>` in the destination | 12 | yes |
| `reconcileLanding` | a nested repository in the destination | 4 | no |
| `pruneLandedBranches` | a stale `refs/heads/<name>.lock` | 1 | yes |
| `authorizeLandingSource` | the public landing fixture | 4 | no |
| `ignoredLstat` | an ignored file in a directory without search permission | 1 | yes |
| `liveBinaryWarnings` | the verb's own stderr through the ambient value | 2 | no |
| `planLandedExplicit` | the real planner and its shape reason | 1 | no |
| `reauthorizeUnlock` | an admin directory at mode 0500 | 1 | no |
| `mergeReconcile` | a stale `index.lock` in the target checkout | 2 | yes |
| `resetMove` | a stale `HEAD.lock` in the admin directory, and the ignore-rule drift for the move that does not land | 3 | no |

A field site is one assignment of the field, as the classification asset counts it. A
helper site can serve several tests: `stubbedLiveBinaryJoins` sets `liveBinaryWarnings`
once, and four live-binary tests call it. The build re-derives each count from the tree.

A probe is a `bench probe` run that omits the field's error branch in the production file
and turns the converted test red. A probe fails when it stays green, or when the fixture
cannot make the converted test pass. After a failed probe, the field stays in the joins
value with its tests unchanged. The build records one `bench learning` entry for reviewer
veto and does not stop.

A field without a probe has no such fallback. If a real fixture cannot make one of its
converted tests pass, the build stops and names the test, per decision 3.

`stubLandJoins` returns a fake commit, so each of its users moves to a real landing with a
shell gate. Its marker stub, its reconcile stub, and its source-authority stub leave with
it. The dead `planLandedExplicitWithOptions` declaration leaves with `planLandedExplicit`.

`TestResetApplyExitsThreeWhenTheMoveDidNotLand` stubs a move that exits 0 and leaves the
checkout apart from the checkpoint. The ignore-rule drift is the real fixture. The move
keeps an ignored file that the checkpoint's rules do not ignore, so the post-move check
exits 3. `TestResetApplyKeepsIgnoredBytesAcrossAnIgnoreRuleChange` uses the same drift,
and decision 3 accepts that near duplicate.

The permission fixtures follow the precedent in `classifier_shape_test.go`. Under the root
user, each such test calls `capability.Capability` with `capability.Privilege`.

### The single-read census

A census entry is a function or a method declaration whose name is exported. A read is one
reference to a read-set name. The census derives the read set from three sources:

- each function that `effects.go` declares;
- each package-qualified function that an `effects.go` body calls;
- each exported gate function whose body reaches `KitValue` through calls inside the gate package.

A kind is one package-qualified call that the census reaches from a read-set name. One
reference counts as one read of each kind that its name reaches. The ambient constructor
reaches `time.Now` and `gate.KitValue`, so the constructor and a later `currentTime` read
`time.Now` twice. A `gate.LaneForCommit` reference reads the kind `gate.KitValue`.

The census parses the non-test files of the package directory and skips `effects.go`. It
parses the gate directory for the third source only, and it follows gate calls by bare
identifier. The census counts only an entry's own reads, and it does not follow a call
from one entry to another. The package holds no such call.

The census accepts a read only in an entry's own body, outside each function literal and
loop body. It accepts only the first read of each kind there. It reports each other read
in one of three messages:

- `<file>:<line>: <declaration> reads <name> below a census entry`
- `<file>:<line>: <declaration> reads <kind> a second time`
- `<file>:<line>: <declaration> calls <Entry>, which reads <kind>`

`<name>` is the reference as the source spells it. The third message applies when an
unexported declaration calls an entry whose own body reads. One census formatter renders
each message. A census test calls that formatter with literal file, line, declaration,
and name inputs, never with the census's own findings.

`ClaimRecordedLease` reads the clock and passes the instant to `claimRecordedLease`.
`CreateCommand` reads the instant and passes it to `createAttributed`, with no net line
growth in `worktree.go`. `ApplyAutomatic` reads the clock once. `Subshell` reads the clock once and calls
`releaseCommandWith` in place of `ReleaseCommand`.

### The serial ceiling

`worktreeSerialCeiling` drops to the live serial count. The ceiling check becomes exact.
A set below the ceiling is refused with
`the package holds <n> serial tests, below the ceiling of <c>: lower worktreeSerialCeiling to <n> in this change`.
The two resume-clean tests drop their `BENCH_HOME` bind, so the count is at most 44. The
ticket moves `serialSet` and `serialCeilingBreach` out of `parallel_census_test.go`, so
that over-budget file does not grow.

A pure renderer, `belowCeilingRefusal(n, c)`, holds
the refusal text, and `serialCeilingBreach` calls it. The test asserts a non-empty breach
that equals the renderer's output for a count of 1 and a ceiling of 2.

`worktreeTestCount` is exact. Each ticket that adds a top-level worktree test raises the
pin in its own commit: tickets 2, 12, 13, 14, and 15.

### The table names

Each `toon.Table` and `toon.TableTyped` call in the package's non-test source names its
table with a package constant. The test-only `cleanupTable` and `selectedTable` constants
move to production.
A test census reports a test string literal that equals a table name and is the block
argument of `mustRows`, `readVerbRows`, `Rows`, `toon.Table`, or `toon.TableTyped`. It also
reports a test string literal that begins with a table name and `[`.

### The joins route

The verb result carries a `viaJoins` field. The verb runner sets it when the run took the
internal form with the call's joins value. A `mustViaJoins` form fails the test when the
field is false. Each test that asserts that a joins stub was not called also calls
`mustViaJoins` on that run. The name stays apart from the `joined` field of `verbForm`.

### The line budgets

The structure growth ratchet reds a file that is over its budget and grew. Each ticket
keeps each over-budget file at or below its base line count, or it moves code out. New
tests go into new test files when a sibling file is near its budget. These files are over
budget at the decision commit:

- `worktree.go`, `ownership.go`, `resume.go`, `lifecycle.go`, `subshell.go`, and `classifier.go`;
- `merge_test.go`, `worktree_test.go`, `resume_test.go`, and `parallel_census_test.go`;
- `identity_component_test.go`, `exec_test.go`, and `land_journey_test.go`;
- `lifecycle_test.go`, `ownership_test.go`, and `pool_reclaim_test.go`.

## Implementation chunks

| stable chunk ID / tickets | delivered outcome | acceptance rows | tests | harder chunk |
| --- | --- | --- | --- | --- |
| SR-C1 / `1-read-the-kit-value-once-in-gate.md` | The gate has one kit reader and kit-taking lane and kit-source forms. | WS1, WS2, WS3, WS4, WS5, WS6, WS7, WS8 | `bench test --package ./internal/gate` | no |
| SR-C2 / `2-carry-an-ambient-value-below-each-verb-entry.md` | Each verb entry passes one ambient value down, and the verb runner passes a kit value and a clock value. | WS9, WS10, WS11, WS12, WS13, WS14, WS15, WS16 | `bench test --package ./internal/worktree` | yes |
| SR-C3 / `3-pass-the-kit-value-to-merge-and-land.md` | The merge and land verbs read the kit at their entries, and the two kit fields leave the joins value. | WS17, WS18, WS19, WS20, WS21, WS22, WS23, WS24, WS25 | `bench test --package ./internal/worktree` | no |
| SR-C4 / `4-interrupt-the-landing-marker-with-a-gate-script.md`, `5-fault-the-landing-follow-on-steps-with-real-fixtures.md`, `6-land-the-stubbed-landing-tests-for-real.md` | The landing tests fault real landings, and the four landing fields leave the joins value. | WS26, WS27, WS28, WS29, WS30, WS31, WS32, WS33, WS34, WS35, WS36, WS37, WS38, WS39, WS40, WS41, WS42, WS43, WS44 | `bench test --package ./internal/worktree` | yes |
| SR-C5 / `7-fault-the-cleanup-reads-with-real-fixtures.md`, `8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md`, `9-fault-the-reset-move-with-real-fixtures.md`, `10-fault-the-merge-reconcile-with-a-stale-index-lock.md` | The cleanup, reset, merge-reconcile, and reauthorize tests use real fixtures, and their six fields leave the joins value. | WS46, WS47, WS48, WS49, WS50, WS84, WS85, WS51, WS45, WS52, WS53, WS54, WS55, WS56, WS81 | `bench test --package ./internal/worktree` | no |
| SR-C6 / `11-lift-each-read-below-a-census-entry.md`, `12-refuse-a-read-below-a-census-entry.md`, `13-make-the-serial-ceiling-exact.md` | Each read sits at a census entry, the census refuses every read below one, and the serial ceiling is exact. | WS57, WS58, WS59, WS60, WS61, WS62, WS63, WS64, WS65, WS66, WS67, WS68, WS69, WS70, WS82, WS83, WS86 | `bench test --package ./internal/worktree` | yes |
| SR-C7 / `14-name-each-table-in-production.md`, `15-require-the-joins-route-of-a-not-called-stub.md` | The tests read the table names from production, and each not-called assertion proves its route. | WS71, WS72, WS73, WS74, WS75, WS76, WS77, WS78, WS79, WS80 | `bench test --package ./internal/worktree` | no |

SR-C1 and SR-C2 each create a seam that later tickets consume, so each is its own chunk.
Each review closes before a consumer ticket starts. Every ticket writes the worktree
package, so the tickets land in series.

## Testing decisions

A good test drives a verb through the verb runner against a real fixture repository. It
asserts the exit code, the rows, and the repository state. The census tests drive the
census function over a synthetic file set, as `parallel_census_test.go` does, and one
live-tree test drives it over the package. The gate tests drive the kit forms over
temporary directories, and they create no repository and start no process.

Prior art: `TestParallelCensusOnTheLiveTree` and the synthetic census tests for the census
shape, `TestVerbCallCensusOnTheLiveTree` for a second live census, and
`TestKitSourceCheckoutMatchesThroughASymlinkSpelling` for the kit forms. The gate seam is
the ordinary `test` phase, which runs both packages.

Each converted test keeps its name. Its row names the probe that proves the fixture
reaches the production error branch. The probe is a build obligation, and the ticket
records the probe command and its red in its verification note.

### Seam diagram

    trigger: a worktree test, through the verb runner
        │
        ▼
    verb call {root, home, kit, clock, joins, args}
        │
        ▼
    [ verb entry: one read of the kit and the clock → ambient value ]
        │                                   ◀ the single-read census attaches here:
        ▼                                     it parses the source and refuses a read below
    [ internal verb form (joins, ambient) ]   an entry
        │
        ▼
    [ real git, real gate script, real planner ]  ──▶  verb result {exit, rows, viaJoins}
                      ◀ tests attach here: a real fixture plants the fault, and the
                        test asserts the exit, the rows, and the repository state

    trigger: a gate or adopt caller
        │
        ▼
    [ gate.KitValue ] ──▶ LaneForCommitAtKit / KitSourceCheckoutAtKit ──▶ lane, verdict
                      ◀ gate tests attach here: temporary directories and an explicit kit

### Acceptance coverage map

| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| WS1 | 1 | The gate non-test source holds one `os.Getenv` or `os.LookupEnv` call with the argument `"BENCH_KIT"`, inside `KitValue`. | planned TestKitValueIsTheOneBenchKitRead in internal/gate | A second inline read in `kitRoot` or `KitDir` makes the count two. |
| WS2 | 2 | `LaneForCommitAtKit(root, kit)` with a kit apart from the root returns a lane whose `Kit` equals the kit argument and whose `Selective` is false. | planned TestLaneForCommitAtKitReadsTheRootManifest in internal/gate | A form that reads the process kit returns the `Kit` that `BENCH_KIT` names, which the gate's phases set to the source root. |
| WS3 | 5 | `LaneForCommitAtKit(root, "")` returns a selective lane for a root that declares a manifest lane. | planned TestLaneForCommitAtKitFallsBackToTheRoot in internal/gate | A fallback to the executable's parent returns the root's manifest lane, which is not selective. |
| WS4 | 3 | `KitSourceCheckoutAtKit(root, kit)` is true when the kit names the root through a symbolic link. | planned TestKitSourceCheckoutAtKitMatchesAnotherSpelling in internal/gate | A form that ignores its kit argument compares the root with the executable's parent and returns false. |
| WS5 | 3 | `KitSourceCheckoutAtKit(root, kit)` is false when the kit names another directory. | planned TestKitSourceCheckoutAtKitRefusesAnotherDirectory in internal/gate | A form that answers true for each non-empty kit fails the false expectation. |
| WS6 | 5 | `KitSourceCheckoutAtKit(root, "")` is false for a temporary root apart from the executable's parent. | planned TestKitSourceCheckoutAtKitFallsBackToTheExecutableParent in internal/gate | A fallback to the root itself answers true. |
| WS7 | 4 | `KitSourceCheckout(link)` under a `BENCH_KIT` that names the kit by another spelling returns true. | `internal/gate/kit_source_test.go` (`TestKitSourceCheckoutMatchesThroughASymlinkSpelling`) | A wrapper that stops passing `KitValue()` ignores the bound kit and returns false. |
| WS8 | 4 | `KitDir()` under a bound `BENCH_KIT` returns the bound value. | planned TestKitDirReturnsTheBoundKitValue in internal/gate | A `KitDir` that stops calling `KitValue` returns the executable's parent. |
| WS9 | 6, 12 | A release under an explicit home drops the released assignment's census record under that home. | `internal/worktree/worktree_test.go` (`TestReleaseDropsTheCensusRecords`) | A retirement path that reads `Home()` drops the record under the operator's home, so the explicit-home record survives. |
| WS10 | 12 | A clean under an explicit home drops the cleaned assignment's census record under that home. | `internal/worktree/worktree_test.go` (`TestCleanDropsTheCensusRecords`) | A clean path that loses the ambient home drops the record under the operator's home. |
| WS11 | 12 | A landing under an explicit home states the census count and drops the records under that home. | `internal/worktree/land_census_test.go` (`TestLandCommandStatesTheCensusCountAndDropsTheRecords`) | A land path that loses the ambient home keeps the records under the explicit home. |
| WS12 | 10 | resume-clean through the verb runner under a non-default home retires a cleanup-pending owned assignment and drops its census record under that home. | planned TestResumeCleanDropsTheCensusRecordUnderTheCallerHome in internal/worktree | The current resume-clean drops the record under `Home()`, so the record under the caller home survives. |
| WS13 | 8, 11 | resume-clean with a clock value eight days after an active assignment's creation names the clean command for that assignment's path. | planned TestResumeCleanJudgesTheClockValue in internal/worktree | A plan that reads the real clock judges the new assignment fresh and names no clean command. |
| WS14 | 8 | `clean --discard-branch` with the clock value at `discardDay` writes the discarded ref that `intent.DiscardedRef(discardDay, ref)` names. | `internal/worktree/clean_discard_test.go` (`TestDiscardTargetWritesTheDiscardedRefFirst`) | A runner or a verb that drops the clock value dates the ref by the real day. |
| WS15 | 9 | `checkVerbCall` refuses a kit value for the `path` key with `verb runner: the path verb takes no kit value`. | `internal/worktree/verb_runner_check_test.go` (`TestVerbCallRefusesAKitValue`) | A runner that runs the public entry with a kit value accepts the call silently. |
| WS16 | 9 | `checkVerbCall` refuses a clock value for the `path` key with `verb runner: the path verb takes no clock value`. | `internal/worktree/verb_runner_check_test.go` (`TestVerbCallRefusesAClockValue`) | A runner that runs the public entry with a clock value accepts the call silently. |
| WS17 | 14 | A merge whose target's committed manifest declares a tally lane, run with a kit apart from the target, appends one byte to the tally. | `internal/worktree/merge_test.go` (`TestMergeRunsTheDeclaredLaneOnTheComposedTree`) | A merge that resolves no declared lane from the target's manifest leaves the tally empty. |
| WS18 | 15 | A declared lane check that exits 1 refuses the merge at exit 1 and leaves the branch tip unchanged. | `internal/worktree/merge_test.go` (`TestMergeRefusesAFailingLaneCheck`) | A merge that publishes after a lane fail moves the branch tip. |
| WS19 | 15 | A fast-forward whose declared lane fails refuses at exit 1. | `internal/worktree/merge_test.go` (`TestMergeRefusesAFastForwardTheLaneFails`) | A fast-forward that skips the declared lane publishes the incoming commit. |
| WS20 | 14 | A declared lane check that edits the checkout refuses the merge. | `internal/worktree/merge_test.go` (`TestMergeRefusesACheckoutEditedDuringTheLane`) | A merge that skips the declared lane never sees the edit and publishes. |
| WS21 | 14 | A declared lane that names the prose placeholder receives the incoming Markdown path. | `internal/worktree/merge_test.go` (`TestMergeResolvesTheProsePlaceholderToTheIncomingMarkdown`) | A lane resolution that drops the placeholder passes no path to the check. |
| WS22 | 14 | A merge run from the primary checkout grades the incoming prose from the target's composed tree. | `internal/worktree/merge_caller_root_test.go` (`TestMergeGradesIncomingProseFromTheComposedTreeWhateverTheCallerRoot`) | A lane anchored at the caller root reads the target checkout, where the incoming file is absent. |
| WS23 | 14 | The delegated integration journey folds each contribution through the declared lane. | `internal/worktree/delegated_integration_test.go` (`TestDelegatedIntegrationJourney`) | A journey whose lane resolution fails refuses the first fold. |
| WS24 | 7, 16 | A landing of a broker-changing diff with the kit value at the destination names the install step. | `internal/worktree/land_effects_test.go` (`TestLandCommandReportsInstallStepForABrokerChangingDiff`) | A land path that reads the process kit never matches the fixture destination, so the install step is absent. |
| WS25 | 16 | A landing with a kit value apart from the destination names the installed repair route. | `internal/worktree/land_effects_test.go` (`TestLandCommandNamesTheInstalledRepairRouteOffTheKitCheckout`) | A land path that treats every destination as the kit names the install step instead. |
| WS26 | 17 | A spec-less landing that the marker-deleting gate script interrupts resumes to completion. | `internal/worktree/land_specless_test.go` (`TestResumeLandCommandSpecLessCompletesAnInterruptedLanding`) | The probe that omits the marker error branch in `land.go` lets the landing complete without an interrupt, so the test is red. |
| WS27 | 17 | A spec-backed landing that the marker-deleting gate script interrupts resumes without `--spec`. | `internal/worktree/land_specless_test.go` (`TestResumeLandCommandWithoutSpecCompletesASpecBackedLanding`) | The same marker probe turns the test red. |
| WS28 | 17 | An interrupted landing resumes with either the spec slug or the spec path. | `internal/worktree/land_resume_test.go` (`TestResumeLandCommandAcceptsSpecSlugAndPath`) | The same marker probe turns the test red. |
| WS29 | 17 | A landing interrupted at the marker resumes and advances the marker to the published commit. | `internal/worktree/land_resume_test.go` (`TestResumeLandCommandCompletesAnInterruptedMarker`) | A resume that skips the marker leaves `refs/bench/green/<branch>` absent. |
| WS30 | 17 | A resume after a marker interrupt allows local capture files in the destination. | `internal/worktree/land_local_capture_test.go` (`TestResumeLandCommandAllowsLocalCaptureInDestination`) | The same marker probe turns the test red. |
| WS31 | 17 | A tickets-only landing that the marker-deleting gate script interrupts resumes its close. | `internal/worktree/land_tickets_only_test.go` (`TestResumeLandCommandTicketsOnlySpecCompletesAnInterruptedClose`) | The same marker probe turns the test red. |
| WS32 | 17 | A resume after a marker interrupt names each failed identity component. | `internal/worktree/identity_component_test.go` (`TestResumeLandCommandNamesEachIdentityComponent`) | A helper that no longer interrupts the landing gives no resume to grade. |
| WS33 | 17 | A resume with an unknown request after a marker interrupt names the reauthorize recovery. | `internal/worktree/land_reauthorization_test.go` (`TestResumeLandCommandUnknownRequestNamesReauthorizeRecovery`) | The same marker probe turns the test red. |
| WS34 | 18 | A published checkout that a nested destination repository leaves unreconciled resumes to a reconciled checkout. | `internal/worktree/land_resume_test.go` (`TestResumeLandCommandReconcilesAnUnreconciledPublishedCheckout`) | The probe that omits the reconcile error branch in `land.go` lets the landing reconcile, so the test is red. |
| WS35 | 17, 18 | Each post-publication failure resumes without a second publication. | `internal/worktree/land_freshness_test.go` (`TestLandCommandResumesEveryPostPublicationFailureWithoutRepublishing`) | A resume that republishes moves the destination past the first published commit. |
| WS36 | 19 | A landing whose prune meets a stale branch lock reports an incomplete prune. | `internal/worktree/land_prunes_landed_siblings_test.go` (`TestLandCommandReportsIncompletePrune`) | The probe that omits the prune error branch in `land.go` reports a complete prune. |
| WS37 | 20 | A release refusal through the public landing fixture lists each colliding path. | `internal/worktree/land_release_refusal_test.go` (`TestLandCommandRefusalListsCollidingPaths`) | A landing whose source authority refuses the fixture never reaches the release refusal. |
| WS38 | 20 | A release refusal keeps a control-bearing path in one table row. | `internal/worktree/land_release_refusal_test.go` (`TestLandCommandRefusalKeepsControlBearingPathInOneTableRow`) | A row writer that splits the escaped path prints two rows. |
| WS39 | 21 | A digest-shaped request token authenticates through a real landing. | `internal/worktree/land_reauthorization_test.go` (`TestLandCommandAuthenticatesDigestShapedRequestToken`) | A landing that refuses the digest-shaped token exits 1. |
| WS40 | 21 | An abbreviated source tip expands to the full tip in a real landing. | `internal/worktree/land_reauthorization_test.go` (`TestLandCommandExpandsAbbreviatedSourceTip`) | A landing that keeps the abbreviated tip refuses the reviewed source. |
| WS41 | 20, 21 | An abbreviated base expands to the full base in a real landing. | `internal/worktree/land_reauthorization_test.go` (`TestLandCommandExpandsAbbreviatedBase`) | A landing that keeps the abbreviated base prints the abbreviated `source_base`, not the full base. |
| WS42 | 17, 18, 21 | Each post-swap failure of a real landing prints its terminal row. | `internal/worktree/land_flags_test.go` (`TestLandCommandPostCASTerminalTable`) | A terminal table that drops a failure row prints no row for that fault. |
| WS43 | 21 | A release diagnostic in a real landing cannot forge a terminal line. | `internal/worktree/land_flags_test.go` (`TestLandCommandReleaseDiagnosticCannotForgeTerminalLines`) | A writer that copies the diagnostic unescaped prints a second terminal line. |
| WS44 | 22 | A real landing advances the green marker to the published commit after the release. | `internal/worktree/land_flags_test.go` (`TestLandCommandProjectGreenOrderTable`) | A landing that advances the marker before the publication leaves the marker at the prior tip. |
| WS45 | 32 | The joins value declares the 15 named fields and each field whose probe failed, and no other field. | review-owned | A review of `joins.go` against the decision list and the recorded probe learnings finds a surviving field. |
| WS46 | 23 | An ignored file in a directory without search permission makes the explicit plan retain with the reason `uncertain`. | `internal/worktree/worktree_test.go` (`TestIgnoredInventoryStatRaceRetains`) | The probe that omits the stat error branch in `clean.go` plans a removal. |
| WS47 | 24 | A retirement that cannot remove the handoff section prints the file and the line on the verb's stderr. | `internal/worktree/worktree_test.go` (`TestRetirementPrintsTheSectionRemovalError`) | A retirement that writes the warning to `os.Stderr` leaves the verb result's stderr empty. |
| WS48 | 24 | The residue guard warns on the verb's stderr before it removes the live binary. | `internal/worktree/live_binary_test.go` (`TestResidueGuardWarnsBeforeRemovingTheLiveBinary`) | A guard that writes to `os.Stderr` leaves the captured writer empty. |
| WS49 | 24 | The residue guard removes a foreign binary with no warning. | `internal/worktree/live_binary_test.go` (`TestResidueGuardRemovesForeignBinariesWithoutWarning`) | A guard that warns for every binary writes to the captured writer. |
| WS50 | 25 | `clean --landed` retains each special assignment path with the real planner's shape reason. | `internal/worktree/clean_landed_hostile_test.go` (`TestCleanLandedSpecialPathsRetainedWithoutOpening`) | A planner that opens the special path reports another reason or blocks. |
| WS51 | 26 | A reauthorize whose admin directory denies writes rolls back and leaves the retained state unchanged. | `internal/worktree/reauthorize_test.go` (`TestReauthorizeCommandRollsBackLockRefreshAndCASLoss`) | The probe that omits the unlock error branch in `reauthorize.go` changes the retained state. |
| WS52 | 27 | A merge whose target holds a stale `index.lock` exits 3 after the branch moves. | `internal/worktree/merge_test.go` (`TestMergeExitsThreeWhenTheReconcileFails`) | The probe that omits the reconcile error branch in `merge.go` exits 0. |
| WS53 | 27 | A reset apply reconciles the merge that a stale `index.lock` left unfinished. | `internal/worktree/reset_repair_test.go` (`TestResetApplyReconcilesAnUnfinishedMerge`) | A fixture that leaves no unfinished merge gives the reset nothing to reconcile, and the probe shows it. |
| WS54 | 28 | A reset apply whose admin directory holds a stale `HEAD.lock` exits 3 and preserves the envelope. | `internal/worktree/reset_apply_test.go` (`TestResetApplyExitsThreeOnAMoveFault`) | Without the `HEAD.lock` plant the apply exits 0, so the plant is what turns the exit to 3. |
| WS55 | 28 | A detached reset apply with a stale `HEAD.lock` exits 3 with `preserved=none`. | `internal/worktree/reset_apply_test.go` (`TestResetApplyExitsThreeWithoutAnEnvelope`) | Without the `HEAD.lock` plant the apply exits 0. |
| WS56 | 28 | A reset apply whose move keeps an ignore-rule drift file exits 3 and names the preserved ref. | `internal/worktree/reset_apply_test.go` (`TestResetApplyExitsThreeWhenTheMoveDidNotLand`) | The probe that omits the post-move check in `reset_apply.go` exits 0. |
| WS57 | 33 | A synthetic entry that reads `currentTime` once outside each loop and literal draws no report. | planned TestSingleReadCensusAcceptsOneReadInAnEntry in internal/worktree | A census that refuses each read reports the entry. |
| WS58 | 34 | A synthetic unexported function that reads `currentTime` draws the one below-entry report that the census formatter renders for it. | planned TestSingleReadCensusRefusesAReadInAnUnexportedFunction in internal/worktree | A census that checks only entries passes the helper. |
| WS59 | 35 | A synthetic entry whose function literal reads `currentTime` draws the formatter's below-entry report. | planned TestSingleReadCensusRefusesAReadInAFunctionLiteral in internal/worktree | A census that walks each literal as the entry's body accepts the read. |
| WS60 | 35 | A synthetic entry whose loop body reads `Home` draws the formatter's below-entry report. | planned TestSingleReadCensusRefusesAReadInALoopBody in internal/worktree | A census that counts syntactic reads only accepts the loop read. |
| WS61 | 36 | A synthetic entry that reads `currentTime` twice draws the formatter's second-read report for `time.Now`. | planned TestSingleReadCensusRefusesASecondRead in internal/worktree | A census that counts nothing passes the second read. |
| WS62 | 37 | A synthetic unexported function that calls a reading entry draws the formatter's helper-call report for `time.Now`. | planned TestSingleReadCensusRefusesAHelperCallToAReadingEntry in internal/worktree | A census that inspects only direct reads passes the call. |
| WS63 | 38 | A synthetic composite literal that holds the value `currentTime` in an unexported function draws the formatter's below-entry report. | planned TestSingleReadCensusRefusesAFunctionValue in internal/worktree | A census that matches only calls passes the value. |
| WS64 | 39 | A synthetic package-level variable that reads `Home()` draws the formatter's below-entry report with the variable's name. | planned TestSingleReadCensusRefusesAPackageLevelRead in internal/worktree | A census that walks only function bodies passes the variable. |
| WS65 | 40 | A function that the synthetic `effects.go` adds is a read with no census edit. | planned TestSingleReadCensusDerivesTheReadSetFromEffects in internal/worktree | A census with a fixed name list passes the new read. |
| WS66 | 40 | A synthetic unexported call to `benchhome.Dir` outside `effects.go` draws the below-entry report. | planned TestSingleReadCensusRefusesAQualifiedRead in internal/worktree | A census that reads only the `effects.go` names passes the qualified call. |
| WS67 | 40 | A synthetic unexported call to an exported gate function that calls `KitValue` draws the formatter's below-entry report. | planned TestSingleReadCensusRefusesAKitWrapperCall in internal/worktree | A census that ignores the gate directory passes the wrapper call. |
| WS68 | 6, 41, 42, 43 | The census over the live package reports nothing. | planned TestSingleReadCensusOnTheLiveTree in internal/worktree | A read left in `claimRecordedLease`, `createAttributed`, `subshellAt`, `releaseAssignment`, or `defaultJoins` draws a report. |
| WS69 | 13, 44 | The live serial set equals `worktreeSerialCeiling`, and the ceiling is below 46. | `internal/worktree/parallel_census_test.go` (`TestSerialSetStaysBelowTheCeiling`) | A ceiling left at 46 after the resume-clean tests leave the serial set draws the below-ceiling refusal. |
| WS70 | 44 | A synthetic serial set of one under a ceiling of two draws a non-empty breach equal to `belowCeilingRefusal(1, 2)`. | planned TestCensusRefusesASerialSetBelowTheCeiling in internal/worktree | A one-sided ceiling check returns an empty breach, which differs from the renderer's text. |
| WS71 | 45 | Each `toon.Table` or `toon.TableTyped` call in the package's non-test source names its table with an identifier. | planned TestTableNamesAreProductionConstants in internal/worktree | A literal `"worktree_cleanup"` left in `classifier.go` draws a report. |
| WS72 | 46 | A synthetic test literal that equals a table name in a `mustRows` block argument draws a report with its file and line. | planned TestTableNameLiteralCensusReportsABlockArgument in internal/worktree | A census that skips block arguments passes the literal. |
| WS73 | 46 | A synthetic test literal that begins with `worktree_cleanup[` draws a report. | planned TestTableNameLiteralCensusReportsARenderedHeader in internal/worktree | A census that matches only whole literals passes the header text. |
| WS74 | 46 | The table-name census over the live test files reports nothing. | planned TestTableNameLiteralCensusOnTheLiveTree in internal/worktree | A table-name literal left in a test file draws a report. |
| WS75 | 47 | A joins-form run sets the verb result's `viaJoins` field, and a public-entry run leaves it false. | planned TestVerbResultReportsTheJoinsRoute in internal/worktree | A runner that never sets the field fails the joins-form half. |
| WS76 | 48 | The reset cleanup-lock test requires the joins route on its plan run. | `internal/worktree/reset_apply_test.go` (`TestResetApplyTakesTheCleanupLock`) | The probe that drops the joins value from the plan call turns the test red. |
| WS77 | 48 | Each refresh test that asserts no build call requires the joins route. | `internal/worktree/land_effects_test.go` (`TestLandSkipsTheRefreshWithoutBuildInputs`) | The probe that drops the joins value from the landing call turns the test red. |
| WS78 | 48 | The resume refresh test requires the joins route on its interrupted landing. | `internal/worktree/land_effects_test.go` (`TestResumeReadsEffectStateFromTheTree`) | The probe that drops the joins value from the interrupted call turns the test red. |
| WS79 | 48 | Each pre-gate refusal test that asserts no landing call requires the joins route. | `internal/worktree/land_flags_test.go` (`TestLandCommandRefusesDestinationAndSourceStateBeforeGate`) | The probe that drops the joins value from the landing call turns the test red. |
| WS80 | 29 | The package declares exactly the pinned count of top-level tests. | `internal/worktree/parallel_census_test.go` (`TestPackageTestCountPin`) | A conversion that deletes or merges a test drops the count below the pin. |
| WS81 | 30, 31 | A failed probe keeps its field with one recorded learning, and a test that cannot convert stops the build with its name. | review-owned | A review of the ticket verification notes and `capture/learnings.md` finds a silent field loss or a silent skip. |
| WS82 | 40 | A synthetic unexported call to an exported gate function that reaches `KitValue` through another gate function draws the formatter's below-entry report. | planned TestSingleReadCensusRefusesAnIndirectKitRead in internal/worktree | A census that reads only direct `KitValue` calls passes a wrapper such as `KitRoot`. |
| WS83 | 11 | The census over the live package reports no clock read in a replan closure. | planned TestSingleReadCensusOnTheLiveTree in internal/worktree | The function-literal read of `currentTime` at the replan in `applyAutomaticWithTerminal` draws a report. |
| WS84 | 24 | The running-binary check fails safe with no warning writer in the joins value when the resolution is unknown. | `internal/worktree/live_binary_test.go` (`TestIsRunningBinaryFailsSafeWhenResolutionIsUnknown`) | An `isRunningBinary` that answers false for an unresolvable or unstattable running binary fails the test. |
| WS85 | 24 | The running-binary check resolves a symbolic link with no warning writer in the joins value. | `internal/worktree/live_binary_test.go` (`TestIsRunningBinaryResolvesThroughASymlink`) | A guard that compares unresolved paths misses the live binary. |
| WS86 | 36 | A synthetic entry that calls the ambient constructor and then `currentTime` draws the formatter's second-read report for `time.Now`. | planned TestSingleReadCensusCountsEachKindOfAConstructor in internal/worktree | A census that counts the constructor as one kind passes the second clock read. |

Not covered: story 49 — the positional fixture tuple census is priced under Out of scope.
Not covered: story 50 — the lock-hook fold is priced under Out of scope.
Not covered: story 51 — the decision source puts the other packages' harness out of scope.

### Edge inventory

The canonical edge classes at the census seam and the fixture seam:

- Absent versus empty kit: an unset `BENCH_KIT` and an empty kit value both select each consumer's fallback. WS3 and WS6 cover each fallback.
- A path with a space: the resume-clean fixtures use `auto clean` and `auto dirty` paths. WS12 reuses that fixture shape.
- A symbolic-link spelling of the kit: WS4 covers it.
- A special file in the parsed directory: the census reuses `parseGoFiles`, which skips each non-regular file. `TestCensusSkipsSpecialFile` covers that walk.
- The root user: each permission fixture calls `capability.Capability` with `capability.Privilege`, as `TestClassifyPathShapeUnknownUnreadableGitEntry` does.
- A missing `git`: every fixture already requires `git`. No row adds a new tool.
- A destructive worktree state: each fault fixture plants its fault inside the test's own temporary repository.

**Won't handle** lines:

- Two `effects.go` functions that call one qualified function for two values — `subshellShell` is the only `os.Getenv` caller, and WS65 reds a new name.
- A clock read in another package, such as `beginVerbSpan` timing — the rule covers one package, and each entry keeps its span.
- A read through reflection or through a method value of another package — the census parses one package's syntax, and WS68 grades every spelled read.
- A test that runs while another process holds a ref lock in the fixture repository — each fixture owns its temporary repository.

The hostile-input checklist walk: control bytes, quoted patch paths, numeric-looking
cells, whitespace predicates, TOON escapes, operator runs, flag values, and JSON-escaped
separators reach no surface that this spec changes. The table census reads Go string
literals and parses no rendered text.

## Ownership fences

- `cmd/bench/command_registry.go`
- `cmd/bench/command_registry_test.go`
- `cmd/bench/help_inventory_test.go`
- `internal/conformance/axi_query_registry_test.go`
- `internal/conformance/subcommand_routing_table_test.go`
- `internal/gate/kit_source.go`
- `internal/gate/kit_source_test.go`
- `internal/gate/kit_value_test.go`
- `internal/gate/lane_select.go`
- `internal/gate/lane_test.go`
- `internal/gate/phases.go`
- `internal/worktree/`
- `internal/worktree/build.go`
- `internal/worktree/build_test.go`
- `internal/worktree/classifier.go`
- `internal/worktree/clean.go`
- `internal/worktree/clean_landed.go`
- `internal/worktree/clean_landed_apply_test.go`
- `internal/worktree/clean_landed_hostile_test.go`
- `internal/worktree/clean_landed_test.go`
- `internal/worktree/clean_set_test.go`
- `internal/worktree/clean_unclaimed_test.go`
- `internal/worktree/delegated_integration_test.go`
- `internal/worktree/effects.go`
- `internal/worktree/eligibility_test.go`
- `internal/worktree/identity_component_test.go`
- `internal/worktree/joins.go`
- `internal/worktree/journey_race_test.go`
- `internal/worktree/land.go`
- `internal/worktree/land_effects_test.go`
- `internal/worktree/land_fixtures_test.go`
- `internal/worktree/land_flags_test.go`
- `internal/worktree/land_freshness_test.go`
- `internal/worktree/land_identity.go`
- `internal/worktree/land_local_capture_test.go`
- `internal/worktree/land_marker_fixture_test.go`
- `internal/worktree/land_prunes_landed_siblings_test.go`
- `internal/worktree/land_reauthorization_test.go`
- `internal/worktree/land_release_refusal_test.go`
- `internal/worktree/land_resume.go`
- `internal/worktree/land_resume_test.go`
- `internal/worktree/land_specless_test.go`
- `internal/worktree/land_surface_test.go`
- `internal/worktree/land_tickets_only_test.go`
- `internal/worktree/lifecycle.go`
- `internal/worktree/list.go`
- `internal/worktree/list_actions_test.go`
- `internal/worktree/list_selected.go`
- `internal/worktree/list_selected_test.go`
- `internal/worktree/live_binary.go`
- `internal/worktree/live_binary_test.go`
- `internal/worktree/merge.go`
- `internal/worktree/merge_caller_root_test.go`
- `internal/worktree/merge_from_sha_test.go`
- `internal/worktree/merge_test.go`
- `internal/worktree/orphan_render_test.go`
- `internal/worktree/ownership.go`
- `internal/worktree/parallel_census_test.go`
- `internal/worktree/path_identifier_test.go`
- `internal/worktree/pool_reclaim.go`
- `internal/worktree/pool_reclaim_test.go`
- `internal/worktree/reauthorize.go`
- `internal/worktree/reauthorize_test.go`
- `internal/worktree/reset.go`
- `internal/worktree/reset_apply.go`
- `internal/worktree/reset_apply_test.go`
- `internal/worktree/reset_plan_test.go`
- `internal/worktree/reset_refusal_test.go`
- `internal/worktree/reset_repair_test.go`
- `internal/worktree/reset_restore_refusal_test.go`
- `internal/worktree/resume.go`
- `internal/worktree/serial_ceiling_test.go`
- `internal/worktree/single_read_census_test.go`
- `internal/worktree/snapshot.go`
- `internal/worktree/snapshot_test.go`
- `internal/worktree/subshell.go`
- `internal/worktree/subshell_test.go`
- `internal/worktree/table_name_census_test.go`
- `internal/worktree/verb_fixture_test.go`
- `internal/worktree/verb_result_route_test.go`
- `internal/worktree/verb_runner_check_test.go`
- `internal/worktree/verb_runner_test.go`
- `internal/worktree/worktree.go`
- `internal/worktree/worktree_test.go`
- `reviews/worktree-seam-reduction.md`

## Ticket graph

| ticket | blocked by | chunk |
| --- | --- | --- |
| `1-read-the-kit-value-once-in-gate.md` | none | SR-C1 |
| `2-carry-an-ambient-value-below-each-verb-entry.md` | `1-read-the-kit-value-once-in-gate.md` | SR-C2 |
| `3-pass-the-kit-value-to-merge-and-land.md` | `2-carry-an-ambient-value-below-each-verb-entry.md` | SR-C3 |
| `4-interrupt-the-landing-marker-with-a-gate-script.md` | `3-pass-the-kit-value-to-merge-and-land.md` | SR-C4 |
| `5-fault-the-landing-follow-on-steps-with-real-fixtures.md` | `4-interrupt-the-landing-marker-with-a-gate-script.md` | SR-C4 |
| `6-land-the-stubbed-landing-tests-for-real.md` | `5-fault-the-landing-follow-on-steps-with-real-fixtures.md` | SR-C4 |
| `7-fault-the-cleanup-reads-with-real-fixtures.md` | `6-land-the-stubbed-landing-tests-for-real.md` | SR-C5 |
| `8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md` | `7-fault-the-cleanup-reads-with-real-fixtures.md` | SR-C5 |
| `9-fault-the-reset-move-with-real-fixtures.md` | `8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md` | SR-C5 |
| `10-fault-the-merge-reconcile-with-a-stale-index-lock.md` | `9-fault-the-reset-move-with-real-fixtures.md` | SR-C5 |
| `11-lift-each-read-below-a-census-entry.md` | `10-fault-the-merge-reconcile-with-a-stale-index-lock.md` | SR-C6 |
| `12-refuse-a-read-below-a-census-entry.md` | `11-lift-each-read-below-a-census-entry.md` | SR-C6 |
| `13-make-the-serial-ceiling-exact.md` | `12-refuse-a-read-below-a-census-entry.md` | SR-C6 |
| `14-name-each-table-in-production.md` | `13-make-the-serial-ceiling-exact.md` | SR-C7 |
| `15-require-the-joins-route-of-a-not-called-stub.md` | `14-name-each-table-in-production.md` | SR-C7 |

Each ticket writes `joins.go` or a shared test file, so the graph is a chain. Ticket 2
consumes the kit reader of ticket 1 in the ambient constructor. Ticket 11 lifts the reads
after every field removal. Ticket 12 then adds the census, so its live census grades the
final source. Ticket 13 changes only the serial ceiling, apart from the census logic.

## Out of scope

- A census refusal for a positional fixture tuple: 17 test helpers in the package return two or more values. Estimate: 20 edits, 2 gate runs.
- The fold of `cleanupLockAttempt` and `creationLockAttempt` into `Fault` steps. The hook tests pass a target or a digest, and a `Fault` step carries neither. Estimate: 6 edits, 1 gate run.
- The test harness of every other package, per the decision source. Estimate: not priced, because each package needs its own map.

## Further notes

### Source trace

| source sentence | rows |
| --- | --- |
| Ticket 2: a field stays for a fixture-proof fault or a real compile or gate run | WS45 |
| Ticket 3: a converted test keeps its name, and the count pin holds | WS80 |
| Ticket 3: a test that cannot convert stops the build and names the test | WS81 |
| Ticket 4: the census refuses a read below a verb entry | WS57 to WS68, WS82, WS83, WS86 |
| Ticket 4: the serial ceiling drops to the new serial count | WS69, WS70 |
| Ticket 10: the 15 kept fields | WS45 |
| Ticket 10: `advanceLandingMarker` goes through a gate-script fixture | WS26 to WS33, WS42, WS44 |
| Ticket 10: the 12 other GO verdicts stand | WS17 to WS25, WS34 to WS44, WS46 to WS56 |
| Ticket 10: a failed probe keeps the field, and the build does not stop | WS81 |
| Ticket 10: the spec writer may fold the lock hooks | story 50, Not covered |
| Ticket 11: each exported function is an entry and reads each value once | WS57, WS61 |
| Ticket 11: a read in an unexported function is refused | WS58 |
| Ticket 11: an unexported call to a reading exported function is refused | WS62 |
| Ticket 11: the read in `claimRecordedLease` lifts into `ClaimRecordedLease` | WS68 |
| Ticket 11: one gate reader reads `BENCH_KIT`, and `kitRoot` and `KitDir` use it | WS1, WS8 |
| Ticket 11: kit-taking lane and kit-source forms, each with its own fallback | WS2 to WS6 |
| Ticket 11: the present functions stay as wrappers | WS7, WS8 |
| Ticket 11: the three latent defects get fixed | WS12, WS68, WS83 |
| Ticket 11: one row runs resume-clean under a non-default home | WS12 |
| `roadmap/FT356.md`: a production constant names each table | WS71 to WS74 |
| `roadmap/FT356.md`: a not-called assertion proves reachability | WS75 to WS79 |
| `roadmap/FT356.md`: the positional fixture tuple census | story 49, Not covered |
| The verb runner spec: the runner accepts the kit and the clock without a signature change | WS13 to WS16, WS24 |

### Reader sweep and proof checklist

Readers of each fact that the spec changes:

- The joins fields: each production use is in `land.go`, `land_resume.go`, `land_identity.go`, `merge.go`, `reset_apply.go`, `reauthorize.go`, `clean.go`, `clean_landed.go`, `clean_set.go`, `clean_discard.go`, `live_binary.go`, `lifecycle.go`, and `worktree.go`. Each test site is in the converted-field table, and the build re-derives it.
- The `now` field: `clean_set.go` and `clean_discard.go` read it, and `discardJoins` in `clean_discard_test.go` replaces it.
- The `home` field: `executeCleanup` in `lifecycle.go` reads it, and `land.go`, `land_resume.go`, and `worktree.go` set it.
- The `checkVerbCall` messages: `verb_runner_test.go` lines 149 and 152 and `verb_runner_check_test.go` lines 367 and 375. No other file holds the text.
- `worktreeSerialCeiling`: `parallel_census_test.go` alone. `roadmap/FT356.md` names the number 46 in prose.
- The `cleanupTable` constant: `verb_fixture_test.go` declares it, and nine test files read it. The `selectedTable` constant: `list_selected_test.go` declares it, and `verb_runner_check_test.go` reads it.
- Test literals of table names: 32 lines in 14 test files hold a table name and `[`. Four block-argument literals sit in `path_identifier_test.go`, `orphan_render_test.go`, and `build_test.go`. The census re-derives the list.
- `gate.LaneForCommit`: `internal/commit/commit.go` and `joins.go`. `gate.KitSourceCheckout`: `cmd/bench/command_registry.go`, `internal/adopt/doctor.go`, `internal/adopt/link.go`, `internal/adopt/doctor_rows.go`, and `joins.go`. `gate.KitDir`: `cmd/bench/command_registry.go` and four files in `internal/adopt`. `gate.KitRoot`: `internal/coverage/citations.go`. These callers keep their calls.
- The census entries that other packages call: `Create`, `Acquire`, `PlanAutomatic`, `ClaimRecordedLease`, `ClassifyRegisteredWorktrees`, and `Pool`. Their signatures do not change.
- The verb call census reads `defaultJoins()` as the first argument of a joins form. Each internal form keeps that first argument.

Proof checklist:

- Cited symbols: each symbol in the table below resolves at `c4f9de71`.
- Import edges: `internal/worktree` already imports `internal/gate`. No new edge.
- Source-row clauses and occurrences: the source trace table above.
- Promised field labels: the verb result field `viaJoins`, the verb call fields `kit` and `clock`, and the census and refusal messages in Implementation decisions.
- Changed-function callers: `ClaimRecordedLease` has callers in other packages, and its signature stays. `ReleaseCommand` has one caller in this package, `subshellAt`, which moves to `releaseCommandWith`. `LaneForCommit` and `KitSourceCheckout` keep their callers through the wrappers.
- Copy survival: WS68 reds a surviving read below an entry. WS1 reds a surviving second `BENCH_KIT` read. WS74 reds a surviving table-name literal.
- Rendered-shape readers: no rendered output changes. The two `checkVerbCall` messages change, and the reader sweep lists their readers.

| package | cited symbols |
| --- | --- |
| `internal/worktree` production | `joins`, `defaultJoins`, `Home`, `currentTime`, `subshellShell`, `claimRecordedLease`, `ClaimRecordedLease`, `applyAutomaticWithTerminal`, `releaseAssignment`, `resumeCleanCommandWith`, `createAttributed`, `subshellAt`, `executeCleanup`, `warnSectionKept`, `warnBeforeRemovingLiveBinary`, `moveResetCheckout`, `reconcileMergeCheckout`, `refreshReauthorizeLock` |
| `internal/worktree` tests | `stubLandJoins`, `mergeFixture`, `kitCheckoutJoins`, `discardJoins`, `stubbedLiveBinaryJoins`, `delegatedJourneyJoins`, `interruptLandingAtMarker`, `checkVerbCall`, `verbCall`, `verbResult`, `verbForms`, `serialSet`, `serialCeilingBreach`, `worktreeSerialCeiling`, `worktreeTestCount`, `parseGoFiles`, `cleanupTable` |
| `internal/gate` | `KitDir`, `KitSourceCheckout`, `kitRoot`, `KitRoot`, `LaneForCommit`, `LaneFor`, `manifestPath` |
| `internal/gate/authorization` | `AdvanceMarker`, `AuthorizeWithWriters` |

Sources re-read in the authoring session:

- `roadmap/FT356.md`, the decision map, its 11 tickets, and the classification asset;
- `joins.go`, `effects.go`, `parallel_census_test.go`, and `effect_census_test.go`;
- `verb_runner_test.go`, `verb_runner_check_test.go`, `verb_fixture_test.go`, and `verb_call_census_test.go`;
- `gate/kit_source.go`, `gate/phases.go` lines 140 to 230, and `gate/lane_select.go` lines 104 to 130;
- `gate/manifest.go` lines 26 to 64 and 171 to 188, and `gate/lane.go` lines 256 to 290;
- `gate/authorization/authorization.go` lines 67 to 194.

Drift since the classification commit `8ab52861`: the joins field set and every cited
default are unchanged. The test-site line numbers moved when the verb runner spec landed.
The count pin is now 699, not 664. The asset's `KitDir` caller list is stale: the adopt
callers are now `setup.go`, `init.go`, `link.go`, and `upgrade.go`.

Not exercised in this session, and so marked `uncertain` until a ticket records its probe red:

- the four probe-dependent fixtures;
- the marker-deleting gate script and the nested-repository reconcile fixture;
- the stale `HEAD.lock` move fixture;
- the resume-clean stale-active summary in WS13.

### Fence disposition

Reviewer disposition of the ownership fences: open. The fence equals the union of the
ticket `Writes:` lines. Ticket 2 holds the whole worktree package, because the ambient
value changes every internal verb form. The other worktree tickets name their files. The
fence holds five gate files for ticket 1.

The binding registry binds the worktree package to five files. Three are in `cmd/bench`,
and two are conformance registry tests. Build preflight requires each worktree ticket to name them. No ticket plans an edit
there, because no verb grammar changes. Each chunk review confirms that the five files
stay unchanged. The spec authoring commit adds the census entry
and ambient value terms to the glossary, so no ticket writes `CONTEXT.md`.

### Completion plan

```bench-completion-plan
{"version":2,"execution":{"mode":"delegate","run_id":"worktree-seam-reduction-20261002","orchestrator_session":"claude:session-018nyJAsDqW5oX9xoqL3vvFk","author_limit":1,"assignments":{"1-read-the-kit-value-once-in-gate.md":[{"session":"claude:bench-writer/sr-t1-author","assignment":"sr-t1-author","model":"opus","effort":"high","source":"6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69","native_ref":"claude:agent/sr-t1-author-20261002@6ea6b7e6fe86e3fee0d0fc9a1ff01b6808486a69"},{"session":"claude:bench-writer/sr-t1-repair1","assignment":"sr-t1-repair1","model":"opus","effort":"high","source":"d3637be279ad27405c6621bc3d22e356df777d1e","native_ref":"claude:agent/sr-t1-repair1-20261002@d3637be279ad27405c6621bc3d22e356df777d1e","predecessor":"claude:bench-writer/sr-t1-author","trigger":"user-directed","stopped":"the author session returned its final report at 0d7b4066 and its record commit 01bbe3da","preserved":"0d7b40665f232bb348aa81d52e4c1360a9650c83"}],"2-carry-an-ambient-value-below-each-verb-entry.md":[{"session":"claude:bench-writer/sr-t2-author","assignment":"sr-t2-author","model":"opus","effort":"high","source":"3ecdaffff22967433363e9495c8807f14e2f85ea","native_ref":"claude:agent/sr-t2-author-20261002@3ecdaffff22967433363e9495c8807f14e2f85ea"},{"session":"claude:bench-writer/sr-t2-repair1","assignment":"sr-t2-repair1","model":"opus","effort":"high","source":"d68dd224133551cdf60152916b9f86842c5ab911","native_ref":"claude:agent/sr-t2-repair1-20261002@d68dd224133551cdf60152916b9f86842c5ab911","predecessor":"claude:bench-writer/sr-t2-author","trigger":"user-directed","stopped":"the author session returned its final report at f620ccd6 and its verification entry at d68dd224","preserved":"f620ccd6e416501f2d4cf504d5090f3cb5104c80"}],"3-pass-the-kit-value-to-merge-and-land.md":[{"session":"claude:bench-writer/sr-t3-author","assignment":"sr-t3-author","model":"opus","effort":"high","source":"1ab5800b8fb669b26a247d8f585bf21e946863a9","native_ref":"claude:agent/sr-t3-author-20261002@1ab5800b8fb669b26a247d8f585bf21e946863a9"},{"session":"claude:bench-writer/sr-t3-repair1","assignment":"sr-t3-repair1","model":"opus","effort":"high","source":"021180277772df06c294bf7e9dba19af786fb963","native_ref":"claude:agent/sr-t3-repair1-20261002@021180277772df06c294bf7e9dba19af786fb963","predecessor":"claude:bench-writer/sr-t3-author","trigger":"user-directed","stopped":"the author session returned its final report at c87a0c3f and its record commit 02118027","preserved":"c87a0c3f7953210e3d22c356e948e9e15c302be0"}],"4-interrupt-the-landing-marker-with-a-gate-script.md":[{"session":"claude:bench-writer/sr-t4-author","assignment":"sr-t4-author","model":"opus","effort":"high","source":"2af88a0b14516dee1019fba94abadc7c85dbd083","native_ref":"claude:agent/sr-t4-author-20261002@2af88a0b14516dee1019fba94abadc7c85dbd083"}],"5-fault-the-landing-follow-on-steps-with-real-fixtures.md":[{"session":"claude:bench-writer/sr-t5-author","assignment":"sr-t5-author","model":"opus","effort":"high","source":"8ffc347ad583e2121f9ad6f3fb88a3c4754b66a1","native_ref":"claude:agent/sr-t5-author-20261002@8ffc347ad583e2121f9ad6f3fb88a3c4754b66a1"},{"session":"claude:bench-writer/sr-t5-repair1","assignment":"sr-t5-repair1","model":"opus","effort":"high","source":"a14627c15a3abdb527eaee9cae4e79218ee0671e","native_ref":"claude:agent/sr-t5-repair1-20261002@a14627c15a3abdb527eaee9cae4e79218ee0671e","predecessor":"claude:bench-writer/sr-t5-author","trigger":"user-directed","stopped":"the author session returned its final report at 2c83701a and its verification entry at a14627c1","preserved":"a2ec1cc86f72fd1161557d74154a39c9ce72e654"}],"6-land-the-stubbed-landing-tests-for-real.md":[{"session":"claude:bench-writer/sr-t6-author","assignment":"sr-t6-author","model":"opus","effort":"high","source":"2c83701a25ec00d47e6d19dafbb58fe44f2d9fa8","native_ref":"claude:agent/sr-t6-author-20261002@2c83701a25ec00d47e6d19dafbb58fe44f2d9fa8"},{"session":"claude:bench-writer/sr-t6-repair1","assignment":"sr-t6-repair1","model":"opus","effort":"high","source":"a14627c15a3abdb527eaee9cae4e79218ee0671e","native_ref":"claude:agent/sr-t6-repair1-20261002@a14627c15a3abdb527eaee9cae4e79218ee0671e","predecessor":"claude:bench-writer/sr-t6-author","trigger":"user-directed","stopped":"the author session returned its final report at a2ec1cc8 and its record commit a14627c1","preserved":"a2ec1cc86f72fd1161557d74154a39c9ce72e654"}],"7-fault-the-cleanup-reads-with-real-fixtures.md":[{"session":"claude:bench-writer/sr-t7-author","assignment":"sr-t7-author","model":"opus","effort":"high","source":"fb8da4b28980ed5262c6fa9e1eb7ace5ead91d71","native_ref":"claude:agent/sr-t7-author-20261002@fb8da4b28980ed5262c6fa9e1eb7ace5ead91d71"}],"8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md":[{"session":"claude:bench-writer/sr-b1-author/t8","assignment":"sr-b1-author/t8","model":"opus","effort":"high","source":"eeeaa7594e1fdbe79e378d85975bddd34bdfb93c","native_ref":"claude:agent/sr-b1-author-20261002@eeeaa7594e1fdbe79e378d85975bddd34bdfb93c"}],"9-fault-the-reset-move-with-real-fixtures.md":[{"session":"claude:bench-writer/sr-b1-author/t9","assignment":"sr-b1-author/t9","model":"opus","effort":"high","source":"eeeaa7594e1fdbe79e378d85975bddd34bdfb93c","native_ref":"claude:agent/sr-b1-author-20261002@eeeaa7594e1fdbe79e378d85975bddd34bdfb93c"}],"10-fault-the-merge-reconcile-with-a-stale-index-lock.md":[{"session":"claude:bench-writer/sr-b1-author/t10","assignment":"sr-b1-author/t10","model":"opus","effort":"high","source":"eeeaa7594e1fdbe79e378d85975bddd34bdfb93c","native_ref":"claude:agent/sr-b1-author-20261002@eeeaa7594e1fdbe79e378d85975bddd34bdfb93c"}],"11-lift-each-read-below-a-census-entry.md":[{"session":"claude:bench-writer/sr-b2-author/t11","assignment":"sr-b2-author/t11","model":"opus","effort":"high","source":"564f4d46cc32fb52362fd6fae155650eff6a2386","native_ref":"claude:agent/sr-b2-author-20261002@564f4d46cc32fb52362fd6fae155650eff6a2386"}],"12-refuse-a-read-below-a-census-entry.md":[{"session":"claude:bench-writer/sr-b2-author/t12","assignment":"sr-b2-author/t12","model":"opus","effort":"high","source":"564f4d46cc32fb52362fd6fae155650eff6a2386","native_ref":"claude:agent/sr-b2-author-20261002@564f4d46cc32fb52362fd6fae155650eff6a2386"}],"13-make-the-serial-ceiling-exact.md":[{"session":"claude:bench-writer/sr-b2-author/t13","assignment":"sr-b2-author/t13","model":"opus","effort":"high","source":"564f4d46cc32fb52362fd6fae155650eff6a2386","native_ref":"claude:agent/sr-b2-author-20261002@564f4d46cc32fb52362fd6fae155650eff6a2386"}],"14-name-each-table-in-production.md":[],"15-require-the-joins-route-of-a-not-called-stub.md":[]}},"chunks":[{"id":"SR-C1","tickets":["1-read-the-kit-value-once-in-gate.md"],"verification":[{"id":"1-gate","command":"bench test --package ./internal/gate","ticket":"1-read-the-kit-value-once-in-gate.md"}]},{"id":"SR-C2","tickets":["2-carry-an-ambient-value-below-each-verb-entry.md"],"verification":[{"id":"2-worktree","command":"bench test --package ./internal/worktree","ticket":"2-carry-an-ambient-value-below-each-verb-entry.md"}]},{"id":"SR-C3","tickets":["3-pass-the-kit-value-to-merge-and-land.md"],"verification":[{"id":"3-worktree","command":"bench test --package ./internal/worktree","ticket":"3-pass-the-kit-value-to-merge-and-land.md"}]},{"id":"SR-C4","tickets":["4-interrupt-the-landing-marker-with-a-gate-script.md","5-fault-the-landing-follow-on-steps-with-real-fixtures.md","6-land-the-stubbed-landing-tests-for-real.md"],"verification":[{"id":"4-worktree","command":"bench test --package ./internal/worktree","ticket":"4-interrupt-the-landing-marker-with-a-gate-script.md"},{"id":"5-worktree","command":"bench test --package ./internal/worktree","ticket":"5-fault-the-landing-follow-on-steps-with-real-fixtures.md"},{"id":"6-worktree","command":"bench test --package ./internal/worktree","ticket":"6-land-the-stubbed-landing-tests-for-real.md"}]},{"id":"SR-C5","tickets":["7-fault-the-cleanup-reads-with-real-fixtures.md","8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md","9-fault-the-reset-move-with-real-fixtures.md","10-fault-the-merge-reconcile-with-a-stale-index-lock.md"],"verification":[{"id":"7-worktree","command":"bench test --package ./internal/worktree","ticket":"7-fault-the-cleanup-reads-with-real-fixtures.md"},{"id":"8-worktree","command":"bench test --package ./internal/worktree","ticket":"8-fault-the-reauthorize-unlock-with-a-denied-admin-directory.md"},{"id":"9-worktree","command":"bench test --package ./internal/worktree","ticket":"9-fault-the-reset-move-with-real-fixtures.md"},{"id":"10-worktree","command":"bench test --package ./internal/worktree","ticket":"10-fault-the-merge-reconcile-with-a-stale-index-lock.md"}]},{"id":"SR-C6","tickets":["11-lift-each-read-below-a-census-entry.md","12-refuse-a-read-below-a-census-entry.md","13-make-the-serial-ceiling-exact.md"],"verification":[{"id":"11-worktree","command":"bench test --package ./internal/worktree","ticket":"11-lift-each-read-below-a-census-entry.md"},{"id":"12-worktree","command":"bench test --package ./internal/worktree","ticket":"12-refuse-a-read-below-a-census-entry.md"},{"id":"13-worktree","command":"bench test --package ./internal/worktree","ticket":"13-make-the-serial-ceiling-exact.md"}]},{"id":"SR-C7","tickets":["14-name-each-table-in-production.md","15-require-the-joins-route-of-a-not-called-stub.md"],"verification":[{"id":"14-worktree","command":"bench test --package ./internal/worktree","ticket":"14-name-each-table-in-production.md"},{"id":"15-worktree","command":"bench test --package ./internal/worktree","ticket":"15-require-the-joins-route-of-a-not-called-stub.md"}]}],"final_verification":[{"id":"coverage","command":"bench coverage --check specs/worktree-seam-reduction/spec.md"},{"id":"gate","command":"bench test --package ./internal/gate"},{"id":"worktree","command":"bench test --package ./internal/worktree"},{"id":"commit","command":"bench test --package ./internal/commit"},{"id":"adopt","command":"bench test --package ./internal/adopt"}]}
```

### Flagged additions

Each addition below is not in a decision answer. The reviewer can veto each one.

- The names `KitValue`, `LaneForCommitAtKit`, and `KitSourceCheckoutAtKit`, under the map's naming discretion.
- The `ambient` value with a warnings writer, so that `liveBinaryWarnings` leaves through the same value.
- The verb runner rule: a call with any kit, clock, or joins value runs the internal form.
- The `checkVerbCall` refusal for a verb without an internal form, which keeps the names of the two refusal tests.
- The census refusals for a read in a function literal, a loop body, a function value, and a package-level declaration.
- The census read set from each exported gate function that reaches `KitValue`, which catches a hidden kit read through a wrapper.
- The exact census messages and the below-ceiling refusal message.
- The exact serial ceiling check, which refuses a set below the ceiling.
- The `viaJoins` verb result field and the `mustViaJoins` form, as the reachability proof for a not-called assertion.
- The test census for table-name literals, as the enforcement of the one-source table names.
- The removal of the dead `planLandedExplicitWithOptions` declaration.
- The ignore-rule drift fixture for `TestResetApplyExitsThreeWhenTheMoveDidNotLand`.
- The census entry and ambient value glossary terms in `CONTEXT.md`.
- The rule that the `AtKit` forms reach no `KitValue` read, and that the empty-kit fallback does not call `KitDir`.
- The pure `belowCeilingRefusal` renderer, so that the ceiling test expectation does not come from the function under test.
- The spec-stage edit of the `ROADMAP.md` recommended sequence, which repoints its first step to this spec.
