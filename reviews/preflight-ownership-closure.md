# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-ownership-closure/spec.md",
  "plan_digest": "sha256:fa4a88a7ab799bed77ecdce96efa188b078dbdacb71b33a4b4d540237f415c27",
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
    },
    {
      "id": "OC-C2",
      "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
      "tip": "cd54d245ab7a1aa2fdf0ca4431125adb6264f205",
      "plan_digest": "sha256:9de99bb1f6e996574324f2d278b8357129972bb0143631cd3ffa66f0300512e3",
      "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
      "acceptance_rows": [
        "OC14",
        "OC15",
        "OC16",
        "OC17",
        "OC18",
        "OC19",
        "OC20",
        "OC21",
        "OC22",
        "OC23",
        "OC24",
        "OC25",
        "OC26",
        "OC27",
        "OC28",
        "OC29",
        "OC30",
        "OC31",
        "OC82"
      ],
      "verification": [
        {
          "id": "ft293_t03-oc-03-focused",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:b6ae8b1f0d7f16f325525d203e3c3779471ee2637b395051e688df9dc20ac218",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,26,1,0,null,yes\n  2,bit,swap,failed,1,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1994\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,4277,26\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,4502,26\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",16\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03",
              "digest": "sha256:b6ae8b1f0d7f16f325525d203e3c3779471ee2637b395051e688df9dc20ac218",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,26,1,0,null,yes\n  2,bit,swap,failed,1,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,26,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1994\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,4277,26\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,4502,26\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:96: closure requires internal/merge/c_test.go = false, want true:\",16\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03-oc-03-check-1",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:c3127d1aa47384337f951a1bce7f8bd4f56ebd8ce00cc4d820202c8270115dbb",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,26821,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,18505,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9688,100\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-2",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:2e9a04fec9cd2537bc6c46831f26777355a052836477ce4d3a67c5ec693d39fa",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2031,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-3",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-4",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:c5a9cf8b89cb727157d3caa1dbc5225b010fa0f5be96242c7cc3f79b69edbb58",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,5581,160\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-5",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:5347f0c26a602c2dd4473b22e8ff6b0f02599b25d4514052cda66afa86792985",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,16362,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-6",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:76d1ca56cc57e9df53de30050b810fa1c7c0c3c8865fffa68a3751e2541c0f05",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16087,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-7",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:6c133a58afa8aa8eede00032ee98b829fd2c4e59c44ccaa6e8008cf4cef29615",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4428,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-8",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:374c88073279b8b718dca7843d1171265b053b5b1896ff30c636cc4b885c6f7d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15029,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03-oc-03-check-9",
          "performer": "claude:ft293_t03",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03",
            "digest": "sha256:80854353d04fe4f6107f2b2ea3d86aaab6d0fc130a936241911c0449d89ffbb7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,54337,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02-oc-02-focused",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:4ad2e7421764c53c909c43fa9c0039992487383fe17d39a7759560a546664cf8",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,10,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,10,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1759\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,622,10\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:96: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,662,10\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:102: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02",
              "digest": "sha256:4ad2e7421764c53c909c43fa9c0039992487383fe17d39a7759560a546664cf8",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,10,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,10,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1759\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,622,10\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:96: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,662,10\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:102: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02-oc-02-check-1",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:a8f6a267e4e66b063f51d5ac551e879f3dac6f53e584b7927ab17c54e9a58069",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,26645,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,18566,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,9282,100\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02-oc-02-check-2",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02-oc-02-check-3",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02-oc-02-check-4",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:9facd5c4f6ef3e20ace24aee209a1a42bcab768df1cf2fecc75dbfecfdea39e5",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14985,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02-oc-02-check-5",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:be9e2ab0ab32ca97f98cd9e227c7201b23ead61a2bab2539580f4ac0d81f2d0d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,3865,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02-oc-02-check-6",
          "performer": "claude:ft293_t02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02",
            "digest": "sha256:16398a8f129e8c3cd2c5ab7dba0f48ddaa73be498c07ad613e166aa0f5eb8664",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15084,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-focused",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:d6e995d14c5a4d20712da486ffccc9fd23112fb42cad82f38ac0ed935c9e6ccd",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,32,1,0,null,yes\n  2,bit,swap,failed,1,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,32,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1998\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,9,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,6751,32\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,6665,32\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",16\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03_rep1",
              "digest": "sha256:d6e995d14c5a4d20712da486ffccc9fd23112fb42cad82f38ac0ed935c9e6ccd",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,32,1,0,null,yes\n  2,bit,swap,failed,1,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,32,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1998\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,9,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,6751,32\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,6665,32\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:117: closure requires internal/merge/c_test.go = false, want true:\",16\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-1",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:4f71406dac5e08d73767c4262fec611d2969b32b060cd8341d06f07a61a227d8",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,35378,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23798,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,14330,114\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-2",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:54ec426f52a7afd740d174b274f2069949f26a0c44633da640da46b5d371b4e1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2332,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-3",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-4",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:e103ed0e9f3b04fa3d50960374c389160dafd55775011ebd25253edda2b88ec6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6471,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-5",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:388629b08ade60e9544b75ebc7cf2c7029b58451230323ba17e93236d9fa8a34",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,19479,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-6",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:b09d261be1f478975a5781bcc9efbf00a0a7d47e218af6829ce1a34580ef1af3",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18334,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-7",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:c4e5d9e211d0c17789bd82f5d77c9845cec61a985e1e22bea2c76fa95dadfff2",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4782,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-8",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:367d37c06d26bed348d053612b019903dde190674cf87ba09cf7d2b18f0e1048",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17410,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep1-oc-03-check-9",
          "performer": "claude:ft293_t03_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep1",
            "digest": "sha256:a0495e13e85751acbece601ab1f158ac90ff37232672172f63fc267ffe76016f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,60311,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1-oc-02-focused",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:e692a7634c26c0047e4942d96e81188aadd1af693e15075d0278e649aa4501bb",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,18,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2050,18\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:117: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:117: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2186,18\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:123: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02_rep1",
              "digest": "sha256:e692a7634c26c0047e4942d96e81188aadd1af693e15075d0278e649aa4501bb",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,18,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,18,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2050,18\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:117: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:117: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2186,18\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:123: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02_rep1-oc-02-check-1",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:8e845f3d68ac2202617ef475e4d73ffa49ad9c13bfd319d74158fafb9e8c7f6a",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,36429,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23815,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,14481,114\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1-oc-02-check-2",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1-oc-02-check-3",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1-oc-02-check-4",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:537fda177fca884a34ec71aaed4eba950971056f40c884e7de8ce0c415991163",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16332,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1-oc-02-check-5",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:25004b0a2ebeb87d91a2f19f756524b936fe7f8fbc070ac7acaf7461819b6758",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4502,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1-oc-02-check-6",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:7220f9fea8e756bc3b2f1902b0ca8b197406d2dde5fb77f6ca0eea2c0efe6e53",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17186,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-focused",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:d4fcead17ec29ab803726512e44f9da880d250ebd265e5f7db72704c73c69b2d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,33,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,33,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2443\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,7064,33\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,7024,33\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03_rep2",
              "digest": "sha256:d4fcead17ec29ab803726512e44f9da880d250ebd265e5f7db72704c73c69b2d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,33,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,33,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2443\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,7064,33\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,7024,33\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-1",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:a5651bb000e8cb96f8a70a57f11edd9b4dff37d5f6e1f5e3d4564e0414001a6e",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,33349,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,22790,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,15354,116\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-2",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:e3b73a61beb5b5bcd7b0ab624fb4598dd9036f8f7ad56e711f53ba3dd54937b0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2096,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-3",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-4",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:6ac918931bb67aba5025f3c17870b1838bb468c273f18884d341d04e830c65ce",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6257,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-5",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:e4d854fa35979e1181ba9ea2ff7d1576f08b121f5d9ef91e88e1fe8f8769ca9f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,17742,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-6",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:dd9a064a72b22366121975452f07d2ec28bc3eb4aefe3ccb52b763543e53521b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17360,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-7",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:e9ca5ceae6e9c4f86dd0b0c1cec53b142fbc63320e6c803d8c00521973c3d2a6",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4373,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-8",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:fee8cb086b696bd852cfae2dd4a51806401e7f3605cd14c48034a2878a0d6b93",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16329,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep2-oc-03-check-9",
          "performer": "claude:ft293_t03_rep2",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep2",
            "digest": "sha256:5a872474f7cdef1762460d6759835383883a23edb5d0c3311482461e3b4dcb2b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,58841,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1b-oc-02-focused",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:6185dc20c35cae96dbf2fdad9831b3715c3a5cfbfc2872a4250caad49b49bd3a",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2152,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2075,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02_rep1",
              "digest": "sha256:6185dc20c35cae96dbf2fdad9831b3715c3a5cfbfc2872a4250caad49b49bd3a",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2152,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2075,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02_rep1b-oc-02-check-1",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:b1f225575203ff696b7cf629d26efc1c32da19be0778e4b2e6103842ebbfc041",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,33182,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,22957,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,15255,116\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1b-oc-02-check-2",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1b-oc-02-check-3",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1b-oc-02-check-4",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:83c0ef91a1002f9aebd22acfc8f690a5014c72d93b0b6f210ce3b2c227f7f4d9",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15943,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1b-oc-02-check-5",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:db5d1d07ec2d8c517fd5516659e7b6fc0df7494783f3ef7d31167a67dce1aa16",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4455,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1b-oc-02-check-6",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:24490aa2d1571dd41b99628d7e4c051cb2b5668460b5a402f36225fa70e080d9",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16017,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-focused",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:d2ee3e58412533b74e09fb4c692987cbe9e6510ca1214aca830a611adb1d6bb4",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,43,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,43,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2445\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,12349,43\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,12379,43\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03_rep3",
              "digest": "sha256:d2ee3e58412533b74e09fb4c692987cbe9e6510ca1214aca830a611adb1d6bb4",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,43,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,43,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2445\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,12349,43\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,12379,43\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-1",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:9fa01d10dcc377ae34575db31a8939b236bfcf3c165beba72d78dca709e98d1d",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,34203,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,24541,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,22126,126\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-2",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:9371e6c16a2f5c99915f46e5656c26cd3343a0c24c5322ad89ee1dcfa8ea1fa8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2170,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-3",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-4",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:0463703cc64a2f269e6a4dc2fe462bd30b0bf60f11bb74c65074175a58f568e8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6378,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-5",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:e467ddbe3d15989abf375ab494bef284def0fb643d82a6dc126e9f58642c0f84",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,17760,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-6",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:6c2537e9f516861ca69d466efd2efa71556138f041278a13063f096fc2f45ad1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17169,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-7",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:e5084f56dc2f21580b889240ca77dc7fb8ba0b46d2aa5c744782798acc6e97ce",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4381,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-8",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:8f6b68dac5fd87caa04c669e5c81b4e9846c842d2bdb664e481d4e67f8cb23c0",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16170,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep3-oc-03-check-9",
          "performer": "claude:ft293_t03_rep3",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep3",
            "digest": "sha256:42f2cc8b9c842d63909a6a7ecd0904e0e4f023c1c07249990e8641cd72632b06",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,58868,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1c-oc-02-focused",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:554e05d6c3747bb16450bb0d9dca9733e962fc4c7af1f41a5853fe3e435a5e4d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2117,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2081,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02_rep1",
              "digest": "sha256:554e05d6c3747bb16450bb0d9dca9733e962fc4c7af1f41a5853fe3e435a5e4d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2117,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2081,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02_rep1c-oc-02-check-1",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:cd7a4611efed814a5d704959ec7e423cfc35b4bc12ff4d1ff5a81cf1d8b3ef54",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,34311,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,24211,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,22267,126\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1c-oc-02-check-2",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1c-oc-02-check-3",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1c-oc-02-check-4",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:b924d744b4acf7cbe242c1c20de31526dbedd0b00a0df8a1fa3fb49990c1758b",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15715,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1c-oc-02-check-5",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:e2615c261a1c3a71c48a7c394dc2acb30d4f0ae9a7f37fb16e1a339bfdd514c8",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4509,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1c-oc-02-check-6",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:0dc9d639c942b55e6248607748b89c88821e3f22a71243a8d1ca0105b7d90c25",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16118,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-focused",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:0c558b6ec1ebda59eaddd170c2888062353d23abcec1e19cf5154a7165499fd5",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,51,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,51,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2445\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,17071,51\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,17019,51\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03_rep4",
              "digest": "sha256:0c558b6ec1ebda59eaddd170c2888062353d23abcec1e19cf5154a7165499fd5",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,51,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,51,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2445\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,17071,51\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,17019,51\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-1",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:01b8c15bcf268886bb24b93f324208b051a1c8064e408127a45c199a9546f464",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,36356,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,25164,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,28445,134\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-2",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:f09526987f65cdbfd33f41d2155c15ca8e70207449b7f46121236083de4d4cc2",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2186,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-3",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-4",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:50a65aa88a382747d2f0fb0497ba0f649f554fa9278d11c807a0c14d3815ebb1",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6673,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-5",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:ebab5c1baa164131602934a8c453a689d4084a8419a19b8f868d6dc6d66e904c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,18568,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-6",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:40e3400e6c481a7a63ac4b858e6ab4c66345cd86aab5cb160f8bc5c9038d4252",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,19199,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-7",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:7c2cb2aed999554105fad99d26614ce2751d8d2dde37d10d16f4367c2ec5aa0e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4725,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-8",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:6fecc44ab6e9bb6df6475fec11dd2719c1d082e0adb183df4e2d7c2df0b0a38d",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17478,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep4-oc-03-check-9",
          "performer": "claude:ft293_t03_rep4",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep4",
            "digest": "sha256:7b6152eb1863379668363395bf07df4149c37d5547402d8f23430054e4828fa8",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,62455,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1d-oc-02-focused",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:bbf5826f6054a95d71e63e3a1a9318dfd8d88dda101a6b84ef535d1629c2191f",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2256,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2356,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02_rep1",
              "digest": "sha256:bbf5826f6054a95d71e63e3a1a9318dfd8d88dda101a6b84ef535d1629c2191f",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,6,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2256,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2356,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02_rep1d-oc-02-check-1",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:a7613653ae5f00dc451d01619330173b84cb6b5153fd5bed818692f254ea3078",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,38185,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,26392,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,29303,134\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1d-oc-02-check-2",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1d-oc-02-check-3",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1d-oc-02-check-4",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:ccab55d0352ecd0d753290c154efba16b711351338553b34dfe8d925438d0cd6",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,16349,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1d-oc-02-check-5",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:a1f344ecbef80a1b42c3b95aefe44d81d5b9ed45ad916d3e532aefd1dd226ebb",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4821,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1d-oc-02-check-6",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:900ec92be8cb97207577793cf3ccfa62720209312e6ffcc35d2c6d84c0420209",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17338,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-focused",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:8420de53fd87f5a4c6623e9d5682daadf0865fabfdd0946c38507d67bf33f44d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,57,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,57,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2445\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,22355,57\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,21556,57\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03_rep5",
              "digest": "sha256:8420de53fd87f5a4c6623e9d5682daadf0865fabfdd0946c38507d67bf33f44d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,57,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,57,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2445\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,22355,57\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,21556,57\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-1",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:ac15f739a98850600cb45cfc640d7669351b0896b8595fa24282e83ea1a29a42",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,38483,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,26352,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,33662,140\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-2",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:574cbe112182f1cb33fab5c3de161e535b5f82dca41f8102c10ec2212ea24c37",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2317,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-3",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-4",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:fb7c3fa590826bacd1bf282638487d30dbdb1dc4573ee39c978bf96940ef171c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,7006,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-5",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:ca6a83d662569b805515dede76bc61a319e336a49092d12c17db1ee4f5ad8855",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,19964,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-6",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:19b06cac7c7a11545112c9b76ec51a0002198eb74e953f8dbc29c94735df2f06",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,18670,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-7",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:9f15b945556ec07c0b5abfb21b0ed5a1de44d38ef6b26010c08e40454940d50f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4792,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-8",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:9b8e51af122bb51cfdb8b23792b6e5a158717464e91b60d20eb13ba736c4947e",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17972,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep5-oc-03-check-9",
          "performer": "claude:ft293_t03_rep5",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep5",
            "digest": "sha256:843bef966a34fe3ed64ae075feec6b52f10f704d101d9c34e05c0cc4d12f95a0",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,65784,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1e-oc-02-focused",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:95f5478f00c0ec301b71b4c018e6decb0dd38aa6ee0d3d4366ccc19f19f04e1a",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2003\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2264,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires internal/conformance/subcommand_routing_table_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2450,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02_rep1",
              "digest": "sha256:95f5478f00c0ec301b71b4c018e6decb0dd38aa6ee0d3d4366ccc19f19f04e1a",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2003\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2264,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires internal/conformance/subcommand_routing_table_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2450,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02_rep1e-oc-02-check-1",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:2d733b0df8dc46c76548a4b3b283d81742d8b696406edc71716ae01bc8a5bb70",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,39436,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,27084,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,34732,140\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1e-oc-02-check-2",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1e-oc-02-check-3",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1e-oc-02-check-4",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:053fe15976a37db7cdf27c0d01e7a033d0c4392e705816ad3a4343c12d47bd36",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17395,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1e-oc-02-check-5",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:106b5dcc7085feb9f57da6091ac38c5003ae0221dc09e2039112c3867117e597",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4979,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1e-oc-02-check-6",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:7a79f50bf7a187df17ca36a10d6309a42a3a7c1cb82d58bc84f769a5c8d43281",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,17550,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-focused",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:c974cfc96bc2652233b0ad5afe628e29cd0a8277a5bbe1ec26f7fed6dc6017e0",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,62,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,62,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2446\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,23591,62\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,15,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,24091,62\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestHelperCallers|TestDeletedEntry|TestOwnerPremises'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-03-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t03_rep6",
              "digest": "sha256:c974cfc96bc2652233b0ad5afe628e29cd0a8277a5bbe1ec26f7fed6dc6017e0",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,3,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,62,1,0,null,yes\n  2,bit,swap,failed,3,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,62,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestHelperCallers|TestDeletedEntry|TestOwnerPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2446\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,8,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,5,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,23591,62\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC20_signature_change_proposes_each_static_test_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",14\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,15,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,24091,62\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[3]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC22_relocation_retains_the_before-state_caller,\"effects_test.go:120: closure requires internal/merge/c_test.go = false, want true:\",16\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deleted_fixture_helper_retains_its_external_before-state_caller,\"effects_test.go:120: closure requires internal/other/other_test.go = false, want true:\",9\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestHelperCallers/OC23_deletion_retains_the_before-state_callers,\"effects_test.go:120: closure requires internal/merge/b_test.go = false, want true:\",18\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-1",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:a4f8a4295fbc503218166dc4d222cdbc06cda27d047aee940017369a644c7849",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,40032,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,27377,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,37969,145\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-2",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:b9fdebd891a03b7e197a0d78e0f4f1b447dfb9528f2ea14236d902f7dfd6e229",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/consumers,pass,2856,66\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-2",
          "command": "bench test --package ./internal/consumers",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-3",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:88e523cdedf9efe73ab2c21b4389e19f812e14d597f51de6a81ad7ab19ed4742",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,4,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-3",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-4",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:e4a13b3f49776475e4ff46ef4d705346e772bbe4bee33e63c4237dee3662e573",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/commitment/repository,pass,6760,161\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-4",
          "command": "bench test --package ./internal/commitment/repository",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-5",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:d1d4650f849d5609548375a7496b346dab30e172c51ede6af0f540857ea22985",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,17098,254\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-5",
          "command": "bench test --package ./internal/reviewrecord/recordcmd",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-6",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:85b4cf24da87c43902a758055a2910cbe3bc087ba814c8edfcf130a7ba6d3774",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,17042,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-6",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-7",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:ffe676b027c1509cd1e224193a23ada68bb6bbc4a1543ff7d0559145ae3f3a20",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4354,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-7",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-8",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:015de64be2f0db0b84a7c5c2c994768786752853a7c34ec7ec7e9f39409080cb",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16255,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-03-check-8",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t03_rep6-oc-03-check-9",
          "performer": "claude:ft293_t03_rep6",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t03_rep6",
            "digest": "sha256:5fb4f2933cfaf845fe2f660b354dba7805e3e72e5986db2be8dbf2a45acb68c7",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,58183,1694\nfailures[0]{package,test,line,lines}:\nskips[1]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n"
          },
          "requirement": "oc-03-check-9",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1f-oc-02-focused",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:6a44b4f7c7829aaa29702316e62bb1ae8c9037ded97c6ef3d36dd3a926919c0d",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2097,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2049,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-focused",
          "command": "bench test --package ./internal/preflight/... --run TestGrammarEffects",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-02-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t02_rep1",
              "digest": "sha256:6a44b4f7c7829aaa29702316e62bb1ae8c9037ded97c6ef3d36dd3a926919c0d",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,19,1,0,null,yes\n  2,bit,swap,failed,1,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,19,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestGrammarEffects,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,1982\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2097,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC14_body-only_edit_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",94\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/handler_literal_body_adds_no_grammar_obligation,\"effects_test.go:120: closure requires cmd/bench/help_inventory_test.go = true, want false:\",74\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,2049,19\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[1]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestGrammarEffects/OC17_planned_unbound_form_requires_the_root_help_inventory,\"effects_test.go:126: evidence for cmd/bench/help_inventory_test.go = [map[kind:registry line:0 path:cmd/bench/help_inventory_test.go reason:Grammar binding cmd/bench requires cmd/bench/help_inventory_test.go beside cmd/bench/command_registr… (635 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t02_rep1f-oc-02-check-1",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:465e7cd23e82631780b6c9f64953bc1cd639cfa9bb13db8e2d01e452b7405937",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,34949,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23921,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,32959,145\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1f-oc-02-check-2",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:95d52b294f98d073b76e9962d185b8982e753f19036f3e01460ff209e98b3278",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/tickets,pass,3,48\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-2",
          "command": "bench test --package ./internal/tickets",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1f-oc-02-check-3",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:9dc3a2106ba5d02c9c72884fc4cdb6f759767003b2ad7ad20a784047d4059575",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/usage,pass,2,31\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-3",
          "command": "bench test --package ./internal/usage",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1f-oc-02-check-4",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:c62b73d054c0271ba9870a78fb6a46eb986bf1dcbbce55b1fde11f059735a299",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15259,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-4",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1f-oc-02-check-5",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:4bc7c06f27de461a02190e719210865c6852d9cf1c6202805dfd861ead1122fc",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  ticket-grammar,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,4321,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-5",
          "command": "bench test --check ticket-grammar",
          "exit_code": 0
        },
        {
          "id": "ft293_t02_rep1f-oc-02-check-6",
          "performer": "claude:ft293_t02_rep1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t02_rep1",
            "digest": "sha256:ca956bc8314dfb8b2568da6cc757975b2227ddcbdd1a037fcee2a029c8de5d7f",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15750,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-02-check-6",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "ft293-oc-c2-standards-r1",
          "performer": "claude:ft293_oc-c2_standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards",
            "digest": "sha256:252bc9f3fc3bfbbb4595d525fbb417698167add00f8a306f4b33f3f04c5886ff",
            "excerpt": "## Standards axis, OC-C2 (tickets 02 and 03)\n\nThe evidence check bound. `bench preflight evidence sha256:3a1d…e2e --check-current` returned `current=true` and `dirty_record=false` at HEAD `0e108087`. The only change after `b73e3715` is in `reviews/preflight-ownership-closure.md`, so the worktree files I cite are the frozen code.\n\n### Findings\n\n- **S1** (ticket 02, `auto-fix`, confidence 8). The new producer walker writes its own copy of how a Go declaration is named.\n  - `internal/preflight/ownership/producers.go:121-134` (`receiverName`) copies `internal/consumers/rows.go:329-350` almost line for line.\n  - `producers.go:92-118` (`index`) re-derives the function, `Type.Method`, type, and var/const naming. `rows.go:295-296` calls `declNames` \"the one derivation of the declaration-name forms\", and `internal/consumers/blast.go:51-52` says both walks read it.\n  - This breaks the one-source-per-fact rule in `AGENTS.md`. The const-group-as-one-unit choice can stay; the receiver name and the method spelling should come from consumers. The repair goes to ticket 03, whose `Writes:` line holds both `internal/preflight/ownership (new)` and `internal/consumers/rows.go`.\n- **S2** (tickets 02 and 03, `auto-fix`, confidence 7). The \"local name of an import\" rule (use `spec.Name` if set, else `path.Base` of the path) is written three times.\n  - `producers.go:84-90` and `producers.go:157-165` (ticket 02).\n  - `internal/consumers/rows.go:211-219` (ticket 03, in `firstUse`).\n  - `producers.go:159` also matches the import path by suffix with its own code, beside `internal/consumers/resolve.go:148-154` (`pkgSuffixMatches`). Same rule as S1.\n- **S3** (ticket 03, `auto-fix`, confidence 6). `localName` (`internal/preflight/ownership/helpers.go:153-157`) reads the symbol grammar a second time, and the two readings disagree.\n  - `blast.go:99-101` says Resolve owns that grammar. `resolve.go:21-23,55-69` accepts a bare name and a package suffix with several segments.\n  - `localName` assumes exactly one qualifier segment. A bare `Kind` becomes `\"\"`, and `preflight.ownership.Kind` becomes `ownership.Kind`.\n  - Result: the tip side resolves through Resolve, while the base census (`deletion.go:195`) and `declaredRoot` (`effects.go:432-434`) look for a different name. They report a wrong \"declares no … at the source base\" or \"no registered grammar producer\" diagnostic. It fails closed, but on the wrong reason.\n  - `ownerKeys` decoding (`manifest.go:138`) does not check the symbol's shape.\n- **S4** (ticket 03, `auto-fix`, confidence 6). The set of member-census kinds is written twice: in the switch at `helpers.go:184` and in the `shape` table at `helpers.go:266-268`. A kind added to the switch only gets a zero shape and a misleading key error.\n- **S5** (ticket 03, `auto-fix`, confidence 5). There is a second per-kind table beside `closureFamily`.\n  - `internal/preflight/proposal.go:149-158` (`proposalSource`) maps each kind to its label in its own switch. Ticket 03 added the `caller` case there, and any unknown writing kind silently falls to \"registry\".\n  - Ticket 03 also rewrote the comment at `internal/preflight/closure.go:18-22` to say \"a new kind joins each projection from this one table\". With `proposalSource` outside the table, that is now false.\n  - The switch already existed in part at the base.\n\n### Advice (no IDs, not counted)\n\n- `helpers.go:359`: the parameter `consumers` shadows the imported `consumers` package.\n- `closure.go:21`: \"has no verdict row yet\" points to a future change. The timeless-present register in `bench-craft-comments` prefers \"has no verdict row\".\n- `registry_data.go` `rowKey` joins with \".\" separately from `GrammarRoot.String`. `TestBindingLinesLocateEveryRow` catches any drift, so this is acceptable.\n- `helpers.go:83-86`: `declaration` drops a Resolve \"malformed symbol\" error and treats the symbol as absent.\n\n### Author-flagged calls\n\n- T02, AST references in `producers.go`: after ticket 03, the reference walk (`producers.go:329-364`) answers a different question than `consumers.Rows`. The typed `Load` can read only the clean tip, not the base, so this walk is not a duplicate. The naming and import knowledge is, which is S1 and S2.\n- T02, separate `producers.go`: no finding; it is 364 lines.\n- T03, `SourceUses` in `rows.go` and the shared `localName`: it reuses `parseDecls` and `topLevelDecls` correctly, but `localName` raises S3. The disposition list and owner-key list do each have one source.\n- T03, `owner` kind and the `writes` flag: no defect, apart from S5.\n\n### Other hunt items, all clean\n\n- There is no second write-entry parser outside `internal/tickets/writes.go`.\n- The `ownership` package imports `internal/preflight` only in the test files.\n- Moving the disposition decoding into `metadata.go` kept its validation.\n- The new expectations have recorded reds: `oc-02-omission` and `oc-03-omission` are in the review record.\n\n**Count:** 5 findings: 0 `ask-user`, 5 `auto-fix`. The worst is S1.\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "6a9f1a46e203cb965abf89f3795e176308c21a07",
          "finding_ids": [
            "S1",
            "S2",
            "S3",
            "S4",
            "S5"
          ],
          "supersedes": []
        },
        {
          "id": "ft293-oc-c2-spec-r1",
          "performer": "claude:ft293_oc-c2_spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec",
            "digest": "sha256:e9ebb45129f696838ded8cb18ded143e4ebc91073292ccf5728cd2103c0e5365",
            "excerpt": "## Spec axis: OC-C2 (tickets 02 and 03), base 10c458e5 to code b73e3715\n\nThe evidence check bound (`current=true`, `dirty_record=false`).\n\nResult: 7 findings, all on ticket 02 or 03. All 19 rows hold. The worst issue is P1.\n\n## Findings\n- **P1** (ticket 02) | Spec 318 \"walked through typed references\", 321 \"An unresolved argument or producer dependency stays diagnostic\". Code: `internal/preflight/ownership/producers.go:329-364`. `deps` matches same-package names by syntax only. An identifier that is not a package unit is skipped silently (354), and so is a selector on an imported package (348-352). Traced input: `h := helpText(); usage.Parse(usage.Grammar{Help: h}, args)` — an edit to `helpText` classifies as body-only and drops the Grammar obligation with no diagnostic. A typed walk is possible at the tip now that `consumers.Load` exists; the base side still has no typed loader. Disposition: ask-user, confidence 7.\n- **P2** (ticket 02) | Story 5, spec 312, 317, 333. Code: `effects.go:77`, `92-94`. A changed `usage.Parse` grammar in a package with no Grammar row adds no inventory obligation (19 production packages). A custom root in such a package does. Disposition: ask-user, confidence 5.\n- **P3** (ticket 03) | Spec 246, 262, 264. Code: `helpers.go:266-268`, `294-297`, `metadata.go:84-96`. A posture or printed-field item that follows the spec table is refused because of an unknown key. `writer_probe` accepts any object. Disposition: ask-user, confidence 7.\n- **P4** (ticket 03) | Spec 264, 370. Code: `helpers.go:317-331`. Code that only reads the field counts as a producer, so a correct producer-only census never turns green. Disposition: ask-user, confidence 6.\n- **P5** (ticket 03) | Spec 476 \"Path in closure is the originating write entry\". Code: `ownership/closure.go:75`. For an owner requirement, `Entry` holds the premise ID, so `closure.path` shows `home-owner`. Disposition: ask-user, confidence 6.\n- **P6** (ticket 02, spec text) | Spec 241-242 contradicts the change rows at 225-226. The code follows the table and takes the inventories from `tickets.GrammarFiles()`. Non-behavioral; disposition: no-op, confidence 7.\n- **P7** (ticket 03) | Spec 367, 371. Code: `helpers.go:112`, `effects.go:168`. Nothing grades `changes[].dispositions` of a liveness change. Disposition: ask-user, confidence 4.\n\n## Author-flagged calls\nT02-1 holds; T02-3 shortfall (P1); T02-4 gap (P2); T02-5 holds; T02-6 holds; T03-1 holds except P5; T03-2 see P3; T03-3 holds for this chunk.\n\n## Advice\n- The OC17 test asserts 2 of the 5 holders; looping over `tickets.GrammarFiles()` would cover all five.\n- A Writes entry that no committed change covers keeps the full Grammar obligation.\n- A changed non-Go file under a Grammar package is classified as body-only (`effects.go:77`).\n- `helpers.go:83-86` discards a `Resolve` error.\n- The fixture uses `internal/merge/`, where spec 508 says `internal/example/`.\n- Ticket 08 must make dispatch refuse when a proposed caller is missing.\n\n## Rows\nOC14-OC31 and OC82 held (OC17 asserts 2 of 5 holders; OC31 held with P4 separate). Every case goes through the real `preflight.Command` on committed git fixtures. I ran no test or probe.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "6a9f1a46e203cb965abf89f3795e176308c21a07",
          "finding_ids": [
            "P1",
            "P2",
            "P3",
            "P4",
            "P5",
            "P7"
          ],
          "supersedes": []
        },
        {
          "id": "ft293-oc-c2-coverage-r1",
          "performer": "claude:ft293_oc-c2_coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7a4200509101962b10089aa0e9004f4e9320e5e6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage",
            "digest": "sha256:941621eb447d1f3eb57caaff079a9daf613ffc90d87f661b095dfa0ebde231dd",
            "excerpt": "## Coverage axis: OC-C2 (tickets 02 and 03)\n\nThe evidence `sha256:3a1da1eb…` binds: `current=true`, `dirty_record=false`, head `0e108087`. The worktree is clean.\n\n### Findings\n\n- **C1.** `internal/commitment/repository/light_path.go:165`. A light-path ticket with an `x.go (deleted)` entry. Ticket 03 says this consumer shares the marker rule (ticket 03, line 23). No test in `light_path_test.go` uses a `(deleted)` entry. My probe replaced `WritesPath` with the old strip that only removes `(new)`, and the probe stayed silent. The record consumer (`internal/reviewrecord/plan.go:116`) has no marker test either. Ticket 03, `auto-fix`, confidence 8.\n- **C2.** `internal/preflight/fence_writes.go:47`, `paths_authorized.go:65`, `ownership/closure.go:296`, `testreport/floor.go:51`. The `(deleted)` projection of `WritesPath` in the delegating consumers. Every TestDeletedEntry case uses the directory-only `helperFence` (`effects_test.go:203`), and the proposal tolerates a red `fence-writes`. A `WritesPath` that returns the raw deleted entry would probably pass. Not probed. Ticket 03, `auto-fix`, confidence 5.\n- **C3.** `internal/preflight/ownership/producers.go:318-320`. \"Ambiguous symbols … cannot certify completion\" (spec line 727). No ownership test exercises the ambiguous-declaration diagnostic. My probe omitted it, and all tests under `./internal/preflight/...` still passed (silent). Malformed Go at the base or tip (`producers.go:63-66`) also has no test fixture. Ticket 02, `auto-fix`, confidence 8.\n- **C4.** `internal/preflight/ownership/producers.go:348-352`. A grammar root that depends on a value in another package. `deps` stops at a selector on an imported package and records no diagnostic. Example: `cmd/bench/main.go:131` reads `commit.HelpRowSuffix` and `commit.LaneClause` from unbound `internal/commit`, and `help_inventory_test.go:145` pins that help row. The fix is a design decision. Ticket 02, `ask-user`, confidence 6.\n- **C5.** `internal/preflight/ownership/helpers.go:135` with `internal/consumers/rows.go:191`. A signature change or relocation of a method that still exists at the tip, where one caller stops using it. `SourceUses` returns no base uses for a method name, and the method diagnostic fires only when `obj == nil`. No test uses a method symbol. Ticket 03, `auto-fix`, confidence 6.\n\n### Advice\n- OC24 and OC26 only assert exit 0.\n- `SourceUses` skips a base file that does not parse (`rows.go:196`).\n\n### Probes (both restored `yes`)\n| Site | Mutation | Command | Verdict | first_failure | restored |\n|---|---|---|---|---|---|\n| `producers.go` | omit the ambiguous-declaration diagnostic | `--package ./internal/preflight/...` | silent | null | yes |\n| `light_path.go` | swap `WritesPath` for a strip that only removes `(new)` | `--package ./internal/commitment/... --package ./internal/worktree` | silent | null | yes |\n\n**Count:** 5 findings. **Worst:** C1, then C3.\n\n### `git status --short`\n(empty)\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "6a9f1a46e203cb965abf89f3795e176308c21a07",
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
          "id": "ft293-oc-c2-standards-r2",
          "performer": "claude:ft293_oc-c2_standards_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r2",
            "digest": "sha256:155599f78ec082bb70156ca08d959c5e3211a6c99fa25b52d40fb68ca115600e",
            "excerpt": "## Result\nAll five first-round findings (S1-S5) are closed. I found one new low-confidence finding (S6). The evidence binds: `bench preflight evidence ... --check-current` returned `current=true`, `dirty_record=false`, head `6f42627e`, and the tree is clean.\n\n## Findings\n- S1 closed. `internal/consumers/resolve.go:269` holds `declNames`, and `DeclNames` at line 296 wraps it. `receiverName` at line 302 is the only remaining definition in `internal`. `internal/preflight/ownership/producers.go:99` and `:106` call `consumers.DeclNames`.\n- S2 closed. `consumers.ImportName` is defined at `internal/consumers/rows.go:301`; its callers are `rows.go` `firstUse`, `producers.go:92` and `producers.go:145`. `PathHasSuffix` at `resolve.go:183` backs `pkgSuffixMatches` and `producers.go:145`. `ownership/imports.go` adds no new import-name or suffix derivation.\n- S3 closed. `readings` (`resolve.go:79`) is shared by `Resolve` and `LocalName` (line 95). `ownership/helpers.go:161-162` calls `consumers.LocalName`. `ownership/metadata.go:140` `identity` refuses a symbol with no declaration spelling.\n- S4 closed. `ownership/members.go:48-52` `memberCensuses` is the one table.\n- S5 closed. `internal/preflight/closure.go:23` `closureFamily` carries a `source` func for each kind; `proposalSource` is deleted.\n- `consumers.WriteRows` complies with the one-source rule: `rows.go:73` composes the one walker, the shared `sortRows`, and one write derivation (`markWrites`, `rows.go:188`).\n\n### S6: two shapes for \"a read at one revision, or its fault\"\n- `ownership/imports.go:56-60` adds the generic `revRead[T]`. Beside it, `ownership/effects.go:30-31` keeps `producerRead`, which has the same shape and an almost identical doc comment. `producerRead` could be `revRead[producers]`.\n- AGENTS.md, one source per fact. Ticket 02 (repair 05639910/4268cb7c). Disposition auto-fix. Confidence 4/10.\n\n## Advice\n- `ownership/members.go:89-94` and `:114-121` branch on literal kind strings beside the `memberShape` table.\n- The comment at `ownership/manifest.go:141` says the writer probe \"also\" takes `reusedTestKeys`; nothing else reads that list.\n- `internal/conformance/registry/packages.go:36` keeps its own import-suffix match (pre-existing, outside the chunk).\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "f5a06b8d024eeb0ada0e7ff7617e53228c5b9723",
          "finding_ids": [
            "S6"
          ],
          "supersedes": [
            "ft293-oc-c2-standards-r1"
          ]
        },
        {
          "id": "ft293-oc-c2-spec-r2",
          "performer": "claude:ft293_oc-c2_spec_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec_r2",
            "digest": "sha256:b4166aea9913c766de0c6ce927dcce70927a7b813e558397ee3325e1f58bb18e",
            "excerpt": "## Result\nThe evidence binds: `--check-current` returned `current=true`, `dirty_record=false` at tip `6f42627e`. P1, P2, P3, P4, P5 and P7 are closed. I found one new finding, P8.\n\n## Closures\n- P1/C4 closed. Per-file import table `producers.go:297`; selector-argument roots `producers.go:203-209`; in-module edges only `producers.go:521-529`; no FuncLit walk `producers.go:507-508`; dot import diagnostic `producers.go:498-500`; undeclared or ambiguous target diagnostic `imports.go:134-143`; referrers by transitive importers `imports.go:163-196`; a referrer contributes only resolved roots `imports.go:104-116`; base and tip read `effects.go:101-110`, `171-178`; test `effects_test.go:174`.\n- P2 closed. `grammarScope` filter removed (`effects.go:79-87`); unbound grammar change adds `tickets.GrammarFiles()` (`closure.go:230-237`); test `effects_test.go:181`.\n- P3 closed. Exact key sets `members.go:49-51`; `writer_probe` decoded with `reusedTestKeys` (`members.go:93`, `manifest.go:133`); subject and expected lists one set (`members.go:98-106`).\n- P4 closed. `consumers.WriteRows` keeps only static writes (`rows.go:207-209`, `markWrites` `rows.go:272-307`); census uses it at `members.go:51`.\n- P5 closed. Owner requirement uses `Entry: entry` for each write (`ownership/closure.go:73-75`); OC28 checks the write path at `effects_test.go:330`.\n- P7 closed. `census.liveness` grades an existing liveness change against the retirement members (`members.go:18-33`), wired at `ownership/closure.go:62`; actions enum `metadata.go:130`; exact-one count `members.go:199-227`; test `effects_test.go:254-257`.\n\n## New finding\n- P8. Spec 264, 371, 481. `internal/consumers/rows.go:294-305` (ticket 03). `markWrites` counts a field only through a `KeyValueExpr` key. An unkeyed struct literal such as `row{\"a\", reason}` writes the field without naming it, so neither `WriteRows` nor `Rows` sees that producer. A printed-field premise that leaves out such a producer still grades green, but spec 481 allows green only for complete proof. Disposition auto-fix: make an unkeyed literal of the field's struct type a census diagnostic (spec 271), or count it as a write. Confidence 5/10.\n\n## Advice\n- A grammar root that calls a method on an in-module external value records an edge to the value but never walks the method body.\n- A parser argument `otherpkg.Fn()` stays diagnostic through `resolvedCall`.\n- The census does not treat an address taken from a field as a write.\n- An unresolved producer in an unchanged referrer never fires its diagnostic.\n\n## Row status\nOC14-OC31 and OC82 hold at their cited seams.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "f5a06b8d024eeb0ada0e7ff7617e53228c5b9723",
          "finding_ids": [
            "P8"
          ],
          "supersedes": [
            "ft293-oc-c2-spec-r1"
          ]
        },
        {
          "id": "ft293-oc-c2-coverage-r2",
          "performer": "claude:ft293_oc-c2_coverage_r2",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f0b6c72029965212a1c0721e3fa7e177faa971a",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage_r2",
            "digest": "sha256:d81c2b0b6db688f5686c3ebd4c0475bdf952e13974f906f2b704049f06f9eb3a",
            "excerpt": "## Coverage confirming axis, OC-C2\n\nThe evidence binds. `--check-current` returned `current=true`, `dirty_record=false` at tip 6f42627. My one probe bit, and `restored=yes`.\n\n- C1: closed. `internal/commitment/repository/light_path_test.go:89` adds the `deleted-marker` case. `light_path.go:165` reads the entry through `tickets.WritesPath`, which strips both markers (`internal/tickets/writes.go:66-71`).\n- C2: closed. `internal/preflight/fence_writes_test.go:66` adds `internal/x/gone.go (deleted)` against an exact fence file. `unionSide` (`fence_writes.go:47`) and `paths_authorized.go:65` take the same split.\n- C3: closed. The ambiguous case reason comes from `ambiguity()` (`imports.go:52`). I probed the two malformed-Go cases: the mutant swapped `return producers{}, err` for `continue` after `parser.ParseFile` (`producers.go:68-71`). Verdict `bit`, failed_tests 2, first_failure `TestGrammarEffects/malformed_base_source_stays_diagnostic`, restored `yes`.\n- C5: closed. The case \"a present method's former caller stays diagnostic\" expects \"merge.T.Merge is a method…\". The repair removed the `obj == nil &&` guard in `helpers.go`.\n- P1/C4, P2, P3, P4, P5, P7: unknown (not matched to cases).\n\n## New finding\n- C6. `producers.go:361` sets `edge{..., called: called[n]}`, and `imports.go:138` uses it to reach the target function's body. The only cross-package call in the fixtures (`names.Count()`) sits inside a handler `FuncLit`, which returns before the edge is recorded; the other cross-package references are a const and a var. A mutant `called: false` should survive every case. Missing test: a case like `Help: names.Help()` where the tip edits the body of `names.Help` and the test expects the external file to be required; also a reach through two packages (`follow`). Disposition auto-fix. Confidence 5/10.\n\n## Advice\nThe multi-hop external walk in `follow` could get a fixture where package A's grammar reads B and B reads C.\n\n`bench worktree exec ft293-build -- git status --short` printed nothing.\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "f5a06b8d024eeb0ada0e7ff7617e53228c5b9723",
          "finding_ids": [
            "C6"
          ],
          "supersedes": [
            "ft293-oc-c2-coverage-r1"
          ]
        },
        {
          "id": "ft293-oc-c2-standards-r3",
          "performer": "claude:ft293_oc-c2_standards_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r3",
            "digest": "sha256:66b15cbbb55df581c2d812af261dbb8dbd3370ef3ff28b92799aa7f54bbcfb58",
            "excerpt": "## Result\n\nS6 is closed. I found no new blocking Standards defect in the cycle 2 delta. The evidence check bound: `current=true`, `dirty_record=false`, tip `6fcce6a2`.\n\n## S6: closed\n\n- The tree no longer contains `producerRead`. An `rg` of `internal/preflight/ownership` and `internal/consumers` found no match.\n- `effects.go:22`, `:54`, `:163` and `:173` all use `revRead[producers]`, the one generic shape at `imports.go:56-60`.\n- The local is now named `read` (`effects.go:169-171`).\n- Every call site reads `.value`. No stale comment mentions the removed type.\n\n## New findings\n\nNone. `consumers.UnkeyedWrites` (`internal/consumers/rows.go:77-110`) reuses `declKey`, `relPath` and `sortRows`, and duplicates no other walker. The `unnamed` reader on `memberShape` (`members.go:42`, `:137-140`, `:271`) is set only for printed-field.\n\n## Advice\n\n- Three comments say roughly the same thing about the unkeyed-literal rule.\n- `UnkeyedWrites` calls `pkg.Info.TypeOf(lit)` twice (`rows.go:90`, `:96`).\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "e2981232cd2cee9da7fe77c1d0ea63332bda0355",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-standards-r2"
          ]
        },
        {
          "id": "ft293-oc-c2-spec-r3",
          "performer": "claude:ft293_oc-c2_spec_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec_r3",
            "digest": "sha256:73e7f1a7e5078e2980e811cc9233796dbec2964a566faef06f41292374164acd",
            "excerpt": "## Result\nP8 is mostly closed, but one case is still open: a literal whose `&T` is elided. I found no other new Spec defect. The evidence bound: `--check-current` returned current=true, dirty_record=false, tip 6fcce6a2.\n\nP8: closed for explicit and value-elided literals, open for pointer-elided ones. Spec 271, 371 and 481 apply.\n- Fixed: `memberCensuses[\"printed-field\"]` now has `unnamed: consumers.UnkeyedWrites` (internal/preflight/ownership/members.go:271). `members` returns the diagnostic before it counts rows (members.go:137-141), and the test checks for `unknown` and the reason (effects_test.go:364-372).\n- Still open: internal/consumers/rows.go:96 only counts a literal whose recorded type's `Underlying()` is `*types.Struct`. When the `&T` is elided, as in `[]*row{{\"c\"}}` or `map[string]*row{\"k\": {\"c\"}}`, go/types records the hint type `*row`, not `row` (GOROOT src/go/types/literals.go:131-139 and :347). So that literal is skipped, the unkeyed write is not reported, and the premise can still read green. Disposition: repair. Dereference a `*types.Pointer` before the struct assertion, and add a fixture for the elided-pointer case. Confidence 7; reasoned from the go/types source, no probe.\n\nOC28-OC31 still hold at their seams by reading.\n\n## Advice\n- The `producerRead` to `revRead[producers]` change folds two copies of one type and changes no behavior.\n- The new cross-package called producer case adds coverage and conflicts with nothing in the spec.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "e2981232cd2cee9da7fe77c1d0ea63332bda0355",
          "finding_ids": [
            "P8"
          ],
          "supersedes": [
            "ft293-oc-c2-spec-r2"
          ]
        },
        {
          "id": "ft293-oc-c2-coverage-r3",
          "performer": "claude:ft293_oc-c2_coverage_r3",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "663abbaceab1f130ecb229fafd94eaf72bffd69b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage_r3",
            "digest": "sha256:1c276882995b645283159cf99e89946d8148397e1c781822d604ce7b68a28c51",
            "excerpt": "## Coverage, round 3, OC-C2: C6 is closed and I found no new IDs\n\nEvidence binds. `--check-current` returned `current=true` and `dirty_record=false`, with tip `6fcce6a2fd894930675d8096f311a9cb03296307` clean.\n\nC6: closed. The text `called: called[n]` appears once, at `internal/preflight/ownership/producers.go:361`. The new case at `internal/preflight/ownership/effects_test.go:166-171` is the first failure in a recorded probe.\n\n| verdict | subject | mutation | failed_tests | first_failure | restored |\n|---|---|---|---|---|---|\n| bit | internal/preflight/ownership/producers.go | swap `called: called[n]` -> `called: false` | 1 | TestGrammarEffects/cross-package_called_producer_requires_its_root_and_the_shared_inventories (`effects_test.go:120: closure requires internal/harnesses/called.go = false, want true`) | yes |\n\nP8 case (read only): it bites. The case at `effects_test.go:364-372` adds an unkeyed `row{\"c\"}` to a producer set that the `complete-field` row grades green. Removing the `shape.unnamed` call, setting `unnamed` to nil for printed-field, or making `UnkeyedWrites` return nothing would each fail it. The case checks the exact file, line and wording of the reason.\n\nNew findings: none.\n\n## Advice\n- The field-match comparison at `internal/consumers/rows.go:30` is not tested; no binding requirement covers this precision.\n\n## git status --short\n(empty)\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "e2981232cd2cee9da7fe77c1d0ea63332bda0355",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-coverage-r2"
          ]
        },
        {
          "id": "ft293-oc-c2-standards-r4",
          "performer": "claude:ft293_oc-c2_standards_r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r4",
            "digest": "sha256:fd4f56427f3915d54ab7f158bddc802e92a999a4aa33f0b602149546e5a7e585",
            "excerpt": "## Result\nI found nothing for chunk OC-C2, cycle 3 delta (`7cf83f37..81a33375 -- internal`, commits bbbc45e4 and 81a33375). The evidence bound: `--check-current` returned `current=true` and `dirty_record=false` at tip `ca7eac2b`.\n\n## Details\n- One source per fact: `literalStructs` (`internal/consumers/resolve.go:260-289`) is the only code that works out which struct a literal writes. The old struct type assertion was replaced, not left alongside. A union whose terms share one underlying struct appends the same row twice; `sortRows` removes it with `dedupe` (`rows.go:150`).\n- Hand-rolled stdlib: none. Go 1.25 `go/types` exports no core-type API. Confidence 8.\n- Comments: they state the reason, in the present tense, with no provenance.\n- Test table (`internal/preflight/ownership/effects_test.go:382-406`): one harness and one expectation string inside a single loop.\n\n## Advice\n- The test comment at `effects_test.go:382-384` repeats the \"elided element literal records its hint type\" fact from the `literalStructs` doc comment. Optional.\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "b5a07c8078e734e739a718f9fbe83ed9216d62be",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-standards-r3"
          ]
        },
        {
          "id": "ft293-oc-c2-spec-r4",
          "performer": "claude:ft293_oc-c2_spec_r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec_r4",
            "digest": "sha256:c4bd3986bd6e58cddfb4f89095bc52f5259f9b4ae0cef8960363ccac92a38081",
            "excerpt": "P8 is still open for one form. Confidence 7. The evidence bound: `--check-current` returned current=true, dirty_record=false.\n\nOpen form: a type-parameter literal whose constraint is a tilde term over an anonymous struct with the same layout as the field's struct. Example: `func renderC[T interface{ ~struct{ Reason string } }]() T { return T{\"c\"} }` with a caller `renderC[row]()`. The fixture declares `type row struct{ Reason string }` (`internal/preflight/ownership/merge_fixture_test.go:55`), so `row` is in T's type set and the literal writes `row.Reason` by position.\n- `literalStructs` (`internal/consumers/resolve.go`) visits the union term and returns the anonymous `*types.Struct`, whose fields have their own declaration positions.\n- `declKey` (`resolve.go:229-240`) builds its key from the position, so the check at `internal/consumers/rows.go:96-98` never matches. The premise then reads green, which breaks spec 264, 371 and 481.\n- The keyed form `T{Reason: \"c\"}` is missed the same way.\n- Every other listed form is caught, each with a case in `effects_test.go:368-378`.\n\nP9: converting from an identical struct type writes the field. Examples: `row(struct{ Reason string }{\"c\"})`, or `row(other{\"c\"})`. UnkeyedWrites and WriteRows both miss this, keyed or unkeyed. Under spec 481 this is binding if a conversion counts as a write form. Confidence 5.\n\nOC28-OC31 hold on reading.\n\nAdvice: a fix could compare struct identity (`types.Identical` against the target struct's underlying type) instead of only declaration position.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "b5a07c8078e734e739a718f9fbe83ed9216d62be",
          "finding_ids": [
            "P8",
            "P9"
          ],
          "supersedes": [
            "ft293-oc-c2-spec-r3"
          ]
        },
        {
          "id": "ft293-oc-c2-coverage-r4",
          "performer": "claude:ft293_oc-c2_coverage_r4",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "1f8bd0b3c8506232401552a27f867390695082c7",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage_r4",
            "digest": "sha256:583d43970a4d722849cc7ae03b075c6ab0d4f9703fc0279961a6ae8b7db18dfa",
            "excerpt": "## Result\nCoverage, round 4, chunk OC-C2: no finding. The evidence `sha256:73e5cd74…` binds (`current=true`, `dirty_record=false`, tip `ca7eac2b`).\n\n| file | mutation | verdict | failed_tests | first_failure | restored |\n|---|---|---|---|---|---|\n| internal/consumers/resolve.go | omit the `case *types.Pointer: if deref { visit(u.Elem(), false) }` branch (lines 273-276) | bit | 5 | TestOwnerPremises/a_pointer-elided_map_element | yes |\n\nThe other four red subtests were the pointer-elided slice element, the named-pointer element, the pointer-elided type-parameter element, and the pointer-cored type-parameter element.\n\nCase table bite by reading (`internal/preflight/ownership/effects_test.go:367-388`): the struct, interface, and union branches of `literalStructs` are each covered.\n\n## Advice\n- A `deref=true` mutant at `resolve.go:275` only widens the answer; no test needed.\n- A pointer to a type parameter whose own constraint is a type parameter is handled but not tested. Optional.\n\n## git status --short\n(empty)\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "b5a07c8078e734e739a718f9fbe83ed9216d62be",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-coverage-r3"
          ]
        },
        {
          "id": "ft293-oc-c2-standards-r5",
          "performer": "claude:ft293_oc-c2_standards_r5",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r5",
            "digest": "sha256:f54db1c38a9ec6d2774ae7b9315246fd5eda4d2e7e61d0adb0afe42200959f97",
            "excerpt": "## Result\nNo finding on the Standards axis for OC-C2 cycle 4. The evidence check bound: `current=true`, `dirty_record=false`, tip `84096486439cd7f7b417c23e0b11b6cf149ab4f0`.\n\nThe move into `loader.go`: no binding requirement is broken. The only binding text is `.agents/skills/bench-craft-seams/SKILL.md:96`: \"Split along responsibility, never line count, and never fragment to dodge the limit\". The four symbols (SourceUse, SourceUses, firstUse, ImportName) are one lexical, untyped reading, and they moved whole without change. `blast.go` (338 lines) cannot take them under the cap, a new file is outside the fence, and `resolve.go` and `rows.go` have no room, so `loader.go` is the only fence file that fits.\n\n## Advice\n- `internal/consumers/loader.go:79-161` puts lexical source reading next to the typed loader.\n- `Row.Via` now carries two vocabularies: the three column values (`rows.go:10-16`) and the prose fragments `viaUnkeyed`/`viaConversion` (`rows.go:72-76`).\n- `resolve.go:269-271`: \"nil when target is no field\" leaves out a field whose struct is not a package-scope type.\n- `fieldwrites_test.go:18` repeats the `viaConversion` text.\n- Not tested: a selector write like `s.Reason = …` on a variable of the same unnamed struct type, later assigned as `row`, seems to escape both `identicalKeys` and `UnnamedWrites`.\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "449b26b52db1864bedada8c26129df197b70b6e7",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-standards-r4"
          ]
        },
        {
          "id": "ft293-oc-c2-spec-r5",
          "performer": "claude:ft293_oc-c2_spec_r5",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec_r5",
            "digest": "sha256:18b7e588fb8196ebd5a2e14e9454b0537fae9e0b3c9bb0066b016a98b833bfbe",
            "excerpt": "Evidence sha256:398232e0… binds: current=true, dirty_record=false, tip 84096486. Tip equals 92bed0be for internal/.\n\nP8: closed. Rule at spec.md:373-375. Keyed identical-layout literals become producers through identicalKeys (rows.go:167-189, called at rows.go:151-153). Unkeyed ones become diagnostics at rows.go:102-108. writesStruct (rows.go:132) and literalStructs (resolve.go:311-334) cover the `~struct{...}` constraint terms. Tests at fieldwrites_test.go:74-77.\n\nP9: closed. Conversions are caught at rows.go:109-117. Nil and same-type operands are exempt (rows.go:115). Tests at fieldwrites_test.go:78-81.\n\nOC28-OC31 still hold. Moving SourceUses, firstUse and ImportName to loader.go:79-105 is a verbatim relocation.\n\nP10. The code is narrower than the rule for a field of a generic struct. fieldStruct returns the origin struct `struct{Reason T}` (resolve.go:277, 292-301), and structKey (resolve.go:265) gives an instantiation `row[string]` a different key. So `row[string]{\"c\"}` and `row[string](o)` no longer reach the diagnostic, and the census can read green. Disposition: fix. Confidence 7.\n\nP11. An implicit-assignability path to the struct leaves the census green. Example: `var s struct{Reason string}; s.Reason = \"c\"; var r row = s`. The selector write is to a different field object, and it is neither a literal nor a conversion. Disposition: fix, or a reviewer ruling. Confidence 6.\n\nAdvice: structKey treats unexported-field structs from different packages as one struct (resolve.go:260-267), stricter only toward the diagnostic.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "449b26b52db1864bedada8c26129df197b70b6e7",
          "finding_ids": [
            "P10",
            "P11"
          ],
          "supersedes": [
            "ft293-oc-c2-spec-r4"
          ]
        },
        {
          "id": "ft293-oc-c2-coverage-r5",
          "performer": "claude:ft293_oc-c2_coverage_r5",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "7c776f6d667e5158847a8281567695aa4cc0fdfa",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage_r5",
            "digest": "sha256:f780b9e2e41fc13fab409f9480fa9ead96063971e05c8e30556920ba2e25a05e",
            "excerpt": "## Result\nNo finding. The evidence binds: `--check-current` reports current=true, dirty_record=false, at tip 84096486.\n\nProbe row: omitted `add(e, viaConversion)` at `internal/consumers/rows.go:116`, ran `--package ./internal/preflight/... --run TestOwnerPremises`. verdict `bit`; failed, 3 tests (all three conversion rows went green); first_failure `TestOwnerPremises/a_conversion_from_another_struct_type`; restored `yes`.\n\nWrite-form table (`internal/preflight/ownership/fieldwrites_test.go:20-40`), by reading: identical-layout unkeyed rows bite under the author's `structKey` probe; identical-layout keyed rows bite by reasoning (not probed); conversion rows bite (probe); the nil and same-type row expects green, and dropping either guard at `rows.go:115` turns it unknown.\n\n## Advice\n- The `from.Type == nil` branch at `rows.go:115` has no test row; it looks unreachable for a type-checked conversion.\n- Splitting the nil and same-type row into two rows would make a red name its guard.\n\n## git status --short\n(empty)\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "449b26b52db1864bedada8c26129df197b70b6e7",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-coverage-r4"
          ]
        },
        {
          "id": "ft293-oc-c2-standards-r6",
          "performer": "claude:ft293_oc-c2_standards_r6",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r6",
            "digest": "sha256:7571b2e9756b496e93170270567b064920a915ec069b3ed3f5ad88c89d034d00",
            "excerpt": "## Result\nNo finding. The evidence binds: `current=true`, `dirty_record=false`, tip `75f73194`. There is no `internal/` change between `1395b49d` and the tip.\n\nThe removed path is fully gone. A search of the whole tree finds no remaining `identicalKeys`, `viaConversion`, `*ast.CallExpr` conversion branch, or `writesOnly` hook in `useRows`. The comments that described the old behavior were rewritten in present tense: `WriteRows` (rows.go:66-68), `UnnamedWrites` (rows.go:80-86), `writesStruct` (rows.go:170-171), `literalStructs` (resolve.go:305-306), `fieldStruct` (resolve.go:269-271), and `memberCensuses` (members.go:45-48). Struct identity still comes only from `structKey`, and `identicalType` uses it rather than a second comparison.\n\n## Advice\n- `identicalType` returns a predicate, but its doc says it \"reports whether\" (rows.go:131).\n- The generic-struct check at rows.go:102-105 does not depend on the package but runs inside the per-package loop.\n- \"The three Via values\" (rows.go:73) is a count that has to change whenever the constant block changes.\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "e382cdac4398cf16f0ef929b637fbb41dbc71edb",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-standards-r5"
          ]
        },
        {
          "id": "ft293-oc-c2-spec-r6",
          "performer": "claude:ft293_oc-c2_spec_r6",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec_r6",
            "digest": "sha256:c528f80370a7b8ee2363aabc62ce176d8d91284eefb4c2144332f7d80c88d2a8",
            "excerpt": "Evidence sha256:280f913a… was current (tip 75f7319, record clean). The internal diff from 1395b49 to the tip is empty, so I graded the committed cycle 5 delta.\n\nP10 is closed. `internal/consumers/rows.go:102-105` makes a generic named struct diagnostic at its declaration. Tests at `internal/preflight/ownership/fieldwrites_test.go:55-57`.\n\nP11 is closed. `identicalType` (`rows.go:134-168`), fed by `Info.Types` and `Info.Defs` (`rows.go:107-116`), reports the unnamed identical struct. Test at `fieldwrites_test.go:52`.\n\nKeyed writes of the field's own struct stay producers (`fieldwrites_test.go:40`). An unkeyed literal stays diagnostic (`rows.go:119-123`).\n\nP12 (blocking, rule sentence 3, in scope, confidence 8). A conversion from a struct type that differs only in field tags reads green. Go ignores tags in a conversion, so `type other struct{ Reason string `json:\"r\"` }; func renderC(o other) row { return row(o) }` writes `row.Reason`. But `structKey` includes tags (`types.TypeString`, `internal/consumers/resolve.go:265-267`), so `identicalType` (`rows.go:154`) never matches. Rule sentence 3 says \"a conversion from another type\" stays diagnostic.\n\nP13 (blocking, rule sentence 4, in scope, confidence 4). A generic alias, `type row[T any] = struct{ Reason T }`, skips the generic check: `own.Type()` is a `*types.Alias`, not a `*types.Named` (`rows.go:102`). An unkeyed `row[string]{\"c\"}` then misses the `structKey` match and reads green.\n\nOC28-OC31 hold by reading.\n\nAdvice: add a tag-differing conversion row to the `fieldwrites_test.go` table.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "e382cdac4398cf16f0ef929b637fbb41dbc71edb",
          "finding_ids": [
            "P12",
            "P13"
          ],
          "supersedes": [
            "ft293-oc-c2-spec-r5"
          ]
        },
        {
          "id": "ft293-oc-c2-coverage-r6",
          "performer": "claude:ft293_oc-c2_coverage_r6",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "c7979e3262180f27c9a94c7de4b040a41a8532b6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage_r6",
            "digest": "sha256:60283ab3a79a666160bde977a8a62c4039992c62db1a2a594a4ccf43d14407df",
            "excerpt": "Coverage round 6 for OC-C2: no finding. The evidence binds: `current=true`, `dirty_record=false`, tip 75f7319464ad.\n\nProbe: removed `add(own.Pos(), viaGeneric)` from `internal/consumers/rows.go:103`, ran `--package ./internal/preflight/... --run TestOwnerPremises`. Verdict `bit`, 3 tests failed, first_failure `TestOwnerPremises/a_generic_struct_with_no_added_write` (also `an_instantiation_conversion`, `an_unkeyed_instantiation_literal`); restored `yes`.\n\nCase grading (`internal/preflight/ownership/fieldwrites_test.go`): generic instantiation literal (line 56) and conversion (line 57) expect `unknown`; implicit assignment (line 52) and bare identical type (line 53) expect `identical`; plain keyed census of a non-generic struct (line 54) stays green with 2 members.\n\nAdvice: the generic branch at rows.go:104 returns inside the package loop and skips `sortRows`; harmless.\n\n`git status --short` printed nothing.\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "e382cdac4398cf16f0ef929b637fbb41dbc71edb",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-coverage-r5"
          ]
        },
        {
          "id": "ft293-oc-c2-standards-r7",
          "performer": "claude:ft293_oc-c2_standards_r7",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r7",
            "digest": "sha256:ac313e07cbb7c06c863329e565d091f75b4e7bf0ae796e0ba35e93e819f24bc2",
            "excerpt": "Standards round 7, OC-C2, cycle 6 delta (1395b49d..24279ede, internal). Evidence sha256:1b1e62c1… bound, `current=true`, tip ee1360c3.\n\nS7. `internal/preflight/ownership/fieldwrites_test.go:26-27` says \"Go ignores tags and alias spellings when it compares two structs, so a type that differs from the struct only in those is identical.\" That is false for tags. Go's type identity includes field tags; only conversion ignores them. `sameStruct`'s own comment at `internal/consumers/rows.go:140-141` has it right (\"Go converts between…\").\n- Binding requirement: `bench-craft-comments` \"What survives merge is only what's true of the code as it stands\", plus one source per fact.\n- Disposition: fix the wording (say Go converts between them), or drop the sentence and let `sameStruct` own the fact.\n- Confidence: 6.\n\nNo `structKey` caller or definition is left in `internal/`. `sameStruct` is the only source of struct identity at `rows.go:145`, used at `:172` and `:193`.\n\nAdvice:\n- `rows.go:143`: \"Either test can join two structs Go keeps apart\" is loose.\n- `rows.go:371` has its own check for generic named types; a different policy, optional to unify.\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "cd54d245ab7a1aa2fdf0ca4431125adb6264f205",
          "finding_ids": [
            "S7"
          ],
          "supersedes": [
            "ft293-oc-c2-standards-r6"
          ]
        },
        {
          "id": "ft293-oc-c2-spec-r7",
          "performer": "claude:ft293_oc-c2_spec_r7",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_spec_r7",
            "digest": "sha256:7763389ed9dddce4baf1cd6e01433933a35f47e0e3bf42d3290f4e1803dc2d98",
            "excerpt": "P12 and P13 are both closed. The evidence check bound (current=true, dirty_record=false, tip ee1360c3), and the only commit after 24279ede changes `reviews/preflight-ownership-closure.md`.\n\n- P12 (closed): `sameStruct` at `internal/consumers/rows.go:145-147` now uses `types.IdenticalIgnoreTags`. Both struct-identity checks call it: `identicalType` (`rows.go:172`) and `writesStruct` (`rows.go:193`). Fixture rows at `internal/preflight/ownership/fieldwrites_test.go:201-203`. `structKey` is gone.\n- P13 (closed): `genericDecl` at `rows.go:132-138` catches an `*types.Alias` that has type params. `fieldStruct` (`internal/consumers/resolve.go:286`) finds the struct through the alias's underlying type. Fixtures at `fieldwrites_test.go:130` and `:204-205`.\n\nResiduals\n- R1 (confirmed as the author states): the struct is declared in a test-variant load, differs in tags or nested aliases, and its field type is a named type with distinct objects across loads. Won't handle: a test-variant struct that differs in tags or alias spelling over a named field type — Go's identity test separates the two loads' objects, so the census can read green and needs manual review.\n- R2: `(*row)(unsafe.Pointer(&o))`, where `o` has a different layout. Won't handle: unsafe pointer conversions to the field's struct — the census sees no identical type, and the write needs manual review.\n- R3: reflection writes, such as `json.Unmarshal(b, &r)` or `reflect.Value.Field(i).SetString`. The existing reflection entry covers callers only. Won't handle: reflection writes to a printed field — no static form exists, and the write needs manual review.\n\nFindings: none. OC28-OC31 still hold.\n\nAdvice: the wider identity test could flip the real repo's OC31 preflight run to diagnostic.\n"
          },
          "axis": "Spec",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "cd54d245ab7a1aa2fdf0ca4431125adb6264f205",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-spec-r6"
          ]
        },
        {
          "id": "ft293-oc-c2-coverage-r7",
          "performer": "claude:ft293_oc-c2_coverage_r7",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_coverage_r7",
            "digest": "sha256:6c38318a1b4753b7e0e8f7199a82bc1610a47f46bdf8c2b43fe6375f30c6ea94",
            "excerpt": "The evidence binds (current=true, dirty_record=false, tip ee1360c3b), and the one probe went red and restored cleanly. No finding, one residual.\n\nProbe: in `internal/consumers/rows.go:133` swapped `ok && alias.TypeParams().Len() > 0 {` to `< 0`, disabling `genericDecl`'s alias branch; `--package ./internal/preflight/... --run TestOwnerPremises`.\n\n| verdict | first_failure | failed_tests | restored |\n|---|---|---|---|\n| bit | TestOwnerPremises/a_generic_alias_with_no_added_write | 2 | yes |\n\nRow grades (`internal/preflight/ownership/fieldwrites_test.go`): tag-differing (:62) needs `IdenticalIgnoreTags` (author probe bit); alias-spelled field (:63) graded by reading; alias of the struct (:64) harmless; two generic-alias rows (:65-66) bite under my probe.\n\nFindings: none.\n\nResiduals: R1: the `TypeString` disjunct in `sameStruct` (`rows.go:146`) has no row with a named field type loaded both plainly and as a test variant. Confidence 6; not probed.\n\nAdvice: the `Unalias`→`*types.Named` arm (`rows.go:136`) also counts a non-generic alias of an instance as generic; untested but only widens toward the diagnostic.\n\n`git status --short`: empty.\n"
          },
          "axis": "Coverage",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "cd54d245ab7a1aa2fdf0ca4431125adb6264f205",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-coverage-r6"
          ]
        },
        {
          "id": "ft293-oc-c2-standards-r8",
          "performer": "claude:ft293_oc-c2_standards_r7",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "27d3621dd5a6cbded15bb40405c808669fc74b1b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_oc-c2_standards_r7",
            "digest": "sha256:ff726e6602eb37cb2c73d9281967675750f08fc4c2cac544481b8997aee2e0bd",
            "excerpt": "Reaffirmation after the evidence-only correction at fed2ca8bf.\n\nS7 is closed. In commit fed2ca8bf, `internal/preflight/ownership/fieldwrites_test.go:20-21` now says: \"A type that differs from the struct only in tags or alias spellings counts as identical, because Go converts between the two.\" It no longer claims that Go's identity rule ignores tags, and the reason it gives matches `sameStruct` (`internal/consumers/rows.go:140-141`).\n\nThe commit changes no code in `internal/`. That file has one hunk (-2/+2), and only comment lines change. The commit touches one other file, `reviews/preflight-ownership-closure.md`, a review record.\n\nNo other Standards finding of round 7 changes.\n"
          },
          "axis": "Standards",
          "base": "10c458e5cd49c50d8c3c30e7cfd13c353569d266",
          "tip": "cd54d245ab7a1aa2fdf0ca4431125adb6264f205",
          "finding_ids": [],
          "supersedes": [
            "ft293-oc-c2-standards-r7"
          ]
        }
      ]
    },
    {
      "id": "OC-C3",
      "base": "cd54d245ab7a1aa2fdf0ca4431125adb6264f205",
      "tip": "2b9d9f49108cbca54cbbb930818c01d82910f382",
      "plan_digest": "sha256:fa4a88a7ab799bed77ecdce96efa188b078dbdacb71b33a4b4d540237f415c27",
      "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
      "acceptance_rows": [
        "OC32",
        "OC33",
        "OC34",
        "OC35",
        "OC36",
        "OC37",
        "OC38",
        "OC39",
        "OC40",
        "OC41",
        "OC42",
        "OC43",
        "OC44",
        "OC45",
        "OC46",
        "OC47",
        "OC48",
        "OC49",
        "OC50",
        "OC87",
        "OC51",
        "OC52",
        "OC53",
        "OC54",
        "OC88",
        "OC89",
        "OC93",
        "OC94",
        "OC95"
      ],
      "verification": [
        {
          "id": "ft293_t07-oc-07-focused",
          "performer": "agy:ft293_t07",
          "role": "author-verification",
          "model": "gemini-4-argon-medium",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agy:ft293_t07",
            "digest": "sha256:5dad21a491c016c533e4120ca033af0263a3359767459390839f3b070a820c2e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/ownership/premises.go,swap,failed,2,TestChargePremises/OC88_omitted_second_seam_citation,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight/...,TestChargeManualApplicability|TestChargePremises|TestPredicatePremises|TestHelperCallers,passed,31\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,904\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,6701,31\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestChargePremises/OC88_omitted_second_seam_citation,\"premises_test.go:342: charge succeeded, want refusal:\",3\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestPredicatePremises/OC88_omitted_second_seam_citation,\"premises_test.go:305: premises rows matching map[id:metadata] = 0, want 1:\",8\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-07-focused",
          "command": "bench test --package ./internal/preflight/... --run 'TestChargeManualApplicability|TestChargePremises|TestPredicatePremises|TestHelperCallers'",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-07-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "agy:ft293_t07",
              "digest": "sha256:5dad21a491c016c533e4120ca033af0263a3359767459390839f3b070a820c2e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/preflight/ownership/premises.go,swap,failed,2,TestChargePremises/OC88_omitted_second_seam_citation,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/preflight/...,TestChargeManualApplicability|TestChargePremises|TestPredicatePremises|TestHelperCallers,passed,31\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,904\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,7,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,6701,31\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestChargePremises/OC88_omitted_second_seam_citation,\"premises_test.go:342: charge succeeded, want refusal:\",3\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestPredicatePremises/OC88_omitted_second_seam_citation,\"premises_test.go:305: premises rows matching map[id:metadata] = 0, want 1:\",8\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t07-oc-07-check-1",
          "performer": "agy:ft293_t07",
          "role": "author-verification",
          "model": "gemini-4-argon-medium",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agy:ft293_t07",
            "digest": "sha256:121b34d5e69586176bba5e079ed3a2a4cd3eff6846867ab3f92fee770408c8c5",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,36527,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23301,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,35062,225\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-07-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t07-oc-07-check-2",
          "performer": "agy:ft293_t07",
          "role": "author-verification",
          "model": "gemini-4-argon-medium",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agy:ft293_t07",
            "digest": "sha256:016ca0756faa4e838fa97a11a94f1d0e5b9deb97b405f19aa5148517f45fb4ad",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/coverage,pass,1754,145\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-07-check-2",
          "command": "bench test --package ./internal/coverage",
          "exit_code": 0
        },
        {
          "id": "ft293_t07-oc-07-check-3",
          "performer": "agy:ft293_t07",
          "role": "author-verification",
          "model": "gemini-4-argon-medium",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agy:ft293_t07",
            "digest": "sha256:2e354af432285a621d82b7894c104795d4f1c94a98945a5e9f63b311a0c5f823",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,2759,110\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-07-check-3",
          "command": "bench test --package ./internal/reviewrecord",
          "exit_code": 0
        },
        {
          "id": "ft293_t07-oc-07-check-4",
          "performer": "agy:ft293_t07",
          "role": "author-verification",
          "model": "gemini-4-argon-medium",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agy:ft293_t07",
            "digest": "sha256:0e43ecbb80d0d28f03f985c7b27dd61be591641578efffed8974c1254ebafcbc",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15655,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-07-check-4",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t07-oc-07-check-5",
          "performer": "agy:ft293_t07",
          "role": "author-verification",
          "model": "gemini-4-argon-medium",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "agy:ft293_t07",
            "digest": "sha256:d0d134bf56c787b506d0819c2c130e150eb11aa2fb894f639b7ed17e750e738f",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,15264,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-07-check-5",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ft293_t04a-oc-04-focused",
          "performer": "claude:ft293_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t04",
            "digest": "sha256:bf5b272cd3d1ec639ad0d749da8504e8240615b26d62c77e3d65b8cf291118dd",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestSizingPremises/OC33_glob-complete,21,1,0,null,yes\n  2,bit,swap,failed,2,TestSizingPremises/OC37_headroom-ignored-deletion,21,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSizingPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2555\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,134,21\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC33_glob-complete,\"premises_test.go:155: premise glob-complete = map[id:glob-complete reason:glob-members premise: the member list names internal/glob/.hidden.go, which is no member source:3dff62e008ac8bc5bedd17f611366b4499f17dba state:unknown], want green an… (279 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC33_glob-visible-only,\"premises_test.go:155: premise glob-visible-only = map[id:glob-visible-only reason:the glob internal/glob/*.go holds 3 tracked member(s) at the source tip 3dff62e008ac8bc5bedd17f611366b4499f17dba source:3dff62e008ac8bc5bedd17f611366b4499f17d… (322 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,133,21\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC37_headroom-ignored-deletion,\"premises_test.go:155: premise headroom-ignored-deletion = map[id:headroom-ignored-deletion reason:headroom premise: the plan of net 1 exceeds the limit 3 of internal/sized, which no grant exempts source:c6c1aad19576b30900af526d1a8f89709b3d1… (317 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC37_headroom-swap,\"premises_test.go:155: premise headroom-swap = map[id:headroom-swap reason:headroom premise: expected net 0, the plan nets 1 source:c6c1aad19576b30900af526d1a8f89709b3d1781 state:unknown], want green and a reason naming \\\"net 0\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-04-focused",
          "command": "bench test --package ./internal/preflight/... --run TestSizingPremises",
          "exit_code": 0,
          "probe": {
            "mutation": "oc-04-omission",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft293_t04",
              "digest": "sha256:bf5b272cd3d1ec639ad0d749da8504e8240615b26d62c77e3d65b8cf291118dd",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,2,TestSizingPremises/OC33_glob-complete,21,1,0,null,yes\n  2,bit,swap,failed,2,TestSizingPremises/OC37_headroom-ignored-deletion,21,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/preflight/...,TestSizingPremises,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,2555\nmutant_report[1]{mutant}:\n  1\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,3,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,134,21\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC33_glob-complete,\"premises_test.go:155: premise glob-complete = map[id:glob-complete reason:glob-members premise: the member list names internal/glob/.hidden.go, which is no member source:3dff62e008ac8bc5bedd17f611366b4499f17dba state:unknown], want green an… (279 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC33_glob-visible-only,\"premises_test.go:155: premise glob-visible-only = map[id:glob-visible-only reason:the glob internal/glob/*.go holds 3 tracked member(s) at the source tip 3dff62e008ac8bc5bedd17f611366b4499f17dba source:3dff62e008ac8bc5bedd17f611366b4499f17d… (322 bytes)\",1\nskips[0]{package,test,reason}:\nmutant_report[1]{mutant}:\n  2\npackages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,4,0\n  github.com/gibbonmi/bench/internal/preflight/ownership,fail,133,21\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC37_headroom-ignored-deletion,\"premises_test.go:155: premise headroom-ignored-deletion = map[id:headroom-ignored-deletion reason:headroom premise: the plan of net 1 exceeds the limit 3 of internal/sized, which no grant exempts source:c6c1aad19576b30900af526d1a8f89709b3d1… (317 bytes)\",1\n  github.com/gibbonmi/bench/internal/preflight/ownership,TestSizingPremises/OC37_headroom-swap,\"premises_test.go:155: premise headroom-swap = map[id:headroom-swap reason:headroom premise: expected net 0, the plan nets 1 source:c6c1aad19576b30900af526d1a8f89709b3d1781 state:unknown], want green and a reason naming \\\"net 0\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        },
        {
          "id": "ft293_t04a-oc-04-check-1",
          "performer": "claude:ft293_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t04",
            "digest": "sha256:27314c1c6013717324d1061472f6bf8c2c07763463055b1b05983d9c0b13f615",
            "excerpt": "packages[6]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/preflight,pass,36240,408\n  github.com/gibbonmi/bench/internal/preflight/chargesource,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/evidencecmd,pass,23223,232\n  github.com/gibbonmi/bench/internal/preflight/ownership,pass,35475,225\n  github.com/gibbonmi/bench/internal/preflight/ownership/ownershiptest,no-tests,0,0\n  github.com/gibbonmi/bench/internal/preflight/preflighttest,no-tests,0,0\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-04-check-1",
          "command": "bench test --package ./internal/preflight/...",
          "exit_code": 0
        },
        {
          "id": "ft293_t04a-oc-04-check-2",
          "performer": "claude:ft293_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t04",
            "digest": "sha256:844502e9e700df12ec096bd937d360ee3f346c8b9ad7a6cb3ae97267dc24769c",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/structure,pass,250,34\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-04-check-2",
          "command": "bench test --package ./internal/structure",
          "exit_code": 0
        },
        {
          "id": "ft293_t04a-oc-04-check-3",
          "performer": "claude:ft293_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t04",
            "digest": "sha256:d398a68e3ffe5314e742a13c133aa531fc43d534078acd7cadfc9b77dd05b950",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15245,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-04-check-3",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "ft293_t04a-oc-04-check-4",
          "performer": "claude:ft293_t04",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "e034437d316941fd9fa5cc132b97c85eb58fda48",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft293_t04",
            "digest": "sha256:1836bacdb3eb9f19c1b672c5a9a625f664b4e8a27837191989da49e2c8a50f54",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/cmd/bench,pass,14586,346\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "requirement": "oc-04-check-4",
          "command": "bench test --package ./cmd/bench",
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
      "from": "sha256:9a5a573300d6ec1ac2977c518eb468f936fdbcf261e48fa72be1b495dcf98b66",
      "to": "sha256:9e4941f4243f14167ae9a17c2f09edf8d5b4ae67633d4f0d597fa7017edfb0cb",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ]
      }
    },
    {
      "from": "sha256:9e4941f4243f14167ae9a17c2f09edf8d5b4ae67633d4f0d597fa7017edfb0cb",
      "to": "sha256:033866f64e250d307058897dad79a90c996dcbe1e261a89b6bc8a824a18ceac9",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ]
      }
    },
    {
      "from": "sha256:033866f64e250d307058897dad79a90c996dcbe1e261a89b6bc8a824a18ceac9",
      "to": "sha256:21c3a6ea30e0fd4aaa576ae906a28c2b2901795ad361b30ddfc0acfad95a321e",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
        ]
      }
    },
    {
      "from": "sha256:21c3a6ea30e0fd4aaa576ae906a28c2b2901795ad361b30ddfc0acfad95a321e",
      "to": "sha256:ce79cb1b86b2675c6faf4835cfc20371b09a870554e42d374afeffd7625664da",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
        ]
      }
    },
    {
      "from": "sha256:ce79cb1b86b2675c6faf4835cfc20371b09a870554e42d374afeffd7625664da",
      "to": "sha256:dfce7ae2a24bff33b3a95fc94408494cdc5c7953056cc06770d02e7b96d93918",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
        ]
      }
    },
    {
      "from": "sha256:dfce7ae2a24bff33b3a95fc94408494cdc5c7953056cc06770d02e7b96d93918",
      "to": "sha256:80d6dcf3d05276b3e84cdbc12395aebb95bb9cdbc1ca5beddddf08a208043058",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
        ]
      }
    },
    {
      "from": "sha256:80d6dcf3d05276b3e84cdbc12395aebb95bb9cdbc1ca5beddddf08a208043058",
      "to": "sha256:78434ecbdf9c9baaefedd9eb9ffc55b6bf2177c716f531deb6f6ae0296e98ab3",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
        ]
      }
    },
    {
      "from": "sha256:78434ecbdf9c9baaefedd9eb9ffc55b6bf2177c716f531deb6f6ae0296e98ab3",
      "to": "sha256:9de99bb1f6e996574324f2d278b8357129972bb0143631cd3ffa66f0300512e3",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
        ]
      }
    },
    {
      "from": "sha256:9de99bb1f6e996574324f2d278b8357129972bb0143631cd3ffa66f0300512e3",
      "to": "sha256:fa4a88a7ab799bed77ecdce96efa188b078dbdacb71b33a4b4d540237f415c27",
      "chunk_ids": {
        "OC-C1": [
          "OC-C1"
        ],
        "OC-C2": [
          "OC-C2"
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

## OC-C2 repair state

OC-C2: repair cycle 1 of 2 is complete (ticket 02 at 4268cb7c, ticket 03 at 76cbf471). Its confirming round closed every first-round finding and found S6, P8, and C6.
OC-C2: repair cycle 2 of 2 is complete at 7cf83f37. Its confirming round closed S6 and C6, and P8 stays open for a pointer-elided literal.

OC-C2: the reviewer's auto-approve direction extends the allowance by one cycle. Repair cycle 3 opens for P8 only, in one ticket 03 repair session.
An Opus/high consultant decided each ask-user finding under the auto-approve direction; each decision is open to reviewer veto.

## OC-C2 Standards

OC-C2: 5 findings. The worst issue is S1. Ticket 03's repair takes all five, because its `Writes:` line holds every cited path.

- S1 (auto-fix, confidence 8): `internal/preflight/ownership/producers.go:92-134` derives declaration names and receiver names again beside `internal/consumers/rows.go` `declNames`.
- S2 (auto-fix, confidence 7): the import local-name rule appears at `producers.go:84-90`, `producers.go:157-165`, and `rows.go:211-219`; `producers.go:159` matches import suffixes beside `resolve.go:148-154`.
- S3 (auto-fix, confidence 6): `helpers.go:153-157` `localName` reads the symbol grammar a second time and disagrees with `consumers.Resolve`.
- S4 (auto-fix, confidence 6): the member-census kinds appear twice, at `helpers.go:184` and `helpers.go:266-268`.
- S5 (auto-fix, confidence 5): `internal/preflight/proposal.go:149-158` keeps a second per-kind table beside `closureFamily`.

## OC-C2 Spec

OC-C2: 6 actionable findings. The worst issue is P1.

- P1 (auto-fix after consultation, confidence 7, ticket 02): `producers.go:343-352` skips a selector on an in-module import with no edge and no diagnostic. The repair resolves it through the import table, roots it in the target package, and keeps the walk out of function literals. Flagged: the walk stays syntactic, because the typed loader reads only the tip.
- P2 (auto-fix after consultation, confidence 5, ticket 02): `effects.go:77` skips a changed production package with no Grammar row; every changed production Go package becomes classifiable.
- P3 (auto-fix after consultation, confidence 7, ticket 03): posture and printed-field items take the table shapes with exact key sets. `writer_probe` decodes as a reused-test subject. Flagged: the key spellings `callers` and `producers` follow the retirement precedent.
- P4 (auto-fix after consultation, confidence 6, ticket 03): a printed-field producer is a static write only, not a read.
- P5 (auto-fix after consultation, confidence 6, ticket 03): an owner requirement names each originating write entry in `closure.path`.
- P7 (auto-fix after consultation, confidence 4, ticket 03): each existing liveness change grades its dispositions against the retirement census. Flagged: the allowed actions follow the manifest table.
- P6 (no-op, confidence 7): spec 241-242 and the change rows at 225-226 disagree; the code follows the table. Flagged for reviewer veto.

## OC-C2 Coverage

OC-C2: 5 findings. The worst issue is C1.

- C1 (auto-fix, confidence 8, ticket 03): no light-path or record test uses a `(deleted)` entry; a probe mutant at `light_path.go` stayed silent.
- C2 (auto-fix, confidence 5, ticket 03): no `fence-writes` or authorization test needs the `(deleted)` strip.
- C3 (auto-fix, confidence 8, ticket 02): no test reaches the ambiguous-declaration diagnostic or malformed Go; a probe mutant at `producers.go` stayed silent.
- C4 (merged into P1, ticket 02): a cross-package grammar dependency adds no edge and no diagnostic.
- C5 (auto-fix, confidence 6, ticket 03): no test changes a method symbol whose former caller stops using it.

## OC-C2 confirming round

- S6 (auto-fix, confidence 4): `ownership/effects.go:30-31` `producerRead` repeats the shape of the new `revRead[T]` at `ownership/imports.go:56-60`.
- P8 (auto-fix, confidence 5): `internal/consumers/rows.go:294-305` misses an unkeyed struct literal write, so a printed-field premise can read green without complete proof.
- C6 (auto-fix, confidence 5): no test reaches the cross-package call edge at `producers.go:361`, so a mutant `called: false` would stay silent.

## OC-C2 second confirming round

- P8 (auto-fix, confidence 7): `internal/consumers/rows.go:96` asserts a struct type before it removes a pointer, so a pointer-elided literal such as `[]*row{{"c"}}` escapes the unkeyed-write diagnostic.

## OC-C2 third confirming round

Repair cycle 3 closed every P8 form that it listed. The round found two more: an identical-layout struct constraint (P8) and a conversion from an identical struct type (P9).
The reviewer chose the fail-closed route in this session. Repair cycle 4 matches fields by struct identity and makes any other path to the struct diagnostic.
The spec gains one rule in the helper and owner closure section that defines the static write forms.

- P8 (auto-fix, confidence 7): `internal/consumers/rows.go:96-98` matches by declaration position, so an anonymous identical-layout struct escapes.
- P9 (auto-fix, confidence 5): a conversion such as `row(other{"c"})` writes the field with no census match.

## OC-C2 fourth confirming round

Repair cycle 4 closed P8 and P9, but it added conversions to a list instead of closing the class. The round found P10 and P11.
Repair cycle 5 closes the class by type presence, under the reviewer's fail-closed decision and the auto-approve direction. The spec rule gains one sentence for the generic and identical-type cases.

- P10 (auto-fix, confidence 7): `internal/consumers/resolve.go:265` keys an instantiation of a generic struct apart from its origin, so `row[string]{"c"}` escapes.
- P11 (auto-fix, confidence 6): an implicit assignment from an unnamed identical struct type, such as `var r row = s`, leaves the census green.

## OC-C2 fifth confirming round

Repair cycle 5 closed P10 and P11 by type presence. The round found two defects in the identity test itself, and the reviewer chose repair cycle 6 for both.

- P12 (auto-fix, confidence 8): `internal/consumers/resolve.go:265-267` keys a struct with its field tags, but Go ignores tags in a conversion, so a tag-differing conversion reads green.
- P13 (auto-fix, confidence 4): `internal/consumers/rows.go:102` misses a generic alias, because its type is a `*types.Alias`.

## OC-C2 sixth confirming round

Repair cycle 6 closed P12 and P13, and the round found no repair finding. The orchestrator fixed S7, a false test comment, as a comment-only evidence correction with no plan assignment.

Three residuals of the printed-field census get Won't-handle lines in the chunk-gap plan commit:

- R1: a test-variant struct that differs in tags or alias spelling over a named field type.
- R2: an unsafe pointer conversion to the field's struct.
- R3: a reflection write to a printed field.

OC-C2: repair cycles are closed at six. No reviewer extension remains open.
