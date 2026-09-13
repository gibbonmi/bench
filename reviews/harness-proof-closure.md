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
  "plan_digest": "sha256:5386adf22ec42126210eac16c632653aec0980a91e5f4205a91daa6b0d1dcfdd",
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
    },
    {
      "id": "HP-C2",
      "base": "4dd18eeca9d43048d35be82dcc0eb87ef1360bfe",
      "tip": "0362040e173b37e371e0cb8f45cd324382a076fb",
      "plan_digest": "sha256:9d9dfa2e2789cc8782246b8fc697718828baca8a3a3fd2a7d1e722a77e65052a",
      "source_digest": "fcb28ac247a76bd20af6b9130b3610c3246e5377",
      "acceptance_rows": ["HP4", "HP5"],
      "verification": [
        {
          "id": "hp-c2-root-entry",
          "performer": "/root/ft120_c2",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "fcb28ac247a76bd20af6b9130b3610c3246e5377",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c2",
            "digest": "sha256:f01117bd505732ea02611da315f27513c555f8577cfc0a4ba892d6397e9595a1",
            "excerpt": "pass: bench test --package ./internal/conformance --run 'TestRootConformance|TestHarnessDefaultsToCurrentGitRoot'; 0 failures, 0 skips"
          },
          "requirement": "root-entry",
          "command": "bench test --package ./internal/conformance --run 'TestRootConformance|TestHarnessDefaultsToCurrentGitRoot'",
          "exit_code": 0,
          "probe": {
            "mutation": "restore the unset-root environment skip",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "collaboration:/root/ft120_c2",
              "digest": "sha256:9d7d99b49fb8be1b617f618d72dae8203e2aa9c68ab1305a9eb40608f9cdda85",
              "excerpt": "bit: restored unset-root skip caused 2 failures; source restored and focused selection passed"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "hp-c2-standards-1",
          "performer": "/root/hpc2_standards_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "fcb28ac247a76bd20af6b9130b3610c3246e5377",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc2_standards_review",
            "digest": "sha256:db92fff372745b6208821a844ec64283805e0f9bc4d22d80b1b5cb3c376f3364",
            "excerpt": "result: completed; axis: Standards; findings: 0; tip: 0362040e173b37e371e0cb8f45cd324382a076fb"
          },
          "axis": "Standards",
          "base": "4dd18eeca9d43048d35be82dcc0eb87ef1360bfe",
          "tip": "0362040e173b37e371e0cb8f45cd324382a076fb",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "hp-c2-spec-1",
          "performer": "/root/hpc2_spec_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "fcb28ac247a76bd20af6b9130b3610c3246e5377",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc2_spec_review",
            "digest": "sha256:5f8863fc67e16615c34a7391d81ade7f90a949fd52a4317c8f0ba5ab74b102ac",
            "excerpt": "result: completed; axis: Spec; findings: 0; tip: 0362040e173b37e371e0cb8f45cd324382a076fb"
          },
          "axis": "Spec",
          "base": "4dd18eeca9d43048d35be82dcc0eb87ef1360bfe",
          "tip": "0362040e173b37e371e0cb8f45cd324382a076fb",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "hp-c2-coverage-1",
          "performer": "/root/hpc1_spec_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "fcb28ac247a76bd20af6b9130b3610c3246e5377",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc1_spec_review",
            "digest": "sha256:d19ab2cffbd6d14587a601994938b4578cc61f07e970d3345acef84ebea29a1d",
            "excerpt": "source_tip: 0362040e173b37e371e0cb8f45cd324382a076fb\npackages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,30\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "axis": "Coverage",
          "base": "4dd18eeca9d43048d35be82dcc0eb87ef1360bfe",
          "tip": "0362040e173b37e371e0cb8f45cd324382a076fb",
          "finding_ids": [],
          "supersedes": []
        }
      ]
    },
    {
      "id": "HP-C3",
      "base": "d1b34aa2c1dc3bf599851ac3f8df75c31827d710",
      "tip": "176a9a5e3449d7ebec0b40183be3549b1602a516",
      "plan_digest": "sha256:5386adf22ec42126210eac16c632653aec0980a91e5f4205a91daa6b0d1dcfdd",
      "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
      "acceptance_rows": ["HP6", "HP7", "HP8", "HP9", "HP10"],
      "verification": [
        {
          "id": "hp-c3-anchor-owner",
          "performer": "/root/ft120_c3",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c2",
            "digest": "sha256:dff05546c3ac85735b6b7f7bdd6a4dd732482eeb4bbce616edb56fc532ea3d4d",
            "excerpt": "pass: bench test --package ./internal/anchors; 0 failures, 0 skips"
          },
          "requirement": "anchor-owner",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0,
          "probe": {
            "mutation": "return locations without the path's registry diagnostics",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "collaboration:/root/ft120_c2",
              "digest": "sha256:77867608714fffb54b2bce1985f415d25c2be9abe37b0bbc615ce7963bcbc4c7",
              "excerpt": "bit: diagnostics-omission mutation failed 9 of 11 tests; restored yes; 0 skips"
            }
          }
        },
        {
          "id": "hp-c3-anchor-command",
          "performer": "/root/ft120_c3",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c2",
            "digest": "sha256:c0a05e6d1b385d50fc060ea96a053e8f7d7bc111239f7a466817f3081fa688b9",
            "excerpt": "pass: bench test --package ./cmd/bench --run TestAnchors; 0 failures, 0 skips"
          },
          "requirement": "anchor-command",
          "command": "bench test --package ./cmd/bench --run TestAnchors",
          "exit_code": 0
        },
        {
          "id": "hp-c3-anchor-lane",
          "performer": "/root/ft120_c3",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c2",
            "digest": "sha256:dd3edeef35e98c77e62f7a213119a3723eba31029d3a7ae01a25c4b9f26e1d49",
            "excerpt": "pass: bench test --package ./internal/gate --run TestSelectLaneByClass; 0 failures, 0 skips"
          },
          "requirement": "anchor-lane",
          "command": "bench test --package ./internal/gate --run TestSelectLaneByClass",
          "exit_code": 0
        },
        {
          "id": "hp-c3-docs-check",
          "performer": "/root/ft120_c3",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/ft120_c2",
            "digest": "sha256:368d0cec92a4dbb90a102f2725ad895480b3071d4c8ac610c3dd8cf662b93f47",
            "excerpt": "pass: bench test --check docs-currency-workflow; 0 failures, 0 skips"
          },
          "requirement": "docs-check",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "hp-c3-standards-1",
          "performer": "/root/hpc2_standards_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc2_standards_review",
            "digest": "sha256:25986fd2301533cc44f8d40519f14b4e2c2b10805f94d41ddf1aa673210c9c01",
            "excerpt": "result: completed; axis: Standards; findings: 0; tip: 176a9a5e3449d7ebec0b40183be3549b1602a516"
          },
          "axis": "Standards",
          "base": "d1b34aa2c1dc3bf599851ac3f8df75c31827d710",
          "tip": "176a9a5e3449d7ebec0b40183be3549b1602a516",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "hp-c3-spec-1",
          "performer": "/root/hpc2_spec_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc2_spec_review",
            "digest": "sha256:50b90510a3a75a66a9bd0359782cead8261638fb6a1702f2adfb963bfb7c5345",
            "excerpt": "result: completed; axis: Spec; findings: 0; tip: 176a9a5e3449d7ebec0b40183be3549b1602a516"
          },
          "axis": "Spec",
          "base": "d1b34aa2c1dc3bf599851ac3f8df75c31827d710",
          "tip": "176a9a5e3449d7ebec0b40183be3549b1602a516",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "hp-c3-coverage-1",
          "performer": "/root/hpc1_spec_review",
          "role": "independent-review",
          "model": "gpt-5.6-terra",
          "effort": "high",
          "source_digest": "5e5d0986651f9c7723a9144846d64453a67cd7cd",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "collaboration:/root/hpc1_spec_review",
            "digest": "sha256:73773f84b269199a9fbe51a2784c27737276c0ee191e16390a7b958f5ea11489",
            "excerpt": "result: completed; axis: Coverage; findings: 0; tip: 176a9a5e3449d7ebec0b40183be3549b1602a516"
          },
          "axis": "Coverage",
          "base": "d1b34aa2c1dc3bf599851ac3f8df75c31827d710",
          "tip": "176a9a5e3449d7ebec0b40183be3549b1602a516",
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
      "from": "sha256:9d9dfa2e2789cc8782246b8fc697718828baca8a3a3fd2a7d1e722a77e65052a",
      "to": "sha256:5386adf22ec42126210eac16c632653aec0980a91e5f4205a91daa6b0d1dcfdd",
      "chunk_ids": {
        "HP-C1": ["HP-C1"],
        "HP-C2": ["HP-C2"],
        "HP-C3": ["HP-C3"],
        "HP-C4": ["HP-C4"],
        "HP-C5": ["HP-C5"]
      }
    }
  ]
}
```
