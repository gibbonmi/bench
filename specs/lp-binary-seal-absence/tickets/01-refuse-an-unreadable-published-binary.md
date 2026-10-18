# Refuse an unreadable published binary in the binary-seal row

Blocked by: none
Writes: internal/preflight/binary_seal.go, internal/preflight/binary_seal_test.go
Covers: none

## What to build

This ticket comes from drain d-359accf567c1, learning L182.

The build preflight gathers the binary-seal facts with an lstat of `dist/bench`.
At the moment, each lstat error makes the binary absent.
An absent binary makes the binary-seal row not applicable, and the audit treats a not-applicable row as exempt.
Thus an unreadable `dist/bench` lets an audit omit manual review.

Make only a "does not exist" lstat error an absence.
Any other lstat error makes the binary present, with a refusal that holds that error as its cause.
Use the refusal shape that the verifier fallback already uses, and keep one construction site for it.
Keep the lstat, because a stat follows a dangling link and reports a broken binary as absent.

## Acceptance

- [ ] A root with no `dist/bench` gives a not-applicable binary-seal row.
- [ ] A root where the lstat of `dist/bench` fails for a cause other than "does not exist" gives a present binary with a refusal, so the build binary-seal row is red.
- [ ] A sealed, a stale, and a dangling `dist/bench` keep their current results.
