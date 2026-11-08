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
