---
name: craft-spec
description: The spec-authoring discipline — the acceptance-coverage-map row schema, what counts as a real red signal, the canonical edge-inventory classes, how stories and scope cuts are sized, and how a build is sliced for delegates. Use when authoring or auditing a spec, writing coverage rows, asking "what's the red signal", walking an edge inventory, judging whether a cut is genuinely out of scope, or slicing a build for delegates.
index: coverage-map rows, edge inventories, story sizing, and delegate slicing for a spec
---

# Spec: synthesize, don't interview

Turn the authorized decision source and what you already know of the codebase
into `specs/<slug>/spec.md` — synthesize, with at most two late questions.

1. **Explore the repo**; use the glossary's terms and respect the area's ADRs.
2. **Sketch the seams** (`craft-seams`): existing over new, the highest that
   still shows the failure, ideally one; confirm them with the reviewer first.
3. **Write the spec** from `bench-write-spec`'s template and run
   `bench coverage --check`. The spec file is the published artifact.

## User stories

A long, numbered, extensive breadth floor, grouped by outcome: one story per
actor-want-benefit — `As an <actor>, I want <feature>, so that <benefit>` — for
every behavior, edge, and reviewed exclusion the source promises. Partial
redundancy is the point: each story restates one behavior downstream slices and
reviews must not drop. A story is a want, never an engineering layer
(`craft-tickets` owns slice sizing); each group carries one `Line:`.

## The acceptance coverage map

Each row ties a story to one observable behavior at a seam: `story`,
`behavior`, `seam`, `red signal`, `why it catches the failure`. An optional
leading `row` column opts the spec into ticket covers traceability (new specs
default to it). The red signal is the command or test already run and failed
because the behavior is absent — say "already covered" or "not TDD-able" rather
than pretend. Name the cheapest wrong implementation per story and the row that
goes red on it (across fences, the composition degenerate through the real
producer); enumerate every quantifier; every source behavior becomes a row or a
stated exception.

## The edge inventory

Walk each behavior through error path, empty/absent input, boundary values,
malformed input, interrupted or partial state, re-run idempotency, process-boundary lifecycle, hostile
environment, plus the profile's hostile-input checklist. Every edge lands in a
row ("edge of N") or a one-line **Won't handle** with a surviving in-scope
caller; a control resolving a class must exercise the **new** surface.

## Bootstrap authority before execution

A trusted-execution or refusal-before-execution claim traces every executable
hop, naming how each validator authenticates the next executable
before launching the next executable. A path, record, digest, or executable
cannot authenticate itself. Without an independent trust root the design is
incomplete; see `references/bootstrap-authority.md`.

## Scope cuts

Price every cut as `<n> edits, <n> gate runs`; a cut must be a separate capability with its own future spec, never "the rest of this feature".

## Slicing a build for delegates

Record **who-writes-where** ownership fences at spec time, checkable at charge
time. A fence entry is an exact repo-relative file or path prefix, never a glob or an
implementation ticket, and an empty or invalid fence section is incomplete.
`craft-tickets` owns the build-time **what-lands-green-next** unit; each ticket
receives the spec's fence. Each fence carries value contracts across it: a contract between tickets is
stated in the ticket's `What to build` and `Acceptance`, re-derived from the
tree by review rather than trusted from the ticket's account.

After a pass touching many sections, reread the complete artifact end to end and reconcile contradictions before handing off.
