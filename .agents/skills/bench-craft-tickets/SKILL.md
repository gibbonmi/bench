---
name: craft-tickets
description: How to break a spec or small change into independently-green tracer tickets with explicit blockers and one fresh write-delegate charge each. Use at build entry, when deriving tickets from stories and seams, when deciding what lands green next, or when a wide refactor needs an expand–contract sequence.
index: breaking a build into independently-green tickets
---

# Tickets: what lands green next

A ticket is the **smallest independently-green** vertical unit: a tightly
related story group that cuts a narrow but complete path through the layers and
can land committed on a green project gate by itself. It is a tracer bullet,
demoable or verifiable on its own. A horizontal layer, tests without behavior,
or behavior without its tests is not a ticket.

Context fit is a split heuristic: when one fresh context cannot hold the unit,
split it further. It never grades the ticket. The green landing commit is the
grading rule.

## Classify before slicing

First decide whether the work is an ordinary build or a wide refactor.

A wide refactor is one mechanical change whose blast radius breaks many call
sites at once, so no ordinary tracer ticket can stay green. Sequence it as:

1. **Expand:** add the new form beside the old so current callers keep working.
2. **Migrate:** move callers in green batches cut by ownership and blast radius.
3. **Contract:** remove the old form after every migrate batch lands.

Every expand, migrate, and contract ticket must land green independently. If a
migrate batch cannot, the expansion or prefactor is incomplete, or the batch
is too wide; repair the preparation or split the batch before proceeding.

For an ordinary build, derive tickets from the spec's stories and seams. Before
the consuming tickets, prefactor any shared primitive that two or more tickets
need. One primitive gets one owning ticket; consuming tickets block on it.

`craft-spec` owns the spec-time **who-writes-where** fence. This skill owns the
build-time **what-lands-green-next** unit. Apply the fence by name; do not
restate or redraw it here.

## Draft the breakdown

For each candidate ticket:

1. Group the smallest set of related stories that produces end-to-end behavior
   at a declared seam.
2. Confirm the group is independently green. Split any group that needs an
   unrelated later ticket before the project gate can pass.
3. Name every real blocker by sibling ticket title. A ticket with all blockers
   done is on the **frontier**.
4. Order blockers before consumers and present the complete breakdown for the
   build's approval surface.

Work the unblocked frontier. One ticket equals **one write-delegate charge** in
a fresh context. Independent frontier tickets may run in parallel only where
the spec-time ownership fences permit it; dependent tickets run sequentially.
The coordinator retains only the parent spec and frontier, then resets the
author context at the next ticket.

Run focused seam checks during the ticket; do not run a standalone full gate
before landing. The path-scoped `bench commit` is the only per-ticket
full-project-gate boundary and commits only on green. If it goes red, repair
from that output and retry. The normal green path runs one full gate. The
ticket carries behavioral acceptance checkboxes, not a project-gate checkbox:
the green landing commit is the one source for that verdict.
`/bench-final-check` remains the final full gate over the composed feature.

## Write one file per ticket

Write each ticket under `specs/<slug>/tickets/` with a verb-first title and this
shape:

```markdown
# <Verb-first title>

Blocked by: <sibling ticket titles, or none>

## What to build

<The end-to-end behavior this ticket makes work.>

## Acceptance

- [ ] <Observable behavioral criterion>
- [ ] <Observable behavioral criterion>
```

Good:

```markdown
# Render cancelled jobs in status

Blocked by: Parse cancelled job records

## What to build

Users see a cancelled row, its reason, and the next recovery action.

## Acceptance

- [ ] Status renders the cancelled row and recovery action.
```

This is a verb-first, end-to-end outcome with an observable acceptance check.

Bad:

```markdown
# Status changes

Blocked by: none

## What to build

Update the parser, renderer, and tests.

## Acceptance

- [ ] The project gate is green.
```

This is a layer-by-layer implementation list, and its checkbox duplicates the
landing commit instead of naming behavior.
