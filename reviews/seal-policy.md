# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/seal-policy/spec.md",
  "plan_digest": "sha256:3c1cee34c7c9a3075fef21a5195d02542c760600dc156ca41e77c174f773dea5",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SP-C1",
      "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
      "tip": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
      "plan_digest": "sha256:60a5b3e8633ad83402c59bbeb24ea24d1256f62e61f95abbe3b299c14fa4d8fd",
      "source_digest": "14177f4d8269826387407edd2d0c969b34829718",
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
        },
        {
          "id": "ft482-01-r1-release-seal-order",
          "performer": "ft482-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "14177f4d8269826387407edd2d0c969b34829718",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-repair-01-r1",
            "digest": "sha256:e953062b98b5990198e05594ae6ced831979d92ff915158e3fb0f7faa315dd3a",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseKeepsALiveUsersGenerationOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,all,passed,51\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,469\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,278,51\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseKeepsALiveUsersGenerationOpen,\"protection_test.go:167: a new hold after the release = {Store: ResourceID: Kind: ResourcePath: Binding:} (protection: the resource generation is sealed), want generation b354e9f8b68d85ff0318d6fdc361b9cf\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "requirement": "release-seal-order",
          "command": "bench test --package ./internal/gocache",
          "exit_code": 0,
          "probe": {
            "mutation": "seal before the grade in the release observation",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-repair-01-r1",
              "digest": "sha256:e953062b98b5990198e05594ae6ced831979d92ff915158e3fb0f7faa315dd3a",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseKeepsALiveUsersGenerationOpen,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gocache,all,passed,51\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,469\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gocache,fail,278,51\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/gocache,TestReleaseKeepsALiveUsersGenerationOpen,\"protection_test.go:167: a new hold after the release = {Store: ResourceID: Kind: ResourcePath: Binding:} (protection: the resource generation is sealed), want generation b354e9f8b68d85ff0318d6fdc361b9cf\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft482-01-r1-protected-release-write",
          "performer": "ft482-repair-01-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "14177f4d8269826387407edd2d0c969b34829718",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-repair-01-r1",
            "digest": "sha256:716fa6979780f64aabffebd62cd0ca71c0929f402899f4be4dfd7fd198546373",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseObservationKeepsAnEarlierSeal,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,81\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,61,81\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestReleaseObservationKeepsAnEarlierSeal,\"recovery_test.go:158: Prepare after the release observation = \u003cnil>, want protection: the resource generation is sealed\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "source_commit": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "requirement": "protected-release-write",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "write admission open on a protected release",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-repair-01-r1",
              "digest": "sha256:716fa6979780f64aabffebd62cd0ca71c0929f402899f4be4dfd7fd198546373",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/recovery.go,swap,failed,1,TestReleaseObservationKeepsAnEarlierSeal,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,81\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,61,81\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestReleaseObservationKeepsAnEarlierSeal,\"recovery_test.go:158: Prepare after the release observation = \u003cnil>, want protection: the resource generation is sealed\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
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
        },
        {
          "id": "ft482-c1-standards-r1",
          "performer": "ft482-confirm-c1-standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "14177f4d8269826387407edd2d0c969b34829718",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-confirm-c1-standards",
            "digest": "sha256:ac01ba0a97a52bf831c5e7cc27cc820a1a6bb5cfe9b756612cadc54162409063",
            "excerpt": "Standards axis, SP-C1 confirming round, evidence sha256:61c93d51f1beb6c825ceea6cd6720e80dbfdb61638add330c533a2d3b45370ed (current, head 61d7a6b8a). Repair delta 9b21102ef..5698f2ebc.\nSP-C1-S1 closed: focusedRun (internal/systemtest/process_lifetime_callers_test.go:202) owns the run body; both tests call it (:235, :256) and keep their assertions.\nSP-C1-S2 closed: the Reopen comment (internal/subprocess/protection/store.go:93-96) states in the present tense why Reopen exists.\nNew findings: none. Advice: the same env-append pattern near line 305 is outside the delta.\n"
          },
          "axis": "Standards",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "finding_ids": [],
          "supersedes": [
            "ft482-c1-standards-accepted"
          ]
        },
        {
          "id": "ft482-c1-spec-r1",
          "performer": "ft482-confirm-c1-spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "14177f4d8269826387407edd2d0c969b34829718",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-confirm-c1-spec",
            "digest": "sha256:2c999d64c743c2dc80663dd33a05c7acd531b7bceb2f24067482f1b3c5792c7a",
            "excerpt": "Spec axis, SP-C1 confirming round, evidence sha256:61c93d51f1beb6c825ceea6cd6720e80dbfdb61638add330c533a2d3b45370ed (current, tip 61d7a6b8a). Repair delta 9b21102ef..5698f2ebc. Reaffirmed, no new finding.\nSP10: TestFocusedTestAfterALeakedDescendant still runs two bench test calls on one HOME, leaked shape first; focusedRun proves each go test ran; the test fails on ErrSealed text.\nThe Reopen comment (store.go:93-97) agrees with the Implementation decisions \"The release observation\".\n"
          },
          "axis": "Spec",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "finding_ids": [],
          "supersedes": [
            "ft482-c1-spec"
          ]
        },
        {
          "id": "ft482-c1-coverage-r1",
          "performer": "ft482-confirm-c1-coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "14177f4d8269826387407edd2d0c969b34829718",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-confirm-c1-coverage",
            "digest": "sha256:03858b00b45fa8f886efa2aff2092abb0109066c9cab6b2686d9744370b997c3",
            "excerpt": "Coverage axis, SP-C1 confirming round, evidence sha256:61c93d51f1beb6c825ceea6cd6720e80dbfdb61638add330c533a2d3b45370ed (current, tree clean). Repair delta 9b21102ef..5698f2ebc. Reaffirmed.\nNo lost or weakened assertion: the mark-file check and groupOf moved into focusedRun; the exit-code check, the group checks, the shared HOME, the shape order, and the sealed-generation check stay. The coordinator probe on focusedRun bit (6 failed tests).\nResidual (no-op, confidence 9): the groupOf fatal now fires before the sealed check; ordering only.\n"
          },
          "axis": "Coverage",
          "base": "ea40264fb165f257f306bb1de5e89dc0b951ccd7",
          "tip": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "finding_ids": [],
          "supersedes": [
            "ft482-c1-coverage"
          ]
        }
      ]
    },
    {
      "id": "SP-C2",
      "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
      "tip": "c2abc657d8eefe515a1a540861945575fcdb6b50",
      "plan_digest": "sha256:3c1cee34c7c9a3075fef21a5195d02542c760600dc156ca41e77c174f773dea5",
      "source_digest": "6fbdbd9e7c00fa3c879a29dead17e6fef11ce98b",
      "acceptance_rows": [
        "SP11",
        "SP12",
        "SP13",
        "SP14",
        "SP15",
        "SP16",
        "SP17",
        "SP18",
        "SP19",
        "SP20",
        "SP33",
        "SP35"
      ],
      "verification": [
        {
          "id": "ft482-02-zombie-leader-only",
          "performer": "ft482-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d9a8090801e0dc2eb6d941dd50e6f144cecd36eb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-author-02",
            "digest": "sha256:5105e6767478e9938138e7bfeb6d5e1132a37aa6988ebaf8e9a1b38b30a923c7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/root_linux.go,swap,failed,3,TestProcLiveTaskBesideAZombieLeader,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,92\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1135\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,97,92\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcLiveTaskBesideAZombieLeader,\"observe_linux_test.go:121: a process with a sleeping task beside its zombie leader counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/malformed_task_record,\"observe_linux_test.go:150: target 7 counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/unreadable_task_record,\"observe_linux_test.go:150: target 7 counts as absent, want present\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "source_commit": "ba078b819de64b9b85a0c257fd3ad282617c1352",
          "requirement": "zombie-leader-only",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "read only the leader task state",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-author-02",
              "digest": "sha256:5105e6767478e9938138e7bfeb6d5e1132a37aa6988ebaf8e9a1b38b30a923c7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/root_linux.go,swap,failed,3,TestProcLiveTaskBesideAZombieLeader,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,92\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1135\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,97,92\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcLiveTaskBesideAZombieLeader,\"observe_linux_test.go:121: a process with a sleeping task beside its zombie leader counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/malformed_task_record,\"observe_linux_test.go:150: target 7 counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/unreadable_task_record,\"observe_linux_test.go:150: target 7 counts as absent, want present\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
            }
          }
        },
        {
          "id": "ft482-02-known-issues-row",
          "performer": "ft482-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "77a20dabbcac4ef1f232452570fb0afc20844b9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-author-02",
            "digest": "sha256:a54e7065c7be4effc05a7712c10bac852517baa1e1ae9a6d792fd6d7ebb7c6f8",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,docs/known-issues.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2342,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: cache release seal: the known-issues file kept a row that FT482 retires\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "1bc4083c06e409ac67ff33b5016e4a6a49f39a5b",
          "requirement": "known-issues-row",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "restore the FT482 known-issues row",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-author-02",
              "digest": "sha256:a54e7065c7be4effc05a7712c10bac852517baa1e1ae9a6d792fd6d7ebb7c6f8",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,docs/known-issues.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2342,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: cache release seal: the known-issues file kept a row that FT482 retires\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft482-02-r1-zombie-leader-only",
          "performer": "ft482-repair-02-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6fbdbd9e7c00fa3c879a29dead17e6fef11ce98b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-repair-02-r1",
            "digest": "sha256:675522a28a111189d8640f2304a4c15e5527c822d32c1f8e7ae6fa0c925f48d6",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/root_linux.go,swap,failed,3,TestProcLiveTaskBesideAZombieLeader,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,96\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1135\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,99,96\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcLiveTaskBesideAZombieLeader,\"observe_linux_test.go:121: a process with a sleeping task beside its zombie leader counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/malformed_task_record,\"observe_linux_test.go:155: target 7 counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/unreadable_task_record,\"observe_linux_test.go:155: target 7 counts as absent, want present\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
          },
          "source_commit": "c2abc657d8eefe515a1a540861945575fcdb6b50",
          "requirement": "zombie-leader-only",
          "command": "bench test --package ./internal/subprocess/protection",
          "exit_code": 0,
          "probe": {
            "mutation": "read only the leader task state",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-repair-02-r1",
              "digest": "sha256:675522a28a111189d8640f2304a4c15e5527c822d32c1f8e7ae6fa0c925f48d6",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/subprocess/protection/root_linux.go,swap,failed,3,TestProcLiveTaskBesideAZombieLeader,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/subprocess/protection,all,passed,96\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1135\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,fail,99,96\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcLiveTaskBesideAZombieLeader,\"observe_linux_test.go:121: a process with a sleeping task beside its zombie leader counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/malformed_task_record,\"observe_linux_test.go:155: target 7 counts as absent, want present\",1\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestProcUnreadableStatKeepsPresent/unreadable_task_record,\"observe_linux_test.go:155: target 7 counts as absent, want present\",1\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/boot_root,\"capability: privilege: only a root user can give a directory another owner\"\n  github.com/gibbonmi/bench/internal/subprocess/protection,TestStoreRootRefusal/other_owner/user_root,\"capability: privilege: only a root user can give a directory another owner\"\n"
            }
          }
        },
        {
          "id": "ft482-02-r1-known-issues-row",
          "performer": "ft482-repair-02-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "6fbdbd9e7c00fa3c879a29dead17e6fef11ce98b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-repair-02-r1",
            "digest": "sha256:681e0c1aea0a7fe0be17c9e88a2d36676279e8877c2c8e060856c9dc28d1a23f",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,docs/known-issues.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2355,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: cache release seal: the known-issues file kept a row that FT482 retires\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "c2abc657d8eefe515a1a540861945575fcdb6b50",
          "requirement": "known-issues-row",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "restore the FT482 known-issues row",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agent:ft482-repair-02-r1",
              "digest": "sha256:681e0c1aea0a7fe0be17c9e88a2d36676279e8877c2c8e060856c9dc28d1a23f",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,docs/known-issues.md,insert,failed,1,TestRootConformance,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,435\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,2355,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:30: gate: cache release seal: the known-issues file kept a row that FT482 retires\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ft482-c2-standards",
          "performer": "ft482-review-c2-standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "77a20dabbcac4ef1f232452570fb0afc20844b9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-review-c2-standards",
            "digest": "sha256:7c54956830bf5fd095bdcd0890956223bf35eee7fef334f4e94b9c36e2951cd9",
            "excerpt": "Standards axis, SP-C2, evidence sha256:ae45d5731fb62027be27d284966f527341b47b8be044c4d49fc056919e918953 (current, head 83a9121f0). Pair 5698f2ebc..ba078b819. Raw count 6, all judgment calls, no hard violation; accepted 0.\n1 reviews/seal-policy.md: the ticket commit wrote the author witness reds (conf 4). No-op: the orchestrator charge directed the author to write them; the record is the run's evidence path.\n2 specs/seal-policy/spec.md:180-186 coverage-cell citations (conf 3). No-op: the ticket-commit citation rule requires them.\n3 internal/subprocess/protection/observe.go:29-31 Probe doc restates code (conf 3). Optional advice.\n4 root_darwin.go:11-13 and root_other.go:18-20 repeat systemProbe (conf 3). No-op: honest repetition under the file budget.\n5 CHANGELOG.md:12 sentence of about 27 words (conf 3). No-op: the prose lane passed.\n6 observe_linux_test.go spec-row tags in comments (conf 3). No-op: repo convention.\nClean: one owner of the stat parse (statFields) and of the zombie rule; the test export exposes only procAbsent; independent expectations carry demonstrated reds.\n"
          },
          "axis": "Standards",
          "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "tip": "1bc4083c06e409ac67ff33b5016e4a6a49f39a5b",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ft482-c2-spec",
          "performer": "ft482-review-c2-spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "77a20dabbcac4ef1f232452570fb0afc20844b9b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-review-c2-spec",
            "digest": "sha256:c5777013ee64f8d6b143e3cac197dd8d6aefe2a175df95b72657ad354b211310",
            "excerpt": "Spec axis, SP-C2, evidence sha256:ae45d5731fb62027be27d284966f527341b47b8be044c4d49fc056919e918953 (current, head 83a9121f0). Pair 5698f2ebc..ba078b819. Count 0, worst none.\nEvery \"The zombie rule\" and \"The retirement\" decision is met; SP11-SP20 and SP33 map to named tests; SP35 met (ADR 0030 Platform scope states the Linux zombie rule).\nDivergence (a), an empty task or member set stays present: not a contradiction; the spec's unreadable-keeps-present and ENOENT sentences cover it; no-op, confidence 8. Spec wording could say \"finds at least one task\".\nDivergence (b), the actor_pid rewrite through plant: no-op, confidence 8. Divergence (c), mutation reds before the commit: no-op, confidence 7.\n"
          },
          "axis": "Spec",
          "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "tip": "1bc4083c06e409ac67ff33b5016e4a6a49f39a5b",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ft482-c2-coverage",
          "performer": "ft482-review-c2-coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "77a20dabbcac4ef1f232452570fb0afc20844b9b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "agent:ft482-review-c2-coverage",
            "digest": "sha256:20f93d0fc86f494525b56d13b876976d5f964e0ddc0cd63ef8f067c31ebfe611",
            "excerpt": "Coverage axis, SP-C2, evidence sha256:ae45d5731fb62027be27d284966f527341b47b8be044c4d49fc056919e918953 (current, tree clean after probes). Pair 5698f2ebc..ba078b819. Count 3, worst SP-C2-C1.\nSP-C2-C1 internal/subprocess/protection/root_linux.go:83 the empty-group case (a group with no member found stays present) has no test; probe swap \"return found\" to \"return true\": silent (92/92 pass, restored); confidence 8; auto-fix: a negative-target row whose root holds no member.\nSP-C2-C2 root_linux.go:81 no group-level witness for a member with a stat record and no tasks; probe found = true: invalid (unused variable); confidence 6; auto-fix: one group row with a stat-only member.\nSP-C2-C3 root_linux.go:~89 and ~110 the non-ENOENT ReadDir error branch and the len(fields) \u003c 3 guard: combined probe silent; a record such as \"7 (a) Z\" is untested; confidence 5; auto-fix.\nNo gap: pgrp field (fields[2] to fields[3] bit), ENOENT vs unreadable, pid-vs-group dispatch, real-zombie wait.\n"
          },
          "axis": "Coverage",
          "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "tip": "1bc4083c06e409ac67ff33b5016e4a6a49f39a5b",
          "finding_ids": [
            "SP-C2-C1",
            "SP-C2-C2",
            "SP-C2-C3"
          ],
          "supersedes": []
        },
        {
          "id": "ft482-c2-coverage-r1",
          "performer": "ft482-confirm-c2-coverage",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6fbdbd9e7c00fa3c879a29dead17e6fef11ce98b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-confirm-c2-coverage",
            "digest": "sha256:57c5a5f6729856ca03f74378772db3b0189702f2b062cfe9c327b3d2090c192e",
            "excerpt": "Coverage axis, SP-C2 confirming round, evidence sha256:bc8df8f35e163e9047e76e76ac98410ff5ca9567c7e203d01cc23fb60b255b73 (current, tree clean). Repair delta 1bc4083c0..c2abc657d.\nSP-C2-C1 closed: the no-member row probes group -7 with only pid 9 in group 3; root_linux.go:83 keeps it present; the coordinator probe (return true) bit.\nSP-C2-C2 closed: the member-with-no-task row keeps found false; the || dead mutant is red.\nSP-C2-C3 closed: the short task record \"7 (a) Z\" fails the guard at line 117; the task-as-file row makes ReadDir return ENOTDIR at line 92; each mutant is red.\nNew findings: none. Confidence 8.\n"
          },
          "axis": "Coverage",
          "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "tip": "c2abc657d8eefe515a1a540861945575fcdb6b50",
          "finding_ids": [],
          "supersedes": [
            "ft482-c2-coverage"
          ]
        },
        {
          "id": "ft482-c2-spec-r1",
          "performer": "ft482-confirm-c2-spec",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6fbdbd9e7c00fa3c879a29dead17e6fef11ce98b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-confirm-c2-spec",
            "digest": "sha256:c8c8dda25eb6a3758385fa60ec97cdc4d00249b28fc1f52f661718970eb17a29",
            "excerpt": "Spec axis, SP-C2 confirming round, evidence sha256:bc8df8f35e163e9047e76e76ac98410ff5ca9567c7e203d01cc23fb60b255b73 (current, head 34b2cdc1b). Repair delta 1bc4083c0..c2abc657d. Reaffirmed, no new finding.\nThe spec prose at spec.md:93 matches root_linux.go (a pid is absent only when dead and found; a group only when found and no member is present). The plan hunk adds only the repair assignment; no acceptance row changes. The four new rows agree with the spec and the code.\n"
          },
          "axis": "Spec",
          "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "tip": "c2abc657d8eefe515a1a540861945575fcdb6b50",
          "finding_ids": [],
          "supersedes": [
            "ft482-c2-spec"
          ]
        },
        {
          "id": "ft482-c2-standards-r1",
          "performer": "ft482-confirm-c2-standards",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6fbdbd9e7c00fa3c879a29dead17e6fef11ce98b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agent:ft482-confirm-c2-standards",
            "digest": "sha256:051e788d8b708d858836f400eebd16eadeec9d52858bc80fad1e356c1b4da110",
            "excerpt": "Standards axis, SP-C2 confirming round, evidence sha256:bc8df8f35e163e9047e76e76ac98410ff5ca9567c7e203d01cc23fb60b255b73 (current). Repair delta 1bc4083c0..c2abc657d.\nThe four new rows extend the existing table with no second harness; the spec prose is plain STE; the plan entry follows the assignment shape.\nOne prose finding (confidence 6, auto-fix): specs/seal-policy/tickets/2-count-a-zombie-user-as-dead.md:15 still states the old empty-set wording after the spec correction at spec.md:93. Route: a prose-only spec correction in the next plan commit before the ticket 3 dispatch, under the build rule for review prose corrections. No code or test change.\n"
          },
          "axis": "Standards",
          "base": "5698f2ebcaa90ef4aace6323c4c5bc4447fb161e",
          "tip": "c2abc657d8eefe515a1a540861945575fcdb6b50",
          "finding_ids": [],
          "supersedes": [
            "ft482-c2-standards"
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
      "from": "sha256:8358c5964ffe01801052a999f35c28ac46c0a5403eeeabd0dda1ba73f89033ae",
      "to": "sha256:60a5b3e8633ad83402c59bbeb24ea24d1256f62e61f95abbe3b299c14fa4d8fd",
      "chunk_ids": {
        "SP-C1": [
          "SP-C1"
        ]
      }
    },
    {
      "from": "sha256:60a5b3e8633ad83402c59bbeb24ea24d1256f62e61f95abbe3b299c14fa4d8fd",
      "to": "sha256:dbcfe96db64f0ed22730bf2af883fb1696813bffcdba233c62695e844403c29e",
      "chunk_ids": {
        "SP-C1": [
          "SP-C1"
        ]
      }
    },
    {
      "from": "sha256:dbcfe96db64f0ed22730bf2af883fb1696813bffcdba233c62695e844403c29e",
      "to": "sha256:3c1cee34c7c9a3075fef21a5195d02542c760600dc156ca41e77c174f773dea5",
      "chunk_ids": {
        "SP-C1": [
          "SP-C1"
        ],
        "SP-C2": [
          "SP-C2"
        ]
      }
    }
  ]
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

