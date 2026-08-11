# Make the family-home next cells identify its help action

Blocked by: none
Ownership fence: `internal/specbuild/render.go`, `internal/specbuild/render_test.go`
Integration surfaces: retained-run family home→`RenderRunsHome`; first retained slug→existing `buildAction`/`ActionFact.action` lifecycle action seam in `internal/specbuild/disclosure.go`; materialized action command→every `spec_build_runs.next` cell and the appended `help[0]` command
Contracts: none crosses

## What to build

Close the accepted Terra Spec finding P1-runs-home-next-help-mismatch. Build the
family's first retained-run inspection fact through the existing lifecycle action
seam before table encoding, materialize that action once, and reuse its rendered
command as every populated home row's `next` value and as the first help command.
Do not copy the retained run's lifecycle `Next` into this aggregate projection:
lifecycle detail remains available through `status <slug> --full`.
When the first slug cannot construct a safe action, emit empty `next` cells and
`help[0]`, preserving the existing fail-closed hostile-slug behavior.

## Acceptance

- [ ] [HM1] (covers SB4) (P1-runs-home-next-help-mismatch) every `next` cell in a populated family-home response is byte-equal to its first help command; if no safe first action can be constructed, every `next` cell is empty and the help envelope is empty.
