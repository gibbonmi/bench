## Outcome

The build-directions spec landed as `287f23d8` on 2026-10-10 from the source pair `dcf7fad0..ee538f7c`. The whole-project gate was green, and the landing retired FT410 and no other row. It states nine reviewer directions in the guidance and adds one git-ignored build log. Four chunks and six tickets landed, with 64 acceptance rows and one added story.

By reviewer direction, the build also ran an Opus effort benchmark. Tickets 01 and 04 took medium and low shadow authors, and ticket 06 took xhigh, medium, and low shadow authors. Each shadow started from the canonical author's source tip in its own worktree and never landed. Three blind Fable/high graders scored every first-pass commit under shuffled labels.

## Gate-stage timings

- landing: commit 287f23d8935168f2acdcfb98421b07e82c17ed42, trace 991e471b447d57c0f29a773848e8cce2
- gofmt: 144 ms
- vet: 1392 ms
- test: 167028 ms
- race: 2393 ms
- system: 59590 ms
- shellcheck: 582 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | expected 2 repair rounds | refuted | 6 | 3 repair cycles by the repair-attribution table | opus / high / orchestrator |
| ticket 01 | every coverage row red then green | held | 9 | coordinator swap probe bit at b2f396fb | opus / high / implementer |
| ticket 02 | every coverage row red then green | held | 9 | coordinator swap probe bit at 486fd1ef | opus / high / implementer |
| ticket 03 | every coverage row red then green | held | 9 | coordinator swap probe bit at d3bd78a9 | opus / high / implementer |
| ticket 04 | every coverage row red then green | held | 9 | coordinator swap probe bit at 823bca72 | opus / high / implementer |
| ticket 05 | every coverage row red then green | held | 9 | coordinator swap probe bit at ed901b25 | opus / high / implementer |
| ticket 06 | every coverage row red then green | held | 9 | coordinator omit probe bit at 4ed02fb8 | opus / high / implementer |

Brier mean: 0.06 over 7 pairs; 0 abstained.

Every canonical author stayed inside its fence and committed on a lane pass. Every ticket author put the required sentences byte-exact; the slices matched the chunk table.

The benchmark gave this cost and blind grade for each first-pass commit. Tokens are the harness-reported subagent count, which can be the final context size.

| ticket | effort | tokens | wall time | blind grade (rank) |
|---|---|---|---|---|
| 01 | high | 87.8k | 566 s | 8 (2nd) |
| 01 | medium | 80.4k | 436 s | 7 (3rd) |
| 01 | low | 56.5k | 376 s | 8 (1st) |
| 04 | high | 116.8k | 466 s | 9 (1st) |
| 04 | medium | 99.7k | 433 s | 9 (2nd) |
| 04 | low | 58.8k | 300 s | 8 (3rd) |
| 06 | xhigh | 265.2k | 1370 s | 9 (1st) |
| 06 | high | 168.3k | 977 s | 8 (4th) |
| 06 | medium | 138.4k | 1004 s | 9 (2nd) |
| 06 | low | 111.0k | 919 s | 8 (3rd) |

No variant had a blocking defect, and every variant met every acceptance row. On the guidance tickets, low put the same production bytes as high at about half the tokens, with thinner red evidence. On ticket 06, only xhigh avoided the stale doc comment and the pasted fixture harness. No effort level caught the two spec-stage gaps.

The reviewer adopted the result as routing. Low serves an exact-spec guidance ticket, and medium is the code-ticket default. Xhigh serves a code ticket where a repair round costs more than about 100k tokens.

## Coordinator catches

- All three ticket 01 authors found that the planned probe selector graded a synthetic harness tree. The orchestrator moved each anchor-row probe to the live check.
- The orchestrator recorded the worktree owner digest as the first assignment id and corrected it from the evidence binding before dispatch.
- The orchestrator took the NRO record-commit base for BD-C3; the checkpoint required the BD-C2 tip, because two comment-only corrections followed it.
- A blind grader found a pasted fingerprint fixture harness that the BD-C4 review missed; the orchestrator added it to the ticket 06 repair.
- The ticket 06 repair charge missed the package test-count pin file; a plan commit widened the fence.
- The orchestrator committed the P21 story fix after the BD-C4 chunk tip, which forced a refreeze and a fresh pass on every axis.

The reviewer asked for a comparison of this spec against the recent write-spec changes. The bare post-ticket preflight at implement-spec:50 (BD62) falls under the pre-review checklist items `Consolidated rules` and `Quantified obligations`, which landed before this spec. The spec answered `Consolidated rules: none`, so the checklist was designed to catch it and missed it.

The probe selector on a synthetic-tree test falls outside every write-spec change. The grader trace binds each row to a grader. No item checks that a planned probe selector reads the probed file. The source-cleanup caller (BD64) falls under `Changed-function callers`, and the spec listed the caller. That item asks for a list, not a row for each caller whose behavior changes, so the gap is outside its scope.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-fold-only-at-chunk-boundaries.md | 1 | spec-row |
| 02-hold-light-paths-during-a-build.md | 0 | none |
| 03-close-the-last-chunk-through-the-complete-checkpoint.md | 1 | spec-row |
| 04-cap-witnesses-and-keep-moved-guards.md | 1 | other |
| 05-scope-review-rounds.md | 0 | none |
| 06-keep-the-build-log.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Add a read-only record status verb that lists the required chunk base, each stale or missing entry, and each commit that breaks the chain.
  Feeds: FT379
- Add a record next verb that prints the next required record or plan action and its preconditions.
  Feeds: new
- Give bench probe and bench test an --out option that writes the complete response for the record verb.
  Feeds: new
- Make each refusal name its field and its route, such as the assessment field, the worktree exec form, or a wrapped swap.
  Feeds: new
- Make bench record assignment derive the assignment from its worktree, or refuse a mismatch.
  Feeds: new
- Make bench test list each failed want line of a multi-rule harness, and make bench commit print the published commit hash.
  Feeds: new
- Name the Grep tool for worktree search in review and author charges, because the census counted 51 raw calls, 42 of them rg.
  Feeds: new

### Skills

- Add a pre-review checklist item that confirms each planned probe selector reads the file that the probe mutates.
  Feeds: FT380
- Make the Changed-function callers item give each caller whose behavior changes an acceptance row or a Won't handle line.
  Feeds: new
- When two guidance files must agree on a rule, make the spec point one file to the other instead of pinning two copies.
  Feeds: new

### Process

- Fold a review's spec-only fix into the repair commit, never into a plan commit after the chunk tip, because the late commit forced a refreeze.
  Feeds: FT380
- Take the predecessor chunk tip as the next chunk base whenever a comment-only correction follows that tip.
  Feeds: FT380
- Name the package test-count pin file in a charge that adds a top-level test.
  Feeds: FT293
- Let one repair session record every ticket of its chunk, because eight verification re-record sessions cost about 25 to 30 minutes after refreezes.
  Feeds: FT405
- Add an evidence-only confirmation form, because three axis passes after the BD-C3 rebase and the story 47 refreeze graded no code.
  Feeds: FT379