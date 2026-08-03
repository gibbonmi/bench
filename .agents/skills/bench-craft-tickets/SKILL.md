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
2. **Migrate:** move callers in green batches, each batch sized by exactly one
   ownership fence.
3. **Contract:** remove the old form after every migrate batch lands. The
   contract ticket's `Blocked by:` names every migration ticket basename, so no
   contract runs while a migration is still open.

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

1. Classify the work first against `Classify before slicing`. A wide refactor
   takes the expand–migrate–contract sequence instead of ordinary grouping;
   otherwise group the smallest set of related stories that produces end-to-end
   behavior at a declared seam.
2. Confirm the group is independently green. Split any group that needs an
   unrelated later ticket before the project gate can pass. Concurrent
   eligibility is fence disjointness: two tickets run at once only when their
   ownership fences share no path. A one-line change pays at most one shared
   test-harness line: below that ceiling it takes no fresh worktree, no fresh
   delegate, and no full gate by default.
3. Name every real blocker by sibling ticket file basename. A ticket with all
   blockers done is on the **frontier**.
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
`/bench-final-check` remains the final full gate over the composed feature. A
ticket that changes gate cadence names which command authors gate evidence —
`bench gate`, the canonical producing entry — and which phase consumes it: a
bare `gate-run --fresh` prints a valid phase result without publishing the
project-green evidence promotion consumes.

## Discover the contracts before writing files

With the breakdown drafted and no ticket file written yet, name what crosses
each ownership fence. `craft-spec` places the fences; this step names their
traffic, so a cross-fence mismatch is a sentence at slicing time instead of a
composed red six tickets later.

Every value crossing an ownership fence names four facts: its type, its
membership or domain rule, its ordering, and its absence semantics. A fence
crossing with an unnamed fact is an undeclared contract.

Land each discovered invariant as an acceptance row on the *consumer* ticket,
asserted against the real producer and the whole enumerated family — never
against a fixture standing in for either, which is how both halves of a
mismatch look green alone.

When neither side can assert an invariant alone, add a junction ticket that
can. A junction row discovered more than one ticket downstream moves a narrower
copy of the row to the junction where it belongs, so the red surfaces at the
mismatch rather than six tickets past it.

Re-derive each contract, and every claim a ticket makes about it, from the tree
after earlier tickets land — never from the spec's account of the base.

This step's output is each ticket's `Contracts:` line (next section). A
discovery that ran leaves either named crossings or the literal claim
`none crosses`; a ticket file with no `Contracts:` line is the visible
signature of a skipped discovery, and a multi-fence breakdown whose every
ticket claims `none crosses` is a claim the review grades, not a default.

## Write one file per ticket

Write each ticket under `specs/<slug>/tickets/` with a verb-first title and this
shape:

```markdown
# <Verb-first title>

Blocked by: <sibling ticket file basenames, or none>
Ownership fence: `<path prefix>`, `<path prefix>`
Contracts: <value> crossing <fence>→<fence>, asserted by <row ID> against the real producer; or none crosses
Assumptions: <clause>; <clause>

## What to build

<The end-to-end behavior this ticket makes work.>

## Acceptance

- [ ] [AB1] <observable behavioral criterion>
- [ ] [AB2] <observable behavioral criterion>

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| <ID> | <the concrete mutation> | <the independent owner> | <the public operation sequence that proves the red> |
```

Every field is one line. The parser reads the prefixed line alone, so a
continuation wrapped onto the next line is dropped without a word.

- **Acceptance rows** are single-line `- [ ] [ID] <behavior>`. The ID is
  ticket-local: a short uppercase tag plus a number, unique within the ticket.
  Give every row its own explicit ID — only an `R`-prefixed ID range-expands
  (`[R1-R3]`), so any other tag written as a range stays one literal row.
- **`Ownership fence:`** enumerates every path the ticket writes,
  comma-separated, each entry backticked. An entry is a path prefix: a package
  directory, or an exact file. Checkpoint enforcement is a whitelist, so a path
  left off the line is a path the ticket may not touch, and scoping prose ("the
  shellcheck entries of…") names nothing the fence can hold.
- **`Assumptions:`** separates its clauses with semicolons, because the parser
  splits the line on commas and a comma-joined sentence shatters into fragments.
- **`Blocked by:`** is `none`, or the file basenames of the sibling tickets that
  must land first. A basename survives a retitle and is what `--ticket` already
  names; a title does not.
- **`Contracts:`** is the contract-discovery step's landing site: each value
  crossing this ticket's fence, with the acceptance row that asserts it
  against the real producer or the junction ticket that will — never a
  fixture standing in for either. `none crosses` is a falsifiable claim, not
  a default; writing it on every ticket of a multi-fence build asserts the
  fences exchange nothing, which the review checks.
- **`## Red mutations`** binds one row per acceptance ID: the concrete mutation
  that breaks that criterion, an owner independent of the code under test, and
  the public operation sequence that shows the red. Re-derive every claim in the
  ticket from the tree after earlier tickets land, never from the spec's account
  of the base.

Good:

<!-- ticket-example:begin -->
```markdown
# Render cancelled jobs in status

Blocked by: parse-cancelled-job-records.md
Ownership fence: `internal/status`, `internal/render/rows.go`
Contracts: the cancelled record with its reason crosses `internal/parse`→`internal/render`, asserted by RC1 against the real parser's output
Assumptions: the parser already emits a cancelled record carrying its reason; the recovery action is derived at render time and stored nowhere; claims re-derived from the tree at pickup

## What to build

Users see a cancelled row, its reason, and the next recovery action.

## Acceptance

- [ ] [RC1] status renders the cancelled row with its reason.
- [ ] [RC2] status renders the recovery action beside a cancelled row.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RC1 | render the cancelled row with an empty reason | the cancelled-row render test | blank the field, run `go test ./internal/render`, expect the missing-reason failure |
| RC2 | return no recovery action for a cancelled record | the recovery-action render test | return the empty action, run `go test ./internal/render`, expect the missing-action failure |
```
<!-- ticket-example:end -->

This is a verb-first, end-to-end outcome, written in the shape the parser
accepts: distinct row IDs, every written path fenced, semicolon-separated
assumptions, and one mutation per ID.

Bad:

```markdown
# Status changes

Blocked by: none
Ownership fence: `internal/status`

## What to build

Status has three consumer classes — operator, CI, and dashboard — each with its
own latency budget and failure tolerance, so the render path has to …

The renderer grew out of the log formatter, which is where the column widths
come from. In last week's review we went back and forth on whether …

## Acceptance

- [ ] [S1] the cancelled row renders.
- [ ] [S2] the reason renders.
- [ ] [S3] the recovery action renders.
…fourteen more rows, one per consumer class and render mode…
```

This is credible prose in the wrong artifact: a spec fragment, a taxonomy, and
a review thread pasted into a unit that no longer fits one fresh context, with
a row count no single green landing can carry. A ticket holds only what a fresh
context needs to land this behavior green.
