# Let a chunk checkpoint wait on a background grade

Blocked by: 08-queue-and-supersede-background-grades.md
Writes: internal/gate/run_transaction.go, internal/gate/checkpoint.go, internal/gate/background_wait.go (new), internal/gate/background_wait_test.go (new), internal/gate/backgroundstate (new), CHANGELOG.md
Covers: BG57, BG58, BG59, BG60, BG61, BG62, BG63, BG64, BG65, BG66

## What to build

Today a chunk checkpoint looks up its key, then takes the checkout lock, and refuses when another run holds it (`internal/gate/run_transaction.go:78-104`).

When a chunk checkpoint finds no green, it looks for a live claim on its own key.
If it finds one, it waits until the claim ends, before it takes the checkout lock.
When the wait starts, the I/O class of the waited grade returns from idle to best-effort. Its nice value stays at 19, because Linux refuses an unprivileged lower nice value.
The wait counts against the checkpoint's own gate timeout.
After the wait, the checkpoint looks up its key again. A green answers it. Otherwise, it runs its own gate.

Before the checkpoint runs its own gate, it does three things:

- A queued grade of its own key leaves the queue.
- It stops the running grade of its worktree that holds another key, and it waits until that grade has exited.
- It drops the queued grade of its worktree that holds another key.

Two foreground runs in one checkout keep today's refusal. `run_transaction.go` has 27 lines of headroom, so put the wait in `background_wait.go`.

## Acceptance

- [ ] A chunk checkpoint that finds a live claim on its key waits, then reuses the green, with one oracle run in total.
- [ ] When the wait starts, the waited grade's I/O class is best-effort, and its nice value stays 19.
- [ ] A wait that passes the checkpoint's gate timeout ends the checkpoint with exit 124.
- [ ] A waiting checkpoint holds no checkout lock.
- [ ] A waited grade that ends red makes the checkpoint run its own oracle.
- [ ] A queued grade of the checkpoint's key leaves the queue when the checkpoint starts its own run.
- [ ] Before its own run, the checkpoint stops the other-key grade of its worktree, waits for its exit, and drops the other-key queued grade.
- [ ] A second foreground run in one checkout still refuses with `gate execution already in progress`.
