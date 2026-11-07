# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/retired-store-sweep/spec.md",
  "plan_digest": "sha256:19ee2fb6487fbd0b39017eb392ac3d31bdab79d1b14798b737b4797eb880470f",
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
      "reviews": [
        {
          "id": "rs-c1-review-standards",
          "performer": "ft446-review-rs-c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c1",
            "digest": "sha256:e9d9461e21e186d34aca206fba70a815da52eef2a4cbb6f8ee51df9ad6a8f96c",
            "excerpt": "Standards: 0 findings. Worst issue: none. Advice: reflow the long comment at internal/worktree/protection_test.go:32; a retiredOf(d) constructor could hold the field mapping that recovery.go:67 builds by hand.\n"
          },
          "axis": "Standards",
          "base": "63d5231343049fa9fb1b3b3905c395f22dc6de65",
          "tip": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "rs-c1-review-spec",
          "performer": "ft446-review-rs-c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c1",
            "digest": "sha256:cb827d0ec2a1414f35ee0584269a4325dac3b827f46d4c64b3583892ad0c79b2",
            "excerpt": "Spec: 0 findings. Worst issue: none. D1 and D7 hold; retiredRecord has schema, kind, resource_path, binding; valid reuses Descriptor.checkOwner; the publish sits after removeUser and before the resource.json removal; Recover, NoGeneration, join, and ObserveRecorded are unchanged. RS1, RS2, RS3, RS25, RS26 audited.\n"
          },
          "axis": "Spec",
          "base": "63d5231343049fa9fb1b3b3905c395f22dc6de65",
          "tip": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "rs-c1-review-coverage",
          "performer": "ft446-review-rs-c1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "e620d451fb5418cf22b640b57fbe2a3ed25619b4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c1",
            "digest": "sha256:1f83bc3cf16e3ac79e748e18284dfa13ddd664452f7796504630cccbfa6494d9",
            "excerpt": "Coverage: 0 findings. Worst issue: none. RS26 keeps a red-capable expectation: a Recover refusal of a retired-only store bit TestCheckoutProtectionPredicate at protection_test.go:61 (restored yes). Advice: a swallowed retired-publish error stays green (probe silent); a fault-router test would pin it (hardening, severity 2-3).\n"
          },
          "axis": "Coverage",
          "base": "63d5231343049fa9fb1b3b3905c395f22dc6de65",
          "tip": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "RS-C2",
      "base": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
      "tip": "d7641b284d3daa0e8a7a7232508b444f9990e914",
      "plan_digest": "sha256:99155751836a1f9efa511515c6441606fcb5eabd4f3f81b8d4aed28f9e0a65b5",
      "source_digest": "4fedf648b82862ca8b8a15db6f224420d346c5ac",
      "acceptance_rows": [
        "RS4",
        "RS5",
        "RS6",
        "RS7",
        "RS8",
        "RS14",
        "RS15",
        "RS31"
      ],
      "verification": [
        {
          "id": "rs-c2-02-retired-sweep-omission",
          "performer": "ft446-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-02",
            "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "requirement": "retired-sweep-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "make discardNested skip each store with no resource.json, as it does today; require TestGenerationRemove to fail on the retired gone store",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-02",
              "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "rs-c2-02-user-guard-omission",
          "performer": "ft446-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-02",
            "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "requirement": "user-guard-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the empty-users condition from the retired-store predicate; require TestRetiredSweepKeepsAUser to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-02",
              "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "rs-c2-02-lock-skip-omission",
          "performer": "ft446-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-02",
            "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "requirement": "lock-skip-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "make a contended store lock join the sweep error; require TestRetiredSweepSkipsALockedStore to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-02",
              "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "rs-c2-02-r1-retired-sweep-omission",
          "performer": "ft446-repair-02-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4fedf648b82862ca8b8a15db6f224420d346c5ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-repair-02-r1",
            "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "requirement": "retired-sweep-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "make discardNested skip each store with no resource.json, as it does today; require TestGenerationRemove to fail on the retired gone store",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-repair-02-r1",
              "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "rs-c2-02-r1-user-guard-omission",
          "performer": "ft446-repair-02-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4fedf648b82862ca8b8a15db6f224420d346c5ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-repair-02-r1",
            "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "requirement": "user-guard-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the empty-users condition from the retired-store predicate; require TestRetiredSweepKeepsAUser to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-repair-02-r1",
              "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "rs-c2-02-r1-lock-skip-omission",
          "performer": "ft446-repair-02-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "4fedf648b82862ca8b8a15db6f224420d346c5ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-repair-02-r1",
            "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "requirement": "lock-skip-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "make a contended store lock join the sweep error; require TestRetiredSweepSkipsALockedStore to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-repair-02-r1",
              "digest": "sha256:59585ba7aea2d654c734f26436fb27fb1c7612e80044f3f67ba8b91767208c16",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,omit,failed,1,TestGenerationRemove,1,1,0,null,yes\n  2,bit,swap,failed,1,TestRetiredSweepKeepsAUser,1,1,0,null,yes\n  3,bit,swap,failed,1,TestRetiredSweepSkipsALockedStore,1,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/subprocess/protection,^TestGenerationRemove$|^TestRetiredSweepKeepsAUser$|^TestRetiredSweepSkipsALockedStore$,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "rs-c2-review-standards",
          "performer": "ft446-review-rs-c2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c2",
            "digest": "sha256:eaa8109402c01efa8470873deaa6ce6cdf3880babf0e730464d1738e6520ce37",
            "excerpt": "Standards: 0 findings. Worst issue: none. lock, tryLock, and lockStores share one routine; the ADR and DATA_HANDLING.md hold no path or code. Advice: the ADR line 'A retired store that names a user stays' could say 'A retired store whose users directory holds an entry stays'.\n"
          },
          "axis": "Standards",
          "base": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "tip": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "rs-c2-review-spec",
          "performer": "ft446-review-rs-c2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c2",
            "digest": "sha256:c72e63bf29e9a4a103d09f4d5e64eaba53d799868d572cfadafdd498186b3fb3",
            "excerpt": "Spec: 0 findings. Worst issue: none. RS4-RS8, RS14, RS15, RS31 met inside the ticket 02 fence; the D2 predicate matches directory.go:68-94; discardNested adds the dir+/ prefix test; the live branch is unchanged; discardRetired(store, selected) serves ticket 03 with no second predicate.\n"
          },
          "axis": "Spec",
          "base": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "tip": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "rs-c2-review-coverage",
          "performer": "ft446-review-rs-c2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c2",
            "digest": "sha256:d383c8e16aae3ad1ac58c2641b83ea7a3f4fcdae473309bea3963f64615eb572",
            "excerpt": "Coverage: 1 finding. Worst issue: RS-C2-C1. RS-C2-C1 | record_test.go and directory.go:83 | severity 3 | confidence 8 | auto-fix | no test grades a retired record that decodes but fails valid() (wrong schema, relative path, non-clean path, control byte); dropping the readRetired error check stayed silent. Add retired cases through requireRetiredKept. Safety audit held: user entry, prefix sibling, and non-blocking lock probes bit.\n"
          },
          "axis": "Coverage",
          "base": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "tip": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "finding_ids": [
            "RS-C2-C1"
          ],
          "supersedes": []
        }
      ],
      "confirmations": [
        {
          "id": "rs-c2-confirm-r1",
          "performer": "ft446-orchestrator",
          "role": "coordinator-confirmation",
          "model": "opus",
          "effort": "high",
          "source_digest": "4fedf648b82862ca8b8a15db6f224420d346c5ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "probe:record.go-valid-returns-nil",
            "digest": "sha256:5587f6bc54347974137f0506ebadc8203b8fa642a8be5e1ae30d08996ffa14e5",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/record.go,swap,failed,2,TestRecordRetiredPath/control_byte,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestRecordRetiredPath$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,890\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,3,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestRecordRetiredPath/control_byte,\"record_test.go:236: store of /dev/shm/MR322K/t/TestRecordRetiredPathcontrol_byte420970758/002/gone after the sweep: lstat /dev/shm/MR322K/t/TestRecordRetiredPathcontrol_byte420970758/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c… (328 bytes)\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestRecordRetiredPath/unclean,\"record_test.go:236: store of /dev/shm/MR322K/t/TestRecordRetiredPathunclean2460236597/002/gone after the sweep: lstat /dev/shm/MR322K/t/TestRecordRetiredPathunclean2460236597/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/binar… (320 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "base": "5ceb326b7e240fdc1b2af3bf76b7d37a71db8d13",
          "reviewed_tip": "bef1597e66827888f6fbdd964ac136d50f035aed",
          "repaired_tip": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "finding_ids": [
            "RS-C2-C1"
          ],
          "supersedes": [],
          "probe": {
            "mutation": "make retiredRecord.valid return nil after the schema check, so the path owner check never runs",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "probe:record.go-valid-returns-nil",
              "digest": "sha256:5587f6bc54347974137f0506ebadc8203b8fa642a8be5e1ae30d08996ffa14e5",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/record.go,swap,failed,2,TestRecordRetiredPath/control_byte,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestRecordRetiredPath$,passed,4\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,890\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,3,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestRecordRetiredPath/control_byte,\"record_test.go:236: store of /dev/shm/MR322K/t/TestRecordRetiredPathcontrol_byte420970758/002/gone after the sweep: lstat /dev/shm/MR322K/t/TestRecordRetiredPathcontrol_byte420970758/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c… (328 bytes)\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestRecordRetiredPath/unclean,\"record_test.go:236: store of /dev/shm/MR322K/t/TestRecordRetiredPathunclean2460236597/002/gone after the sweep: lstat /dev/shm/MR322K/t/TestRecordRetiredPathunclean2460236597/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/binar… (320 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ]
    },
    {
      "id": "RS-C3",
      "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
      "tip": "2220bf07a832c326e1184acb54419fc11798badb",
      "plan_digest": "sha256:19ee2fb6487fbd0b39017eb392ac3d31bdab79d1b14798b737b4797eb880470f",
      "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
      "acceptance_rows": [
        "RS9",
        "RS10",
        "RS11",
        "RS12",
        "RS13",
        "RS27",
        "RS28",
        "RS30"
      ],
      "verification": [
        {
          "id": "rs-c3-03-cache-check-omission",
          "performer": "ft446-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-03",
            "digest": "sha256:e9d8b089027f5895e2fd7fe61c1fb04a37c7f5ff0d66615860ab0b6e5ab66910",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gocache/lock.go,swap,failed,1,TestReleaseRunsTheRetiredCheck,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,^TestReleaseRunsTheRetiredCheck$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,15,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseRunsTheRetiredCheck,\"protection_test.go:189: the last uncontended release ran the check with retired generations [], want one call after the retirement\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2220bf07a832c326e1184acb54419fc11798badb",
          "requirement": "cache-check-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "remove the absent-resource check call from the cache retire; require TestReleaseRunsTheRetiredCheck to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-03",
              "digest": "sha256:e9d8b089027f5895e2fd7fe61c1fb04a37c7f5ff0d66615860ab0b6e5ab66910",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gocache/lock.go,swap,failed,1,TestReleaseRunsTheRetiredCheck,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,^TestReleaseRunsTheRetiredCheck$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,15,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseRunsTheRetiredCheck,\"protection_test.go:189: the last uncontended release ran the check with retired generations [], want one call after the retirement\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c3-03-kind-scope-omission",
          "performer": "ft446-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-03",
            "digest": "sha256:3175d3bf7a1b3d9fbb6a8af1e9f46c9c318b6e84809a4a275d5bc7793cbd3cf1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/directory.go,swap+swap,failed,1,TestDiscardRetired,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestDiscardRetired$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,525\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestDiscardRetired,\"directory_test.go:246: store of other kind after the cache check: lstat /dev/shm/XUMFBL/t/TestDiscardRetired2120579337/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/binary-03227c8793dd2dcf7b15e4a38651ff1d: no such file or dire… (264 bytes)\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2220bf07a832c326e1184acb54419fc11798badb",
          "requirement": "kind-scope-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "drop both the store-name prefix filter and the record kind comparison from the absent-resource check; require TestDiscardRetired to fail on the retired store of another kind",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-03",
              "digest": "sha256:3175d3bf7a1b3d9fbb6a8af1e9f46c9c318b6e84809a4a275d5bc7793cbd3cf1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/directory.go,swap+swap,failed,1,TestDiscardRetired,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestDiscardRetired$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,525\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestDiscardRetired,\"directory_test.go:246: store of other kind after the cache check: lstat /dev/shm/XUMFBL/t/TestDiscardRetired2120579337/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/binary-03227c8793dd2dcf7b15e4a38651ff1d: no such file or dire… (264 bytes)\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c3-03-registry-row",
          "performer": "ft446-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-03",
            "digest": "sha256:e758ae7587cead58273837153ff99792631bf0477760840f58a3286a53c06a60",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_registry_test.go,omit,failed,2,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,all,passed,1218\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1095\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,44474,1218\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,\"injected_ports_test.go:247: a reintroduced canary audit: want its zero-port diagnostic, got [injected port derivation found no ports: internal/canary is present in the graded tree and the audit records it as declaring them, so an empty inve… (607 bytes)\",1\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: test seam has no registry row: internal/gocache.discardRetiredCaches is replaced at internal/gocache/protection_test.go:153 — add a variable row to injectedPortRegistry naming its production owner, consumer test… (284 bytes)\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "source_commit": "2220bf07a832c326e1184acb54419fc11798badb",
          "requirement": "registry-row",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "delete the discardRetiredCaches seamVariable row from the injected-port registry; require the registry check to refuse the unregistered variable",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-03",
              "digest": "sha256:e758ae7587cead58273837153ff99792631bf0477760840f58a3286a53c06a60",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_registry_test.go,omit,failed,2,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,all,passed,1218\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1095\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,44474,1218\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,\"injected_ports_test.go:247: a reintroduced canary audit: want its zero-port diagnostic, got [injected port derivation found no ports: internal/canary is present in the graded tree and the audit records it as declaring them, so an empty inve… (607 bytes)\",1\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: test seam has no registry row: internal/gocache.discardRetiredCaches is replaced at internal/gocache/protection_test.go:153 — add a variable row to injectedPortRegistry naming its production owner, consumer test… (284 bytes)\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
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
  },
  "amendments": [
    {
      "from": "sha256:d56366661c354b73477e981cfa5c28d2c881d5b2337c36cde8f0fc5412b7f6dc",
      "to": "sha256:2d22011a8556aa7de0f62d1d5720662c09790fa7124e337c6d6df82185339255",
      "chunk_ids": {
        "RS-C1": [
          "RS-C1"
        ]
      }
    },
    {
      "from": "sha256:2d22011a8556aa7de0f62d1d5720662c09790fa7124e337c6d6df82185339255",
      "to": "sha256:99155751836a1f9efa511515c6441606fcb5eabd4f3f81b8d4aed28f9e0a65b5",
      "chunk_ids": {
        "RS-C1": [
          "RS-C1"
        ],
        "RS-C2": [
          "RS-C2"
        ]
      }
    },
    {
      "from": "sha256:99155751836a1f9efa511515c6441606fcb5eabd4f3f81b8d4aed28f9e0a65b5",
      "to": "sha256:19ee2fb6487fbd0b39017eb392ac3d31bdab79d1b14798b737b4797eb880470f",
      "chunk_ids": {
        "RS-C1": [
          "RS-C1"
        ],
        "RS-C2": [
          "RS-C2"
        ]
      }
    }
  ]
}
```

## RS-C1 author witness reds

The ticket 01 author ft446-author-01 reported five witness reds, and each probe reported a restored subject:

1. With no retired publish, TestRetiredStoreRecordsItsPath and TestRetiredStoreHasNoGeneration were red before the publish edit.
2. With no production edit, TestCheckoutProtectionPredicate was red at its retired.json requirement.
3. A `Recover` refusal of a store that holds retired.json made TestRetiredStoreHasNoGeneration red.
4. The same mutation made TestProtectionCompletedHold red at the clean after the release.
5. A `join` refusal of a store that holds retired.json made TestProtectionCompletedHold red at the next hold.

The orchestrator wrote these reds after the RS-C1 review charge, not before it. The unified reviewer read the author return through the charge.

## RS-C1 review pickup

Unified review by ft446-review-rs-c1 on sonnet at high effort over the pair 63d5231..5ceb326. Repair cycles used: 0 of 2.

### Standards

Findings: 0. Worst issue: none.

### Spec

Findings: 0. Worst issue: none.

### Coverage

Findings: 0. Worst issue: none. The RS26 expectation is red-capable at its cleanup assertion: a `Recover` refusal of a store with only retired.json made TestCheckoutProtectionPredicate fail at that assertion.

### Advice

- Reflow the long comment line in the checkout protection test.
- A constructor of the retired record from a descriptor could hold the field mapping that `SealAndObserve` builds by hand.
- A fault-router test could pin a swallowed retired-publish error, which no test catches now. This is hardening, not a correctness defect.

## RS-C2 author witness reds

The ticket 02 author ft446-author-02 reported these witness reds, and each probe reported a restored subject:

1. RS5: TestGenerationRemove was red before any production edit, because the retired gone store stayed.
2. RS6: with no resource-absence condition in the predicate, TestGenerationRemove was red.
3. RS7: with no separator in the inside test, TestGenerationRemove was red on the prefix sibling.
4. RS4: with a lax reader in place of the strict record reader, TestRecordUnknownField was red.
5. The non-blocking lock: with a blocking lock in its place, TestRetiredSweepSkipsALockedStore timed out.

The three plan probes ran as one mutant set, and each row bit with a restored subject.

## RS-C2 review pickup

Unified review by ft446-review-rs-c2 on sonnet at high effort over the pair 5ceb326..bef1597. Raw findings: 1. Repair targets: 1. Repair cycles used: 0 of 2 before the repair below.

### Standards

Findings: 0. Worst issue: none.

### Spec

Findings: 0. Worst issue: none.

### Coverage

Findings: 1. Worst issue: RS-C2-C1.

- RS-C2-C1 (auto-fix, confidence 8): no test grades a retired record that decodes but fails its validity check. A wrong schema, a relative path, a path that is not clean, or a control byte falls in this class. The review probe that drops the retired-read error check in the predicate stayed silent. The spec's hostile-input line says that RS4 grades a relative path and control bytes. The fix adds retired-record cases to the record tests through the existing kept-store helper. The fix is test-only, and it goes to the fresh repair session ft446-repair-02-r1.

### Advice

- The ADR sentence about a retired store that names a user could name the users directory entry instead.
- No test reaches the pre-lock skip of a live store or the re-check after the lock. The first is a cost contract, and the second is a race guard.

## RS-C2 repair r1

The repair session ft446-repair-02-r1 made a test-only repair at d7641b2. It reported these witness reds, and each probe reported a restored subject:

1. With no retired-read error check in the predicate, TestRecordWrongSchema was red.
2. With the same mutation, TestRecordRetiredPath was red in its control-byte and unclean subtests.
3. With no retired-read error check and no inside-directory condition, TestRecordRetiredPath was red in all three subtests.

The orchestrator probe made the retired validator return no error after its schema check. TestRecordRetiredPath was red, and the subject was restored. The confirmation rs-c2-confirm-r1 closes RS-C2-C1, so RS-C2 takes no confirming round. Repair cycles used: 1 of 2.
