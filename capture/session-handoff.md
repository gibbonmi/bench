# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `18c7c212`, clean tree, pushed through `f1135fd6`; `036c1761` and later are unpushed
Spec: none staged; `specs/` is empty
Gate: green at `18c7c212` (fresh, via the landing `bench commit` itself)

## State

**Reviewer override standing for this repo:** the 2026-08 capability audit's
own priority order
(`docs/audits/2026-08-bench-capability/results-fable-high/proposed-roadmap.md`)
supersedes `ROADMAP.md`'s `## Recommended sequence` until the audit's active
items (A1–A11; A12 is a decision, not a build) are exhausted. This is recorded
in `ROADMAP.md` itself, not just here.

A1 (live-root conformance in the dev gate) is landed — `a2914fd5` plus the
follow-up gate-defect fix `f1135fd6` (a prospective/composed-tree checkout
never carries the gitignored `dist/` `package.json` `files[]` entry;
`package-core-guard` now exempts it below ship tier). **A2 is next, not yet
started** — a prior attempt in this session edited
`internal/gate/verdict.go` and was explicitly reverted at the reviewer's
request to run in a fresh session instead; nothing of that attempt survives
on disk.

**A2 — one staleness rule for gate verdicts** (P0, no dependencies, audit
estimate ~20 lines + test). Full spec:
`docs/audits/2026-08-bench-capability/results-fable-high/action-items.yaml`,
entry `id: A2`.

- **Problem:** `internal/gate/verdict.go`'s `inspectSubjectAt` checks
  `rec.Status != "green"` (line ~218) and returns *before* checking
  `rec.Tree != s.Tree` (tree drift) or `rec.Oracle != s.Oracle`. So a red
  verdict recorded against a tree the working tree has since moved away from
  (e.g. green@T0 → red@T1 → revert→T0) is reported as a red for the *current*
  tree, when it's really drift — the record doesn't describe T0 at all.
- **`internal/status/status.go`'s `GateVerdict`** compounds this:
  `Stale` is defined only when `in.Status == "green"` (`nonReusableGreen`), so
  a drifted red never gets marked stale and `appendGateInfo` reports
  "red: fix before commit" unconditionally on `gv.Status == "red"`.
- **Proposed fix** (per the audit, worth re-deriving fresh rather than trusting
  verbatim): reorder `inspectSubjectAt` to check tree/oracle drift before the
  status branch; extend `GateVerdict`'s `Stale` derivation to cover a drifted
  record of any status, not just green.
- **Acceptance criteria** (from action-items.yaml, authoritative): the
  green@T0→red@T1→revert→T0 sequence yields "gate green (reused)" or
  "stale — re-run", never "red"; the handoff pin block never writes a red for
  a tree the evidence store holds green; a regression test in
  `internal/gate` or `internal/status` bites when either fix is reverted.
- **Non-goals:** no second verdict store; no change to the evidence key or
  reuse semantics.

## Next command

Fresh session. Recommended: Sonnet orchestrates, delegates the implementation
itself to an Opus subagent (write-delegate, isolated worktree per
`craft-delegate`), Sonnet handles the line declaration, gate, commit, and this
handoff's close. `/bench-debug` is the right phase — A2 is a defect in
`inspectSubjectAt`/`GateVerdict`, not new-feature work, and the audit already
supplies the reproduction shape (green@T0→red@T1→revert→T0); build the Phase 1
repro loop from that before touching either file.

## Shape

Rewritten in full at every phase close, pruned rather than accreted: a fresh
session pays for every line it reads cold, so drop anything it would not act on.
Operational gotchas are placed by lifetime, not copied here: one that recurs across
phases belongs in `projects/benchkit.md`'s cold-session notes, and one scoped to a
build belongs in that spec's coverage rows. This file names at most when you'll hit
one, never the command — a second copy drifts from the source.
Keep the three sections above — **State** (what is true now, including anything
uncommitted), **Next command** (the exact harness-native invocation, not a
description of it), and this one.

The handoff carries no date of its own. `bench status` computes its age from the
commit that last wrote this file and reports a `handoff` row once anything has
landed since. Where this document and the tree disagree, the tree wins.
