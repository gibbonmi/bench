# Review evidence file pages review record

## RE1 author evidence

Ticket 1 had a fresh `bench-writer` author on opus at high effort. The author
started at `7bfc1182` and committed `260a62ea` on a lane pass. The author read
the metadata, the ticket, and the spec sources of the prepared evidence itself.

The guidance now states one order: the author record commit, then the review
charge, then the axis dispatch. A new chunk-chain anchor enforces the order.
Two new workflow fixtures swap the order and omit the new sentences.
`TestReviewRecordChargeOrder` binds a charge after a record commit, then
requires the stale-tip refusal after a later record commit.

### Probe verdicts

Each probe ran through `bench probe`, and each restore reads `yes`.

| File | Mutation | Check or test | Verdict |
|---|---|---|---|
| `internal/preflight/charge_pack.go` | swap: accept the stale source tip | TestReviewRecordChargeOrder | bit |
| `internal/anchors/registry_chunk_chain.go` | omission: the charge-order anchor | docs-currency-workflow | silent |
| `internal/anchors/registry_chunk_chain.go` | omission: the charge-order anchor | TestChunkChainAnchors | bit |
| `internal/anchors/registry_chunk_chain.go` | omission: the charge-order anchor | TestEveryRetainedFixtureBitesThroughRegisteredOwner | bit |
| `.agents/commands/bench-implement-spec.md` | swap: the charge before the record commit | docs-currency-workflow | bit |
| order-swap fixture `MUTATE.json` | swap: the fixture keeps the correct order | TestEveryRetainedFixtureBitesThroughRegisteredOwner | bit |

The orchestrator ran the last probe as the independent coordinator probe.

The silent row is not a defect. The `docs-currency-workflow` check grades
only the live tree, and the live tree obeys the order when the anchor is
absent. The spec names the workflow fixture owner as the RE11 proof. That
owner bit on both new fixtures when the anchor was absent.

### Verification

The author's two ticket checks passed at `260a62ea`. The JSON payload holds
each result. The build preflight is green at `260a62ea`.

## RE1 chunk review, round 1

The frozen pair is base `c8c444ffae2fb1578cfa54a22fa632590ffbc322` and tip
`260a62ea1e807c1989645b894c1a029f5a4b491e`. The record commit `f215f284`
follows the tip. The shared evidence is
`sha256:426844e754c9cc8b438b49e3aaccd54e5b114491f72ce30f11ba4e06b17026ae`.
Each axis ran in a fresh `bench-reviewer` session on fable at high effort, by
user direction. Only the Coverage axis ran probes, and it left the tree clean.

Raw findings: Standards 2, Spec 0, Coverage 1. Repair targets: 2 open for a
reviewer decision.

## Standards

Findings: 2. Worst: R1.

- R1 (ask-user, confidence 6): `.agents/commands/bench-implement-spec.md:21`
  joins two Land lines with no word change. The join offsets the new Build
  line against the 80-line prose budget. The axis cites invariant 4, the
  smallest diff. A revert reds the budget, so the fix is a budget decision or
  a real cut.
- R2 (no-op, confidence 2): the diagnostic and the needle occur in the
  registry, the test, and the fixtures. The record shows the omission red for
  the independent expectation, and the fixture copies obey the family
  convention.

Advice: the new Build paragraph carries the record order and the
plan-expansion rule. `commitRecord` in the new test restates the sequence of
`preflighttest.LegacyCommitted`. The new test restates the stale-tip text
without a declaration of independence.

## Spec

Findings: 0. RE11 and RE12 are covered. The delta holds no RE2 work.

Advice: the swap fixture keeps the sentence "Prepare the review charge from
that record commit.", so the mutated text contradicts itself. It still reds.

## Coverage

Findings: 1. Worst: R3.

- R3 (ask-user, confidence 6): `internal/anchors/registry_chunk_chain.go:8`
  is a substring check. A second sentence that states the reverse order in the
  same Build section passes `docs-currency-workflow`. The probe that appends
  "Prepare the review charge before the author record commit." came back
  silent, and its restore reads `yes`. The spec row names only omission and
  order-swap mutations, and every `RequireInSection` anchor shares this
  property.

