## Outcome

The seal-policy build (FT482) landed at a8c7cd0fe on a green whole gate. A cache release now grades every user before it seals, and it seals and retires only a safe generation. On Linux, a process whose tasks are all zombies counts as absent, and a group with only such members counts as absent. A Markdown-only commit or merge in a kit worktree starts no background lane pass. The FT482 known-issues row is retired, and ADR 0030 states the release policy and the zombie rule.

The reviewer signed off the spec at the session start, and the spec stage landed at 0a142aeac. The build ran three serial chunks on Opus authors and three Sonnet/high review axes for each round.

## Cost

- landing: a8c7cd0fef3f0443c0501fc88e30443bea9b1047, range a8c7cd0fef3f0443c0501fc88e30443bea9b1047^1..a8c7cd0fef3f0443c0501fc88e30443bea9b1047^2
- wall: 162 min
- commits: 32 total, 26 bookkeeping, 6 change
- graph: 3 tickets, chain 2, 3 chunks, 2 serial checkpoints
- dispatch: 3 authors, 3 repairs
- trace: unknown
- wall per ticket: 54.0 min against 60 min
- bookkeeping per ticket: 8.7 against 7.1, over budget

Delegate sessions: 2 staleness slices, 3 authors, 3 repairs, 15 review axes in 5 rounds. Provider usage and charges are unknown.

## Gate-stage timings

- landing: commit a8c7cd0fef3f0443c0501fc88e30443bea9b1047, trace 4a099c307af630a47224da8fb143feb0
- SP-C1 chunk checkpoint: green, system phase 49.8 s.
- SP-C2 chunk checkpoint: green, system phase 50.9 s.
- Complete checkpoint: red twice on FT372 landing tests with evidence absent, then green with the system phase at 48.4 s.
- Landing: the gate reused the fresh complete-checkpoint verdict.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| build line | Each chunk takes at most one repair round. | held: each chunk took one round | 6 | true | opus / high / orchestrator |
| ticket 1 | The author lands SP1 to SP10 and SP34 in one attempt. | held, with two Standards findings | unknown | abstained | opus / high / implementer |
| ticket 2 | The author lands SP11 to SP20, SP33, and SP35 in one attempt. | held after a plan probe amendment | unknown | abstained | opus / high / implementer |
| ticket 3 | The author lands SP21 to SP32 and SP36 in one attempt. | held after the SP23 seam moved to a system test | unknown | abstained | opus / medium / implementer |

Brier mean: 0.16 over 1 pair; 3 abstained.

The ticket 1 author landed in one attempt. The review found a duplicated focused-run harness in the SP10 system test and a Reopen comment that told a change history.

The ticket 2 author stopped correctly when the `known-issues-row` plan probe was silent, because no `./internal/anchors` test reads the real file. It kept an empty task or member set present, which is the fail-closed reading, and the Spec axis accepted it.

The ticket 3 author stopped correctly when the BG6 merge fixture could not run the kit lane for a Markdown-only merge. It also found that the BG7 to BG10 no-grade helper can miss a pass, and it wrote positive controls.

## Coordinator catches

- The SP-C3 repair session ran a substitute mutation for the `rename-detection` plan probe. The coordinator ran the named mutation, and the repair session then recorded it as the owed performer.
- The coordinator recorded the named probe under its own identity first, so the completion entry refused it. The check reads only the owed performer of the last entry.
- The SP-C1 repair bundled a comment fix with a test refactor, and the SP-C2 confirmation carried a spec prose correction in the chunk delta. Each voided a test-only confirmation and cost one confirming round.
- A full system run with an uncommitted record file failed TestLifecyclePublicationCompatibility, and a clean-tree rerun was green.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-grade-before-the-cache-release-seals.md | 1 | one-source |
| 2-count-a-zombie-user-as-dead.md | 1 | delegate-error |
| 3-start-no-grade-for-a-markdown-only-commit.md | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Make `bench record verification` refuse a chunk entry whose performer is not the owed performer of its requirement, at write time.
  Feeds: new
- Make `bench spec widen` mark an absent path `(new)`, and print the kit-pin sentence that a system-tagged path needs.
  Feeds: new
- Census account: 88 raw calls, rg 87 and cat 1. The rg calls are sanctioned reads of response spills, cited lines, and gate logs, and the cat call read the handoff, which `bench read` serves. Count rg apart from other heads in the census close, so the one missed `bench read` form shows.
  Feeds: new

### Skills

- Make ticket slicing in `bench-write-spec` run each plan probe command against its mutation, and confirm that each named fixture reaches the row's branch.
  Feeds: new

### Process

- Make the build entry of `bench-implement-spec` require the scorecard read before the line declaration, because this build repeated three recorded decisions.
  Feeds: new
- Commit a comment-only Go correction as the orchestrator's evidence-only commit, and charge the repair session with test changes only.
  Feeds: none