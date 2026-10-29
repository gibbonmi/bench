# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/git-query-cache/spec.md",
  "plan_digest": "sha256:8efe8b8ecb27c8ac7599d34e7a97cb2080000afc20578bb0a68f72888bc5e43a",
  "implementation_session": "",
  "chunks": [
    {
      "id": "GQ-C1",
      "base": "f979bae6807f82e9bc4a259b78076c3170710597",
      "tip": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
      "plan_digest": "sha256:9ee3720fffdcc27a823ed7f88d3b68730848a2c596d6d31b43c28feb25a84638",
      "source_digest": "a0dcc1a8b07b80d1a167ccf311371857f81d72d2",
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
        },
        {
          "id": "t01r1-probe-01-hit",
          "performer": "claude:gq-t01-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a0dcc1a8b07b80d1a167ccf311371857f81d72d2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gq-t01-repair1-20261029",
            "digest": "sha256:a729b89ae05f6961fda3b870427330984c0fbdda5036be06d75a7e9f364e5704",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,9,TestRunObjectReads/absent_path,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReads,passed,14\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1791\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,78,14\nfailures[9]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/absent_path,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/ancestor,\"run_object_test.go:122: starts of \\\"--is-ancestor aa11d40c84a97b3fe77fcd450a78dc42b779ba4b\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/blob,\"run_object_test.go:122: starts of \\\"cat-file blob ce013625030ba8dba906f756967f9e9ca394464a\\\" = 4, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/commit,\"run_object_test.go:122: starts of \\\"rev-parse --verify --quiet aa11d40c84a97b3fe77fcd450a78dc42b779ba4b\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/flag_operand,\"run_object_test.go:122: starts of \\\"ls-tree -z -r 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/non-ancestor,\"run_object_test.go:122: starts of \\\"--is-ancestor d14ac56ea34f2ed64d264eb2ed5f42a4aba23ad3\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_file,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_path,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_suffix,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1^{tree}\\\" = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "requirement": "probe-01-hit",
          "command": "bench probe internal/git/run.go --swap 'if raw, ok := r.objects[key]; ok {' --with 'if raw, ok := r.objects[key]; false && ok {' --package ./internal/git --run TestRunObjectReads",
          "exit_code": 0,
          "probe": {
            "mutation": "Never answer an object read from the store. Each second read then starts git, so the object-read test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gq-t01-repair1-20261029",
              "digest": "sha256:a729b89ae05f6961fda3b870427330984c0fbdda5036be06d75a7e9f364e5704",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,9,TestRunObjectReads/absent_path,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReads,passed,14\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,1791\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,78,14\nfailures[9]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/absent_path,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/ancestor,\"run_object_test.go:122: starts of \\\"--is-ancestor aa11d40c84a97b3fe77fcd450a78dc42b779ba4b\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/blob,\"run_object_test.go:122: starts of \\\"cat-file blob ce013625030ba8dba906f756967f9e9ca394464a\\\" = 4, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/commit,\"run_object_test.go:122: starts of \\\"rev-parse --verify --quiet aa11d40c84a97b3fe77fcd450a78dc42b779ba4b\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/flag_operand,\"run_object_test.go:122: starts of \\\"ls-tree -z -r 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/non-ancestor,\"run_object_test.go:122: starts of \\\"--is-ancestor d14ac56ea34f2ed64d264eb2ed5f42a4aba23ad3\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_file,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_path,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1\\\" = 2, want 1\",1\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReads/tree_suffix,\"run_object_test.go:122: starts of \\\"ls-tree -z 2e81171448eb9f2ee3821e3d447aa6b2fe3ddba1^{tree}\\\" = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t01r1-probe-01-live",
          "performer": "claude:gq-t01-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "a0dcc1a8b07b80d1a167ccf311371857f81d72d2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:gq-t01-repair1-20261029",
            "digest": "sha256:c3bf202bf8f466dacf89f4202063bd4e1a65571f20bc9b51717f3f2aad38f072",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunObjectReadsStayLive,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReadsStayLive,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,33,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReadsStayLive,\"run_object_test.go:183: reads after the commit and the update-ref = \\\"hello\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630f2e0adb550\\\\\\\\ta.txt|5ee7e8923224fadbe38ce6633ee045d79c4435c8|true|false\\\", want \\\"next\\\\\\\\n|100644 blob 5b69858846f588e863c6d82f… (450 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "requirement": "probe-01-live",
          "command": "bench probe internal/git/run.go --swap 'if !FullObjectID(object) {' --with 'if false {' --package ./internal/git --run TestRunObjectReadsStayLive",
          "exit_code": 0,
          "probe": {
            "mutation": "Cache a blob read whose operand is a ref. The second read of the moved ref then returns the old bytes, so the live-read test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:gq-t01-repair1-20261029",
              "digest": "sha256:c3bf202bf8f466dacf89f4202063bd4e1a65571f20bc9b51717f3f2aad38f072",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/git/run.go,swap,failed,1,TestRunObjectReadsStayLive,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/git,TestRunObjectReadsStayLive,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,510\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/git,fail,33,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/git,TestRunObjectReadsStayLive,\"run_object_test.go:183: reads after the commit and the update-ref = \\\"hello\\\\\\\\n|100644 blob 5b69858846f588e863c6d82fbcf630f2e0adb550\\\\\\\\ta.txt|5ee7e8923224fadbe38ce6633ee045d79c4435c8|true|false\\\", want \\\"next\\\\\\\\n|100644 blob 5b69858846f588e863c6d82f… (450 bytes)\",1\nskips[0]{package,test,reason}:\n"
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
        },
        {
          "id": "gq-c1-standards-2",
          "performer": "claude:gq-c1-standards-confirm-20261029",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a0dcc1a8b07b80d1a167ccf311371857f81d72d2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c1-standards-confirm-20261029",
            "digest": "sha256:a9c1e1e4ce33b7c32425dd2f4208abdf30ee71d2fb568ed2301ebbf7134237ea",
            "excerpt": "axis: Standards (confirming round, repair delta 98fabb293..0cc386a93)\nS1: closed. The ADR names the three reads that check again and states that the top-level read does not; it holds no path and no code.\nThe corrected Run.AdminPath comment is accurate. No new finding.\n"
          },
          "axis": "Standards",
          "base": "f979bae6807f82e9bc4a259b78076c3170710597",
          "tip": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "finding_ids": [],
          "supersedes": [
            "gq-c1-standards-1"
          ]
        },
        {
          "id": "gq-c1-spec-2",
          "performer": "claude:gq-c1-spec-confirm-20261029",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a0dcc1a8b07b80d1a167ccf311371857f81d72d2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c1-spec-confirm-20261029",
            "digest": "sha256:3d39c9ffc23d416bc98629a54fa59c52a50762a276c63bee37c079b89723bad7",
            "excerpt": "axis: Spec (confirming round, repair delta 98fabb293..0cc386a93)\nP1: closed. The run.go change is comment-only and states the decided rule: every read on a non-nil run checks the parent; a nil run answers as the package form.\nNo new finding; the added cases do not contradict GQ22 or GQ23.\n"
          },
          "axis": "Spec",
          "base": "f979bae6807f82e9bc4a259b78076c3170710597",
          "tip": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "finding_ids": [],
          "supersedes": [
            "gq-c1-spec-1"
          ]
        },
        {
          "id": "gq-c1-coverage-2",
          "performer": "claude:gq-c1-coverage-confirm-20261029",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "a0dcc1a8b07b80d1a167ccf311371857f81d72d2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c1-coverage-confirm-20261029",
            "digest": "sha256:097e9682ab6692b1459abc8ebda4d1b1e75361e96b300f10b1dc3e0c3cca5058",
            "excerpt": "axis: Coverage (confirming round, repair delta 98fabb293..0cc386a93)\nC1 closed: a second AdminPath name and the root+\"/.\" spelling each force a start; the root+\"/.\" blob read gives two starts.\nC2 closed: IsAncestor(refs/heads/side, commit) moves from true to false after update-ref.\nC3 closed: the ^{tree} and -r cases cache; the commit-plus-path listing stays live.\nC4 closed: abbreviatedIdentity expands a 7-digit prefix to a 64-digit id.\nNo new finding.\n"
          },
          "axis": "Coverage",
          "base": "f979bae6807f82e9bc4a259b78076c3170710597",
          "tip": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "finding_ids": [],
          "supersedes": [
            "gq-c1-coverage-1"
          ]
        }
      ]
    },
    {
      "id": "GQ-C2A",
      "base": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
      "tip": "01d699aaa16c6fb51b3ad2e4ecc0f233c4f337e9",
      "plan_digest": "sha256:8efe8b8ecb27c8ac7599d34e7a97cb2080000afc20578bb0a68f72888bc5e43a",
      "source_digest": "8c88f856f6eaa8102be9eaff1dc2ab985bc04582",
      "acceptance_rows": [
        "GQ27",
        "GQ28",
        "GQ29",
        "GQ30",
        "GQ31",
        "GQ32",
        "GQ33"
      ],
      "verification": [
        {
          "id": "t02-probe-02",
          "performer": "claude:gq-t02-author-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "1ad021511d104377070112847e04c2af460ee150",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:gq-t02-author-20261029",
            "digest": "sha256:e9dd7b70a306f91ca0f8755538db66d5c58889b8d8beed76bb3885c6b5ca8e47",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,TestRepoReadsTheCommonDirOnce/run,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestRepoReadsTheCommonDirOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,308\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,11,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestRepoReadsTheCommonDirOnce/run,\"repo_test.go:33: common-directory starts = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2d4bc88f94616517c99110528d230229cda74665",
          "requirement": "probe-02",
          "command": "bench probe internal/intent/intent.go --swap 'common, err := r.Run.CommonDir(r.Root)' --with 'common, err := git.CommonDir(r.Root)' --package ./internal/intent --run TestRepoReadsTheCommonDirOnce",
          "exit_code": 0,
          "probe": {
            "mutation": "Read the ledger address without the run. Two ledger reads then start the common-directory query twice, so the ledger test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:gq-t02-author-20261029",
              "digest": "sha256:e9dd7b70a306f91ca0f8755538db66d5c58889b8d8beed76bb3885c6b5ca8e47",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,TestRepoReadsTheCommonDirOnce/run,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestRepoReadsTheCommonDirOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,308\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,11,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestRepoReadsTheCommonDirOnce/run,\"repo_test.go:33: common-directory starts = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t02r1-probe-02",
          "performer": "claude:gq-t02-repair1-20261029",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "8c88f856f6eaa8102be9eaff1dc2ab985bc04582",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-xhigh:gq-t02-repair1-20261029",
            "digest": "sha256:e9dd7b70a306f91ca0f8755538db66d5c58889b8d8beed76bb3885c6b5ca8e47",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,TestRepoReadsTheCommonDirOnce/run,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestRepoReadsTheCommonDirOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,308\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,11,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestRepoReadsTheCommonDirOnce/run,\"repo_test.go:33: common-directory starts = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "01d699aaa16c6fb51b3ad2e4ecc0f233c4f337e9",
          "requirement": "probe-02",
          "command": "bench probe internal/intent/intent.go --swap 'common, err := r.Run.CommonDir(r.Root)' --with 'common, err := git.CommonDir(r.Root)' --package ./internal/intent --run TestRepoReadsTheCommonDirOnce",
          "exit_code": 0,
          "probe": {
            "mutation": "Read the ledger address without the run. Two ledger reads then start the common-directory query twice, so the ledger test fails.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-xhigh:gq-t02-repair1-20261029",
              "digest": "sha256:e9dd7b70a306f91ca0f8755538db66d5c58889b8d8beed76bb3885c6b5ca8e47",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/intent/intent.go,swap,failed,1,TestRepoReadsTheCommonDirOnce/run,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/intent,TestRepoReadsTheCommonDirOnce,passed,3\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,308\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/intent,fail,11,3\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/intent,TestRepoReadsTheCommonDirOnce/run,\"repo_test.go:33: common-directory starts = 2, want 1\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "gq-c2a-standards-1",
          "performer": "claude:gq-c2a-standards-20261029",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "1ad021511d104377070112847e04c2af460ee150",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c2a-standards-20261029",
            "digest": "sha256:c92f6a0d4eda86851dc00c4b25e662b7557d48863e94c643acda17acbec6a3cf",
            "excerpt": "axis: Standards | verdict: findings | 1 | worst: S1\nS1 | ask-user | confidence 4 | internal/commitment/repository/run_test.go:74-86 | The malformed-listing subtest builds its own git shim beside gittest.StubGitDir; AGENTS.md names a pasted fixture harness as a defect. Fix: a gittest mode, after a Writes widen.\nAdvice: several new exported Repo methods carry no doc comment.\n"
          },
          "axis": "Standards",
          "base": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "tip": "2d4bc88f94616517c99110528d230229cda74665",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "gq-c2a-spec-1",
          "performer": "claude:gq-c2a-spec-20261029",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "1ad021511d104377070112847e04c2af460ee150",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c2a-spec-20261029",
            "digest": "sha256:9e61598cfe593181a8b01df203d622bf43e3c1cc25580cd9962407c8f5ef496e",
            "excerpt": "axis: Spec | verdict: findings | 2 | worst: P1\nP1 | ask-user | confidence 4 | internal/commitment/repository/readiness.go:111 | inheritance builds Store{Root: owner.Worktree} with no Run; the ticket requires store reads through the run.\nP2 | no-op | confidence 3 | internal/roadmap/tree_validation.go:60-91 | the roadmap run forms refuse a malformed listing record; spec line 157 orders the move onto Run.ListTree.\nGQ27-GQ33 each assert the row behavior.\n"
          },
          "axis": "Spec",
          "base": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "tip": "2d4bc88f94616517c99110528d230229cda74665",
          "finding_ids": [
            "P1",
            "P2"
          ],
          "supersedes": []
        },
        {
          "id": "gq-c2a-coverage-1",
          "performer": "claude:gq-c2a-coverage-20261029",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "1ad021511d104377070112847e04c2af460ee150",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:gq-c2a-coverage-20261029",
            "digest": "sha256:0d074fce33c0a8a823ec091e22ecb9ae73d8d173d46ea58aaa3089726011a94b",
            "excerpt": "axis: Coverage | verdict: findings | 6 | worst: C1\nC1 | auto-fix | 8 | internal/roadmap/tree_validation.go:78,85,99 | RevisionDocumentIn has no test; a nil run there survives.\nC2 | auto-fix | 7 | internal/commitment/repository/repository.go:29 | a hidden shared run in ledger() for a nil-run store survives; no nil-run ledger case.\nC3 | auto-fix | 7 | internal/intent/intent.go:355 | r.Run.Worktrees to git.Worktrees survives; only Read is tested through Repo.\nC4 | ask-user | 6 | internal/commitment/repository/lease.go:45,49 | a nil run in the lease survives in this chunk; only GQ34 in GQ-C2B can catch it.\nC5 | ask-user | 5 | internal/roadmap/tree_validation.go:90 | the --name-only listing stays raw; no test pins it.\nC6 | ask-user | 6 | reviews/git-query-cache.md | the record holds only the probe-02 witness; the GQ29-GQ33 reds are not recorded.\n"
          },
          "axis": "Coverage",
          "base": "0cc386a93fed6a55bbc47e31648557aa2f8532b0",
          "tip": "2d4bc88f94616517c99110528d230229cda74665",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4",
            "C5",
            "C6"
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
      "from": "sha256:7d3a811db66c80945a137d70bfb8b0c7995202a3a5bf521a5e9bca38fb700356",
      "to": "sha256:9ee3720fffdcc27a823ed7f88d3b68730848a2c596d6d31b43c28feb25a84638",
      "chunk_ids": {
        "GQ-C1": [
          "GQ-C1"
        ]
      }
    },
    {
      "from": "sha256:9ee3720fffdcc27a823ed7f88d3b68730848a2c596d6d31b43c28feb25a84638",
      "to": "sha256:e3db4c3db95e2710475c15f84fca211605a654f3490985a286b20ceebec9fdd1",
      "chunk_ids": {
        "GQ-C1": [
          "GQ-C1"
        ]
      }
    },
    {
      "from": "sha256:e3db4c3db95e2710475c15f84fca211605a654f3490985a286b20ceebec9fdd1",
      "to": "sha256:8efe8b8ecb27c8ac7599d34e7a97cb2080000afc20578bb0a68f72888bc5e43a",
      "chunk_ids": {
        "GQ-C1": [
          "GQ-C1"
        ],
        "GQ-C2A": [
          "GQ-C2A"
        ]
      }
    }
  ]
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

