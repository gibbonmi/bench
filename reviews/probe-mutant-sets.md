# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/probe-mutant-sets/spec.md",
  "plan_digest": "sha256:bb382785f355201ccac68fe63e38675f5033c694a0f62850ef28f7860ba80ee9",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PM-A",
      "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
      "tip": "868975d7ed65806ddf627da3821948b075a9cb74",
      "plan_digest": "sha256:bb382785f355201ccac68fe63e38675f5033c694a0f62850ef28f7860ba80ee9",
      "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
      "acceptance_rows": [
        "PM45",
        "PM46",
        "PM47",
        "PM48",
        "PM49",
        "PM50"
      ],
      "verification": [
        {
          "id": "pm-t01-probe-family",
          "performer": "claude:pm_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
            "digest": "sha256:8b44d0e3886e6458b0438e11b162ccaf452b4b4c3da9aa10aed308fdc5b2584d",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,868975d7ed65806ddf627da3821948b075a9cb74,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26266,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72427,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59960,229\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,71307,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t01-testreport",
          "performer": "claude:pm_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
            "digest": "sha256:8b44d0e3886e6458b0438e11b162ccaf452b4b4c3da9aa10aed308fdc5b2584d",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,868975d7ed65806ddf627da3821948b075a9cb74,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26266,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72427,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59960,229\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,71307,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t01-command-package",
          "performer": "claude:pm_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
            "digest": "sha256:8b44d0e3886e6458b0438e11b162ccaf452b4b4c3da9aa10aed308fdc5b2584d",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,868975d7ed65806ddf627da3821948b075a9cb74,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26266,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72427,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59960,229\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,71307,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t01-conformance-package",
          "performer": "claude:pm_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
            "digest": "sha256:8b44d0e3886e6458b0438e11b162ccaf452b4b4c3da9aa10aed308fdc5b2584d",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,868975d7ed65806ddf627da3821948b075a9cb74,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26266,335\n  github.com/gibbonmi/bench/internal/conformance,pass,72427,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59960,229\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,71307,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t01-root-conformance",
          "performer": "claude:pm_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
            "digest": "sha256:90e207b0f0d0a38c9fc7971793e6d608967d92c18f0861434b7b37b084249731",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,868975d7ed65806ddf627da3821948b075a9cb74,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16008,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t01-central-property",
          "performer": "claude:pm_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
            "digest": "sha256:1ca2580ff2102b585d0bdaa85e0cba26686f29567c386e0401dd04263a36a603",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/probe.go,swap,failed,1,TestProbeRepeatsThePackage,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestProbeRepeatsThePackage$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,579\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,283,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./other,^TestClampNegat… (817 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Pass only the last --package value from selectionArgs. TestProbeRepeatsThePackage must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "868975d7ed65806ddf627da3821948b075a9cb74",
              "digest": "sha256:1ca2580ff2102b585d0bdaa85e0cba26686f29567c386e0401dd04263a36a603",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/probe.go,swap,failed,1,TestProbeRepeatsThePackage,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestProbeRepeatsThePackage$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,579\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,283,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:29: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./other,^TestClampNegat… (817 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
