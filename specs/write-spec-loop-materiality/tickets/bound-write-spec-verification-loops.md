# Bound the write-spec verification loops

Blocked by: none
Writes: .agents/commands/bench-write-spec.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors/write-spec-materiality-exit/, tests/canary/workflow-guidance-anchors/write-spec-degenerate-cheapness/

## What to build

Reviewer-approved (2026-08-14, conversation) governance for `/bench-write-spec`'s
verification loops, motivated by the 33-round FT189 loop: a finding blocks only
when it changes observable behavior, an ownership fence, or the ticket graph —
a prose-or-accounting-only round is the acceptance round; a revision may not
add promises beyond the decision source unless a blocking finding demands it;
and loop 1's degenerate standard is the cheapest *plausible* wrong
implementation, routing contrived degenerates to the build's mutation probes
instead of new spec rows. Both new clauses land outside the existing pinned
anchor needles and get their own anchor rows plus workflow-guidance canary
fixtures, so the fixture-bite proof is the red that shows each anchor bites.
Test seam: the workflow-guidance conformance family (`docs-currency-workflow`
owner through the auto-discovered fixture universe).

## Acceptance

- [ ] `.agents/commands/bench-write-spec.md` carries the materiality exit, the promise-inflation guard, and the cheapest-plausible degenerate standard, with every pre-existing anchor needle byte-intact
- [ ] Each new clause has a Require anchor row and a canary fixture whose mutation makes `docs-currency-workflow` emit the row's diagnostic and whose restore clears it