### SP-C1 repair state

Repair round r1 (ft482-repair-01-r1, opus/high) closed SP-C1-S1 and SP-C1-S2 at 5698f2ebc. SP-C1 used 1 of its 2 repair cycles. The coordinator probe on focusedRun bit with 6 failed tests. The confirming round closed both findings, and Spec and Coverage reaffirmed.

Advice: the env-append pattern near line 305 of `internal/systemtest/process_lifetime_callers_test.go` is outside the SP-C1 delta.

## SP-C2 author witness reds

The ticket 2 author ft482-author-02 observed these witness reds. Each base red came before the production edit, through `bench test --package ./internal/subprocess/protection --run <regex>`. Each mutation red came through `bench probe` on `internal/subprocess/protection/root_linux.go`, and each probe restored its subject.

1. SP11: with the default observer still on `syscall.Kill`, TestReleaseObservationRetiresAZombieUser was red, because the release observation of the zombie user returned protected (`observe_linux_test.go:47`).
2. SP12: with the default observer still on `syscall.Kill`, TestZombiePendingActorIsReleased was red, because the release observation of the zombie actor returned protected (`observe_linux_test.go:69`).
3. SP13: with a group rule that read only the leader process, TestZombieLeaderWithALiveMemberIsProtected was red, because the release observation of the live group returned safe (`observe_linux_test.go:80`).
4. SP14: with a process rule that read only the leader task state, TestProcLiveTaskBesideAZombieLeader was red, because the process counted as absent (`observe_linux_test.go:121`).
5. SP15: with a parse that split at the first `)`, TestProcStatParsesAHostileCommandName was red, because the sleeping task counted as absent (`observe_linux_test.go:130`).
6. SP16: with a task read failure that counted as a zombie, TestProcUnreadableStatKeepsPresent was red in its unreadable and malformed task record cases (`observe_linux_test.go:150`).
7. SP16, the found branch: with a process rule that counted a process with no task as absent, TestProcUnreadableStatKeepsPresent was red in its no-task case (`observe_linux_test.go:150`).
8. SP17: with a gone member record that kept the group present, TestProcVanishedMemberIsSkipped was red, because the group counted as present (`observe_linux_test.go:168`).
9. SP33: before the registry held the FT482 Forbid entry, TestTicketSlicingAnchors was red, because the tree with the FT482 row raised no diagnostic (`registry_ticket_passes_test.go:37`). The row deletion came after this red.

