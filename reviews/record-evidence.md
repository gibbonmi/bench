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
      "verification": [
        {
          "id": "re-c1-v-reviewrecord",
          "performer": "claude:bench-writer/re-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "688d9e5eb06e01776f1b3a2ff14d7c96781c707e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-author-20261002/verify-1-reviewrecord",
            "digest": "sha256:2871ae29133d5fcd9cf1dbe4f3662a4571a4571d329a4f80c45c3e43467947e9",
            "excerpt": "packages[2]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,1591\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-reviewrecord",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "re-c1-v-preflight",
          "performer": "claude:bench-writer/re-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "688d9e5eb06e01776f1b3a2ff14d7c96781c707e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-author-20261002/verify-1-preflight",
            "digest": "sha256:3198f19f4b0c51a895f0b3f57f07bcb3afbec77181d36964aa7f78a75dacd2bf",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/preflight,pass,17799\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-preflight",
          "command": "bench test --package ./internal/preflight",
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

## RE-C1 freeze

The orchestrator froze RE-C1 after ticket 1, with base `11aeb8e3` and tip `2af46532`. The coordinator probe swapped the duplicate-fence guard in the `parse.go` locator, and `bench probe` returned `bit` with 2 failed tests and `restored=yes`.

## RE-C1 ticket 1 author evidence

The author is `claude:bench-writer/re-t1-author` on opus at medium effort. The author used 1 of 3 attempts. Each probe below ran through `bench probe` on `internal/reviewrecord/write.go`, returned `bit`, and returned `restored=yes`.

| Row | Test | Red route | Status |
|---|---|---|---|
| RE1 | `TestRenderKeepsProseAroundTheFence` | Probe: swap the suffix `document[span.end:]` for a bare closing line. | verified |
| RE2 | `TestRenderAppendsAFenceAfterUnterminatedProse` | Probe: omit the added final newline. | verified |
| RE3 | `TestRenderStartsANewDocument` | Probe: swap the heading for an empty string. | verified |
| RE4 | `TestRenderRoundTripsAHostileExcerpt` | Probe: swap the payload so that it writes a raw tab. `Read` refused the JSON. | verified |
| RE5 | `TestRenderWritesAReadablePayload` | Probe: omit `SetEscapeHTML(false)`. | verified |
| RE6 | `TestRenderRefusesAnUnterminatedFence` | Probe: swap the refusal guard so that Render appends a fence. | verified |
| RE102 | `TestRenderRefusesADuplicateFence` | Probe: the same swap as RE6. | verified |
| RE103 | `TestFixtureSaveRendersThroughRender` | Red before the first production edit, then green. | verified |
| RE9 | `TestReviewRecordSource`, with no assertion changed | Probe: swap the closing line for an unparseable one. The fence became unterminated. | verified |
| RE10 | `TestDelegatedEvidenceProjection`, with no assertion changed | Probe: make Render call `Parse` first. The cases for version 9 and the smuggled identity failed. | verified |
| RE7 | The RE7 search | The search found no hit. | verified |
| RE8 | The RE8 search | The search found one hit, the `recordFence` constant in `files.go`. | verified |

The first red of the `record_test.go` rows was a compile failure, because `Render` did not exist. The probes above give the behavioral reds.

The central probe omits `SetEscapeHTML(false)` and runs the whole `internal/reviewrecord` package. It failed 2 tests, RE5 and RE103.

The duplicated-facts sweep found the fence delimiter in two places, the locator and `Render`. One constant, `fenceMarker`, now holds it beside the locator. The comment sweep removed a comment that stated what the code does and a paragraph that repeated the doc of `Render`.

Focused checks at the chunk tip:

- `bench test --package ./internal/reviewrecord/...` passed in 1512 ms.
- `bench test --package ./internal/preflight` passed in 17639 ms.
- `TestRootConformance` passed in 5854 ms.
- `bench structure` reported no issue in `internal/reviewrecord`. The lane check of structure growth passed.

No check skipped a test.
