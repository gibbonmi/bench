# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/commit-merge-parent/spec.md",
  "plan_digest": "sha256:7654187ea92c112db47b0160f0dd2cc0ea2e1645dc5648227d777f5de457dd7f",
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
      "tip": "cca1013b089a78e230bb6f15d9f3266ce0ab0ffa",
      "plan_digest": "sha256:7654187ea92c112db47b0160f0dd2cc0ea2e1645dc5648227d777f5de457dd7f",
      "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
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
        },
        {
          "id": "t3-landing-focused-cmp_t3_rep1",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/t3-landing-focused.txt",
            "digest": "sha256:69b288c285abc8e5a2c5a30c1663f40cda51b4609a3794137ae3cc7d8cb92922",
            "excerpt": "$ bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,627,28\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n\n$ bench probe internal/landing/landing.go --swap 'append([]string{r.Expected}, snapshot.mergeHeads...)' --with '[]string{r.Expected}' --package ./internal/landing --run '^TestLandRecordsEachMergeHeadAsAParent$'\nexit_code=1\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRecordsEachMergeHeadAsAParent$,passed,3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,123,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/octopus_merge,\"pending_merge_test.go:113: commit and parents = \\\"8b85c33f3ac25d8cb56c91a3615c2ccefaf83eba e7bf41ab64d9c45b59043cd857d35560bb3496d3\\\", want \\\"8b85c33f3ac25d8cb56c91a3615c2ccefaf83eba e7bf41ab64d9c45b59043cd857d35560bb3496d3 b5f323234c3fe9780e2… (303 bytes)\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/one_merge_head,\"pending_merge_test.go:113: commit and parents = \\\"657256f7881fd7eb2767f8bbd13ca4d44f25e6fb e7bf41ab64d9c45b59043cd857d35560bb3496d3\\\", want \\\"657256f7881fd7eb2767f8bbd13ca4d44f25e6fb e7bf41ab64d9c45b59043cd857d35560bb3496d3 1f21e0cce56cd255e9b… (262 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:cmp_t3_rep1/t3-landing-focused.txt",
              "digest": "sha256:69b288c285abc8e5a2c5a30c1663f40cda51b4609a3794137ae3cc7d8cb92922",
              "excerpt": "$ bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,627,28\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n\n$ bench probe internal/landing/landing.go --swap 'append([]string{r.Expected}, snapshot.mergeHeads...)' --with '[]string{r.Expected}' --package ./internal/landing --run '^TestLandRecordsEachMergeHeadAsAParent$'\nexit_code=1\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRecordsEachMergeHeadAsAParent$,passed,3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,123,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/octopus_merge,\"pending_merge_test.go:113: commit and parents = \\\"8b85c33f3ac25d8cb56c91a3615c2ccefaf83eba e7bf41ab64d9c45b59043cd857d35560bb3496d3\\\", want \\\"8b85c33f3ac25d8cb56c91a3615c2ccefaf83eba e7bf41ab64d9c45b59043cd857d35560bb3496d3 b5f323234c3fe9780e2… (303 bytes)\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/one_merge_head,\"pending_merge_test.go:113: commit and parents = \\\"657256f7881fd7eb2767f8bbd13ca4d44f25e6fb e7bf41ab64d9c45b59043cd857d35560bb3496d3\\\", want \\\"657256f7881fd7eb2767f8bbd13ca4d44f25e6fb e7bf41ab64d9c45b59043cd857d35560bb3496d3 1f21e0cce56cd255e9b… (262 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-commit-focused-cmp_t3_rep1",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/t3-commit-focused.txt",
            "digest": "sha256:9b8d9fb2cd0c350c8c3908596f06c2af6103b452c991974e0023722441a6b26d",
            "excerpt": "$ bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,224,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-commit-focused",
          "command": "bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'",
          "exit_code": 0
        },
        {
          "id": "t3-floor-cmp_t3_rep1",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/t3-floor.txt",
            "digest": "sha256:5eaf7eeb599bf29000b9a0c82547e1275d18379a1e755afceef526ef96c4af6b",
            "excerpt": "$ bench test --package ./internal/git --package ./internal/landing --package ./internal/commit\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,9767,117\n  github.com/gibbonmi/bench/internal/git,pass,1311,171\n  github.com/gibbonmi/bench/internal/landing,pass,8218,220\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-floor",
          "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "t3-root-conformance-cmp_t3_rep1",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/t3-root-conformance.txt",
            "digest": "sha256:112ccba215724ede4a68753a9b350703416469bc6976220731a9e9b6083b139c",
            "excerpt": "$ bench test --check root-conformance\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15671,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-git-plumbing-owner-cmp_t3_rep1",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/t3-git-plumbing-owner.txt",
            "digest": "sha256:84f92cf01355486a1ac29227818e2308358e3d76e8cc0b3d5db0d996e66cec10",
            "excerpt": "$ bench test --check git-plumbing-owner\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,96,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t3-ordinary-build-census-cmp_t3_rep1",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/t3-ordinary-build-census.txt",
            "digest": "sha256:50aaf5dd5a580b3093585a881d98e60ee57d70f62b5761960891f9da5d4a37ad",
            "excerpt": "$ bench test --check ordinary-build-census\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,83d43fbaf4cbf79f571cafdfeb054cc5ad434120,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,318,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t2-landing-focused-cmp_t2-r2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip2/t2-landing-focused.txt",
            "digest": "sha256:cd0d76454af9c3d951d23ebf7fdd4f5d05fd4248130c332af3e64cf1b8173c44",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,335,9\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/attribution.go --omit 'return nil, NamedPathError{errors.New(\"pending merge changed paths that the commit does not name: \" + strings.Join(omitted, \" \"))}' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$'\nexit 1 (the mutant bit)\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/attribution.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,7\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/conflicted_path_omitted,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"conflict.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/deleted_path_omitted,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"gone.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/directory_does_not_cover_a_prefix_sibling,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"pkgx/a.go\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/octopus_side_path_omitted,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"side-b\\\\\\\\\\\"\\\"\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:cmp_t2/tip2/t2-landing-focused.txt",
              "digest": "sha256:cd0d76454af9c3d951d23ebf7fdd4f5d05fd4248130c332af3e64cf1b8173c44",
              "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,335,9\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/attribution.go --omit 'return nil, NamedPathError{errors.New(\"pending merge changed paths that the commit does not name: \" + strings.Join(omitted, \" \"))}' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$'\nexit 1 (the mutant bit)\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/attribution.go,omit,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,7\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/conflicted_path_omitted,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"conflict.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/deleted_path_omitted,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"gone.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/directory_does_not_cover_a_prefix_sibling,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"pkgx/a.go\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/octopus_side_path_omitted,\"pending_merge_test.go:83: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"side-b\\\\\\\\\\\"\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-commit-focused-cmp_t2-r2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip2/t2-commit-focused.txt",
            "digest": "sha256:9445587a6aeeadb8f60eb8b3fca99c5e018c0b40ebb0210b7605e4ba362f48bc",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,284,7\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-commit-focused",
          "command": "bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'",
          "exit_code": 0
        },
        {
          "id": "t2-floor-cmp_t2-r2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip2/t2-floor.txt",
            "digest": "sha256:cedbb577665f6d1d27224768c6d50dfa4c0ff08950d7ebc620092fd50ab65682",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,10065,117\n  github.com/gibbonmi/bench/internal/git,pass,1353,171\n  github.com/gibbonmi/bench/internal/landing,pass,8434,220\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,17\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-floor",
          "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance-cmp_t2-r2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip2/t2-root-conformance.txt",
            "digest": "sha256:fd09ef17d0f4cf6b76b4cfb3da229ab141dd1aaafea90663b1ceee1396c0bd72",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15223,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-git-plumbing-owner-cmp_t2-r2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip2/t2-git-plumbing-owner.txt",
            "digest": "sha256:b4951e9056477caec3514134afd93b01d1cea27d787025195f373f37c0886272",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check git-plumbing-owner\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,95,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t2-ordinary-build-census-cmp_t2-r2",
          "performer": "claude:cmp_t2",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2/tip2/t2-ordinary-build-census.txt",
            "digest": "sha256:4b091d7bef2191abaf4b46b78d4e4f91670a6e8a1ad83a3e6858f6558e61c4bd",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check ordinary-build-census\nexit 0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,bc50e71e7019b1c84766fd059a5e4d85ee047276,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,317,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t4-guidance-cmp_t4-r2",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip2/t4-guidance.txt",
            "digest": "sha256:99b3cbfb205819a04f3c101e255d1d3f4a559f92bacb46aeaee8163a08c02ec7",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check docs-currency-workflow\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,27aaedbafb32747ea3f06d06fa87a7f8986e033d,false\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5013,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench anchors .bench/BENCH-reference.md\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,27aaedbafb32747ea3f06d06fa87a7f8986e033d,false\nanchors[47]{kind,section,step,needle,line}:\n  require,\"\",0,ship only in the Bench kit repository,25\n  require,\"\",0,a linked repo upgrades with `bench upgrade`,25\n  require,\"\",0,`.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only),98\n  require,\"\",0,Do not append to,43\n  require,\"\",0,The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.,39\n  require,\"\",0,/bench-final-check` writes `capture/retros/<spec-slug>.md`,45\n  require,\"\",0,/bench-drain` owns their reviewed drain,47\n  require,\"\",0,A reviewed spec-backed build keeps its serial ticket commits in one retained integration source,277\n  require,\"\",0,\"no separate conformance phase, driver, or per-check evidence partition\",400\n  require,\"\",0,The spec is optional on the landing and on its resume,289\n  require,\"\",0,constructs through the shared refusal-route registry,296\n  require-in-section,Command Notes,0,\"runs the fast lane on a private checkout of the composed snapshot, and a lane pass\",280\n  forbid,\"\",0,may inherit exact per-check evidence,0\n  require,\"\",0,bench harnesses codex,537\n  require,\"\",0,The destructive-git guard allows an agent push to any branch other than the default branch.,527\n  require,\"\",0,\"A red run prints one `failures[N]{phase,line}` table, and then `gate: red`.\",468\n  require,\"\",0,\"guidance, not rules\",19\n  require,\"\",0,with authority you do not have,27\n  require,\"\",0,Read `CONTEXT.md` for the mental model,13\n  require,\"\",0,Setup and adoption connect a repository to the kit,216\n  require,\"\",0,Context commands expose current state,217\n  require,\"\",0,Oracle commands inspect or enforce readiness,218\n  require,\"\",0,Work commands own isolated execution,219\n  forbid,\"\",0,progressive loading,0\n  require,\"\",0,the plan amendment declares version 2 with a delegate execution block. The amendment splits each chunk verification into one verification for each ticket.,277\n  forbid,\"\",0,\"the rule is guidance, not a hook\",0\n  forbid,\"\",0,\"`bench commit` enforces this boundary: it refuses the primary checkout\",0\n  require,\"\",0,\"Every phase lands this way, and `.bench/BENCH.md` states how `bench commit` enforces that rule.\",291\n  forbid,\"\",0,\"A stale Bench executable is rebuilt, and the landing re-runs under it.\",0\n  require,\"\",0,the reviewer merges the destination into the source worktree with raw Git,310\n  forbid,\"\",0,owns their reviewed drain and its capture commit,0\n  require,\"\",0,`/bench-final-check` states when a tracked or an ignored retro commits.,48\n  forbid,\"\",0,A tracked retro commits with the phase close,0\n  require,\"\",0,The primary checkout owns the `main` section.,361\n  require,\"\",0,The verb keeps a non-empty Next command.,362\n  require,\"\",0,refuses a State that pins a commit outside the tip's ancestry,363\n  require-in-section,Plumbing subcommands,0,`bench gate-prose` is the one internal verb that a session runs directly; Command Notes gives its forms.,383\n  require,\"\",0,\"A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table\",472\n  require,\"\",0,`bench handoff` rewrites only the calling worktree's assignment section.,360\n  forbid,\"\",0,\"`bench handoff [--harness <name>]\",0\n  forbid,\"\",0,`bench worktree reset --to <commit> <target>` plans,0\n  forbid,\"\",0,`bench retro <slug> (--body <markdown>,0\n  require,\"\",0,\"Before the first dispatch, the plan amendment makes the authored version 1 fence a version 2 plan.\",409\n  forbid,\"\",0,`.claude/skills/` carries only the,0\n  forbid,\"\",0,because Claude already has each phase as a command,0\n  forbid,\"\",0,`.claude/skills/` links every `.agents/skills/` skill that has no same-named command.,0\n  require,\"\",0,`.claude/README.md` states which skills `.claude/skills/` links.,73\nhelp[0]{cmd,why}:\n\nNote: every require row has a nonzero line and every forbid row reads 0. The route needle sits at line 310 because the three new physical lines of ticket 4 sit above it; the needle text is intact.\n"
          },
          "requirement": "t4-guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t4-skills-index-cmp_t4-r2",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip2/t4-skills-index.txt",
            "digest": "sha256:9e3d7a7e1e4c60aee02c81ed9c90aa3204a6fef043b62e02cfa2b9b66f785f4f",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check skills-index-command-adapters\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,27aaedbafb32747ea3f06d06fa87a7f8986e033d,false\ncheck[1]{name,kind,tests_run,subjects}:\n  skills-index-command-adapters,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-skills-index",
          "command": "bench test --check skills-index-command-adapters",
          "exit_code": 0
        },
        {
          "id": "t4-floor-cmp_t4-r2",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip2/t4-floor.txt",
            "digest": "sha256:cc95b4722e6600852771ebe51d3254a484bc002fffbb4bd6677dc0cf071c6ffb",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,27aaedbafb32747ea3f06d06fa87a7f8986e033d,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16568,333\n  github.com/gibbonmi/bench/internal/anchors,pass,2297,179\n  github.com/gibbonmi/bench/internal/conformance,pass,56147,1685\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-floor",
          "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t4-root-conformance-cmp_t4-r2",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip2/t4-root-conformance.txt",
            "digest": "sha256:7f484c3d191929ce750a7c44b0d731e0f52ee72ae13a0c779cf7a5ece81aa0b2",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,27aaedbafb32747ea3f06d06fa87a7f8986e033d,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15184,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-landing-focused-cmp_t2_rep1",
          "performer": "claude:cmp_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2_rep1/t2-landing-focused.txt",
            "digest": "sha256:e3a5ea94da4214855b753b5291be5a661e3ce029a50fe2e2379032b141489fb0",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,372,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/attribution.go --omit 'return nil, NamedPathError{errors.New(\"pending merge changed paths that the commit does not name: \" + strings.Join(omitted, \" \"))}' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/attribution.go,omit,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1560\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,377,8\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/conflicted_path_omitted,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"conflict.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/deleted_path_omitted,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"gone.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/directory_does_not_cover_a_prefix_sibling,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"pkgx/a.go\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/glob-shaped_name_covers_only_itself,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"side-a\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/octopus_side_path_omitted,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"side-b\\\\\\\\\\\"\\\"\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:cmp_t2_rep1/t2-landing-focused.txt",
              "digest": "sha256:e3a5ea94da4214855b753b5291be5a661e3ce029a50fe2e2379032b141489fb0",
              "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet|TestLandRefusesAnUnreadableMergeState|TestLandComposesOnlyNamedPathsAndCASesExpectedBase'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,372,10\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench probe internal/landing/attribution.go --omit 'return nil, NamedPathError{errors.New(\"pending merge changed paths that the commit does not name: \" + strings.Join(omitted, \" \"))}' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/attribution.go,omit,failed,5,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,8\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,1560\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,377,8\nfailures[5]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/conflicted_path_omitted,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"conflict.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/deleted_path_omitted,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"gone.txt\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/directory_does_not_cover_a_prefix_sibling,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"pkgx/a.go\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/glob-shaped_name_covers_only_itself,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"side-a\\\\\\\\\\\"\\\"\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/octopus_side_path_omitted,\"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \\\"pending merge changed paths that the commit does not name: \\\\\\\\\\\"side-b\\\\\\\\\\\"\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t2-commit-focused-cmp_t2_rep1",
          "performer": "claude:cmp_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2_rep1/t2-commit-focused.txt",
            "digest": "sha256:03506c2d3f983b2b436b532f2e98f6ca2495a1b142f6c822a864a118ee2e222b",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,281,7\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-commit-focused",
          "command": "bench test --package ./internal/commit --run 'TestCommitRefusesAnIncompleteMergeNameSet|TestCommitRefusesAnUnreadableMergeState|TestDryRunRefusesAnIncompleteMergeNameSet|TestMergeNameSetRefusalPrecedesCommitment'",
          "exit_code": 0
        },
        {
          "id": "t2-floor-cmp_t2_rep1",
          "performer": "claude:cmp_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2_rep1/t2-floor.txt",
            "digest": "sha256:a3cd890ab6f39776c2932a5c47671dac00d5a0acae62dd6e468ca122f312d24a",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\npackages[4]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,9463,117\n  github.com/gibbonmi/bench/internal/git,pass,1338,171\n  github.com/gibbonmi/bench/internal/landing,pass,8111,221\n  github.com/gibbonmi/bench/internal/refusalroute,pass,4,17\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t2-floor",
          "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit --package ./internal/refusalroute",
          "exit_code": 0
        },
        {
          "id": "t2-root-conformance-cmp_t2_rep1",
          "performer": "claude:cmp_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2_rep1/t2-root-conformance.txt",
            "digest": "sha256:e6ee2c7d06a62e30e2c1ac36b14e12759d5478569437f67eb9d00c03d7fd9375",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15551,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t2-git-plumbing-owner-cmp_t2_rep1",
          "performer": "claude:cmp_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2_rep1/t2-git-plumbing-owner.txt",
            "digest": "sha256:832d7c74a14cca512cd48398f2dce8cd517cbb6ecac198fd5bba55f56dc67c03",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check git-plumbing-owner\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,97,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t2-ordinary-build-census-cmp_t2_rep1",
          "performer": "claude:cmp_t2_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t2_rep1/t2-ordinary-build-census.txt",
            "digest": "sha256:c9a952114ab6dc140c87de88e2b871b4e18ec10014221de92b810378194970e2",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check ordinary-build-census\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,cca1013b089a78e230bb6f15d9f3266ce0ab0ffa,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,319,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t2-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t3-landing-focused-cmp_t3_rep1-r2",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/tip2/t3-landing-focused.txt",
            "digest": "sha256:4ec1715fc2a473371f9a324c4a50cd92407460435300eb95343ac662430720fd",
            "excerpt": "$ bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,635,28\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n\n$ bench probe internal/landing/landing.go --swap 'append([]string{r.Expected}, snapshot.mergeHeads...)' --with '[]string{r.Expected}' --package ./internal/landing --run '^TestLandRecordsEachMergeHeadAsAParent$'\nexit_code=1\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRecordsEachMergeHeadAsAParent$,passed,3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,125,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/octopus_merge,\"pending_merge_test.go:115: commit and parents = \\\"8d575701b7d8dfe54a5d6363754d955a0baaf2c4 6d918bf832e43746291f6f7a550ca6c1fbb3ad64\\\", want \\\"8d575701b7d8dfe54a5d6363754d955a0baaf2c4 6d918bf832e43746291f6f7a550ca6c1fbb3ad64 223fbcf0928208fa507… (303 bytes)\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/one_merge_head,\"pending_merge_test.go:115: commit and parents = \\\"c3d7f99af215352a27da4e91754e0b9c69b5ae6e 6d918bf832e43746291f6f7a550ca6c1fbb3ad64\\\", want \\\"c3d7f99af215352a27da4e91754e0b9c69b5ae6e 6d918bf832e43746291f6f7a550ca6c1fbb3ad64 24be87c4f63efe7e29b… (262 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
              "ref": "claude-agent:cmp_t3_rep1/tip2/t3-landing-focused.txt",
              "digest": "sha256:4ec1715fc2a473371f9a324c4a50cd92407460435300eb95343ac662430720fd",
              "excerpt": "$ bench test --package ./internal/landing --run 'TestLandRecordsEachMergeHeadAsAParent|TestComposeAfterMergeCommitReadsTheMergedBase|TestLandRefusesAMergeStateChangeDuringAuthorization|TestLandPublishesAMergeWithAnUnchangedTree|TestLandEndsThePendingMerge|TestLandReportsPublishedMergeWhenReconciliationFails|TestLandDetachedHEADHasOneExpectedParent|TestLandPreAuthorizationRefusalTable'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,pass,635,28\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n\n$ bench probe internal/landing/landing.go --swap 'append([]string{r.Expected}, snapshot.mergeHeads...)' --with '[]string{r.Expected}' --package ./internal/landing --run '^TestLandRecordsEachMergeHeadAsAParent$'\nexit_code=1\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/landing/landing.go,swap,failed,2,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/landing,^TestLandRecordsEachMergeHeadAsAParent$,passed,3\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/landing,fail,125,3\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/octopus_merge,\"pending_merge_test.go:115: commit and parents = \\\"8d575701b7d8dfe54a5d6363754d955a0baaf2c4 6d918bf832e43746291f6f7a550ca6c1fbb3ad64\\\", want \\\"8d575701b7d8dfe54a5d6363754d955a0baaf2c4 6d918bf832e43746291f6f7a550ca6c1fbb3ad64 223fbcf0928208fa507… (303 bytes)\",1\n  github.com/gibbonmi/bench/internal/landing,TestLandRecordsEachMergeHeadAsAParent/one_merge_head,\"pending_merge_test.go:115: commit and parents = \\\"c3d7f99af215352a27da4e91754e0b9c69b5ae6e 6d918bf832e43746291f6f7a550ca6c1fbb3ad64\\\", want \\\"c3d7f99af215352a27da4e91754e0b9c69b5ae6e 6d918bf832e43746291f6f7a550ca6c1fbb3ad64 24be87c4f63efe7e29b… (262 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t3-commit-focused-cmp_t3_rep1-r2",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/tip2/t3-commit-focused.txt",
            "digest": "sha256:4c6588229b67e450ab1ce842a02a317c9b84de3812618add2f41f07937be663b",
            "excerpt": "$ bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,234,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-commit-focused",
          "command": "bench test --package ./internal/commit --run 'TestCommitRecordsTheMergeParent'",
          "exit_code": 0
        },
        {
          "id": "t3-floor-cmp_t3_rep1-r2",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/tip2/t3-floor.txt",
            "digest": "sha256:365258b156074394eaef1ac66c606da1996ca9651af79e126db0359b63e7b3f3",
            "excerpt": "$ bench test --package ./internal/git --package ./internal/landing --package ./internal/commit\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commit,pass,9715,117\n  github.com/gibbonmi/bench/internal/git,pass,1353,171\n  github.com/gibbonmi/bench/internal/landing,pass,8336,221\nfailures[0]{package,test,line,lines}:\nskips[2]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/descendant-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/landing,TestLandPreAuthorizationRefusalTable/direct-device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t3-floor",
          "command": "bench test --package ./internal/git --package ./internal/landing --package ./internal/commit",
          "exit_code": 0
        },
        {
          "id": "t3-root-conformance-cmp_t3_rep1-r2",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/tip2/t3-root-conformance.txt",
            "digest": "sha256:4167816e73f537d63ec8bf7fd2d5605b06fbbee2d175843d53ded079dbdecba1",
            "excerpt": "$ bench test --check root-conformance\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15151,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t3-git-plumbing-owner-cmp_t3_rep1-r2",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/tip2/t3-git-plumbing-owner.txt",
            "digest": "sha256:5f5ad13ef3e187df20e91d57212579acb35347707c9b3b28f7e85bfef0c17f1b",
            "excerpt": "$ bench test --check git-plumbing-owner\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,98,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t3-ordinary-build-census-cmp_t3_rep1-r2",
          "performer": "claude:cmp_t3_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t3_rep1/tip2/t3-ordinary-build-census.txt",
            "digest": "sha256:eb6e5508e8627f98f5fc6b61da75fc22c8442a2789409af23f1612140c69defd",
            "excerpt": "$ bench test --check ordinary-build-census\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,a026c7d19bc6d4458266691778c068e7d4c9b26e,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,319,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t3-ordinary-build-census",
          "command": "bench test --check ordinary-build-census",
          "exit_code": 0
        },
        {
          "id": "t4-guidance-cmp_t4-r3",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip3/t4-guidance.txt",
            "digest": "sha256:0fac0b9516a21d8dd59ff12e4ff8090a660a2900c04a26c5bca1323e81cf284f",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check docs-currency-workflow\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,64840c504e88b81504a8f895bb3adee6d3d62aec,false\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4953,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench worktree exec \"commit-merge-parent-build\" -- bench anchors .bench/BENCH-reference.md\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,64840c504e88b81504a8f895bb3adee6d3d62aec,false\nanchors[47]{kind,section,step,needle,line}:\n  require,\"\",0,ship only in the Bench kit repository,25\n  require,\"\",0,a linked repo upgrades with `bench upgrade`,25\n  require,\"\",0,`.agents/skills/bench-craft-synthesis/SKILL.md` (kit-only),98\n  require,\"\",0,Do not append to,43\n  require,\"\",0,The `census` signal counts raw calls per assignment from `$BENCH_HOME/census/<repo-key>/`.,39\n  require,\"\",0,/bench-final-check` writes `capture/retros/<spec-slug>.md`,45\n  require,\"\",0,/bench-drain` owns their reviewed drain,47\n  require,\"\",0,A reviewed spec-backed build keeps its serial ticket commits in one retained integration source,277\n  require,\"\",0,\"no separate conformance phase, driver, or per-check evidence partition\",400\n  require,\"\",0,The spec is optional on the landing and on its resume,289\n  require,\"\",0,constructs through the shared refusal-route registry,296\n  require-in-section,Command Notes,0,\"runs the fast lane on a private checkout of the composed snapshot, and a lane pass\",280\n  forbid,\"\",0,may inherit exact per-check evidence,0\n  require,\"\",0,bench harnesses codex,537\n  require,\"\",0,The destructive-git guard allows an agent push to any branch other than the default branch.,527\n  require,\"\",0,\"A red run prints one `failures[N]{phase,line}` table, and then `gate: red`.\",468\n  require,\"\",0,\"guidance, not rules\",19\n  require,\"\",0,with authority you do not have,27\n  require,\"\",0,Read `CONTEXT.md` for the mental model,13\n  require,\"\",0,Setup and adoption connect a repository to the kit,216\n  require,\"\",0,Context commands expose current state,217\n  require,\"\",0,Oracle commands inspect or enforce readiness,218\n  require,\"\",0,Work commands own isolated execution,219\n  forbid,\"\",0,progressive loading,0\n  require,\"\",0,the plan amendment declares version 2 with a delegate execution block. The amendment splits each chunk verification into one verification for each ticket.,277\n  forbid,\"\",0,\"the rule is guidance, not a hook\",0\n  forbid,\"\",0,\"`bench commit` enforces this boundary: it refuses the primary checkout\",0\n  require,\"\",0,\"Every phase lands this way, and `.bench/BENCH.md` states how `bench commit` enforces that rule.\",291\n  forbid,\"\",0,\"A stale Bench executable is rebuilt, and the landing re-runs under it.\",0\n  require,\"\",0,the reviewer merges the destination into the source worktree with raw Git,310\n  forbid,\"\",0,owns their reviewed drain and its capture commit,0\n  require,\"\",0,`/bench-final-check` states when a tracked or an ignored retro commits.,48\n  forbid,\"\",0,A tracked retro commits with the phase close,0\n  require,\"\",0,The primary checkout owns the `main` section.,361\n  require,\"\",0,The verb keeps a non-empty Next command.,362\n  require,\"\",0,refuses a State that pins a commit outside the tip's ancestry,363\n  require-in-section,Plumbing subcommands,0,`bench gate-prose` is the one internal verb that a session runs directly; Command Notes gives its forms.,383\n  require,\"\",0,\"A green run prints one `phases[N]{phase,verdict,elapsed_ms}` table\",472\n  require,\"\",0,`bench handoff` rewrites only the calling worktree's assignment section.,360\n  forbid,\"\",0,\"`bench handoff [--harness <name>]\",0\n  forbid,\"\",0,`bench worktree reset --to <commit> <target>` plans,0\n  forbid,\"\",0,`bench retro <slug> (--body <markdown>,0\n  require,\"\",0,\"Before the first dispatch, the plan amendment makes the authored version 1 fence a version 2 plan.\",409\n  forbid,\"\",0,`.claude/skills/` carries only the,0\n  forbid,\"\",0,because Claude already has each phase as a command,0\n  forbid,\"\",0,`.claude/skills/` links every `.agents/skills/` skill that has no same-named command.,0\n  require,\"\",0,`.claude/README.md` states which skills `.claude/skills/` links.,73\nhelp[0]{cmd,why}:\n\nNote: every require row has a nonzero line and every forbid row reads 0. The route needle sits at line 310 because the three new physical lines of ticket 4 sit above it; the needle text is intact.\n"
          },
          "requirement": "t4-guidance",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "t4-skills-index-cmp_t4-r3",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip3/t4-skills-index.txt",
            "digest": "sha256:967854a60e7b2789c2640555bfdba613294b2c241d911709e5482a703df8b2d8",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check skills-index-command-adapters\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,64840c504e88b81504a8f895bb3adee6d3d62aec,false\ncheck[1]{name,kind,tests_run,subjects}:\n  skills-index-command-adapters,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,13,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-skills-index",
          "command": "bench test --check skills-index-command-adapters",
          "exit_code": 0
        },
        {
          "id": "t4-floor-cmp_t4-r3",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip3/t4-floor.txt",
            "digest": "sha256:9202a0746d0553bdc28ade4b79459bd976dd1b80211542e53ec15809e8c4d39b",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,64840c504e88b81504a8f895bb3adee6d3d62aec,false\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15832,333\n  github.com/gibbonmi/bench/internal/anchors,pass,2158,179\n  github.com/gibbonmi/bench/internal/conformance,pass,56507,1685\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t4-floor",
          "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "t4-root-conformance-cmp_t4-r3",
          "performer": "claude:cmp_t4",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_t4/tip3/t4-root-conformance.txt",
            "digest": "sha256:8edbd717b801a758721cf954f4c694e42193d655d9c8a13a8071b3024ebfeddb",
            "excerpt": "$ bench worktree exec \"commit-merge-parent-build\" -- bench test --check root-conformance\nexit_code=0\ntree[1]{target,head,dirty}:\n  commit-merge-parent-build,64840c504e88b81504a8f895bb3adee6d3d62aec,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15137,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t4-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "mp2-standards-r1",
          "performer": "claude:cmp_mp2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_standards",
            "digest": "sha256:58acb62628a80ab18a6be670211d273ea584549e388308243a1d323836477abf",
            "excerpt": "MP2 Standards axis, claude:cmp_mp2_standards, opus/high, one iteration, frozen pair 004e4e94..fabb55f3 (chunk tip eb7c791c), evidence sha256:93aa46bc81afba73fd184bf0df7e7ab4a4592e2dd2de5aac6fed6eeae1e4016b current=true, clean tree.\n\nFindings: 8 (S1-S4 auto-fix, S8 ask-user, S5-S7 no-op). Worst issue: S3. None changes behavior.\n\nS1 (auto-fix, confidence 6). craft-comments follows ste-prose.md (25-word descriptive sentences; a vertical list for three or more parallel items). The edited commit-named-path doc comment at internal/refusalroute/faces_commit.go:44-47 is one 40-word sentence that lists five causes in a comma run. Ticket 2 line 39 requires the cause, not the form.\nS2 (auto-fix, confidence 4). The same 25-word limit on new comment prose: internal/commit/merge_parent_test.go:16-19 (45 words), :35-39 (about 30 and 27), :158-162 (about 29); internal/landing/pending_merge_fixture_test.go:59-62 (about 35); internal/landing/pending_merge_test.go:43-46 (about 32). Nearby repo comments often break the same limit.\nS3 (auto-fix, confidence 5). AGENTS.md one source per fact; Duplicated Code. internal/commit/lane_test.go:87-90 owns head(t, root) with 8 call sites; merge_parent_test.go re-derives it as strings.TrimSpace(string(runGit(t, ..., \"rev-parse\", \"HEAD\"))) at lines 57, 77, 153 and as revParse(\"HEAD\") at 176 and 193.\nS4 (auto-fix, confidence 4). craft-comments \"One source owns a fact\": the composedSnapshot doc at landing.go:58 (\"the merge heads that its name-set check read\") and the composeAuthorized sentence at landing.go:126-127 (\"The snapshot keeps the merge heads that the attribution read\") state the same fact; ticket 3 line 37 asks that doc to state the waiver, recheck, parents, and completion step.\nS5 (no-op, confidence 8). One merge fixture and its constants per package (startPendingMerge vs mergeCheckout/pendingMerge; mergeConflictPath, mergeIncoming, mergeNameSet, mergeNameSetRefusal) follow ticket 2 lines 46-47 and 60 and ticket 3 lines 47 and 51. stagedPaths is the only listing helper (attribution.go:70, :74); mergeStatePath has no earlier copy; the commit side reuses adminPath (refusal_route_test.go:242).\nS6 (no-op, confidence 3). Long Function: Owner.Land at landing.go:168-208 is 40 lines in six straight steps; ticket 3 line 33 keeps commit-tree there, and line 40 moves a helper out only past 400 lines (393 now).\nS7 (no-op, confidence 3). Primitive Obsession and Data Clumps: attributedPaths returns ([]string, []string, error) (attribution.go:35); ticket 3 line 16 and git.PendingMerge's []string answer decide the shape.\nS8 (ask-user, confidence 4). AGENTS.md: an independently authored test expectation needs a recorded and demonstrated red. mergeNameSetRefusal (merge_parent_test.go:31, pending_merge_test.go:16), const changed = \"pending merge changed during the commit\" (pending_merge_test.go:131), and the \"read pending merge:\" prefixes spell refusal sentences. The record holds one matching red (the landing omission at reviews/commit-merge-parent.md:588) and none for the recheck sentence or the commit-side expectation. Spec line 619 asks one recorded red for each independent expectation; tickets 2 line 22 and 3 line 26 fix each sentence as a spec contract.\n\nAdvice: refusedLine (merge_parent_test.go:70-75) builds its own name-to-face lookup like refusal_route_test.go:117-120; \"tracked.txt\" appears 12 times in merge_parent_test.go; the stagedPaths error text ties a generic helper to one caller; both fixture comments end with \"no shared helper crosses the package boundary\", a reviewer-facing clause; the plant closure in TestLandRefusesAMergeStateChangeDuringAuthorization has no t.Helper().\n\nRule checks: the two new .bench/BENCH-reference.md sentences (lines 285-287) are 15 and 13 words, active, present tense, word for word from ticket 4 line 15, paragraph of six sentences, no other doc states the rule (only .bench/commitment.json:456 planning data); whether \"further\" is an approved STE word was not checked. Unchanged surfaces pass: commit.go and landing_reviewed_test.go empty diff stat; faces_commit.go comment only; new literals hold no next=, next[, run bench , or ; then . git-plumbing-owner pass: no --git-path in the delta; compositionState now reads through git.AdminPath. Budgets pass: landing.go 393, landing_helpers_test.go 388, pending_merge_test.go 255, pending_merge_fixture_test.go 132, merge_parent_test.go 196, attribution.go 278, commitment_test.go 247, faces_commit.go 70; landing_reviewed_test.go 437 with no line added. No FT<n> story tag. MP56 verdict: pass (11 changed files, none under ROADMAP.md or roadmap/, commit.go untouched).\n\nExamined: the full delta for internal and .bench/BENCH-reference.md; tickets 2, 3, 4; spec lines 600-629; AGENTS.md, .bench/BENCH.md, projects/benchkit.md (Seams, Gate), craft-review and finding-discipline.md, craft-comments, ste-prose.md; checkout.go, staged.go, gitexec.go, landing.go:1-215, attribution.go:1-111, refusal_route_test.go:95-254, routetest.go:115-169, lane_test.go:78-93, the git-plumbing-owner check; rg searches; the mutation entries of reviews/commit-merge-parent.md.\nEvidence cursors fetched: none beyond --check-current; ten spilled command outputs read (one listing only, one not opened).\n"
          },
          "axis": "Standards",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "eb7c791c83709ec159c6e8dd98592170b8bedc7e",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5",
            "S6",
            "S7",
            "S8"
          ],
          "supersedes": []
        },
        {
          "id": "mp2-spec-r1",
          "performer": "claude:cmp_mp2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_spec",
            "digest": "sha256:1f84413eec0c0d67196c25759401632d0d0d58286df95f5d0efcb849fe20f71e",
            "excerpt": "MP2 Spec axis, claude:cmp_mp2_spec, opus/high, one iteration, frozen pair 004e4e94..fabb55f3 (chunk tip eb7c791c), evidence sha256:93aa46bc81afba73fd184bf0df7e7ab4a4592e2dd2de5aac6fed6eeae1e4016b current=true, clean tree.\n\nFindings: 1. Worst issue: P1 (documentation only).\n\nP1 (auto-fix, confidence 8). Spec line 350 (MP42 why cell): \"A newline-framed git diff-index read C-quotes that name, so the correct name set refuses.\" Both listings come from one helper (stagedPaths at attribution.go:97-109; merged at :70, covered at :74), as spec line 171 requires. Trace: drop -z and split on newlines; both listings carry the same C-quoted string, covers[path] still matches at :83-87, the name set publishes, and the MP42 row stays green. MP41 catches that mutation (merge_parent_test.go:162). The row is met; only its rationale is wrong. Fix: reword the MP42 why cell to what it guards (a correct literal spelling of a hostile merged path publishes) and point the framing catch at MP41; a spec-text plan commit with no behavior change.\n\nAdvice: spec line 233 \"the route anchor at line 307\" is now 310 (reviews/commit-merge-parent.md:510 retains the anchors output); spec lines 172 (attribution.go:154, :177) and 186 (landing.go:139-141, now :146-148) cite pre-edit line numbers. For Coverage: no row seems to catch a dropped :(literal) prefix at attribution.go:100 (a glob name such as * would cover every merged path; unverified reading of Git's matcher). Two plain-error sentences are not in the promised labels (list the paths of the pending merge: at attribution.go:104; end the pending merge: at landing.go:196); both wrap required failure modes and carry no route token; not scope creep.\n\nRequirement trace, no further defect: the D1 listing flags at attribution.go:98 read raw through outputRaw (gitexec.go:40-42) and split on NUL (:107-108); the covered listing appends :(literal)+name (:99-101); the complement walks the merged listing in order (:83-87); NamedPathError with the exact prefix, strconv.Quote per path, one-space joins (:85, :89); a reader refusal returns a plain error (:66-68) that landingFace maps to handback (commit.go:215-229); D1 runs after the name checks and before compose (attribution.go:53-58, landing.go:133-137) and before AuthorizeCandidate (commit.go:136, :140, :144); the empty-diff waiver applies only during a merge (landing.go:146) and DryRun reuses composeAuthorized (:158-161); the recheck sentence, plain error, and reader refusal all refuse before commit-tree (:173-179); parents -p expected then each head in file order (:181-185); git merge --quit after updateRef and before reconcile, only when heads exist (:189-199); every failure after the ref update returns *PublishedUnreconciledError (:196, :201-205). Scope: the diff file list holds no commit.go, help, record, or route text; every path sits inside a ticket Writes: line or the review pickup.\n\nRow audit: MP19-MP50, MP52-MP54, MP62 met; each seam cell cites an existing test in citation form; MP22 (landing_test.go:336), MP31 (state_test.go:63-71), MP43 (landing_test.go:69) exist. MP40 discriminates the ordering (conflict.txt is not a planning path, commitment/model.go:118). Review-owned: MP51 met (run(r.Root, \"merge\", \"--quit\") at landing.go:194-198; the only os.Remove calls in non-test landing code clean temp directories). MP55 met (.bench/BENCH-reference.md:285-287 holds the two sentences after the fast-lane note, six-sentence paragraph; the retained anchors excerpt shows every require row nonzero and every forbid row 0, Command Notes needle at 280, route needle at 310). MP56 met (no ROADMAP.md, roadmap/, or internal/commit/commit.go in the diff or the fence).\n\nExamined: spec lines 1-260 and 296-640; tickets 2, 3, 4 in full and ticket 5's Writes line; attribution.go, landing.go:100-209, gitexec.go, checkout.go, commit.go:105-234, state_test.go:18-97, commitment_test.go:40-129, commitment/model.go:105-139, BENCH-reference.md:272-311; the three new test files; reviews/commit-merge-parent.md:485-574; craft-review Spec axis; finding-discipline.md. No test or probe run.\nEvidence cursors fetched: the full delta spill; the diff name list; the symbol grep; the coverage projection (lines 60-113); the planningCommitRepo uses. Not read: the spec heading list and the review-record grep spills.\n"
          },
          "axis": "Spec",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "eb7c791c83709ec159c6e8dd98592170b8bedc7e",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "mp2-coverage-r1",
          "performer": "claude:cmp_mp2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d8900dd4fbecce79b12367bfe698ef54d459d9ad",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_coverage",
            "digest": "sha256:595369c2c812d7dfe516f99f9374e878e52a4afa139c6cd180bb326eb99360c1",
            "excerpt": "MP2 Coverage axis, claude:cmp_mp2_coverage, opus/high, one iteration, frozen pair 004e4e94..fabb55f3 (chunk tip eb7c791c), evidence sha256:93aa46bc81afba73fd184bf0df7e7ab4a4592e2dd2de5aac6fed6eeae1e4016b current=true, dirty=false.\n\nEnumeration: git.PendingMerge gives absent, heads in file order, or a read pending merge: refusal; coverPendingMerge runs diff-index --cached --name-only --no-renames -z twice (merged listing with no names, covered listing with one :(literal) pathspec per name) and builds the refusal with strconv.Quote (attribution.go:65-109); names are cleaned by repositoryPath and de-duplicated and sorted by unique (:127-141, :267-278); the heads pass through composedSnapshot.mergeHeads to the waiver (landing.go:146), the recheck (:173-179), the commit-tree parents (:181-185), merge --quit (:194-198), then reconcile (:199). Write set: every changed file sits inside the Writes: line of the ticket that changed it; ticket 2 did not touch landing.go.\n\nC1 (auto-fix, confidence 5). State: a pending octopus merge whose name set omits a merged path. Bypass: swap `if err != nil || len(heads) == 0 { return nil, err` for `if err != nil || len(heads) != 1 { return heads, err`: silent, restored yes, 22 landing tests pass; every commit-package merge journey uses the one-head startPendingMerge. Needed: an octopus refusal row in TestLandRefusesAnIncompleteMergeNameSet (octopusMerge naming only side-a and expecting \"side-b\").\nC2 (no-op, confidence 8). Newline framing in both listings: the MP42 reason is false (both listings C-quote alike, coverage still matches, MP42 stays green; landing run silent, commit run bit on MP41 two_paths_omitted). MP42 still catches a dropped :(literal) (the ticket 2 probe). Fix: amend the MP42 reason to cite the :(literal) swap; non-behavioral spec inaccuracy flagged for veto.\nC3 (auto-fix, confidence 4). State: during authorization PendingMerge returns the same heads in a different order; the spec says another merge head list refuses and the code does (slices.Equal, landing.go:177). Bypass: compare sorted copies: silent, restored yes, 8 tests. Needed: a \"reordered merge heads\" row in TestLandRefusesAMergeStateChangeDuringAuthorization from octopusMerge with MERGE_HEAD lines reversed.\n\nAdvice: the command-level waiver holds (an ours resolution in TestCommitRecordsTheMergeParent still published a two-parent merge; silent exploration), and no MP29-style row exists at the Command seam; Owner.DryRun and bench commit --dry-run have no positive row during a pending merge (no non-contrived mutation separates it from Land); a directory that covers a deleted path below it was not run; a merged path that the caller unstages drops out of the index listing and the merge commit records ours for it (observation B hazard), and no Won't handle line names this; the AuthorizeCandidate read of the heads (commit.go:144) precedes composeAuthorized and the recheck does not cover that window (for the MP3 review); a stale MERGE_HEAD from a pre-fix one-parent commit makes the next ordinary commit a merge commit (allowed by .pending, no row pins it); a merge of thousands of paths produces one unbounded refusal line.\n\nProbes (all restored yes): 1 landing.go waiver == 0 to >= 0, run UnchangedTree and PreAuthorizationRefusalTable: bit (MP29, MP30), 2 failed. 2 waiver to < 0: bit (MP31 symlink row), 1 failed. 3 attribution.go D1 skipped when len(heads) != 1 with heads kept: silent (C1). 4 merge_parent_test.go ours resolution (exploration): silent. 5 attribution.go newline framing in both listings, landing IncompleteMergeNameSet: silent (C2). 6 same in commit: bit (MP41 two_paths_omitted), 1 failed. 7 attribution.go D1 always on, TestLandComposesOnlyNamedPathsAndCASesExpectedBase: bit (MP43), 1 failed. 8 landing.go recheck compares sorted copies: silent (C3). Two device rows skipped for privilege. Final git status --short: empty.\n\nFindings: 3 (C1 auto-fix, C2 no-op, C3 auto-fix). Worst issue: C1.\n\nConsumer rows touched=false walked: commit.go:136, :140, :187 (faces map NamedPathError to commit-named-path and a plain error to commit-handback); the Land callers in landing_test.go, state_test.go, composition_test.go, landing_reviewed_test.go (no pending merge, D1 does nothing); composedSnapshot uses at landing_test.go:289 and :315 (parameter type only); compositionState at composition_test.go:221 and :229 (AdminPath, same path in a primary checkout); refusalroute/registry.go:204 (comment only); the landingRepo callers (label value unchanged).\n\nExamined: the internal delta and per-commit file list; spec rows MP19-MP56 and MP62, Edge inventory, Testing decisions, Observations; tickets 2, 3, 4; the hostile-input checklist; craft-review Coverage axis and finding-discipline.md; attribution.go (full), landing.go:100-208, gitexec.go, lane_owner.go, checkout.go:28-73, commit.go:95-229; pending_merge_test.go, pending_merge_fixture_test.go, merge_parent_test.go, landing_test.go (MP43 at :69, MP22 at :336).\nEvidence cursors fetched: the summary, manifest page v1.93aa46bc...4016b.m.0.0, source s20 page 0 (consumers), s.20.1 (stream_end). Not fetched: m.0.1 and any other source page.\n"
          },
          "axis": "Coverage",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "eb7c791c83709ec159c6e8dd98592170b8bedc7e",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
          ],
          "supersedes": []
        },
        {
          "id": "mp2-coverage-r2",
          "performer": "claude:cmp_mp2_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_coverage_r2",
            "digest": "sha256:d90c47e4d9006616aa211133009ae6d7e8d485d3c0e78177e237d046e6c4a76f",
            "excerpt": "MP2 Coverage axis round 2 (confirming), claude:cmp_mp2_coverage_r2, opus/high, one iteration, frozen pair 004e4e94..54dc1636 (chunk tip 83d43fba), evidence sha256:1bb2d23a76b8f3572503775882a30c0c1a9f2b77f1ecdf3e65d5a55cd19e0bec current=true, clean tree. The repair delta fabb55f3..54dc1636 changes test files and comments only; the one production edit is the comment at landing.go:125-127.\n\nC1 confirmed closed: row \"octopus side path omitted\" at pending_merge_test.go:68 names side-a and expects \"side-b\", the two paths octopusMerge adds (pending_merge_fixture_test.go:80). Probe (swap attribution.go:88 `if len(omitted) > 0 {` for `if len(omitted) > 0 && len(heads) == 1 {`, run ^TestLandRefusesAnIncompleteMergeNameSet$): bit, 1 failed (octopus_side_path_omitted), restored yes.\nC3 confirmed closed: row \"reordered merge heads\" at pending_merge_test.go:183 starts from octopusMerge; the reverse stub (:150-154) reads the real MERGE_HEAD lines, reverses them, and writes them back through plantHeads. Probe (swap landing.go:176 slices.Equal for an order-blind set comparison with slices.ContainsFunc, run ^TestLandRefusesAMergeStateChangeDuringAuthorization$): bit, 1 failed (reordered_merge_heads), restored yes.\n\nFocused runs on the tip: landing IncompleteMergeNameSet|MergeStateChange passed, 13 tests; commit RecordsTheMergeParent|RefusesAnIncompleteMergeNameSet passed, 5 tests (the S3 head helper swap kept them green; the rev-parse mergeIncoming call at merge_parent_test.go:78 stays inline because head reads only HEAD).\n\nNew findings on the repair delta: 0. Worst issue: none. Final git status --short: empty.\n\nAdvice: the C1 row writes side-a and \"side-b\" as literals while the octopus closure discards the added names that octopusMerge returns (:52-55); a rename turns the row red rather than silent, so not a coverage gap; a Standards call.\n\nExamined: the full repair delta for internal; reviews/commit-merge-parent.md:1256-1360 and the MP2 Coverage round 1 excerpt at :1183; attribution.go:30-114, landing.go:160-199, pending_merge_fixture_test.go (whole), pending_merge_test.go:40-269, lane_test.go:84-90.\nEvidence cursors fetched: --check-current only; the diff spill read in full; the two probe spills not read (inline verdicts complete).\n"
          },
          "axis": "Coverage",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "83d43fbaf4cbf79f571cafdfeb054cc5ad434120",
          "finding_ids": [],
          "supersedes": [
            "mp2-coverage-r1"
          ]
        },
        {
          "id": "mp2-standards-r2",
          "performer": "claude:cmp_mp2_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_standards_r2",
            "digest": "sha256:4920f6ebd63b1fdde9ef7d21aa7a7a2cc00b6d72541505b0e540f6d0c87089dc",
            "excerpt": "MP2 Standards axis round 2 (confirming), claude:cmp_mp2_standards_r2, opus/high, one iteration, frozen pair 004e4e94..54dc1636 (chunk tip 83d43fba), evidence sha256:1bb2d23a76b8f3572503775882a30c0c1a9f2b77f1ecdf3e65d5a55cd19e0bec current=true, clean tree. Repair delta fabb55f3..54dc1636 read for internal and the spec.\n\nS1 confirmed closed: faces_commit.go:44-47 now holds five sentences of 11, 14, 17, 11, and 5 words; commit 188a3f55 changes those four comment lines only.\nS2 confirmed closed: every sentence in the five comments is 25 words or fewer (merge_parent_test.go:16-19, :33-38, :157-161; pending_merge_fixture_test.go:57-61; pending_merge_test.go:43-48, six sentences with the coverage parenthetical).\nS3 confirmed closed: merge_parent_test.go calls head(t, ...) at lines 56, 76, 152, 174, 191 (lane_test.go:87-90); revParse is gone; the one rev-parse left (:168) resolves the branch mergeIncoming, not HEAD.\nS4 confirmed closed: the merge-heads fact appears only in the composedSnapshot doc at landing.go:58; the composeAuthorized doc at :125-127 keeps the steps and the waiver.\nS8 confirmed closed (evidence-only): reviews/commit-merge-parent.md:1299-1329 keeps both sentence-swap probes on the clean tip fabb55f3, each red on 3 tests with the subject restored.\n\nScope: record commits 19a56959, bdf5f4ba, bc50e71e, 27aaedba, 54dc1636 touch only reviews/commit-merge-parent.md; plan commit 1063cbb1 changes one plan JSON line; spec amendment eb0321e5 changes the MP42 why cell (spec:350), five citation lines (spec:154, 156, 172, 186, 233), and one Won't handle line (spec:398); none changes production code. The five citations are accurate at the tip. Budgets: merge_parent_test.go 194, landing.go 392, pending_merge_fixture_test.go 131, pending_merge_test.go 279; landing_reviewed_test.go 437 with an empty diff. One source per fact: plantHeads (pending_merge_test.go:142-146) is the only MERGE_HEAD writer in the landing tests; the start field's three starts (:156-168) are the old inline branch moved into the table; the two octopus adapters (:52-55, :165-168) repeat text, not knowledge.\n\nNew finding: S21 (auto-fix, confidence 3). ste-prose.md:23 (no telegraphic style) through craft-comments: merge_parent_test.go:16 \"The paths and the branch of a pending merge journey.\" has no verb after the S2 split. Fix: give it a verb.\n\nAdvice: faces_commit.go:45-46 could take a five-item vertical list; reverse (pending_merge_test.go:150-154) splits MERGE_HEAD with strings.Fields beside plantHeads and could read through benchgit.PendingMerge; plantHeads and plant have no t.Helper(); the octopus adapter at :52 ignores added; landing.go:58 \"its name-set check\" could name composeAuthorized; pending_merge_test.go:137-138 could make Owner.Land the subject.\n\nFindings: 1 new (S21). Worst issue: S21, no behavior effect.\nExamined: the repair-delta spill; the per-commit log with stats and the diffs of eb0321e5 and 1063cbb1; reviews/commit-merge-parent.md:1256-1355 and the round 1 excerpts; craft-comments and ste-prose.md; merge_parent_test.go, pending_merge_test.go, pending_merge_fixture_test.go in full; landing.go:50-151, faces_commit.go:38-53, lane_test.go:84-90, landing_helpers_test.go:290-319, land_refusal.go:74-85, BENCH-reference.md:278-311; rg searches; an empty diff stat for landing_reviewed_test.go.\nEvidence cursors fetched: none beyond --check-current; seven spilled outputs read, one not opened.\n"
          },
          "axis": "Standards",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "83d43fbaf4cbf79f571cafdfeb054cc5ad434120",
          "finding_ids": [
            "S21"
          ],
          "supersedes": [
            "mp2-standards-r1"
          ]
        },
        {
          "id": "mp2-spec-r2",
          "performer": "claude:cmp_mp2_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_spec_r2",
            "digest": "sha256:4edc09f83e63656b1759df54cdb48e2fc8492b5a37e7a45e879f219c0bd5351f",
            "excerpt": "MP2 Spec axis round 2 (confirming), claude:cmp_mp2_spec_r2, opus/high, one iteration, frozen pair 004e4e94..54dc1636 (chunk tip 83d43fba), evidence sha256:1bb2d23a76b8f3572503775882a30c0c1a9f2b77f1ecdf3e65d5a55cd19e0bec current=true, clean tree.\n\n1. P1/C2 (MP42 why cell): confirmed, with the P22 caveat. The cell is accurate: TestLandRefusesAnIncompleteMergeNameSet passes the hostile name as written and expects a publication (pending_merge_test.go:51, :67), so a rule that changes the name's bytes before the :(literal) pathspec refuses it; the framing catch sits with MP41 (merge_parent_test.go:116, :124; Coverage probe 6 at reviews/commit-merge-parent.md:1183).\n2. Line citations: four confirmed (land_refusal.go:80-83 and :76-79; attribution.go:213 and :234; BENCH-reference.md:310 holds the route needle of registry_retained_workflow.go:206); one finding P21.\n3. The new Won't handle line (spec.md:398): confirmed. A merged path is what git diff-index --cached lists against the base (Terms, spec.md:111); line 118 says a path staged after the merge joins, so an unstaged path leaves by the same rule; D1 (spec.md:130) measures through that listing; stories 20-23 and 32 and spec.md:389 hold; .pending is untouched.\n4. Production behavior unchanged: confirmed. Only landing.go (one doc sentence removed at :125-127) and faces_commit.go (comment, :44-47) are production files in the delta. The octopus refusal row (pending_merge_test.go:68) takes side-a and side-b from octopusMerge (fixture :80-89); the reordered-heads row (:183) reverses MERGE_HEAD during authorization (:150-154) and expects the recheck sentence that slices.Equal at landing.go:176 produces; merge_parent_test.go calls the existing head helper.\n5. Seam cells: confirmed. MP19-MP54 and MP62 cite 16 existing tests in the citation form (merge_parent_test.go :85, :101, :133, :145, :162; pending_merge_test.go :31, :49, :96, :121, :139, :208, :243, :266; landing_test.go :69, :336; state_test.go:21). MP51 met (git merge --quit at landing.go:193-197; os.RemoveAll only on temp directories), MP55 met (BENCH-reference.md:285-287), MP56 met.\n\nP21 (auto-fix, confidence 9). Spec line 186 cites landing.go:146-148; at the tip the refusal is lines 145-147 (the repair commit 83d43fba removed one doc line above the block). Fix: a spec-text change to landing.go:145-147.\nP22 (no-op on the spec, confidence 5). The MP42 cell is true as worded, but the record's C2 excerpt says \"MP42 still catches a dropped :(literal) (the ticket 2 probe)\" and no recorded probe shows MP42 red under a :(literal) swap or drop; on the axis's reading of Git's matcher (match_pathspec_item, from memory), a plain drop leaves the hostile row publishing. Disposition: leave the cell; record the probe or remove the claim.\n\nAdvice: a dropped :(literal) on the covered listing (attribution.go:100) looks uncaught; a refusal row where a glob name covers another merged path (side-* in the octopus fixture) would catch it. Spec line 197 (landing.go:247-257) and line 17 (landing.go:163) look like base-state citations the refresh did not touch.\n\nFindings: 2 (P21 auto-fix, P22 no-op). Worst issue: P21. Neither is behavioral.\nExamined: --check-current; the repair delta for internal and the spec with its file list; the git log base to tip; attribution.go:1-269; landing.go:120-261; land_refusal.go:70-89; BENCH-reference.md:276-313; the anchor needles; spec.md:68-183 and 327-399; pending_merge_test.go and pending_merge_fixture_test.go in full; merge_parent_test.go:14-148; lane_test.go:85-90; the probe entries of reviews/commit-merge-parent.md; an os.Remove sweep. No test or probe run.\nEvidence cursors fetched: none beyond --check-current; every spilled command output read in full.\n"
          },
          "axis": "Spec",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "83d43fbaf4cbf79f571cafdfeb054cc5ad434120",
          "finding_ids": [
            "P21",
            "P22"
          ],
          "supersedes": [
            "mp2-spec-r1"
          ]
        },
        {
          "id": "mp2-coverage-r2b",
          "performer": "claude:cmp_mp2_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "d75d13e7419a494ab08ca5090d3bf6b660deff16",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_coverage_r2/amended",
            "digest": "sha256:963ed9d6fef92a461d2d0123d803ca28efa82b60cbabdb95fbe652bc4974db6c",
            "excerpt": "MP2 Coverage axis round 2, amended read-only pass, claude:cmp_mp2_coverage_r2, opus/high, frozen pair 004e4e94..54dc1636 (chunk tip 83d43fba), evidence sha256:1bb2d23a76b8f3572503775882a30c0c1a9f2b77f1ecdf3e65d5a55cd19e0bec. No probe and no test run this pass; C1 and C3 stay confirmed closed as in the first pass.\n\nC21 (auto-fix, confidence 7): no row catches a covered listing that drops :(literal). Input: during a pending merge the incoming commit adds a file whose literal name looks like a glob (side-*) and a sibling that the glob also matches (side-a); the name set names side-* and not side-a. Expected break: with the prefix dropped at attribution.go:100, Git reads side-* as a glob in the covered listing, side-a counts as covered, and D1 stays silent; compose still adds each name as :(literal) (attribution.go:213), so the merge commit takes side-* and the base version of side-a, and the next landing undoes side-a (the observation B hazard). The orchestrator's probe (swap `\":(literal)\"+name` for `name`, run ^TestLandRefusesAnIncompleteMergeNameSet$) was silent with 7 passed and restored yes; the hostile name a b[x]*\\né.txt stays covered without the prefix because Git tries an exact match first (from the axis's knowledge of Git, not read in this tree). Why a witness is needed: spec.md:171-172 decides that one :(literal) pathspec for each name decides coverage, the rule compose and reconcile use (attribution.go:213, :234-253); projects/benchkit.md:181 names glob characters; spec line 375 claims MP41 and MP42 cover glob characters, but no row tests glob matching. Row to add: a refusal row in TestLandRefusesAnIncompleteMergeNameSet, \"glob-shaped name does not cover a sibling\", with the incoming commit adding side-* and side-a, the name set naming side-* and not side-a, and the refusal naming \"side-a\" (spec observation 7 at spec.md:604 shows :(literal)name* printing nothing). Confirm the fold by rerunning the prefix-drop swap; it must bite this row.\n\nFinding count: 1 (C21). Worst issue: C21.\nRead this pass: spec.md:160-179 and the literal and glob hits at :78, :171, :349-350, :375, :604, :638; attribution.go:200-219 and :100, :213, :234-253; projects/benchkit.md:181; from the earlier pass pending_merge_test.go:40-90 and pending_merge_fixture_test.go:35-71. Evidence cursors fetched this pass: none.\n"
          },
          "axis": "Coverage",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "83d43fbaf4cbf79f571cafdfeb054cc5ad434120",
          "finding_ids": [
            "C21"
          ],
          "supersedes": [
            "mp2-coverage-r2"
          ]
        },
        {
          "id": "mp2-standards-r3",
          "performer": "claude:cmp_mp2_standards_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_standards_r3",
            "digest": "sha256:9688cc878f9a65fc9cb3008df54dba741db65c936c2d0afe3d5f624ef283dfb1",
            "excerpt": "MP2 Standards axis round 3, claude:cmp_mp2_standards_r3, opus/high, one iteration, frozen pair 004e4e94..f10aeedf (chunk tip cca1013b), evidence sha256:b697ee6e8a770961a6a50c28416889ce2042b5f6c73c2ffa5622b822d12cbc68 current=true, clean tree. Round 3 delta 54dc1636..f10aeedf read for internal and the spec.\n\nS21 confirmed closed: merge_parent_test.go:16 opens with \"These constants name the paths and the branch of a pending merge journey.\" (a verb, present tense, active); lines 16-19 are 87, 89, 85, 80 columns inside the file's width; 2df5483b changes comment lines only.\n\nNew findings on the round 3 delta: 0. Worst issue: none.\n\nChecks on the C21 repair cca1013b (pending_merge_test.go only): the new doc sentence at :46-47 (15 words, present tense, active, timeless, no narration) meets the STE and craft-comments rules and the file's width; the new row at :69 writes the fixture paths and the quoted omission as literals like the pkgx/a.go row at :67, reuses mergeConflictPath, mergeDeletedPath, and mergeNameSetRefusal, and adds no second derivation of a production fact; pending_merge_test.go is 281 lines. Plan commits 1539f297 and 91fcb8be change one spec line each; record commits f930814b, a026c7d1, 64840c50, f10aeedf change reviews/commit-merge-parent.md only.\n\nAdvice: the row's fixture path side-a shares its name with the octopus fixture's side path in the next row; a different stem would keep a reader from linking the rows (confidence 2).\n\nExamined: the round 3 diff for internal and the spec; the per-commit file stat of every commit in the range; pending_merge_test.go:1-110; craft-comments in full; comment widths and line counts of merge_parent_test.go and pending_merge_test.go.\nEvidence cursors fetched: the round 3 diff spill, the commit stat log spill, and two width-measurement spills; no preflight evidence page beyond --check-current.\n"
          },
          "axis": "Standards",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "cca1013b089a78e230bb6f15d9f3266ce0ab0ffa",
          "finding_ids": [],
          "supersedes": [
            "mp2-standards-r2"
          ]
        },
        {
          "id": "mp2-coverage-r3",
          "performer": "claude:cmp_mp2_coverage_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_coverage_r3",
            "digest": "sha256:dd5700001f3571a87ad93d35bed9aea09371bd00c59cdf06f6b351e92db46be3",
            "excerpt": "MP2 Coverage axis round 3, claude:cmp_mp2_coverage_r3, opus/high, one iteration, frozen pair 004e4e94..f10aeedf (chunk tip cca1013b), evidence sha256:b697ee6e8a770961a6a50c28416889ce2042b5f6c73c2ffa5622b822d12cbc68 current=true, clean tree. Round 3 delta 54dc1636..f10aeedf read for internal.\n\nC21 confirmed closed. Probe: bench probe internal/landing/attribution.go --swap '\":(literal)\"+name' --with '\":(glob)\"+name' --package ./internal/landing --run '^TestLandRefusesAnIncompleteMergeNameSet$': bit, 1 failed (glob-shaped_name_covers_only_itself at pending_merge_test.go:85, \"Land error = <nil>, want the named-path refusal ... \\\"side-a\\\"\"), restored yes. The row at pending_merge_test.go:69 adds side-* and side-a through pendingMerge, names only side-* with mergeConflictPath and mergeDeletedPath, and expects \"side-a\". The swap target attribution.go:100 is the only \":(literal)\"+name occurrence (lines 213 and 234 spell path). Focused run on the tip: bench test --package ./internal/landing --run 'TestLandRefusesAnIncompleteMergeNameSet' pass, 336 ms, 8 tests, 0 skips.\n\nNew findings on the round 3 delta: 0 (a comment-only rewording at merge_parent_test.go:9-16 and the C21 row with its doc lines). Worst issue: none. Final git status --short: empty.\n\nAdvice: the doc comment at pending_merge_test.go:49 still cites \"(Coverage rows MP35-MP38 and MP42.)\" and the glob row has no row ID of its own; a spec-stage call on the map at spec line 375.\n\nExamined: --check-current; the round 3 internal diff; the C21 entries of reviews/commit-merge-parent.md (lines 1588 and 1785); pending_merge_test.go:40-119; attribution.go:80-119 and the literal hits at 100, 213, 234-253; one probe and one focused run. pending_merge_fixture_test.go not read this pass.\nEvidence cursors fetched: the round 3 diff spill (read in full); the probe spill not opened (inline rows complete); no preflight evidence page beyond --check-current.\n"
          },
          "axis": "Coverage",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "cca1013b089a78e230bb6f15d9f3266ce0ab0ffa",
          "finding_ids": [],
          "supersedes": [
            "mp2-coverage-r2b"
          ]
        },
        {
          "id": "mp2-spec-r3",
          "performer": "claude:cmp_mp2_spec_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "aeb6aa53b0b6b6b07aff0a70df38bd3cecc724b0",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:cmp_mp2_spec_r3",
            "digest": "sha256:f19488a922ace58a2c799a89070bdf6900114ed33b03dbf742fc9113e00eb24a",
            "excerpt": "MP2 Spec axis round 3, claude:cmp_mp2_spec_r3, opus/high, one iteration, frozen pair 004e4e94..f10aeedf (chunk tip cca1013b), evidence sha256:b697ee6e8a770961a6a50c28416889ce2042b5f6c73c2ffa5622b822d12cbc68 current=true, clean tree.\n\n1. P21 confirmed: spec line 186 cites landing.go:145-147; at the tip line 145 is the empty-diff condition, 146 the NamedPathError return, 147 the close; plan commit 1539f297 made the change.\n2. P22 record evidence: finding P31 (the record does not retain the probe run).\n3. C21 repair row confirmed: pending_merge_test.go:69 with its doc sentence at :46-47; mergeCheckout writes each added path on the incoming branch (fixture :44-46), so the merged paths are side-* and side-a; under spec lines 171-174 side-a is the only omitted path, quoted \"side-a\"; every expected value comes from the fixture names; cca1013b changes only pending_merge_test.go (4 added, 2 removed); no production file, no MP row, no seam cell changes.\n4. Seam cells MP19-MP56 and MP62 confirmed: every test-seam cell uses the citation form; the 16 named tests exist (merge_parent_test.go :85, :101, :133, :145, :162; pending_merge_test.go :31, :50, :98, :123, :141, :210, :245, :268; landing_test.go :69, :336; state_test.go :21); MP51, MP55, MP56 keep their review-owned form.\n\nP31 (evidence-only, confidence 8). Spec line 350 (MP42 why cell) rests on the :(literal) spelling swap, but reviews/commit-merge-parent.md retains no run of it: no `Quote(name)` text in the file or its history; the only retained --swap probe is the landing.go parents swap; lines 1689 and 1766 state the claim in prose only. The AGENTS.md exception requires a recorded red. Disposition: the orchestrator reruns the swap at the current tip, retains the full output block under the P22 line in the Spec section, and reworks the prose to cite it; no spec or code change.\n\nAdvice: the C21 fold check (the prefix-drop swap biting the new row) is not retained either; the record keeps only the omit probe at cca1013b, which fails 5 tests including the glob row. Spec line 375 says MP41 and MP42 cover glob characters while the glob witness is the unmapped row :69; final reconciliation decides. The delta also refreshes the plan JSON (claude:cmp_t2_rep1) and the S21 comment.\n\nFindings: 1 (P31). Worst issue: P31, a record-evidence gap with no behavior effect.\nExamined: --check-current; the round 3 delta for internal and the spec; the commit log with stats; landing.go:125-159; pending_merge_test.go:40-119; pending_merge_fixture_test.go:28-75 and :77-90; spec lines 159-188 and 327-370; reviews/commit-merge-parent.md:1675-1788 plus searches for Quote(name), spelling, swap, mutation/probe, --swap; an existence sweep of the 16 cited tests. No test or probe run.\nEvidence cursors fetched: none beyond --check-current; eight spilled outputs read in full, two not opened.\n"
          },
          "axis": "Spec",
          "base": "004e4e94019144653de1ec8e179abe015b704242",
          "tip": "cca1013b089a78e230bb6f15d9f3266ce0ab0ffa",
          "finding_ids": [
            "P31"
          ],
          "supersedes": [
            "mp2-spec-r2"
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
    },
    {
      "from": "sha256:26c4daa7f1486d87e7341202b339d5fe7ef02fdf1a4d1fdfbc7b0f8f723304f3",
      "to": "sha256:a04f42ffc9dfb7f7f9d0fa07008db4015b2556030e54abef1763c39e19541c70",
      "chunk_ids": {
        "MP1": [
          "MP1"
        ],
        "MP2": [
          "MP2"
        ]
      }
    },
    {
      "from": "sha256:a04f42ffc9dfb7f7f9d0fa07008db4015b2556030e54abef1763c39e19541c70",
      "to": "sha256:7654187ea92c112db47b0160f0dd2cc0ea2e1645dc5648227d777f5de457dd7f",
      "chunk_ids": {
        "MP1": [
          "MP1"
        ],
        "MP2": [
          "MP2"
        ]
      }
    }
  ]
}
```

## Repair state

Chunk MP1 round 1 ran three opus/high axes in parallel at `f7b5027b` over the frozen pair `342975e3..f7b5027b`. Only the Coverage axis ran probes, and the tree was clean after every return. MP1 uses repair cycle 1 of 2 through one fresh repair session for ticket 1, which takes S1, C1, C2, C3, C4, and C5. S2 is a record gap with no code change, so the orchestrator closes it below as evidence-only work. Under `--auto-approve`, the reviewer approved every repair round and fence expansion in advance.

The ticket 1 repair session `claude:cmp_t1_rep1` closed S1 and C1 to C5 in `a329d3ab` with tests only. It refroze MP1 at that tip in `4e6fd43a`. MP1 round 2 ran three opus/high axes in parallel over the repair delta `f7b5027b..4e6fd43a`. Every axis confirmed each fold, the Coverage axis bit each fold with an independent probe, and no axis found a new finding. MP1 closes with repair cycle 1 of 2 used and no open finding.

Chunk MP2 round 1 ran three opus/high axes in parallel at `fabb55f3` over the frozen pair `004e4e94..fabb55f3`. The chunk tip is `eb7c791c`. Only the Coverage axis ran probes, and the tree was clean after every return. MP2 uses repair cycle 1 of 2 through one fresh repair session for ticket 3, which takes S2, S3, S4, C1, and C3. Every path of those findings sits inside the ticket 3 `Writes:` line.

S1 is a comment-only Go correction, so the orchestrator commits it as evidence-only work. P1 and C2 are one spec-text amendment, which the orchestrator commits as a plan commit. S8 closes evidence-only below. S5, S6, and S7 are `no-op`.

The ticket 3 repair session `claude:cmp_t3_rep1` closed S2, S3, S4, C1, and C3 in `83d43fba` with tests and comments only. It refroze MP2 at that tip in `bc50e71e`, and the ticket 2 and ticket 4 authors re-recorded their entries at that source. MP2 round 2 ran three opus/high axes in parallel over the repair delta `fabb55f3..54dc1636`. Every axis confirmed each fold.

Standards raised S21, a verbless comment sentence, which the orchestrator corrected in `2df5483b` as comment-only work. Spec raised P21, a citation one line off after the repair, which the orchestrator corrected in the plan commit `1539f297`. Spec also raised P22, a `no-op` on the record's unproven MP42 probe claim.

The orchestrator then ran two probes at `54dc1636`: the `:(literal)` spelling swap bit MP42, and the dropped `:(literal)` prefix stayed silent. The Coverage axis graded the silent probe as C21. MP2 uses repair cycle 2 of 2 through one fresh repair session for ticket 2, which takes C21. After that repair, every MP2 author re-records at the refrozen tip, and a round 3 confirms C21, S21, and P21.

The ticket 2 repair session `claude:cmp_t2_rep1` closed C21 in `cca1013b` with one test row. It refroze MP2 at that tip in `a026c7d1`, and the ticket 3 repair session and the ticket 4 author re-recorded their entries at that source. MP2 round 3 ran three opus/high axes over the delta `54dc1636..f10aeedf`. Standards confirmed S21 closed, and Coverage confirmed C21 closed with a `:(glob)` probe.

Spec confirmed P21 and raised P31, a record gap: the MP42 spelling probe had prose but no retained run. The orchestrator retains both probe runs at the frozen tip below, and the Spec axis reaffirms P31. MP2 then closes with repair cycle 2 of 2 used and no open finding.

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

MP2 findings: 8. Worst issue: S3.

- S1 (`auto-fix`, confidence 6): the `commit-named-path` doc comment in `internal/refusalroute/faces_commit.go` is one 40-word sentence with five causes in a comma run. Rule: `craft-comments` follows the STE sentence rules. Orchestrator reflow.
- S2 (`auto-fix`, confidence 4): five new test comments pass 25 words (`merge_parent_test.go:16-19`, `:35-39`, `:158-162`; `pending_merge_fixture_test.go:59-62`; `pending_merge_test.go:43-46`). Repair: split them.
- S3 (`auto-fix`, confidence 5): `merge_parent_test.go` re-derives the `head(t, root)` helper of `internal/commit/lane_test.go:87-90` at five sites. Repair: call `head`.
- S4 (`auto-fix`, confidence 4): the `composedSnapshot` doc and the `composeAuthorized` doc in `landing.go` state the same merge-heads fact. Repair: keep it once.
- S5 (`no-op`, confidence 8): one merge fixture per package follows the ticket keep decisions. Ticket 2 lines 46-47 and 60 and ticket 3 lines 47 and 51 state them.
- S6 (`no-op`, confidence 3): `Owner.Land` holds six straight steps in 40 lines. Ticket 3 line 33 keeps `commit-tree` there, and line 40 moves a helper out only past 400 lines.
- S7 (`no-op`, confidence 3): the `[]string` merge heads follow ticket 3 line 16 and the `git.PendingMerge` answer.
- S21 (`auto-fix`, confidence 3, round 2): the const-block comment at `merge_parent_test.go:16` had no verb after the S2 split. Closed in `2df5483b`.
- S8 (`ask-user`, confidence 4): the refusal-sentence expectations `pending merge changed during the commit` and the commit-side `pending merge changed paths that the commit does not name:` had no recorded red. The AGENTS.md exception needs a named mutation and a recorded red, so the orchestrator ran both sentence swaps on the clean tip `fabb55f3`:

```text
tree[1]{target,head,dirty}:
  commit-merge-parent-build,fabb55f385f1b69e66983f9245175369eff3a435,false
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/landing/landing.go,swap,failed,3,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/landing,^TestLandRefusesAMergeStateChangeDuringAuthorization$,passed,5
spilled{lines=13,bytes=1375,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/f80fcedec92401bd6cf2acfce17ca107/1791700649488204910-267ca0fb63d13b80.out}
failures[3]{package,test,line,lines}:
  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAMergeStateChangeDuringAuthorization/another_merge_head,"pending_merge_test.go:171: Land error = pending merge differs after the commit authorization, want the plain refusal \"pending merge changed during the commit\"",1
  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAMergeStateChangeDuringAuthorization/merge_appears,"pending_merge_test.go:171: Land error = pending merge differs after the commit authorization, want the plain refusal \"pending merge changed during the commit\"",1
  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAMergeStateChangeDuringAuthorization/merge_disappears,"pending_merge_test.go:171: Land error = pending merge differs after the commit authorization, want the plain refusal \"pending merge changed during the commit\"",1
