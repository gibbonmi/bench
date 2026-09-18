## Outcome

The landing replaced the full charge response with bounded, content-addressed charge evidence.
A preparation command now publishes one evidence artifact and returns its identity.
A consumer reads the manifest and each source in pages of at most 8,192 bytes, verifies them, and binds them to the current assignment.
The landing also added explicit cleanup with an exact fingerprinted plan, the review preparation route, and the guidance for retained reuse and native handoff.
It retired the full charge route of the build phase and of the review phase.

Two transport records prove the delivery in the Codex harness and in the Claude harness.
Each consumer read 23 production pages through its own shell tool and rebuilt six sources with its own decoder.
The landing published commit `2e30834b78eb828d19a5fca6bd41450ab1489643` from reviewed source `16e6584b2bcdfec0732273c7273a3d15ff618361` on base `27d3a8eac8d8cb334d691a6815199d92f80e562b`.
Rows CE94 and CE117 reconcile as covered by a reviewer decision of 2026-09-18.
Row CE94 still waits for a privileged native run of its device case, and row CE117 still waits for macOS storage evidence.
The retirement removed the review pickup, so this retro is the record of those two pending notes.

## Gate-stage timings

- landing: commit `2e30834b78eb828d19a5fca6bd41450ab1489643`
- gofmt: 117 ms
- vet: 1168 ms
- test: 115055 ms
- race: 2918 ms
- system: 36433 ms
- shellcheck: 548 ms

## Ticket-versus-spec-slice and delegate performance

Seven serial tickets landed in seven chunks, and each chunk held one ticket.
Every chunk needed at least one repair cycle, and chunk CE-C1D needed four under two reviewer extensions.
The ticket authors and the repair writers ran Opus at high and at medium effort, and every review axis ran Opus at medium effort.

The coordinator changed between sessions. The first sessions ran on Opus and forked their ticket authors.
The last session ran on Fable, so ticket 7 used a fresh Opus writer in place of a fork.
One Sonnet writer ran the first CE-C1A repair, because an Opus author session ended on an API error.
A second cycle closed what that repair left open.

The authors of tickets 5 and 6 each used their whole 200-turn budget, which shows two slices that were too large for one author session.
Three writers reported a red `fixture-closure` check that was green, because they ran the preflight against the primary checkout.

The ticket 7 author returned five rows, and all five held under the coordinator's probes and the Spec axis.
Its first return still carried a second phase inventory, which the Standards axis found.
The CE-C3 repair writer closed 13 targets with 16 recorded reds in one cycle. Its one wrong judgment left the worktree seal red.

Six Sonnet readers read digests of earlier sessions for this retro, and three of those digests held this build.
The earlier authors and axes stated many confidences, mostly 7 to 10, but the review pickup retained none of them.
This table therefore holds the claims of the last session and the three earlier claims that the readers found with a label.
The charges of the last session did not ask for a confidence, which is a coordinator omission.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| review finding | CE-C3-ST6 | verified | 9 | held | opus / medium / Standards |
| review finding | CE-C3-CV7 | verified | unknown | held | opus / medium / Coverage |
| review finding | CE-C3-CV8 | verified | unknown | held | opus / medium / Coverage |
| review finding | CE-C4-ST1 | verified | unknown | held | opus / medium / Standards |
| review finding | CE-C4-ST2 | verified | unknown | held | opus / medium / Standards |
| review finding | CE-C4-ST3 | verified | unknown | held | opus / medium / Standards |
| review finding | CE-C4-SP1 | verified | unknown | held | opus / medium / Spec |
| review finding | CE-C4-CV1 | verified | unknown | held | opus / medium / Coverage |
| review finding | SS-CV1 | verified | 7 | held | opus / high / Coverage |
| line declaration | ticket 3, expected repair rounds: 1 | verified | 6 | refuted | opus / unknown / orchestrator |
| line declaration | ticket 4, expected repair rounds: 1 | verified | 5 | refuted | opus / unknown / orchestrator |
| delegate return | CE-C1B rename repair, `fixture-closure` red | verified | 8 | refuted | opus / low / repair writer |
| delegate return | CE-C3 repair cycle 1, seal left red | verified | unknown | refuted | opus / medium / repair writer |
| delegate return | CE-C3 repair cycle 2 | verified | unknown | held | opus / medium / repair writer |
| delegate return | ticket 7, rows CE104 to CE113 | verified | unknown | held | opus / high / author |
| delegate return | CE-C4 repair cycle 1 | verified | unknown | held | opus / medium / author |

