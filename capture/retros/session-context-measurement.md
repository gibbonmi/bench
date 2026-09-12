# Retro: session-context-measurement

## Outcome

The spec landed at `e577ca53cc9ed20dca4860819039b0b9330063a7` on 2026-09-12 with `Status: implemented`.
`bench harnesses <harness> --record <path> --format codex-rollout-2026-09-11` reads one pinned Codex rollout and reports fourteen typed observations.
Each observation carries a value, a unit, an availability, a source, and a boundary sentence.
Both compiled views stay byte-identical to the baseline.
The budget-evidence report records 20 representative cases and proposes per-surface byte bounds; every proposal awaits reviewer approval.
The review record holds two chunks, eight plus nine review occurrences, the completed reconciliation of 25 rows, and the final verification.

## Gate-stage timings

- landing: commit e577ca53cc9ed20dca4860819039b0b9330063a7, trace 964d17da8d340aad4a118981c62091f3
- gofmt: 162 ms
- vet: 1706 ms
- test: 165053 ms
- race: 4128 ms
- system: 56978 ms
- shellcheck: 1468 ms

## Ticket-versus-spec-slice and delegate performance

One Opus/high delegate retained authorship through both tickets, three repair rounds, and the final reconciliation.
Ticket 1 delivered 21 rows in one pass. Sixteen rows went red against a skeleton reader, and five were proven red-capable by a named probe.
The first ME-C1 review returned 4, 4, and 4 findings and seven repair targets; the repair landed in one commit and every reaffirmation passed.

Ticket 2 delivered the report in one pass with 19 cases. The first ME-C2 review returned 3, 2, and 5 findings and eight repair targets, then one path cell.
The chunk boundaries held: each ticket was one chunk, one commit, one review, and one checkpoint.
Three Opus/medium axes reviewed each chunk in separate venues and converged in two rounds per chunk.
A Sonnet/xhigh delegate ran the debug path on the learnings inbox and falsified the two-parser hypothesis before it fixed the writer.

## Coordinator catches

- The staged spec carried no completion plan and no dogfood record; the coordinator added both under the plan-expansion policy before the first checkpoint.
- The author claimed per-dimension source-absence tests that did not exist; SPEC-2 exposed the gap and the repair added them.
- The coordinator's omission probes bit at two sites the author never probed: the orphan-output half of the unmatched census and the skipped-line mark.
- The author's case 2 capture used `git show --format=`, which two axes measured 825 bytes short; the repair named one command per cell.
- A mid-build main merge moved the landing base; the first landing refused, and the second landed on the merged tip.
- The sibling-merge refusal route would have squashed a hand-resolved merge; the landing's `git merge --continue` route kept both parents.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-inspect-record.md | 1 | delegate-error |
| 2-compare-budget-cases.md | 2 | delegate-error, delegate-error |

## Agent-experience improvements

### Bench CLI

- Add a review-record verb that writes the fenced record from a JSON file, per the census entry of 8 raw calls.
  Feeds: new
- Make `bench worktree merge`'s conflict refusal route to `git merge --continue`, the route the landing refusal already names, so a hand-resolved merge keeps both parents.
  Feeds: FT208
- Hash only the plan fence and the ticket bytes for the plan digest, so a prose amendment does not force a record amendment.
  Feeds: new
- Give `bench test --check` a cell that proves the selected check ran, and give `bench probe` a `--swap-file` pair for a multi-line site.
  Feeds: new

### Skills

- Make `craft-spec` require the completion-plan fence beside the chunk table, so a staged spec cannot reach its first checkpoint without one.
  Feeds: new
- Make `craft-gate` require that a new gate obligation lands with its entry check, its reference sentence, and its refusal route.
  Feeds: FT200

### Process

- Write and commit the review pickup before the repair charge goes out; this build sent one repair first.
  Feeds: none
- Merge the latest main tip into the source before the last chunk's review, and pass that tip as the landing base.
  Feeds: none
- Move the repair fence to the chunk, the union of its tickets' writes plus the paths the review names.
  Feeds: new