skips[0]{package,test,reason}:
```

```text
tree[1]{target,head,dirty}:
  commit-merge-parent-build,fabb55f385f1b69e66983f9245175369eff3a435,false
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/landing/attribution.go,swap,failed,3,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/commit,^TestCommitRefusesAnIncompleteMergeNameSet$,passed,4
spilled{lines=13,bytes=1687,omitted_lines=4,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/f80fcedec92401bd6cf2acfce17ca107/1791700664748895789-11dd68777e96e21e.out}
failures[3]{package,test,line,lines}:
  github.com/gibbonmi/bench/internal/commit,TestCommitRefusesAnIncompleteMergeNameSet/added_path_omitted,"merge_parent_test.go:126: refusal = (\"error: pending merge left paths that the commit does not name: \\\\\"added.txt\\\\\"\", [\"correct the paths or the files that the refusal names\" \"bench commit --in 'commit-fixture' -m 'm' -- 'conflict.txt' 'trac… (455 bytes)",1
  github.com/gibbonmi/bench/internal/commit,TestCommitRefusesAnIncompleteMergeNameSet/modified_path_omitted,"merge_parent_test.go:126: refusal = (\"error: pending merge left paths that the commit does not name: \\\\\"tracked.txt\\\\\"\", [\"correct the paths or the files that the refusal names\" \"bench commit --in 'commit-fixture' -m 'm' -- 'conflict.txt' 'ad… (455 bytes)",1
  github.com/gibbonmi/bench/internal/commit,TestCommitRefusesAnIncompleteMergeNameSet/two_paths_omitted,"merge_parent_test.go:126: refusal = (\"error: pending merge left paths that the commit does not name: \\\\\"esc\\\\\\\\x1b name*\\\\\\\\n.txt\\\\\" \\\\\"tracked.txt\\\\\"\", [\"correct the paths or the files that the refusal names\" \"bench commit --in 'commit-fixture' -m… (493 bytes)",1
