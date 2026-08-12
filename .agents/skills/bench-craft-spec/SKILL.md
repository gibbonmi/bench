---
name: craft-spec
description: The spec-authoring discipline — the acceptance-coverage-map row schema, what counts as a real red signal, the canonical edge-inventory classes, how stories and scope cuts are sized, and how a build is sliced for delegates. Use when authoring or auditing a spec, writing coverage rows, asking "what's the red signal", walking an edge inventory, judging whether a cut is genuinely out of scope, or slicing a build for delegates.
index: coverage-map rows, edge inventories, story sizing, and delegate slicing for a spec
---

# The spec discipline: rows, edges, and scope

A spec fixes the build's target: stories set the breadth, seams set where tests
live, and the coverage map makes "done" checkable by the gate instead of by
belief. This skill is the one source for the map's row schema, the edge-class
list, the sizing rules, and the delegate-slicing rule.

## The acceptance coverage map

Each row ties a user story to one observable behavior at a chosen seam, with
these fields: `story`, `behavior`, `seam`, `red signal`, and
`why it catches the failure`. An optional leading `row` column opts the spec
into ticket covers traceability: a 6-cell map whose header leads with `row`
gives each row a spec-local ID that ticket acceptance rows name in their
`covers` annotations; a 5-cell map grades exactly as today. The red signal is
the command or test already run and failed because the mapped behavior is
absent or wrong — if it's already covered or cannot start red, say so in the
row instead of pretending it is TDD coverage. Newly authored specs default
to the 6-cell shape.

Three rules keep a map honest:

- **The degenerate-implementation check.** Name the cheapest wrong
  implementation of each story and confirm a row goes red on it. Across
  fences, also name the **composition degenerate** — per-fence rows stay
  green while the end-to-end command fails — confirm a row through the real
  producer catches it; per-fence fixture rows cannot see it.
- **Enumerate every quantifier.** When a promise quantifies over a set ("each
  check", "every parser"), enumerate it or state the granularity; an
  unenumerated "each" lets the build pick the cheapest reading.
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
defects appearing only once state is serialized and a fresh process reloads
it. Every edge lands in a coverage row (the story column may read "edge of
N") or a one-line **Won't handle** entry directly under the map; both are
veto surface, and a silently untested edge is the failure the walk prevents.

A control resolving a class must name how it exercises the **new** surface;
a control predating the change covers the old code by default.

## Bootstrap authority before execution

When a story claims trusted or authenticated execution, or refusal before
execution, trace the real process from the raw OS entrypoint to the first
candidate-controlled instruction and through every executable hop, naming at
each hop the already-trusted validator and how it authenticates the next
executable before launching the next executable. A path, record, digest, or
executable cannot authenticate itself. Without an independent trust root, the
design is incomplete unless a reviewer-visible trust assumption says otherwise.
Coverage places markers in candidate-controlled executables, corrupts or
replaces the next authority, and asserts no marker runs before refusal;
slicing names who publishes, locates, validates, and invokes the first
trusted executable, and no complete owner is a pre-build slicing defect.

Two guards on the exclusions: no **amputated callers** (a **Won't handle**
line needs one surviving in-scope caller) and **compatibility proven, not
promised** (divergence from a named external format is a reviewer decision,
never a silent promise).

## Story sizing and scope cuts

User stories are the breadth **floor** — the guard against a loop that does
the minimum and stops. Price every proposed cut in *agent* time, derived
rather than guessed: state it as `<n> edits, <n> gate runs`. Size each story
as an independently deliverable and demonstrable tracer outcome — one
complete behavior travels through the system and can be shown on its own —
and reject a horizontal engineering layer wearing a story name.

- **No deferral under the threshold.** Anything under ~30 minutes of agent work
  with no new architectural decision is part of this build; do it and state
  the estimate. Binds the author, never the reviewer.
- **A cut must be a separate capability, not the rest of this one.** Only
  legitimately out of scope with its own future spec — "the rest of *this*
  feature" is an acceptance criterion, so move it into the stories and the map.

**Check the story partition before locking scope.** Disjoint package sets
with no shared seam are a split signal for the reviewer; a bundle is
legitimate only when chosen, never defaulted. `craft-tickets` owns build-time
wide-refactor classification; point to that rule rather than restating it.

## Slicing a build for delegates

At spec time, record explicit **who-writes-where** ownership fences. Each
fence names every path one writer may edit and is checkable at charge time. A
fence entry is an exact repo-relative file or path prefix, never a glob or an
implementation ticket, and an empty or invalid fence section is incomplete.

`craft-tickets` owns the build-time **what-lands-green-next** unit. This
section owns only the spec-time **who-writes-where** fence; point to the
ticket rule by name rather than restating it here. `craft-tickets` derives
complete tracer tickets from the locked stories, and each resulting ticket
receives the spec's applicable ownership fence.

Each fence carries value contracts across it, and `craft-tickets` owns naming
them in `Discover the contracts before writing files`; this section points at
that step by name rather than restating what it requires.

## Reread the whole artifact after wide edits

A pass touching many sections of the same spec, skill, or doctrine file can
leave one section's claim contradicting another's even though each edit
looked correct alone. After such a pass, reread the complete artifact end to
end and reconcile any contradiction before handing off: reviewing just the
changed hunks misses a stale sentence the edit silently falsified.
