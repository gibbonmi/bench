# The gate fixture test runs its scripts without a text-file-busy race

Blocked by: none
Writes: internal/testrepo/gate_fixture_test.go

## What to build

`TestGateFixtureWritesOnlyRequestedPathsAndCanRepeat` runs each written gate
script through `/bin/sh`. It does not execute the script file directly.

A parallel subtest can fork while another subtest holds a write handle to its
script. The forked child keeps that handle until its own exec. A direct exec of
the script then fails with `text file busy`. A shell reads the script and does
not execute the file, so the race cannot occur. The test still asserts that
each written script has its executable bits.

Source: the 2026-09-17 drain learning "A text-file-busy test flake was
attributed to a markdown-only landing". The drain reproduced the flake with
`go test ./internal/testrepo -run TestGateFixtureWritesOnlyRequestedPathsAndCanRepeat -count=400`:
12 of 1,600 subtests failed with `text file busy`. The attribution half of that
learning goes to FT141.

## Acceptance

- [ ] The repro command above passes with `-count=400`.
- [ ] The test fails when a written script loses its executable bits.
- [ ] The test fails when a script prints the wrong body.
- [ ] `bench gate` green.
