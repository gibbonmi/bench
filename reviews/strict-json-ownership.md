# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/strict-json-ownership/spec.md",
  "plan_digest": "sha256:352f91ffc5cbebf5221070c81e6c1f2898a87851ee1d00dd37ec1ab08a479802",
  "implementation_session": "",
  "chunks": [
    {
      "id": "SJ-C1",
      "base": "835f5491b7aa50042f5f65bcb6ac941592e8181b",
      "tip": "e32079f005169589a02dbf3826fde9250a988aff",
      "plan_digest": "sha256:f12fdb6d19f0b2ae0d0689ed9dd450dd26628d69f7f864a711955a0cf02e43ff",
      "source_digest": "3304f73fea53f227fdb311af4751cde912612bfc",
      "acceptance_rows": [
        "SJ1",
        "SJ2",
        "SJ3",
        "SJ4",
        "SJ5",
        "SJ6",
        "SJ7",
        "SJ8",
        "SJ9",
        "SJ10",
        "SJ11",
        "SJ12",
        "SJ13",
        "SJ14",
        "SJ15"
      ],
      "verification": [
        {
          "id": "sj-c1-policy",
          "performer": "ft354-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3304f73fea53f227fdb311af4751cde912612bfc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft354-author-01",
            "digest": "sha256:ac4ab6210898d21f5921f0c3f8f96b8d330d4d049215980123d725dac0673b4e",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/releasepreflight/vulnerability.go,swap,failed,14,TestVulnerabilityPolicyDocumentContract/empty_array_then_array,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/releasepreflight,TestVulnerabilityPolicyDocumentContract,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,3706\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/releasepreflight,fail,3,16\nfailures[14]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_array,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_boolean,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_malformed_suffix,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_null,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_number,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_object,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_string,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/extra_typed_field,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/repeated_expires_key,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/repeated_id_key,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/repeated_reason_key,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/used_array_then_malformed_suffix,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/used_array_then_number,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/used_array_then_object,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "e32079f005169589a02dbf3826fde9250a988aff",
          "requirement": "policy",
          "command": "bench test --package ./internal/releasepreflight --package ./internal/jsonfile",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the policy bytes for an empty array at the jsonfile.DecodeDocument call in ValidateVulnerabilityPolicy, and record that the suffix case of TestVulnerabilityPolicyDocumentContract turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft354-author-01",
              "digest": "sha256:ac4ab6210898d21f5921f0c3f8f96b8d330d4d049215980123d725dac0673b4e",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/releasepreflight/vulnerability.go,swap,failed,14,TestVulnerabilityPolicyDocumentContract/empty_array_then_array,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/releasepreflight,TestVulnerabilityPolicyDocumentContract,passed,16\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,3706\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/releasepreflight,fail,3,16\nfailures[14]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_array,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_boolean,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_malformed_suffix,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_null,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_number,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_object,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/empty_array_then_string,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = \u003cnil>, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/extra_typed_field,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/repeated_expires_key,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/repeated_id_key,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/repeated_reason_key,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/used_array_then_malformed_suffix,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/used_array_then_number,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\n  github.com/gibbonmi/bench/internal/releasepreflight,TestVulnerabilityPolicyDocumentContract/used_array_then_object,\"vulnerability_test.go:50: ValidateVulnerabilityPolicy error = uncovered vulnerabilities: GO-1, want containing \\\"exception file is malformed\\\"\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "sj-c1-standards-1",
          "performer": "ft354-review-c1-standards-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3304f73fea53f227fdb311af4751cde912612bfc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c1-standards-1",
            "digest": "sha256:2e94b9ce545b7f7d89df0f05fdec955111c481f7c910bba78531105d6e7a3bb1",
            "excerpt": "## Standards\n0 findings. No issue reaches finding level. The worst candidate is A1, which is advice only.\nDuplicated code: the delta removes a second hand-written strict decoder and calls the shared owner at internal/releasepreflight/vulnerability.go:90.\nTest expectations: the refusal strings in vulnerability_test.go:80-115 check caller error output that SJ1-SJ15 require; they copy no implementation logic.\nComments: one added comment at vulnerability_test.go:37 passes craft-comments, apart from wording point A2.\nAdvice:\n- A1. vulnerability.go:87 keeps its own zero-byte check before DecodeDocument; it owns the SJ8 domain error. An optional why-comment could protect it.\n- A2. The comment at vulnerability_test.go:37 says \"requires that refusal text\", but line 70 checks containment.\n- A3. spec.md:358 states history (\"The build-entry staleness pass refreshed ...\"); state the current state instead.\n- A4. The SJ-C5 canary list in the spec repeats the ticket 5 Writes: line.\n- A5. ValidateVulnerabilityPolicy (vulnerability.go:84) has no doc comment; old code.\nImplementation command: did not contribute; no change is necessary.\nSteps 1-7 done.\n"
          },
          "axis": "Standards",
          "base": "835f5491b7aa50042f5f65bcb6ac941592e8181b",
          "tip": "e32079f005169589a02dbf3826fde9250a988aff",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "sj-c1-spec-1",
          "performer": "ft354-review-c1-spec-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3304f73fea53f227fdb311af4751cde912612bfc",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c1-spec-1",
            "digest": "sha256:c7af86e3ccaa43d713e174d0eb4db116afed8d092d5062b9b50632b14289d436",
            "excerpt": "## Spec\n1 finding. Worst issue (P1): rows SJ5 and SJ7-SJ15 have no recorded red; the spec and the ticket require one for each row.\n- P1, severity 5, confidence 8, auto-fix, rows SJ5 and SJ7-SJ15. spec.md:145 \"Each acceptance row needs its own red-to-green evidence.\"; tickets/1-policy-document.md:42 \"Record one red and one green for each covered row.\"; reviews/strict-json-ownership.md:89 records no red for these rows. The plan probe leaves SJ5's whitespace case green, and the preservation rows need domain mutations (for example \u003c to \u003c= at vulnerability.go:110 for SJ11). Repair: run and record a red for each row; no code changes.\nSJ1-SJ15 map to cases at vulnerability_test.go:81-95 and 102-116; all 15 seam cells cite existing tests (lines 77, 99); vulnerability.go:90 uses jsonfile.DecodeDocument; findingIDs unchanged; the trailing message change is acceptable.\nAdvice: the plan narrowing of census verification to final verification needs a bench learning entry; the SJ-C5 canary narrowing is not explained in the spec; SJ2/SJ3 check only the malformed prefix.\nImplementation command contributed: the ticket's Mutation proof names one probe that catches only refusal rows; name a probe for SJ5 and one domain mutation for each SJ7-SJ15 row.\nSteps 1-8 done.\n"
          },
          "axis": "Spec",
          "base": "835f5491b7aa50042f5f65bcb6ac941592e8181b",
          "tip": "e32079f005169589a02dbf3826fde9250a988aff",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "sj-c1-coverage-1",
          "performer": "ft354-review-c1-coverage-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3304f73fea53f227fdb311af4751cde912612bfc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c1-coverage-1",
            "digest": "sha256:e2a824c0f15e08aad5f409ae17df5d906295ea7c8bcf67bd1882ceb31fd12336",
            "excerpt": "## Coverage\n0 findings. The only open evidence item was SJ5 with no recorded red, and the axis probe bit it this round, so it closes as evidence-only.\nWrite set: only internal/releasepreflight/vulnerability.go and vulnerability_test.go change, inside the SJ-C1 fence. The decode goes through jsonfile.DecodeDocument (vulnerability.go:90).\nTraced mutations: SJ4 duplicate check, SJ6 permissive decoder, SJ7/SJ8 guards, SJ11 \u003c to \u003c=, SJ13 seen check, SJ15 three guards each go red.\nEdge inventory (spec 243-244) is covered at vulnerability_test.go:60-69. Hostile-input classes: SJ5 no trailing newline; SJ7/SJ8 absent vs present-but-empty.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/releasepreflight/vulnerability.go,swap,failed,1,TestVulnerabilityPolicyDocumentContract/surrounding_whitespace_without_final_newline,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/releasepreflight,TestVulnerabilityPolicyDocumentContract,passed,16\nAdvice: no test covers a present whitespace-only file; both decoders refuse it, so it is a regression pin and not a defect.\nImplementation command: contributed to no issue.\nSteps 1-8 done.\n"
          },
          "axis": "Coverage",
          "base": "835f5491b7aa50042f5f65bcb6ac941592e8181b",
          "tip": "e32079f005169589a02dbf3826fde9250a988aff",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "sj-c1-spec-2",
          "performer": "ft354-review-c1-spec-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "3304f73fea53f227fdb311af4751cde912612bfc",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c1-spec-2",
            "digest": "sha256:df2a9c9580df682c7c01909bcb8111fc01371bf32267271078f2b554882ac332",
            "excerpt": "## Spec\nReaffirmation of P1 by ft354-review-c1-spec-2. P1 closes for SJ5, SJ8, and SJ10-SJ15: each row's recorded first-failing case is that row's own case.\nResidual P1-r, severity 3, confidence 7: for SJ7 and SJ9 the recorded first failure belonged to SJ10 cases (reviews/strict-json-ownership.md:164,166); repair is to record the SJ7 and SJ9 cases as red, no code change.\nThe internal/ tree is the same at 1b12ee8b and e32079f0, so the mutants ran on the frozen code.\nSteps 1-4 done.\nAdvice (bounded repair policy: an initial-round finding at severity 3 or lower that is not a correctness defect): P1-r. The coordinator closed it with a narrowed mutant set at b68631c1:\nmutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestVulnerabilityPolicyDomainCompatibility/absent_policy_without_findings,4,1,0,null,yes\n  2,bit,swap,failed,2,TestVulnerabilityPolicyDomainCompatibility/empty_array_without_findings,4,1,0,null,yes\n"
          },
          "axis": "Spec",
          "base": "835f5491b7aa50042f5f65bcb6ac941592e8181b",
          "tip": "e32079f005169589a02dbf3826fde9250a988aff",
          "finding_ids": [],
          "supersedes": [
            "sj-c1-spec-1"
          ]
        }
      ]
    },
    {
      "id": "SJ-C2",
      "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
      "tip": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
      "plan_digest": "sha256:4eb2f514d61a626059587c5896e808e64450264ba12146938389c4f1e1673f57",
      "source_digest": "c3b2a959e0f36fdafa82ade2a1e1a7e3e7967248",
      "acceptance_rows": [
        "SJ16",
        "SJ17",
        "SJ18",
        "SJ19",
        "SJ20",
        "SJ21",
        "SJ22",
        "SJ23",
        "SJ24",
        "SJ25",
        "SJ26",
        "SJ27",
        "SJ28",
        "SJ29",
        "SJ30",
        "SJ31"
      ],
      "verification": [
        {
          "id": "sj-c2-owner-gate",
          "performer": "ft354-author-02",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "cee4624207f1395d09c8a38efa3288218981b4ec",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft354-author-02",
            "digest": "sha256:cfbc7e954a1df4afa108f3b521c91bf1083e6abb75e6365febcf6318de6174fd",
            "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDocumentStructuralContract/unknown_nested_field,15,1,0,null,yes\n  2,bit,swap,failed,1,TestDocumentStructuralContract/whitespace_after_the_value,15,1,0,null,yes\n  3,bit,swap,failed,2,TestDocumentStructuralContract/nested_duplicate_in_an_array_element,15,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/jsonfile,TestDocumentStructuralContract,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "4f88b2d46289a69c396185d1a40de5f0f6182326",
          "requirement": "owner-gate",
          "command": "bench test --package ./internal/jsonfile --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap 'seen[name] = true' for 'seen[name] = false' in internal/jsonfile/decode.go, and record that the nested duplicate case of TestDocumentStructuralContract turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft354-author-02",
              "digest": "sha256:cfbc7e954a1df4afa108f3b521c91bf1083e6abb75e6365febcf6318de6174fd",
              "excerpt": "mutants[3]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestDocumentStructuralContract/unknown_nested_field,15,1,0,null,yes\n  2,bit,swap,failed,1,TestDocumentStructuralContract/whitespace_after_the_value,15,1,0,null,yes\n  3,bit,swap,failed,2,TestDocumentStructuralContract/nested_duplicate_in_an_array_element,15,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/jsonfile,TestDocumentStructuralContract,passed,4\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "sj-c2-r1-owner-gate",
          "performer": "ft354-repair-02-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "c3b2a959e0f36fdafa82ade2a1e1a7e3e7967248",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:ft354-repair-02-1",
            "digest": "sha256:aa9b68fe26da7312f4967c1b8801ebe18c8e8172699cfa647e002dee6d29c700",
            "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap+swap+swap,failed,1,TestDocumentStructuralContract/repeated_key_with_an_escaped_spelling,15,1,0,null,yes\n  2,bit,swap,failed,2,TestDocumentStructuralContract/nested_duplicate_in_an_array_element,15,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/jsonfile,TestDocumentStructuralContract,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "requirement": "owner-gate",
          "command": "bench test --package ./internal/jsonfile --package ./internal/gate",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap 'seen[name] = true' for 'seen[name] = false' in internal/jsonfile/decode.go, and record that the nested duplicate case of TestDocumentStructuralContract turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-high:ft354-repair-02-1",
              "digest": "sha256:aa9b68fe26da7312f4967c1b8801ebe18c8e8172699cfa647e002dee6d29c700",
              "excerpt": "mutants[2]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap+swap+swap+swap,failed,1,TestDocumentStructuralContract/repeated_key_with_an_escaped_spelling,15,1,0,null,yes\n  2,bit,swap,failed,2,TestDocumentStructuralContract/nested_duplicate_in_an_array_element,15,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/jsonfile,TestDocumentStructuralContract,passed,3\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "sj-c2-standards-1",
          "performer": "ft354-review-c2-standards-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "cee4624207f1395d09c8a38efa3288218981b4ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c2-standards-1",
            "digest": "sha256:db178731cee7a4a46d57cea88f7290914f890ec8e7588c5522f75dbecfd41373",
            "excerpt": "## Standards\n1 finding. Worst issue: a test case name claims a fixture the case does not contain.\n- S1 | sev 5 | conf 9 | auto-fix | internal/jsonfile/decode_test.go:60 | Mysterious Name against SJ17: the case \"repeated key with an escaped spelling\" has the body {\"name\":\"x\",\"name\":\"y\"}; neither key is escaped, so SJ17's mapped evidence is a plain duplicate. Coordinator confirmed the raw bytes.\nAdvice:\n- verdict.go:341 \"The typed decode stays here\" narrates a move.\n- decode.go:14-15 the group comment does not start with a symbol name.\n- json_compatibility_test.go:22 implies a run that writes the baseline.\n- decode_test.go:131-133 repeats earlier bodies in a second table loop; incidental repetition.\n- strictJSON runs the suffix rule three times per call; one source, cost note only.\nHarness reuse: no second harness is pasted; json_compatibility_test.go reuses writeLaneFile, outcomeWrite, phaseTable, manifestPath, defect, phaseNames, inspectJSON, readyVerdict.\nRejectTrailing passes the one-source rule: one suffix owner at decode.go:115; the gate keeps only its typed decode for precedence and maps owner refusals to its sentinels.\nImplementation command: no change. Steps 1-8 done.\n"
          },
          "axis": "Standards",
          "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
          "tip": "4f88b2d46289a69c396185d1a40de5f0f6182326",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "sj-c2-spec-1",
          "performer": "ft354-review-c2-spec-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "cee4624207f1395d09c8a38efa3288218981b4ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c2-spec-1",
            "digest": "sha256:1ef7e87b8e641a6ff498da06b7a3b4913fb58b48d03626d81607d309b214c65b",
            "excerpt": "## Spec\n1 finding. Worst issue: row SJ17 is falsely classified as proven.\n- P1, severity 5, confidence 9, auto-fix, SJ17, internal/jsonfile/decode_test.go:60. The case \"repeated key with an escaped spelling\" has the body {\"name\":\"x\",\"name\":\"y\"} with no escape; a raw-byte walker would pass it. Production is correct, because json.Decoder.Token returns unescaped names (decode.go:78-89). The recorded SJ17 red came only from the missing category. Fix: an escaped body, inside ticket 2's Writes.\nRows SJ16 and SJ18-SJ31 hold at the tip; all 16 seam cells name existing tests. strictJSON (verdict.go:343-357) keeps typed, then trailing, then duplicate precedence, and the baseline matches. rejectDuplicateNames is gone; decodeManifest and unknownField are unchanged; the other callers keep their calls; the baseline test only reads its file and normalizes only $ROOT.\nAdvice: record jsonfile.RejectTrailing as a Terms-and-API addition; strictJSON could map trailing data by errors.Is(err, jsonfile.ErrTrailingData) rather than any non-nil result.\nImplementation command: did not contribute. Steps 1-9 done.\n"
          },
          "axis": "Spec",
          "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
          "tip": "4f88b2d46289a69c396185d1a40de5f0f6182326",
          "finding_ids": [
            "P1"
          ],
          "supersedes": []
        },
        {
          "id": "sj-c2-coverage-1",
          "performer": "ft354-review-c2-coverage-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "cee4624207f1395d09c8a38efa3288218981b4ec",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c2-coverage-1",
            "digest": "sha256:ee8647cd7178cdc5830d499fe3a00665567f4c884e08e7dfaff5abd070e420e1",
            "excerpt": "## Coverage\n1 finding. Worst issue: SJ17 (escaped-spelling duplicate keys) has a named case but no test.\n- C1, severity 5, confidence 9, auto-fix, SJ17 (spec.md:200), internal/jsonfile/decode_test.go:60. The body {\"name\":\"x\",\"name\":\"y\"} has no escape, so a raw-spelling duplicate check would pass. Production keys on the decoded Token() string (decode.go:82-89). The recorded SJ17 red came from the missing ErrDuplicateKey category.\nClass check (json_compatibility_test.go:98-121, 150-164): present are framed, unframed, whitespace, null, unknown, duplicate root and nested, syntax error, second value, malformed suffix, unknown-plus-suffix, duplicate-plus-suffix, size, mode, regular-file, absent. Unreadable is partial: \"under a file\" covers the Lstat failure; the open/read failure (\"cache unavailable\", verdict.go:285/291) has no case.\nCallers: all six route through strictJSON; tests reach decodeManifest, loadVerdict, and loadEvidence; a keyword sample found no duplicate or trailing case for loadManifest, parseReaderSchedule, loadReaderTable.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  silent,internal/gate/verdict.go,swap,passed,0,null,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateJSONCompatibility,passed,7\nThe mutation swapped occurrence 1 of `Unavailable, \"cache unavailable\"` for `Ready, \"\"`.\nAdvice: the unreadable cache state has no open/read-failure case; one duplicate row per remaining caller would pin each refusal; the hostile-input checklist does not apply.\nImplementation command: did not contribute. Steps 1-10 done.\n"
          },
          "axis": "Coverage",
          "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
          "tip": "4f88b2d46289a69c396185d1a40de5f0f6182326",
          "finding_ids": [
            "C1"
          ],
          "supersedes": []
        },
        {
          "id": "sj-c2-standards-2",
          "performer": "ft354-review-c2-standards-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c3b2a959e0f36fdafa82ade2a1e1a7e3e7967248",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c2-standards-2",
            "digest": "sha256:c74a307e8521802e49ff70289adccfc643a0709c60900d173a56152fac8157a2",
            "excerpt": "## Standards\n0 new findings; S1 is closed. At bc1d783e, internal/jsonfile/decode_test.go:61 holds a raw-string body with a \\u escape, so the test has two spellings of one decoded key.\nRepair code passes AGENTS.md and craft-comments: hostCacheFiles at json_compatibility_test.go:180 is the one fixture source and removes its entries from the stored baseline (lines 94-100); the skip uses capability.Environment inside the subtest (line 101); comments are present tense with no tags.\nAdvice: the comment at json_compatibility_test.go:177-180 names root, but the real condition is ptrace access; the review record showed no prose reds for the new rows.\nSteps 1-5 done.\n"
          },
          "axis": "Standards",
          "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
          "tip": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "finding_ids": [],
          "supersedes": [
            "sj-c2-standards-1"
          ]
        },
        {
          "id": "sj-c2-spec-2",
          "performer": "ft354-review-c2-spec-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c3b2a959e0f36fdafa82ade2a1e1a7e3e7967248",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c2-spec-2",
            "digest": "sha256:2b959825a9e7eb4db6c53ba7a826fa6c3c96afa56fe59720f1a3c3abf910a538",
            "excerpt": "## Spec\nFindings: 0. P1 is closed, and SJ27 is covered.\nP1: SJ17 (spec.md:201); decode_test.go:61 holds a \\u escape that decodes to the same key, expecting ErrDuplicateKey; repair mutant 1 names this case first.\nSJ27 (spec.md:211; edge inventory spec.md:248): json_compatibility_test.go:180-197 adds /proc/1/mem and /proc/self/mem, compared at :104 with unavailable / cache unavailable entries; both pass the metadata check and reach the open and read refusals at verdict.go:285 and :291.\nCoverage cells for SJ17 and SJ27 cite tests that exist at bc1d783e.\nAdvice: as root or with permissive ptrace, the unopenable case reaches the read refusal instead; the record lacked prose reds for SJ27 and a repair prose section.\nSteps 1-5 done.\n"
          },
          "axis": "Spec",
          "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
          "tip": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "finding_ids": [],
          "supersedes": [
            "sj-c2-spec-1"
          ]
        },
        {
          "id": "sj-c2-coverage-2",
          "performer": "ft354-review-c2-coverage-2",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "c3b2a959e0f36fdafa82ade2a1e1a7e3e7967248",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c2-coverage-2",
            "digest": "sha256:a1cad1a1ec57e1983db8da4a55020e68e5a20744e7d91efcfe3715fe932724fa",
            "excerpt": "## Coverage\n0 findings. C1 is closed and SJ27 is closed.\nC1: internal/jsonfile/decode_test.go:61 holds the escaped body with ErrDuplicateKey; the raw-span mutant turned only this case red.\nSJ27: occurrence 1 (verdict.go:285, failed open) and occurrence 2 (verdict.go:291, failed read or close) are witnessed. The round 1 silent mutation now bites:\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/gate/verdict.go,swap,failed,1,TestGateJSONCompatibility/unreadable_cache_baseline,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/gate,TestGateJSONCompatibility,passed,8\nHost and CI run as non-root (native-runtime.yml:20,38), so both are red-capable. As root, only the occurrence-1 witness goes silent and the test cannot fail falsely. Non-Linux skips through capability.Environment (json_compatibility_test.go:55).\nAdvice: a non-root container with ptrace_scope 0 could open /proc/1/mem; consider a capability skip when the open succeeds (unverified).\nSteps 1-5 done; gate 8 and jsonfile 15 tests pass.\n"
          },
          "axis": "Coverage",
          "base": "533160a782e35a64f68d5944f35f5103e0eadacd",
          "tip": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "finding_ids": [],
          "supersedes": [
            "sj-c2-coverage-1"
          ]
        }
      ]
    },
    {
      "id": "SJ-C3",
      "base": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
      "tip": "7e962ba222ab4879f10957a60bfe8ff777630ea7",
      "plan_digest": "sha256:352f91ffc5cbebf5221070c81e6c1f2898a87851ee1d00dd37ec1ab08a479802",
      "source_digest": "72f25b7ffe8bd44566ddd58fa02ea55423a80d45",
      "acceptance_rows": [
        "SJ32",
        "SJ33",
        "SJ34",
        "SJ35"
      ],
      "verification": [
        {
          "id": "sj-c3-review",
          "performer": "ft354-author-03",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2bcb79a2c20eef4cd362a1b0c870ab1972ba958c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft354-author-03",
            "digest": "sha256:a7446f26d49f68013b277518f5cd3951358c9c9a37009e99105b9a8d6aab4bd2",
            "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewJSONCompatibility/unframed_review_document,2,1,0,null,yes\n  2,bit,swap,failed,1,TestReviewJSONCompatibility/fenced_completion_plan,2,1,0,null,yes\n  3,bit,swap,failed,1,TestReviewJSONCompatibility/record_baseline,2,1,0,null,yes\n  4,bit,swap,failed,1,TestReviewJSONCompatibility/byte_limits,2,1,0,null,yes\n  5,bit,swap,failed,1,TestReviewJSONCompatibility/unframed_review_document,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord,^TestReviewJSONCompatibility$/^unframed_review_document$|^TestReviewJSONCompatibility$/^fenced_completion_plan$|^TestReviewJSONCompatibility$/^record_baseline$|^TestReviewJSONCompatibility$/^byte_limits$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "2c72f6478c064633d4d6879c82e54178b835b3aa",
          "requirement": "review",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the review bytes for null at the jsonfile.DecodeDocument call in internal/reviewrecord/parse.go, and record that the valid review-document case of TestReviewJSONCompatibility turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft354-author-03",
              "digest": "sha256:a7446f26d49f68013b277518f5cd3951358c9c9a37009e99105b9a8d6aab4bd2",
              "excerpt": "mutants[5]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewJSONCompatibility/unframed_review_document,2,1,0,null,yes\n  2,bit,swap,failed,1,TestReviewJSONCompatibility/fenced_completion_plan,2,1,0,null,yes\n  3,bit,swap,failed,1,TestReviewJSONCompatibility/record_baseline,2,1,0,null,yes\n  4,bit,swap,failed,1,TestReviewJSONCompatibility/byte_limits,2,1,0,null,yes\n  5,bit,swap,failed,1,TestReviewJSONCompatibility/unframed_review_document,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord,^TestReviewJSONCompatibility$/^unframed_review_document$|^TestReviewJSONCompatibility$/^fenced_completion_plan$|^TestReviewJSONCompatibility$/^record_baseline$|^TestReviewJSONCompatibility$/^byte_limits$,passed,6\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        },
        {
          "id": "sj-c3-r1-review",
          "performer": "ft354-repair-03-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "72f25b7ffe8bd44566ddd58fa02ea55423a80d45",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-medium:ft354-repair-03-1",
            "digest": "sha256:50327d5db54caea570b0ac94cc63f53e6dd43294ddc38969f217e7138968fa66",
            "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewJSONCompatibility/record_baseline,2,1,0,null,yes\n  2,bit,swap,failed,1,TestReviewJSONCompatibility/plan_baseline,2,1,0,null,yes\n  3,bit,swap,failed,1,TestReviewJSONCompatibility/record_baseline,2,1,0,null,yes\n  4,bit,swap,failed,1,TestReviewJSONCompatibility/unframed_review_document,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord,^TestReviewJSONCompatibility$/^record_baseline$|^TestReviewJSONCompatibility$/^plan_baseline$|^TestReviewJSONCompatibility$/^unframed_review_document$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
          },
          "source_commit": "7e962ba222ab4879f10957a60bfe8ff777630ea7",
          "requirement": "review",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0,
          "probe": {
            "mutation": "Swap the review bytes for null at the jsonfile.DecodeDocument call in internal/reviewrecord/parse.go, and record that the valid review-document case of TestReviewJSONCompatibility turns red.",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:bench-writer-medium:ft354-repair-03-1",
              "digest": "sha256:50327d5db54caea570b0ac94cc63f53e6dd43294ddc38969f217e7138968fa66",
              "excerpt": "mutants[4]{mutant,verdict,mutation,cause,failed_tests,first_failure,ran,focused_exit_code,exit_code,first_build_line,restored}:\n  1,bit,swap,failed,1,TestReviewJSONCompatibility/record_baseline,2,1,0,null,yes\n  2,bit,swap,failed,1,TestReviewJSONCompatibility/plan_baseline,2,1,0,null,yes\n  3,bit,swap,failed,1,TestReviewJSONCompatibility/record_baseline,2,1,0,null,yes\n  4,bit,swap,failed,1,TestReviewJSONCompatibility/unframed_review_document,2,1,0,null,yes\nselection[1]{form,target,run,baseline,builds}:\n  package,./internal/reviewrecord,^TestReviewJSONCompatibility$/^record_baseline$|^TestReviewJSONCompatibility$/^plan_baseline$|^TestReviewJSONCompatibility$/^unframed_review_document$,passed,5\nprobe_set_execution[1]{schema,command_exit_code,trailing_bytes}:\n  bench-probe-set-v1,0,0\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "sj-c3-standards-1",
          "performer": "ft354-review-c3-standards-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2bcb79a2c20eef4cd362a1b0c870ab1972ba958c",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c3-standards-1",
            "digest": "sha256:9a1faf0a5959acb8e0c60c5bccb8f2b22f1dcceb3aeef3e5ed411c311e07626f",
            "excerpt": "## Standards\n1 finding. Worst issue: S1, the gate baseline-family harness is pasted a second time into the review package, and SJ-C4 would paste a third copy.\n- S1 (severity 5, confidence 8, auto-fix), internal/reviewrecord/json_compatibility_test.go:120, :197, :210 at 2c72f647. AGENTS.md: \"Knowledge duplication is a defect ... a fixture harness pasted N times.\" reviewFamily (:120-145) repeats the 18-case table of jsonFamily (gate json_compatibility_test.go:109-133); compareReviewJSON (:210-232) copies compareGateJSON[T] (gate :208-230); readReviewJSONBaseline (:197-208) copies readGateJSONBaseline (gate :195-206). spec.md:249 makes the case list one fact. Owner: an exported test-support package internal/jsonfile/jsonfiletest (new) with Family, ReadBaseline[T], and Compare[T], consumed by the gate, review, and release tests; the path sits inside the directory of the signed-off SJ-C2 entry internal/jsonfile/decode.go.\nAdvice: the trailing branch of decode (parse.go:239-241) runs a second strict decode to rank diagnostics; its comment states why. Comments pass craft-comments. WritePlanPayload (fixture.go:86) removes the old planDocument duplication.\nImplementation command contributed: no ticket's Writes line holds a shared owner for the harness; ticket 3 forbade only the fence and repository harness copies.\nSteps 1-7 done.\n"
          },
          "axis": "Standards",
          "base": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "tip": "2c72f6478c064633d4d6879c82e54178b835b3aa",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "sj-c3-spec-1",
          "performer": "ft354-review-c3-spec-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2bcb79a2c20eef4cd362a1b0c870ab1972ba958c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c3-spec-1",
            "digest": "sha256:a4911c3c79b1ab455329132d13f76ab8c43051566014f4d1f427b99e5cad30e6",
            "excerpt": "## Spec\n0 findings. Worst issue is advice A1: the review family lacks the nested cases of the gate family.\nbounds.ClassifyBytes runs before jsonfile.DecodeDocument (parse.go:228-232). uniqueJSON is removed; decode maps ErrDuplicateKey (parse.go:236) and re-decodes on ErrTrailingData (parse.go:238-246) so a type error outranks the suffix, which spec.md:113-114 allows.\nCompeting inputs hold against a hand trace of the old code: unknown-plus-suffix, duplicate-plus-suffix, whitespace-only, null, syntax, wrong type, and suffix variants.\nfenced, locate, plan.go, parsePlan, ReadPlan, checkPlan, and the Parse domain checks are untouched. SJ32-SJ35 cells cite TestReviewJSONCompatibility (json_compatibility_test.go:36).\nAdvice: A1 the review family lacks nested duplicate and nested unknown (gate json_compatibility_test.go:82-83), severity 3; A2 the mutant selection leaves out plan_baseline; A3 the pre-edit capture is not visible in history.\nSteps 1-8 done.\n"
          },
          "axis": "Spec",
          "base": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "tip": "2c72f6478c064633d4d6879c82e54178b835b3aa",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "sj-c3-coverage-1",
          "performer": "ft354-review-c3-coverage-1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "2bcb79a2c20eef4cd362a1b0c870ab1972ba958c",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer:ft354-review-c3-coverage-1",
            "digest": "sha256:f51bf29eb9f757cf5bfc6de8d625c483351d7c721e1825e96d6ff3dac4eb4abc",
            "excerpt": "## Coverage\n0 findings. Worst issue is advice only: the review baseline has no nested-duplicate, escaped-duplicate, or empty-array case; old and new code give the same result on all three.\nPresent classes (json_compatibility_test.go:114-151 against spec.md:240-249): valid record and plan seeds, null, unsupported version, wrong type, invalid evidence, invalid spec path, version 2 session, oversize, invalid UTF-8, framed, unframed, whitespace, whitespace only, empty, unknown, root duplicate, syntax, second value, framed second value, malformed suffix, unknown plus suffix, duplicate plus suffix, wrong-type and unknown-plus-duplicate combinations.\nBoth callers: Parse and ReadPlan (through parsePlan, plan.go:64) run the duplicate, suffix, and precedence classes; byte-limit and domain cases run through Parse, and the bounds check sits in the shared decode.\nprobe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  invalid,internal/jsonfile/decode.go,swap,build-failed,0,null,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/reviewrecord,TestReviewJSONCompatibility,passed,0\nAdvice: add nested-duplicate and escaped-duplicate rows to the review family (severity 3, not a correctness defect); the plan side has no domain-refusal rows beyond null and unsupported version.\nImplementation command: did not contribute. Steps 1-9 done (invalid probe).\n"
          },
          "axis": "Coverage",
          "base": "bc1d783e99c01ae3629aa6c40dbd947e4fb7b7fb",
          "tip": "2c72f6478c064633d4d6879c82e54178b835b3aa",
          "finding_ids": [],
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
      "from": "sha256:f12fdb6d19f0b2ae0d0689ed9dd450dd26628d69f7f864a711955a0cf02e43ff",
      "to": "sha256:85b39d6036837dbfefc2ff28e2bb5a87ee2030867effe236f7d8b8a82b7442c4",
      "chunk_ids": {
        "SJ-C1": [
          "SJ-C1"
        ]
      }
    },
    {
      "from": "sha256:85b39d6036837dbfefc2ff28e2bb5a87ee2030867effe236f7d8b8a82b7442c4",
      "to": "sha256:4eb2f514d61a626059587c5896e808e64450264ba12146938389c4f1e1673f57",
      "chunk_ids": {
        "SJ-C1": [
          "SJ-C1"
        ],
        "SJ-C2": [
          "SJ-C2"
        ]
      }
    },
    {
      "from": "sha256:4eb2f514d61a626059587c5896e808e64450264ba12146938389c4f1e1673f57",
      "to": "sha256:6d7996b4cd67f9532568e966e905982078cf48a880d96d3e0ae0df593646db72",
      "chunk_ids": {
        "SJ-C1": [
          "SJ-C1"
        ],
        "SJ-C2": [
          "SJ-C2"
        ]
      }
    },
    {
      "from": "sha256:6d7996b4cd67f9532568e966e905982078cf48a880d96d3e0ae0df593646db72",
      "to": "sha256:352f91ffc5cbebf5221070c81e6c1f2898a87851ee1d00dd37ec1ab08a479802",
      "chunk_ids": {
        "SJ-C1": [
          "SJ-C1"
        ],
        "SJ-C2": [
          "SJ-C2"
        ],
        "SJ-C3": [
          "SJ-C3"
        ]
      }
    }
  ],
  "folds": [
    {
      "id": "F1",
      "after": "SJ-C1",
      "from": "d02ea82677295e38832d2eaa65862ab1128ca5cd",
      "commit": "533160a782e35a64f68d5944f35f5103e0eadacd",
      "incoming": "e6939da810757e1dc463d36a4d62e80fdee1c60a",
      "source_digest": "3fc51898cf3629890a792750a2dc88a9da91c782",
      "resolution": [],
      "resolution_digest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "test_packages": [],
      "reviews": []
    }
  ]
}
```

## SJ-C1 ticket 1 author witness reds

The ticket 1 author `ft354-author-01` (opus / medium) wrote the new tests first. It ran them in one `bench test` call against the unchanged production code at `6da62689`.

| subject | condition | failing witness | row |
| --- | --- | --- | --- |
| `internal/releasepreflight/vulnerability.go` | second decode refuses only a nil error | `TestVulnerabilityPolicyDocumentContract/empty_array_then_malformed_suffix` and `used_array_then_malformed_suffix` | SJ1 |
| `internal/releasepreflight/vulnerability.go` | fixed second target | the number, boolean, string, and array cases of `TestVulnerabilityPolicyDocumentContract` | SJ3 |
| `internal/releasepreflight/vulnerability.go` | no duplicate refusal | the `repeated_id_key`, `repeated_reason_key`, and `repeated_expires_key` cases | SJ4 |

The object and null cases of SJ2 and SJ3 went red on the message only, because the old code already refused them. The plan probe `policy` answered `bit` on SJ1 to SJ4 and SJ6 after the ticket commit `e32079f0`.
SJ5 has no observed red, and its claimed red is a swap to `jsonfile.Decode` framing. SJ6 and the preservation rows SJ7 to SJ15 passed on the old code, and the author claims their mutations.

## SJ-C1 round 1 review

The three axes ran on opus / medium over the frozen pair `835f5491..e32079f0`. Standards and Coverage pass with zero finding IDs. Spec raised P1 at severity 5: rows SJ5 and SJ7 to SJ15 had no recorded red.
P1 is an evidence gap, and the code and tests stay unchanged. The repair count of SJ-C1 is 0 of 2.

The Coverage probe of this round bit SJ5. Its swap to `jsonfile.Decode` framing turned `surrounding_whitespace_without_final_newline` red, with restored `yes`.
The coordinator ran one mutant set of ten domain mutations at `1b12ee8b` over `TestVulnerabilityPolicyDomainCompatibility`. Each row answered `bit` with restored `yes`.

| mutant | mutation in `ValidateVulnerabilityPolicy` | first failing case | row |
| --- | --- | --- | --- |
| 1 | `if present {` to `if true {` | `absent_policy_without_findings` in the narrowed rerun | SJ7 |
| 2 | `if len(data) == 0 {` to `if false {` | `zero-byte_policy` | SJ8 |
| 3 | `if len(data) == 0 {` to `if len(data) < 5 {` | `empty_array_without_findings` and `null_without_findings` in the narrowed rerun | SJ9 |
| 4 | `if len(wanted) > 0 {` to `if false {` | `absent_policy_with_a_finding` | SJ10 |
| 5 | `<` to `<=` in the expiry check | `expiry_equal_to_today` | SJ11 |
| 6 | the expiry check to `if false {` | `expiry_before_today` | SJ12 |
| 7 | `if seen[item.ID] {` to `if false {` | `duplicate_exception_entries` | SJ13 |
| 8 | `if !wanted[item.ID] {` to `if false {` | `unused_exception` | SJ14 |
| 9 | `strings.TrimSpace(item.ID) == ""` to `item.ID == ""` | `blank_id` | SJ15 |
| 10 | omit the full-date match | `partial_date` | SJ15 |

The Spec reaffirmation found that mutants 1 and 3 named SJ10 cases first. A narrowed rerun at `b68631c1` turned the SJ7 and SJ9 cases red, with restored `yes`.

## SJ-C1 advice

- Standards: an optional why-comment on the empty-file check, and the containment wording of the test helper comment.
- Standards: the history wording of the inventory sentence, the repeated SJ-C5 fence list, and the old missing doc comment.
- Spec: the census narrowing has a `bench learning` entry, and the SJ-C5 canary narrowing removed paths that a landed commit deleted.
- Coverage: no test pins a present whitespace-only policy.

## SJ-C2 ticket 2 author witness reds

The ticket 2 author `ft354-author-02` (opus / high) wrote the new tests first. It ran them in one `bench test` call against the unchanged production code at `55a5813c`, with 21 failures.
The witness rows are SJ16, SJ17, SJ19, SJ20, and SJ25 in `TestDocumentStructuralContract` and `TestDocumentErrorCategories`.

The author ran three mutant sets after the ticket commit, and every row answered `bit` with restored `yes`.

| subject | set | rows |
| --- | --- | --- |
| `internal/jsonfile/decode.go` | framing set over `TestDocumentFramingCompatibility` | SJ22, SJ23, SJ24 |
| `internal/jsonfile/decode.go` | structural set over `TestDocumentStructuralContract`, with the plan probe | SJ18, SJ21, SJ16 |
| `internal/gate/verdict.go` | one `--expect` set over `TestGateJSONCompatibility` | SJ28, SJ29, SJ30, SJ31, SJ26, SJ27 |

The plan probe `owner-gate` swapped `seen[name] = true` and turned the nested duplicate case red. The ticket adds the export `jsonfile.RejectTrailing`, so the gate keeps one suffix owner.

## SJ-C2 round 1 review

The three axes ran on opus / medium over the frozen pair `533160a7..4f88b2d4`. Each axis raised one finding at severity 5, and all three name one repair target.

## Standards

Findings: 1. Worst issue: S1, auto-fix, confidence 9. The case `repeated key with an escaped spelling` at `internal/jsonfile/decode_test.go:60` has no escaped key.

## Spec

Findings: 1. Worst issue: P1, auto-fix, confidence 9. Row SJ17 has no test, because the case body repeats one raw spelling.

## Coverage

Findings: 1. Worst issue: C1, auto-fix, confidence 9. SJ17 is untested, for the same reason as S1 and P1.
The Coverage probe of the `cache unavailable` result at `internal/gate/verdict.go` was silent. The orchestrator adds that result to the repair, because SJ27 names the unavailable result.

The de-duplicated repair target count is 1, and the orchestrator adds the SJ27 gap to the same session. The repair count of SJ-C2 is 1 of 2.

## SJ-C2 advice

- Standards: the move narration at `verdict.go:341`, the group comment of the categories, the baseline write wording, the second table loop, and the triple suffix check.
- Spec: record `jsonfile.RejectTrailing` in Terms and API, and map trailing data by category.
- Coverage: one duplicate row for each of loadManifest, parseReaderSchedule, and loadReaderTable.

## SJ-C2 repair 1 witness reds

The repair session `ft354-repair-02-1` (opus / high) changed tests and testdata only, at `f4eb1f9b` and `bc1d783e`. Each row below answered `bit` with restored `yes`.

| subject | mutation | failing witness | row |
| --- | --- | --- | --- |
| `internal/gate/verdict.go` | occurrence 1 of `Unavailable, "cache unavailable"` to `Ready, ""` | `TestGateJSONCompatibility/unreadable_cache_baseline` | SJ27 |
| `internal/gate/verdict.go` | occurrence 2 of the same result | `TestGateJSONCompatibility/unreadable_cache_baseline` | SJ27 |
| `internal/jsonfile/decode.go` | duplicate memory keyed on the name and its raw escape span | `TestDocumentStructuralContract/repeated_key_with_an_escaped_spelling`, the only failure | SJ17 |
| `internal/jsonfile/decode.go` | plan probe `owner-gate`, `seen[name] = false` | `TestDocumentStructuralContract/nested_duplicate_in_an_array_element` | SJ16 |

Before the two new baseline entries existed, the unreadable cache subtest failed on the missing entries. The coordinator probe swapped the read-failure condition to `readErr == io.EOF`, and the unreadable case turned red.

## SJ-C2 round 2 confirming review

The three axes ran on opus / medium over the repair delta `4f88b2d4..bc1d783e`. Each axis passes with zero finding IDs, and S1, P1, and C1 close. The repair count of SJ-C2 is 1 of 2.

## SJ-C2 round 2 advice

- Standards and Coverage: the `/proc/1/mem` case depends on ptrace access, not on root alone. With ptrace access, the case reaches the read refusal and the open witness goes silent.

## SJ-C3 ticket 3 author witness reds

The ticket 3 author `ft354-author-03` (opus / medium) captured the review baseline from the unchanged code. All four rows are preservation rows, so the one pre-edit call passed.
After the ticket commit `f7c06590`, one mutant set over `internal/reviewrecord/parse.go` turned each row red, with restored `yes`.

| mutant | mutation | failing witness | row |
| --- | --- | --- | --- |
| 1 | `DecodeDocument` to `Decode` | `TestReviewJSONCompatibility/unframed_review_document` | SJ32 |
| 2 | the fence span start moved by one byte | `TestReviewJSONCompatibility/fenced_completion_plan` | SJ33 |
| 3 | the precedence decode to `new(any)` | `TestReviewJSONCompatibility/record_baseline` | SJ34 |
| 4 | the byte-limit state check to `false` | `TestReviewJSONCompatibility/byte_limits` | SJ35 |
| 5 | plan probe `review`, the real bytes to `null` | `TestReviewJSONCompatibility/unframed_review_document` | SJ32 |

The author flagged that three review test helpers repeat three gate test helpers. A shared helper needs a path outside the ticket 3 fence.

## SJ-C3 round 1 review

The three axes ran on opus / medium over the frozen pair `bc1d783e..2c72f647`. Spec and Coverage pass with zero finding IDs. The repair count of SJ-C3 is 1 of 2.

## Standards

Findings: 1. Worst issue: S1, severity 5, auto-fix, confidence 8. The review test pastes the gate baseline-family harness at `internal/reviewrecord/json_compatibility_test.go:120`, `:197`, and `:210`.
The single owner is a new test-support package `internal/jsonfile/jsonfiletest`. The gate, review, and release compatibility tests consume it.

## Spec

Findings: 0. Worst issue: none. Advice A1 says the review family lacks the nested cases of the gate family, and one shared family closes it.

## Coverage

Findings: 0. Worst issue: none. The one Coverage probe did not compile, so it proves nothing.

## SJ-C3 advice

- Standards: the second strict decode in the trailing branch ranks the diagnostics, and its comment states why.
- Spec: the mutant selection leaves out `plan_baseline`, and history does not show the pre-edit capture.
- Coverage: the plan side has no domain refusal beyond null and the unsupported version.

## SJ-C3 repair 1 witness reds

The repair session `ft354-repair-03-1` (opus / medium) moved the hostile-document family, the baseline reader, and the generic compare into `internal/jsonfile/jsonfiletest` at `7e962ba2`. The gate and review tests consume it.
The review baseline gained four nested entries, captured against the base `parse.go`. The gate entry `nested unknown field` now names the key `extra`, and the repair claims the old result from a code read.

| mutant | mutation in `internal/reviewrecord/parse.go` | first failure | restored |
| --- | --- | --- | --- |
| 1 | the duplicate-key case to `false &&` | `TestReviewJSONCompatibility/record_baseline` | yes |
| 2 | the same mutation | `TestReviewJSONCompatibility/plan_baseline` | yes |
| 3 | the unwrap skips an unknown-field cause | `TestReviewJSONCompatibility/record_baseline` | yes |
| 4 | plan probe `review`, the real bytes to `null` | `TestReviewJSONCompatibility/unframed_review_document` | yes |
