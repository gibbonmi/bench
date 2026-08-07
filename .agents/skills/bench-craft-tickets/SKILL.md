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

The green landing commit is the grading rule.

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

Any temporary guard or oracle introduced by an expand ticket names in
`Integration surfaces:` the dependent contract ticket that retires or migrates
it and the single final owner that survives contraction. That map is the
sequence's inventory source; dependent tickets point to it instead of carrying
independent current and final inventories.

Every expand, migrate, and contract ticket must land green independently. If a
migrate batch cannot, the expansion or prefactor is incomplete, or the batch
is too wide; repair the preparation or split the batch before proceeding.

`craft-spec` owns the spec-time **who-writes-where** fence. This skill owns the
build-time **what-lands-green-next** unit. A repair ticket derived from a
debug receipt takes its ownership fence from the receipt's required paths —
never from the blocked ticket's fence, which the repair must stay out of. Apply the fence by name; do not
restate or redraw it here.

## Draft the breakdown

A seam is a reason to inspect a ticket boundary, not an automatic ticket
boundary. Split until splitting further would leave no independently-green
landing. Keeping a group whole requires naming the specific red the thinner cut
would strand; review re-derives that falsifiable claim by attempting the split.
Feature wholeness does not justify the grouping.

| Signal | Response |
|---|---|
| A proposed split strands a specific project-gate red | Keep it together; name that red for review to reproduce by attempting the split. |
| Two behaviors are independently useful and independently green | Split them. |
| Two or more tickets need the same primitive | Prefactor one owning ticket and block the consumers on it. |
| Neither side of a contract can prove the invariant alone | Create the junction ticket described in `Discover the contracts before writing files`. |
| Stories partition into disjoint package groups | Return to `craft-spec`'s `Check the story partition before locking scope`; that spec-time decision owns the split. |
| A mechanical refactor cannot keep an ordinary tracer green | Use `Classify before slicing`; its expand–migrate–contract sequence owns the cut. |

For each candidate ticket:

1. Classify the work first against `Classify before slicing`. A wide refactor
   takes the expand–migrate–contract sequence instead of ordinary grouping;
   otherwise take the ordinary-build branch and apply the sizing rule above at
   a declared seam.
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

The first response to a discovered crossing is to remove it. Prefer regrouping
so one vertical ticket owns both sides of the value; a mechanical change whose
crossing comes from a form in flux takes the expand–migrate–contract sequence
instead. A producer/consumer split that survives both preferences is the
exception, and it carries its traffic explicitly: the producer's `Integration
surfaces:` name the dependent, the consumer's `Blocked by:` names the producer,
and assignment refuses the consumer until the edge is reciprocal. A crossing
kept without naming why neither regrouping applies is a slicing defect, not a
contract to document — a split ticket pair maintains its shared value in prose
across every mid-run repair, which is exactly the surface that goes stale.

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

Before discovery is complete, enumerate every required integration surface for
each file, registry member, command, or fixture the ticket adds, changes, or
removes. Start from the nearest existing member of the same family and find every
path that registers or classifies it, records it in a manifest or inventory,
routes it, derives a count from it, or consumes it. A new fixture in a classified
family therefore carries both its classifier and its real consumer as surfaces.

An extraction or single-source ticket also enumerates consumers by the moved
fact's exact literal values and exported symbols, not only by the abstraction
category being extracted. A residue sweep for tables or matchers is incomplete
until it accounts for bespoke consumers of those facts. An exported symbol left
with zero current consumers is residue unless an acceptance row explicitly owns
it as public API or a named dependent ticket consumes it next.

This step's outputs are each ticket's `Ownership fence:`, `Blocked by:`,
`Integration surfaces:`, and `Contracts:` lines (next section). Discovery is
complete only when `Integration surfaces:` resolves every required surface to
exactly one owner: a fence entry; an existing unchanged path exercised by a named
acceptance row; a blocker basename; or a dependent basename whose `Blocked by:`
names this ticket. Every cross-fence fact also resolves to a named contract. The
literal `none` is a falsifiable claim that review re-derives from the tree; a
missing line is the visible signature of skipped integration discovery.
`Contracts:` likewise carries either named crossings or `none crosses`; a missing
line skipped value-contract discovery, and a multi-fence breakdown whose every
ticket claims `none crosses` is a claim the review grades, not a default.

Before assignment, close the discovery audit: every mechanical fact named in
`Contracts:`, and every applicable mechanical promise from the spec's coverage
map or edge inventory, appears explicitly in `## Acceptance` and carries a
machine-checkable identity through its red-mutation row or rows. A fact found
only in `Contracts:`, `What to build`, or other surrounding prose leaves the
audit open. If the ticket or schema cannot represent that traceability, stop
ticketing and report the missing seam; prose or review cannot substitute for
it. A genuinely semantic reviewer-only claim remains an explicit exception in
the acceptance row and coverage map.