## Chunk GQ-C1 repair cycle 1

Repair cycles used: 1 of 2. The repair commit is 0cc386a93, and the refreeze record is e31006c59. The repair corrected the ADR sentence (S1) and the `Run.AdminPath` doc comment (P1), and it added the cases of C1 to C4.

The repair session reported these witness reds:

- C1: an `AdminPath` key without the name turns `TestRunAdminPathRevalidates` red.
- C1: an `AdminPath` key without the root spelling turns `TestRunAdminPathRevalidates` red.
- C1: an object key without the root spelling turns `TestRunObjectReads/root_spellings` red.
- C2: an `IsAncestor` guard without the ancestor term turns `TestRunObjectReadsStayLive` red.
- C3: a listing rule without the `^{tree}` trim turns `TestRunObjectReads/tree_suffix` red.
- C3: a listing rule without the flag skip turns `TestRunObjectReads/flag_operand` red.
- C4: a 40-digit rule in `abbreviatedIdentity` turns `TestReauthorizeArgAcceptsAFullObjectID` red.

The coordinator probe omitted the stop at the `--` operand in the listing rule. Four cases of `TestRunObjectReads` went red, and `tree_suffix` was one of them.

## Chunk GQ-C2A review pickup

The review read the frozen pair 0cc386a93 to 2d4bc88f9. Repair cycles used: 0 of 2. Four accepted findings fall on ticket 02, so one combined repair charge carries them.

