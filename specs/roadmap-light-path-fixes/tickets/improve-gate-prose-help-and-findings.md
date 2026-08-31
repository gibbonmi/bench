# Improve gate-prose help and findings

Blocked by: none
Writes: internal/gate/gate_prose.go, internal/gate/gate_prose_test.go
Covers: LF28

## What to build

Make bench gate-prose help exit zero on stdout. Include the offending sentence
with each finding's line number.

## Acceptance

- [ ] Help exits zero and writes usage to stdout.
- [ ] Unknown flags remain usage errors on stderr.
- [ ] Every finding includes its line number and offending sentence.

