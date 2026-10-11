## Outcome

The commit-merge-parent spec landed as `85f2b05c` on 2026-10-11 from the source pair `342975e3..d23f5eff`. The whole-project gate was green. During a pending merge, `bench commit` now records the previous tip and each merge head as parents. It refuses a name set that omits a merged path, ends the merge through Git, and grades the branch's own change. Three chunks and five tickets landed with 62 acceptance rows, and the delivery closed no roadmap row by design.

By reviewer direction, the build ran an Opus effort assessment. Tickets 1, 2, and 3 each took an opus/max shadow author beside the canonical opus/xhigh author. Each shadow started from the canonical author's source tip in its own worktree and never landed. The reviewer gave advance approval for every repair round and fence expansion and deferred every question to an opus/xhigh consultant; no consultation was needed.

## Gate-stage timings

- landing: commit 85f2b05ca18da90ed7c2f7033fb8bc51a2783de0, trace 3b472a4b746a8e2e4a68fc099ca7dfbd
- gofmt: 179 ms
- vet: 1333 ms
- test: 152592 ms
- race: 7393 ms
- system: 55323 ms
- shellcheck: 536 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| line declaration | expected 2 repair rounds | refuted | 6 | 4 repair cycles by the repair-attribution table | fable / high / orchestrator |
| ticket 1 | every coverage row red then green | held | 9 | coordinator swap probe bit at dfacb0be | opus / xhigh / implementer |
| ticket 2 | every coverage row red then green | held | 9 | coordinator swap probe bit at 7fa14194 | opus / xhigh / implementer |
| ticket 3 | every coverage row red then green | held | 9 | coordinator swap probe bit at 1a3e8992 | opus / xhigh / implementer |
| ticket 4 | the two sentences and every anchor intact | held | 9 | author anchor probe bit at eb7c791c; Spec axis read MP55 met | opus / xhigh / implementer |
| ticket 5 | every coverage row red then green | held | 9 | coordinator omit probe bit at b82ba53f | opus / xhigh / implementer |

Brier mean: 0.07 over 6 pairs; 0 abstained.

Every author stayed inside its fence and committed on a lane pass in one attempt. Every ticket matched its chunk table row, and no ticket needed a fence expansion. Each author observed its rows red before the production edit, and each named probe bit with a restored subject.

The assessment gave this cost for each first-pass commit. Tokens are the harness-reported subagent count.

| ticket | effort | tokens | harness time | outcome |
|---|---|---|---|---|
| 1 | xhigh | 210k | 1275 s | landed; review found one duplicated fixture helper and five missing edge rows |
| 1 | max | 301k | 1768 s | not landed; composed the existing fixtures and found the annotated-tag edge |
| 2 | xhigh | 283k | 1374 s | landed; review found one missing glob-shaped row |
| 2 | max | 399k | 3111 s | not landed; the same design and the same MP42 wording flag |
| 3 | xhigh | 257k | 1119 s | landed; review found comment form and a duplicate helper |
| 3 | max | 369k | 1742 s | not landed; the same design with two extra recheck witnesses |

Max cost 1.41 to 1.44 times the tokens of xhigh on every ticket and gave the same design each time. The one extra edge that max found, the annotated-tag merge head, the MP1 Coverage axis also found. No review graded a max tree, so the comparison is descriptive.

## Coordinator catches

- Two coordinator probes came back `invalid`, because a swap left a variable unused; the retried swap at the variable's source line bit.
- The record's C2 excerpt claimed an MP42 probe that no run retained. The orchestrator reran the spelling swap and the prefix drop at the frozen tip and retained both blocks.
- The orchestrator's prefix-drop probe stayed silent on MP2; the Coverage axis graded it as C21, and a ticket 2 repair added the glob-shaped row.
- Each refreeze made every sibling ticket's verification stale. The author sessions of tickets 2, 3, and 4 were resumed five times in total to re-record at the new source.
- The ticket 5 freeze refused with a stale plan digest, because the seam-cell rewrite moved the digest. The author ran the amendment first, as a prior build did.
- A review entry with only `no-op` findings records as a fail, so each such axis needed a confirming entry before the checkpoint.
- The MP3 last-head-only candidate graded as optional advice under the hardening-cycle rule, and the Coverage axis accepted that label.
- Four repair sessions each closed every finding in one attempt, and every independent probe of the confirming rounds bit.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-read-the-pending-merge-through-one-owner.md | 1 | spec-row |
| 2-refuse-an-incomplete-merge-name-set.md | 1 | spec-row |
| 3-record-each-merge-head-and-end-the-merge.md | 1 | one-source |
| 4-document-the-merge-contract-of-bench-commit.md | 0 | none |
| 5-grade-a-pending-merge-by-the-branch-change.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Let bench record verification import a bench test or bench probe spill directly, because every author wrote 6 to 7 excerpt files by hand.
  Feeds: FT405
- Let bench probe accept repeated --package, and state the multi-line --swap form in its help.
  Feeds: new
- Print every failure line in bench test when the lines cell is above one, or point to --full.
  Feeds: new
- Let bench record review carry a finding with a no-op disposition as a pass, so a no-op-only round needs no confirming entry.
  Feeds: new
- Let one refreeze re-stamp or re-run the sibling verification entries of a chunk, because five re-record passes followed three refreezes.
  Feeds: FT405
- The landing census counted 95 raw calls: rg 87, python3 4, ls 2, sed 1, wc 1. A symbol search has bench consumers, and the rest have no Bench form.
  Feeds: new

### Skills

- Make the implement-spec author charge name bench record amendment before the freeze of a last-chunk ticket that rewrites seam cells.
  Feeds: new
- Make the author charge ask one witness for each load-bearing pathspec magic, with a fixture whose glob-shaped name also matches a sibling.
  Feeds: none
- Widen the known-issues probe row to an unused variable, because two orchestrator swaps and one author omit came back invalid for that cause.
  Feeds: none

### Process

- Keep each author session open until its chunk closes, because the resumed sessions re-recorded at the chunk tip in 2 to 4 minutes each.
  Feeds: none
- Keep opus/xhigh as the code-ticket line for an owner-seam ticket, because opus/max gave the same design at 1.4 times the tokens.
  Feeds: none
- Record an orchestrator probe run in the review record at the moment it runs, because a prose claim without the block costs a reaffirmation pass.
  Feeds: none