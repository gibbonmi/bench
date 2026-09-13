# Harness proof closure review

## HP-C1 — gate fixture ownership

The implementation centralizes gate scripts, their canonical manifest, and the
private fixture path in `internal/testrepo`. The author verification passed both
planned commands and demonstrated that an undeclared ambient command makes the
owner test fail before restoring the source. Independent Standards, Spec, and
Coverage reviewers each examined the frozen pair and returned a pass with no
findings in the single permitted review round.

```bench-review-record
{
  "version": 2,
  "spec": "specs/harness-proof-closure/spec.md",
  "plan_digest": "sha256:9d9dfa2e2789cc8782246b8fc697718828baca8a3a3fd2a7d1e722a77e65052a",
  "chunks": [
    {
      "id": "HP-C1",
      "base": "01d8b3c3fded6dc1bff7768faae1d4fbd48d2ae7",
      "tip": "42d44425656d170d22be7ba3f4de8277df84899d",
      "plan_digest": "sha256:9d9dfa2e2789cc8782246b8fc697718828baca8a3a3fd2a7d1e722a77e65052a",
      "source_digest": "497dea0db26729d73da34ec21330abb99d10a6c6",
      "acceptance_rows": ["HP1", "HP2", "HP3"],
      "verification": [
        {
          "id": "hp-c1-fixture-owner",
          "performer": "/root/ft120_c1",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "497dea0db26729d73da34ec21330abb99d10a6c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c1",
            "digest": "sha256:b72b1a0d0eca594d11783ce2c7a1560a4ef22d15fca50848040db3d179a83480",
            "excerpt": "pass: bench test --package ./internal/testrepo; probe bit after omitting manifest command token; restore pass"
          },
          "requirement": "fixture-owner",
          "command": "bench test --package ./internal/testrepo",
          "exit_code": 0,
          "probe": {
            "mutation": "write an ambient command token without recording it in the manifest",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "collaboration:/root/ft120_c1",
              "digest": "sha256:e742cd85ddbb18ba19f4236defab180390ac55bc259de4292109373480ed2155",
              "excerpt": "bit: undeclared ambient command was unavailable under the private fixture PATH; source restored and owner test passed"
            }
          }
        },
        {
          "id": "hp-c1-fixture-consumers",
          "performer": "/root/ft120_c1",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "497dea0db26729d73da34ec21330abb99d10a6c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c1",
            "digest": "sha256:61a1dcd016024d7c760723c25784534084b7b99c3166193df549af8510738fa5",
            "excerpt": "pass: bench test --package ./internal/...; 94 packages, zero failures"
          },
          "requirement": "fixture-consumers",
          "command": "bench test --package ./internal/...",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "hp-c1-standards-1",
          "performer": "/root/hpc1_standards_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "497dea0db26729d73da34ec21330abb99d10a6c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc1_standards_review",
            "digest": "sha256:7dc2acf9275bbe74ea728258c85e452c5490885be1dc85d4c8eaab4a7a93adef",
            "excerpt": "result: completed; axis: Standards; findings: 0; tip: 42d44425656d170d22be7ba3f4de8277df84899d"
          },
          "axis": "Standards",
          "base": "01d8b3c3fded6dc1bff7768faae1d4fbd48d2ae7",
          "tip": "42d44425656d170d22be7ba3f4de8277df84899d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "hp-c1-spec-1",
          "performer": "/root/hpc1_spec_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "497dea0db26729d73da34ec21330abb99d10a6c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc1_spec_review",
            "digest": "sha256:8f9d8b3b92ca46a59f111361fac076015623688796d9a09592d2cd13810ef092",
            "excerpt": "result: completed; axis: Spec; findings: 0; tip: 42d44425656d170d22be7ba3f4de8277df84899d"
          },
          "axis": "Spec",
          "base": "01d8b3c3fded6dc1bff7768faae1d4fbd48d2ae7",
          "tip": "42d44425656d170d22be7ba3f4de8277df84899d",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "hp-c1-coverage-1",
          "performer": "/root/hpc1_coverage_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "497dea0db26729d73da34ec21330abb99d10a6c6",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc1_coverage_review",
            "digest": "sha256:747aeb31f85bd0c17420b69090035a113fd3c5fb497dfda39dc24df4b3b2fa16",
            "excerpt": "result: completed; axis: Coverage; findings: 0; tip: 42d44425656d170d22be7ba3f4de8277df84899d"
          },
          "axis": "Coverage",
          "base": "01d8b3c3fded6dc1bff7768faae1d4fbd48d2ae7",
          "tip": "42d44425656d170d22be7ba3f4de8277df84899d",
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
  }
}
```
