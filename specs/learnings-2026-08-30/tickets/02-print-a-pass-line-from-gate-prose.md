# Print a pass line from gate-prose

Blocked by: none
Writes: internal/gate/gate_prose.go, internal/gate/gate_prose_test.go, internal/commit/lane_test.go

## What to build

`bench gate-prose <root> -- <path...>` exits 0 with no output on a pass. A
caller cannot tell a clean list from a path list that graded nothing, and
must infer the verdict from the exit code alone. Bench responses are bounded
and complete, so the verb states its verdict.

On a pass, `GateProseCommand` in `internal/gate/gate_prose.go` prints one
`prose[N]{path,verdict}` table to stdout through the `toon` package, one row
per named path with the verdict `pass`. This is the shape `roadmap/FT270.md`
decides for the verb. An empty path list prints the table header with `N`
equal to 0. A red list keeps its current shape: one finding line per finding
and no trailer. Exit codes do not change.

The word `green` never appears in the pass output. The lane composes this
verb as its `prose` check and writes each check's stdout through
`prefixedPhaseWriters`. `internal/commit/lane_test.go` states that a lane
pass is not a graded green. A `green` token in the lane stream would
contradict that property.

A `bench commit` lane pass may carry the table under the `[prose]` check
prefix before the `lane{outcome=pass,...}` line. That extra output is
accepted. The `internal/commit/lane_test.go` prose check is a shell stand-in,
so those tests never see the verb's output; do not edit them.

## Acceptance

- [ ] `GateProseCommand([root, "--", "docs/notes.md"])` on a clean file exits 0 and prints a `prose[1]{path,verdict}` table whose one row is `docs/notes.md,pass`, and nothing else to stdout.
- [ ] `GateProseCommand([root])` exits 0 and prints the `prose[0]{path,verdict}` header alone.
- [ ] A red list prints only the finding lines and exits 1; no `prose[` table appears.
- [ ] The pass output contains no `green` token.
- [ ] `TestGateProseCommandCleanList` and `TestGateProseCommandEmptyPathList` in `internal/gate/gate_prose_test.go` assert the pass line, and the delegate records them red before the fix.
- [ ] `go test ./internal/gate/ -run 'GateProse|Lane' ./internal/commit/ -parallel 2` passes.
