# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `5f9e997`, clean tree, 1 unpushed commit
Spec: none staged.
Gate: green at `a9aeffc` — stale, work tree `448815b`

## State

- **The drain closed and is unpushed.** Commit `5f9e997` reconciled the board,
  emptied `IDEAS.md`, cleared `.bench/learnings.md`, and retired
  `specs/ft148-worktree-orphan-retirement.md`. Both capture sources are at zero
  and `bench roadmap --context` parses clean. The reviewer approved the batch and
  owns the push.
- **FT148 shipped and its residue moved to FT98.** Orphan retirement works, but
  the session-start wall is only partly gone: 21 recovery refs and the 17
  assignment rows holding them survive because compaction correctly declines rows
  that preserve work. `bench worktree recovery <ref>` returns `retain … unlanded`
  on a sampled ref. That is FT98's landed-proof, now recorded in its face one as
  current evidence — not a defect in what FT148 built.
- **Two learnings became FT107 clauses seven and eight.** The review-pickup
  commit is reworded from a buried subordinate clause into its own ordered step,
  and "Route the venue" gets a precedence clause for a harness that *may not*
  spawn a write subagent — `craft-delegate` covers only *cannot*. Both are kit
  prose, built later under `craft-synthesis`. The drained `bench gate pin`
  discoverability idea is FT150 (LOW).
- **`ft91-gate-phase-split` stays unretired on purpose,** so `bench status` will
  keep reporting one spec awaiting retirement until the reviewer rules. Retiring
  it destroys the veto surface on stories 4, 5, and 9: stories 4 and 5 shipped as
  *probed* phases instead of the kit-owned `.bench/phases.json` the spec named,
  and story 9 — that manifest — is unbuilt. The matching roadmap row is the
  `roadmap` signal's "1 row for merged work"; both clear together when the
  reviewer decides.
- **Three FT148 review findings were accepted as residual risk, not fixed,** and
  each is a design change the reviewer owns rather than a regression: `lineSafe`
  admits display-hostile non-control runes, the sweep grades `residualAssignment`
  against a snapshot taken outside the ledger lock, and interrupt-mid-sweep has no
  fault injection.
- **One recovery ref still wants a by-hand look:**
  `refs/bench/recovery/incident-20260712-ambient-probe` matches no assignment.

## Next command

`/bench-debug` — FT91's next arm, the top of the refreshed sequence. Slice C's
premise was falsified: the whole gate is unchanged at ~4m51s and the critical
path is `internal/contract/surface/artifact` (~207 s) and
`internal/contract/surface` (~178 s), untouched by any of the six arms. Ask why
those two cost what they cost before assuming a pipeline shape.

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
