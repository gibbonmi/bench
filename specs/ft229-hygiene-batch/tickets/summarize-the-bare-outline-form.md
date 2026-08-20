# Summarize the bare outline form

Blocked by: none
Writes: internal/outline

## What to build

Bare `bench outline` stops emitting 200 arbitrary symbol rows. It emits meta
plus one row per scanned directory carrying that directory's symbol count, so a
cold probe costs a screen and the output is complete rather than a prefix. A
path argument keeps emitting symbol rows for that path. A new `--full` flag
emits symbol rows repository-wide, so the old capability moves rather than
disappears. The 200-row cap goes with the form that needed it. The output stays
AXI-conformant TOON with a definitive empty state.

## Acceptance

- [ ] bare `bench outline` emits meta and one row per scanned directory with its symbol count.
- [ ] `bench outline <path>` emits symbol rows for that path.
- [ ] `bench outline --full` emits symbol rows repository-wide.
- [ ] the bare form's row count equals the scanned directory count, with no cap applied.
- [ ] a tree with no scannable source yields the definitive empty state.
