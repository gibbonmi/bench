---
description: Turn a loose idea into a sequenced decision map. Use only when the idea needs more than one session of decisions before it can be built.
---

# /start-ideation — push back the fog

Use this when an idea is too unresolved to spec yet — when there are open
questions whose answers change what gets built. If the idea is already clear,
skip straight to `/spec`.

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
- **Grill** — converse to surface the decision. Use `/grill`. The default.

## Two modes

**Bootstrap** (loose idea in): run a `/grill` pass to surface the open decisions,
write the map with the frontier identified and the trivially-decidable tickets
resolved inline, then stop. Building the map is one session's work — do not also
resolve tickets.

**Resume** (map + ticket number in): load the whole map, resolve that one ticket,
record the answer in its body (current state only), add any newly-discovered
tickets with correct `Blocked by` edges, then stop. If a resolution invalidates
other tickets, update or delete them.

## The exit

The map is deliberately incomplete beyond the frontier. You are done when the
path to the finish line is clear — no unresolved tickets blocking the build. At
that point the map is closed; offer to move to `/spec`.

A natural seam in the work is recorded as a decision in the map for *me* to make —
it is never a reason to close the map early or to spin off a separate map or PRD on
your own. Surfacing "this could be two slices" is useful; deciding to slice and
deferring the rest is mine.

If the first grill surfaces no real fog — no multi-session decisions — say so and
recommend skipping straight to `/spec`. Don't manufacture tickets.