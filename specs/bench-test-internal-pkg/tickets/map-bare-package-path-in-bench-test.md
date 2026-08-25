# Map a bare package path in bench test

Blocked by: none
Writes: internal/testreport/testreport.go, internal/testreport/testreport_test.go (new)

## What to build

`bench test internal/usage` passes the operand to `go test` unchanged, so Go
resolves it as the standard-library path and fails with `package
internal/usage is not in std`. The verb grades the wrong package for a name
every agent in this repo types. A bare operand that names a directory under
the repo root becomes the relative pattern `./<operand>` before the
invocation. An operand that already starts with `./`, `../`, or `/`, or that
names no directory under the root, passes through unchanged. Put the rule in
one pure function that takes the root and the operand, so the test grades the
mapping without a Go subprocess.

## Acceptance

- [ ] `bench test internal/usage` runs `./internal/usage` and reports `internal/usage` as the package row.
- [ ] `bench test internal/usage/...` maps to `./internal/usage/...`.
- [ ] `bench test ./internal/usage`, `bench test ../x`, and `bench test github.com/gibbonmi/bench/internal/usage` pass through unchanged.
- [ ] An operand that names no directory under the root passes through unchanged, so Go reports its own error.
- [ ] The default with no operand stays `./...`.