SP18, SP19, and SP20 stayed green on their existing tests.

## SP-C2 review

### Standards

Count: 0 accepted. Worst issue: none. The axis raised six judgment calls at confidence 4 or lower. Five of them are no-op:

- The author wrote the witness reds because the orchestrator charge told it to.
- The coverage-cell citations obey the ticket-commit rule.
- The two platform probes repeat because of the file budget.
- The CHANGELOG sentence passed the prose lane.
- Spec-row tags in test comments are a repository convention.

Advice: the `Probe` doc comment at `internal/subprocess/protection/observe.go:29-31` restates what the code shows.

### Spec

Count: 0. Worst issue: none. Each decision of "The zombie rule" and "The retirement" holds, and SP35 holds. The axis ruled that an empty task or member set that stays present agrees with the spec's unreadable-keeps-present and ENOENT sentences.

### Coverage

Count: 3 accepted. Worst issue: SP-C2-C1.

- SP-C2-C1 (auto-fix, confidence 8): `internal/subprocess/protection/root_linux.go:83`. No test pins that a group with no member found stays present. A probe that answers absent there is silent. Add a group row whose root holds no member of the group.
- SP-C2-C2 (auto-fix, confidence 6): `root_linux.go:81`. No group row grades a member that has a stat record and no tasks. Add one group row with such a member.
- SP-C2-C3 (auto-fix, confidence 5): `root_linux.go:~89` and `~110`. A combined probe of the non-ENOENT task-directory error and the short-record guard is silent. Add a short record such as `7 (a) Z`, and add an unreadable task directory of a group member where the fixture can make one.

