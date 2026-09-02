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

## Where an axis under-reads

- A test-deleting Standards finding names the surviving assertion or file as coverage.
  A deletion with no named survivor is an open coverage hole.
- An axis refutes a strong finding with a real run before the axis reports the finding.
  A strong finding is one that names a defect, a gap, or a violation without a hedge.
- An environment-variable Coverage finding cites the producer before it claims absence.
  The consumer alone does not show which producer binds the variable.

## When a seam cannot reach the state

- An unreachable row seam amends the row's seam column. The build records the helper seam
  it adds as a decision, so the reviewer sees the new surface.
