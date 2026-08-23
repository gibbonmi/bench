# Learnings — usage journal

<!-- entries below -->

## 2026-08-22 - write-spec review round took three iterations  [open]
- **What happened:** The spec-and-tickets round for `asd-ste100-progressive-disclosure` took three iterations; the first returned "revise" from both reviewers. Review caught three gap classes. First, coverage rows with no gate seam: the moved core passages, the exclusion file's shrink rule, and four embedded templates. Second, rows attached to a seam that could not red them: walk and refusal rows on a parse-only package. Third, tickets past one context window: the check ticket, eight comment batches of ten thousand lines each, and a handoff named in every ticket.
- **Right behavior:** Before the round, the author walks each row and names the gate check or test that reds it. A row whose only seam is "the orchestrator reads" is review-owned and says so, or gets a gate twin. A comment ticket is sized by the lines the delegate must read, not the lines it edits. A phase-close artifact appears in no per-ticket `Writes:`. The author missed these because the map was written from the stories down and the tickets from the areas up. The citation audit checked ids, not seam reachability.
- **Proposed rule change:** `craft-spec` review rubric gains one question: "name the gate check or test that reds this row, or mark the row review-owned". `craft-tickets` gains one sentence: "size a rewrite ticket by the lines the delegate must read".

## 2026-08-22 - review round looped without a cap or a scoped charge  [open]
- **What happened:** The spec-and-tickets review round ran three iterations. The author declared no iteration cap at launch. Each re-review charge listed the folds to verify and then asked the delegates to apply the full four-lens pass again. A fresh full pass surfaces new small findings every time, so the round could not converge on its own. The reviewer stopped the loop.
- **Right behavior:** Declare the round's cap in the line before the first iteration; one iteration plus one fold is the default. Scope a re-review charge to the named folds. Ask only for a finding that changes observable behavior, an ownership fence, or the ticket graph. Treat everything below that bar as acceptance-round folds and record them in the verification log. Never relaunch a full-lens pass after the first iteration.
- **Proposed rule change:** `bench-write-spec.md` gains one sentence under the review step: "Declare the round's iteration cap before the first charge. A re-review charge names the folds and asks only for findings above the blocking bar." `craft-delegate` gains: "A re-review charge names what changed and asks only for regressions above the blocking bar."

## 2026-08-22 - three prose budgets lifted by two lines each, temporary  [open]
- **What happened:** The reviewer-accepted STE sweep turned three long sentences into vertical lists. The lists pushed `bench-craft-spec/SKILL.md` to 152 lines, and `bench-craft-delegate/SKILL.md` and `bench-craft-tdd/SKILL.md` to 122 lines. The profile's budget table gained one exact row for each file at the new number.
- **Right behavior:** Treat the three lifts as temporary. FT100 owns the budget measure and the cut line; its audit decides whether each file earns the lines or cuts back to the old bound.
- **Proposed rule change:** none. The drain routes this entry to FT100 with the three files and the old numbers (150, 120, 120).
