---
description: Roadmap maintenance — reconcile ROADMAP.md against the tree, drain capture/IDEAS.md, implementation retros, and the learnings journal into it, refresh the recommended sequence, and propose the whole pass as one batch diff for reviewer approval. The single exit for parked ideas, pending retros, and open learnings. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-drain — reconcile the roadmap, drain the capture

## Entry orientation

This is the single roadmap-maintenance phase. `bench status` and `bench roadmap`
point here whenever `capture/IDEAS.md` has parked lines, `capture/learnings.md` has open
entries, or `capture/retros/` has pending implementation retros. One run
reconciles the roadmap against the tree, drains all three capture sources,
refreshes the recommended sequence, and hands the reviewer one diff to approve.

The "not a workflow phase" label above is scope, not an exemption. `main`
still takes writes only through a landing. Implement-now writers use their
own bench worktrees. Create the batch worktree only when the protocol below
allows it, and close it through `bench worktree land`.

At entry, invoke `bench roadmap --context` exactly once. Its successful schema-4
index is the complete local inventory for every step below. It covers every
roadmap row and capture unit, each capture path, every true body byte count,
and all cross-check blocks. Accept only `context.schema = 4`; every other
schema stops the phase before any batch mutation.

Do not guess recurrence facts from an older schema. The index proves what exists; targeted fetches
and named body reads prove content.

Fetch complete roadmap detail only for rows the reconcile touches with
`bench roadmap --context --row <ids>`. Read idea, learning, and retro bodies
from the paths the index names. If the query fails, stop the phase and report its
error. A targeted fetch or named body read that fails also stops the phase:
manual evidence reconstruction is a different, partial input, not a fallback.

Before beginning `## 1. Reconcile first`, read `context.sequence_trusted` from that
snapshot. When it is false, stop before any batch mutation. Report every
`occurrence_discrepancies` row together with the complete index snapshot. The
structural evidence remains visible for reviewer diagnosis. Do not infer a
ledger or a sequence from partial sources.

## Delegate the evidence

Use at most three read-only delegates after the snapshot passes its trust check.
Assign at most one delegate to each scope: roadmap reconcile, idea and journal mapping, and retro analysis.
Skip a scope when its indexed source set is empty. Charge every delegate under
`craft-delegate` and `craft-line`.

Read delegates edit nothing, take no new inventory, and use only the snapshot and its named paths.
Each read delegate returns these fields: proposed owner, classification, occurrence, evidence, and reviewer decision.

Resolve duplicate incidents and reviewer decisions before any batch writer starts.
Verify that the tree stayed unchanged. Keep ignored capture removal, the handoff,
verification, and landing with the coordinator. An implement-now writer may run
while other reads continue. Route each implement-now writer through
`craft-delegate` isolation and `craft-line` routing.

If an implement-now item exists, create the batch worktree only after every such item lands green on `main`.
If no implement-now item exists, create the batch worktree after all reads finish and the coordinator resolves duplicate incidents and reviewer decisions.
If tracked changes remain, one later write delegate authors the complete tracked batch.

## Exit handoff

Report the reconcile verdicts (rows removed or reworded), the drained idea
count, each retro recommendation disposition, each journal verdict, and the
refreshed sequence. Flag judgment calls for veto. Run `bench roadmap --flow`
once and quote its flow block in the exit. On approval, commit on green,
once, over everything the pass touched. The recommended next command is the
top line of the refreshed `## Recommended sequence`.

## 1. Reconcile first

The board is split, so every edit below has one owner.
`ROADMAP.md` holds the priority ordering and one heading line per row.
A row's body, its `Occurrence:` ledger, and its `Sources:` line are edited in `roadmap/FT<n>.md`.
Retiring a row deletes its `ROADMAP.md` heading line and its `roadmap/FT<n>.md` together in the same batch.

Before you drain anything, verify every `ROADMAP.md` row against the tree. When a
row's spec may have shipped, use `bench spec history <slug>` for the shipped-row
check. Remove shipped work and correct stale wording. Row presence is
status, so this pass is the backstop for anything spec-retire missed. The
empty-state recommendation is only trustworthy if the roadmap is current. Write
no completion markers; history lives in git.

