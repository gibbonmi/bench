# Record evidence review record

```bench-review-record
{
  "version": 2,
  "spec": "specs/record-evidence/spec.md",
  "plan_digest": "sha256:eedace2d50c0412d6dbe225f07c12e8c17ebf3575c1e5f6f19a53dee8a558012",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RE-C1",
      "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
      "tip": "2af46532eb386bc7e4eae65cc504e79e8243e5d7",
      "plan_digest": "sha256:eedace2d50c0412d6dbe225f07c12e8c17ebf3575c1e5f6f19a53dee8a558012",
      "source_digest": "688d9e5eb06e01776f1b3a2ff14d7c96781c707e",
      "acceptance_rows": [
        "RE1",
        "RE2",
        "RE3",
        "RE4",
        "RE5",
        "RE6",
        "RE7",
        "RE8",
        "RE9",
        "RE10",
        "RE102",
        "RE103"
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
  }
}
```

## RE-C1 freeze

The orchestrator froze RE-C1 after ticket 1, with base `11aeb8e3` and tip `2af46532`. The coordinator probe swapped the duplicate-fence guard in the `parse.go` locator, and `bench probe` returned `bit` with 2 failed tests and `restored=yes`.
