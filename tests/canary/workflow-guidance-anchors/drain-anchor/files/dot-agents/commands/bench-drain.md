---
description: Roadmap maintenance — reconcile ROADMAP.md against the tree, drain capture/IDEAS.md and the learnings journal into it, refresh the recommended sequence, and propose the whole pass as one batch diff for reviewer approval. The single exit for parked ideas and open learnings. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-drain — reconcile the roadmap, drain the capture

## Entry orientation

This is the single roadmap-maintenance phase. `bench status` and `bench roadmap`
point here whenever `capture/IDEAS.md` has parked lines or `capture/learnings.md` has open
entries. One run reconciles the roadmap against the tree, drains both capture
sources, refreshes the recommended sequence, and hands the reviewer one diff to
approve.

## Exit handoff

Close by reporting the reconcile verdicts (rows removed or reworded), the drained
idea count, each journal verdict, and the refreshed sequence — with judgment
calls flagged for veto. On approval, commit on green; the recommended next
command is the top line of the refreshed `## Recommended sequence`.

## 1. Reconcile the roadmap

Before draining anything, verify every `ROADMAP.md` row against the tree —
shipped work is removed, stale wording is corrected. Row presence is status, so
this pass is the backstop for anything spec-retire missed; the empty-state
recommendation is only trustworthy if the roadmap is current. No completion
markers — history lives in git.

## 2. Drain the inbox

`capture/IDEAS.md` is a pure inbox: every run empties to zero. Each parked line gets one
disposition in the roadmap — a new prioritized row, a merge into the row that
already covers it, or a drop as already-triaged — and the parked-pending-evidence
tier is a valid destination for items awaiting a real trigger. No line stays
parked in the inbox; partial drains would kill the empty-state trigger.

## 3. Verdict the journal

Read `capture/learnings.md` itself, not just the open-entry count — a malformed
entry still needs a verdict. Every open entry gets a verdict in the batch diff:
work-shaped becomes a roadmap item; rule-shaped becomes a roadmap item whose
next action is the kit edit (built later under the `craft-synthesis` discipline,
gated as usual); the rest are dismissed with one line of why. Nothing leaves the
journal without the reviewer's approval of that diff.

## 4. Refresh the sequence

Rewrite the `## Recommended sequence` section: two or three numbered lines, each
naming the item and the phase command to run. This is the format contract
`bench roadmap` extracts verbatim once both capture sources are empty — the CLI
does no judgment, so this section is where the judgment lands.

## 5. Batch-propose, then commit on green

Draft the full pass — reconciled roadmap, emptied inbox, journal verdicts
including dismissals — as one uncommitted batch diff. That diff is the verdict
sheet: the reviewer approves or adjusts it once, and there are no per-item
interactive sign-offs. On approval, commit on green. Never commit the drain
without that approval; a standing batch approval (the AGENTS.md rule) counts,
with contestable calls flagged for post-hoc veto.
