# Pin every TOON byte class

Blocked by: pin-default-full-and-empty-classes.md
Ownership fence: `internal/axi/compatibility`, `cmd/bench/axi_compatibility_test.go`, `specs/axi-compatibility-oracle/testdata`
Integration surfaces: pair and empty matrix plus the candidate-rebuild harness→pin-default-full-and-empty-classes.md; byte-class case fixtures→`specs/axi-compatibility-oracle/testdata`; byte-class derivation→`internal/axi/compatibility`; TOON encoder adapter mutated by the BY1 rows→`internal/toon/toon.go` exercised unchanged by every BY1 row
Contracts: the byte-class case IDs `<member>-<class>` cross `internal/toon/toon.go`→`specs/axi-compatibility-oracle/testdata`; their type is one baseline observation record per case ID, membership is the five classes the TOON adapter owns — control-rune escaping, separator quoting, numeric-string typing, the final newline, and block order — order is stable case ID ascending, and a class with no case refuses the load; asserted by BY1 against the really rebuilt candidate executable
Closure: BY1/control-escape, BY1/quote-escape, BY1/numeric-string-typing, BY1/final-newline, BY1/block-order

## What to build

The TOON adapter owns five byte facts that no cap can express and no visible-text
comparison can see: control-rune escaping, quoting a cell that carries the field
separator, keeping a numeric-looking string a string, the trailing newline, and
the order blocks are emitted in. Each is a delta an agent parsing stdout would
trip over while the rendering still looks right to a human, so each gets its own
case and its own mutation.

The truncation bounds are the sibling ticket `pin-truncation-bound-edges.md`:
those break in a cap owner rather than in the encoder, so the two land green
independently.

Mutations are applied to a scratch copy of the tree and the candidate executable is
rebuilt through `scripts/go-build.sh`, exactly as in the blocking ticket. The
rebuild is bounded at 180s and each case child at 30s.

## Acceptance

- [ ] [BY1] (covers CO5) each of the five TOON byte classes compares byte-exact against the pinned baseline, and a candidate rebuilt with any one of those classes changed reports a raw stdout delta on the case that owns it.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| BY1/control-escape | stop escaping control runes in TOON cells in the candidate rebuild | the byte-class test | run `go test ./cmd/bench -run TestExactMatrixPreservesByteClasses/control -timeout 900s`; it fails at the raw stdout equality assertion for case `root-status-control-bytes`, reporting the raw control byte where the baseline carries its escape; the rebuild is bounded at 180s and each case child at 30s |
| BY1/quote-escape | stop quoting cells containing the field separator in the candidate rebuild | the byte-class test | run `go test ./cmd/bench -run TestExactMatrixPreservesByteClasses/quoting -timeout 900s`; it fails at the raw stdout equality assertion for case `root-anchors-quoted-cell`, reporting the unquoted cell that changes the row's field count; bounded by the 180s rebuild and 30s case deadlines |
| BY1/numeric-string-typing | encode a numeric-looking string cell as a TOON number in the candidate rebuild | the byte-class test | run `go test ./cmd/bench -run TestExactMatrixPreservesByteClasses/numeric_string -timeout 900s`; it fails at the raw stdout equality assertion for case `root-maps-numeric-string`, reporting the unquoted number where the baseline keeps the string spelling including its leading zero; bounded by the 180s rebuild and 30s case deadlines |
| BY1/final-newline | drop the trailing newline from the TOON table renderer in the candidate rebuild | the byte-class test | run `go test ./cmd/bench -run TestExactMatrixPreservesByteClasses/final_newline -timeout 900s`; it fails at the raw stdout equality assertion for case `root-guards-final-newline`, reporting a one-byte length difference with identical visible text; bounded by the 180s rebuild and 30s case deadlines |
| BY1/block-order | emit a multi-block command's blocks in reverse order in the candidate rebuild | the byte-class test | run `go test ./cmd/bench -run TestExactMatrixPreservesByteClasses/block_order -timeout 900s`; it fails at the raw stdout equality assertion for case `root-roadmap-context-blocks`, reporting the swapped block headers while every block body is byte-identical; bounded by the 180s rebuild and 30s case deadlines |
