# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/probe-mutant-sets/spec.md",
  "plan_digest": "sha256:9991026f819b1b4561ced20dd66fcea641105979c3797228b93b3e4eeac3f255",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PM-A",
      "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
      "tip": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
      "plan_digest": "sha256:faebeaa6ca8f219e2c4e6b61652d1e446266fd81a33e37d53f9b07e8d56f925e",
      "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
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
        },
        {
          "id": "pm-t01-rep1-probe-family",
          "performer": "claude:pm_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
            "digest": "sha256:46f717be1c2b5f3f945e1a6079f75bbb9b5dc823e509506d7fbdef4ecb403371",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cc2abf92b66b5430ed598624ba3c86aaa7f19d90,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22053,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71312,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59838,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,68328,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t01-rep1-testreport",
          "performer": "claude:pm_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
            "digest": "sha256:46f717be1c2b5f3f945e1a6079f75bbb9b5dc823e509506d7fbdef4ecb403371",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cc2abf92b66b5430ed598624ba3c86aaa7f19d90,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22053,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71312,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59838,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,68328,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t01-rep1-command-package",
          "performer": "claude:pm_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
            "digest": "sha256:46f717be1c2b5f3f945e1a6079f75bbb9b5dc823e509506d7fbdef4ecb403371",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cc2abf92b66b5430ed598624ba3c86aaa7f19d90,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22053,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71312,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59838,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,68328,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t01-rep1-conformance-package",
          "performer": "claude:pm_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
            "digest": "sha256:46f717be1c2b5f3f945e1a6079f75bbb9b5dc823e509506d7fbdef4ecb403371",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cc2abf92b66b5430ed598624ba3c86aaa7f19d90,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22053,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71312,1690\n  github.com/gibbonmi/bench/internal/probe,pass,59838,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,68328,232\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t01-rep1-root-conformance",
          "performer": "claude:pm_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
            "digest": "sha256:44d03a4c7dc4e308001731be08f5005747e79e17b9c0d394e7f2b0e30dfcfb85",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cc2abf92b66b5430ed598624ba3c86aaa7f19d90,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15606,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t01-rep1-central-property",
          "performer": "claude:pm_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
            "digest": "sha256:b26f673263bf9a79ceb0bc6a677d2708fcfe7bb9eed489e06ea94655afa4c1dd",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/probe.go,swap,failed,1,TestProbeRepeatsThePackage,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestProbeRepeatsThePackage$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,579\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,313,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./other,^TestClampNegat… (817 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
              "digest": "sha256:b26f673263bf9a79ceb0bc6a677d2708fcfe7bb9eed489e06ea94655afa4c1dd",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/probe.go,swap,failed,1,TestProbeRepeatsThePackage,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestProbeRepeatsThePackage$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,579\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,313,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestProbeRepeatsThePackage,\"package_test.go:31: stdout = (\\\"probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\\\\\\\\n  bit,clamp.go,swap,failed,1,TestClampNegative,yes\\\\\\\\nselection[1]{form,target,run,baseline,ran}:\\\\\\\\n  package,./other,^TestClampNegat… (817 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "pm-a-standards-r1",
          "performer": "claude:pm_a_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_a_standards",
            "digest": "sha256:0d6fa2e47b2ffd5fcd13b1544a3d6e8a29a3692a1eb8eaf302fe251028677077",
            "excerpt": "Standards axis claude:pm_a_standards (opus/high), PM-A round 1, frozen pair a6739258..14bec38c (chunk tip 868975d7), evidence sha256:d787d670 current.\ncount: 1 (no-op only; recorded with no finding ID)\nworst: S1\nS1 | no-op | 6 | refusal_test.go:225 holds a hand copy of PackageExpressionNote (producer selection_facts.go:93). Ticket 1 line 20 directs \"Update the independent note copy in internal/probe/refusal_test.go\", and the copy's red is recorded (TestProbeHelpSpellings test-first). No change.\nRefuted: comment prose (command.go:28-30, probe.go:140-141, selection_facts.go:90-92, selection_facts_test.go:148, package_test.go) is current-code and within 25 words; packageValueRule and probePackageNote have no references; BENCH-reference.md:260 staleness is owned by ticket 6; bothPackages red is in the central probe excerpt; the review record shape matches recordcmd/verification.go.\nAdvice: four floor rows share one combined excerpt; selection_facts_test.go:148 could say \"the shared package note\"; package_test.go:33 and :47 both count two starts; the argv hook repeats an incidental shell idiom.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
          "tip": "868975d7ed65806ddf627da3821948b075a9cb74",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "pm-a-spec-r1",
          "performer": "claude:pm_a_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_a_spec",
            "digest": "sha256:cb3230bd4ccef8c25c791116251c71032217d9d62a76791bca332861d6da140c",
            "excerpt": "Spec axis claude:pm_a_spec (opus/high), PM-A round 1, frozen pair a6739258..14bec38c (chunk tip 868975d7), evidence sha256:d787d670 current.\ncount: 0\nworst: none\nPM45 met: TestProbeRepeatsThePackage (internal/probe/package_test.go:69-92) requires 2 starts, each ending \" ./clamp ./other\"; selectionArgs loops parsed.Repeated[\"--package\"] (probe.go:142-150).\nPM46 met: TestProbePackagesShareOneBuildPerRun (package_test.go:95-108) asserts startCount == 2 and len(buildRecords) == 2.\nPM47 met: TestProbeRepeatsThePackage requires selectionRow(..., \"./clamp ./other\", ...).\nPM48 met: TestProbeHelpSpellings whole-help equality against the updated copy (refusal_test.go:224-244); \"run one call per package\" survives in no help, code, or guidance.\nPM49 met: probePackageNote and packageValueRule deleted; ProbeNotes starts with PackageExpressionNote (selection_facts.go:93-101); TestProbeNotesDeriveTheirTokens prefix and count 2.\nPM50 met: probeHelpSuffix derives from probe.Usage (cmd/bench/test_command.go:14); TestHelpInventoryIsComplete pins [--package <expr>]...\nScope: the spec edit rewrites only the PM45-PM47 seam cells, as bench-implement-spec.md:52 requires. No scope creep.\nAdvice: the PM45 suffix check passes if an argv repeats a package before the final \" ./clamp ./other\".\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
          "tip": "868975d7ed65806ddf627da3821948b075a9cb74",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "pm-a-coverage-r1",
          "performer": "claude:pm_a_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7f2262ac58ba52b36216c7d449600a955fb84029",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_a_coverage",
            "digest": "sha256:ca069ff566d046b630b5d0e1ea60617f6f55e4f4646807a9782251981c6cca6a",
            "excerpt": "Coverage axis claude:pm_a_coverage (opus/high), PM-A round 1, frozen pair a6739258..14bec38c (chunk tip 868975d7), evidence sha256:d787d670 current.\ncount: 2\nworst: C1\nC1 | auto-fix | 8 | Input: selectionArgs joins the repeated values into one --package value \"./clamp ./other\". Real Go reads one argument and fails as one missing directory (selection_facts.go:90-93). The stub hook records \"$*\" (package_test.go:72), which erases argument boundaries, and the assertion is HasSuffix(start, \" ./clamp ./other\") (package_test.go:86-89). Probe: swap the loop at probe.go:147 to join later values into args[1]; bench probe --package ./internal/probe read silent, 229 tests, 0 failed, restored yes. PM45 is not proven. Needed: record each argv element on its own line and assert the last two elements equal ./clamp and ./other.\nC2 | auto-fix | 6 | Input: --package ./clamp --package \"\". No test pins the empty-value refusal for a later occurrence of a repeatable flag (usage/parse.go:179; parse_test.go:419-447; refusal_test.go). Probe: guard swapped to skip a repeated occurrence; silent over probe, usage, testreport, worktree; restored yes. Current behavior is correct (exit 2, usage refusal). Needed: a refusal_test.go row with exit 2, subject unchanged, no run child. Hardening cycle.\nAdvice: duplicate --package values are undecided; --check beside a repeated --package is covered by the two-selections row.\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
          "tip": "868975d7ed65806ddf627da3821948b075a9cb74",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "pm-a-standards-r2",
          "performer": "claude:pm_a_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_a_standards_r2",
            "digest": "sha256:a9872d0749047d9a3be98fa89fffa7e42a347ae63a43e46b3459249a78cefbb9",
            "excerpt": "Standards axis claude:pm_a_standards_r2 (opus/high), PM-A confirming round 2, chunk pair a6739258..57336246, repair delta 14bec38c..cc2abf92 -- internal/, evidence sha256:da630e2d current.\ncount: 0\nworst: none\nC1 repair (package_test.go:20-49): the per-argument expectation strings.Fields(bothPackages) shares its source with the target cell; the hook uses the documented probetest $n contract; comments are timeless, tag-free, under 25 words.\nC2 repair (refusal_test.go:168-174): the expectation derives from usage.EmptyFlagValue and toon.Usage, as tree precedents do; requireRefused asserts exit 2, subject unchanged, empty home, no run child; both paths are on ticket 1 Writes.\nS1 stays closed as no-op: the delta does not touch probeHelpNotes, and ticket 1 line 20 directs the copy.\nAdvice: two hand-built argv hooks now exist (package_test.go:26, subject_check_test.go:81); a probetest helper fits if a third appears.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
          "tip": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
          "finding_ids": [],
          "supersedes": [
            "pm-a-standards-r1"
          ]
        },
        {
          "id": "pm-a-spec-r2",
          "performer": "claude:pm_a_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_a_spec_r2",
            "digest": "sha256:a8abc1bcdee8a807ed5d33607921335806d3276bed0e272a98b67c9aef6e42f6",
            "excerpt": "Spec axis claude:pm_a_spec_r2 (opus/high), PM-A confirming round 2, chunk pair a6739258..57336246, repair delta 14bec38c..57336246 (repair cc2abf92), evidence sha256:da630e2d current.\ncount: 0\nworst: none\nPM45 met: one argument per line per start (package_test.go:26), exactly 2 starts (:36), last two arguments equal [\"./clamp\",\"./other\"] (:39-48); C1 closed from the Spec side.\nPM47 met: selectionRow with bothPackages \"./clamp ./other\" (package_test.go:15, 30-31).\nC2 adds no behavior: TestProbeRefusesAnEmptyLaterPackage (refusal_test.go:168-174) pins the existing NoEmptyValue refusal (command.go:35, usage/parse.go:179-180) through usage.EmptyFlagValue.\nRecords: pm-t01-rep1-{probe-family,testreport,command-package,conformance-package,root-conformance,central-property} at cc2abf92, source digest da6caedf equals the chunk digest.\nAdvice: the central probe red shows the target cell, not the argument-boundary assertion; four floor entries share one excerpt.\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
          "tip": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
          "finding_ids": [],
          "supersedes": [
            "pm-a-spec-r1"
          ]
        },
        {
          "id": "pm-a-coverage-r2",
          "performer": "claude:pm_a_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "da6caedf244c2a6605118277ac7c11b57402a419",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_a_coverage_r2",
            "digest": "sha256:04bbbc81f5aee4b00478d63515104405a080db358da13cf4354e68c335c17b80",
            "excerpt": "Coverage axis claude:pm_a_coverage_r2 (opus/high), PM-A confirming round 2, chunk pair a6739258..57336246, repair delta 14bec38c..cc2abf92 -- internal/, evidence sha256:da630e2d current.\ncount: 0\nworst: none\nC1 closed: own probe at internal/testreport/command.go:195, operands append strings.Join(request.packageExprs, \" \"); the target cell is unchanged, so only the argv record can catch it. bench probe: bit, failed 1, TestProbeRepeatsThePackage at package_test.go:47, restored yes.\nC2 closed and red-capable: (1) drop NoEmptyValue on --package at internal/probe/command.go:35: bit at refusal_test.go:173, restored yes; (2) limit the empty-value check at internal/usage/parse.go:179 to the first occurrence: bit, restored yes.\nNo new defect: per-start files keyed by $n; the len guard prevents a slice panic; new imports are used.\nAdvice: a companion row with the empty value first would pin the other order; optional.\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "a6739258f810d1dcdd1efc2a028a0fe9726d8083",
          "tip": "cc2abf92b66b5430ed598624ba3c86aaa7f19d90",
          "finding_ids": [],
          "supersedes": [
            "pm-a-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "PM-B",
      "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
      "tip": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
      "plan_digest": "sha256:7770b3b86644c8d0349158481463350d60a3662bd9064a30fc86f0469d1a9c86",
      "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
      "acceptance_rows": [
        "PM27",
        "PM28",
        "PM29",
        "PM31",
        "PM70",
        "PM76",
        "PM79",
        "PM80"
      ],
      "verification": [
        {
          "id": "pm-t02-testreport",
          "performer": "claude:pm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
            "digest": "sha256:8deb2d6f37afc340e1ffa4ff9765a373088566697c1f3652c23c667f0d16abe9",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,27346d1166cecef1d115a4823ec30cadf3e3d2c8,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26414,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70888,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27531,29\n  github.com/gibbonmi/bench/internal/testreport,pass,65133,238\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\nelapsed 79s\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t02-runbinary",
          "performer": "claude:pm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
            "digest": "sha256:8deb2d6f37afc340e1ffa4ff9765a373088566697c1f3652c23c667f0d16abe9",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,27346d1166cecef1d115a4823ec30cadf3e3d2c8,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26414,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70888,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27531,29\n  github.com/gibbonmi/bench/internal/testreport,pass,65133,238\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\nelapsed 79s\n"
          },
          "requirement": "runbinary",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "pm-t02-command-package",
          "performer": "claude:pm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
            "digest": "sha256:8deb2d6f37afc340e1ffa4ff9765a373088566697c1f3652c23c667f0d16abe9",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,27346d1166cecef1d115a4823ec30cadf3e3d2c8,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26414,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70888,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27531,29\n  github.com/gibbonmi/bench/internal/testreport,pass,65133,238\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\nelapsed 79s\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t02-conformance-package",
          "performer": "claude:pm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
            "digest": "sha256:8deb2d6f37afc340e1ffa4ff9765a373088566697c1f3652c23c667f0d16abe9",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,27346d1166cecef1d115a4823ec30cadf3e3d2c8,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,26414,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70888,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27531,29\n  github.com/gibbonmi/bench/internal/testreport,pass,65133,238\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\nelapsed 79s\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t02-root-conformance",
          "performer": "claude:pm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
            "digest": "sha256:8b86f5e40d129e394258fc07d65e180e59d139522261de6a53565ef2ab57b98c",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,27346d1166cecef1d115a4823ec30cadf3e3d2c8,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15450,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit 0\nelapsed 17s\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t02-central-property",
          "performer": "claude:pm_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
            "digest": "sha256:c599e3c79c472f6a7952c69adc27e35a4522507fdb8433b46afe90738a12edd2",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/held_runs.go,swap,failed,1,TestHeldRunsOwnAFreshBuildWhenStale,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestHeldRunsOwnAFreshBuildWhenStale$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,3018,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestHeldRunsOwnAFreshBuildWhenStale,\"held_runs_test.go:24: builds read trees [\\\"v1\\\"], want [\\\"v1\\\" \\\"v2\\\"]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Reuse the held run binary without its validation before a later run. TestHeldRunsOwnAFreshBuildWhenStale must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
              "digest": "sha256:c599e3c79c472f6a7952c69adc27e35a4522507fdb8433b46afe90738a12edd2",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/held_runs.go,swap,failed,1,TestHeldRunsOwnAFreshBuildWhenStale,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestHeldRunsOwnAFreshBuildWhenStale$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,3018,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestHeldRunsOwnAFreshBuildWhenStale,\"held_runs_test.go:24: builds read trees [\\\"v1\\\"], want [\\\"v1\\\" \\\"v2\\\"]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pm-t02-rep1-testreport",
          "performer": "claude:pm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
            "digest": "sha256:ec0e50292d911698b1c826ccafacae3c589ebd2e5324ac4e7ea274d3d9500334",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,78584bdee070d6eff1aaf9ff8fe9da81961e0e66,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23887,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70760,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27070,29\n  github.com/gibbonmi/bench/internal/testreport,pass,66380,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t02-rep1-runbinary",
          "performer": "claude:pm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
            "digest": "sha256:ec0e50292d911698b1c826ccafacae3c589ebd2e5324ac4e7ea274d3d9500334",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,78584bdee070d6eff1aaf9ff8fe9da81961e0e66,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23887,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70760,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27070,29\n  github.com/gibbonmi/bench/internal/testreport,pass,66380,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\n"
          },
          "requirement": "runbinary",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0
        },
        {
          "id": "pm-t02-rep1-command-package",
          "performer": "claude:pm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
            "digest": "sha256:ec0e50292d911698b1c826ccafacae3c589ebd2e5324ac4e7ea274d3d9500334",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,78584bdee070d6eff1aaf9ff8fe9da81961e0e66,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23887,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70760,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27070,29\n  github.com/gibbonmi/bench/internal/testreport,pass,66380,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t02-rep1-conformance-package",
          "performer": "claude:pm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
            "digest": "sha256:ec0e50292d911698b1c826ccafacae3c589ebd2e5324ac4e7ea274d3d9500334",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/testreport --package ./internal/runbinary --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,78584bdee070d6eff1aaf9ff8fe9da81961e0e66,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23887,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70760,1690\n  github.com/gibbonmi/bench/internal/runbinary,pass,27070,29\n  github.com/gibbonmi/bench/internal/testreport,pass,66380,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nexit 0\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t02-rep1-root-conformance",
          "performer": "claude:pm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
            "digest": "sha256:088995b0e8b0b94bbc4ef8fb8a5b2f190f11e11744cee5e807e5aa4f24448afc",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,78584bdee070d6eff1aaf9ff8fe9da81961e0e66,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15501,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nexit 0\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t02-rep1-central-property",
          "performer": "claude:pm_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
            "digest": "sha256:ceec8cf03671acd16d6b9b979328c2bbaf6a53be5680bdee1d48843d45063b36",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/held_runs.go,swap,failed,1,TestHeldRunsOwnAFreshBuildWhenStale,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestHeldRunsOwnAFreshBuildWhenStale$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,3018,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestHeldRunsOwnAFreshBuildWhenStale,\"held_runs_test.go:23: builds read trees [\\\"v1\\\"], want [\\\"v1\\\" \\\"v2\\\"]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "Reuse the held run binary without its validation before a later run. TestHeldRunsOwnAFreshBuildWhenStale must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
              "digest": "sha256:ceec8cf03671acd16d6b9b979328c2bbaf6a53be5680bdee1d48843d45063b36",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/testreport/held_runs.go,swap,failed,1,TestHeldRunsOwnAFreshBuildWhenStale,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/testreport,^TestHeldRunsOwnAFreshBuildWhenStale$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,338\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/testreport,fail,3018,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/testreport,TestHeldRunsOwnAFreshBuildWhenStale,\"held_runs_test.go:23: builds read trees [\\\"v1\\\"], want [\\\"v1\\\" \\\"v2\\\"]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "pm-b-standards-r1",
          "performer": "claude:pm_b_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_b_standards",
            "digest": "sha256:4a4eeed43e46b3969426e0d5161dd28747c246ced6e3f95eb9664e081328ef49",
            "excerpt": "Standards axis claude:pm_b_standards (opus/high), PM-B round 1, frozen pair e5ca6b0c..1b8f4f66 (chunk tip 27346d11), evidence sha256:c5d89a03 current.\ncount: 1\nworst: S1\nS1 | auto-fix | 4 | AGENTS.md one source per fact (\"a fixture harness pasted N times\"). TestHeldRunsInterruptACancelledBuild (held_runs_test.go:105-118: pid file, t.Setenv(cancelGoPIDEnv), t.Setenv(cancelParkingPath), waitForPID, SIGKILL cleanup) is a third copy of the arm-and-wait wiring at cancel_test.go:100-122 and :229-251. Put the env contract and the wait in one cancel_test.go helper and call it from all three sites.\nClean: no copied runbinary or freshness rule (Factory.Validate, Selection.Owned); the command.go split keeps one owner per step; runBinaryFactory follows the laneRunBinary precedent; comments carry no tags; the interrupted line prefix is incidental text; record entries follow the PM-A precedent.\nAdvice: an interruptedSelection helper beside selectionFailure in outcome.go; the Validate doc names its callers; TestHeldRunsReuseOnAFreshVerify hardcodes that Own validates once; Run uses pointer identity for a fresh selection.\nImplementation command: the ticket line \"reuses the parking builder\" read as the builder only; name its arm-and-wait wiring too.\n"
          },
          "axis": "Standards",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "pm-b-spec-r1",
          "performer": "claude:pm_b_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_b_spec",
            "digest": "sha256:059be10b4939b1670cf8d4f3e25db0bd6cb20c099cf1fbd5bd13894558cb58ed",
            "excerpt": "Spec axis claude:pm_b_spec (opus/high), PM-B round 1, frozen pair e5ca6b0c..1b8f4f66 (chunk tip 27346d11), evidence sha256:c5d89a03 current.\ncount: 0\nworst: none\nPM27 met: TestHeldRunsOwnAFreshBuildWhenStale, read == [v1 v2], Builds 2.\nPM28 met: TestHeldRunsReuseOnAFreshVerify, verified == [path x3], one build.\nPM29 met: TestHeldRunsRefuseAStaleInheritedBinary, refused, no builder call, Builds 0.\nPM31 met: TestHeldRunsRemoveEachOwnedBuild, only the held directory after a stale run, none after Close.\nPM70 met: TestHeldRunsInterruptACancelledBuild.\nPM76 met: TestHeldRunsInterruptACancelledValidation, the cancel fires on the held validation.\nHeld runs and Trust chain hold: selectRunBinary first, Factory.Validate later (no copy), Selection.Owned branch, the fresh build closes on return; Validate reaches freshness.VerifyAt.\nFlags: (a) Builds counts built binaries, consistent with \"counts each binary that it builds\"; (b) the interrupt rule on the first selection is required by \"each selection, build, and Go child\"; (c) undecided by the spec and unobservable at the probe seam; (d) Execute keeps its order and mapping.\nScope creep: none.\nAdvice: Run enters runSelected when the context was cancelled after a successful validation or build; a ctx.Err check before runSelected would close it. No test cancels during the first held selection.\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "pm-b-coverage-r1",
          "performer": "claude:pm_b_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "731f556a019b125c5317003c4dd37772ab9090ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_b_coverage",
            "digest": "sha256:2504ff924329e43b77889b1c6c33efc015d76e51377e6a8b22e2ac62a291ac1b",
            "excerpt": "Coverage axis claude:pm_b_coverage (opus/high), PM-B round 1, frozen pair e5ca6b0c..1b8f4f66 (chunk tip 27346d11), evidence sha256:c5d89a03 current.\ncount: 2\nworst: C2\nC2 | auto-fix | 8 | held_runs.go:78: a stale owned held binary whose fresh build fails without a cancel is unpinned. Bypass --swap \"return nil, err\" --with \"return h.held, nil\" --occurrence 3 (the run uses the stale held binary) read silent, passed, restored yes over ./internal/testreport (238 tests). Spec line 248 and story 24; the D1 false-silent hazard. Needed: a held_runs_test.go row where the second run's stub builder returns runbinary.BuildError with no cancel; assert build-failed and no Go child.\nC1 | auto-fix | 8 | held_runs.go:34: a cancel during the first run's selection or build is unpinned (spec line 234, ticket line 25). Bypass --swap \"if h.ctx.Err() != nil {\" --with \"if h.held != nil && h.ctx.Err() != nil {\" read silent, passed, restored yes over the package. Needed: a row that parks the first build, cancels, and asserts interrupted, the parked directory gone, and Builds() 0.\nBit (restored yes): selection != h.held to selection == nil (TestHeldRunsRemoveEachOwnedBuild); Execute NotifyCancel replaced by a plain WithCancel (TestCommandDrainsBuilderGroupOnCancelSignal, 3 subtests).\nAdvice: a cancel during a stale owned validation still starts the fresh build; build count 0 for prose is pinned only at the probe seam (PM32).\nImplementation command: no command change; the gaps come from the spec map and edge inventory.\n"
          },
          "axis": "Coverage",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "27346d1166cecef1d115a4823ec30cadf3e3d2c8",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "pm-b-standards-r2",
          "performer": "claude:pm_b_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_b_standards_r2",
            "digest": "sha256:d042be7875726b5dc942e6465c956e96f22741f8da85e57692d8974919442dbc",
            "excerpt": "Standards axis claude:pm_b_standards_r2 (opus/high), PM-B confirming round 2, chunk pair e5ca6b0c..c20478f5, repair delta 1b8f4f66..78584bde -- internal/, evidence sha256:086431d9 current.\ncount: 0\nworst: none\nS1 confirmed fixed: armParking (cancel_test.go:102-105) owns the pid file; parking.env (:108-110) is the one pairing of cancelGoPIDEnv and cancelParkingPath; parking.descendant (:122-127) the one pairing of waitForPID and the SIGKILL cleanup; signalGoHopCancelHelper, signalCancelHelper, and parkedBuild.interrupt all call it; an rg sweep over internal/ finds no other parent-side copy.\nNew helpers: replaceAfter (held_runs_test.go:231-239) and park/parkedBuild/interrupt (:242-284) each hold one fact; the two new tests reuse them; TestHeldRunsFailAStaleRunWhoseBuildFails uses the production runbinary.BuildError type. Comments are timeless and tag-free.\nAdvice: parking.setenv round-trips env() through strings.Cut; \"sealed\" names both a count and a tree state in two doc comments.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
          "finding_ids": [],
          "supersedes": [
            "pm-b-standards-r1"
          ]
        },
        {
          "id": "pm-b-spec-r2",
          "performer": "claude:pm_b_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_b_spec_r2",
            "digest": "sha256:9e77dc48e58874b26439f7eefb79ad218eda353ea98a278d68292df0ae1bc11a",
            "excerpt": "Spec axis claude:pm_b_spec_r2 (opus/high), PM-B confirming round 2, chunk pair e5ca6b0c..c20478f5, delta 1b8f4f66..c20478f5, evidence sha256:086431d9 current.\ncount: 1\nworst: P1\nP1 | auto-fix | 6 | .bench/BENCH.md:124 \"Before dispatch, record every plan or gate expansion with bench learning\". No learning exists for plan expansion 2c88de8f (PM79, PM80, ticket 2 Covers, PM-B chunk row). Needed: the orchestrator runs bench learning with what changed, why, and verification.\nPlan expansion graded correct and in scope: PM79 follows spec.md:248 and D1, maps build-failed as Execute does, and does not duplicate PM44; PM80 follows spec.md:232/:234 and is a separate branch from PM70; Builds 0 matches spec.md:252.\nPM79 met by TestHeldRunsFailAStaleRunWhoseBuildFails (held_runs_test.go:118-134). PM80 met by TestHeldRunsInterruptACancelledFirstBuild (:105-116, Builds 0 at :113).\npm-t02-rep1 entries: all six PM-B requirements at 78584bde, source digest d0542038, pass; the central probe bit, restored yes.\nAdvice: five non-probe entries share one combined excerpt.\nImplementation command: P1 is the orchestrator's omission; the plan-expansion step should run bench learning with the plan commit, before the dispatch.\n"
          },
          "axis": "Spec",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
          "finding_ids": [
            "P1"
          ],
          "supersedes": [
            "pm-b-spec-r1"
          ]
        },
        {
          "id": "pm-b-coverage-r2",
          "performer": "claude:pm_b_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_b_coverage_r2",
            "digest": "sha256:8377cf917972a8d69b8ab0839178a0423d361578aa8d4b306f1061880169cb91",
            "excerpt": "Coverage axis claude:pm_b_coverage_r2 (opus/high), PM-B confirming round 2, chunk pair e5ca6b0c..c20478f5, repair delta 1b8f4f66..78584bde -- internal/ (test files only), evidence sha256:086431d9 current.\ncount: 0\nworst: none\nC2/PM79 closed: held_runs.go:37 swap runs runSelected on the stale binary before selectionFailure: bit at held_runs_test.go:132 \"Go child marker after the failed build\", restored yes.\nC1/PM80 closed: held_runs.go:59 h.builds++ before selectRunBinary: bit at held_runs_test.go:114 \"Builds() = 1, want 0\", restored yes.\nS1 red-capable: command.go:337 omit drainGoProcessGroup: bit (TestFocusedGoTestDrainsGroupOnCancelSignal); runbinary.go:309/:311 pid-only kill and no drain: bit on all 3 signal subtests. All restored yes.\nTwo own probe attempts read invalid (a removed brace; a tab mismatch) and were retried in a valid form.\nAdvice: omitting drainBuilderGroup alone at runbinary.go:307 or :311 leaves the cancel test green; older code outside the delta.\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
          "finding_ids": [],
          "supersedes": [
            "pm-b-coverage-r1"
          ]
        },
        {
          "id": "pm-b-spec-r2b",
          "performer": "claude:pm_b_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d054203893828c026c8c0357a76eea9b48f31919",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_b_spec_r2",
            "digest": "sha256:e857acd1b45d3f867dc3acb89b24d04eb41f223ffa20834b46fd0a2ff242d976",
            "excerpt": "Spec axis claude:pm_b_spec_r2 (opus/high), PM-B round 2 narrow reaffirmation of P1, chunk pair e5ca6b0c..c20478f5.\ncount: 0\nworst: none\nP1 closed, confidence 8: capture/learnings.md:272-275 \"2026-10-13 — Plan expansion for the PM-B repair of probe-mutant-sets [open]\" states what changed (plan 2c88de8f: PM79, PM80, ticket 2 lines, chunk row, edge sentence), why (two Held runs behaviors with no row), and verification (repair 78584bde, each bypass probe bit). The entry records its own lateness and proposes the process fix. No repository source changed.\n"
          },
          "axis": "Spec",
          "base": "e5ca6b0c50fd357d5233499495af5774ae566689",
          "tip": "78584bdee070d6eff1aaf9ff8fe9da81961e0e66",
          "finding_ids": [],
          "supersedes": [
            "pm-b-spec-r2"
          ]
        }
      ]
    },
    {
      "id": "PM-C1",
      "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
      "tip": "1c524805cea9dd907c4e306d94f1523463a97fc1",
      "plan_digest": "sha256:da6c45a4a51006e8483f681cd55a50dbb09d4c8fff29378f6724059c54511d88",
      "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
      "acceptance_rows": [],
      "verification": [
        {
          "id": "pm-t03-probe-family",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
            "digest": "sha256:a15736d3aaf23038867eda7064b9a44ef1a04c9c134f6509f38645fdff294e9e",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,dd3cdaecb1d9d1a3a57bcd89707629380f27d177,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17966,335\n  github.com/gibbonmi/bench/internal/conformance,pass,62116,1690\n  github.com/gibbonmi/bench/internal/probe,pass,49989,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t03-command-package",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
            "digest": "sha256:a15736d3aaf23038867eda7064b9a44ef1a04c9c134f6509f38645fdff294e9e",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,dd3cdaecb1d9d1a3a57bcd89707629380f27d177,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17966,335\n  github.com/gibbonmi/bench/internal/conformance,pass,62116,1690\n  github.com/gibbonmi/bench/internal/probe,pass,49989,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t03-conformance-package",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
            "digest": "sha256:a15736d3aaf23038867eda7064b9a44ef1a04c9c134f6509f38645fdff294e9e",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,dd3cdaecb1d9d1a3a57bcd89707629380f27d177,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17966,335\n  github.com/gibbonmi/bench/internal/conformance,pass,62116,1690\n  github.com/gibbonmi/bench/internal/probe,pass,49989,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t03-root-conformance",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
            "digest": "sha256:2dd823505a465d5e216ca8f4e80ac6205200037ebfd5feaa3913ab29e269b002",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,dd3cdaecb1d9d1a3a57bcd89707629380f27d177,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15374,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t03-central-property",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
            "digest": "sha256:ebe2cd2f37c013c2500a56e3303ef76c2d923d06c8a87d22b562f7c21264be30",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/command.go,swap,failed,1,TestExpectRequiresTheNamedWitness/prose_check,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestExpectRequiresTheNamedWitness$,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,597\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,1852,7\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:69: refusal = (\\\"error: probe subject not graded — clamp.go is not a prose subject\\\\\\\\n\\\", 1), want (\\\"usage: bench probe <file> (--swap <old> (--with <new> | --with-file <path>) [--occurrence <n>] | --omit <old> [--occurrence <n>… (589 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the prose --expect case from the cross-half rule. TestExpectRequiresTheNamedWitness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
              "digest": "sha256:ebe2cd2f37c013c2500a56e3303ef76c2d923d06c8a87d22b562f7c21264be30",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/command.go,swap,failed,1,TestExpectRequiresTheNamedWitness/prose_check,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestExpectRequiresTheNamedWitness$,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,597\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,1852,7\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:69: refusal = (\\\"error: probe subject not graded — clamp.go is not a prose subject\\\\\\\\n\\\", 1), want (\\\"usage: bench probe <file> (--swap <old> (--with <new> | --with-file <path>) [--occurrence <n>] | --omit <old> [--occurrence <n>… (589 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pm-t03-r2-probe-family",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "1c524805cea9dd907c4e306d94f1523463a97fc1",
            "digest": "sha256:0d6180ded17932123c3b09f64b01e2ebaec5d67793aabb66e24af6eec3bf8c81",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,1c524805cea9dd907c4e306d94f1523463a97fc1,true\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17402,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61501,1690\n  github.com/gibbonmi/bench/internal/probe,pass,49102,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t03-r2-command-package",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "1c524805cea9dd907c4e306d94f1523463a97fc1",
            "digest": "sha256:0d6180ded17932123c3b09f64b01e2ebaec5d67793aabb66e24af6eec3bf8c81",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,1c524805cea9dd907c4e306d94f1523463a97fc1,true\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17402,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61501,1690\n  github.com/gibbonmi/bench/internal/probe,pass,49102,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t03-r2-conformance-package",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "1c524805cea9dd907c4e306d94f1523463a97fc1",
            "digest": "sha256:0d6180ded17932123c3b09f64b01e2ebaec5d67793aabb66e24af6eec3bf8c81",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,1c524805cea9dd907c4e306d94f1523463a97fc1,true\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17402,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61501,1690\n  github.com/gibbonmi/bench/internal/probe,pass,49102,230\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t03-r2-root-conformance",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "1c524805cea9dd907c4e306d94f1523463a97fc1",
            "digest": "sha256:528e9a29c9d299e4e7f2c0d8093b5ad682dfaafdc044661cd0d8144c5317b07b",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,1c524805cea9dd907c4e306d94f1523463a97fc1,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15411,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t03-r2-central-property",
          "performer": "claude:pm_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "1c524805cea9dd907c4e306d94f1523463a97fc1",
            "digest": "sha256:843df0049bb48e1e14f5bf71c24c975925a0e937fe2d5a1de2e8391a02cb5c5d",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/command.go,swap,failed,1,TestExpectRequiresTheNamedWitness/prose_check,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestExpectRequiresTheNamedWitness$,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,597\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,1557,7\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:69: refusal = (\\\"error: probe subject not graded — clamp.go is not a prose subject\\\\\\\\n\\\", 1), want (\\\"usage: bench probe <file> (--swap <old> (--with <new> | --with-file <path>) [--occurrence <n>] | --omit <old> [--occurrence <n>… (589 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop the prose --expect case from the cross-half rule. TestExpectRequiresTheNamedWitness must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "1c524805cea9dd907c4e306d94f1523463a97fc1",
              "digest": "sha256:843df0049bb48e1e14f5bf71c24c975925a0e937fe2d5a1de2e8391a02cb5c5d",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/command.go,swap,failed,1,TestExpectRequiresTheNamedWitness/prose_check,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestExpectRequiresTheNamedWitness$,passed,7\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,597\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,1557,7\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestExpectRequiresTheNamedWitness/prose_check,\"verdict_test.go:69: refusal = (\\\"error: probe subject not graded — clamp.go is not a prose subject\\\\\\\\n\\\", 1), want (\\\"usage: bench probe <file> (--swap <old> (--with <new> | --with-file <path>) [--occurrence <n>] | --omit <old> [--occurrence <n>… (589 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "pm-c1-standards-r1",
          "performer": "claude:pm_c1_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c1_standards",
            "digest": "sha256:368369e2d9a8d9b518b3c5ac2ace3dbc6df96039e7b44f8dbf4e856a1e5258c7",
            "excerpt": "Standards axis claude:pm_c1_standards (opus/high), PM-C1 round 1, frozen pair 32dc3858..c4ec0f13 (chunk tip dd3cdaec), evidence sha256:0ae93234 current.\ncount: 1\nworst: S1\nS1 | auto-fix | 5 | internal/probe/command.go:186-187: the expectFitsSelection doc comment says \"an --expect in mutationArgs must name a test that the selection runs\", but the body (189-190) only refuses --expect under --check prose (spec line 164). craft-comments: a comment describes the code as it stands. Reword to the narrow rule. Comment-only, evidence-only under bounded-repair-policy line 47.\nClean: each mutation flag declared once (mutationFlags 22-32, expectFlag 35); no rule copied between the halves; moved comments accurate; behavior preserved through short-circuit order; no tags.\nAdvice: mutation_input.go:25 goes stale in ticket 4; mixed flag-name lookup styles; the mutationFlags comment anticipates ticket 4's grammar; parameter names mutationArgs/selection.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
          "tip": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "pm-c1-spec-r1",
          "performer": "claude:pm_c1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c1_spec",
            "digest": "sha256:040d825d542d90a7b7e5268eb5f0cd27b2feb04232ef53ba328c3de71e907389",
            "excerpt": "Spec axis claude:pm_c1_spec (opus/high), PM-C1 round 1, frozen pair 32dc3858..c4ec0f13 (chunk tip dd3cdaec), evidence sha256:0ae93234 current.\ncount: 0\nworst: none\nOne flag list: mutationFlags (command.go:22-32) declares each mutation flag once; grammar reads it through slices.Concat (:44); expectFlag (:35) stands alone for ticket 4's entry grammar.\nFlag set and order unchanged; Usage not in the diff, so the help inventory line cannot change.\nTwo halves: mutationFormFits (136-168) reads only mutation flags; selectionFormFits (172-184) reads only --package, --run, --check; gradeForms runs both.\nRefusals: every old case condition carries over; every failure returns usageLine() at exit 2.\nCross-half rule: expectFitsSelection(mutationArgs, selection) (188-191) reproduces the prose --expect refusal; called as (parsed, parsed) at the top level; fits ticket 4 line 16.\nNothing of ticket 4 built; no test file changed.\nAdvice: the expectFitsSelection doc comment's first sentence promises more than the function checks.\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
          "tip": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "pm-c1-coverage-r1",
          "performer": "claude:pm_c1_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "843ac87c4e23efe5b5de4e73516218b7a953557c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c1_coverage",
            "digest": "sha256:18cdf58dadade5a4dec88d6cfe23955c25794ac6b57c0fee1dd50091ae1640ad",
            "excerpt": "Coverage axis claude:pm_c1_coverage (opus/high), PM-C1 round 1, frozen pair 32dc3858..c4ec0f13 (chunk tip dd3cdaec), evidence sha256:0ae93234 current.\ncount: 0\nworst: none\nRules of base gradeForms and their pins: R1 !oneGroup (refusal_test.go:145-158); R3 insert texts != 1 (edit_test.go:84, zero texts only); R4 (edit_test.go:83,88-90; refusal_test.go:154-157); R5 combined && withFile (combined_test.go:66-72); R6-R9 occurrence (edit_test.go:66-71, 85-87); R10 (refusal_test.go:142-143); R11 (refusal_test.go:144, outcome_test.go:200); R12 prose --expect (verdict_test.go:66-70). R2 withs > 0 && withFile is covered by R3-R5 and cannot be the only refusing rule.\nProbes (restored yes): R3 upper bound dropped: silent (unpinned at base); R2 to false: silent (redundant rule); R8 control, remove deleteGroup from the occurrence case: bit.\nRefactor safety: short-circuit order kept; has(x) equals the Flags lookup; slices.Concat keeps order; Usage unchanged.\nAdvice: no row for an insert with two texts (pre-existing); R2 is redundant; expectFitsSelection argument roles are not pinned until ticket 4.\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
          "tip": "dd3cdaecb1d9d1a3a57bcd89707629380f27d177",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "pm-c1-standards-r1b",
          "performer": "claude:pm_c1_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c1_standards",
            "digest": "sha256:f61c816394f77ae3ff009af7022d01a6f441ff616d08674f3b91ee968b60df54",
            "excerpt": "Standards axis claude:pm_c1_standards (opus/high), PM-C1 narrow reaffirmation of S1 at the correction tip, delta dd3cdaec..1c524805.\ncount: 0\nworst: none\nS1 closed, confidence 9: internal/probe/command.go:186-187 now reads \"a prose run names no test, so an --expect in mutationArgs refuses under a prose selection\", which matches the body at 189-190 and spec line 164. The delta is a two-line comment rewrite with no code change and no new defect.\n"
          },
          "axis": "Standards",
          "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
          "tip": "1c524805cea9dd907c4e306d94f1523463a97fc1",
          "finding_ids": [],
          "supersedes": [
            "pm-c1-standards-r1"
          ]
        },
        {
          "id": "pm-c1-spec-r1b",
          "performer": "claude:pm_c1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c1_spec",
            "digest": "sha256:99357e7af39387988db20c0e182a0d149f8811546279f0763644143dc8f78494",
            "excerpt": "Spec axis claude:pm_c1_spec (opus/high), PM-C1 reaffirmation at the correction tip, delta dd3cdaec..1c524805.\ncount: 0\nworst: none\nRound 1 result holds, confidence 9: the only code change is the two-line doc comment at internal/probe/command.go:186-187; the body at 188-191 is unchanged, so the flag-list, two-half split, byte-identity, and ticket 4 seam checks still apply. The new comment resolves the round 1 advice.\n"
          },
          "axis": "Spec",
          "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
          "tip": "1c524805cea9dd907c4e306d94f1523463a97fc1",
          "finding_ids": [],
          "supersedes": [
            "pm-c1-spec-r1"
          ]
        },
        {
          "id": "pm-c1-coverage-r1b",
          "performer": "claude:pm_c1_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "06e1253ebab968f10343e16bf7b195a954f3062d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c1_coverage",
            "digest": "sha256:d314416f5e94d415555eb3bc17c2d2cd81c424ee72cb68bf0b916db278f16943",
            "excerpt": "Coverage axis claude:pm_c1_coverage (opus/high), PM-C1 reaffirmation at the correction tip, delta dd3cdaec..1c524805.\ncount: 0\nworst: none\nRound 1 result holds, confidence 9: internal/probe/command.go has one hunk @@ -183,8 +183,8 @@ that changes only the doc comment at 186-187; the body at 188-190 is byte-identical; no executable line changed, so the rule-to-test mapping and the three probe verdicts still apply.\n"
          },
          "axis": "Coverage",
          "base": "32dc3858c6e91b286869293d5a86ecec7f3ece79",
          "tip": "1c524805cea9dd907c4e306d94f1523463a97fc1",
          "finding_ids": [],
          "supersedes": [
            "pm-c1-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "PM-C2",
      "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
      "tip": "485fa5e295bec6471009169ba0ee9bfd140580dd",
      "plan_digest": "sha256:12c94ca56db2698862ce62c097467259f0b90432a2f66febcac4f339bb145b72",
      "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
      "acceptance_rows": [
        "PM1",
        "PM2",
        "PM3",
        "PM4",
        "PM5",
        "PM6",
        "PM7",
        "PM8",
        "PM9",
        "PM10",
        "PM11",
        "PM12",
        "PM13",
        "PM14",
        "PM15",
        "PM16",
        "PM17",
        "PM18",
        "PM19",
        "PM20",
        "PM21",
        "PM22",
        "PM23",
        "PM24",
        "PM25",
        "PM26",
        "PM30",
        "PM32",
        "PM33",
        "PM35",
        "PM36",
        "PM37",
        "PM38",
        "PM39",
        "PM40",
        "PM41",
        "PM43",
        "PM44",
        "PM51",
        "PM67",
        "PM69",
        "PM71",
        "PM72",
        "PM77",
        "PM78",
        "PM81",
        "PM82",
        "PM83",
        "PM84",
        "PM85",
        "PM86"
      ],
      "verification": [
        {
          "id": "pm-t04-probe-family",
          "performer": "claude:pm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
            "digest": "sha256:0e36f96848601b67695e7f470b8a9f0587bc528770f97f471ecf0c537fcdb147",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21343,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71491,1690\n  github.com/gibbonmi/bench/internal/probe,pass,72263,303\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74660,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t04-testreport",
          "performer": "claude:pm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
            "digest": "sha256:0e36f96848601b67695e7f470b8a9f0587bc528770f97f471ecf0c537fcdb147",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21343,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71491,1690\n  github.com/gibbonmi/bench/internal/probe,pass,72263,303\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74660,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t04-command-package",
          "performer": "claude:pm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
            "digest": "sha256:0e36f96848601b67695e7f470b8a9f0587bc528770f97f471ecf0c537fcdb147",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21343,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71491,1690\n  github.com/gibbonmi/bench/internal/probe,pass,72263,303\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74660,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t04-conformance-package",
          "performer": "claude:pm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
            "digest": "sha256:0e36f96848601b67695e7f470b8a9f0587bc528770f97f471ecf0c537fcdb147",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21343,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71491,1690\n  github.com/gibbonmi/bench/internal/probe,pass,72263,303\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74660,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-root-conformance",
          "performer": "claude:pm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
            "digest": "sha256:6ae14dfdd20d620bcd986b6ddec5666f89958191ebf9b0d270b839b41d104781",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15544,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-central-property",
          "performer": "claude:pm_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
            "digest": "sha256:a73d0e87a905037e033cfac67d2ec1ef02350a4c916f4233803c11341de11537",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,630\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,387,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:143: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (1006 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the restore between two mutants, so the second mutant stacks on the first. TestMutantSetRunsEachMutantAlone must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
              "digest": "sha256:a73d0e87a905037e033cfac67d2ec1ef02350a4c916f4233803c11341de11537",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,630\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,387,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:143: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (1006 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pm-t04-rep1-probe-family",
          "performer": "claude:pm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
            "digest": "sha256:5591f0000bf83fb0294d09886c933c4e98e981512d4eaa23a880624f8dbd85bb",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,4b9543cf396f0709957f44e812eafcdb5b37775a,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22264,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71722,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74105,311\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,72746,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nelapsed=79s\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep1-testreport",
          "performer": "claude:pm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
            "digest": "sha256:5591f0000bf83fb0294d09886c933c4e98e981512d4eaa23a880624f8dbd85bb",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,4b9543cf396f0709957f44e812eafcdb5b37775a,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22264,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71722,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74105,311\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,72746,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nelapsed=79s\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep1-command-package",
          "performer": "claude:pm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
            "digest": "sha256:5591f0000bf83fb0294d09886c933c4e98e981512d4eaa23a880624f8dbd85bb",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,4b9543cf396f0709957f44e812eafcdb5b37775a,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22264,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71722,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74105,311\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,72746,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nelapsed=79s\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep1-conformance-package",
          "performer": "claude:pm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
            "digest": "sha256:5591f0000bf83fb0294d09886c933c4e98e981512d4eaa23a880624f8dbd85bb",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,4b9543cf396f0709957f44e812eafcdb5b37775a,false\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,22264,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71722,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74105,311\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,72746,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\nelapsed=79s\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep1-root-conformance",
          "performer": "claude:pm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
            "digest": "sha256:2fc46727e4db748a9c69a59fb5e1207284f4a2e9848aa6b91c5a248b20f206f2",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,4b9543cf396f0709957f44e812eafcdb5b37775a,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15331,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\nelapsed=19s\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep1-central-property",
          "performer": "claude:pm_t04_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
            "digest": "sha256:d17351d0518adcd9d4970dcda2e292c9de513c13827b051cb459aa1414e54382",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,629\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,380,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:149: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (869 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the restore between two mutants, so the second mutant stacks on the first. TestMutantSetRunsEachMutantAlone must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "4b9543cf396f0709957f44e812eafcdb5b37775a",
              "digest": "sha256:d17351d0518adcd9d4970dcda2e292c9de513c13827b051cb459aa1414e54382",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,629\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,380,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:149: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (869 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pm-t04-rep2-probe-family",
          "performer": "claude:pm_t04_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
            "digest": "sha256:85725c7b6501a2770cc17034ba253bfd3e40fc4f9d2d952e1f11e497fc1eb991",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,9d484121013c213d6d9f7808a23d157d8408dedd,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20380,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70508,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74750,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,75040,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep2-testreport",
          "performer": "claude:pm_t04_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
            "digest": "sha256:85725c7b6501a2770cc17034ba253bfd3e40fc4f9d2d952e1f11e497fc1eb991",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,9d484121013c213d6d9f7808a23d157d8408dedd,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20380,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70508,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74750,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,75040,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep2-command-package",
          "performer": "claude:pm_t04_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
            "digest": "sha256:85725c7b6501a2770cc17034ba253bfd3e40fc4f9d2d952e1f11e497fc1eb991",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,9d484121013c213d6d9f7808a23d157d8408dedd,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20380,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70508,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74750,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,75040,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep2-conformance-package",
          "performer": "claude:pm_t04_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
            "digest": "sha256:85725c7b6501a2770cc17034ba253bfd3e40fc4f9d2d952e1f11e497fc1eb991",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,9d484121013c213d6d9f7808a23d157d8408dedd,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20380,335\n  github.com/gibbonmi/bench/internal/conformance,pass,70508,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74750,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,75040,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep2-root-conformance",
          "performer": "claude:pm_t04_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
            "digest": "sha256:b1beac0b2ecbd73e40ea0a58936d0cb8fd2707bca4895c73f6b2d14aec046ed0",
            "excerpt": "$ bench test --check root-conformance\ncheck root-conformance pass 15124ms tests_run=1 failures[0]\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep2-central-property",
          "performer": "claude:pm_t04_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
            "digest": "sha256:aa00380c3aa0b4b7564abe0de38d9569214a24de32d089bfe15848779d918dae",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,629\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,361,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:149: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (869 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the restore between two mutants, so the second mutant stacks on the first. TestMutantSetRunsEachMutantAlone must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
              "digest": "sha256:aa00380c3aa0b4b7564abe0de38d9569214a24de32d089bfe15848779d918dae",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,629\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,361,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:149: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (869 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pm-t04-rep3-probe-family",
          "performer": "claude:pm_t04_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
            "digest": "sha256:1add2e53de3001beac4bb60fdd61fb9f7888804e2fcfb11a77d7c90955303055",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,2d4436147f7d2332ace522eef0e366d2ee306355,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20345,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71003,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74214,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74596,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep3-testreport",
          "performer": "claude:pm_t04_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
            "digest": "sha256:1add2e53de3001beac4bb60fdd61fb9f7888804e2fcfb11a77d7c90955303055",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,2d4436147f7d2332ace522eef0e366d2ee306355,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20345,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71003,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74214,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74596,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "testreport",
          "command": "bench test --package ./internal/testreport",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep3-command-package",
          "performer": "claude:pm_t04_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
            "digest": "sha256:1add2e53de3001beac4bb60fdd61fb9f7888804e2fcfb11a77d7c90955303055",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,2d4436147f7d2332ace522eef0e366d2ee306355,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20345,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71003,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74214,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74596,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep3-conformance-package",
          "performer": "claude:pm_t04_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
            "digest": "sha256:1add2e53de3001beac4bb60fdd61fb9f7888804e2fcfb11a77d7c90955303055",
            "excerpt": "$ bench test --package ./internal/probe/... --package ./internal/testreport --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,2d4436147f7d2332ace522eef0e366d2ee306355,true\npackages[5]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20345,335\n  github.com/gibbonmi/bench/internal/conformance,pass,71003,1690\n  github.com/gibbonmi/bench/internal/probe,pass,74214,312\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/testreport,pass,74596,240\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep3-root-conformance",
          "performer": "claude:pm_t04_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
            "digest": "sha256:e75f913fa7b0e645da5b9881c163e5b6568cd3492fde7ce856e920d7f0ef3157",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,485fa5e295bec6471009169ba0ee9bfd140580dd,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15181,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t04-rep3-central-property",
          "performer": "claude:pm_t04_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
            "digest": "sha256:375a4cafd6d132a62e50f533a4b571077ebccb50d1f0981de624dd457a396233",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,629\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,370,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:149: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (869 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the restore between two mutants, so the second mutant stacks on the first. TestMutantSetRunsEachMutantAlone must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "485fa5e295bec6471009169ba0ee9bfd140580dd",
              "digest": "sha256:375a4cafd6d132a62e50f533a4b571077ebccb50d1f0981de624dd457a396233",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/probe/set.go,swap,failed,1,TestMutantSetRunsEachMutantAlone,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/probe/...,^TestMutantSetRunsEachMutantAlone$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,629\npackages[2]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/probe,fail,370,1\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/probe,TestMutantSetRunsEachMutantAlone,\"set_test.go:149: each replacement replaced [\\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\nfunc Clamp(n int) int {\\\\\\\\n\\\\\\\\tif n < 0 {\\\\\\\\n\\\\\\\\t\\\\\\\\treturn 0\\\\\\\\n\\\\\\\\t}\\\\\\\\n\\\\\\\\treturn n\\\\\\\\n}\\\\\\\\n\\\" \\\"package probefixture\\\\\\\\n\\\\\\\\n// Clamp keeps n at or above zero.\\\\\\\\… (869 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "pm-c2-standards-r1",
          "performer": "claude:pm_c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c2_standards",
            "digest": "sha256:1c7289194c9e1cf41495e1eb2a6b95579d76a0dd9dc9d84fa91dfaf5ca962fb2",
            "excerpt": "Standards axis claude:pm_c2_standards (opus/high), PM-C2 round 1, frozen pair 6584e93e..dfc5b339 (chunk tip b647725e), evidence sha256:40515554 current.\ncount: 3\nworst: S1\nS1 | auto-fix | 6 | AGENTS.md one source per fact: the single form's refusal order after the mutation input is spelled twice (set.go:111-139 vs probe.go:94-114: Prepare, gradeProseSubject, miss branch, gradeGateLock, lockSubject with deferred release, rereadSubject); the mutationInput-then-mutate pair is composed twice (set_input.go:77-80 vs probe.go:86-93). Needed: shared probe.go helpers for Prepare plus prose check, gate plus lock plus reread, and mutation resolution, called by both forms.\nS2 | auto-fix | 5 | baseline.go:14-16 and :50-58: baselineCells repeats the invalid/untouched pair that untouchedCells now owns. Build baselineCells from untouchedCells.\nS3 | auto-fix | 5 | lockProcessArgv (lock_holder_test.go:32) copies the helper re-entry argv of baseline_test.go:242; the ready-marker wait in lockHolderEnv (:45) is a third deadline poll loop after awaitFile (:151-159) and awaitStubStart (baseline_test.go:293-300). One reentryArgv helper and one poll helper.\nClean: one *refusal path for mutation titles; runRestored, recoverMutation, failRestore, keptCopy, mutatedCells single-sourced; setTables apart from resultTables; maxMutants, setSchema, resultSchema constants; InstallStarts and InstallSequence share install and answer; interruptedProbe and startLockHolder generalized; probeHelpSuffix derives from Usage; main.go 429 lines at base and tip; files under 400 lines; comments within bounds.\nAdvice: set.go:236-240 special-cases exit 2; set.go:147 and :227 restate the invalid exit; set_test.go:312 wrap width; \"probe subject in use\" literal twice; subjectFromPath comment; set helpers live in set_input_test.go.\nImplementation command: S1, the ticket's \"put the loop in set.go\" invites copying the admission sequence; name shared probe.go helpers in such a ticket.\n"
          },
          "axis": "Standards",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "pm-c2-spec-r1",
          "performer": "claude:pm_c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c2_spec",
            "digest": "sha256:594311ee6078cec119ac65afa937496ef7fa53931f15cb67fb61ab0dc58f903d",
            "excerpt": "Spec axis claude:pm_c2_spec (opus/high), PM-C2 round 1, frozen pair 6584e93e..dfc5b339 (chunk tip b647725e), evidence sha256:40515554 current.\ncount: 3\nworst: P1\nP1 | ask-user | 6 | set.go stopFrom and the baseline loop vs spec lines 195, 202, 213, 227: \"An unmatched entry marks its row invalid with the cause substring-miss\" collides with \"A failed baseline gives every row...\" and \"each later row reads set-stopped\". The code lets the stop or baseline cause overwrite the miss; the result depends on position ([m, miss, m] with a gate run during m1 gives row 2 substring-miss and row 3 gate-held; an interrupt gives row 2 set-stopped). The spec states no precedence.\nP2 | ask-user | 6 | set.go grade preserve branch vs spec line 218 (\"preserves, writes, runs, and restores that mutant as the single form does\"): the single form answers a preserve failure with \"error: probe preservation failed — <err>\" at exit 1 (subject.go:104-106, probe.go:183-186); the set prints a mutation-failed + untouched row and drops the error text. No set test covers this path.\nP3 | no-op | 7 | spec PM20 why-it-catches cell (line 424) says a skipped restore \"stacks the edits, and the hook refuses them\"; mutants come from the start bytes (line 187), so the replaced-bytes sequence is the real catch. Non-behavioral spec inaccuracy, flagged for veto.\nRows: PM1-PM26, PM30, PM32, PM33, PM35-PM41, PM43, PM44, PM51, PM67, PM69, PM71, PM72, PM77, PM78 all met; lines, cells, cause tokens, exits, refusal order, table headers, command-exit rule, report suffix, and trailing_bytes match.\nAuthor flags: (a) P1; (b) P2; (c) wording follows the <path>: <reason> shape, no finding; (d) PM20 met.\nScope: no creep.\nAdvice: an interrupted fresh build (PM71) gets a non-null focused_exit_code, so a mutant_report prints for a run whose Go child never started; spec line 277 does not decide this.\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
          "finding_ids": [
            "P1",
            "P2",
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "pm-c2-coverage-r1",
          "performer": "claude:pm_c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7246a8b8542d6fc74225bfb6f69a693d842b935b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c2_coverage",
            "digest": "sha256:fb5b95ccd0ecf95fd0fb227023cd15e2ce80868f19756941f050c26cd306d417",
            "excerpt": "Coverage axis claude:pm_c2_coverage (opus/high), PM-C2 round 1, frozen pair 6584e93e..dfc5b339 (chunk tip b647725e), evidence sha256:40515554 current.\ncount: 5\nworst: C1\nC1 | auto-fix | 8 | set.go:179-181: a published mutant write failure whose restore also fails should read restore-failed/no, a preserved table, and exit 2. Probe omit \"s.rows[i].exit = s.rows[i].cells.failRestore()\" (occurrence 1): silent over 47 ^TestMutantSet tests, restored yes. Needed: a TestMutantSetStopsOnAMutationFailure case with mutation fault directory-sync and restore fault rename.\nC2 | auto-fix | 7 | set.go:205: a signal after a run completes (during its restore or between runs) should stop the next row with set-stopped/untouched (spec 212-214). Probes that disabled the s.ctx.Err() check and the OutcomeInterrupted branch (set.go:195) were each silent. Needed: an in-process row that raises SIGINT during mutant 1's restore and asserts row 2 set-stopped/untouched and 2 Go starts.\nC3 | auto-fix | 7 | set.go:254: trailing_bytes must count every byte after the execution row (spec 278); only report suffixes and all-miss 0 are asserted. Probe e.trailing = len(reports): silent. Needed: assert setExecution with the byte count in TestMutantSetBaselineEndsTheSet and TestMutantSetStopsOnARestoreFailure.\nC4 | auto-fix | 5 | set.go:149: on a failed baseline over an owned binary the builds cell should read 1; probe omitting the assignment: silent. Needed: an owned case of the failed-baseline test.\nC5 | ask-user | 4 | set.go:168-173: a preserve failure reads mutation-failed/untouched and drops \"probe preservation failed — <err>\" (the single form refuses with that line, refusal_test.go:211). Probe swapping the cause to set-stopped: silent. Decision needed, then a row.\nBit (restored yes): swapped expectFitsSelection arguments (PM72 red); cap compared with >= (PM11 red); last-row command exit without the 2-rank (PM22 red).\nAdvice: a write-failure row drops err.Error() unless the restore also fails; the row.exit == 2 break adds nothing to max.\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "b647725ed4f1a49e5635ed7e9bb9ea3fe8df118f",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5"
          ],
          "supersedes": []
        },
        {
          "id": "pm-c2-standards-r2",
          "performer": "claude:pm_c2_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c2_standards_r2",
            "digest": "sha256:83ef13515cb3b3e443584125167a9d25886c7a7b593de0b47fd00bd51414ec31",
            "excerpt": "Standards axis claude:pm_c2_standards_r2 (opus/high), PM-C2 confirming round 2, chunk pair 6584e93e..28d4a4e4, repair delta b647725e..4b9543cf -- internal/, evidence sha256:82d4d742 current.\ncount: 1\nworst: S4\nS1 closed: resolveMutation (probe.go:112), prepareSelection (:123), holdSubject (:136) serve runWith, runSet, and resolveEntry; no step sequenced twice; holdSubject's own release on a refused reread is sound; probe.go 384 lines.\nS2 closed: baselineCells builds on untouchedCells (baseline.go:50-54).\nS3 closed: reentryArgv (baseline_test.go:375) and pollUntil serve every site; settleDelay replaces the 50 ms literal.\nRemoved OutcomeInterrupted set branch: sound; every interrupted producer requires the held context done (held_runs.go:34, command.go:283, :298), and the next matching entry's stopCause gives set-stopped.\nS4 | auto-fix | 6 | AGENTS.md one source per fact: sleepingRun (baseline_test.go:289) is still spelled out at baseline_test.go:174 and lock_test.go:114. Replace both with sleepingRun.\nAdvice: holdRecordLock keeps its own 10 ms poll interval; requireTrailingBytes check order; requirePreserved repeats older preserved-row assertions.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "4b9543cf396f0709957f44e812eafcdb5b37775a",
          "finding_ids": [
            "S4"
          ],
          "supersedes": [
            "pm-c2-standards-r1"
          ]
        },
        {
          "id": "pm-c2-spec-r2",
          "performer": "claude:pm_c2_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c2_spec_r2",
            "digest": "sha256:e53c0914c74a020abd9dcb0af36ea47e7ca1435404983fca22f41db4e800befd",
            "excerpt": "Spec axis claude:pm_c2_spec_r2 (opus/high), PM-C2 confirming round 2, chunk pair 6584e93e..28d4a4e4, delta dfc5b339..28d4a4e4, evidence sha256:82d4d742 current. The axis named its findings P1 and P2; the record uses P21 and P22.\ncount: 2\nworst: P21\nP21 | auto-fix | 6 | spec.md:232 \"A set answer prints no failure text, except the reason that the preserved table gives for a kept copy\" says too much: baseline reports and interrupted mutant reports print error lines (spec 206, 283; testreport/command.go:287; set.go:255-256). Code correct. Narrow the sentence to a failed preservation, a failed write, and a stop.\nP22 | no-op | 6 | PM24 why cell (spec.md:442) \"A loop that ignores the interrupted kind starts mutant 2's run\" names a mutant that is now the green code; the context check (set.go:193) catches the interrupt. Reword to the skipped context check.\nPlan expansion 06b238e6 correct and in scope: D3, the Sequence sentences, PM81-PM86, ticket 4 lines, chunk row, PM20 cell (P3 closed); seam cells name the real tests.\nP1, P2 of round 1 closed: PM81-PM86 audited against their tests.\nThe removed OutcomeInterrupted branch is sound: every interrupted outcome comes from the held context, so the next matching entry stops with set-stopped.\nrep1 entries complete at 4b9543cf, source digest b9185855.\nAdvice: D3 sits under Reviewer decisions as a consultant decision (vetoable); spec lines 231 and 233 both restate the later-row rule.\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "4b9543cf396f0709957f44e812eafcdb5b37775a",
          "finding_ids": [
            "P21",
            "P22"
          ],
          "supersedes": [
            "pm-c2-spec-r1"
          ]
        },
        {
          "id": "pm-c2-coverage-r2",
          "performer": "claude:pm_c2_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "b91858554d677c063d660c207bafdb298326a7c0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c2_coverage_r2",
            "digest": "sha256:80ce8cefd6a6460875a4557ae2df7458a58db105d90d14b26df4caa3dde22a57",
            "excerpt": "Coverage axis claude:pm_c2_coverage_r2 (opus/high), PM-C2 confirming round 2, chunk pair 6584e93e..28d4a4e4, repair delta b647725e..4b9543cf -- internal/, evidence sha256:82d4d742 current.\ncount: 0\nworst: none\nClosed: PM83/C1 (set_test.go:631, restore-failed, no, exit 2, preserved table); PM84/C2 (baseline_test.go:278); PM85/C3 (requireTrailingBytes in two tests, derived count); PM86/C4 (owned case, builds 1, set_test.go:573-581); PM82/C5 (whole-answer compare, set_test.go:657); PM81 outcome asserted on both paths.\nOwn probes (restored yes): omit s.preserved = keptCopy(...) (set.go:173): bit; swap preserveFailedCause for mutationFailedCause (set.go:162): bit; index the miss guard by j-i (set.go:210): silent over 55 ^TestMutantSet tests.\n(a) No test pins a miss row after a non-zero stop index; graded covered, because every stopFrom caller goes through the one guard and removing the guard bit in /baseline.\n(b) No interrupted outcome is possible without a cancelled set context (command.go:287, :301, held_runs.go:35 all require the held ctx; signal.NotifyContext over Background has no other canceller).\n(c) PM84's 50 ms settleDelay has a low false-red flake risk on a starved host; no deterministic wait exists without a production seam.\nAdvice: a [match, miss, match] subtest that fails mutant 1's restore would pin the guard past index 0; a test hook for the set cancel context would make PM84 deterministic (production seam, reviewer call).\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "4b9543cf396f0709957f44e812eafcdb5b37775a",
          "finding_ids": [],
          "supersedes": [
            "pm-c2-coverage-r1"
          ]
        },
        {
          "id": "pm-c2-standards-r3",
          "performer": "claude:pm_c2_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_c2_standards_r2",
            "digest": "sha256:e1d8f5cbf74dd337c30b24b84c9364e421cbe8849d225b60720371377d64f340",
            "excerpt": "Standards axis claude:pm_c2_standards_r2 (opus/high), PM-C2 narrow confirming round 3, delta 9d484121..ae0f9b8b -- internal/.\ncount: 1\nworst: S5\nS4 closed: sleepingRun replaces the local at baseline_test.go:174 and the literal at lock_test.go:114; rg for \"sleep 60\" finds only the definition (baseline_test.go:297) and an unrelated build script.\nS5 | auto-fix | 7 | AGENTS.md one source per fact: the restore-failed row {1, \"restore-failed\", \"swap\", \"failed\", 1, \"TestClampNegative\", 1, 1, 2, nil, \"no\"} is authored by hand at baseline_test.go:256 (new), set_test.go:177, and set_test.go:309; the file already collapses row shapes into biteRow, missRow, untouchedRow, redRow (set_input_test.go:85-103). Add restoreFailedRow(n) beside them and use it at all three sites.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
          "finding_ids": [
            "S5"
          ],
          "supersedes": [
            "pm-c2-standards-r2"
          ]
        },
        {
          "id": "pm-c2-spec-r3",
          "performer": "claude:pm_c2_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c2_spec_r2",
            "digest": "sha256:c75bc2062904d867d5cf7baf8cae78679aad1dddb891e0d1e83fa090f568f48d",
            "excerpt": "Spec axis claude:pm_c2_spec_r2 (opus/high), PM-C2 narrow confirming round 3, delta 28d4a4e4..ccd8c12c -- specs/ internal/ (plan 58e6eca9, small fix ae0f9b8b).\ncount: 0\nworst: none\nP21 closed: spec.md:232 now reads \"A failed preservation, a failed write, and a stop print no error line, and a kept copy's reason prints only in the preserved table.\" No conflict with lines 206 and 283; the code matches (set.go:156-176, keptCopy).\nP22 closed: the PM24 why cell names the context check (stopCause, set.go:193).\nThe restore-failure subtest of TestMutantSetMissKeepsItsCause (baseline_test.go:250-258) meets PM81 and adds no behavior; the sleepingRun replacements are neutral refactors.\nAdvice: the opus/low small-fix assignment follows the reviewer direction for small-fix delegates.\n"
          },
          "axis": "Spec",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
          "finding_ids": [],
          "supersedes": [
            "pm-c2-spec-r2"
          ]
        },
        {
          "id": "pm-c2-coverage-r3",
          "performer": "claude:pm_c2_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "5e5f4291a5286ba17939b513998d18a1483e3621",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c2_coverage_r2",
            "digest": "sha256:ddd5856716bc4c0b694f8fe50c8e6c9e97d9d514818cf1068a861e350d626caa",
            "excerpt": "Coverage axis claude:pm_c2_coverage_r2 (opus/high), PM-C2 narrow confirming round 3, delta 9d484121..ae0f9b8b -- internal/, head ccd8c12c clean.\ncount: 0\nworst: none\nQuestion (a) closed: the restore-failure subtest of TestMutantSetMissKeepsItsCause (baseline_test.go:250-258) calls stopFrom(1) over the miss row and asserts restore-failed, missRow(2), set-stopped, exit 2; the j-i guard swap bit on /restore_failure.\nsleepingRun replacements keep the tests red-capable. Own probe: baseline_test.go:297 sleepingRun to probetest.Start{Events: cannedPass}, run TestProbeRestoresOnInterrupt, TestSubjectLockReleasesOnInterrupt, TestMutantSetStopsOnInterrupt: bit, 4 failed, restored yes.\nNo command change.\n"
          },
          "axis": "Coverage",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "ae0f9b8b96384ce2136f4b437c2004bef6296a39",
          "finding_ids": [],
          "supersedes": [
            "pm-c2-coverage-r2"
          ]
        },
        {
          "id": "pm-c2-standards-r3b",
          "performer": "claude:pm_c2_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c2_standards_r2",
            "digest": "sha256:bedd6b44399df4ccbf747d2a7ac3df04e836b64b06603b94357ca461f59cc9ac",
            "excerpt": "Standards axis claude:pm_c2_standards_r2 (opus/high), PM-C2 reaffirmation of S5 at the cycle 3 tip, delta ae0f9b8b..485fa5e2 -- internal/.\ncount: 0\nworst: none\nS5 closed, confidence 9: restoreFailedRow(n) (set_input_test.go:105-108) is the one source of the restore-failed row; baseline_test.go:256, set_test.go:177, and set_test.go:317 call it; rg for \"restore-failed\", \"swap\" finds only the helper. No new defect.\n"
          },
          "axis": "Standards",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "485fa5e295bec6471009169ba0ee9bfd140580dd",
          "finding_ids": [],
          "supersedes": [
            "pm-c2-standards-r3"
          ]
        },
        {
          "id": "pm-c2-spec-r3b",
          "performer": "claude:pm_c2_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c2_spec_r2",
            "digest": "sha256:a78c922c0b29988657403d45352b231b88fced43dae877f911e466208f5ed59a",
            "excerpt": "Spec axis claude:pm_c2_spec_r2 (opus/high), PM-C2 reaffirmation at the cycle 3 tip, delta ae0f9b8b..485fa5e2 -- internal/.\ncount: 0\nworst: none\nRound 3 result holds, confidence 9: set_input_test.go:105-108 adds restoreFailedRow(n) returning the row that baseline_test.go:256 (n = 1), set_test.go:177 (n = 1), and set_test.go:317 (n = 2) spelled inline; no expectation changed; test files only; P21 and P22 stay closed and PM81 stays met.\n"
          },
          "axis": "Spec",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "485fa5e295bec6471009169ba0ee9bfd140580dd",
          "finding_ids": [],
          "supersedes": [
            "pm-c2-spec-r3"
          ]
        },
        {
          "id": "pm-c2-coverage-r3b",
          "performer": "claude:pm_c2_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "edfc7908b69827412682c80dba49c276d8e5c3cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_c2_coverage_r2",
            "digest": "sha256:36325dae601ed480b17adaf028a94e7d139a1295a99cc251e20c73651fde1ee3",
            "excerpt": "Coverage axis claude:pm_c2_coverage_r2 (opus/high), PM-C2 reaffirmation at the cycle 3 tip, delta ae0f9b8b..485fa5e2 -- internal/.\ncount: 0\nworst: none\nRound 3 result holds, confidence 9: restoreFailedRow(n) (set_input_test.go:106-108) matches the three replaced rows cell for cell with the same entry numbers (1, 1, 2); assertions, faults, and exits are unchanged, so the earlier bites still apply. No probe needed.\n"
          },
          "axis": "Coverage",
          "base": "6584e93e6df7826a38d93d05bcfe0094afe82971",
          "tip": "485fa5e295bec6471009169ba0ee9bfd140580dd",
          "finding_ids": [],
          "supersedes": [
            "pm-c2-coverage-r3"
          ]
        }
      ]
    },
    {
      "id": "PM-D",
      "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
      "tip": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
      "plan_digest": "sha256:2d2a1f9466ffc76362c39594e0c13eddc9994c0fa59884d007540acae89abbba",
      "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
      "acceptance_rows": [
        "PM34",
        "PM42",
        "PM52",
        "PM53",
        "PM54",
        "PM55",
        "PM56",
        "PM57",
        "PM58",
        "PM59",
        "PM60",
        "PM74",
        "PM87"
      ],
      "verification": [
        {
          "id": "pm-t05-probe-family",
          "performer": "claude:pm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6b447a69eabfae707aa27b72f08034495070386e",
            "digest": "sha256:94b7490f5527ad3b26f849a3dc1fda52ccd3b35dbb11fc7a9d3cb1f29eb8cda8",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,6b447a69eabfae707aa27b72f08034495070386e,false\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17856,335\n  github.com/gibbonmi/bench/internal/conformance,pass,63081,1690\n  github.com/gibbonmi/bench/internal/probe,pass,67606,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3066,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,11566,191\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(exit 0; wall 73.8 s)\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t05-record-forms",
          "performer": "claude:pm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6b447a69eabfae707aa27b72f08034495070386e",
            "digest": "sha256:94b7490f5527ad3b26f849a3dc1fda52ccd3b35dbb11fc7a9d3cb1f29eb8cda8",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,6b447a69eabfae707aa27b72f08034495070386e,false\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17856,335\n  github.com/gibbonmi/bench/internal/conformance,pass,63081,1690\n  github.com/gibbonmi/bench/internal/probe,pass,67606,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3066,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,11566,191\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(exit 0; wall 73.8 s)\n"
          },
          "requirement": "record-forms",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "pm-t05-command-package",
          "performer": "claude:pm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6b447a69eabfae707aa27b72f08034495070386e",
            "digest": "sha256:94b7490f5527ad3b26f849a3dc1fda52ccd3b35dbb11fc7a9d3cb1f29eb8cda8",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,6b447a69eabfae707aa27b72f08034495070386e,false\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17856,335\n  github.com/gibbonmi/bench/internal/conformance,pass,63081,1690\n  github.com/gibbonmi/bench/internal/probe,pass,67606,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3066,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,11566,191\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(exit 0; wall 73.8 s)\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t05-conformance-package",
          "performer": "claude:pm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6b447a69eabfae707aa27b72f08034495070386e",
            "digest": "sha256:94b7490f5527ad3b26f849a3dc1fda52ccd3b35dbb11fc7a9d3cb1f29eb8cda8",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,6b447a69eabfae707aa27b72f08034495070386e,false\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17856,335\n  github.com/gibbonmi/bench/internal/conformance,pass,63081,1690\n  github.com/gibbonmi/bench/internal/probe,pass,67606,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3066,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,11566,191\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n(exit 0; wall 73.8 s)\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t05-root-conformance",
          "performer": "claude:pm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6b447a69eabfae707aa27b72f08034495070386e",
            "digest": "sha256:c1624ea32bb34f36986fd73170c370cc8833935592e000c3158fdc52bc5cf5a0",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,6b447a69eabfae707aa27b72f08034495070386e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16253,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n(exit 0; wall about 23.854166800 s including date)\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t05-central-property",
          "performer": "claude:pm_t05",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "6b447a69eabfae707aa27b72f08034495070386e",
            "digest": "sha256:ea686e2890afdb62daab4bc294533b3ca1f561cc58db6f651f8f13c564252654",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/recordcmd/verification.go,swap,failed,1,TestVerificationChecksTheRowExit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/recordcmd,^TestVerificationChecksTheRowExit$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,554\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,52,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationChecksTheRowExit,\"verification_import_test.go:171: verification = exit 0, output \\\"verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}:\\\\\\\\n  chunk,\\\\\\\\\\\"1\\\\\\\\\\\",tests-1,tests,author-verification,fail,bdf67c11e413902be858a01300e96d4553d6… (411 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Check --exit-code against the command exit instead of the selected row's exit. TestVerificationChecksTheRowExit must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "6b447a69eabfae707aa27b72f08034495070386e",
              "digest": "sha256:ea686e2890afdb62daab4bc294533b3ca1f561cc58db6f651f8f13c564252654",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/recordcmd/verification.go,swap,failed,1,TestVerificationChecksTheRowExit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/recordcmd,^TestVerificationChecksTheRowExit$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,554\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,52,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationChecksTheRowExit,\"verification_import_test.go:171: verification = exit 0, output \\\"verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}:\\\\\\\\n  chunk,\\\\\\\\\\\"1\\\\\\\\\\\",tests-1,tests,author-verification,fail,bdf67c11e413902be858a01300e96d4553d6… (411 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "pm-t05-rep1-probe-family",
          "performer": "claude:pm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
            "digest": "sha256:90ad6c5489ad315c5051b3dc176980f70b68db7a2c45213c0ffc3391e3b0214c",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cddaffed69db8b8fb6a0d2bb58c8dfa3fe19b828,true\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18176,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61826,1690\n  github.com/gibbonmi/bench/internal/probe,pass,65916,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3051,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,15394,192\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "probe-family",
          "command": "bench test --package ./internal/probe/...",
          "exit_code": 0
        },
        {
          "id": "pm-t05-rep1-record-forms",
          "performer": "claude:pm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
            "digest": "sha256:90ad6c5489ad315c5051b3dc176980f70b68db7a2c45213c0ffc3391e3b0214c",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cddaffed69db8b8fb6a0d2bb58c8dfa3fe19b828,true\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18176,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61826,1690\n  github.com/gibbonmi/bench/internal/probe,pass,65916,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3051,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,15394,192\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "record-forms",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "pm-t05-rep1-command-package",
          "performer": "claude:pm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
            "digest": "sha256:90ad6c5489ad315c5051b3dc176980f70b68db7a2c45213c0ffc3391e3b0214c",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cddaffed69db8b8fb6a0d2bb58c8dfa3fe19b828,true\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18176,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61826,1690\n  github.com/gibbonmi/bench/internal/probe,pass,65916,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3051,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,15394,192\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "pm-t05-rep1-conformance-package",
          "performer": "claude:pm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
            "digest": "sha256:90ad6c5489ad315c5051b3dc176980f70b68db7a2c45213c0ffc3391e3b0214c",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --package ./internal/probe/... --package ./internal/reviewrecord/... --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cddaffed69db8b8fb6a0d2bb58c8dfa3fe19b828,true\npackages[7]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18176,335\n  github.com/gibbonmi/bench/internal/conformance,pass,61826,1690\n  github.com/gibbonmi/bench/internal/probe,pass,65916,315\n  github.com/gibbonmi/bench/internal/probe/probetest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,3051,99\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,15394,192\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "conformance-package",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t05-rep1-root-conformance",
          "performer": "claude:pm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
            "digest": "sha256:6e30f91f0de8131c561384305ea00be17fb5bef10b4f35ec5455fa230667fcd1",
            "excerpt": "$ bench worktree exec probe-mutant-sets-build -- bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  probe-mutant-sets-build,cddaffed69db8b8fb6a0d2bb58c8dfa3fe19b828,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15777,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "pm-t05-rep1-central-property",
          "performer": "claude:pm_t05_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
            "digest": "sha256:1335aca59b1b80a21fe4817a0a70b934710117a4fc513d4c3217726fb865d32b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/recordcmd/verification.go,swap,failed,3,TestVerificationChecksTheRowExit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/recordcmd,all,passed,192\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1144\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,11001,192\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationChecksTheRowExit,\"verification_import_test.go:216: verification = exit 0, output \\\"verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}:\\\\\\\\n  chunk,\\\\\\\\\\\"1\\\\\\\\\\\",tests-1,tests,author-verification,fail,bdf67c11e413902be858a01300e96d4553d6… (411 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationImportsARestoreFailedRow,\"verification_import_test.go:193: verification = exit 1, output \\\"error: bench record verification refused — --exit-code 2 differs from the imported mutants.exit_code 1\\\\\\\\n\\\"; want exit 0\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationImportsOneMutantRow,\"verification_import_test.go:182: verification = exit 1, output \\\"error: bench record verification refused — --exit-code 0 differs from the imported mutants.exit_code 1\\\\\\\\n\\\"; want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "central-property",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0,
          "probe": {
            "mutation": "Check --exit-code against the selected row's focused exit instead of its row exit. TestVerificationChecksTheRowExit must turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
              "digest": "sha256:1335aca59b1b80a21fe4817a0a70b934710117a4fc513d4c3217726fb865d32b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/reviewrecord/recordcmd/verification.go,swap,failed,3,TestVerificationChecksTheRowExit,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord/recordcmd,all,passed,192\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1144\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,fail,11001,192\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationChecksTheRowExit,\"verification_import_test.go:216: verification = exit 0, output \\\"verification[1]{list,chunk,id,requirement,role,outcome,source_digest,excerpt_digest}:\\\\\\\\n  chunk,\\\\\\\\\\\"1\\\\\\\\\\\",tests-1,tests,author-verification,fail,bdf67c11e413902be858a01300e96d4553d6… (411 bytes)\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationImportsARestoreFailedRow,\"verification_import_test.go:193: verification = exit 1, output \\\"error: bench record verification refused — --exit-code 2 differs from the imported mutants.exit_code 1\\\\\\\\n\\\"; want exit 0\",1\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,TestVerificationImportsOneMutantRow,\"verification_import_test.go:182: verification = exit 1, output \\\"error: bench record verification refused — --exit-code 0 differs from the imported mutants.exit_code 1\\\\\\\\n\\\"; want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "pm-d-standards-r1",
          "performer": "claude:pm_d_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_d_standards",
            "digest": "sha256:22f0652391ea483aae7430d2883a3ad49d11c00dcd9e3d6b4f8901cb24aad474",
            "excerpt": "Standards axis claude:pm_d_standards (opus/high), PM-D round 1, frozen pair 1870269d..520f4899 (chunk tip 6b447a69), evidence sha256:d4c5162f current.\ncount: 1\nworst: S1\nS1 | ask-user | 7 | result.go:42 CommandExit in MutantResult, set at result.go:105 (ReadMutant), ignored by RenderMutantSet (set.go:109-111), read only by result_test.go mutantResult (~191). Speculative Generality and spec.md:303 / ticket 5 line 13 (\"MutantResult carries Verdict, Restored, RowExit, and FocusedExit, so no field named for the command exit carries a row exit\"). The only production .CommandExit reader is verification.go:115 on probe.Result. Needed: remove the field, the command parameter of mutantResult, and the \"CommandExit is ignored\" sentence; re-record the central probe on a mutation that needs no new surface.\nClean: ReadResult and ReadMutant share readForm, splitEnvelope, decodeTables, requireSchema; FieldSource and SetFieldSource share fieldSource; importRefusal owns the three lines; no --mutant wording in internal/probe; no second header-count parser.\nAdvice: RenderMutantSet is approved by spec.md:389 and reuses setRun.answer; probe.ResultSchema has a production reader; the producer trailing count is derived at probe.go:332 and set.go:278 (pre-existing); set.go:116 repeats \"no\"; result.go:204 wording \"want one mutants row\"; result_test.go setAnswer drops its exit code; TestReadMutantSharesTheReaderRules recomputes the trailing bytes; mutantSet only wraps RenderMutantSet; missedRow has an empty cause.\nImplementation command: charge \"Add no production surface that only a test or a probe reads; if a named mutation needs one, report it.\"\n"
          },
          "axis": "Standards",
          "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
          "tip": "6b447a69eabfae707aa27b72f08034495070386e",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "pm-d-spec-r1",
          "performer": "claude:pm_d_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_d_spec",
            "digest": "sha256:e0c1bbe0f6f9796a74a61fd9765f076830fd4486f5a0b76b39b011f49474ec1b",
            "excerpt": "Spec axis claude:pm_d_spec (opus/high), PM-D round 1, frozen pair 1870269d..520f4899 (chunk tip 6b447a69), evidence sha256:d4c5162f current.\ncount: 2\nworst: P1\nP1 | ask-user | 7 | spec.md:303 \"MutantResult carries Verdict, Restored, RowExit, and FocusedExit, so no field named for the command exit carries a row exit\"; result.go:41-42 adds CommandExit, which no production code reads (only result_test.go:191). Options: accept as a flagged deviation, or remove it and re-author the central probe on the declared API (for example bind the row import's exit to FocusedExit).\nP2 | ask-user | 6 | spec.md:477 PM59 \"--mutant 1 of a row that reads substring-miss refuses with the null focused-exit line\"; the test's missedRow renders through RenderMutantSet with cause \"\", restored no, mutation \"\", not substring-miss/untouched. The null-to-0 mutation is still caught. Options: record the renderer decision and amend PM59, or build PM59's answer through a probe run (spec.md:390).\nRows: PM34, PM42, PM52-PM58, PM60, PM74 met; PM59 partial. All four bindings match spec.md:307-309. Exact refusal lines match.\nFlags: header row count for both forms is not a finding; PM58's coordinator red is valid by trace.\nScope creep: none beyond P1.\nAdvice: readImported uses mutant 0 for \"no --mutant\"; an explicit presence bool would decouple it from verificationValid.\nImplementation command: phrase central-property probes over the spec's declared API; a renderer-built row should carry the cells its acceptance row names.\n"
          },
          "axis": "Spec",
          "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
          "tip": "6b447a69eabfae707aa27b72f08034495070386e",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "pm-d-coverage-r1",
          "performer": "claude:pm_d_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "77e3cd4b4bb3e9ac327168404ef131ad7a0fb3ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:pm_d_coverage",
            "digest": "sha256:34b51dbaf4b3e994008f07da36a2e03a7cde5dc59325c068f9bf77cfcac5d146",
            "excerpt": "Coverage axis claude:pm_d_coverage (opus/high), PM-D round 1, frozen pair 1870269d..520f4899 (chunk tip 6b447a69), evidence sha256:d4c5162f current.\ncount: 1\nworst: C1\nC1 | auto-fix | 8 | recordcmd verification.go readImported mutant branch (restored: result.Restored) vs spec.md:308 \"It binds ... the restore to its restored cell\": an import of a restore-failed row is unpinned. Probe restored: true over ./internal/reviewrecord/... and ./cmd/bench (625 tests): silent, restored yes. Needed: a recordcmd case importing {Verdict: restore-failed, Restored: false, RowExit: 2, FocusedExit: 1} with --exit-code 2, asserting outcome restore-failed and restore fail.\nProbes (restored yes): ReadMutant Restored: true: bit; ReadMutant rows[len-n]: silent in probe, bit in recordcmd (3 tests); readForm trailing > check: bit.\nInput family covered: bit, silent, invalid null focus, restore-failed on the reader side; cut and extended answers; spill frames; quoted cell; repeated table; schema substitution; single form with and without --mutant; --mutant 0, -1, x, out of range, alone; command-exit vs row-exit. Header row-count variants all refuse. Each sentinel line pinned. Hostile-input walk clean.\nAdvice: setAnswer rows 1 and 3 project to the same MutantResult; no ReadMutant test reads an answer with a preserved table; integer() accepts +1 and 01 (pre-existing).\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
          "tip": "6b447a69eabfae707aa27b72f08034495070386e",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "pm-d-standards-r2",
          "performer": "claude:pm_d_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_d_standards_r2",
            "digest": "sha256:4200d819c51b223099a574ffee35f3879da4071de8b7ef8ee68436bbe3fd847d",
            "excerpt": "Standards axis claude:pm_d_standards_r2 (opus/high), PM-D confirming round 2, chunk pair 1870269d..d4b2a7f1, repair delta 6b447a69..fc4c772e -- internal/ cmd/, evidence sha256:5140b0c8 current.\ncount: 0\nworst: none\nS1 closed: MutantResult has exactly Verdict, Restored, RowExit, FocusedExit (result.go:32-41); ReadMutant fills only those (:103); RenderMutantSet has no remaining reference; the remaining CommandExit is the single form's, read at verification.go:115.\nNo new defect: the probetest clamp harness is one source for the journey and recordcmd (no import cycle); probe_test.go newFixture is a second clamp harness that predates the base and stays outside ticket 5's fence; the restore-failure hook targets the preserved copy; the kept literals are CLI inputs or the PM87 oracle; comments match the code.\nAdvice: collapse newFixture onto probetest after this spec; the hook could use filepath.Base(probetest.ClampPath); ExitCode 1 could read the focused exit.\nImplementation command: no command change.\n"
          },
          "axis": "Standards",
          "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
          "tip": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
          "finding_ids": [],
          "supersedes": [
            "pm-d-standards-r1"
          ]
        },
        {
          "id": "pm-d-spec-r2",
          "performer": "claude:pm_d_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_d_spec_r2",
            "digest": "sha256:1949599aa0634a22171f744f45f308ad1387b35712f76ef3961d57a48f02bc09",
            "excerpt": "Spec axis claude:pm_d_spec_r2 (opus/high), PM-D confirming round 2, chunk pair 1870269d..d4b2a7f1, delta 520f4899..d4b2a7f1, evidence sha256:5140b0c8 current.\ncount: 0\nworst: none\nPlan commit 579ceadf correct and in scope: the central probe mutates the declared RowExit and FocusedExit and bites TestVerificationChecksTheRowExit; PM87 (story 48) pins the fail side of spec.md:308; ticket 5 lines and Covers match; the chunk row adds PM87.\nP1 closed: MutantResult has exactly Verdict, Restored, RowExit, FocusedExit (result.go:32-41); RenderMutantSet removed.\nP2 closed: the missed row comes from a real probe run and asserts substring-miss and untouched (verification_import_test.go:234-236), then the null focused-exit refusal (:243).\nPM87 met by TestVerificationImportsARestoreFailedRow (:187-194); PM34, PM42, PM52-PM60, PM74 stay met.\npm-t05-rep1 entries complete for six requirements at fc4c772e, source digest d88eb672.\nAdvice: spec.md:389 and ticket 5 line 22 still allow an exported renderer; the TestVerificationChecksTheRowExit comment names only the command exit; four entries share one combined excerpt.\nImplementation command: no command change.\n"
          },
          "axis": "Spec",
          "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
          "tip": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
          "finding_ids": [],
          "supersedes": [
            "pm-d-spec-r1"
          ]
        },
        {
          "id": "pm-d-coverage-r2",
          "performer": "claude:pm_d_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d88eb672abcf7273668e74f584b006b59838673f",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:pm_d_coverage_r2",
            "digest": "sha256:b84052c46a749012b441c7c6094e98c446fdd713b618dc6494e230a64ec4b8d6",
            "excerpt": "Coverage axis claude:pm_d_coverage_r2 (opus/high), PM-D confirming round 2, chunk pair 1870269d..d4b2a7f1, repair delta 6b447a69..fc4c772e -- internal/ cmd/, evidence sha256:5140b0c8 current.\ncount: 0\nworst: none\nC1/PM87 closed: TestVerificationImportsARestoreFailedRow (verification_import_test.go:187-194) builds its answer from a real probe run and requires restore-failed, exit 1, restore fail.\nOwn probes, all bit with restored yes: (1) PM59 null mapped to 0 at verification.go:154: TestVerificationMutantRefusals/null_focused_exit; (2) ReadMutant FocusedExit to &row.exit (result.go:103): 3 failed; (3) ReadMutant row index to e.rows[0] (result.go:102): 2 failed; (4) preserved copy renamed (subject.go:96): ImportsARestoreFailedRow failed; (5) bypass replacer.Replace (subject.go:130): TestNativeRecordCLIJourney failed; (6) single-form focused bound to CommandExit (verification.go:115): journey failed.\nNo new defect: RenderMutantSet, MutantResult.CommandExit, the journey constants, and the old row fixtures have no remaining reference.\nAdvice: a precondition that row 1 reads restore-failed and no before the PM87 import would make a red plain.\nImplementation command: no command change.\n"
          },
          "axis": "Coverage",
          "base": "1870269dc00a43134d20f02af0a7afef77f08eec",
          "tip": "fc4c772eab1d55cbe6fc8e46faf6cccc6709cd1c",
          "finding_ids": [],
          "supersedes": [
            "pm-d-coverage-r1"
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
  },
  "amendments": [
    {
      "from": "sha256:bb382785f355201ccac68fe63e38675f5033c694a0f62850ef28f7860ba80ee9",
      "to": "sha256:faebeaa6ca8f219e2c4e6b61652d1e446266fd81a33e37d53f9b07e8d56f925e",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ]
      }
    },
    {
      "from": "sha256:faebeaa6ca8f219e2c4e6b61652d1e446266fd81a33e37d53f9b07e8d56f925e",
      "to": "sha256:b0608123d05d27cb3a1d51026e87515bdeb95b1cae077671a1e62d28fa7e32db",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ]
      }
    },
    {
      "from": "sha256:b0608123d05d27cb3a1d51026e87515bdeb95b1cae077671a1e62d28fa7e32db",
      "to": "sha256:9c3e83f7d2d4d29468af0b47d2c90b7b85daf5523c6c203f77cf0286be49944b",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ]
      }
    },
    {
      "from": "sha256:9c3e83f7d2d4d29468af0b47d2c90b7b85daf5523c6c203f77cf0286be49944b",
      "to": "sha256:8875b5f170b67bff1a7399af0ff34be749c7aadeb398e94fe08406189fdc3786",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ]
      }
    },
    {
      "from": "sha256:8875b5f170b67bff1a7399af0ff34be749c7aadeb398e94fe08406189fdc3786",
      "to": "sha256:7770b3b86644c8d0349158481463350d60a3662bd9064a30fc86f0469d1a9c86",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ]
      }
    },
    {
      "from": "sha256:7770b3b86644c8d0349158481463350d60a3662bd9064a30fc86f0469d1a9c86",
      "to": "sha256:da6c45a4a51006e8483f681cd55a50dbb09d4c8fff29378f6724059c54511d88",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ]
      }
    },
    {
      "from": "sha256:da6c45a4a51006e8483f681cd55a50dbb09d4c8fff29378f6724059c54511d88",
      "to": "sha256:04122c8b7b657e7f8438574e1824d8b05956b51a79ab10aa96f0e68e80ca75cf",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ]
      }
    },
    {
      "from": "sha256:04122c8b7b657e7f8438574e1824d8b05956b51a79ab10aa96f0e68e80ca75cf",
      "to": "sha256:2c1499fc0954d59fe2ec2a1a1729145a83b273eb3c017753746c4434314f505c",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ]
      }
    },
    {
      "from": "sha256:2c1499fc0954d59fe2ec2a1a1729145a83b273eb3c017753746c4434314f505c",
      "to": "sha256:3a2a2ee9dc4559305f9a21adbc4ed12b5e01c29c4312d42a743cd4df9ff8d1de",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ]
      }
    },
    {
      "from": "sha256:3a2a2ee9dc4559305f9a21adbc4ed12b5e01c29c4312d42a743cd4df9ff8d1de",
      "to": "sha256:57daa087b89f5465ae290d62d1bb53a635c69d98c241f756cdef99a7fe850e42",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ]
      }
    },
    {
      "from": "sha256:57daa087b89f5465ae290d62d1bb53a635c69d98c241f756cdef99a7fe850e42",
      "to": "sha256:545f48b74b97aa58d0569865213d72c52812405f94bebcdf658ed962803cbda3",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ]
      }
    },
    {
      "from": "sha256:545f48b74b97aa58d0569865213d72c52812405f94bebcdf658ed962803cbda3",
      "to": "sha256:12c94ca56db2698862ce62c097467259f0b90432a2f66febcac4f339bb145b72",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ]
      }
    },
    {
      "from": "sha256:12c94ca56db2698862ce62c097467259f0b90432a2f66febcac4f339bb145b72",
      "to": "sha256:2347e1808b88930f9c4ed566a14dac1375ecac2014ec2348e85fa3ba1e294de0",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ]
      }
    },
    {
      "from": "sha256:2347e1808b88930f9c4ed566a14dac1375ecac2014ec2348e85fa3ba1e294de0",
      "to": "sha256:ecee5ab8c00cd885f6d3bc8c0be2f1e014eff67d3b48b412d894647d12e5d6a4",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ]
      }
    },
    {
      "from": "sha256:ecee5ab8c00cd885f6d3bc8c0be2f1e014eff67d3b48b412d894647d12e5d6a4",
      "to": "sha256:4562e6b248984cef019b825bc9527d2d51abbf843d19f08486c863eeb3776e55",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ],
        "PM-D": [
          "PM-D"
        ]
      }
    },
    {
      "from": "sha256:4562e6b248984cef019b825bc9527d2d51abbf843d19f08486c863eeb3776e55",
      "to": "sha256:2d2a1f9466ffc76362c39594e0c13eddc9994c0fa59884d007540acae89abbba",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ],
        "PM-D": [
          "PM-D"
        ]
      }
    },
    {
      "from": "sha256:2d2a1f9466ffc76362c39594e0c13eddc9994c0fa59884d007540acae89abbba",
      "to": "sha256:9991026f819b1b4561ced20dd66fcea641105979c3797228b93b3e4eeac3f255",
      "chunk_ids": {
        "PM-A": [
          "PM-A"
        ],
        "PM-B": [
          "PM-B"
        ],
        "PM-C1": [
          "PM-C1"
        ],
        "PM-C2": [
          "PM-C2"
        ],
        "PM-D": [
          "PM-D"
        ]
      }
    }
  ]
}
```

## PM-A round 1

Frozen pair: a6739258 to 14bec38c, with the chunk tip 868975d7.
Repair allowance: cycle 1 of 2 starts for C1 and C2, and it is also the hardening cycle.

### Standards

Count 0 after disposition. Worst issue: none.
S1 closes as no-op at confidence 6, because ticket 1 directs the independent note copy in `internal/probe/refusal_test.go`.

### Spec

Count 0. Worst issue: none. Rows PM45 to PM50 are met.

### Coverage

Count 2. Worst issue: C1.

- C1, auto-fix, confidence 8: the stub hook records `"$*"`, so `TestProbeRepeatsThePackage` cannot tell two arguments from one joined argument. A joined-value probe read silent. Record each argument on its own line and assert the last two elements.
- C2, auto-fix, confidence 6: no row pins the empty-value refusal for a later `--package`. Add a row to `internal/probe/refusal_test.go` with exit 2, the subject unchanged, and no run child.

### Advice

- Four floor rows share one combined excerpt, and each row names one package.
- A duplicate `--package` value has no decided behavior.

## PM-A round 2

The confirming round read the repair delta from 14bec38c to 57336246, with the repair commit cc2abf92.
PM-A closes with repair cycle 1 of 2 used, and that cycle was also the hardening cycle.

### Standards

Count 0. Worst issue: none. S1 stays closed as no-op.

### Spec

Count 0. Worst issue: none. PM45 and PM47 stay met, and C2 adds no behavior.

### Coverage

Count 0. Worst issue: none. C1 and C2 are closed, and each has an independent probe that bit.

### Advice

- A companion row with the empty value first would pin the other order.
- A shared `probetest` argv recorder fits if a third hook appears.

## PM-B round 1

Frozen pair: e5ca6b0c to 1b8f4f66, with the chunk tip 27346d11.
Repair allowance: cycle 1 of 2 starts for S1, C1, and C2. The plan adds rows PM79 and PM80 for C2 and C1, so the cycle is a repair, not a hardening cycle.

### Standards

Count 1. Worst issue: S1.

- S1, auto-fix, confidence 4: `TestHeldRunsInterruptACancelledBuild` makes a third copy of the parking-builder arm-and-wait wiring of `cancel_test.go`. Put that wiring in one `cancel_test.go` helper for all three sites.

### Spec

Count 0. Worst issue: none. Rows PM27 to PM31, PM70, and PM76 are met.

### Coverage

Count 2. Worst issue: C2.

- C2, auto-fix, confidence 8: no row pins a stale owned binary whose fresh build fails with no cancel. A mutant that reuses the stale binary read silent. Row PM79 pins `build-failed` with no Go child.
- C1, auto-fix, confidence 8: no row pins a cancel during the first held selection. A mutant that applies the interrupt rule only to later runs read silent. Row PM80 pins `interrupted`, the parked directory gone, and a build count of 0.

### Advice

- A context check before the fresh build would avoid one builder start after a cancel.
- A helper beside `selectionFailure` could hold the interrupted selection line.

## PM-B round 2

The confirming round read the delta from 1b8f4f66 to c20478f5, with the repair commit 78584bde.
PM-B closes with repair cycle 1 of 2 used.

### Standards

Count 0. Worst issue: none. S1 is closed, because one `parking` helper in `cancel_test.go` serves all three sites.

### Spec

Count 1. Worst issue: P1.

- P1, auto-fix, confidence 6: the plan expansion 2c88de8f had no `bench learning` entry before the dispatch. The orchestrator captured the entry, and the Spec session reaffirmed P1 closed in entry pm-b-spec-r2b. This correction is evidence-only and uses no repair cycle.

The plan expansion is correct and in scope, and rows PM79 and PM80 are met.

### Coverage

Count 0. Worst issue: none. C1 and C2 are closed, and the `cancel_test.go` refactor keeps the cancel tests red-capable.

### Advice

- An omission of one `drainBuilderGroup` call in `runbinary.go` leaves the cancel test green. That code is older than this build.

## PM-C1 round 1

Frozen pair: 32dc3858 to c4ec0f13, with the chunk tip dd3cdaec.
PM-C1 uses no repair cycle, because its one finding is a comment-only correction.

### Standards

Count 1. Worst issue: S1.

- S1, auto-fix, confidence 5: the `expectFitsSelection` doc comment claims a test-name check that the body does not make. The orchestrator rewrites the comment to the narrow prose rule, and the issuing axis reaffirms it.

### Spec

Count 0. Worst issue: none. The prefactor keeps the usage line, the flag set, and every refusal, and its seam fits ticket 4.

### Coverage

Count 0. Worst issue: none. The existing tests pin every moved rule except two gaps that are older than this build.

### Advice

- No row pins an `--insert-before` with two texts. That gap is older than this build.
- The `withs > 0 && withFile` rule is redundant, because later rules refuse every case that it refuses.
- Ticket 4 must pin the argument roles of `expectFitsSelection`.

## PM-C1 close

The orchestrator committed the comment-only correction 1c524805 and refroze PM-C1 from 32dc3858 to 1c524805.
The ticket 3 author re-recorded its five requirements at that tip, and the central probe bit again.
Each axis reaffirmed its result at that tip, and S1 is closed. PM-C1 used no repair cycle.

## PM-C2 round 1

Frozen pair: 6584e93e to dfc5b339, with the chunk tip b647725e.
Repair allowance: cycle 1 of 2 starts for S1 to S3, P1, P2, and C1 to C5.
The orchestrator took reds for PM6, PM67, and PM69 with three probes, and each probe bit with the `restored` cell `yes`.

### Standards

Count 3. Worst issue: S1.

- S1, auto-fix, confidence 6: `set.go` repeats the single form's refusal order after the mutation input, and it composes `mutationInput` and `mutate` a second time. Share both through `probe.go` helpers.
- S2, auto-fix, confidence 5: `baselineCells` repeats the cell pair that `untouchedCells` owns.
- S3, auto-fix, confidence 5: the lock holder copies the helper re-entry argv, and it adds a third poll loop.

### Spec

Count 3. Worst issue: P1.

- P1, ask-user, confidence 6: a stop or a failed baseline overwrote a missed row. Decision D3 keeps `substring-miss`, and row PM81 pins it.
- P2, ask-user, confidence 6: a failed preservation read `mutation-failed`. Decision D3 adds the cause `preserve-failed`, and row PM82 pins it.
- P3, no-op, confidence 7: the PM20 why cell described stacked edits. The plan commit fixes that cell, and the code is correct.

### Coverage

Count 5. Worst issue: C1.

- C1, auto-fix, confidence 8: no row pinned a published write failure whose restore fails. Row PM83 pins it.
- C2, auto-fix, confidence 7: no row pinned a signal during a restore. Row PM84 pins it.
- C3, auto-fix, confidence 7: `trailing_bytes` was not asserted for the `preserved` table or the baseline report. Row PM85 pins it.
- C4, auto-fix, confidence 5: no row pinned `builds` 1 for a failed owned baseline. Row PM86 pins it.
- C5, ask-user, confidence 4: this finding is P2, and D3 decides it.

### Advice

- An interrupted fresh build prints a `mutant_report` for a run whose Go child never started.
- The `row.exit == 2` break adds nothing to the maximum.

## PM-C2 round 2

The confirming round read the delta from dfc5b339 to 28d4a4e4, with the repair commit 4b9543cf.
Repair allowance: cycle 2 of 2 starts for S4, P21, and P22.

### Standards

Count 1. Worst issue: S4. S1, S2, and S3 are closed.

- S4, auto-fix, confidence 6: two older spellings of `sleepingRun` remain in `baseline_test.go` and `lock_test.go`.

### Spec

Count 2. Worst issue: P21. P1, P2, and P3 of round 1 are closed.

- P21, auto-fix, confidence 6: the new sentence about failure text says too much, because baseline and interrupted reports print error lines. The orchestrator narrows it in a plan commit.
- P22, no-op, confidence 6: the PM24 why cell names the removed interrupted branch. The orchestrator rewords it in the same plan commit.

### Coverage

Count 0. Worst issue: none. C1 to C5 and row PM81 are closed.

### Advice

- No test pins a missed row after a stop at a later index. The S4 repair adds that subtest to row PM81.
- The PM84 test waits 50 ms for the cancel, so a starved host can give a false red. A deterministic wait needs a production seam, so it stays a reviewer call.

## PM-C2 round 3

The narrow round read the delta from 9d484121 to ccd8c12c, with the small-fix commit ae0f9b8b.
Repair allowance: both cycles are used. Under `--auto-approve`, the reviewer approved repair rounds past the allowance in advance, so cycle 3 starts for S5 alone.

### Standards

Count 1. Worst issue: S5. S4 is closed.

- S5, auto-fix, confidence 7: the `restore-failed` row is now spelled by hand in three tests. Add one `restoreFailedRow` helper beside `biteRow` and use it at all three sites.

### Spec

Count 0. Worst issue: none. P21 and P22 are closed, and the new subtest meets row PM81.

### Coverage

Count 0. Worst issue: none. A miss row after a later stop is now pinned, and the `sleepingRun` tests stay red-capable.

## PM-C2 close

The cycle 3 small fix 485fa5e2 adds one `restoreFailedRow` helper, and PM-C2 is refrozen from 6584e93e to 485fa5e2.
Each axis reaffirmed its result at that tip, and S5 is closed.
PM-C2 used three repair cycles, and the third cycle used the reviewer's `--auto-approve` extension.

### Advice for the reviewer

- Decision D3 came from the build consultant, so the reviewer can veto it.
- The PM84 test waits 50 ms for the cancel. A deterministic wait needs a production test seam for the set cancel context.

## PM-D round 1

Frozen pair: 1870269d to 520f4899, with the chunk tip 6b447a69.
Repair allowance: cycle 1 of 2 starts for S1, P1, P2, and C1.
The orchestrator took the PM58 red with a probe that bit, and the `restored` cell read `yes`.

### Standards

Count 1. Worst issue: S1.

- S1, ask-user, confidence 7: `MutantResult` gained a `CommandExit` field that no production code reads, against the four fields of the spec. The spec decides it, so the repair removes the field. The plan probe becomes a mutation over the declared fields.

### Spec

Count 2. Worst issue: P1.

- P1, ask-user, confidence 7: this finding is S1, and the same repair closes it.
- P2, ask-user, confidence 6: the PM59 record row has empty cells, not the producer's `substring-miss` row. The repair builds that row from the producer, through a probe run or the producer's own row constructor.

### Coverage

Count 1. Worst issue: C1.

- C1, auto-fix, confidence 8: no test imports a row whose restore failed. A mutant that always records `pass` read silent. Row PM87 pins the `fail` side of the restore binding.

### Advice

- The reader test rows 1 and 3 project to the same result, so only the record tests catch a reversed row index.

## PM-D round 2

The confirming round read the delta from 520f4899 to d4b2a7f1, with the repair commit fc4c772e.
PM-D closes with repair cycle 1 of 2 used.

### Standards

Count 0. Worst issue: none. S1 is closed, because `MutantResult` holds exactly the four spec fields.

### Spec

Count 0. Worst issue: none. P1 and P2 are closed, and row PM87 is met.

### Coverage

Count 0. Worst issue: none. C1 is closed, and six independent probes bit with the `restored` cell `yes`.

### Advice

- The probe package's own clamp harness in `probe_test.go` is a second harness, and a learning records it for a later change.
- The spec still allows an exported renderer for record rows, but the build now uses probe runs only.

## PM-D checkpoint red

The PM-D checkpoint gate at d5b3784e is red in `TestFreshnessPublicationTopology`.
The check rejects the sealed-publication call that commit fc4c772e moved into `internal/probe/probetest/fixture.go`.
The check skips only `_test.go` files, so the call must return to the test callers.
Repair allowance: cycle 2 of 2 starts for this red, and the fix changes no oracle.
