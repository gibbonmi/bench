# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/retired-store-sweep/spec.md",
  "plan_digest": "sha256:d56366661c354b73477e981cfa5c28d2c881d5b2337c36cde8f0fc5412b7f6dc",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RS-C1",
      "base": "63d5231343049fa9fb1b3b3905c395f22dc6de65",
      "tip": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
      "plan_digest": "sha256:d56366661c354b73477e981cfa5c28d2c881d5b2337c36cde8f0fc5412b7f6dc",
      "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
      "acceptance_rows": [
        "RS1",
        "RS2",
        "RS3",
        "RS25",
        "RS26"
      ],
      "verification": [
        {
          "id": "rs-c1-01-retired-record-written",
          "performer": "ft446-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-01",
            "digest": "sha256:4cd95b7a55f520200a8f357ef036c3bb072636e8a009728acd9376fe95f6c5fb",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,omit,failed,1,TestRetiredStoreRecordsItsPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestRetiredStoreRecordsItsPath$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,537\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestRetiredStoreRecordsItsPath,\"lifecycle_test.go:72: retired record = {Store: ResourceID: Kind: ResourcePath: Binding:} (lstat /dev/shm/OSPGWH/t/TestRetiredStoreRecordsItsPath3599939175/001/cache/store/retired.json: no such file or directory), want {Store: ResourceID: Ki… (481 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "requirement": "retired-record-written",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "remove the retired-record publish from SealAndObserve; require TestRetiredStoreRecordsItsPath to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-01",
              "digest": "sha256:4cd95b7a55f520200a8f357ef036c3bb072636e8a009728acd9376fe95f6c5fb",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,omit,failed,1,TestRetiredStoreRecordsItsPath,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestRetiredStoreRecordsItsPath$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,537\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,2,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestRetiredStoreRecordsItsPath,\"lifecycle_test.go:72: retired record = {Store: ResourceID: Kind: ResourcePath: Binding:} (lstat /dev/shm/OSPGWH/t/TestRetiredStoreRecordsItsPath3599939175/001/cache/store/retired.json: no such file or directory), want {Store: ResourceID: Ki… (481 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c1-01-cache-no-generation",
          "performer": "ft446-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-01",
            "digest": "sha256:cef4d6be147e3d0a719fc57b918f95eebeb19ffb49f85e798490e5e64944b1d7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,insert,failed,1,TestProtectionCompletedHold,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,^TestProtectionCompletedHold$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,413\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionCompletedHold,\"protection_test.go:139: clean of a cache with a retired store = (\\\"error: cache protected — protection: the store holds a retired record\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "requirement": "cache-no-generation",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "make Recover return a non-ErrNotExist error when the store holds retired.json; require TestProtectionCompletedHold to fail on its clean after the release",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-01",
              "digest": "sha256:cef4d6be147e3d0a719fc57b918f95eebeb19ffb49f85e798490e5e64944b1d7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,insert,failed,1,TestProtectionCompletedHold,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,^TestProtectionCompletedHold$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,413\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestProtectionCompletedHold,\"protection_test.go:139: clean of a cache with a retired store = (\\\"error: cache protected — protection: the store holds a retired record\\\\\\\\n\\\", 1), want exit 0\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c1-01-checkout-no-generation",
          "performer": "ft446-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-01",
            "digest": "sha256:9272af06f0f2db186f382d561d5042251986b161a2c7d4ee0bc9be159e8ba0c6",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,insert,failed,1,TestCheckoutProtectionPredicate,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,^TestCheckoutProtectionPredicate$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,319\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,13,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCheckoutProtectionPredicate,\"protection_test.go:58: an unused generation kept the checkout\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "requirement": "checkout-no-generation",
          "command": "bench test --package ./internal/worktree --run '^TestCheckoutProtectionPredicate$'",
          "exit_code": 0,
          "probe": {
            "mutation": "make Recover return a non-ErrNotExist error when the store holds retired.json; require TestCheckoutProtectionPredicate to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-01",
              "digest": "sha256:9272af06f0f2db186f382d561d5042251986b161a2c7d4ee0bc9be159e8ba0c6",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,insert,failed,1,TestCheckoutProtectionPredicate,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/worktree,^TestCheckoutProtectionPredicate$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,319\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,fail,13,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/worktree,TestCheckoutProtectionPredicate,\"protection_test.go:58: an unused generation kept the checkout\",1\nskips[0]{package,test,reason}:\n"
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
