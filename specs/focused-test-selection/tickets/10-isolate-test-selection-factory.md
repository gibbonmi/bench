# Isolate the test selection factory

Blocked by: 09-drain-before-pipe-completion.md
Writes: internal/testreport/runbinary_test.go, internal/testreport/check_test.go

## What to build

Make the shared test factory own its selection environment. Prevent an outer gate selection from bypassing the installed test factory.

Keep tests that explicitly install an inherited selection able to do so after the helper returns.

## Acceptance

- [x] K03 — the named-check environment test uses its installed selection under the gate environment.
- [x] N02 — tests that explicitly set an inherited selection keep their refusal and reuse coverage.

Delivered outcome: the shared test factory clears the outer gate selection.
Tests can install deliberate inherited selections after that boundary.
