## Outcome

The workflow-guidance-repair spec landed at `a00e593a` from the reviewed source `03e3a931` over the base `78602e89`. Eight tickets in chunks GR-A, GR-B, and GR-C repaired the 62 audit findings, the `bench commit` help strings, and the ADR 0014 correction. Each restated rule now points to its owner, and Forbid rows keep the retired copies out. The capture writer's tracked-or-ignored rule now sets when a retro commits.

Eight fresh Opus/high authors wrote the tickets in `Blocked by:` order, and nine fresh Opus/low repair sessions wrote five repair cycles. One cycle changed only ticket text. Eight review rounds of three fresh Opus/high axes graded the three chunks. The orchestrator wrote only plan, spec, and record commits. Two items stay open for the reviewer: the drain's "commit on green" wording and a stale field-guide sentence. A `bench idea` entry parks each one.

## Gate-stage timings

- landing: commit a00e593ab62dc1c3490a01a268df680f455ee925, trace e5354c11e29315b0b4bb3f688e17ce3d, census 15
- gofmt: 109 ms
- vet: 1081 ms
- test: 111413 ms
- race: 2952 ms
- system: 37145 ms
- shellcheck: 503 ms

## Ticket-versus-spec-slice and delegate performance

The build kept the three approved chunks and their ticket order. Three plan expansions changed the approved plan. New anchor rows moved from `registry_data.go` to `registry_retained_workflow.go`, because the first file is over its structure budget. The ticket 3 fence gained the convergence helper file. Ticket 7 took two GR-B findings on paths that it owns.

Each author worked inside its fence and committed on a lane pass. The ticket 3 author stopped once before any edit, because a helper outside its fence required three sentences that the ticket removes.

The harness reported these figures for each subagent. The token figure is the harness total for each subagent, and a resumed session reports a cumulative figure.

| session group | role | tokens | tool calls | minutes |
|---|---|---|---|---|
| ticket authors 1 to 8 | implementer | 1462k | 804 | 69.9 serial |
| GR-A repairs (4 sessions over 2 cycles) | implementer | 279k | 189 | 13.1 serial |
| GR-B repairs (2 sessions) | implementer | 152k | 103 | 7.1 serial |
| GR-C repairs (3 sessions) | implementer | 170k | 129 | 8.4 serial |
| GR-A axes, rounds 1 to 3 | reviewer | 634k | 192 | 7.2 parallel |
| GR-B axes, rounds 1 and 2 | reviewer | 445k | 149 | 5.9 parallel |
| GR-C axes, rounds 1 to 3 | reviewer | 558k | 191 | 6.8 parallel |

The build took 164 minutes from the plan amendment to the landing, which is about 20 minutes for each ticket. The orchestrator token count is unknown.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| `bench-craft-delegate/SKILL.md` | GR6 needle reds a dropped coverage-row duty | verified | 9 | held | Opus / high / implementer |
| `bench-review-implementation.md` | GR41 to GR44 rows keep the review copies out | verified | 9 | refuted | Opus / high / implementer |
| `internal/commit/dry_run_test.go` | GR53 test reds a dry-run line without the lane | verified | 9 | refuted | Opus / high / implementer |
| `bench-final-check.md` | GR66 and GR67 rows keep the scaffold copies out | verified | 9 | refuted | Opus / high / implementer |
| `bench-craft-line/SKILL.md` | GR75 needle reds a weakened step 5 | verified | 8 | refuted | Opus / high / implementer |
| `cmd/bench/anchor_help_test.go` | GR107 and GR108 tests red a wrong kind name | verified | 8 | refuted | Opus / high / implementer |
| `delegation-discipline.md` | GR17 to GR32 rows hold under review | verified | 9 | held | Opus / high / implementer |
| `bench-drain.md` | GR92 to GR105 rows hold under review | verified | 9 | held | Opus / high / implementer |
| `bench-craft-synthesis/SKILL.md` | GR45 row reds the commit-gate claim | verified | 9 | held | Opus / high / implementer |

Brier mean: 0.42 over 9 pairs. Abstentions: 0.

## Coordinator catches

- The ticket 1 lane refused growth of `registry_data.go`, so the orchestrator amended the Testing decisions and routed new rows to another registry file.
- The ticket 3 author found a helper outside its fence that required removed sentences. The orchestrator expanded the fence before any edit.
- The ticket 4 return left a stale worktree binary, so the binary-seal preflight row went red until the orchestrator ran `bench worktree build`.
- One delegate put backticks inside double quotes, and the shell started a `bench gate` that refused at once, so later charges required single quotes.
- One delegate left an `rg` call without a path waiting on its input, so the orchestrator stopped the task and later charges required a path.
- The orchestrator promoted two advice items to finding G11, then withdrew it when existing rows, one of them an older acceptance anchor, pinned that wording.
- The confirming Spec axis found that the orchestrator's own plan-expansion box for ticket 7 contradicted the G5 repair, so a text-only cycle amended the box.
- Each plan change after a chunk tip needed a plan-digest amendment entry, so the record carries a three-step amendment chain.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-align-the-delegation-skill-with-fresh-authors.md | 1 | delegate-error |
| 2-align-the-delegation-discipline-with-fresh-authors.md | 1 | spec-row |
| 3-route-phase-command-repairs-to-fresh-sessions.md | 2 | spec-row, delegate-error |
| 4-state-the-lane-and-landing-gate-once.md | 1 | spec-row |
| 5-capture-the-retro-by-the-tracked-rule.md | 1 | spec-row |
| 6-bind-each-ticket-author-to-the-declared-line.md | 1 | spec-row |
| 7-point-the-guides-at-each-fact-owner.md | 2 | spec-row, other |
| 8-correct-the-drain-and-craft-skill-references.md | 1 | delegate-error |

## Agent-experience improvements

### Bench CLI

- Make `bench anchors` read the registry of the worktree that `bench worktree exec` names, as the census entry "workflow-guidance-repair landing census: 15 raw calls" proposes.
  Feeds: new
- Add a Bench verb that renders the review-record payload and its plan amendments from the returns, so that the orchestrator writes no JSON script.
  Feeds: new

### Skills

- Tell a spec author to check the structure budget of each registry file that a coverage row names before the author places new rows there.
  Feeds: new
- Tell a spec author to plan a Forbid row for the owner sentence at each former copy site through the owner's named constant.
  Feeds: new

### Process

- Make the spec reader sweep search each removed sentence in normalized form across the conformance helpers and every phase command.
  Feeds: new
- Record every ticket assignment in one plan commit before the first chunk freeze, so the plan digest changes only for repair assignments.
  Feeds: none