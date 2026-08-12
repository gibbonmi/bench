# Learnings — usage journal

## 2026-08-12 — delegate worktree released by clean instead of release [open]

What happened: after landing the first repair delegate's diff, the coordinator
removed its worktree with `bench worktree clean`, leaving the assignment
`recovered` in the intent ledger; every subsequent worktree create was refused
with "conflicts with its existing assignment" until `bench resume-clean`
reconciled it. Right behavior: a delegate assignment ends through the release
path (or clean followed immediately by `bench resume-clean`); the pool grants
one delegate worktree per session, so a lingering record blocks the next
delegate. Proposed rule change: none — `craft-delegate` could name the
clean-then-reconcile pair for harness-created delegate worktrees.

## 2026-08-12 — one-off "prospective authorization refused: inherited" [open]

What happened: one `bench commit` attempt (the second landing of a serial
repair run) reported `gate: red` / `prospective authorization refused:
inherited`; an immediately following direct `bench gate` on the identical tree
was green and the retried commit landed on the reused fresh verdict. Not
reproduced across seven later landings. Right behavior: on this refusal, run
`bench gate` directly and retry rather than treating the tree as red. Proposed
rule change: none yet — needs a repro before it is a defect ticket; the
verdict-class plumbing around "inherited" records is the suspect.