## SP-C2 repair r1 witness reds

The repair session ft482-repair-02-r1 added four rows to TestProcUnreadableStatKeepsPresent and observed these witness reds. Each red came through one `bench probe --mutants` call on `internal/subprocess/protection/root_linux.go`, and each probe restored its subject.

1. SP-C2-C1, the no member row: with `return found` changed to `return true`, the no member and member with no task cases were red (`observe_linux_test.go:155`).
2. SP-C2-C2, the member with no task row: with the group rule changed from `found = found || member` to `found = found || member || dead`, the member with no task case was red (`observe_linux_test.go:155`).
3. SP-C2-C3, the short task record row: with the stat guard changed from `end < 0 || len(fields) < 3` to `end < 0`, the short task record case was red (`observe_linux_test.go:155`).
4. SP-C2-C3, the unreadable member task directory row: with `return errors.Is(err, fs.ErrNotExist), false` changed to `return true, false`, that case was red (`observe_linux_test.go:155`).

### SP-C2 repair state

Repair round r1 (ft482-repair-02-r1, opus/high) closed SP-C2-C1, SP-C2-C2, and SP-C2-C3 at c2abc657d with test rows only. SP-C2 used 1 of its 2 repair cycles. The coordinator probe on the empty-group branch bit with 2 failed tests. The confirming round was necessary because the spec prose correction for the empty-set reading was in the chunk delta. Coverage closed all three findings, and Spec reaffirmed.

