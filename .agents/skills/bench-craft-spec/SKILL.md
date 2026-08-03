---
name: craft-spec
description: The spec-authoring discipline — the acceptance-coverage-map row schema, what counts as a real red signal, the canonical edge-inventory classes, how stories and scope cuts are sized, and how a build is sliced for delegates. Use when authoring or auditing a spec, writing coverage rows, asking "what's the red signal", walking an edge inventory, judging whether a cut is genuinely out of scope, or slicing a build for delegates.
index: coverage-map rows, edge inventories, story sizing, and delegate slicing for a spec
---

# The spec discipline: rows, edges, and scope

A spec fixes the build's target: stories set the breadth, seams set where tests
live, and the coverage map makes "done" checkable by the gate instead of by
belief. This skill is the one source for the map's row schema, the edge-class
list, the sizing rules, and the delegate-slicing rule — `/bench-write-spec`
composes it, and `craft-tdd`, `craft-review`, and `craft-delegate` point here
when they need the schema off the phase path.

## The acceptance coverage map

Each row ties a user story to one observable behavior at a chosen seam, with
these fields: `story`, `behavior`, `seam`, `red signal`, and
`why it catches the failure`. The red signal is the command or test that has
already been run and failed because the mapped behavior is absent or wrong. If
the behavior is already covered or cannot start red, say so in the row instead
of pretending it is TDD coverage.

Three rules keep a map honest:

- **The degenerate-implementation check.** Before locking the map, name the
  cheapest wrong implementation of each story — the sequential port, the
  always-green stub — and confirm a row goes red on it; a map the degenerate
  implementation passes has not pinned the behavior. When the map's seams span
  more than one package or ownership fence, also name the **composition
  degenerate**: the cheapest implementation that keeps every per-fence row
  green while the end-to-end command still fails — one side softened, the
  other still refusing — and confirm a row driven through the real producer
  goes red on it. Per-fence rows asserted against fixtures cannot see this
  degenerate; that is how both halves of a mismatch look green alone.
- **Enumerate every quantifier.** When a behavior or red-signal promise
  quantifies over a set ("each check", "every parser"), enumerate the set or
  state the granularity explicitly — per item or per class. An unenumerated
  "each" lets the build pick the cheapest reading, and review is the wrong
  place to catch that.
- **Source behaviors become rows.** Every observable behavior promised by the
  authorized decision source lands as a coverage row or a stated exception; a
  promised behavior with no row is missing coverage, not an editorial choice.

## The edge inventory

Stories are happy-path shaped; the edge walk generates the cases nobody
declared. For each mapped behavior, walk the canonical edge classes — error
path, empty/absent input, boundary values, malformed input, interrupted or
partial state, re-run idempotency, process-boundary lifecycle, hostile
environment — plus the project profile's hostile-input checklist when one
exists. Process-boundary lifecycle is the class unit-level success hides:
defects that appear only once state is serialized and a fresh process reloads
it, and recomposition suites that stop at the first success instead of walking
every recomposition path. Every edge lands in exactly one of two places: a
coverage row (the story column may read "edge of N"), or a one-line
**Won't handle** entry directly under the map. Both are veto surface; a
silently untested edge is the failure the walk exists to prevent.

A class resolved by pointing at an existing control must name how that
control exercises the **new** surface. A control that predates the change
covers the old code by default — "resolved by the existing re-entry test"
over a newly added action is a claim about a test that has never run the
action, and it rots into an untested wedge. If the control cannot be shown
to reach the new code, the class becomes a row.

Two guards on the exclusions:

- **No amputated callers.** Before writing a **Won't handle** line about an
  interface, verify at least one in-scope caller can still exercise the
  feature under that exclusion — a cut that amputates the surface's primary
  calling convention is a spec defect, not a scope cut.
- **Compatibility is proven, not promised.** If the spec names an external
  format or protocol, check whether an official implementation exists and
  whether the current output conforms; divergence is a reviewer decision, not
  a silent compatibility promise. If a map or research asset claims byte or
  wire compatibility, require a runnable probe against the caller's own edge
  outputs as the evidence.

## Story sizing and scope cuts

User stories are the breadth **floor** — the guard against a loop that does
the minimum and stops. Price every proposed cut in *agent* time, derived
rather than guessed: agent time is dominated by verification, so state it as
`<n> edits, <n> gate runs` — a vibes number can't pass as a price. Two rules:

- **No deferral under the threshold.** Anything under ~30 minutes of agent
  work that introduces no new architectural decision is not a candidate for
  deferral — it is part of this build; do it and state the estimate. The rule
  binds the author, never the reviewer: scope stays the reviewer's, and a cut
  the reviewer chooses is recorded with its estimate and no argument.
- **A cut must be a separate capability, not the rest of this one.** Something
  is only legitimately out of scope if it has its own future spec — a distinct
  feature. "The rest of *this* feature" (the error cases, the other half of
  the CRUD, the edge handling) is an acceptance criterion: move it into the
  stories and the map so the gate enforces it. You can't ticket your way out
  of the spec's own breadth.

`craft-tickets` owns build-time wide-refactor classification and its
expand–migrate–contract sequence. A spec names that the build is wide and
keeps the ownership fences below explicit; it points to that rule rather than
restating the sequence.

## Slicing a build for delegates

When a spec splits a build into delegate slices, the slice boundary and the
ownership fence are the same line: each slice owns every file it must edit and
edits every file it owns. Cut slices by ownership, not by behavior theme — a
theme names what the work is about; a fence names who writes where, and only
the fence is checkable at charge time. A slice that reaches outside its fence
collides with another writer, and one that stops short returns with its own
work unwritten; either way the coordinator pays a round trip.

`craft-tickets` also owns prefactoring shared primitives before their
consumers. This skill gives the resulting primitive ticket one ownership fence;
it does not repeat the ordering rule.

Fence alignment is tier-independent: a future cheap-tier retest changes
delegate routing, never this rule.

`craft-tickets` owns the build-time **what-lands-green-next** unit. This
section owns only the spec-time **who-writes-where** fence; point to the ticket
rule by name rather than restating it here.

Each fence carries value contracts across it, and `craft-tickets` owns naming
them in `Discover the contracts before writing files`; this section points at
that step by name rather than restating what it requires.

```
Slice per package — internal/scan, internal/report, internal/bounds — one
delegate each, owning every file in its package.
```
Good — the package boundary is the ownership fence, so no two slices write
the same file and the merges cannot conflict.

```
Slice by theme — one delegate takes "fail-closed behavior" across the tree.
```
Bad — a theme touches every package that implements it, so this slice reaches
into `internal/bounds` — another delegate's fence — and every crossing costs a
round trip.
