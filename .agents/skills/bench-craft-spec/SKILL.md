---
name: craft-spec
description: The spec-authoring discipline — the acceptance-coverage-map row schema, which edges a spec must dispose of, how stories and scope cuts are sized, and how a build is sliced for delegates. Use when authoring or auditing a spec, writing coverage rows, deciding an edge is a Won't handle, judging whether a cut is genuinely out of scope, or slicing a build for delegates.
index: coverage-map rows, edge inventories, story sizing, and delegate slicing for a spec
---

# Spec: synthesize, don't interview

Turn the authorized decision source and what you know of the codebase into
`specs/<slug>/spec.md` — synthesize, with at most two late questions.

1. **Explore the repo**; use the glossary's terms and respect the area's ADRs.
2. **Sketch the seams** (`craft-seams`): existing over new, the highest that
   still shows the failure, ideally one; confirm them with the reviewer first.
3. **Write the spec** from the template below, in ASD-STE100 prose per
   `references/ste-prose.md`, and run `bench coverage --check`. The spec file
   is the published artifact.

## User stories

Write a long, numbered list grouped by outcome, with an extensive breadth floor. One story per
actor-want-benefit — `As an <actor>, I want <feature>, so that <benefit>` — covers every behavior, edge,
and reviewed exclusion the source promises. Partial redundancy is the point. A story is a want, never an
engineering layer (`craft-tickets` owns slice sizing). Each group carries one `Line:`.

## The acceptance coverage map

Each row ties a story to one observable behavior at a seam: `story`, `behavior`, `seam`, `why it catches the failure`.
An optional leading `row` column opts the spec into ticket covers traceability (new specs
default to it). `bench coverage --check` refuses a row that references more than four stories, and it
refuses a row that states two predicates (`;`). It also refuses a declared story that no row references,
unless a `Not covered: story <n> — <reason>` line sits under the map.
Name the cheapest wrong implementation per story. Name the row that goes red on it, across fences and through the composition
degenerate to the real producer.

Enumerate every quantifier; every source behavior becomes a row or an exception.

## The edge inventory

`craft-tdd` walks the canonical edge classes at the seam. Attach the profile's hostile-input checklist to
that walk. Give each edge the reviewer deliberately excludes a one-line **Won't handle** with a surviving
in-scope caller.

## Bootstrap authority before execution

A trusted-execution or refusal-before-execution claim traces every executable hop, naming how each
validator authenticates the next executable before launching the next executable. A path, record, digest,
or executable cannot authenticate itself. Without an independent trust root the design is incomplete; see
`references/bootstrap-authority.md`.

## Scope cuts

Price every cut as `<n> edits, <n> gate runs`. A cut must be a separate capability with its own future
spec, never "the rest of this feature".

## Slicing a build for delegates

Record **who-writes-where** ownership fences at spec time, checkable at charge
time. A fence entry is an exact repo-relative file or path prefix, never a glob
or an implementation ticket. An empty or invalid fence section is incomplete.

`craft-tickets` owns the build-time **what-lands-green-next** unit; each ticket receives the spec's fence. Each
fence carries value contracts across it. A contract between tickets is stated in the ticket's `What to build`
and `Acceptance`. Review re-derives that contract from the tree; it does not trust the ticket's account.

After a pass touching many sections, reread the complete artifact end to end and reconcile contradictions
before handing off.

## Review rubric

The round asks five questions. Would the cheapest wrong implementation pass? Does every source behavior
have a red-capable row? Does every line match cached routing? Does any behavior, why-it-catches clause,
or decision answer name an outcome family instead of an exact predicate? Are the source and observed reds
sound even when the source is same-session, conflicting, or mostly not observed?

The degenerate standard is the cheapest plausible wrong implementation — a degenerate that needs deliberate
contrivance is the build's mutation-probe target, never a new spec row. A finding blocks only when it changes
observable behavior, an ownership fence, or the ticket graph. A round that returns only prose or accounting
findings is the acceptance round. Fold those fixes into the acceptance instead of another round. A revision
may not add a promise beyond the decision source unless a blocking finding demands it. The review flags an
unflagged addition for removal rather than demanding rows for it.

## Template

```markdown
# <feature>

Status: staged

Decision source: <one ready compiled map, reviewer-confirmed conversation with date, or named reviewed artifact>

Introduces commands: <optional — the `/bench-…` and `$bench-…` phases this spec ships, valid in its own directory while staged>

Verification log: <n> iteration(s) to accept — <note>

## Problem
The problem, from the user's point of view.

## Solution
The solution, from the user's point of view.

## User stories
A long, numbered, extensive breadth floor grouped by outcome — one
`As an <actor>, I want <feature>, so that <benefit>` per behavior, edge, and
reviewed exclusion, partially redundant on purpose. Each group opens with its
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
| row | story | behavior | seam | why it catches the failure |
|---|---|---|---|---|
| <unique spec-local ID> | <story #> | <observable behavior> | <test seam> | <why this fails when behavior is missing> |

### Edge inventory
Each deliberately excluded edge takes a
**Won't handle** line: `<edge> — <one-clause why the exclusion is safe>`.

## Ownership fences
List each exact repo-relative file or path prefix a writer may edit. `craft-spec`
owns the fence rules; an empty section is incomplete, not unrestricted authority.

## Out of scope
Each genuine separate capability includes its derived
`<n> edits, <n> gate runs` estimate.

## Further notes
```

Before a build starts, emit a scannable approval table. The table covers
stories and their lines, seam diagrams, acceptance coverage including edge
dispositions, ownership fences with an explicit reviewer disposition, and out
of scope. Pause for sign-off. The user stories set breadth, engineering seams
place tests, and the gate defines done.

