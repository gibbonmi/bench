# Teach fence derivation from the Contracts crossing

Blocked by: none
Ownership fence: `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/conformance/docs_workflow_helpers_test.go`
Contracts: the taught fence-anchoring clause crosses `.agents/skills/bench-craft-tickets/SKILL.md`→`internal/conformance/docs_workflow_helpers_test.go`, asserted by FD2 against the real needle check
Assumptions: `.claude/skills/bench-craft-tickets` is a symlink to the `.agents` copy, so one edit serves both harnesses; the `Contracts:` field already teaches the `<fence>→<fence>` grammar and this ticket constrains it rather than inventing it; claims re-derived from the tree at pickup

## What to build

`craft-tickets` states that a ticket's ownership fence is derived from its
`Contracts:` line rather than from the lines that prompted the ticket, and that
each crossing anchors at least one backticked path inside the ticket's own
fence. A ticket author reading the skill can tell, before writing the file, that
an artifact merely *advertising* a fact the ticket changes is a path the ticket
writes.

## Acceptance

- [ ] [FD1] `craft-tickets` states the fence-derivation rule and the anchoring requirement in the ticket-template field notes.
- [ ] [FD2] the docs needle check pins the new clause, so deleting it turns the conformance phase red.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| FD1 | delete the fence-derivation sentence from the `Ownership fence:` field note | the docs needle check | remove the sentence, run the scoped conformance check, expect the dropped-clause diagnostic |
| FD2 | delete the needle registration for the new clause | the needle check's own count | remove the `requireCollapsed` call, restore the sentence, confirm the mutation in FD1 no longer goes red |
