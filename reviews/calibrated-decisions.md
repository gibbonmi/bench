# Calibrated decisions

## CD1 ticket 1 author evidence

The ticket adds one `Claim schema` section to the delegation discipline reference. It adds one pointer sentence to the delegate skill. Seven anchors and seven omission canary fixtures grade the seven sentences.

### Scenario rows

These two rows show the claim schema shape the ticket describes. The table is illustrative: its label cell shows the scenario the ticket states, not a coordinator probe. These rows do not count as pairs.

| surface | claim | status | confidence | label |
| --- | --- | --- | --- | --- |
| delegate return | the named check runs green on the returned tree | claimed | 7 | held |
| delegate return | the guidance reads clearly for a cold reader | abstained |  |  |

The coordinator probes the first row's named check on the exact tree and labels it `held`. The coordinator labels nothing on the second row and counts one abstention.

### Done-claim table

The author did not write the label cell. The coordinator wrote each label. The label source is the gate lane at the merge, green at tip 09f26779. One independent swap probe on the CR13 sentence also bit.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR1 | verified | 9 | held |
| CR2 | verified | 9 | held |
| CR3 | verified | 9 | held |
| CR13 | verified | 9 | held |
| CR14 | verified | 9 | held |
| CR32 | verified | 9 | held |
| CR34 | verified | 9 | held |
| CR37 | verified | 9 | held |

### Red-then-green log

Each row below except CR14 uses `bench probe` with the omission kind against `docs-currency-workflow`. The row for CR14 uses a line-growth swap against `guidance-prose-budgets`. The verb records the green baseline, then the red under the mutation, then the proved restore.

| row | red line under the mutation | green |
| --- | --- | --- |
| CR1 | `gate: calibration: a done-claim row needs its status and stated confidence` | baseline passed, restored yes |
| CR2 | `gate: calibration: an unstatable confidence needs an abstained row` | baseline passed, restored yes |
| CR3 | `gate: calibration: the coordinator's tree probe labels each done-claim row` | baseline passed, restored yes |
| CR13 | `gate: calibration: the charge must name the claim schema section` | baseline passed, restored yes |
| CR14 | `gate: prose-budget exceeded: .agents/skills/bench-craft-delegate/SKILL.md is 128 lines, over its 126-line budget` | baseline passed, restored yes |
| CR32 | `gate: calibration: a claim must carry no free-text field` | baseline passed, restored yes |
| CR34 | `gate: calibration: verified and claimed need their evidence definitions` | baseline passed, restored yes |
| CR37 | `gate: calibration: an abstention must never become a refuted claim` | baseline passed, restored yes |

`TestEveryRetainedFixtureBitesThroughRegisteredOwner` proves each of the seven new canaries bites through its registered owner. The same run proves the 525 earlier fixtures keep their planted diagnostics.

### Probe verdict

The self-probe omits the CR3 sentence from the delegation discipline reference and runs `docs-currency-workflow`. The verdict line reads `bit`, and `restored` reads `yes`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | pass | 890 ms |
| `bench test --check guidance-prose-budgets` | pass | 4 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 15 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 9299 ms |
| `go vet ./...` | pass | no output |
| `bench gate-prose` on both edited files | pass | two pass rows |

## CD2 ticket 2 author evidence

The ticket adds one `What a confidence states` section to the finding discipline reference. It adds one second sentence to the review skill's pointer line and one sentence to the review-implementation pickup step. Six anchors and six omission canary fixtures grade the six sentences.

### Scenario finding

This row shows the finding shape the ticket describes. The label cell holds the scenario the ticket states, not a reviewer disposition on real work. The row does not count as a pair.

| surface | claim | citation | status | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- | --- |
| review finding | the behavior under the cited line is wrong | the diff line the axis read this pass | claimed | 3 | no-op | refuted |

The finding blocks the round, because its kind and its citation decide. The confidence of 3 changes nothing about that block. The reviewer disposes the finding `no-op`, and the fixed mapping labels it `refuted`.

### Done-claim table

The author did not write a label cell. The coordinator wrote each label. The label source is the gate lane at the merge, green at tip 1ecb33ec. One independent swap probe on the CR6 mapping also bit.

| row | status | confidence | label |
| --- | --- | --- | --- |
| CR4 | verified | 9 | held |
| CR5 | verified | 9 | held |
| CR6 | verified | 9 | held |
| CR15 | verified | 9 | held |
| CR16 | verified | 9 | held |
| CR22 | verified | 9 | held |
| CR23 | verified | 8 | held |
| CR30 | verified | 9 | held |

### Red-then-green log

Each row below except CR16 and CR23 uses `bench probe` with the omission kind against `docs-currency-workflow`. The row for CR16 uses a line-growth swap against `guidance-prose-budgets`. The row for CR23 uses a diff of the package against the base commit. The verb records the green baseline, then the red under the mutation, then the proved restore.

