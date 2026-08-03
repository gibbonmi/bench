# Learnings — usage journal

Append one entry when you deviate from the workflow, make a judgment call you're
unsure about, catch a should-have-asked in hindsight, or catch yourself assembling
the same ad-hoc check a second time (a codification candidate — name the `bench`
subcommand it wants to be). You capture; the reviewer
decides. `/bench-what-next` verdicts every open entry in its reviewed batch
diff: work-shaped and rule-shaped entries become roadmap items (rule-shaped ones
built later under the synthesis discipline), the rest are dismissed with one line
of why. A resolved entry leaves this file, and its verdict is recorded in the
drain commit. The journal holds open entries only; history lives in git. Never
rewrite a kit rule yourself — that is the whole point of capturing here instead.

Format per entry. Heading: `## YYYY-MM-DD — short title  [open]`

- **What happened:** …
- **Right behavior:** …
- **Proposed rule change:** … (or "none")

## 2026-08-03 — Skipped two charged craft-tickets steps at FT181 ticket derivation  [open]

- **What happened:** The coordinator derived four tickets without executing "Discover the contracts before writing files" (no junction row for the specbuild→worktree abandon contract) and without the prefactor rule (the shape classifier was a shared primitive of two tickets, stated in prose to both delegates, each of which derived its own copy). Ticket 3's Assumptions line actively encoded the gap: "specbuild abandon tests drive a counting fake owner, so no internal/worktree change is needed here" — a spec ambiguity converted into a build instruction. The semantic review then paid two repair rounds for what the skipped steps would have caught at derivation (findings P1/C1, S1).
- **Right behavior:** Execute both steps at derivation and land their output in the ticket artifact, so a skip is visible.
- **Proposed rule change:** Partially landed 2026-08-03 in 2892501 (reviewer-directed): the ticket template now requires a `Contracts:` field; craft-spec gained the composition degenerate and the existing-control edge rule. Residual to verdict: the skip was possible because prose-only steps leave no artifact — any remaining prose-only steps in the phase commands deserve the same landing-site treatment (prefactor has no artifact slot yet).

## 2026-08-03 — Auto-repaired every review finding without triaging blocking vs contestable  [open]

- **What happened:** The round-1 review returned 8 blocking findings; the coordinator routed all into repair delegates without distinguishing must-fix production defects (P1, P2) from judgment-priced hardening (S1, C2, C4) and without surfacing the triage to the reviewer before spending the repair rounds.
- **Right behavior:** When a review returns material blocking findings, present the triage — "N must-fix, M contestable, promote-now or harden-first" — before charging repairs; the finding's disposition may be the review delegate's, but build-priority is the reviewer's.
- **Proposed rule change:** /bench-implement-spec's finding-disposition paragraph could require the triage be shown when a review returns more than a threshold of blocking findings; reviewer decides the threshold or rejects the rule.

## 2026-08-03 — Did not audit my own process compliance when the extra work appeared  [open]

- **What happened:** Two direct signals (the ticket-3 delegate's fake-owner caveat at integration; round-1's "pinned against fakes" findings) pointed at coordinator process skips, and the coordinator absorbed the repair work without asking where it came from until the reviewer forced the question.
- **Right behavior:** When a review produces findings the build's own discipline claims to prevent, the first question is "which charged step did not run," answered against the skill text before the repair pass starts — invariant 1 applied to process, not just code.
- **Proposed rule change:** none (behavioral; the Contracts-field landing makes the specific skip visible, and the retro carries the pattern).

## 2026-08-03 — Review-receipt disposition "accepted" means an open defect and blocks promote  [open]

- **What happened:** The coordinator wrote a review receipt marking closed findings "accepted" (intending endorsed/closed); `hasAcceptedFinding` reads that as an unrepaired defect, promote refused, and clearing the recorded review required a recomposition round.
- **Right behavior:** Closed findings carry a non-"accepted" disposition (resolved/endorsed); "accepted" is reserved for a defect the review accepts as real and unrepaired.
- **Proposed rule change:** The receipt-skeleton helper (parked in capture/IDEAS.md 2026-08-03) should enumerate the disposition vocabulary, or the receipt schema should validate the field so the wrong word is refused at submission.
