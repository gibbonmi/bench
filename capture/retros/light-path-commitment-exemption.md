## Outcome

The build delivers FT391 in two chunks and five tickets, and it landed at `63680c21` with a green gate. An unbound commit or landing is now admitted when one tickets-only folder holds one ticket whose `Writes:` line covers every production path. Commit mode reads every tickets-only folder, and publication mode reads only the path that `--spec` names. One owner in `internal/tickets` holds the `Writes:` grammar. The guide, the drain, the implementation command, and ADR 0028 state the light-path exemption and the drain delegate route. All 56 acceptance rows reconcile as covered, and the landing closed FT391.

Each ticket had a fresh Opus author at high effort. The reviewer directed Fable at high effort for every review axis. LP-C1 used both repair cycles, and LP-C2 used one repair cycle of two.

## Gate-stage timings

- landing: commit 63680c21cacf58f5575d6d48f6f05df2ca0472d6, trace 08c1c83ad51fefbf988fe6a3450440a1
- gofmt: 164 ms
- vet: 1437 ms
- test: 266089 ms
- race: 2976 ms
- system: 81242 ms
- shellcheck: 566 ms

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| ticket 1 | LP27-LP33 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 2 | LP1, LP2, LP7-LP23, LP49, LP51, LP56 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 3 | LP3-LP6, LP24-LP26, LP50, LP52, LP53 meet the spec | verified | 9 | missed | claude-opus-5-5 / high / implementation |
| ticket 4 | LP34-LP46, LP54, LP55 red then green | verified | 9 | held | claude-opus-5-5 / high / implementation |
| ticket 5 | LP47, LP48 stated in the ADRs | verified | 9 | held | claude-opus-5-5 / high / implementation |
| LP-C1 Coverage C1 | the `TicketsOnly` clause is untested | verified | 9 | held | claude-fable-5-1 / high / review |
| LP-C1 Coverage C2 | the publication folder scope is untested | verified | 8 | held | claude-fable-5-1 / high / review |
| LP-C1 Coverage C3 | the first-ticket operand is untested | verified | 8 | held | claude-fable-5-1 / high / review |
| LP-C1 Coverage C4 | an executable ticket refuses itself | verified | 6 | held | claude-fable-5-1 / high / review |
| LP-C1 Standards S1 | the ticket body template has two copies | verified | 8 | held | claude-fable-5-1 / high / review |
| LP-C1 Spec P1 | publication mode reads every folder | verified | 6 | held | claude-fable-5-1 / high / review |
| LP-C1 Coverage C5 | `ls-tree` matches a glob without the literal flag | verified | 8 | missed | claude-fable-5-1 / high / review |
| LP-C1 Standards S3 | the specs root constant has two copies | verified | 7 | held | claude-fable-5-1 / high / review |
| LP-C1 Spec P3 | spec line 93 contradicts the fence | verified | 9 | held | claude-fable-5-1 / high / review |
| LP-C2 Standards S5 | ADR 0028 copies the tier binding | verified | 5 | held | claude-fable-5-1 / high / review |
| LP-C2 Standards S6 | the recurrence diagnostic needs new words | verified | 7 | missed | claude-fable-5-1 / high / review |
| LP-C2 Standards S7 | the registry comment omits restrictions | verified | 6 | held | claude-fable-5-1 / high / review |
| LP-C2 Standards S9 | the registry comment omits grant holders | verified | 6 | held | claude-fable-5-1 / high / review |

Brier mean: 0.172 over 18 pairs. Abstentions: 0.

## Coordinator catches

- The Spec axis proposed `no-op` for P1. Spec line 99 is exact, so the coordinator routed P1 to a ticket 3 repair.
- The cycle 2 repair stopped on C5. The coordinator confirmed in a scratch tree that `ls-tree` turns off wildcard matching, and closed C5. The same probe found a closed-failure `git show` glob defect, which is parked as an idea.
- The ticket 4 repair stopped at its fence on S6. Spec line 292 pins that diagnostic, so the coordinator closed S6 against the acceptance row.
- The coordinator's own plan commit expanded the fence into `internal/spec` and left spec line 93 stale. The checkpoint refused the Spec axis, and the spec fix moved the frozen source. All 12 LP-C1 verifications ran again.
- The coordinator froze a chunk before it recorded the amendment, which the FD retro also found. The record verbs refused, and the coordinator restored the entry by hand.
- The ticket 1 author ran in a cleared session. A `session-lost` assignment gave ticket 1 a fresh verification session.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| 1-own-writes-grammar.md | 0 | none |
| 2-admit-light-path-commit.md | 1 | one-source |
| 3-admit-light-path-landing.md | 2 | delegate-error, one-source |
| 4-state-light-path-guidance.md | 1 | other |
| 5-record-commitment-scope-adr.md | 1 | one-source |

## Agent-experience improvements

### Bench CLI

- Add a `bench record assignment` verb that appends and validates one plan assignment, because the census counted 12 python3 calls and most edited the plan block.
  Feeds: new
- Make `bench record chunk` refuse a plan digest change until `bench record amendment` maps it, because the wrong order recurred from FD.
  Feeds: new
- Let the checkpoint accept a spec-prose change after the chunk tip as evidence-only, so a one-sentence fix does not force every ticket to verify again.
  Feeds: new

### Skills

- Make the Standards axis read the spec acceptance rows before it proposes new words for a pinned text.
  Feeds: new
- Make a review axis prove a claim about Git pathspec behavior in a scratch tree before it files the finding.
  Feeds: new

### Process

- When a plan commit expands a fence into a package, update each spec sentence that says that package needs no edit in the same commit.
  Feeds: new
- When a ticket adds a fixture helper, name each existing helper that writes the same artifact and put its path on the `Writes:` line.
  Feeds: new