Retire here a spec whose work has landed but whose directory still sits under
`specs/`; do not leave it for a later invocation. Run
`bench spec retire <slug>` during this pass so its deletions join the batch
below. Promote whatever of the spec stays durable onto its roadmap row first.
Leave that row with no spec path named; the row survives only as the residual
work the spec did not ship.

## 2. Drain occurrence evidence

An entry feeds a row only when it changes the row's priority, scope, or `Next:`.
Dismiss an occurrence-only entry with one line of why.

### Normalize touched rows

Keep current-state core prose concise. Each drained historical incident, build,
retro, or learning merged into a row becomes exactly one physical
`Occurrence: <when/source> — <short situation>.` line. Occurrence lines contain
event evidence only. New feature faces and decisions remain concise core prose.
Normalize every touched row before batch proposal.

**Good — event-only evidence:** `Occurrence: 2026-08-15 gate build — primary-checkout preflight failed on a stale base.`
**Bad — remedy derivation:** `Occurrence: 2026-08-15 gate build — change preflight selection to fix the stale base.`

For every `pending` owner/incident pair in `capture_occurrences`, add its incident
key to that owner's `Occurrences:` line in `ROADMAP.md` before removing any source
unit. Show the owning-row edit in the batch before the corresponding
source removal. Every `already-recorded` source already has that key: remove its
source unit without adding another key. This procedure applies to ideas, retro
recommendations, and learning entries before their source-specific drain removes
them.

## 3. Drain the inbox

`capture/IDEAS.md` is a pure inbox: every run empties to zero. Each parked line gets one
disposition in the roadmap: a new prioritized row, a merge into the row that
already covers it, or a drop as already-triaged. The parked-pending-evidence
tier is a valid destination for items awaiting a real trigger. No line stays
parked in the inbox. A partial drain would kill the empty-state trigger.

A new row needs a `Next:` token and a class before it opens. Write the token as
the row's `Next:` line in `roadmap/FT<n>.md`. Step 6 gives the class. Each
token maps to one phase:

| token | phase |
|---|---|
| `shape` | `/bench-shape-idea` |
| `spec` | `/bench-write-spec` |
| `ticket` | the light-path ticket |
| `decide` | a reviewer decision |
| `kit-edit` | a `craft-synthesis` kit edit |

## 4. Drain implementation retros

The index is the sole retro inventory this run reads. Do not re-enumerate
`capture/retros/` into a second, potentially different listing. The run reads
each body at the path its index row names; that is required detail retrieval,
not another inventory. Give every actionable recommendation in every body one
explicit disposition. The dispositions are a merge into an existing roadmap
row, a new roadmap row, a learning-or-rule disposition, or an explicit
dismissal. Give the dismissal one line of why.

When a learning-or-rule
disposition adds journal material, carry it immediately through the journal
verdict step below. This way the run does not create a fresh open entry
behind itself.

When drained bodies carry repair-attribution tables, also report tickets
total, one-shots, and per-cause counts, reading causes only from the drained
tables. The cause vocabulary belongs to the retro template
`/bench-final-check` owns, so this run neither restates nor extends it. The
tally still counts a term it does not recognize, as written. The tally only
reports; it adds no roadmap row grammar of its own. A body with no
attribution table simply contributes nothing to it.

After every recommendation has a disposition, remove every drained
`capture/retros/*.md` file in the same reviewer-approved batch. A partial
retro drain is not allowed. The pending count must reach zero, and remove no
source file before its dispositions are present for review.

## 5. Verdict the journal

Read `capture/learnings.md` itself, not just the open-entry count. A malformed
entry still needs a verdict. Every open entry gets a verdict in the batch
diff. A work-shaped entry becomes a roadmap item. A rule-shaped entry becomes
a roadmap item whose next action is the kit edit, built later under the
`craft-synthesis` discipline and gated as usual. Dismiss the rest with one
line of why.

Nothing leaves the journal without the reviewer's approval of
that diff.

A defect-shaped entry claims a sanctioned command misbehaved. It becomes a
roadmap row only after its red signal reproduces through the accused command
itself, invoked as the entry quotes it. A repro through a lookalike surface —
a raw `git add` standing in for `bench commit` — proves nothing about the
accused path. Without the real repro, dismiss the entry as unreproduced or
re-park it. A re-parked entry names the missing repro as its graduation trigger.

