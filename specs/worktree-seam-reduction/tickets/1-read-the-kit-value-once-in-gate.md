# 1. Read the kit value once in the gate

Blocked by: none
Writes: internal/gate/kit_source.go, internal/gate/kit_source_test.go, internal/gate/kit_value_test.go (new), internal/gate/phases.go, internal/gate/lane_select.go
Covers: WS1, WS2, WS3, WS4, WS5, WS6, WS7, WS8

## What to build

Chunk: SR-C1.

Add `gate.KitValue`, the one function that reads the raw `BENCH_KIT` value.
Make `kitRoot` in `phases.go` and `KitDir` call it.

Add `gate.LaneForCommitAtKit(root, kit)` and `gate.KitSourceCheckoutAtKit(root, kit)` to
`kit_source.go`. An empty kit makes the lane form fall back to the graded root. An empty
kit makes the kit-source form fall back to the parent of the running executable, as
`KitDir` does. That fallback does not call `KitDir`, so neither new form reaches
`KitValue`. Make `LaneForCommit` and `KitSourceCheckout` wrappers that pass
`KitValue()`, so that their callers in other packages do not change.

Put the new tests in `kit_value_test.go`. The tests use temporary directories and an
explicit kit, and they create no repository and start no process.

## Acceptance

- [ ] The gate non-test source holds one `BENCH_KIT` read, inside `KitValue`.
- [ ] `LaneForCommitAtKit` returns the root's manifest lane for a kit apart from the root.
- [ ] `LaneForCommitAtKit` with an empty kit returns a selective lane.
- [ ] `KitSourceCheckoutAtKit` matches a kit that names the root through a symbolic link, and refuses another directory.
- [ ] `KitSourceCheckoutAtKit` with an empty kit compares the root with the executable's parent.
- [ ] `KitDir` and `KitSourceCheckout` keep their answers under a bound `BENCH_KIT`.
