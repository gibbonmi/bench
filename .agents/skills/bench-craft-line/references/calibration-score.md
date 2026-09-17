# Calibration score

## What the declaration states

The line declaration states the expected repair-round count and a stated confidence as an integer from 0 to 10.
The count is the number of repair cycles the author expects after the initial review.
The confidence states how sure the author is of that count. A real-number confidence is not permitted.

## How one claim scores

One claim's calibration score is `(p - label)^2`, with `p = n / 10` and label 1 for `held` or 0 for `refuted`.
Here `n` is the stated confidence. The Brier mean is the mean of these scores over the labeled claims.
An abstention scores 0, stays out of the Brier mean, and is counted apart.
The log rule is excluded, because a wrong 0 or a wrong 10 scores to negative infinity under it.

## Who writes a label

A label source is the gate, the coordinator's probe of the exact tree, or the reviewer's disposition.
A model judgment is never a label source.
The repair-attribution table's actual round count labels the expectation `held` when it equals the expected count and `refuted` otherwise.
