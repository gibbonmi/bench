## Outcome

FT451 landed at `82f6df3b7`. A lane pass now detaches a background grade, and a chunk checkpoint or a landing reuses its green. The BG-C6 and BG-C7 checkpoints reused a background green with no oracle run. The census ratio was 1.258 against the 1.2 bound, and the reviewer waived the bound.

The build cost too much. It ran from 2026-10-27 15:23 to 2026-10-28 07:20, about 16 hours of wall time. 71 delegate sessions made 4,540 tool calls in 12.3 hours, at about 10 seconds for each call. About 104 of the 135 build commits were plan, record, or review commits.

| role | sessions | tool calls | time |
|---|---|---|---|
| ticket authors | 23 | about 2,500 | about 7.3 h |
| repair sessions | 11 | about 1,100 | about 3.4 h |
| reviews, confirmations, and consultations | 33 | about 800 | about 1.3 h |

Lane passes took 25 minutes in total, and 11 whole-gate runs took 59 minutes. Most time went to model turns, not to commands. No two authors overlapped, although the plan allowed three.

## Gate-stage timings

- landing: commit 82f6df3b7454ff8107edeaf6ac7243736384f59a, trace 69949aab55573c12ae5ff96d02f5a86f
- gofmt: 258 ms
- vet: 2396 ms
- test: 323238 ms
- race: 2954 ms
- system: 61819 ms
- shellcheck: 631 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| unknown | unknown | unknown | unknown | unknown | unknown |

Brier mean: unknown. Pairs: 0. Abstentions: 0.

The opus authors at high effort held the ticket boundaries in BG-C0 to BG-C3. The opus authors at medium effort held them in BG-C5 and BG-C7 to BG-C8. Tickets 06, 08, 09, and 10 each took about 40 minutes and 150 to 240 calls. A trivial repair still took 38 to 81 calls, because each writer session repeats the same record ceremony. The sonnet reviews at high effort took 1 to 6 minutes, and each caught a real defect.

## Coordinator catches

- The BG-C8 review and the live run showed that a stopped grade recorded `aborted`, and that each new grade erased the log of the earlier grade.
- The final `bench gate` went red because the checkout guard counted the files that a detached grade writes.
- The F1 confirming round found that the repair had deleted 15 checkout guard tests. The orchestrator had marked that file new in the widen.
- The first landing went red because the decision map named the closing roadmap row as a source.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-measure-the-background-census.md | 1 | other |
| 02-derive-the-run-binary-source-once.md | 1 | delegate-error |
| 03-read-the-baseline-members-from-commits.md | 0 | none |
| 04-add-the-background-width-input.md | 0 | none |
| 05-record-background-claims-and-outcomes.md | 1 | delegate-error |
| 06-grade-a-chunk-subject-in-the-background.md | 1 | one-source |
| 07-complete-the-background-run-outputs.md | 1 | edge-found |
| 08-start-a-background-grade-after-a-lane-pass.md | 0 | none |
| 09-queue-and-supersede-background-grades.md | 1 | edge-found |
| 10-let-a-chunk-checkpoint-wait-on-a-background-grade.md | 0 | none |
| 11-grade-a-landing-subject-in-the-background.md | 0 | none |
| 12-let-a-landing-wait-on-a-background-grade.md | 1 | delegate-error |
| 13-stop-background-grades-when-work-goes-away.md | 4 | edge-found, delegate-error, edge-found, other |
| 14-show-background-outcomes.md | 0 | none |
| 15-retire-the-prepared-run.md | 0 | none |
| 16-record-the-background-grade-rules.md | 0 | none |

## Agent-experience improvements

### Bench CLI

- Census account: 32 `rg` calls were sanctioned reads, so propose no form for them; 9 `cd`, 2 `diff`, and 1 other call were unknown.
  Feeds: none
- Add a verb that reads and edits the plan fence and the coverage map, because 6 `python3` calls edited or parsed that Bench-owned state.
  Feeds: new
- Make `bench spec widen` refuse the `(new)` label on a path that git tracks.
  Feeds: new
- Make `bench retro --scaffold` derive a cost section from git and the trace store: wall time, sessions, tool calls, and commits by kind.
  Feeds: new
- Let `bench assessment record` derive its attempts from the review record and the trace store, so the orchestrator types no attempt by hand.
  Feeds: new

### Skills

- Make `/bench-implement-spec` dispatch every eligible ticket up to the author limit, with one worktree for each concurrent chunk.
  Feeds: new
- Stop the recorded author runs of checks that the checkpoint gate runs again, and keep only the witness reds in the record.
  Feeds: new

### Process

- Replace the per-ticket assignment, widen, excerpt, and record steps with one record step and one record commit for each ticket.
  Feeds: FT414
- Give the drain a cost budget, and make each overrun propose one step removal before any new rule is admitted.
  Feeds: new
- Admit the bookkeeping cut as one outcome before FT402, as the 2026-10-28 opus sweep recommends.
  Feeds: new