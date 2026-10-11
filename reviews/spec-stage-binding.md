# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/spec-stage-binding/spec.md",
  "plan_digest": "sha256:e279e6de12a31db14190481ef8ba8a08db4095c51c8229dfa30516079660fb3b",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SB1",
      "base": "8c8a6ca52da262197b60a93236a5fd7ad8fcb38d",
      "tip": "fccb67ae75db0a0f462262fff5f599a4ee2dc9fa",
      "plan_digest": "sha256:e279e6de12a31db14190481ef8ba8a08db4095c51c8229dfa30516079660fb3b",
      "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
      "acceptance_rows": [
        "SB01",
        "SB02",
        "SB03",
        "SB04",
        "SB05",
        "SB06",
        "SB07",
        "SB08",
        "SB09",
        "SB10",
        "SB11",
        "SB12",
        "SB13",
        "SB14",
        "SB15",
        "SB16",
        "SB17",
        "SB18",
        "SB19",
        "SB20",
        "SB21",
        "SB22"
      ],
      "verification": [
        {
          "id": "t1-approval-focused-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-approval-focused.txt",
            "digest": "sha256:dfe4ac9aa1214fcffcf98b080e13bd20c4bb6cc316175860298f193d2eeb2910",
            "excerpt": "$ bench test --package ./internal/commitment/repository --run 'TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree'\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,121,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/commitment/repository/candidate.go --swap 'readTree(current, *candidate, source, tree, revision)' --with 'readTree(current, *candidate, source, revision, revision)' --package ./internal/commitment/repository --run 'TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree'\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/candidate.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree,passed,5\nspilled{lines=14,bytes=2063,omitted_lines=5,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/117ec43b4dcf3574d4c69d0aecedb569/1791730870017675952-9cc9879da2ab79ce.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAdmitPublicationGradesANewBindingAtTheComposedTree/approved,\"head_binding_test.go:103: AdmitPublication(the composed approved binding) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want admission\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAdmitPublicationGradesANewBindingAtTheComposedTree/changed-after-approval,\"head_binding_test.go:108: AdmitPublication(a composed changed spec) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want a refusal naming \\\"commitment source \\\\\\\\\\\"spec\\\\\\\\\\\" changed from sha256:c0724cf0547d5… (367 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAuthorizeCandidateAdmitsAHeadHeldBinding,\"head_binding_test.go:67: AuthorizeCandidate(the approved binding beside its spec) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want admission\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAuthorizeCandidateRefusesAChangedBoundSpec,\"head_binding_test.go:80: AuthorizeCandidate(a changed bound spec) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want a refusal naming \\\"commitment source \\\\\\\\\\\"spec\\\\\\\\\\\" changed from sha256:c0724cf0547d5f8… (365 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-approval-focused",
          "command": "bench test --package ./internal/commitment/repository --run 'TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree'",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the candidate tree for the default-branch revision in the new-deliverable branch of approvedTransition in internal/commitment/repository/candidate.go, so that the approval check reads a new deliverable at the default-branch tip. TestAuthorizeCandidateAdmitsAHeadHeldBinding and TestAdmitPublicationGradesANewBindingAtTheComposedTree turn red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ssb_t1/t1-approval-focused.txt",
              "digest": "sha256:dfe4ac9aa1214fcffcf98b080e13bd20c4bb6cc316175860298f193d2eeb2910",
              "excerpt": "$ bench test --package ./internal/commitment/repository --run 'TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree'\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,121,5\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n\n$ bench probe internal/commitment/repository/candidate.go --swap 'readTree(current, *candidate, source, tree, revision)' --with 'readTree(current, *candidate, source, revision, revision)' --package ./internal/commitment/repository --run 'TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree'\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,internal/commitment/repository/candidate.go,swap,failed,4,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/commitment/repository,TestAuthorizeCandidateAdmitsAHeadHeldBinding|TestAuthorizeCandidateRefusesAChangedBoundSpec|TestAdmitPublicationGradesANewBindingAtTheComposedTree,passed,5\nspilled{lines=14,bytes=2063,omitted_lines=5,cut_lines=0,path=/home/mgibs/.bench/responses/bench-2826441890/117ec43b4dcf3574d4c69d0aecedb569/1791730870017675952-9cc9879da2ab79ce.out}\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAdmitPublicationGradesANewBindingAtTheComposedTree/approved,\"head_binding_test.go:103: AdmitPublication(the composed approved binding) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want admission\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAdmitPublicationGradesANewBindingAtTheComposedTree/changed-after-approval,\"head_binding_test.go:108: AdmitPublication(a composed changed spec) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want a refusal naming \\\"commitment source \\\\\\\\\\\"spec\\\\\\\\\\\" changed from sha256:c0724cf0547d5… (367 bytes)\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAuthorizeCandidateAdmitsAHeadHeldBinding,\"head_binding_test.go:67: AuthorizeCandidate(the approved binding beside its spec) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want admission\",1\n  github.com/gibbonmi/bench/internal/commitment/repository,TestAuthorizeCandidateRefusesAChangedBoundSpec,\"head_binding_test.go:80: AuthorizeCandidate(a changed bound spec) = commitment source \\\"spec\\\" refused: missing or nonregular tree file specs/stage/spec.md, want a refusal naming \\\"commitment source \\\\\\\\\\\"spec\\\\\\\\\\\" changed from sha256:c0724cf0547d5f8… (365 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "t1-plan-focused-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-plan-focused.txt",
            "digest": "sha256:f65a5527872a711359f3f40f0cd9cb48ef349e4c5d92a3f28fa795e12d5c8220",
            "excerpt": "$ bench test --package ./internal/commitment --run 'TestHeadBinding'\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment,pass,381,11\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-plan-focused",
          "command": "bench test --package ./internal/commitment --run 'TestHeadBinding'",
          "exit_code": 0
        },
        {
          "id": "t1-journey-focused-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-journey-focused.txt",
            "digest": "sha256:8b1e375d58d3895e69dd9d8d9f8f0666cce748b7c313e13f8eb20b4f66adb365",
            "excerpt": "$ bench test --package ./internal/worktree --run 'TestCommitmentSpecStageLandsSpecAndBinding|TestPackageTestCountPin'\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/worktree,pass,903,2\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-journey-focused",
          "command": "bench test --package ./internal/worktree --run 'TestCommitmentSpecStageLandsSpecAndBinding|TestPackageTestCountPin'",
          "exit_code": 0
        },
        {
          "id": "t1-floor-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-floor.txt",
            "digest": "sha256:0accbb611f88341810dcb1b1339b828c617d14123f4ebfbb0cdae6bf73fe6b38",
            "excerpt": "$ bench test --package ./internal/commitment --package ./internal/commitment/repository --package ./internal/commitment/commitcmd --package ./internal/commitment/commitmenttest --package ./internal/commit --package ./internal/worktree --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\npackages[9]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,21901,333\n  github.com/gibbonmi/bench/internal/anchors,pass,2126,179\n  github.com/gibbonmi/bench/internal/commit,pass,12973,117\n  github.com/gibbonmi/bench/internal/commitment,pass,5320,172\n  github.com/gibbonmi/bench/internal/commitment/commitcmd,pass,788,19\n  github.com/gibbonmi/bench/internal/commitment/commitmenttest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,4924,104\n  github.com/gibbonmi/bench/internal/conformance,pass,73694,1685\n  github.com/gibbonmi/bench/internal/worktree,pass,68745,1399\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "t1-floor",
          "command": "bench test --package ./internal/commitment --package ./internal/commitment/repository --package ./internal/commitment/commitcmd --package ./internal/commitment/commitmenttest --package ./internal/commit --package ./internal/worktree --package ./internal/anchors --package ./cmd/bench --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "t1-root-conformance-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-root-conformance.txt",
            "digest": "sha256:0b8e9b4d57e4ff5d10a79fa47e1da8d3a0dd1d7a08edfabd6d7743d4f3bb59fd",
            "excerpt": "$ bench test --check root-conformance\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15623,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-root-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "t1-git-plumbing-owner-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-git-plumbing-owner.txt",
            "digest": "sha256:48e2fbf350ea1499e81f8071e27cce6e80104e65e8c39f4447de3b6b1bbcdfe7",
            "excerpt": "$ bench test --check git-plumbing-owner\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\ncheck[1]{name,kind,tests_run,subjects}:\n  git-plumbing-owner,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,96,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-git-plumbing-owner",
          "command": "bench test --check git-plumbing-owner",
          "exit_code": 0
        },
        {
          "id": "t1-ordinary-build-census-ssb_t1",
          "performer": "claude:ssb_t1",
          "role": "author-verification",
          "model": "opus",
          "effort": "xhigh",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ssb_t1/t1-ordinary-build-census.txt",
            "digest": "sha256:cc52a32af2bfb7619c67339280bc87f4ac5327c3241b645e5fbeea823371dedb",
            "excerpt": "$ bench test --check ordinary-build-census\ntree[1]{target,head,dirty}:\n  spec-stage-binding-build,fccb67ae75db0a0f462262fff5f599a4ee2dc9fa,false\ncheck[1]{name,kind,tests_run,subjects}:\n  ordinary-build-census,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,321,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "t1-ordinary-build-census",
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
  }
}
```
