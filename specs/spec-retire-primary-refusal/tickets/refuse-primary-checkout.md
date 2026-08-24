# Refuse `bench spec retire` on the primary checkout

Blocked by: none
Writes: internal/spec/spec.go, internal/spec/spec_test.go

## What to build

`bench spec retire` refuses the primary checkout before any deletion, the
same way `bench commit` does. The check reuses `git.IsPrimaryCheckout` and
prints `usage.PrimaryCheckoutRefusal()`. A retire from a linked Bench
worktree keeps its current behavior. When the checkout identity is unknown,
retire refuses without a deletion. A retire outside a Git repository keeps
its current behavior.

## Acceptance

- [ ] On the primary checkout, `bench spec retire <slug>` exits 1, prints the shared primary-checkout refusal, and deletes nothing.
- [ ] The refusal string has one source: the test proves the retire path emits `usage.PrimaryCheckoutRefusal()`, and no second copy of the string appears.
- [ ] From a linked worktree, a retire of a merged-implemented spec still deletes and exits 0.
- [ ] A mutation that removes the new check turns a named test red; record the probe and restore production exactly.