### Standards

Findings: 1. Worst issue: S1.

- S1, auto-fix, confidence 4: the malformed-listing case builds its own git shim beside the shared `gittest` stub (internal/commitment/repository/run_test.go:74-86). The one-source rule names a pasted fixture harness as a defect. The repair adds a `gittest` mode, and the `--auto-approve` default of this run covers the widen.

### Spec

Findings: 2. Worst issue: P1.

- P1, auto-fix, confidence 4: `inheritance` builds an inner store without the run (internal/commitment/repository/readiness.go:111). The ticket requires each store read to go through its run.
- P2, no-op, confidence 3: the roadmap run forms refuse a malformed listing record. Spec line 157 orders the move onto the run's listing, so the behavior is part of the spec.

### Coverage

Findings: 6. Worst issue: C1.

- C1, auto-fix, confidence 8: no test counts the starts of `RevisionDocumentIn` (internal/roadmap/tree_validation.go:78). The repair adds that count. The store call sites stay with GQ34, which counts the whole landing.
- C2, auto-fix, confidence 7: no ledger case uses a store with no run (internal/commitment/repository/repository.go:29). So a hidden shared run survives, and the repair adds that case.
- C3, no-op, confidence 7: each other root form of the ledger is a one-line delegation, and GQ34 counts a lost cache. A test for each form is optional advice.
- C4, no-op, confidence 6: the spec puts the lease count in the GQ34 landing row of chunk GQ-C2B.
- C5, no-op, confidence 5: the spec's Won't-handle line keeps the `--name-only` listing live.
- C6, auto-fix, confidence 6: the record held only the probe-02 witness. The coordinator records the author's other witness reds below, and this change is evidence-only.

