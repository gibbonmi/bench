# Worktree verb runner review record

## VR-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `vr-t1-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `7872c5e4` and committed `137fb303` on a lane pass in the first attempt. The author then committed this record in a second commit.

The chunk pair is `0c95c944..137fb303`. The base is the `main` tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The coordinator moves the tip to the record commit. The source digest does not change, because the record file is outside the graded source.

The author wrote each test before the code that it grades, in three stages. The red and green log for each row follows:

- Stage 1: every runner function and reader was a stub, and the runner returned exit -1. All 22 planned tests failed, and VR1 failed in each of its 14 key subtests, which gives 36 failure rows. The log records 35 failures. The stub source is not in the history, so the count is not reproducible.
- Stage 2: the runner and `checkVerbCall` were complete, and the readers and must forms were stubs. VR1 to VR5, VR17 to VR19, VR21, and VR22 passed. VR6 to VR14, VR16, and VR59 failed for the reader reason, with 11 failures. The log does not record the stage 2 result of VR15, and that result is not reproducible.
- Stage 3: after the readers and the must forms, all 22 tests passed.
- VR15: in stage 1, the test failed in its fixture only, because the stub runner gave no plan. A probe of the finished must form gave the row its own red. The probe made every reader error fail the recorder, and the test failed on the `none` plan.
- VR20 is review-owned. The call value declares the `kit` and `clock` fields beside `root` and `home`.

The author reported these deviations from the ticket. The Spec axis grades each one.

- The `show` and `build` grammars refuse an extra argument with the command name only. So the VR1 expectation for those two keys reads the command name from the first three words of the usage constant.
- The VR59 test builds the faulted unclaimed set with the four fixture lines of `TestCleanUnclaimedErrorRowRefusesTheSet`, as the ticket tells. That file is outside the fence, so the author did not extract a shared builder. This is a second copy of one fixture.
- The runner tests do not call `newOwnedAssignment`. That builder returns a positional tuple, ticket 3 changes it, and the ticket 3 `Writes:` line does not hold the runner files. A one-value builder, `runnerRepo`, returns a call value, and each test creates its assignment with `mustCreate`. A first draft had a local tuple builder, and the duplicated-facts sweep removed it, because VR40 refuses a tuple in any test file.
- VR17 stubs the `buildSubject` field of the joins value, so the default build is never reached.
- VR5 compares a direct `ExecCommand` call with the runner, because a child that writes both streams and exits 3 grades all three values at once.
- The runner also fails the test for an unknown key. The fingerprint reader returns an error for a cell that does not decode as text. No test grades that branch.
- No check required an edit to the five registry paths on the `Writes:` line, so the diff leaves them unchanged.

The author computed the plan digest from the `ReadPlan` rule. The input is the JSON array of the spec bytes and the eleven ticket bytes, in plan order. The same rule gives the digest `7355203a` for the tree targets record at `a08359f9`, which confirms the method. The source digest is the tree of `137fb303`, which does not hold this record file.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The second row is the plan probe `1-probe`, and the third row is the author's own probe of the joins dispatch. The JSON payload holds the exact command and output of those two rows. The first row ran on the source before the ticket commit, so only this table holds it.

| File | Mutation | Test | Row | Verdict |
|---|---|---|---|---|
| `internal/worktree/verb_runner_test.go` | swap: `if errors.Is(err, errNoVerbFingerprint) {` to `if err == nil {` | TestVerbResultMustNoFingerprintAcceptsAnErrorPlan | VR15 | bit |
| `internal/worktree/verb_runner_test.go` | swap: `unapplicableFingerprint` to `"probe-placeholder"` | TestVerbResultFingerprintTreatsAPlaceholderAsAbsent | VR59 | bit |
| `internal/worktree/verb_runner_test.go` | swap: `if call.joins != nil {` to `if false {` | TestVerbRunnerPassesTheJoinsValue | VR17 | bit |

### Verification

The author ran each VR-C1 verification on the source of `137fb303`, and each passed. The JSON payload holds each result. The package excerpt omits its two skip rows. Each skip is an environment capability skip for unix sockets.

