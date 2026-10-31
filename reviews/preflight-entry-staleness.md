# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-entry-staleness/spec.md",
  "plan_digest": "sha256:61921ea7fbf2edf0094bea10f19540d2e33106ef92d140f0c1b3de8d63ee5bf9",
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
  },
  "amendments": [
    {
      "from": "sha256:9cce72026531bfc0a7a156d706d2a333f67c68a342dfeb0c584aae3679f2ab3b",
      "to": "sha256:61921ea7fbf2edf0094bea10f19540d2e33106ef92d140f0c1b3de8d63ee5bf9",
      "chunk_ids": {
        "ES-C1": [
          "ES-C1"
        ]
      }
    }
  ]
}
```

## ES-C1 ticket 01 author witness reds

The ticket 01 author `ft375-author-01` (opus / medium) reported these witness reds at ticket commit `a7cfab215`.
The author applied each mutation by hand, saw the named test fail, and restored each subject byte for byte.
The two plan probes `writes-landed` and `precedence` answered `bit` in `bench probe`.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/preflight/decision.go` | drop the `writes-landed` row from `Decide` | `TestWritesLandedBareEntry` | ES24 |
| `internal/preflight/fence_writes.go` | also send a landed entry to the `writes-resolve` red | `TestWritesLandedBareEntry` | ES25 |
| `internal/preflight/writes_landed.go` | drop `--first-parent -m` from the deleting-commit query | `TestWritesLandedShapes` | ES26, ES27 |
| `internal/preflight/writes_landed.go` | replace `--no-renames` with `--find-renames --full-diff` | `TestWritesLandedShapes` | ES28 |
| `internal/preflight/writes_landed.go` | accept an empty deleting-commit answer | `TestWritesLandedShapes` | ES29 |
| `internal/preflight/fence_writes.go` | classify the whole red by the landed set | `TestWritesLandedPrecedence/ES30` | ES30 |
| `internal/preflight/writes_landed.go` | ignore the tip `ls-tree` listing | `TestWritesLandedShapes` | ES31 |
| `internal/preflight/fence_writes.go` | drop `BuildDeletedPaths` from step 1 | `TestWritesLandedPrecedence/ES32` | ES32 |
| `internal/preflight/fence_writes.go` | grade `LandedDeletions` before `DeletedPaths` | `TestWritesLandedPrecedence/ES33` | ES33 |
| `internal/preflight/writes_landed.go` | replace the default-branch lookup with `HEAD` | `TestWritesLandedPrecedence/ES34` | ES34 |
| `internal/preflight/writes_landed.go` | use the raw path in place of the `:(literal)` pathspec | `TestWritesLandedShapes` | ES35 |
| `internal/preflight/fence_writes.go` | change the kept `writes-resolve` detail text | `TestTicketGrammarRowsCarryOwnDetails` | ES36 |
| `internal/preflight/fence_writes.go` | drop the commit from the `writes-landed` detail | `TestWritesLandedBareEntry` | ES37 |
| `internal/preflight/decision.go` | add the row in every mode | `TestWritesLandedReviewMode` | ES38 |
| `internal/preflight/decision.go` | add the row outside `ticketRow` | `TestTicketGatedRowsNotApplicableWithoutTickets` | ES39 |
| `internal/preflight/writes_landed.go` | remove the build-mode guard of the gatherer | `TestWritesLandedReviewMode` | ES40 |
| `internal/preflight/decision.go` | omit the row from the audit form | `TestWritesLandedKeepsAuditManual` | ES45 |

An `exempt` flag on a red row does not change `requiredGreen`, so the ES45 witness omits the row from the audit form. The build amendment of the ES45 and ES46 cells records that divergence.
