# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `2915fcc0`, clean tree, 3 unpushed commits
Spec: none staged; `specs/` is empty
Gate: green at `2915fcc0` (fresh, via the landing `bench commit` itself)

## State

**Reviewer override standing for this repo:** the 2026-08 capability audit's
own priority order
(`docs/audits/2026-08-bench-capability/results-fable-high/proposed-roadmap.md`)
supersedes `ROADMAP.md`'s `## Recommended sequence` until the audit's active
items (A1–A11; A12 is a decision, not a build) are exhausted. This is recorded
in `ROADMAP.md` itself, not just here.

A1 (live-root conformance in the dev gate) is landed — `a2914fd5` plus the
follow-up gate-defect fix `f1135fd6`. **A2 (one staleness rule for gate
verdicts) is landed** — `2915fcc0`. `inspectSubjectAt` now grades tree/oracle
drift via one `driftReason` derivation ahead of the record's own status, and
`Inspection.Drifted` carries that out; `status.GateVerdict` marks a drifted
record stale whatever verdict it carries, and `appendGateInfo` checks `Stale`
before `Status == "red"/"timeout"`. `handoff`'s `gateField` needed no change —
it already deferred staleness entirely to `status.GateVerdict.Stale`.
Verified independently (not just the delegate's claim): full `bench gate
--fresh` green in the delegate's worktree and again after fast-forwarding
main; the green@T0→red@T1→revert→T0 sequence re-run end-to-end against a
freshly built binary in a separate throwaway worktree, confirming `bench
status` now reports `stale (gated tree …, work tree …) → re-run the gate`
and `bench handoff` reports `Gate: red at \`…\` — stale, work tree \`…\`` —
never a bare "red" or a false "— current" for a drifted record.

**A3 is next** — the `/bench` front door (`bench status --route`,
`/bench-drain` rename, staged-spec/un-adopted signals). Both its dependencies
(A1, A2) are now landed. Full spec:
`docs/audits/2026-08-bench-capability/results-fable-high/action-items.yaml`,
entry `id: A3`. Priority P1, medium complexity (one flag, two new status
signals, action normalization, two thin harness adapters, one rename touching
anchors/tests — mitigated by A1's conformance coverage). Read the full entry
before starting; it is not summarized further here to avoid a second
derivation that can drift from the source.

## Next command

Fresh session. `/bench-write-spec` is the right phase — A3 is genuinely
multi-seam (status signals, a new flag, two harness adapters, a rename) and
crosses the light-path threshold in `.bench/BENCH.md`, unlike A1/A2 which were
single-seam defects. Read `action-items.yaml`'s A3 entry and
`proposed-roadmap.md` for the full shape before drafting stories.

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
