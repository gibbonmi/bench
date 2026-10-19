# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/plan-commit/spec.md",
  "plan_digest": "sha256:e58e068c33f002cf04c2e560f13931b7ea46e5951ad50b6021fdc72f29c97a32",
  "implementation_session": "",
  "chunks": [
    {
      "id": "PC-C1",
      "base": "7b45186c91c875273fb4866096d8977b235c2bc5",
      "tip": "5daa190037d8cc2503ff1e3d01802c9e7cf886f5",
      "plan_digest": "sha256:e58e068c33f002cf04c2e560f13931b7ea46e5951ad50b6021fdc72f29c97a32",
      "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
      "acceptance_rows": [
        "PC1",
        "PC2",
        "PC3",
        "PC4",
        "PC5",
        "PC6",
        "PC7",
        "PC8",
        "PC9",
        "PC10",
        "PC11",
        "PC12",
        "PC13",
        "PC14",
        "PC15",
        "PC16",
        "PC17",
        "PC18",
        "PC37",
        "PC38",
        "PC19",
        "PC20",
        "PC21"
      ],
      "verification": [
        {
          "id": "t2-t2-preflight-packages",
          "performer": "claude:ft380-t02-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t02-author-20261019",
            "digest": "sha256:742deac3fb1c77c6c68a93adc1c83673d6ba6a78507007ad9e1ca60a09193ec4",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,37140,379\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,26652,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,15879,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,6336,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-preflight-packages",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "t2-t2-review-charge-omission",
          "performer": "claude:ft380-t02-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t02-author-20261019",
            "digest": "sha256:a4633625defb6bbc23b933e806b9406d23f426c064164e7f3d21e669b1b7bed1",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/review.go,swap,failed,2,TestCheckFloorRefusesTheReviewCharge,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,379\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,514\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,35815,379\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestCheckFloorRefusesTheReviewCharge,\"check_floor_command_test.go:87: review charge = (1):\",3\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceReviewRefusals/mismatched_pair,\"review_charge_test.go:260: review refusal = (0), want \\\"tip-current\\\" without a published handle:\",3\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-review-charge-omission",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "remove the verdict.Red refusal of the review charge; TestCheckFloorRefusesTheReviewCharge must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft380-t02-author-20261019",
              "digest": "sha256:a4633625defb6bbc23b933e806b9406d23f426c064164e7f3d21e669b1b7bed1",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/review.go,swap,failed,2,TestCheckFloorRefusesTheReviewCharge,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,379\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,514\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,35815,379\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestCheckFloorRefusesTheReviewCharge,\"check_floor_command_test.go:87: review charge = (1):\",3\n  github.com/gibbonmi/bench/internal/preflight,TestEvidenceReviewRefusals/mismatched_pair,\"review_charge_test.go:260: review refusal = (0), want \\\"tip-current\\\" without a published handle:\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-t2-root-conformance",
          "performer": "claude:ft380-t02-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t02-author-20261019",
            "digest": "sha256:b6174b51a50f5d0a66aa164d72c4cc9fde88e7e8a1cb5dbd4740d2cdff015df6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18144,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-coverage-package",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:56682f0df24036022fd225eaa6365299685e19df33283917485ec775dbdab2fe",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,2899,153\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-coverage-package",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "t1-t1-preflight-packages",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:947f8db7ec7ca54ef9c0b8a0d758d747ac93509975a4ec422a84482b3f815416",
            "excerpt": "packages[8]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,40467,379\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,27604,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,13703,177\n  github.com/gibbonmi/bench/internal/preflight/ownership/inventorytest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/staleness,pass,6757,70\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-preflight-packages",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "t1-t1-command-package",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:beeb77818880f445e150d20cc097c94190a7acd778800c39e5ee10ba2ad003ac",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19541,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-command-package",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-t1-gate-package",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:59305d10616e958be754da24dcd3b7e66c815fafccfe6d791662f64f176d881e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,31322,490\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-gate-package",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "t1-t1-diff-package",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:1a80acba813ea016e77a5518c0101014ae7d8441695fecfbbd423ab38afc2f12",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/diff,pass,7899,101\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-diff-package",
          "command": "bench test --package ./internal/diff",
          "exit_code": 0
        },
        {
          "id": "t1-t1-recordtest-package",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:7dd58c33b5a23a0efad8d9151e7542c7202e6f4f5a3c139ceff787447096a540",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-recordtest-package",
          "command": "bench test --package ./internal/reviewrecord/recordtest",
          "exit_code": 0
        },
        {
          "id": "t1-t1-system",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:8b1567774451180ceb8a0144e5e21c21f701a4774e00cb71e5e78a17dbd107a5",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  system,system,175,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/systemtest,pass,86809,175\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-system",
          "command": "bench test --check system",
          "exit_code": 0
        },
        {
          "id": "t1-t1-root-conformance",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:1ab301ff8f97a5387d58595af298ee448e660fd349134a94e816e7769e21e528",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17590,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-t1-unlisted-row-omission",
          "performer": "claude:ft380-t01-author-20261019",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "90f8a4ec255c9007eb7128a71347f0ffbebd2b21",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
            "digest": "sha256:83a11f8b88424c91713c02ec26fad249b283b18d49fbbdc443a486701d8fabae",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/chunk_rows.go,omit,failed,2,TestChunkRowsRedsAnUnlistedRow,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,379\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,498\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,35968,379\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestChunkRowsRedsAnUnlistedRow,\"chunk_rows_test.go:43: chunk-rows = {Check:chunk-rows Verdict:green Detail: Next: exempt:false}, want red\",1\n  github.com/gibbonmi/bench/internal/preflight,TestChunkRowsRefusesTheBuildCharge,\"chunk_rows_test.go:93: build charge = (0):\",4\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-unlisted-row-omission",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0,
          "probe": {
            "mutation": "drop the declared-row-to-cell direction in chunkRowsCheck; TestChunkRowsRedsAnUnlistedRow must fail",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft380-t01-author-20261019",
              "digest": "sha256:83a11f8b88424c91713c02ec26fad249b283b18d49fbbdc443a486701d8fabae",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/chunk_rows.go,omit,failed,2,TestChunkRowsRedsAnUnlistedRow,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight,all,passed,379\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,498\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,fail,35968,379\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight,TestChunkRowsRedsAnUnlistedRow,\"chunk_rows_test.go:43: chunk-rows = {Check:chunk-rows Verdict:green Detail: Next: exempt:false}, want red\",1\n  github.com/gibbonmi/bench/internal/preflight,TestChunkRowsRefusesTheBuildCharge,\"chunk_rows_test.go:93: build charge = (0):\",4\nskips[0]{package,test,reason}:\n"
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
