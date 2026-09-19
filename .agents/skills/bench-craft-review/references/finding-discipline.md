# Finding discipline

Charged from `craft-review` when an axis writes a finding. Each rule below settles one
question the finding must answer. `craft-review` keeps the three-axis split, the smell
baseline, and the universal-claim rule.

## What a string expectation proves

- A generated script's independently authored string expectation is the mutation catch.
  An expectation that reads the generator's own output grades the generator against
  itself.

## What a citation points at

- A finding cites the line the axis read this pass, or the symbol instead. A line number
  from an earlier pass points at bytes that moved.
- A review-owned row names the artifact that stores its evidence. Evidence that lives
  only in a delegate return is not citable at the landing.
- A review of gate-anchored prose names the anchor state of each sentence that it proposes to change. A change to an anchored sentence moves its anchor row, its test expectation, and its canary in the same repair.

## Where an axis under-reads

- A test-deleting Standards finding names the surviving assertion or file as coverage.
  A deletion with no named survivor is an open coverage hole.
- For a runnable defect claim, attempt refutation with a real run before reporting a strong finding.
  A strong finding is one that names a defect, a gap, or a violation without a hedge.
- For a mandatory standard without an automated check, cite the exact requirement and the violating source.
  State why executable refutation is unavailable.
  Inspect contrary source evidence and applicable exceptions before retaining the finding.
- An environment-variable Coverage finding cites the producer before it claims absence.
  The consumer alone does not show which producer binds the variable.
- The Coverage axis probes a test's fixture source, not only its assertion. A fixture
  that names a symbol the production file declares can stay green while the assertion
  never runs.
- A Coverage finding describes the tree before the probes of the axis. Behavior that those probes added is never a finding.
- An axis reads the seam cell of a row before it judges a review-owned row unmet. A review-owned seam places the evidence in the review record, not in a test.

## What an axis return carries

- An axis return lists the evidence cursors that it fetched. The coordinator then finds a partial retrieval before it accepts the result.

## When a ticket already decided

- A finding that contradicts a ticket's explicit keep decision is a `no-op`. The
  coordinator cites the ticket line in the disposition.

## When a seam cannot reach the state

- An unreachable row seam amends the row's seam column. The build records the helper seam
  it adds as a decision, so the reviewer sees the new surface.

## What a confidence states

- A finding carries a stated confidence as an integer from 0 to 10.
- The confidence never changes whether a finding blocks. The kind and the citation
  decide the round.
- `auto-fix` and `ask-user` label a finding `held`, and `no-op` labels it `refuted`.
  The reviewer's disposition is the label source.
- Optional advice carries no confidence.
