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

## RI-C1a chunk review, round 1

The frozen pair is base `72a749a3` and tip `c6b2a74e`. The record commit
`2e460e05` follows the tip, and the shared evidence is `sha256:5f24a383…`.
Each axis ran in a fresh `bench-reviewer` session on opus at high effort, the
conditional review line for an opus implementation. Only the Coverage axis ran
probes, and it left the tree clean.

Raw findings: Standards 6, Spec 2, Coverage 4. Repair targets after the fold:
one ticket 1 repair session with nine items, and two spec wording repairs by the
orchestrator. Two questions went to the reviewer route, `codex exec` on
`gpt-6-astra` at xhigh, and both answers are recorded below as decisions.

## Standards

Findings: 6. Worst: R1.

- R1 (auto-fix, confidence 7): AGENTS.md test-expectation exception, "someone
  must record and demonstrate that red". `clean_classes_test.go:49-50`,
  `clean_unclaimed_test.go:67`, and `clean_unclaimed_test.go:301` hand-spell
  `class=`, ` holder=`, and `retained: content main lacks`, which production
  owns at `clean_classes.go:27-29`. The record holds no red for those
  spellings. The repair records one probe red per spelling.
- R2 (auto-fix, confidence 5): craft-comments register. The comment at
  `clean_classes_test.go:44-46` argues its own case and claims a red no record
  holds. It folds into R1.
- R3 (auto-fix, confidence 5): AGENTS.md one source per fixture harness.
  `squashIntoMain` at `clean_classes_test.go:20-26` repeats `squashLand` at
  `squash_landed_test.go:30-36`.
- R4 (auto-fix, confidence 4): craft-comments one source per fact. Three
  comments restate that a faulted set has no fingerprint and no apply:
  `clean_classes.go:32-33`, `clean_unclaimed.go:33-35`, `clean_unclaimed.go:124-126`.
- R5 (auto-fix, confidence 3): craft-comments contract in one place. The doc
  comment at `clean_classes.go:79-84` restates the holder rule the glossary and
  the spec own; the why clauses stay.
- R6 (no-op, confidence 3): the pool-path heredoc and the `sed -i` swap of the
  ticket 1 author. No tree defect, and the committed bytes passed the lane. The
  orchestrator captured the bypass as a learning.

Advice: the unclaimed row fields are built by hand at four sites, unchanged in
count by the diff. The active-or-cleanup-pending check is the fourth hand-written
one in the package. `reachableFromAHead` reads oddly.

## Spec

Findings: 2. Worst: R7. Every RI-C1a row is delivered, and the holder rule
matches both reviewer decisions.

- R7 (auto-fix, confidence 8): row RI5 read "one commit", but a one-commit
  squash lands by patch containment, so a classifier that stops at cherry prints
  `landed`. The author's two-commit fixture follows story 5. The reviewer route
  confirmed the row now reads "two commits".
- R8 (auto-fix, confidence 4): row RI13 and story 13 read "every unclaimed
  row", but an error row carries no class prefix by RI11. The reviewer route
  confirmed both now read "classified row".

Advice: 32 rows still carry the `planned` seam although their tests exist. The
orchestrator updates the seam cells at the chunk close. The axis also asked
whether a content-landed recorded active branch should hold; that question
became decision D2 below.

## Coverage

Findings: 4. Worst: R9.

- R9 (ask-user, confidence 9): an unrecorded symbolic ref in a Bench namespace
  is selected, classifies through its target, and the apply's `update-ref -d`
  dereferences it. Probes 4 and 5 deleted a unique root and the checked-out
  `main`. No row decided symrefs. Decision D1 below routes it to an error row
  inside ticket 1's fence, and the landing prune's shared delete is captured
  as an idea.
- R10 (auto-fix, confidence 9): a probe that dropped `class` from the
  fingerprint parts stayed silent, because the RI63 fixture changes both the
  class and the holder. New row RI85 pins a class-only change.
- R11 (auto-fix, confidence 8): a probe that admitted every record state as a
  holder stayed silent. New row RI86 pins a complete record's branch.
- R12 (auto-fix, confidence 8): a probe that limited holders to active records
  stayed silent. New row RI87 pins a cleanup-pending holder.

Probes, each with restore `yes`:

- swap of the fingerprint class part: silent
- swap of the holder filter to active only: silent
- swap of the holder filter to every state: silent
- a symref fixture at a unique root: bit
- a symref fixture at `main`: bit

Advice: `PruneLandedBranches` shares the dereferencing delete, outside this
spec. `LandedInDefault` takes the default branch as a short name.

## RI-C1a decisions and repair cycle 1

The reviewer route decided D1: a symbolic ref produces an error row, never a
silent exclusion, and `DeleteBranchExact` stays unchanged in FT199 (row RI83).
It decided D2: a recorded assignment branch that is landed by content only is
not a holder (rule sentence and row RI84). It confirmed R7 and R8. The plan
commit that follows this record carries RI83 to RI87, the wording repairs, and
the repair session assignment.

One fresh repair session for ticket 1 on opus at high consumes repair cycle 1
of 2. It owns R1 to R5, R9 to R12, and D2. The
implementation command contributed to R10, because the author's probe list
named a fingerprint-binding swap without a field. The orchestrator captured the
per-field probe rule as a learning.

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
        "RI62", "RI63", "RI66", "RI72", "RI73", "RI81", "RI82", "RI83", "RI84", "RI85", "RI86", "RI87"
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
      "reviews": [
        {
          "id": "ri-c1a-r1-standards",
          "performer": "claude:bench-reviewer/ri-c1a-standards",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-standards@c6b2a74e",
            "digest": "sha256:d6f3966b440ea17dcb22ec20b97b50cb877761e05779ba5204d1264ed001e961",
            "excerpt": "Standards: 6 findings. Worst: the new tests spell class=, holder=, and retained: content main lacks independently of the production constants with no recorded red."
          },
          "axis": "Standards",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "finding_ids": ["R1", "R2", "R3", "R4", "R5", "R6"],
          "supersedes": []
        },
        {
          "id": "ri-c1a-r1-spec",
          "performer": "claude:bench-reviewer/ri-c1a-spec",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-spec@c6b2a74e",
            "digest": "sha256:e41ea229c8be77d3e2137c89f35bab643cf354dfdc9efa96af7473a9dbf8bf87",
            "excerpt": "Spec: 2 findings. Worst: row RI5 reads one commit, but a one-commit squash lands by patch containment; the two-commit fixture follows story 5."
          },
          "axis": "Spec",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "finding_ids": ["R7", "R8"],
          "supersedes": []
        },
        {
          "id": "ri-c1a-r1-coverage",
          "performer": "claude:bench-reviewer/ri-c1a-coverage",
          "role": "independent-review",
          "model": "opus",
          "effort": "high",
          "source_digest": "be01c3ac5a00497a0c7e6631aa38ee47a4c2b80e",
          "state": "completed",
          "outcome": "fail",
          "native_ref": {
            "ref": "claude:agent/ri-c1a-coverage@c6b2a74e",
            "digest": "sha256:a09849151032402867990fe921a2563d7548a14599b378b9ca010743682b501f",
            "excerpt": "Coverage: 4 findings. Worst: a symbolic ref in a Bench namespace makes the bulk apply delete its target, observed on a unique root and on the checked-out main."
          },
          "axis": "Coverage",
          "base": "72a749a35dc4c37de87f56534959ac9c137099e1",
          "tip": "c6b2a74effa57cbafdd8bc449c9fbbb24d1468eb",
          "finding_ids": ["R9", "R10", "R11", "R12"],
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