`TestPackageTestCountPin` passed at 686 tests, and `TestSerialSetStaysBelowTheCeiling` passed at the ceiling of 46. The author also ran these checks, and each passed: `bench test --package ./internal/conformance` and `bench structure --growth 0c95c944`. The commit chain reported 15 green checks, 1 check that does not apply, and 0 red checks.

### Coordinator verification

The coordinator found a clean tree, and the ticket commit touches only paths on the `Writes:` line. The coordinator ran an independent omission probe at a site that no author probe used. The probe removed the two-fingerprints refusal in `readVerbFingerprint`, and `TestVerbResultFingerprintRefusesConflictingCells` failed. The restore reads `yes`.

## VR-C1 chunk review, round 1

The frozen pair is base `0c95c9447c20189f3f2155719ef965bffc339856` and tip `4479b2fa62bcb334f4b65aedfcd155cde31d66a6`. The coordinator moved the chunk tip from the ticket commit to the record commit, and the source digest stays the same. The shared evidence is `sha256:e3159ba832b81092e567ee0470a52d17f5da030e7694a5fb0fa17e4318cc1ef4`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, on the conditional review line. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 9: Standards 4, Spec 3, and Coverage 2. The Standards literal finding and the Spec literal finding name the same fix. Two findings are `no-op`. So 6 repair targets remain, and they take repair cycle 1 of 2.

The chunk record for VR44 to VR46 follows. The PASS set at the tip is the base set plus the 22 added tests, and all 686 top-level tests pass. The SKIP set holds the two socket capability subtests, and the base has the same two. No existing test changed except the `worktreeTestCount` value, so no assertion fell.

## Standards

Findings: 4. The worst issue is a second copy of the faulted unclaimed fixture.

- R1: `internal/worktree/verb_runner_check_test.go:268-278` copies the five fixture lines and the reason comment of `clean_unclaimed_test.go:312-321`. `AGENTS.md` names a fixture harness pasted N times as duplicated knowledge. Extract one builder beside `addUnclaimedBranch` in `clean_set_apply_test.go`, and let both tests call it. Both files are in the spec fence, so the repair takes a plan commit for the ticket 1 `Writes:` line. `auto-fix`. Confidence 8.
- R2: `internal/worktree/verb_runner_check_test.go:83` restates the reason text `"invalid invocation; run "` from the producer at `worktree.go:247`. The reader sweep says that no literal moves into a test. The `clean` subtest can render `cleanInvocationError` into a buffer and compare the whole output. `auto-fix`. Confidence 7.
- R3: `internal/worktree/verb_runner_check_test.go:39`, the `usageCommand` helper, derives the refusal command name a second time from the first three words of the usage constant. The verbs take that name from `worktreeShowGrammar.Cmd` and `buildGrammar.Cmd`. Read those two fields instead. The production duplication between the `Cmd` literals and the usage constants is outside the fence and goes to the ideas inbox. `auto-fix`. Confidence 5.
- R4: `internal/worktree/verb_runner_check_test.go:206` asserts that the fingerprint is the last record cell, which is a layout fact of `reset.go:67` that the reader contract does not need. Assert that the record holds the cell instead. `auto-fix`. Confidence 4.

## Spec

Findings: 3. The worst issue is a why-clause mutant that the VR9 test does not kill alone.

- `internal/worktree/verb_runner_check_test.go:182`: a reader that takes the first 64-hex run passes VR9, because a real plan has no 64-hex text before the cell. VR14 kills that mutant, so the suite still catches it. `no-op`. Confidence 7.
- R2: the same literal as the Standards finding R2. `auto-fix`. Confidence 6.
- R5: `reviews/worktree-verb-runner.md:11-12`: the stage 2 log does not account for VR15. The stage 1 count reads 35, but 22 tests and 14 subtests give 36. Correct the log from the run evidence. This correction is evidence-only. `auto-fix`. Confidence 5.

## Coverage

Findings: 2. The worst issue is the untested record branch of the fingerprint reader.

