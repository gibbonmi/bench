# Learnings — usage journal

- 2026-08-20 — The land-executable-freshness spec review took 2 iterations.
  Stage that missed: spec authoring. What review caught: the spec promised the
  freshness refusal fires "before any landing proof", but no row could go red
  on a check placed just before publication, because every proof between the
  argument parse and `landReviewed` is read-only — the no-mutation row (LF2)
  cannot distinguish first from last read-only position. Why it was missed:
  the author rowed the promise's *effect* (nothing mutated) instead of its
  *predicate* (which refusal wins when two proofs would both refuse), even
  though the tree already held the idiom for ordering rows
  (`TestLandCommandRefusesDestinationAndSourceStateBeforeGate`'s call
  counter). Proposed rule change: `craft-spec`'s map discipline gains one
  line — an ordering or before/after promise needs a row where two refusals
  compete and the map names which message wins; a row that only asserts
  absence of side effects does not cover ordering.
