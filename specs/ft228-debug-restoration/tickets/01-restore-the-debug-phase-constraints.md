# Restore the debug phase's upstream constraints

Blocked by: none
Writes: .agents/commands/bench-debug.md, .agents/skills/bench-debug/references/loop-constructions.md, CHANGELOG.md

## What to build

Edit `.agents/commands/bench-debug.md` in place, keeping its `## Entry
orientation` / `## Exit handoff` contract, every preserved upstream constraint,
and every Bench addition word-for-word, and restore what the compression lost
(upstream `reference-skill-repos/skills@d574778f`,
`skills/engineering/diagnosing-bugs/SKILL.md`, is the source text):

- Phase 1 points at a new reference file,
  `.agents/skills/bench-debug/references/loop-constructions.md`, carrying the
  ten loop constructions in upstream's try-in-order framing plus the structured
  human-in-the-loop last resort (described, no executable template). The
  five-construction inline sentence is replaced by the pointer.
- Phase 1 ends with the exact stop-gate `No red-capable command, no Phase 2`
  over the four completion criteria as `- [ ]` checkbox lines (red-capable,
  deterministic, fast, agent-runnable).
- `Tighten the loop` returns as a named step inside Phase 1: treat the loop as
  a product — faster, sharper, more deterministic.
- Phase 2 carries the three confirmations as `- [ ]` checkbox lines (the user's
  failure mode, a debuggable reproduction rate, the captured symptom), the
  exact stop-gate `Do not proceed until you have reproduced and minimised`,
  and the reproduction-economics rule: the sentence beginning
  `A green proxy only narrows a hypothesis`, requiring a load- or
  environment-sensitive failure to be reproduced through the accused command
  under exposing conditions before another stand-in is trusted.
- Phase 6's close-out returns as a `- [ ]` checkbox form; upstream's dead
  `/improve-codebase-architecture` pointer is not imported.

Each sentence ticket 02 will pin must occur exactly once in the file. One
CHANGELOG entry under Unreleased covers the restoration alone — the
invocation settle is documented by the ticket that lands it. The
contract between this ticket and ticket 02: the exact restored sentences above
are the anchor needles ticket 02 registers, re-derived from the tree by review.

## Acceptance

- [ ] Phase 1 names `.agents/skills/bench-debug/references/loop-constructions.md`, and that file lists all ten constructions plus the HITL last resort (DR1, DR2).
- [ ] The two stop-gate sentences appear verbatim, one in Phase 1, one in Phase 2 (DR3, DR4).
- [ ] The Phase 1 completion criteria and Phase 2 confirmations are `- [ ]` checkbox lines (DR5, DR6).
- [ ] `Tighten the loop` is a named step naming faster, sharper, and more deterministic (DR7).
- [ ] The `A green proxy only narrows a hypothesis` sentence is in Phase 2 (DR8).
- [ ] Every pre-existing bench-debug anchor still passes and the Forbid needle stays absent — the kit gate is green (DR9).
- [ ] The CHANGELOG carries the restoration entry.
