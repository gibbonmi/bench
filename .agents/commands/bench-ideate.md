---
description: Turn a loose idea into a sequenced decision map. Use only when the idea needs more than one session of decisions before it can be built.
---

# /bench-ideate — push back the fog

Use this when an idea is too unresolved to spec yet — when there are open
questions whose answers change what gets built. If the idea is already clear,
skip straight to `/bench-spec`.

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
- **Grill** — converse to surface the decision. Use `/bench-craft-grill`. The default.

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

**Bootstrap** (loose idea in): run a `/bench-craft-grill` pass to surface the open decisions,
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
declaring the map closed, scan it for unwritten answers — any ticket still holding a
`— (open` / `— (deferred` placeholder or a `GRILL DEFERRED` banner — and refuse to
close while any remain; a decision made in conversation but not written into the map
is not recorded. Then close the map and lead with the recommended next action
(usually `/bench-spec`) and a one-clause why.

A natural seam in the work is recorded as a decision in the map for *me* to make —
it is never a reason to close the map early or to spin off a separate map or PRD on
your own. Surfacing "this could be two slices" is useful; deciding to slice and
deferring the rest is mine.

If the first grill surfaces no real fog — no multi-session decisions — say so and
recommend skipping straight to `/bench-spec`. Don't manufacture tickets.