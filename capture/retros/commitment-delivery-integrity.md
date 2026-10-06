## Outcome

The build delivers FT390 in two chunks and four tickets, and it landed at `97f96c9d` with a green gate. A plan now refuses a new or changed binding that names no obligation, and the completion landing refuses that binding before the gate. A plan binds only the sources that it keeps open, and one function orders the plan sources in byte order. All 26 acceptance rows reconcile as covered, and the landing recorded the FT390 delivery fact.

Each ticket had a fresh Opus author at high effort. Sonnet at high effort ran every review axis. Each chunk used one repair cycle of two, and each repair added tests only.

## Gate-stage timings

- landing: commit 97f96c9d0d3de74994217d1ef682163445b69500, trace 82a8dccb499318a91bfc7abff56928e5
- gofmt: 125 ms
- vet: 1266 ms
- test: 243191 ms
- race: 2902 ms
- system: 75088 ms
- shellcheck: 484 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 | FD1-FD5, FD9, FD23, FD24, FD26 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 2 | FD6-FD8, FD10, FD11 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 3 | FD12-FD19, FD25 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 4 | FD20-FD22 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| FD-C1 Coverage | four retention and scan edges untested | verified | 8 | held | claude-sonnet-5-5 / high / review |
| FD-C2 Coverage | two comparator terms untested | verified | 5 | held | claude-sonnet-5-5 / high / review |

Brier mean: 0.055 over 6 pairs. Abstentions: 0.

## Coordinator catches

- Ticket 2 stopped at its fence before a test-count pin bump. A plan commit added the pin file to the fence, and the same author finished.
- The orchestrator froze FD-C1 again before it recorded the plan amendment, so the record could not map the chunk. The scorecard already held the rule for that order. The record verbs restored the entry.
- Two authors recorded the exit of the probe verb, not the exit of the mutated test run. The checkpoint refused, and the authors recorded the entries again.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-refuse-obligation-free-plan.md | 1 | spec-row |
| 2-refuse-obligation-free-delivery.md | 0 | none |
| 3-bind-only-open-plan-sources.md | 0 | none |
| 4-order-plan-sources-canonically.md | 1 | spec-row |

## Agent-experience improvements

### Bench CLI

- Add a verb that edits the plan execution block, because the census counted 37 raw calls and seven were python3 plan edits.
  Feeds: new
- Make the commitment outlook name the plan step with the path, identity, and obligation when an eligible outcome binds no deliverable.
  Feeds: new
- Let `bench commitment start` bind against the approved policy on the branch, so a binding change needs no separate landing and second worktree.
  Feeds: new
- Make `bench preflight build` run the plan validation that the ticket charge runs, so a repeated verification ID turns the preflight red.
  Feeds: new
- Make `bench record chunk` refuse a plan digest change until `bench record amendment` maps it.
  Feeds: new
- Print the exit of the mutated test run in `bench probe`, and print a test count for each package in `bench test`.
  Feeds: new

### Skills

- Give each term of a comparator that a spec names, and its compare mode, a coverage row of its own.
  Feeds: new

### Process

- List the count-pin file in the `Writes:` line of a ticket that adds a top-level test to a pinned package.
  Feeds: new