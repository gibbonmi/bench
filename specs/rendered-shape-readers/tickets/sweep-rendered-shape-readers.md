# Sweep rendered-shape readers into the ticket Writes line

Blocked by: none
Writes: .agents/skills/bench-craft-tickets/references/slicing-checks.md, .agents/skills/bench-craft-spec/references/map-discipline.md, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: none

## What to build

In the ft336 build, four tickets changed a rendered output shape, and each one widened its
fence after the build started. The slicing rules sweep readers by symbol and by decision
fact. No rule found the tests, the help inventories, the generated references, or the anchor
registries that pinned the old text. One new slicing rule runs one `rg` for the old text and
puts each hit on the ticket's own `Writes:` line. The map-discipline proof checklist gains a
`Rendered-shape readers` class that cites the needle and the hits, and it refers to the
slicing rule. One anchor row pins the rule sentence, and one independent test expectation
pins the anchor row.

## Acceptance

- [ ] `bench test --check docs-currency-workflow` fails on the base guidance with the new anchor row, and passes after the guidance edit.
- [ ] A probe that omits the new rule sentence turns `docs-currency-workflow` red, and the restore returns it to green.
- [ ] No other file states the rendered-shape sweep rule, and `map-discipline.md` refers to it without a restatement.
