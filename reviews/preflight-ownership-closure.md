# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/preflight-ownership-closure/spec.md",
  "plan_digest": "sha256:9a5a573300d6ec1ac2977c518eb468f936fdbcf261e48fa72be1b495dcf98b66",
  "implementation_session": "",
  "chunks": [
    {
      "id": "OC-C1",
      "base": "1a8fece7d1173873ce9a32955b4b1d9a25600a4d",
      "tip": "a54b85b3cadc7c89241c5a950a90ae69765ad177",
      "plan_digest": "sha256:9a5a573300d6ec1ac2977c518eb468f936fdbcf261e48fa72be1b495dcf98b66",
      "source_digest": "69a2ca4dcac7f1c908a79f785d4b288fc32b4d48",
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