skips[0]{package,test,reason}:
```

## Spec

Findings: 0. Worst issue: none. Every row MP01 to MP18 is met, and every MP01 to MP15 seam cell cites an existing test.

Advice: the owner section cites `land_refusal.go:82-91` and `:78-81` in the present tense, and the function now sits at `:76-83`; final reconciliation decides a currency touch.

MP2 findings: 1. Worst issue: P1 (documentation only).

- P1 (`auto-fix`, confidence 8): the MP42 why cell says a newline-framed listing makes the correct name set refuse. Both listings share one framing, so that mutation leaves MP42 green, and MP41 catches it. Repair: a spec-text plan commit rewords the cell to the `:(literal)` coverage it guards and points the framing catch at MP41. The same commit refreshes the stale line citations that the axis listed.

- P21 (`auto-fix`, confidence 9, round 2): the `landing.go` empty-diff citation was one line off after the repair. Closed in `1539f297`.
- P22 (`no-op`, confidence 5, round 2): the record claimed an MP42 probe that no run retained. The orchestrator ran the `:(literal)` spelling swap at `54dc1636`, and it bit the hostile-path row.
- P31 (evidence-only, confidence 8, round 3): the MP42 spelling probe had no retained run. The orchestrator reran it on the clean frozen tip `f10aeedf` (source `cca1013b`), and it bit the hostile-path row with the subject restored:

```text
tree[1]{target,head,dirty}:
  commit-merge-parent-build,f10aeedf2dc9151d660516d1c7a33f9eefdec7a7,true
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/landing/attribution.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,8
spilled{lines=11,bytes=791,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/f80fcedec92401bd6cf2acfce17ca107/1791703477562998030-fe03b600eb0a1623.out}
packages[1]{package,status,elapsed_ms,tests_run}:
  github.com/gibbonmi/bench/internal/landing,fail,295,8