- R6: `internal/worktree/verb_runner_test.go:192-203`: no test grades the absent rule or the conflict rule on a record line. A no-op `reset` plan writes `fingerprint=none` in its record. A swap that returns the first record cell was silent, and the restore reads `yes`. Add a real no-op `reset` record test that expects the no-fingerprint error, add a conflicting-record case, and raise `worktreeTestCount` for each added test. `auto-fix`. Confidence 8.
- `internal/worktree/verb_runner_test.go:95-108`: only VR17 grades a `joined:` entry. A swap of the `land` joins form was silent. Tickets 2 to 10 run stubbed joins values through each of these keys, so a wrong form fails those migrated tests. `no-op`. Confidence 6.

### Advice

- The `"usage: "` prefix in the VR1 expectations repeats each grammar's help composition. A test can read the grammar's help field where a grammar variable exists.
- The length check of 64 in the VR9 test restates the digest width.
- A shared record encoder in `toon` would give the producer and the reader one source.

### Command contribution

The implementation command did not contribute. The ticket told the author to build the fault as an outside test does, but its `Writes:` line did not hold the owning file. Ticket slicing owns that fix: the first ticket that needs a shared builder names it and writes its owning file.

## VR-C1 repair 1

A fresh `bench-writer` repair session, `vr-t1-repair-1`, ran on opus at medium effort, with a cap of 2 attempts. The session started at `c71d932b` and committed `3a318142` on a lane pass in the first attempt. The lane passed with 15 green checks, 1 check that does not apply, and 0 red checks. The source digest of the repair is `0dee07c82fe40e1420c5629289bd195d4e8ff5cd`, the tree of `3a318142` without this record file. This record is a second commit.

The repair changes only test files. The red and green route for each target follows:

- R1: a new builder, `addBrokenUnclaimedBranch`, owns the broken unclaimed ref and its reason comment. `TestCleanUnclaimedErrorRowRefusesTheSet` and `TestVerbResultFingerprintTreatsAPlaceholderAsAbsent` both call it, and the assertions of the clean test do not change. The builder is in `verb_runner_check_test.go`, not beside `addUnclaimedBranch`. In `clean_set_apply_test.go`, the builder made the file 417 lines, and the structure lane refused it at its budget of 400. Red: a swap in the builder that makes the ref name a commit failed both tests.
- R2: the `clean` row of VR1 renders `cleanInvocationError` into a buffer, and the expectation is that whole output. Red: a swap that sends the `clean` key to the reclaim entry failed the `clean` subtest.
- R3: the `show` and `build` rows read `worktreeShowGrammar.Cmd` and `buildGrammar.Cmd`. The `usageCommand` helper is deleted. Red: a swap that sends the `show` key to the build entry failed the `show` subtest.
- R4: the VR10 test asserts that the record holds the `fingerprint=` cell with the value, and does not assert its position. Red: a reader that cuts the cell at the first `0` failed the test.
- R5: the stage 1 and stage 2 lines of the author evidence now state what this record cannot reproduce. The stub source is not in the history.
- R6: two new tests read the record branch. `TestVerbResultFingerprintReadsANoOpRecordAsAbsent` reads a real no-op `reset` plan and expects the no-fingerprint error. `TestVerbResultFingerprintRefusesConflictingRecords` adds a second copy of a real `reset` record with a different fingerprint and expects a conflict error. A shared builder, `resetRunnerPlan`, gives both tests and VR10 a real plan. `worktreeTestCount` is now 688.

### Repair probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`. The first two rows ran on the committed source, and the JSON payload holds the plan probe.

| File | Mutation | Test | Target | Verdict |
|---|---|---|---|---|
| `internal/worktree/verb_runner_test.go` | swap: `values = recordFingerprints(stdout)` to a return of the first record cell | both R6 tests | R6 | bit, 2 failed |
| `internal/worktree/verb_runner_test.go` | swap: `unapplicableFingerprint` to `"probe-placeholder"` | TestVerbResultFingerprintTreatsAPlaceholderAsAbsent | VR59 | bit |
| `internal/worktree/verb_runner_test.go` | swap: `if value != values[0] {` to `if false {` | TestVerbResultFingerprintRefusesConflictingRecords | R6 | bit |
| `internal/worktree/verb_runner_check_test.go` | swap: `"hash-object", "-w", "tracked.txt"` to `"rev-parse", "HEAD"` | both R1 tests | R1 | bit, 2 failed |
| `internal/worktree/verb_runner_test.go` | swap: the `clean` entry to `ReclaimCommand` | TestVerbRunnerKeyReachesItsOwnVerb/clean | R2 | bit |
| `internal/worktree/verb_runner_test.go` | swap: the `show` entry to `BuildCommand` | TestVerbRunnerKeyReachesItsOwnVerb/show | R3 | bit |
| `internal/worktree/verb_runner_test.go` | swap: the record cell cut at `,` to a cut at `0` | TestVerbResultFingerprintReadsTheRecordCell | R4 | bit |

