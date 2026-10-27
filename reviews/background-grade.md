# Review outcomes

```bench-review-record
{
  "version": 2,
  "spec": "specs/background-grade/spec.md",
  "plan_digest": "sha256:d3cd8ebcd312d6538a9422dac4b8cf86c3760a8c9d9b0e67f841d3ca08e2c596",
  "implementation_session": "",
  "chunks": [
    {
      "id": "BG-C0",
      "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
      "tip": "2bb0564d3585caf53620bd8d754269fcbcbd032b",
      "plan_digest": "sha256:143f753bd2b8510ff891b8076aa2e2f1bc7d8ddafaf78f72760bf7da7abc0076",
      "source_digest": "d1bc0f6d850ed1a654f9ef579c4b01d937668eb2",
      "acceptance_rows": [
        "BG93",
        "BG94",
        "BG95"
      ],
      "verification": [
        {
          "id": "bg-c0-t01-conformance",
          "performer": "claude:bg-t01-author-20261027",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "3d4db9dbba5bcfcfee48a3710cd1382dad2b240b",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bg-t01-author-20261027",
            "digest": "sha256:301696dc3cd894eb5e64634e001f59b6de0e7f205ada8cc506c3f22b616cd2a0",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,14733,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "7bdc95c72c6ef570407a4105cfd8d3597b69e989",
          "requirement": "conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        },
        {
          "id": "bg-c0-t01-r1-conformance",
          "performer": "claude:bg-t01-repair1-20261027",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "d1bc0f6d850ed1a654f9ef579c4b01d937668eb2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bg-t01-repair1-20261027",
            "digest": "sha256:e5ce57c0166dc636a56848eefa1bf03cd0ad399c17ad71fd0675eeb88c459385",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,15241,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2bb0564d3585caf53620bd8d754269fcbcbd032b",
          "requirement": "conformance",
          "command": "bench test --check root-conformance",
          "exit_code": 0
        }
      ],
      "reviews": [
        {
          "id": "bg-c0-review-standards",
          "performer": "claude:bg-c0-standards-20261027",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "3d4db9dbba5bcfcfee48a3710cd1382dad2b240b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bg-c0-standards-20261027",
            "digest": "sha256:9f5cb1503ff2492bc7c130defafeed80ff400fdb6dbd447a8b4238900ae7aeb2",
            "excerpt": "Standards axis BG-C0 (84e9b1d37..7bdc95c72): 5 findings, no blocking defect. Worst: census.md:83-85 Affinity column with no stated read method (ask-user, conf 5). Others: inline list of four packages at :127 (auto-fix), waiver results without numbers at :142-145 (ask-user), present tense for a past run at :64 (auto-fix), session-local absolute paths at :45-46 (no-op).\n"
          },
          "axis": "Standards",
          "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
          "tip": "7bdc95c72c6ef570407a4105cfd8d3597b69e989",
          "finding_ids": [
            "BG-C0-S1",
            "BG-C0-S2",
            "BG-C0-S3",
            "BG-C0-S4"
          ],
          "supersedes": []
        },
        {
          "id": "bg-c0-review-spec",
          "performer": "claude:bg-c0-spec-20261027",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "3d4db9dbba5bcfcfee48a3710cd1382dad2b240b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bg-c0-spec-20261027",
            "digest": "sha256:ac0d6619fe7a643db910259c9492c9b8dddff87edaeffd45cf04399d3898bb90",
            "excerpt": "Spec axis BG-C0 (84e9b1d37..7bdc95c72): BG93 met, BG94 met, BG95 met; medians 237.73 and 298.98, ratio 1.2577, bound 285.28, Verdict: stop correct. 2 findings. Worst: census.md:148-151 the waiver's sealed-rerun claim has no recorded data (ask-user, conf 5). Other: unsealed proxy limit disclosed at :130-139 (no-op).\n"
          },
          "axis": "Spec",
          "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
          "tip": "7bdc95c72c6ef570407a4105cfd8d3597b69e989",
          "finding_ids": [
            "BG-C0-P1"
          ],
          "supersedes": []
        },
        {
          "id": "bg-c0-review-coverage",
          "performer": "claude:bg-c0-coverage-20261027",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "3d4db9dbba5bcfcfee48a3710cd1382dad2b240b",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bg-c0-coverage-20261027",
            "digest": "sha256:744212fd3cf3c9f2bf2bcf848b87613ae80a97937f1ed2ef56df642e3d49318a",
            "excerpt": "Coverage axis BG-C0 (84e9b1d37..7bdc95c72): arithmetic correct, 0 red foreground runs. 6 findings. Worst: census.md:136 and :148-151 waiver and contention claims rest on unrecorded runs (ask-user, conf 7). Others: Affinity column source not named :89-91 (auto-fix); race and system phase priority unsampled :86-87 (ask-user); restart time unrecorded :80-81; sample counts do not reconcile :91-95 (auto-fix); shared worktree and cache not stated :51 (ask-user).\n"
          },
          "axis": "Coverage",
          "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
          "tip": "7bdc95c72c6ef570407a4105cfd8d3597b69e989",
          "finding_ids": [
            "BG-C0-C1",
            "BG-C0-C2",
            "BG-C0-C4",
            "BG-C0-C5",
            "BG-C0-C6"
          ],
          "supersedes": []
        },
        {
          "id": "bg-c0-confirm1-standards",
          "performer": "claude:bg-c0-confirm1-20261027",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d1bc0f6d850ed1a654f9ef579c4b01d937668eb2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bg-c0-confirm1-20261027",
            "digest": "sha256:af2e58b6b88fd3ca09773fd02f711ff17e8dbc38b6849afdebccd5230e289ff9",
            "excerpt": "Unified confirming round BG-C0, repair delta 7bdc95c72..2bb0564d3, Standards axis: T1-T5 closed; last line Verdict: stop at census.md:207; new findings 0.\n"
          },
          "axis": "Standards",
          "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
          "tip": "2bb0564d3585caf53620bd8d754269fcbcbd032b",
          "finding_ids": [],
          "supersedes": [
            "bg-c0-review-standards"
          ]
        },
        {
          "id": "bg-c0-confirm1-spec",
          "performer": "claude:bg-c0-confirm1-20261027",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d1bc0f6d850ed1a654f9ef579c4b01d937668eb2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bg-c0-confirm1-20261027",
            "digest": "sha256:2bc56b9a29895c067616f10334f6040aa64d003cb1324be80ba1c952d8d87b7b",
            "excerpt": "Unified confirming round BG-C0, repair delta 7bdc95c72..2bb0564d3, Spec axis: T1-T5 closed; last line Verdict: stop at census.md:207; new findings 0.\n"
          },
          "axis": "Spec",
          "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
          "tip": "2bb0564d3585caf53620bd8d754269fcbcbd032b",
          "finding_ids": [],
          "supersedes": [
            "bg-c0-review-spec"
          ]
        },
        {
          "id": "bg-c0-confirm1-coverage",
          "performer": "claude:bg-c0-confirm1-20261027",
          "role": "independent-review",
          "model": "sonnet",
          "effort": "high",
          "source_digest": "d1bc0f6d850ed1a654f9ef579c4b01d937668eb2",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-reviewer-high:bg-c0-confirm1-20261027",
            "digest": "sha256:ba2acc777a09957ef45f026b9dd3e3c41ec5bba08ed040626e4a20b9e81158d6",
            "excerpt": "Unified confirming round BG-C0, repair delta 7bdc95c72..2bb0564d3, Coverage axis: T1-T5 closed; last line Verdict: stop at census.md:207; new findings 0.\n"
          },
          "axis": "Coverage",
          "base": "84e9b1d372f0931e75db1072013e63cd4f94ed00",
          "tip": "2bb0564d3585caf53620bd8d754269fcbcbd032b",
          "finding_ids": [],
          "supersedes": [
            "bg-c0-review-coverage"
          ]
        }
      ]
    },
    {
      "id": "BG-C1",
      "base": "2bb0564d3585caf53620bd8d754269fcbcbd032b",
      "tip": "2f788c7e089aeb26a394761866c52b08caaba7a7",
      "plan_digest": "sha256:d3cd8ebcd312d6538a9422dac4b8cf86c3760a8c9d9b0e67f841d3ca08e2c596",
      "source_digest": "659c103cfb6dc07b4bb03234f3115fc0d6b2c29e",
      "acceptance_rows": [
        "BG12",
        "BG13",
        "BG14",
        "BG15",
        "BG109",
        "BG16",
        "BG17",
        "BG18"
      ],
      "verification": [
        {
          "id": "bg-c1-t02-gate",
          "performer": "claude:bg-t02-author-20261027",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "659c103cfb6dc07b4bb03234f3115fc0d6b2c29e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bg-t02-author-20261027",
            "digest": "sha256:730b17822bc10e06e7ce8fa357ffdd0a9c1ee913cb675ce4ba96c58de9ad466a",
            "excerpt": "packages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/gate,pass,143009,734\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2f788c7e089aeb26a394761866c52b08caaba7a7",
          "requirement": "gate",
          "command": "bench test --package ./internal/gate",
          "exit_code": 0
        },
        {
          "id": "bg-c1-t02-conformance",
          "performer": "claude:bg-t02-author-20261027",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "659c103cfb6dc07b4bb03234f3115fc0d6b2c29e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude-agent:bench-writer-high:bg-t02-author-20261027",
            "digest": "sha256:bb6d7f040a4cf30f5446eaaaf8bce7e76988870dd6a5c415e43f707f4bd5e36a",
            "excerpt": "check[1]{name,kind,tests_run,subjects}:\n  root-conformance,conformance,1,0\npackages[1]{package,status,elapsed_ms,tests_run}:\n  github.com/gibbonmi/bench/internal/conformance,pass,16945,1\nfailures[0]{package,test,line,lines}:\nskips[0]{package,test,reason}:\n"
          },
          "source_commit": "2f788c7e089aeb26a394761866c52b08caaba7a7",
          "requirement": "conformance",
          "command": "bench test --check root-conformance",
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
  },
  "amendments": [
    {
      "from": "sha256:b4a2344dfd68af6535b28d63451857b2045ac3d53de755aa665f8ab08e68ba5a",
      "to": "sha256:143f753bd2b8510ff891b8076aa2e2f1bc7d8ddafaf78f72760bf7da7abc0076",
      "chunk_ids": {
        "BG-C0": [
          "BG-C0"
        ]
      }
    },
    {
      "from": "sha256:143f753bd2b8510ff891b8076aa2e2f1bc7d8ddafaf78f72760bf7da7abc0076",
      "to": "sha256:d3cd8ebcd312d6538a9422dac4b8cf86c3760a8c9d9b0e67f841d3ca08e2c596",
      "chunk_ids": {
        "BG-C0": [
          "BG-C0"
        ]
      }
    }
  ]
}
```

