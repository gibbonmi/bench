---
description: Periodic platform assessment — verify the previous assessment's backlog landed against the tree, fan out read-only area sweeps on the mid tier, synthesize adversarially on the top tier, and produce a dated ASSESSMENT.md that replaces its predecessor with a ranked, agent-time-sized backlog. Deliberately invoked on cadence or on ask; surfaces findings, never merges. Maintenance, not a workflow phase.
disable-model-invocation: true
---

# /bench-assess — assess the platform, replace the assessment

## Entry orientation

This is the periodic assessment phase. Run it on the reviewer's cadence or explicit
ask — never fire it autonomously; it is a deliberately-invoked phase, not a workflow
step and not a build. One run re-baselines the whole platform: it confirms the last
assessment's backlog actually shipped, sweeps the codebase area by area, synthesizes
the findings under adversarial re-verification, and produces one dated `ASSESSMENT.md`
at the repo root for the reviewer to read.

## Exit handoff

Close by reporting the new assessment's headline: how the prior backlog verified, the
count of high/med/low findings, and the top of the ranked backlog. Route findings by
kind — operational items (drain a learning, delete a salvage branch) go to
`/bench-what-next`; new backlog items enter `ROADMAP.md` only through that reviewed
drain, so park them with `bench idea` (or into `IDEAS.md`) rather than editing the
roadmap here. Recommend the next command in this harness's invocation form.

## 1. Verify the previous assessment landed

Before any new sweep, verify the previous assessment's backlog landed — story by
story, against the current tree, not against commit messages. Each prior ranked item
is FIXED, PARTIAL, or open, with a one-line evidence cite; prior high/med findings are
confirmed closed or carried forward. This reconciled baseline is what the new file
builds on, so a shipped item is never re-listed as a fresh finding. Treat a missing
predecessor as a first run — there is nothing to reconcile.

## 2. Fan out the area sweeps (mid tier)

Fan out read-only area sweeps on the mid tier, one delegate per area, each reading
source and running live commands but writing nothing. The six areas:

1. **Adoption and packaging** — the install paths, the first hour, the tarball shape.
2. **Workflow commands and skills** — the phase files, the craft skills, the docs and
   vocabulary that carry the method.
3. **Enforcement layer** — the gate, the hooks, harness parity, the guards.
4. **CLI and Go core** — `bin/bench.sh` and the `internal/` packages behind it.
5. **Gate authority, tests, records** — seam coverage, the canary, ADRs, the
   user-facing changelog, and relevant Git history.
6. **Live operational state** — a real gate run, the roadmap, open learnings, parked
   ideas, stray branches.

The tier→model binding stays reviewer-owned in `projects/<name>.md` per invariant 2;
this command names the tier, not the model.

## 3. Synthesize adversarially (top tier)

Synthesize adversarially on the top tier. A delegate's finding is a claim, not a
result: re-verify every load-bearing claim against source or live output before it
enters the file, and mark the ones you re-verified with a ✓. A claim that survives
only in a delegate's summary does not make the assessment; a claim you could not
verify is recorded as an unknown in the verification notes, not asserted.

## 4. Write the assessment (the output contract)

Produce one dated `ASSESSMENT.md` at the repo root that **replaces its predecessor** —
git history is the archive, so overwrite the file rather than appending; do not keep a
dated series or an archive folder. The file carries:

- **A severity grammar**: **high** = an invariant or advertised guarantee is not
  actually held; **med** = a real defect or reachable unowned state; **low** =
  friction, drift risk, or hygiene.
- **Findings by area**, each cited to the source that proves it.
- **A ranked improvement backlog**, ordered by platform leverage, each item sized in
  rough agent-time.
- **A verification-notes section** recording what was re-verified (✓) and the known
  coverage limits — the unknowns the sweep could not execute or reproduce, named
  honestly rather than hidden.

## 5. Re-confirm the parked, don't re-file it

Items already tracked — a parked roadmap row, an accepted posture in an ADR, an
upstream-blocked feature — are re-confirmed still-correct and cited, not re-filed as
new findings. Re-filing tracked work churns the backlog and hides the real deltas; the
value of each run is the change since the last one.
