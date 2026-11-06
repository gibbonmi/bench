# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/production-test-seam-policy/spec.md",
  "plan_digest": "sha256:d60cbebbc23104cb18f056422e79d5ec33c22e8db5f7f151705f4e208f69a747",
  "implementation_session": "",
  "chunks": [
    {
      "id": "TS1",
      "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
      "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
      "plan_digest": "sha256:5c08a62494fc1fdd742618a16f4d9ccab37c50326fb6f594f7fef3f23db32a4b",
      "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
      "acceptance_rows": [
        "PS01",
        "PS02",
        "PS03",
        "PS12",
        "PS13",
        "PS14",
        "PS41",
        "PS42",
        "PS43",
        "PS50",
        "PS51",
        "PS52",
        "PS53",
        "PS54",
        "PS55",
        "PS61"
      ],
      "verification": [
        {
          "id": "ts1-01-wrapper-bypass",
          "performer": "ft343-author-01",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-author-01",
            "digest": "sha256:5ff4c883590de2b57efe65fadff342ba3c996a22b74d4c1e1136f5892cbb5329",
            "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_test.go,swap,failed,2,TestEveryRetainedFixtureBitesThroughRegisteredOwner,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestEveryRetainedFixtureBitesThroughRegisteredOwner$,passed,289\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4108,289\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner,\"fixture_bite_test.go:33: completed fixture proofs do not equal the current producer: got 287, want 288\",1\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner/unregistered-port,\"fixture_bite_test.go:27: unregistered-port did not bite through owner injected-port-registry; want \\\"injected port has no registry row: internal/gitguard.Auditor is injected at a parameter of auditVerdict in internal/gitguard/gitguard.go\\\" in… (253 bytes)\",1\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "requirement": "wrapper-bypass",
          "command": "bench test --package ./internal/conformance --run '^TestEveryRetainedFixtureBitesThroughRegisteredOwner$'",
          "exit_code": 0,
          "probe": {
            "mutation": "make checkInjectedPortRegistry return without calling the relocated owner; require the injected-ports canary bite to fail with its exact expected diagnostic",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "claude-agent:ft343-author-01",
              "digest": "sha256:5ff4c883590de2b57efe65fadff342ba3c996a22b74d4c1e1136f5892cbb5329",
              "excerpt": "probe[1]{verdict,subject,mutation,cause,failed_tests,first_failure,restored}:\n  bit,internal/conformance/injected_ports_test.go,swap,failed,2,TestEveryRetainedFixtureBitesThroughRegisteredOwner,yes\nselection[1]{form,target,run,baseline,ran}:\n  package,./internal/conformance,^TestEveryRetainedFixtureBitesThroughRegisteredOwner$,passed,289\nprobe_execution[1]{schema,command_exit_code,focused_exit_code,trailing_bytes}:\n  bench-probe-result-v2,0,1,771\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,fail,4108,289\nfailures[2]{package,test,line,lines}:\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner,\"fixture_bite_test.go:33: completed fixture proofs do not equal the current producer: got 287, want 288\",1\n  github.com/gibbonmi/bench/internal/conformance,TestEveryRetainedFixtureBitesThroughRegisteredOwner/unregistered-port,\"fixture_bite_test.go:27: unregistered-port did not bite through owner injected-port-registry; want \\\"injected port has no registry row: internal/gitguard.Auditor is injected at a parameter of auditVerdict in internal/gitguard/gitguard.go\\\" in… (253 bytes)\",1\nskips[0]{package,test,reason}:\n"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "ts1-review-standards",
          "performer": "ft343-review-ts1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts1",
            "digest": "sha256:8e17e424b12352d673152240db0d7ebd8855fe12de99136be6a455dde63a24b7",
            "excerpt": "Standards: no findings. The relocation leaves one AST scanner in internal/conformance/injectedports, and the old file keeps only the wrapper and the legacy bite test. The tests use the shared writeFixtureFile in conformance and one local writeRoot in the owner, as the spec allows. slices.Sort/Compact replaces an unexported test helper the package cannot reach, so it is not duplicated knowledge. Advice: put the per-row mutant reds in the record so the independent-expectation exception holds.\n"
          },
          "axis": "Standards",
          "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
          "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ts1-review-spec",
          "performer": "ft343-review-ts1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts1",
            "digest": "sha256:1a43e6e5bfbfcf2bffe0784b5bfbbee9bd244df71035ea37485b20fbf71c96d7",
            "excerpt": "Spec: no findings. All 16 TS1 rows map to live seams, the legacy messages stay byte-exact as exported constants, and checkInjectedPortRegistry calls Check once. Deriving only over Required follows the TS1 chunk row and ticket 04:29, and it conflicts with ticket 01:12, which goes to reviewer veto. Silence on a row outside Required stays legacy behavior until the TS4 census.\n"
          },
          "axis": "Spec",
          "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
          "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "ts1-review-coverage",
          "performer": "ft343-review-ts1",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "09e08a98b1a98ed7686c7fe73ebf7582a9ced4ac",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:ft343-review-ts1",
            "digest": "sha256:b27ac4b95b8318a293563a59ffbc51892a22c1887f14edd53c6b981190ae46a3",
            "excerpt": "Coverage: no findings. Every mapped row's why-clause mutation changes its fixture's diagnostic heads, and the wrapper-bypass probe bit with the canary's exact diagnostic. Advice for TS4: ENOTDIR as absence, and symlinked package directories outside the root.\n"
          },
          "axis": "Coverage",
          "base": "77de12e23e7fc3f4745ac7132e1ceeae02ce396c",
          "tip": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "TS2",
      "base": "d626a7a9966a65e54e8aa17089fdaadbe5e69897",
      "tip": "cdeb6e383ca0d1e2095ea1bb1308dbe062e3a8d8",
      "plan_digest": "sha256:d60cbebbc23104cb18f056422e79d5ec33c22e8db5f7f151705f4e208f69a747",
      "source_digest": "8534ee5c398a9f6fb2e6632a94406be03591bf1c",
      "acceptance_rows": [
        "PS04",
        "PS05",
        "PS06",
        "PS07",
        "PS08",
        "PS15",
        "PS16",
        "PS17",
        "PS18",
        "PS19",
        "PS20",
        "PS21",
        "PS44",
        "PS57"
      ],
      "verification": [],
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
      "from": "sha256:5c08a62494fc1fdd742618a16f4d9ccab37c50326fb6f594f7fef3f23db32a4b",
      "to": "sha256:d60cbebbc23104cb18f056422e79d5ec33c22e8db5f7f151705f4e208f69a747",
      "chunk_ids": {
        "TS1": [
          "TS1"
        ]
      }
    }
  ]
}
```

## TS1 author witness reds

The author `ft343-author-01` reported these witness reds at `d626a7a`. Each mutant reported `bit`, and each file restored byte-identically.

| row | file | mutation | failing test |
|---|---|---|---|
| PS01 | `derive.go` | the interface arm returns false | `TestSeamDerivationForms` |
| PS02 | `derive.go` | the function-type arm returns false | `TestSeamDerivationForms` |
| PS03 | `derive.go` | the struct arm returns false | `TestSeamDerivationForms` |
| PS50 | `derive.go` | the pointer arm returns an empty name | `TestSeamDerivationForms` |
| PS51 | `derive.go` | the type-assertion arm returns early | `TestSeamDerivationForms` |
| PS52 | `derive.go` | the case-clause arm returns early | `TestSeamDerivationForms` |
| PS53 | `derive.go` | the interface arm returns true | `TestSeamDerivationForms` |
| PS54 | `derive.go` | a data field is skipped | `TestSeamDerivationForms` |
| PS55 | `derive.go` | the injection test is always true | `TestSeamDerivationForms` |
| PS12 | `derive.go` | the parse diagnostic is omitted | `TestSeamAdmissionEvidence` |
| PS13 | `check.go` | the orphan guard is always true | `TestSeamAdmissionEvidence` |
| PS14 | `check.go` | the zero-inventory branch is disabled | `TestSeamAdmissionEvidence` |
| PS42 | `check.go` | the package path joins the live checkout | `TestSeamAdmissionEvidence` |
| PS43 | `check.go` | every stat error counts as absence | `TestSeamAdmissionEvidence` |
| PS61 | `check.go` | a blank exemption compares as valid | `TestSeamAdmissionEvidence` |

## TS2 author witness reds

The author `ft343-author-02` reported these witness reds before `cdeb6e3`. Each mutant reported `bit`, each file restored byte-identically, and each baseline was green.

| file | mutation | failing case |
|---|---|---|
| `variables.go` | the test-replacement census is disabled | PS04, PS05, PS06, PS57 |
| `variables.go` | cross-package forwarding is disabled | PS07 |
| `variables.go` | same-package forwarding is disabled | PS08 |
| `variables.go` | the reader filter is disabled | the test-support negative case |
| `variables.go` | a `ForTest` name rule is added | PS44 |
| `variables.go` | setter reach and helper reach are disabled | the admitted-consumer cases |
| `variables.go` | the `Test` prefix check is omitted | the PS16 helper case |
| `variables.go` | the setter-owner exclusion is omitted | the PS15 setter-as-owner case |
| `variables.go` | the initializer reader is dropped | the initializer-owner case |
| `scope.go` | the parameter test is bypassed | the production-data negative case |
| `scope.go` | scope resolution is disabled | the shadowing negative case |
| `scope.go` | only the first assignment target is read | PS57 |
| `scope.go` | import resolution is disabled | PS04, PS07 |
| `seams.go` | the owner, purpose, and isolation checks are removed one at a time | PS15, PS18, PS19 |
| `seams.go` | the consumer checks are inverted or emptied | PS16, PS17 |
| `seams.go` | the duplicate, alias, and orphan checks are removed one at a time | PS20, PS21, the stale-row case |
| `seams.go` | the unregistered-setter diagnostic is removed | PS07, PS08 |

The reviewer decided on 2026-11-06 to retire `indexEncoder`, its setter, and its forwarder, because the seam had no consumer and no injectable failure.

## Standards

Findings: 0. Worst issue: none.

## Spec

Findings: 0. Worst issue: none.
The reviewer flags ticket 01 line 12 for reviewer veto: it claims the wide derivation that the chunk table gives to TS4.
Carry to the TS4 charge: a row outside the required packages stays silent until the TS4 census grades it.

## Coverage

Findings: 0. Worst issue: none.
Advice for TS4: decide ENOTDIR as absence, and decide a symlinked package directory that points outside the root.

Repair cycles consumed for TS1: 0 of 2.
