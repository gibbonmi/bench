## Outcome

The fresh-ticket-authors spec landed at `c58e5aa5` from the reviewed source `304190cd` over the base `d23694e9`. Every spec-backed build now gives each ticket a fresh author session on the declared line. The version 2 delegate plan with an author limit of 1 records the authors. A post-review repair goes to a fresh repair session, and the orchestrator reconciles the final source. ADR 0023 records the decision, and ADR 0021 points to it.

The build used its own new shape. Three fresh Opus/high authors wrote the three tickets in `Blocked by:` order, and four fresh repair sessions wrote two repair cycles. Three review rounds of three fresh Opus/high axes graded chunk FA. The orchestrator wrote only plan, spec, and record commits.

## Gate-stage timings

- landing: commit c58e5aa5739405418cb259ce13b97bda2d410d0d, trace 9be80c530443ed757db0977f29d05505, census 7
- gofmt: 140 ms
- vet: 1169 ms
- test: 117157 ms
- race: 2912 ms
- system: 37971 ms
- shellcheck: 539 ms

## Ticket-versus-spec-slice and delegate performance

The build kept the one approved chunk. The first review added rows FA27 to FA30 and widened the ticket 3 fence to two out-of-fence paths by reviewer decision. The confirming round added FA31 by reviewer decision. Each ticket author worked inside its fence and committed on a green lane.

The harness reported these figures for each subagent. The token figure is the harness total for each subagent, and its cumulative or final-context meaning is not documented here.

| session | role | tokens | tool calls | minutes |
|---|---|---|---|---|
| ticket 1 author | implementer | 222k | 112 | 15.2 |
| ticket 2 author | implementer | 164k | 87 | 9.4 |
| ticket 3 author | implementer | 147k | 118 | 9.4 |
| round 1 axes (3) | reviewer | 312k | 79 | 3.5 parallel |
| cycle 1 repairs (3) | implementer | 226k | 158 | 15.7 serial |
| round 2 axes (3) | reviewer | 185k | 50 | 1.8 parallel |
| cycle 2 repair | implementer | 48k | 39 | 2.8 |
| round 3 axes (3) | reviewer | 102k | 28 | 0.7 parallel |

The build took about 92 minutes from the plan amendment to the final record. That is about 31 minutes for each ticket, the same rate as FT71, and it includes reviewer question pauses. The orchestrator token count is unknown.

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| `.bench/BENCH-reference.md` | FA14 needle reds a version 1 amendment | verified | 9 | held | Opus / high / implementer |
| `.agents/commands/bench-implement-spec.md` | FA10 needle reds a dropped handoff refresh | verified | 9 | refuted | Opus / high / implementer |
| `.agents/skills/bench-craft-tickets/SKILL.md` | FA18 needle reds a retained-session size | verified | 9 | held | Opus / high / implementer |
| `docs/field-guide.html` | FA22 page states the current authorship rule | verified | 9 | refuted | Opus / high / implementer |
| `.agents/commands/bench-implement-spec.md` | R5 wording agrees with the narrow author read | verified | 9 | held | Opus / high / implementer |
| `delegation-discipline.md` | FA30 needle reds the transfer entry | verified | 9 | held | Opus / high / implementer |
| `.agents/commands/bench-write-spec.md` | R10 scan reads the write-spec file again | verified | 9 | held | Opus / high / implementer |
| `docs/field-guide.html` | FA31 reds the retired sentence anywhere | verified | 9 | held | Opus / high / implementer |

Brier mean: 0.21 over 8 pairs. Abstentions: 0.

## Coordinator catches

- The ticket 3 return left a stale worktree binary, so the binary-seal preflight row went red until the orchestrator ran `bench worktree build`.
- A coordinator probe at the site the author already probed was vacuous, so the orchestrator ran a second swap at a new site.
- The ticket 1 and ticket 2 repair verification named earlier sources, so each repair session reran its verification at the final tip before the checkpoint.
- The fixture-closure preflight refused the ticket 3 fence until it named the three canary directories that pin `projects/benchkit.md`.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-route-each-ticket-to-a-fresh-author.md | 1 | delegate-error |
| 2-align-the-line-and-delegation-skills.md | 1 | spec-row |
| 3-align-the-phase-commands-and-docs.md | 2 | spec-row, delegate-error |

## Agent-experience improvements

### Bench CLI

- Let `bench anchors` take a directory, let `bench test` list every failure line, and add `bench probe --insert-before` and `--expect`, as the census entry "fresh-ticket-authors: census 7" proposes.
  Feeds: new
- Add a Bench verb that renders the review-record payload from verification and review returns, so the orchestrator writes no JSON builder script.
  Feeds: new
- Make the harness token figure for a subagent state whether it is cumulative, so a retro can compare it with a cached-read total.
  Feeds: none

### Skills

- Tell a ticket author to pair each rewritten sentence with a Forbid row for its retired text, because two rounds found unpinned retired sentences.
  Feeds: new
- State in the review phase whether an amended coverage row takes a separate repair ticket or joins its owner ticket's Covers line.
  Feeds: new

### Process

- Make the spec reader sweep grep every public doc for each retired sentence, because the sweep missed the field-guide callout, `projects/benchkit.md`, and `CONTEXT.md`.
  Feeds: new