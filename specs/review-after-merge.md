# review-after-merge — a diff surface for work that already landed

Status: staged

Source: `ASSESSMENT.md` backlog 3 (findings §1 high, §3 med-high; both hit live
on 2026-07-06 reviewing the ft9 landing commit).
Drafted without a decision map under the reviewer's 2026-07-06 batch approval;
default calls are flagged in the implementation decisions for post-hoc veto.

## Problem

`/bench-review-implementation` sources its diff from `bench diff --full` and
demands it be non-empty — but on this repo's documented happy path (the build
commits to the working branch, then review runs), the work is already at HEAD
of the default branch, merge-base equals HEAD, and the diff is empty every
time. The phase's step 1 names no remedy, so a review requested after the work
lands has no lifecycle representation. This was hit live this session.

## Solution

`bench diff --commit <sha>` bounds the review context to one landing commit:
base becomes the commit's first parent, so the files table, log table, and (with
`--full`) the raw diff body describe exactly what that commit landed — for a
merge commit, everything the merge brought in. `/bench-review-implementation`
step 1 documents the fallback: when the diff is empty because the work already
landed, review the landing commit with `bench diff --full --commit <sha>`.

## User stories

1. As a review agent, I want `bench diff --commit <sha>` to emit the standard
   preamble and `files[N]{status,path}` table for the range from the commit's
   first parent to the commit, so I can pin a review to a landed change with
   the same output contract the branch path has.
   Line: claude-sonnet-5 / medium. This is a thin range-selection addition at
   the existing diff command seam, fully observable by the AXI contract tests,
   mirroring the ft9 routing for the same seam.

2. As a review agent, I want `--commit` to compose with `--full` (log table =
   commits the landing commit brought in; diff body = first-parent diff), so a
   review-after-merge gets the whole bundle in one call.
   Line: claude-sonnet-5 / medium. Same seam and observability as story 1.

3. As a review agent, I want the edge postures pinned: a root commit (no
   parent) and an unresolvable sha exit 1 with a structured error; `--commit`
   without an argument or with extra arguments exits 2; a merge commit's bundle
   reflects the full merged-in change.
   Line: claude-sonnet-5 / medium. Edge behavior at the same gate-observable
   seam.

