# Route spec arguments through usage

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/spec_test.go
Covers: LF6

## What to build

Replace specArg hand parsing with usage.Parse while preserving supported
operands and subcommands.

## Acceptance

- [ ] Help exits zero and writes usage to stdout.
- [ ] Unknown flags and missing operands exit two on stderr.
- [ ] Existing valid spec invocations retain their behavior.

