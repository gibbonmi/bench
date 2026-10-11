# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commit-merge-parent/spec.md",
  "plan_digest": "sha256:3b3b51e2c87d4449eb3044f3d9d1fb6c5eb414dd8b48e9d711f48822bfe6eda6",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MP1",
      "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
      "tip": "dfacb0be4c0870e191535bb104620168643d4bb1",
      "plan_digest": "sha256:3b3b51e2c87d4449eb3044f3d9d1fb6c5eb414dd8b48e9d711f48822bfe6eda6",
      "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
      "acceptance_rows": [
        "MP01",
        "MP02",
        "MP03",
        "MP04",
        "MP05",
        "MP06",
        "MP07",
        "MP08",
        "MP09",
        "MP10",
        "MP11",
        "MP12",
        "MP13",
        "MP14",
        "MP15",
        "MP16",
        "MP17",
        "MP18"
      ],
      "verification": [
        {
          "id": "t1-reader-focused-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-reader-focused.txt",
            "digest": "sha256:12c6ae34dd56ce695c24be4972fdf2a5af0a5073556169f582bd987fed2a783c",
            "excerpt": "$ bench test --package ./internal/git --run 'TestPendingMerge'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,90,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/git/checkout.go --omit 'return nil, fmt.Errorf(\"%q is not a regular file\", path)' --package ./internal/git --run '^TestPendingMergeRefusesUnreadableState$'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/git/checkout.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,^TestPendingMergeRefusesUnreadableState$,passed,6\n  github.com/gibbonmi/bench/internal/git,fail,20027,6\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/fifo,\"checkout_test.go:154: timeout: waited 20s for git.PendingMerge\",1\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/live_symbolic_link,\"checkout_test.go:154: PendingMerge = [\\\"f181b104448a2b9ed40d31670626888100605f71\\\"], <nil>; want no merge head and a read pending merge refusal\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-reader-focused",
          "command": "bench test --package ./internal/git --run 'TestPendingMerge'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the regular-file refusal after os.Lstat in git.PendingMerge in internal/git/checkout.go, so that a directory at MERGE_HEAD reaches the read. TestPendingMergeRefusesUnreadableState turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:cmp_t1/t1-reader-focused.txt",
              "digest": "sha256:12c6ae34dd56ce695c24be4972fdf2a5af0a5073556169f582bd987fed2a783c",
              "excerpt": "$ bench test --package ./internal/git --run 'TestPendingMerge'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,90,17\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/git/checkout.go --omit 'return nil, fmt.Errorf(\"%q is not a regular file\", path)' --package ./internal/git --run '^TestPendingMergeRefusesUnreadableState$'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/git/checkout.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,^TestPendingMergeRefusesUnreadableState$,passed,6\n  github.com/gibbonmi/bench/internal/git,fail,20027,6\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/fifo,\"checkout_test.go:154: timeout: waited 20s for git.PendingMerge\",1\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/live_symbolic_link,\"checkout_test.go:154: PendingMerge = [\\\"f181b104448a2b9ed40d31670626888100605f71\\\"], <nil>; want no merge head and a read pending merge refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-route-focused-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-route-focused.txt",
            "digest": "sha256:bc40eafeafbdf25967197092144637c51ef55e5a4203b1f2e18cab5ee50432c2",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestSourceMergePending|TestLandCommandPublicConflictRepairRequiresNewReviewedTip|TestPackageTestCountPin'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1076,4\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-route-focused",
          "command": "bench test --package ./internal/worktree --run 'TestSourceMergePending|TestLandCommandPublicConflictRepairRequiresNewReviewedTip|TestPackageTestCountPin'",
          "exit_code": 0
        },
        {
          "id": "t1-floor-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-floor.txt",
            "digest": "sha256:b731a182bd52a49ed551ce381c79454561a58859b3be5e0dbdbcdbe7a1326a40",
            "excerpt": "$ bench test --package ./internal/git --package ./internal/worktree --package ./internal/conformance --package ./cmd/bench\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19882,333\n  github.com/gibbonmi/bench/internal/conformance,pass,73782,1685\n  github.com/gibbonmi/bench/internal/git,pass,1485,165\n  github.com/gibbonmi/bench/internal/worktree,pass,68233,1395\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-floor",
          "command": "bench test --package ./internal/git --package ./internal/worktree --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-root-conformance.txt",
            "digest": "sha256:531d479d5c77879b434b93b1f2e03886103dd534e73ad2d375ff542cc1c6b784",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14325,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-git-plumbing-owner-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-git-plumbing-owner.txt",
            "digest": "sha256:6ad7235fb2ab0c5317bbfb2475458b52915aaae2106ce60077668df77c25826a",
            "excerpt": "$ bench test --check git-plumbing-owner\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,99,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t1-ordinary-build-census-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-ordinary-build-census.txt",
            "digest": "sha256:ab5f175d42fbc28d83842a85356b706de26ae3a1ab07315c3fdaf329351f2156",
            "excerpt": "$ bench test --check ordinary-build-census\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,335,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t1-skip-ownership-cmp_t1",
          "performer": "claude:cmp_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1/t1-skip-ownership.txt",
            "digest": "sha256:bd95de563ff079a0c3a8c60cb063827b6b8b8ece2f4720742001cc5e96a3feed",
            "excerpt": "$ bench test --check skip-ownership\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,dfacb0be4c0870e191535bb104620168643d4bb1,false\ncheck[1]{name,kind,tests_run,subjects}:\n  skip-ownership,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,48,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-skip-ownership",
          "command": "bench test --check skip-ownership",
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
