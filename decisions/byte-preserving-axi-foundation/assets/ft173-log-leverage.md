# FT173 harness-log leverage ledger

Ledger ID: `ft173-log-leverage`

Subject: `64d85df9f3ff051fe4482934381a3cc4a5ca1489`

Review status: **approved.** This is the QD5 first-ticket
checkpoint. It blocks assignment of QD1--QD4 behavior tickets; it does not
approve an implementation or amend coverage. The approved L4 route is captured
separately for roadmap review.

## Corpus ledger

Selection method: start with the parent session IDs named by the FT173 shaping,
foundation, roadmap, forward-build, and coherent-diff records; do not include
child transcripts merely because their parent prompt is copied into them. Add
three later Bench sessions that respectively exercised an active AXI lifecycle,
its fence repair, and a committed guidance spec. Search each selected transcript
for (1) missing or wrong `help[]`, (2) repeated commands that form one query,
and (3) agent-side output shaping (`head`, `tail`, `awk`, `sort`, and equivalent
projections). The source directories and the selected IDs below are the complete
corpus, not popularity samples.

| id | harness | source directory | transcript count | session-id time range and exact IDs | purpose |
|---|---|---|---:|---|---|
| C1 | Codex | `/home/devuser/.codex/sessions/2026/08/09` and `/home/devuser/.codex/sessions/2026/08/10` | 7 (3,675 lines) | `2026-08-09T10:20:43`--`2026-08-10T12:54:31`; `019fe6e5-bc97-7862-af02-32de96b4d37f`, `019fe752-223b-7091-9677-45a2e475e14c`, `019fe790-2e4f-7821-b3e2-52f77af84f29`, `019fe7f1-e414-70d1-a910-92ca5f58ae87`, `019feb9b-a4e2-70d3-b33a-730e757d434a`, `019fec90-f106-7ab3-a84a-877917518df9`, `019fec98-e588-7161-941b-9498f455f9de` | FT173 parent runs: shaping, foundation, roadmap/log-review requirement, forward-build, and reject/accept coherent-diff review. |
| C2 | Claude | `/home/devuser/.claude/projects/-home-devuser-workspace-bench` | 1 (804 lines) | `2026-08-10T14:55:23`--session end; `13e4bde3-f5c0-4d3a-99be-265a58ce8cf9` | Fable restructuring handoff: it names contextual `help[]`, coherent `bench diff`, and the bounded AXI surface as the delivered behavior. |
| C3 | Codex | `/home/devuser/.codex/sessions/2026/08/11` | 3 (4,927 lines) | `2026-08-11T04:27:36`--`2026-08-11T18:34:05`; `019fefef-2b9e-7de0-b269-75f5b9945744`, `019ff095-6e78-7bc2-9057-cba0891a344e`, `019ff2f6-2646-7f83-9c57-141c82e3355c` | Representative recent Bench work: active-build orientation, fence diagnosis, and a full reviewed spec authoring path. |
| C4 | Corpus total | C1--C3 | 11 (9,406 lines) | `019fe6e5-bc97-7862-af02-32de96b4d37f`--`019ff2f6-2646-7f83-9c57-141c82e3355c` by selected time range | Complete selected corpus. |

| id | unavailable source | disposition |
|---|---|---|
| U1 | `/home/devuser/.claude/sessions/` contains no JSONL transcript archive; the Claude project archive in C2 is the available source. | named unavailable; not substituted. |
| U2 | No reviewer-supplied Codex usage-log manifest, session-ID list, or separate archive was available. | named unavailable; C1 and C3 are explicit local-session selection, not a substitute for an external manifest. |

## Opportunity ledger

The two current FT173 evidence assets establish the relevant ownership boundary:
the command-help inventory assigns useful next actions to the remaining query
surfaces and `bench diff`; the FT173 contract keeps `bench diff` as the sole
coherent Git inspection owner. A disposition is deliberately one value only.

| id | class | observed opportunity | evidence in the selected corpus | disposition |
|---|---|---|---|---|
| L1 | missing or wrong `help[]` | The query-oriented runs inspect `maps`, `guards`, `coverage`, and `worktree list`, then formulate the actionable follow-up themselves. The output has the fact but not the typed, state-carrying action. | C1's shaping and roadmap sessions name missing contextual actions; C2's Fable handoff retains contextual `help[]` for the approved query surface. | already-owned: QD1 of `specs/axi-query-disclosure/spec.md:55-59` and the per-command action inventory at `decisions/byte-preserving-axi-foundation/assets/ft173-command-help-inventory.md:71-75,105-110`. |
| L2 | missing or wrong `help[]` | Lifecycle sessions repeatedly reconstruct the next `bench spec build` operation from status, assignment, and candidate facts. | C1 and C3 include active-build and fence-repair sessions; the observed pattern is the feature's original high-frequency lifecycle case. | already-owned: `axi-spec-build-complete`, as fixed by `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md:223-228`; it is outside this ticket's QD1--QD4 fence. |
| L3 | repeated tool-call sequence replaceable by one query | Agents compose `git status --short`, diff statistics/whitespace checks, log inspection, and a patch view to orient a review subject. | Repeated in C1--C3; the selected FT173 review sessions show the same orientation sequence before assessing a candidate. | already-owned: `axi-coherent-diff`; `bench diff` is the sole coherent snapshot owner in `decisions/byte-preserving-axi-foundation/ft173-axi-contract.md:197-202`, with bounded and full views in `ft173-command-help-inventory.md:74`. |
| L4 | output shaping agents performed themselves | Agents reduce transcript and command output with `head`, `tail`, `awk`, and `sort` to obtain a bounded corpus or evidence projection. A query over harness transcripts would be a separate operational surface, not a projection of an approved QD query. | Present in C1--C3 while locating commands, IDs, counts, and focused evidence. | routed: reviewer approved `harness-transcript-query`, recorded in `capture/IDEAS.md` for roadmap review. |

Disposition counts: already-owned 3; routed 1; fold 0; decline 0.

## Completeness rule and sign-off

This ledger is complete only when C1--C4 and U1--U2 each occur exactly once,
L1--L4 each occur exactly once, and every L row has exactly one disposition from
`fold`, `already-owned`, `decline`, or `routed`. The L4 route is not a
coverage change. Were an observed opportunity to require a QD1--QD4 coverage
change, its only permitted disposition here would be **proposed amendment** for
reviewer approval; no spec edit would be made by this ticket.

Reviewer sign-off: **approved.** The reviewer approved the complete ledger and
the L4 `harness-transcript-query` roadmap route before QD1--QD4 assignment.
