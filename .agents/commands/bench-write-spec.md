---
description: Turn a reviewed decision source into a build spec with user stories, engineering seams, and testing decisions. Use before any build past the lighter-path threshold in .bench/BENCH.md.
---

# /bench-write-spec — lock the seams before the loop runs

## Entry orientation

This is the feature-spec phase. It turns reviewed product intent and current
codebase evidence into `specs/<feature>/spec.md`: user stories, implementation
decisions, engineering seams, testing decisions, and the gate that defines
done.

## Exit handoff

Close by showing the approval table for user stories, seams with diagrams,
acceptance coverage, and out of scope. The spec carries `Status: staged`
(staged → implemented at the green gate → promote-then-delete on merge). Stop
for sign-off, then recommend a fresh mid-tier build session on one retained
integration source; review its frozen base and tip, and hand the accepted source
to `bench worktree land`.

## Entry contract

Accept exactly one of three decision sources: a ready compiled map, the
reviewer-confirmed current conversation, or a named reviewed artifact. No
unnamed memory, unreviewed note, or fourth override authorizes a draft.

- **Ready compiled map.** If the source is a ready top-level map, validate it,
  move it and its owned assets into `specs/<slug>/decisions/`, update their
  references, and continue from that compiled path. A re-run reads the
  already-compiled spec-local map; it never recreates a top-level copy.
- **Reviewer-confirmed current conversation.** Confirm that every load-bearing
  product or scope fork is closed in this conversation. Record the date; do not
  manufacture a map to restate it.
- **Named reviewed artifact.** Name the exact approved assessment, findings
  document, issue, or other artifact and verify that it contains settled
  decisions rather than unresolved prompts.

Record exactly one `Decision source:` line in the spec. For a map-backed spec,
name the compiled map path. For the current-conversation source, name the
reviewer confirmation and date. For a reviewed artifact, name that artifact.
The line records provenance; it is not a second research manifest.

For a map-backed source, re-read and re-verify every structured `## Sources`
entry before choosing engineering seams. Confirm repository paths against the
current tree and revisit cited external evidence; disclose anything that could
not be re-read. Consume the map's Sources in place without copying a research
manifest into the spec.

## Late uncertainty

Resolve ordinary late uncertainty here. Ask at most two late clarification
questions, one at a time, each with a recommended answer; route a dependency
tree or multi-session fog to `$bench-shape-idea`. Wait for each reviewer answer
and do not let spec authoring become an unbounded grill.

## Ownership

Spec authoring owns engineering seams, deep-versus-thin design, tests,
acceptance coverage, hostile-input attachment, and gate attachment. Shaping
sources may constrain observable behavior, scope, compatibility, or a
reviewer-chosen architectural seam; they do not pre-author the rest of the
engineering plan.

## Who runs this phase

The session holding the decision source authors the spec and tickets at whatever
tier it runs. A same-session current-conversation source is
authorized by the entry contract, not by inherited memory: verify the reviewer
confirmation and write its provenance line. The fresh implementation session
starts only after ticket approval. Top tier remains an explicit escalation that
pauses for reviewer approval under `craft-line`.

## Process

1. **Read the current state.** Explore the repo and use its vocabulary. Read
   settled ADRs and `projects/<name>.md`. Resolve the one authorized decision
   source under the entry contract. Top-level `decisions/` holds pre-spec
   working maps; compiled maps live under `specs/<slug>/decisions/`. Charge
   `bench-craft-domain` here: pin canonical terms, derive producer partitions,
   and resolve code-versus-claim conflicts before stories lock.

   Before locking stories, read `craft-tickets`, `craft-delegate`, `craft-tdd`,
   and `craft-seams` so the outcomes reflect how implementation will be sliced,
   charged, tested, and fenced.

   When the spec will edit gate-enforced artifacts — anchored clauses, canary
   fixtures, counted or byte-pinned substrings — read the enforcement
   surface's content before locking rows: every fixture `EXPECT` and mutation
   string, every bespoke check that greps or counts the target files, and the
   body of any test cited in an `already covered` red-signal cell. A claim
   built from file, fixture, or test names alone ships wrong coverage.

2. **Derive the engineering seams.** Reconcile the decision source with the
   current repo. Prefer an existing seam and use the highest seam that
   exercises real behavior; fewer seams are better. Use `craft-seams`. State
   each seam for reviewer veto, then draw the template's compact ASCII flow:
   inputs, the unit behind the seam, outputs, trigger, and test-attach point.
   A seam explicitly chosen by the reviewer remains a constraint; every other
   engineering seam is this phase's responsibility.

   After the seams are explicit, lock stories by applying `craft-spec`'s
   canonical `Story sizing and scope cuts` rule.

3. **Price every cut.** Apply `bench-craft-spec`'s breadth floor, derived
   `<n> edits, <n> gate runs` pricing, under-threshold no-deferral rule, and
   separate-capability test. Scope remains the reviewer's decision.

4. **Map acceptance coverage.** For non-trivial work, add an
   acceptance coverage map before implementation. `bench-craft-spec` owns the row schema,
   honest red-signal classifications, degenerate-implementation check, and
   quantifier enumeration. Every observable behavior promised by the decision
   source becomes a row or an explicit exception.

