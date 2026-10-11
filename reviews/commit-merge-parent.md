# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commit-merge-parent/spec.md",
  "plan_digest": "sha256:06b9db25ba9dba1eee9e58f5e427693c6ff3473a78c4ce1e7ad3e6aa8c2f09dd",
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
      "from": "sha256:3b3b51e2c87d4449eb3044f3d9d1fb6c5eb414dd8b48e9d711f48822bfe6eda6",
      "to": "sha256:06b9db25ba9dba1eee9e58f5e427693c6ff3473a78c4ce1e7ad3e6aa8c2f09dd",
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