For a drained item that meets the light-path observables, build the item in this session ("implement now") by default.
Write its one ticket file. Spawn a write-delegate charged with
that ticket under `craft-delegate` isolation and `craft-line` routing. Then verify
the returned diff in the main session against the ticket's acceptance rows and
the gate. Open a `ROADMAP.md` row only when the reviewer declines.
Items needing a reviewer decision, a new seam, or spec-level design
still graduate to `ROADMAP.md`.

## 6. Classify every run; restructure on request

While you walk the rows, classify each row. Use fix (a defect in existing
behavior, with evidence), feature (new capability or guidance), or
decision-only. Report the classification in the exit rather than write it
into the row grammar. It exists to steer the sequence. Every run classifies;
no invocation skips it.

The board-restructuring pass is opt-in. It runs only when the reviewer
invokes the phase with `--restructure`, because the whole-board pass is
heavy. Most runs rightly stop once they drain and update the docs.

When
invoked, propose the reducing moves in the same batch diff; never apply them
silently. Merge rows that edit the same owner surface into one batched row.
Collapse rows that are faces of one missing primitive. Fold a leftover clause
into its parent row. Group rows sharing one failure class under a theme
header. Otherwise the board accretes near-duplicates that each pay a full
pipeline and gate separately.

A default run that spots a restructure
candidate names it in the exit. It does not apply the change. When the flow
report shows a positive net delta, propose reducing moves in the next batch
diff. That obligation does not wait for `--restructure`.

## 7. Refresh the sequence

Rewrite the `## Recommended sequence` section: two or three numbered lines, each
naming the item and the phase command to run. This is the format contract
`bench roadmap` extracts verbatim once all capture sources are empty. The CLI
does no judgment, so this section is where the judgment lands. Rank rows by
severity. Within an equal-severity class, choose actionable work over blocked
work.

Only when rows are equally actionable, apply literal dependencies, then explicit
reviewer pricing. Only when all four stronger inputs tie, rank by descending
occurrence count. When occurrence count also ties, apply the
existing reproduced defect-over-feature rule, then cheapest-first cost rule.
No CLI command sorts or rewrites `ROADMAP.md`. This is reviewed maintenance
judgment, not global sorting.

## 8. Batch-propose, then commit once on green

Follow `## Delegate the evidence`. If tracked changes remain, charge the one
later write delegate to draft the complete tracked pass as one uncommitted batch diff.
If no tracked changes
remain, start no batch writer.

The reviewer batch contains the tracked diff, proposed ignored-source removals,
and every journal verdict. The tracked diff contains roadmap dispositions, retro removals, earned `bench spec retire` work, and provider scorecards.
`/bench-final-check` refreshes those scorecards with landing evidence. Unlike
per-spec retros, scorecards persist after this batch.

Everything tracked in the pass lands in one diff and one commit. Ignored local
changes do not enter that diff or commit. The gate is what a commit costs, and
this pass is bookkeeping. A split of other tracked work buys nothing and pays
the oracle twice. The later required per-spec commit exceptions still apply.
Leave no tracked part for a later run.

The diff includes the run's concise `CHANGELOG.md` entry only when the pass
changes notable user-facing behavior. The approved commit and Git history own
reconcile verdicts, dismissed learnings, and promotion evidence. Do not
mirror that history in a second ledger. That diff is the verdict sheet: the
reviewer approves or adjusts it once, and there are no per-item interactive
sign-offs.

On approval, commit on green. Never commit the drain without that
approval. A standing batch approval (the AGENTS.md rule) counts, with
contestable calls flagged for post-hoc veto.

Three constraints shape the drain's commits.
An item completed through "implement now" lands as its own commit on green before the drain's batch commit.
This is the second exception to the one-batch-commit rule.
The other exception is the extra per-spec commits required when one pass retires two or more specs.

When the
pass retires a spec, the commit subject **ends with** `spec-retire: <slug>`.
That suffix is the exact grammar `bench spec history` matches. A subject that
merely mentions the slug loses the retirement evidence a later run's
shipped-row check reads. The suffix carries one slug, so a pass retiring two
or more specs takes one extra commit per additional slug. The exit says so.
That is rare, and correctness of the history query outranks the saved
gate.

After approval, the coordinator empties ignored inbox and journal sources. It
writes ignored `capture/session-handoff.md` last. Its pin block names the
pre-commit HEAD by construction, which is correct rather than stale. `bench
status` dates the ignored handoff by its write time. The tree wins wherever the
two disagree.
