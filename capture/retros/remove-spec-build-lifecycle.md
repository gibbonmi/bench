# Retro — remove-spec-build-lifecycle

Shipped 2026-08-11: five serial commit-on-green tickets plus one
composed-review repair round, four write delegates and two read-only review
delegates, ~16k lines net deleted, first build run under the Pocock-alignment
cadence (no provisional lifecycle, direct-to-main, gate at the terminal
landing).

**What worked.** Serial ticket-at-a-time with a fresh delegate per ticket and
coordinator verification against the tree caught every done-claim gap the
same hour it was made. Both independent review passes paid for themselves
concretely: spec falsification caught the draft reconcile deleting the gate's
own `refs/bench/green` verdict store and a nonexistent sweep script cited as
"already covered"; the composed review caught two anchor rows deleted while
their prose survived and 350 lines of stranded dead chain. Red-first
discipline held everywhere it was asked (removed-verb sweep red on 68 lines;
RM10 guard red before the reconcile existed).

**What dragged.** Harness diagnostics repeatedly showed stale compiler errors
on deleted files after each delegate returned; only real `go build`/`go test`
runs settled the question — verify against the tree, never the diagnostic
panel. The ~20-symbol dead-chain sweep should have been in the deletion
ticket's charge (the composed review had to catch it): a wholesale package
deletion wants an explicit orphaned-caller sweep step.

**Open thread carried to the reviewer.** `worktree clean --apply` still
preserves into a namespace resume sweeps — refuse-dirty follow-up
recommended; three accepted deviations flagged in the handoff.
