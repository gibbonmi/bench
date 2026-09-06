# FT183 per-component-scoping review residuals

Status: ready

## Destination

Dispositions for the two faces left by the shipped per-component-gate-scoping
review: the unreachable whole-changeset reduced fallback in `internal/gate`'s
component scoping, and the unbound derivation-source check in
`internal/conformance` that grades `internal/gate`'s registry.

## Notes

## Decisions so far

- [Retire the whole-changeset reduced fallback](ft183-gate-scoping-residuals/tickets/1.md): Resolved 2026-08-02: remove it.
- [Observe which derivation a registry row resolves through](ft183-gate-scoping-residuals/tickets/2.md): Resolved 2026-08-03: the summary exists at `decisions/ft183-gate-scoping-residuals/assets/ft183-derivation-binding.md`.
- [Which binding mechanism, if any](ft183-gate-scoping-residuals/tickets/3.md): Resolved 2026-08-03: candidate A, function identity.
- [The hand-declared exemption leaves canary's resolver ungraded](ft183-gate-scoping-residuals/tickets/4.md): Resolved 2026-08-03: bind the canary row too.
- [The orphaned Reduced verdict record class](ft183-gate-scoping-residuals/tickets/5.md): Resolved 2026-08-03: retire it fully — fields, readers, and record class go with the path.

## Not yet specified

## Spec-writer discretion

- Exact placement and naming of the identity check within `internal/gate`'s
  test files, and the shape of the `Source → function` expectation table.
  These stay open provided the demonstrated-red requirement, the
  method-expression guard, and the refuse-unknown-rows exhaustiveness
  survive.

## Out of scope

- Reintroducing any whole-changeset reduced path for the kit root; #1 closed
  that direction.

## Sources

- Path: `decisions/ft183-gate-scoping-residuals/assets/ft183-derivation-binding.md`
  Supports: #2's summary and the factual premises of #3 and #4. Two resolver swaps verified to pass the derivation-source check. The asset also prices five candidate mechanisms. Produced 2026-08-03 by two read-only research delegates, and corrected after the doc review.
  Drift: re-verify if `internal/gate/component_inputs.go`'s registry shape or `internal/conformance/derivation_source_test.go` changes before the spec reads this map.