4. As a reviewer running `/bench-review-implementation`, I want step 1 to name
   the empty-diff remedy — landed work is reviewed via
   `bench diff --full --commit <sha>` — so the phase no longer dead-ends on its
   own happy path.
   Line: claude-fable-5 / high. This is net-new guidance prose in a workflow
   command file, which the profile's doc-authoring leverage override routes to
   the top tier (unlike ft9's mechanical call-site swap, which deviated down).

## Implementation decisions

- **Range semantics.** `--commit <sha>`: base = `<sha>^`, head = `<sha>` — files
  via `git diff --name-status --no-renames -z <sha>^ <sha>`, log via
  `git log <sha>^..<sha>`, body via `git diff <sha>^ <sha>`. For a merge commit
  the first-parent range shows everything the merge landed; for an ordinary
  commit it is that commit's change. No three-dot form is needed — the base is
  an exact parent, not a divergent branch. (Default call, flagged: an
  alternative `--merge-only` view of second-parent history is out of scope.)
- **Preamble stays the contract**: `branch:` unchanged, `base: <sha^>`,
  `method: commit <sha>` — the method line is how a reader (and a test) sees
  that recorded-base resolution was bypassed. `benchBase` resolution is skipped
  entirely under `--commit`; the flag is an explicit override.
- **Sha validation is loud**: the argument is verified as a commit
  (`rev-parse --verify <sha>^{commit}`) before any section renders; a root
  commit's missing parent is its own structured error naming the case, not a
  git stderr leak.
- **Arg parsing grows real iteration.** `diff.Command`'s current
  positional-switch parsing can't host a value-taking flag; it becomes a small
  loop (the shape `bench commit`'s parseArgs already uses). This touches the
  same lines as the cli-contract-accuracy spec's unknown-argument-attribution
  fix — build that spec first or land this one with the attribution fix
  included; the two specs note the ordering mutually.
- **Doc anchor follows the ft9 pattern**: the step-1 fallback phrase is pinned
  by tightening the existing `bench diff --full` anchor in
  `internal/conformance/docs_workflow_helpers_test.go` to also require the
  `--commit` remedy token in `bench-review-implementation.md`.

## Testing decisions

- **What a good test is here:** drive `bench diff --commit` / `--full --commit`
  in a real git fixture (linear commits plus one merge) and assert stdout
  sections and exit codes at the command seam. Prior art:
  `internal/contract/axi/axi_wave2_test.go` (`testAXIDiffRecordedBase`,
  `testAXIDiffFallbackShape`, `testAXIDiffErrorPosture`).
- **Seam:** one — `diff.Command` stdout, exercised by AXI contract tests.
- **Gate:** the project gate (`bench gate`); the story-4 anchor lands in the
  conformance phase.

### Seam diagram

    trigger: review agent (or /bench-review-implementation step 1 fallback) runs
             `bench diff [--full] --commit <sha>`
        │
        ▼
    args: [--full] [--commit <sha>]  ──▶  [ diff.Command                    ]  ──▶  stdout:
    repo state                       ──▶  [  verify sha, base = sha^        ]        branch / base: <sha^> /
                                          [  files table (shared renderer)  ]        method: commit <sha>
                                          [  +full: log[N], raw diff body   ]        files[N]{status,path}
                                          [                                  ]        log[N] + diff_body (--full)
                      ◀ tests attach here: AXI contract drives the built binary against a
                        fixture with a landed merge and asserts the preamble, section
                        lines, `@@` markers, and the exit postures for root/bogus shas.

### Acceptance coverage map

| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| 1 | `bench diff --commit <sha>` emits base = first parent, method `commit <sha>`, and the commit's files table | diff.Command stdout (AXI contract) | `bench diff --commit HEAD` in a fixture — today exits 2 `unknown argument`, so RequireExit(0) and the `method: commit` assertion fail | a flag that is unrecognized, or that ignores the sha and reuses branch-base resolution, fails the method/base assertions |
| 2 | `--full --commit <sha>` appends the landing commit's log table and raw diff body | diff.Command stdout (AXI contract) | same fixture with `--full --commit` — exits 2 today | a compose failure (either flag disabling the other) drops the `log[` header or the `@@` tail and fails the row |
| 3 | a merge commit's bundle shows the full merged-in change (files + brought-in commits + body) | diff.Command stdout (AXI contract) | fixture with a two-branch merge — no assertion can pass today | pins the first-parent semantics; diffing against the wrong parent or merge-base yields a different file set and fails |
| 3 | root commit exits 1 with a structured no-parent error; bogus sha exits 1; `--commit` missing its value or given twice exits 2 | diff.Command stdout (AXI contract) | probes against the fixture's root commit and a garbage sha — all exit 2 (unknown argument) today with the wrong text | wrong exit postures or leaked git stderr fail the exact assertions, extending testAXIDiffErrorPosture |
| 1 | bare `bench diff` and `bench diff --full` behavior is byte-unchanged | diff.Command stdout (AXI contract) | already covered — the existing contract rows pin both paths; re-run green after the parser rewrite | the arg-parser rewrite is the risk; the existing suite is the regression net, stated openly rather than re-promised as new TDD |
| 4 | `/bench-review-implementation` step 1 names the `--commit` fallback for landed work | conformance docs anchor | tighten the docs anchor to require the `--commit` token in the step-1 prose — red until the doc carries the remedy | the phase doc is where the dead-end lives; the anchor pins the remedy the same way ft9 pinned the original repoint |

### Edge inventory

- error path → rows: root commit, unresolvable sha, flag-without-value.
- empty/absent input → covered: an empty commit (no file changes) renders
  `files[0]` and empty sections at exit 0 through the shared renderers already
  pinned for empty-since-base; asserted in the story-1 fixture with an
  `--allow-empty` commit.
- boundary values → row: merge commit vs ordinary commit (story 3).
- malformed input → **Won't handle** beyond sha verification: ref names, `@{}`
  syntax, and abbreviations are resolved by `rev-parse` exactly as git does —
  the command adds no ref grammar of its own.
- control bytes in subjects → **Won't handle** here: the log table's refusal
  posture is owned by the cli-contract-accuracy spec's row; `--commit` reuses
  the same renderer, so one owner keeps the fact single-sourced.
- interrupted/partial state, re-run idempotency — **Won't handle**: read-only
  and deterministic for a fixed tree (same posture the ft9 spec recorded).
- hostile environment (detached HEAD) → covered: `branch:` already renders
  `(detached)` via the existing preamble path; the story-1 fixture asserts the
  preamble under a detached checkout since landed-commit review often runs
  there.

## Out of scope

- **`--commit` accepting a range or multiple shas** (review several landings in
  one bundle) — a separate aggregation capability; one landing commit per
  review call matches the phase's unit of work. Estimate: ~3 edits, 2 gate
  runs.
- **A second-parent (`--merge-only`) view** — a distinct forensic capability
  for archaeology on merge topology, not review-after-merge. Estimate: ~2
  edits, 2 gate runs.
