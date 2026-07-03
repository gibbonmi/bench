---
description: Turn the current conversation into a build spec — user stories, pre-agreed test seams, and testing decisions. Use before any build that is more than a trivial change. No interview; synthesize what we've already discussed.
---

# /bench-write-spec — lock the seams before the loop runs

## Entry orientation

This is the feature-spec phase. It turns the current conversation, decision map,
and codebase context into `specs/<feature>.md`: user stories, implementation
decisions, pre-agreed seams, testing decisions, and the gate that defines done.

## Exit handoff

Close by showing the approval table for user stories, seams (each with its seam
diagram), acceptance coverage, and out of scope. Stop for reviewer sign-off. After approval, the recommended next
command is `/bench-implement-spec` for an interactive build, or `bench shift` only
when the spec is locked and it passes `craft-line`'s venue-routing test: every
story's line is cheap and the coverage map is fully gate-observable.

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
   present, plus any `projects/<name>.md`. `decisions/` holds working maps, not
   settled records: treat an open map as questions to respect, never as decisions
   already made.

2. **Pick the seams.** Sketch where this feature will be tested. Prefer an
   existing seam to a new one. Use the highest seam that exercises the real
   behavior — the fewer seams, the better; one is ideal. (See the `craft-seams` skill.)
   State the seams explicitly and check they match my expectation before writing
   the spec. Then make each seam legible: draw the small ASCII data-flow the
   template's seam diagram section shows — inputs, the unit the seam fronts,
   outputs, the trigger that invokes it, and the marked test-attach point — so I
   can veto a wrong seam by looking at a picture instead of reconstructing the
   flow from prose.

3. **Price every cut, in your time — not a human's.** Before anything is deferred,
   estimate it in *agent* time, because that's the real cost here. The instinct to
   defer is calibrated to humans who can't spare the afternoon; you can. Derive the
   estimate instead of guessing: agent time is dominated by verification, so state
   it as `<n> edits, <n> gate runs` — a vibes number can't pass as a price. Two
   rules follow, and they are the point of this step:
   - **No deferral under the threshold — this binds you, not me.** Anything under
     ~30 minutes of your work that introduces no new architectural decision is
     something you do *not* get to propose deferring — it's just part of this build,
     so do it now and state the estimate out loud. But scope is mine: if *I* choose
     to defer something small for my own reasons, record it in Out of scope with its
     estimate and no argument. The rule stops you pitching the dodge; it never
     overrides my call.
   - **A cut must be a separate capability, not the rest of this one.** Something is
     only legitimately out of scope if it has its own future spec — a distinct
     feature. If it's "the rest of *this* feature" (the error cases, the other half
     of the CRUD, the edge handling), it is not out of scope; it's an acceptance
     criterion. Move it into the user stories and testing decisions so the **gate**
     enforces it, rather than leaving it in a prose aside that disappears. You can't
     ticket your way out of the spec's own breadth.

4. **Map acceptance coverage.** For non-trivial feature work, add an
   **acceptance coverage map** to Testing decisions before implementation begins.
   Each row ties a user story to one observable behavior at a chosen seam, with
   these fields: `story`, `behavior`, `seam`, `red signal`, and
   `why it catches the failure`. The red signal is the command or test that has
   already been run and failed because the mapped behavior is absent or wrong. If
   the behavior is already covered or cannot start red, say so in the row instead
   of pretending it is TDD coverage. When a behavior or red-signal promise
   quantifies over a set ("each check", "every parser"), enumerate the set or
   state the granularity explicitly — per item or per class. An unenumerated
   "each" lets the build pick the cheapest reading, and review is the wrong
   place to catch that.

5. **Walk the edge inventory.** Stories are happy-path shaped; this step generates
   the cases nobody declared. For each mapped behavior, walk the edge classes —
   error path, empty/absent input, boundary values, malformed input,
   interrupted/partial state, re-run idempotency, hostile environment (this is
   the canonical edge-class list; `craft-tdd` and `/bench-review-implementation`
   point here) — and
   consult the project
   profile's hostile-input checklist (`projects/<name>.md`) when one exists, so
   domain edges recur instead of being rediscovered per defect. When the profile
   has no checklist for the surface being touched, quarry the kit's hostile-input
   library (`.agents/skills/bench-craft-seams/references/hostile-input-library.md`)
   and propose adding the tuned section to the profile. Every edge lands
   in exactly one of two places: a coverage row (story column may read "edge of
   N"), or a one-line **Won't handle** entry directly under the map. Both are
   veto surface; a silently untested edge is the failure this step exists to
   prevent. Before writing a **Won't handle** line about an interface, verify at
   least one in-scope caller can still exercise the feature under that exclusion —
   a cut that amputates the surface's primary calling convention is a spec defect,
   not a scope cut.

6. **Route each story.** Give every user story its line — the resolved model id
   and effort from the `craft-line` decision table, judged per story on spec
   precision, seam certainty, and gate coverage. Work the gate fully observes
   routes cheap; prose and semantics the gate can't grade bump a tier. This is
   the spec-time half of invariant #2: the build inherits routing I approved,
   instead of picking models mid-loop.

7. **Write `specs/<feature>.md`** using the template below.

## Template

```markdown
# <feature>

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
