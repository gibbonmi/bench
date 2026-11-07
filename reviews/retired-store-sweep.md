# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/retired-store-sweep/spec.md",
  "plan_digest": "sha256:2d22011a8556aa7de0f62d1d5720662c09790fa7124e337c6d6df82185339255",
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
      "tip": "bef1597e66827888f6fbdd964ac136d50f035aed",
      "plan_digest": "sha256:2d22011a8556aa7de0f62d1d5720662c09790fa7124e337c6d6df82185339255",
      "source_digest": "910b248be12293a5b95f20d1c24d6d8e144f9156",
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
