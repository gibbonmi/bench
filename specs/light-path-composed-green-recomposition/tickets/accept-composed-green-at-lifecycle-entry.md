# Accept composed green at lifecycle entry

Blocked by: none
Ownership fence: `internal/gate`, `internal/specbuild`, `projects/benchkit.md`, `.bench/BENCH-reference.md`, `internal/contract`, `internal/conformance`, `tests/canary`
Assumptions: reduced verdicts inherit only from retained full-green ancestors; partial verdict validation already cross-checks every skipped component's evidence; prep-release retains its stricter ship-tier rule; claims re-derived from the tree at pickup

## What to build

Spec-build lifecycle entry accepts an exact-tip green verdict whose executed and inherited evidence together cover the whole tree, while every incomplete or invalid evidence chain still fails closed with an executable sufficient remediation.

## Acceptance

- [ ] [CG1] The gate package exposes one predicate that accepts full green and valid exact-tip reduced or partial green only when composed evidence covers every skipped phase or component.
- [ ] [CG2] `bench spec build start` succeeds on a reduced-green exact tip with intact inherited evidence and refuses absent, red, or broken-chain evidence.
- [ ] [CG3] Every spec-build lifecycle precondition stating the same whole-tree-green fact routes through the gate-owned predicate rather than re-deriving verdict classes.
- [ ] [CG4] A refusal that requires a fresh whole-tree run directs the operator to `bench gate --fresh`, while prep-release's distinct ship-tier demand remains unchanged.
- [ ] [CG5] Gate and lifecycle documentation describe composed dev green without weakening reduced-run inheritance or per-component validation.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| CG1 | remove inherited coverage for one skipped phase or component | gate predicate unit tests | construct each verdict class through the gate fixture, run `go test ./internal/gate`, expect incomplete composed evidence to be refused |
| CG2 | restore the exact full-verdict-only lifecycle check | specbuild start regression test | create full-green ancestor and reduced-green specs-only tip, invoke Start, run `go test ./internal/specbuild`, expect the old refusal |
| CG3 | bypass the gate predicate in one lifecycle bootstrap caller | specbuild lifecycle tests | exercise start and recomposition bootstrap paths, run `go test ./internal/specbuild`, expect the inconsistent authorization failure |
| CG4 | change the sufficient remediation back to plain `bench gate` | lifecycle and anchor contract tests | trigger the refusal through the public lifecycle seam, run the owning package tests, expect the missing `--fresh` assertion |
| CG5 | restore the old claim that every reduced dev verdict is lifecycle-ineligible | conformance and canary anchors | run the focused contract and canary checks, expect the pinned prose mismatch |
