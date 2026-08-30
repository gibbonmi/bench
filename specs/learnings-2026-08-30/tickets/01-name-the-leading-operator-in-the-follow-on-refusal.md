# Name the leading operator in the follow-on refusal

Blocked by: none
Writes: internal/benchguard/benchguard.go, internal/benchguard/benchguard_test.go

## What to build

The follow-on guard refuses `cd X && bench gate` with the fixed sentence
`Run the Bench command without a shell follow-on.` and `operator=&&`. The
operator there comes before the Bench segment, so the sentence points at the
wrong side of the line. The reader is told to remove something that is not
there.

`Verdict` in `internal/benchguard/benchguard.go` learns which side the named
operator sits on. `operatorFor` already resolves the token by one precedence:
a redirection inside the span, else the control operator before the span,
else the one after it. Record that choice on the verdict as one field, and
let `Message` render one of three sentences from it.

- Operator after the span: the current sentence stays byte-for-byte.
- Operator before the span: `Run the Bench command from the current directory; it resolves the worktree itself.`
- Redirection inside the span: `Run the Bench command without a redirection.`

The prefix `BLOCKED: Bench response is bounded, complete, and self-contained.`
stays on every sentence, because `internal/systemtest/bench_follow_on_test.go`
and the operator's memory read it. The `segment=` and `operator=` suffix stays
unchanged. `BlockMessage()` keeps its current text for the follow-on case, so
any caller that still reads it sees no change.

## Acceptance

- [ ] `Classify("cd /tmp && bench gate", r).Message()` renders the current-directory sentence with `operator=&&`.
- [ ] `Classify("bench gate && echo done", r).Message()` renders the unchanged follow-on sentence with `operator=&&`.
- [ ] `Classify("bench gate 2>&1", r).Message()` renders the redirection sentence.
- [ ] Every rendered message starts with `BLOCKED: Bench response is bounded, complete, and self-contained.`.
- [ ] A new table test in `internal/benchguard/benchguard_test.go` covers the three sides, and the delegate records the leading-operator row red before the fix.
- [ ] `go test ./internal/benchguard/ ./internal/guards/ -parallel 2` passes.
