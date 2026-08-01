# Close the ambient-state review gaps

Blocked by: Keep ambient git state local and capture-neutral

## What to build

The project profile states the ambient dashboard contract that FT113 and FT145
implemented, and regression coverage proves that excluding a sibling worktree or
the handoff capture never suppresses other dirt in the named checkout.

The former FT121 spec-status-flip policy remains excluded.

## Acceptance

- [x] The profile's exact capture-only allowlist includes `capture/session-handoff.md`.
- [x] The profile distinguishes checkout-local dirty paths from repository-wide
      unpushed-commit and unique-branch counts.
- [x] A named checkout and a sibling can both be dirty while `git.LandedState`
      reports exactly the named checkout's dirty path.
- [x] A tracked handoff rewrite plus another dirty path still reports that path
      and derives `commit on green`.
