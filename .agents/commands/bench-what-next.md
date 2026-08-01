---
description: Roadmap maintenance — reconcile ROADMAP.md against the tree, drain capture/IDEAS.md, implementation retros, and the learnings journal into it, refresh the recommended sequence, and propose the whole pass as one batch diff for reviewer approval. The single exit for parked ideas, pending retros, and open learnings. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-what-next — reconcile the roadmap, drain the capture

## Entry orientation

This is the single roadmap-maintenance phase. `bench status` and `bench roadmap`
point here whenever `capture/IDEAS.md` has parked lines, `capture/learnings.md` has open
entries, or `capture/retros/` has pending implementation retros. One run
reconciles the roadmap against the tree, drains all three capture sources,
refreshes the recommended sequence, and hands the reviewer one diff to approve.

At entry, invoke `bench roadmap --context` exactly once. Its successful schema-3
snapshot is the complete local evidence for every step below. Accept only
`context.schema = 3`; every other schema stops the phase before any batch mutation.
Do not guess recurrence facts from an older schema. If the query fails, stop the
phase and report its error; manual evidence reconstruction would create a different,
partial input and is not a fallback.

Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that
snapshot. When it is false, stop before any batch mutation and report every
`occurrence_discrepancies` row together with the complete context snapshot. The
structural evidence remains visible for reviewer diagnosis; do not infer a ledger
or a sequence from partial sources.

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

## 2. Drain occurrence evidence

For every `pending` owner/incident pair in `capture_occurrences`, add its incident
key to that owner's `Occurrences:` line in `ROADMAP.md` before removing any source
unit for the pair. Show the owning-row edit in the batch before the corresponding
source removal. Every `already-recorded` source already has that key: remove its
source unit without adding another key. This procedure applies to ideas, retro
recommendations, and learning entries before their source-specific drain removes
them.

## 3. Drain the inbox

`capture/IDEAS.md` is a pure inbox: every run empties to zero. Each parked line gets one
disposition in the roadmap — a new prioritized row, a merge into the row that
already covers it, or a drop as already-triaged — and the parked-pending-evidence
tier is a valid destination for items awaiting a real trigger. No line stays
parked in the inbox; partial drains would kill the empty-state trigger.

## 4. Drain implementation retros

The snapshot's `retros` bodies are the only retro evidence this run reads; do
not re-read the directory into a second, potentially different snapshot. For
every actionable recommendation in every body, record one explicit disposition:
merge into an existing roadmap row, a new roadmap row, a learning-or-rule
disposition, or an explicit dismissal with one line of why. When a
learning-or-rule disposition adds journal material, carry it immediately
through the journal verdict step below so the run does not create a fresh open
entry behind itself.

After every recommendation has a disposition, remove every drained
`capture/retros/*.md` file in the same reviewer-approved batch. A partial retro
drain is not allowed: the pending count must reach zero, and no source file is
removed before its dispositions are present for review.

## 5. Verdict the journal

Read `capture/learnings.md` itself, not just the open-entry count — a malformed
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

## 6. Classify every run; restructure on request

While walking the rows, classify each as a fix (a defect in existing
behavior, with evidence), a feature (new capability or guidance), or
decision-only; report the classification in the exit rather than writing it
into the row grammar — it exists to steer the sequence. Every run classifies;
no invocation skips it.

The board-restructuring pass is opt-in: it runs only when the reviewer
invokes the phase with `--restructure`, because the whole-board pass is heavy
and most runs rightly stop at draining and updating the docs. When invoked,
propose the reducing moves in the same batch diff, never applied silently:
merge rows that edit the same owner surface into one batched row, collapse
rows that are faces of one missing primitive, fold a leftover clause into its
parent row, and group rows sharing one failure class under a theme header —
the board otherwise accretes near-duplicates that each pay a full pipeline
and gate separately. A default run that spots a restructure candidate names
it in the exit rather than applying it.

## 7. Refresh the sequence

Rewrite the `## Recommended sequence` section: two or three numbered lines, each
naming the item and the phase command to run. This is the format contract
`bench roadmap` extracts verbatim once all capture sources are empty — the CLI
does no judgment, so this section is where the judgment lands. Rank rows by
severity. Within an equal-severity class, choose actionable work over blocked work;
only when rows are equally actionable, apply literal dependencies, then explicit
reviewer pricing. Only when all four stronger inputs tie, rank by descending
occurrence count. When occurrence count also ties, apply the existing reproduced
defect-over-feature rule, then cheapest-first cost rule. No CLI command sorts or
rewrites `ROADMAP.md`; this is reviewed maintenance judgment, not global sorting.

## 8. Batch-propose, then commit on green

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
