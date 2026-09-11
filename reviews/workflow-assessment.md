# Workflow assessment review

Chunk 1 frozen pair: `48da9cdf0aabebe39eac7131122ff9d4448b8167..0b9ff07536a7c7862d9a07ce6282703e4e72c559`.

Eight raw findings and eight accepted repair targets. User waived the additional cross-harness pass: “No just have the SOL review”. Implementation and repairs remain in the retained author session.

Code citations below are relative to internal/assessment; spec citations refer to specs/workflow-assessment/spec.md.

## Standards

3 findings; worst severity high.

- S1 — high; auto-fix. Reuse canonical JSON parser: store.go:59-78 duplicates internal/jsonfile/decode.go:13-98; AGENTS.md:34-47.

- S2 — medium; auto-fix. Centralize schema vocabulary: validate.go:28-31 and README.md:17-29; AGENTS.md:34-47.

- S3 — medium; auto-fix. Demonstrate remaining independent expectations with biting reds and observable idempotence: record_test.go:178-235; spec.md:423-427; AGENTS.md:41-47.

## Spec

3 findings; worst severity high.

- P1 — high; auto-fix. Cumulative snapshots 10, absent, 15 produce 25: record.go:100-110; spec.md:82,134.

- P2 — medium; auto-fix. Run-wide ambiguity contaminates unrelated attempts: record.go:157-168; spec.md:130.

- P3 — medium; auto-fix. Quality and applicable unknown charges lack provenance: validate.go:60; cost.go:69; spec.md:88,139.

## Coverage

2 findings; worst severity high.

- C1 — high; auto-fix. Epoch 0 then 1 then 0 is counted as complete: record.go:70; spec.md:82.

- C2 — medium; auto-fix. ESC/BEL task and chunk IDs poison query output: validate.go:12; command.go:77; projects/benchkit.md:183.

All three native Sol reviewers inspected read-only isolated worktrees. Native usage and costs are unknown. No delegated tests, probes, or authorship ran.

