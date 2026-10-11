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
      "reviews": [
        {
          "id": "sb1-standards-r1",
          "performer": "claude:ssb_sb1_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ssb_sb1_standards",
            "digest": "sha256:ca50d9f0902e7c8a51758df523995e1d7d70284425a9ce7d2b69a197223e99c9",
            "excerpt": "claude:ssb_sb1_standards (opus/high) — SB1 round 1 Standards axis over 8c8a6ca5..caa476a6 (chunk tip fccb67ae)\n\nFindings: 3. Worst: S1.\n\nS1 (ask-user, confidence 6): the independently authored expectations of the rows that pass before the edit have no red in the record. Rule: AGENTS.md, an independent expectation needs a recorded and demonstrated red; ticket lines 59-77 and 101 require the SB01 pre-edit red and one probe per row that passes before the edit (SB07, SB10, SB11, SB13 at HEAD, SB15 identity omission, SB16, SB19, SB20, SB21, SB22). reviews/spec-stage-binding.md names one mutation only, the central readTree swap (lines 53 and 66). The new expectations are the exact refusal strings at head_binding_test.go:125,161,197 and repository/head_binding_test.go:57. The author's claude-agent:ssb_t1/* native refs are outside the tree; a delegate return is not citable at the landing. Repair: record the ticket's probes and the SB01 pre-edit red in the review record, or have the reviewer accept the reduction.\n\nS2 (auto-fix, confidence 8): craft-comments with ste-prose, 25 words per sentence. internal/commitment/head_binding_test.go:44-45, the TestHeadBindingPlansAStagedSpecThatHEADHolds comment is 27 words; internal/commitment/repository/head_binding_test.go:22-24, the second sentence of the specStage comment is 29 words. Repair: split each sentence in two.\n\nS3 (auto-fix, confidence 4): craft-comments, a comment states the current code. internal/commitment/head_binding_test.go headBinding (line 39) says it \"approves the spec bytes body\"; the function only builds a DeliveryBinding. Repair: \"headBinding binds the spec bytes body at headSpec with no obligation.\"\n\nRefuted: no existing rev-parse HEAD helper, so checkoutHead is no duplicate; namesDeliverable and light_path.go approvedDeliverable are different predicates; Bound is single-sourced with three callers; the commitmenttest to repository import already existed (milestone.go:11); the store_test.go:147 edit is permitted by ticket line 124; the owner's proposed parameter is required by the spec term.\n\nAdvice: three loops share \"for each source: validate at readTree(...)\" (validateDeliverables, validateSources, approvedTransition line 182); one validateSourcesAt could own two, and Plan resolves sourceRevision twice. The \"changed from\" refusal format is rebuilt in two new test files (head_binding_test.go:161, repository/head_binding_test.go:57). A SourceBinding{ID, Path, Identity(body)} literal appears in four places; a commitmenttest.SpecBinding helper would fit beside TicketsBinding. checkoutHead and validateDeliverable comments mostly restate the code.\n\nExamined: check-current current=true at caa476a6; git diff --stat and git diff -- internal, the spec diff, git log (three spill files); AGENTS.md, .bench/BENCH.md, projects/benchkit.md, craft-review SKILL and finding-discipline, craft-comments; the ticket and spec lines 117-191; the review record in full with a search for every bench probe call; sources.go, candidate.go:160-254, repository.go:25-224, delivery.go, authority.go:105-149, light_path.go:165-189, model.go:50-74, commitmenttest/repo.go:60-236, the new test files through the diff. No native-ref files fetched; no probe or test run. Tokens 131k, 41 tool uses.\n"
          },
          "axis": "Standards",
          "base": "8c8a6ca52da262197b60a93236a5fd7ad8fcb38d",
          "tip": "fccb67ae75db0a0f462262fff5f599a4ee2dc9fa",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "sb1-spec-r1",
          "performer": "claude:ssb_sb1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ssb_sb1_spec",
            "digest": "sha256:5a3f8f3da2c791de4136b68af97377cbe15f20c8ed4d08be9db377cd9b31e98b",
            "excerpt": "claude:ssb_sb1_spec (opus/high) — SB1 round 1 Spec axis over 8c8a6ca5..caa476a6 (chunk tip fccb67ae)\n\nFindings: 1. Worst: P1.\n\nP1 (auto-fix, confidence 6): the frozen tree keeps only the central probe's red. Spec line: \"The build first observes SB01 red at `Store.Plan`, before the first production edit. Each independent expectation then records one red through a named omission or a behavioral swap.\" (Validation plan). reviews/spec-stage-binding.md t1-approval-focused-ssb_t1 keeps the central probe only; commits fccb67ae and caa476a6 carry no red log for the SB01 pre-edit red, SB07, SB10/SB11, the SB13 HEAD swap, the SB13/SB15 identity omission, SB16, SB19, SB20, SB21, SB22. The author's native refs claude-agent:ssb_t1/*.txt are outside the tree and were not read. Repair: rerun those probes and keep them in the review record.\n\nAuthor gaps: gap 1 (owner takes the proposed policy) no-op, confidence 9: the Terms define a new deliverable source over both policies, and readTree and NewDeliverable follow the spec; the ticket line 22 is the non-behavioral omission. Gap 2 (approve and show in-process) no-op, confidence 8: cmd/bench/commitment_command.go:12-14 routes bench commitment to commitcmd.Command(boundaryRoot(), args), the same handler. Gap 3 (SB11 red before the edit) no-op, confidence 8 provided P1 is repaired: the all-sources-at-HEAD mutation still bites SB11 through B.spec at specs/main/spec.md.\n\nPredicates confirmed: E1/SB08 refuses with no tip fallback; every other source reads revision; validateDeliverable keeps sourceRevision for its two callers admission.go:41 and readiness.go:60; inventory roadmap rows read revision; refusal sentences unchanged at sources.go:105,109,160 and candidate.go:189; commit.go, commitcmd, help rows untouched; --request <request> stays at commitcmd/command.go:31.\n\nRows: SB01-SB22 all met (SB10, SB11, SB13, SB16, SB19 with the probe red not kept, P1); SB16, SB20, SB22 pins unchanged.\n\nAdvice: journey step 3 binds Outcomes[0] through commitmenttest.Bound, not \"the outcome that bench commitment show names\" (probably equal in the fixture, unchecked). Ticket line 62 could add SB11 to the rows that fail before the edit.\n\nExamined: check-current current=true at caa476a6; one git diff 8c8a6ca5..caa476a6 (four spill files) and git log; the whole spec, ticket 1, the review record, both head_binding_test.go files, sources.go, candidate.go:25-190, repository.go:85-234, authority.go:48-102, land_specless_test.go:205-256, cmd/bench/commitment_command.go; callers of every changed function and the three pinned tests. Tokens 129k, 35 tool uses.\n"
          },
          "axis": "Spec",
          "base": "8c8a6ca52da262197b60a93236a5fd7ad8fcb38d",
          "tip": "fccb67ae75db0a0f462262fff5f599a4ee2dc9fa",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "sb1-coverage-r1",
          "performer": "claude:ssb_sb1_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "bbc0d279dc095adf4e8c293e919dc5365f418eb8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ssb_sb1_coverage",
            "digest": "sha256:fa04f8353674faf2f9432408772e92209d031f084c3c65307656b3bfdb5ad760",
            "excerpt": "claude:ssb_sb1_coverage (opus/high) — SB1 round 1 Coverage axis over 8c8a6ca5..caa476a6 (chunk tip fccb67ae)\n\nFindings: 3. Worst: C1.\n\nC1 (ask-user, confidence 8): state: the spec-stage worktree's base is behind main because main landed a production commit during the spec stage, and the worktree then commits its approved binding. Break: bench commit refuses with `assignment has no current delivery binding`; when the policy changes, branchReference returns the default-branch tip, so the diff from the tip to the candidate shows main's production change in reverse, which reaches readyFor and then lightPath. The Edge inventory decides only a moved policy; the close order has no fold step; SB12 and SB17 run only against an unmoved main. Needed: a decision (a fold step before the binding commit, or a grading change) plus an SB12 variant with landProduction before the commit. Read: candidate.go (branchReference, authorizeCandidate), light_path.go (lightPath). Reproduced by probe 5.\n\nC2 (auto-fix, confidence 6): state: main holds no policy (nil current) and a worktree plans the first policy that binds a spec that HEAD holds. Break: with the nil clause of NewDeliverable dropped, each source reads main and the plan refuses; no test notices (probes 1-3 silent). Needed: a bootstrap row in internal/commitment/head_binding_test.go from a worktree whose HEAD holds the bound spec. Read: delivery.go NewDeliverable.\n\nC3 (auto-fix, confidence 5): state: main's policy binds path P at identity I1; the worktree re-stages P at I2 and plans a binding at I2. Break: a path-only match reads main and refuses with \"changed from I2 to I1\"; no test notices (probe 4 silent). Needed: a re-binding row in TestHeadBindingKeepsTheDefaultTipForEveryOtherSource or a sibling. Read: delivery.go namesDeliverable.\n\nRefuted: Approve replay after the landing; a detached HEAD; a delivered binding (Completions refuses a proposal that changes a delivery fact); the competing refusals at plan, commit, and landing (SB07-SB09, SB13, SB15).\n\nProbes (all through bench worktree exec spec-stage-binding-build -- env GOFLAGS=\"-p=4 -parallel=4\"):\n1. bench probe internal/commitment/delivery.go --swap '(current == nil || !namesDeliverable(*current, source))' --with '(current != nil && !namesDeliverable(*current, source))' --package ./internal/commitment/... : silent, 0, restored yes\n2. probe 1 on ./internal/worktree --run 'Commitment|Planning|Bootstrap': silent, 0, restored yes\n3. probe 1 on ./internal/commit: silent, 0, restored yes\n4. bench probe internal/commitment/delivery.go --swap 'return binding.Source == source' --with 'return binding.Source.Path == source.Path' --package ./internal/commitment/... : silent, 0, restored yes (not run on worktree or commit, a sample)\n5. C1 refutation, a fixture change in internal/commitment/repository/head_binding_test.go that adds landProduction(t, root) before the binding commit, --run TestAuthorizeCandidateAdmitsAHeadHeldBinding: bit, 1 (`assignment has no current delivery binding`), restored yes\n\nAdvice: no row covers a tickets-only binding at the candidate or composed tree (story 7 covers Plan only); namesDeliverable counts settled bindings while validateDeliverables skips them, with no break possible today.\n\nExamined: input family (current nil-able and proposed policies with Unsettled and policySources; checkoutHead and sourceRevision; spec bytes at HEAD, candidate, and composed trees through sourceAt; ticketsOnlyAt; roadmap rows; ledger receipts); write set is the ticket 1 Writes line and the 10 code files are inside it; call-site postures of readTree, validateDeliverableAt, checkoutHead; the spec, ticket 1, the delta for internal/commitment and internal/worktree, sources.go, candidate.go, repository.go:90-224, delivery.go:100-190, authority.go:250-264, verification.go:187-204, readiness.go:40-70, light_path.go:27-57, stale_base_test.go:1-173, land_specless_test.go:235-280, spec/tickets_only.go; cursors: the identity binding and the review record lines 53 and 100-190. git status clean at the end.\n"
          },
          "axis": "Coverage",
          "base": "8c8a6ca52da262197b60a93236a5fd7ad8fcb38d",
          "tip": "fccb67ae75db0a0f462262fff5f599a4ee2dc9fa",
          "finding_ids": [
            "C1",
            "C2",
            "C3"
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
  }
}
```

## Repair state

Chunk SB1 round 1 ran three opus/high axes in parallel at `caa476a6` over the frozen pair `8c8a6ca5..caa476a6` (chunk tip `fccb67ae`). Only the Coverage axis ran probes, and the tree was clean after every return. SB1 uses repair cycle 1 of 2 through one fresh repair session for ticket 1. That session takes S2, S3, C2, C3, and the C1 row SB45. C2 and C3 are the chunk's one hardening cycle. S1 and P1 are one record gap with no code change, so the orchestrator closes them below as evidence-only work.

The reviewer is away and defers each question to an opus/high consultant. Under `--auto-approve`, the reviewer approved every repair round, fence expansion, and spec or ticket expansion in advance.

## Standards

Findings: 3. Worst issue: S1.

- S1 (`ask-user`, confidence 6): the independent expectations of the rows that pass before the edit had no red in the record. The record held only the central probe. Rule: AGENTS.md, an independent expectation needs a recorded and demonstrated red. Closed evidence-only: the orchestrator reran the author probes 2 to 10 on the clean tip `caa476a6`. The section "Retained red evidence" below keeps them beside the author's pre-edit reds. No reduction was accepted.
- S2 (`auto-fix`, confidence 8): two test comments hold a sentence over the 25-word bound. They are `internal/commitment/head_binding_test.go:44-45` (27 words) and `internal/commitment/repository/head_binding_test.go:22-24` (29 words). Rule: `craft-comments` with ste-prose. Repair: split each sentence in two.
- S3 (`auto-fix`, confidence 4): `headBinding` in `internal/commitment/head_binding_test.go:39` says it "approves the spec bytes body". The function builds a `DeliveryBinding`. Rule: `craft-comments`, a comment states the current code. Repair: "headBinding binds the spec bytes body at headSpec with no obligation."

Advice: one `validateSourcesAt` could own two of the three "for each source: validate at readTree" loops. A `commitmenttest` helper could own the "changed from" expectation and the `SourceBinding{ID, Path, Identity(body)}` literal.

## Spec

Findings: 1. Worst issue: P1.

- P1 (`auto-fix`, confidence 6): the frozen tree keeps only the central probe's red. The Validation plan says "Each independent expectation then records one red through a named omission or a behavioral swap." Closed evidence-only with S1: the reruns below are source-bound to `caa476a6`.
- The author's three flags are `no-op`. The owner takes the proposed policy because the Terms define a new deliverable source over both policies, and ticket line 22 is the non-behavioral omission. The journey calls the same `commitcmd.Command` handler that `cmd/bench/commitment_command.go:12-14` routes. SB11 was red before the edit, and its probe red (probe 4) stays its witness.

Advice: journey step 3 binds `Outcomes[0]` through `commitmenttest.Bound`. Ticket line 62 could add SB11 to its pre-edit list.

## Coverage

Findings: 3. Worst issue: C1.

- C1 (`ask-user`, confidence 8): when `main` lands a production commit during the spec stage, the planning worktree's binding commit refuses. The refusal is `assignment has no current delivery binding`. For a policy change, `branchReference` (`internal/commitment/repository/candidate.go`) returns the default-branch tip, so the reversed production diff reaches `readyFor` and `lightPath`. The Coverage probe 5 reproduced it. Consultant `claude:ssb_consult_c1` (opus/high) verified the cause and decided route A. The orchestrator commits the spec and ticket amendment as a plan commit.
- Route A: the close order folds `main` with `bench worktree merge --from main <target>` before the plan step. The Edge inventory marks a `main` move between the fold and the binding commit as Won't handle, because a second fold clears it. The SB1 repair proves the fold route through the new row SB45.
- C2 (`auto-fix`, confidence 6): with the nil clause of `NewDeliverable` dropped, a bootstrap plan from a worktree whose HEAD holds the bound spec refuses. No test notices (Coverage probes 1 to 3 silent). Repair: a bootstrap row in `internal/commitment/head_binding_test.go`.
- C3 (`auto-fix`, confidence 5): with a path-only match in `namesDeliverable`, a re-binding of a path that `main` binds at another identity refuses with "changed from". No test notices (Coverage probe 4 silent). Repair: a re-binding row beside `TestHeadBindingKeepsTheDefaultTipForEveryOtherSource`.

Advice: no row covers a tickets-only binding at the candidate or composed tree. `namesDeliverable` counts settled bindings while `validateDeliverables` skips them.

## Retained red evidence

The author observed these reds before the first production edit, on source `2a503764` plus the new tests. The logs are `claude-agent:ssb_t1/pre-edit-commitment.txt` and `pre-edit-journey.txt`.

- SB01: `head_binding_test.go:53: Plan(a binding of the spec that HEAD holds) = commitment source "A.spec" refused: ... want acceptance`.
- SB03: `head_binding_test.go:83: worktree inventory = [{Kind:run ...}], want one deliverable row`.
- SB06: `head_binding_test.go:104: ... is not an approved spec or tickets-only folder`.
- SB08: `head_binding_test.go:135: main-only: Plan = <nil>, want a refusal`.
- SB09: `head_binding_test.go:135: not-staged: ... want ... "deliverable \"specs/head/spec.md\" is not staged"`.
- SB11: `head_binding_test.go:176: retained-binding: ... refused ... want acceptance`.
- SB12 to SB15 and SB21: red at Plan inside their fixtures.
- SB17 and SB18: `land_specless_test.go:281: plan of the HEAD binding = commitment source "spec" refused`.

The orchestrator reran the author's probes 2 to 10 on the clean tip `caa476a6`. Each verdict is `bit` with `restored=yes`. The bodies are `claude-agent:ssb-orchestrator-20261011/sb1-review/rerun-probe-<n>.txt`.

| # | command | failed | failure line |
|---|---|---|---|
| 2 | `bench probe internal/commitment/repository/candidate.go --swap 'readTree(current, *candidate, source, tree, revision)' --with 'readTree(current, *candidate, source, "HEAD", revision)' --package ./internal/commitment/repository --run 'TestAuthorizeCandidateRefusesAChangedBoundSpec'` | 1 | `head_binding_test.go:80: AuthorizeCandidate(a changed bound spec) = <nil>, want a refusal naming "commitment source \"spec\" changed from ...` (SB13) |
| 3 | `bench probe internal/commitment/repository/candidate.go --swap 'if err := store.validateSourceAt(readTree(current, *candidate, source, tree, revision), source); err != nil {' --with '... err != nil && !(commitment.NewDeliverable(current, *candidate, source) && strings.Contains(err.Error(), " changed from ")) {' --package ./internal/commitment/repository --run 'TestAuthorizeCandidateRefusesAChangedBoundSpec\|TestAdmitPublicationGradesANewBindingAtTheComposedTree'` | 2 | `head_binding_test.go:108: AdmitPublication(a composed changed spec) = <nil>` and `:80` (SB13, SB15) |
| 4 | `bench probe internal/commitment/repository/sources.go --swap 'if commitment.NewDeliverable(current, proposed, source) {' --with 'if true \|\| commitment.NewDeliverable(current, proposed, source) {' --package ./internal/commitment --run 'TestHeadBindingKeepsTheDefaultTipForEveryOtherSource'` | 2 | `:163: ... commitment source "FT1" changed from sha256:6cb8... to sha256:adab...` and `:176: ... "B.spec" refused ... want acceptance` (SB10, SB11) |
| 5 | `bench probe internal/commitment/repository/sources.go --swap 'head, err := git.Output("-C", store.Root, "rev-parse", "HEAD")' --with '_, _ = git.Output("-C", store.Root, "add", "-A"); head, err := git.Output("-C", store.Root, "write-tree")' --package ./internal/commitment --run 'TestHeadBindingRefusesAnUnheldDeliverable'` | 1 | `:135: working-tree-only Plan = <nil>, want a refusal` (SB07) |
| 6 | `bench probe internal/commitment/repository/candidate.go --swap 'if receipt.Plan != expected.ID \|\| ... {' --with 'if false && (...) {' --package ./internal/commit --run 'TestCommitmentProposalIdentity'` | 1 | `commitment_test.go:125: commit = 0 committed{...path_count=2}` (SB16) |
| 7 | `bench probe internal/commitment/repository/sources.go --swap 'return store.validateDeliverableAt(revision, source)' --with 'return fmt.Errorf("deliverable %q refused at %s", source.Path, revision)' --package ./internal/commitment --run 'TestCommitmentDeliverableTypes'` | 3 | `deliverable_test.go:53: ... start refused — deliverable "specs/light/spec.md" refused at ...` (SB20) |
| 8 | `bench probe internal/commitment/repository/admission.go --swap 'selected, err := activeOutcome(policy, outcome)' --with 'if head, _, headErr := store.policyAt("HEAD"); headErr == nil && head != nil { policy = *head }; selected, err := activeOutcome(policy, outcome)' --package ./internal/commitment --run 'TestHeadBindingStartRefusesBeforeTheLanding'` | 1 | `:199: Start(the HEAD-bound spec) = commitment source "A.spec" refused ..., want ... is not approved for outcome "A"` (SB21) |
| 9 | `bench probe internal/commitment/commitcmd/command.go --swap '"--deliverable": outlook.Deliverable} {' --with '"--deliverable": outlook.Deliverable, "--request": "fixed"} {' --package ./cmd/bench --run 'TestCommitmentNextProjection'` | 1 | `commitment_test.go:75: ... --request fixed ..., want ... --request <request> ...` (SB22) |
| 10 | `bench probe internal/commitment/repository/outlook.go --swap 'policy, exists, err := store.Policy()' --with 'head, _, err := store.policyAt("HEAD"); policy, exists := commitment.Policy{}, head != nil; if head != nil { policy = *head }' --package ./internal/worktree --run 'TestCommitmentSpecStageLandsSpecAndBinding'` | 1 | `land_specless_test.go:291: show in .../worktrees/... before the landing = (...)` (SB19) |

The author's probes 1 (central, rerun at `fccb67ae`), 11, 12, 13, and 14 ran at the committed tip. Its verification entries and logs hold them. The orchestrator's independent probe ran at `caa476a6`: `bench probe internal/commitment/delivery.go --omit 'return true' --package ./internal/commitment --run '^TestHeadBindingPlansAStagedSpecThatHEADHolds$'`. It bit with 1 failed test (`head_binding_test.go:53 ... refused ..., want acceptance`) and `restored=yes`.
