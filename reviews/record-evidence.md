# Record evidence review record

```bench-review-record
{
  "version": 2,
  "spec": "specs/record-evidence/spec.md",
  "plan_digest": "sha256:833ccd1de096637009430bda209d327b710beca8fea9eaf774250694705d73fc",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RE-C1",
      "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
      "tip": "6d7f3969e9061e8022630706c4430f64a7248312",
      "plan_digest": "sha256:91c7c928f595b8551cae347d20bbf32ab50668ff972c6c398aa619324e16d192",
      "source_digest": "f9f850d71e25d2db0ecc3c4d49a3452fbecea0d8",
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
        },
        {
          "id": "re-c1-v2-reviewrecord",
          "performer": "claude:bench-writer/re-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f9f850d71e25d2db0ecc3c4d49a3452fbecea0d8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-repair-1-20261002/verify-1-reviewrecord",
            "digest": "sha256:51a07bc7cafc63a2a37cbc894795c24afdca90ab3bce333c66ae2476519b601e",
            "excerpt": "packages[2]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,1541\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-reviewrecord",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "re-c1-v2-preflight",
          "performer": "claude:bench-writer/re-t1-repair-1",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "f9f850d71e25d2db0ecc3c4d49a3452fbecea0d8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-repair-1-20261002/verify-1-preflight",
            "digest": "sha256:6bf20f60103de04b25a5e738b9763024cb2d99ffce0d5e21bc38002be92c2ca2",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/preflight,pass,18056\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-preflight",
          "command": "bench test --package ./internal/preflight",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "re-c1-r1-standards",
          "performer": "claude:bench-reviewer/re-c1-r1-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "688d9e5eb06e01776f1b3a2ff14d7c96781c707e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/re-c1-r1-standards-20261002@2af46532eb386bc7e4eae65cc504e79e8243e5d7",
            "digest": "sha256:53e8418a8e135e39231ca577b2bbcf3c27fccb8a55dd41aa378a95a6f24530cc",
            "excerpt": "Standards axis, RE-C1, pair 11aeb8e3..2af46532: fail, 1 finding (S1).\nS1: source_test.go:246 keeps the RE103 red record in a code comment (craft-comments: the spec owns the red record); auto-fix.\nOne fence scanner (locate), one encoder (Render), and one fence-name constant (RE8: one hit) confirmed; the hand-built malformed fences are exempt by spec rule.\nAdvice: the write.go:11 forward reference to \"the write transaction\", the parallel record builder and render harness in record_test.go, and the open-line grammar derived in two places."
          },
          "axis": "Standards",
          "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
          "tip": "2af46532eb386bc7e4eae65cc504e79e8243e5d7",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "re-c1-r1-spec",
          "performer": "claude:bench-reviewer/re-c1-r1-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "688d9e5eb06e01776f1b3a2ff14d7c96781c707e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-c1-r1-spec-20261002@2af46532eb386bc7e4eae65cc504e79e8243e5d7",
            "digest": "sha256:0cac302f58027364f7198d93c583c10795cc47f872173eb3d02557d7ffe9842c",
            "excerpt": "Spec axis RE-C1 at 11aeb8e3..2af46532: pass, no finding IDs.\nRE1-RE10, RE102 and RE103 hold. The RE7 search has no hit. The RE8 search has one hit, the constant at files.go:90.\n`Render` uses the one `locate` scanner in parse.go, with two-space indent, no HTML escape and no validation. `Save` and `recordFence` render through it.\nThe version 2 plan matches the spec's chunk table.\nAdvice only: rename the package constant `recordFence` to avoid sharing a name with the preflight helper."
          },
          "axis": "Spec",
          "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
          "tip": "2af46532eb386bc7e4eae65cc504e79e8243e5d7",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "re-c1-r1-coverage",
          "performer": "claude:bench-reviewer/re-c1-r1-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "688d9e5eb06e01776f1b3a2ff14d7c96781c707e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-c1-r1-coverage-20261002@2af46532eb386bc7e4eae65cc504e79e8243e5d7",
            "digest": "sha256:4828f29b5b501ad80ed931d702adda4991882af0e4d60cb4f66475f130aa87ef",
            "excerpt": "Coverage RE-C1 (fable/high): pass, no finding IDs.\nProbes: plan fence beside the record fence, inline name, trailing-space opener, unterminated-newline closer, and fence-like lines plus U+2028 in the excerpt all round-trip (silent, restored=yes).\nTwo mutations survive with no correctness defect (advice only): the conditional newline at write.go:26 and the exact-match closer at parse.go:250.\nCRLF and empty files are closed by Won't handle and RE31."
          },
          "axis": "Coverage",
          "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
          "tip": "2af46532eb386bc7e4eae65cc504e79e8243e5d7",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "re-c1-r2-standards",
          "performer": "claude:bench-reviewer/re-c1-r2-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "f9f850d71e25d2db0ecc3c4d49a3452fbecea0d8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-c1-r2-standards-20261002@6d7f3969e9061e8022630706c4430f64a7248312",
            "digest": "sha256:a5dbc0647affba241ef5d866d5711e69ec9490fa5be68eff11084c62178847a0",
            "excerpt": "Standards confirming round, RE-C1, repair delta 2af46532..6d7f3969: pass, no finding IDs.\nS1 closed: source_test.go:72-73 keeps only the reason for the external package; craft-comments:49.\nAdvice: the round 1 S1 line cite (246) is wrong; the actual line is 72."
          },
          "axis": "Standards",
          "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
          "tip": "6d7f3969e9061e8022630706c4430f64a7248312",
          "finding_ids": [],
          "supersedes": [
            "re-c1-r1-standards"
          ]
        },
        {
          "id": "re-c1-r2-spec",
          "performer": "claude:bench-reviewer/re-c1-r2-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "f9f850d71e25d2db0ecc3c4d49a3452fbecea0d8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-c1-r2-spec-20261002@6d7f3969e9061e8022630706c4430f64a7248312",
            "digest": "sha256:0abb81f1bd262f7943422f532a500f16d38479ba9b0bf041a7533368ac50cae2",
            "excerpt": "Spec axis RE-C1 confirming round, delta 2af46532..6d7f3969: pass, no finding IDs.\nS1 is closed by a comment-only edit, and RE103's body is unchanged.\nThe plan adds re-t1-repair-1 (opus/medium, user-directed) and leaves chunks, tickets and commands unchanged."
          },
          "axis": "Spec",
          "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
          "tip": "6d7f3969e9061e8022630706c4430f64a7248312",
          "finding_ids": [],
          "supersedes": [
            "re-c1-r1-spec"
          ]
        },
        {
          "id": "re-c1-r2-coverage",
          "performer": "claude:bench-reviewer/re-c1-r2-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "f9f850d71e25d2db0ecc3c4d49a3452fbecea0d8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-c1-r2-coverage-20261002@6d7f3969e9061e8022630706c4430f64a7248312",
            "digest": "sha256:bbba085b35e7dae4d6590228541242a5aaee42018839959da9011e61cd33c21a",
            "excerpt": "RE-C1 round 2 Coverage: pass.\nS1 fold confirmed; RE103 still bites (probe bit, restored=yes).\nFindings: none."
          },
          "axis": "Coverage",
          "base": "11aeb8e316c82b08c3e77be6391309a9ffd0fb9b",
          "tip": "6d7f3969e9061e8022630706c4430f64a7248312",
          "finding_ids": [],
          "supersedes": [
            "re-c1-r1-coverage"
          ]
        }
      ]
    },
    {
      "id": "RE-C2",
      "base": "6d7f3969e9061e8022630706c4430f64a7248312",
      "tip": "bf3ddc94aa2108fc225d1c946167ea150ae6041c",
      "plan_digest": "sha256:833ccd1de096637009430bda209d327b710beca8fea9eaf774250694705d73fc",
      "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
      "acceptance_rows": [
        "RE11",
        "RE12",
        "RE13",
        "RE14",
        "RE15",
        "RE16",
        "RE17",
        "RE18",
        "RE19",
        "RE20",
        "RE21",
        "RE22",
        "RE23",
        "RE24",
        "RE25",
        "RE26",
        "RE27",
        "RE28",
        "RE29",
        "RE30",
        "RE31",
        "RE32",
        "RE33",
        "RE34",
        "RE35",
        "RE37",
        "RE38",
        "RE39",
        "RE40",
        "RE41",
        "RE105",
        "RE106"
      ],
      "verification": [
        {
          "id": "re-c2-v-reviewrecord",
          "performer": "claude:bench-writer/re-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t2-author-20261002/verify-2-reviewrecord",
            "digest": "sha256:43529b8bc492762f2b4d095fe2bce1858f83ec7486338159d8869d7c30a3d0aa",
            "excerpt": "packages[3]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/reviewrecord,pass,1613\n  github.com/gibbonmi/bench/internal/reviewrecord/recordcmd,pass,1314\n  github.com/gibbonmi/bench/internal/reviewrecord/recordtest,no-tests,0\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-reviewrecord",
          "command": "bench test --package ./internal/reviewrecord/...",
          "exit_code": 0
        },
        {
          "id": "re-c2-v-cmd",
          "performer": "claude:bench-writer/re-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t2-author-20261002/verify-2-cmd",
            "digest": "sha256:09c5964fb55dc55fc358fa14709f625913987e167eada453669b8c356f486c05",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/cmd/bench,pass,13084\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "2-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "re-c2-v-conformance",
          "performer": "claude:bench-writer/re-t2-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "medium",
          "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t2-author-20261002/verify-2-conformance",
            "digest": "sha256:63dc9935c6310261a24f195011318b2b4df41880854cfd399607242ba5b0a5bb",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,35474\nfailures[0]{package,test,line}:\nskips[3]{package,test,reason}:\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceProseBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/JANYB4/t/TestGuidanceProseBudgetRefusesNonRegularSubjectssocket2429873957/001/.agents/skills/bench-craft-linked/SKILL.md: bind: invalid argument\"\n  github.com/gibbonmi/bench/internal/conformance,TestGuidanceSweepRejectsNonRegularEntriesBeforeReading/character_device,\"capability: privilege: cannot create a character device: operation not permitted\"\n  github.com/gibbonmi/bench/internal/conformance,TestSkillDescriptionBudgetRefusesNonRegularSubjects/socket,\"capability: fifo: unix sockets unavailable on this filesystem: listen unix /tmp/JANYB4/t/TestSkillDescriptionBudgetRefusesNonRegularSubjectssocket2420658603/001/.agents/skills/bench-craft-planted/SKILL.md: bind: invalid argument\""
          },
          "requirement": "2-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "re-c2-r1-standards",
          "performer": "claude:bench-reviewer/re-c2-r1-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/re-c2-r1-standards-20261002@bf3ddc94aa2108fc225d1c946167ea150ae6041c",
            "digest": "sha256:e3fb4b8f30ed668470828f5029c7f48e608f24f65110b729f4dfbc758ecbb86a",
            "excerpt": "Standards RE-C2, pair 6d7f3969..bf3ddc94: fail, 1 finding (S1).\nS1: recordcmd/command.go:119 writes its own control-rune check instead of using sanitize.LineSafe (one source per fact; compose an existing seam); auto-fix.\nOne help source (forms -> HelpRows -> formHelpRows), one write transaction, and one refusal path confirmed. The independent TOON field list has its recorded red.\nAdvice: HelpRow is declared in two packages; the spec-path layout is rebuilt by hand at command.go:123; the usage-line shape differs from preflight; the fixture body is repeated at fixture.go:62."
          },
          "axis": "Standards",
          "base": "6d7f3969e9061e8022630706c4430f64a7248312",
          "tip": "bf3ddc94aa2108fc225d1c946167ea150ae6041c",
          "finding_ids": [
            "S1"
          ],
          "supersedes": []
        },
        {
          "id": "re-c2-r1-spec",
          "performer": "claude:bench-reviewer/re-c2-r1-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/re-c2-r1-spec-20261002@bf3ddc94aa2108fc225d1c946167ea150ae6041c",
            "digest": "sha256:9a28e5e2fb185d6871eeea75ca056eaf66001cf0055bf31b9705938c5910d840",
            "excerpt": "Spec axis, RE-C2, frozen pair 6d7f3969..bf3ddc94: approve with one finding.\nAll 32 RE-C2 rows hold, and registry, output, refusal order and budgets match the spec.\nF1 (P3, auto-fix, confidence 5): write.go:964 re-parses the rendered bytes without the reader's whole-document bound, so a document over ControlRecordLimit can be written that Read then refuses.\nThe plan expansion (help-projection allowlist) is in scope."
          },
          "axis": "Spec",
          "base": "6d7f3969e9061e8022630706c4430f64a7248312",
          "tip": "bf3ddc94aa2108fc225d1c946167ea150ae6041c",
          "finding_ids": [
            "P3"
          ],
          "supersedes": []
        },
        {
          "id": "re-c2-r1-coverage",
          "performer": "claude:bench-reviewer/re-c2-r1-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "0fa91f6aaaa602460b5a707256004373bd936892",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/re-c2-r1-coverage-20261002@bf3ddc94aa2108fc225d1c946167ea150ae6041c",
            "digest": "sha256:10b12412c3760fc86f1528638468dcc81d872bcf0e59d68d17dca5da556511ca",
            "excerpt": "Coverage RE-C2 round 1: fail, findings C1 and C2.\nC1: the step 3 control-character refusal on `--chunk` is untested, and the probe was silent.\nC2: the rendered-bytes check bounds the payload, not the whole document, so a write can leave a record over 2 MiB that `Read` refuses. Trace evidence only.\nAdvice: the temp-file cleanup after `CreateTemp` is untested (probe silent, no feasible input); the unknown-form echo is raw.\nBoth probes returned restored=yes."
          },
          "axis": "Coverage",
          "base": "6d7f3969e9061e8022630706c4430f64a7248312",
          "tip": "bf3ddc94aa2108fc225d1c946167ea150ae6041c",
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
  },
  "amendments": [
    {
      "from": "sha256:91c7c928f595b8551cae347d20bbf32ab50668ff972c6c398aa619324e16d192",
      "to": "sha256:833ccd1de096637009430bda209d327b710beca8fea9eaf774250694705d73fc",
      "chunk_ids": {
        "RE-C1": [
          "RE-C1"
        ]
      }
    }
  ]
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

## RE-C1 chunk review, round 1

Three fresh fable / high sessions reviewed the frozen pair `11aeb8e3..2af46532`. Each axis bound the review evidence `sha256:6223b5b6` with `--check-current`. The raw finding count is 1, and the repair-target count is 1. The repair allowance of RE-C1 is 2 cycles, and 0 cycles are used.

### Standards

Finding count: 1. Worst issue: S1.

- S1, auto-fix, confidence 6. The doc comment on `TestFixtureSaveRendersThroughRender` in `internal/reviewrecord/source_test.go` line 72 keeps the RE103 red record. The `craft-comments` skill states that the spec owns the red record. Remove that sentence and keep the reason for the external test package.

### Spec

Finding count: 0. Worst issue: none. Each RE-C1 row holds, and the version 2 plan agrees with the chunk table of the spec.

### Coverage

Finding count: 0. Worst issue: none. The axis ran five probes through `bench probe`, and each probe returned `restored=yes`.

### Advice

- The `Render` doc comment in `write.go` names a write transaction that ticket 2 adds.
- The `renderRecord` and `render` helpers in `record_test.go` are parallel to the `recordtest` fixture, because an internal test cannot import `recordtest`.
- The opening line of a fence comes from the same constants in `parse.go` and `write.go`. A small helper can make it one source.
- The constant `recordFence` has the same name as the preflight test helper.
- No test covers a document with no fence that ends in a newline. No test pins the exact match of the closing line. The code is correct for both edges.

### Command contribution

The Standards axis suggests that the author charge restate the `craft-comments` rule: a red record stays in the spec and not in a code comment. The Spec and Coverage axes found no contribution.

## RE-C1 repair 1

The plan commit `155c28c5` assigned the fresh repair session `claude:bench-writer/re-t1-repair-1` on opus at medium effort. The repair commit `6d7f3969` removes the red-record sentence from the comment on `TestFixtureSaveRendersThroughRender`. The diff changes comment lines only. RE-C1 has used 1 of its 2 repair cycles.

The orchestrator re-froze RE-C1 with base `11aeb8e3` and tip `6d7f3969`. The chunk entry now holds the source digest and the plan digest of that tip, and the earlier results stay in the record.

## RE-C1 repair 1 verification

The repair session `claude:bench-writer/re-t1-repair-1` ran the two verification commands of ticket 1 at commit `2bcc9d21`. The source digest of that run is `f9f850d7`.

- `bench test --package ./internal/reviewrecord/...` passed in 1541 ms.
- `bench test --package ./internal/preflight` passed in 18056 ms.

No check failed, and no check skipped a test.

## RE-C1 chunk review, round 2

Three fresh fable / high sessions ran the confirming round on the repair delta `2af46532..6d7f3969`. Each axis bound the review evidence `sha256:5c762237` with `--check-current`. Each axis passed with no finding, so the raw finding count is 0. S1 is closed. RE-C1 used 1 of its 2 repair cycles.

### Standards, round 2

Finding count: 0. Worst issue: none. The axis found no duplicated knowledge in the repair delta.

### Spec, round 2

Finding count: 0. Worst issue: none. The plan change adds only the repair assignment, and each RE-C1 row still holds.

### Coverage, round 2

Finding count: 0. Worst issue: none. A probe swapped the `Render` call in `Save` for `json.MarshalIndent`. The verdict was `bit` on RE103, with `restored=yes`.

### Advice, round 2

- The round 1 S1 text named line 246, and the comment was at line 72. This record now names line 72. The correction is evidence-only.
- The repair range also holds the record commits that add this file.

## RE-C2 freeze

The orchestrator froze RE-C2 after ticket 2, with base `6d7f3969` and tip `bf3ddc94`. The worktree build of `bench record chunk` wrote the chunk entry, and an independent script gave the same source and plan digests.

The plan commit `6079c896` assigned ticket 2. The ticket 2 author stopped on a fence gap: the help-row projection allowlist in `internal/conformance/help_inventory_single_source_test.go`. The plan commit `4ad69dd9` added that path to the ticket 2 `Writes:` line and to the spec fences. The orchestrator could not regenerate the charge evidence, because the uncommitted delta of the author made the tree dirty. The author received the one changed path in its follow-up charge.

The record amendment from the RE-C1 plan digest to the RE-C2 plan digest maps RE-C1 to itself. The orchestrator wrote it by hand, because ticket 5 adds `bench record amendment`.

The coordinator probe swapped the version 2 plan rule of the chunk change in `write.go`. `bench probe` returned `bit` on RE17 with `restored=yes`.

## RE-C2 ticket 2 author evidence

The author is `claude:bench-writer/re-t2-author` on opus at medium effort. The author used 1 of 4 attempts. The ticket commit is `bf3ddc94`.

Before the implementation, a compile-only stub `Command` returned exit 0 with no output, and `HelpRows` returned no rows. All 30 tests in `recordcmd` failed on assertions against the stub. In the table below, "stub" names that red. Each probe ran through `bench probe` and returned `restored=yes`.

| Row | Test | Red route |
|---|---|---|
| RE11 | `TestRecordChunkWritesFullCommitIDs` | stub |
| RE12 | `TestRecordChunkSourceDigestExcludesTheRecord` | stub. Probe: swap the source digest for the tip tree ID. Verdict `bit`. |
| RE13 | `TestRecordChunkPlanDigestReadsThePlan` | stub |
| RE14 | `TestRecordChunkCopiesAcceptanceRows` | stub |
| RE15 | `TestRecordChunkCreatesAVersionTwoRecord` | stub |
| RE16 | `TestRecordChunkCreatesInAnEmptyDirectory` | stub. Probe: swap `os.MkdirAll` for `os.Mkdir`. Verdict `bit`. |
| RE17 | `TestRecordChunkRefusesAVersionOnePlan` | stub |
| RE18 | `TestRecordChunkAppendsAfterTheRecordedChunk` | stub |
| RE19 | `TestRecordChunkUpdateKeepsResults` | stub. Probe: replace the whole recorded entry. Verdict `bit`. |
| RE20 | `TestRecordChunkRepeatIsByteIdentical` | stub. The RE19 probe also failed this test. |
| RE21 | `TestRecordChunkRefusesAnUnplannedChunk` | stub |
| RE22 | `TestRecordChunkRefusesAnUnknownRevision` | stub |
| RE23 | `TestRecordChunkReportsCreated` | stub. Probe: rename the `source_digest` field. Verdict `bit`. |
| RE24 | `TestRecordChunkReportsAdded` | stub. Probe: swap the action `added` for `created`. Verdict `bit`. |
| RE25 | `TestRecordChunkReportsUpdated` | stub. The RE23 probe also failed this test. |
| RE26 | `TestRecordKeepsTheProseAroundTheFence` | stub |
| RE27 | `TestRecordChangesOnlyTheRecordPath` | stub |
| RE28 | `TestRecordRefusesADanglingRecordLink` | stub. Probe: read the record with `os.ReadFile`. Verdict `bit`. |
| RE29 | `TestRecordRefusesALiveRecordLink` | stub |
| RE30 | `TestRecordRefusesASpecialRecordFile` | stub |
| RE31 | `TestRecordRefusesAnEmptyRecordFile` | stub |
| RE32 | `TestRecordRefusesAnUnterminatedFence` | stub |
| RE33 | `TestRecordRefusesAnInvalidRecord` | stub |
| RE34 | `TestRecordRefusesThePrimaryCheckout` | stub |
| RE35 | `TestRecordPrimaryRefusalComesFirst` | stub |
| RE106 | `TestRecordFailedTemporaryWriteChangesNothing` | stub. Probe: write the record in place. Verdict `bit` on this test only. |
| RE37 | `TestRecordRefusesOutsideARepository` | stub |
| RE38 | `TestRecordGrammarRefusals` | stub |
| RE39 | `TestRecordHelpPrintsEachForm` | stub |
| RE105 | `TestRecordHelpSpellings` | stub |
| RE40 | `TestRecordRouteAnswersItsUsage` | Before the registry row existed, the dispatcher exited 2. |
| RE41 | `TestHelpInventoryIsComplete` | The golden row came before the registry row. |

The `os.ReadFile` probe did not fail RE29. The parser refused the bytes outside the tree, so that probe is not the writer that the RE29 row names.

A probe omitted the second parse of the rendered bytes in the write transaction. The verdict was `silent`. No ticket 2 row reaches a valid record that renders invalid. RE66 in ticket 3 is the row that grades this omission.

The fence of ticket 2 did not hold `internal/conformance/help_inventory_single_source_test.go`. Its projection allowlist refused the `record` help rows. The plan commit `4ad69dd9` added the path, and the allowlist now names `recordHelpRows`. The refusal of `TestRootConformance` before that edit is the red. After that edit, `TestRootConformance` refused a duration literal in the FIFO test. The test now waits `bounds.TestDeadline(0)`, and that refusal is its red.

The duplicated-facts sweep made one projection, `formHelpRows`, for the preflight and record help rows. The help spellings come from `usage.Parse`. The usage line, the grammar, and the help row of a form come from one declaration. Each independent test expectation has a recorded red: the TOON field list, the golden row, and the usage prefix. The comment sweep found no comment that holds a red record or a test result.

Focused checks at commit `f712f8ad`:

- `bench test --package ./internal/reviewrecord/...` passed in 2927 ms.
- `bench test --package ./cmd/bench` passed in 13084 ms.
- `bench test --package ./internal/conformance` passed in 35474 ms. It skipped 3 tests that the author did not write. These tests need a unix socket or a character device, which this filesystem cannot create.

At commit `bf3ddc94`, these checks passed: `TestRootConformance`, `skip-ownership`, `subcommand-routing`, `axi-query-registry`, and `package-core-guard`. `bench structure` reported no issue in the changed paths. `cmd/bench/main.go` has 444 lines, and `cmd/bench/command_registry_test.go` did not change.

## RE-C2 chunk review, round 1

Three fresh fable / high sessions reviewed the frozen pair `6d7f3969..bf3ddc94`. Each axis bound the review evidence `sha256:7eccb1a7` with `--check-current`. The raw finding count is 4. P3 and C2 name the same fix, so the repair-target count is 3. The repair allowance of RE-C2 is 2 cycles, and 0 cycles are used.

### Standards

Finding count: 1. Worst issue: S1.

- S1, auto-fix, confidence 5. The single-line flag check in `internal/reviewrecord/recordcmd/command.go` line 119 writes its own control-rune test. `sanitize.LineSafe` already owns that predicate. AGENTS.md makes duplicated knowledge a defect, and invariant 4 says to compose an existing seam.

### Spec

Finding count: 1. Worst issue: P3. Each of the 32 RE-C2 rows holds.

- P3, auto-fix, confidence 5. The transaction parses the rendered payload, but it does not grade the whole rendered document against the bound of the reader. A chunk form can then write a record that `Read` refuses. The spec says that the transaction parses the rendered bytes through the same reader path. Refusal step 7 refuses a rendered record that fails the reader.

### Coverage

Finding count: 2. Worst issue: C2.

- C1, auto-fix, confidence 6. No test grades refusal step 3 for `--chunk`. A probe that turned off the check for `--chunk` returned `silent`. The repair adds this test as the one hardening cycle of RE-C2.
- C2, auto-fix, confidence 6. This is the same defect as P3, from a trace of a record near 2 MiB.

### Repair targets

- T1 closes P3 and C2. Grade the rendered document through the bound of the reader before the rename, and add a red-capable test.
- T2 closes S1. Compose `sanitize.LineSafe` in the single-line flag check.
- T3 closes C1. Add a chunk-form test that sends a control byte in `--chunk` with an unknown `--tip` and expects a refusal that names `--chunk`.

### Advice

- `HelpRow` has a declaration in `recordcmd` and in `evidencecmd`, and `formHelpRows` bridges them.
- `command.go` builds the spec path by hand, and `reviewrecord.Slug` grades it.
- The usage-line shape differs from the preflight verb.
- `NewLinked` repeats the default fixture body.
- The temporary-file cleanup after `CreateTemp` has no feasible test input.
- An unknown form name echoes raw through the shared `toon.Usage` helper.
- An empty `reviews/` directory stays after a failed first write.

### Command contribution

The Spec axis suggests that the write-spec registry sweep bind the help-row projection allowlist to any new projected help row. `bench learning` holds that entry. The Coverage axis suggests that the author run one probe for each single-line flag of each form. The Standards axis found no contribution.