The axis closed the RE12 refusal cause, the binding-success half, and the
prepared tip. Its independent bypass probe on `charge_pack.go` bit.

Advice: a needle inside a fenced block in the Build section also satisfies
the locator.

## RE1 repair cycle 1

The reviewer decided both open findings. R3 takes the chunk's one hardening
cycle. R1 reverts the reflow and raises the prose budget row. Repair cycle 1
of 2 is consumed.

A fresh `bench-writer` repair session on opus at low effort started at
`639e3a59` and committed `304913d7`. A `ForbidInSection` anchor now reds a
reverse-order statement in the Build section. A new reversed fixture proves it
through the workflow owner. The Land lines keep the base layout, and the
budget row reads 81.

Two plan expansions preceded and followed the repair. The first added
`projects/benchkit.md` to the `Writes:` line of ticket 1. The second added the
seven fixtures that pin that file, after `fixture-closure` went red. The
chunk tip is now `c024a237`. The repair session ran the two ticket checks
again at that tip, and each check passed.

| File | Mutation | Check or test | Verdict |
|---|---|---|---|
| `.agents/commands/bench-implement-spec.md` | swap: add a reverse-order sentence, before the repair | docs-currency-workflow | silent |
| `.agents/commands/bench-implement-spec.md` | swap: add a reverse-order sentence, after the repair | docs-currency-workflow | bit |
| `internal/anchors/registry_chunk_chain.go` | omission: the forbid anchor | TestEveryRetainedFixtureBitesThroughRegisteredOwner | bit |
| `projects/benchkit.md` | swap: the budget row back to 80 | guidance-prose-budgets | bit |
| reversed fixture `MUTATE.json` | swap: "before" to "after" | TestEveryRetainedFixtureBitesThroughRegisteredOwner | bit |

The orchestrator ran the last probe as the independent coordinator probe.
Each restore reads `yes`.

The build preflight reports `base-current` red, because `main` gained a
landing after the build started. The review chain forbids a `main` merge
after the first chunk, and the landing composes that commit.

