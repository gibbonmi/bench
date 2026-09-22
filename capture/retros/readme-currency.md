# README currency retrospective

## Outcome

The README now explains why AI-written code still needs maintainable seams, consumer inspection, and external verification.
It also corrects the current workflow, setup, installation, maintenance, example-profile, and shift-note guidance.
Four Mermaid diagrams show research and specification, chunk review and repair,
final reconciliation and landing, and capture as visible workflows.
The change modifies no runtime behavior.

## Gate-stage timings

At source close, no landing timings were available for this README work.
The scaffold's commit `5a7fb32f2d9d19cbd093f6f6ad386e93002c7884` and trace belong to the previous drain.
This retrospective does not attribute those timings to this change.

- `docs-currency-workflow`: 603 ms package; 3.33 s wall; no skips
- `load-validity-metadata`: 104 ms package; 1.79 s wall; no skips
- `prose-mechanics`: 339 ms package; 1.89 s wall; no skips

The coordinator will retain the whole-project gate and landing evidence after publication.

## Ticket-versus-spec-slice and delegate performance

| surface | claim | status | confidence | label | model / effort / role |
|---|---|---|---|---|---|
| `update-readme.md` | The README candidate matches the cited sources and passes the focused checks. | claimed | 9 | unknown | gpt-5.6-sol / high / implementation |

## Coordinator catches

The coordinator verified the command help, the consumer query, and the outline query.
The coordinator also caught two distribution details before author verification.
The packaged README needs a canonical release-status URL, and the durable install example must use `bench setup`.
The first Mermaid render caught a reserved node identifier before publication.
The coordinator accepted the revised semantics and rendered all four final diagrams.

## Repair attribution

| ticket | rounds | causes |
|---|---|---|
| `update-readme.md` | 1 | other |

## Agent-experience improvements

### Bench CLI

None. The read-only commands supplied bounded evidence for the README claims.
Feeds: none

### Skills

None. The seam, domain, ticket, ADR, and synthesis guidance covered the work.
Feeds: none

### Process

None. The light path matched this prose-only change.
No scorecard change is due because the run supplied no comparative routing evidence.
Feeds: none
