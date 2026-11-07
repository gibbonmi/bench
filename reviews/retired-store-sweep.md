# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/retired-store-sweep/spec.md",
  "plan_digest": "sha256:ce20aca3c19ac69f01f324f724ef82aed6df6f7437403fb7fbfdc15715ccad84",
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
      "tip": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
      "plan_digest": "sha256:c10f32d8ce7a14d274b01166ca9d0f66f743e88e7d0a5fd13afa01bf51ff5af9",
      "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
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
        },
        {
          "id": "rs-c3-03-r1-cache-check-omission",
          "performer": "ft446-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-repair-03-r1",
            "digest": "sha256:944c857819b3e0784e31fdc4d88ed2cab7aa1ba1787665e18d45b72897bfc9ca",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gocache/lock.go,swap,failed,1,TestReleaseRunsTheRetiredCheck,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,^TestReleaseRunsTheRetiredCheck$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,14,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseRunsTheRetiredCheck,\"protection_test.go:189: the last uncontended release ran the check with retired generations [], want one call after the retirement\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
          "requirement": "cache-check-omission",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "remove the absent-resource check call from the cache retire; require TestReleaseRunsTheRetiredCheck to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-repair-03-r1",
              "digest": "sha256:944c857819b3e0784e31fdc4d88ed2cab7aa1ba1787665e18d45b72897bfc9ca",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gocache/lock.go,swap,failed,1,TestReleaseRunsTheRetiredCheck,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,^TestReleaseRunsTheRetiredCheck$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,385\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,14,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseRunsTheRetiredCheck,\"protection_test.go:189: the last uncontended release ran the check with retired generations [], want one call after the retirement\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c3-03-r1-kind-scope-omission",
          "performer": "ft446-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-repair-03-r1",
            "digest": "sha256:c90788c64d491ac37b27ac87b87091a40827cb253ae8cfa99de454921068d95c",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/directory.go,swap,failed,1,TestDiscardRetired,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestDiscardRetired$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,525\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestDiscardRetired,\"directory_test.go:235: store of other kind after the cache check: lstat /dev/shm/43ZAMY/t/TestDiscardRetired4229170059/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/binary-2b11dcecaa32a7adf66bafa88d254043: no such file or dire… (264 bytes)\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
          "requirement": "kind-scope-omission",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "drop both the store-name prefix filter and the record kind comparison from the absent-resource check; require TestDiscardRetired to fail on the retired store of another kind",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-repair-03-r1",
              "digest": "sha256:c90788c64d491ac37b27ac87b87091a40827cb253ae8cfa99de454921068d95c",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/directory.go,swap,failed,1,TestDiscardRetired,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,^TestDiscardRetired$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,525\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestDiscardRetired,\"directory_test.go:235: store of other kind after the cache check: lstat /dev/shm/43ZAMY/t/TestDiscardRetired4229170059/001/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/binary-2b11dcecaa32a7adf66bafa88d254043: no such file or dire… (264 bytes)\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c3-03-r1-registry-row",
          "performer": "ft446-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-repair-03-r1",
            "digest": "sha256:a0d732b5d6defd5d81a0f5928b6646ad341e08aa9ab7abdcb5cd6aa54bd2b2d3",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_registry_test.go,omit,failed,2,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,all,passed,1218\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1095\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,43392,1218\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,\"injected_ports_test.go:247: a reintroduced canary audit: want its zero-port diagnostic, got [injected port derivation found no ports: internal/canary is present in the graded tree and the audit records it as declaring them, so an empty inve… (607 bytes)\",1\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: test seam has no registry row: internal/gocache.discardRetiredCaches is replaced at internal/gocache/protection_test.go:153 — add a variable row to injectedPortRegistry naming its production owner, consumer test… (284 bytes)\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "source_commit": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
          "requirement": "registry-row",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0,
          "probe": {
            "mutation": "delete the discardRetiredCaches seamVariable row from the injected-port registry; require the registry check to refuse the unregistered variable",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-repair-03-r1",
              "digest": "sha256:a0d732b5d6defd5d81a0f5928b6646ad341e08aa9ab7abdcb5cd6aa54bd2b2d3",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_registry_test.go,omit,failed,2,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,all,passed,1218\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1095\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,43392,1218\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestInjectedPortRegistryCheckBites/a_retired_live-tree_audit_fires,\"injected_ports_test.go:247: a reintroduced canary audit: want its zero-port diagnostic, got [injected port derivation found no ports: internal/canary is present in the graded tree and the audit records it as declaring them, so an empty inve… (607 bytes)\",1\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: test seam has no registry row: internal/gocache.discardRetiredCaches is replaced at internal/gocache/protection_test.go:153 — add a variable row to injectedPortRegistry naming its production owner, consumer test… (284 bytes)\",1\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "rs-c3-review-standards",
          "performer": "ft446-review-rs-c3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c3",
            "digest": "sha256:5e5ebc69a6a930a0fcd2719f27e28ce889c6c5243db22f82bc3438113178dc1c",
            "excerpt": "Standards: 0 native findings; advice item (a), severity 3, confidence 7, ask-user: goneStore (directory_test.go:197-205) repeats the retire-then-delete pair of goneRetired (fixture_test.go:209-218); retiredCaches repeats its SetStoreBase and Protocol setup; the fix is goneRetired calling goneStore, which needs fixture_test.go. The orchestrator accepts it as RS-C3-S1 (auto-fix) under the AGENTS.md one-source rule for fixture harnesses. ADR sentence and comments hold.\n"
          },
          "axis": "Standards",
          "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "tip": "2220bf07a832c326e1184acb54419fc11798badb",
          "finding_ids": [
            "RS-C3-S1"
          ],
          "supersedes": []
        },
        {
          "id": "rs-c3-review-spec",
          "performer": "ft446-review-rs-c3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c3",
            "digest": "sha256:ddc0dd69eec0d05c9df1ff116e65d24bb7064f4efc97824c0f4309f73b6704a7",
            "excerpt": "Spec: 0 findings. Worst issue: none. RS9, RS10, RS30, RS11, RS12, RS13, RS27, RS28 mapped; the registry row matches D3 field for field; the release calls the check only after a nil SealAndObserve; the check keeps both kind tests; lock and read skips stay silent; a failed deletion or unlistable root joins the error.\n"
          },
          "axis": "Spec",
          "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "tip": "2220bf07a832c326e1184acb54419fc11798badb",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "rs-c3-review-coverage",
          "performer": "ft446-review-rs-c3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "de044a104a19b71f6e5de1b473254b270f98bd89",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-review-rs-c3",
            "digest": "sha256:98d9c26cfb633793c7c8d5cac619cbedcc15f659b2a19816f4cc955c4384e942",
            "excerpt": "Coverage: 0 findings. Worst issue: none. Probes bit with restored yes: check call removed, check in the contended branch, kind comparison swapped, prefix filter swapped, registry row omitted.\n"
          },
          "axis": "Coverage",
          "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "tip": "2220bf07a832c326e1184acb54419fc11798badb",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "rs-c3-confirm-standards",
          "performer": "ft446-confirm-rs-c3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-confirm-rs-c3",
            "digest": "sha256:d0e11dedddee4eadac4ebfe5a97db1b68606d4031c60b49e5397128e50fc4741",
            "excerpt": "Standards (confirming round): RS-C3-S1 closed; 0 new findings. goneStore has one definition at fixture_test.go:219 and goneRetired calls it; privateProtocol at fixture_test.go:207-215 serves goneRetired, retiredCaches, and openPrivateRun; no other SetStoreBase call in the package tests.\n"
          },
          "axis": "Standards",
          "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "tip": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
          "finding_ids": [],
          "supersedes": [
            "rs-c3-review-standards"
          ]
        },
        {
          "id": "rs-c3-confirm-spec",
          "performer": "ft446-confirm-rs-c3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-confirm-rs-c3",
            "digest": "sha256:0f8ce15c24939fc135344bf05e93eb6ba18480e5db41e5dc856cd9e18f7737de",
            "excerpt": "Spec (confirming round): reaffirmed; 0 new findings. The repair delta changes no production file; the fixture_test.go widen sits in the directory of ticket 03's signed-off Writes entries; validator kinds are unchanged.\n"
          },
          "axis": "Spec",
          "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "tip": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
          "finding_ids": [],
          "supersedes": [
            "rs-c3-review-spec"
          ]
        },
        {
          "id": "rs-c3-confirm-coverage",
          "performer": "ft446-confirm-rs-c3",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "eb2fea37d68f7aa2826ec949f44a653e2e6053cc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-confirm-rs-c3",
            "digest": "sha256:f5bd471242f083759c88419b6670eb0daf8d51fd55b26af70e7f98f33f727d37",
            "excerpt": "Coverage (confirming round): reaffirmed; 0 new findings. No assertion removed or weakened; protection package 75 tests pass; probe goneStore keeps the resource bit TestDiscardRetired at directory_test.go:235, restored yes.\n"
          },
          "axis": "Coverage",
          "base": "d7641b284d3daa0e8a7a7232508b444f9990e914",
          "tip": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
          "finding_ids": [],
          "supersedes": [
            "rs-c3-review-coverage"
          ]
        }
      ]
    },
    {
      "id": "RS-C4",
      "base": "e6cde66c6c37581723731dc6b7c558ac8edfa7a6",
      "tip": "398691f9cf8dbf2bb1d60c747bc2f7585dc57b6a",
      "plan_digest": "sha256:ce20aca3c19ac69f01f324f724ef82aed6df6f7437403fb7fbfdc15715ccad84",
      "source_digest": "8e8882a8a7fff5820dc1d01f98d70398480d9f2d",
      "acceptance_rows": [
        "RS16",
        "RS17",
        "RS19",
        "RS20",
        "RS21",
        "RS22",
        "RS23",
        "RS24",
        "RS18",
        "RS29"
      ],
      "verification": [
        {
          "id": "rs-c4-04-protocol-forwarding",
          "performer": "ft446-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8e8882a8a7fff5820dc1d01f98d70398480d9f2d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-04",
            "digest": "sha256:9d31f833fb1481937177621ab8ce0a97480b47d6d193a2e6c03f618126bceeac",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/inherited.go,swap,failed,1,TestRunInheritedRegistersOnLocalProtocol,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess,^TestRunInheritedRegistersOnLocalProtocol$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,525\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestRunInheritedRegistersOnLocalProtocol,\"startgate_test.go:162: child vector = [], \u003cnil>; want the local resource [{Store:/dev/shm/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/kit-test-2d7b4f0a75140a8b8124ff82624cd749 ResourceID:530b9b53a0ae8960e06278e1673bc941 Kind:kit… (356 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "398691f9cf8dbf2bb1d60c747bc2f7585dc57b6a",
          "requirement": "protocol-forwarding",
          "command": "bench test --package ./internal/subprocess",
          "exit_code": 0,
          "probe": {
            "mutation": "make RunInherited use the installed protocol whatever the registration names; require TestRunInheritedRegistersOnLocalProtocol to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-04",
              "digest": "sha256:9d31f833fb1481937177621ab8ce0a97480b47d6d193a2e6c03f618126bceeac",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/inherited.go,swap,failed,1,TestRunInheritedRegistersOnLocalProtocol,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess,^TestRunInheritedRegistersOnLocalProtocol$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,525\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess,fail,3,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess,TestRunInheritedRegistersOnLocalProtocol,\"startgate_test.go:162: child vector = [], \u003cnil>; want the local resource [{Store:/dev/shm/bench-process-1000/f8cb4960-efae-461c-83a0-7bdf62294c67/kit-test-2d7b4f0a75140a8b8124ff82624cd749 ResourceID:530b9b53a0ae8960e06278e1673bc941 Kind:kit… (356 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c4-04-build-vector-refusal",
          "performer": "ft446-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8e8882a8a7fff5820dc1d01f98d70398480d9f2d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-04",
            "digest": "sha256:b24a8d769b651a5922d46c092ddeda40b8af41bcb7d0a1479342dbe6da854908",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/inherited.go,omit,failed,1,TestBuildScriptMalformedResources,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,^TestBuildScriptMalformedResources$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,519\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,17,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuildScriptMalformedResources,\"lifetime_resources_test.go:88: Own under a malformed vector = selected Bench executable \\\"/dev/shm/XIBFPO/t/TestBuildScriptMalformedResources258308181/003/bench-run-2517604797/bench\\\" is unavailable: lstat /dev/shm/XIBFPO/t/TestBuildScriptMal… (347 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "398691f9cf8dbf2bb1d60c747bc2f7585dc57b6a",
          "requirement": "build-vector-refusal",
          "command": "bench test --package ./internal/runbinary",
          "exit_code": 0,
          "probe": {
            "mutation": "make RunInherited ignore an unreadable inherited vector; require TestBuildScriptMalformedResources to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-04",
              "digest": "sha256:b24a8d769b651a5922d46c092ddeda40b8af41bcb7d0a1479342dbe6da854908",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/inherited.go,omit,failed,1,TestBuildScriptMalformedResources,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/runbinary,^TestBuildScriptMalformedResources$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,519\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/runbinary,fail,17,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/runbinary,TestBuildScriptMalformedResources,\"lifetime_resources_test.go:88: Own under a malformed vector = selected Bench executable \\\"/dev/shm/XIBFPO/t/TestBuildScriptMalformedResources258308181/003/bench-run-2517604797/bench\\\" is unavailable: lstat /dev/shm/XIBFPO/t/TestBuildScriptMal… (347 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "rs-c4-04-runner-callers",
          "performer": "ft446-author-04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8e8882a8a7fff5820dc1d01f98d70398480d9f2d",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft446-author-04",
            "digest": "sha256:b1c1f2e8d136b8c33df64b2f374e905bd27889f250094fc0e60aad558132bac1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,1,TestBoundsFixedShutdown,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/bounds,^TestBoundsFixedShutdown$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,319\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,fail,3014,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/bounds,TestBoundsFixedShutdown,\"lifetime_test.go:80: canceled runner status = timeout: context canceled\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "398691f9cf8dbf2bb1d60c747bc2f7585dc57b6a",
          "requirement": "runner-callers",
          "command": "bench test --package ./internal/bounds --package ./internal/gate --package ./internal/preprelease --package ./internal/releasepreflight --package ./internal/testreport",
          "exit_code": 0,
          "probe": {
            "mutation": "make the bounds cancel request answer a deadline for a cancelled parent; require TestBoundsFixedShutdown to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft446-author-04",
              "digest": "sha256:b1c1f2e8d136b8c33df64b2f374e905bd27889f250094fc0e60aad558132bac1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/bounds/bounds.go,swap,failed,1,TestBoundsFixedShutdown,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/bounds,^TestBoundsFixedShutdown$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,319\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,fail,3014,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/bounds,TestBoundsFixedShutdown,\"lifetime_test.go:80: canceled runner status = timeout: context canceled\",1\nskips[0]{package,test,reason}:\n"
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
    },
    {
      "from": "sha256:19ee2fb6487fbd0b39017eb392ac3d31bdab79d1b14798b737b4797eb880470f",
      "to": "sha256:c10f32d8ce7a14d274b01166ca9d0f66f743e88e7d0a5fd13afa01bf51ff5af9",
      "chunk_ids": {
        "RS-C1": [
          "RS-C1"
        ],
        "RS-C2": [
          "RS-C2"
        ],
        "RS-C3": [
          "RS-C3"
        ]
      }
    },
    {
      "from": "sha256:c10f32d8ce7a14d274b01166ca9d0f66f743e88e7d0a5fd13afa01bf51ff5af9",
      "to": "sha256:ce20aca3c19ac69f01f324f724ef82aed6df6f7437403fb7fbfdc15715ccad84",
      "chunk_ids": {
        "RS-C1": [
          "RS-C1"
        ],
        "RS-C2": [
          "RS-C2"
        ],
        "RS-C3": [
          "RS-C3"
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

## RS-C3 author witness reds

The ticket 03 author ft446-author-03 reported nine witness reds, and each probe reported a restored subject:

1. RS9: with a kind predicate that is always false, TestDiscardRetired was red.
2. RS10: with no resource-absence condition, TestDiscardRetired was red.
3. RS30: with no record kind comparison, TestDiscardRetired was red on a binary store under the cache prefix.
4. RS30: with no store-name prefix filter, TestDiscardRetired was red on a cache store under the binary prefix.
5. RS12: with an unreadable record that returns its error, TestDiscardRetiredFailsClosed was red.
6. RS13: with the deletion error dropped, TestDiscardRetiredFailsClosed was red.
7. RS13: with a root listing error dropped, TestDiscardRetiredFailsClosed was red.
8. RS11: with a check call at the contended return, TestReleaseRunsTheRetiredCheck was red.
9. RS28: with the check call before the retirement, TestReleaseRunsTheRetiredCheck was red.

The author reported one duplicated fixture outside its fence. The new helper in the directory tests repeats two lines of the protection fixture helper.

## RS-C3 review pickup

Unified review by ft446-review-rs-c3 on sonnet at high effort over the pair d7641b2..2220bf0. Raw findings: 1. Repair targets: 1. Repair cycles used: 0 of 2 before the repair below.

### Standards

Findings: 1. Worst issue: RS-C3-S1.

- RS-C3-S1 (auto-fix, confidence 7): the new directory-test helper repeats the retire-then-delete pair of the protection fixture helper. It also repeats the store base and protocol setup. The reviewer returned this item as advice that needs a widen. The orchestrator accepts it as a finding, because the working agreement keeps each fixture harness single-sourced. The fixture file is in the directory of the ticket 03 fence, so the widen takes a plan commit and no stop. The fix is test-only and goes to the fresh repair session ft446-repair-03-r1.

### Spec

Findings: 0. Worst issue: none.

### Coverage

Findings: 0. Worst issue: none.

### Advice

- The unlisted-root subtest fails through the root mode check, not through the directory listing. That listing branch is hard to reach.

## RS-C3 repair r1 and confirming round

The repair session ft446-repair-03-r1 made a test-only repair at e6cde66. One fixture helper now owns the retire-then-delete fact, and one helper owns the private store base and protocol setup. The repair added no test expectation, so it reported no witness red.

The orchestrator probe made the shared helper retire a store without a resource deletion. TestDiscardRetired was red, and the subject was restored. The coordinator confirmation refused, because the ticket widen puts a ticket file in the chunk delta. A unified confirming round by ft446-confirm-rs-c3 then closed RS-C3-S1 with no new finding. Repair cycles used: 1 of 2.

The test that grades a protected removal still builds its own protocol with another kind set. That code predates this chunk, so it stays outside the repair.

## RS-C4 ticket 04 author witness reds

The ticket 04 author ft446-author-04 reported three witness reds, and each probe reported a restored subject:

1. RS16: when the launch used the installed protocol for every registration, TestRunInheritedRegistersOnLocalProtocol was red.
2. RS23: when the launch ignored an unreadable inherited vector, TestBuildScriptMalformedResources was red.
3. RS17: when the launch had no nil-protocol default, TestProcessLifetimeCacheSecondProcessClean was red.

The bounds runner passes an empty registration literal, as the detached launch does under D4. The orchestrator reads the acceptance item as a ban on a hand-built inherited registration, and flags that reading for reviewer veto.
