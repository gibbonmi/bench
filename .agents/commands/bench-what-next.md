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

At entry, invoke `bench roadmap --context` exactly once. Its successful schema-4
index is the complete local inventory for every step below: every roadmap row and
capture unit, each capture path, every true body byte count, and all cross-check
blocks. Accept only `context.schema = 4`; every other schema stops the phase before
any batch mutation. Do not guess recurrence facts from an older schema. The index
proves what exists; targeted fetches and named body reads prove content.

Fetch complete roadmap detail only for rows the reconcile touches with
`bench roadmap --context --row <ids>`. Read idea, learning, and retro bodies
from the paths the index names. If the query fails, stop the phase and report its
error. A targeted fetch or named body read that fails also stops the phase; manual
evidence reconstruction would create a different, partial input and is not a
fallback.

Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that
snapshot. When it is false, stop before any batch mutation and report every
`occurrence_discrepancies` row together with the complete index snapshot. The
structural evidence remains visible for reviewer diagnosis; do not infer a ledger
or a sequence from partial sources.

## Exit handoff

Close by reporting the reconcile verdicts (rows removed or reworded), the drained
idea count, each retro recommendation disposition, each journal verdict, and the
refreshed sequence — with judgment calls flagged for veto. On approval, commit
on green, once, over everything the pass touched; the recommended next command
is the top line of the refreshed `## Recommended sequence`.

## 1. Reconcile first

Before draining anything, verify every `ROADMAP.md` row against the tree. When a
row's spec may have shipped, use `bench spec history <slug>` for the shipped-row
check. Shipped work is removed and stale wording is corrected. Row presence is
status, so this pass is the backstop for anything spec-retire missed; the
empty-state recommendation is only trustworthy if the roadmap is current. No
completion markers — history lives in git.

A spec whose work has landed but whose directory still sits under `specs/` is
retired here rather than left for a later invocation: run
`bench spec retire <slug>` during this pass so its deletions join the batch
below. Promote whatever of the spec stays durable onto
its roadmap row first, and leave that row naming no spec path — the row survives
only as the residual work the spec did not ship.

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

The index is the sole retro inventory this run reads; do not re-enumerate
`capture/retros/` into a second, potentially different listing. Reading each body
at the path its index row names is required detail retrieval, not another inventory.
For every actionable recommendation in every body, record one explicit disposition:
merge into an existing roadmap row, a new roadmap row, a learning-or-rule
disposition, or an explicit dismissal with one line of why. When a
learning-or-rule disposition adds journal material, carry it immediately
through the journal verdict step below so the run does not create a fresh open
entry behind itself.

When the drained bodies carry repair-attribution tables, the exit also reports
the running one-shot tally: report tickets total, one-shots, and per-cause
counts, reading causes only from the drained tables. The cause vocabulary
belongs to the retro template `/bench-final-check` owns, so this run neither
restates nor extends it — a term it does not recognize is still counted as
written. The tally is reporting only: it adds no roadmap row grammar of its own,
and a body with no attribution table simply contributes nothing to it.

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

For a drained item that meets the light-path observables, offer the reviewer
"implement now" instead of a `ROADMAP.md` row: write its one ticket file, spawn a
write-delegate charged with that ticket under `craft-delegate` isolation and
`craft-line` routing, then verify the returned diff in the main session against
the ticket's acceptance rows and the gate. Items needing a reviewer decision, a
new seam, or spec-level design still graduate to `ROADMAP.md`.

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

## 8. Batch-propose, then commit once on green

Draft the full pass as one uncommitted batch diff: the reconciled roadmap, the
emptied inbox, retro dispositions and removals, journal verdicts including
dismissals, every `bench spec retire` the reconcile earned, and the rewritten
`capture/session-handoff.md`. Everything the pass touches lands in that one
diff and one commit. The gate is what a commit costs, and this pass is
bookkeeping — splitting it across a drain commit, a retire commit, and a
handoff commit buys nothing and pays the oracle three times. Leave no part of
the pass for a follow-up invocation to commit.

The diff includes the run's concise `CHANGELOG.md` entry only when the pass
changes notable user-facing behavior. The approved commit and Git history own
reconcile verdicts, dismissed learnings, and promotion evidence; do not mirror
that history in a second ledger. That diff is the verdict sheet: the reviewer
approves or adjusts it once, and there are no per-item interactive sign-offs. On
approval, commit on green. Never commit the drain without that approval; a
standing batch approval (the AGENTS.md rule) counts, with contestable calls
flagged for post-hoc veto.

Three constraints shape the drain's commits. An item completed through
"implement now" lands as its own commit on green before the drain's batch commit;
this is the second exception to the one-batch-commit rule, beside the extra
per-spec commits required when one pass retires two or more specs. When the pass
retires a spec, the commit subject **ends with** `spec-retire: <slug>` — that suffix
is the exact grammar `bench spec history` matches, so a subject that merely
mentions the slug loses the retirement evidence a later run's shipped-row check
reads. The suffix carries one slug, so a pass retiring two or more specs takes
one extra commit per additional slug and says so in the exit; that is rare, and
correctness of the history query outranks the saved gate. And the handoff is
written last, immediately before the commit: its pin block names the pre-commit HEAD by
construction, which is correct rather than stale — `bench status` dates the
handoff by the commit that wrote it, and the tree wins wherever the two
disagree.
