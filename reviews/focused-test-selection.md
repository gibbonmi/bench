# Focused test selection final repair review

Frozen base: `0dae83df5372727d608c6fb2b29fa118d732562f`

Prior reviewed tip: `dd049e235aa7f9643105c5d5d975da3edfe14683`

Reviewed repair tip: `282956a54549f96dff6f89e0b51b8491174a6f01`

Raw findings: 2. De-duplicated repair targets: 1.

## Standards

Count: 0. Worst issue: none.

## Spec

Count: 1. Worst issue: P1 a retained stdout pipe can block the final drain.

- **P1 — auto-fix — drain before a retained pipe can block completion.** N02
  at `spec.md:257` requires process-group cancellation. In
  `internal/testreport/command.go:209-231`, both completion orders can wait for
  decode or process completion before the final drain. The escalation at lines
  259-266 kills only the leader. A SIGINT-resistant group descendant that keeps
  stdout open blocks the awaited completion and makes the drain unreachable.

## Coverage

Count: 1. Worst issue: P1 the cancellation fixture does not retain stdout.

- **P1 — auto-fix — cover a resistant descendant that inherits stdout.**
  `internal/testreport/cancel_test.go:60-94` starts the parking child with nil
  stdout, which routes output to the null device. Make the child retain the Go
  stdout pipe. The test must require the structured interruption, no partial
  tables, and descendant removal before its deadline.