failures[1]{package,test,line,lines}:
  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/hostile_path_spelled_literally,"pending_merge_test.go:78: Land error = pending merge changed paths that the commit does not name: \"a b[x]*\\\\né.txt\", want a publication",1
skips[0]{package,test,reason}:
```

Review-owned rows: MP51 is met, because `Land` runs `git merge --quit` and Bench removes no admin entry file. MP55 is met, because the two sentences sit after the fast-lane note and the retained anchors excerpt shows every require row nonzero. MP56 is met, because no roadmap path and no `internal/commit/commit.go` appear in the diff or the fence.

## Coverage

Findings: 5. Worst issue: C2.

- C1 (`auto-fix`, confidence 8): no test runs the reader through a relative `gitdir:` pointer, which `spec.md` claims MP01 covers. Repair: a `TestPendingMergeReadsEveryHead` step that rewrites `linked/.git` to a relative `gitdir:` and expects the merge head.
- C2 (`auto-fix`, confidence 7): the octopus order pin depends on random SHAs, and an ascending-sort mutant stayed silent. Repair: an MP03 fixture whose file order differs from both sorted orders.
- C3 (`auto-fix`, confidence 7): no row plants an annotated tag object ID, and a `len(commit) == len(line)` mutant stayed silent. Repair: an annotated-tag case in `TestPendingMergeRefusesMalformedState`.
- C4 (`auto-fix`, confidence 6): CRLF, trailing space, and a double final newline refuse but have no pin; a `TrimSpace` mutant stayed silent. Repair: three malformed cases.
- C5 (`auto-fix`, confidence 5): a file-presence copy of the reader in `land_refusal.go` would pass MP15 to MP18. Repair: an MP15 subtest with a non-empty malformed MERGE_HEAD that expects false.

MP2 findings: 3. Worst issue: C1.

- C1 (`auto-fix`, confidence 5): a D1 check gated on one head stays silent, because no refusal row uses an octopus merge. Repair: an octopus refusal row in `TestLandRefusesAnIncompleteMergeNameSet`.
- C2 (`no-op`, confidence 8): the newline-framing mutation leaves MP42 green; the MP42 why cell is inexact, and P1 carries the spec fix.
- C3 (`auto-fix`, confidence 4): a recheck that compares sorted copies stays silent. Repair: a reordered-heads row in `TestLandRefusesAMergeStateChangeDuringAuthorization`.
- C21 (`auto-fix`, confidence 7, round 2): a covered listing without the `:(literal)` prefix stays silent. No row puts a glob-shaped merged name beside a sibling that the glob matches. Repair: a refusal row where the incoming commit adds `side-*` and `side-a`, the name set names only `side-*`, and the refusal names `side-a`.

  Closed in `cca1013b`. The orchestrator reran the prefix-drop swap on the clean frozen tip `f10aeedf`, and it bit the new row with the subject restored:

```text
tree[1]{target,head,dirty}:
  commit-merge-parent-build,f10aeedf2dc9151d660516d1c7a33f9eefdec7a7,true
probe[1]{verdict,subject,mutation,cause,failed_tests,restored}:
  bit,internal/landing/attribution.go,swap,failed,1,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/landing,^TestLandRefusesAnIncompleteMergeNameSet$,passed,8
spilled{lines=11,bytes=811,omitted_lines=2,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/f80fcedec92401bd6cf2acfce17ca107/1791703492507481259-8c4fb8a25dd700a8.out}
packages[1]{package,status,elapsed_ms,tests_run}:
  github.com/gibbonmi/bench/internal/landing,fail,349,8
failures[1]{package,test,line,lines}:
  github.com/gibbonmi/bench/internal/landing,TestLandRefusesAnIncompleteMergeNameSet/glob-shaped_name_covers_only_itself,"pending_merge_test.go:85: Land error = <nil>, want the named-path refusal \"pending merge changed paths that the commit does not name: \\\\\"side-a\\\\\"\"",1
skips[0]{package,test,reason}:
```

Advice kept for the final reconciliation: no positive `DryRun` row exists during a pending merge. A merged path that the caller unstages drops out of the index listing, and the Edge inventory now names that case. The `AuthorizeCandidate` read of the heads precedes `composeAuthorized`, and the MP3 review reads that window.
