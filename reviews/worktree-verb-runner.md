# Worktree verb runner review record

## VR-C1 author evidence

Ticket 1 had a fresh `bench-writer` author, `vr-t1-author`, on opus at medium effort, with a cap of 3 attempts. The author started at `7872c5e4` and committed `137fb303` on a lane pass in the first attempt. The author then committed this record in a second commit.

The chunk pair is `0c95c944..137fb303`. The base is the `main` tip, because a plan commit is never a chunk base. The payload names the ticket commit as the tip, because a record cannot name its own commit. The coordinator moves the tip to the record commit. The source digest does not change, because the record file is outside the graded source.

The author wrote each test before the code that it grades, in three stages. The red and green log for each row follows:

- Stage 1: every runner function and reader was a stub, and the runner returned exit -1. The package run had 35 failures. All 22 planned tests failed, and VR1 failed in each of its 14 key subtests.
- Stage 2: the runner and `checkVerbCall` were complete, and the readers and must forms were stubs. VR1 to VR5, VR17 to VR19, VR21, and VR22 passed. VR6 to VR14, VR16, and VR59 failed for the reader reason, with 11 failures.
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

```bench-review-record
{
  "version": 2,
  "spec": "specs/worktree-verb-runner/spec.md",
  "plan_digest": "sha256:77c5099bae42000b631c569a15e3af6008579e9a60e93b691bf485f6636e9d2f",
  "implementation_session": "",
  "chunks": [
    {
      "id": "VR-C1",
      "base": "0c95c9447c20189f3f2155719ef965bffc339856",
      "tip": "137fb303135702af0819807d648c228b2ea33312",
      "plan_digest": "sha256:77c5099bae42000b631c569a15e3af6008579e9a60e93b691bf485f6636e9d2f",
      "source_digest": "db49c9755719fa79a394a329fdd5dedeb5372136",
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
        "VR59"
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
            "digest": "sha256:654259c95231dd11208253be28771fef52f30606fa4f246587e037cc4c1161fd",
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
            "digest": "sha256:9281fe083e8f7d3d4ed2405a2b9d283816264eb40866402580136ae909ba6122",
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
              "digest": "sha256:9281fe083e8f7d3d4ed2405a2b9d283816264eb40866402580136ae909ba6122",
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
            "digest": "sha256:e5057e77d067575ac77ddc8c33bd59bc2f2039837e80396a88e8daa674bdac0b",
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
              "digest": "sha256:e5057e77d067575ac77ddc8c33bd59bc2f2039837e80396a88e8daa674bdac0b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/worktree/verb_runner_test.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,TestVerbRunnerPassesTheJoinsValue,passed,1\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/worktree,fail,60"
            }
          }
        }
      ],
      "reviews": []
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
