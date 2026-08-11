# Repair the ticket-provenance test comment

Blocked by: none
Ownership fence: `internal/specbuild/reclaim_test.go`
Integration surfaces: none crosses
Contracts: none crosses
Closure: TP1/timeless-comment-register

## What to build

Close the accepted Standards finding (S2) from the Terra/xhigh review of
candidate `399dca908c7b1e1a4162eb7625e497dfb6786750`:
`internal/specbuild/reclaim_test.go:481-482` opens
`TestReclaimApplyFailureRendersOneHelpBlockOverTheSpentDeletions` with

    // A partial apply renders the spent-deletion tables and the refusal as one response, not
    // two help envelopes stapled together (DH1).

The trailing `(DH1)` names the repair ticket that produced the test. That is
change history addressed to the reviewer of a merged diff, not a fact about
the code as it stands, and `.agents/skills/bench-craft-comments/SKILL.md`
rules it out directly: the register is timeless present, with "no provenance
('per the spec', 'as requested in review')", because that prose "becomes noise
the day the diff merges".

Reword the comment so it states only what the test proves — a partial apply
renders the spent-deletion tables and the refusal inside one response envelope
— carrying no ticket ID, review round, candidate, or repair-history token. The
test body, its name, its assertions, and every other line of the file stay
byte-identical.

The red signal is worth naming honestly rather than manufacturing. This is a
prose-only edit with no runtime-observable behavior, and a Go test that
regex-matched this comment for a bracketed ticket ID would be a second source
for a rule `bench-craft-comments` already owns: it would grade one hand-picked
comment while every other comment in the tree stays ungraded, and it would
turn a review judgment into an assertion that cannot distinguish a provenance
label from any other parenthesized capital. So the mutation below is verified
by re-inspection against the skill, and the existing focused reclaim render
test is what proves the edit touched nothing but prose.

## Acceptance

- [ ] [TP1] (covers local) (S2) the comment above
  `TestReclaimApplyFailureRendersOneHelpBlockOverTheSpentDeletions` describes
  the behavior the test pins in the timeless present, naming no ticket ID,
  review round, candidate, or repair history; the test's name, body, and
  assertions and every other line of `internal/specbuild/reclaim_test.go` are
  unchanged, and the test stays green.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| TP1/timeless-comment-register | reinstate a ticket-provenance token in the comment — the literal `(DH1)`, a "repair round" phrase, or any equivalent change-history label | re-inspection against `.agents/skills/bench-craft-comments/SKILL.md`, with the existing focused reclaim render test proving the edit is comment-only | read the comment against the skill's register rule and require it to name no ticket, review, or change; then run `go test ./internal/specbuild -run TestReclaimApplyFailureRendersOneHelpBlockOverTheSpentDeletions` and require it green with its assertions untouched |
