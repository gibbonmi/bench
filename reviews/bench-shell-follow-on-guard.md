# Guard Bench commands from shell follow-ons

Frozen base: `f05c5e7a4a4e03b7539455822d045aa2cd0a9e92`

Reviewed tip: `ff488a6aac4384480399d2d11fceb47ad1c575e6`

Prior reviewed tip: `dc8c0880adee49315506a8c884879222d5470446`

## Standards

Finding count: 2. Worst issue: P1 duplicated shell-prefix policy.

- P1, `auto-fix`: Bench guard and Git guard independently project command words
  and encode the same routine-prefix set. The repair already made those copies
  diverge. Move the shared projection and prefix facts into `internal/shellcommand`
  while preserving each guard's policy-specific result. Citation:
  `internal/benchguard/benchguard.go` and `internal/gitguard/scan.go`.
- P2, `auto-fix`: Ticket 06 still claims `Writes: internal/gate`, but its proof
  moved to the approved `cmd/bench` seam. Correct the ticket fence to describe
  the resulting state. Citation: ticket 06 and `cmd/bench/main_test.go`.

## Spec

Finding count: 1. Worst issue: P1 value-taking prefix options hide Bench.

- P1, `auto-fix`: `env -u X`, `timeout -s KILL`, and `xargs -n 1` stop prefix
  resolution at the option value, so a following Bench call is allowed to carry
  shell syntax. Consume supported value-taking option forms without treating
  query-only `command` forms as execution. Citation: spec routine-prefix edge,
  ticket 04, and `internal/benchguard/benchguard.go`.

## Coverage

Finding count: 3. Worst issue: P1 prefix option gaps.

- P1, `auto-fix`: Add process cases for value-taking `env`, `timeout`, and
  `xargs` options; exact probes currently exit 0. Citation:
  `internal/systemtest/bench_follow_on_test.go` prefix matrix.
- P1, `auto-fix`: `command -V bench | cat` is a query, not a Bench execution,
  and must remain allowed. Add a process allowance case. Citation: story 10 and
  `internal/benchguard/benchguard.go`.
- P2, `auto-fix`: The wrapper acceptance names outer pipeline, `&&`, and
  redirection, but the process suite asserts only the pipeline. Add the other
  two outer forms. Citation: ticket 04 and the FOG18 process matrix.
