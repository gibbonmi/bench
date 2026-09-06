# Fold and move the FT298 residual rules

Blocked by: none
Writes: .agents/skills/bench-craft-spec/references/map-discipline.md, .agents/skills/bench-craft-delegate/references/delegation-discipline.md, .agents/commands/bench-write-spec.md
Covers: none

## What to build

`roadmap/FT298.md` records five reviewer verdicts on the sixteen rules landed
at `4f86a7cf`. Apply the five verdicts. Leave the sixth item open for the
reviewer. That item is the batch-approval amendment decision.

1. **Fold.** Delete the `delegation-discipline.md` headroom-route bullet
   under `## In the charge`. `craft-tickets` already states this rule. Keep
   `craft-tickets` as the one source.

2. **Fold.** `map-discipline.md`'s "Before the map locks" section carries
   two bullets: a repo-wide literal-bytes search, and a changed-message test
   enumeration. A changed message is a moved literal. Join the two bullets
   into one. Keep every existing obligation.

3. **Fold.** `map-discipline.md`'s "Per row" section carries two bullets
   that grade the same producer path. One traces a failure message to its
   producer and confirms reachability. One names each guard the input
   passes before the producer. Join the two bullets into one. Keep every
   existing obligation.

4. **Move.** Delete the bench-signal bullet from `map-discipline.md`'s
   "Before the map locks" section. That bullet states a phase obligation of
   `/bench-write-spec`, not a coverage-map rule. Add the same obligation to
   `bench-write-spec.md`'s `## Process` step 1 ("Author").

5. **Fold.** `map-discipline.md`'s "At ticket slicing" section carries two
   one-line bullets: a pasted-operand quoting rule, and an anchor-needle
   quoting rule. An anchor needle is a pasted operand. Join the two bullets
   into one.

Every kept or written sentence obeys `references/ste-prose.md`. Keep each
sentence to 25 words or fewer. Keep each bullet's paragraph to six
sentences or fewer.

Check `.agents/commands/bench-write-spec.md`'s guidance-prose-budget row in
`projects/benchkit.md`. The row currently reads 73 lines. If the edit
pushes the file over budget, raise that row's number by the exact overage.
Make that change in the same commit. State whether the headroom route was
needed.

## Acceptance

- [ ] `delegation-discipline.md` carries no headroom-route bullet. `craft-tickets`' SKILL.md stays unchanged as the one source.
- [ ] `map-discipline.md`'s "Before the map locks" section carries one bullet for the moved-literal-bytes search and the changed-message enumeration.
- [ ] `map-discipline.md`'s "Per row" section carries one bullet for the producer-reachability and guards-before-the-producer criteria.
- [ ] `map-discipline.md` carries no bench-signal bullet. `bench-write-spec.md`'s `## Process` step 1 states the bench-signal quoting obligation.
- [ ] `map-discipline.md`'s "At ticket slicing" section carries one bullet for the pasted-operand and anchor-needle quoting rules.
- [ ] `bench gate-prose` passes on every touched file.
- [ ] `bench test --check guidance-prose-budgets` passes.
