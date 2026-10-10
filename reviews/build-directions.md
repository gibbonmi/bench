# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/build-directions/spec.md",
  "plan_digest": "sha256:8c0367a9d0b21d1603f3e78be5d804444dbff9d59ad19095af13ddfb3a1b744d",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BD-C1",
      "base": "dcf7fad065c259c7013bcd21138431360a15e33f",
      "tip": "d3bd78a98f4834886c8714d9bf00930e89c3ef23",
      "plan_digest": "sha256:8c0367a9d0b21d1603f3e78be5d804444dbff9d59ad19095af13ddfb3a1b744d",
      "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
      "acceptance_rows": [
        "BD1",
        "BD2",
        "BD3",
        "BD4",
        "BD5",
        "BD6",
        "BD60",
        "BD61",
        "BD7",
        "BD8",
        "BD9",
        "BD10",
        "BD11",
        "BD12",
        "BD13",
        "BD57",
        "BD14",
        "BD15",
        "BD16",
        "BD17",
        "BD18",
        "BD19",
        "BD20",
        "BD21",
        "BD50",
        "BD51",
        "BD52",
        "BD55",
        "BD58"
      ],
      "verification": [
        {
          "id": "bd-01-focused-t01",
          "performer": "claude:bd_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t01",
            "digest": "sha256:f157b43568b676159d859a14ac65483a3b347e7d215e54b36a8259d95fbf2b45",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,384,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-01-focused",
          "command": "bench test --package ./internal/anchors --run TestChunkChainAnchors",
          "exit_code": 0
        },
        {
          "id": "bd-01-live-t01",
          "performer": "claude:bd_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t01",
            "digest": "sha256:5e8129df6bbda5f7fe46629729edb599eedb0cda13cff8000d67e95093a93fb5",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,.agents/commands/bench-implement-spec.md,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,420\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,6716,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: chunk chain: no main fold between two tickets of a chunk\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-01-live",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the sentence \"It never folds `main` between two tickets of a chunk.\" from .agents/commands/bench-implement-spec.md. The docs-currency-workflow check turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bd_t01",
              "digest": "sha256:5e8129df6bbda5f7fe46629729edb599eedb0cda13cff8000d67e95093a93fb5",
              "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,.agents/commands/bench-implement-spec.md,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,420\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,6716,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: chunk chain: no main fold between two tickets of a chunk\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "bd-01-floor-t01",
          "performer": "claude:bd_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t01",
            "digest": "sha256:0b9e6fa2c484e845379065c3afdc08a11903ad1ff6b100a2e21248aba5ad2fe2",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17566,333\n  github.com/gibbonmi/bench/internal/anchors,pass,1920,179\n  github.com/gibbonmi/bench/internal/conformance,pass,63694,1684\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "bd-01-floor",
          "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "bd-01-conformance-t01",
          "performer": "claude:bd_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t01",
            "digest": "sha256:54ee48231e0d2207a2af0e301f279eed65cea347d46b4369f6acd8900ab92849",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18551,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-01-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "bd-02-focused-t02",
          "performer": "claude:bd_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t02",
            "digest": "sha256:02c88f60c29af72cbb56eefb6e4b97f9b6d065cedc6707aaeb07635fb0072d43",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,.bench/BENCH.md,omit,failed,19,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommitmentGuidance,passed,49\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,4692\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3108,49\nfailures[19]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/active-outcome_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/binding-only_limit,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/blocker_report,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/build_rule_pointer,\"commitment_guidance_test.go:95: .bench/BENCH.md holds \\\" `/bench-implement-spec` states the build rule.\\\" 0 times, want once\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/canonical_rule,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/drain_light-path_dispatch,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/explicit_displacement,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/in-build_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_Writes_boundary,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_commitment_start,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_observable,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_row_rule,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/mid-work_defect_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/next-outcome_start,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/outside-a-build_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/purpose_priority,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/start_route,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/unadmitted_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/uncommitted_intake,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-02-focused",
          "command": "bench test --package ./internal/conformance --run TestCommitmentGuidance",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the clause \"`/bench-implement-spec` states the build rule.\" from .bench/BENCH.md. TestCommitmentGuidance turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bd_t02",
              "digest": "sha256:02c88f60c29af72cbb56eefb6e4b97f9b6d065cedc6707aaeb07635fb0072d43",
              "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,.bench/BENCH.md,omit,failed,19,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,TestCommitmentGuidance,passed,49\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,4692\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,3108,49\nfailures[19]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/active-outcome_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/binding-only_limit,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/blocker_report,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/build_rule_pointer,\"commitment_guidance_test.go:95: .bench/BENCH.md holds \\\" `/bench-implement-spec` states the build rule.\\\" 0 times, want once\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/canonical_rule,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/drain_light-path_dispatch,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/explicit_displacement,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/in-build_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_Writes_boundary,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_commitment_start,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_observable,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/light-path_row_rule,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/mid-work_defect_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/next-outcome_start,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/outside-a-build_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/purpose_priority,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/start_route,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/unadmitted_learning_fix,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\n  github.com/gibbonmi/bench/internal/conformance,TestCommitmentGuidance/uncommitted_intake,\"commitment_guidance_test.go:99: live .bench/BENCH.md raised [\\\"commitment guidance: operating guide dropped the pointer to the build rule\\\"]\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "bd-02-floor-t02",
          "performer": "claude:bd_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t02",
            "digest": "sha256:e14ab074ffedd17297a5d89e99140a531df80818406d49e590dfd47085b231e6",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18100,333\n  github.com/gibbonmi/bench/internal/anchors,pass,1896,179\n  github.com/gibbonmi/bench/internal/conformance,pass,63555,1684\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "bd-02-floor",
          "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "bd-02-conformance-t02",
          "performer": "claude:bd_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t02",
            "digest": "sha256:b88cc99971ff81ad447cf64f9892575a7627f4797af2775b5e33a7013d3c93b6",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,19192,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-02-conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "bd-03-focused-t03",
          "performer": "claude:bd_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t03",
            "digest": "sha256:37b94a975779038220cd7526a2cb227257116ba50e56aaa4e1db2598284247c6",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,346,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-03-focused",
          "command": "bench test --package ./internal/anchors --run TestChunkChainAnchors",
          "exit_code": 0
        },
        {
          "id": "bd-03-live-t03",
          "performer": "claude:bd_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t03",
            "digest": "sha256:cd02a87db949b0bd99feaf3ee15efff1abd0fc00a5368ab994b957e6cbfa7c34",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,.agents/commands/bench-implement-spec.md,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,441\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5846,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: chunk chain: the last chunk keeps its chunk, verification, and review entries\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-03-live",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0,
          "probe": {
            "mutation": "Omit the sentence \"The last chunk keeps its chunk entry, its verification entries, and its three review entries.\" from .agents/commands/bench-implement-spec.md. The docs-currency-workflow check turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bd_t03",
              "digest": "sha256:cd02a87db949b0bd99feaf3ee15efff1abd0fc00a5368ab994b957e6cbfa7c34",
              "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\nprobe[1]{verdict,subject,mutation,cause,failed_tests,restored}:\n  bit,.agents/commands/bench-implement-spec.md,omit,failed,1,yes\nselection[1]{form,target,run,baseline,ran}:\n  check,docs-currency-workflow,^TestRootConformance$,passed,1\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v1,0,1,441\ncheck[1]{name,kind,tests_run,subjects}:\n  docs-currency-workflow,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,5846,1\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestRootConformance,\"gate_entry_test.go:29: gate: chunk chain: the last chunk keeps its chunk, verification, and review entries\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "bd-03-floor-t03",
          "performer": "claude:bd_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t03",
            "digest": "sha256:16438fcf0bb323c1ade2e80b9f1f61e203ff5f1ac3e9a28cb2112e663d175ba9",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\npackages[3]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,20078,333\n  github.com/gibbonmi/bench/internal/anchors,pass,2521,179\n  github.com/gibbonmi/bench/internal/conformance,pass,69366,1684\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "bd-03-floor",
          "command": "bench test --package ./internal/anchors --package ./internal/conformance --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "bd-03-conformance-t03",
          "performer": "claude:bd_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "abbd803bc52d56e9be73ca6b490ec45bd0c9cd03",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bd_t03",
            "digest": "sha256:26e3f8acb3f09a3c43a74ca28c03c19465dae932b2a3851fcafc4caa8b9880f8",
            "excerpt": "tree[1]{target,head,dirty}:\n  build-directions,d3bd78a98f4834886c8714d9bf00930e89c3ef23,true\ncheck[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18148,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "bd-03-conformance",
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
