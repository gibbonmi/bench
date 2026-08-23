# Learnings — usage journal

<!-- entries below -->

## 2026-08-23 - appended a pipeline to a bench call  [open]
- **What happened:** The coordinator ran `bench commit ... 2>&1 | tail -12` during the `worktree-test-latency` build. The CLI contract forbids appended redirection and pipelines on `bench` calls.
- **Right behavior:** Run the bare `bench` command and read its complete output. Long output is CLI-owned projection work.
- **Proposed rule change:** None; the rule exists. This entry records the violation.
- **Repeat:** The coordinator appended `2>&1 | tail -3` to a second `bench commit` in the same session, after this entry existed. A habit-level fix is needed: the harness memory should carry a no-pipe-on-bench rule.

## 2026-08-23 - universal coverage claims need an inventory owner  [open]
- **What happened:** Iteration 1 of the `worktree-test-latency` spec review found three universal coverage rows without a complete inventory. They also lacked a concrete ledger owner.
- **Right behavior:** A universal row names one authoritative inventory and its enforcement seam. It also names the mutation that proves omission turns red.
- **Proposed rule change:** Add this universal-claim check to `craft-spec`'s acceptance coverage review rubric.
