# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `f325816` plus this commit, unpushed (21 commits ahead of origin)
Spec: `specs/exact-prospective-landing/spec.md`, `specs/ft187-communication-surface-cut/spec.md`,
and `specs/pre-push-guard-visibility/spec.md` are staged and untouched. `recovery-discard`
retired in this commit.
Gate: green, run by `bench commit` for this drain.

## State

**Phase reached: roadmap drain complete; the staged spec frontier is untouched.**
All three capture sources are empty — `capture/IDEAS.md` at zero, no open learnings, no
pending retros — so the next `bench status` should recommend build work rather than a drain.

- **This drain verdicted five learnings and one idea, all into existing rows.** The idea
  and the release-preserves-unlanded learning both close against the shipped
  `recovery-discard` route and leave FT98 owning only the residue drain. The lost-review-
  findings learning extends FT107's sixth clause (write `reviews/<spec-slug>.md`, required,
  delegated reviews included); the debug-session worktree learning is FT107's new
  seventeenth clause; the `Assumptions:` field question and the assign-time fence check are
  FT174's; the receipt-generator half re-priced **FT184 from LOW to MEDIUM**, which is the
  one contestable call in the batch.
- **`bench resume` still re-preserves the whole recovery backlog at session start.** The
  route to clear it now exists (`bench worktree recovery <ref> --discard <fingerprint>`,
  one ref per invocation, and `bench spec build reclaim <slug>`), but the drain is a
  reviewer judgment, not automatable — it is FT98's remaining work, not this pass's. Only
  seven refs under `refs/bench/recovery/` still exist; most preserved rows name refs that
  are already gone.
- **The `recovery-discard` build never wrote a retro**, because it landed through
  `bench commit` after its lifecycle run was abandoned rather than through
  `promote` + `/bench-final-check`. Nothing is owed; the evidence that mattered reached
  the roadmap through the two learnings that build produced.

## Next command

`/bench-write-spec` for FT156, taking the anchor-mechanism ruling as the grill at spec
entry — it gates prose batch 1 on both goal tracks and reviewer latency is the binding
constraint. If that ruling is not ready, `/bench-implement-spec` for
`specs/exact-prospective-landing/spec.md` (FT188) instead; it removes the writer lock the
other two staged specs pay. Push is also outstanding: 21 unpushed commits on `main`.

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
