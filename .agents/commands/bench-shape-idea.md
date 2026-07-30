---
description: Turn a multi-session unresolved decision tree into a reviewed decision map. Use only when the idea needs more than one session of decisions before it can be specified.
---

# /bench-shape-idea — push back the fog

## Entry orientation

This is the decision-mapping phase. Decision maps are situational, not
mandatory: use one when unresolved reviewer choices form a dependency tree or
must survive across sessions. A clear idea with no such fog proceeds directly
to `/bench-write-spec`.

The phase produces or resumes a compact top-level `decisions/<topic>.md` map.
Run `bench maps --template` for the canonical paste-ready schema; the CLI
template and parser share one schema owner, so this command does not restate
the exact Markdown grammar.

## Exit handoff

Close by reporting which decision map changed, which decision tickets are
resolved, blocked, deferred, or on the frontier, and whether the map is ready
for a spec. Recommend `/bench-write-spec` when it is ready; otherwise recommend
one focused `/bench-shape-idea` resume on the next frontier decision ticket.

## Ownership

Shaping owns reviewer decisions, constraints, exclusions, research objects,
rejected alternatives, and bounded discretion. Bounded discretion is limited
to reversible technical choices that do not change observable behavior, scope,
an architectural seam, compatibility, or what the gate proves.

Spec authoring owns engineering seams, tests, acceptance coverage,
hostile-input attachment, and gate attachment. A reviewer may explicitly
choose an engineering seam while shaping; record that choice in its decision
ticket answer, but do not manufacture seam decisions simply to complete a map.

The map stays top-level while shaping is open and is loaded whole into each
planning session, so keep it tight. Link to research assets instead of inlining
them. `/bench-write-spec` moves a ready map and its owned assets into
`specs/<slug>/decisions/`; compiled maps there are settled provenance, not the
active shaping frontier.

## Decision tickets

A decision ticket is one unresolved reviewer choice or evidence-producing
question, sized for a focused session and connected through its `Blocked by`
IDs. Use the four schema-owned types:

- **Research** — read primary docs, APIs, or local code and produce a short
  cited summary asset. Include a runnable compatibility probe when the answer
  claims byte or wire compatibility.
- **Prototype** — write throwaway code to make a reviewer choice concrete.
- **Grill** — use `craft-grill` to surface and record the reviewer decision.
- **Task** — complete manual work needed before a decision can be made, naming
  who owns the work.

Grill and Prototype decision tickets resolve only through live exchange with
me. Research runs agent-alone. When the harness can delegate and another
frontier decision ticket can proceed concurrently, route Research through
`craft-delegate` as a read-only delegation; otherwise resolve it inline.
Before asking me about a fact, look it up in the tree — reviewer attention is
for decisions.

## Map content

Use `bench maps --template`, then fill the map as current state rather than a
deliberation log:

- **Destination** fixes the outcome and scope.
- **Decision tickets** record each question, dependency, type, and answer.
- **Not yet specified** holds honest in-scope fog that is not sharp enough to
  become a decision ticket.
- **Spec-writer discretion** lists only the bounded discretion shaping grants.
- **Out of scope** records reviewed exclusions so later sessions do not reopen
  them.
- **Sources** records structured research objects and the decision each
  supports. An all-grilled map may leave this section empty.

The map is ready only when its status is `ready`, every decision ticket is
resolved, fog is empty, and every research object remains valid. Run
`bench maps` before declaring readiness; a diagnostic or any row for this map
means it is not ready.

## Starting from the roadmap

`ROADMAP.md` is the working prioritization document `/bench-what-next`
maintains. Raw capture lives in `IDEAS.md` and reaches the roadmap only through
a reviewed `/bench-what-next` drain.

When invoked cold, read `ROADMAP.md` and offer its top items, recommended
sequence first, asking which to pull up. When the conversation already carries
a fresh idea, proceed with it without interrupting for the roadmap. If the
roadmap is empty or absent, say so, note that `/bench-what-next` rebuilds it,
and continue.

Pulling an item leaves its roadmap row in place: row presence is the item's
status as current open work, and the row remains through shaping, spec, and
build until shipped retirement removes it. This command never edits
`ROADMAP.md`.

## Two modes

**Bootstrap** (loose idea in): use `craft-grill` to discover whether the idea
really contains a multi-session dependency tree. If it does, create the map
from `bench maps --template`, record the frontier, and stop after that first
focused shaping pass. If it does not, create no map and recommend
`/bench-write-spec` from the reviewer-confirmed current conversation.

**Resume** (map + decision-ticket number in): load the whole map, resolve that
decision ticket, record the current answer, and add any newly discovered
decision tickets with correct `Blocked by` edges. If an answer invalidates
other tickets, update or delete them. While the reviewer is present, carry a
Grill straight into newly unblocked decision tickets; never pause for
permission or a re-prompt between questions. Stop when no reviewer-answerable
frontier remains or the reviewer stops answering. Leave Research, Prototype,
or Task evidence work explicit instead of force-grilling it.

## The exit

Run `bench maps` and re-read every decision ticket, constraint, exclusion,
research object, and discretion item against the current conversation and
tree. A conversation answer not written into the map is not recorded.

A possible scope split is a reviewer decision, not an excuse to close shaping
early. Surface it as a decision ticket and continue until the reviewer either
chooses the split or rejects it. When the map reaches `ready`, recommend a
fresh mid-tier `/bench-write-spec` session with one clause explaining why the
decision source is ready.
