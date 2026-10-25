## Outcome

FT445 landed as `5183804f2f29706779dbf1e5e3b307d22fc7bf99` from source tip `6889a3de4` on review base `59386368b`. The build ran 12 tickets in 10 chunks as a `--full --delegate --auto-approve` run with an author limit of 3. Ticket authors ran on Opus at high or xhigh effort. Review axes ran on Sonnet/high, and ticket 12 ran on Sonnet/high as a cheap-tier trial.

Each ticket's evidence now binds to its own commit, and a repair no longer restages. A test-only repair closes through `bench record confirmation`. The orchestrator prepares the checkpoint gate at a freeze or a repair commit, and the checkpoint reuses that green. A test-only delta after a whole green takes a delta checkpoint, and its green has its own record class.

The reviewer dropped ticket 13, the exec reuse of the worktree binary, after the ticket 12 measurement. The net saving was about 1.5 s for each repair. The reuse would also break `bench probe` under exec at its first mutant. The final verification passed 9 of 9 rows, and the completion reconciled 80 coverage rows.

## Gate-stage timings

- landing: commit 5183804f2f29706779dbf1e5e3b307d22fc7bf99, trace 6522cda3232a4e0d2c6997291b51541a
- gofmt: 196 ms
- vet: 1730 ms
- test: 207198 ms
- race: 2721 ms
- system: 47825 ms
- shellcheck: 620 ms

The FT362 landing gate took a test phase of 196 s and a system phase of 48 s, so FT445 adds about 11 s of unit tests. The RC-C6B checkpoint ran whole with `--fresh` in 5 min 40 s, by the bootstrap rule. RC81: one test-only repair cycle took 40 min 23 s, below the 45-to-60-minute baseline and above the 11-to-14-minute target. The repair author's own runs took about 26 min of that cycle.

## Ticket-versus-spec-slice and delegate performance

Each ticket fit its slice, and no author crossed a chunk boundary. The ticket 11 author session was lost after its commit, and a `session-lost` replacement verified the preserved source with no fix. The ticket 12 author measured the premise and found the probe hazard that dropped ticket 13. The table holds the RC-C6B and RC-C7 review claims, because earlier sessions kept no claim table.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| RC-C6B r1 Coverage | C1 a red delta can write a delta record with no red | finding | 9 | true | sonnet / high / reviewer |
| RC-C6B r1 Coverage | C2 the tree term of the delta reuse is unpinned | finding | 9 | true | sonnet / high / reviewer |
| RC-C6B r1 Coverage | C3 the Ready term of the delta reuse is unpinned | finding | 8 | true | sonnet / high / reviewer |
| RC-C6B r1 Coverage | C4 the stored run binary value is unpinned | finding | 7 | true | sonnet / high / reviewer |
| RC-C6B r1 Standards | the delta validator repeats the lane validator tail | finding | 6 | true | sonnet / high / reviewer |
| RC-C6B r1 Standards | the digest read relies on a commented invariant | finding | 4 | false | sonnet / high / reviewer |
| RC-C7 r1 Standards | S1 the spec still describes the dropped exec reuse | finding | 9 | true | sonnet / high / reviewer |
| RC-C7 r1 Standards | the asset still names ticket 13 as next | finding | 8 | true | sonnet / high / reviewer |
| RC-C7 r1 Spec | P1 three spec passages contradict the decision | finding | 8 | true | sonnet / high / reviewer |
| RC-C7 r1 Spec | the asset and ticket 12 treat ticket 13 as live | finding | 7 | true | sonnet / high / reviewer |
| RC-C7 r1 Spec | the empty source root claim is too broad | finding | 5 | true | sonnet / high / reviewer |
| RC-C7 r1 Spec | the asset's probe facts lack a run | finding | 6 | false | sonnet / high / reviewer |
| RC-C7 r1 Spec | the roadmap row still lists change 6 | finding | 6 | true | sonnet / high / reviewer |
| RC-C7 r1 Coverage | the asset states the pre-drop decision | finding | 7 | true | sonnet / high / reviewer |
| RC-C7 r1 Coverage | the roadmap row still lists change 6 | finding | 6 | true | sonnet / high / reviewer |

Brier mean 0.111 over 15 pairs; 0 abstentions.

## Coordinator catches

- The coordinator found the lost ticket 11 author session at resume, and it recorded a `session-lost` replacement before any record step.
- The coordinator's one probe on the ticket 11 tree found the unpinned tree term of the delta reuse before any review axis ran.
- The coordinator confirmed the ticket 12 hazard in `held_runs.go` and stopped for the reviewer, who dropped ticket 13.
- The coordinator caught that the RC-C7 r2 Coverage axis ran its evidence check outside the worktree, and it had the axis bind its result.
- Review caught the coordinator's own descope commit, which left five spec passages that described the dropped exec reuse.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 01-exempt-the-assignment-block.md | 1 | delegate-error |
| 02-grade-entries-at-their-own-commit.md | 1 | one-source |
| 03-record-entries-at-the-own-commit.md | 1 | delegate-error |
| 04-grade-the-floor-per-chunk.md | 1 | delegate-error |
| 05-one-plan-commit-and-preflight-per-dispatch.md | 0 | none |
| 06-refuse-a-dirty-chunk-checkout.md | 0 | none |
| 07-key-the-chunk-checkpoint-without-the-record.md | 1 | edge-found |
| 08-prepare-the-checkpoint-gate.md | 1 | spec-row |
| 09-list-the-system-suite-in-the-selector.md | 1 | delegate-error |
| 10-run-the-delta-checkpoint.md | 1 | one-source |
| 11-record-and-reuse-the-delta-green.md | 1 | delegate-error |
| 12-measure-the-run-binary-cost.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Add a read of one chunk of the version 2 plan, because 4 raw `python3` calls decoded the plan fence to read chunk rows and assignments (owned gap).
  Feeds: new
- Add a verb that drops a ticket from the plan fence, the coverage map, and the chunk table and lists its other references, because a hand descope missed five passages.
  Feeds: new
- Make `bench preflight evidence` name `bench worktree exec <label>` in its refusal from the primary checkout, because five review axes hit the ancestry refusal first.
  Feeds: new
- Make `bench record ticket` name the repair prefix form in its duplicate-ID refusal, because a refreeze author had to guess `t12r1`.
  Feeds: new
- Census account: 4 `python3` calls are an owned gap, the plan-fence reads above. 3 `rg` and 2 `ls` calls are sanctioned reads of the review record and the tickets. 3 `cat`, 2 `cp`, and 1 `sed` call came from delegates with unknown purpose.
  Feeds: none

### Skills

- Make the repair charge in implement-spec name only the plan's probes, and let the one coordinator probe carry the finding mutants, because the RC-C6B repair ran a duplicate mutant set.
  Feeds: new
- Make craft-spec require a measurement ticket to test the proposed change against each exec child that mutates source, such as `bench probe`, because that hazard surfaced only at the measurement.
  Feeds: new

### Process

- Measure RC81 again on the first build after FT445, because the bootstrap `--fresh` rule kept the delta checkpoint from serving this build's repair cycle.
  Feeds: new
- Prioritize the four gate-redesign ideas in the closing drain, by reviewer direction, because the author's own runs, not the gate, now dominate the repair cycle.
  Feeds: new