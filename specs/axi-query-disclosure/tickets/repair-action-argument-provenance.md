# Repair action argument provenance

Blocked by: none
Writes: `internal/axi/action.go`, `internal/axi/action_test.go`, `internal/maps/maps_test.go`, `internal/coverage/coverage_test.go`, `internal/worktree/list_actions_test.go`

## What to build

Remove the test-fitted literal `unknown` rejection from the shared executable-action owner. Keep structural validation there; prove exact known values and rejection of guessed replacements at the call-site derivation tests where provenance is observable.

## Acceptance

- [ ] [RP1] (covers QD1) `KnownArgument("unknown")` is a valid literal value, while empty, control-bearing, placeholder-bearing, and unsafe executable-name values remain refused by the shared owner.
- [ ] [RP2] (covers QD1) maps, coverage, and worktree derivation fixtures independently pin their source-derived paths, ids, and argv, so replacing one carried value with a guessed literal turns its owning test red.
- [ ] [RP3] (covers QD1) the shared owner contains no literal-value blacklist masquerading as provenance validation.