### Ticket 02 author witness reds

- GQ27: a ledger address without the run turns `TestRepoReadsTheCommonDirOnce/run` red.
- GQ28: a root `Read` through a shared run turns `TestRepoReadsTheCommonDirOnce/root` red.
- GQ29: a policy listing through the package form turns `TestStoreReadsAPolicyOnce` red.
- GQ30: a store ledger with a nil run turns `TestStoreReadsAPolicyOnce/ledger` red.
- GQ31: a policy listing through a hidden shared run turns `TestStoreReadsAPolicyOnce/no_run` red.
- GQ32: an empty listing read as present turns `TestStoreReadsAPolicyOnce/absent` red.
- Malformed record: a listing error read as absence turns `TestStoreReadsAPolicyOnce/malformed_listing` red.
- GQ33: a revision listing through the package form turns `TestRevisionIndexReadsOnce` red.

## Chunk GQ-C2A repair cycle 1

Repair cycles used: 1 of 2. The repair commit is 01d699aaa, and the refreeze record is 0b8a235bc. The repair added a `gittest` stub mode for a malformed listing (S1). It also passed the run to the inner store of `inheritance` (P1), and it added the cases of C1 and C2.

The repair session reported these witness reds:

- S1: a listing error read as absence turns `TestStoreReadsAPolicyOnce/malformed_listing` red through the new stub mode.
- C2: a hidden shared run for a store with no run turns `TestStoreReadsAPolicyOnce/ledger_no_run` red.
- C1: a nil run inside `RevisionDocumentIn` turns `TestRevisionIndexReadsOnce/document` red.