```bench-review-record
{
  "version": 2,
  "spec": "specs/review-evidence-file-pages/spec.md",
  "plan_digest": "sha256:d7852c577c37e073b66469c12432eaeced47dabe23f549758c73f9554f99f0b3",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RE1",
      "base": "c8c444ffae2fb1578cfa54a22fa632590ffbc322",
      "tip": "c024a23711443fb649bbd6063a70f70459e37511",
      "plan_digest": "sha256:d7852c577c37e073b66469c12432eaeced47dabe23f549758c73f9554f99f0b3",
      "source_digest": "65fd30ec52ec589771bba326bc17068488c199a1",
      "acceptance_rows": [
        "RE11",
        "RE12"
      ],
      "verification": [
        {
          "id": "re1-1-workflow-r1",
          "performer": "claude:bench-writer/re-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f695ba5f19389ccf2b45f93462410180a429370",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-author-20260926/1-workflow@260a62ea",
            "digest": "sha256:90222f23de517c6dc4717d6778e330befc75a5204f58e78b91032ae5ae90be46",
            "excerpt": "internal/conformance,pass,556; failures[0]; skips[0]"
          },
          "requirement": "1-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "re1-1-record-order-r1",
          "performer": "claude:bench-writer/re-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "8f695ba5f19389ccf2b45f93462410180a429370",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-author-20260926/1-record-order@260a62ea",
            "digest": "sha256:1e368a6b9cb67d3508da92d8dbc44e4996f4fefd9419399ef2798cb13eaff0a3",
            "excerpt": "ok github.com/gibbonmi/bench/internal/preflight/evidencecmd 7.776s"
          },
          "requirement": "1-record-order",
          "command": "go test -count=1 -parallel=2 ./internal/preflight/evidencecmd",
          "exit_code": 0
        },
        {
          "id": "re1-1-workflow-r2",
          "performer": "claude:bench-writer/re-t1-repair-c1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "65fd30ec52ec589771bba326bc17068488c199a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-repair-c1-20260926/1-workflow@c024a237",
            "digest": "sha256:374ff3cb73050cfa13dc3a30b7cb61d8fc586d391e9585ad3bf66a03db3354d2",
            "excerpt": "packages[1]{package,status,elapsed_ms}:\n  github.com/gibbonmi/bench/internal/conformance,pass,619\nfailures[0]{package,test,line}:\nskips[0]{package,test,reason}:"
          },
          "requirement": "1-workflow",
          "command": "bench test --check docs-currency-workflow",
          "exit_code": 0
        },
        {
          "id": "re1-1-record-order-r2",
          "performer": "claude:bench-writer/re-t1-repair-c1",
          "role": "author-verification",
          "model": "opus",
          "effort": "low",
          "source_digest": "65fd30ec52ec589771bba326bc17068488c199a1",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re-t1-repair-c1-20260926/1-record-order@c024a237",
            "digest": "sha256:fc3f21bce946d69e94977ac14b7e01a324a007132768827d8f496103db93bb2f",
            "excerpt": "ok  \tgithub.com/gibbonmi/bench/internal/preflight/evidencecmd\t8.174s"
          },
          "requirement": "1-record-order",
          "command": "go test -count=1 -parallel=2 ./internal/preflight/evidencecmd",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "re1-r1-standards",
          "performer": "claude:bench-reviewer/re1-standards",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "8f695ba5f19389ccf2b45f93462410180a429370",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/re1-standards@260a62ea",
            "digest": "sha256:725548af93241f2d871b66aa614c954f563370ded359156a546b46facc08b28e",
            "excerpt": "Standards: 2 findings. Worst: the Land paragraph join is an unrelated cosmetic edit driven by a zero-headroom prose budget, not by RE11/RE12."
          },
          "axis": "Standards",
          "base": "c8c444ffae2fb1578cfa54a22fa632590ffbc322",
          "tip": "260a62ea1e807c1989645b894c1a029f5a4b491e",
          "finding_ids": [
            "R1",
            "R2"
          ],
          "supersedes": []
        },
        {
          "id": "re1-r1-spec",
          "performer": "claude:bench-reviewer/re1-spec",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "8f695ba5f19389ccf2b45f93462410180a429370",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/re1-spec@260a62ea",
            "digest": "sha256:50f558b1cd6fa8ce3e3478b8de1161b6bdcd0ddd81dd5ba600dc96d01a25b8d0",
            "excerpt": "Spec: 0 findings. Worst: none; both mapped rows are delivered as the spec states them, and no RE2 work leaked."
          },
          "axis": "Spec",
          "base": "c8c444ffae2fb1578cfa54a22fa632590ffbc322",
          "tip": "260a62ea1e807c1989645b894c1a029f5a4b491e",
          "finding_ids": [],
          "supersedes": []
        },
        {
          "id": "re1-r1-coverage",
          "performer": "claude:bench-reviewer/re1-coverage",
          "role": "independent-review",
          "model": "fable",
          "effort": "high",
          "source_digest": "8f695ba5f19389ccf2b45f93462410180a429370",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/re1-coverage@260a62ea",
            "digest": "sha256:ea09b045a8e1778070bf8ae1fdb846839693f667c40510e74b0301c879e077f9",
            "excerpt": "Coverage: 1 finding. Worst: RE11's anchor is a substring presence check, so a contradicting order sentence in the same Build section passes `docs-currency-workflow` silently."
          },
          "axis": "Coverage",
          "base": "c8c444ffae2fb1578cfa54a22fa632590ffbc322",
          "tip": "260a62ea1e807c1989645b894c1a029f5a4b491e",
          "finding_ids": [
            "R3"
          ],
          "supersedes": []
        }
      ]
    }
  ],
  "completion": {
    "state": "pending"
  }
}
```
