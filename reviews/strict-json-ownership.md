# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/strict-json-ownership/spec.md",
  "plan_digest": "sha256:85b39d6036837dbfefc2ff28e2bb5a87ee2030867effe236f7d8b8a82b7442c4",
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
      "tip": "4f88b2d46289a69c396185d1a40de5f0f6182326",
      "plan_digest": "sha256:85b39d6036837dbfefc2ff28e2bb5a87ee2030867effe236f7d8b8a82b7442c4",
      "source_digest": "cee4624207f1395d09c8a38efa3288218981b4ec",
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
      "from": "sha256:f12fdb6d19f0b2ae0d0689ed9dd450dd26628d69f7f864a711955a0cf02e43ff",
      "to": "sha256:85b39d6036837dbfefc2ff28e2bb5a87ee2030867effe236f7d8b8a82b7442c4",
      "chunk_ids": {
        "SJ-C1": [
          "SJ-C1"
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
