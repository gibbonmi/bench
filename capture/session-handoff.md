# Session handoff

Repository: `bench` (origin `https://github.com/gibbonmi/bench.git`)
Path: `~/workspace/bench`
Branch: `main` at pre-drain `77f048d5c44687544c37341fd92be233922a91b4`; 20 commits ahead of upstream before this commit
Spec: `specs/ft187-communication-surface-cut/spec.md` (Status: staged), `specs/pre-push-guard-visibility/spec.md` (Status: staged)
Gate: the roadmap drain commit authors the current verdict

## State

**Phase reached: roadmap capture drained and reconciled after FT195 retirement.**

Both implementation retros, all three parked ideas, and all four open learnings
have explicit roadmap dispositions. FT169 now carries only its residual landing work;
the retired `go-build-cache-footprint` spec and FT195 sequence entry are removed.
FT197 owns Go-managed gate invocation and process lifetime. FT198 owns the
reviewer-requested progressively loaded roadmap design, with `ROADMAP.md` as its
concise index.

Closed decisions stay closed: FT195's Go publication child owns the sealed
transaction; the topology invariant is one all-package entry; artifact promotion
is recorded by staged-entry consumption. The craft-tickets
thinnest-independently-green slicing rule already landed at `1691017` and is not
queued again. Three FT195 review judgments remain available for reviewer veto in
the retained lifecycle record.

Unpushed commits await the reviewer's push decision.

## Next command

`$bench-implement-spec`

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