### Repair verification

The session ran each check on the source of `3a318142`, and each passed: `bench test --package ./internal/worktree`, the plan probe `1-probe`, `bench test --package ./internal/conformance`, and `bench structure --growth 0c95c944`. `TestPackageTestCountPin` passed at 688 tests, and `TestSerialSetStaysBelowTheCeiling` passed at the ceiling of 46. The two new tests call `t.Parallel()` and bind no environment. The digest of the `1-worktree` result is of the text that the session received, because that output did not spill to a file.

## VR-C1 chunk review, round 2

This round confirms repair 1. The frozen pair is base `0c95c9447c20189f3f2155719ef965bffc339856` and tip `d44d28ce6ff1964ffc4cb77d969cecdb881356eb`. The shared evidence is `sha256:565cf66f6d2860a4079406d195130860fd2993adebbac2089ebf7de385f56926`. Each axis ran in a fresh `bench-reviewer` session on opus at high effort, and each read the repair delta `179e2eb8..d44d28ce`. Only the Coverage axis ran tests and probes, and it left the tree clean.

The raw finding count is 2: Standards 0, Spec 2, and Coverage 0. Every fold of R1 to R6 holds. Both Spec findings name spec text, not code, so they use no repair cycle. The consumed allowance stays at 1 of 2 repair cycles.

The coordinator made three evidence-only corrections to this record. The repair verification entries named the full tree `bef4351236549fad5b2cef6e5958d23b685aabb8` as their source digest. The source digest excludes this record file, so the entries now name `0dee07c82fe40e1420c5629289bd195d4e8ff5cd`. The plan digest now reads the plan at `f3b6ad76`, which adds VR60 and VR61. Each author and repair result digest hashed the whole command output, so each digest now hashes its embedded excerpt, as the record parser requires.

The chunk record for VR44 to VR46 follows. All 688 top-level tests pass, which is the base 664 plus 24 added tests. The SKIP set holds the two socket capability subtests. No existing test lost a failure call. The only edits to existing files are setup lines in `clean_unclaimed_test.go` and the count constant.

### Standards, round 2

Findings: 0. All five folds hold. The shared fault builder sits in the runner check file, because `clean_set_apply_test.go` is at its 400-line budget. The axis rates that placement a judgment call and not a binding defect.

### Spec, round 2

Findings: 2. The worst issue is a coverage map that did not name the two record-branch tests.

- R7: `specs/worktree-verb-runner/spec.md` had no row for the record-branch absent rule and conflict rule. Plan commit `f3b6ad76` adds VR60 and VR61, and ticket 1 covers both. `auto-fix`. Confidence 5.
- R8: the VR1 row said that every expectation reads its usage constant, but the `show` and `build` rows read the grammar `Cmd` field. Plan commit `f3b6ad76` states that source in VR1, in the ticket, and in the reader sweep. This change is non-behavioral, and the reviewer can veto it. `auto-fix`. Confidence 6.

### Coverage, round 2

Findings: 0. Each fold bit under a new probe site. The probes removed the placeholder rule on the record branch and broke the shared fault builder. They also moved the `show` and `clean` keys to other verbs and cut the record cell at the wrong byte. Each restore reads `yes`.

### Advice, round 2

- A new file can hold the whole unclaimed branch fixture family and free `clean_set_apply_test.go` from its budget.
- The VR1 `clean` row does not refuse an empty expectation. A check that the rendered refusal is not empty closes that gap.
- In the conflict test, a second `value` shadows the first one.

