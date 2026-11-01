# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/strict-json-ownership/spec.md",
  "plan_digest": "sha256:f12fdb6d19f0b2ae0d0689ed9dd450dd26628d69f7f864a711955a0cf02e43ff",
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
