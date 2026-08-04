# Teach the zero-matched-test probe failure in craft-tickets

Blocked by: teach-authorizing-input-mutations.md
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`
Contracts: none crosses

## What to build

The red-mutations field bullet in `craft-tickets` makes a probe that matches
zero tests a failed probe: red-mutation evidence carries the matched-test count
per probe, zero is a hard failure of the operation sequence, and authors verify
`-run` selectors against real top-level test names at authoring time — subtests
registered inside one function do not match a bare selector.

This is the hung-run hazard's silent sibling, and the bullet says so beside the
bounded-failure rule that already owns that hazard: `go test` exits `ok` on a
selector that matches nothing, so the evidence channel silently carries no
evidence. The recovery-discard case: `-run Recovery` against a package whose
only top-level function is `TestRuntimeWorktreeContracts` matched nothing,
while the neighboring ticket's `-run SpecBuild` selector matched — the hazard
is per-row and invisible to a skim.

Blocked by the input-mutations ticket only because both edit the same skill
file; the coordinator lands them serially.

## Acceptance

- [ ] [PC1] the red-mutations bullet requires the matched-test count per probe in red evidence and names zero a hard failure of the operation sequence.
- [ ] [PC2] the bullet tells authors to verify `-run` selectors against real top-level test names, and sits beside the bounded-failure rule.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC1 | revert the matched-count rule | the reviewer, at review | remove the sentences, re-read the bullet against `specs/authoring-hardening/spec.md` story 4, expect the promised rule to be absent |
| PC2 | keep the count rule but drop the authoring-time selector check | the reviewer, at review | remove the verify-selectors sentence, re-read against the story, expect the authoring half of the defense to be missing |
