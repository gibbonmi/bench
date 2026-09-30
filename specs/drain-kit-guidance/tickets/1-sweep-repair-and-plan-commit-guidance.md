# Name the repair retention and the plan-commit sweeps in the build guidance

Blocked by: none
Writes: .claude/commands/bench-implement-spec.md, .agents/skills/bench-craft-tickets/references/slicing-checks.md, internal/anchors/registry_data.go, internal/anchors/registry_data_test.go, internal/anchors/registry_ticket_passes.go, internal/anchors/registry_ticket_passes_test.go
Covers: none

## What to build

The tree-targets build found five gaps in the build guidance. Close each gap at one site, and add no second copy of a rule that another file already states.

- A repair session retains its ticket verification entries at the final source before the chunk checkpoint. The implement-spec command retains the author verification only, so the first TT-C4 checkpoint refused.
- A plan commit that renames a spec symbol searches every ticket for the old symbol. Ticket 5 kept the old quoting function after the spec fold.
- A fence expansion reads the fence disposition text of the spec and corrects each sentence that names the added package. Two sentences went false after TT-C4 fence expansions.
- A posture change for new output names each stream that the output reaches. Fold this clause into the existing slicing check about the call sites of a helper that the change turns red. A helper that asserts on stderr went red after the posture list marked its file green.
- A ticket that moves a test expectation names the file that holds the old expectation in `Writes:`. Put this clause beside the slicing check about a changed rendered output shape.

Put the first rule in the implement-spec command beside its author retention sentence. Put the other four rules in the slicing checks. Keep each anchored sentence, or update its anchor registry entry in the same commit.

## Acceptance

- [ ] The implement-spec command names the retention of repair verification entries at the final source before the checkpoint.
- [ ] The slicing checks name the symbol sweep of a renaming plan commit and the sentence sweep of a fence expansion.
- [ ] The existing helper call-site check names each stream that the new output reaches, and no sibling bullet restates it.
- [ ] A slicing check names the file of a moved test expectation in `Writes:`.
- [ ] The anchor registry and the prose gate pass on each edited guidance file.
