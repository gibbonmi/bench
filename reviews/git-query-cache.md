# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/git-query-cache/spec.md",
  "plan_digest": "sha256:7d3a811db66c80945a137d70bfb8b0c7995202a3a5bf521a5e9bca38fb700356",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GQ-C1",
      "base": "f979bae6807f82e9bc4a259b78076c3170710597",
      "tip": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
      "plan_digest": "sha256:7d3a811db66c80945a137d70bfb8b0c7995202a3a5bf521a5e9bca38fb700356",
      "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
      "acceptance_rows": [
        "GQ1",
        "GQ2",
        "GQ3",
        "GQ4",
        "GQ5",
        "GQ6",
        "GQ7",
        "GQ8",
        "GQ9",
        "GQ10",
        "GQ13",
        "GQ14",
        "GQ15",
        "GQ16",
        "GQ17",
        "GQ18",
        "GQ19",
        "GQ20",
        "GQ21",
        "GQ22",
        "GQ23",
        "GQ24",
        "GQ25",
        "GQ26",
        "GQ43",
        "GQ44",
        "GQ45",
        "GQ47"
      ],
      "verification": [
        {
          "id": "t01-probe-01-hit",
          "performer": "claude:gq-t01-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
            "digest": "sha256:f8f661e0cfed631a1c568308ddfda0ddde5bc474b41dc8ecf4e4c59b31d25213",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,7,TestRunObjectReads/absent_path,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReads,passed,11\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1424\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,67,11\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/absent_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/blob,\"run_object_test.go:110: starts of \\\"cat-file blob ce013625030ba8dba906f756967f9e9ca394464a\\\" = 4, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/commit,\"run_object_test.go:110: starts of \\\"rev-parse --verify --quiet 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/non-ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 9324b9253cfba87b42cd827ede089247cbcfad3f\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_file,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "requirement": "probe-01-hit",
          "command": "bench probe internal/git/run.go --swap 'if raw, ok := r.objects[key]; ok {' --with 'if raw, ok := r.objects[key]; false && ok {' --package ./internal/git --run TestRunObjectReads",
          "exit_code": 0,
          "probe": {
            "mutation": "Never answer an object read from the store. Each second read then starts git, so the object-read test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
              "digest": "sha256:f8f661e0cfed631a1c568308ddfda0ddde5bc474b41dc8ecf4e4c59b31d25213",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,7,TestRunObjectReads/absent_path,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReads,passed,11\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1424\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,67,11\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/absent_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/blob,\"run_object_test.go:110: starts of \\\"cat-file blob ce013625030ba8dba906f756967f9e9ca394464a\\\" = 4, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/commit,\"run_object_test.go:110: starts of \\\"rev-parse --verify --quiet 8103df5ae7d83f834a8d236b36d62eaee5a4ead6\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/non-ancestor,\"run_object_test.go:110: starts of \\\"--is-ancestor 9324b9253cfba87b42cd827ede089247cbcfad3f\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_file,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_path,\"run_object_test.go:110: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01-probe-01-live",
          "performer": "claude:gq-t01-author-20261028",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
            "digest": "sha256:fed66f682741bdb6213eefbd1c8c4b1420aec6c0cb6dd49364a3ce1eca64313b",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunObjectReadsStayLive,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReadsStayLive,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,31,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReadsStayLive,\"run_object_test.go:166: reads after the commit and the update-ref = \\\"hello\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630f2e0adb550\\\\\\\\ta.txt|18c9d8cb4d427298e0689ab1cb3c2ca91d0d6ff8|true\\\", want \\\"next\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630… (433 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "requirement": "probe-01-live",
          "command": "bench probe internal/git/run.go --swap 'if !FullObjectID(object) {' --with 'if false {' --package ./internal/git --run TestRunObjectReadsStayLive",
          "exit_code": 0,
          "probe": {
            "mutation": "Cache a blob read whose operand is a ref. The second read of the moved ref then returns the old bytes, so the live-read test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gq-t01-author-20261028",
              "digest": "sha256:fed66f682741bdb6213eefbd1c8c4b1420aec6c0cb6dd49364a3ce1eca64313b",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunObjectReadsStayLive,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReadsStayLive,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,31,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReadsStayLive,\"run_object_test.go:166: reads after the commit and the update-ref = \\\"hello\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630f2e0adb550\\\\\\\\ta.txt|18c9d8cb4d427298e0689ab1cb3c2ca91d0d6ff8|true\\\", want \\\"next\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630… (433 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gq-c1-standards-1",
          "performer": "claude:gq-c1-standards-20261028",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c1-standards-20261028",
            "digest": "sha256:522a21ad204fb6e20273e50e9869c1397fd977a82e7f763b06d9616b382f3971",
            "excerpt": "axis: Standards\nverdict: findings\nfinding count: 1 (worst: S1, ADR 0033 overstates the revalidation rule)\nS1 | auto-fix | confidence 7 | docs/adr/0033-a-command-run-caches-its-pure-git-queries.md:31 | \"Each location read checks the stored directory again on the filesystem\" is wrong for the top-level location: Run.RootAt stores its answer and runs no validateCommonDir check; only CommonDir, AdminDir, and AdminPath revalidate. Fix: name the three directory-valued reads.\nGQ43 pass. GQ44 pass except S1.\nAdvice: FullObjectID hex loop repeats hexIdentity's loop; the ticket keeps hexIdentity unchanged.\n"
          },
          "axis": "Standards",
          "base": "f979bae6807f82e9bc4a259b78076c3170710597",
          "tip": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "gq-c1-spec-1",
          "performer": "claude:gq-c1-spec-20261028",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c1-spec-20261028",
            "digest": "sha256:9feb7f34cd046ac6fd645450436e16a5d0d0a00e54145ea3ce60621f70d5b140",
            "excerpt": "axis: Spec | verdict: findings | 2 findings | worst: P1\nP1 | ask-user | confidence 5 | internal/git/run.go:77-94 | Run.AdminPath checks the parent directory on the first read too; the nil run and the package form return the path. The spec row says the read checks the parent \"again\"; the flagged addition says \"each cached read\". Confirm the first-read check; if kept, the doc comment must state it.\nP2 | no-op | confidence 4 | specs/git-query-cache/spec.md:247-271 | The seam-cell citation rewrite sits outside the Writes: line; it is bookkeeping with no behavior change.\nEvery ticket 01 row has a test that asserts its behavior.\nAdvice: GQ9 moves the branch by commit, not update-ref; GQ19 checks only cat-file for ReadTreeFile; GQ23 refuses on the third read.\n"
          },
          "axis": "Spec",
          "base": "f979bae6807f82e9bc4a259b78076c3170710597",
          "tip": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "gq-c1-coverage-1",
          "performer": "claude:gq-c1-coverage-20261028",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "ab9506e1be46c9e86a6302b435dd4a45219d6a7e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c1-coverage-20261028",
            "digest": "sha256:ce2cae0997b2278751432de0c93ee5c8a9bad23ae5baabf5c5ced944d83e2f99",
            "excerpt": "axis: Coverage | verdict: findings | 4 findings | worst: C1\nC1 | auto-fix | confidence 8 | internal/git/run.go:182, :318 | AdminPath key (name, root) and object key (-C root) have no test that tells two keys apart; dropping name or root from the key survives (spec Edge inventory: spellings key as exact bytes).\nC2 | auto-fix | confidence 8 | internal/git/run.go:273 | IsAncestor purity guard: deleting the !FullObjectID(ancestor) term survives; no ref-ancestor, full-descendant case next to GQ15.\nC3 | auto-fix | confidence 8 | internal/git/run.go:237 | pureListing: removing the ^{tree} trim or the flag skip survives; no <fullid>^{tree}, flag, or two-revision case.\nC4 | auto-fix | confidence 6 | internal/worktree/land_identity.go:291 | abbreviatedIdentity 64-digit expansion (flagged addition) has no test; a 40-only check survives.\n"
          },
          "axis": "Coverage",
          "base": "f979bae6807f82e9bc4a259b78076c3170710597",
          "tip": "de8410c563e7ea8e2fd2d7a1b2bebac97fd646f1",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
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

## Chunk GQ-C1 review pickup

The review read the frozen pair f979bae68 to de8410c56. Repair cycles used: 0 of 2. All six accepted findings fall on ticket 01, so one combined repair charge carries them.

## Standards

Findings: 1. Worst issue: S1.

- S1, auto-fix, confidence 7: the ADR 0033 sentence "Each location read checks the stored directory again on the filesystem" is wrong for the top-level read. Only the common-directory, administration-directory, and administration-path reads check again (docs/adr/0033-a-command-run-caches-its-pure-git-queries.md:31).

## Spec

Findings: 2. Worst issue: P1.

- P1, auto-fix, confidence 5: the doc comment of `Run.AdminPath` says that the method answers as the package form does (internal/git/run.go:77-94). But a non-nil run checks the parent directory on each read, the first read included. Of three consultants, opus/high and fable/high chose "check every read", and sonnet/high chose "check only a stored answer". The tiebreak cites spec story 9 and the edge inventory, which bind the check to each read. The repair corrects the comment only. A nil run keeps the package behavior under story 7.
- P2, no-op, confidence 4: the seam-cell citation rewrite in the spec sits outside the `Writes:` line. The build rule puts each coverage-cell citation in the same commit as its cited test, so no repair target remains.

## Coverage

Findings: 4. Worst issue: C1.

- C1, auto-fix, confidence 8: no test separates two keys of `AdminPath` or of the object store (internal/git/run.go:182 and :318). So a key without the name or the root spelling survives.
- C2, auto-fix, confidence 8: no test reads a ref ancestor with a full-id descendant (internal/git/run.go:273; row GQ15). So a guard without the ancestor term survives.
- C3, auto-fix, confidence 8: no test lists a full id with the `^{tree}` suffix, a flag operand, or two revisions (internal/git/run.go:237). So a listing rule without the suffix trim or the flag skip survives.
- C4, auto-fix, confidence 6: no test expands a prefix to a 64-digit id through `abbreviatedIdentity`, so a 40-digit rule there survives (internal/worktree/land_identity.go:291; a flagged addition).
