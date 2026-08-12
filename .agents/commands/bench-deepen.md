---
description: Architecture-deepening survey — scan for shallow modules and present deepening candidates as a visual HTML report, then grill through whichever one the reviewer picks. Scopes from a named direction, the latest ASSESSMENT.md, or commit-history hot spots. Deliberately invoked; surfaces candidates, never refactors. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-deepen — survey deepening opportunities, grill the pick

## Entry orientation

This is the deepening-survey phase. Run it on the reviewer's ask — it is a
deliberately-invoked phase, not a workflow step and not a build. One run surfaces
architectural friction as **deepening opportunities** — refactors that turn shallow
modules into deep ones, aimed at testability and agent-navigability — renders them as
a visual report, and walks the reviewer's pick through a grilling loop. It proposes;
it never refactors in-phase.

Charge `craft-seams` before exploring: it owns the vocabulary (module, interface,
depth, seam, adapter, leverage, locality), the deletion test, and "the interface is
the test surface". Use those terms exactly in every candidate — never "component,"
"service," "API," or "boundary." Domain names come from `CONTEXT.md` per
`craft-domain`; ADRs record decisions this phase does not re-litigate.

## 1. Scope before you scan

Deepening pays off where change keeps happening, so decide *where* to look before
looking:

- A reviewer-named direction — a module, a subsystem, a pain point — wins; take it
  and skip the inference below.
- Otherwise, if an `ASSESSMENT.md` exists at the repo root, take its
  architecture-shaped findings and backlog rows as candidate territory — this is the
  `/bench-assess` conjunction: the audit finds the friction, this phase designs the
  deepening.
- Otherwise, walk back a good stretch of `git log --oneline` for hot spots — the
  files and areas that keep coming up. Scattered changes with no hot spot widen the
  net.

Read `CONTEXT.md` and the ADRs for the area first.

## 2. Explore (mid tier, read-only)

Spawn one read-only delegate on the mid tier (the tier→model binding stays
reviewer-owned in `projects/<name>.md` per invariant 2) to walk the scoped area
organically — no rigid heuristics — noting friction:

- Where does understanding one concept require bouncing between many small modules?
- Where are modules shallow — interface nearly as complex as the implementation?
- Where were pure functions extracted just for testability while the real bugs hide
  in how they're called (no locality)?
- Where do tightly-coupled modules leak across their seams?
- What is untested, or hard to test through its current interface?

Apply the deletion test from `craft-seams` to anything suspected shallow. A
delegate's finding is a claim: before a candidate enters the report, confirm against
source that the named files exist and the shallowness reads as described.

## 3. Present candidates as an HTML report

Write one self-contained HTML file to the OS temp directory (`$TMPDIR`, falling back
to `/tmp`), named `architecture-review-<timestamp>.html`, open it (`xdg-open` /
`open` / `start`), and report the absolute path — nothing lands in the repo. The
scaffold, card layout, diagram patterns, and report tone live in
`.agents/skills/bench-deepen/HTML-REPORT.md`; follow it. Each candidate card carries
files, problem, solution, wins framed as locality and leverage, a before/after
diagram, and a recommendation-strength badge; the report ends with a top
recommendation.

A candidate that contradicts an ADR appears only when the friction is real enough to
warrant reopening the decision, marked with a callout naming the ADR; the rest stay
unlisted. Propose no concrete interfaces yet. Then ask the reviewer which candidate
to explore.

## 4. Grill the pick

Walk the chosen candidate with the reviewer through `craft-grill` frontier rounds —
constraints, dependencies, the shape of the deepened module, what sits behind the
seam, what tests survive. Side effects land inline as decisions crystallize, per
`craft-domain`: a new or sharpened concept updates `CONTEXT.md` there and then; a
candidate rejected for a load-bearing reason earns one offer of an ADR (`craft-adr`)
so a future survey doesn't re-suggest it — skip ephemeral or self-evident reasons.
For a genuinely uncertain interface, use design-it-twice from `craft-seams`.

## Exit handoff

A settled deepening becomes work through the normal workflow, right-sized per
`.bench/BENCH.md` — the light path for a one-ticket deepening, `/bench-write-spec`
past that threshold; recommend the next command in this harness's invocation form.
Candidates worth keeping but not pursuing now are parked with `bench idea`, never
written into `ROADMAP.md` here.
