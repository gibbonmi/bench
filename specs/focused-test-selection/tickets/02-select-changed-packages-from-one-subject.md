# Select changed packages from one subject

Blocked by: 01-focus-explicit-package-and-test-runs

Writes: internal/diff/range.go, internal/diff/explicit_base_test.go, internal/diff/source_tip_pair_test.go, internal/testreport/command.go (new), internal/testreport/selection.go (new), internal/testreport/selection_test.go (new), cmd/bench/main.go, cmd/bench/command_registry_test.go, specs/focused-test-selection/

## What to build

Expose one root-bound coherent changed-subject seam from `internal/diff`.
Resolve the current Go package and embed graph once, add the transitive reverse
closure over production and both test edge classes. Run the deterministic
package union. Return explicit empty output for proven non-Go subjects and
refuse every Go-relevant path the current graph cannot map safely. Support the
default live, explicit-base live, and frozen base/source-tip forms.

## Acceptance checklist

- [ ] C01 — default live selection composes committed, staged, tracked-worktree, and untracked paths.
- [ ] C02 — explicit base stays live and base/source-tip stays immutable.
- [ ] C03 — one movement retry is allowed and a second drift refuses without partial output.
- [ ] C04 — Go files, exact embeds, and all three reverse-import edge classes select the complete closure.
- [ ] C05 — module metadata and mixed known inputs select one sorted, deduplicated union.
- [ ] C06 — empty and non-Go-only subjects render explicit zero-row result tables.
- [ ] C07 — control-byte and unsafe paths refuse without either Go child; surviving-package deletion still maps.
- [ ] C08 — `--run` applies unchanged to the complete changed-package union.
- [ ] N03 — changed runs leave every gate-owned record absent or byte-identical.

Delivered outcome: agents can turn a coherent diff into the smallest
provably complete current-package test run without gaining gate authority.
