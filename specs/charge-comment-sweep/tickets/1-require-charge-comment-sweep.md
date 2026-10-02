# Require a comment sweep in each author and repair charge

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md, specs/charge-comment-sweep/tickets/1-require-charge-comment-sweep.md
Covers: none

## What to build

Review rounds found stale or duplicated code comments after the authors
committed their tickets. The delegation discipline requires a duplicated-facts
sweep of the delegate's own delta before the ticket commit. That sweep does not
name the comments in the delta.

The duplicated-facts sweep in "In the charge" also examines each comment that
the delta adds or changes against `craft-comments`. Thus each author charge and
each repair charge gets a comment sweep before the ticket commit. The final
repair attempt repeats the duplicated-facts sweep, so it repeats the comment
sweep too. The rule refers to `craft-comments` and does not restate its rules.

## Acceptance

- [x] The "In the charge" rule for the two results before the ticket commit names a sweep of each added or changed comment against `craft-comments`.
- [x] The anchored sentences of that rule and of the final-attempt rule do not change, and the anchor check stays green.

## Verification

The `docs-currency-workflow` check passes. A `bench probe` run changed "two
results" to "three results" in the anchored sentence, and the root conformance
test was red on the duplicated-facts sweep anchor. The probe restored the file.
A `bench probe` run removed the new comment-sweep sentence, and the check stayed
green, because no anchor holds that sentence.
