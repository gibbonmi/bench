# Scaffold declared-input hygiene for consumers

Blocked by: none
Writes: internal/adopt/setup.go, internal/adopt/setup_test.go, internal/conformance/validity_checks_test.go
Covers: LF2

## What to build

Ship the gitignored-declared-input check in the consumer gate scaffold. Reuse
the declared-input grammar and keep ordinary ignored files outside the check.

## Acceptance

- [ ] A linked consumer rejects a declared input that Git ignores.
- [ ] An undeclared ignored file remains allowed.
- [ ] Paths with spaces or glob characters are treated literally.

