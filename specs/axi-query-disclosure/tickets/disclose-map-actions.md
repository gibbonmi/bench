# Disclose map frontier and repair actions

Blocked by: record-ft173-log-leverage.md, disclose-learnings-actions.md
Writes: `internal/maps/`

## What to build

Make the public `maps` query derive actions from each rendered map row: a frontier names `/bench-shape-idea`, and an invalid map names `bench maps --template` together with its diagnostic path. Carry each row's values, deduplicate exact templates without reordering, and keep empty or terminal results honestly action-free.

## Acceptance

- [ ] [MH1] (covers QD1) every actionable frontier and invalid-map row yields its own correctly valued action, including a many-row fixture that rejects first-match-only derivation and guessed paths.
- [ ] [MH2] (covers QD1) empty and complete map results append the honest zero-row help block, with no repair or shaping busywork.
- [ ] [MH3] (covers QD6) old-to-new public-command fixtures prove that each named maps state changes only by its appended help block; primary bytes, stream, exit, and argv behavior remain byte-equal.
