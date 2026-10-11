# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commit-merge-parent/spec.md",
  "plan_digest": "sha256:26c4daa7f1486d87e7341202b339d5fe7ef02fdf1a4d1fdfbc7b0f8f723304f3",
  "implementation_session": "",
  "chunks": [
    {
      "id": "MP1",
      "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
      "tip": "a329d3ab9339c2c88dad79d95f48de811e0e0737",
      "plan_digest": "sha256:06b9db25ba9dba1eee9e58f5e427693c6ff3473a78c4ce1e7ad3e6aa8c2f09dd",
      "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
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
        },
        {
          "id": "t1-reader-focused-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-reader-focused.txt",
            "digest": "sha256:7ad4cb9031aaa92759916d77c7d0c254f4c999a4d7ebc88f2f894cae1f72038d",
            "excerpt": "$ bench test --package ./internal/git --run 'TestPendingMerge'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,120,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/git/checkout.go --omit 'return nil, fmt.Errorf(\"%q is not a regular file\", path)' --package ./internal/git --run '^TestPendingMergeRefusesUnreadableState$'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/git/checkout.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,^TestPendingMergeRefusesUnreadableState$,passed,6\n  github.com/gibbonmi/bench/internal/git,fail,20026,6\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/fifo,\"checkout_test.go:163: timeout: waited 20s for git.PendingMerge\",1\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/live_symbolic_link,\"checkout_test.go:163: PendingMerge = [\\\"1bf766ee2035832490c56523abfe04b50676b93c\\\"], <nil>; want no merge head and a read pending merge refusal\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:cmp_t1_rep1/t1-reader-focused.txt",
              "digest": "sha256:7ad4cb9031aaa92759916d77c7d0c254f4c999a4d7ebc88f2f894cae1f72038d",
              "excerpt": "$ bench test --package ./internal/git --run 'TestPendingMerge'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,pass,120,23\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/git/checkout.go --omit 'return nil, fmt.Errorf(\"%q is not a regular file\", path)' --package ./internal/git --run '^TestPendingMergeRefusesUnreadableState$'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/git/checkout.go,omit,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,^TestPendingMergeRefusesUnreadableState$,passed,6\n  github.com/gibbonmi/bench/internal/git,fail,20026,6\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/fifo,\"checkout_test.go:163: timeout: waited 20s for git.PendingMerge\",1\n  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/live_symbolic_link,\"checkout_test.go:163: PendingMerge = [\\\"1bf766ee2035832490c56523abfe04b50676b93c\\\"], <nil>; want no merge head and a read pending merge refusal\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-route-focused-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-route-focused.txt",
            "digest": "sha256:e956c541df1c698d298a3ae9e3ea763865779ae9185fc7039fef900a1d3cc94a",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestSourceMergePending|TestLandCommandPublicConflictRepairRequiresNewReviewedTip|TestPackageTestCountPin'\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,1081,7\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-route-focused",
          "command": "bench test --package ./internal/worktree --run 'TestSourceMergePending|TestLandCommandPublicConflictRepairRequiresNewReviewedTip|TestPackageTestCountPin'",
          "exit_code": 0
        },
        {
          "id": "t1-floor-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-floor.txt",
            "digest": "sha256:89095dc78a8e85390c4012629036364adb2eaa68370041924fbab5d2b416a4ea",
            "excerpt": "$ bench test --package ./internal/git --package ./internal/worktree --package ./internal/conformance --package ./cmd/bench\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19328,333\n  github.com/gibbonmi/bench/internal/conformance,pass,71380,1685\n  github.com/gibbonmi/bench/internal/git,pass,1410,171\n  github.com/gibbonmi/bench/internal/worktree,pass,66134,1398\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-floor",
          "command": "bench test --package ./internal/git --package ./internal/worktree --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-root-conformance.txt",
            "digest": "sha256:5101cffaf168b4be4452bed5c6a1f4510b1132376c3b8211d596c149a7dc954b",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15725,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-git-plumbing-owner-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-git-plumbing-owner.txt",
            "digest": "sha256:4e6fc1df11319885edb7acdb0ed3db9e5503f3858037683b08040712168bd5a6",
            "excerpt": "$ bench test --check git-plumbing-owner\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,96,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t1-ordinary-build-census-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-ordinary-build-census.txt",
            "digest": "sha256:3c208c334d71b82937e934334f1d757ebdf695074c2d2b6e9ae0b4a1004e27bd",
            "excerpt": "$ bench test --check ordinary-build-census\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,318,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t1-skip-ownership-cmp_t1_rep1",
          "performer": "claude:cmp_t1_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t1_rep1/t1-skip-ownership.txt",
            "digest": "sha256:1620a53acc1838b37297f8ba81926ea90b6bb9622a28902e523f7f35a512e59e",
            "excerpt": "$ bench test --check skip-ownership\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a329d3ab9339c2c88dad79d95f48de811e0e0737,false\ncheck[1]{name,kind,tests_run,subjects}:\n  skip-ownership,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,46,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-skip-ownership",
          "command": "bench test --check skip-ownership",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "mp1-standards-r1",
          "performer": "claude:cmp_mp1_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp1_standards",
            "digest": "sha256:9cc2df94d269c58c6ff71a602bde717d9e7e82a190750ddb2f0ca27fcb4148c2",
            "excerpt": "MP1 Standards axis, claude:cmp_mp1_standards, opus/high, one iteration, frozen pair 342975e3..f7b5027b, evidence sha256:50ef098f0dcf0b23ccd0b2b9e0bb78b1cba498584daecb6f7b9dacf1a4dd33fc current=true.\n\nFindings: 2, both judgment-level auto-fix. Worst issue: S2.\n\nS1 (auto-fix, confidence 5). Rule: AGENTS.md \"Knowledge duplication is a defect ... a fixture harness pasted N times\"; smell Duplicated Code. The new helper newLinkedCheckout at internal/git/checkout_test.go:17-24 creates a repository and runs runGit(t, root, \"worktree\", \"add\", \"-q\", \"-b\", \"linked\", linked); the package already has linkedWorktree at internal/git/admin_readers_test.go:35-40 for the same fact. No new test reads the branch \"linked\"; the merges in TestPendingMergeReadsEveryHead run on a detached HEAD as well. Fix: newLinkedCheckout returns newRepo(t) and linkedWorktree(t, root).\n\nS2 (auto-fix, confidence 5). Rule: the AGENTS.md exception for an independently authored test expectation requires a named mutation that turns the gate red and a recorded demonstration of that red; spec.md:618 \"Each independent expectation then records one red\". requirePendingMergeRefusal at checkout_test.go:57 restates the production prefix \"read pending merge:\" (checkout.go:30). The independence is necessary (the prefix is a promised label, spec.md:152, :530). The review record's only probe omits the regular-file refusal, so no recorded run shows a prefix mutation red. Fix: record a prefix-swap red; no code change.\n\nJudgments with no finding: pendingMergeWithin against worktreesWithin is incidental text (the deadline and verdict facts each have one source in internal/bounds; the goroutine-select shape appears at six sites in the package; ticket line 46 directs \"as worktreesWithin does\"). git-plumbing-owner: only checkout_test.go:29 spells --git-path, inside internal/git and a test file, which GR28 exempts. Budgets hold: checkout_test.go 180, admin_readers_test.go 143, land_journey_test.go 422 with one line changed in place, parallel_census_test.go 1096 with the pin changed in place to 773. Comments follow craft-comments; no FT<n> story tag. The spec seam-cell flip inside dfacb0be is required by map-discipline.md:126. PendingMerge/mergeHeads is not a Middle Man; no duplicate admin reader or hex validator in internal/git.\n\nAdvice (no id, no confidence): mergeHeadPath could call independentRevParse (admin_readers_test.go:19); the admin_readers_test.go header rows omit MP15; the fullCommitID hex pre-check's purpose (keep option-shaped text out of the rev-parse argv) is unstated; a later light-path ticket could add one bounded-call helper for the six deadline-select sites.\n\nExamined: the full frozen delta; the commit attribution; the ticket; AGENTS.md, .bench/BENCH.md, projects/benchkit.md Seams and Gate, craft-review and finding-discipline.md, craft-comments; testhelpers_test.go, admin_readers_test.go, worktree_admin_hostile_test.go:200-232, git.go:218-234, land_journey_test.go:285-319, spec lines 520-624, git_plumbing_owner_test.go, map-discipline.md.\nEvidence cursors fetched: the diff spill (670 lines, all) and the log spill (18 lines, all); no preflight evidence page beyond --check-current.\n"
          },
          "axis": "Standards",
          "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
          "tip": "dfacb0be4c0870e191535bb104620168643d4bb1",
          "finding_ids": [
            "S1",
            "S2"
          ],
          "supersedes": []
        },
        {
          "id": "mp1-spec-r1",
          "performer": "claude:cmp_mp1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp1_spec",
            "digest": "sha256:9f70008937ddd927be0f5b33402eb858bba5611760e599a1262bc799c94b1e82",
            "excerpt": "MP1 Spec axis, claude:cmp_mp1_spec, opus/high, one iteration, frozen pair 342975e3..f7b5027b, evidence sha256:50ef098f0dcf0b23ccd0b2b9e0bb78b1cba498584daecb6f7b9dacf1a4dd33fc current=true.\n\nFindings: 0. Worst issue: none.\n\nRow audit: MP01 to MP18 met. MP01-MP04 TestPendingMergeReadsEveryHead (absent, merge, octopus, unterminated last line; linked worktree, real git merge --no-commit --no-ff). MP05-MP09 TestPendingMergeRefusesUnreadableState (outside every repository, directory, FIFO under bounds.TestDeadline with capability.Fifo skip, dangling link, live link). MP10-MP14 TestPendingMergeRefusesMalformedState (empty, abbreviated ID, main as a real branch, unknown object, tree; each asserts the prefix and no heads). MP15 TestSourceMergePendingReadsTheSharedOwner (t.Parallel, empty MERGE_HEAD through mustAdminPath, expects false). MP16 unchanged TestSourceMergePendingIsFalseWhenUndecided. MP17, MP18 unchanged LRS6 and LRS22 in TestLandCommandPublicConflictRepairRequiresNewReviewedTip (land_journey_test.go:299-318). Every MP01-MP15 seam cell uses the citation form and names an existing test (spec.md:309-323).\n\nRequirement trace: AdminPath resolution (spec.md:141) met; Lstat before read, absent gives nil,nil (142-143) met; non-regular refuses before read (144) met; optional final newline and empty refusal (147, 150) met, \"\\n\" alone refuses; lowercase hex plus ResolveCommit unchanged (148) met: abbreviated, main, tree, unknown object, and tag all refuse; file order met; read pending merge: prefix wraps every error including AdminPath failure; sourceMergePending returns err == nil && len(heads) > 0 (land_refusal.go:80-83), doc comment names the reader and keeps the undecided policy; os.ReadFile and the os and path/filepath imports are gone. No scope creep: the only new export is git.PendingMerge; worktreeTestCount 772 to 773; the LRS22 comment edit is an in-place change that ticket line 63 allows.\n\nAdvice (no id, no confidence): the spec's \"The pending-merge owner\" cites land_refusal.go:82-91 and :78-81 in the present tense, and the function now sits at :76-83; the MP15 why-cell describes the removed read. The helper mergeHeadPath uses --path-format=absolute, which resolves an existing symlink; safe because every caller resolves the path before a link is planted.\n\nExamined: the frozen diff; spec.md 1-160, 305-324, 372-397, 525-551; the ticket in full; land_refusal.go:1-110, admin_readers_test.go, land_journey_test.go:270-324, git.go:210-234, worktree_admin.go:195-229; checkout.go and checkout_test.go through the diff; craft-review SKILL.md:40-109 and finding-discipline.md; the bench coverage projection.\nEvidence cursors fetched: the diff response, the spec heading map response, the coverage summary response, and the coverage rows response lines 1-25; no preflight evidence page beyond --check-current.\n"
          },
          "axis": "Spec",
          "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
          "tip": "dfacb0be4c0870e191535bb104620168643d4bb1",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "mp1-coverage-r1",
          "performer": "claude:cmp_mp1_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "13b4c96a6f5f097d6ee83647525079d6c6f7c401",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp1_coverage",
            "digest": "sha256:aff48426f5c7b10edb0450a39f468d725befe92f4abe4d81889d619d3ab42e3c",
            "excerpt": "MP1 Coverage axis, claude:cmp_mp1_coverage, opus/high, one iteration, frozen pair 342975e3..f7b5027b, evidence sha256:50ef098f0dcf0b23ccd0b2b9e0bb78b1cba498584daecb6f7b9dacf1a4dd33fc current=true.\n\nFindings: 5, all auto-fix. Worst issue: C2.\n\nEnumeration: the input family is the MERGE_HEAD bytes Git writes (one lowercase full commit ID and \"\\n\" per line, argument order), a hand edit (unterminated last line, CRLF, padding whitespace, extra blank lines), a planted special file (directory, FIFO, dangling or live symlink), a planted non-commit object ID (tree, annotated tag that peels to a commit), and the linked worktree's gitdir: pointer (absolute or relative). Write set: every changed file is inside the ticket Writes: line except the orchestrator plan and record files. Consumers capture s17: 22 rows, all touched=true.\n\nC1 (auto-fix, confidence 8). Relative gitdir: pointer. spec.md:379 says \"MP01 covers the relative gitdir: pointer\", but newLinkedCheckout runs a plain git worktree add (checkout_test.go:22), which writes an absolute gitdir: on Git 2.43.0; no test runs the reader through a relative pointer (projects/benchkit.md:214-221 names that input). Expected break unknown. Needed: a TestPendingMergeReadsEveryHead step that rewrites linked/.git to a relative gitdir: and expects the merge head, or amend line 379.\nC2 (auto-fix, confidence 7). The octopus and unterminated-last-line steps put two heads in one file order whose sort position depends on random SHAs (checkout_test.go:103-109). An ascending-sort probe was silent; a descending-sort probe bit in a run where file order was ascending. A sorting reader breaks MP03 and E3 and MP21 inherits the wrong parent order. Needed: an MP03 fixture whose file order differs from both sorted orders (two octopus steps with reversed argument order, or three heads in a non-monotonic order).\nC3 (auto-fix, confidence 7). A full annotated-tag object ID in MERGE_HEAD peels through ^{commit} to a commit; only commit == line (checkout.go:78) refuses it, and no row plants a tag. The mutant err == nil && len(commit) == len(line) kept MP11-MP14 green (silent). Needed: an \"annotated tag\" case in TestPendingMergeRefusesMalformedState (spec.md:149 rule).\nC4 (auto-fix, confidence 6). \"sha \\n\", \"sha\\r\\n\", \"sha\\n\\n\": the reader refuses them (checkout.go:57, :74) but no case pins the refusal; swapping TrimSuffix for TrimSpace was silent. Needed: CRLF, trailing-space, and double-newline cases in TestPendingMergeRefusesMalformedState.\nC5 (auto-fix, confidence 5). A copy of the reader in land_refusal.go (AdminPath + os.ReadFile + len(data) > 0) passes MP15 (empty file), MP16, LRS6, and LRS22; the ticket acceptance \"land_refusal.go holds no read of MERGE_HEAD of its own\" has no test. Needed: an MP15 subtest with a non-empty malformed MERGE_HEAD (main\\n, a tree ID, or a live symlink to a valid record) that expects false; a subtest leaves worktreeTestCount unchanged.\n\nAdvice (no id, no confidence): removing the lowercase-hex prefilter (checkout.go:74) was silent and the mutant is equivalent; a comment that states its purpose (keep option-like lines away from git) would stop a later removal. No test calls PendingMerge in a primary checkout; the Won't handle line at spec.md:397 and AdminPath's own tests own that branch.\n\nProbes (all internal/git/checkout.go, --package ./internal/git, restored yes): TrimSuffix to TrimSpace, run ^TestPendingMerge: silent. Ascending sort after append, run ^TestPendingMergeReadsEveryHead$: silent. Descending sort, same run: bit, 2 failed (octopus, unterminated last line). Prefilter condition to `if line == \"\" {`, run ^TestPendingMerge: silent (equivalent mutant). commit == line to len(commit) == len(line), run ^TestPendingMerge: silent. Focused run bench test --package ./internal/worktree --run '^TestSourceMergePending' passed with 2 tests. Final git status --short: empty.\n\nExamined: the frozen diff and its spill; spec stories 11-19, the pending-merge owner, Testing decisions, rows MP01-MP18, Edge inventory; ticket 1; the hostile-input checklist; craft-review Coverage axis; finding-discipline.md; checkout.go, checkout_test.go, testhelpers_test.go, worktree_admin.go:193-243, git.go:220-234, land_refusal.go:60-109, admin_readers_test.go, land_journey_test.go:285-319, land_refusal_fixture_test.go:150-169; the repository-wide MERGE_HEAD reference sweep.\nEvidence cursors fetched: the summary, manifest page v1.50ef098f...dd33fc.m.0.0, source s17 (consumers) page 0 whole. Unfetched: manifest page m.0.1 and sources s1-s16 and s18.\n"
          },
          "axis": "Coverage",
          "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
          "tip": "dfacb0be4c0870e191535bb104620168643d4bb1",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5"
          ],
          "supersedes": []
        },
        {
          "id": "mp1-standards-r2",
          "performer": "claude:cmp_mp1_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp1_standards_r2",
            "digest": "sha256:86b7af5200687c3ed985be27f7d1572b87c1f11466bce00d99f79636d3529b90",
            "excerpt": "MP1 Standards axis round 2 (confirming), claude:cmp_mp1_standards_r2, opus/high, one iteration, frozen pair 342975e3..4e6fd43a (chunk tip a329d3ab), evidence sha256:f7212784ac42b5449c862c82e97af8a19cbbe7db9685196db1f9a0d6c68b94cf current=true, clean tree.\n\nS1: confirmed closed. newLinkedCheckout no longer exists under internal/ (rg finds nothing). Each former call site builds its checkout from newRepo plus linkedWorktree (checkout_test.go:54-55, :71-72, :134, :173-174); linkedWorktree stays the only linked-worktree fixture (admin_readers_test.go:34-40). The merges at checkout_test.go:84-87 run on a detached HEAD.\nS2: confirmed closed (evidence only). The expectation stays independently authored at checkout_test.go:48; reviews/commit-merge-parent.md:415-430 keeps the prefix-swap probe on the clean tip f7b5027b: bit, swap, failed 10, restored yes, selection ^TestPendingMergeRefuses with a passing baseline of 12.\n\nNew findings on the repair delta f7b5027b..4e6fd43a: 0. Worst issue: none.\n\nGrade of the repair delta: doc comments at checkout_test.go:67-69, :168-171, and admin_readers_test.go:70-73 are present tense with no narration and no FT tag; the why comment at checkout_test.go:88-90 states the constraint (both octopus orders stay). No new helper; six raw rev-parse spellings now call independentRevParse and three new calls use it; no raw \"rev-parse\" spelling is left. The -c user.email flags match the package idiom (testhelpers_test.go:62, facts_test.go, refs_test.go). The MP15 table at admin_readers_test.go:78-82 reuses three malformed inputs as incidental fixture data. Budgets: checkout_test.go 197, admin_readers_test.go 153; land_journey_test.go and parallel_census_test.go unchanged by the repair. Plan commit 28590065 touches only spec.md (the repair assignment and the Won't handle line at spec.md:397); 2f304925, 4606f4ae, and 4e6fd43a touch only reviews/commit-merge-parent.md; a329d3ab touches only the two _test.go files.\n\nAdvice: the repair report says seven rev-parse conversions, the axis counts six plus three new calls; the doc comment at checkout_test.go:168-171 names \"trailing white space\" and could name the two-final-newlines case.\n\nExamined: the repair delta and its log with stats; checkout_test.go (all); admin_readers_test.go:1-60 (git); worktree admin_readers_test.go:55-94; journey_test.go:155-163; reviews/commit-merge-parent.md:404-447; craft-comments; rg sweeps for newLinkedCheckout, \"rev-parse\", independentRevParse, and the identity flags.\nEvidence cursors fetched: the repair-delta diff spill (197 lines, all), the log spill (16 lines, all), the rg spill (15 lines, all); no preflight evidence page beyond --check-current.\n"
          },
          "axis": "Standards",
          "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
          "tip": "a329d3ab9339c2c88dad79d95f48de811e0e0737",
          "finding_ids": [],
          "supersedes": [
            "mp1-standards-r1"
          ]
        },
        {
          "id": "mp1-spec-r2",
          "performer": "claude:cmp_mp1_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp1_spec_r2",
            "digest": "sha256:2b79fa8c1c4935980e52754cd021fca55bdfb770db9464cd73585a9e2193ade2",
            "excerpt": "MP1 Spec axis round 2 (confirming), claude:cmp_mp1_spec_r2, opus/high, one iteration, frozen pair 342975e3..4e6fd43a (chunk tip a329d3ab), evidence sha256:f7212784ac42b5449c862c82e97af8a19cbbe7db9685196db1f9a0d6c68b94cf current=true, clean tree.\n\n1. Production behavior unchanged, no unrequested behavior: confirmed. The repair delta lists four files: two test files (internal/git/checkout_test.go, internal/worktree/admin_readers_test.go) and two records (reviews/commit-merge-parent.md; specs/commit-merge-parent/spec.md with the claude:cmp_t1_rep1 assignment and one Won't handle line). No production .go file changes. The linked fixture moves to a detached HEAD through linkedWorktree, a fixture detail.\n2. New cases state approved rules with fixture-derived expectations: confirmed. Relative gitdir: step (checkout_test.go:109-118) computes the path with filepath.Rel from --show-toplevel and --git-dir and expects the fixture heads (spec.md:379). Octopus orders (:98, :99-102) expect the argument order of each merge call, so a sorting reader fails one step for any IDs (spec.md:151, MP03). Annotated tag (:178, :185) plants git rev-parse annotated, the unpeeled tag object (spec.md:148-149). Trailing space, CRLF, and \"\\n\\n\" (:186-188) match spec.md:147-148. Every refusal goes through requirePendingMergeRefusal (:45-51) with its own literal prefix (spec.md:152). MP15 (worktree admin_readers_test.go:78-89): subtests empty, main\\n, and a tree ID from gitOutput HEAD^{tree}, each expects false (spec.md:155-156).\n3. New Won't handle line (spec.md:397): confirmed. It follows from the reader rule at spec.md:148-149 and story 12; the refusal authors no one-parent commit, so story 3 and the .pending criterion hold; it contradicts none of stories 11-19, D1, or D2. The Git fact (git merge <annotated tag> writes the tag object ID to MERGE_HEAD) comes from the axis's memory of builtin/merge.c write_merge_heads, not from a run in this pass.\n4. Rows MP01-MP18 cite existing tests in the citation form: confirmed (spec.md:309-326); the six named tests exist at checkout_test.go:70, :133, :172, worktree admin_readers_test.go:62, :74, and land_journey_test.go:278.\n\nFindings: 0. Worst issue: none.\n\nAdvice: spec.md:323 describes MP15 as only \"an empty MERGE_HEAD\" and its why-cell describes the removed read; final reconciliation could widen the behavior cell. spec.md:397 could name the remedy (merge the tag's peeled commit). The doc comment of TestSourceMergePendingIsFalseWhenUndecided cites \"Coverage row GR12\" while MP16 also cites that test; this predates the delta.\n\nExamined: --check-current output; the full repair delta with stat and log; spec.md:1-170 and 300-399; checkout_test.go in full; worktree admin_readers_test.go:59-90; the helpers independentRevParse, linkedWorktree, gitOutput, mustAdminPath; an existence sweep of the six cited tests.\nEvidence cursors fetched: the diff spill (501 lines, all) and the helper spill (49 lines, all). Not fetched: the spec heading-grep spill and any preflight evidence page beyond --check-current.\n"
          },
          "axis": "Spec",
          "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
          "tip": "a329d3ab9339c2c88dad79d95f48de811e0e0737",
          "finding_ids": [],
          "supersedes": [
            "mp1-spec-r1"
          ]
        },
        {
          "id": "mp1-coverage-r2",
          "performer": "claude:cmp_mp1_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "ab34cb06d3a7b1e5a3167d0bde0283b276f72ac9",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp1_coverage_r2",
            "digest": "sha256:25f993b472905014e769d70b784181fa7a270c371c6fbf54add4e894651a8642",
            "excerpt": "MP1 Coverage axis round 2 (confirming), claude:cmp_mp1_coverage_r2, opus/high, one iteration, frozen pair 342975e3..4e6fd43a (chunk tip a329d3ab), evidence sha256:f7212784ac42b5449c862c82e97af8a19cbbe7db9685196db1f9a0d6c68b94cf current=true, dirty=false. The repair changed no production file.\n\nFocused runs on the tip: bench test --package ./internal/git --run '^TestPendingMerge' passed, 23 tests; bench test --package ./internal/worktree --run '^TestSourceMergePending' passed, 5 tests.\n\nC1 confirmed closed: step \"relative gitdir pointer\" at checkout_test.go:109-118 expects heads from independentRevParse of the fixture commits. Probe (swap AdminPath for a reader that joins the raw gitdir: value unresolved, run ^TestPendingMergeReadsEveryHead$): bit, failed 1 (relative_gitdir_pointer), restored yes.\nC2 confirmed closed: step \"octopus in reverse order\" at :99-102 expects [heads[1], heads[0]]. Probe (append then descending insertion sort): bit, failed 3 (octopus, unterminated last line, relative gitdir pointer), restored yes.\nC3 confirmed closed: case \"annotated tag\" at :185 plants the unpeeled tag object ID from independentRevParse(t, linked, \"annotated\"). Probe (swap commit == line for commit != \"\"): bit, failed 2 (annotated_tag, abbreviated_commit_ID), restored yes.\nC4 confirmed closed: cases \"trailing space\", \"carriage return\", \"two final newlines\" at :186-188. Probe (swap strings.Split(content, \"\\n\") for strings.Fields(content)): bit, failed 3 (exactly those cases), restored yes.\nC5 confirmed closed: subtests \"branch name\" and \"tree\" at worktree admin_readers_test.go:78-89, the tree bytes from gitOutput rev-parse HEAD^{tree}. Probe (land_refusal.go:82, swap the return for one that also answers true when the error holds \"holds the line\", run ^TestSourceMergePending): bit, failed 2 (branch_name, tree), restored yes; the empty case and GR12 stayed green.\n\nNew findings on the repair delta: 0. Worst issue: none. Final git status --short: empty.\n\nAdvice: the steps of TestPendingMergeReadsEveryHead depend on each other (checkout_test.go:88), so a run filtered to one later step alone fails; this shape comes from round 1. No real git merge of an annotated tag checks the Won't handle claim at spec.md:397; the test plants the bytes, and the behavior is fail-closed either way.\n\nExamined: the full repair diff (501 lines); checkout.go, checkout_test.go, worktree admin_readers_test.go (whole); land_refusal.go:1-25 and :60-99; git.go:220-234; reviews/commit-merge-parent.md:406-446.\nEvidence cursors fetched: --check-current only (no evidence pages); the diff spill (whole). The five probe spill files were not read; each inline verdict and failure row was complete.\n"
          },
          "axis": "Coverage",
          "base": "342975e31979d4a8ea1708282aa4a13ce7a8ef8a",
          "tip": "a329d3ab9339c2c88dad79d95f48de811e0e0737",
          "finding_ids": [],
          "supersedes": [
            "mp1-coverage-r1"
          ]
        }
      ]
    },
    {
      "id": "MP2",
      "base": "004e4e94019144653de1ec8e179abe015b704242",
      "tip": "eb7c791c83709ec159c6e8dd98592170b8bedc7e",
      "plan_digest": "sha256:26c4daa7f1486d87e7341202b339d5fe7ef02fdf1a4d1fdfbc7b0f8f723304f3",
      "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
      "acceptance_rows": [
        "MP24",
        "MP25",
        "MP32",
        "MP33",
        "MP34",
        "MP35",
        "MP36",
        "MP37",
        "MP38",
        "MP39",
        "MP40",
        "MP41",
        "MP42",
        "MP43",
        "MP62",
        "MP19",
        "MP20",
        "MP21",
        "MP22",
        "MP23",
        "MP26",
        "MP27",
        "MP28",
        "MP29",
        "MP30",
        "MP31",
        "MP44",
        "MP45",
        "MP46",
        "MP47",
        "MP48",
        "MP49",
        "MP50",
        "MP51",
        "MP52",
        "MP53",
        "MP54",
        "MP56",
        "MP55"
      ],
      "verification": [
        {
          "id": "t4-guidance-cmp_t4",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/t4-guidance.txt",
            "digest": "sha256:7d84934d23da9414fadfeaedfe5bdd166c3405cf0fd161de56a1da00eb9198b8",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check docs-currency-workflow\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,9410ea4139b8528a6541f211d2dfe31a45219a94,true\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4946,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench anchors .bench/BENCH-reference.md\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,9410ea4139b8528a6541f211d2dfe31a45219a94,true\nanchors[47]{kind,section,step,needle,line}:\n  require,\"\",0,ship only in the Bench kit repository,25\n  require,\"\",0,a linked repo upgrades with `bench upgrade`,25\n  require,\"\",0,`.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only),98\n  require,\"\",0,Do not append to,43\n  require,\"\",0,The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.,39\n  require,\"\",0,/bench-final-check` writes `capture/retros/<spec-slug>.md`,45\n  require,\"\",0,/bench-drain` owns their reviewed drain,47\n  require,\"\",0,A reviewed spec-backed build keeps its serial ticket commits in one retained integration source,277\n  require,\"\",0,\"no separate conformance phase, driver, or per-check evidence partition\",400\n  require,\"\",0,The spec is optional on the landing and on its resume,289\n  require,\"\",0,constructs through the shared refusal-route registry,296\n  require-in-section,Command Notes,0,\"runs the fast lane on a private checkout of the composed snapshot, and a lane pass\",280\n  forbid,\"\",0,may inherit exact per-check evidence,0\n  require,\"\",0,bench harnesses codex,537\n  require,\"\",0,The destructive-git guard allows an agent push to any branch other than the default branch.,527\n  require,\"\",0,\"A red run prints one `failures[N]{phase,line}` table, and then `gate: red`.\",468\n  require,\"\",0,\"guidance, not rules\",19\n  require,\"\",0,with authority you do not have,27\n  require,\"\",0,Read `CONTEXT.md` for the mental model,13\n  require,\"\",0,Setup and adoption connect a repository to the kit,216\n  require,\"\",0,Context commands expose current state,217\n  require,\"\",0,Oracle commands inspect or enforce readiness,218\n  require,\"\",0,Work commands own isolated execution,219\n  forbid,\"\",0,progressive loading,0\n  require,\"\",0,the plan amendment declares version 2 with a delegate execution block. The amendment splits each chunk verification into one verification for each ticket.,277\n  forbid,\"\",0,\"the rule is guidance, not a hook\",0\n  forbid,\"\",0,\"`bench commit` enforces this boundary: it refuses the primary checkout\",0\n  require,\"\",0,\"Every phase lands this way, and `.bench/BENCH.md` states how `bench commit` enforces that rule.\",291\n  forbid,\"\",0,\"A stale Bench executable is rebuilt, and the landing re-runs under it.\",0\n  require,\"\",0,the reviewer merges the destination into the source worktree with raw Git,310\n  forbid,\"\",0,owns their reviewed drain and its capture commit,0\n  require,\"\",0,`/bench-final-check` states when a tracked or an ignored retro commits.,48\n  forbid,\"\",0,A tracked retro commits with the phase close,0\n  require,\"\",0,The primary checkout owns the `main` section.,361\n  require,\"\",0,The verb keeps a non-empty Next command.,362\n  require,\"\",0,refuses a State that pins a commit outside the tip's ancestry,363\n  require-in-section,Plumbing subcommands,0,`bench gate-prose` is the one internal verb that a session runs directly; Command Notes gives its forms.,383\n  require,\"\",0,\"A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table\",472\n  require,\"\",0,`bench handoff` rewrites only the calling worktree's assignment section.,360\n  forbid,\"\",0,\"`bench handoff [--harness <name>]\",0\n  forbid,\"\",0,`bench worktree reset --to <commit> <target>` plans,0\n  forbid,\"\",0,`bench retro <slug> (--body <markdown>,0\n  require,\"\",0,\"Before the first dispatch, the plan amendment makes the authored version 1 fence a version 2 plan.\",409\n  forbid,\"\",0,`.claude/skills/` carries only the,0\n  forbid,\"\",0,because Claude already has each phase as a command,0\n  forbid,\"\",0,`.claude/skills/` links every `.agents/skills/` skill that has no same-named command.,0\n  require,\"\",0,`.claude/README.md` states which skills `.claude/skills/` links.,73\nhelp[0]{cmd,why}:\n\nNote: every require row has a nonzero line and every forbid row reads 0. The route needle moved from line 307 to line 310 because the three new physical lines sit above it; the needle text is intact.\n"
          },
          "requirement": "t4-guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t4-skills-index-cmp_t4",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/t4-skills-index.txt",
            "digest": "sha256:68eaad6ef1677791e3fa0793edfb52e87a8628c577cb560e5b53f371efaf806f",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check skills-index-command-adapters\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,9410ea4139b8528a6541f211d2dfe31a45219a94,true\ncheck[1]{name,kind,tests_run,subjects}:\n  skills-index-command-adapters,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-skills-index",
          "command": "bench test --check skills-index-command-adapters",
          "exit_code": 0
        },
        {
          "id": "t4-floor-cmp_t4",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/t4-floor.txt",
            "digest": "sha256:61c8c362fb9fa09172621197a00d1bd63eaba093f2265d0b504af84eae3b47da",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,9410ea4139b8528a6541f211d2dfe31a45219a94,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,23732,333\n  github.com/gibbonmi/bench/internal/anchors,pass,2460,179\n  github.com/gibbonmi/bench/internal/conformance,pass,60166,1685\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-floor",
          "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t4-root-conformance-cmp_t4",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/t4-root-conformance.txt",
            "digest": "sha256:e2d0f33981051f90e9c89cbba2d8f33a30ac50729e5d12fc98c2455844b206e0",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,9410ea4139b8528a6541f211d2dfe31a45219a94,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17236,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-landing-focused-cmp_t2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip/t2-landing-focused.txt",
            "digest": "sha256:e61f9da0ea01af2997b8fcd6c7ecd9e5cebfe4bce9c0c698b40e6c5e06d833b8",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,272,8\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/attribution.go --omit 'return nil, NamedPathError{errors.New(\"pending merge changed paths that the commit does not name: \" + strings.Join(omitted, \" \"))}' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$'\nexit 1 (the mutant bit)\nnote: ticket 3 changed the return to two values, so the omitted text gained \"nil, \"; the mutation is still the omission of the D1 refusal.\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/attribution.go,omit,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/conflicted_path_omitted,\"pending_merge_test.go:75: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"conflict.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/deleted_path_omitted,\"pending_merge_test.go:75: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"gone.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/directory_does_not_cover_a_prefix_sibling,\"pending_merge_test.go:75: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"pkgx/a.go\\\\\\\\\\\"\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-landing-focused",
          "command": "bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the NamedPathError refusal that attributedPaths raises for the omitted merged paths in internal/landing/attribution.go, so that an incomplete name set composes. TestLandRefusesAnIncompleteMergeNameSet turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:cmp_t2/tip/t2-landing-focused.txt",
              "digest": "sha256:e61f9da0ea01af2997b8fcd6c7ecd9e5cebfe4bce9c0c698b40e6c5e06d833b8",
              "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,272,8\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/attribution.go --omit 'return nil, NamedPathError{errors.New(\"pending merge changed paths that the commit does not name: \" + strings.Join(omitted, \" \"))}' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$'\nexit 1 (the mutant bit)\nnote: ticket 3 changed the return to two values, so the omitted text gained \"nil, \"; the mutation is still the omission of the D1 refusal.\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/attribution.go,omit,failed,3,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,6\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/conflicted_path_omitted,\"pending_merge_test.go:75: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"conflict.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/deleted_path_omitted,\"pending_merge_test.go:75: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"gone.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/directory_does_not_cover_a_prefix_sibling,\"pending_merge_test.go:75: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"pkgx/a.go\\\\\\\\\\\"\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-commit-focused-cmp_t2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip/t2-commit-focused.txt",
            "digest": "sha256:b440317a21d3b8595256428597832e60f11cead150d6a621baf0615ccddab7ba",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,284,7\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-commit-focused",
          "command": "bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'",
          "exit_code": 0
        },
        {
          "id": "t2-floor-cmp_t2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip/t2-floor.txt",
            "digest": "sha256:225c790de56cd1171a5994255cbe7f63f09142e824e5cdeaa8476d63060db556",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,9680,117\n  github.com/gibbonmi/bench/internal/git,pass,1343,171\n  github.com/gibbonmi/bench/internal/landing,pass,8045,218\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,17\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-floor",
          "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance-cmp_t2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip/t2-root-conformance.txt",
            "digest": "sha256:6b98953f8718046e4752631311d1735a23f077a917fb9e3eca0db0e17d1cdb92",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15895,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-git-plumbing-owner-cmp_t2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip/t2-git-plumbing-owner.txt",
            "digest": "sha256:f710b6641139e7bed4fbe9dc394af0c2380bcd575ca42ada126cdde49088f882",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check git-plumbing-owner\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,98,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t2-ordinary-build-census-cmp_t2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip/t2-ordinary-build-census.txt",
            "digest": "sha256:4a20c1666589e693e309a465c14204103f88beb7ef8aaabe28f27f846ccc0f0b",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check ordinary-build-census\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,c58f0c15d78288b1c53ed0a10bba4f963f275cf2,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,317,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t3-landing-focused-cmp_t3",
          "performer": "claude:cmp_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3/tip/t3-landing-focused.txt",
            "digest": "sha256:8b7196faeac70db0ceea660012d036ec8f3421d556728c5e93b73a8301ebc668",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,582,27\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/landing.go --swap 'append([]string{r.Expected}, snapshot.mergeHeads...)' --with '[]string{r.Expected}' --package ./internal/landing --run '^TestLandRecordsEachMergeHeadAsAParent$'\nexit 1\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRecordsEachMergeHeadAsAParent$,passed,3\n  github.com/gibbonmi/bench/internal/landing,fail,124,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/octopus_merge,\"pending_merge_test.go:105: commit and parents = \\\"cedc46390bb571a71e6fc55b800c4b8ffdd73449 95408da9f3789e0b7ae3b11400181b0188ec96c4\\\", want \\\"cedc46390bb571a71e6fc55b800c4b8ffdd73449 95408da9f3789e0b7ae3b11400181b0188ec96c4 458b76d472d83da3d4a… (303 bytes)\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/one_merge_head,\"pending_merge_test.go:105: commit and parents = \\\"94020d0e16dd12801a102c9bbfac8a2eb78dda20 95408da9f3789e0b7ae3b11400181b0188ec96c4\\\", want \\\"94020d0e16dd12801a102c9bbfac8a2eb78dda20 95408da9f3789e0b7ae3b11400181b0188ec96c4 559add3dfab5e89c778… (262 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-landing-focused",
          "command": "bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the merge-head -p arguments from the commit-tree call of Owner.Land in internal/landing/landing.go, so that the publication records one parent during a pending merge. TestLandRecordsEachMergeHeadAsAParent turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:cmp_t3/tip/t3-landing-focused.txt",
              "digest": "sha256:8b7196faeac70db0ceea660012d036ec8f3421d556728c5e93b73a8301ebc668",
              "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,582,27\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/landing.go --swap 'append([]string{r.Expected}, snapshot.mergeHeads...)' --with '[]string{r.Expected}' --package ./internal/landing --run '^TestLandRecordsEachMergeHeadAsAParent$'\nexit 1\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRecordsEachMergeHeadAsAParent$,passed,3\n  github.com/gibbonmi/bench/internal/landing,fail,124,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/octopus_merge,\"pending_merge_test.go:105: commit and parents = \\\"cedc46390bb571a71e6fc55b800c4b8ffdd73449 95408da9f3789e0b7ae3b11400181b0188ec96c4\\\", want \\\"cedc46390bb571a71e6fc55b800c4b8ffdd73449 95408da9f3789e0b7ae3b11400181b0188ec96c4 458b76d472d83da3d4a… (303 bytes)\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/one_merge_head,\"pending_merge_test.go:105: commit and parents = \\\"94020d0e16dd12801a102c9bbfac8a2eb78dda20 95408da9f3789e0b7ae3b11400181b0188ec96c4\\\", want \\\"94020d0e16dd12801a102c9bbfac8a2eb78dda20 95408da9f3789e0b7ae3b11400181b0188ec96c4 559add3dfab5e89c778… (262 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-commit-focused-cmp_t3",
          "performer": "claude:cmp_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3/tip/t3-commit-focused.txt",
            "digest": "sha256:5233021b53ac9630e5b605c3f1efe01919340c33b4aa202ccc1a938c458959ea",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,233,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-commit-focused",
          "command": "bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'",
          "exit_code": 0
        },
        {
          "id": "t3-floor-cmp_t3",
          "performer": "claude:cmp_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3/tip/t3-floor.txt",
            "digest": "sha256:486568439002c4cc43b3d7f33aa656a007ae2d7a15f4433834c25b6231f3b646",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/git --package ./internal/landing --package ./internal/commit\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,9573,117\n  github.com/gibbonmi/bench/internal/git,pass,1310,171\n  github.com/gibbonmi/bench/internal/landing,pass,8043,218\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-floor",
          "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "t3-root-conformance-cmp_t3",
          "performer": "claude:cmp_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3/tip/t3-root-conformance.txt",
            "digest": "sha256:007b2b1e7ff3a43ab1a16a0da8778575d24ff99d614c229b773d567aafa6f97a",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15400,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-git-plumbing-owner-cmp_t3",
          "performer": "claude:cmp_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3/tip/t3-git-plumbing-owner.txt",
            "digest": "sha256:4bd7b0179825eb01ce318cccf79e8a02f37f1e67d28c6cddd7b98c58c08fd627",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check git-plumbing-owner\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,97,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t3-ordinary-build-census-cmp_t3",
          "performer": "claude:cmp_t3",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3/tip/t3-ordinary-build-census.txt",
            "digest": "sha256:51e86f07ca98f62bd4d268aaae757990575f783d0a043c16bccff0b641adfb86",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check ordinary-build-census\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc5312c81de0dcf8d42e6384a71ef69346fc9a1d,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,318,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
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
  },
  "amendments": [
    {
      "from": "sha256:3b3b51e2c87d4449eb3044f3d9d1fb6c5eb414dd8b48e9d711f48822bfe6eda6",
      "to": "sha256:06b9db25ba9dba1eee9e58f5e427693c6ff3473a78c4ce1e7ad3e6aa8c2f09dd",
      "chunk_ids": {
        "MP1": [
          "MP1"
        ]
      }
    },
    {
      "from": "sha256:06b9db25ba9dba1eee9e58f5e427693c6ff3473a78c4ce1e7ad3e6aa8c2f09dd",
      "to": "sha256:ab73d3743fe659a20eb212cb44af6b2b0cdee6be76c7a5a1bcb2f22d924d2c7c",
      "chunk_ids": {
        "MP1": [
          "MP1"
        ]
      }
    },
    {
      "from": "sha256:ab73d3743fe659a20eb212cb44af6b2b0cdee6be76c7a5a1bcb2f22d924d2c7c",
      "to": "sha256:3790bad197ff15582b28a99980728308b34727636ec7bba008c0a26ccfdaa72d",
      "chunk_ids": {
        "MP1": [
          "MP1"
        ]
      }
    },
    {
      "from": "sha256:3790bad197ff15582b28a99980728308b34727636ec7bba008c0a26ccfdaa72d",
      "to": "sha256:26c4daa7f1486d87e7341202b339d5fe7ef02fdf1a4d1fdfbc7b0f8f723304f3",
      "chunk_ids": {
        "MP1": [
          "MP1"
        ]
      }
    }
  ]
}
```

## Repair state

Chunk MP1 round 1 ran three opus/high axes in parallel at `f7b5027b` over the frozen pair `342975e3..f7b5027b`. Only the Coverage axis ran probes, and the tree was clean after every return. MP1 uses repair cycle 1 of 2 through one fresh repair session for ticket 1, which takes S1, C1, C2, C3, C4, and C5. S2 is a record gap with no code change, so the orchestrator closes it below as evidence-only work. Under `--auto-approve`, the reviewer approved every repair round and fence expansion in advance.

The ticket 1 repair session `claude:cmp_t1_rep1` closed S1 and C1 to C5 in `a329d3ab` with tests only. It refroze MP1 at that tip in `4e6fd43a`. MP1 round 2 ran three opus/high axes in parallel over the repair delta `f7b5027b..4e6fd43a`. Every axis confirmed each fold, the Coverage axis bit each fold with an independent probe, and no axis found a new finding. MP1 closes with repair cycle 1 of 2 used and no open finding.

## Standards

Findings: 2. Worst issue: S2.

- S1 (`auto-fix`, confidence 5): `newLinkedCheckout` in `internal/git/checkout_test.go` repeats the linked-worktree fixture that `linkedWorktree` in `internal/git/admin_readers_test.go` owns. Rule: AGENTS.md, one source per fact. Repair: compose `newRepo` and `linkedWorktree`.
- S2 (`auto-fix`, confidence 5): the independent `read pending merge:` expectation in `requirePendingMergeRefusal` had no recorded red. Rule: the AGENTS.md exception for an independently authored test expectation. Closed evidence-only: the orchestrator ran the prefix-swap probe on the clean tip `f7b5027b`, and it bit with 10 failed tests and a restored subject:

```text
tree[1]{target,head,dirty}:
  commit-merge-parent-build,f7b5027ba5a82de74869ddd3cc343bb53b0e0550,false
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/git/checkout.go,swap,failed,10,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/git,^TestPendingMergeRefuses,passed,12
spilled{lines=20,bytes=3881,omitted_lines=11,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/f80fcedec92401bd6cf2acfce17ca107/1791693993797782733-e4d5a241f34fdcf8.out}
  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/directory,"checkout_test.go:154: PendingMerge = [], \"/dev/shm/H5YTDC/t/TestPendingMergeRefusesUnreadableState3925670371/001/.git/worktrees/linked/MERGE_HEAD\" is not a regular file; want no merge head and a read pending merge refusal",1
  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/fifo,"checkout_test.go:154: PendingMerge = [], \"/dev/shm/H5YTDC/t/TestPendingMergeRefusesUnreadableState3925670371/001/.git/worktrees/linked/MERGE_HEAD\" is not a regular file; want no merge head and a read pending merge refusal",1
  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/live_symbolic_link,"checkout_test.go:154: PendingMerge = [], \"/dev/shm/H5YTDC/t/TestPendingMergeRefusesUnreadableState3925670371/001/.git/worktrees/linked/MERGE_HEAD\" is not a regular file; want no merge head and a read pending merge refusal",1
  github.com/gibbonmi/bench/internal/git,TestPendingMergeRefusesUnreadableState/outside_every_repository,"checkout_test.go:154: PendingMerge = [], checkout administration path resolution failed (git -C /dev/shm/H5YTDC/t/TestPendingMergeRefusesUnreadableState3925670371/004 rev-parse --git-path MERGE_HEAD: exit status 128: fatal: not a git reposi… (283 bytes)",2
