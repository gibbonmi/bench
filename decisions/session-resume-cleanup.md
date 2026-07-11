## Destination

Make a cold session resume cost nothing: provably-safe stale worktrees are cleaned
automatically at SessionStart, work intent survives session death via a
deterministic ledger, and one `bench status` verdict answers "is everything
committed and pushed?" — no agent re-investigation, no manual cleanup ritual.

## #1: Where does the cleanup fire?

Blocked by: none
Type: Grill

### Question

Auto at SessionStart, prevention at Stop, or a sharper manual prompt — where does
the burden actually disappear?

### Answer

Auto at SessionStart. The session-start hook runs the safe subset of
`bench worktree clean` before rendering `bench status`, so the dashboard reflects
post-clean state, and injects a one-line report of what was cleaned and what was
kept. The work is fully codified in the CLI — the hook adds no agent tokens beyond
the report line. Stop-time prevention alone was rejected: the failing case is the
session that dies uncleanly, where no end-of-session step ever runs.

## #2: What captures the handoff before the session goes cold?

Blocked by: #1
Type: Grill

### Question

When a session dies before its final verification, git state is recomputable at
resume but intent is lost — which worktree was doing what, and what step never ran.
What carries that across?

### Answer

A deterministic ledger, written by the tooling at the *start* of work (so it
survives unclean death) and refreshed by the Stop hook. It lives in the git dir
next to the gate cache (`bench-last-gate` precedent) — machine-written, never
tracked, never model-authored. At resume, `bench status` joins ledger intent with
live git state into decision-grade rows ("worktree X: '<objective>', landed —
cleaned" / "branch Y: unique commits, next: /bench-final-check"). Model-written
prose handoffs were rejected: they cost tokens every stop and are lost exactly in
the unclean-death case.

## #3: How aggressive is the unattended clean?

Blocked by: #1
Type: Grill

### Question

Manual `bench worktree clean` salvage-commits dirty WIP and removes the checkout.
May an unattended SessionStart hook do the same?

### Answer

No — conservative auto. The auto path removes only *clean* out-of-pool worktrees
and sweeps `worktree-*` branches proven landed in the default branch (ancestry or
patch containment, the existing `LandedInDefault` rule). Dirty, locked (live
agent), and leased worktrees are never touched unattended — they surface as ledger
rows with a next action. Salvage-commit stays exclusive to the manual
`bench worktree clean`. An unattended hook destroys nothing and commits nothing.

## #4: Which surfaces write intent into the ledger?

Blocked by: #2
Type: Grill

### Question

The pain-source worktrees (`.claude/worktrees/agent-*`) are created by the
harness's Agent tool — bench never sees them born. What are the ledger's writers?

### Answer

Three: `bench shift` records its objective (it already has one), `bench worktree`
gains an optional objective argument, and the existing `check-agent-line`
PreToolUse hook additionally records each Agent delegation's objective — it
already intercepts every Agent call in Claude Code and Codex, so the pain source
is covered with no new hook. Recording commits, spec associations, or gate
verdicts in the ledger was rejected: the ledger carries only what git cannot
(intent), per the one-source rule.

## #5: What does the resume surface do about committed/pushed?

Blocked by: #1
Type: Grill

### Question

Half the resume burden is verifying everything landed. Push is reviewer-owned by
guard posture — how far does the surface go?

### Answer

A consolidated landed-state verdict in `bench status` (the single renderer — no
new command): unpushed commit count, branches still holding unique commits, dirty
paths — one glance, one remaining manual action (the push). Auto-push was
rejected: it would reopen the closed guard decision that pushing is the
reviewer's act.

## #6: Which modules and envelopes own each fact?

Blocked by: #3, #4, #5
Type: Research

### Question

Trace the exact seams before the spec locks: (a) what the `check-agent-line`
PreToolUse envelope actually carries — can a delegation be correlated with the
`worktree-agent-*` branch/worktree it later creates, or only recorded as
uncorrelated intent; (b) whether the Stop-hook envelope supports the ledger
refresh; (c) which Go package owns the ledger (new vs `internal/worktree` vs
`internal/status`) and where the conservative-clean subset lives relative to
`CleanCommand`; (d) how the new status rows fit the five-row budget and severity
ladder contract; (e) gate attachment for hook-driven behavior, with the observable
red signal per seam. Cite each claim to the owning source; probe the hook
envelopes against the live harness.

### Answer

— (open)

## Not yet specified

- Ledger row lifecycle: when a row expires (on clean, on landed, after N days) —
  sharpens after #6 fixes the ledger's shape.
- Report wording and TOON shape of the injected clean report and verdict row.
- OpenCode / headless-shift parity for the ledger refresh (session-start.sh is
  already shared; the Stop side may differ per harness).

## Out of scope

- Auto-push or any hook-initiated push — push stays the reviewer's act.
- Model-written prose handoffs at session end.
- Salvage-committing WIP from the unattended path.
- Pool sizing or tightening changes — the warm pool is self-managing.
- Interactive confirmation flows in hooks — the auto path acts or surfaces, never asks.
