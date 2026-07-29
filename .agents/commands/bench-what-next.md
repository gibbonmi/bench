---
description: Roadmap maintenance — reconcile ROADMAP.md against the tree, drain IDEAS.md, implementation retros, and the learnings journal into it, refresh the recommended sequence, and propose the whole pass as one batch diff for reviewer approval. The single exit for parked ideas, pending retros, and open learnings. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-what-next — reconcile the roadmap, drain the capture

## Entry orientation

This is the single roadmap-maintenance phase. `bench status` and `bench roadmap`
point here whenever `IDEAS.md` has parked lines, `.bench/learnings.md` has open
entries, or `.bench/retros/` has pending implementation retros. One run
reconciles the roadmap against the tree, drains all three capture sources,
refreshes the recommended sequence, and hands the reviewer one diff to approve.

At entry, invoke `bench roadmap --context` exactly once. Its successful schema-2
snapshot is the complete local evidence for every step below. If the query fails,
stop the phase and report its error; manual evidence reconstruction would create a
different, partial input and is not a fallback.

## Exit handoff

Close by reporting the reconcile verdicts (rows removed or reworded), the drained
idea count, each retro recommendation disposition, each journal verdict, and the
refreshed sequence — with judgment calls flagged for veto. On approval, commit
on green; the recommended next command is the top line of the refreshed
`## Recommended sequence`.

## 1. Reconcile first

Before draining anything, verify every `ROADMAP.md` row against the tree. When a
row's spec may have shipped, use `bench spec history <slug>` for the shipped-row
check. Shipped work is removed and stale wording is corrected. Row presence is
status, so this pass is the backstop for anything spec-retire missed; the
empty-state recommendation is only trustworthy if the roadmap is current. No
completion markers — history lives in git.

## 2. Drain the inbox

`IDEAS.md` is a pure inbox: every run empties to zero. Each parked line gets one
disposition in the roadmap — a new prioritized row, a merge into the row that
already covers it, or a drop as already-triaged — and the parked-pending-evidence
tier is a valid destination for items awaiting a real trigger. No line stays
parked in the inbox; partial drains would kill the empty-state trigger.

## 3. Drain implementation retros

The snapshot's `retros` bodies are the only retro evidence this run reads; do
not re-read the directory into a second, potentially different snapshot. For
every actionable recommendation in every body, record one explicit disposition:
merge into an existing roadmap row, a new roadmap row, a learning-or-rule
disposition, or an explicit dismissal with one line of why. When a
learning-or-rule disposition adds journal material, carry it immediately
through the journal verdict step below so the run does not create a fresh open
entry behind itself.

After every recommendation has a disposition, remove every drained
`.bench/retros/*.md` file in the same reviewer-approved batch. A partial retro
drain is not allowed: the pending count must reach zero, and no source file is
removed before its dispositions are present for review.

## 4. Verdict the journal

Read `.bench/learnings.md` itself, not just the open-entry count — a malformed
entry still needs a verdict. Every open entry gets a verdict in the batch diff:
work-shaped becomes a roadmap item; rule-shaped becomes a roadmap item whose
next action is the kit edit (built later under the `craft-synthesis` discipline,
gated as usual); the rest are dismissed with one line of why. Nothing leaves the
journal without the reviewer's approval of that diff.

A defect-shaped entry — one claiming a sanctioned command misbehaved — becomes
a roadmap row only after its red signal reproduces through the accused command
itself, invoked as the entry quotes it. A repro through a lookalike surface (a
raw `git add` standing in for `bench commit`) proves nothing about the accused
path: without the real repro, dismiss the entry as unreproduced, or re-park it
with the missing repro named as its graduation trigger.

## 5. Restructure the board

The drain adds rows; this step subtracts them, so the board does not accrete
near-duplicates that each pay a full pipeline and gate separately. Propose the
reducing moves in the same batch diff, never applied silently: merge rows that
edit the same owner surface into one batched row, collapse rows that are faces
of one missing primitive, fold a leftover clause into its parent row, and
group rows sharing one failure class under a theme header. While walking the
rows, classify each as a fix (a defect in existing behavior, with evidence), a
feature (new capability or guidance), or decision-only; report the
classification in the exit rather than writing it into the row grammar — it
exists to steer the sequence.

## 6. Refresh the sequence

Rewrite the `## Recommended sequence` section: two or three numbered lines, each
naming the item and the phase command to run. This is the format contract
`bench roadmap` extracts verbatim once all capture sources are empty — the CLI
does no judgment, so this section is where the judgment lands. Fixes lead: at
equal priority a reproduced defect outranks a feature, cheapest first, and a
feature tops the sequence only on explicit reviewer pricing, named in its
line.

## 7. Batch-propose, then commit on green

Draft the full pass — reconciled roadmap, emptied inbox, retro dispositions and
removals, journal verdicts including dismissals — as one uncommitted batch diff.
The diff includes the run's concise `CHANGELOG.md` entry only when the pass
changes notable user-facing behavior. The approved commit and Git history own
reconcile verdicts, dismissed learnings, and promotion evidence; do not mirror
that history in a second ledger. That diff is the verdict sheet: the reviewer
approves or adjusts it once, and there are no per-item interactive sign-offs. On
approval, commit on green. Never commit the drain without that approval; a
standing batch approval (the AGENTS.md rule) counts, with contestable calls
flagged for post-hoc veto.
