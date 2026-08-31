# Review pickup — otel-seam-record

Base: effc457b7a73c8b81af808c1e312ba24b1a0aaa4. Reviewed tip: 14a905925f60cdf4eda33422ed50ef5ee08bef7c.
Raw findings: Standards 6, Spec 2, Coverage 4. De-duplicated repair targets: 9 — 6 auto-fix, 3 ask-user.

## Standards

Count: 6. Worst issue: S1, the exec exit regression.

- S1 — auto-fix (same target as Spec F1). `internal/worktree/exec.go:84-96`: `execSpanExit` maps each nonzero child exit to 1, and `ExecCommand` returns the mapped value. The verb's own grammar help promises the child's own exit ("any other exit 2 is the child's own"). The base returned `runWorktreeChild(...)` verbatim. Repair: give the span closer the mapped outcome, and return the child's own code. Writes: internal/worktree/exec.go, internal/worktree/exec_test.go.
- S2 — auto-fix through the C3 test (same target as Coverage C3). The instrumented hook set is declared twice: six `Hook: true` rows in `cmd/bench/main.go` and six `hook.*` rows in `internal/otelrecord/registry.go`. No check reconciles the two sets, and the registry's own comment says a second list does not exist. The C3 set-equality test is the reconciliation. The deeper choice — derive one list from the other — is a seam decision; see the reviewer notes.
- S3 — auto-fix. Five files each spell the outcome vocabulary and an exit policy: `hookSpanOutcome`, `gateSpanOutcome`, `commitSpanOutcome`, `worktreeSpanOutcome`, `phaseSpanOutcome`, plus `execSpanExit`. Two policies exist: zero is green, and zero-or-three is green. Repair: move the outcome constants and the two named mappers into `internal/otelrecord`; keep the call sites. Writes: internal/otelrecord/attributes.go, cmd/bench/main.go, internal/gate/gate.go, internal/gate/runner.go, internal/commit/commit.go, internal/worktree/worktree.go, internal/worktree/exec.go.
- S4 — ask-user. The begin-span open/close sequence is pasted six times (`beginHookSpan`, `beginCommitSpan`, `beginGateSpan`, `beginLaneSpan`, `beginLandingSpan`, `beginVerbSpan`). An `otelrecord.Begin` primitive would own the protocol. This is a seam decision.
- S5 — no-op. `worktree.Home` and `homeEnv` are pass-through shims over `internal/benchhome`, kept for seven live callers. Flagged so the shim does not become permanent.
- S6 — auto-fix. Comment provenance: the duplicated three-line comment pair in `internal/worktree/main_test.go:97-99,143-145`, and the "row OT22 is review-owned" clause in `internal/otelrecord/attributes.go:3-5`. Writes: internal/worktree/main_test.go, internal/otelrecord/attributes.go.

## Spec

Count: 2. Worst issue: F1, the exec exit regression.

- F1 — auto-fix. Same target as S1. Traced execution: `bench worktree exec <label> -- sh -c 'exit 3'` returns 1; the base returned 3.
- F2 — ask-user. Story 17 / row OT20 is partial. The worktree dispatch has twelve subverbs; five carry spans (`create`, `exec`, `merge`, `release`, `land`). `build` and `reauthorize` act on one named assignment and carry no span; `clean` and `reclaim` are bulk verbs with no span. The reviewer decides the instrumented set.

Reviewer note: the spec lists cross-process span parenting as a Won't handle, and `withGateSpanEnv` injects a traceparent. These two disagree on paper; the reviewer disposes.

## Coverage

Count: 4. Worst issue: C1, the FIFO hang.

- C1 — ask-user. A FIFO or a device at `<home>/otel/<key>/traces.jsonl` blocks the first `OnStart`, so each recorded verb hangs. `Writer.Append` checks only the directory, not the file. The profile's hostile-input checklist names this class. The spec never decided the refusal semantics for the file itself. Proposed row: the writer refuses a non-regular file at the record path.
- C2 — auto-fix. A symlinked `<home>/otel` parent defeats the symlink refusal, because the Lstat covers only the leaf. `os.MkdirAll` follows the parent link, so the record writes outside the Bench home. Story 8 names this escape. Repair: grade every level below the home, with a writer test that symlinks the parent. Writes: internal/otelrecord/writer.go, internal/otelrecord/writer_test.go.
- C3 — auto-fix. No test binds the six `hook.*` registry seams to the `Hook: true` dispatch rows. All six registry entries name `beginHookSpan`, so the conformance check cannot see a dropped flag. Repair: a `cmd/bench` unit test that asserts set equality between the `hook.` registry seams and the `Hook: true` rows. Add the analog for the four `worktree.*` seam constants. Writes: cmd/bench/main_test.go (or a new cmd/bench test file), internal/worktree/worktree_test.go if the analog lands there.
- C4 — auto-fix, test-only. The env pair `BENCH_OTEL_ROOT` / `BENCH_OTEL_TRACEPARENT` is untested on both sides. No test asserts the `gateEnv()` strip. No test covers an ambient `BENCH_OTEL_ROOT` from the operator's shell. Repair: an `internal/gate` unit test for the strip, and a system test for the ambient value. Writes: internal/gate/gate_test.go (or the package's existing test file), internal/systemtest/otel_gate_test.go.
