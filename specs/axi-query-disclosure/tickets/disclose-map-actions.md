# Disclose map frontier and repair actions

Blocked by: record-ft173-log-leverage.md
Writes: `internal/axi/`, `internal/maps/`

## What to build

Extend the typed AXI action owner with the reusable executable-invocation and harness-phase kinds and kind-specific renderer dispatch, then consume both through the public `maps` query: a frontier names `/bench-shape-idea` without a `bench ` prefix, and an invalid map names `bench maps --template` together with its diagnostic path. Carry every known argv token, use only explicitly unknown future-input slots, carry each row's values, deduplicate exact templates without reordering, and keep empty or terminal results honestly action-free. Existing diff action bytes remain unchanged.

## Acceptance

- [ ] [MH1] (covers QD1) the reusable invocation owner rejects dropped known argv, guessed values, undeclared placeholders, and prose-as-command; the phase owner and renderer emit the canonical `/bench-shape-idea` phase without a `bench ` prefix while preserving existing diff action bytes; and the public `maps` query consumes both kinds.
- [ ] [MH2] (covers QD1) every actionable frontier and every invalid-map row whose disclosure cells are representable yields its own correctly valued action, including a many-row fixture that rejects first-match-only derivation and guessed paths. An unsupported diagnostic-only `why` value preserves the primary diagnostic and exit and appends honest empty help as specified by `repair-control-bearing-action-values.md`.
- [ ] [MH3] (covers QD1) empty and complete map results append the honest zero-row help block, with no repair or shaping busywork.
- [ ] [MH4] (covers QD6) old-to-new public-command fixtures prove that each named maps state changes only by its appended help block; primary bytes, stream, exit, and argv behavior remain byte-equal.
