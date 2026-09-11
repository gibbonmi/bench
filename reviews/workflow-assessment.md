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
  "plan_digest": "sha256:c0cbc6c20f7fe0cb79cf83be763055af1e73a0184914a17dbded87a2a908e85a",
  "implementation_session": "01a0920d-3021-73c1-9ed3-9980d5decc71",
  "chunks": [
    {
      "id": "1",
      "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
      "tip": "0b9ff07536a7c7862d9a07ce6282703e4e72c559",
      "plan_digest": "sha256:c0cbc6c20f7fe0cb79cf83be763055af1e73a0184914a17dbded87a2a908e85a",
      "source_digest": "2abd08dbd42c4a84b076f8826bfc4e7e6f90fd5f",
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
      "verification": [],
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
