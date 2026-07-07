# Kit Rule Edits (FT35, FT36)

## #1: What exactly changes in craft-delegate? (FT35)

Blocked by: —
Type: Grill

### Question
The 2026-07-06 learning: Agent-tool worktrees were repeatedly cut several
commits behind main — one delegate stalled, one was permission-denied its own
fast-forward, one built against the stale base. The adopted mid-run fix held.
What is the rule text?

### Answer
`craft-delegate`'s charge section gains the stale-base opener as a standing
requirement for worktree-isolated delegates: the charge opens with "run
`git merge --ff-only main`, verify HEAD equals main, stop and report if the
merge is denied or diverges" — and the orchestrator's side of the contract is
stated too: a blocked worktree is fast-forwarded by the orchestrating
session, which then resumes the same delegate. Read-only delegates are
unaffected. The skill's good-charge example gains the opener line so the
template teaches it.

## #2: What exactly changes in /bench-write-spec? (FT36)

Blocked by: —
Type: Grill

### Question
The 2026-07-06 learning: the entry contract says refuse without a closed map,
but an explicit reviewer batch drain produced ten map-less specs under a
defensible override that had to be reconciled ad hoc. What is the contract
clause?

### Answer
The entry contract gains one sentence: an explicit reviewer-directed batch
drain (an assessment or reviewed findings doc into specs) may substitute for
per-spec maps, with every defaulted decision flagged in-spec for post-hoc
veto; absent that explicit instruction, the map gate stands. This records the
override path without weakening the default.

## Handoff

1. **Module boundaries.** Two files: `bench-craft-delegate/SKILL.md` and
   `bench-write-spec.md`. Rule text only; no code, no index changes
   (descriptions unchanged).
2. **Contracts.** Both edits add-only: no existing rule is weakened; the
   delegate opener applies to worktree-isolated charges only; the write-spec
   clause requires the reviewer's explicit instruction.
3. **Deep vs thin.** Thin — the decisions were made in the drained learnings;
   this writes them into the owning artifacts under craft-synthesis.
4. **Black-box assertables.** Text presence; canary fixtures embedding either
   file's text refreshed in the same diff if touched.
5. **Gate attachment.** Docs conformance scans; skills-index `--check`
   unaffected.
6. **Hostile-input owners.** n/a — prose.
7. **Uncertainty flags.** n/a.
8. **Rejected alternatives.** A permission-layer change so delegates can
   fast-forward themselves (wider write surface for a narrow need); leaving
   the override unwritten (the ad-hoc reconciliation repeats).
9. **Domain watch-outs.** Leverage artifacts — every delegation and every
   spec session loads these; wording stays as terse as the surrounding text.

Dependency order: n/a — single spec, two independent one-file edits.