## Write one file per ticket

Write each ticket under `specs/<slug>/tickets/` with a verb-first title and this
shape:

```markdown
# <Verb-first title>

Blocked by: <sibling ticket file basenames, or none>
Ownership fence: `<path prefix>`, `<path prefix>`
Integration surfaces: <surface>→<fence path | existing path + row ID | blocker basename | dependent basename>; or none
Contracts: <value> crossing <fence>→<fence>, asserted by <row ID> against the real producer; or none crosses
Closure: <acceptance ID>/<atomic fact>, <acceptance ID>/<atomic fact>

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

Every field is one line. Lifecycle-parsed fields are read from their prefixed
line alone, so a continuation wrapped onto the next line is dropped without a
word.

- **Acceptance rows** are single-line `- [ ] [ID] <behavior>`. The ID is
  ticket-local: a short uppercase tag plus a number, unique within the ticket.
  Give every row its own explicit ID — only an `R`-prefixed ID range-expands
  (`[R1-R3]`), so any other tag written as a range stays one literal row.
  Before keeping independently failing classes or members under one ID, apply
  the compound-claim rule in `## Red mutations` below.
  Under a spec whose coverage map opts into row IDs (`craft-spec`'s leading
  `row` column), every row also carries `(covers <ID>)` or `(covers local)`
  after the row-ID bracket, naming exactly one map row — the mapping is 1:1.
  `local` marks a ticket-time discovery or repair row: assign accepts it,
  promote's totality ignores it, and review grades whether the marker is
  honest. The annotation attaches to single-ID rows only, so a range line's
  expanded rows are unannotated — and under an opted-in spec a row missing its
  annotation refuses at assign, which is what keeps omission from dodging the
  mapping.
