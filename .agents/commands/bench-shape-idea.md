---
description: Turn a multi-session unresolved decision tree into a reviewed decision map. Use only when the idea needs more than one session of decisions before it can be specified.
---

# /bench-shape-idea — push back the fog

## Entry orientation

This is the decision-mapping phase. Decision maps are situational, not
mandatory: use one when unresolved reviewer choices form a dependency tree or
must survive across sessions. A clear idea with no such fog proceeds directly
to `/bench-write-spec`.

Charge `bench-craft-domain` on entry so decision tickets use canonical terms
and concept-edge scenarios rather than overloaded vocabulary.

The phase produces or resumes a top-level `decisions/<topic>.md` map beside its
`decisions/<topic>/tickets/` folder.
The map is an index: it lists the decisions made and links the ticket that holds each one.
A decision ticket is one file under the map's tickets folder, named by its number.
The answer lives only in the ticket file.
Run `bench maps --template` for the canonical paste-ready schema, and
`bench maps --ticket-template` for one ticket file. The CLI templates and the
parser share one schema owner, so this command does not restate the exact
Markdown grammar.

## Exit handoff

Report which decision map changed. Report each decision ticket's status:
resolved, blocked, deferred, or on the frontier. Report whether the map is
ready for a spec. Recommend `/bench-write-spec` when it is ready. Otherwise
recommend one focused `/bench-shape-idea` resume on the next frontier decision
ticket.

## Ownership

Shaping owns reviewer decisions, constraints, exclusions, research objects,
rejected alternatives, and bounded discretion. Bounded discretion covers only reversible technical choices. Such a choice
does not change observable behavior, scope, an architectural seam,
compatibility, or what the gate proves.

Shaping resolves scope into deliverable outcomes. It does not inventory
engineering seams as decomposition units.

Spec authoring owns engineering seams, tests, acceptance coverage,
hostile-input attachment, and gate attachment. A reviewer may explicitly
choose an engineering seam while shaping. Record that choice in its decision
ticket answer. Do not manufacture a seam decision simply to complete a map.

The topic folder stays top-level while shaping is open. The index loads whole
into each planning session, so keep it tight. Link to a research asset instead
of an inline copy. A map-owned asset stays in the map's assets folder,
`decisions/<topic>/assets/`.
`/bench-write-spec` moves a ready topic folder into `specs/<slug>/decisions/`.
Compiled maps there are settled provenance, not the active shaping frontier.

## Decision tickets

A decision ticket is one unresolved reviewer choice or evidence-producing
question, sized for a focused session and connected through its `Blocked by`
IDs. Use the four schema-owned types:

- **Research** — read primary docs, APIs, or local code and produce a short
  cited summary asset. Include a runnable compatibility probe when the answer
  claims byte or wire compatibility.
- **Prototype** — write throwaway code to make a reviewer choice concrete;
  charge the `prototype` skill.
- **Grill** — run `craft-grill` frontier rounds to record the reviewer decision.
- **Task** — complete the manual work needed before the reviewer can decide.

Name a ticket by its title, with its number beside it.
The shaping worktree lease is the claim, and no owner field enters a ticket.

Grill and Prototype decision tickets resolve only through live exchange with
me. Research runs agent-alone. When the harness can delegate and another
frontier decision ticket can proceed concurrently, route Research through
`craft-delegate` as a read-only delegation. Otherwise resolve it inline.
Before asking me about a fact, look it up in the tree — reviewer attention is
for decisions.
A grill recommendation that asserts current-code behavior names the evidence read in the current session.

## Map content

Use the two templates, then fill the index and each ticket as current state
rather than a deliberation log.
Decisions so far holds one gist line per resolved ticket, with a link to its file.
Notes holds the domain, the skills a session consults, and the standing preferences of that map.

- **Destination** fixes the outcome and scope.
- **Not yet specified** holds honest in-scope fog that is not sharp enough to
  become a decision ticket.
- **Spec-writer discretion** lists only the bounded discretion shaping grants.
- **Out of scope** records reviewed exclusions so later sessions do not reopen
  them.
- **Sources** records structured research objects and the decision each
  supports. An all-grilled map may leave this section empty.

Read one ready decision map's `## Sources` block before the first write. A live
record shows the field grammar. Run `bench maps` and `bench gate-prose` on the
first skeleton. One verb alone leaves one lane unrun.

The map is ready only when its status is `ready`, every decision ticket is
resolved, fog is empty, and every research object remains valid. Run
`bench maps` before declaring readiness. A diagnostic for this map means it is
not ready. A ready map projects a `ready` row and a `/bench-write-spec <path>`
action.

## Starting from the roadmap

`ROADMAP.md` is the working prioritization document `/bench-drain`
maintains. Raw capture lives in `capture/IDEAS.md` and reaches the roadmap only through
a reviewed `/bench-drain` drain.

When invoked cold, read the index `ROADMAP.md`: one heading line per row, no
bodies. Offer its top items, recommended sequence first, and ask which to
pull up. Fetch detail per row only for the ones in play, with
`bench roadmap --context --row <ids>` or from the row's own `roadmap/FT<n>.md`.
When the conversation already carries a fresh idea, proceed with it without
interrupting for the roadmap. If the roadmap is empty or absent, say so. Note
that `/bench-drain` rebuilds it, and continue.

A pulled item keeps its roadmap row in place. Row presence is the
item's status as current open work. The row remains through shaping, spec,
and build until shipped retirement removes it. This command never edits
`ROADMAP.md` or `roadmap/`.

## Two modes

**Bootstrap** (loose idea in): use `craft-grill` to discover whether the idea
really contains a multi-session dependency tree. If it does, create the index
from `bench maps --template` and each ticket from `bench maps --ticket-template`,
then record the frontier. Stop after that first focused shaping pass. If it does
not, create no map. Recommend `/bench-write-spec` from the reviewer-confirmed
current conversation.

**Resume** (map name in): run `bench maps <map>` to resume one map, because
every row names the file to read next. Read the index and that one ticket file,
not the whole folder. Record the answer in the ticket, then add its gist line to
Decisions so far. Add any newly discovered decision tickets with correct
`Blocked by` edges. If an answer invalidates other tickets, update or delete
them.

While the reviewer is present, carry a
Grill straight into newly unblocked decision tickets as the next numbered
frontier round. Do not stop between rounds: never pause for permission or a
re-prompt. Stop when no reviewer-answerable frontier remains or the reviewer
stops answering. Leave Research, Prototype, or Task work explicit, not
force-grilled.

## The exit

Run `bench maps` and re-read every decision ticket, constraint, exclusion,
research object, and discretion item against the current conversation and
tree. A conversation answer not written into the map is not recorded.

When shaped scope contains two independently useful behaviors, the possible
split is a reviewer decision, not an excuse to close shaping early. Surface it
as a decision ticket before one bundled spec becomes the default. Continue
until the reviewer chooses or rejects the split. When the map reaches `ready`, recommend
`/bench-write-spec` from the session holding the ready decision source, with one
clause explaining why the source is ready.
