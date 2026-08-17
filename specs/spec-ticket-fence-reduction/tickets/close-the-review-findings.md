# Close the review findings

Blocked by: anchor-the-realigned-consumers.md
Writes: internal/coverage/coverage.go, internal/coverage/coverage_test.go, CONTEXT.md, internal/anchors/registry_data.go, tests/canary/workflow-guidance-anchors, reviews/spec-ticket-fence-reduction.md

## What to build

Repair ticket for the accepted findings in `reviews/spec-ticket-fence-reduction.md`
(deleted by this ticket's commit). The `projection()` comment in `coverage.go`
stops narrating the deleted opt-in flag. `TestCommand` gains two cases for the
behavior cell now reaching the TOON sink: a control byte is refused with exit 1
and the `unrepresentable TOON cell` line, and a comma- or tab-bearing behavior
renders as one quoted three-field row. `CONTEXT.md`'s **acceptance row** entry
keeps its definition and Avoid list and drops the restated grading sentence,
with its anchor and canary following the trimmed text. The reduced-schema text
in `docs/field-guide.html`, `docs/reporesident-distillation.md`, and the CHANGELOG
entry gains `Require` anchors and canaries so SR29 is gate-graded.

## Acceptance

- [ ] `(covers SR9, SR10)` A control-bearing behavior cell makes `Command` return exit 1 with the `unrepresentable TOON cell` error; a delimiter-bearing behavior renders as one quoted row — both asserted literally.
- [ ] `(covers SR28)` `CONTEXT.md` **acceptance row** carries no grading rule; its anchor and canary bite on the trimmed text.
- [ ] `(covers SR29)` Reverting the field guide's map card, the distillation sentence, or the CHANGELOG entry turns the gate red through a canary each.
- [ ] `(covers SR26)` `bench canary` and `bench gate` green; `reviews/spec-ticket-fence-reduction.md` is deleted.
