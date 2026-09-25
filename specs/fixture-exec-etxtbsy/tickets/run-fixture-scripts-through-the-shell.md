# Run the gate fixture scripts through the shell

Blocked by: none
Writes: internal/testrepo/gate_fixture_test.go
Covers: none

## What to build

The gate fixture tests write a gate script and then run it. A parallel test's fork can hold a
write handle to the script. A direct exec of the script then fails with `text file busy`. One
test helper runs each fixture script through `/bin/sh`, so the kernel never executes the file.
Every gate fixture test that runs a script uses this helper.

## Acceptance

- [ ] `go test -count=300 -run TestGateFixture ./internal/testrepo/` passes, where the base tree failed with `text file busy`.
- [ ] No gate fixture test executes a fixture script directly.
- [ ] One helper comment holds the reason for the shell route.
