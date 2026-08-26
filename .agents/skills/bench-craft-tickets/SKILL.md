---
name: craft-tickets
description: How to break a spec or small change into tracer-bullet tickets — complete vertical slices, demoable alone, one context window each — with explicit blockers and, for spec-backed builds, one fresh write-delegate charge each. Use during spec authoring, when deriving tickets from stories and seams, when deciding what lands green next, or when a wide refactor needs an expand–contract sequence.
index: breaking a build into tracer-bullet tickets
---

# Tickets: what lands green next

Break a spec into **tracer-bullet** tickets. Each ticket cuts a narrow but COMPLETE
vertical path through every layer: schema, command, output, and tests. Each ticket is
demoable or verifiable on its own, and is sized to one fresh context window.
A horizontal layer, tests without behavior, or behavior without its tests is
not a ticket. A coverage row that only adds a test to a seam its parent slice
already opened is that slice's acceptance row. Its green integration-source
commit is the grading rule.

## Draft the breakdown

Gather context: the spec, or the conversation. Explore the codebase if you have not. Put any
prefactoring that makes the change easy first, as its own ticket. Then draft the vertical slices.

Name every real blocker by sibling ticket file basename. A ticket with all blockers
landed is on the **frontier**, and blockers order before consumers. A wide mechanical
refactor can break every ordinary tracer ticket. It instead
sequences as expand (new form beside the old), migrate (move callers
in green batches), then contract. Contract removes the old form once every migrate
ticket lands, `Blocked by:` naming them all.

**Reviewer-approved breakdown**: before the coordinator assigns a spec-backed ticket, it
presents the reviewer a numbered list — title, `Blocked by:`, and delivered outcome — for
every ticket. The coordinator asks: does the granularity feel right, too coarse or too fine,
and are the blocking edges correct? Does each ticket depend only on tickets that genuinely
gate it, and should any ticket be merged or split? Iterate until the reviewer approves, and
record approval.
For spec-backed builds, this is the only route onto the frontier; the batch-approval AFK carve-out
in `.bench/BENCH.md` is the sole no-round-trip exception. The light path is the exception:
`.bench/BENCH.md`'s right-size table is the one ticket's standing approval, and the main session implements it inline.

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

Write the prose in ASD-STE100 per `craft-spec`'s `references/ste-prose.md`.
`Blocked by:` is `none`, or the sibling basenames that must land first. A
basename survives a retitle and is what `--ticket` already names. `What to
build` states the end-to-end behavior. It also states any contract shared with a
sibling: a meaningful crossing lives in this prose and in `Acceptance`, never in a
separate schema field. Review re-derives the crossing from the tree rather than trusting
the ticket's account alone.

`Acceptance` rows are observable behavioral criteria, not a project-gate checkbox. Cite each
coverage row id in full, because preflight reads ids, not ranges. `Writes:` is advisory only. It
helps judge whether two frontier tickets are disjoint enough to parallelize, and nothing enforces
it. Each path it names exists in the tree, or the ticket marks the path new.

Good:

<!-- ticket-example:begin -->
```markdown
# Render cancelled jobs in status

Blocked by: parse-cancelled-job-records.md
Writes: internal/status, internal/render/rows.go

## What to build

Users see a cancelled row, its reason, and the next recovery action — one
demoable path from parsed record to rendered row, sized to a fresh context.

## Acceptance

- [ ] status renders the cancelled row with its reason.
- [ ] status renders the recovery action beside a cancelled row.
```
<!-- ticket-example:end -->

## Land the frontier

Spec-backed builds work the unblocked frontier. One ticket equals **one fresh write-delegate charge**. Frontier tickets run in parallel by default, one worktree each.
The coordinator serializes a pair only for a conflict it names in the charge. A conflict is an overlapping `Writes:` note, a contract one ticket reads from the other's diff, or the test budget. Dependent tickets run sequentially.

Run focused checks during the ticket, not a standalone full gate. For a reviewed spec chain, the coordinator folds each
returned diff onto one retained integration source in `Blocked by:` order, one gate per commit. Review freezes that source's
base and tip. `bench worktree land` composes and gates it, and final-check reports the evidence.
