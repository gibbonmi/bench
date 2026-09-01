# Stop truncating package operands at a -test. flag

Blocked by: none
Writes: internal/gate/tag_census.go, internal/gate/tag_census_test.go
Covers: none

## What to build

`testPackagesOf` in `internal/gate/tag_census.go` stops the package scan at
the first `-test.`-prefixed argument. A phase argv can write a
`-test.`-prefixed flag, for example `-test.run=Pattern`, ahead of its package
operand. The scan then drops that operand. The census then records no
package for the phase.

Only `-args` truly ends package-operand evidence. Everything after `-args`
passes to the test binary uninterpreted. Skip a `-test.`-prefixed argument
instead of ending the scan on it. Keep `-args` as the one true scan
terminator.

## Acceptance

- [x] A `-test.`-prefixed flag ahead of a package operand does not drop that
      operand from the census.
- [x] An `-args` flag still excludes every following argument from the
      package census.