| row | red line under the mutation | green |
| --- | --- | --- |
| CR4 | `gate: calibration: a finding needs its stated confidence` | baseline passed, restored yes |
| CR5 | `gate: calibration: the confidence must never change whether a finding blocks` | baseline passed, restored yes |
| CR6 | `gate: calibration: the dispositions need their label mapping` | baseline passed, restored yes |
| CR15 | `gate: calibration: the review skill must point at the finding confidence` | baseline passed, restored yes |
| CR16 | `gate: prose-budget exceeded: .agents/skills/bench-craft-review/SKILL.md is 123 lines, over its 122-line budget` | baseline passed, restored yes |
| CR22 | `gate: calibration: the pickup line must carry its stated confidence` | baseline passed, restored yes |
| CR23 | no mutation: `git diff --stat <base> -- internal/reviewrecord` gives empty output | the package is byte-identical to the base |
| CR30 | `gate: calibration: optional advice must carry no confidence` | baseline passed, restored yes |

`TestEveryRetainedFixtureBitesThroughRegisteredOwner` proves each of the six new canaries bites through its registered owner. A control run with one wrong `EXPECT` byte made that fixture fail, which shows the run reads the new fixtures. The same run proves every earlier fixture keeps its planted diagnostic.

### Probe verdict

The self-probe omits the CR15 sentence from the review skill and runs `docs-currency-workflow`. The verdict line reads `bit`, and `restored` reads `yes`.

### Verification table

Each run reports no skip. The elapsed time is the package time the verb reports.

| check | verdict | elapsed |
| --- | --- | --- |
| `bench test --check docs-currency-workflow` | pass | 1185 ms |
| `bench test --check guidance-prose-budgets` | pass | 5 ms |
| `bench test --package ./internal/anchors/... --run 'TestCalibration'` | pass | 29 ms |
| `bench test --package ./internal/conformance --run TestEveryRetainedFixtureBitesThroughRegisteredOwner` | pass | 11923 ms |
| `go vet ./...` | pass | no output |
| `git diff --stat <base> -- internal/reviewrecord` | empty | no output |
| `bench gate-prose` on the three edited files and this pickup | pass | four pass rows |

## CD1 review

The frozen pair is base `9148850200714f000eec2fbf44cddea6182f95c7` and tip `09f26779f65b7938f313cff9ec877fabe9d009f5`. The reviewer directed the review line. The first pass of every chunk review runs fable / medium. Every later pass on the same chunk runs sonnet / xhigh. Each axis ran in its own read-only worktree.

The raw finding count is 2, both on the Standards axis. The de-duplicated repair target count is 1. CD1 consumed 0 repair cycles. The one accepted repair is an evidence-only prose correction inside this pickup.

### Standards

Finding count: 2 judgment calls, 0 hard violations. Worst issue: the registry test restates the seven needles with no recorded independence red.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD1-S1 | `internal/anchors/registry_calibration_test.go` restates the seven needles and diagnostics from the registry file, and the pickup records omission reds through `bench probe`, not a registry-row deletion the test alone catches. | AGENTS.md, one source per fact, the test-expectation exception. | 5 | no-op | refuted |
| CD1-S2 | The red-then-green paragraph in this pickup stated a universal rule, then retracted it for CR14 in the next sentence. | ASD-STE100 prose rules; a universal statement that the next sentence retracts. | 6 | auto-fix | held |

CD1-S1 is `no-op`. The exemplar pair `registry_ticket_passes.go` and its test hold the same shape. The Coverage axis showed the red route through `anchorHarness.check` and the fixture-bite run. The coordinator repaired CD1-S2 in this pickup as an evidence-only correction. The paragraph now names the CR14 exception in its first sentence.

A sonnet / xhigh reaffirmation pass in its own worktree confirmed the correction. It found CD1-S2 gone, CD1-S1 still `no-op`, no other claim changed, and no new finding.

The axis gave two items of optional advice with no confidence. Replace the derived count "the 525 earlier fixtures" with "every earlier fixture". Prefer "for every status" over "whatever its status" in the reference. Both stay open. The reference bytes are the spec's pasted needle, and the count is run evidence.

### Spec

The finding count is 0, and there is no worst issue. The axis audited every CD1 row and found each held.

The exact needle bytes sit in the named section. Each anchor is `require-in-section`. Each fixture exists under its fence name. The delegate skill holds 125 lines. The package `internal/reviewrecord` is unchanged, and every pre-existing anchored sentence keeps its bytes.

The axis gave one item of optional advice with no confidence. The scenario table's label is illustrative. The coordinator added that note above the table.

### Coverage

The finding count is 0. No worst issue holds against the tree. The axis ran eight independent probes. Five bit and three stayed silent.

