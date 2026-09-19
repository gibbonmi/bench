# 6. Guard the live checkout in the gate's kit phases

Blocked by: 1-open-kit-test-run.md, 4-run-build-scripts-on-kit-copy.md
Writes: internal/gate/checkout_guard.go (new), internal/gate/checkout_guard_test.go (new), internal/gate/phases.go, internal/gate/phases_test.go
Covers: TD23, TD24, TD25, TD26, TD27, TD28, TD29, TD30, TD31, TD32

## What to build

Chunk: TD-C2.

The gate's phase command reads the checkout state of a kit root before the kit phases start and after they end. The state is each untracked and ignored path with its size and modification time. It also includes each Bench-owned file in the root's git directory with its size and modification time. A difference outside the declared set makes the gate red. One failure row names each changed path, in byte order.

The declared set holds the conformance timing file and the gate's own run log, run stream, lock, and owner record. Each declared path comes from the function that owns it.

A root whose checkout state cannot be read makes the gate red with one row that names the unreadable state. A linked root gets no guard. `fixturePhaseRoot` gains a git directory, because the guard reds a root with none.

Ticket 4 moved the one writer that the census confirmed. The census compared path lists only, so the first gate run with the guard is the complete census. If that run finds another undeclared writer, move its write under its own temporary directory. Add the writer's file to this ticket's `Writes:` line before the fix, as the operating guide requires for an in-scope expansion. A new declared path returns to the reviewer.

## Acceptance

- [ ] A kit fixture phase that adds a path, changes an ignored file, or removes an untracked file makes the gate red. One row names the path.
- [ ] A kit fixture phase that adds a Bench-owned file to the git directory makes the gate red with a row that names the file.
- [ ] A kit fixture phase that writes the timing file, or that runs under an open gate run log, leaves the gate green.
- [ ] Two changed paths print in byte order.
- [ ] An unreadable checkout state makes the gate red with one row that names it.
- [ ] A linked-root fixture phase that adds a path leaves the gate green.
- [ ] The whole gate is green with the guard on.
