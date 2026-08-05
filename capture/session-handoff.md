# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` — HEAD `8945287`, 13 dirty paths, 1 unpushed commit
Spec: `specs/exact-prospective-landing/spec.md` (Status: staged), `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/go-build-cache-footprint/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: green at `8945287` — stale only by this capture-and-roadmap batch

## State

**Phase reached: `/bench-what-next` batch drafted for approval.** The approved
two-commit lifecycle repair first landed `authoring-hardening`'s stranded
implemented flip as `8945287`; the sanctioned retirement now supplies this
batch's spec deletions. FT193 leaves the roadmap, FT194 records the reproduced
project-green promotion wedge, and FT195 accounts for the staged
`go-build-cache-footprint` spec.

- Capture drains to zero: no ideas, open learnings, or pending retros remain.
- The four journal entries fold into FT188, FT194, and FT195. The retro's receipt
  generator recommendation folds into FT184; its assignment-path request is
  already served, and its incomplete trailing fragment is dismissed.
- Closed decisions stay closed: exact prospective landing remains the first
  ordered parallel-session scope; FT156 still takes its mechanism ruling at spec
  entry; promotion remains the sole project-green publisher for reviewed builds.
- The pending commit subject must end with `spec-retire: authoring-hardening` so
  `bench spec history` retains the retirement evidence.

## Next command

`$bench-write-spec`

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
