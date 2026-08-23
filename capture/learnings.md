# Learnings — usage journal

<!-- entries below -->

## 2026-08-23 - universal coverage claims need an inventory owner  [open]
- **What happened:** Iteration 1 of the `worktree-test-latency` spec review found three universal coverage rows without a complete inventory. They also lacked a concrete ledger owner.
- **Right behavior:** A universal row names one authoritative inventory and its enforcement seam. It also names the mutation that proves omission turns red.
- **Proposed rule change:** Add this universal-claim check to `craft-spec`'s acceptance coverage review rubric.
