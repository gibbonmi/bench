# 2. Derive the kit-source path through the canonical owner

Blocked by: none
Writes: internal/gate/kit_source.go, internal/gate/kit_source_test.go
Covers: TT51, TT52

## What to build

Chunk: TT-C2.

Make the gate's `resolvedPath` helper delegate to `canonicalpath.Resolve`. `KitSourceCheckout` then compares two absolute canonical spellings. A resolve error answers no match.

Add a test with `BENCH_KIT` set to `.` and the current directory at the root, through `t.Chdir`. Keep `TestKitSourceCheckoutMatchesThroughASymlinkSpelling` green without a change.

## Acceptance

- [ ] With `BENCH_KIT` set to `.` and the current directory at the root, `KitSourceCheckout(root)` answers true.
- [ ] A symlinked root spelling still matches the kit, and an unrelated directory does not.