- **`Ownership fence:`** enumerates every path the ticket writes,
  comma-separated, each entry backticked. An entry is a path prefix: a package
  directory, or an exact file. Checkpoint enforcement is a whitelist, so a path
  left off the line is a path the ticket may not touch, and scoping prose ("the
  shellcheck entries of…") names nothing the fence can hold. Derive the fence
  from the `Contracts:` line below and the `Integration surfaces:` line below,
  rather than from the lines that prompted the ticket: a value this ticket
  changes is advertised wherever it is restated — a spec's count, a profile row,
  a sibling ticket — and each of those
  advertisements is a path this ticket writes. A repair fenced to a finding's
  cited lines cannot maintain an advertisement it may not touch, so the
  contradiction surfaces a review round later instead of at slicing time.
- **`Integration surfaces:`** is the review-owned integration-discovery landing
  site; lifecycle does not infer fence completeness from it, and reads it for
  exactly one refusal: naming a sibling basename as a dependent obliges that
  sibling's `Blocked by:` to name this ticket back, and assign and refresh
  refuse the sibling until it does. Whether the named value is really what this
  ticket exports stays review's question. Resolve
  every required surface through the completion criterion above. `none` asserts
  that the sibling-family and literal/symbol searches found no integration
  surface; review repeats those searches rather than accepting the word as proof.
- A genuinely unverifiable-at-authoring-time claim belongs in the ticket's
  What-to-build prose; a checkable precondition belongs on `Blocked by:`.
- **`Blocked by:`** is `none`, or the file basenames of the sibling tickets that
  must land first. A basename survives a retitle and is what `--ticket` already
  names; a title does not.
- **`Contracts:`** is the contract-discovery step's landing site: each value
  crossing this ticket's fence, with the acceptance row that asserts it
  against the real producer or the junction ticket that will — never a
  fixture standing in for either. `none crosses` is a falsifiable claim, not
  a default; writing it on every ticket of a multi-fence build asserts the
  fences exchange nothing, which the review checks. Each crossing anchors at
  least one backticked path inside this ticket's own fence — the side this
  ticket writes; the other side may name a surface no path holds (`bash`, every
  audited package). A crossing written entirely in concepts
  (`registry`→`derived inventory`) anchors nothing and `assign` refuses it.
- **`Closure:`** is the machine-checkable inventory for the mechanical facts
  discovered above. Give every independently failing class or member a unique
  `<acceptance ID>/<lowercase-kebab-name>` token. Every acceptance ID owns at
  least one token, and every token appears as the criterion of a red-mutation
  row. `bench spec build assign` refuses a modern ticket whose inventory is
  absent or open. The checker proves the declared graph is closed; review still
  compares the inventory to `Contracts:`, the coverage map, and the edge
  inventory so an undeclared fact cannot hide behind a mechanically valid graph.
- **`## Red mutations`** binds every acceptance ID to one or more rows: each
  names a concrete mutation that breaks the criterion, an owner independent of
  the code under test, and the public operation sequence that shows the red.
  Split independently failing classes or members into separate acceptance IDs,
  or repeat the same ID in the mutation table until every class or member is
  exercised. One representative mutation cannot prove a compound or quantified
  claim. An outcome label ("stale or absent", "invalid") is a compound claim in
  disguise: name the exact predicate the code must preserve — the specific error
  class, the specific comparison — so each closure token owns one singular red.
  Mutate the **subject**, never
  the assertion: weakening a shared check to always-pass is invisible to a suite
  whose subjects already satisfy it, so the row reads green and proves nothing.
  A ticket touching a value that **authorizes** an action — a fingerprint, a
  digest, a token — carries at least one row that mutates the value's *inputs*:
  revert the authorizing constant, or drop a field from the hash preimage.
  Control-flow mutations grade only the plumbing that carries the value; a
  reverted constant or a dropped hashed field passes every one of them, and
  only an input mutation grades what the value commits to.
  The red must be a **bounded failure**, because a hung run and a broken harness
  are the same observation at the gate — bound whatever the mutation can stall
  (the wait, the poll, the child that outlives its test) before claiming the row.
  Re-derive every claim in the ticket from the tree after earlier tickets land,
  never from the spec's account of the base.

Good:

<!-- ticket-example:begin -->
```markdown
# Render cancelled jobs in status

Blocked by: parse-cancelled-job-records.md
Ownership fence: `internal/status`, `internal/render/rows.go`
Integration surfaces: cancelled-record producer→parse-cancelled-job-records.md; status row→internal/status; render schema→internal/render/rows.go
Contracts: the cancelled record with its reason crosses `internal/parse`→`internal/render`, asserted by RC1 against the real parser's output
Closure: RC1/reason, RC2/recovery-action

## What to build

Users see a cancelled row, its reason, and the next recovery action. The
existing cancelled-status contract rejects a cancelled row missing either its
reason or its recovery action, so a thinner field-by-field cut strands that
contract red.

## Acceptance

- [ ] [RC1] (covers CJ1) status renders the cancelled row with its reason.
- [ ] [RC2] (covers CJ2) status renders the recovery action beside a cancelled row.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| RC1/reason | render the cancelled row with an empty reason | the cancelled-row render test | blank the field, run `go test ./internal/render`, expect the missing-reason failure |
| RC2/recovery-action | return no recovery action for a cancelled record | the recovery-action render test | return the empty action, run `go test ./internal/render`, expect the missing-action failure |
```
<!-- ticket-example:end -->

This is a verb-first, end-to-end outcome, written in the shape the parser
accepts: distinct row IDs, a covers annotation per row, every written path
fenced, an atomic closure inventory, and mutation coverage for every fact.

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
no visible contract or integration-surface discovery and scope no single green
landing can carry. A ticket holds only what a fresh context needs to land
this behavior green.
