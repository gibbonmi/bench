# Keep ambient git state local and capture-neutral

Blocked by: none

## What to build

The ambient dashboard and session handoff report actionable state for the current
checkout. A sanctioned handoff capture does not stale the gate or change the next
action it derives, and dirt in another registered worktree does not inflate or hide
the current checkout's dirty-path count. Repository-wide unpushed-commit and unique
branch facts remain repository-wide.

The former FT121 spec-status-flip policy is excluded because `ROADMAP.md` leaves its
remedy to a reviewer decision.

## Acceptance

- [x] A tracked `capture/session-handoff.md` is capture-only for gate-staleness projection.
- [x] Rewriting a tracked handoff excludes that file from both its dirty-path count
      and its derived next action, with one status-query exclusion mechanism rather
      than handoff-only count arithmetic.
- [x] `git.LandedState` counts dirty paths only in the checkout named by its root
      argument while retaining repository-wide unpushed-commit and unique-branch
      counts.
- [x] A dirty or unreadable sibling worktree cannot inflate or suppress the current
      checkout's ambient git signal.
- [x] The user-visible behavior is recorded once in `CHANGELOG.md`.
