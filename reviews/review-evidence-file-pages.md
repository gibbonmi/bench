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

```bench-review-record
{
  "version": 2,
  "spec": "specs/review-evidence-file-pages/spec.md",
  "plan_digest": "sha256:137fb3616ac828784f2984e8c6a918d4f069bf9e993fb8a854cb8d9d4fd41c3c",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RE1",
      "base": "c8c444ffae2fb1578cfa54a22fa632590ffbc322",
      "tip": "260a62ea1e807c1989645b894c1a029f5a4b491e",
      "plan_digest": "sha256:137fb3616ac828784f2984e8c6a918d4f069bf9e993fb8a854cb8d9d4fd41c3c",
      "source_digest": "8f695ba5f19389ccf2b45f93462410180a429370",
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
        }
      ],
      "reviews": []
    }
  ],
  "completion": {
    "state": "pending"
  }
}
```
