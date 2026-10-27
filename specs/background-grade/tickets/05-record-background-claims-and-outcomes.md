# Record background claims and outcomes

Blocked by: 01-measure-the-background-census.md
Writes: internal/gate/backgroundstate (new), internal/subprocess, CHANGELOG.md
Covers: BG32, BG33, BG34, BG35, BG112, BG113

## What to build

Add the leaf package `internal/gate/backgroundstate`. It imports no `internal/gate` code, so `internal/gocache`, `internal/worktree`, `internal/status`, and `internal/commit` can import it.

The package owns three records:

- A claim names a grade key, a worktree, an owner pid, and the boot id. A claim is live only while its owner is alive and its boot id is the current boot id.
- A background record of a worktree names the latest outcome, the grade key, the tip, and the log path. The outcome is `green`, `red`, `aborted`, or `stopped`.
- The machine queue and its slot hold the queued grades and the one running grade. Ticket 09 gives the queue its rules.

An absent record reads as empty. A present record that does not read retains its grade and starts none, as ADR 0030 retains a resource.
The author picks the storage places, with the claims beside the evidence store and the queue under the Bench home.
Write every record through the durable replacement that the gate store uses.

The package also gives one stop operation. It asks the owner of a claim to end, and it returns only after that owner has exited.
The `process-lifetime-closure` check refuses a signal that code outside `internal/subprocess` sends (`internal/conformance/process_lifetime_closure_test.go:55-62`).
So the stop sends its signal through an operation of `internal/subprocess`, and this ticket adds that operation there if none fits.

## Acceptance

- [ ] A written claim names its key, its worktree, its owner pid, and the boot id.
- [ ] A claim whose owner pid is not alive is not live.
- [ ] A claim whose boot id differs from the current boot is not live.
- [ ] A background record names the outcome, the grade key, the tip, and the log path.
- [ ] The stop operation returns only after the owner process has exited.
- [ ] No code in `internal/gate/backgroundstate` sends a signal itself, so the lifetime closure check finds no signal outside `internal/subprocess`.
- [ ] `go list -deps ./internal/gate/backgroundstate` names no `internal/gate` package other than itself.