The probes were a sentence moved above its heading, a duplicate heading, a fenced copy, and a suffix weakening. They were also a case-only rewrite, a renamed skill section, a duplicated sentence, and a contradicting sibling. The three silent probes are matcher-wide semantics of `internal/anchors/locate.go`. The spec routes those semantics to review, so they are not defects of this delta.

The axis gave one item of optional advice with no confidence. The CR1 why-it-catches clause could say "omitted or byte-changed under fold" instead of "reworded". A fenced copy and a case-only rewrite stay green under the section-scoped matcher.

### Pairs recorded for CD5

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| delegate return | CR1 | verified | 9 | held | opus / high / author |
| delegate return | CR2 | verified | 9 | held | opus / high / author |
| delegate return | CR3 | verified | 9 | held | opus / high / author |
| delegate return | CR13 | verified | 9 | held | opus / high / author |
| delegate return | CR14 | verified | 9 | held | opus / high / author |
| delegate return | CR32 | verified | 9 | held | opus / high / author |
| delegate return | CR34 | verified | 9 | held | opus / high / author |
| delegate return | CR37 | verified | 9 | held | opus / high / author |
| review finding | CD1-S1 | claimed | 5 | refuted | fable / medium / Standards |
| review finding | CD1-S2 | claimed | 6 | held | fable / medium / Standards |

## CD2 review

The frozen pair is base `87869759bf975c1f76a6e12ccc638e8471d40dfd` and tip `1ecb33ec4cede2be8bde4bc07bb648786e906d13`. Main moved during the review, so the integration source merged main after that tip. The chunk keeps the ticket-only pair, because the checkpoint compares the recorded tip's own tree. Each axis ran fable / medium in its own read-only worktree.

The raw finding count is 4. The de-duplicated repair target count is 1. CD2 consumed 0 repair cycles. The one accepted finding is a new seam that the reviewer routed to a new chunk, CD2b, not a repair of this delta.

### Standards

Finding count: 2 judgment calls, 0 hard violations. Worst issue: two guidance sentences carry the 0 to 10 range.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD2-S1 | The finding discipline reference and the review skill each state the 0 to 10 range. | AGENTS.md, one source per fact; the spec requires both sentences verbatim as CR4 and CR15. | 8 | no-op | refuted |
| CD2-S2 | The registry test adds a second function instead of extending the first. | Ticket 2, "its harness test"; the exemplar keeps one function per topic. | 4 | no-op | refuted |

Both findings are `no-op`. The spec's pasted needles decide CD2-S1, and the ticket-passes exemplar decides CD2-S2. The axis noted that the scenario finding table omits the model, effort, and role cell; that table is illustrative and not a pair.

### Spec

The finding count is 0, and there is no worst issue. The axis audited every CD2 row and found each held. It accepted the author's heading judgment: `## What a confidence states` follows the reference's one-question-per-heading shape.

The axis gave two items of optional advice with no confidence. The CR22 anchor is section-scoped, so it cannot pin step 6. The disposition-label fixture removes two lines where one sentence would match its siblings.

### Coverage

Finding count: 2. Worst issue: the CR22 sentence can leave the pickup step while the check stays green.

| id | finding | citation | confidence | disposition | label |
| --- | --- | --- | --- | --- | --- |
| CD2-C1 | Omitting the step 6 heading moves the CR22 sentence into step 5, and `docs-currency-workflow` stays green, because `RequireInSection` is heading-scoped. | Spec row CR22, "the pickup step"; probe 8, `silent`. | 7 | ask-user | held |
| CD2-C2 | A contradicting sentence appended after the CR5 needle stays green. | The known bound of every sentence anchor; no row promises it. | 3 | no-op | refuted |

The reviewer decided CD2-C1. A step-scoped anchor kind is a new seam, and the reviewer expanded this run with one ticket for it. The ticket forms chunk CD2b, re-pins CR22 on step 6, and takes its own three-axis review. CD2 closes on the section anchor, and the finding is resolved by that decision. The axis ran eight probes: six bit and two stayed silent.

### Pairs recorded for CD5

| surface | claim | status | confidence | label | model / effort / role |
| --- | --- | --- | --- | --- | --- |
| delegate return | CR4 | verified | 9 | held | opus / high / author |
| delegate return | CR5 | verified | 9 | held | opus / high / author |
| delegate return | CR6 | verified | 9 | held | opus / high / author |
| delegate return | CR15 | verified | 9 | held | opus / high / author |
| delegate return | CR16 | verified | 9 | held | opus / high / author |
| delegate return | CR22 | verified | 9 | held | opus / high / author |
| delegate return | CR23 | verified | 8 | held | opus / high / author |
| delegate return | CR30 | verified | 9 | held | opus / high / author |
| review finding | CD2-S1 | claimed | 8 | refuted | fable / medium / Standards |
| review finding | CD2-S2 | claimed | 4 | refuted | fable / medium / Standards |
| review finding | CD2-C1 | claimed | 7 | held | fable / medium / Coverage |
| review finding | CD2-C2 | claimed | 3 | refuted | fable / medium / Coverage |

