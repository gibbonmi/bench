# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-entry-staleness/spec.md",
  "plan_digest": "sha256:9cce72026531bfc0a7a156d706d2a333f67c68a342dfeb0c584aae3679f2ab3b",
  "implementation_session": "",
  "chunks": [
    {
      "id": "ES-C1",
      "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
      "tip": "a7cfab215f73c7c62d2db3249470ee666ca97eff",
      "plan_digest": "sha256:9cce72026531bfc0a7a156d706d2a333f67c68a342dfeb0c584aae3679f2ab3b",
      "source_digest": "7552417da12a6904c2e4fdc2e580c4df9d3c5ec5",
      "acceptance_rows": [
        "ES24",
        "ES25",
        "ES26",
        "ES27",
        "ES28",
        "ES29",
        "ES30",
        "ES31",
        "ES32",
        "ES33",
        "ES34",
        "ES35",
        "ES36",
        "ES37",
        "ES38",
        "ES39",
        "ES40",
        "ES45",
        "ES13",
        "ES14"
      ],
      "verification": [
        {
          "id": "es-c1-t01-writes-landed",
          "performer": "ft375-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7552417da12a6904c2e4fdc2e580c4df9d3c5ec5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft375-author-01",
            "digest": "sha256:c65978ee475cd9dceb1d1c435297d992219849ea9657cb8ec77756480cc21926",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/writes_landed.go,omit,failed,1,TestWritesLandedShapes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedShapes$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,130,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedShapes,\"writes_landed_test.go:125: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: shapes/side.md (499ab32accd90e116a573120ae1301c5633b60a0), one.md: shapes/old.md (5aeddc6c64ebea5ede791c5… (561 bytes)\",14\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a7cfab215f73c7c62d2db3249470ee666ca97eff",
          "requirement": "writes-landed",
          "command": "bench test --package ./internal/preflight --run 'WritesLanded|TicketGrammarRowsCarryOwnDetails|TicketGatedRowsNotApplicableWithoutTickets|DecideAllGreen|DecideChargeRows|SourceTip|WritesResolveEarlierChunkDeletion|CommittedDeletion|CommandBuildRendersFourGrammarRows'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop --first-parent -m from the deleting-commit query, and record that TestWritesLandedShapes turns red on the landing-merge case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft375-author-01",
              "digest": "sha256:c65978ee475cd9dceb1d1c435297d992219849ea9657cb8ec77756480cc21926",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/writes_landed.go,omit,failed,1,TestWritesLandedShapes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedShapes$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,130,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedShapes,\"writes_landed_test.go:125: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: shapes/side.md (499ab32accd90e116a573120ae1301c5633b60a0), one.md: shapes/old.md (5aeddc6c64ebea5ede791c5… (561 bytes)\",14\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "es-c1-t01-precedence",
          "performer": "ft375-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7552417da12a6904c2e4fdc2e580c4df9d3c5ec5",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft375-author-01",
            "digest": "sha256:40082b7197f8c77c0edc3a0959da951d36e1cac68a28e1cd63ae003dd1175d23",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,1,TestWritesLandedPrecedence/ES33_deletion_inside_the_explicit_range,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedPrecedence$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,524\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,300,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedPrecedence/ES33_deletion_inside_the_explicit_range,\"writes_landed_test.go:168: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: docs/gone.md (4441837e613cd6a6bfa43d6bb2b718682d1b4f37)\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\"\\\", want \\\"\\\":\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "a7cfab215f73c7c62d2db3249470ee666ca97eff",
          "requirement": "precedence",
          "command": "bench test --package ./internal/preflight --run 'WritesLandedPrecedence'",
          "exit_code": 0,
          "probe": {
            "mutation": "Grade LandedDeletions before DeletedPaths in the check, and record that TestWritesLandedPrecedence turns red on the folded explicit-range case of ES33.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft375-author-01",
              "digest": "sha256:40082b7197f8c77c0edc3a0959da951d36e1cac68a28e1cd63ae003dd1175d23",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,1,TestWritesLandedPrecedence/ES33_deletion_inside_the_explicit_range,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedPrecedence$,passed,5\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,524\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,300,5\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedPrecedence/ES33_deletion_inside_the_explicit_range,\"writes_landed_test.go:168: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: docs/gone.md (4441837e613cd6a6bfa43d6bb2b718682d1b4f37)\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\"\\\", want \\\"\\\":\",8\nskips[0]{package,test,reason}:\n"
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
