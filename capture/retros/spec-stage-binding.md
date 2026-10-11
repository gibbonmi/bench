## Outcome

The spec-stage-binding spec landed as `88de8d9d` on 2026-10-11 from the source pair `8c8a6ca5..7130e46e`. The whole-project gate was green. A spec-stage close now lands its staged spec and its deliverable binding in one spec-less landing. Plan, approve, and the inventory read a new deliverable at the checkout HEAD. The approval check of `bench commit` and `bench worktree land` reads it at the candidate tree.

`bench handoff` writes one `Spec binding` line for each live spec, and the write-spec Exit handoff states the close order. Three chunks and three tickets landed with 46 acceptance rows, and the delivery settles FT395.

The reviewer was away for the whole build and deferred every question to an opus/high read-only consultant. Repair rounds and spec or ticket expansions had advance approval. One consultation ran: it verified the stale-base refusal of a binding commit and chose the fold of `main` before the plan. The spec gained two rows during the build, SB45 for the fold route and SB46 for the changelog section.

## Gate-stage timings

- landing: commit 88de8d9d16ea1ab595687515504cbf006e66a5f2, trace 9dfcb97c362d328215a02808218c039d
- gofmt: 159 ms
- vet: 1335 ms
- test: 157134 ms
- race: 3244 ms
- system: 58858 ms
- shellcheck: 565 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | expected 3 repair rounds | held | 6 | 2 repair cycles by the repair-attribution table | fable / medium / orchestrator |
| ticket 1 | every coverage row red then green | held | 9 | coordinator omit probe bit at caa476a6 | opus / xhigh / implementer |
| ticket 2 | every coverage row red then green | held | 9 | coordinator omit probe bit at 1350865e | opus / xhigh / implementer |
| ticket 3 | every anchored sentence red then green | held | 9 | coordinator live-tree omit probe bit at 33f0ce69 | opus / xhigh / implementer |

Brier mean: 0.05 over 4 pairs; 0 abstained.

Every author stayed inside its fence and committed on a lane pass in one attempt. Every ticket matched its chunk table row, and no ticket needed a fence expansion. Each author observed its rows red before the production edit, and each named probe bit with a restored subject. Each author ran its probes but kept only the central one in the record, so every chunk's first review found the same record gap.

The harness reported this cost for each first-pass commit and repair.

| ticket | role | effort | tokens | harness time | outcome |
|---|---|---|---|---|---|
| 1 | author | xhigh | 396k | 1772 s | landed; review found two comment sentences, a stale comment, two silent mutants, and the stale-base refusal |
| 1 | repair | xhigh | 245k | 1204 s | closed S2, S3, C2, C3, and SB45 in one attempt, plus one re-record |
| 2 | author | xhigh | 280k | 1115 s | landed; review found a stale pins comment and a silent milestone-scoped mutant |
| 2 | repair | xhigh | 180k | 700 s | closed S2 and C1 in one attempt, plus one re-record |
| 3 | author | xhigh | 266k | 1069 s | landed; review found one long comment sentence and no coverage row for the changelog section |

## Coordinator catches

- The review record of every chunk kept only the central probe. The orchestrator reran each author's probes at the clean tip and retained them, and the later repair charges asked for one file per probe.
- A binding commit on a planning base behind `main` refuses, which the SB1 Coverage axis reproduced. The consultant verified the cause in `branchReference`, and the spec gained the fold step and SB45.
- A prose-only plan commit after the repair's refreeze moved the SB1 source digest, so the chunk needed a second refreeze and a verification rerun. The lesson went to `bench learning`.
- The SB1 round 2 Spec axis found that the Won't-handle line described a second fold on a dirty tree, which `bench worktree merge` refuses. The line now names `git restore`, a fold, and a new plan.
- Two comment-only corrections, a provenance tag in SB2 and a 31-word sentence in SB3, each cost a refreeze and an author re-record. The source digest covers test and comment bytes.
- The lane's prose check refused the review pickup four times on sentence and paragraph bounds. The pickup prose now splits every long sentence before the commit.
- The `--check-current` bind of an axis refused once outside `bench worktree exec`, and the rerun inside the worktree passed.
- Three repair sessions closed every finding in one attempt, and every independent probe of the confirming rounds bit.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-bind-a-new-deliverable-from-its-publishing-tree.md | 1 | spec-row |
| 2-name-the-binding-state-in-the-handoff.md | 1 | spec-row |
| 3-state-the-one-landing-close.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Add `bench probe --occurrence <n>` and echo the mutation text in the probe response, because a repeated text needed a no-op swap and the logs lost their commands.
  Feeds: new
- Let `bench record verification` read the last spilled response, because every author copied each response into an excerpt file by hand.
  Feeds: FT405
- Print the whole packages table of `bench test` when the body is under about one kilobyte, because a nine-package floor spilled after three rows.
  Feeds: new
- Let `bench record chunk` take `--amend`, because a seam-cell rewrite moved the plan digest on two of three freezes.
  Feeds: new
- Let one refreeze re-stamp the author verification entries of a chunk when only comment or prose bytes changed, because three re-record passes followed comment-only corrections.
  Feeds: FT405
- The landing census counted 201 raw calls: rg 194, awk 2, cat 2, head 1, mkdir 1, python3 1. A symbol search has `bench consumers`, and the rest have no Bench form.
  Feeds: new

### Skills

- Make the author charge require one excerpt file per probe with its command line first, because the record keeps only the probes the author files.
  Feeds: none
- Make the implement-spec command batch prose-only spec corrections before the repair's refreeze, because a later plan commit moves the source digest.
  Feeds: new
- Make the ticket-slicing checks decide a fold step for a close whose worktree may trail `main`, because the stale-base refusal reached review, not the spec.
  Feeds: none

### Process

- Keep the opus/high consultant as the reviewer's stand-in for a build judgment call, because its one consultation verified the cause and chose the cheaper route.
  Feeds: none
- Write the review pickup in sentences of 25 words or fewer from the first draft, because the lane graded the record four times.
  Feeds: none
- Keep each author session open until its chunk closes, because every re-record after a refreeze took 3 to 4 minutes.
  Feeds: none