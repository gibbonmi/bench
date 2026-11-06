# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/production-test-seam-policy/spec.md",
  "plan_digest": "sha256:7e678524a4f7c5a93e7d074c672fefdfa367d97bc7f7b34529fca1bfb6891c96",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TS1",
      "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
      "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
      "plan_digest": "sha256:5c08a62494fc1fdd742618a16f4d9ccab37c50326fb6f594f7fef3f23db32a4b",
      "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
      "acceptance_rows": [
        "PS01",
        "PS02",
        "PS03",
        "PS12",
        "PS13",
        "PS14",
        "PS41",
        "PS42",
        "PS43",
        "PS50",
        "PS51",
        "PS52",
        "PS53",
        "PS54",
        "PS55",
        "PS61"
      ],
      "verification": [
        {
          "id": "ts1-01-wrapper-bypass",
          "performer": "ft343-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-author-01",
            "digest": "sha256:5ff4c883590de2b57efe65fadff342ba3c996a22b74d4c1e1136f5892cbb5329",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_test.go,swap,failed,2,TestEveryRetainedFixtureBitesThroughRegisteredOwner,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestEveryRetainedFixtureBitesThroughRegisteredOwner$,passed,289\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4108,289\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner,\"fixture_bite_test.go:33: completed fixture proofs do not equal the current producer: got 287, want 288\",1\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner/unregistered-port,\"fixture_bite_test.go:27: unregistered-port did not bite through owner injected-port-registry; want \\\"injected port has no registry row: internal/gitguard.Auditor is injected at a parameter of auditVerdict in internal/gitguard/gitguard.go\\\" in… (253 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "requirement": "wrapper-bypass",
          "command": "bench test --package ./internal/conformance --run '^TestEveryRetainedFixtureBitesThroughRegisteredOwner$'",
          "exit_code": 0,
          "probe": {
            "mutation": "make checkInjectedPortRegistry return without calling the relocated owner; require the injected-ports canary bite to fail with its exact expected diagnostic",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-author-01",
              "digest": "sha256:5ff4c883590de2b57efe65fadff342ba3c996a22b74d4c1e1136f5892cbb5329",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_test.go,swap,failed,2,TestEveryRetainedFixtureBitesThroughRegisteredOwner,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestEveryRetainedFixtureBitesThroughRegisteredOwner$,passed,289\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4108,289\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner,\"fixture_bite_test.go:33: completed fixture proofs do not equal the current producer: got 287, want 288\",1\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner/unregistered-port,\"fixture_bite_test.go:27: unregistered-port did not bite through owner injected-port-registry; want \\\"injected port has no registry row: internal/gitguard.Auditor is injected at a parameter of auditVerdict in internal/gitguard/gitguard.go\\\" in… (253 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ts1-review-standards",
          "performer": "ft343-review-ts1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts1",
            "digest": "sha256:8e17e424b12352d673152240db0d7ebd8855fe12de99136be6a455dde63a24b7",
            "excerpt": "Standards: no findings. The relocation leaves one AST scanner in internal/conformance/injectedports, and the old file keeps only the wrapper and the legacy bite test. The tests use the shared writeFixtureFile in conformance and one local writeRoot in the owner, as the spec allows. slices.Sort/Compact replaces an unexported test helper the package cannot reach, so it is not duplicated knowledge. Advice: put the per-row mutant reds in the record so the independent-expectation exception holds.\n"
          },
          "axis": "Standards",
          "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
          "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ts1-review-spec",
          "performer": "ft343-review-ts1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts1",
            "digest": "sha256:1a43e6e5bfbfcf2bffe0784b5bfbbee9bd244df71035ea37485b20fbf71c96d7",
            "excerpt": "Spec: no findings. All 16 TS1 rows map to live seams, the legacy messages stay byte-exact as exported constants, and checkInjectedPortRegistry calls Check once. Deriving only over Required follows the TS1 chunk row and ticket 04:29, and it conflicts with ticket 01:12, which goes to reviewer veto. Silence on a row outside Required stays legacy behavior until the TS4 census.\n"
          },
          "axis": "Spec",
          "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
          "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ts1-review-coverage",
          "performer": "ft343-review-ts1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts1",
            "digest": "sha256:b27ac4b95b8318a293563a59ffbc51892a22c1887f14edd53c6b981190ae46a3",
            "excerpt": "Coverage: no findings. Every mapped row's why-clause mutation changes its fixture's diagnostic heads, and the wrapper-bypass probe bit with the canary's exact diagnostic. Advice for TS4: ENOTDIR as absence, and symlinked package directories outside the root.\n"
          },
          "axis": "Coverage",
          "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
          "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "TS2",
      "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
      "tip": "25f5571dd24d6232974f986c7c601153caf5349a",
      "plan_digest": "sha256:4517bbfeca3a80a55cb305321af5ce4e7f580dc43a8d3f33d1af10f266fc08c4",
      "source_digest": "5017bb8a031464c90a925d040c7c25c57b89b955",
      "acceptance_rows": [
        "PS04",
        "PS05",
        "PS06",
        "PS07",
        "PS08",
        "PS15",
        "PS16",
        "PS17",
        "PS18",
        "PS19",
        "PS20",
        "PS21",
        "PS44",
        "PS57"
      ],
      "verification": [],
      "reviews": [
        {
          "id": "ts2-review-standards",
          "performer": "ft343-review-ts2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8534ee5c398a9f6fb2e6632a94406be03591bf1c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts2",
            "digest": "sha256:28208e50dc5a564ff9d8401f0f57be37fe12354e71375281be21f4745ee41571",
            "excerpt": "Standards: 0 findings. Worst: none. The whole-tree loader in source.go repeats derive.go's per-package parse and read path and its diagnostic prose. Carry it to TS4 as a scanner-copy collapse. Writes fence: every changed path is inside ticket 02's Writes: line.\n"
          },
          "axis": "Standards",
          "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "tip": "cdeb6e383ca0d1e2095ea1bb1308dbe062e3a8d8",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ts2-review-spec",
          "performer": "ft343-review-ts2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8534ee5c398a9f6fb2e6632a94406be03591bf1c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts2",
            "digest": "sha256:9fdb799ce909cd05361eec863aacabe37e7c7cc915caa72b6670f67c9d26e00e",
            "excerpt": "Spec: 2 findings. Worst: P2. Gate laneRunBinary and testreport runBinaryFactory are admitted while their consumers replace Verify with unconditional success, and no hook-active refusal witness exists (prospective has one at prospective_owner_test.go:98). P1: a same-package field, index, or pointer replacement of a package collaborator is not derived (scope.go assignedRoot), against spec line 119.\n"
          },
          "axis": "Spec",
          "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "tip": "cdeb6e383ca0d1e2095ea1bb1308dbe062e3a8d8",
          "finding_ids": [
            "TS2-P1",
            "TS2-P2"
          ],
          "supersedes": []
        },
        {
          "id": "ts2-review-coverage",
          "performer": "ft343-review-ts2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8534ee5c398a9f6fb2e6632a94406be03591bf1c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts2",
            "digest": "sha256:28590579b04ae5ce0579acd6250f81ef5ebcceaeb5ebe7ce5cfd437c0d829fbc",
            "excerpt": "Coverage: 1 finding. Worst: C1. No fixture covers a same-package field or index replacement of a package collaborator, so the P1 omission stays green. Sampled 14 of 67 rows: consumer reach, purpose, and serial isolation hold. Git join repair: still red when the bound regresses, but only at the go test timeout (600 s stub).\n"
          },
          "axis": "Coverage",
          "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "tip": "cdeb6e383ca0d1e2095ea1bb1308dbe062e3a8d8",
          "finding_ids": [
            "TS2-C1"
          ],
          "supersedes": []
        },
        {
          "id": "ts2-review-standards-r1",
          "performer": "ft343-review-ts2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5017bb8a031464c90a925d040c7c25c57b89b955",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts2",
            "digest": "sha256:6438b5ef3c3fba43b6b10a4a1408697cc761f2d179fde39c3de1f7e49bd84fd8",
            "excerpt": "Standards reaffirmed at 06ec9df: the repair delta adds no second derivation. The assignedRoot comment now states the field, index, and pointer behavior.\n"
          },
          "axis": "Standards",
          "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "tip": "25f5571dd24d6232974f986c7c601153caf5349a",
          "finding_ids": [],
          "supersedes": [
            "ts2-review-standards"
          ]
        },
        {
          "id": "ts2-review-spec-r1",
          "performer": "ft343-review-ts2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5017bb8a031464c90a925d040c7c25c57b89b955",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts2",
            "digest": "sha256:ff8e37dac8ad7f6045c960944a7662f810ca1df50aa3cee03503dcc4332eda48",
            "excerpt": "Spec: P1 closed (scope.go:103-112, plus three PS04 cases at derive_test.go:60-71). P2 closed (lane_run_test.go:123 and runbinary_test.go:157 keep the real verifier with the hook installed). The phases row is serial in substance, because the seam is read only on the test goroutine.\n"
          },
          "axis": "Spec",
          "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "tip": "25f5571dd24d6232974f986c7c601153caf5349a",
          "finding_ids": [],
          "supersedes": [
            "ts2-review-spec"
          ]
        },
        {
          "id": "ts2-review-coverage-r1",
          "performer": "ft343-review-ts2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "5017bb8a031464c90a925d040c7c25c57b89b955",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts2",
            "digest": "sha256:ca338f9995c0efd34929a1480f16ec65cc6b08a80a4feb519661ebe914a6f315",
            "excerpt": "Coverage: C1 closed with three red-witnessed PS04 forms. Open advice for TS4: an index or field store of a parameter in production code now makes a setter-derived seam without any test override.\n"
          },
          "axis": "Coverage",
          "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "tip": "25f5571dd24d6232974f986c7c601153caf5349a",
          "finding_ids": [],
          "supersedes": [
            "ts2-review-coverage"
          ]
        }
      ]
    },
    {
      "id": "TS3",
      "base": "25f5571dd24d6232974f986c7c601153caf5349a",
      "tip": "fa17d5a5e5b31cc8b2e32b4abf7ff2a05fb169c3",
      "plan_digest": "sha256:7e678524a4f7c5a93e7d074c672fefdfa367d97bc7f7b34529fca1bfb6891c96",
      "source_digest": "afdb682fa5ca025aa8ad69f46dd24e2d336782f3",
      "acceptance_rows": [
        "PS09",
        "PS10",
        "PS11",
        "PS22",
        "PS23",
        "PS24",
        "PS25",
        "PS26",
        "PS27",
        "PS28",
        "PS29",
        "PS30",
        "PS31",
        "PS32",
        "PS33",
        "PS34",
        "PS35",
        "PS36",
        "PS37",
        "PS38",
        "PS39",
        "PS40",
        "PS56",
        "PS58",
        "PS59",
        "PS60",
        "PS66"
      ],
      "verification": [
        {
          "id": "ts3-03-restore-omission",
          "performer": "ft343-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9c0ab83192d6112ce19d04aa55f3ac8ed4d58768",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-author-03",
            "digest": "sha256:7ef77f0cfa792c3db14abad976f402b26c87b03e4c1ac05145600aac4fe15e09",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/diff/command_test.go,omit,failed,1,TestSnapshotOverrideRestoresAfterFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/diff,all,passed,101\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,352\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,fail,6022,101\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/diff,TestSnapshotOverrideRestoresAfterFailure,\"command_test.go:149: after the fatal path the child observed, want prior=1 replacement=1:\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "3345204ff2b51418035702472dc1efa6fbc755df",
          "requirement": "restore-omission",
          "command": "bench test --package ./internal/diff",
          "exit_code": 0,
          "probe": {
            "mutation": "omit failure-safe restoration; require the distinguishable installed replacement witness to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-author-03",
              "digest": "sha256:7ef77f0cfa792c3db14abad976f402b26c87b03e4c1ac05145600aac4fe15e09",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/diff/command_test.go,omit,failed,1,TestSnapshotOverrideRestoresAfterFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/diff,all,passed,101\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,352\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,fail,6022,101\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/diff,TestSnapshotOverrideRestoresAfterFailure,\"command_test.go:149: after the fatal path the child observed, want prior=1 replacement=1:\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ts3-03-join-omission",
          "performer": "ft343-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9c0ab83192d6112ce19d04aa55f3ac8ed4d58768",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-author-03",
            "digest": "sha256:f4d4b060a70ae7d4b693eb66e36b43abf9394d57f576feed4808bccb8df992b7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/guards/guards_cleanup_test.go,omit,failed,1,TestScanOverrideJoinsBeforeRestore,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/guards,all,passed,57\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,374\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6343,57\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanOverrideJoinsBeforeRestore,\"guards_cleanup_test.go:107: the held worker observed \\\"prior\\\" after its release, want the replacement it ran with\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "3345204ff2b51418035702472dc1efa6fbc755df",
          "requirement": "join-omission",
          "command": "bench test --package ./internal/guards",
          "exit_code": 0,
          "probe": {
            "mutation": "remove join before restoration; require the held live consumer witness to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-author-03",
              "digest": "sha256:f4d4b060a70ae7d4b693eb66e36b43abf9394d57f576feed4808bccb8df992b7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/guards/guards_cleanup_test.go,omit,failed,1,TestScanOverrideJoinsBeforeRestore,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/guards,all,passed,57\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,374\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6343,57\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanOverrideJoinsBeforeRestore,\"guards_cleanup_test.go:107: the held worker observed \\\"prior\\\" after its release, want the replacement it ran with\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ts3-03-r1-restore-omission",
          "performer": "ft343-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b6cee2ad8cad7d40aedec71b95b9d3169ff942b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-repair-03-r1",
            "digest": "sha256:1136aa7548746c2ac498258ddeccd9cb7ce1794e923d74770aa465d4753ee6d7",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/diff/command_test.go,omit,failed,1,TestSnapshotOverrideRestoresAfterFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/diff,all,passed,101\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,352\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,fail,5988,101\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/diff,TestSnapshotOverrideRestoresAfterFailure,\"command_test.go:148: after the fatal path the child observed, want prior=1 replacement=1:\",3\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "3c4d8e1a724f262db9494d109957d78e4e944642",
          "requirement": "restore-omission",
          "command": "bench test --package ./internal/diff",
          "exit_code": 0,
          "probe": {
            "mutation": "omit failure-safe restoration; require the distinguishable installed replacement witness to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-repair-03-r1",
              "digest": "sha256:1136aa7548746c2ac498258ddeccd9cb7ce1794e923d74770aa465d4753ee6d7",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/diff/command_test.go,omit,failed,1,TestSnapshotOverrideRestoresAfterFailure,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/diff,all,passed,101\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,352\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,fail,5988,101\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/diff,TestSnapshotOverrideRestoresAfterFailure,\"command_test.go:148: after the fatal path the child observed, want prior=1 replacement=1:\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ts3-03-r1-join-omission",
          "performer": "ft343-repair-03-r1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7b6cee2ad8cad7d40aedec71b95b9d3169ff942b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-repair-03-r1",
            "digest": "sha256:302441baf4c856c0d0fb645f2d56f2196dd2dd3c938adeefcbb9991244a6b2a6",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/guards/guards_cleanup_test.go,omit,failed,1,TestScanOverrideJoinsBeforeRestore,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/guards,all,passed,57\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,374\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6343,57\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanOverrideJoinsBeforeRestore,\"guards_cleanup_test.go:106: the held worker observed \\\"prior\\\" after its release, want the replacement it ran with\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "3c4d8e1a724f262db9494d109957d78e4e944642",
          "requirement": "join-omission",
          "command": "bench test --package ./internal/guards",
          "exit_code": 0,
          "probe": {
            "mutation": "remove join before restoration; require the held live consumer witness to fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-repair-03-r1",
              "digest": "sha256:302441baf4c856c0d0fb645f2d56f2196dd2dd3c938adeefcbb9991244a6b2a6",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/guards/guards_cleanup_test.go,omit,failed,1,TestScanOverrideJoinsBeforeRestore,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/guards,all,passed,57\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,374\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/guards,fail,6343,57\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/guards,TestScanOverrideJoinsBeforeRestore,\"guards_cleanup_test.go:106: the held worker observed \\\"prior\\\" after its release, want the replacement it ran with\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ts3-review-standards",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9c0ab83192d6112ce19d04aa55f3ac8ed4d58768",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:71ca53e84edaa1e54bb88ebf756fd79498238cd1a952cffa26651fada45bad15",
            "excerpt": "Standards: 0 findings. gate.DeclaredEnvironment wraps loadManifest, keeps its validation, and stays the one manifest reader. The phase family comes only from readReleaseRegistries, and the owner holds no phase list. PS66 holders agree. Advice: align or state the empty-verify versus empty-manifest asymmetry in gradedRootFamilies.\n"
          },
          "axis": "Standards",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3345204ff2b51418035702472dc1efa6fbc755df",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ts3-review-spec",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9c0ab83192d6112ce19d04aa55f3ac8ed4d58768",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:732164a3bbf8f38b69636396124d20c43fb50969b6605f72b17edc0ec4b43b20",
            "excerpt": "Spec: 1 finding. TS3-P1: the BENCH_OFFLINE and PATH rows name Go owners, but their consumers exercise only bin/bench.sh. Review-owned PS27, PS39, PS60, PS66 hold; PS38 is vacuous. Both orchestrator decisions are sound. Advice: record a PS34 omission red at final verification.\n"
          },
          "axis": "Spec",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3345204ff2b51418035702472dc1efa6fbc755df",
          "finding_ids": [
            "TS3-P1"
          ],
          "supersedes": []
        },
        {
          "id": "ts3-review-coverage",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "9c0ab83192d6112ce19d04aa55f3ac8ed4d58768",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:a92ebefb50c7e5e154780875794ebcb29d96885ae86ef30ddac17b6ae9215d52",
            "excerpt": "Coverage: 1 finding. TS3-C1: an environment-row consumer passes when any outermost NAME= string appears in its body, including an assertion literal, and without reaching the owner, against spec line 181. Sampled 13 of 55 rows: 11 true, 2 fail as TS3-P1. PS22, PS23, PS24 witnesses hold.\n"
          },
          "axis": "Coverage",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3345204ff2b51418035702472dc1efa6fbc755df",
          "finding_ids": [
            "TS3-C1"
          ],
          "supersedes": []
        },
        {
          "id": "ts3-review-standards-r1",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b6cee2ad8cad7d40aedec71b95b9d3169ff942b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:292edd80af86116d367d206f4cda917fbbee06e273ddc08a859e35f09c90615d",
            "excerpt": "Standards reaffirmed at d8a77f1: envflow.go is the only flow derivation, and consumer reach and hook derivation share one successors walk.\n"
          },
          "axis": "Standards",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3c4d8e1a724f262db9494d109957d78e4e944642",
          "finding_ids": [],
          "supersedes": [
            "ts3-review-standards"
          ]
        },
        {
          "id": "ts3-review-spec-r1",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b6cee2ad8cad7d40aedec71b95b9d3169ff942b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:4eeb8d745615a26b0a66e0c7f453112c4e72fca7f41802189c04be1121c34870",
            "excerpt": "Spec: TS3-P1 closed; BENCH_OFFLINE and PATH name consumers that reach their Go owners. The reach-based hook decision is sound against spec lines 135 and 181. Sampled 7 child-process rows; 5 confirmed, 2 likely.\n"
          },
          "axis": "Spec",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3c4d8e1a724f262db9494d109957d78e4e944642",
          "finding_ids": [],
          "supersedes": [
            "ts3-review-spec"
          ]
        },
        {
          "id": "ts3-review-coverage-r1",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b6cee2ad8cad7d40aedec71b95b9d3169ff942b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:6463811d5f9546bbd66cbf220c41cc80d88f4e53e39182a5ed528e2befaa7adc",
            "excerpt": "Coverage: TS3-C1 closed through env-flow positions. TS3-C2 open: an in-process set reached only through a function value, or inherited by a child with no explicit entry, derives no hook and no diagnostic, so the census fails open.\n"
          },
          "axis": "Coverage",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3c4d8e1a724f262db9494d109957d78e4e944642",
          "finding_ids": [
            "TS3-C2"
          ],
          "supersedes": [
            "ts3-review-coverage"
          ]
        },
        {
          "id": "ts3-review-coverage-r2",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7b6cee2ad8cad7d40aedec71b95b9d3169ff942b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:3b7a3e8e64a8265fe4ed092780608b97db76e616b30061691fd43b3328069a7d",
            "excerpt": "Coverage reaffirmed at fa17d5a: TS3-C2 is closed as the reviewer-accepted Won't handle line at spec line 478, which gives those consumers to review. No code changed since the last round, and no new finding.\n"
          },
          "axis": "Coverage",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "3c4d8e1a724f262db9494d109957d78e4e944642",
          "finding_ids": [],
          "supersedes": [
            "ts3-review-coverage-r1"
          ]
        },
        {
          "id": "ts3-review-standards-r2",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afdb682fa5ca025aa8ad69f46dd24e2d336782f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:c98064070ab19e752a25e0f38fe60023511287746059daedb0fa9d61ca383a93",
            "excerpt": "Standards pass at 25f5571..fa17d5a: the delta changes no code, only one spec line and record prose, and adds no second source of any fact.\n"
          },
          "axis": "Standards",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "fa17d5a5e5b31cc8b2e32b4abf7ff2a05fb169c3",
          "finding_ids": [],
          "supersedes": [
            "ts3-review-standards-r1"
          ]
        },
        {
          "id": "ts3-review-spec-r2",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afdb682fa5ca025aa8ad69f46dd24e2d336782f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:0d8c6df05bf60bc9580bbf041f1f3b705b6fdfc964b239923b1848a54c500206",
            "excerpt": "Spec pass at 25f5571..fa17d5a: spec.md:478 adds the reviewer-decided Won't handle line, narrows no acceptance row, and TS3-P1 stays closed.\n"
          },
          "axis": "Spec",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "fa17d5a5e5b31cc8b2e32b4abf7ff2a05fb169c3",
          "finding_ids": [],
          "supersedes": [
            "ts3-review-spec-r1"
          ]
        },
        {
          "id": "ts3-review-coverage-r3",
          "performer": "ft343-review-ts3",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "afdb682fa5ca025aa8ad69f46dd24e2d336782f3",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts3",
            "digest": "sha256:28a5221484927f4d8e269d39cc037a481d3bf0f42b37a8a77e83deb62c48e0dd",
            "excerpt": "Coverage pass at 25f5571..fa17d5a: spec.md:478 marks both paths as Won't handle, so TS3-C2 is closed; TS3-C1 and its witnesses are unchanged.\n"
          },
          "axis": "Coverage",
          "base": "25f5571dd24d6232974f986c7c601153caf5349a",
          "tip": "fa17d5a5e5b31cc8b2e32b4abf7ff2a05fb169c3",
          "finding_ids": [],
          "supersedes": [
            "ts3-review-coverage-r2"
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
      "from": "sha256:5c08a62494fc1fdd742618a16f4d9ccab37c50326fb6f594f7fef3f23db32a4b",
      "to": "sha256:d60cbebbc23104cb18f056422e79d5ec33c22e8db5f7f151705f4e208f69a747",
      "chunk_ids": {
        "TS1": [
          "TS1"
        ]
      }
    },
    {
      "from": "sha256:d60cbebbc23104cb18f056422e79d5ec33c22e8db5f7f151705f4e208f69a747",
      "to": "sha256:4517bbfeca3a80a55cb305321af5ce4e7f580dc43a8d3f33d1af10f266fc08c4",
      "chunk_ids": {
        "TS1": [
          "TS1"
        ],
        "TS2": [
          "TS2"
        ]
      }
    },
    {
      "from": "sha256:4517bbfeca3a80a55cb305321af5ce4e7f580dc43a8d3f33d1af10f266fc08c4",
      "to": "sha256:f751c3cde724773c3d43fc845aaf16735229d50bb7885bfe5aecf4f35304510c",
      "chunk_ids": {
        "TS1": [
          "TS1"
        ],
        "TS2": [
          "TS2"
        ]
      }
    },
    {
      "from": "sha256:f751c3cde724773c3d43fc845aaf16735229d50bb7885bfe5aecf4f35304510c",
      "to": "sha256:cb88367034249883e1ef69c35305e63340bbd2ae08f03f5f2f942449807d3698",
      "chunk_ids": {
        "TS1": [
          "TS1"
        ],
        "TS2": [
          "TS2"
        ],
        "TS3": [
          "TS3"
        ]
      }
    },
    {
      "from": "sha256:cb88367034249883e1ef69c35305e63340bbd2ae08f03f5f2f942449807d3698",
      "to": "sha256:7e678524a4f7c5a93e7d074c672fefdfa367d97bc7f7b34529fca1bfb6891c96",
      "chunk_ids": {
        "TS1": [
          "TS1"
        ],
        "TS2": [
          "TS2"
        ],
        "TS3": [
          "TS3"
        ]
      }
    }
  ]
}
```

## TS1 author witness reds

The author `ft343-author-01` reported these witness reds at `d626a7a`. Each mutant reported `bit`, and each file restored byte-identically.

| row | file | mutation | failing test |
|---|---|---|---|
| PS01 | `derive.go` | the interface arm returns false | `TestSeamDerivationForms` |
| PS02 | `derive.go` | the function-type arm returns false | `TestSeamDerivationForms` |
| PS03 | `derive.go` | the struct arm returns false | `TestSeamDerivationForms` |
| PS50 | `derive.go` | the pointer arm returns an empty name | `TestSeamDerivationForms` |
| PS51 | `derive.go` | the type-assertion arm returns early | `TestSeamDerivationForms` |
| PS52 | `derive.go` | the case-clause arm returns early | `TestSeamDerivationForms` |
| PS53 | `derive.go` | the interface arm returns true | `TestSeamDerivationForms` |
| PS54 | `derive.go` | a data field is skipped | `TestSeamDerivationForms` |
| PS55 | `derive.go` | the injection test is always true | `TestSeamDerivationForms` |
| PS12 | `derive.go` | the parse diagnostic is omitted | `TestSeamAdmissionEvidence` |
| PS13 | `check.go` | the orphan guard is always true | `TestSeamAdmissionEvidence` |
| PS14 | `check.go` | the zero-inventory branch is disabled | `TestSeamAdmissionEvidence` |
| PS42 | `check.go` | the package path joins the live checkout | `TestSeamAdmissionEvidence` |
| PS43 | `check.go` | every stat error counts as absence | `TestSeamAdmissionEvidence` |
| PS61 | `check.go` | a blank exemption compares as valid | `TestSeamAdmissionEvidence` |

## TS2 author witness reds

The author `ft343-author-02` reported these witness reds before `cdeb6e3`. Each mutant reported `bit`, each file restored byte-identically, and each baseline was green.

| file | mutation | failing case |
|---|---|---|
| `variables.go` | the test-replacement census is disabled | PS04, PS05, PS06, PS57 |
| `variables.go` | cross-package forwarding is disabled | PS07 |
| `variables.go` | same-package forwarding is disabled | PS08 |
| `variables.go` | the reader filter is disabled | the test-support negative case |
| `variables.go` | a `ForTest` name rule is added | PS44 |
| `variables.go` | setter reach and helper reach are disabled | the admitted-consumer cases |
| `variables.go` | the `Test` prefix check is omitted | the PS16 helper case |
| `variables.go` | the setter-owner exclusion is omitted | the PS15 setter-as-owner case |
| `variables.go` | the initializer reader is dropped | the initializer-owner case |
| `scope.go` | the parameter test is bypassed | the production-data negative case |
| `scope.go` | scope resolution is disabled | the shadowing negative case |
| `scope.go` | only the first assignment target is read | PS57 |
| `scope.go` | import resolution is disabled | PS04, PS07 |
| `seams.go` | the owner, purpose, and isolation checks are removed one at a time | PS15, PS18, PS19 |
| `seams.go` | the consumer checks are inverted or emptied | PS16, PS17 |
| `seams.go` | the duplicate, alias, and orphan checks are removed one at a time | PS20, PS21, the stale-row case |
| `seams.go` | the unregistered-setter diagnostic is removed | PS07, PS08 |

The repair session `ft343-repair-02-r1` reported these witness reds at `25f5571`. Each probe row restored the file.

| file | mutation | failing test |
|---|---|---|
| `scope.go` | the pointer case matches a different node type | `TestSeamDerivationForms/PS04_pointer_replacement` |
| `scope.go` | the index case matches a different node type | `TestSeamDerivationForms/PS04_index_replacement` |
| `scope.go` | the old import-only selector check returns | `TestSeamDerivationForms/PS04_same-package_field_replacement` |
| `internal/runbinary/runbinary.go` | the verifier is forced to success | `TestLaneRunBinaryHookKeepsFreshnessRefusal` |
| `internal/runbinary/runbinary.go` | the verifier is forced to success | `TestSelectRunBinaryHookKeepsFreshnessRefusal` |

The reviewer decided on 2026-11-06 to retire `indexEncoder`, its setter, and its forwarder, because the seam had no consumer and no injectable failure.

## TS3 author witness reds

The author `ft343-author-03` reported these witness reds before `3345204`. Each mutant reported `bit`, and each file restored byte-identically.

| row | file | mutation | failing test |
|---|---|---|---|
| PS09 | `source.go` | the alias import binding is dropped | `TestSeamDerivationForms` |
| PS56 | `controls.go` | the forwarding predicate callee is removed | `TestSeamDerivationForms` |
| PS10 | `controls.go` | the case fold is disabled | `TestSeamDerivationForms` |
| PS11 | `controls.go` | every identity counts as resolved | `TestSeamDerivationForms` |
| PS31, PS32 | `controls.go` | the family lookup is disabled | `TestSeamDerivationForms` |
| PS31, PS32 | `controls.go` | the environment entries are dropped | `TestSeamDerivationForms` |
| PS11 | `controls.go` | the empty-family split is removed | `TestSeamDerivationForms` |
| PS65 | `variables.go` | the test-support classification is empty | `TestSeamDerivationForms` |
| PS26, PS63 | `controls.go` | the runtime production reach is disabled | `TestSeamAdmissionEvidence` |
| PS22 | `internal/diff/command_test.go` | the failure-safe cleanup is omitted | `TestSnapshotOverrideRestoresAfterFailure` |
| PS23 | `internal/guards/guards_cleanup_test.go` | the join is omitted | `TestScanOverrideJoinsBeforeRestore` |
| PS24 | `internal/assessment/store.go` | a shared package global replaces the instance collaborator | `TestInstanceCollaboratorsStayLocal` |
| PS25 | `internal/releasepreflight/vulnerability.go` | the scanner override returns success | `TestSeamAuthorityPreserved` |
| PS33 | `internal/releasepreflight/vulnerability.go` | the DATE read returns an empty value | `TestVulnerabilityEntryDateOverride` |
| PS58 | `internal/bounds/bounds.go` | the context drops its parent | `TestUnboundedParentCancellation` |

The repair session `ft343-repair-03-r1` reported these witness reds before `3c4d8e1`. Each mutant reported `bit`, and each file restored.

| file | mutation | failing test |
|---|---|---|
| `internal/gate/testlog_exec.go` | the root comes from the working directory, not the environment | `TestCreateTestLogKeysThePackageByTheConformanceRoot` |
| `controls.go` | the reader reach check is always true | `TestSeamDerivationForms` |
| `controls.go` | the test-method successors are omitted | `TestSeamAdmissionEvidence` |
| `controls.go` | the production method successors are dropped | `TestSeamAdmissionEvidence` |
| `controls.go` | the entry reach filter is disabled | `TestSeamAdmissionEvidence` |
| `envflow.go` | an unknown position counts as an environment | `TestSeamDerivationForms` |
| `envflow.go` | a parameter slot, method key, keyed field, or return flow is dropped | `TestSeamDerivationForms` |
| `seams.go` | the owner reach is ignored | `TestSeamAdmissionEvidence` |
| `seams.go` | the child acceptance is dropped | `TestSeamAdmissionEvidence` |

During the repair, the orchestrator made a decision under spec line 135. A hook derives only from a test set that reaches its reader, or from a child-process entry. That decision replaces the earlier literal rule, and it removed seven rows with no reaching test.

The author did not demonstrate a PS34 readiness omission, because `bench probe` has no route to the ship-tier root test.

The orchestrator decided two spec gaps on 2026-11-06, and the reviewer can veto both. The literal environment-hook rule of spec line 135 bounds the admission set. The readiness witness lives in the native release probe.

## Standards

TS1 findings: 0. Worst issue: none.

TS2 findings: 0. Worst issue: none.
Carry to TS4 (advice): `source.go` repeats the per-package parse and read path of `derive.go`; TS4 collapses the scanner copy.

## Spec

TS1 findings: 0. Worst issue: none.
The reviewer flags ticket 01 line 12 for reviewer veto: it claims the wide derivation that the chunk table gives to TS4.
Carry to the TS4 charge: a row outside the required packages stays silent until the TS4 census grades it.

TS2 findings: 2. Worst issue: TS2-P2.

- TS2-P1, severity 3, confidence 8, auto-fix: `assignedRoot` in `internal/conformance/injectedports/scope.go` drops a same-package field, index, or pointer replacement of a package variable, against spec line 119.
- TS2-P2, severity 4, confidence 5, auto-fix under `--auto-approve` and spec line 208. The consumers of the `laneRunBinary` and `runBinaryFactory` rows replace `Verify` with unconditional success. No test keeps the real validation active with the hook installed and still observes refusal.

The orchestrator also accepts three evidence corrections from the TS2 advice for the same repair. The repair records the publication `clock` exclusion that ticket 02 requires. It also corrects the `models` isolation text and the `sessioninspect.phases` purpose text.

## Coverage

TS1 findings: 0. Worst issue: none.
Advice for TS4: decide ENOTDIR as absence, and decide a symlinked package directory that points outside the root.

TS2 findings: 1. Worst issue: TS2-C1.

- TS2-C1, severity 3, confidence 8, auto-fix with TS2-P1. No derivation case covers a same-package field or index replacement of a package variable.

Advice for TS4: decide whether a production setter with no test override is a seam. Qualify the owner identity with its package, and decide a symlinked `.go` file.

TS3 Standards findings: 0. Worst issue: none.

TS3 Spec findings: 1. Worst issue: TS3-P1.

- TS3-P1, severity 3, confidence 8, auto-fix. The `BENCH_OFFLINE` and `PATH` rows name Go owners, but their consumers run only `bin/bench.sh`.

TS3 Coverage findings: 1. Worst issue: TS3-C1.

- TS3-C1, severity 3, confidence 7, auto-fix by orchestrator decision under spec line 181. An environment-row consumer passes on any outermost `NAME=` string, and it need not reach the owner.

The orchestrator decision for TS3-C1 has two parts. First, an environment set counts only where the value flows into an environment. Second, an in-process consumer must reach the owner through production calls. A child-process consumer stays accepted, and review owns its truth under spec lines 182 and 183.

Advice for the final verification: record a PS34 readiness omission red through the ship-tier route, if a probe route exists.

Repair cycles consumed for TS1: 0 of 2.
Repair cycles consumed for TS2: 1 of 2 (repair round r1).
Repair cycles consumed for TS3: 2 of 2 (repair rounds r1 and r2).
The reviewer closed TS3-C2 on 2026-11-06 as an accepted limitation, and the spec states it as a Won't handle line. The r2 diff is discarded, and its recovery ref is kept.
The TS3 confirming round closes TS3-P1 and TS3-C1, and it opens TS3-C2.

- TS3-C2, severity 3, confidence 6, auto-fix under `--auto-approve` in repair round r2. The reach-based census fails open in two cases. One is an in-process set that reaches its reader only through a function value. The other is a child that inherits the set with no explicit entry.

The TS2 confirming round closes TS2-P1, TS2-P2, and TS2-C1, and it adds no new finding.
Advice for TS4: a production store of a parameter into a package map or field now derives a setter seam without a test override. A negative control must decide it.
