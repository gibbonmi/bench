# Learnings — usage journal

<!-- entries below -->

## 2026-08-23 - a debug fix ran in the main checkout, not an isolated worktree  [open]
- **What happened:** The FT169 debug phase wrote its repro loop and fix in the main checkout. `bench worktree land` requires a staged `spec.md` with an ownership fence, and a bug fix has none. So an isolated worktree would have had no verb to land through. The session was the only writer, so the tree stayed attributable.
- **Right behavior:** A bug fix that crosses three packages still wants its own worktree. The blocker is the verb, not the discipline: a spec-less landing would let the debug path isolate and land like a build.
- **Proposed rule change:** none here. The parked spec-less-landing idea removes the blocker; until it ships, the debug skill's isolation rule is satisfied by the sole-writer condition, stated in the close.

## 2026-08-23 - the landing's capture policy takes one side per file  [open]
- **What happened:** The FT169 fix composes a conflicted `capture/` path by policy: the source wins `capture/session-handoff.md`, the destination wins every other capture file. On a conflicted `capture/learnings.md`, the source session's appended entries are dropped in favor of the destination's.
- **Right behavior:** The handoff rule is right (the closing session's state wins). For the append-only journals, a union merge keeps both sides' entries. The row's text named the one-side rule, so that rule shipped; the union is a reviewer decision.
- **Proposed rule change:** the phase-owned-file merge-rules spec decides per file: one-side for the handoff, union for the journals, and a rule for `ROADMAP.md` and `roadmap/`.