The Brier mean is 0.27 over 5 labeled pairs, and the abstention count is 11.
The two review findings give 0.05, and the three refuted claims give 0.417.

## Coordinator catches

The coordinator caught the red `binary-seal` check that the CE-C3 repair writer left, and it rebuilt the worktree binary with the sanctioned verb.
It caught a stale git-ignored broker manifest that fails a hand-run system check after any worktree rebuild.
Its first reading blamed a mutation probe, and a second occurrence on a clean tip corrected that reading.

The coordinator caught the CE-C3 checkpoint refusal. The chunk base was a plan commit after the CE-C2 close, so a fourth review round covered that delta under the corrected base.
It caught the completion checkpoint refusal. The reconciliation changed the spec after the last chunk review, so CE-C4 extended to the reconciled tip with one more round.
It refused a `bench doctor --fix` inside a worktree, because that verb also rewrites the user shim.

The coordinator reproduced each accepted finding with its own probe before it charged a repair.
It refuted one advice item of the CE-C3 Spec axis during the reconciliation. The cleanup tests that the axis reported absent exist in the store package.
It ran the first Codex transport decoder itself, and the Spec axis correctly returned that record.
The contract requires the consumer to rebuild the evidence.

## Repair attribution

The causes of the rounds before chunk CE-C3 come from the session text of the earlier coordinator, which the Sonnet readers quoted.

| ticket | rounds | causes |
|---|---|---|
| 1-validate-legacy-prepared-packs.md | 3 | spec-row; delegate-error; spec-row |
| 2-publish-bounded-build-evidence.md | 2 | spec-row; other |
| 3-verify-current-evidence.md | 2 | spec-row; spec-row |
| 4-activate-bounded-build-guidance.md | 4 | spec-row; other; other; other |
| 5-freeze-bounded-review-evidence.md | 2 | other; delegate-error |
| 6-clean-evidence-explicitly.md | 2 | delegate-error; delegate-error |
| 7-preserve-consumer-context.md | 1 | delegate-error |

Sixteen repair rounds landed seven tickets.
The `spec-row` rounds closed a coverage clause that the first return did not grade.
The `other` rounds closed duplicated knowledge and weak grading, and the fourth round of ticket 4 repaired checks that earlier rounds had added.
Both line declarations expected one round, and the tickets took two and four.

## Agent-experience improvements

### Bench CLI

- The `bounded-charge-evidence` census recorded 156 raw calls, led by `cp` 55 and `python3` 31, and its learning entry proposes two verb changes.
  Feeds: new
- Add a `bench review record` verb for a chunk entry, a round, and the completion entry, because each coordinator session scripted that fence.
  Feeds: new
- Make `bench worktree build` republish the git-ignored broker manifest beside the wrapper, because a hand-run system check fails after every rebuild.
  Feeds: new
- Make the chunk checkpoint name the expected base commit when it reports a chain gap.
  Feeds: new
- Make the completion checkpoint name the unreviewed commits when it reports a stale reviewed source.
  Feeds: new
- Export `BENCH_HOME` in the landing route when the wrapper can derive it, because a green landing gate refused as infrastructure.
  Feeds: new
- Let a sandboxed evidence reader take its shared lock with no write access to the store, because the Codex sandbox refused the read.
  Feeds: new

### Skills

- State in the implement-spec phase that the chunk base is the close commit of the predecessor chunk.
  Feeds: new
- State in the implement-spec phase that the final reconciliation commit joins the last chunk review, because the completion checkpoint requires that source.
  Feeds: new
- State in `craft-delegate` that every write charge and every axis charge asks for a stated confidence, because 11 of 13 labeled claims carry none.
  Feeds: new
- State in `craft-line` that a recorded fork line resolves to a fresh writer on the recorded model when the coordinator model differs.
  Feeds: new

- Give `craft-tickets` a size signal from the author turn budget, because the authors of tickets 5 and 6 each used all 200 turns.
  Feeds: new

### Process

- Retain each stated confidence in the review pickup at return time, because the earlier sessions stated dozens and the pickup kept none.
  Feeds: new
- Run every delegate preflight through the worktree, because three writers reported a stale red from the primary checkout.
  Feeds: new

- Carry the transport consumer protocol into the charge evidence reference, because two Codex runs failed on a hook and on a sandbox.
  Feeds: new
- Decide at spec approval how a row with evidence that needs an absent capability reconciles, because the completion checkpoint accepts only one value.
  Feeds: new
