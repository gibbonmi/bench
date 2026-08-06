# Teach the covers schema in the skills and pin it in conformance

Blocked by: parse-ticket-covers-annotations.md
Ownership fence: `.agents/skills/bench-craft-spec/SKILL.md`, `.agents/skills/bench-craft-tickets/SKILL.md`, `internal/conformance/example_agreement_test.go`, `internal/anchors/registry_data.go`
Integration surfaces: production parser for the taught example→parse-ticket-covers-annotations.md; needle evaluation→existing internal/conformance/docs_workflow_helpers_test.go checkWorkflowAnchors exercised by SK3; skill frontmatter and index conformance→existing checks over both SKILL.md files, content-only edits
Contracts: the taught example's covers annotations crossing `.agents/skills/bench-craft-tickets/SKILL.md`→`internal/conformance/example_agreement_test.go`, asserted by SK2 against the production ParseTicket

Fence breadth (four directories, kept whole): the example-agreement check and
the taught example it grades must change in one green landing, and the needle
registry pins the same clauses the skill edits introduce.

## What to build

`craft-spec` teaches the optional leading `row` column and its opt-in meaning
in the acceptance-coverage-map field enumeration; `craft-tickets` teaches the
covers grammar — `(covers <ID>)` / `(covers local)` after the row-ID bracket,
the single-ID-rows-only rule, and `local` as graded honesty — in the
Acceptance-rows bullet, and its taught example (the ticket-example marker
block) carries covers annotations on RC1 and RC2.

The example-agreement check grows a `taughtExampleCovers` expected literal and
a third `ticketExampleFieldDiag` comparison in `gradeTicketExample`, so
stripping the annotations turns the check red rather than thinning the prose.
Both new schema clauses get `anchors.Anchor` entries in `registry_data.go`
following the existing craft-spec/craft-tickets needle shape. Frontmatter and
the skills index are untouched; `.claude/skills/` entries are symlinks, so the
`.agents/` edit is the whole edit.

## Acceptance

- [ ] [SK1] craft-spec teaches the optional `row` column and its opt-in meaning, and craft-tickets teaches the covers grammar and the local rule with a taught example carrying covers annotations
- [ ] [SK2] the example-agreement check grades the taught example's covers against grown expected literals, and stripping the annotations from the example turns the check red
- [ ] [SK3] the new schema clauses in both skills are pinned by docs needles, and deleting either clause turns the conformance docs phase red

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SK1 | delete the local-rule sentence from craft-tickets | its docs needle | apply to a temp copy, run the needle evaluation, expect its diagnostic |
| SK2 | strip the covers annotations from the taught example rows | the example-agreement covers literal | apply to a temp copy, run `go test ./internal/conformance -run TestRootConformance`, expect the covers-agreement diagnostic |
| SK3 | delete the craft-spec row-column clause | its docs needle | apply to a temp copy, run the needle evaluation, expect its diagnostic |