skips[0]{package,test,reason}:
```

## Spec

Findings: 0. Worst issue: none. Every row MP01 to MP18 is met, and every MP01 to MP15 seam cell cites an existing test.

Advice: the owner section cites `land_refusal.go:82-91` and `:78-81` in the present tense, and the function now sits at `:76-83`; final reconciliation decides a currency touch.

## Coverage

Findings: 5. Worst issue: C2.

- C1 (`auto-fix`, confidence 8): no test runs the reader through a relative `gitdir:` pointer, which `spec.md` claims MP01 covers. Repair: a `TestPendingMergeReadsEveryHead` step that rewrites `linked/.git` to a relative `gitdir:` and expects the merge head.
- C2 (`auto-fix`, confidence 7): the octopus order pin depends on random SHAs, and an ascending-sort mutant stayed silent. Repair: an MP03 fixture whose file order differs from both sorted orders.
- C3 (`auto-fix`, confidence 7): no row plants an annotated tag object ID, and a `len(commit) == len(line)` mutant stayed silent. Repair: an annotated-tag case in `TestPendingMergeRefusesMalformedState`.
- C4 (`auto-fix`, confidence 6): CRLF, trailing space, and a double final newline refuse but have no pin; a `TrimSpace` mutant stayed silent. Repair: three malformed cases.
- C5 (`auto-fix`, confidence 5): a file-presence copy of the reader in `land_refusal.go` would pass MP15 to MP18. Repair: an MP15 subtest with a non-empty malformed MERGE_HEAD that expects false.
