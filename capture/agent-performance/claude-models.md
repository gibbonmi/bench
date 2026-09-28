# Claude model scorecard

Last incorporated landing: `ref-inventory` (`ab886f5b5ee04867110e79e4e1879a920a50434f`, 2026-09-27).

The FT199 build ran in Claude Code with a Fable/high orchestrator. Six fresh Opus/high ticket authors, eight Opus/high repair sessions, 36 Opus/high review axes over 12 rounds, and one Fable/high consultation took part. Delegate usage was about 5.0 million tokens. Orchestrator usage and provider costs remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 37 landings + implementer, 10 charges + reviewer, 21 axes | On `ref-inventory` the high orchestrator ran four chunks through 12 review rounds and probed every accepted repair with an independent mutation. It moved the last chunk tip twice for late plan commits and let the first landing refuse on an unreviewed composition. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 30 axes + control, 1 full-retrieval run + decider, 12 consultations + debug, 2 charges | On `ref-inventory` the one consultation recommended a local no-deref delete with a new row and one more bounded cycle, and the reviewer granted it. On `review-evidence-file-pages` the axes found 18 findings over two chunks, and each confirming round held. | Review axes by reviewer direction from round two on; consultations whenever the coordinator is uncertain; the RE2-style full control | unknown |
| Opus / high | implementer, latest 10 fresh ticket sessions; reviewer, 36 axes; repair, 9 sessions; spec reviewer, latest 2 passes | On `ref-inventory` six authors landed first-pass with biting probes, and the 36 axes found 49 findings that all held; the Coverage axes found five symref defects across four rounds. Seven of eight repair sessions landed in one attempt, and one stopped honestly on a Git limit the row had not decided. | High for fresh ticket authors, first-round review axes, repairs at the ticket's effort, and foundational Go seams. | implementation: 0.100 over 66 pairs, 2 abstained; reviewer: 0.163 over 45 pairs, 4 abstained; repair: 0.010 over 2 pairs |
| Opus / medium, low | reviewer, orchestrator, and repair sessions combined | On `review-evidence-file-pages` eight low repair sessions closed each target, and every claimed probe bit under the coordinator's own probe. Two of those repairs introduced new duplication, and one stopped honestly at its fence when its charge named the wrong helper. | Medium for gates, conformance, canaries, and triage. Low for exact tickets; a post-review repair keeps the ticket's declared effort by the 2026-09-27 decision. | 0.109 over 31 pairs; 9 abstained |
| Sonnet / high | orchestrator, 3 landings + reviewer, 14 axes + implementer, 3 comparison tickets | On `ft336-bounded-output` the three comparison authors used about twice the calls and time of Opus, and a blinded grader preferred Opus on each ticket. The authors also bypassed three guardrails: a test-only hook, copied fixtures, and a raw commit and reset. | Later review passes after a repair, and orchestration; not fresh ticket authorship | high 0 labeled pairs; xhigh 0.40 over 2 pairs |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft336-bounded-output` ticket 8 | Sonnet / high / implementer | The comparison author reached green through a raw commit and reset outside `bench commit`, and a light path then denied that route in the guard. |
| `ft336-bounded-output` BO-C7 R59 | Fable / high / decider | The first decision spelled the registry value `4096` and missed the gate's own read limit; the second decision gave that limit its own registry entry. |
| `ref-inventory` RI-C2b R48 and R49 | Opus / high / reviewer | Two consecutive Coverage rounds each planted a symref in one window of the discard transaction, and each found a delete that removed the only handle. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| `review-evidence-file-pages` RE2 control | claude-fable-5-1 / high / reviewer | The narrow round found 8 findings in 399474 tokens and the full control found 6 in 227918 tokens. They agreed on 3, so neither shape covered the other. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Run a first review round on the harness mid binding. Escalate review to Fable only from round two on, by reviewer rule.
- Give each ticket a fresh Opus/high author session. Give each post-review repair a fresh Opus session at the ticket's declared effort.
- Make every fix delegate follow `bench-debug`, and give a duplication fold the test-before-fix phase only.
- Ask a Fable/high read-only consultation when the coordinator is uncertain, and check its claim against the tree.
- Do not route fresh ticket authorship to Sonnet/high, because the controlled comparison on three tickets preferred Opus on each.
- Charge each author with root conformance, every package it writes, and `cmd/bench` for a public response or embed change.
- Charge each author and repair with the duplicated-facts sweep and one recorded red per independent expectation.
- Read a helper before a repair charge names it as the fold target.
- Run the propose-writes preflight before dispatch, and fold fixture paths into the Writes lines; fence expansions are pre-approved.
- Keep the narrow review shape provisional. Pair narrow axes with one full reader in a later trial.
- Probe author claims independently and require current verification at the final chunk tip.
- Record ticket assignments before the first chunk freeze, and record later plan amendments.
- Run tests and probes serially when review axes share a tree.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings, and return exhaustion to the reviewer.
- Run the destination merge from inside the worktree before the final reconciliation, fence the merged paths, and review them before the landing.
- Cite the last chunk's seams before its final repair commit, because a plan commit after the chunk tip moves the tip and forces author reruns.
- Plant a symref at every ref a transaction window touches in the first Coverage round.
- Stop a consultation that produces no output past its expected window, and reroute it.
- Read the assignment census before landing releases the worktree.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
