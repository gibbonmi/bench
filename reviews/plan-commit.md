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
