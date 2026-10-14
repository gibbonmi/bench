# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-ownership-closure/spec.md",
  "plan_digest": "sha256:9e4941f4243f14167ae9a17c2f09edf8d5b4ae67633d4f0d597fa7017edfb0cb",
  "implementation_session": "",
  "chunks": [
    {
      "id": "OC-C1",
      "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
      "tip": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
      "plan_digest": "sha256:9e4941f4243f14167ae9a17c2f09edf8d5b4ae67633d4f0d597fa7017edfb0cb",
      "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
      "acceptance_rows": [
        "OC1",
        "OC2",
        "OC3",
        "OC4",
        "OC5",
        "OC6",
        "OC7",
        "OC8",
        "OC9",
        "OC10",
        "OC11",
        "OC12",
        "OC13",
        "OC55",
        "OC56",
        "OC57",
        "OC58",
        "OC59",
        "OC60",
        "OC71",
        "OC72",
        "OC73",
        "OC74",
        "OC75",
        "OC76",
        "OC77",
        "OC78",
        "OC79",
        "OC80",
        "OC81",
        "OC83",
        "OC84",
        "OC85",
        "OC90",
        "OC91",
        "OC92"
      ],
      "verification": [
        {
          "id": "ft293_t01-oc-01-focused",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:d8cf4e70353d81a274ac0d10e5cd07c0ea5019d9b281c6ab35b0836562b4907a",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestProposalFixedPoint/OC3_second-hop_pin,29,1,0,null,yes\n  2,bit,swap,failed,1,TestProposalProvenance/OC5_two_reasons,29,1,0,null,yes\n  3,bit,swap,failed,1,TestProposalProvenance/OC8_registry_owner,29,1,0,null,yes\n  4,bit,swap,failed,6,TestProposalFixedPoint/OC1_missing_pin,29,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,5043\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1541,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:109: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:114: writes_proposal rows matching map[path:internal/anchors/cycle.go] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:133: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:140: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1484,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:133: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",47\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1545,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC8_registry_owner,\"proposal_test.go:156: registry reason = map[kind:registry line:67 path:internal/toon/toon_test.go reason:binding internal/toon requires internal/toon/toon_test.go beside internal/toon/toon.go at kit source 8cd3e918cfedae2f74e5f387587b2ca2d4… (369 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1670,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC1_missing_pin,\"proposal_test.go:106: writes_proposal rows matching map[path:tests/canary/fam/pin-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:109: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:114: writes_proposal rows matching map[path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:133: closure_evidence rows matching map[kind:fixture line:1 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:140: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC7_included_BASE_source,\"proposal_test.go:146: closure_evidence rows matching map[kind:fixture path:tests/canary/fam/include-fixture] = 0, want 1:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t01",
              "digest": "sha256:d8cf4e70353d81a274ac0d10e5cd07c0ea5019d9b281c6ab35b0836562b4907a",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestProposalFixedPoint/OC3_second-hop_pin,29,1,0,null,yes\n  2,bit,swap,failed,1,TestProposalProvenance/OC5_two_reasons,29,1,0,null,yes\n  3,bit,swap,failed,1,TestProposalProvenance/OC8_registry_owner,29,1,0,null,yes\n  4,bit,swap,failed,6,TestProposalFixedPoint/OC1_missing_pin,29,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,5043\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1541,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:109: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:114: writes_proposal rows matching map[path:internal/anchors/cycle.go] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:133: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:140: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1484,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:133: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",47\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1545,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC8_registry_owner,\"proposal_test.go:156: registry reason = map[kind:registry line:67 path:internal/toon/toon_test.go reason:binding internal/toon requires internal/toon/toon_test.go beside internal/toon/toon.go at kit source 8cd3e918cfedae2f74e5f387587b2ca2d4… (369 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,1670,29\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[6]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC1_missing_pin,\"proposal_test.go:106: writes_proposal rows matching map[path:tests/canary/fam/pin-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:109: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:114: writes_proposal rows matching map[path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:133: closure_evidence rows matching map[kind:fixture line:1 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:140: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC7_included_BASE_source,\"proposal_test.go:146: closure_evidence rows matching map[kind:fixture path:tests/canary/fam/include-fixture] = 0, want 1:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t01-oc-01-check-1",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:3288aec9c7536ab18dfd942d49e7f21c0ffd7fd2aa83fd25e033a3253d6326e4",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,28325,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,19119,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,3020,56\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-2",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:a593973c216a1b2e94fb10fe12f847692201793805c2a16be7888682320b8dc5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,1998,181\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-2",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-3",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:728181303db3d34458306edb69275900e21dd207b9e0f50fd3737e000ce3c1e8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/canary,pass,8,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-3",
          "command": "bench test --package ./internal/canary",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-4",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:ce9d1917db2487d48bfd8111cc431a886347204de9cea5d7ce2ed81780ddf598",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,1787,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-4",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-5",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:0709dd52e61dc6ae1b11d2e598c50edea3226a53f175ef9a9aa160e543689d0c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/maps,pass,93,97\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-5",
          "command": "bench test --package ./internal/maps",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-6",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:24354281e7e1601b811b62ccde26b0239c321464b14d7c221a8e1c879f438bab",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,593,61\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-6",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-7",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:3e2b2ebc2c9a6bf1cc6e87319e941b831608341b3c124784b677d0dac09897d2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,pass,1041,111\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-7",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-8",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-8",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-9",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:e40c409edd5f000d2d6983eb75f6ed13dcdca994e112e11a186da8cafc3e3e3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14726,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-9",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-10",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:161e4c0891f858c2c387d8b9eb4e5e2711eef9d3295cca22671192b645370716",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4613,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-10",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-11",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:c0efd3489394bdbc45b3084a51d8fe61270b463db0fddc937af269f9bd246cb6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16907,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-11",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-12",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:c2fb9022e89ee0e42801103180d341813bf234e60e573bba2761885ae2a3b36d",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,55903,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-01-check-12",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-13",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:c9eb8ae6de1c44a7dea8bf2e89c1df75f6274721739ba23191590976376f6f3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-13",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-14",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:0570173198b71d68ffdbfaa4facbe957045d752e3f35aaf27496ec6af277bf05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-14",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "ft293_t01-oc-01-check-15",
          "performer": "claude:ft293_t01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01",
            "digest": "sha256:0aaabee70dfcea4b354fdc4790e3db36125abe22a576723f2b6660764c8834f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-15",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-focused",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:c12436754acc3030bb4bbf1ec5345d1646105b754426e84c631e6a6c66da350b",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestProposalFixedPoint/OC3_second-hop_pin,36,1,0,null,yes\n  2,bit,swap,failed,1,TestProposalProvenance/OC5_two_reasons,36,1,0,null,yes\n  3,bit,swap,failed,1,TestProposalProvenance/OC8_registry_owner,36,1,0,null,yes\n  4,bit,swap,failed,7,TestProposalFixedPoint/OC1_missing_pin,36,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,5255\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2192,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:110: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:115: writes_proposal rows matching map[path:internal/anchors/cycle.go] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:149: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:156: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2189,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:149: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",47\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2206,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC8_registry_owner,\"proposal_test.go:172: registry reason = map[kind:registry line:68 path:internal/toon/toon_test.go reason:binding internal/toon requires internal/toon/toon_test.go beside internal/toon/toon.go at kit source 8cd3e918cfedae2f74e5f387587b2ca2d4… (374 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2194,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC1_missing_pin,\"proposal_test.go:107: writes_proposal rows matching map[path:tests/canary/fam/pin-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC1_path_with_a_space,\"proposal_test.go:137: writes_proposal rows matching map[path:tests/canary/fam/space fixture] = 0, want 1:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:110: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:115: writes_proposal rows matching map[path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:149: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:156: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC7_included_BASE_source,\"proposal_test.go:162: closure_evidence rows matching map[kind:fixture path:tests/canary/fam/include-fixture] = 0, want 1:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-01-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t01_rep1",
              "digest": "sha256:c12436754acc3030bb4bbf1ec5345d1646105b754426e84c631e6a6c66da350b",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,4,TestProposalFixedPoint/OC3_second-hop_pin,36,1,0,null,yes\n  2,bit,swap,failed,1,TestProposalProvenance/OC5_two_reasons,36,1,0,null,yes\n  3,bit,swap,failed,1,TestProposalProvenance/OC8_registry_owner,36,1,0,null,yes\n  4,bit,swap,failed,7,TestProposalFixedPoint/OC1_missing_pin,36,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestProposalFixedPoint|TestProposalProvenance|TestMetadataGrammar,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,5255\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2192,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[4]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:110: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:115: writes_proposal rows matching map[path:internal/anchors/cycle.go] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:149: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",26\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:156: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",26\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2189,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:149: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",47\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  3\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2206,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC8_registry_owner,\"proposal_test.go:172: registry reason = map[kind:registry line:68 path:internal/toon/toon_test.go reason:binding internal/toon requires internal/toon/toon_test.go beside internal/toon/toon.go at kit source 8cd3e918cfedae2f74e5f387587b2ca2d4… (374 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  4\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2194,36\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[7]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC1_missing_pin,\"proposal_test.go:107: writes_proposal rows matching map[path:tests/canary/fam/pin-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC1_path_with_a_space,\"proposal_test.go:137: writes_proposal rows matching map[path:tests/canary/fam/space fixture] = 0, want 1:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC3_second-hop_pin,\"proposal_test.go:110: writes_proposal rows matching map[path:tests/canary/fam/registry-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalFixedPoint/OC4_cycle,\"proposal_test.go:115: writes_proposal rows matching map[path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC5_two_reasons,\"proposal_test.go:149: closure_evidence rows matching map[kind:fixture line:2 path:tests/canary/fam/cycle-fixture] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC6_anchor_holder_line,\"proposal_test.go:156: closure_evidence rows matching map[kind:anchor path:internal/anchors/cycle.go] = 0, want 1:\",18\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestProposalProvenance/OC7_included_BASE_source,\"proposal_test.go:162: closure_evidence rows matching map[kind:fixture path:tests/canary/fam/include-fixture] = 0, want 1:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-1",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:57407454c6f84b6dc6348a8127ddb6d84f0287a6b1c4cd481d64552e745016f1",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,35960,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23109,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,4257,64\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-2",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:7c606975adfc4f0e661f6fb0871829141b0fb394e2b80d97172f9da7dd057582",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/anchors,pass,2335,181\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-2",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-3",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:57663b4d057e8bf1d36191dab341a70394ebc1db11803e03784b629a45794331",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/canary,pass,9,32\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-3",
          "command": "bench test --package ./internal/canary",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-4",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:1b785647227ec534c5b9721350755f53026ef4ffc2db3ba66543ffa642a54626",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,2648,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-4",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-5",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:4e9b84fa055982a66e9f2f9ab0c1fe695e1f57ac27db895aca1a877fde591b23",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/maps,pass,118,97\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-5",
          "command": "bench test --package ./internal/maps",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-6",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:e72756f233057d117d1155bcf556399fb6e5db2ef504f808292a61f8f10fae6c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/bounds,pass,684,61\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-6",
          "command": "bench test --package ./internal/bounds",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-7",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:16564278e5eac2733a2d44d7d36ea97b74cc485e5cdc70af39bfbe507adfbc72",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/spec,pass,1223,111\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-7",
          "command": "bench test --package ./internal/spec",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-8",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:5203d455b3fdb4a061354ccd4af581c67c80986f2d286fb730d2980836b7ef92",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,5,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-8",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-9",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:7d1db1e1043ca1e7f1f1e46fee731d72565d42ef51a5bf96c63bfb7afc7fc7dc",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18520,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-9",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-10",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:158c43efd354b3459fa44247ac15ac1ea07c02e4c52a5de3696fc800e9577ff9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5428,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-10",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-11",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:bcc05a4e56110697d12405140eb090d1ff979d23bbfe7b8733f5c3237c72a5c3",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,18762,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-11",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-12",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:48d61003d94f0344f6d50db4a55ef5211697a9eac67749d7da05b88b046521fd",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,67061,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-01-check-12",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-13",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:c9eb8ae6de1c44a7dea8bf2e89c1df75f6274721739ba23191590976376f6f3e",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-13",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-dot-import-package-alias/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-14",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:0570173198b71d68ffdbfaa4facbe957045d752e3f35aaf27496ec6af277bf05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-14",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-multiple-dot-import-wait/files/internal/sample",
          "exit_code": 0
        },
        {
          "id": "ft293_t01_rep1-oc-01-check-15",
          "performer": "claude:ft293_t01_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t01_rep1",
            "digest": "sha256:0aaabee70dfcea4b354fdc4790e3db36125abe22a576723f2b6660764c8834f7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-01-check-15",
          "command": "bench test --package ./tests/canary/package-core-guard/bounds-parenthesized-wait/files/internal/sample",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ft293-oc-c1-standards-r1",
          "performer": "claude:ft293_oc-c1_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c1_standards",
            "digest": "sha256:24dd6a4e20da29772eac17fe07fd14ae9d5abb1c704ec2851fe75fe994623e6f",
            "excerpt": "## Standards axis: OC-C1 (base 1a8fece7d, chunk tip a54b85b3c)\n\nThe evidence check bound: `current=true`, `dirty_record=false`, tree head `8d4c023c`, not dirty.\n\n## Findings\n\n- **S1:** `internal/preflight/ownership/metadata.go:84-127` (`duplicateKey`, called at `manifest.go:147`) adds a new token walker that rejects duplicate JSON keys. The repo already has a strict decoder that owns this fact: `internal/jsonfile/decode.go:24` (`DecodeDocument`), whose `scanValue` at `:69-80` rejects duplicate keys recursively and also rejects trailing values. This breaks the AGENTS.md rule \"one source per fact\": the diff adds another derivation of strict JSON decoding instead of composing the existing owner. Older copies already exist (`reviewrecord/parse.go:242`, `releaseevidence/json_validation.go`), but this diff adds a fresh one. The per-key-set `fields()` check can stay, because the manual-claim key set depends on applicability. Disposition `auto-fix` (decode `top` through `jsonfile.DecodeDocument`; the OC56 reason text changes with it). Confidence 6.\n- **S2:** `internal/tickets/registry_data_test.go:127-139` (`TestBindingRulesProjectBoundFiles`) re-implements `BoundFiles`' dedupe-and-sort over `BindingRules` (`registry_data.go:120-127`) step for step, then compares the two results. This breaks the AGENTS.md rule on independent test expectations. Its independence is not needed to turn a named mutation red, and no red is recorded: `oc-01-omission` names only the second-hop pin, one independent reason, and the fence role. Disposition `auto-fix` (assert the projection's properties, or check it against fixed rows). Confidence 5.\n- **S3:** `internal/preflight/ownership/proposal_test.go:185-189` (`sealedKit`) finds the expected registry line by searching the bytes for `{Prefix: \"internal/toon\",`. That is a second test-side locator of the registry row: the owner test `registry_data_test.go:115` already pins `BindingLines` to the physical line by its own text search. The OC8 case could take `tickets.BindingLines` as its oracle. This is the same AGENTS.md rule as S2, and no red is recorded for this independence. Disposition `ask-user` (the reviewer decides whether to keep the independent locator and record its red, or derive the line from the owner). Confidence 4.\n\n## Author-flagged judgment calls\n\n1. **The `ownershiptest` subpackage:** not a finding. `ownershiptest/repo.go` imports only `preflighttest` and `chargesource`, and `preflight/proposal_test.go:26-49` now delegates to it, so the proposal seeder exists once.\n2. **The `..`/escape refusal in every mode:** not a Standards finding. The check is a single source in the shared collector (`ownership/closure.go:44-50, 68-70`). Whether its reach beyond the proposal is correct is a Spec-axis question against OC76.\n3. **The `… refused` bootstrap kinds:** not a finding. I searched the whole tree (`rg --hidden`) and found no doc, test, or registry that spells the old `fixture inventory not readable` or `anchor registry not readable` kinds. Bootstrap kinds have no single enforced convention.\n\n## Advice (not findings)\n\n- `ownership/closure.go:68` `repositoryPath` and `canary/mutation.go:250` `safeBasePath` both define a clean, local, repository-relative spelling, with different edge rules. The Writes-entry rule may belong in `tickets` beside `WritesPath`, before ticket 3 adds `ParseWrite`.\n- `preflight/closure.go:105-107`: the `\"ownership not readable\"` branch can never run, because `Collect` returns only `*ownership.Error`.\n- `preflight/proposal.go:100`: the `kit-pin` evidence kind is a bare literal outside the `ownership.Kind` vocabulary.\n- `preflight/gather.go:98`: `collected` holds a `Source`, not a collection (unclear name).\n- `ownershiptest/repo.go:61-72` repeats the c0/feature/c1 seeding of `preflighttest/fixture.go:195-204`. The same repetition was in the base `proposal_test.go`; the diff moved it and did not add it.\n- For the Spec axis: the inventory refusal now comes after the changed-files and default-tip refusals (`gather.go:96-102`). Before, it came inside the ticket gather at line 68. It also no longer runs when no ticket has Writes (`ownership/closure.go:41`).\n\n## Count\n\n3 findings. The worst is S1: a fresh strict-JSON duplicate-key walker beside `internal/jsonfile`.\n"
          },
          "axis": "Standards",
          "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
          "tip": "a54b85b3cadc7c89241c5a950a90ae69765ad177",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "ft293-oc-c1-spec-r1",
          "performer": "claude:ft293_oc-c1_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c1_spec",
            "digest": "sha256:396341a675c88c0e5c8432f9f40ebc7ec51f054e2f7a1f348f236e19837421fd",
            "excerpt": "**Spec axis: OC-C1 (ft293-build, pair 1a8fece7..a54b85b3)**\n\nThe evidence handle bound: `current=true`, `dirty_record=false`, head 8d4c023c.\n\nI found 2 issues. Of the 36 rows, 35 held and 1 is unverified. Neither issue blocks a row. The worst one is P1.\n\n## Findings\n\n- **P1.** Spec lines 85 and 121-122 say \"Complete means full supported premise certification\" and \"Completeness states what that author declared\". The defect is at `internal/preflight/ownership/ownership.go:165`. That line sets `snapshot.Complete = manifest.Present && len(snapshot.Diagnostics) == 0` and never reads `manifest.Completeness`.\n  - Traced case: a manifest that declares `\"completeness\":\"partial\"` with empty arrays and no diagnostics gives `Complete=true`. A partial declaration cannot be full certification. By the spec's own rule at line 141, a missing fence is incomplete, so a partial one must be too.\n  - Nothing reads `Snapshot.Complete` in this checkpoint (I grepped `internal/preflight`), so no output changes yet. The value contract is still wrong.\n  - Disposition: `auto-fix`. Confidence 6.\n- **P2.** Spec lines 219-220 say \"The needle must occur exactly once at the declared line. Moved or ambiguous needles are diagnostics.\" Line 564 schedules a refresh \"before OC-C3 activates\".\n  - The ticket commit (`a54b85b3`, 36 changed lines, all in the coverage map) rewrote rows OC1-OC13, OC55-OC60, OC71-OC81, OC83-OC85 and OC90-OC92 from `planned …` to backticked citations. The spec's `bench-planning-facts` seams still quote the old `| planned internal/… |` text as their needles, so those needles no longer match.\n  - The manifest also still maps OC77 to `TestExplicitAuthority` (planned, rows OC69, OC70, OC77, OC86). Its `TestSourceSafety` rows are only OC71-OC76, but the map now cites OC77 at `source_test.go`.\n  - When ticket 05 turns on the OC88 cross-check, this stale inventory will refuse.\n  - Disposition: `no-op` for this checkpoint. The refresh is owed through the enabling-plan authority before OC-C3 starts. Confidence 6.\n\n## Advice\n\n- Author call 1 (OC77 moved to `source_test.go`): acceptable. Ticket 01 Covers OC77, and the planned file belongs to ticket 08. Pair the move with the P2 refresh.\n- Author call 2 (a `..` or escaping Writes entry refuses gather in every mode): consistent with spec line 730, \"Paths escaping the repository are refused\".\n- Author call 3 (`ownershiptest` subpackage): inside the `internal/preflight/ownership` fence and Writes entry. Acceptable.\n- Refusal order: the fixture and anchor inventory refusals moved from `gradeTickets` to after the changed-paths and default-tip refusals (`gather.go:86-101`). Neither refusal is on the line-207 list, but the precedence did change.\n- `InventoryComplete` and `CurrentComplete` come out true for an empty declared-complete inventory, without the manual-location or seam-citation cross-checks that spec lines 146 and 189 describe. Nothing consumes them, so successors must not read them before those checks land.\n- OC78 builds `responsebound.New` by hand rather than going through the CLI wrapper (`cmd/bench/census_output.go:22`).\n- OC85's system test is required at the first hop, not the second, despite the row's \"second-hop\" wording.\n\n## Rows\n\nOC1-OC13 held; OC55-OC60 held; OC71-OC77 held; OC78 unverified (hand-built bound owner, not the CLI route); OC79-OC81 held; OC83 held; OC84 held; OC85 held (weak second-hop witness); OC90-OC92 held.\n\nCount: 2 findings (P1 auto-fix, P2 no-op). Worst: P1.\n"
          },
          "axis": "Spec",
          "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
          "tip": "a54b85b3cadc7c89241c5a950a90ae69765ad177",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "ft293-oc-c1-coverage-r1",
          "performer": "claude:ft293_oc-c1_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c1_coverage",
            "digest": "sha256:ac9a1e02c20ae77fd50861c8f9398976d5adebbd6b0e2cd3548b7790ffd68783",
            "excerpt": "## Coverage axis: OC-C1 (ft293-build, pair 1a8fece7d..a54b85b3c, HEAD 8d4c023c)\n\nThe evidence bound: `--check-current` returned current=true and dirty_record=false.\n\n### Findings\n- **C1.** `internal/preflight/ownership/premises_test.go:36`, row OC56 (\"Last-key-wins decoding cannot pass\"). The only duplicate-key case is the top-level `\"version\"` key. The check walks the whole payload (`internal/preflight/ownership/metadata.go:115-117`), but no test has a duplicate key inside a nested object (location, subject, expected, owner). For nested objects, `fields()` at `metadata.go:132` decodes into a map, so the last key wins silently. My probe 1 limited the duplicate check to the top-level object, and it went silent. Needed test: a TestMetadataGrammar case with a nested duplicate key, for example `\"expected\":{\"count\":1,\"count\":2}`. Disposition: auto-fix. Confidence 8.\n- **C2.** `internal/preflight/ownership/proposal_test.go:318-343` and `internal/canary/inventory_test.go:297-318`, row OC83 and the edge \"present-unreadable inventory\". Both witnesses only exercise the wrong-type states: a regular file and a symlink at `tests/canary`. Neither exercises a `tests/canary` directory that exists but cannot be read, which `bounds.classifyDir` returns as `StateUnreadable` (`internal/bounds/classify.go:94-96`). My probe 2 made `Pins` treat `StateUnreadable` as empty, and it went silent. Needed test: a chmod-000 `tests/canary` subcase that skips under root. Disposition: auto-fix. Confidence 6.\n- **C3.** `internal/preflight/ownership/premises_test.go:38`, row OC58. Spec line 127 says \"Each item has a unique id\". The decoder collects IDs across all arrays (`internal/preflight/ownership/manifest.go:211,221`), but the only test duplicates an ID inside `premises`. A per-array ID set would therefore pass. Needed test: a cross-array duplicate-ID case. Disposition: auto-fix. Confidence 6.\n- **C4.** Spec line 729 names spaces and CRLF documents as path cases. No ticket-01 test puts a space in a fixture, anchor, or Writes path, and none uses a CRLF spec carrying the planning-facts fence. No row maps either edge. Disposition: ask-user. Confidence 4.\n\n### Advice\n- OC72 (`source_test.go:128`) checks only the substring `\"spec\"`, while OC71, OC73 and OC74 check `\"spec not readable\"` or `\"spec not found\"`.\n- The 2 MiB payload bound at `metadata.go:77` cannot be reached through the command, because the spec reader applies the same `ControlRecordLimit` first.\n- `canary.FixturePins` and `anchors.ReferencingFiles` now have no production callers, only tests.\n\n### Probes\n1. `internal/preflight/ownership/metadata.go`: swapped `if stack[top].keys[value] {` for `if top == 0 && stack[top].keys[value] {`. Ran `--package ./internal/preflight/... --run \"TestMetadataGrammar|TestProposal|TestSource\"` (69 tests). Verdict silent, first_failure null, restored yes.\n2. `internal/canary/pins.go`: swapped `case bounds.StateAbsent, bounds.StateEmpty:` for `case bounds.StateAbsent, bounds.StateEmpty, bounds.StateUnreadable:`. Ran `--package ./internal/preflight/... --package ./internal/canary/... --run \"TestProposal|TestSource|TestPins|TestFixturePins\"` (58 tests). Verdict silent, first_failure null, restored yes.\n\n### Summary\n4 findings (C1-C4). The worst is C1. `git status --short` printed nothing.\n"
          },
          "axis": "Coverage",
          "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
          "tip": "a54b85b3cadc7c89241c5a950a90ae69765ad177",
          "finding_ids": [
            "C1",
            "C2",
            "C3",
            "C4"
          ],
          "supersedes": []
        },
        {
          "id": "ft293-oc-c1-standards-r2",
          "performer": "claude:ft293_oc-c1_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c1_standards_r2",
            "digest": "sha256:bac4220f3abaf3922ccf5f251b0ba629d708b8f758fd62f28eae8e06b9de22ee",
            "excerpt": "## Result\n\nThe evidence check bound: `current=true`, `dirty_record=false`, head `ad54e4f1`, not dirty. S1 and S3 are closed. S2 is closed in the code, but its record is open. I found no new Standards finding in the repair delta.\n\n## S1-S3\n\n- S1: closed. `duplicateKey` is gone from the package. `rg` finds no `duplicateKey`, `json.NewDecoder` or `Token()` walker under `internal/preflight/ownership`. `internal/preflight/ownership/manifest.go:151` now calls `jsonfile.DecodeDocument`. The `scanValue` function at `internal/jsonfile/decode.go:69-80` rejects a repeated key at every nesting level.\n- S2: closed in code, record open. `internal/tickets/registry_data_test.go:128-146` no longer runs `BoundFiles`' dedupe-and-sort a second time. It now checks against fixed rows. The refreeze probe in the review record is still the four `oc-01-omission` mutants; none targets this test. AGENTS.md accepts such an independent expectation only when its red is recorded and demonstrated. Fix: record a red for this test. Disposition `auto-fix` (record only), confidence 6.\n- S3: closed. `internal/preflight/ownership/proposal_test.go:200-205` takes the line from `tickets.BindingLines` (`registry_data.go:134`). The byte search and the `bytes` import are gone. Refreeze mutant 3 (`OC8_registry_owner`) still goes red.\n\n## New findings\n\nNone.\n\n## Advice\n\n- The S2 fixed rows do not show that `BoundFiles` sorts or dedupes.\n- `inventory_test.go:26` and the proposal \"unreadable directory\" case skip on `os.Geteuid() == 0`; the existing pattern at `internal/bounds/classify_test.go:294-299` checks readability after the chmod instead.\n- `manifest.go:152` wraps a message that `decode.go:53` already prefixes with \"decode persisted JSON\".\n\n## Reaffirmation after the evidence-only correction\n\nS2 is now closed in both the code and the record, with one condition: the record section is not committed yet. In the working copy of `reviews/preflight-ownership-closure.md`, lines 779-796 (\"## Independent expectation reds\") record the probe at `ad54e4f1`. The verdict is `bit` and `first_failure` is `TestBindingRulesProjectBoundFiles`. The test failed on its fixed-row check at `registry_data_test.go:139`, and the file was restored. That red shows what AGENTS.md requires: the fixed rows catch a swapped file in the toon row. The record counts as evidence only once it lands with the chunk. Nothing else in my round changes.\n"
          },
          "axis": "Standards",
          "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
          "tip": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c1-standards-r1"
          ]
        },
        {
          "id": "ft293-oc-c1-spec-r2",
          "performer": "claude:ft293_oc-c1_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c1_spec_r2",
            "digest": "sha256:0a23366869afc18b59301a94c103e104e3141a9f88d683299e8974bc60f3b6f2",
            "excerpt": "## Result\nSpec axis, OC-C1 repair delta `a54b85b3..10c458e5`: P1 is closed and P2 is closed. I found no new Spec finding. The evidence `sha256:c363d8bf…` bound: head `ad54e4f1`, current true, not dirty.\n\n## P1: closed\n- `internal/preflight/ownership/ownership.go:165` now reads `snapshot.Complete = manifest.Present && manifest.Completeness == completenessDeclared && len(snapshot.Diagnostics) == 0`. The constant is defined at `manifest.go:134`. This matches spec lines 85-86 and 121-122.\n- The Snapshot doc comment at `ownership.go:146-150` now names \"a partial inventory or any manual claim\".\n- The test is the new `Complete …` subtests in `premises_test.go` (TestMetadataGrammar). They call `ownership.Collect` directly and cover three cases: declared-complete gives true, partial gives false, and a manual claim gives false.\n\n## P2: closed\nI parsed the facts line (spec:1008) with a script.\n- Needles: all 159 spec-sourced locations hold their needle at the declared line. Each needle occurs exactly once in the spec outside the facts line itself.\n- Seams: 13 seams are `existing`. Each has exactly one `func <name>(` in its cited file, and each row that names a seam cites that path and function.\n- OC77: it is under `TestSourceSafety` (`source_test.go`).\n- OC92: it has three seams, all of which exist: `TestMetadataGrammarFenceRoles`, plus `TestLiftedFieldScanReportsDuplicateField` and `TestLiftedFieldScanSkipsFencedLines` in `internal/maps/maps_parse_test.go`. This matches row spec:716.\n\n## C2 / OC83: holds\n- Row spec:707 cites `TestProposalUnknown` through the actual command. `proposal_test.go:362-374` shows that an unreadable `tests/canary` refuses through the real proposal (the checkout walk) and never gives an empty proposal.\n- The collector's own refusal is tested at `internal/canary/inventory_test.go:179-188` through `FixturePins`, which projects `Pins` (`pins.go:81-83`); `Pins` refuses any state other than absent, empty, or parsed (`pins.go:48-55`).\n\n## C4: satisfies the Edge inventory (spec:729)\n- Spaces: the subtest `OC1 path with a space` (`proposal_test.go:267-279`) uses a spaced fixture directory and a spaced Writes path.\n- CRLF: the case `CRLF fence section` (`premises_test.go:210`) runs a full CRLF section through to the manual-claim row.\n\n## Advice\n- The space case asserts nothing about `internal/example/with space.go` itself.\n- Row OC83 could also cite `TestPinsRefuseASpecialInventory`.\n- `TestBindingRulesProjectBoundFiles` now hand-writes the `internal/toon` rules; it grades a hardening-cycle check, so it is advice only.\n- The switch to `jsonfile.DecodeDocument` keeps OC56 intact.\n\nCount: 0 new findings; P1 and P2 closed.\n"
          },
          "axis": "Spec",
          "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
          "tip": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c1-spec-r1"
          ]
        },
        {
          "id": "ft293-oc-c1-coverage-r2",
          "performer": "claude:ft293_oc-c1_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "9ff2d86d94d6e271723b57a02a2d85a3df083996",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c1_coverage_r2",
            "digest": "sha256:4fa8571708b2b88f61d958c60b18b2c749c1210194d3919705c498b9558ba9ed",
            "excerpt": "## Result\nThe Coverage confirming axis for OC-C1 is clean. C1 through C4 are closed, P1 bites, and the repair delta has no new coverage finding. Evidence `sha256:c363d8bf…` bound: `current=true`, `dirty_record=false`, head `ad54e4f1b`.\n\n## Findings\n- C1 (OC56): closed. `premises_test.go:199` puts a nested `\"expected\":{\"count\":1,\"count\":2}` in `OC56 nested duplicate key`. The case expects `duplicate object field \"count\"`, and the top-level-only mutant (probe 2) bites it.\n- C2 (OC83): closed. `internal/canary/inventory_test.go:25-34` adds `unreadable` (mode 000, with a capability skip under root). Probe 1 bites it. `proposal_test.go:316-328` asserts the checkout-walk refusal `tests/canary: permission denied`.\n- C3 (OC58): closed. `premises_test.go:202` adds a `scope` item and a `premises` item that both carry ID `p1`, and expects `duplicate item id \"p1\"`. `manifest.go:191-192` keeps item IDs across every array. Read, not probed.\n- C4 (edge inventory): closed. `proposal_test.go:267-279` uses a space in both the fixture directory and a written path, and asserts the `writes_proposal` and `closure` rows. `premises_test.go:210` turns the whole fenced section to CRLF and expects the same `semantic`/`manual` result as the LF case. Read, not probed.\n- P1 Complete cases: bite. `premises_test.go:218-232`, probe 3.\n- New findings: none.\n\n## Advice\n- `manifest.go:154` (`if top == nil`) has no case behind it, and probe 4 survives it. This is a near-equivalent mutant on behavior that existed before the repair. It is not a finding.\n\n## Probes (all restored=yes)\n1. `pins.go:50`: add `bounds.StateUnreadable` to the absent/empty case; canary + preflight, `TestPins|TestProposal`; bit; `TestPinsRefuseASpecialInventory/unreadable`.\n2. `decode.go:87`: swap `if err := scanValue(decoder, child); err != nil {` with `_ = child; var raw json.RawMessage; if err := decoder.Decode(&raw); err != nil {`; preflight, `TestMetadataGrammar`; bit; `TestMetadataGrammar/OC56_nested_duplicate_key`.\n3. `ownership.go:159`: omit ` && manifest.Completeness == completenessDeclared`; preflight, `TestMetadataGrammar`; bit; `TestMetadataGrammar/Complete_partial_inventory`.\n4. `manifest.go:154`: swap `if top == nil {` with `if false {`; preflight, `TestMetadataGrammar|TestProposal`; silent; null.\n\n## git status --short\n(empty)\n"
          },
          "axis": "Coverage",
          "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
          "tip": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c1-coverage-r1"
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
      "from": "sha256:9a5a573300d6ec1ac2977c518eb468f936fdbcf261e48fa72be1b495dcf98b66",
      "to": "sha256:9e4941f4243f14167ae9a17c2f09edf8d5b4ae67633d4f0d597fa7017edfb0cb",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ]
      }
    }
  ]
}
```


## Repair state

OC-C1: repair cycle 1 of 2 is complete at 10c458e5 (claude:ft293_t01_rep1, opus, high). No reviewer extension exists.

## Standards

OC-C1: 3 findings. The worst issue is S1.

- S1 (auto-fix, confidence 6): `internal/preflight/ownership/metadata.go:84-127` adds a duplicate-key token walker beside the strict decoder `internal/jsonfile/decode.go:24`, which owns that fact.
- S2 (auto-fix, confidence 5): `internal/tickets/registry_data_test.go:127-139` derives the `BoundFiles` projection again in the test, with no recorded red for that independence.
- S3 (auto-fix, confidence 4): `internal/preflight/ownership/proposal_test.go:185-189` locates the registry line a second time; the orchestrator takes `tickets.BindingLines` as the one source under the auto-approve direction.

## Spec

OC-C1: 1 actionable finding. The worst issue is P1.

- P1 (auto-fix, confidence 6): `internal/preflight/ownership/ownership.go:165` sets `Complete` without a read of the declared completeness, so a partial manifest reads complete.
- P2 (no-op, confidence 6): the planning-facts needles quote the planned row text. The orchestrator refreshes them in the enabling plan commit of this repair.

## Coverage

OC-C1: 4 findings. The worst issue is C1.

- C1 (auto-fix, confidence 8): no OC56 case puts a duplicate key in a nested object; a top-level-only mutant stays silent.
- C2 (auto-fix, confidence 6): no OC83 case reads a present but unreadable `tests/canary` directory; a mutant that treats it as empty stays silent.
- C3 (auto-fix, confidence 6): no OC58 case repeats one item ID across two arrays.
- C4 (auto-fix, confidence 4): no case covers a path with a space or a CRLF spec with the planning-facts fence. The orchestrator assigns both edges to ticket 01 under the auto-approve direction.

## Independent expectation reds

S2: the fixed rows of TestBindingRulesProjectBoundFiles are an independent expectation. The orchestrator demonstrated their red at ad54e4f1 with this probe.

```text
bench worktree exec ft293-build -- bench probe internal/tickets/registry_data.go --swap '"internal/toon/toon_test.go",' --with '"internal/toon/toon_fixture_test.go",' --expect TestBindingRulesProjectBoundFiles --package ./internal/tickets --run TestBindingRulesProjectBoundFiles
probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:
  bit,internal/tickets/registry_data.go,swap,failed,1,TestBindingRulesProjectBoundFiles,yes
selection[1]{form,target,run,baseline,ran}:
  package,./internal/tickets,TestBindingRulesProjectBoundFiles,passed,1
probe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:
  bench-probe-result-v2,0,1,512
packages[1]{package,status,elapsed_ms,tests_run}:
  github.com/gibbonmi/bench/internal/tickets,fail,2,1
failures[1]{package,test,line,lines}:
  github.com/gibbonmi/bench/internal/tickets,TestBindingRulesProjectBoundFiles,"registry_data_test.go:139: BindingRules(internal/toon/toon.go) = [{internal/toon internal/conformance/data_handling_test.go} {internal/toon internal/toon/toon_fixture_test.go}] and BoundFiles = [internal/conformance/data_handling_test.go in… (389 bytes)",2
skips[0]{package,test,reason}:
```
