---
name: craft-tickets
description: How to break a spec or small change into independently-green tracer tickets with explicit blockers and one fresh write-delegate charge each. Use at build entry, when deriving tickets from stories and seams, when deciding what lands green next, or when a wide refactor needs an expand–contract sequence.
index: breaking a build into independently-green tickets
---

# Tickets: what lands green next

A ticket is the **smallest independently-green** vertical unit: a tightly
related story group that cuts a narrow but complete path through the layers
and can land committed on a green project gate by itself — a tracer bullet,
demoable or verifiable alone. A horizontal layer, tests without behavior, or
behavior without its tests is not a ticket. Its green integration-source commit is the grading rule.

## Draft the breakdown

Split until splitting further would leave no independently-green landing.
Keeping a group whole requires naming the specific red a thinner cut would
strand — review re-derives that claim by attempting the split. Feature
wholeness alone does not justify the grouping.

Name every real blocker by sibling ticket file basename; a ticket with all
blockers landed is on the **frontier**, and blockers order before consumers.
A wide mechanical refactor that would break every ordinary tracer ticket
instead sequences as expand (new form beside the old), migrate (move callers
in green batches), then contract (remove the old form once every migrate
ticket lands, `Blocked by:` naming them all).

**Reviewer-approved breakdown.** Before any ticket is assigned, the
coordinator presents the reviewer a numbered list — title, `Blocked by:`, and
delivered outcome — for every proposed ticket, iterates it with the
reviewer, and records approval. This is the only route onto the frontier;
the batch-approval AFK carve-out in `.bench/BENCH.md` is the sole
no-round-trip exception.

## Write one file per ticket

Write each ticket under `specs/<slug>/tickets/` with a verb-first title:

```markdown
# <Verb-first title>

Blocked by: <sibling ticket file basenames, or none>
Writes: <paths this ticket expects to touch, advisory>

## What to build

<The end-to-end behavior this ticket makes work.>

## Acceptance

- [ ] <observable behavioral criterion>
- [ ] <observable behavioral criterion>
```

`Blocked by:` is `none`, or the sibling basenames that must land first — a
basename survives a retitle and is what `--ticket` already names. `What to
build` states the end-to-end behavior, including any contract shared with a
sibling: a meaningful crossing lives in this prose and in `Acceptance`, never
a separate schema field, re-derived from the tree by review rather than
trusted from the ticket's account. `Acceptance` rows are observable
behavioral criteria, not a project-gate checkbox. `Writes:` is advisory only,
judging whether two frontier tickets are disjoint enough to parallelize —
never enforced.

Good:

<!-- ticket-example:begin -->
```markdown
# Render cancelled jobs in status

Blocked by: parse-cancelled-job-records.md
Writes: internal/status, internal/render/rows.go

## What to build

Users see a cancelled row, its reason, and the next recovery action. The
existing cancelled-status contract rejects a cancelled row missing either its
reason or its recovery action, so a thinner field-by-field cut strands that
contract red.

## Acceptance

- [ ] status renders the cancelled row with its reason.
- [ ] status renders the recovery action beside a cancelled row.
```
<!-- ticket-example:end -->

A verb-first, end-to-end outcome a fresh delegate can land alone.

## Land the frontier

Work the unblocked frontier. One ticket equals **one fresh write-delegate
charge**; independent frontier tickets run in parallel only where their
`Writes:` notes are disjoint, and dependent tickets run sequentially.

Run focused checks during the ticket, not a standalone full gate. For a
reviewed spec chain, commit tickets serially on one retained integration source,
one full-project gate per commit. Review freezes that source's base and tip;
`bench worktree land` composes and gates it, and final-check reports the evidence.
