# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/record-ticket/spec.md",
  "plan_digest": "sha256:e8251dd1728182f946187f20d4eaba706b5aef43e7b9466d453dc0a7f24489e1",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RT-RETAIN",
      "base": "53e60dc58f0432fdadf2b31b78d4e7a25f85570d",
      "tip": "9d989e126b4cbeac507b12ea3137250fd9c6b294",
      "plan_digest": "sha256:e8251dd1728182f946187f20d4eaba706b5aef43e7b9466d453dc0a7f24489e1",
      "source_digest": "e3239b89055da6012a766088982793175b60df59",
      "acceptance_rows": [
        "RT1",
        "RT2",
        "RT3",
        "RT4",
        "RT5",
        "RT6",
        "RT7",
        "RT8",
        "RT65"
      ],
      "verification": [
        {
          "id": "ft405_t01-retention",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:478e609c4e4b03228b63522691ab5613be1bce9db5810171f2c3779a005b2cde",
            "excerpt": "packages[11]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18816,345\n  github.com/gibbonmi/bench/internal/conformance,pass,81955,1694\n  github.com/gibbonmi/bench/internal/gate,pass,47840,489\n  github.com/gibbonmi/bench/internal/git,pass,1417,172\n  github.com/gibbonmi/bench/internal/handoff,pass,2960,56\n  github.com/gibbonmi/bench/internal/responsebound,pass,266,61\n  github.com/gibbonmi/bench/internal/responsebound/responseboundtest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/roadmap,pass,1905,176\n  github.com/gibbonmi/bench/internal/shift,pass,6375,116\n  github.com/gibbonmi/bench/internal/status,pass,8259,142\n  github.com/gibbonmi/bench/internal/worktree,pass,77044,1433\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "retention",
          "command": "bench test --package ./cmd/bench --package ./internal/responsebound/... --package ./internal/git --package ./internal/gate --package ./internal/handoff --package ./internal/roadmap --package ./internal/shift --package ./internal/status --package ./internal/worktree --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft405_t01-retain-root",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:2a33e4364ef1af5ff2ab111c6dbf515049beae83a45ab8986775641c82cbd5b1",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft405-build,9d989e126b4cbeac507b12ea3137250fd9c6b294,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15332,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "retain-root",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft405_t01-nested-exec",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:d01d525aa63ccb3d71293d82464fc70e3ad0a218a4800ce4fd159c3776a9262d",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft405-build,9d989e126b4cbeac507b12ea3137250fd9c6b294,true\ncheck[1]{name,kind,tests_run,subjects}:\n  system,system,2,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,539,2\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "nested-exec",
          "command": "bench test --check system --run TestExecNested",
          "exit_code": 0
        },
        {
          "id": "ft405_t01-record-skip",
          "performer": "claude:ft405_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e3239b89055da6012a766088982793175b60df59",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft405_t01",
            "digest": "sha256:b048b8cab438469381e9018b42eb30add7d60968249659c6b570e8c9a518fad2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,cmd/bench/retain.go,omit,failed,1,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,345\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,527\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,14822,345\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,\"retain_test.go:102: retained after [\\\"record\\\" \\\"chunk\\\"] = {Verb:record Args:[chunk] Exit:0 Head:726829397659fa99a2294cdda885d1cf0a718bd0 Dirty:[] Output:[108 105 110 101 32 48 49 10 108 105 110 101 32 48 50 10]}, want the test response {Verb:… (352 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "record-skip",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the record-family skip from retention",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft405_t01",
              "digest": "sha256:b048b8cab438469381e9018b42eb30add7d60968249659c6b570e8c9a518fad2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,cmd/bench/retain.go,omit,failed,1,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./cmd/bench,all,passed,345\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,527\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,fail,14822,345\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/cmd/bench,TestRetainSkipsRecordWorktreeAndRetiringCalls/record,\"retain_test.go:102: retained after [\\\"record\\\" \\\"chunk\\\"] = {Verb:record Args:[chunk] Exit:0 Head:726829397659fa99a2294cdda885d1cf0a718bd0 Dirty:[] Output:[108 105 110 101 32 48 49 10 108 105 110 101 32 48 50 10]}, want the test response {Verb:… (352 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
