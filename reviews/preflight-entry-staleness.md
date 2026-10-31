# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-entry-staleness/spec.md",
  "plan_digest": "sha256:de3f820102ea29b8180f39876d5892962f0165e94df8b6612c2e12b7bcefcb4b",
  "implementation_session": "",
  "chunks": [
    {
      "id": "ES-C1",
      "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
      "tip": "07922d5437750b982cb74a31792286b933c6e463",
      "plan_digest": "sha256:de3f820102ea29b8180f39876d5892962f0165e94df8b6612c2e12b7bcefcb4b",
      "source_digest": "6ed05a0723fbfe51f5f39953e48115a85ad3cf8a",
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
        },
        {
          "id": "es-c1-t02-closure-prediction",
          "performer": "ft375-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a45b01c8b6c245008776cbb22fb6a64dcfce5af7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft375-author-02",
            "digest": "sha256:a0b2674ada3bf3e8b7369b7382f562f9e50cb19fda32279472a9c380c6c3c185",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/commitment/repository/closure.go,swap,failed,1,TestClosedRowPathsMatchClosure/refused-source,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,^TestClosedRowPathsMatchClosure$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,400\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,49,6\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestClosedRowPathsMatchClosure/refused-source,\"closure_test.go:162: ClosedRowPaths = [roadmap/FT1.md decisions/DEC1.md], \u003cnil>; want [roadmap/FT1.md]\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d9d130308c0b42a4118cbf4b63ad43534c395d61",
          "requirement": "closure-prediction",
          "command": "bench test --package ./internal/commitment --package ./internal/commitment/repository --run 'CommitmentCloses|ClosedRowPathsMatchClosure|CommitmentDeliver|CommitmentLightClosure|CommitmentClosureRefusesObligationFreeBinding'",
          "exit_code": 0,
          "probe": {
            "mutation": "Skip the RowOwner step in the store method, and record that TestClosedRowPathsMatchClosure turns red on the refused obligation source of ES14.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft375-author-02",
              "digest": "sha256:a0b2674ada3bf3e8b7369b7382f562f9e50cb19fda32279472a9c380c6c3c185",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/commitment/repository/closure.go,swap,failed,1,TestClosedRowPathsMatchClosure/refused-source,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,^TestClosedRowPathsMatchClosure$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,400\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,fail,49,6\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/commitment/repository,TestClosedRowPathsMatchClosure/refused-source,\"closure_test.go:162: ClosedRowPaths = [roadmap/FT1.md decisions/DEC1.md], \u003cnil>; want [roadmap/FT1.md]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "es-c1-r1-writes-landed",
          "performer": "ft375-repair-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6ed05a0723fbfe51f5f39953e48115a85ad3cf8a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft375-repair-01",
            "digest": "sha256:6c556a75f36f7922475e24ca7d100ef9a07e516ddb77cafe024402ab283eb1cf",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/writes_landed.go,omit,failed,1,TestWritesLandedShapes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedShapes$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,130,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedShapes,\"writes_landed_test.go:114: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: shapes/side.md (731414e4d2647b1db18362fd43ebdc59c048140c), one.md: shapes/old.md (fef7d76255b05c60a0cc972… (561 bytes)\",14\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "07922d5437750b982cb74a31792286b933c6e463",
          "requirement": "writes-landed",
          "command": "bench test --package ./internal/preflight --run 'WritesLanded|TicketGrammarRowsCarryOwnDetails|TicketGatedRowsNotApplicableWithoutTickets|DecideAllGreen|DecideChargeRows|SourceTip|WritesResolveEarlierChunkDeletion|CommittedDeletion|CommandBuildRendersFourGrammarRows'",
          "exit_code": 0,
          "probe": {
            "mutation": "Drop --first-parent -m from the deleting-commit query, and record that TestWritesLandedShapes turns red on the landing-merge case.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft375-repair-01",
              "digest": "sha256:6c556a75f36f7922475e24ca7d100ef9a07e516ddb77cafe024402ab283eb1cf",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/writes_landed.go,omit,failed,1,TestWritesLandedShapes,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedShapes$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,513\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,130,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedShapes,\"writes_landed_test.go:114: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: shapes/side.md (731414e4d2647b1db18362fd43ebdc59c048140c), one.md: shapes/old.md (fef7d76255b05c60a0cc972… (561 bytes)\",14\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "es-c1-r1-precedence",
          "performer": "ft375-repair-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "6ed05a0723fbfe51f5f39953e48115a85ad3cf8a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft375-repair-01",
            "digest": "sha256:d873e5ef436338d327bfb6f4cc09279787d5c3f9c269987a335761fc3d3cc717",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,1,TestWritesLandedPrecedence/deletion_inside_the_explicit_range,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedPrecedence$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,519\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,362,6\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedPrecedence/deletion_inside_the_explicit_range,\"writes_landed_test.go:167: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: docs/gone.md (20b60a23ec5ed532d25d8442b8a56ed72b865d8d)\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\"\\\", want \\\"\\\":\",8\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "07922d5437750b982cb74a31792286b933c6e463",
          "requirement": "precedence",
          "command": "bench test --package ./internal/preflight --run 'WritesLandedPrecedence'",
          "exit_code": 0,
          "probe": {
            "mutation": "Grade LandedDeletions before DeletedPaths in the check, and record that TestWritesLandedPrecedence turns red on the folded explicit-range case of ES33.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft375-repair-01",
              "digest": "sha256:d873e5ef436338d327bfb6f4cc09279787d5c3f9c269987a335761fc3d3cc717",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/fence_writes.go,swap,failed,1,TestWritesLandedPrecedence/deletion_inside_the_explicit_range,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,^TestWritesLandedPrecedence$,passed,6\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,519\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,362,6\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestWritesLandedPrecedence/deletion_inside_the_explicit_range,\"writes_landed_test.go:167: writes-landed row = \\\"  writes-landed,red,\\\\\\\\\\\"Writes: entry names a path that a landed commit deleted: one.md: docs/gone.md (20b60a23ec5ed532d25d8442b8a56ed72b865d8d)\\\\\\\\\\\",\\\\\\\\\\\"\\\\\\\\\\\"\\\", want \\\"\\\":\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "es-c1-r1-standards",
          "performer": "ft375-review-standards-1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a45b01c8b6c245008776cbb22fb6a64dcfce5af7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft375-review-standards-1",
            "digest": "sha256:407754341f678a2373d05d3fd0afc5b5d7bb82a1f6d80be7fa36b8b1795bd19f",
            "excerpt": "Standards axis, ES-C1, frozen tip d9d130308, evidence 602131f7 current.\nS1 internal/preflight/writes_landed.go:42 defaultBranchTip repeats the ResolvedDefault+ResolveCommit pair of internal/git/default_branch.go:39-47, internal/preflight/default_branch.go:95-104, internal/spec/widen.go:213-220. Rule: one source per fact. Confidence 7. ask-user.\nS2 internal/preflight/writes_landed_test.go:30 commitEverything equals commitAll at internal/preflight/staleness/drift_test.go:44; writes_landed_test.go:226 manualRequired re-reads the staleness row that table()/decision() read. Rule: one source per fact, fixture harnesses single-sourced. Confidence 7. ask-user.\nS3 internal/preflight/writes_landed.go:24 candidate filter repeats walkWrites cases; unkillable; spec mandates it. Confidence 5. no-op.\nS4 internal/preflight/writes_landed_test.go:7-8,70,88,133,183,206,138-156 spec row IDs in comments and test names. Rule: craft-comments provenance. Confidence 6. auto-fix.\nS5 internal/commitment/delivery.go:46 closing returns 5 values. Judgment call. Confidence 3. no-op.\nCount 5. Worst S1.\n"
          },
          "axis": "Standards",
          "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
          "tip": "d9d130308c0b42a4118cbf4b63ad43534c395d61",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5"
          ],
          "supersedes": []
        },
        {
          "id": "es-c1-r1-spec",
          "performer": "ft375-review-spec-1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a45b01c8b6c245008776cbb22fb6a64dcfce5af7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft375-review-spec-1",
            "digest": "sha256:dac654fc3f2ecffe9a4402a9657b11d6b53f32cd09d3386b25604f1594eed4a9",
            "excerpt": "Spec axis, ES-C1, frozen tip d9d130308, evidence 602131f7 current.\nP1 internal/preflight/writes_landed_test.go:155-170 ES32 fixture never puts the entry in LandedDeletions, so it does not show step 1 wins over step 4 when both hold; code correct at fence_writes.go:70-73; ES32 met as written. Confidence 3. no-op.\nCount 1. Worst P1. Audited ES13, ES14, ES24-ES40, ES45 against the frozen diff; no other miss; no scope creep.\n"
          },
          "axis": "Spec",
          "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
          "tip": "d9d130308c0b42a4118cbf4b63ad43534c395d61",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "es-c1-r1-coverage",
          "performer": "ft375-review-coverage-1",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a45b01c8b6c245008776cbb22fb6a64dcfce5af7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft375-review-coverage-1",
            "digest": "sha256:e9200500f274d859c0ba4b2cff433348268d420382f6a110d650ea55ead83c1b",
            "excerpt": "Coverage axis, ES-C1, frozen tip d9d130308, evidence 602131f7 current.\nC1 internal/preflight/fence_writes.go:70 seen[named] repeated-entry guard untested; probe omit '|| seen[named]' over ./internal/preflight: silent. Confidence 7. auto-fix.\nC2 internal/preflight/writes_landed.go:509 exists[entry] filter equivalent at row level; probe silent. Confidence 4. no-op.\nC3 internal/preflight/writes_landed.go:549 failed-query guard; probe drop 'err == nil &&': silent; near-equivalent. Confidence 4. no-op.\nCount 3. Worst C1. Probes 3: silent, silent, silent. Tree clean after probes.\n"
          },
          "axis": "Coverage",
          "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
          "tip": "d9d130308c0b42a4118cbf4b63ad43534c395d61",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "es-c1-r2-standards",
          "performer": "ft375-review-standards-2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6ed05a0723fbfe51f5f39953e48115a85ad3cf8a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft375-review-standards-2",
            "digest": "sha256:d1eda7980628adfc3c63d330e7d617296362bc4da82adb3fa1e47b7936f6a389",
            "excerpt": "Standards axis, ES-C1 confirming round after repair 1, evidence b26574f8 current, tip a5ab7ce4b.\nS1 closed: git.DefaultTip internal/git/default_branch.go:44 serves DefaultContains, landedDeletions, defaultTipEqualPaths, signedOffScope; each keeps its handling.\nS2 closed: preflighttest/history.go:92,101,118; drift_test.go:16-21 alias is single-source.\nS4 closed: no spec row IDs remain.\nNew findings: none. Count 0.\n"
          },
          "axis": "Standards",
          "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
          "tip": "07922d5437750b982cb74a31792286b933c6e463",
          "finding_ids": [],
          "supersedes": [
            "es-c1-r1-standards"
          ]
        },
        {
          "id": "es-c1-r2-spec",
          "performer": "ft375-review-spec-2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6ed05a0723fbfe51f5f39953e48115a85ad3cf8a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft375-review-spec-2",
            "digest": "sha256:f9de822f6cd079c236815b27d60a1fb35af8309d95029edcb1ecab31080afd0c",
            "excerpt": "Spec axis, ES-C1 reaffirmation over repair delta d9d130308..07922d543, evidence b26574f8 current.\nNo spec predicate changed answer: ES34 fail-closed, DefaultContains messages, defaultTipEqualPaths and signedOffScope answers, detail texts, precedence, query flags, coverage citations all hold.\nP2 (process, confidence 4): widen outside signed-off directories needs a decision. Orchestrator disposition no-op: bench spec widen --decision recorded both paths as decided (commit f90e68848).\nReaffirmed. Count 1 (no-op).\n"
          },
          "axis": "Spec",
          "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
          "tip": "07922d5437750b982cb74a31792286b933c6e463",
          "finding_ids": [],
          "supersedes": [
            "es-c1-r1-spec"
          ]
        },
        {
          "id": "es-c1-r2-coverage",
          "performer": "ft375-review-coverage-2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "6ed05a0723fbfe51f5f39953e48115a85ad3cf8a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:ft375-review-coverage-2",
            "digest": "sha256:32b05fbad9db0d188b6071ecf95e5df53e78e8716330f32ee6a168c11a14e111",
            "excerpt": "Coverage axis, ES-C1 confirming round after repair 1, evidence b26574f8 current, tree clean after probes.\nC1 closed: TestWritesLandedPrecedence/repeated_entry at internal/preflight/writes_landed_test.go:137.\nC4 internal/preflight/default_branch.go:41 unresolved arm of defaultTipEqualPaths untested; probe swap errors.Is guard: silent over ./internal/preflight/... Confidence 7. auto-fix.\nC5 internal/spec/widen.go:214 unresolved refusal of signedOffScope unpinned; probe same swap: silent over ./internal/spec/... ./cmd/... Confidence 6. auto-fix.\nProbe 3 internal/git/default_branch.go DefaultContains same swap: bit (TestMergeRefusesAnUnresolvedDefaultBranch).\nCount 2. Worst C4.\n"
          },
          "axis": "Coverage",
          "base": "fbcd16dfa1cff57bfd5590fca61aa9df5ac385f4",
          "tip": "07922d5437750b982cb74a31792286b933c6e463",
          "finding_ids": [
            "C4",
            "C5"
          ],
          "supersedes": [
            "es-c1-r1-coverage"
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
      "from": "sha256:9cce72026531bfc0a7a156d706d2a333f67c68a342dfeb0c584aae3679f2ab3b",
      "to": "sha256:61921ea7fbf2edf0094bea10f19540d2e33106ef92d140f0c1b3de8d63ee5bf9",
      "chunk_ids": {
        "ES-C1": [
          "ES-C1"
        ]
      }
    },
    {
      "from": "sha256:61921ea7fbf2edf0094bea10f19540d2e33106ef92d140f0c1b3de8d63ee5bf9",
      "to": "sha256:00eccbd9db08f8b3abfebba2213f8183f1f1fd339b94f9a790bf6c2bee72f41c",
      "chunk_ids": {
        "ES-C1": [
          "ES-C1"
        ]
      }
    },
    {
      "from": "sha256:00eccbd9db08f8b3abfebba2213f8183f1f1fd339b94f9a790bf6c2bee72f41c",
      "to": "sha256:de3f820102ea29b8180f39876d5892962f0165e94df8b6612c2e12b7bcefcb4b",
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

## ES-C1 ticket 02 author witness reds

The ticket 02 author `ft375-author-02` (opus / medium) reported these witness reds at ticket commit `d9d130308`.
The author applied each mutation by hand, saw the named test fail, and restored each subject.
The plan probe `closure-prediction` answered `bit` in `bench probe`.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/commitment/delivery.go` | force the obligation filter of `closing` to true | `TestCommitmentCloses/partial` | ES13 |
| `internal/commitment/delivery.go` | return an error for a deliverable that the milestone does not approve | `TestClosedRowPathsMatchClosure/unapproved` | ES14 branch |
| `internal/commitment/repository/closure.go` | return an error when the tree holds no policy | `TestClosedRowPathsMatchClosure/no-policy` | ES14 branch |
| `internal/commitment/repository/closure.go` | drop the error of `Closes` | `TestClosedRowPathsMatchClosure/obligation-free` | ES14 branch |
| `internal/commitment/repository/closure.go` | skip the row-owner step in `ClosedRowPaths` | `TestClosedRowPathsMatchClosure/refused-source` | ES14 |

## ES-C1 round 1 pickup

The three axes reviewed the frozen pair `fbcd16dfa..d9d130308` on sonnet / high. The round found 9 raw findings and 4 repair targets, and every target falls on ticket 01.

## Standards

Count: 5. Worst issue: S1.

- S1 (auto-fix, confidence 7): `internal/preflight/writes_landed.go:42` repeats the default-branch tip lookup of `internal/git/default_branch.go:39-47`, `internal/preflight/default_branch.go:95-104`, and `internal/spec/widen.go:213-220`. One shared lookup in `internal/git` replaces each copy, and each caller keeps its own failure handling.
- S2 (auto-fix, confidence 7): `internal/preflight/writes_landed_test.go:30` and `:226` repeat the commit helper and the `manual_required` read of `internal/preflight/staleness/drift_test.go:30-44`. One shared helper in `internal/preflight/preflighttest` replaces both copies.
- S4 (auto-fix, confidence 6): `internal/preflight/writes_landed_test.go` names spec row IDs in comments and test names. Name the behavior instead.
- S3 (no-op): the spec requires the gatherer filter at `writes_landed.go:24`.
- S5 (no-op): the five return values of `closing` are a judgment call.

## Spec

Count: 1. Worst issue: P1.

- P1 (no-op, confidence 3): the ES32 fixture does not hold both step 1 and step 4 for one entry. The code grades step 1 first, and ES32 is met as written.

## Coverage

Count: 3. Worst issue: C1.

- C1 (auto-fix, confidence 7): no test covers a repeated `Writes:` entry. A probe that omits the `seen[named]` guard at `internal/preflight/fence_writes.go:70` answered `silent`.
- C2 (no-op): the `exists[entry]` filter is an equivalent mutant at row level.
- C3 (no-op): the failed-query guard is a near-equivalent mutant.

## ES-C1 repair cycle 1 witness reds

The repair session `ft375-repair-01` (opus / medium) repaired S1, S2, S4, and C1 at repair commit `07922d543`. ES-C1 has used 1 of its 2 repair cycles.
The plan probes `writes-landed` and `precedence` answered `bit` again at that commit.

| subject | mutation | failing witness | finding |
| --- | --- | --- | --- |
| `internal/preflight/fence_writes.go` | omit the `seen[named]` repeated-entry guard | `TestWritesLandedPrecedence/repeated_entry` | C1 |

## ES-C1 confirming round 1 pickup

Standards closed S1, S2, and S4 and found no new finding. Spec reaffirmed the repair delta; its process finding P2 is a no-op, because `bench spec widen --decision` recorded the two outside paths as decided.
Coverage closed C1 and found two new gaps that the repair delta reaches. ES-C1 takes repair cycle 2 of 2 as a test-only repair.

- C4 (auto-fix, confidence 7): no test pins the unresolved arm of `defaultTipEqualPaths` at `internal/preflight/default_branch.go:41`. A probe that turns an unresolved default branch into a bootstrap failure answered `silent`.
- C5 (auto-fix, confidence 6): no test pins the unresolved refusal of `signedOffScope` at `internal/spec/widen.go:214`. The same probe answered `silent`.
