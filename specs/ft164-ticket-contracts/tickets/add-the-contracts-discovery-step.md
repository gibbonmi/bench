# Add the contracts-discovery step to the breakdown

Blocked by: open-the-breakdown-on-blast-radius.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/conformance/docs_workflow_helpers_test.go`, `internal/conformance/fixture_bite_test.go`
Assumptions: the step lands between the breakdown procedure and the write-one-file-per-ticket section and owns its own H2 so `scopedSection` can resolve it; the ownership-fence cross-pointer to `craft-spec` is pointed at rather than restated; each needle takes one byte-exact row in `TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting`. Re-derive from the tree at pickup.

## What to build

FT164 story 4: between drafting the breakdown and writing the ticket files, the
coordinator names what crosses each ownership fence, so a cross-fence mismatch
is a sentence at slicing time instead of a composed red six tickets later — the
exact failure the last build paid for, with two tickets each green alone and red
together over an undeclared domain mismatch.

Every value crossing a fence names four facts: its type, its membership or
domain rule, its ordering, and its absence semantics. Each discovered invariant
lands as an acceptance row on the *consumer* ticket, asserted against the real
producer and the whole enumerated family — not against a fixture and not against
the members that happen to exist in one worktree, which is how both halves
looked green while composed behavior was broken.

When neither side can assert the invariant alone, it gets a junction ticket; and
when a junction row sits more than one ticket downstream of the fence it
describes, a narrower copy moves to the junction so the red surfaces where the
mismatch is, not six tickets later. Ticket claims are re-derived from the tree
after earlier tickets land, never from the spec's account of the base — one
ticket in the last build asserted three defects its predecessor had already
fixed.

Enforcement is anchor-first in the same diff: four needles scoped to the new
step's section, each with one mutation-table row proving its diagnostic fires.

## Acceptance

- [ ] [CD1] the step requires every fence-crossing value to name its type, membership or domain rule, ordering, and absence semantics.
- [ ] [CD2] each discovered invariant is taught as a consumer-ticket acceptance row asserted against the real producer and the whole enumerated family.
- [ ] [CD3] the junction rule is taught: neither side able to assert alone means a junction ticket, and a junction row more than one ticket downstream moves a narrower copy to the junction.
- [ ] [CD4] ticket claims are taught as re-derived from the tree after earlier tickets land, never from the spec's account of the base.
- [ ] [CD5] every needle this ticket registers has a byte-exact entry in `TestSpecBuildCadenceAnchorsRejectDeletionSwapAndRawGitRouting` whose named diagnostic fires.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CD1 | reduce the four-fact sentence to "check the fences agree" | the `contracts four facts` mutation subtest | swap the sentence, run the anchor check, expect the four-fact diagnostic |
| CD2 | drop "the whole enumerated family" from the consumer-row sentence | the `contracts consumer row` mutation subtest | swap the sentence, run the anchor check, expect the real-producer-and-family diagnostic |
| CD3 | delete the downstream-copy half of the junction rule | the `contracts junction` mutation subtest | delete the clause, run the anchor check, expect the junction diagnostic |
| CD4 | swap re-derivation from the tree for re-reading the spec | the `contracts re-derivation` mutation subtest | swap the sentence, run the anchor check, expect the re-derivation diagnostic |
| CD5 | register a needle and land no mutation entry for it | `review` plus the mutation harness | delete one entry, run the harness and watch the remaining rows still pass, then read this ticket's needle list against the table |
