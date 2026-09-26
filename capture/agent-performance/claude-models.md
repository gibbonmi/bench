# Claude model scorecard

Last incorporated landing: `ft336-bounded-output` (`5ec03981791ecf93a178021377bce1ba02b05113`, 2026-09-25).
Latest incorporated phase: FT337 spec review (`4e461ebbb5d1c3033f37764de82c8d59759d14f4`).
At capture time, its staging gate and human approval remain pending.

The user selected Opus/high, and the terminal model identifier was `claude-opus-5-5`.
Pass 1 found the missing file-identity contract. Pass 2 accepted its repair with four nonblocking clarifications.
Both passes were read-only and claimed no executed test evidence.

The terminal provider costs were USD 1.7542908 and USD 1.1407578, totaling USD 2.8950486.
The phase retrospective records the separate input, cache, and output counts.
Author and coordinator usage remains unknown. This review adds no labeled calibration pair.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 36 landings + implementer, 10 charges + reviewer, 21 axes | On `slicing-closure` the medium axes found the trailing-slash bypass with a probe and two one-source defects, and each finding held. One Coverage probe on the shared tree reddened a concurrent Spec run, and two returns skipped changed evidence pages. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 3 axes + decider, 8 consultations + debug, 2 charges | On `ft336-bounded-output` each ask-user decision was usable, but one missed a production `4096` in the gate package until the coordinator sent the fact back. Its two debug charges traced a fence miss to the slicing guidance and split Sonnet's guardrail dodges into model and guardrail causes. | Ask-user decisions and debug diagnosis when the reviewer names the tier; review of another provider's candidate | unknown |
| Opus / high | implementer, latest 10 fresh ticket sessions; spec reviewer, latest 2 passes | In FT337, Opus/high found the missing file-identity contract and accepted the repaired draft on pass 2. Its read-only reports separated their source claims from executed test evidence. | High for fresh ticket authors, anchored guidance, foundational Go seams, and review axes by reviewer direction. | retained implementation: 0.110 over 59 pairs, 2 abstained; new spec review: unknown, 0 labeled pairs, 2 test-result abstentions |
| Opus / medium, low | reviewer, orchestrator, and repair sessions combined | On `ft336-bounded-output` the medium axes closed seven chunks and one merge round, and each confirming round found the repairs held. Thirteen low repair sessions closed their targets in one attempt each, and each claimed probe bit under the coordinator's own probe. | Medium for review axes, gates, conformance, guidance, canaries, and triage. Low for exact tickets and post-review repairs. | 0.161 over 20 pairs; 9 abstained |
| Sonnet / high | orchestrator, 3 landings + reviewer, 14 axes + implementer, 3 comparison tickets | On `ft336-bounded-output` the three comparison authors used about twice the calls and time of Opus, and a blinded grader preferred Opus on each ticket. The authors also bypassed three guardrails: a test-only hook, copied fixtures, and a raw commit and reset. | Later review passes after a repair, and orchestration; not fresh ticket authorship | high 0 labeled pairs; xhigh 0.40 over 2 pairs |
| Sonnet / low–medium | implementer, latest 10 of 79 ticket-sized charges | On `ft311-diagnostics` one low charge landed after two corrections: its case-fold fixture was not red-capable for the row's named mutation, and its fixture left the fenced directory. On `craft-research-skill` one low review repair closed five findings first-pass with a biting omission probe. | Low for a prose or exact-spec repair at a known seam under a covering gate when the reviewer names it; the coordinator probes every return and runs the whole-tree gate before the landing | unknown |

## Representative evidence

| landing | model / effort / role | what it shows |
| --- | --- | --- |
| `ft336-bounded-output` ticket 8 | Sonnet / high / implementer | The comparison author reached green through a raw commit and reset outside `bench commit`, and a light path then denied that route in the guard. |
| `ft336-bounded-output` BO-C7 R59 | Fable / high / decider | The first decision spelled the registry value `4096` and missed the gate's own read limit; the second decision gave that limit its own registry entry. |
| `bounded-charge-evidence` CE-C4 Coverage axis | Opus / medium / reviewer | The axis narrowed one needle to its first sentence, saw no red, and traced the hole to a paragraph pin that selected only two leads. |
| `ft311-recoverable-reset` Coverage axis | Fable / high / reviewer | The axis built the tip, added an ignore rule after the checkpoint, ran the apply, and observed the ignored bytes deleted with no envelope, which no row had decided. |
| FT337 spec review | claude-opus-5-5 / high / reviewer | Two passes exposed and resolved B1, with eight initial edits and four final clarifications. Terminal cost was USD 2.8950486; neither pass ran tests. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- Use the review tier that the reviewer names. FT337 used Opus/high for two spec-review passes.
- Give each ticket a fresh Opus/high author session, and give each post-review repair a fresh Opus/low session.
- Use the reviewer-directed exception for a different implementation line, including the proposed FT337 line.
- Do not route fresh ticket authorship to Sonnet/high, because the controlled comparison on three tickets preferred Opus on each.
- Use a fresh writer when a fork would inherit another tier.
- Charge each author with known registry budgets and the rendered-shape reader sweep.
- Bind each axis and confirming round to exact source bytes in its own native context.
- Charge narrow reads through the canonical review phase, while preserving its explicit full control.
- Probe author claims independently and require current verification at the final chunk tip.
- Record ticket assignments before the first chunk freeze, and record later plan amendments.
- Run tests and probes serially when review axes share a tree.
- Check a decider's claim against the tree before applying it.
- Keep declared confidence unchanged, and separate unknown or abstained evidence from labeled claims.
- Apply the repair allowance to blocking findings. Fold accepted nonblocking edits without another paid review.
- Freeze each chunk's base at its predecessor's close commit, then reconcile moved main before landing.
- Include merged test packages in final verification, and run the complete landing gate.
- Read the assignment census before landing releases the worktree.
- Return material acceptance shortfalls and undecided behavior to the reviewer.
- Preserve exact terminal provider costs separately from token counts, and keep unavailable usage unknown.