5. **Walk the edge inventory.** Apply `bench-craft-spec`'s canonical edge
   classes and row-or-**Won't handle** rule to every behavior. Consult the
   project's hostile-input checklist. When none covers the surface, quarry
   `.agents/skills/bench-craft-seams/references/hostile-input-library.md` and
   propose a tuned profile addition. Apply `craft-spec`'s named
   `Bootstrap authority before execution` rule.

6. **Route each story.** Give every user story its resolved model and effort
   from `craft-line`, checking the profile's cached routings first. Fully
   gate-observable work routes cheap; prose and semantics the gate cannot grade
   bump a tier.

7. **Write `specs/<feature>/spec.md`.** Use the template below and run
   `bench coverage --check` on the draft. The stale-command-reference sweep
   remains fail-closed across staged specs. If compiling a ready top-level map,
   move (do not copy) the source map and any map-owned assets from top-level
   `decisions/` into `specs/<slug>/decisions/` and update every reference to
   the moved paths in the same green change.

8. **Retire superseded work by promotion then deletion.** Do not leave a
   `Superseded by` marker. Promote durable decisions to their current-state
   owner, then delete the old spec under a `spec-retire: <name>` commit and
   repair every reference. The same
   promote-then-delete commit removes the spec's `ROADMAP.md` row. Git is the
   archive; recover provenance with `bench spec history <slug>`. Whole-folder
   retirement removes the compiled maps and map-owned assets, plus
   implementation tickets, with the spec.

9. **Verification loops before sign-off.** Every draft gets two verification
   loops through the harness's native agent surface: a read-only same-family
   mid-tier delegate at high effort takes the spec alone in loop 1 before
   slicing, then takes the breakdown alone in loop 2 after slicing and checks it
   against `craft-tickets`. Each loop repeats author-fix/re-review uncapped until
   no blocking findings — a recorded reviewer exception to the declared
   iteration cap — with every round reported in one line. The loop verdicts stay
   advisory; reviewer sign-off remains the hard stop.

   `--reviewer <tier-or-model> [effort]` overrides both loop delegates: a tier
   resolves through the invoking harness's own `.bench/lines.env` column
   (`--reviewer mid xhigh` under Codex resolves `BENCH_CODEX_MID` at xhigh), a
   model id must already be bound in `.bench/lines.env` and an unbound id is
   refused, an own-family id runs through the native agent surface, and a
   cross-family id uses the exact recipe in
   `.agents/skills/bench-craft-delegate/references/cross-harness-reviewers.md`.

   Loop 1 asks the falsification questions: would the cheapest wrong
   implementation pass, does every source behavior have a red-capable row, does
   every line match cached routing, does any behavior, red signal, or decision
   answer name an outcome family instead of an exact predicate, are the source
   and observed reds sound even when the source is same-session, conflicting, or
   mostly not observed, and — where the stories partition into disjoint package
   or fence sets — could a narrower capability ship
   on its own gate? Apply `craft-spec`'s named
   `Bootstrap authority before execution` rule.

10. **Slice the implementation tickets.** After loop 1 accepts, charge
    `craft-tickets` and write the breakdown under `specs/<slug>/tickets/`. Carry
    its numbered title, `Blocked by:`, and delivered outcome list into the
    approval table so the spec and tickets receive one sign-off, then run loop 2
    on that breakdown.

    At close, write `Verification log: spec <n> + tickets <m> iteration(s) to
    accept — <note>` into the spec. When either loop takes more than one iteration
    to accept, append one `capture/learnings.md` entry naming the stage that
    missed, what review caught, why it was missed, and the proposed rule change.

## Template

```markdown
# <feature>

Status: staged

Decision source: <one ready compiled map, reviewer-confirmed conversation with date, or named reviewed artifact>

Verification log: spec <n> + tickets <m> iteration(s) to accept — <note>

## Problem
The problem, from the user's point of view.

## Solution
The solution, from the user's point of view.

## User stories
A numbered breadth floor. Each story ends with
`Line: <resolved model id> / <effort>.` and one plain sentence explaining why.

## Implementation decisions
Modules, interfaces, schema or contract changes, and architectural calls.
Record decisions rather than file paths or snippets that rot.

## Testing decisions
- What external behavior a good test exercises.
- Which engineering seams receive tests and their prior art.
- Which gate seam observes the feature.

### Seam diagram

    trigger: <who/what invokes this>
        │
        ▼
    <input>  ──▶  [ <unit behind the seam> ]  ──▶  <output>
                      ◀ tests attach here: <how a test drives and observes it>

### Acceptance coverage map
| row | story | behavior | seam | red signal | why it catches the failure |
|---|---|---|---|---|---|
| <unique spec-local ID> | <story #> | <observable behavior> | <test seam> | <observed red, already covered, or not TDD-able with reason> | <why this fails when behavior is missing> |

### Edge inventory
Every edge class lands in a row above or a
**Won't handle** line: `<edge> — <one-clause why the exclusion is safe>`.

## Ownership fences
List each exact repo-relative file or path prefix a writer may edit. `craft-spec`
owns the fence rules; an empty section is incomplete, not unrestricted authority.

## Out of scope
Each genuine separate capability includes its derived
`<n> edits, <n> gate runs` estimate.
```

Before a build starts, emit a scannable approval table covering stories and
their lines, seam diagrams, acceptance coverage including edge dispositions,
ownership fences with an explicit reviewer disposition, and out of scope. Pause
for sign-off. The user stories set breadth, engineering
seams place tests, and the gate defines done.