```bench-review-record
{
  "version": 2,
  "spec": "specs/worktree-verb-runner/spec.md",
  "plan_digest": "sha256:53cb1a1fe5e8e1c000b919dcde0b788f21556ae3e122fd755105b89e68d88e0e",
  "implementation_session": "",
  "chunks": [
    {
      "id": "VR-C1",
      "base": "0c95c9447c20189f3f2155719ef965bffc339856",
      "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
      "plan_digest": "sha256:53cb1a1fe5e8e1c000b919dcde0b788f21556ae3e122fd755105b89e68d88e0e",
      "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
      "acceptance_rows": [
        "VR1",
        "VR2",
        "VR3",
        "VR4",
        "VR5",
        "VR6",
        "VR7",
        "VR8",
        "VR9",
        "VR10",
        "VR11",
        "VR12",
        "VR13",
        "VR14",
        "VR15",
        "VR16",
        "VR17",
        "VR18",
        "VR19",
        "VR20",
        "VR21",
        "VR22",
        "VR59",
        "VR60",
        "VR61"
      ],
      "verification": [
        {
          "id": "vr-c1-1-worktree-r1",
          "performer": "claude:bench-writer/vr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-author-20261001/1-worktree@137fb303",
            "digest": "sha256:e0a99d19673fa071a4aab6ce31a04e581a6233c21ef1905353526dc0ec9e9c6e",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,51282\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c1-1-probe-r1",
          "performer": "claude:bench-writer/vr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-author-20261001/1-probe@137fb303",
            "digest": "sha256:cb76d16fbf406599002d07b0d92643b9337e8419f0362e61c16114efdfd45de0",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,43"
          },
          "requirement": "1-probe",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'unapplicableFingerprint' --with '\"probe-placeholder\"' --package ./internal/worktree --run TestVerbResultFingerprintTreatsAPlaceholderAsAbsent",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-author-20261001/1-probe@137fb303",
              "digest": "sha256:cb76d16fbf406599002d07b0d92643b9337e8419f0362e61c16114efdfd45de0",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,43"
            }
          }
        },
        {
          "id": "vr-c1-author-probe-vr17-joins-r1",
          "performer": "claude:bench-writer/vr-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-author-20261001/vr17-joins-probe@137fb303",
            "digest": "sha256:52ab180cc662e1d83145275c16aae2873f295ae08b8395e8fcf063087facffc7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbRunnerPassesTheJoinsValue,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,60"
          },
          "requirement": "author-probe-VR17-joins",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'if call.joins != nil {' --with 'if false {' --package ./internal/worktree --run TestVerbRunnerPassesTheJoinsValue",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-author-20261001/vr17-joins-probe@137fb303",
              "digest": "sha256:52ab180cc662e1d83145275c16aae2873f295ae08b8395e8fcf063087facffc7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbRunnerPassesTheJoinsValue,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,60"
            }
          }
        },
        {
          "id": "vr-c1-1-worktree-r2",
          "performer": "claude:bench-writer/vr-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-repair-1-20261001/1-worktree@3a318142",
            "digest": "sha256:72cb722cbb113d863e4c1bf0439bc8e8c0c50a0b6739059a5f5af238f7370600",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,pass,50325\nfailures[0]{package,test,line}:\nskips[2]{package,test,reason}:"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "vr-c1-1-probe-r2",
          "performer": "claude:bench-writer/vr-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-t1-repair-1-20261001/1-probe@3a318142",
            "digest": "sha256:d4c62cc4f91e4a5a9a0f48e10c988ef86c408c95989af5576fede4b1d6e87bda",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,36"
          },
          "requirement": "1-probe",
          "command": "bench probe internal/worktree/verb_runner_test.go --swap 'unapplicableFingerprint' --with '\"probe-placeholder\"' --package ./internal/worktree --run TestVerbResultFingerprintTreatsAPlaceholderAsAbsent",
          "exit_code": 0,
          "probe": {
            "mutation": "swap",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude:agent/vr-t1-repair-1-20261001/1-probe@3a318142",
              "digest": "sha256:d4c62cc4f91e4a5a9a0f48e10c988ef86c408c95989af5576fede4b1d6e87bda",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbResultFingerprintTreatsAPlaceholderAsAbsent,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,36"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "vr-c1-standards-r1",
          "performer": "claude:bench-reviewer/vr-c1-standards-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-standards@4479b2fa",
            "digest": "sha256:19f02fddf42f5f19396b117b6a2e9a6b6ef8edf11c1b72e04404a52e97b6b28f",
            "excerpt": "Standards: 4 findings. Worst: the VR59 test pastes a second copy of the faulted-unclaimed fixture and its comment from clean_unclaimed_test.go."
          },
          "axis": "Standards",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "4479b2fa62bcb334f4b65aedfcd155cde31d66a6",
          "finding_ids": [
            "R1",
            "R2",
            "R3",
            "R4"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c1-spec-r1",
          "performer": "claude:bench-reviewer/vr-c1-spec-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-spec@4479b2fa",
            "digest": "sha256:6f57a4c2208052b9c43203b4c35da81b90400c0409e5d68580a15592d70382c3",
            "excerpt": "Spec: 3 findings. Worst: the VR9 first-64-hex mutant survives its own test; VR14 kills it at suite level. All 22 planned tests exist with exact names."
          },
          "axis": "Spec",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "4479b2fa62bcb334f4b65aedfcd155cde31d66a6",
          "finding_ids": [
            "R2",
            "R5"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c1-coverage-r1",
          "performer": "claude:bench-reviewer/vr-c1-coverage-r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-coverage@4479b2fa",
            "digest": "sha256:b054f944c1d41ef4c3d785f4257ce6521b29fb826df474d43a963f49e4013569",
            "excerpt": "Coverage: 2 findings. Worst: the record branch of readVerbFingerprint has no test for the absent rule or the conflict rule; a swap that returns the first record cell is silent."
          },
          "axis": "Coverage",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "4479b2fa62bcb334f4b65aedfcd155cde31d66a6",
          "finding_ids": [
            "R6"
          ],
          "supersedes": []
        },
        {
          "id": "vr-c1-standards-r2",
          "performer": "claude:bench-reviewer/vr-c1-standards-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-standards-r2@d44d28ce",
            "digest": "sha256:02c08c774b1445068209559be3c4397268999c530c936632f712774dd155e382",
            "excerpt": "Standards: 0 findings. All five folds hold; the shared fault builder in the runner check file is a judgment call, not a binding defect."
          },
          "axis": "Standards",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-standards-r1"
          ]
        },
        {
          "id": "vr-c1-spec-r2",
          "performer": "claude:bench-reviewer/vr-c1-spec-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/vr-c1-spec-r2@d44d28ce",
            "digest": "sha256:333dfc09f314f4275d1726d255de0f5ae78f45ede30466295b171d7182d8a1d5",
            "excerpt": "Spec: 2 findings on spec text. Folds R2, R5, R6 and the plan commit hold; the R6 tests lacked coverage rows, and VR1 wording names only the usage constant."
          },
          "axis": "Spec",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
          "finding_ids": [
            "R7",
            "R8"
          ],
          "supersedes": [
            "vr-c1-spec-r1"
          ]
        },
        {
          "id": "vr-c1-coverage-r2",
          "performer": "claude:bench-reviewer/vr-c1-coverage-r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "0dee07c82fe40e1420c5629289bd195d4e8ff5cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/vr-c1-coverage-r2@d44d28ce",
            "digest": "sha256:a9eacdfca84623ee5f8281963a2cf3e4333538d549aea8f4dabf925f6464fe75",
            "excerpt": "Coverage: 0 findings. All folds hold, every probe bit and restored, 688 top-level tests pass, and the skip set is the two socket subtests."
          },
          "axis": "Coverage",
          "base": "0c95c9447c20189f3f2155719ef965bffc339856",
          "tip": "d44d28ce6ff1964ffc4cb77d969cecdb881356eb",
          "finding_ids": [],
          "supersedes": [
            "vr-c1-coverage-r1"
          ]
        }
      ]
    }
  ],
  "completion": {
    "state": "pending",
    "source_digest": "",
    "performer": "",
    "reconciliation": {},
    "verification": []
  }
}
```
