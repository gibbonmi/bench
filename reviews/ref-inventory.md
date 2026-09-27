# Ref inventory review record

## RI-C1a author evidence

The plan amendment `2d3535b6` declared the version 2 execution plan on the
`main` tip `72a749a3`, which is the RI-C1a base. The serial order is ticket 6,
then ticket 1. The chunk tip is `c6b2a74e`.

Ticket 6 had a fresh `bench-writer` author on opus at high effort. The author
started at `2d3535b6` and committed `f3ff2543` on a lane pass with the build
preflight green. The delta is two glossary lines in `CONTEXT.md`. The eight
anchors kept their needles and lines. Rows RI57 and RI82 are review-owned, so
they carry no test red; the author reported both as claimed at confidence 9.
The orchestrator read the diff, confirmed a clean tree, and ran the coordinator
probe: an omission of the reader-sweep anchor line under
`docs-currency-workflow`, verdict `bit`, restore `yes`.

Ticket 1 had a fresh `bench-writer` author on opus at high effort from the plan
commit `e3b00c6d`. The author stopped before its commit on a behavioral spec
contradiction: ancestry is transitive, so a ref one commit under an
ancestry-landed ref is itself landed, and RI9 as written was unreachable. The
omission probe of that holder kind was an equivalent mutant. The orchestrator
routed the question to the reviewer through `codex exec` on `gpt-6-astra` at
xhigh reasoning, by user direction on 2026-09-27. The decision was option B:
drop the ancestry-landed holder kind everywhere. The plan commit `e790e56a`
amended story 9, story 63, RI9, RI82, the holder rule, ticket 1, and ticket 6.

The ticket 1 charge evidence `sha256:d03d7a73…` predates that amendment. The
author read the amended spec and ticket bytes from the tree at `e790e56a`,
because the dirty source refused a regenerated charge. The same author session
then committed `2ace6d83` on a lane pass with the build preflight green, and it
reported all 34 rows verified. The named probe, subsumed tested before landed,
bit six cases including RI62 and RI9, restore `yes`. The orchestrator ran the
coordinator probe: an omission of the retained-row skip in the unclaimed apply
loop over `TestCleanUnclaimed`, verdict `bit` with four tests red, restore `yes`.

Ticket 6 reopened for the glossary `holder` sentence. The same author session
committed `c6b2a74e` on a lane pass with the build preflight green. The eight
anchors kept their needles and lines, and the `holder` entry names only a
recorded assignment branch or a unique root.

Two author-discipline notes carry no tree defect. The ticket 1 author wrote one
test block into the pool path with a shell heredoc before it moved to the Edit
tool. It also made one string swap with `sed -i`. The committed bytes passed
the lane and the checks. The ticket 1 charge named `ledger.go` as the
shift-prefix source. The constant lives in `internal/intent/ledger/validate.go`.

### Verification

Each author reran its plan checks at the chunk tip `c6b2a74e` on a clean tree.
The JSON payload holds each result. The `internal/conformance` run skipped the
same three capability tests as the baseline: two need a unix socket and one
needs a character device. The first ticket 1 conformance rerun failed because
the orchestrator's untracked draft of this record was over the prose bound at
that moment. The draft was repaired, and the retained rerun passed.

```bench-review-record
{
  "version": 2,
  "spec": "specs/ref-inventory/spec.md",
  "plan_digest": "sha256:dedf49a52da69c6ba1418a46f9d9f3586cba01ac634110edf19f32d62ad5ee99",
  "implementation_session": "",
  "chunks": [
    {
      "id": "RI-C1a",
      "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
      "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
      "plan_digest": "sha256:dedf49a52da69c6ba1418a46f9d9f3586cba01ac634110edf19f32d62ad5ee99",
      "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
      "acceptance_rows": [
        "RI1", "RI2", "RI3", "RI4", "RI5", "RI6", "RI7", "RI8", "RI9", "RI10", "RI11", "RI12", "RI13", "RI14", "RI15",
        "RI17", "RI18", "RI19", "RI20", "RI21", "RI22", "RI23", "RI25", "RI28", "RI55", "RI57", "RI59", "RI60", "RI61",
        "RI62", "RI63", "RI66", "RI72", "RI73", "RI81", "RI82"
      ],
      "verification": [
        {
          "id": "ri-c1a-6-anchors-r1",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-anchors@c6b2a74e",
            "digest": "sha256:80be124f957a7a1d26cd487fc038f3d197227ddf17c6bd7937649af23f57371d",
            "excerpt": "github.com/gibbonmi/bench/internal/anchors,pass,937"
          },
          "requirement": "6-anchors",
          "command": "bench test --package ./internal/anchors",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-6-conformance-r1",
          "performer": "claude:bench-writer/ri-t6-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t6-author-20260927/6-conformance@c6b2a74e",
            "digest": "sha256:1fad3031dc2da3fcafd69f0df5d0adb2237368d60772997c3ec6c354ef6b8682",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,43922"
          },
          "requirement": "6-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-worktree-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-worktree@c6b2a74e",
            "digest": "sha256:df17c9b605a27341a7b94b29354dcd35fa2bd39ba05dc2fb28244fa164e9f9b8",
            "excerpt": "github.com/gibbonmi/bench/internal/worktree,pass,52743"
          },
          "requirement": "1-worktree",
          "command": "bench test --package ./internal/worktree",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-status-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-status@c6b2a74e",
            "digest": "sha256:8c89907f1a191473cd2b9f9ad6edefca805768ef41cc8727ebbc472a99bd11cb",
            "excerpt": "github.com/gibbonmi/bench/internal/status,pass,12750"
          },
          "requirement": "1-status",
          "command": "bench test --package ./internal/status",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-cmd-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-cmd@c6b2a74e",
            "digest": "sha256:59f6131664735a60e2038849a565fb2922d7ed0a58afd1e3e560d13be7606bd0",
            "excerpt": "github.com/gibbonmi/bench/cmd/bench,pass,8319"
          },
          "requirement": "1-cmd",
          "command": "bench test --package ./cmd/bench",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-conformance-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-conformance@c6b2a74e",
            "digest": "sha256:0d336fac346b0b7813079779d353daffd315eafee9b0af3a05b5d16e61e83f97",
            "excerpt": "github.com/gibbonmi/bench/internal/conformance,pass,36775"
          },
          "requirement": "1-conformance",
          "command": "bench test --package ./internal/conformance",
          "exit_code": 0
        },
        {
          "id": "ri-c1a-1-system-r1",
          "performer": "claude:bench-writer/ri-t1-author",
          "role": "author-verification",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "pass",
          "native_ref": {
            "ref": "claude:agent/ri-t1-author-20260927/1-system@c6b2a74e",
            "digest": "sha256:ab56672c31a42be0c763cb5a3e5345a5cec9edfcfd6b6e17bb7cabafd68f2196",
            "excerpt": "github.com/gibbonmi/bench/internal/systemtest,pass,38997"
          },
          "requirement": "1-system",
          "command": "bench test --check system",
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