The Standards axis found that ticket 2 still states the old empty-set wording. The next plan commit corrects that ticket prose before the ticket 3 dispatch, as a prose-only spec correction.

## SP-C3 author witness reds

The ticket author ft482-author-03 observed these witness reds. The rows SP21, SP24, and SP26 were red on the base launch, before the production edit. Each other red came through one `bench probe` call on `internal/gate/background_launch.go`, and each probe restored its subject.

1. SP21, `TestMarkdownOnlyCommitStartsNoGrade`: on the base launch, the Markdown commit and the next commit each started a lane pass (`background_test.go:187`). The duplicated-facts sweep then merged that helper into `requireOnlyHeadPass`. The `commit-markdown` plan probe shows this red again on the committed test.
2. SP22, `TestMixedCommitStartsAGrade`: with the Markdown test changed to `!prose.IsSubjectName(change.Path) && false`, the mixed commit started no lane pass (`background_test.go:218`).
3. SP23, `TestMarkdownOnlyMergeStartsNoGrade`: with ` || markdownOnly(root, tip)` omitted, the base rule, the Markdown merge and the next merge each started a lane pass (`markdown_merge_background_test.go:92`).
4. SP24, `TestMarkdownDeletionStartsNoGrade`: on the base launch, the deletion tip started a lane pass before the control launch (`background_launch_test.go:250`).
5. SP25, `TestRenameToMarkdownStartsAGrade`: with each deleted path skipped, so that only the added `notes.md` counts, the rename started no lane pass (`background_launch_test.go:262`).
6. SP26, `TestQuotedMarkdownPathStartsNoGrade`: on the base launch, the `a "b".md` tip started a lane pass before the control launch (`background_launch_test.go:271`).
7. SP27, `TestMarkdownBuildInputStartsAGrade`: with the build-input diff read from the tip to the tip, the build-input change started no lane pass (`background_launch_test.go:286`).
8. SP28, `TestRootCommitStartsAGrade`: with the guard changed to `err == nil && len(changes) == 0`, so that a failed first-parent diff reads as empty, the root commit started no lane pass (`background_launch_test.go:298`).
9. SP29, `TestUnreadableTipStartsAGrade`: with the same guard change, the tip that names no commit started no lane pass (`background_launch_test.go:306`).