## BG-C0 review

The three axes reviewed the frozen pair `84e9b1d37..7bdc95c72`. They found 13 raw findings and no blocking defect. The coordinator collapsed them into 5 repair targets on ticket 01, and it sends all of them to one combined repair session.

## Standards

Count: 5. Worst issue: the Affinity column of the priority table names no read method.

- BG-C0-S1, auto-fix, confidence 5: `specs/background-grade/census.md:83` names no source for the Affinity column. The sampler read it with `taskset -cp`. Target T2.
- BG-C0-S2, auto-fix, confidence 6: `specs/background-grade/census.md:127` lists four packages in one sentence. The record must use a vertical list. Target T3.
- BG-C0-S3, auto-fix, confidence 4: `specs/background-grade/census.md:142` states the debug-loop results with no figures. Target T1.
- BG-C0-S4, auto-fix, confidence 3: `specs/background-grade/census.md:64` uses the present tense for a past run. Target T3.
- No-op, confidence 3: `specs/background-grade/census.md:45` holds session-local paths. No rule forbids them in a measurement record.

## Spec

Count: 2. Worst issue: the waiver cites a sealed-proxy rerun with no recorded data.

- BG-C0-P1, auto-fix, confidence 5: `specs/background-grade/census.md:148` claims a sealed-proxy rerun with no wall time, command, or exit. Target T1.
- No-op, confidence 7: `specs/background-grade/census.md:130` discloses the unsealed proxy limit, and the limit does not change the verdict.

BG93, BG94, and BG95 are met. The medians recompute to 237.73 s and 298.98 s, and the verdict `stop` is correct.

## Coverage

Count: 6. Worst issue: the waiver and the contention claim rest on runs that the record does not hold.

- BG-C0-C1, auto-fix, confidence 7: `specs/background-grade/census.md:136` and `:148` rest on unrecorded runs. Target T1.
- BG-C0-C2, auto-fix, confidence 7: `specs/background-grade/census.md:89` names no affinity read. Target T2.
- No-op, confidence 6: `specs/background-grade/census.md:86` samples no race or system phase. Those phases are children of the sampled `proxy.sh`, and nice and the I/O class pass to each child.
- BG-C0-C4, auto-fix, confidence 5: `specs/background-grade/census.md:80` gives no restart time. `proxy.cmds` records pass 2 at epoch 1793131699.40. Target T4.
- BG-C0-C5, auto-fix, confidence 4: `specs/background-grade/census.md:91` gives sample counts that do not reconcile. Target T4.
- BG-C0-C6, auto-fix, confidence 4: `specs/background-grade/census.md:51` does not state that the foreground and the proxy shared one worktree and one Go cache. Target T5.

The implementation command contributed to no finding.
