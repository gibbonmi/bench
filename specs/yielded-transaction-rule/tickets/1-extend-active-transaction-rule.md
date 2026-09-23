# Extend the active-transaction rule to every Bench transaction

Blocked by: none
Writes: .agents/skills/bench-craft-delegate/references/delegation-discipline.md
Covers: none

## What to build

The delegation discipline says that only a running `bench commit` stays active until its process exits. Two drained learnings show the same fault with other transactions. A review venue update overlapped an active gate, and a focused run overlapped a yielded session. The rule names every Bench transaction and a yielded session, so a coordinator waits for the terminal result of each one.

## Acceptance

- [ ] The "Before the landing" rule names `bench commit`, `bench gate`, `bench worktree merge`, `bench worktree land`, and a focused test run as active until their process exits.
- [ ] The rule states that a yielded session stays active until it reports its terminal result.
