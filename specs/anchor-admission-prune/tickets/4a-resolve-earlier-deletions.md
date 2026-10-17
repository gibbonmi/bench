# Resolve an earlier chunk's deletion in a later chunk charge

Blocked by: 4-prune-phase-registries.md
Writes: internal/preflight/gather.go, internal/preflight/closure.go, internal/preflight/decision.go, internal/preflight/fence_writes.go, internal/preflight/fence_writes_test.go, internal/diff/range.go, internal/diff/review_base_test.go
Covers: AP48

## What to build

A ticket charge grades the `Writes:` entries of every ticket in the spec with `writes-resolve`. With an explicit chunk `--base`, the deletion facts read only that chunk range and that base. So an entry for a file that an earlier chunk of the same build deleted resolves neither as an unmarked entry nor as a `(deleted)` entry. Today the six files that tickets 3 and 4 deleted block the charge of every later chunk.

Make such an entry resolve when a commit of the build's integration range deleted the path. That range starts at the build's landing base: the `main` tip merged before the first chunk, or the `main` tip of the last recorded fold. Compose the existing seam that derives that base; do not write a second derivation. An entry for a path that no commit of the build deleted stays red, and so does a never-tracked path.

## Acceptance

- [ ] The ticket 5 charge, `bench preflight build anchor-admission-prune --charge --ticket 5-prune-hand-pins.md`, passes `writes-resolve` with the AP-C7 chunk tip as `--base`.
- [ ] A regression test shows that an earlier-chunk deletion resolves and that a never-tracked path stays red.
- [ ] `bench test --package ./internal/preflight` and `bench test --check root-conformance` pass.
