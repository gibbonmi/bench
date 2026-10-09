# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/adoption-lifecycle-consumption/spec.md",
  "plan_digest": "sha256:a5fa48d1160a4f45c76b2dfe2e6c33e052685518fa0aa6f87d8478db4799d8dd",
  "implementation_session": "",
  "chunks": [
    {
      "id": "ALC-OWNER",
      "base": "aabd462d52349b00cfb3d66f67a6dd24891468d2",
      "tip": "518c52eec4fe9edf5f8f6baad50bc8404e5a9f32",
      "plan_digest": "sha256:a5fa48d1160a4f45c76b2dfe2e6c33e052685518fa0aa6f87d8478db4799d8dd",
      "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
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
        },
        {
          "id": "t1r1-decision-preview",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/04a-decision-preview.txt",
            "digest": "sha256:ba8fe694228fc2eb21c4fd56f3b28af5d0a92bac0a46ccc2314d05aaddf134a1",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,pass,209664,62\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n[exited with code 0]\n"
          },
          "requirement": "decision-preview",
          "command": "bench test --package ./internal/adopt --run 'Test(PlanLifecycleFromImmutableInventories|CompleteLifecycleDecisions|LifecycleDecisionOrder|LifecycleUpgradePreviewCompatibility|LifecycleUpgradeVersionGuards)'",
          "exit_code": 0
        },
        {
          "id": "t1r1-owner-floor-conformance",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/04b-owner-floor-conformance.txt",
            "digest": "sha256:7a66c5b9b48e4ace0d84aff7fa8ba1ceb32872ac542434aa992235d039d2262b",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,28145,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n[exited with code 0]\n"
          },
          "requirement": "owner-floor-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1r1-owner-floor-adopt",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/04c-owner-floor-adopt.txt",
            "digest": "sha256:5397f7108043a48030a095e1ee534960683c7ec87e6fd1f3c6e5dd20c4d1e053",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,pass,408780,165\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n[exited with code 0]\n"
          },
          "requirement": "owner-floor-adopt",
          "command": "bench test --package ./internal/adopt",
          "exit_code": 0
        },
        {
          "id": "t1r1-owner-floor-cmd",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/04d-owner-floor-cmd.txt",
            "digest": "sha256:f4ef3d8214bca576bea83d7af780001e6310e94f9af93ee9ffb85a6efca569a5",
            "excerpt": "tree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,31369,333\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n[exited with code 0]\n"
          },
          "requirement": "owner-floor-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1r1-owner-preview-probes",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/p-owner-preview-probes.txt",
            "digest": "sha256:64fd5954cd78db253c88f3c92f67435b0ec494d555ee99b73a386e4d8a4c4527",
            "excerpt": "Repair probes at 518c52ee, each through bench probe; every verdict bit, restored yes.\n\nC1: bench probe internal/adopt/decision.go --swap 'if o.FileKind != \"symlink\" && o.Perm != o.StagedPerm {' --with 'if o.Perm != o.StagedPerm {' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/link,\"lifecycle_decision_test.go:132: decision 22 = {Path:symlink-mode Kind:preserve Reason:modified-managed ManifestAction:retain Fingerprint:old Order:19}, want {Path:symlink-mode Kind:preserve Reason:converged ManifestAction:record Fingerprint… (255 bytes)\",2\nskips[0]{package,test,reason}:\n\nC2: bench probe internal/adopt/decision.go --swap 'case o.Recorded != \"\" && o.Observed != \"\" && o.Observed != o.Recorded:' --with 'case o.Recorded != \"\" && o.Observed != o.Recorded:' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/link,\"lifecycle_decision_test.go:132: decision 6 = {Path:claude-owned-special Kind:preserve Reason:modified-managed ManifestAction:retain Fingerprint:old Order:20}, want {Path:claude-owned-special Kind:preserve Reason:project-owned ManifestAction… (268 bytes)\",2\nskips[0]{package,test,reason}:\n\nC3: bench probe internal/adopt/decision.go --swap 'return !(o.Strict && o.FileKind == \"regular\" && o.StagedKind == \"regular\" && o.Perm != o.StagedPerm)' --with 'return !(o.Strict && o.Perm != o.StagedPerm)' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/link,\"lifecycle_decision_test.go:132: decision 19 = {Path:strict-link-dest Kind:preserve Reason:modified-managed ManifestAction:retain Fingerprint:old Order:21}, want {Path:strict-link-dest Kind:change Reason: ManifestAction:record Fingerprint:ne… (251 bytes)\",3\nskips[0]{package,test,reason}:\n\nC4: bench probe internal/adopt/decision.go --omit 'return nil, fmt.Sprintf(\"managed path %q is invalid\", name)' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/invalid_preserve,\"lifecycle_decision_test.go:196: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation(nil), Refusal:\\\"\\\"}, want the refusal \\\"managed path \\\\\\\\\\\"../x\\\\\\\\\\\" is invalid\\\"\",1\nskips[0]{package,test,reason}:\n\nP1: bench probe internal/adopt/decision.go --swap 'return LifecycleDecision{}, \"refused non-regular dropped asset: \" + o.Path' --with 'return LifecycleDecision{}, fmt.Sprintf(\"dropped asset %q has no fingerprint\", o.Path)' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_special,\"lifecycle_decision_test.go:196: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision(nil), Operations:[]adopt.Operation(nil), Refusal:\\\"dropped asset \\\\\\\\\\\"dir\\\\\\\\\\\" has no fingerprint\\\"}, want the refusal \\\"refused non-regular dr… (257 bytes)\",1\nskips[0]{package,test,reason}:\n\nS2: bench probe internal/adopt/upgrade.go --swap 'block, err := toon.TableTyped(upgradeTable, upgradeColumns,' --with 'block, err := toon.TableTyped(upgradeTable, []string{\"from\", \"to\", \"added\", \"altered\", \"removed\"},' --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/upgrade.go,swap,failed,11,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$,passed,28\nspilled{lines=19,bytes=4088,omitted_lines=10,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/32297bb2a2b17d4d68b8c88121928c21/1791549446919451104-9847b31ba544fbf4.out}\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/preserved_CLAUDE.md,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_adapter,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_file,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_inline,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preview-probes",
          "command": "bench probe internal/adopt/decision.go --omit 'decisions = append(decisions, preserved(name, \"requested-preserve\", manifestRetain, before))' --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Preserve-row omission, owner-order omission, and real Upgrade synthetic-to-exact preview swap must bite and restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft217_t1_repair1/p-owner-preview-probes.txt",
              "digest": "sha256:64fd5954cd78db253c88f3c92f67435b0ec494d555ee99b73a386e4d8a4c4527",
              "excerpt": "Repair probes at 518c52ee, each through bench probe; every verdict bit, restored yes.\n\nC1: bench probe internal/adopt/decision.go --swap 'if o.FileKind != \"symlink\" && o.Perm != o.StagedPerm {' --with 'if o.Perm != o.StagedPerm {' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/link,\"lifecycle_decision_test.go:132: decision 22 = {Path:symlink-mode Kind:preserve Reason:modified-managed ManifestAction:retain Fingerprint:old Order:19}, want {Path:symlink-mode Kind:preserve Reason:converged ManifestAction:record Fingerprint… (255 bytes)\",2\nskips[0]{package,test,reason}:\n\nC2: bench probe internal/adopt/decision.go --swap 'case o.Recorded != \"\" && o.Observed != \"\" && o.Observed != o.Recorded:' --with 'case o.Recorded != \"\" && o.Observed != o.Recorded:' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/link,\"lifecycle_decision_test.go:132: decision 6 = {Path:claude-owned-special Kind:preserve Reason:modified-managed ManifestAction:retain Fingerprint:old Order:20}, want {Path:claude-owned-special Kind:preserve Reason:project-owned ManifestAction… (268 bytes)\",2\nskips[0]{package,test,reason}:\n\nC3: bench probe internal/adopt/decision.go --swap 'return !(o.Strict && o.FileKind == \"regular\" && o.StagedKind == \"regular\" && o.Perm != o.StagedPerm)' --with 'return !(o.Strict && o.Perm != o.StagedPerm)' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/link,\"lifecycle_decision_test.go:132: decision 19 = {Path:strict-link-dest Kind:preserve Reason:modified-managed ManifestAction:retain Fingerprint:old Order:21}, want {Path:strict-link-dest Kind:change Reason: ManifestAction:record Fingerprint:ne… (251 bytes)\",3\nskips[0]{package,test,reason}:\n\nC4: bench probe internal/adopt/decision.go --omit 'return nil, fmt.Sprintf(\"managed path %q is invalid\", name)' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/invalid_preserve,\"lifecycle_decision_test.go:196: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision{}, Operations:[]adopt.Operation(nil), Refusal:\\\"\\\"}, want the refusal \\\"managed path \\\\\\\\\\\"../x\\\\\\\\\\\" is invalid\\\"\",1\nskips[0]{package,test,reason}:\n\nP1: bench probe internal/adopt/decision.go --swap 'return LifecycleDecision{}, \"refused non-regular dropped asset: \" + o.Path' --with 'return LifecycleDecision{}, fmt.Sprintf(\"dropped asset %q has no fingerprint\", o.Path)' --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$,passed,25\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,25\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestCompleteLifecycleDecisions/dropped_special,\"lifecycle_decision_test.go:196: PlanLifecycle() = adopt.LifecyclePlan{Decisions:[]adopt.LifecycleDecision(nil), Operations:[]adopt.Operation(nil), Refusal:\\\"dropped asset \\\\\\\\\\\"dir\\\\\\\\\\\" has no fingerprint\\\"}, want the refusal \\\"refused non-regular dr… (257 bytes)\",1\nskips[0]{package,test,reason}:\n\nS2: bench probe internal/adopt/upgrade.go --swap 'block, err := toon.TableTyped(upgradeTable, upgradeColumns,' --with 'block, err := toon.TableTyped(upgradeTable, []string{\"from\", \"to\", \"added\", \"altered\", \"removed\"},' --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/upgrade.go,swap,failed,11,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$,passed,28\nspilled{lines=19,bytes=4088,omitted_lines=10,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/32297bb2a2b17d4d68b8c88121928c21/1791549446919451104-9847b31ba544fbf4.out}\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/preserved_CLAUDE.md,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_adapter,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_file,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_inline,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,altered,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-owner-order-probe",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/p-s1-owner-order-omission.txt",
            "digest": "sha256:ff4c43ad6d0ef7aee393a1916c78ee7e960047dc163bab9300de25aad9b15245",
            "excerpt": "S1 owner-order omission at 518c52ee: bench probe internal/adopt/decision.go --swap 'sort.Slice(decisions, func(i, j int) bool { return decisions[i].Path < decisions[j].Path })' --with 'sort.SliceStable(decisions, func(i, j int) bool { return false })' --package ./internal/adopt --run '^TestLifecycleDecisionOrder$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleDecisionOrder$,passed,4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,5,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/execution,\"lifecycle_decision_test.go:215: reordered observations changed the plan:\",9\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/legacy,\"lifecycle_decision_test.go:231: reordered inventories changed the plan:\",3\nskips[0]{package,test,reason}:\n\nAuthor excerpt at a2e997fb (scratchpad/t1/p02-b-owner-order-omission.txt):\ntree[1]{target,head,dirty}:\n  ft217-build,a2e997fb27e0382d6c465a8190c14d21f056fd14,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleDecisionOrder$,passed,4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/execution,\"lifecycle_decision_test.go:194: reordered observations changed the plan:\",7\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/legacy,\"lifecycle_decision_test.go:210: reordered inventories changed the plan:\",3\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preview-probes",
          "command": "bench probe internal/adopt/decision.go --omit 'decisions = append(decisions, preserved(name, \"requested-preserve\", manifestRetain, before))' --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Preserve-row omission, owner-order omission, and real Upgrade synthetic-to-exact preview swap must bite and restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft217_t1_repair1/p-s1-owner-order-omission.txt",
              "digest": "sha256:ff4c43ad6d0ef7aee393a1916c78ee7e960047dc163bab9300de25aad9b15245",
              "excerpt": "S1 owner-order omission at 518c52ee: bench probe internal/adopt/decision.go --swap 'sort.Slice(decisions, func(i, j int) bool { return decisions[i].Path < decisions[j].Path })' --with 'sort.SliceStable(decisions, func(i, j int) bool { return false })' --package ./internal/adopt --run '^TestLifecycleDecisionOrder$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleDecisionOrder$,passed,4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,5,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/execution,\"lifecycle_decision_test.go:215: reordered observations changed the plan:\",9\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/legacy,\"lifecycle_decision_test.go:231: reordered inventories changed the plan:\",3\nskips[0]{package,test,reason}:\n\nAuthor excerpt at a2e997fb (scratchpad/t1/p02-b-owner-order-omission.txt):\ntree[1]{target,head,dirty}:\n  ft217-build,a2e997fb27e0382d6c465a8190c14d21f056fd14,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/decision.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleDecisionOrder$,passed,4\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/adopt,fail,4,4\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/execution,\"lifecycle_decision_test.go:194: reordered observations changed the plan:\",7\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleDecisionOrder/legacy,\"lifecycle_decision_test.go:210: reordered inventories changed the plan:\",3\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1r1-preview-swap-probe",
          "performer": "claude:ft217_t1_repair1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "a4ceabeedebed5984c9fa95ad56129fcd2a33cc8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft217_t1_repair1/p-s1-preview-swap.txt",
            "digest": "sha256:ee26100faddf1454b82169be51ec01109f3d58f62aaff5fb973e227e6320b12e",
            "excerpt": "S1 synthetic-to-exact preview swap at 518c52ee: bench probe internal/adopt/upgrade.go --swap 'if entry.kind == \"file\" {' --with 'if entry.kind != \"seed\" {' --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/upgrade.go,swap,failed,16,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$,passed,28\nspilled{lines=24,bytes=5250,omitted_lines=15,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/32297bb2a2b17d4d68b8c88121928c21/1791549537755320700-0aa162a85bb469cf.out}\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/preserved_CLAUDE.md,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,47,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_adapter,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_file,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,47,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_inline,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\nskips[0]{package,test,reason}:\n\nAuthor excerpt at a2e997fb (scratchpad/t1/p03-c-synthetic-to-exact-adapter-swap.txt):\ntree[1]{target,head,dirty}:\n  ft217-build,a2e997fb27e0382d6c465a8190c14d21f056fd14,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/upgrade.go,swap,failed,12,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$,passed,28\nspilled{lines=20,bytes=4329,omitted_lines=11,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/32297bb2a2b17d4d68b8c88121928c21/1791546189834605584-10fa61093e9f9043.out}\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/preserved_CLAUDE.md,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_adapter,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,45,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_file,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_inline,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "owner-preview-probes",
          "command": "bench probe internal/adopt/decision.go --omit 'decisions = append(decisions, preserved(name, \"requested-preserve\", manifestRetain, before))' --package ./internal/adopt --run '^TestCompleteLifecycleDecisions$'",
          "exit_code": 0,
          "probe": {
            "mutation": "Preserve-row omission, owner-order omission, and real Upgrade synthetic-to-exact preview swap must bite and restore.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft217_t1_repair1/p-s1-preview-swap.txt",
              "digest": "sha256:ee26100faddf1454b82169be51ec01109f3d58f62aaff5fb973e227e6320b12e",
              "excerpt": "S1 synthetic-to-exact preview swap at 518c52ee: bench probe internal/adopt/upgrade.go --swap 'if entry.kind == \"file\" {' --with 'if entry.kind != \"seed\" {' --package ./internal/adopt --run '^TestLifecycleUpgradePreviewCompatibility$'\ntree[1]{target,head,dirty}:\n  ft217-build,518c52eec4fe9edf5f8f6baad50bc8404e5a9f32,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/upgrade.go,swap,failed,16,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$,passed,28\nspilled{lines=24,bytes=5250,omitted_lines=15,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/32297bb2a2b17d4d68b8c88121928c21/1791549537755320700-0aa162a85bb469cf.out}\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/preserved_CLAUDE.md,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,47,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_adapter,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_file,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,47,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_inline,\"lifecycle_preview_test.go:114: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\nskips[0]{package,test,reason}:\n\nAuthor excerpt at a2e997fb (scratchpad/t1/p03-c-synthetic-to-exact-adapter-swap.txt):\ntree[1]{target,head,dirty}:\n  ft217-build,a2e997fb27e0382d6c465a8190c14d21f056fd14,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/adopt/upgrade.go,swap,failed,12,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/adopt,^TestLifecycleUpgradePreviewCompatibility$,passed,28\nspilled{lines=20,bytes=4329,omitted_lines=11,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/32297bb2a2b17d4d68b8c88121928c21/1791546189834605584-10fa61093e9f9043.out}\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/preserved_CLAUDE.md,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,0,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_adapter,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,45,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_file,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\n  github.com/gibbonmi/bench/internal/adopt,TestLifecycleUpgradePreviewCompatibility/unowned_inline,\"lifecycle_preview_test.go:115: Upgrade --check = (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,46,0\\\\\\\\n\\\", \\\"\\\"), want (0, \\\"upgrade[1]{from,to,added,changed,removed}:\\\\\\\\n  1.0.0,1.0.1,1,0,0\\\\\\\\n\\\", \\\"\\\")\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "owner-standards-r1",
          "performer": "claude:ft217_owner_standards_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft217_owner_standards_r1/review-owner-standards-r1.txt",
            "digest": "sha256:83954ff811c0aee02408de96521d97d8377f8c7e5a373e3727b19ce9f25f4db2",
            "excerpt": "ALC-OWNER Standards axis, round 1. Session claude:ft217_owner_standards_r1, opus / medium. Frozen pair aabd462d..12680fe8, evidence sha256:58865964 current=true.\n\nFindings: 6. Worst: S1, most independently authored expectations have no recorded red in the review record.\n\nHard violations\nS1. AGENTS.md one-source rule: an independent expectation is duplicated knowledge unless its red is recorded and demonstrated. The record holds one bite (requested-preserve omission, reviews/adoption-lifecycle-consumption.md:109) and no red for the owner-order omission or the synthetic-to-exact preview swap that ticket 01 lines 46 and 52-53 require. Affected: TestLifecycleDecisionOrder, TestLifecycleUpgradePreviewCompatibility. Confidence 6. auto-fix: record the two probes as verification entries.\nS2. Same rule: previewTable (lifecycle_preview_test.go:20) copies the table name \"upgrade\" and the column list {\"from\",\"to\",\"added\",\"changed\",\"removed\"} from upgrade.go:69 with no recorded header red. Confidence 5. auto-fix: share one header source with upgrade.go, or record a header red.\n\nJudgment calls\nS3. Reason literals repeat inside decision.go (\"modified-managed\" 217/272, \"project-owned\" 219/274, \"absent\" 283/299, \"refused\" 297/301) while Kind and ManifestAction are constants; the drifting duplicate is the reason vocabulary across link_transaction.go:166-168,197,231, unlink.go:300-303, broker.go:156-158. Ticket line 62 accepts this until callers consume Decision.Reason. Confidence 5. no-op now; carry to ALC-LINK and ALC-UNLINK Standards review.\nS4. converged and ownedClean repeat the mode rule (link_stage.go:218) and the strict-mode rule (link_stage.go:312); ticket lines 16-17, 21, 62 make this transitional. Confidence 6. no-op under the ticket keep decision; carry forward so later chunks remove the caller copy.\nS5. upgrade.go:127-130 comment argues correctness (\"so it decides nothing a second time\"), against craft-comments. Confidence 4. auto-fix: state only the constraint.\nS6. The ALC-OWNER probe prose in the plan fence still says \"The ticket commit replaces the placeholder with the exact omitted source text\" after the placeholder is gone; history, not current state (invariant 3). Confidence 5. auto-fix.\n\nAdvice (no confidence): ticket 01 Writes omits spec.md though the ticket commit edits it; previewManifestPath adds one more \".bench/link-manifest.tsv\" literal; LifecycleObservation has 20 fields; fixture reuse is correct (consumerRepo, linkConsumer, readFile, writeFixtureFile); 400-line rule passes (349, 240, 224, 230).\n\nCursors fetched: metadata summary; v1.5886…a07.m.0.0; s15 page 0; v1.5886…a07.s.15.1. Not fetched: s15 page 2, s1. Examined: the full diff and its log, craft-review, finding-discipline, craft-comments, ticket 01, link_stage.go:196-318, compatibility_test.go:100-244, link_transaction_test.go:75-114, reason and manifest-path literal sweeps, the review record probe sections. Not read: projects/benchkit.md, setup_test.go body, coverage source s16, smell-baseline reference. Tree clean; wrote nothing.\n"
          },
          "axis": "Standards",
          "base": "aabd462d52349b00cfb3d66f67a6dd24891468d2",
          "tip": "12680fe8d5da397fff87f7d978a92e12ef7d43a6",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5",
            "S6"
          ],
          "supersedes": []
        },
        {
          "id": "owner-spec-r1",
          "performer": "claude:ft217_owner_spec_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft217_owner_spec_r1/review-owner-spec-r1.txt",
            "digest": "sha256:abb48516d36a24cc0252713dd36b555a6323b61a3e67cabb6c66886042d83ff1",
            "excerpt": "ALC-OWNER Spec axis, round 1. Session claude:ft217_owner_spec_r1, opus / medium. Frozen pair aabd462d..12680fe8, evidence sha256:58865964 current=true.\n\nFindings: 3. None blocks the chunk. Worst: P1, owner-emitted hard refusal text differs from the caller diagnostics.\n\nP1. Spec line 261 \"Keep hard checks at their current callers and in their current order\"; the table puts ALC11/ALC19 at \"caller hard refusal\". The owner returns its own whole-plan refusal strings (decision.go:212, :251, :280-285, unknown action, duplicate). The withdrawal text \"refused manifest path %q\" (decision.go:281) differs from the caller \"refused manifest path: \"+row.rel (link_transaction.go:214), and the owner merges the caller special-file refusal (link_transaction.go:218) into \"dropped asset ... has no fingerprint\". Nothing consumes it in this chunk. Confidence 5. ask-user: should the owner signal a refusal kind only, with the caller keeping the text? Orchestrator disposition under --auto-approve: auto-fix the text to the caller's exact bytes now; the single-source question carries to the ALC-LINK charge, where ticket 02 removes one copy.\nP2. The ALC08 owner-half seam cell still reads \"planned internal/adopt/lifecycle_decision_test.go (`TestCompleteLifecycleDecisions`)\" (spec.md:417) though the test exists (lifecycle_decision_test.go:98, :106). Confidence 7. auto-fix.\nP3. The owner-preview-probes probe text still says \"The ticket commit replaces the placeholder...\" (spec.md:704) after the command carries the exact text. Confidence 6. auto-fix.\n\nGraded with no finding: the observation field list (symlink target merged into ResolvedEquivalent/CanonicalAdapter; special file as Present with FileKind \"other\") satisfies spec.md:162-170; every partition table row traces (seed parent before presence skip, converged before symlink-parent refusal, converged/ownedClean/strict match link_stage.go:206-230 and :298-317, CLAUDE.md precedence matches link_transaction.go:193-199, unlink residuals match unlink.go:118-141, seed old-row retention matches link_transaction.go:201-206); Operations from add/change/remove only; nil vs present-empty Execution tested; manifest dispositions match; upgrade.go comment-only change satisfies ticket line 28 because upgradePlanCounts counts decision.Operations (upgrade.go:152-162) with sentinels and the CLAUDE.md Preserve unchanged; decision_test.go unchanged; ALC03/04/31/32 cells cited; ALC05 refusal order kept.\n\nAdvice (no confidence): claudeDecision returns change for a known-form CLAUDE.md equal to the staged form, ticket 03 confirms it publishes nothing new; executionDecisions checks duplicates in slice order; the ALC-LINK plan entries were edited in the ticket commit (dropped-probes added), check the bench learning record.\n\nCursors fetched: summary page; v1.5886…a07.m.0.0; s15 page 0. Examined: whole spec, ticket 01, the delta, decision.go, both new tests, upgrade.go:30-164, link_transaction.go:90-239, link_stage.go:55-110 and :206-317, unlink.go:110-159, decision_test.go, finding-discipline.md, bench coverage rows. Not read: craft-review SKILL.md itself.\n"
          },
          "axis": "Spec",
          "base": "aabd462d52349b00cfb3d66f67a6dd24891468d2",
          "tip": "12680fe8d5da397fff87f7d978a92e12ef7d43a6",
          "finding_ids": [
            "P1",
            "P2",
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "owner-coverage-r1",
          "performer": "claude:ft217_owner_coverage_r1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2e62aeea3b8f964a772e0ae6da4843d26c1daaeb",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft217_owner_coverage_r1/review-owner-coverage-r1.txt",
            "digest": "sha256:b5530cb1b87d72c6babb9325fc2930bcf2fddbf1722364e7c75b7ea1a0a4f18a",
            "excerpt": "ALC-OWNER Coverage axis, round 1. Session claude:ft217_owner_coverage_r1, opus / medium. Frozen pair aabd462d..12680fe8, evidence sha256:58865964 current=true.\n\nFindings: 4. Worst: finding 1, the convergence mode guard can be removed and no cited test turns red.\n\n1. Symlink destination, regular staged entry, equal fingerprint, different Perm (owned). Dropping `o.FileKind != \"symlink\" &&` from converged (decision.go:233) returns change instead of converged-preserve, against link_stage.go:218. Missing: a symlink-destination converged fixture in linkObservations (ALC08 decision level). Confidence 8. auto-fix.\n2. CLAUDE.md with a recorded row and a special or unfingerprintable destination (Recorded set, Observed \"\"). Dropping `o.Observed != \"\" &&` (decision.go:271) turns project-owned/omit into modified-managed/retain, against link_transaction.go:195. Missing: an owned claude-special fixture in TestCompleteLifecycleDecisions. Confidence 8. auto-fix.\n3. Strict repair where the destination or staged entry is not a regular file and modes differ. Shortening ownedClean (decision.go:245) to `!(o.Strict && o.Perm != o.StagedPerm)` reports a symlink or adapter as modified-managed, against link_stage.go:312. Missing: a strict fixture with a non-regular destination (ALC09, ALC25 decision level). Confidence 7. auto-fix.\n4. Legacy Preserve entry with an invalid path (`../x`). Omitting the Preserve validation return (decision.go:133-134) accepts hostile input; ALC05 and the edge row \"Invalid immutable legacy input\" need an invalid-Preserve case and a refusal-order case (Current before Desired before Preserve) in the new decision test; TestPlanLifecycleFromImmutableInventories stays unchanged. Confidence 8. auto-fix.\n\nProbes (all --package ./internal/adopt --run '^Test(CompleteLifecycleDecisions|LifecycleDecisionOrder|PlanLifecycleFromImmutableInventories)$', each verdict silent, 21 tests, restored=yes):\n1 --swap 'if o.FileKind != \"symlink\" && o.Perm != o.StagedPerm {' --with 'if o.Perm != o.StagedPerm {'\n2 --swap 'case o.Recorded != \"\" && o.Observed != \"\" && o.Observed != o.Recorded:' --with 'case o.Recorded != \"\" && o.Observed != o.Recorded:'\n3 --omit 'return nil, fmt.Sprintf(\"managed path %q is invalid\", name)'\n4 --swap 'return !(o.Strict && o.FileKind == \"regular\" && o.StagedKind == \"regular\" && o.Perm != o.StagedPerm)' --with 'return !(o.Strict && o.Perm != o.StagedPerm)'\n\nPassing checks: cited test names exist; decision_test.go unchanged across the pair; write set matches the Writes line plus spec.md; edge inventory rows covered (manifest variants, multiple preserves, Preserve-only path, Preserve precedence, empty vs nil Execution, earliest caller refusal).\n\nAdvice (no confidence): an empty inventory with a bad action returns an empty plan before the action check, untested; the Escapes fact is consulted only for withdrawal and unlink rows, so the ticket 02 caller keeps the desired-row escape check.\n\nCursors fetched: evidence summary page only. git status clean at the end.\n"
          },
          "axis": "Coverage",
          "base": "aabd462d52349b00cfb3d66f67a6dd24891468d2",
          "tip": "12680fe8d5da397fff87f7d978a92e12ef7d43a6",
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
  },
  "amendments": [
    {
      "from": "sha256:313d7279f6b6809daf20133edb5c43f3ca1b21f88bf5e601f1ad82d7f00bcc23",
      "to": "sha256:270f819ab599b99b9a4ccb3d3d39e4c50efa5c6228e8e1e594ce1b03164a65f5",
      "chunk_ids": {
        "ALC-OWNER": [
          "ALC-OWNER"
        ]
      }
    },
    {
      "from": "sha256:270f819ab599b99b9a4ccb3d3d39e4c50efa5c6228e8e1e594ce1b03164a65f5",
      "to": "sha256:a5fa48d1160a4f45c76b2dfe2e6c33e052685518fa0aa6f87d8478db4799d8dd",
      "chunk_ids": {
        "ALC-OWNER": [
          "ALC-OWNER"
        ]
      }
    }
  ]
}
```
