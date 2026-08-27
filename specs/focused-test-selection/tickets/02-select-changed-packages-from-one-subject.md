# Select changed packages from one subject

Blocked by: 01-focus-explicit-package-and-test-runs

Writes: internal/diff/range.go, internal/diff/explicit_base_test.go, internal/diff/source_tip_pair_test.go, internal/testreport/command.go (new), internal/testreport/selection.go (new), internal/testreport/selection_test.go (new), cmd/bench/main.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, specs/focused-test-selection/

## What to build

Expose one root-bound coherent changed-subject seam from `internal/diff`.
Resolve the test-aware current Go graph once. Retain only its ordinary matched
package records. Add the transitive reverse closure over production and both
test edge classes. Run the deterministic package union. Return explicit empty
output for proven non-Go subjects.

Refuse every Go-relevant path that the current graph cannot map safely.
Support the default live, explicit-base live, and frozen base/source-tip
forms. Keep the complete public help inventory current with each runnable
form.

## Acceptance checklist

- [x] C01 — default live selection composes committed, staged, tracked-worktree, and untracked paths.
- [x] C02 — explicit base stays live and base/source-tip stays immutable.
- [x] C03 — one movement retry is allowed and a second drift refuses without partial output.
- [x] C04 — the live graph exposes all three embed-input classes on retained ordinary packages.
- [x] C05 — module metadata and mixed known inputs select one sorted, deduplicated union.
- [x] C06 — empty and non-Go-only subjects render explicit zero-row result tables.
- [x] C07 — control-byte and unsafe paths refuse without either Go child; surviving-package deletion still maps.
- [x] C08 — `--run` applies unchanged to the complete changed-package union.
- [x] C09 — exact Go and embed inputs close over all three reverse-import edge classes.
- [x] N03 — changed runs leave every gate-owned record absent or byte-identical.

Delivered outcome: agents can turn a coherent diff into the smallest
provably complete current-package test run without gaining gate authority.
