# Migrate the shift outcome

Blocked by: none
Ownership fence: `internal/shift`
Integration surfaces: shared outcome carrier→`internal/axi` exercised by SH1; compatibility oracle→`cmd/bench/axi_compatibility_test.go` exercised by SH1; final contraction→contract-outcome-action-routes.md
Contracts: the shift kind, exact exit integer, ordered `shift_result` fields, sanitized detail, and the durable interrupted state cross `internal/shift`'s `Result`→`Result.Emit`, membership is `complete`, `failed`, `usage`, `incomplete`, `no-op`, `interrupted`, order is the current `resultFields` slice, and absence is `branch: none` with `recovery: none`, asserted by SH1 against the real `Loop` exit paths
Closure: SH1/complete, SH1/failed, SH1/usage, SH1/incomplete, SH1/no-op, SH1/interrupted, SH1/field-order, SH1/detail-sanitize, SH1/interrupted-reload, SH1/bypass

## What to build

`shift.Result` constructs a shared outcome before `Emit` renders the `shift_result` block. `exitCodes` stays the single source of the kind-to-exit mapping (0/1/2/3/4/130), `resultFields` stays the single source of field order, `sanitize.Preview` stays the only detail filter, and `evidenceOutcome` keeps splitting a post-mutation failure into `failed` versus `incomplete` by committed count. Public bytes and exits are unchanged.

OA8 is the spec's observed-red row: its recorded red was the absence probe for the shared owner, which is finished once `internal/axi` exists. Every mutation below is a post-implementation subject mutation on `internal/shift`, not a re-run of that probe.

Refresh precondition: `internal/axi` must exist as Go package `axi` with `go doc ./internal/axi Outcome` resolving, from `specs/axi-carriers-and-registry/tickets/introduce-outcome-carrier.md`. `cmd/bench/axi_compatibility_test.go` is an unchanged input this ticket only runs.

Write the owner test `TestShiftOutcomeRoutePerKind` in `internal/shift/result_test.go`: one subtest per closure token, driving the real `Loop` to each terminal outcome with a stub adapter script in `t.TempDir()`, asserting the constructed outcome's kind, carried exit, and emitted block. `SH1/interrupted-reload` runs the fresh-process case through the built `bench` binary.

## Acceptance

- [ ] [SH1] (covers OA8) all six `shift.Outcome` kinds construct a shared outcome carrying `exitCodes`' exact exit, cross that route before `Emit`, and render the current `shift_result` block byte-for-byte.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SH1/complete | map `OutcomeComplete` to exit 3 in `exitCodes` | `TestShiftOutcomeRoutePerKind/complete` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/complete -timeout 120s` with a stub adapter that commits once then reports done; the exit assertion fails with `exit = 3, want 0` and the emitted `exit` cell shows `3`; the stub adapter exits immediately and the loop's iteration cap is 2, so the run cannot spin |
| SH1/failed | replace `evidenceOutcome`'s `committed == 0` test with `committed < 0`, so a zero-commit failure reports `incomplete` | `TestShiftOutcomeRoutePerKind/failed` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/failed -timeout 120s` with a stub adapter that fails before any commit; the kind assertion fails with `outcome = "incomplete", want "failed"` and exit `3, want 1`; iteration cap 2 with an immediately-exiting stub |
| SH1/usage | return the setup-failure path's `Result` with `OutcomeFailed` instead of `OutcomeUsage` | `TestShiftOutcomeRoutePerKind/usage` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/usage -timeout 60s` with a missing adapter; the kind assertion fails with `outcome = "failed", want "usage"` and exit `1, want 2`; the failure precedes any child spawn |
| SH1/incomplete | make the exhausted-cap path report `OutcomeComplete` instead of `OutcomeIncomplete` | `TestShiftOutcomeRoutePerKind/incomplete` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/incomplete -timeout 120s` with a stub adapter that always changes the tree and an iteration cap of 2; the kind assertion fails with `outcome = "complete", want "incomplete"` and exit `0, want 3`; the cap of 2 bounds the loop |
| SH1/no-op | make the no-change path report `OutcomeComplete` instead of `OutcomeNoOp` | `TestShiftOutcomeRoutePerKind/no-op` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/no-op -timeout 120s` with a stub adapter that changes nothing; the kind assertion fails with `outcome = "complete", want "no-op"` and exit `0, want 4`; iteration cap 2 with an immediately-exiting stub |
| SH1/interrupted | map `OutcomeInterrupted` to exit 1 in `exitCodes` | `TestShiftOutcomeRoutePerKind/interrupted` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/interrupted -timeout 120s`, cancelling the loop through its existing signal channel after the first iteration; the exit assertion fails with `exit = 1, want 130`; the checkpoint path's existing outcome deadline reaps the child, so a hung adapter still terminates |
| SH1/field-order | swap `"exit"` and `"branch"` in `resultFields` | `TestShiftOutcomeRoutePerKind/field-order` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/field-order -timeout 120s`; the raw-block assertion fails because the `shift_result` header reads `outcome,branch,exit,...` instead of `outcome,exit,branch,...`; stub adapter with iteration cap 2 |
| SH1/detail-sanitize | emit `r.Detail` raw instead of `sanitize.Preview(r.Detail)` | `TestShiftOutcomeRoutePerKind/detail-sanitize` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/detail-sanitize -timeout 120s` with a stub adapter whose failure detail carries `\x1b[31m`; the raw-block assertion fails on the unescaped escape sequence; stub adapter with iteration cap 2 |
| SH1/interrupted-reload | stop writing the interrupted outcome to the intent entry in `finish`, so a fresh process reads no interrupted state | `TestShiftOutcomeRoutePerKind/interrupted-reload` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/interrupted-reload -timeout 180s`, interrupting one `bench shift` process then reading the ledger from a second; the reload assertion fails with `entry.Outcome = "", want "interrupted"`; both processes run under `bounds.TestDeadline(30*time.Second)` and the first is killed at that deadline if it does not exit |
| SH1/bypass | have `Emit` build its row from `r` directly without constructing the shared outcome | `TestShiftOutcomeRoutePerKind/bypass` | apply the mutation, run `go test ./internal/shift -run TestShiftOutcomeRoutePerKind/bypass -timeout 120s`; the route assertion fails with `observed outcome constructions = 0, want 1` for every one of the six kinds while the emitted block stays byte-identical; stub adapter with iteration cap 2 |
