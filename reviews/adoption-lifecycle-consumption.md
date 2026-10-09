# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/adoption-lifecycle-consumption/spec.md",
  "plan_digest": "sha256:313d7279f6b6809daf20133edb5c43f3ca1b21f88bf5e601f1ad82d7f00bcc23",
  "implementation_session": "",
  "chunks": [
    {
      "id": "ALC-OWNER",
      "base": "aabd462d52349b00cfb3d66f67a6dd24891468d2",
      "tip": "12680fe8d5da397fff87f7d978a92e12ef7d43a6",
      "plan_digest": "sha256:313d7279f6b6809daf20133edb5c43f3ca1b21f88bf5e601f1ad82d7f00bcc23",
      "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
      "acceptance_rows": [
        "ALC02",
        "ALC03",
        "ALC04",
        "ALC05",
        "ALC31",
        "ALC32"
      ],
      "verification": [
        {
          "id": "t1-decision-preview",
          "performer": "claude:ft217_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1/04a-decision-preview.txt",
            "digest": "sha256:b2f7533476fd212f60f4401759169cc0439f5d2775284c11f1269eb94e0f7a92",
            "excerpt": "1. bench worktree exec ft217-build -- bench test --package ./internal/adopt --run 'Test(PlanLifecycleFromImmutableInventories|CompleteLifecycleDecisions|LifecycleDecisionOrder|LifecycleUpgradePreviewCompatibility|LifecycleUpgradeVersionGuards)'\nexit 0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,pass,29375,58\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n"
          },
          "requirement": "decision-preview",
          "command": "bench test --package ./internal/adopt --run 'Test(PlanLifecycleFromImmutableInventories|CompleteLifecycleDecisions|LifecycleDecisionOrder|LifecycleUpgradePreviewCompatibility|LifecycleUpgradeVersionGuards)'",
          "exit_code": 0
        },
        {
          "id": "t1-owner-floor-conformance",
          "performer": "claude:ft217_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1/04b-owner-floor-conformance.txt",
            "digest": "sha256:8984ee5fcf72d7908c9cec92559b40e7e0b77a5c7aa0a1d75036495f96db3dcb",
            "excerpt": "2. bench worktree exec ft217-build -- bench test --check root-conformance\nexit 0\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16751,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n"
          },
          "requirement": "owner-floor-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-owner-floor-adopt",
          "performer": "claude:ft217_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1/04c-owner-floor-adopt.txt",
            "digest": "sha256:892b5755e25c0788ebec94ac0349e554dd1b42e38bac4165695df7fc6a365f26",
            "excerpt": "3. bench worktree exec ft217-build -- bench test --package ./internal/adopt\nexit 0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,pass,112066,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n"
          },
          "requirement": "owner-floor-adopt",
          "command": "bench test --package ./internal/adopt",
          "exit_code": 0
        },
        {
          "id": "t1-owner-floor-cmd",
          "performer": "claude:ft217_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1/04d-owner-floor-cmd.txt",
            "digest": "sha256:6ffb1f36f4f2d18bc8244a26aa97d35d0ab0561c3e51060dc9d4b1ed1db94cd0",
            "excerpt": "4. bench worktree exec ft217-build -- bench test --package ./cmd/bench\nexit 0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13832,333\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-floor-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-owner-preview-probes",
          "performer": "claude:ft217_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1/p01-a-preserve-row-omission.txt",
            "digest": "sha256:a687772086612ad07d1ea673131d3d173482cb1e716939935286a346e81af1a1",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft217-build,a2e997fb27e0382d6c465a8190c14d21f056fd14,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,12\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/legacy,\"lifecycle_decision_test.go:91: decision 0 = {Path:added Kind:add Reason: ManifestAction:record Fingerprint:new Order:0}, want {Path:CLAUDE.md Kind:preserve Reason:requested-preserve ManifestAction:retain Fingerprint:keep Order:0}\",7\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preview-probes",
          "command": "bench probe internal/adopt/decision.go --omit 'decisions = append(decisions, preserved(name, \"requested-preserve\", manifestRetain, before))' --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Preserve-row omission, owner-order omission, and real Upgrade synthetic-to-exact preview swap must bite and restore. The ticket commit replaces the placeholder with the exact omitted source text.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft217_t1/p01-a-preserve-row-omission.txt",
              "digest": "sha256:a687772086612ad07d1ea673131d3d173482cb1e716939935286a346e81af1a1",
              "excerpt": "tree[1]{target,head,dirty}:\n  ft217-build,a2e997fb27e0382d6c465a8190c14d21f056fd14,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestCompleteLifecycleDecisions$,passed,12\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,12\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/legacy,\"lifecycle_decision_test.go:91: decision 0 = {Path:added Kind:add Reason: ManifestAction:record Fingerprint:new Order:0}, want {Path:CLAUDE.md Kind:preserve Reason:requested-preserve ManifestAction:retain Fingerprint:keep Order:0}\",7\nskips[0]{package,test,reason}:\n"
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
