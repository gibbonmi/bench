# Run one named conformance check

Blocked by: 02-select-changed-packages-from-one-subject

Writes: internal/conformance/registry/scope.go (new), internal/conformance/gate_entry_test.go, internal/conformance/tier_test.go, internal/testreport/command.go (new), internal/testreport/testreport.go, internal/testreport/testreport_test.go, internal/testreport/runbinary_test.go, internal/testreport/cancel_test.go, internal/testreport/selection_test.go (new), internal/testreport/check_test.go (new), cmd/bench/main.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, specs/focused-test-selection/

## What to build

Add one registry-owned singular-scope transport to
`TestRootConformance`. Validate a requested check against the dev registry
and scrub ambient conformance controls. Run only the entry test with the root,
tier, scope, kit, and selected executable pinned. Refuse unknown, ship-only,
and conflicting forms before a child starts. Finish the public help text and
prove all focused forms remain non-authoritative.

## Acceptance checklist

- [x] F02 — final selector and subject-flag grammar refuses invalid combinations with exit 2.
- [x] K01 — one registered dev scope executes and timing names no other check.
- [x] K02 — invalid, ship-only, and conflicting check forms refuse before child start.
- [x] K03 — exact conformance controls and freshness checks defeat ambient and inherited redirects.
- [x] N02 — every form shares the validated binary, cancellation, decoder, and renderer chain.
- [x] N04 — named-check runs leave every gate-owned record absent or byte-identical.
- [x] H01 — public help names every focused form and claims no gate verdict.

Delivered outcome: agents can run one named conformance check through
`bench test` with the same environment owner as the gate, but without a gate
verdict.