## Native review record

The fenced payload below retains every terminal return for the checkpoint. The first Standards occurrence lists its two findings; a later occurrence supersedes it after the evidence-only correction is reaffirmed.

```bench-review-record
{
  "version": 2,
  "spec": "specs/calibrated-decisions/spec.md",
  "plan_digest": "sha256:9a5079a914a8c51fea3eccf16a76b5dbb9a687c7eebec78a9bb84c36bfc6ff94",
  "implementation_session": "",
  "chunks": [
    {
      "id": "CD1",
      "base": "9148850200714f000eec2fbf44cddea6182f95c7",
      "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
      "plan_digest": "sha256:9a5079a914a8c51fea3eccf16a76b5dbb9a687c7eebec78a9bb84c36bfc6ff94",
      "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
      "acceptance_rows": ["CR1", "CR2", "CR3", "CR13", "CR14", "CR32", "CR34", "CR37"],
      "verification": [
        {
          "id": "cd1-workflow",
          "performer": "claude:bench-writer/cd-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t1-author-20260917/cd1-workflow@09f26779",
            "digest": "sha256:22335da2557fa9f4e7ee2983872b3ccb666dda93bbc3cd85f52d14fc03f34491",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,817\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "cd1-budgets",
          "performer": "claude:bench-writer/cd-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t1-author-20260917/cd1-budgets@09f26779",
            "digest": "sha256:6f72c99a424331a3381e284e55ffa703b490f81e19b498e848e8bfc5d9ca25ed",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,5\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "budgets",
          "command": "bench test --check guidance-prose-budgets",
          "exit_code": 0
        },
        {
          "id": "cd1-prose",
          "performer": "claude:bench-writer/cd-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-t1-author-20260917/cd1-prose@09f26779",
            "digest": "sha256:4094ea8549b8b3a7e5f00286119ddfa586e3ac1a55cd0a75a8d29fbe59e6bacd",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,156\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "prose",
          "command": "bench test --check prose-mechanics",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "cd1-standards",
          "performer": "claude:bench-reviewer/cd-c1-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/cd-c1-standards-20260917@09f26779",
            "digest": "sha256:80c16504f159fca1ef6d0ca4e647ff755032e64b04b5dd96ddae18608c38f315",
            "excerpt": "Count: 0 hard violations, 2 judgment calls. CD1-S1 duplicated expectation without a recorded independence red, confidence 5, no-op. CD1-S2 STE contradiction inside one pickup paragraph, confidence 6, auto-fix."
          },
          "axis": "Standards",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": ["CD1-S1", "CD1-S2"],
          "supersedes": []
        },
        {
          "id": "cd1-standards-r2",
          "performer": "claude:bench-reviewer/cd-c1-standards-r2",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "xhigh",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c1-standards-r2-20260917@09f26779",
            "digest": "sha256:3bb29329cb87762c6aa5d13b5c81412238dc14bc68c213f8960525cf1e869255",
            "excerpt": "Verdict: pass. CD1-S2 repaired in the pickup; CD1-S1 no-op holds against the tree; no new findings."
          },
          "axis": "Standards",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": [],
          "supersedes": ["cd1-standards"]
        },
        {
          "id": "cd1-spec",
          "performer": "claude:bench-reviewer/cd-c1-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c1-spec-20260917@09f26779",
            "digest": "sha256:1287bdc2526232d4fe9411913e22d0b3d1b65478f1220b82e6f5a2473b238175",
            "excerpt": "Finding count: 0 blocking, 0 auto-fix. Worst issue: none. Every CD1 row held: CR1, CR2, CR3, CR13, CR14, CR32, CR34, CR37."
          },
          "axis": "Spec",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "cd1-coverage",
          "performer": "claude:bench-reviewer/cd-c1-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "medium",
          "source_digest": "e4866f1d8569f350474dcf28334a23f195baedd8",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/cd-c1-coverage-20260917@09f26779",
            "digest": "sha256:86daef777047f41fb6b08515c5cc99f8172e72e245a7d80afb68a2fee2a9a4bb",
            "excerpt": "Finding count: 0 blocking. Worst issue: none that holds against the tree. Eight independent probes: five bit, three silent on matcher-wide semantics routed to review."
          },
          "axis": "Coverage",
          "base": "9148850200714f000eec2fbf44cddea6182f95c7",
          "tip": "09f26779f65b7938f313cff9ec877fabe9d009f5",
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