```bench-review-record
{
  "version": 1,
  "spec": "specs/workflow-assessment/spec.md",
  "plan_digest": "sha256:bebc30452041719f1f95abf6b093687e13ad5ca87665ba14a62a8d50bfb4f811",
  "implementation_session": "01a0920d-3021-73c1-9ed3-9980d5decc71",
  "chunks": [
    {
      "id": "1",
      "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
      "tip": "d5a1dc4728327dcb94317c7dfea24712b98a1e7b",
      "plan_digest": "sha256:bebc30452041719f1f95abf6b093687e13ad5ca87665ba14a62a8d50bfb4f811",
      "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
      "acceptance_rows": [
        "A1",
        "A2",
        "A3",
        "A4",
        "A5",
        "A6",
        "A7",
        "A8",
        "A9",
        "A10",
        "A11",
        "A12",
        "A14",
        "A15",
        "A16",
        "A17",
        "A24",
        "A25",
        "A26",
        "A27",
        "A28",
        "A29"
      ],
      "verification": [
        {
          "id": "chunk1-repair-assessment",
          "performer": "01a0920d-3021-73c1-9ed3-9980d5decc71",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "medium",
          "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native-tool:983108",
            "excerpt": "github.com/gibbonmi/bench/internal/assessment,pass,656",
            "digest": "sha256:af3c7bc53e4875fa2d77cf40ed91ccedd6bc802f4985d58bcf61f187fcbfcfa5"
          },
          "requirement": "assessment",
          "command": "bench test --package ./internal/assessment",
          "exit_code": 0
        },
        {
          "id": "chunk1-repair-dispatcher",
          "performer": "01a0920d-3021-73c1-9ed3-9980d5decc71",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "medium",
          "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native-tool:c72885",
            "excerpt": "github.com/gibbonmi/bench/cmd/bench,pass,14633",
            "digest": "sha256:cf62d7d1627125f9733599e5a1f6eae9aeb5fae0bfa08b92b2f044a0aa81daac"
          },
          "requirement": "dispatcher",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "chunk1-repair-cache-probe",
          "performer": "01a0920d-3021-73c1-9ed3-9980d5decc71",
          "role": "author-verification",
          "model": "gpt-6-astra",
          "effort": "medium",
          "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native-tool:9fe7da",
            "excerpt": "package,./internal/assessment,TestAssessmentRecord,passed,34",
            "digest": "sha256:98d0fffec2d8f6cc61c7dc3d0930048594af5dfb711ada738d13eaaf60d1f055"
          },
          "requirement": "cache-probe",
          "command": "bench test --package ./internal/assessment --run TestAssessmentRecord",
          "exit_code": 0,
          "probe": {
            "mutation": "omit cached-input subtraction",
            "outcome": "bit",
            "exit_code": 1,
            "restore": "pass",
            "native_ref": {
              "ref": "native-tool:9fe7da",
              "excerpt": "bit,internal/assessment/record.go,swap,failed,1,yes",
              "digest": "sha256:f1551004eb42781e7eb5bfe06363bf276267178abd689b3052585165095023b7"
            }
          }
        }
      ],
      "reviews": [
        {
          "id": "chunk1-standards-1",
          "performer": "/root/assessment_1_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "2abd08dbd42c4a84b076f8826bfc4e7e6f90fd5f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native-agent:/root/assessment_1_standards",
            "digest": "sha256:3210ee28b740ad9eba651a132d045cd269954e9c025da6f279b653b0a2a1abf9",
            "excerpt": "Completed with 3 findings. Raw count: 3. Worst severity: high."
          },
          "axis": "Standards",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "0b9ff07536a7c7862d9a07ce6282703e4e72c559",
          "finding_ids": [
            "S1",
            "S2",
            "S3"
          ],
          "supersedes": []
        },
        {
          "id": "chunk1-spec-1",
          "performer": "/root/assessment_1_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "2abd08dbd42c4a84b076f8826bfc4e7e6f90fd5f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native-agent:/root/assessment_1_spec",
            "digest": "sha256:f279cf499647876cfc6a6bd4bda309931adebe464c9202325e701e623c4400e9",
            "excerpt": "Spec review completed with **3 findings**. Raw count: **3**. Worst issue: **high severity**. Disposition totals: **auto-fix 3, ask-user 0, no-op 0**."
          },
          "axis": "Spec",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "0b9ff07536a7c7862d9a07ce6282703e4e72c559",
          "finding_ids": [
            "P1",
            "P2",
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "chunk1-coverage-1",
          "performer": "/root/assessment_1_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "2abd08dbd42c4a84b076f8826bfc4e7e6f90fd5f",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native-agent:/root/assessment_1_coverage",
            "digest": "sha256:439bc8a18b04efea56f7a8067c26ffc38604326d3b3a96d667e28f3cfefbc6bc",
            "excerpt": "Coverage review completed with **2 findings**."
          },
          "axis": "Coverage",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "0b9ff07536a7c7862d9a07ce6282703e4e72c559",
          "finding_ids": [
            "C1",
            "C2"
          ],
          "supersedes": []
        },
        {
          "id": "chunk1-standards-2",
          "performer": "/root/assessment_1_standards",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "native-agent:/root/assessment_1_standards:reaffirmation-1",
            "excerpt": "Completed/pass with **zero findings**.",
            "digest": "sha256:941bc710f2748e9bab14e94cee830d99a31ecf015650cfd8958947ed46bf9e92"
          },
          "axis": "Standards",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "d5a1dc4728327dcb94317c7dfea24712b98a1e7b",
          "finding_ids": [],
          "supersedes": [
            "chunk1-standards-1"
          ]
        },
        {
          "id": "chunk1-spec-2",
          "performer": "/root/assessment_1_spec",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native-agent:/root/assessment_1_spec:reaffirmation-1",
            "excerpt": "Spec reaffirmation completed with **2 later-delta findings**. The original P1\u2013P3 findings are closed.",
            "digest": "sha256:3aacbfbfc6fc8f1bc03350ec5a08af30ea628f3fd5afe16b1245855d5609c01d"
          },
          "axis": "Spec",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "d5a1dc4728327dcb94317c7dfea24712b98a1e7b",
          "finding_ids": [
            "P4",
            "P5"
          ],
          "supersedes": [
            "chunk1-spec-1"
          ]
        },
        {
          "id": "chunk1-coverage-2",
          "performer": "/root/assessment_1_coverage",
          "role": "independent-review",
          "model": "gpt-5.6-sol",
          "effort": "high",
          "source_digest": "795f39f78f7da2633842d66d5aef5d5ff7105cd4",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "native-agent:/root/assessment_1_coverage:reaffirmation-1",
            "excerpt": "**Completed with 1 actionable later-delta finding.** Raw count: **1**. Worst severity: **medium**. Disposition: **auto-fix**.",
            "digest": "sha256:fde941d4f8a1aa6b1c26020b13ed590374268959dfdc5fd17e883e3bebd4db0c"
          },
          "axis": "Coverage",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "d5a1dc4728327dcb94317c7dfea24712b98a1e7b",
          "finding_ids": [
            "C3"
          ],
          "supersedes": [
            "chunk1-coverage-1"
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
      "from": "sha256:c0cbc6c20f7fe0cb79cf83be763055af1e73a0184914a17dbded87a2a908e85a",
      "to": "sha256:bebc30452041719f1f95abf6b093687e13ad5ca87665ba14a62a8d50bfb4f811",
      "chunk_ids": {
        "1": [
          "1"
        ],
        "2": [
          "2"
        ],
        "3": [
          "3"
        ]
      }
    }
  ]
}
```

Repair verification targets the current frozen tip. The plan amendment adds repair ticket 4 inside chunk 1 and preserves all three stable chunk IDs. Sol reaffirmations are pending.

## Standards

Reaffirmation: zero findings; none remain from S1-S3.

## Spec

Reaffirmation: two findings; worst severity medium. P1-P3 are closed.

- P4 — medium; auto-fix. Epoch regression also invalidates independent delta quantities. Citations: internal/assessment/record.go:72,91; spec.md:82.
- P5 — low; auto-fix. Add time_reference to the exhaustive Run and Attempt field lists. Citations: spec.md:76,301; internal/assessment/types.go:58,74.

## Coverage

Reaffirmation: one finding; worst severity medium. C1-C2 are closed.

- C3 — medium; auto-fix. Delta events with epochs 1 then 0 lose their known total. Add an independent delta partition. Citations: internal/assessment/record.go:72-76,133-138; spec.md:80-82.

This round has three raw findings and two repair targets: P4 and C3 name the same repair. All are accepted within the approved scope.
