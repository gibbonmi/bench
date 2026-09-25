# Claude model scorecard

Last incorporated landing: `ft336-bounded-output` (`5ec03981791ecf93a178021377bce1ba02b05113`, 2026-09-25).
Ten fresh Opus/high authors wrote ten tickets, and thirteen fresh Opus/low repair sessions closed the review findings.
Opus/medium ran every review axis, Fable/high decided each ask-user finding, and Sonnet/high authored three comparison tickets that never landed.
The retro recorded 8 labeled pairs from this session, with a Brier mean of 0.010 and 0 abstentions.
The harness token figure for each subagent is known, but the orchestrator tokens and provider costs remain unknown.

## Current routing

| model / effort | role and sample | observed quality | current use | calibration |
| --- | --- | --- | --- | --- |
| Fable / low–high | orchestrator, 36 landings + implementer, 10 charges + reviewer, 21 axes | On `slicing-closure` the medium axes found the trailing-slash bypass with a probe and two one-source defects, and each finding held. One Coverage probe on the shared tree reddened a concurrent Spec run, and two returns skipped changed evidence pages. | Coordination of a delegated build and adversarial spec review; it implements only when the reviewer names it | orchestrator 0.36 over 1 pair; medium reviewer 0.188 over 17 pairs; 0 abstained |
| Fable / high | reviewer, 3 axes + decider, 8 consultations + debug, 2 charges | On `ft336-bounded-output` each ask-user decision was usable, but one missed a production `4096` in the gate package until the coordinator sent the fact back. Its two debug charges traced a fence miss to the slicing guidance and split Sonnet's guardrail dodges into model and guardrail causes. | Ask-user decisions and debug diagnosis when the reviewer names the tier; review of another provider's candidate | unknown |
| Opus / high | implementer, latest 10 fresh ticket sessions; reviewer, latest 24 axes | On `ft336-bounded-output` each of ten fresh authors landed its ticket, but five needed a fence expansion for a file that a rendered shape or a registry required. A blinded grader preferred the Opus author over the Sonnet author on all three compared tickets. | High for fresh ticket authors, anchored guidance, foundational Go seams, and review axes by reviewer direction. | 0.110 over 59 pairs; 2 abstained |
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
| `ft311-landing-completion` Coverage axis | Opus / medium / reviewer | The axis probed a state file whose fence never closes with a scratch test, observed the next run refuse the document, and named the row that should exist. |

## Current decisions

- Change routing only after two comparable runs, one controlled comparison, or explicit user direction.
- The top tier implements only when the reviewer names it for the run.
- Give each ticket a fresh Opus/high author session, and give each post-review repair a fresh Opus/low session.
- Do not route fresh ticket authorship to Sonnet/high, because the controlled comparison on three tickets preferred Opus on each.
- Use a fresh writer when a fork would inherit another tier.
- Charge each author with the known traps of earlier chunks, such as registry budgets and a rendered-shape reader sweep.
- Tell each delegate to put backticks in single quotes, because a double-quoted backtick ran a command.
- Tell each delegate to give each search a path, because a search with no path waited on its input.
- Use the review tier that the reviewer names, and run each axis and each confirming round in a fresh session.
- Charge each author and each axis with a narrow read: one binding check, the needed pages or one diff, targeted reads, and a bounded return.
- Probe each author's done-claim at a site the author did not probe, because a repeat site proves nothing.
- Ask each current ticket author to rerun its verification at the final chunk tip before the checkpoint record.
- Record every ticket assignment in one plan commit before the first chunk freeze.
- Map each later plan change with an amendment entry in the review record.
- Let only one review axis run tests or probes on a shared tree, and start Coverage probes after other test runs end.
- Check each Fable decision against the tree before you apply it, because one decision missed a production literal.
- Keep Standards, Spec, and Coverage in separate native contexts and retain each return independently.
- Ask for a stated confidence in every write charge and every axis charge, and freeze it at return time.
- Reproduce a finding with a probe before you charge its repair.
- Apply the repair allowance to the chunk, and retain optional advice outside blocking findings.
- Freeze a chunk base at the close commit of the predecessor chunk, and freeze the first chunk base at the merged `main` tip.
- Merge the current `main` after the last chunk when the landing base moved, then run a review round and the final checks again.
- Include the package of each merged test file in the final checks.
- Run the whole-tree gate after the last repair and before landing.
- Read the assignment census before landing releases a worktree.
- Return material acceptance shortfalls and undecided behavior to the reviewer.
- Preserve unknown costs and usage, and do not change model defaults from this single run.
