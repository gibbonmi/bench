---
description: Turn the current conversation into a build spec — user stories, pre-agreed test seams, and testing decisions. Use before any build past the lighter-path threshold in .bench/BENCH.md. No interview; synthesize what we've already discussed.
---

# /bench-write-spec — lock the seams before the loop runs

## Entry orientation

This is the feature-spec phase. It turns the current conversation, decision map,
and codebase context into `specs/<feature>/spec.md`: user stories, implementation
decisions, pre-agreed seams, testing decisions, and the gate that defines done.

## Exit handoff

Close by showing the approval table for user stories, seams (each with its seam
diagram), acceptance coverage, and out of scope. The spec carries a
`Status: staged` line (staged → implemented at the green gate → promote-then-delete
on merge). Stop for reviewer sign-off. After approval, recommend a
**new session on the mid tier by default** to orchestrate the build — a fresh session keeps the
build off the top model's big-context iteration, escalating per stage instead. That
session runs
`/bench-implement-spec` for an interactive build, or `bench shift` only when the
spec is locked and it passes `craft-line`'s venue-routing test: every story's line
is cheap and the coverage map is fully gate-observable.

**This phase refuses to run without a complete map.** A new spec is compiled from
a named top-level `decisions/<topic>.md` whose `## Handoff` is placeholder-free —
run `bench maps` and confirm it shows no row for the map. A revision to a live
spec reads its compiled map under `specs/<slug>/decisions/`; settled provenance
there is deliberately outside the top-level `bench maps` query. On a missing
map, or one whose Handoff is still open, name the map to close and stop; never
draft a spec from conversation alone. One override exists: an explicit reviewer-directed
batch drain (an assessment or reviewed findings doc pushed into specs on the
reviewer's instruction) may substitute for per-spec maps, with every defaulted
decision flagged in-spec for post-hoc veto — absent that explicit instruction,
the map gate stands. Default spec authoring starts in a fresh mid-tier session
and reads the complete map plus the repo rather than inheriting the shaping
conversation. The sole same-session exception applies when every load-bearing
fork has already been put to the reviewer and closed in the current session:
write those decisions directly into a new decision map with a complete Handoff,
then continue from that file rather than unwritten grill memory, flag the map in
the spec for reviewer veto, and compile the spec without routing through
`/bench-shape-idea`. This path records settled decisions rather than bypassing
them: if any fork remains open, run `/bench-shape-idea` and keep the normal map
gate. Bugs are the exception: they route to `/bench-debug`, which needs no spec.

**Who runs this phase.** The mid tier authors ordinary specs. The top tier is an
explicit exception only when the Handoff carries uncertainty flags and I approve
the escalation. Ordinary authoring stays in a fresh session. The reviewer-closed
path in the entry contract is the sole same-session exception, and it writes the
complete map before continuing from that file. Either way, write from the map
plus the repo, and treat any reach into unwritten grill memory as a missing
Handoff fact. A same-session mid **delegate** is allowed only on my explicit ask:
it is read-only (or worktree-isolated) and returns the spec text for the invoking
session to write.

Synthesize the current conversation and codebase understanding into a spec the
build loop can run against. Do **not** interview me — use what you already know.
If something load-bearing is genuinely missing, ask one question, then proceed.

The point of this command is to decide the **seams and the tests before any code
is written**, so that when `/bench-implement-spec` runs TDD it tests at a seam I chose, against a
notion of "correct" I defined — not one the agent invents mid-loop and then
over-fits to. This is the single most important step for keeping an autonomous
loop honest.

## Process

1. **Read the current state.** Explore the repo. Use the project's vocabulary.
   Respect settled ADRs in the area you're touching — read `docs/adr/` when
   present, plus any `projects/<name>.md`. Top-level `decisions/` holds
   pre-spec working maps: treat unresolved questions with respect, never as
   decisions already made. A map under `specs/<slug>/decisions/` is settled provenance compiled
   with that spec, not a shaping frontier.

2. **Pick the seams — off the map's Handoff first.** If a decision map produced
   this spec, read its `## Handoff` section and take the seams from there: items
   1–6 name the module boundaries, contracts, deep-vs-thin split, black-box
   assertables, gate attachment, and hostile-input owners. Map-sourced seams are
   **pre-agreed** — approved at the map's close — so verify each still matches the
   current repo, but pause for my sign-off only on a seam you had to invent (the
   map was silent), a seam the Handoff flagged uncertain (escalate per the
   `craft-line` ladder), or a deliberate deviation from what the map named. A spec
   always has a map behind it (the entry contract above), so there is no
   from-scratch source — name the map to close and stop instead. Then:
   sketch where this feature will be tested. Prefer an
   existing seam to a new one. Use the highest seam that exercises the real
   behavior — the fewer seams, the better; one is ideal. (See the `craft-seams` skill.)
   State the seams explicitly and check they match my expectation before writing
   the spec. Then make each seam legible: draw the small ASCII data-flow the
   template's seam diagram section shows — inputs, the unit the seam fronts,
   outputs, the trigger that invokes it, and the marked test-attach point — so I
   can veto a wrong seam by looking at a picture instead of reconstructing the
   flow from prose.

3. **Price every cut, in your time — not a human's.** The sizing discipline —
   the breadth floor, the derived `<n> edits, <n> gate runs` pricing, the
   under-threshold no-deferral rule, and the separate-capability cut test — is
   `bench-craft-spec`'s; apply it to every proposed cut. Scope stays mine: if
   *I* choose to defer something small for my own reasons, record it in Out of
   scope with its estimate and no argument — the rules stop you pitching the
   dodge; they never override my call.

4. **Map acceptance coverage.** For non-trivial feature work, add an
   **acceptance coverage map** to Testing decisions before implementation begins.
   The row schema and the red-signal definition are `bench-craft-spec`'s — that
   skill is the one source of the five fields (including the `red signal` and
   `why it catches the failure` columns), the honest-classification rule, the
   degenerate-implementation check, the quantifier-enumeration rule, and the
   Handoff-assertables-to-rows rule; compose it and apply them here.

5. **Walk the edge inventory.** The canonical edge classes, the
   row-or-**Won't handle** landing rule, the amputated-caller guard, and the
   compatibility-probe rule are `bench-craft-spec`'s — walk them for each
   mapped behavior. Consult the project profile's hostile-input checklist
   (`projects/<name>.md`) when one exists, so domain edges recur instead of
   being rediscovered per defect; when the profile has no checklist for the
   surface being touched, quarry the kit's hostile-input library
   (`.agents/skills/bench-craft-seams/references/hostile-input-library.md`) and
   propose adding the tuned section to the profile.

6. **Route each story.** Give every user story its line — the resolved model id
   and effort from the `craft-line` decision table, judged per story on spec
   precision, seam certainty, and gate coverage. Check the profile's `Lines`
   cached routings first: when a cached row matches the work, the story's line
   matches it, or the story names the deviation and why. Work the gate fully observes
   routes cheap; prose and semantics the gate can't grade bump a tier. This is
   the spec-time half of invariant #2: the build inherits routing I approved,
   instead of picking models mid-loop.

7. **Write `specs/<feature>/spec.md`** using the template below, then run
   `bench coverage --check` on the draft, so map-format defects surface at
   author time instead of at the gate. A spec whose deliverable is a new
   `/bench-*` phase command lands in the same diff as the command, because the
   stale-command-reference sweep remains fail-closed across staged specs rather
   than exempting them. In that same green change, move (do not copy) the source
   map and any map-owned assets from top-level `decisions/` into
   `specs/<slug>/decisions/`, preserving their useful relative layout, and
   update every reference to the moved paths. A re-run reads the already-compiled
   spec-local map; it never recreates a top-level copy.

8. **Retire what this spec supersedes — promote, then delete.** If the new spec
   replaces an existing `specs/*/spec.md` (same feature, new direction), it does **not**
   get a **Superseded by** marker left in place — a superseded spec left live reads
   as a second source of truth. Instead, in the same change: promote anything still
   durable (a decision → an ADR, a hostile edge → the profile checklist, a seam →
   the profile seam list), then **delete** the old file under a `spec-retire: <name>`
   commit, fixing every dangling reference in that same commit so the stale-reference
   sweep stays green. The same promote-then-delete commit removes the spec's
   `ROADMAP.md` row — row presence is status, so a retired spec that leaves its row
   behind keeps shipped work listed as current. git is the archive
   (`git log --grep=spec-retire`), so no marker and no archive folder is kept —
   recover a retired spec's origin with `bench spec history <slug>` rather than
   hand-running that query. Retiring the old spec is part of writing the new one,
   not a cleanup for later.
   The same pass applies when retiring a merged, implemented spec flagged by
   `bench status`. Whole-folder retirement removes the compiled maps and
   map-owned assets under `specs/<slug>/decisions/` with the spec and its tickets;
   there is no separate decision-map cleanup step.

9. **Falsification review before sign-off.** Every draft gets the pass: spawn a
   reviewer sub-agent before the approval table, unconditionally. The pass
   **runs at the mid tier**, read-only, one iteration — at that binding it is
   not a top-tier bump, so it spawns without asking; the project's `Lines`
   carries the routing. These signals argue for running it at the top tier
   instead:

   - the Handoff carries uncertainty flags (item 7);
   - the draft deviates from the map;
   - **no map backs the draft** — the reviewer-directed batch-drain override in the
     entry contract. This is the path with the most undecided content and the
     least prior sign-off;
   - **the map was written in the same session as the draft** — the
     reviewer-closed path in the entry contract. The author is checking its own
     recall, which is the bias the fresh-session default exists to remove;
   - **the coverage map's red signals are mostly not observed reds** — when
     `already covered` and `not TDD-able` outnumber rows with a real red
     command, the map can pass its format check while grading nothing.

   A fired signal never escalates on its own: when one fires and you judge the
   top tier worth it, pause and ask me before spawning — a top-tier pass is an
   ordinary bump under `craft-line`'s ladder, with no standing opt-out. Absent
   my yes, run the pass at the mid tier anyway. A reviewer sub-agent is never a
   licence to spend a tier the project has not granted.

   Give it a fresh small context — the Handoff (or, with no map, the decided
   scope and the finding inventory it was compiled from) plus the draft, nothing
   else — and charge it with falsification questions, never an open "review
   this": would the cheapest wrong implementation pass this map? does every
   Handoff assertable, or every finding, have a row that would actually go red if
   it were left unfixed? does each line match the cached routings? On the
   no-map and same-session paths, charge it additionally at the defaults —
   are the decisions the draft marked for veto the *only* ones it decided
   unilaterally? It returns findings and an advisory recommend/block verdict.

   The verdict is advisory: sign-off stays mine. Verify the reviewer's findings
   against the tree before folding them in — a delegate's finding is a claim,
   not a result.

## Template

```markdown
# <feature>

Status: staged

## Problem
The problem, from the user's point of view.

## Solution
The solution, from the user's point of view.

## User stories
A long, numbered list. "As a <actor>, I want <feature>, so that <benefit>."
Exhaustive — this list defines the breadth the build must cover, which is the
guard against a loop that does the minimum and stops.
Every story displays its model and effort: end the story with
`Line: <resolved model id> / <effort>.` followed by one whole, plain sentence
explaining why that row was chosen. The sentence must read on its own — no
stacked clauses or fragments the reviewer has to decode.

## Implementation decisions
Modules touched, interfaces modified, schema/contract changes, architectural
calls. Decisions, not file paths or snippets — those rot. Exception: a prototype
snippet that encodes a decision more precisely than prose (a state machine, a
schema, a type) may be inlined, trimmed to the decision-rich part.

## Testing decisions
- What a good test is here: exercise external behavior at the seam, not internals.
- Which seams get tested, and the prior art (similar tests already in the repo).
- The gate command this feature must pass (defaults to the project gate).

### Seam diagram
One small ASCII data-flow per tested seam — the seam must be legible as a
picture before the build is approved. Show what flows in, the unit the seam
fronts, what flows out, who or what triggers it, and mark the test-attach point:

    trigger: <who/what invokes this>
        │
        ▼
    <input>  ──▶  [ <unit behind the seam> ]  ──▶  <output>
    <input>  ──▶  [                        ]
                      ◀ tests attach here: <how a test drives and observes it>

Keep it to one screen per seam; a diagram that needs scrolling is usually a seam
placed too low.

### Acceptance coverage map
| story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|
| <story #> | <observable behavior> | <test seam> | <observed red command/test, or already covered / not TDD-able with reason> | <why this signal would fail if the behavior were missing> |

### Edge inventory
The edge classes walked per behavior, each resolved as a coverage row above or a
**Won't handle** line here: `<edge> — <one-clause why it's safe to skip>`.
Exclusions are decisions on the page, not silent omissions.

## Out of scope
Each genuine cut as one line: **what** it is — why it's a *separate capability*
(not just the rest of this feature) — **your** derived estimate to build it later
(`<n> edits, <n> gate runs`).
Anything you can't defend as a separate capability, or that falls under the
threshold, does not belong here — it goes into the user stories above. An empty
section is a fine and common answer; a long one is a signal you're shrinking the
target. (Borrows Pocock's `triage` OUT-OF-SCOPE convention: exclusions are
decisions on the page, not silent omissions.) Rank any high-value cuts worth
preserving, and park concrete future features on the roadmap when they should not
disappear after this spec.
```

Every line in this section is something I can read and veto in one pass — which is
the whole purpose. A deferral with a 15-minute estimate next to it usually argues
against itself.

When the spec is written, the build has a fixed target: the user stories set the
breadth, the seams set where tests live, and the gate sets what "done" means.

Before any build starts, emit a scannable approval table — user stories (each
with its line) / seams (with their seam diagrams) / acceptance coverage (edge
rows and won't-handle lines included) / out of scope —
and pause for my sign-off. The full spec file
stays as written; the table is the at-a-glance veto surface. Only after I approve,
lead with the recommended next action — `/bench-implement-spec` interactively, or a
`bench shift` when the spec passes `craft-line`'s venue-routing test — and a
one-clause why, not a neutral either/or.
