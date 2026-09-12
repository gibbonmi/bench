# Session context measurement review

## Chunk ME-C1

The frozen pair is base `48da9cdf0aabebe39eac7131122ff9d4448b8167` and tip `eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2`.
Three Opus/medium axes ran on 2026-09-11 in separate native contexts and read-only venues.
Raw findings: Standards 4, Spec 4, Coverage 4.
De-duplicated repair targets: 7, of which 6 are auto-fix and 1 was decided by the reviewer at review time.
The record chunk names base `8aac4f9d6f40232ed14803b336501f948ebf905f`, the main tip merged into the source after the review, and the current source tip.

The blast walk found five `touched=false` consumers.
All five read `bench.commandRegistry`: `cmd/bench/command_registry.go:245`, `:261`, `:289`, `cmd/bench/command_registry_test.go:69`, and `cmd/bench/otel_hook_seams_test.go:18`.
The author's green `cmd/bench` package run exercises each of them.

## Standards

Findings: 4. Worst issue: medium.

- ST1 (medium, auto-fix after the reviewer's decision): the new package re-declares the three control-record state spellings that `internal/bounds.FileState` owns. The reviewer chose the one-source shape: `Failure.State` becomes `bounds.FileState`. The reader keeps its own streaming read and `Lstat` check. Citations at tip eb51bc6f: internal/harnesstranscript/transcript.go:30-35; internal/bounds/classify.go:16-21; internal/harnesses/command.go:121.
- ST2 (low, auto-fix): the help suffix `[<harness>] [--record <path> --format <source-id>]` reads as an optional harness. The verb and the spec require one for a record call. Citations at tip eb51bc6f: cmd/bench/main.go:80; cmd/bench/help_inventory_test.go:60; internal/harnesses/command.go:87, :140-142; specs/session-context-measurement/spec.md:51.
- ST3 (low, no-op): the compiled view names `tool calls` and `Read paths`; the record view names `outer calls` and `explicit read paths`. The compiled names are FT204 supplier placeholders outside the fence, and ME11 keeps both compiled views byte-identical. Citations at tip eb51bc6f: internal/harnesses/harnesses.go:119-120.
- ST4 (low, no-op): the independence paragraph appears in both test files, and the inline compiled-view baseline copies every record cell. AGENTS.md permits the independent expectation. The author's `render("harness")` probe recorded the baseline's red. Citations at tip eb51bc6f: internal/harnesses/observed_test.go:14-16, :305-400; internal/harnesstranscript/transcript_test.go:9-11.

## Spec

Findings: 4. Worst issue: medium.

- SPEC-1 (medium, auto-fix): the dogfood run over the pinned record was recorded in no spec artifact. The coordinator repaired this with a Dogfood runs note under Further notes. Citations at tip eb51bc6f: specs/session-context-measurement/spec.md:57.
- SPEC-2 (medium, auto-fix): a usage snapshot that omits one token key reports that dimension as an observed zero. The spec says a missing dimension stays unknown. COV-4 names the same defect. Citations at tip eb51bc6f: specs/session-context-measurement/spec.md:69, :74; internal/harnesstranscript/codex.go:99-104, :308-314; internal/harnesses/observed_test.go:264-272.
- SPEC-3 (low, no-op): ME18's second clause is asserted by the absence of the literal `complete producer output`. The unmatched-call half closes correctly, and the boundary sentence discharges the claim. Citations at tip eb51bc6f: internal/harnesses/observed_test.go:236; internal/harnesstranscript/codex.go:67.
- SPEC-4 (low, auto-fix): two provenance cells misdescribe the source. The `explicit read paths` cell says `no pinned shape` while its boundary names shell text, and the `observed nested calls` cell names a shape the reader never decodes. COV-3 shares the repair. Citations at tip eb51bc6f: internal/harnesstranscript/codex.go:53, :56, :320, :328.

## Coverage

Findings: 4. Worst issue: medium.

- COV-1 (medium, auto-fix): the string-form `function_call_output.output` path reaches no test. A probe that drops its count stayed silent over 48 tests, and the pinned record holds 77 such results. The repair adds one fixture with a `function_call` and its string-form completion. Citations at tip eb51bc6f: internal/harnesstranscript/codex.go:233-238; internal/harnesses/observed_test.go:24-26.
- COV-2 (medium, decided): a line past `maxRecordLine` stops the scan, so later turns, compactions, and the end timestamp vanish. The reviewer decided that the reader skips the oversized line, marks the affected measures incomplete, and keeps reading. The spec's edge inventory records the decision. Citations at tip eb51bc6f: internal/harnesstranscript/codex.go:40-43, :157-167; internal/harnesses/command.go:57.
- COV-3 (low, auto-fix): the result-text provenance cell names only the array field. It shares COV-1's repair. Citations at tip eb51bc6f: internal/harnesstranscript/codex.go:49, :316-318.
- COV-4 (low, no-op as a separate target): a snapshot missing one counter key reads as observed 0. SPEC-2's repair covers it. Citations at tip eb51bc6f: internal/harnesstranscript/codex.go:308-314; internal/harnesstranscript/transcript.go:13.

## Record

```bench-review-record
{
  "version": 1,
  "spec": "specs/session-context-measurement/spec.md",
  "plan_digest": "sha256:13883eacbdd6a65cd0df63f6e087ce509ad0be8584aec32c285d5cf8fc3dc748",
  "implementation_session": "claude:opus-high:ticket-author",
  "chunks": [
    {
      "id": "ME-C1",
      "base": "8aac4f9d6f40232ed14803b336501f948ebf905f",
      "tip": "6d981dca4d04045634da730482ac8f0cd80c50fc",
      "plan_digest": "sha256:13883eacbdd6a65cd0df63f6e087ce509ad0be8584aec32c285d5cf8fc3dc748",
      "source_digest": "48d86d061fd95c6a09fb77a7fd5f19e9a00cde01",
      "acceptance_rows": [
        "ME1", "ME2", "ME3", "ME4", "ME5", "ME6", "ME7", "ME8", "ME9", "ME10", "ME11", "ME12",
        "ME17", "ME18", "ME19", "ME20", "ME21", "ME22", "ME23", "ME24", "ME25"
      ],
      "verification": [
        {
          "id": "me-c1-harness-tests",
          "performer": "claude:opus-high:ticket-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "",
          "state": "pending",
          "outcome": "",
          "native_ref": {"ref": "", "digest": "", "excerpt": ""},
          "requirement": "harness-tests",
          "command": "bench test --package ./internal/harnesses",
          "exit_code": null
        },
        {
          "id": "me-c1-reader-tests",
          "performer": "claude:opus-high:ticket-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "",
          "state": "pending",
          "outcome": "",
          "native_ref": {"ref": "", "digest": "", "excerpt": ""},
          "requirement": "reader-tests",
          "command": "bench test --package ./internal/harnesstranscript",
          "exit_code": null
        },
        {
          "id": "me-c1-command-tests",
          "performer": "claude:opus-high:ticket-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "",
          "state": "pending",
          "outcome": "",
          "native_ref": {"ref": "", "digest": "", "excerpt": ""},
          "requirement": "command-tests",
          "command": "bench test --package ./cmd/bench",
          "exit_code": null
        }
      ],
      "reviews": [
        {
          "id": "me-c1-standards-1",
          "performer": "claude:opus-medium:standards-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "edd67cd0cd0b49322212631b55646023a053170a",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude:agent/me-c1-standards-1",
            "digest": "sha256:cd6dce7013d5cb5132130de12922ca71d593afd901654b9df8dd1ba1b7ea61a5",
            "excerpt": "result: completed; axis: Standards; findings: 4; worst: medium; tip: eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2"
          },
          "axis": "Standards",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2",
          "finding_ids": ["ST1", "ST2", "ST3", "ST4"],
          "supersedes": []
        },
        {
          "id": "me-c1-spec-1",
          "performer": "claude:opus-medium:spec-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "edd67cd0cd0b49322212631b55646023a053170a",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude:agent/me-c1-spec-1",
            "digest": "sha256:7f141b9b0425ccc052552cef19f87e24b9564aa000b12b46a99a97728cb0d85f",
            "excerpt": "result: completed; axis: Spec; findings: 4; worst: medium; tip: eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2"
          },
          "axis": "Spec",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2",
          "finding_ids": ["SPEC-1", "SPEC-2", "SPEC-3", "SPEC-4"],
          "supersedes": []
        },
        {
          "id": "me-c1-coverage-1",
          "performer": "claude:opus-medium:coverage-axis",
          "role": "independent-review",
          "model": "opus",
          "effort": "medium",
          "source_digest": "edd67cd0cd0b49322212631b55646023a053170a",
          "state": "completed",
          "outcome": "findings",
          "native_ref": {
            "ref": "claude:agent/me-c1-coverage-1",
            "digest": "sha256:fa076b6464ba49782644ed654ee3021f8bf6e834d06ebecc2496e711c0f251bd",
            "excerpt": "result: completed; axis: Coverage; findings: 4; worst: medium; tip: eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2"
          },
          "axis": "Coverage",
          "base": "48da9cdf0aabebe39eac7131122ff9d4448b8167",
          "tip": "eb51bc6fd4545c4b0da7b3f294e755d35be0a6f2",
          "finding_ids": ["COV-1", "COV-2", "COV-3", "COV-4"],
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
