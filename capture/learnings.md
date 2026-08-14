# Learnings — usage journal

## 2026-08-14 — write-spec falsification loop took 33 rounds on worktree-enumeration-hang

What happened: /bench-write-spec's loop 1 needed 33 author-fix/re-review rounds
before ACCEPT (tickets loop: 2). The author drafted from the compiled map plus
targeted code reading; the reviewer's tree-verification kept finding real
holes — the linked-worktree producer input, callers that swallow or relabel
the new typed error, per-child bound scoping, process-status partition arms,
and fixture reachability at command surfaces.
Right behavior: most catches trace to one root: the draft reasoned from the
producing seam outward and enumerated consumers late. Enumerating every
caller/renderer of the changed contract *before* writing coverage rows —
call-graph first, rows second — would have removed roughly half the rounds
(list.go, resume.go, dashboard, PruneLandedBranches were each late finds).
Proposed rule change: none — a proposal was made and empirically rejected.
The candidate craft-spec sentence ("enumerate every production consumer of a
changed contract before writing rows") was A/B tested in Codex (2026-08-14,
3 runs per arm, first-draft grading against the FT189 answer key): treatment
medians tied at 3/5, controls did not improve, and no run in either arm
found the intra-package PruneLandedBranches caller or all eight
--git-common-dir sites. Reviewer verdict: dropped. Residual observation for
the drain: the deep consumer surfaces were found only by adversarial
tree-verifying review, not by authoring guidance — the verification loop,
not the skill, is the load-bearing control.
