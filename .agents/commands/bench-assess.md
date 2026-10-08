---
description: Periodic platform assessment — verify the previous backlog against the tree, fan out read-only sweeps on the mid tier, synthesize adversarially on the top tier, and produce an ASSESSMENT.md. Surfaces findings, never merges. Maintenance, not a phase.
disable-model-invocation: true
---

# /bench-assess — assess the platform, replace the assessment

## Entry orientation

This is the periodic assessment phase. Run it on the reviewer's cadence or on an
explicit ask; never fire it on your own. It is a deliberately-invoked phase, not a
workflow step and not a build. One run re-baselines the whole platform. It
confirms that the last assessment's backlog actually shipped and sweeps the
codebase area by area. It synthesizes the findings under adversarial
re-verification and produces one dated `ASSESSMENT.md` at the repo root for
the reviewer to read.

At entry, record the repository, the source commit, the scope, and the requested
stopping point. The default stopping point is the assessment and its
recommendations. Neither an assessment nor a planning approval starts an
implementation.

## Exit handoff

Report the new assessment's headline: how the prior backlog verified, the count of
high, medium, and low findings, and the top of the ranked backlog. Route findings by
kind. An operational item — drain a learning, delete a salvage branch — goes to
`/bench-drain`. A new backlog item enters `ROADMAP.md` only through that reviewed
drain, so park it with `bench idea` rather than editing the roadmap here.

An architecture-shaped finding — a shallow module, a leaky seam, friction the
sweeps hit while reading the code — additionally feeds `/bench-deepen`. That
phase scopes its deepening survey from this file's findings. Recommend it
when the run surfaces one. Recommend the next command in this harness's
invocation form. When the reviewer asks for planning follow-through, run step 6
before the exit.

## 1. Verify the previous assessment landed

Before any new sweep, verify the previous assessment's backlog landed. Compare it
story by story with the current tree, not with commit messages. Mark each
prior ranked item FIXED, PARTIAL, or open, with a one-line evidence cite.
Confirm each prior high or medium finding closed, or carry it forward. This reconciled
baseline is what the new file builds on, so you never re-list a shipped item as a
fresh finding. Treat a missing predecessor as a first run; there is nothing to
reconcile.

## 2. Fan out the area sweeps (mid tier)

Fan out read-only area sweeps on the mid tier, one delegate per area. Charge
`craft-research` for the round, the charges, and the citations. The six areas:

1. **Adoption and packaging** — the install paths, the first hour, the tarball shape.
2. **Workflow commands and skills** — the phase files, the craft skills, the docs and
   vocabulary that carry the method.
3. **Enforcement layer** — the gate, the hooks, harness parity, the guards.
4. **CLI and Go core** — `bin/bench.sh` and the `internal/` packages behind it.
5. **Gate authority, tests, records** — seam coverage, test cost, the canary, ADRs,
   the user-facing changelog, and relevant Git history.
6. **Live operational state** — a real gate run, the roadmap, open learnings, parked
   ideas, stray branches.

The tier-to-model binding stays reviewer-owned in `projects/<name>.md` per
invariant 2. This command names the tier, not the model.

The area 2 charge separates conflicting authority from useful constraints and
lookup detail. A shorter guide is not by itself a better guide.

The area 5 charge assesses test cost through test purpose:

- Read the retained gate logs, timing records, retros, and approved
  test-efficiency specs first. Record the sample dates, source identities,
  conditions, and missing measurements. Run no manual benchmark that the
  reviewer excluded.
- For each expensive test family, name its observable failure and its production
  entry point. Separate fixture construction, repeated dispatch, source
  acquisition, and redundant setup.
- A cheaper test keeps the same failure witness at the right seam. Keep real
  journeys, independent omission expectations, failure assertions, and mutation
  freshness.
- Package durations overlap, so their savings do not add. A shared helper alone
  proves no speedup. Record an unmeasured cost as unknown.
- A speed gain inside the recorded run-to-run spread is no change. Record each
  reverted attempt, so that a later run does not try it again.

## 3. Synthesize adversarially (top tier)

Synthesize adversarially on the top tier. `craft-research` owns the
re-verification, the citations, and the residual unknowns.

Each finding keeps the `craft-research` separation of fact, inference, and tested
result. A defect names its wrong observable result and its smallest
discriminating regression. Keep a defect repair apart from a behavior-preserving
refactor. Correct a stale roadmap premise from current code before you recommend
its implementation.

A removal or cut finding names the consumers of the surface inside and outside
the tree, its replacement, and the reason the surface exists. An in-tree search
that finds no caller does not prove that no consumer exists.

## 4. Reconcile with existing work, then rank

Read each active roadmap row, each relevant staged spec and decision map, and
each literal dependency. Map each finding to its existing owner before you
propose new work. Classify the overlap as complete, partial, separate, or
unresolved. Reuse an approved spec or a closed decision instead of a competing
plan.

Re-confirm an already tracked item — a parked roadmap row, an accepted posture in
an ADR, an upstream-blocked feature — as still correct, and cite it. Do not
re-file it as a new finding. A re-filed tracked item churns the backlog and
hides the real deltas. The value of each run is the change since the last one.

Take a fresh priority view; the current priority labels do not decide it. Rank by
consequence, evidence, recurrence, test value, effort, and prerequisites.
Recurrence includes repeated fix commits in one area, read from `git log`. Order
the work by the purpose order in `.bench/BENCH.md`. Keep each literal blocker
apart from the proposed order.

## 5. Write the assessment (the output contract)

Produce one dated `ASSESSMENT.md` at the repo root that **replaces its predecessor**.
Git history is the archive, so overwrite the file rather than append to it. Do not keep
a dated series or an archive folder. The file carries:

- **A source pin**: the repository, the source commit, the scope, and the
  stopping point from the entry.
- **A severity grammar**: **high** means an invariant or advertised guarantee does
  not actually hold; **med** means a real defect or a reachable unowned state;
  **low** means friction, drift risk, or hygiene.
- **Findings by area**, each cited to the source that proves it, with its existing
  owner and its overlap class.
- **A ranked improvement backlog**, ordered by the step 4 factors, each item sized
  in rough agent-time.
- **A roadmap disposition**: the proposed next action for each active row, and each
  parked or paused row with its return condition.
- **A verification-notes section** that records what you re-verified (✓) and the
  known coverage limits. Name the unknowns the sweep could not execute or
  reproduce honestly rather than hide them.

## 6. Planning follow-through (on request only)

Run this step only when the reviewer asks for it. Hand a multi-session decision
tree to `/bench-shape-idea` and an approved outcome to `/bench-write-spec`, each
as its own phase. Reuse or amend a spec that already owns the outcome; a new map
does not reopen an approved design. Stop at the authorized planning stage.

Report the exact state of each outcome: assessed, mapped, spec-reviewed,
tickets-reviewed, or build-authorized.

**Good:** "The map is ready; the spec and implementation tickets are not written."
This report separates settled decisions from executable work.

**Bad:** "Ready to build," after only a map check. This report claims later
stages without their artifacts or their approval.
