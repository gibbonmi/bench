# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `73c97aa5`; integration worktree `bench worktree path d58e5fed7d1634c79cdeacadfa519a8c` on `bench/assign/…/d58e5fed7d1634c79cdeacadfa519a8c` at `36a74a59`, frozen review base `73c97aa5`
Spec: `specs/spec-ticket-fence-reduction/spec.md` (Status: staged)
Gate: green in the worktree at `36a74a59`; this file is the only uncommitted path there

## State

`/bench-implement-spec --full` is in the build phase, five of eight tickets landed
green on the integration source (last: the loop collapse at `36a74a59`). Remaining,
in `Blocked by:` order: `remake-craft-spec-and-craft-tickets-on-their-sources`, then
`realign-the-consumers-glossary-and-docs`. Each is one write delegate on **`opus`**/high
in the worktree (reviewer instruction: no `fable` for tickets); the coordinator probes
with a different mutation kind and site, then lands with `bench commit`. Review then
runs on `opus`/high over base `73c97aa5` → the source tip; accepted findings land as
repair tickets; from the destination `bench worktree land`, then `/bench-final-check`.

Reviewer decisions already closed this build, do not reopen: the
`internal/conformance/fixture_bite_test.go` fence addition (ratified); the 73-line
budget for `bench-write-spec.md` in place of the map's 60 (accepted, spec and ticket
amended); the narrower-capability falsification question stays in the command, pinned
by an out-of-fence conformance literal (follow-up, not this build). Bare
`go test ./internal/conformance/` skips the live-tree anchor check — trust `bench gate`.

## Next command

`/bench-implement-spec --full spec-ticket-fence-reduction --reviewer opus high`

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
