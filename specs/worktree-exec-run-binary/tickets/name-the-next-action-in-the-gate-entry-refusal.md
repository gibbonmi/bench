# Name the next action in the gate entry refusal

Blocked by: none
Writes: .bench/gate.sh, internal/conformance/gate_entry_test.go

## What to build

The gate entry refuses an absent or relative run-binary value with one line
naming the variable. The variable is wrapper-owned: no operator sets it, nothing
in the message says what to run instead, and neither condition that produces it
— an invocation that skipped the wrapper, or a child whose caller stripped the
routing — is named. A session that reads it concludes the command it just ran is
unusable and stops using it. That is observed behavior, not a hypothesis.

Reword that one refusal in operator terms, naming the wrapper invocation as the
next action. Leave the three later refusals in the block as they are; each
already names a condition an operator can act on.

Two constraints bound the wording. The refusal is not specific to
`bench worktree exec` — any invocation reaching the gate entry without a wrapper
meets it, from any directory — so the message must name no worktree and read
correctly outside one. And the gate-entry contract test asserts the output
carries the run-binary variable's name, so that literal token stays in the
message or the assertion moves with it. The gate script never ships to
consumers, so kit-relative wording is safe.

This ticket and `own-the-run-binary-in-an-exec-child.md` are independent and
write disjoint files. That ticket removes the most common way to reach this
refusal; this one repairs the refusal for every other way. Neither blocks the
other, and they may be built concurrently.

Evidence comes from running the gate entry as a real subprocess under each
condition.

## Acceptance

- [ ] The gate entry with no run-binary variable refuses and names the wrapper invocation.
- [ ] The gate entry with a relative run-binary value refuses with that same message.
- [ ] The refusal names no worktree.
- [ ] The refusal text still contains the run-binary variable's name.
- [ ] The gate entry with a valid absolute executable proceeds exactly as it does today.
- [ ] The non-executable, symbolic-link, and uncleaned-path refusals keep their present wording.
- [ ] The gate-entry test reds when the refusal drops the wrapper invocation, and that red is demonstrated and recorded.
