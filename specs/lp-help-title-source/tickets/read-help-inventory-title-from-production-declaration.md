# Read the help inventory title from its production declaration

Blocked by: none
Writes: internal/conformance/help_inventory_single_source_test.go
Covers: none

## What to build

The constant `helpInventoryTitle` in `cmd/bench/command_registry.go` holds the title that `bench help` prints. The help inventory single-source check counts that title on the production surfaces. The check kept a second literal copy of the title, because `internal/conformance` cannot import package `main`. Thus a title change edited the check too.

The check reads the title from the `helpInventoryTitle` declaration in its existing parse of `cmd/bench/command_registry.go`. The check fails closed when that declaration is absent or is not a non-empty string constant. The copied literal goes, and the fixture declares its own title.

## Acceptance

- [ ] A change to the title text of `helpInventoryTitle` in `cmd/bench/command_registry.go` keeps the root conformance pass green with no edit to the check.
- [ ] The check reds the root conformance pass when the title text appears twice on the production surfaces.
- [ ] The check reds when `cmd/bench/command_registry.go` declares no string constant `helpInventoryTitle`.
