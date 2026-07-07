---
name: craft-spec
description: The spec-authoring discipline — the acceptance-coverage-map row schema, what counts as a real red signal, the canonical edge-inventory classes, and how stories and scope cuts are sized. Use when authoring or auditing a spec, writing coverage rows, asking "what's the red signal", walking an edge inventory, or judging whether a cut is genuinely out of scope.
index: coverage-map rows, edge inventories, and story sizing for a spec
---

# The spec discipline: rows, edges, and scope

A spec fixes the build's target: stories set the breadth, seams set where tests
live, and the coverage map makes "done" checkable by the gate instead of by
belief. This skill is the one source for the map's row schema, the edge-class
list, and the sizing rules — `/bench-write-spec` composes it, and `craft-tdd`
and `craft-review` point here when they need the schema off the phase path.

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
  implementation passes has not pinned the behavior.
- **Enumerate every quantifier.** When a behavior or red-signal promise
  quantifies over a set ("each check", "every parser"), enumerate the set or
  state the granularity explicitly — per item or per class. An unenumerated
  "each" lets the build pick the cheapest reading, and review is the wrong
  place to catch that.
- **Assertables become rows.** When a decision map's Handoff names black-box
  assertables, each one lands as a coverage row or a stated exception; an
  assertable with no row is a missing behavior, not an editorial choice.

## The edge inventory

Stories are happy-path shaped; the edge walk generates the cases nobody
declared. For each mapped behavior, walk the canonical edge classes — error
path, empty/absent input, boundary values, malformed input, interrupted or
partial state, re-run idempotency, hostile environment — plus the project
profile's hostile-input checklist when one exists. Every edge lands in exactly
one of two places: a coverage row (the story column may read "edge of N"), or
a one-line **Won't handle** entry directly under the map. Both are veto
surface; a silently untested edge is the failure the walk exists to prevent.

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
