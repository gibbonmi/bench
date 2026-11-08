# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/seal-policy/spec.md",
  "plan_digest": "sha256:60a5b3e8633ad83402c59bbeb24ea24d1256f62e61f95abbe3b299c14fa4d8fd",
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
