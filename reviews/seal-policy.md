# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/seal-policy/spec.md",
  "plan_digest": "sha256:8358c5964ffe01801052a999f35c28ac46c0a5403eeeabd0dda1ba73f89033ae",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SP-C1",
      "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
      "tip": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
      "plan_digest": "sha256:8358c5964ffe01801052a999f35c28ac46c0a5403eeeabd0dda1ba73f89033ae",
      "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
      "acceptance_rows": [
        "SP1",
        "SP2",
        "SP3",
        "SP4",
        "SP5",
        "SP6",
        "SP7",
        "SP8",
        "SP9",
        "SP10",
        "SP34"
      ],
      "verification": [
        {
          "id": "ft482-01-release-seal-order",
          "performer": "ft482-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-author-01",
            "digest": "sha256:57f8ccbbcaed6a7e073e7021b5577677570a6f211bfe206a41a869ae7e4b2b5a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseKeepsALiveUsersGenerationOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,all,passed,51\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,469\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,275,51\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseKeepsALiveUsersGenerationOpen,\"protection_test.go:167: a new hold after the release = {Store: ResourceID: Kind: ResourcePath: Binding:} (protection: the resource generation is sealed), want generation 29745d0d92d67d9226991b3c4fd3df29\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
          "requirement": "release-seal-order",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "seal before the grade in the release observation",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-author-01",
              "digest": "sha256:57f8ccbbcaed6a7e073e7021b5577677570a6f211bfe206a41a869ae7e4b2b5a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseKeepsALiveUsersGenerationOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,all,passed,51\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,469\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,275,51\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseKeepsALiveUsersGenerationOpen,\"protection_test.go:167: a new hold after the release = {Store: ResourceID: Kind: ResourcePath: Binding:} (protection: the resource generation is sealed), want generation 29745d0d92d67d9226991b3c4fd3df29\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft482-01-protected-release-write",
          "performer": "ft482-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-author-01",
            "digest": "sha256:1f44434fdf230ab892247a2bf86f91852181cb71e05fc58c61b55cee713806e8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseObservationKeepsAnEarlierSeal,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,81\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,62,81\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestReleaseObservationKeepsAnEarlierSeal,\"recovery_test.go:158: Prepare after the release observation = \u003cnil>, want protection: the resource generation is sealed\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "source_commit": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
          "requirement": "protected-release-write",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "write admission open on a protected release",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-author-01",
              "digest": "sha256:1f44434fdf230ab892247a2bf86f91852181cb71e05fc58c61b55cee713806e8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseObservationKeepsAnEarlierSeal,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,81\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,62,81\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestReleaseObservationKeepsAnEarlierSeal,\"recovery_test.go:158: Prepare after the release observation = \u003cnil>, want protection: the resource generation is sealed\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ft482-c1-standards",
          "performer": "ft482-review-c1-standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "agent:ft482-review-c1-standards",
            "digest": "sha256:3a4911a26f9fcaec765e36b8af32bb1fd619688934dce0c798c8f88d7a92d767",
            "excerpt": "Standards axis, SP-C1, evidence sha256:457286fc5169b0d8d01a67f24b07e0aab9de21b1c5180d45152fc38743762bbf (current=true, tip 0da3fab45). Count: 3, all judgment calls. Worst: SP10 test repeats the focused-run harness body (confidence 5).\nSP-C1-S1 internal/systemtest/process_lifetime_callers_test.go:242-260 TestFocusedTestAfterALeakedDescendant re-derives the per-run body of TestProcessLifetimeFocusedDifferential (:215-235) and picks shapes by bare index; confidence 5; auto-fix: extract one focused-run helper.\nSP-C1-S2 internal/subprocess/protection/store.go:93-96 Reopen comment narrates the cache-release change and is vague; confidence 4; auto-fix: state only the sealing cases.\nSP-C1-S3 internal/gocache/protection_test.go:154-157 and :224-226 inline Prepare twice with no pending-user helper; confidence 3; no-op.\n"
          },
          "axis": "Standards",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
          "finding_ids": [
            "SP-C1-S1",
            "SP-C1-S2",
            "SP-C1-S3"
          ],
          "supersedes": []
        },
        {
          "id": "ft482-c1-standards-accepted",
          "performer": "ft482-review-c1-standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "agent:ft482-review-c1-standards",
            "digest": "sha256:3a4911a26f9fcaec765e36b8af32bb1fd619688934dce0c798c8f88d7a92d767",
            "excerpt": "Standards axis, SP-C1, evidence sha256:457286fc5169b0d8d01a67f24b07e0aab9de21b1c5180d45152fc38743762bbf (current=true, tip 0da3fab45). Count: 3, all judgment calls. Worst: SP10 test repeats the focused-run harness body (confidence 5).\nSP-C1-S1 internal/systemtest/process_lifetime_callers_test.go:242-260 TestFocusedTestAfterALeakedDescendant re-derives the per-run body of TestProcessLifetimeFocusedDifferential (:215-235) and picks shapes by bare index; confidence 5; auto-fix: extract one focused-run helper.\nSP-C1-S2 internal/subprocess/protection/store.go:93-96 Reopen comment narrates the cache-release change and is vague; confidence 4; auto-fix: state only the sealing cases.\nSP-C1-S3 internal/gocache/protection_test.go:154-157 and :224-226 inline Prepare twice with no pending-user helper; confidence 3; no-op.\n"
          },
          "axis": "Standards",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
          "finding_ids": [
            "SP-C1-S1",
            "SP-C1-S2"
          ],
          "supersedes": [
            "ft482-c1-standards"
          ]
        },
        {
          "id": "ft482-c1-spec",
          "performer": "ft482-review-c1-spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-review-c1-spec",
            "digest": "sha256:8bcfc205b4d3f064a7dd4245b115ef332d74c11f7924189ce92d8cc7114a8956",
            "excerpt": "Spec axis, SP-C1, evidence sha256:457286fc5169b0d8d01a67f24b07e0aab9de21b1c5180d45152fc38743762bbf (current=true, head 0da3fab45). Count: 0. Worst: none.\nAudited SP1-SP10 and SP34 against the spec and the code at 9b21102ef: one shared observe body with sealFirst as the only difference, one grade loop and one tail, no open-admission write, only heldLock.retire switches, every other SealAndObserve caller and its signature kept, the three named comments changed, ADR 0030 lines 26-28 state the release policy (SP34 met).\nDivergences accepted: the deadUser helper (confidence 8) and the ErrSealed.Error() comparison (confidence 8).\n"
          },
          "axis": "Spec",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ft482-c1-coverage",
          "performer": "ft482-review-c1-coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "954cdf78617349760a2035f9e137c8ff35a1a043",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-review-c1-coverage",
            "digest": "sha256:f4cb94690f51ae8dd8b68ede702fbb67ffd9929d294ec5f533246687f0afc67a",
            "excerpt": "Coverage axis, SP-C1, evidence sha256:457286fc5169b0d8d01a67f24b07e0aab9de21b1c5180d45152fc38743762bbf (current=true, tree clean after probes). Count: 1, low severity. Worst: an unpinned comment claim.\nSP-C1-C1 internal/subprocess/protection/recovery.go:94-100 the seal guard `if r.Admission != admissionOpen { return nil }` is unpinned; probe swap to false: silent (81 tests pass, subject restored); confidence 4; no-op (a repeated durable write, no observable state).\nProbe: drop the result.Users append: bit (TestReleaseObservationNamesTheLiveUser red).\nWalked: every production branch of observe is pinned by SP2-SP8; the SP10 test runs beside a live sleep descendant, and its mark file proves the go test ran.\n"
          },
          "axis": "Coverage",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "9b21102ef8e2e2970b3bdd4eab12c7108ddc9e00",
          "finding_ids": [],
          "supersedes": []
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

