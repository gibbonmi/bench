# Print the broker notice after the refresh effect

Blocked by: none
Writes: internal/worktree/land.go, internal/worktree/land_effects.go, internal/worktree/land_resume.go, internal/worktree/land_effects_test.go, internal/worktree/land_broker_notice_test.go, internal/worktree/parallel_census_test.go, .agents/commands/bench-final-check.md
Covers: none

## What to build

A landing whose reviewed diff changes a promotion broker build input prints a
broker notice. The notice names the manual install step. At the moment, the
landing prints the notice before it publishes. The refresh effect runs later.
In the kit source checkout, a complete refresh already republishes the broker
executable, its seal, and the broker manifest. A coordinator then reads the
printed rebuild step as a step that the landing still owes.

The landing finds out before the release whether the diff changes a broker
build input, because the release removes the source worktree. The landing
prints the notice after the effects row.

In the kit source checkout, the refresh effect is the manual step. The notice
names the stamped rebuild and `bench doctor --fix` only when the refresh
reports `failed`. Off the kit source checkout, the refresh does not replace the
installed broker. The notice names `bench repair` or the release install for
every refresh result. A resumed landing has no source worktree, so it prints no
notice.

The final-check command states which diff changes the broker source. It does
not repeat when the landing prints the notice.

## Acceptance

- [ ] A broker-changing landing in the kit source checkout whose refresh reports `complete` prints no manual rebuild step.
- [ ] After a `failed` refresh in the kit source checkout, the notice names the stamped rebuild and `bench doctor --fix`. A test asserts that the notice follows the refresh diagnostics.
- [ ] A broker-changing landing off the kit source checkout names `bench repair` or the release install after a `complete` refresh.
- [ ] The `internal/worktree` package suite stays green.
