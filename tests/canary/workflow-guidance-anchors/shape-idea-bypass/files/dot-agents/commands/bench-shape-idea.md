---
description: Turn a loose idea into a sequenced decision map. Use only when the idea needs more than one session of decisions before it can be built.
---

# /bench-shape-idea — push back the fog

## Entry orientation

This is the decision-mapping phase. Use it when the idea has unresolved questions
whose answers change what gets built. It produces or resumes a compact
`decisions/<topic>.md` map, with the current frontier recorded and resolved
answers written into the file.

## Exit handoff

Close by reporting which decision map changed, which tickets are resolved or still
open, and whether the map is ready for a spec. The recommended next command is
`/bench-write-spec` on a fresh mid-tier session when no unresolved ticket still
blocks the build; otherwise it is another focused `/bench-shape-idea` resume on the
next open ticket.

Use this when an idea is too unresolved to spec yet — when there are open
questions whose answers change what gets built. If the idea is already clear,
skip straight to `/bench-write-spec`.

The output is a single compact markdown file, `decisions/<topic>.md`, git-tracked.
It is loaded whole into every planning session, so keep it tight. Link to assets;
don't inline them.

## Structure

```markdown
## #1: <the open question, as a question>

Blocked by: #<n>, #<n>
Type: Research | Prototype | Grill

### Question
<what we don't yet know, and why it matters for the build>

### Answer
<filled in when resolved — current state only, no deliberation log>
```

Each ticket is sized to one focused session. Three kinds:

- **Research** — read docs/APIs/local code, produce a short summary asset.
- **Prototype** — write throwaway code to answer "how should it look/behave."
- **Grill** — converse to surface the decision. Use `craft-grill`. The default.

## The Handoff — the seams a closed map hands the spec-writer

Every map closes with a `## Handoff` section: the structure `/bench-write-spec`
reads seams off instead of re-deriving them or escalating for answers the grill
already settled. It is **required on every close** — `bench maps` keeps showing
the map's row until the section is present and placeholder-free — and every item
is answered. An item that does not apply is written `n/a — <one clause>`, so the
exclusion is a decision on the page, not a silent gap.

```markdown
## Handoff

1. **Module boundaries.** Each unit and its responsibility — inside vs outside.
2. **Contracts.** Per boundary: inputs, outputs, exit codes, error posture — the
   observable interface.
3. **Deep vs thin.** Which units hide complexity so the seam attaches at the
   interface; which are pass-throughs with no seam of their own.
4. **Black-box assertables.** What a test can assert at each seam (exit code,
   stdout, file/git state) without reaching inside.
5. **Gate attachment.** Which seam the gate observes and how; flag any seam the
   gate cannot see (needs TDD or manual verify).
6. **Hostile-input owners.** Map each class from the profile's hostile-input
   checklist to the seam that owns it.
7. **Uncertainty flags.** Any seam the grill could not settle — so the spec-writer
   escalates per the `craft-line` ladder instead of guessing on the mid tier.
8. **Rejected alternatives.** So the spec-writer does not reopen a closed decision.
9. **Domain watch-outs.** Hazards stated as domain facts for any reader, never
   model-addressed coaching — operating lessons go through `.bench/learnings.md`
   and `/bench-integrate-learnings`, never per-spec notes.

Dependency order: <recommended build order when the map yields multiple slices;
`n/a — single spec` otherwise>
```

When the map yields more than one buildable slice, the Handoff's seam list is the
menu they are cut from and the **Dependency order** line records the recommended
sequence — a recommendation for me, not a decision you make. Slicing stays my call.

## Starting from the roadmap

`ROADMAP.md` (repo root) is the capture-and-forget sink: ideas parked with `bench idea`
that the user committed to nothing. This command is where a parked idea graduates into
committed work.

When invoked **cold** — no specific idea already in hand from the conversation — read
`ROADMAP.md` and offer the parked items, asking which (if any) to pull up. The chosen
entry seeds the bootstrap below. When you are already carrying a fresh idea from the
conversation, proceed with it and do **not** interrupt with the roadmap prompt; if the
roadmap is empty or absent, say so and continue.

When a pulled idea actually becomes a map — i.e. you write `decisions/<topic>.md` from
it — **remove that entry's line from `ROADMAP.md`** in the same step: promotion means it
is no longer merely parked. A pull the user abandons before any map is written leaves
the line untouched. The roadmap carries no status of its own; the line's presence *is*
its status.

## Two modes

**Bootstrap** (loose idea in): run a `craft-grill` pass to surface the open decisions,
write the map with the frontier identified and the trivially-decidable tickets
resolved inline, then stop. Building the map is one session's work — do not also
resolve tickets.

**Resume** (map + ticket number in): load the whole map, resolve that one ticket,
record the answer in its body (current state only), add any newly-discovered
tickets with correct `Blocked by` edges, then stop. If a resolution invalidates
other tickets, update or delete them.

## The exit

The map is deliberately incomplete beyond the frontier. You are done when the
path to the finish line is clear — no unresolved tickets blocking the build. Before
declaring the map closed, run `bench maps` — it lists every ticket still holding a
`— (open` / `— (deferred` placeholder or a `GRILL DEFERRED` banner, and on a map
with no open tickets a missing or still-placeholdered `## Handoff` section — and
refuse to close while this map still shows a row; a decision made in conversation
but not written into the map, or a Handoff item left unfilled, is not recorded. Then close the map and lead with the recommended next action
(usually `/bench-write-spec`) and a one-clause why.

A natural seam in the work is recorded as a decision in the map for *me* to make —
it is never a reason to close the map early or to spin off a separate map or PRD on
your own. Surfacing "this could be two slices" is useful; deciding to slice and
deferring the rest is mine.

If the first grill surfaces no real fog — no multi-session decisions — say so,
close the short map with its `## Handoff`, and recommend `/bench-write-spec` on a
fresh mid-tier session. Don't manufacture tickets.