## SP-C1 author witness reds

The ticket 1 author ft482-author-01 reported these witness reds. Each red came before the production edit, through `bench test --package <package> --run <regex>`:

1. SP1: TestReleaseKeepsALiveUsersGenerationOpen was red on the base code, because the new hold after the release read `the resource generation is sealed` (`internal/gocache/protection_test.go:167`).
2. SP2 and SP5: with a stub that only called SealAndObserve, TestReleaseObservationLeavesAdmissionOpen and TestReleaseObservationLeavesAdmissionOpenWhenUncertain were red, because Prepare read `the resource generation is sealed` (`recovery_test.go:108` and `:116`).
3. SP3 and SP4: with a stub that graded no user, TestReleaseObservationNamesTheLiveUser and TestReleaseObservationRetiresASafeGeneration were red on an empty disposition.
4. SP6: with a stub that reopened admission, TestReleaseObservationKeepsAnEarlierSeal was red, because Prepare returned nil.
5. SP7: with a stub that removed the users before the seal, TestReleaseObservationSealFailure was red, because the store held no user after the failed seal.
6. SP10: with the cache retire on SealAndObserve again, TestFocusedTestAfterALeakedDescendant was red, because the second focused run refused with `the resource generation is sealed` (`process_lifetime_callers_test.go:256`).

The two plan probes bit, and each probe reported a restored subject.

## SP-C1 review

### Standards

Count: 2 accepted, 1 no-op. Worst issue: SP-C1-S1.

- SP-C1-S1 (auto-fix, confidence 5): `internal/systemtest/process_lifetime_callers_test.go:242-260`. TestFocusedTestAfterALeakedDescendant repeats the focused-run body of TestProcessLifetimeFocusedDifferential (`:215-235`), and it selects shapes by bare index. Extract one focused-run helper that both tests call.
- SP-C1-S2 (auto-fix, confidence 4): `internal/subprocess/protection/store.go:93-96`. The Reopen comment tells how the cache release changed, and its first clause is vague. State only the cases that seal a generation.
- SP-C1-S3 (no-op, confidence 3): `internal/gocache/protection_test.go:154-157` and `:224-226` call Prepare inline. This is incidental text, not duplicated knowledge. The first Standards entry listed it, and the superseding entry drops it.

### Spec

Count: 0. Worst issue: none. The axis audited SP1 to SP10 and SP34, and it accepted the deadUser helper and the ErrSealed comparison.

### Coverage

Count: 0 accepted, 1 no-op. Worst issue: none.

- SP-C1-C1 (no-op, confidence 4): `internal/subprocess/protection/recovery.go:94-100`. A probe that removes the seal guard is silent. The guard only skips a second durable write of the same sealed state, so no test can observe it.
