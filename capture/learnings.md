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
Proposed rule change: craft-spec could add one sentence to the coverage-map
section: "before writing rows for a changed contract, enumerate every
production consumer of that contract and give each a row or an exception" —
the exception discipline existed here but the enumeration came from review
rather than the author.
